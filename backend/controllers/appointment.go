package controllers

import (
	"context"
	"fmt"
	"hospital-management/db"
	"hospital-management/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BookAppointmentRequest struct {
	PatientEmail   string `json:"patient_email"   binding:"required"`
	DoctorID       int    `json:"doctor_id"       binding:"required"`
	ReceptionistID int    `json:"receptionist_id" binding:"required"`
	Date           string `json:"date"            binding:"required"`
	Time           string `json:"time"            binding:"required"`
}

type BookOTAppointmentRequest struct {
	PatientEmail   string `json:"patient_email"   binding:"required"`
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

	var patientID int
	var patientName string
	if err := db.Pool.QueryRow(ctx,
		`SELECT id, name FROM patient WHERE email = $1`, req.PatientEmail,
	).Scan(&patientID, &patientName); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Patient not found"})
		return
	}

	var doctorName, doctorEmail string
	if err := db.Pool.QueryRow(ctx,
		`SELECT name, email FROM doctor WHERE id = $1`, req.DoctorID,
	).Scan(&doctorName, &doctorEmail); err != nil {
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
		req.ReceptionistID, patientID, req.DoctorID, req.Date, req.Time,
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
	// Send email in background
	go func() {
		err := utils.SendEmail(
			req.PatientEmail,
			"Appointment Confirmed",
			fmt.Sprintf(
				"Hello %s,\n\nYour appointment has been booked successfully.\n\nDoctor: %s\nDate: %s\nTime: %s\n\nThank you!",
				patientName, doctorName, req.Date, req.Time,
			),
		)
		if err != nil {
			fmt.Println("Patient email error:", err)
		}

		err = utils.SendEmail(
			doctorEmail,
			"New Appointment Scheduled",
			fmt.Sprintf(
				"Hello Dr. %s,\n\nYou have a new appointment scheduled with patient %s.\n\nDate: %s\nTime: %s\n\nThank you!",
				doctorName, patientName, req.Date, req.Time,
			),
		)
		if err != nil {
			fmt.Println("Doctor email error:", err)
		}
	}()

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

	var patientID int
	var patientName string
	if err := db.Pool.QueryRow(ctx,
		`SELECT id, name FROM patient WHERE email = $1`, req.PatientEmail,
	).Scan(&patientID, &patientName); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Patient not found"})
		return
	}

	var doctorName, doctorEmail string
	if err := db.Pool.QueryRow(ctx,
		`SELECT name, email FROM doctor WHERE id = $1`, req.DoctorID,
	).Scan(&doctorName, &doctorEmail); err != nil {
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
		req.ReceptionistID, patientID, req.DoctorID, req.Date, req.Time, otID,
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
	// 📧 Send email in background
	go func() {
		err := utils.SendEmail(
			req.PatientEmail,
			"OT Appointment Confirmed",
			fmt.Sprintf(
				"Hello %s,\n\nYour OT appointment has been booked.\n\nDoctor: %s\nDate: %s\nTime: %s\nOT ID: %d\n\nThank you!",
				patientName, doctorName, req.Date, req.Time, otID,
			),
		)
		if err != nil {
			fmt.Println("Email error:", err)
		}
		err = utils.SendEmail(
			doctorEmail,
			"New Appointment Scheduled",
			fmt.Sprintf(
				"Hello Dr. %s,\n\nYou have a new appointment scheduled with patient %s.\n\nDate: %s\nTime: %s\nOT ID: %d\n\nThank you!",
				doctorName, patientName, req.Date, req.Time, otID,
			),
		)
		if err != nil {
			fmt.Println("Doctor email error:", err)
		}
	}()

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

type GetAvailableDoctorsRequest struct {
	Date string `json:"date" binding:"required"`
	Time string `json:"time" binding:"required"`
}

func GetAvailableDoctors(c *gin.Context) {
	var req GetAvailableDoctorsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	ctx := context.Background()

	rows, err := db.Pool.Query(ctx,
		`SELECT id, name, specialization, p_no, email
		 FROM   doctor
		 WHERE  id NOT IN (
		            SELECT doctor_id FROM appointment
		            WHERE  date   = $1
		              AND  time   = $2
		              AND  status = 'pending'
		        )
		 ORDER BY name ASC`,
		req.Date, req.Time,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch available doctors: " + err.Error()})
		return
	}
	defer rows.Close()

	type Doctor struct {
		ID             int    `json:"id"`
		Name           string `json:"name"`
		Specialization string `json:"specialization"`
		PNo            string `json:"p_no"`
		Email          string `json:"email"`
	}

	var doctors []Doctor
	for rows.Next() {
		var d Doctor
		if err := rows.Scan(&d.ID, &d.Name, &d.Specialization, &d.PNo, &d.Email); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read doctor data: " + err.Error()})
			return
		}
		doctors = append(doctors, d)
	}

	if doctors == nil {
		doctors = []Doctor{}
	}

	c.JSON(http.StatusOK, gin.H{
		"available_doctors": doctors,
		"count":             len(doctors),
		"date":              req.Date,
		"time":              req.Time,
	})
}

func CancelAppointment(c *gin.Context) {
	id := c.Param("id")
	ctx := context.Background()

	var currentStatus string
	if err := db.Pool.QueryRow(ctx,
		`SELECT status FROM appointment WHERE id = $1`, id,
	).Scan(&currentStatus); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Appointment not found"})
		return
	}

	if currentStatus != "pending" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Cannot cancel appointment with status '%s'. Only pending appointments can be cancelled.", currentStatus),
		})
		return
	}

	_, err := db.Pool.Exec(ctx,
		`UPDATE appointment SET status = 'cancelled' WHERE id = $1`, id,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel appointment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "Appointment cancelled successfully",
		"appointment_id": id,
		"status":         "cancelled",
	})
}
