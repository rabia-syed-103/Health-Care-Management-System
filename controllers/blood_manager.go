package controllers

import (
	"context"
	"hospital-backend/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetDonationHistory(c *gin.Context) {
	ctx := context.Background()

	rows, err := db.Pool.Query(ctx,
		`SELECT
			don.id,
			don.donation_date::text,
			don.status,
			d.name          AS donor_name,
			d.email         AS donor_email,
			d.b_gr          AS donor_blood_group,
			bm.name         AS manager_name,
			b.id            AS blood_id,
			b.b_gr          AS blood_group,
			b.unit          AS units,
			b.expiry_date::text AS expiry_date
		 FROM       donation        don
		 JOIN       donor           d   ON d.id   = don.donor_id
		 JOIN       blood_manager   bm  ON bm.id  = don.manager_id
		 JOIN       blood           b   ON b.id   = don.blood_id
		 ORDER BY   don.donation_date DESC`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch donation history: " + err.Error()})
		return
	}
	defer rows.Close()

	type Donation struct {
		ID              int    `json:"id"`
		DonationDate    string `json:"donation_date"`
		Status          string `json:"status"`
		DonorName       string `json:"donor_name"`
		DonorEmail      string `json:"donor_email"`
		DonorBloodGroup string `json:"donor_blood_group"`
		ManagerName     string `json:"manager_name"`
		BloodID         int    `json:"blood_id"`
		BloodGroup      string `json:"blood_group"`
		Units           int    `json:"units"`
		ExpiryDate      string `json:"expiry_date"`
	}

	var donations []Donation
	for rows.Next() {
		var don Donation
		if err := rows.Scan(
			&don.ID, &don.DonationDate, &don.Status,
			&don.DonorName, &don.DonorEmail, &don.DonorBloodGroup,
			&don.ManagerName,
			&don.BloodID, &don.BloodGroup, &don.Units, &don.ExpiryDate,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to read donation data",
				"details": err.Error(),
			})
			return
		}
		donations = append(donations, don)
	}

	if donations == nil {
		donations = []Donation{}
	}

	c.JSON(http.StatusOK, gin.H{
		"donations": donations,
		"count":     len(donations),
	})
}

// VIEW BLOOD INVENTORY

func GetBloodInventory(c *gin.Context) {
	ctx := context.Background()

	rows, err := db.Pool.Query(ctx,
		`SELECT
			b.id,
			b.b_gr,
			b.unit,
			b.status,
			b.collected_date::text,
			b.expiry_date::text,
			(b.expiry_date - CURRENT_DATE)  AS days_until_expiry,
			CASE
				WHEN b.expiry_date <= CURRENT_DATE          THEN 'EXPIRED'
				WHEN (b.expiry_date - CURRENT_DATE) <= 7    THEN 'CRITICAL'
				WHEN (b.expiry_date - CURRENT_DATE) <= 14   THEN 'WARNING'
				ELSE                                             'OK'
			END                             AS expiry_alert,
			COALESCE(d.name, 'Unknown')     AS donor_name,
			COALESCE(d.email, '-')          AS donor_email
		 FROM       blood b
		 LEFT JOIN  donor d ON d.id = b.donor_id
		 WHERE      b.status IN ('available', 'reserved')
		   AND      b.unit > 0
		 ORDER BY
		 	CASE
				WHEN (b.expiry_date - CURRENT_DATE) <= 7  THEN 1
				WHEN (b.expiry_date - CURRENT_DATE) <= 14 THEN 2
				ELSE                                            3
			END,
			b.expiry_date ASC`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch blood inventory: " + err.Error()})
		return
	}
	defer rows.Close()

	type BloodUnit struct {
		ID              int    `json:"id"`
		BloodGroup      string `json:"blood_group"`
		Units           int    `json:"units"`
		Status          string `json:"status"`
		CollectedDate   string `json:"collected_date"`
		ExpiryDate      string `json:"expiry_date"`
		DaysUntilExpiry int    `json:"days_until_expiry"`
		ExpiryAlert     string `json:"expiry_alert"`
		DonorName       string `json:"donor_name"`
		DonorEmail      string `json:"donor_email"`
	}

	var inventory []BloodUnit
	for rows.Next() {
		var b BloodUnit
		if err := rows.Scan(
			&b.ID, &b.BloodGroup, &b.Units, &b.Status,
			&b.CollectedDate, &b.ExpiryDate, &b.DaysUntilExpiry, &b.ExpiryAlert,
			&b.DonorName, &b.DonorEmail,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to read blood inventory data",
				"details": err.Error(),
			})
			return
		}
		inventory = append(inventory, b)
	}

	if inventory == nil {
		inventory = []BloodUnit{}
	}

	c.JSON(http.StatusOK, gin.H{
		"inventory": inventory,
		"count":     len(inventory),
	})
}

// VIEW PENDING BLOOD REQUESTS
func GetPendingBloodRequests(c *gin.Context) {
	ctx := context.Background()

	rows, err := db.Pool.Query(ctx,
		`SELECT
			br.id,
			br.request_date::text,
			br.quantity_needed,
			br.status,
			p.name          AS patient_name,
			p.email         AS patient_email,
			p.b_gr          AS patient_blood_group,
			d.name          AS doctor_name,
			d.email         AS doctor_email
		 FROM   blood_request br
		 JOIN   patient       p ON p.id = br.patient_id
		 JOIN   doctor        d ON d.id = br.doctor_id
		 WHERE  br.status = 'pending'
		 ORDER BY br.request_date ASC`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch pending blood requests: " + err.Error()})
		return
	}
	defer rows.Close()

	type BloodRequest struct {
		ID                int    `json:"id"`
		RequestDate       string `json:"request_date"`
		QuantityNeeded    int    `json:"quantity_needed"`
		Status            string `json:"status"`
		PatientName       string `json:"patient_name"`
		PatientEmail      string `json:"patient_email"`
		PatientBloodGroup string `json:"patient_blood_group"`
		DoctorName        string `json:"doctor_name"`
		DoctorEmail       string `json:"doctor_email"`
	}

	var requests []BloodRequest
	for rows.Next() {
		var r BloodRequest
		if err := rows.Scan(
			&r.ID, &r.RequestDate, &r.QuantityNeeded, &r.Status,
			&r.PatientName, &r.PatientEmail, &r.PatientBloodGroup,
			&r.DoctorName, &r.DoctorEmail,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to read blood request data",
				"details": err.Error(),
			})
			return
		}
		requests = append(requests, r)
	}

	if requests == nil {
		requests = []BloodRequest{}
	}

	c.JSON(http.StatusOK, gin.H{
		"pending_requests": requests,
		"count":            len(requests),
	})
}

// VIEW EXPIRED BLOOD UNITS
func GetExpiredBlood(c *gin.Context) {
	ctx := context.Background()

	rows, err := db.Pool.Query(ctx,
		`SELECT
			b.id,
			b.b_gr,
			b.unit,
			b.status,
			b.collected_date::text,
			b.expiry_date::text,
			(CURRENT_DATE - b.expiry_date)  AS days_expired,
			COALESCE(d.name, 'Unknown')     AS donor_name,
			COALESCE(d.email, '-')          AS donor_email
		 FROM       blood b
		 LEFT JOIN  donor d ON d.id = b.donor_id
		 WHERE      b.expiry_date <= CURRENT_DATE
		    OR      b.status = 'expired'
		 ORDER BY   b.expiry_date ASC`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch expired blood: " + err.Error()})
		return
	}
	defer rows.Close()

	type ExpiredBlood struct {
		ID            int    `json:"id"`
		BloodGroup    string `json:"blood_group"`
		Units         int    `json:"units"`
		Status        string `json:"status"`
		CollectedDate string `json:"collected_date"`
		ExpiryDate    string `json:"expiry_date"`
		DaysExpired   int    `json:"days_expired"`
		DonorName     string `json:"donor_name"`
		DonorEmail    string `json:"donor_email"`
	}

	var expired []ExpiredBlood
	for rows.Next() {
		var e ExpiredBlood
		if err := rows.Scan(
			&e.ID, &e.BloodGroup, &e.Units, &e.Status,
			&e.CollectedDate, &e.ExpiryDate, &e.DaysExpired,
			&e.DonorName, &e.DonorEmail,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to read expired blood data",
				"details": err.Error(),
			})
			return
		}
		expired = append(expired, e)
	}

	if expired == nil {
		expired = []ExpiredBlood{}
	}

	c.JSON(http.StatusOK, gin.H{
		"expired_blood": expired,
		"count":         len(expired),
	})
}
