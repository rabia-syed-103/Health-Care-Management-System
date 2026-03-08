package controllers

import (
	"context"
	"fmt"
	"hospital-backend/db"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// ─── REQUEST STRUCTS ─────────────────────────────────────────────────────────

type PrescriptionMedicineItem struct {
	MedicineID int `json:"medicine_id" binding:"required"`
	Quantity   int `json:"quantity"    binding:"required,min=1"`
}

type PrescribeMedicinesRequest struct {
	DoctorID  int                        `json:"doctor_id"   binding:"required"`
	PatientID int                        `json:"patient_id"  binding:"required"`
	Medicines []PrescriptionMedicineItem `json:"medicines"   binding:"required,min=1"`
}

type DispenseMedicinesRequest struct {
	PrescriptionID int `json:"prescription_id" binding:"required"`
	PharmacistID   int `json:"pharmacist_id"   binding:"required"`
}

// ─── TRANSACTION 1A — PRESCRIBE MEDICINES ────────────────────────────────────
// Steps:
//   1. Verify patient exists
//   2. Verify doctor exists
//   3. Verify patient has a completed/pending appointment with that doctor
//   4. BEGIN transaction
//   5. Insert into Prescription
//   6. Insert each medicine into Prescription_Medicine
//   7. COMMIT

func PrescribeMedicines(c *gin.Context) {
	var req PrescribeMedicinesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	ctx := context.Background()

	// ── STEP 1: Verify patient exists ────────────────────────────────────────
	var patientName string
	err := db.Pool.QueryRow(ctx,
		`SELECT name FROM patient WHERE id = $1`, req.PatientID,
	).Scan(&patientName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Patient not found"})
		return
	}

	// ── STEP 2: Verify doctor exists ─────────────────────────────────────────
	var doctorName string
	err = db.Pool.QueryRow(ctx,
		`SELECT name FROM doctor WHERE id = $1`, req.DoctorID,
	).Scan(&doctorName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Doctor not found"})
		return
	}

	// ── STEP 3: Verify appointment exists ────────────────────────────────────
	var appointmentID int
	err = db.Pool.QueryRow(ctx,
		`SELECT id FROM appointment
		 WHERE  patient_id = $1
		   AND  doctor_id  = $2
		   AND  status     IN ('pending', 'completed')
		 ORDER BY date DESC, time DESC
		 LIMIT 1`,
		req.PatientID, req.DoctorID,
	).Scan(&appointmentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No valid appointment found for this patient with the specified doctor",
		})
		return
	}

	// ── STEP 4: BEGIN transaction ─────────────────────────────────────────────
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to begin transaction"})
		return
	}
	// ROLLBACK fires automatically if we return before COMMIT
	defer func() {
		if tx != nil {
			tx.Rollback(ctx)
		}
	}()

	// ── STEP 5: Insert into Prescription ─────────────────────────────────────
	var prescriptionID int
	today := time.Now().Format("2006-01-02")

	err = tx.QueryRow(ctx,
		`INSERT INTO prescription (doctor_id, patient_id, date)
		 VALUES ($1, $2, $3)
		 RETURNING id`,
		req.DoctorID, req.PatientID, today,
	).Scan(&prescriptionID)
	if err != nil {
		// ROLLBACK triggered by defer
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create prescription — transaction rolled back",
			"details": err.Error(),
		})
		return
	}

	// ── STEP 6: Insert each medicine into Prescription_Medicine ──────────────
	// The trg_pm_before trigger will fire here and catch duplicates automatically
	for _, med := range req.Medicines {
		// Verify medicine exists and is not expired before inserting
		var medicineName string
		var expiryDate time.Time
		var stock int
		err = tx.QueryRow(ctx,
			`SELECT name, expiry_date, stock FROM medicine WHERE id = $1`,
			med.MedicineID,
		).Scan(&medicineName, &expiryDate, &stock)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Medicine ID " + itoa(med.MedicineID) + " not found — transaction rolled back",
			})
			return
		}

		if expiryDate.Before(time.Now()) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Medicine '" + medicineName + "' (ID: " + itoa(med.MedicineID) + ") is expired — transaction rolled back",
			})
			return
		}

		_, err = tx.Exec(ctx,
			`INSERT INTO prescription_medicine (prescription_id, medicine_id, quantity)
			 VALUES ($1, $2, $3)`,
			prescriptionID, med.MedicineID, med.Quantity,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to add medicine to prescription — transaction rolled back",
				"details": err.Error(),
			})
			return
		}
	}

	// ── STEP 7: COMMIT ────────────────────────────────────────────────────────
	if err = tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}
	tx = nil // prevent defer rollback after successful commit

	c.JSON(http.StatusCreated, gin.H{
		"message":         "Prescription created successfully",
		"prescription_id": prescriptionID,
		"patient":         patientName,
		"doctor":          doctorName,
		"appointment_id":  appointmentID,
		"medicines_count": len(req.Medicines),
	})
}

// ─── TRANSACTION 1B — DISPENSE MEDICINES ─────────────────────────────────────
// Steps:
//   1. Verify pharmacist exists
//   2. Verify prescription exists and retrieve its medicines
//   3. BEGIN transaction
//   4. For each medicine in prescription:
//        a. Check stock availability (server-side, before DB trigger fires)
//        b. Check medicine is not expired
//        c. Insert into Dispensing (trg_dispensing_before fires → validates)
//        d. Stock deduction handled by trg_dispensing_after_insert trigger
//   5. COMMIT

func DispenseMedicines(c *gin.Context) {
	var req DispenseMedicinesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	ctx := context.Background()

	// ── STEP 1: Verify pharmacist exists ─────────────────────────────────────
	var pharmacistName string
	err := db.Pool.QueryRow(ctx,
		`SELECT name FROM pharmacist WHERE id = $1`, req.PharmacistID,
	).Scan(&pharmacistName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pharmacist not found"})
		return
	}

	// ── STEP 2: Fetch prescription medicines ─────────────────────────────────
	type MedRow struct {
		PrescMedID int
		MedicineID int
		Name       string
		Quantity   int
		Stock      int
		Expiry     time.Time
	}

	rows, err := db.Pool.Query(ctx,
		`SELECT pm.id, pm.medicine_id, m.name, pm.quantity, m.stock, m.expiry_date
		 FROM   prescription_medicine pm
		 JOIN   medicine m ON m.id = pm.medicine_id
		 WHERE  pm.prescription_id = $1`,
		req.PrescriptionID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch prescription medicines"})
		return
	}
	defer rows.Close()

	var medicines []MedRow
	for rows.Next() {
		var m MedRow
		if err := rows.Scan(&m.PrescMedID, &m.MedicineID, &m.Name, &m.Quantity, &m.Stock, &m.Expiry); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read medicine row"})
			return
		}
		medicines = append(medicines, m)
	}

	if len(medicines) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Prescription not found or has no medicines"})
		return
	}

	// ── Server-side pre-validation before opening transaction ─────────────────
	for _, m := range medicines {
		if m.Expiry.Before(time.Now()) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Medicine '" + m.Name + "' is expired. Cannot dispense.",
			})
			return
		}
		if m.Stock < m.Quantity {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Insufficient stock for '" + m.Name + "'. Available: " +
					itoa(m.Stock) + ", Needed: " + itoa(m.Quantity),
			})
			return
		}
	}

	// ── STEP 3: BEGIN transaction ─────────────────────────────────────────────
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

	// ── STEP 4: Insert into Dispensing for each medicine ─────────────────────
	// trg_dispensing_before  → validates expiry + stock
	// trg_dispensing_after_insert → deducts stock automatically
	var dispensingIDs []int
	for _, m := range medicines {
		var dispensingID int
		err = tx.QueryRow(ctx,
			`INSERT INTO dispensing (medicine_id, pharmacist_id, prescription_id, quantity)
			 VALUES ($1, $2, $3, $4)
			 RETURNING id`,
			m.MedicineID, req.PharmacistID, req.PrescriptionID, m.Quantity,
		).Scan(&dispensingID)
		if err != nil {
			// DB trigger may have fired with a detailed error message
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Dispensing failed for '" + m.Name + "' — transaction rolled back",
				"details": err.Error(),
			})
			return
		}
		dispensingIDs = append(dispensingIDs, dispensingID)
	}

	// ── STEP 5: COMMIT ────────────────────────────────────────────────────────
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
		"medicines_count": len(medicines),
	})
}

// ─── HELPER ──────────────────────────────────────────────────────────────────

func itoa(n int) string {
	return fmt.Sprintf("%d", n)
}
