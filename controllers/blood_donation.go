package controllers

import (
	"context"
	"fmt"
	"hospital-backend/db"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// ─── REQUEST STRUCT ───────────────────────────────────────────────────────────

type BloodDonationRequest struct {
	DonorID   int `json:"donor_id"   binding:"required"`
	ManagerID int `json:"manager_id" binding:"required"`
}

// ─── TRANSACTION 3 — BLOOD DONATION ──────────────────────────────────────────
// A donor arrives to donate blood. Blood manager verifies eligibility,
// creates a blood inventory entry, records the donation, and updates
// the donor's last donation date.
//
// Flow:
//  1. Verify blood manager exists
//  2. Verify donor exists and fetch their blood group + last donation date
//  3. Server-side 90-day eligibility check (before opening transaction)
//  4. BEGIN transaction
//  5. Lock donor row (SELECT FOR UPDATE) — prevents duplicate donations
//     if donor walks into two locations simultaneously
//  6. Re-check 90-day rule inside transaction (after lock) — ironclad guard
//  7. INSERT into Blood — create a new blood inventory entry for this donation
//     trg_blood_before fires → validates no negative units, auto-expires if needed
//  8. INSERT into Donation — record the donation
//     trg_donor_management_before fires → validates 90-day rule at DB level
//     trg_donor_management_after fires → updates donor's Last_Donate + increments Blood.Unit
//  9. COMMIT

func RecordBloodDonation(c *gin.Context) {
	var req BloodDonationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	ctx := context.Background()
	today := time.Now().Format("2006-01-02")

	// ── STEP 1: Verify blood manager exists ───────────────────────────────────
	var managerName string
	if err := db.Pool.QueryRow(ctx,
		`SELECT name FROM blood_manager WHERE id = $1`, req.ManagerID,
	).Scan(&managerName); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Blood manager not found"})
		return
	}

	// ── STEP 2: Verify donor exists, fetch blood group + last donation date ───
	var donorName string
	var bloodGroup string
	var lastDonate time.Time

	if err := db.Pool.QueryRow(ctx,
		`SELECT name, b_gr, last_donate FROM donor WHERE id = $1`, req.DonorID,
	).Scan(&donorName, &bloodGroup, &lastDonate); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Donor not found"})
		return
	}

	// ── STEP 3: Server-side 90-day eligibility check ─────────────────────────
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

	// ── STEP 5: Lock donor row (SELECT FOR UPDATE) ────────────────────────────
	// Prevents the same donor from donating at two locations at the same time.
	// If another transaction already locked this donor, this will wait/fail.
	var lockedDonorID int
	if err := tx.QueryRow(ctx,
		`SELECT id FROM donor WHERE id = $1 FOR UPDATE`, req.DonorID,
	).Scan(&lockedDonorID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to lock donor record — transaction rolled back",
		})
		return
	}

	// ── STEP 6: Re-check 90-day rule inside transaction (after lock) ──────────
	// The lock guarantees no other transaction updated last_donate between
	// our pre-check (Step 3) and now.
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

	// ── STEP 7: INSERT into Blood ─────────────────────────────────────────────
	// Create a new blood inventory entry for this donation.
	// - unit starts at 0 (trg_donor_management_after will increment it to 1)
	// - collected_date = today
	// - expiry_date = today + 42 days (standard blood shelf life)
	// - status = 'available'
	// trg_blood_before fires here → validates unit >= 0, auto-expires if needed
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

	// ── STEP 8: INSERT into Donation ─────────────────────────────────────────
	// trg_donor_management_before fires → validates 90-day rule at DB level (final guard)
	// trg_donor_management_after fires → increments Blood.unit by 1
	//                                  → updates Donor.last_donate to today
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

	// ── STEP 9: COMMIT ────────────────────────────────────────────────────────
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
