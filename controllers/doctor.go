package controllers

import (
	"context"
	"hospital-backend/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ─── VIEW MY APPOINTMENTS / SCHEDULE ─────────────────────────────────────────
// Doctor sees only their own appointments using their doctor_id from JWT token
func GetMyAppointments(c *gin.Context) {
	// Get doctor_id from JWT token set by AuthMiddleware
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

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
			p.b_gr        AS patient_blood_group,
			p.p_no        AS patient_phone,
			r.name        AS receptionist_name
		 FROM       appointment  a
		 JOIN       patient      p ON p.id = a.patient_id
		 JOIN       receptionist r ON r.id = a.receptionist_id
		 WHERE      a.doctor_id = $1
		 ORDER BY   a.date DESC, a.time DESC`,
		userID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch appointments: " + err.Error()})
		return
	}
	defer rows.Close()

	type Appointment struct {
		ID                int    `json:"id"`
		Date              string `json:"date"`
		Time              string `json:"time"`
		Status            string `json:"status"`
		OTID              string `json:"ot_id"`
		PatientName       string `json:"patient_name"`
		PatientEmail      string `json:"patient_email"`
		PatientBloodGroup string `json:"patient_blood_group"`
		PatientPhone      string `json:"patient_phone"`
		ReceptionistName  string `json:"receptionist_name"`
	}

	var appointments []Appointment
	for rows.Next() {
		var a Appointment
		if err := rows.Scan(
			&a.ID, &a.Date, &a.Time, &a.Status, &a.OTID,
			&a.PatientName, &a.PatientEmail, &a.PatientBloodGroup, &a.PatientPhone,
			&a.ReceptionistName,
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

// ─── VIEW PATIENT HISTORY ─────────────────────────────────────────────────────
// Doctor looks up a patient by email and sees all their:
//   - Basic info
//   - All appointments (with any doctor)
//   - All prescriptions + medicines prescribed
func GetPatientHistory(c *gin.Context) {
	email := c.Param("email")
	ctx := context.Background()

	// ── Fetch patient basic info ──────────────────────────────────────────────
	var patient struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
		BGr   string `json:"b_gr"`
		PNo   string `json:"p_no"`
	}
	if err := db.Pool.QueryRow(ctx,
		`SELECT id, name, email, b_gr, p_no FROM patient WHERE email = $1`, email,
	).Scan(&patient.ID, &patient.Name, &patient.Email, &patient.BGr, &patient.PNo); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Patient not found"})
		return
	}

	// ── Fetch all appointments for this patient ───────────────────────────────
	type AppointmentRow struct {
		ID         int    `json:"id"`
		Date       string `json:"date"`
		Time       string `json:"time"`
		Status     string `json:"status"`
		OTID       string `json:"ot_id"`
		DoctorName string `json:"doctor_name"`
	}

	apptRows, err := db.Pool.Query(ctx,
		`SELECT
			a.id,
			a.date::text,
			a.time::text,
			a.status,
			COALESCE(a.ot_id::text, 'None') AS ot_id,
			d.name AS doctor_name
		 FROM   appointment a
		 JOIN   doctor      d ON d.id = a.doctor_id
		 WHERE  a.patient_id = $1
		 ORDER BY a.date DESC, a.time DESC`,
		patient.ID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch appointment history"})
		return
	}
	defer apptRows.Close()

	var appointments []AppointmentRow
	for apptRows.Next() {
		var a AppointmentRow
		if err := apptRows.Scan(&a.ID, &a.Date, &a.Time, &a.Status, &a.OTID, &a.DoctorName); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to read appointment history",
				"details": err.Error(),
			})
			return
		}
		appointments = append(appointments, a)
	}
	if appointments == nil {
		appointments = []AppointmentRow{}
	}

	// ── Fetch all prescriptions + medicines for this patient ──────────────────
	type MedicineRow struct {
		MedicineName string `json:"medicine_name"`
		Quantity     int    `json:"quantity"`
	}
	type PrescriptionRow struct {
		ID         int           `json:"id"`
		Date       string        `json:"date"`
		DoctorName string        `json:"doctor_name"`
		Medicines  []MedicineRow `json:"medicines"`
	}

	prescRows, err := db.Pool.Query(ctx,
		`SELECT
			pr.id,
			pr.date::text,
			d.name  AS doctor_name
		 FROM   prescription pr
		 JOIN   doctor       d ON d.id = pr.doctor_id
		 WHERE  pr.patient_id = $1
		 ORDER BY pr.date DESC`,
		patient.ID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch prescriptions"})
		return
	}
	defer prescRows.Close()

	var prescriptions []PrescriptionRow
	for prescRows.Next() {
		var p PrescriptionRow
		if err := prescRows.Scan(&p.ID, &p.Date, &p.DoctorName); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to read prescription data",
				"details": err.Error(),
			})
			return
		}

		// Fetch medicines for this prescription
		medRows, err := db.Pool.Query(ctx,
			`SELECT m.name, pm.quantity
			 FROM   prescription_medicine pm
			 JOIN   medicine              m ON m.id = pm.medicine_id
			 WHERE  pm.prescription_id = $1`,
			p.ID,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch prescription medicines"})
			return
		}

		var medicines []MedicineRow
		for medRows.Next() {
			var m MedicineRow
			if err := medRows.Scan(&m.MedicineName, &m.Quantity); err != nil {
				medRows.Close()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read medicine data"})
				return
			}
			medicines = append(medicines, m)
		}
		medRows.Close()

		if medicines == nil {
			medicines = []MedicineRow{}
		}
		p.Medicines = medicines
		prescriptions = append(prescriptions, p)
	}

	if prescriptions == nil {
		prescriptions = []PrescriptionRow{}
	}

	c.JSON(http.StatusOK, gin.H{
		"patient":       patient,
		"appointments":  appointments,
		"prescriptions": prescriptions,
	})
}
