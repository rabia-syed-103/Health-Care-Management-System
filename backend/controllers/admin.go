package controllers

import (
	"context"
	"hospital-management/db"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

//  STRUCTS

type AddDoctorRequest struct {
	Name           string `json:"name"           binding:"required"`
	Email          string `json:"email"          binding:"required"`
	PNo            string `json:"p_no"           binding:"required"`
	Specialization string `json:"specialization" binding:"required"`
	Password       string `json:"password"       binding:"required"`
}
type EditDoctorRequest struct {
	Name           string `json:"name"`
	PNo            string `json:"p_no"`
	Specialization string `json:"specialization"`
}

type AddStaffRequest struct {
	Name     string `json:"name"     binding:"required"`
	Email    string `json:"email"    binding:"required"`
	PNo      string `json:"p_no"     binding:"required"`
	Password string `json:"password" binding:"required"`
}
type EditStaffRequest struct {
	Name string `json:"name"`
	PNo  string `json:"p_no"`
}

type EditMedicineRequest struct {
	Name       string `json:"name"`
	BatchNo    string `json:"batch_no"`
	Stock      int    `json:"stock"`
	ExpiryDate string `json:"expiry_date"`
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// DOCTORS

func AdminGetAllDoctors(c *gin.Context) {
	ctx := context.Background()
	rows, err := db.Pool.Query(ctx,
		`SELECT id, name, email, p_no, specialization FROM doctor ORDER BY id`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch doctors: " + err.Error()})
		return
	}
	defer rows.Close()

	type Doctor struct {
		ID             int    `json:"id"`
		Name           string `json:"name"`
		Email          string `json:"email"`
		PNo            string `json:"p_no"`
		Specialization string `json:"specialization"`
	}
	var doctors []Doctor
	for rows.Next() {
		var d Doctor
		if err := rows.Scan(&d.ID, &d.Name, &d.Email, &d.PNo, &d.Specialization); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read doctor data: " + err.Error()})
			return
		}
		doctors = append(doctors, d)
	}
	if doctors == nil {
		doctors = []Doctor{}
	}
	c.JSON(http.StatusOK, gin.H{"doctors": doctors, "count": len(doctors)})
}

func AdminAddDoctor(c *gin.Context) {
	var req AddDoctorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}
	hashed, err := hashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	ctx := context.Background()
	var id int
	err = db.Pool.QueryRow(ctx,
		`INSERT INTO doctor (name, email, p_no, specialization, password)
		 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		req.Name, req.Email, req.PNo, req.Specialization, hashed,
	).Scan(&id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add doctor: " + err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message":        "Doctor added successfully",
		"id":             id,
		"name":           req.Name,
		"email":          req.Email,
		"specialization": req.Specialization,
	})
}

func AdminEditDoctor(c *gin.Context) {
	email := c.Param("email")
	var req EditDoctorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}
	ctx := context.Background()
	result, err := db.Pool.Exec(ctx,
		`UPDATE doctor SET
			name           = COALESCE(NULLIF($1, ''), name),
			p_no           = COALESCE(NULLIF($2, ''), p_no),
			specialization = COALESCE(NULLIF($3, ''), specialization)
		 WHERE email = $4`,
		req.Name, req.PNo, req.Specialization, email,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update doctor: " + err.Error()})
		return
	}
	if result.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Doctor not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Doctor updated successfully"})
}

func AdminDeleteDoctor(c *gin.Context) {
	email := c.Param("email")
	ctx := context.Background()
	result, err := db.Pool.Exec(ctx, `DELETE FROM doctor WHERE email = $1`, email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete doctor: " + err.Error()})
		return
	}
	if result.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Doctor not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Doctor deleted successfully"})
}

// RECEPTIONISTS

func AdminGetAllReceptionists(c *gin.Context) {
	ctx := context.Background()
	rows, err := db.Pool.Query(ctx,
		`SELECT id, name, email, p_no FROM receptionist ORDER BY id`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch receptionists: " + err.Error()})
		return
	}
	defer rows.Close()

	type Receptionist struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
		PNo   string `json:"p_no"`
	}
	var list []Receptionist
	for rows.Next() {
		var r Receptionist
		if err := rows.Scan(&r.ID, &r.Name, &r.Email, &r.PNo); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read receptionist data: " + err.Error()})
			return
		}
		list = append(list, r)
	}
	if list == nil {
		list = []Receptionist{}
	}
	c.JSON(http.StatusOK, gin.H{"receptionists": list, "count": len(list)})
}

func AdminAddReceptionist(c *gin.Context) {
	var req AddStaffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}
	hashed, err := hashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	ctx := context.Background()
	var id int
	err = db.Pool.QueryRow(ctx,
		`INSERT INTO receptionist (name, email, p_no, password)
		 VALUES ($1, $2, $3, $4) RETURNING id`,
		req.Name, req.Email, req.PNo, hashed,
	).Scan(&id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add receptionist: " + err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Receptionist added successfully", "id": id, "name": req.Name, "email": req.Email})
}

func AdminEditReceptionist(c *gin.Context) {
	email := c.Param("email")
	var req EditStaffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}
	ctx := context.Background()
	result, err := db.Pool.Exec(ctx,
		`UPDATE receptionist SET
			name = COALESCE(NULLIF($1, ''), name),
			p_no = COALESCE(NULLIF($2, ''), p_no)
		 WHERE email = $3`,
		req.Name, req.PNo, email,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update receptionist: " + err.Error()})
		return
	}
	if result.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Receptionist not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Receptionist updated successfully"})
}

func AdminDeleteReceptionist(c *gin.Context) {
	email := c.Param("email")
	ctx := context.Background()
	result, err := db.Pool.Exec(ctx, `DELETE FROM receptionist WHERE email = $1`, email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete receptionist: " + err.Error()})
		return
	}
	if result.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Receptionist not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Receptionist deleted successfully"})
}

// PHARMACISTS

func AdminGetAllPharmacists(c *gin.Context) {
	ctx := context.Background()
	rows, err := db.Pool.Query(ctx,
		`SELECT id, name, email, p_no FROM pharmacist ORDER BY id`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch pharmacists: " + err.Error()})
		return
	}
	defer rows.Close()

	type Pharmacist struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
		PNo   string `json:"p_no"`
	}
	var list []Pharmacist
	for rows.Next() {
		var p Pharmacist
		if err := rows.Scan(&p.ID, &p.Name, &p.Email, &p.PNo); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read pharmacist data: " + err.Error()})
			return
		}
		list = append(list, p)
	}
	if list == nil {
		list = []Pharmacist{}
	}
	c.JSON(http.StatusOK, gin.H{"pharmacists": list, "count": len(list)})
}

func AdminAddPharmacist(c *gin.Context) {
	var req AddStaffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}
	hashed, err := hashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	ctx := context.Background()
	var id int
	err = db.Pool.QueryRow(ctx,
		`INSERT INTO pharmacist (name, email, p_no, password)
		 VALUES ($1, $2, $3, $4) RETURNING id`,
		req.Name, req.Email, req.PNo, hashed,
	).Scan(&id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add pharmacist: " + err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Pharmacist added successfully", "id": id, "name": req.Name, "email": req.Email})
}

func AdminEditPharmacist(c *gin.Context) {
	email := c.Param("email")
	var req EditStaffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}
	ctx := context.Background()
	result, err := db.Pool.Exec(ctx,
		`UPDATE pharmacist SET
			name = COALESCE(NULLIF($1, ''), name),
			p_no = COALESCE(NULLIF($2, ''), p_no)
		 WHERE email = $3`,
		req.Name, req.PNo, email,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update pharmacist: " + err.Error()})
		return
	}
	if result.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pharmacist not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Pharmacist updated successfully"})
}

func AdminDeletePharmacist(c *gin.Context) {
	email := c.Param("email")
	ctx := context.Background()
	result, err := db.Pool.Exec(ctx, `DELETE FROM pharmacist WHERE email = $1`, email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete pharmacist: " + err.Error()})
		return
	}
	if result.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pharmacist not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Pharmacist deleted successfully"})
}

// BLOOD MANAGERS

func AdminGetAllBloodManagers(c *gin.Context) {
	ctx := context.Background()
	rows, err := db.Pool.Query(ctx,
		`SELECT id, name, email, p_no FROM blood_manager ORDER BY id`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch blood managers: " + err.Error()})
		return
	}
	defer rows.Close()

	type BloodManager struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
		PNo   string `json:"p_no"`
	}
	var list []BloodManager
	for rows.Next() {
		var bm BloodManager
		if err := rows.Scan(&bm.ID, &bm.Name, &bm.Email, &bm.PNo); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read blood manager data: " + err.Error()})
			return
		}
		list = append(list, bm)
	}
	if list == nil {
		list = []BloodManager{}
	}
	c.JSON(http.StatusOK, gin.H{"blood_managers": list, "count": len(list)})
}

func AdminAddBloodManager(c *gin.Context) {
	var req AddStaffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}
	hashed, err := hashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	ctx := context.Background()
	var id int
	err = db.Pool.QueryRow(ctx,
		`INSERT INTO blood_manager (name, email, p_no, password)
		 VALUES ($1, $2, $3, $4) RETURNING id`,
		req.Name, req.Email, req.PNo, hashed,
	).Scan(&id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add blood manager: " + err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Blood manager added successfully", "id": id, "name": req.Name, "email": req.Email})
}

func AdminEditBloodManager(c *gin.Context) {
	email := c.Param("email")
	var req EditStaffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}
	ctx := context.Background()
	result, err := db.Pool.Exec(ctx,
		`UPDATE blood_manager SET
			name = COALESCE(NULLIF($1, ''), name),
			p_no = COALESCE(NULLIF($2, ''), p_no)
		 WHERE email = $3`,
		req.Name, req.PNo, email,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update blood manager: " + err.Error()})
		return
	}
	if result.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Blood manager not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Blood manager updated successfully"})
}

func AdminDeleteBloodManager(c *gin.Context) {
	email := c.Param("email")
	ctx := context.Background()
	result, err := db.Pool.Exec(ctx, `DELETE FROM blood_manager WHERE email = $1`, email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete blood manager: " + err.Error()})
		return
	}
	if result.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Blood manager not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Blood manager deleted successfully"})
}
