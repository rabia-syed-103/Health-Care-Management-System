package controllers

import (
	"context"
	"hospital-management/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AddERPatientEntryRequest struct {
	ERPatientID int    `json:"er_patient_id" binding:"required"`
	DoctorID    int    `json:"doctor_id"     binding:"required"`
	ERID        int    `json:"er_id"         binding:"required"`
	Status      string `json:"status"`
}

type AddERPatientRequest struct {
	Name        string `json:"name"         binding:"required"`
	Age         int    `json:"age"          binding:"required"`
	PNo         string `json:"p_no"         binding:"required"`
	ArrivalTime string `json:"arrival_time" binding:"required"`
}

func GetAllAppointments(c *gin.Context) {
	ctx := context.Background()

	rows, err := db.Pool.Query(ctx,
		`SELECT
			a.id,
			a.date::text,
			a.time::text,
			a.status,
			COALESCE(a.ot_id::text, 'None') AS ot_id,
			p.name        AS patient_name,
			p.email       AS patient_email,
			d.name        AS doctor_name,
			r.name        AS receptionist_name
		 FROM       appointment  a
		 JOIN       patient      p ON p.id = a.patient_id
		 JOIN       doctor       d ON d.id = a.doctor_id
		 JOIN       receptionist r ON r.id = a.receptionist_id
		 ORDER BY   a.date DESC, a.time DESC`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch appointments: " + err.Error()})
		return
	}
	defer rows.Close()

	type Appointment struct {
		ID               int    `json:"id"`
		Date             string `json:"date"`
		Time             string `json:"time"`
		Status           string `json:"status"`
		OTID             string `json:"ot_id"`
		PatientName      string `json:"patient_name"`
		PatientEmail     string `json:"patient_email"`
		DoctorName       string `json:"doctor_name"`
		ReceptionistName string `json:"receptionist_name"`
	}

	var appointments []Appointment
	for rows.Next() {
		var a Appointment
		if err := rows.Scan(
			&a.ID, &a.Date, &a.Time, &a.Status, &a.OTID,
			&a.PatientName, &a.PatientEmail,
			&a.DoctorName, &a.ReceptionistName,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to read appointment data",
				"details": err.Error(),
			})
			return
		}
		appointments = append(appointments, a)
	}

	if appointments == nil {
		appointments = []Appointment{}
	}

	c.JSON(http.StatusOK, gin.H{
		"appointments": appointments,
		"count":        len(appointments),
	})
}

func GetAllERPatients(c *gin.Context) {
	ctx := context.Background()

	rows, err := db.Pool.Query(ctx,
		`SELECT id, name, age, p_no, arrival_time::text
		 FROM   er_patient
		 ORDER BY id DESC`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch ER patients: " + err.Error()})
		return
	}
	defer rows.Close()

	type ERPatient struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		Age         int    `json:"age"`
		PNo         string `json:"p_no"`
		ArrivalTime string `json:"arrival_time"`
	}

	var patients []ERPatient
	for rows.Next() {
		var p ERPatient
		if err := rows.Scan(&p.ID, &p.Name, &p.Age, &p.PNo, &p.ArrivalTime); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to read ER patient data",
				"details": err.Error(),
			})
			return
		}
		patients = append(patients, p)
	}

	if patients == nil {
		patients = []ERPatient{}
	}

	c.JSON(http.StatusOK, gin.H{
		"er_patients": patients,
		"count":       len(patients),
	})
}

func AddERPatient(c *gin.Context) {
	var req AddERPatientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	ctx := context.Background()
	var id int
	err := db.Pool.QueryRow(ctx,
		`INSERT INTO er_patient (name, age, p_no, arrival_time)
		 VALUES ($1, $2, $3, $4) RETURNING id`,
		req.Name, req.Age, req.PNo, req.ArrivalTime,
	).Scan(&id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add ER patient: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":      "ER patient registered successfully",
		"id":           id,
		"name":         req.Name,
		"age":          req.Age,
		"p_no":         req.PNo,
		"arrival_time": req.ArrivalTime,
	})
}

func AddERPatientEntry(c *gin.Context) {
	var req AddERPatientEntryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	if req.Status == "" {
		req.Status = "waiting"
	}

	ctx := context.Background()

	var patientName string
	if err := db.Pool.QueryRow(ctx,
		`SELECT name FROM er_patient WHERE id = $1`, req.ERPatientID,
	).Scan(&patientName); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ER patient not found"})
		return
	}

	var doctorName string
	if err := db.Pool.QueryRow(ctx,
		`SELECT name FROM doctor WHERE id = $1`, req.DoctorID,
	).Scan(&doctorName); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Doctor not found"})
		return
	}

	var shiftDate, shiftTime string
	if err := db.Pool.QueryRow(ctx,
		`SELECT date::text, time::text FROM er_shift WHERE id = $1`, req.ERID,
	).Scan(&shiftDate, &shiftTime); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ER shift not found"})
		return
	}

	var assignedCount int
	if err := db.Pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM er_shift_doctor
		 WHERE doctor_id = $1 AND er_id = $2`,
		req.DoctorID, req.ERID,
	).Scan(&assignedCount); err != nil || assignedCount == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Doctor is not assigned to this ER shift",
		})
		return
	}

	var entryID int
	err := db.Pool.QueryRow(ctx,
		`INSERT INTO er_patient_entry (er_patient_id, doctor_id, er_id, status)
		 VALUES ($1, $2, $3, $4) RETURNING id`,
		req.ERPatientID, req.DoctorID, req.ERID, req.Status,
	).Scan(&entryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add ER patient entry: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "ER patient entry created successfully",
		"entry_id":   entryID,
		"patient":    patientName,
		"doctor":     doctorName,
		"shift_date": shiftDate,
		"shift_time": shiftTime,
		"status":     req.Status,
	})
}

func GetAllERPatientEntries(c *gin.Context) {
	ctx := context.Background()

	rows, err := db.Pool.Query(ctx,
		`SELECT
			e.id,
			e.status,
			ep.name              AS patient_name,
			ep.age               AS patient_age,
			ep.p_no              AS patient_phone,
			ep.arrival_time::text,
			d.name               AS doctor_name,
			ers.date::text       AS shift_date,
			ers.time::text       AS shift_time
		 FROM       er_patient_entry e
		 JOIN       er_patient       ep  ON ep.id  = e.er_patient_id
		 JOIN       doctor           d   ON d.id   = e.doctor_id
		 JOIN       er_shift         ers ON ers.id = e.er_id
		 ORDER BY   ers.date DESC, ers.time DESC`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch ER entries: " + err.Error()})
		return
	}
	defer rows.Close()

	type EREntry struct {
		ID           int    `json:"id"`
		Status       string `json:"status"`
		PatientName  string `json:"patient_name"`
		PatientAge   int    `json:"patient_age"`
		PatientPhone string `json:"patient_phone"`
		ArrivalTime  string `json:"arrival_time"`
		DoctorName   string `json:"doctor_name"`
		ShiftDate    string `json:"shift_date"`
		ShiftTime    string `json:"shift_time"`
	}

	var entries []EREntry
	for rows.Next() {
		var e EREntry
		if err := rows.Scan(
			&e.ID, &e.Status,
			&e.PatientName, &e.PatientAge, &e.PatientPhone, &e.ArrivalTime,
			&e.DoctorName,
			&e.ShiftDate, &e.ShiftTime,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to read ER entry data",
				"details": err.Error(),
			})
			return
		}
		entries = append(entries, e)
	}

	if entries == nil {
		entries = []EREntry{}
	}

	c.JSON(http.StatusOK, gin.H{
		"er_entries": entries,
		"count":      len(entries),
	})
}
