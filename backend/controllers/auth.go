package controllers

import (
	"context"
	"hospital-management/db"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Name           string `json:"name"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	PNo            string `json:"p_no"`
	Role           string `json:"role"`
	Specialization string `json:"specialization"` // only for doctor
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func generateToken(userID int, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

// REGISTER
func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if req.Name == "" || req.Email == "" || req.Password == "" || req.Role == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name, email, password, and role are required"})
		return
	}
	if req.Role != "admin" && req.PNo == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "p_no is required for this role"})
		return
	}
	if req.Role == "doctor" && req.Specialization == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "specialization is required for doctors"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	var userID int

	switch req.Role {
	case "admin":

		err = db.Pool.QueryRow(context.Background(),
			`INSERT INTO admin (name, email, p_no, password) VALUES ($1, $2, $3, $4) RETURNING id`,
			req.Name, req.Email, req.PNo, string(hashedPassword),
		).Scan(&userID)

	case "doctor":
		err = db.Pool.QueryRow(context.Background(),
			`INSERT INTO doctor (name, email, p_no, specialization, password) VALUES ($1, $2, $3, $4, $5) RETURNING id`,
			req.Name, req.Email, req.PNo, req.Specialization, string(hashedPassword),
		).Scan(&userID)

	case "receptionist":
		err = db.Pool.QueryRow(context.Background(),
			`INSERT INTO receptionist (name, email, p_no, password) VALUES ($1, $2, $3, $4) RETURNING id`,
			req.Name, req.Email, req.PNo, string(hashedPassword),
		).Scan(&userID)

	case "pharmacist":
		err = db.Pool.QueryRow(context.Background(),
			`INSERT INTO pharmacist (name, email, p_no, password) VALUES ($1, $2, $3, $4) RETURNING id`,
			req.Name, req.Email, req.PNo, string(hashedPassword),
		).Scan(&userID)

	case "blood_manager":
		err = db.Pool.QueryRow(context.Background(),
			`INSERT INTO blood_manager (name, email, p_no, password) VALUES ($1, $2, $3, $4) RETURNING id`,
			req.Name, req.Email, req.PNo, string(hashedPassword),
		).Scan(&userID)

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role. Must be: admin, doctor, receptionist, pharmacist, blood_manager"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user: " + err.Error()})
		return
	}

	token, err := generateToken(userID, req.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"token":   token,
		"role":    req.Role,
		"id":      userID,
	})
}

//LOGIN

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	var storedPassword string
	var userID int
	var queryErr error

	switch req.Role {
	case "admin":
		queryErr = db.Pool.QueryRow(context.Background(),
			`SELECT id, password FROM admin WHERE email = $1`, req.Email,
		).Scan(&userID, &storedPassword)

	case "doctor":
		queryErr = db.Pool.QueryRow(context.Background(),
			`SELECT id, password FROM doctor WHERE email = $1`, req.Email,
		).Scan(&userID, &storedPassword)

	case "receptionist":
		queryErr = db.Pool.QueryRow(context.Background(),
			`SELECT id, password FROM receptionist WHERE email = $1`, req.Email,
		).Scan(&userID, &storedPassword)

	case "pharmacist":
		queryErr = db.Pool.QueryRow(context.Background(),
			`SELECT id, password FROM pharmacist WHERE email = $1`, req.Email,
		).Scan(&userID, &storedPassword)

	case "blood_manager":
		queryErr = db.Pool.QueryRow(context.Background(),
			`SELECT id, password FROM blood_manager WHERE email = $1`, req.Email,
		).Scan(&userID, &storedPassword)

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role"})
		return
	}

	if queryErr != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	token, err := generateToken(userID, req.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
		"role":    req.Role,
		"id":      userID,
	})
}
