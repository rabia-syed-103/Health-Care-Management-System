package controllers

import (
	"context"
	"fmt"
	"hospital-management/db"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type CreateBloodRequestRequest struct {
	DoctorID       int    `json:"doctor_id"       binding:"required"`
	PatientEmail   string `json:"patient_email"   binding:"required"`
	QuantityNeeded int    `json:"quantity_needed" binding:"required,min=1"`
}

type FulfillBloodRequestRequest struct {
	RequestID int `json:"request_id" binding:"required"`
	ManagerID int `json:"manager_id" binding:"required"`
}

func CreateBloodRequest(c *gin.Context) {
	var req CreateBloodRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	ctx := context.Background()

	var doctorName string
	if err := db.Pool.QueryRow(ctx,
		`SELECT name FROM doctor WHERE id = $1`, req.DoctorID,
	).Scan(&doctorName); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Doctor not found"})
		return
	}

	var patientID int
	var patientName, patientBloodGroup string
	if err := db.Pool.QueryRow(ctx,
		`SELECT id, name, TRIM(b_gr) FROM patient WHERE email = $1`, req.PatientEmail,
	).Scan(&patientID, &patientName, &patientBloodGroup); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Patient not found"})
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

	var existingCount int
	if err := tx.QueryRow(ctx,
		`SELECT COUNT(*) FROM blood_request
			WHERE  patient_id = $1 AND status = 'pending'`,
		patientID,
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

	today := time.Now().Format("2006-01-02")
	var requestID int
	if err := tx.QueryRow(ctx,
		`INSERT INTO blood_request (doctor_id, patient_id, quantity_needed, status, request_date)
		VALUES ($1, $2, $3, 'pending', $4)
		RETURNING id`,
		req.DoctorID, patientID, req.QuantityNeeded, today,
	).Scan(&requestID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create blood request — transaction rolled back",
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

func FulfillBloodRequest(c *gin.Context) {
	var req FulfillBloodRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	ctx := context.Background()

	var managerName string
	if err := db.Pool.QueryRow(ctx,
		`SELECT name FROM blood_manager WHERE id = $1`, req.ManagerID,
	).Scan(&managerName); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Blood manager not found"})
		return
	}

	var patientName, patientBloodGroup string
	var quantityNeeded int
	var requestStatus string

	if err := db.Pool.QueryRow(ctx,
		`SELECT p.name, TRIM(p.b_gr), br.quantity_needed, br.status
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

	var lockedStatus string
	if err := tx.QueryRow(ctx,
		`SELECT status FROM blood_request WHERE id = $1 FOR UPDATE`,
		req.RequestID,
	).Scan(&lockedStatus); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to lock request — transaction rolled back"})
		return
	}

	if lockedStatus != "pending" {
		c.JSON(http.StatusConflict, gin.H{
			"error": fmt.Sprintf(
				"Blood request ID %d was just fulfilled by another manager. Status: '%s'.",
				req.RequestID, lockedStatus,
			),
		})
		return
	}

	var selectedBloodID int
	var selectedBloodGroup string
	var selectedUnits int
	var selectedExpiry time.Time

	err = tx.QueryRow(ctx,
		`SELECT b.id, b.b_gr, b.unit, b.expiry_date
		 FROM   blood b
		 WHERE  b.status      IN ('available')
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

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}
	tx = nil

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
