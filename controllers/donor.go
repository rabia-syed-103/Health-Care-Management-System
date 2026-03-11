package controllers

import (
	"context"
	"hospital-backend/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AddDonorRequest struct {
	Name       string `json:"name"        binding:"required"`
	PNo        string `json:"p_no"        binding:"required"`
	Email      string `json:"email"       binding:"required"`
	BGr        string `json:"b_gr"        binding:"required"`
	LastDonate string `json:"last_donate" binding:"required"`
}

type EditDonorRequest struct {
	Name       string `json:"name"`
	PNo        string `json:"p_no"`
	BGr        string `json:"b_gr"`
	LastDonate string `json:"last_donate"`
}

func AddDonor(c *gin.Context) {
	var req AddDonorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	ctx := context.Background()
	var id int
	err := db.Pool.QueryRow(ctx,
		`INSERT INTO donor (name, p_no, email, b_gr, last_donate)
		 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		req.Name, req.PNo, req.Email, req.BGr, req.LastDonate,
	).Scan(&id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add donor: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":     "Donor added successfully",
		"id":          id,
		"name":        req.Name,
		"email":       req.Email,
		"b_gr":        req.BGr,
		"p_no":        req.PNo,
		"last_donate": req.LastDonate,
	})
}

func GetAllDonors(c *gin.Context) {
	ctx := context.Background()
	rows, err := db.Pool.Query(ctx,
		`SELECT id, name, p_no, email, b_gr, last_donate FROM donor ORDER BY id`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch donors"})
		return
	}
	defer rows.Close()

	type Donor struct {
		ID         int    `json:"id"`
		Name       string `json:"name"`
		PNo        string `json:"p_no"`
		Email      string `json:"email"`
		BGr        string `json:"b_gr"`
		LastDonate string `json:"last_donate"`
	}

	var donors []Donor
	for rows.Next() {
		var d Donor
		if err := rows.Scan(&d.ID, &d.Name, &d.PNo, &d.Email, &d.BGr, &d.LastDonate); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read donor data"})
			return
		}
		donors = append(donors, d)
	}

	if donors == nil {
		donors = []Donor{}
	}

	c.JSON(http.StatusOK, gin.H{
		"donors": donors,
		"count":  len(donors),
	})
}

func GetDonor(c *gin.Context) {
	email := c.Param("email")
	ctx := context.Background()

	var donor struct {
		ID         int    `json:"id"`
		Name       string `json:"name"`
		PNo        string `json:"p_no"`
		Email      string `json:"email"`
		BGr        string `json:"b_gr"`
		LastDonate string `json:"last_donate"`
	}

	err := db.Pool.QueryRow(ctx,
		`SELECT id, name, p_no, email, b_gr, last_donate FROM donor WHERE email = $1`, email,
	).Scan(&donor.ID, &donor.Name, &donor.PNo, &donor.Email, &donor.BGr, &donor.LastDonate)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Donor not found"})
		return
	}

	c.JSON(http.StatusOK, donor)
}

func EditDonor(c *gin.Context) {
	email := c.Param("email")
	var req EditDonorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	ctx := context.Background()
	result, err := db.Pool.Exec(ctx,
		`UPDATE donor SET
			name        = COALESCE(NULLIF($1, ''), name),
			p_no        = COALESCE(NULLIF($2, ''), p_no),
			b_gr        = COALESCE(NULLIF($3, ''), b_gr),
			last_donate = CASE WHEN $4 = '' THEN last_donate ELSE $4::date END
		 WHERE email = $5`,
		req.Name, req.PNo, req.BGr, req.LastDonate, email,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update donor: " + err.Error()})
		return
	}
	if result.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Donor not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Donor updated successfully"})
}

func DeleteDonor(c *gin.Context) {
	email := c.Param("email")
	ctx := context.Background()

	result, err := db.Pool.Exec(ctx,
		`DELETE FROM donor WHERE email = $1`, email,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete donor: " + err.Error()})
		return
	}
	if result.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Donor not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Donor deleted successfully"})
}
