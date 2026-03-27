package controllers

import (
	"context"
	"hospital-management/db"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type PrescriptionMedicineItem struct {
	MedicineName string `json:"medicine_name" binding:"required"`
	Quantity     int    `json:"quantity"      binding:"required,min=1"`
}
type PrescribeMedicinesRequest struct {
	DoctorEmail  string                     `json:"doctor_email"  binding:"required"`
	PatientEmail string                     `json:"patient_email" binding:"required"`
	Medicines    []PrescriptionMedicineItem `json:"medicines"     binding:"required,min=1"`
}

type DispenseMedicinesRequest struct {
	PrescriptionID int `json:"prescription_id" binding:"required"`
	PharmacistID   int `json:"pharmacist_id"   binding:"required"`
}

func PrescribeMedicines(c *gin.Context) {
	var req PrescribeMedicinesRequest
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

	var doctorID int
	var doctorName string
	if err := db.Pool.QueryRow(ctx,
		`SELECT id, name FROM doctor WHERE email = $1`, req.DoctorEmail,
	).Scan(&doctorID, &doctorName); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Doctor not found"})
		return
	}

	var appointmentID int
	if err := db.Pool.QueryRow(ctx,
		`SELECT id FROM appointment
		 WHERE  patient_id = $1
		   AND  doctor_id  = $2
		   AND  status     IN ('pending', 'completed')
		 ORDER BY date DESC, time DESC
		 LIMIT 1`,
		patientID, doctorID,
	).Scan(&appointmentID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No valid appointment found for this patient with the specified doctor",
		})
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

	var prescriptionID int
	today := time.Now().Format("2006-01-02")

	if err := tx.QueryRow(ctx,
		`INSERT INTO prescription (doctor_id, patient_id, date)
		 VALUES ($1, $2, $3)
		 RETURNING id`,
		doctorID, patientID, today,
	).Scan(&prescriptionID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create prescription — transaction rolled back",
			"details": err.Error(),
		})
		return
	}

	for _, med := range req.Medicines {
		var medicineID int
		var expiryDate time.Time
		var stock int
		if err := tx.QueryRow(ctx,
			`SELECT id, expiry_date, stock FROM medicine WHERE LOWER(name) = LOWER($1) LIMIT 1`,
			med.MedicineName,
		).Scan(&medicineID, &expiryDate, &stock); err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Medicine '" + med.MedicineName + "' not found — transaction rolled back",
			})
			return
		}

		if expiryDate.Before(time.Now()) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Medicine '" + med.MedicineName + "' is expired — transaction rolled back",
			})
			return
		}

		if _, err := tx.Exec(ctx,
			`INSERT INTO prescription_medicine (prescription_id, medicine_id, quantity)
			 VALUES ($1, $2, $3)`,
			prescriptionID, medicineID, med.Quantity,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to add medicine to prescription — transaction rolled back",
				"details": err.Error(),
			})
			return
		}
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}
	tx = nil

	c.JSON(http.StatusCreated, gin.H{
		"message":         "Prescription created successfully",
		"prescription_id": prescriptionID,
		"patient":         patientName,
		"doctor":          doctorName,
		"appointment_id":  appointmentID,
		"medicines_count": len(req.Medicines),
	})
}
func DispenseMedicines(c *gin.Context) {
	var req DispenseMedicinesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	ctx := context.Background()

	var pharmacistName string
	if err := db.Pool.QueryRow(ctx,
		`SELECT name FROM pharmacist WHERE id = $1`, req.PharmacistID,
	).Scan(&pharmacistName); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pharmacist not found"})
		return
	}

	// Fetch prescribed medicines (name + quantity only) — outside tx is fine, read-only
	type PrescMed struct {
		Name     string
		Quantity int
	}

	rows, err := db.Pool.Query(ctx,
		`SELECT m.name, pm.quantity
		 FROM   prescription_medicine pm
		 JOIN   medicine m ON m.id = pm.medicine_id
		 WHERE  pm.prescription_id = $1`,
		req.PrescriptionID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch prescription medicines"})
		return
	}

	var prescMeds []PrescMed
	for rows.Next() {
		var m PrescMed
		if err := rows.Scan(&m.Name, &m.Quantity); err != nil {
			rows.Close()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read medicine row"})
			return
		}
		prescMeds = append(prescMeds, m)
	}
	rows.Close() // close explicitly before starting the transaction

	if len(prescMeds) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Prescription not found or has no medicines"})
		return
	}

	// Begin transaction — stock check + decrement + dispensing insert all happen inside
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

	var dispensingIDs []int

	for _, pm := range prescMeds {
		// Lock the best non-expired batch with sufficient stock FOR UPDATE
		var medicineID int
		err := tx.QueryRow(ctx,
			`SELECT id FROM medicine
			 WHERE  LOWER(name) = LOWER($1)
			   AND  expiry_date > CURRENT_DATE
			   AND  stock >= $2
			 ORDER BY expiry_date ASC
			 LIMIT 1
			 FOR UPDATE`,
			pm.Name, pm.Quantity,
		).Scan(&medicineID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Medicine '" + pm.Name + "' has no valid non-expired batch with sufficient stock.",
			})
			return
		}

		var dispensingID int
		if err := tx.QueryRow(ctx,
			`INSERT INTO dispensing (medicine_id, pharmacist_id, prescription_id, quantity)
			 VALUES ($1, $2, $3, $4)
			 RETURNING id`,
			medicineID, req.PharmacistID, req.PrescriptionID, pm.Quantity,
		).Scan(&dispensingID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Dispensing failed for '" + pm.Name + "'",
				"details": err.Error(),
			})
			return
		}

		dispensingIDs = append(dispensingIDs, dispensingID)
	}

	if err = tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}
	tx = nil

	c.JSON(http.StatusOK, gin.H{
		"message":         "Medicines dispensed successfully",
		"prescription_id": req.PrescriptionID,
		"pharmacist":      pharmacistName,
		"dispensing_ids":  dispensingIDs,
		"medicines_count": len(dispensingIDs),
	})
}
