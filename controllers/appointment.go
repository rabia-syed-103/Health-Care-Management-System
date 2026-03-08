package controllers

import (
	"context"
	"fmt"
	"hospital-backend/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ─── REQUEST STRUCTS ──────────────────────────────────────────────────────────

type BookAppointmentRequest struct {
	PatientID      int    `json:"patient_id"      binding:"required"`
	DoctorID       int    `json:"doctor_id"       binding:"required"`
	ReceptionistID int    `json:"receptionist_id" binding:"required"`
	Date           string `json:"date"            binding:"required"` // "YYYY-MM-DD"
	Time           string `json:"time"            binding:"required"` // "HH:MM"
}

type BookOTAppointmentRequest struct {
	PatientID      int    `json:"patient_id"      binding:"required"`
	DoctorID       int    `json:"doctor_id"       binding:"required"`
	ReceptionistID int    `json:"receptionist_id" binding:"required"`
	Date           string `json:"date"            binding:"required"`
	Time           string `json:"time"            binding:"required"`
}

// ─── TRANSACTION 2A — PATIENT BOOKS REGULAR APPOINTMENT ─────────────────────
// Patient comes in, receptionist schedules with a doctor. No OT needed.
//
// Flow:
//  1. Verify patient exists
//  2. Verify doctor exists
//  3. Verify receptionist exists
//  4. BEGIN transaction
//  5. Check doctor is free at requested date + time
//  6. INSERT appointment (ot_id = NULL, status = 'pending')
//     trg_appointment_management fires as 2nd guard against double booking
//  7. COMMIT

func BookAppointment(c *gin.Context) {
	var req BookAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	ctx := context.Background()

	// ── STEP 1: Verify patient exists ────────────────────────────────────────
	var patientName string
	if err := db.Pool.QueryRow(ctx,
		`SELECT name FROM patient WHERE id = $1`, req.PatientID,
	).Scan(&patientName); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Patient not found"})
		return
	}

	// ── STEP 2: Verify doctor exists ─────────────────────────────────────────
	var doctorName string
	if err := db.Pool.QueryRow(ctx,
		`SELECT name FROM doctor WHERE id = $1`, req.DoctorID,
	).Scan(&doctorName); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Doctor not found"})
		return
	}

	// ── STEP 3: Verify receptionist exists ───────────────────────────────────
	var receptionistName string
	if err := db.Pool.QueryRow(ctx,
		`SELECT name FROM receptionist WHERE id = $1`, req.ReceptionistID,
	).Scan(&receptionistName); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Receptionist not found"})
		return
	}

	// ── STEP 4: BEGIN transaction ─────────────────────────────────────────────
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to begin transaction"})
		return
	}
	defer func() {
		if tx != nil {
			tx.Rollback(ctx)
		}
	}()

	// ── STEP 5: Check doctor is free at requested date + time ─────────────────
	var conflictCount int
	if err := tx.QueryRow(ctx,
		`SELECT COUNT(*) FROM appointment
		 WHERE  doctor_id = $1
		   AND  date      = $2
		   AND  time      = $3
		   AND  status    = 'pending'`,
		req.DoctorID, req.Date, req.Time,
	).Scan(&conflictCount); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to check doctor availability — transaction rolled back",
		})
		return
	}
	if conflictCount > 0 {
		c.JSON(http.StatusConflict, gin.H{
			"error": fmt.Sprintf(
				"Dr. %s already has a pending appointment on %s at %s. Please pick another slot.",
				doctorName, req.Date, req.Time,
			),
		})
		return
	}

	// ── STEP 6: INSERT appointment (ot_id = NULL, status = 'pending') ─────────
	// trg_appointment_management fires here as the second guard
	var appointmentID int
	if err := tx.QueryRow(ctx,
		`INSERT INTO appointment (receptionist_id, patient_id, doctor_id, date, time, status, ot_id)
		 VALUES ($1, $2, $3, $4, $5, 'pending', NULL)
		 RETURNING id`,
		req.ReceptionistID, req.PatientID, req.DoctorID, req.Date, req.Time,
	).Scan(&appointmentID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create appointment — transaction rolled back",
			"details": err.Error(),
		})
		return
	}

	// ── STEP 7: COMMIT ────────────────────────────────────────────────────────
	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}
	tx = nil

	c.JSON(http.StatusCreated, gin.H{
		"message":        "Appointment booked successfully",
		"appointment_id": appointmentID,
		"status":         "pending",
		"ot_id":          nil,
		"patient":        patientName,
		"doctor":         doctorName,
		"receptionist":   receptionistName,
		"date":           req.Date,
		"time":           req.Time,
	})
}

// ─── TRANSACTION 2B — DOCTOR REQUESTS OT APPOINTMENT FOR A PATIENT ───────────
// Doctor needs an operating theatre for a patient.
// Receptionist finds a free OT at the requested time and books it.
//
// Flow:
//  1. Verify patient exists
//  2. Verify doctor exists
//  3. Verify receptionist exists
//  4. BEGIN transaction
//  5. Check doctor is free at requested date + time
//  6. Find a free OT: is_available = TRUE AND not already booked at same slot
//  7. Lock that OT row (SELECT FOR UPDATE) — prevents two receptionists
//     grabbing the same OT simultaneously
//  8. INSERT appointment with found ot_id, status = 'pending'
//  9. UPDATE ot SET is_available = FALSE
// 10. COMMIT

func BookOTAppointment(c *gin.Context) {
	var req BookOTAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	ctx := context.Background()

	// ── STEP 1: Verify patient exists ────────────────────────────────────────
	var patientName string
	if err := db.Pool.QueryRow(ctx,
		`SELECT name FROM patient WHERE id = $1`, req.PatientID,
	).Scan(&patientName); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Patient not found"})
		return
	}

	// ── STEP 2: Verify doctor exists ─────────────────────────────────────────
	var doctorName string
	if err := db.Pool.QueryRow(ctx,
		`SELECT name FROM doctor WHERE id = $1`, req.DoctorID,
	).Scan(&doctorName); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Doctor not found"})
		return
	}

	// ── STEP 3: Verify receptionist exists ───────────────────────────────────
	var receptionistName string
	if err := db.Pool.QueryRow(ctx,
		`SELECT name FROM receptionist WHERE id = $1`, req.ReceptionistID,
	).Scan(&receptionistName); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Receptionist not found"})
		return
	}

	// ── STEP 4: BEGIN transaction ─────────────────────────────────────────────
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to begin transaction"})
		return
	}
	defer func() {
		if tx != nil {
			tx.Rollback(ctx)
		}
	}()

	// ── STEP 5: Check doctor is free at requested date + time ─────────────────
	var conflictCount int
	if err := tx.QueryRow(ctx,
		`SELECT COUNT(*) FROM appointment
		 WHERE  doctor_id = $1
		   AND  date      = $2
		   AND  time      = $3
		   AND  status    = 'pending'`,
		req.DoctorID, req.Date, req.Time,
	).Scan(&conflictCount); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to check doctor availability — transaction rolled back",
		})
		return
	}
	if conflictCount > 0 {
		c.JSON(http.StatusConflict, gin.H{
			"error": fmt.Sprintf(
				"Dr. %s already has a pending appointment on %s at %s.",
				doctorName, req.Date, req.Time,
			),
		})
		return
	}

	// ── STEPS 6 + 7: Find a free OT and lock it ──────────────────────────────
	// Free OT = is_available TRUE AND not already booked at this date+time
	// FOR UPDATE locks the row so two concurrent receptionists can't
	// grab the same OT at the same time
	var otID int
	if err := tx.QueryRow(ctx,
		`SELECT ot.id
		 FROM   ot
		 WHERE  ot.is_available = TRUE
		   AND  ot.id NOT IN (
		            SELECT a.ot_id
		            FROM   appointment a
		            WHERE  a.ot_id  IS NOT NULL
		              AND  a.date   = $1
		              AND  a.time   = $2
		              AND  a.status = 'pending'
		        )
		 ORDER BY ot.id
		 LIMIT  1
		 FOR UPDATE`,
		req.Date, req.Time,
	).Scan(&otID); err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": fmt.Sprintf(
				"No operating theatre available on %s at %s. All OTs are occupied.",
				req.Date, req.Time,
			),
		})
		return
	}

	// ── STEP 8: INSERT appointment with the found OT ──────────────────────────
	var appointmentID int
	if err := tx.QueryRow(ctx,
		`INSERT INTO appointment (receptionist_id, patient_id, doctor_id, date, time, status, ot_id)
		 VALUES ($1, $2, $3, $4, $5, 'pending', $6)
		 RETURNING id`,
		req.ReceptionistID, req.PatientID, req.DoctorID, req.Date, req.Time, otID,
	).Scan(&appointmentID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create OT appointment — transaction rolled back",
			"details": err.Error(),
		})
		return
	}

	// ── STEP 9: Mark OT as unavailable ───────────────────────────────────────
	if _, err := tx.Exec(ctx,
		`UPDATE ot SET is_available = FALSE WHERE id = $1`, otID,
	); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to reserve OT — transaction rolled back",
			"details": err.Error(),
		})
		return
	}

	// ── STEP 10: COMMIT ───────────────────────────────────────────────────────
	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}
	tx = nil

	c.JSON(http.StatusCreated, gin.H{
		"message":        "OT appointment booked successfully",
		"appointment_id": appointmentID,
		"status":         "pending",
		"ot_id":          otID,
		"patient":        patientName,
		"doctor":         doctorName,
		"receptionist":   receptionistName,
		"date":           req.Date,
		"time":           req.Time,
	})
}
