package controllers

import (
	"context"
	"hospital-backend/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AddMedicineRequest struct {
	BatchNo    string `json:"batch_no"     binding:"required"`
	Name       string `json:"name"         binding:"required"`
	Stock      int    `json:"stock"        binding:"required,min=1"`
	ExpiryDate string `json:"expiry_date"  binding:"required"`
}

func GetPendingPrescriptions(c *gin.Context) {
	ctx := context.Background()

	rows, err := db.Pool.Query(ctx,
		`SELECT
			pr.id                   AS prescription_id,
			pr.date::text           AS prescription_date,
			p.name                  AS patient_name,
			p.email                 AS patient_email,
			d.name                  AS doctor_name,
			COUNT(pm.id)            AS total_medicines,
			STRING_AGG(
				m.name || ' (qty: ' || pm.quantity::text || ')',
				', ' ORDER BY m.name
			)                       AS medicines_list
		 FROM       prescription    pr
		 JOIN       patient         p  ON p.id  = pr.patient_id
		 JOIN       doctor          d  ON d.id  = pr.doctor_id
		 JOIN       prescription_medicine pm ON pm.prescription_id = pr.id
		 JOIN       medicine        m  ON m.id  = pm.medicine_id
		 WHERE      pr.id NOT IN (
		                SELECT DISTINCT prescription_id FROM dispensing
		            )
		 GROUP BY   pr.id, pr.date, p.name, p.email, d.name
		 ORDER BY   pr.date DESC`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch pending prescriptions: " + err.Error()})
		return
	}
	defer rows.Close()

	type Prescription struct {
		ID             int    `json:"prescription_id"`
		Date           string `json:"prescription_date"`
		PatientName    string `json:"patient_name"`
		PatientEmail   string `json:"patient_email"`
		DoctorName     string `json:"doctor_name"`
		TotalMedicines int    `json:"total_medicines"`
		MedicinesList  string `json:"medicines_list"`
	}

	var prescriptions []Prescription
	for rows.Next() {
		var pr Prescription
		if err := rows.Scan(
			&pr.ID, &pr.Date,
			&pr.PatientName, &pr.PatientEmail,
			&pr.DoctorName,
			&pr.TotalMedicines, &pr.MedicinesList,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to read prescription data",
				"details": err.Error(),
			})
			return
		}
		prescriptions = append(prescriptions, pr)
	}

	if prescriptions == nil {
		prescriptions = []Prescription{}
	}

	c.JSON(http.StatusOK, gin.H{
		"pending_prescriptions": prescriptions,
		"count":                 len(prescriptions),
	})
}

func GetMedicineStock(c *gin.Context) {
	ctx := context.Background()

	rows, err := db.Pool.Query(ctx,
		`SELECT
			id,
			batch_no,
			name,
			stock,
			expiry_date::text,
			CASE
				WHEN stock = 0                          THEN 'OUT OF STOCK'
				WHEN expiry_date <= CURRENT_DATE        THEN 'EXPIRED'
				WHEN expiry_date <= CURRENT_DATE + 7    THEN 'CRITICAL'
				WHEN expiry_date <= CURRENT_DATE + 30   THEN 'EXPIRING SOON'
				ELSE                                         'OK'
			END AS alert
		 FROM  medicine
		 ORDER BY
		 	CASE
				WHEN stock = 0                          THEN 1
				WHEN expiry_date <= CURRENT_DATE        THEN 2
				WHEN expiry_date <= CURRENT_DATE + 7    THEN 3
				WHEN expiry_date <= CURRENT_DATE + 30   THEN 4
				ELSE                                         5
			END,
			expiry_date ASC`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch medicine stock: " + err.Error()})
		return
	}
	defer rows.Close()

	type Medicine struct {
		ID         int    `json:"id"`
		BatchNo    string `json:"batch_no"`
		Name       string `json:"name"`
		Stock      int    `json:"stock"`
		ExpiryDate string `json:"expiry_date"`
		Alert      string `json:"alert"`
	}

	var medicines []Medicine
	for rows.Next() {
		var m Medicine
		if err := rows.Scan(
			&m.ID, &m.BatchNo, &m.Name,
			&m.Stock, &m.ExpiryDate, &m.Alert,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to read medicine data",
				"details": err.Error(),
			})
			return
		}
		medicines = append(medicines, m)
	}

	if medicines == nil {
		medicines = []Medicine{}
	}

	c.JSON(http.StatusOK, gin.H{
		"medicines": medicines,
		"count":     len(medicines),
	})
}

func AddMedicine(c *gin.Context) {
	var req AddMedicineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	ctx := context.Background()

	var existingID int
	var existingStock int
	err := db.Pool.QueryRow(ctx,
		`SELECT id, stock FROM medicine WHERE batch_no = $1`, req.BatchNo,
	).Scan(&existingID, &existingStock)

	if err == nil {
		newStock := existingStock + req.Stock
		_, err = db.Pool.Exec(ctx,
			`UPDATE medicine
			 SET stock       = $1,
			     expiry_date = $2
			 WHERE id = $3`,
			newStock, req.ExpiryDate, existingID,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to restock medicine: " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":     "Medicine restocked successfully",
			"action":      "updated",
			"id":          existingID,
			"batch_no":    req.BatchNo,
			"name":        req.Name,
			"added_stock": req.Stock,
			"total_stock": newStock,
			"new_expiry":  req.ExpiryDate,
		})
		return
	}

	var newID int
	err = db.Pool.QueryRow(ctx,
		`INSERT INTO medicine (batch_no, name, stock, expiry_date)
		 VALUES ($1, $2, $3, $4) RETURNING id`,
		req.BatchNo, req.Name, req.Stock, req.ExpiryDate,
	).Scan(&newID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add medicine: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":     "New medicine added successfully",
		"action":      "created",
		"id":          newID,
		"batch_no":    req.BatchNo,
		"name":        req.Name,
		"stock":       req.Stock,
		"expiry_date": req.ExpiryDate,
	})
}

func GetMedicineByName(c *gin.Context) {
	name := c.Param("name")
	ctx := context.Background()

	rows, err := db.Pool.Query(ctx,
		`SELECT
			id,
			batch_no,
			name,
			stock,
			expiry_date::text,
			CASE
				WHEN stock = 0                        THEN 'OUT OF STOCK'
				WHEN expiry_date <= CURRENT_DATE      THEN 'EXPIRED'
				WHEN expiry_date <= CURRENT_DATE + 7  THEN 'CRITICAL'
				WHEN expiry_date <= CURRENT_DATE + 30 THEN 'EXPIRING SOON'
				ELSE                                       'OK'
			END AS alert
		 FROM  medicine
		 WHERE LOWER(name) = LOWER($1)
		 ORDER BY expiry_date ASC`,
		name,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search medicine: " + err.Error()})
		return
	}
	defer rows.Close()

	type Medicine struct {
		ID         int    `json:"id"`
		BatchNo    string `json:"batch_no"`
		Name       string `json:"name"`
		Stock      int    `json:"stock"`
		ExpiryDate string `json:"expiry_date"`
		Alert      string `json:"alert"`
	}

	var medicines []Medicine
	for rows.Next() {
		var m Medicine
		if err := rows.Scan(&m.ID, &m.BatchNo, &m.Name, &m.Stock, &m.ExpiryDate, &m.Alert); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to read medicine data",
				"details": err.Error(),
			})
			return
		}
		medicines = append(medicines, m)
	}

	if len(medicines) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No medicine found with name: " + name})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"medicines": medicines,
		"count":     len(medicines),
	})
}
