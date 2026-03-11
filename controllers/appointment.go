package controllers

import (
	"context"
	"fmt"
	"hospital-backend/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

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

func BookAppointment(c *gin.Context) {
	var req BookAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	ctx := context.Background()

	var patientName string
	if err := db.Pool.QueryRow(ctx,
		`SELECT name FROM patient WHERE id = $1`, req.PatientID,
	).Scan(&patientName); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Patient not found"})
		return
	}

	var doctorName string
	if err := db.Pool.QueryRow(ctx,
		`SELECT name FROM doctor WHERE id = $1`, req.DoctorID,
	).Scan(&doctorName); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Doctor not found"})
		return
	}

	var receptionistName string
	if err := db.Pool.QueryRow(ctx,
		`SELECT name FROM receptionist WHERE id = $1`, req.ReceptionistID,
	).Scan(&receptionistName); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Receptionist not found"})
		return
	}

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

func BookOTAppointment(c *gin.Context) {
	var req BookOTAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	ctx := context.Background()

	var patientName string
	if err := db.Pool.QueryRow(ctx,
		`SELECT name FROM patient WHERE id = $1`, req.PatientID,
	).Scan(&patientName); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Patient not found"})
		return
	}

	var doctorName string
	if err := db.Pool.QueryRow(ctx,
		`SELECT name FROM doctor WHERE id = $1`, req.DoctorID,
	).Scan(&doctorName); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Doctor not found"})
		return
	}

	var receptionistName string
	if err := db.Pool.QueryRow(ctx,
		`SELECT name FROM receptionist WHERE id = $1`, req.ReceptionistID,
	).Scan(&receptionistName); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Receptionist not found"})
		return
	}

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

	if _, err := tx.Exec(ctx,
		`UPDATE ot SET is_available = FALSE WHERE id = $1`, otID,
	); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to reserve OT — transaction rolled back",
			"details": err.Error(),
		})
		return
	}

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
