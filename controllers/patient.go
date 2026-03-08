package controllers

import (
	"context"
	"hospital-backend/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ─── STRUCTS ──────────────────────────────────────────────────────────────────

type AddPatientRequest struct {
	Name  string `json:"name"  binding:"required"`
	Email string `json:"email" binding:"required"`
	BGr   string `json:"b_gr"  binding:"required"`
	PNo   string `json:"p_no"  binding:"required"`
}

type EditPatientRequest struct {
	Name string `json:"name"`
	BGr  string `json:"b_gr"`
	PNo  string `json:"p_no"`
}

// ─── ADD PATIENT ──────────────────────────────────────────────────────────────
func AddPatient(c *gin.Context) {
	var req AddPatientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	ctx := context.Background()
	var id int
	err := db.Pool.QueryRow(ctx,
		`INSERT INTO patient (name, email, b_gr, p_no)
		 VALUES ($1, $2, $3, $4) RETURNING id`,
		req.Name, req.Email, req.BGr, req.PNo,
	).Scan(&id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add patient: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Patient added successfully",
		"id":      id,
		"name":    req.Name,
		"email":   req.Email,
		"b_gr":    req.BGr,
		"p_no":    req.PNo,
	})
}

// ─── VIEW ALL PATIENTS ────────────────────────────────────────────────────────
func GetAllPatients(c *gin.Context) {
	ctx := context.Background()
	rows, err := db.Pool.Query(ctx,
		`SELECT id, name, email, b_gr, p_no FROM patient ORDER BY id`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch patients"})
		return
	}
	defer rows.Close()

	type Patient struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
		BGr   string `json:"b_gr"`
		PNo   string `json:"p_no"`
	}

	var patients []Patient
	for rows.Next() {
		var p Patient
		if err := rows.Scan(&p.ID, &p.Name, &p.Email, &p.BGr, &p.PNo); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read patient data"})
			return
		}
		patients = append(patients, p)
	}

	if patients == nil {
		patients = []Patient{}
	}

	c.JSON(http.StatusOK, gin.H{
		"patients": patients,
		"count":    len(patients),
	})
}

// ─── VIEW SINGLE PATIENT BY EMAIL ─────────────────────────────────────────────
func GetPatient(c *gin.Context) {
	email := c.Param("email")
	ctx := context.Background()

	var patient struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
		BGr   string `json:"b_gr"`
		PNo   string `json:"p_no"`
	}

	err := db.Pool.QueryRow(ctx,
		`SELECT id, name, email, b_gr, p_no FROM patient WHERE email = $1`, email,
	).Scan(&patient.ID, &patient.Name, &patient.Email, &patient.BGr, &patient.PNo)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Patient not found"})
		return
	}

	c.JSON(http.StatusOK, patient)
}

// ─── EDIT PATIENT BY EMAIL ────────────────────────────────────────────────────
// Email is used to identify the patient — you cannot change the email itself
func EditPatient(c *gin.Context) {
	email := c.Param("email")
	var req EditPatientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	ctx := context.Background()
	result, err := db.Pool.Exec(ctx,
		`UPDATE patient SET
			name = COALESCE(NULLIF($1, ''), name),
			b_gr = COALESCE(NULLIF($2, ''), b_gr),
			p_no = COALESCE(NULLIF($3, ''), p_no)
		 WHERE email = $4`,
		req.Name, req.BGr, req.PNo, email,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update patient: " + err.Error()})
		return
	}
	if result.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Patient not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Patient updated successfully"})
}

// ─── DELETE PATIENT BY EMAIL ──────────────────────────────────────────────────
func DeletePatient(c *gin.Context) {
	email := c.Param("email")
	ctx := context.Background()

	result, err := db.Pool.Exec(ctx,
		`DELETE FROM patient WHERE email = $1`, email,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete patient: " + err.Error()})
		return
	}
	if result.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Patient not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Patient deleted successfully"})
}
