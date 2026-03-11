package controllers

import (
	"context"
	"fmt"
	"hospital-backend/db"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type BloodDonationRequest struct {
	DonorID   int `json:"donor_id"   binding:"required"`
	ManagerID int `json:"manager_id" binding:"required"`
}

// TRANSACTION 3 — BLOOD DONATION

func RecordBloodDonation(c *gin.Context) {
	var req BloodDonationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	ctx := context.Background()
	today := time.Now().Format("2006-01-02")

	var managerName string
	if err := db.Pool.QueryRow(ctx,
		`SELECT name FROM blood_manager WHERE id = $1`, req.ManagerID,
	).Scan(&managerName); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Blood manager not found"})
		return
	}

	var donorName string
	var bloodGroup string
	var lastDonate time.Time

	if err := db.Pool.QueryRow(ctx,
		`SELECT name, b_gr, last_donate FROM donor WHERE id = $1`, req.DonorID,
	).Scan(&donorName, &bloodGroup, &lastDonate); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Donor not found"})
		return
	}

	daysSince := int(time.Since(lastDonate).Hours() / 24)
	if daysSince < 90 {
		eligibleDate := lastDonate.AddDate(0, 0, 90).Format("2006-01-02")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf(
				"Donor '%s' is not eligible yet. Last donated %d days ago on %s. Eligible from: %s (90-day rule).",
				donorName,
				daysSince,
				lastDonate.Format("2006-01-02"),
				eligibleDate,
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

	var lockedDonorID int
	if err := tx.QueryRow(ctx,
		`SELECT id FROM donor WHERE id = $1 FOR UPDATE`, req.DonorID,
	).Scan(&lockedDonorID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to lock donor record — transaction rolled back",
		})
		return
	}

	var freshLastDonate time.Time
	if err := tx.QueryRow(ctx,
		`SELECT last_donate FROM donor WHERE id = $1`, req.DonorID,
	).Scan(&freshLastDonate); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to re-verify donor — transaction rolled back"})
		return
	}

	freshDaysSince := int(time.Since(freshLastDonate).Hours() / 24)
	if freshDaysSince < 90 {
		eligibleDate := freshLastDonate.AddDate(0, 0, 90).Format("2006-01-02")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf(
				"Donor '%s' is not eligible. Another donation was just recorded. Eligible from: %s.",
				donorName, eligibleDate,
			),
		})
		return
	}

	expiryDate := time.Now().AddDate(0, 0, 42).Format("2006-01-02")

	var bloodID int
	if err := tx.QueryRow(ctx,
		`INSERT INTO blood (b_gr, collected_date, status, expiry_date, unit, donor_id)
		 VALUES ($1, $2, 'available', $3, 0, $4)
		 RETURNING id`,
		bloodGroup, today, expiryDate, req.DonorID,
	).Scan(&bloodID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create blood inventory entry — transaction rolled back",
			"details": err.Error(),
		})
		return
	}

	var donationID int
	if err := tx.QueryRow(ctx,
		`INSERT INTO donation (donor_id, manager_id, blood_id, donation_date, status)
		 VALUES ($1, $2, $3, $4, 'completed')
		 RETURNING id`,
		req.DonorID, req.ManagerID, bloodID, today,
	).Scan(&donationID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to record donation — transaction rolled back",
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
		"message":       "Blood donation recorded successfully",
		"donation_id":   donationID,
		"blood_id":      bloodID,
		"donor":         donorName,
		"blood_group":   bloodGroup,
		"manager":       managerName,
		"donation_date": today,
		"expiry_date":   expiryDate,
		"units_added":   1,
	})
}
