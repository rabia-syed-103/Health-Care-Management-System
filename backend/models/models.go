package models

import "time"

// Admin

type Admin struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"-"`
	P_no     string `json:"p_no"`
}

// Doctor

type Doctor struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	PNo            string `json:"p_no"`
	Specialization string `json:"specialization"`
	Password       string `json:"-"`
}

//Receptionist

type Receptionist struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	PNo      string `json:"p_no"`
	Password string `json:"-"`
}

// Patient

type Patient struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	BGr   string `json:"b_gr"`
	PNo   string `json:"p_no"`
}

// Appointment

type Appointment struct {
	ID             int       `json:"id"`
	ReceptionistID int       `json:"receptionist_id"`
	PatientID      int       `json:"patient_id"`
	DoctorID       int       `json:"doctor_id"`
	Date           time.Time `json:"date"`
	Time           time.Time `json:"time"`
	Status         string    `json:"status"`
	OTID           *int      `json:"ot_id,omitempty"`
}

// Ot

type OT struct {
	ID          int  `json:"id"`
	IsAvailable bool `json:"is_available"`
}

// Pharmacist

type Pharmacist struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	PNo      string `json:"p_no"`
	Password string `json:"-"`
}

// Medicine

type Medicine struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Stock      int       `json:"stock"`
	ExpiryDate time.Time `json:"expiry_date"`
	BatchNo    string    `json:"batch_no"`
}

//Prescription

type Prescription struct {
	ID        int       `json:"id"`
	DoctorID  int       `json:"doctor_id"`
	PatientID int       `json:"patient_id"`
	Date      time.Time `json:"date"`
}

// PrescriptionMedicine

type PrescriptionMedicine struct {
	ID             int `json:"id"`
	PrescriptionID int `json:"prescription_id"`
	MedicineID     int `json:"medicine_id"`
	Quantity       int `json:"quantity"`
}

// Dispensing

type Dispensing struct {
	ID             int `json:"id"`
	MedicineID     int `json:"medicine_id"`
	PharmacistID   int `json:"pharmacist_id"`
	PrescriptionID int `json:"prescription_id"`
	Quantity       int `json:"quantity"`
}

// Blood Manager

type BloodManager struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	PNo      string `json:"p_no"`
	Password string `json:"-"`
}

// Blood

type Blood struct {
	ID            int       `json:"id"`
	BGr           string    `json:"b_gr"`
	CollectedDate time.Time `json:"collected_date"`
	Status        string    `json:"status"`
	ExpiryDate    time.Time `json:"expiry_date"`
	Unit          int       `json:"unit"`
}

// Blood Request

type BloodRequest struct {
	ID             int       `json:"id"`
	DoctorID       int       `json:"doctor_id"`
	PatientID      int       `json:"patient_id"`
	QuantityNeeded int       `json:"quantity_needed"`
	Status         string    `json:"status"`
	RequestDate    time.Time `json:"request_date"`
}

// Blood Request Fulfillment

type BloodRequestFulfillment struct {
	ID               int       `json:"id"`
	RequestID        int       `json:"request_id"`
	BloodID          int       `json:"blood_id"`
	ManagerID        int       `json:"manager_id"`
	QuantityProvided int       `json:"quantity_provided"`
	FulfillmentDate  time.Time `json:"fulfillment_date"`
}

// Donor

type Donor struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	PNo        string    `json:"p_no"`
	Email      string    `json:"email"`
	BGr        string    `json:"b_gr"`
	LastDonate time.Time `json:"last_donate"`
}

// Donation

type Donation struct {
	ID           int       `json:"id"`
	DonorID      int       `json:"donor_id"`
	ManagerID    int       `json:"manager_id"`
	BloodID      int       `json:"blood_id"`
	DonationDate time.Time `json:"donation_date"`
	Status       string    `json:"status"`
}

// ER Shift

type ERShift struct {
	ID             int       `json:"id"`
	ReceptionistID *int      `json:"receptionist_id,omitempty"`
	Time           time.Time `json:"time"`
	Date           time.Time `json:"date"`
}

// ER Shift Doctor

type ERShiftDoctor struct {
	ID       int `json:"id"`
	DoctorID int `json:"doctor_id"`
	ERID     int `json:"er_id"`
}

// ER Patient

type ERPatient struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Age         int       `json:"age"`
	PNo         string    `json:"p_no"`
	ArrivalTime time.Time `json:"arrival_time"`
}

// er_patient_entry

type ERPatientEntry struct {
	ID          int    `json:"id"`
	ERPatientID int    `json:"er_patient_id"`
	DoctorID    int    `json:"doctor_id"`
	ERID        int    `json:"er_id"`
	Status      string `json:"status"`
}
