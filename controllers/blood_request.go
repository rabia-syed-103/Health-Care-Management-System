package controllers

import (
	"context"
	"fmt"
	"hospital-backend/db"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// ─── REQUEST STRUCTS ──────────────────────────────────────────────────────────

type CreateBloodRequestRequest struct {
	DoctorID       int `json:"doctor_id"       binding:"required"`
	PatientID      int `json:"patient_id"      binding:"required"`
	QuantityNeeded int `json:"quantity_needed" binding:"required,min=1"`
}

type FulfillBloodRequestRequest struct {
	RequestID int `json:"request_id" binding:"required"`
	ManagerID int `json:"manager_id" binding:"required"`
	// No blood_id, no quantity — system reads everything from the request
	// and auto-selects matching compatible blood inventory
}

// ─── TRANSACTION 4A — DOCTOR CREATES BLOOD REQUEST ───────────────────────────
// Doctor requests blood for a patient. Records it in blood_request table.
//
// Flow:
//  1. Verify doctor exists
//  2. Verify patient exists + fetch blood group
//  3. BEGIN transaction
//  4. Check patient has no existing pending blood request
//  5. INSERT into blood_request (status = 'pending')
//  6. COMMIT

func CreateBloodRequest(c *gin.Context) {
	var req CreateBloodRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	ctx := context.Background()

	// ── STEP 1: Verify doctor exists ─────────────────────────────────────────
	var doctorName string
	if err := db.Pool.QueryRow(ctx,
		`SELECT name FROM doctor WHERE id = $1`, req.DoctorID,
	).Scan(&doctorName); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Doctor not found"})
		return
	}

	// ── STEP 2: Verify patient exists + fetch blood group ────────────────────
	var patientName, patientBloodGroup string
	if err := db.Pool.QueryRow(ctx,
		`SELECT name, b_gr FROM patient WHERE id = $1`, req.PatientID,
	).Scan(&patientName, &patientBloodGroup); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Patient not found"})
		return
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

	// ── STEP 4: Check patient has no existing pending blood request ───────────
	var existingCount int
	if err := tx.QueryRow(ctx,
		`SELECT COUNT(*) FROM blood_request
		 WHERE  patient_id = $1 AND status = 'pending'`,
		req.PatientID,
	).Scan(&existingCount); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check existing requests — transaction rolled back"})
		return
	}
	if existingCount > 0 {
		c.JSON(http.StatusConflict, gin.H{
			"error": fmt.Sprintf(
				"Patient '%s' already has a pending blood request. Fulfill or cancel it first.",
				patientName,
			),
		})
		return
	}

	// ── STEP 5: INSERT into blood_request ─────────────────────────────────────
	today := time.Now().Format("2006-01-02")
	var requestID int
	if err := tx.QueryRow(ctx,
		`INSERT INTO blood_request (doctor_id, patient_id, quantity_needed, status, request_date)
		 VALUES ($1, $2, $3, 'pending', $4)
		 RETURNING id`,
		req.DoctorID, req.PatientID, req.QuantityNeeded, today,
	).Scan(&requestID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create blood request — transaction rolled back",
			"details": err.Error(),
		})
		return
	}

	// ── STEP 6: COMMIT ────────────────────────────────────────────────────────
	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}
	tx = nil

	c.JSON(http.StatusCreated, gin.H{
		"message":          "Blood request created successfully",
		"request_id":       requestID,
		"status":           "pending",
		"doctor":           doctorName,
		"patient":          patientName,
		"patient_blood_gr": patientBloodGroup,
		"quantity_needed":  req.QuantityNeeded,
		"request_date":     today,
	})
}

// ─── TRANSACTION 4B — BLOOD MANAGER FULFILLS A BLOOD REQUEST ─────────────────
// Blood manager picks a pending request by request_id.
// System reads the patient's blood group and quantity needed FROM the request.
// Then auto-selects the best matching compatible blood from inventory.
// Records fulfillment and updates blood stock + request status.
//
// Flow:
//  1. Verify blood manager exists
//  2. Fetch blood request — must be 'pending'
//     Read: patient blood group + quantity_needed directly from request
//  3. BEGIN transaction
//  4. Lock blood request row (FOR UPDATE) — prevents two managers
//     fulfilling the same request simultaneously
//  5. Re-verify request is still 'pending' after lock
//  6. AUTO-SELECT + LOCK best compatible blood entry:
//     - Matches patient blood group compatibility rules
//     - Status = 'available' or 'reserved'
//     - Not expired
//     - Has enough units >= quantity_needed
//     - Ordered by expiry_date ASC (use oldest blood first)
//  7. INSERT into blood_request_fulfillment
//     trg_blood_fulfillment_before fires → DB-level validation
//     trg_blood_fulfillment_after fires  → deducts Blood.unit,
//                                          updates Blood.status,
//                                          sets Blood_Request.status = 'complete'
//  8. COMMIT

func FulfillBloodRequest(c *gin.Context) {
	var req FulfillBloodRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	ctx := context.Background()

	// ── STEP 1: Verify blood manager exists ───────────────────────────────────
	var managerName string
	if err := db.Pool.QueryRow(ctx,
		`SELECT name FROM blood_manager WHERE id = $1`, req.ManagerID,
	).Scan(&managerName); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Blood manager not found"})
		return
	}

	// ── STEP 2: Fetch blood request details ───────────────────────────────────
	// Read patient blood group + quantity needed straight from the request
	var patientName, patientBloodGroup string
	var quantityNeeded int
	var requestStatus string

	if err := db.Pool.QueryRow(ctx,
		`SELECT p.name, p.b_gr, br.quantity_needed, br.status
		 FROM   blood_request br
		 JOIN   patient p ON p.id = br.patient_id
		 WHERE  br.id = $1`,
		req.RequestID,
	).Scan(&patientName, &patientBloodGroup, &quantityNeeded, &requestStatus); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Blood request not found"})
		return
	}

	if requestStatus != "pending" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf(
				"Blood request ID %d is already '%s'. Only pending requests can be fulfilled.",
				req.RequestID, requestStatus,
			),
		})
		return
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

	// ── STEP 4: Lock the blood request row ────────────────────────────────────
	// Prevents two managers fulfilling the same request at the same time
	var lockedStatus string
	if err := tx.QueryRow(ctx,
		`SELECT status FROM blood_request WHERE id = $1 FOR UPDATE`,
		req.RequestID,
	).Scan(&lockedStatus); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to lock request — transaction rolled back"})
		return
	}

	// ── STEP 5: Re-verify still pending after lock ────────────────────────────
	if lockedStatus != "pending" {
		c.JSON(http.StatusConflict, gin.H{
			"error": fmt.Sprintf(
				"Blood request ID %d was just fulfilled by another manager. Status: '%s'.",
				req.RequestID, lockedStatus,
			),
		})
		return
	}

	// ── STEP 6: AUTO-SELECT + LOCK best compatible blood entry ────────────────
	// Compatibility rules match the DB trigger exactly.
	// ORDER BY expiry_date ASC = use oldest stock first (reduces waste).
	// FOR UPDATE = locks the selected blood row immediately.
	var selectedBloodID int
	var selectedBloodGroup string
	var selectedUnits int
	var selectedExpiry time.Time

	err = tx.QueryRow(ctx,
		`SELECT b.id, b.b_gr, b.unit, b.expiry_date
		 FROM   blood b
		 WHERE  b.status      IN ('available', 'reserved')
		   AND  b.expiry_date  > CURRENT_DATE
		   AND  b.unit        >= $1
		   AND  b.b_gr = ANY(
		            CASE $2::varchar
		                WHEN 'O-'  THEN ARRAY['O-']
		                WHEN 'O+'  THEN ARRAY['O-','O+']
		                WHEN 'A-'  THEN ARRAY['O-','A-']
		                WHEN 'A+'  THEN ARRAY['O-','O+','A-','A+']
		                WHEN 'B-'  THEN ARRAY['O-','B-']
		                WHEN 'B+'  THEN ARRAY['O-','O+','B-','B+']
		                WHEN 'AB-' THEN ARRAY['O-','A-','B-','AB-']
		                WHEN 'AB+' THEN ARRAY['O-','O+','A-','A+','B-','B+','AB-','AB+']
		                ELSE ARRAY[]::varchar[]
		            END
		        )
		 ORDER BY b.expiry_date ASC
		 LIMIT  1
		 FOR UPDATE`,
		quantityNeeded, patientBloodGroup,
	).Scan(&selectedBloodID, &selectedBloodGroup, &selectedUnits, &selectedExpiry)

	if err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": fmt.Sprintf(
				"No compatible blood available for patient '%s' (blood group: %s). "+
					"Needed: %d units. No matching non-expired inventory found.",
				patientName, patientBloodGroup, quantityNeeded,
			),
		})
		return
	}

	// ── STEP 7: INSERT into blood_request_fulfillment ─────────────────────────
	// quantity_provided = quantity_needed (full fulfillment)
	// trg_blood_fulfillment_before → final DB-level validation
	// trg_blood_fulfillment_after  → Blood.unit -= qty, Blood.status updated,
	//                                Blood_Request.status → 'complete'
	today := time.Now().Format("2006-01-02")
	var fulfillmentID int
	if err := tx.QueryRow(ctx,
		`INSERT INTO blood_request_fulfillment
		 (request_id, blood_id, manager_id, quantity_provided, fulfillment_date)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id`,
		req.RequestID, selectedBloodID, req.ManagerID, quantityNeeded, today,
	).Scan(&fulfillmentID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fulfill blood request — transaction rolled back",
			"details": err.Error(),
		})
		return
	}

	// ── STEP 8: COMMIT ────────────────────────────────────────────────────────
	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}
	tx = nil

	// Fetch final request status (trigger sets it to 'complete')
	var finalStatus string
	db.Pool.QueryRow(ctx,
		`SELECT status FROM blood_request WHERE id = $1`, req.RequestID,
	).Scan(&finalStatus)

	c.JSON(http.StatusCreated, gin.H{
		"message":          "Blood request fulfilled successfully",
		"fulfillment_id":   fulfillmentID,
		"request_id":       req.RequestID,
		"request_status":   finalStatus,
		"patient":          patientName,
		"patient_blood_gr": patientBloodGroup,
		"blood_id_used":    selectedBloodID,
		"blood_group_used": selectedBloodGroup,
		"units_provided":   quantityNeeded,
		"blood_expiry":     selectedExpiry.Format("2006-01-02"),
		"manager":          managerName,
		"fulfillment_date": today,
	})
}
