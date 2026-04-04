# Hospital Management System
### Group 14 | Advanced Database Management System

| Member | Roll Number |
|--------|-------------|
| Rabia Syed | BSCS24103 |
| Maryam Farooq | BSCS24053 |
| Areej Farhan | BSCS24023 |

---

## Project Overview

A full-stack Hospital, Pharmacy, and Blood Bank Management System that digitizes core hospital operations across five staff roles. The system manages patient registration, appointment scheduling, prescription fulfillment, medicine dispensing, blood donation tracking, and blood request fulfillment all enforced with real database transactions (BEGIN/COMMIT/ROLLBACK), role-based access control, and multi-layer data integrity.

**Problem Solved:** Manual hospital workflows are error-prone and don't prevent concurrent conflicts like double-booking a doctor, over-dispensing medicine, a donor donating twice within 90 days. This system enforces all those rules at the database level, not just the application level.

---

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Frontend | React 18 (Vite), React Router v6, Axios, Recharts, Tailwind CSS v3 |
| Backend | Go 1.21+, Gin (HTTP framework) |
| Database | PostgreSQL 15+ |
| Auth | JWT (golang-jwt/jwt v5), bcrypt password hashing |
| ORM / DB Driver | pgx/v5 with pgxpool (raw SQL) |
| Config | godotenv (.env loading) |
| API Spec | OpenAPI 3.0 (swagger.yaml) |
| Deployment | AWS EC2 (Ubuntu) |

---

## System Architecture

```
Browser (React + Vite)
        │
        │  HTTP/JSON (axios)
        ▼
  Go Backend (Gin)
        │
        ├── middleware/auth.go     → JWT verification
        ├── middleware/auth.go     → RBAC
        ├── controllers/          → Raw SQL transactions 
        │
        │  pgxpool
        ▼
  PostgreSQL 15
        │
        ├── 21 tables
        ├── 5 triggers         
        ├── 4 views         
        └── 4 indexes
```

**Request Lifecycle:**
1. React sends HTTP request with `Authorization: <token>` header
2. `AuthMiddleware` parses and validates the JWT; sets `user_id` and `role` in Gin context
3. `RequireRole` checks the role against the allowed list for that route
4. Controller runs pre-checks , then opens a `BEGIN` transaction
5. SQL operations execute and any failure triggers `defer tx.Rollback()`
6. On success, `tx.Commit()` persists all changes atomically
7. JSON response returned to React

---

## UI Examples

### 1. Login Page
Role-based login with JWT token storage. Each role redirects to its own dashboard after authentication. Required because every protected route checks the JWT role claim. Wrong role gets 401/403.

### 2. Receptionist : Book Appointment
Step 1: Receptionist selects date and time, system queries all doctors with no conflicting pending appointment. Step 2: Receptionist picks a doctor from the available list and enters the patient email. Required to demonstrate the available doctors transaction and prevent double booking.

### 3. Blood Manager : Pending Requests + Fulfill
Lists all pending blood requests with patient name, blood group, and quantity needed. One click opens the fulfillment modal. The system automatically selects the oldest compatible blood unit, locks the row with `SELECT FOR UPDATE`, and commits atomically. Required to show Transaction end-to-end and fulfill the blood-managment component of the project.

---

## Setup & Installation

### Prerequisites

| Tool | Version |
|------|---------|
| Go | 1.21+ |
| PostgreSQL | 15+ |
| Node.js | 20+ |
| Git | Any |

---

### Clone the Repository

```bash
git clone https://github.com/rabia-syed-103/Health-Care-Management-System.git
cd Health-Care-Management-System
```

---

### Database Setup

```bash
# Open psql
psql -U postgres

# Inside psql
CREATE DATABASE hospital_db;
\q

# Run schema
psql -U postgres -d hospital_db -f database/schema.sql

# Run seed data
psql -U postgres -d hospital_db -f database/seed.sql
```

---

### Backend Setup

**1. Configure environment variables**

```bash
# Windows
copy backend\.env.example backend\.env

# Mac / Linux
cp backend/.env.example backend/.env
```

Open `backend/.env` and fill in your values:

```env
DB_HOST=localhost        # PostgreSQL host
DB_PORT=5432             # PostgreSQL port (default)
DB_USER=postgres         # PostgreSQL username
DB_PASSWORD=yourpassword # Your PostgreSQL password
DB_NAME=hospital_db      # Database name created above
JWT_SECRET=anyrandomkey  # Any random string (used to sign tokens)
PORT=8080                # Port the Go server listens on
```

**2. Install Go dependencies**

```bash
cd backend
go mod tidy
```

This installs:
- `github.com/gin-gonic/gin` : HTTP router
- `github.com/jackc/pgx/v5` : PostgreSQL driver with connection pooling
- `github.com/golang-jwt/jwt/v5` : JWT auth
- `github.com/joho/godotenv` : .env loading
- `golang.org/x/crypto` : bcrypt password hashing

**3. Start the backend server**

```bash
go run main.go
```

Output:
```
Connected to PostgreSQL successfully!
[GIN-debug] Listening and serving HTTP on :8080
```

**4. Verify the server**

```
GET http://localhost:8080/health
→ { "status": "Server is running" }
```

---

### Frontend Setup

**1. Install dependencies**

```bash
cd frontend
npm install
```

**2. Configure environment**

```bash
# Windows
copy .env.example .env

# Mac / Linux
cp .env.example .env
```

Open `frontend/.env`:

```env
VITE_API_URL=http://localhost:8080/api/v1
```

Change this to your EC2 public IP when deploying:

```env
VITE_API_URL=http://<ec2-public-ip>:8080/api/v1
```

**3. Start the frontend dev server**

```bash
npm run dev
```

Open http://localhost:5173 in your browser.

---

## User Roles

| Role | What They Can Do | What They Cannot Do |
|------|-----------------|-------------------|
| **Admin** | Full CRUD on Doctors, Receptionists, Pharmacists, Blood Managers; view all dashboards; book appointments; all read operations | Add/edit/delete medicine (pharmacist handles that) |
| **Doctor** | Prescribe medicines; create blood requests; view own appointments; view patient history | Book appointments; dispense medicine; record blood donations; manage patients or donors |
| **Receptionist** | Full CRUD on Patients; book regular & OT appointments; register ER patients; add ER entries; view all appointments | Manage donors; prescribe; dispense; blood operations |
| **Pharmacist** | View pending prescriptions; dispense medicines; add/restock medicine; check medicine stock | Manage patients or donors; book appointments; blood operations |
| **Blood Manager** | Full CRUD on Donors; record blood donations; fulfill blood requests; view inventory, history, expired blood | Prescribe medicine; book appointments; manage patients |

### Seed Credentials

| Role | Email | Password |
|------|-------|----------|
| Admin | admin@hospital.com | admin123 |
| Doctor | doctor@hospital.com | doctor123 |
| Receptionist | reception@hospital.com | recep123 |
| Pharmacist | pharmacist@hospital.com | pharma123 |
| Blood Manager | bloodmgr@hospital.com | blood123 |

---

## Feature Walkthrough

### Auth
| Feature | Role | Endpoint |
|---------|------|----------|
| Register | All | `POST /auth/register` |
| Login (returns JWT) | All | `POST /auth/login` |

### Patient Management
| Feature | Role | Endpoint |
|---------|------|----------|
| Add patient | Receptionist, Admin | `POST /patients/` |
| View all patients | Receptionist, Admin | `GET /patients/` |
| View patient by email | Receptionist, Admin | `GET /patients/:email` |
| Edit patient | Receptionist, Admin | `PUT /patients/:email` |
| Delete patient | Receptionist, Admin | `DELETE /patients/:email` |

### Appointments
| Feature | Role | Endpoint |
|---------|------|----------|
| Get available doctors for a slot | Receptionist, Admin | `POST /appointments/available-doctors` |
| Book regular appointment | Receptionist, Admin | `POST /appointments/book` |
| Book OT appointment | Receptionist, Admin | `POST /appointments/book-ot` |
| View all appointments | Receptionist, Admin | `GET /receptionist/appointments` |

### Prescriptions
| Feature | Role | Endpoint |
|---------|------|----------|
| Prescribe medicines | Doctor, Admin | `POST /prescriptions/prescribe` |
| Dispense medicines | Pharmacist, Admin | `POST /dispensing/dispense` |
| View pending prescriptions | Pharmacist, Admin | `GET /pharmacist/pending-prescriptions` |

### Medicine
| Feature | Role | Endpoint |
|---------|------|----------|
| View all medicine stock | Pharmacist, Admin | `GET /pharmacist/medicine-stock` |
| Search medicine by name | Pharmacist, Admin | `GET /pharmacist/medicine/:name` |
| Add / restock medicine | Pharmacist, Admin | `POST /pharmacist/medicine` |

### Blood Donation & Requests
| Feature | Role | Endpoint |
|---------|------|----------|
| Record blood donation | Blood Manager, Admin | `POST /blood/donate` |
| Create blood request  | Doctor, Admin | `POST /blood-request/create` |
| Fulfill blood request | Blood Manager, Admin | `POST /blood-fulfill/fulfill` |
| View donation history | Blood Manager, Admin | `GET /blood-manager/donations` |
| View blood inventory | Blood Manager, Admin | `GET /blood-manager/inventory` |
| View pending requests | Blood Manager, Admin | `GET /blood-manager/pending-requests` |
| View expired blood | Blood Manager, Admin | `GET /blood-manager/expired-blood` |

### Donors
| Feature | Role | Endpoint |
|---------|------|----------|
| Add donor | Blood Manager, Admin | `POST /donors/` |
| View all donors | Blood Manager, Admin | `GET /donors/` |
| View donor by email | Blood Manager, Admin | `GET /donors/:email` |
| Edit donor | Blood Manager, Admin | `PUT /donors/:email` |
| Delete donor | Blood Manager, Admin | `DELETE /donors/:email` |

### ER Management
| Feature | Role | Endpoint |
|---------|------|----------|
| View all ER patients | Receptionist, Admin | `GET /receptionist/er-patients` |
| Register ER patient | Receptionist, Admin | `POST /receptionist/er-patients` |
| View all ER entries | Receptionist, Admin | `GET /receptionist/er-entries` |
| Add ER entry | Receptionist, Admin | `POST /receptionist/er-entries` |

### Doctor Dashboard
| Feature | Role | Endpoint |
|---------|------|----------|
| View my appointments | Doctor, Admin | `GET /doctor/my-appointments` |
| View patient history | Doctor, Admin | `GET /doctor/patient-history/:email` |

### Admin : Staff Management
| Feature | Endpoint |
|---------|----------|
| CRUD Doctors | `GET/POST /admin/doctors`, `PUT/DELETE /admin/doctors/:email` |
| CRUD Receptionists | `GET/POST /admin/receptionists`, `PUT/DELETE /admin/receptionists/:email` |
| CRUD Pharmacists | `GET/POST /admin/pharmacists`, `PUT/DELETE /admin/pharmacists/:email` |
| CRUD Blood Managers | `GET/POST /admin/blood-managers`, `PUT/DELETE /admin/blood-managers/:email` |

---

## Transaction Scenarios

### 1. Prescribe Medicines:
**Trigger:** Doctor submits prescription with patient email, doctor email, and a list of medicine names and quantities.  
**Atomic operations:** Verify patient exists -> verify doctor exists -> verify appointment exists -> `BEGIN` -> `INSERT prescription` -> for each medicine: look up by name, check expiry, `INSERT prescription_medicine` -> `COMMIT`  
**Rollback causes:** Medicine not found; medicine expired; DB trigger detects duplicate same-day prescription for same doctor+patient; any INSERT fails.  
**Endpoint:** `POST /api/v1/prescriptions/prescribe`  
**Code:** `controllers/prescription.go -> PrescribeMedicines()`

### 2. Dispense Medicines
**Trigger:** Pharmacist submits prescription ID and pharmacist ID.  
**Atomic operations:** Fetch all medicines on prescription -> pre-check stock and expiry -> `BEGIN` -> for each medicine: `INSERT dispensing` (trigger `trg_dispensing_after_insert` auto-deducts stock and trigger `trg_dispensing_before` re-validates) -> `COMMIT`  
**Rollback causes:** Insufficient stock; expired medicine; trigger violations.  
**Endpoint:** `POST /api/v1/dispensing/dispense`  
**Code:** `controllers/prescription.go -> DispenseMedicines()`

### 3. Book Regular Appointment
**Trigger:** Receptionist submits patient email, doctor ID, date, and time.  
**Atomic operations:** Verify receptionist, patient, doctor -> `BEGIN` -> check doctor has no pending appointment at same date+time -> `INSERT appointment (ot_id=NULL, status='pending')` -> `COMMIT`  
**Rollback causes:** Doctor already has a pending appointment at that slot.  
**Endpoint:** `POST /api/v1/appointments/book`  
**Code:** `controllers/appointment.go -> BookAppointment()`

### 4. Book OT Appointment
**Trigger:** Receptionist submits patient email, doctor ID, date, and time (doctor has requested an OT surgery).  
**Atomic operations:** Verify all parties -> `BEGIN` -> check doctor free -> find OT where `is_available=TRUE` and not assigned at same date+time → `SELECT FOR UPDATE` on OT (race-condition lock) -> `INSERT appointment (ot_id=found OT, status='pending')` -> `UPDATE ot SET is_available=FALSE` -> `COMMIT`  
**Rollback causes:** Doctor not free; no OT available; concurrent transaction grabbed the OT first.  
**Endpoint:** `POST /api/v1/appointments/book-ot`  
**Code:** `controllers/appointment.go -> BookOTAppointment()`

### 5. Blood Donation
**Trigger:** Blood manager submits donor email and manager ID.  
**Atomic operations:** Verify manager and donor -> pre-check 90-day rule -> `BEGIN` -> `SELECT FOR UPDATE` on donor (prevents duplicate donation at two locations) -> re-check 90 days after lock -> `INSERT blood (unit=0, expiry=today+42 days)` -> `INSERT donation` (triggers `trg_donor_management_after` increment blood unit +1 and update `donor.last_donate=today`) -> `COMMIT`  
**Rollback causes:** Donor donated within last 90 days; concurrent transaction already processed this donor.  
**Endpoint:** `POST /api/v1/blood/donate`  
**Code:** `controllers/blood_donation.go -> RecordBloodDonation()`

### 6. Create Blood Request
**Trigger:** Doctor submits patient email, doctor ID, and quantity needed.  
**Atomic operations:** Verify doctor and patient (by email) -> `BEGIN` -> check patient has no existing pending request -> `INSERT blood_request (status='pending')` -> `COMMIT`  
**Rollback causes:** Patient already has a pending blood request.  
**Endpoint:** `POST /api/v1/blood-request/create`  
**Code:** `controllers/blood_request.go -> CreateBloodRequest()`

### 7. Fulfill Blood Request
**Trigger:** Blood manager submits request ID and manager ID.  
**Atomic operations:** Fetch request (patient blood group + quantity needed) -> `BEGIN` -> `SELECT FOR UPDATE` on blood_request -> re-check still 'pending' -> auto-select compatible non-expired blood ordered by oldest expiry first -> `SELECT FOR UPDATE` on blood row -> `INSERT blood_request_fulfillment` (triggers `trg_blood_fulfillment_before` validates; `trg_blood_fulfillment_after` deducts stock, updates blood status, marks request 'complete') -> `COMMIT`  
**Rollback causes:** Request already fulfilled; no compatible blood available; stock changed between pre-check and lock.  
**Endpoint:** `POST /api/v1/blood-fulfill/fulfill`  
**Code:** `controllers/blood_request.go -> FulfillBloodRequest()`

---

## ACID Compliance

| Property | Implementation |
|----------|---------------|
| **Atomicity** | Every transaction uses `BEGIN` / `COMMIT` / `defer tx.Rollback()` in Go. Any step failure rolls back all changes — no partial writes persist. |
| **Consistency** | PostgreSQL `NOT NULL`, `FOREIGN KEY`, `CHECK` constraints on all tables. 5 triggers enforce domain rules: duplicate prescription check, 90-day donation rule, stock deduction, expiry validation, blood compatibility. Medicine stock cannot go negative (trigger prevents it). |
| **Isolation** | `SELECT FOR UPDATE` row-level locks used in 4,5,7 and to prevent race conditions when two users target the same OT room, donor record, or blood unit simultaneously. |
| **Durability** | PostgreSQL WAL (Write-Ahead Logging) ensures committed transactions survive crashes. `pgxpool` connection pooling maintains persistent connections for reliability. |

---

## Indexing & Performance

The following indexes are defined in `database/performance.sql`:

| Index | Table | Column(s) | Reason |
|-------|-------|-----------|--------|
| `idx_appointment_doctor_date` | appointment | doctor_id, date, time | Speeds up doctor availability check on every appointment booking |
| `idx_appointment_patient` | appointment | patient_id | Speeds up patient history lookups |
| `idx_prescription_patient` | prescription | patient_id | Speeds up prescription history per patient |
| `idx_blood_group_status` | blood | b_gr, status, expiry_date | Speeds up compatible blood auto-selection in Transaction 4B |

Run `database/performance.sql` to apply indexes and view EXPLAIN ANALYZE results comparing before/after query costs.

---

## API Reference

Full specification is in `docs/swagger.yaml`. Quick reference:

| Method | Route | Auth | Purpose |
|--------|-------|------|---------|
| POST | `/auth/register` | No | Register any role |
| POST | `/auth/login` | No | Login, receive JWT |
| GET | `/patients/` | Receptionist, Admin | List all patients |
| POST | `/patients/` | Receptionist, Admin | Add patient |
| PUT | `/patients/:email` | Receptionist, Admin | Edit patient |
| DELETE | `/patients/:email` | Receptionist, Admin | Delete patient |
| POST | `/appointments/available-doctors` | Receptionist, Admin | Get free doctors for a slot |
| POST | `/appointments/book` | Receptionist, Admin | Book regular appointment |
| POST | `/appointments/book-ot` | Receptionist, Admin | Book OT appointment |
| POST | `/prescriptions/prescribe` | Doctor, Admin | Prescribe medicines (Tx 1) |
| POST | `/dispensing/dispense` | Pharmacist, Admin | Dispense medicines (Tx 2) |
| GET | `/pharmacist/pending-prescriptions` | Pharmacist, Admin | Pending prescriptions |
| GET | `/pharmacist/medicine-stock` | Pharmacist, Admin | All medicine with status |
| POST | `/pharmacist/medicine` | Pharmacist, Admin | Add or restock medicine |
| POST | `/blood/donate` | Blood Manager, Admin | Record blood donation (Tx 5) |
| POST | `/blood-request/create` | Doctor, Admin | Create blood request (Tx 6) |
| POST | `/blood-fulfill/fulfill` | Blood Manager, Admin | Fulfill blood request (Tx 7) |
| GET | `/blood-manager/inventory` | Blood Manager, Admin | Active blood inventory |
| GET | `/blood-manager/pending-requests` | Blood Manager, Admin | Pending requests |
| GET | `/doctor/my-appointments` | Doctor, Admin | Own schedule (from JWT) |
| GET | `/doctor/patient-history/:email` | Doctor, Admin | Full patient history |
| GET/POST | `/admin/doctors` | Admin | View / add doctors |
| PUT/DELETE | `/admin/doctors/:email` | Admin | Edit / delete doctor |

All protected routes require `Authorization: Bearer <token>` header. Full request/response schemas, status codes, and examples are in `docs/swagger.yaml`.

---

## Accessing the Deployed API (EC2)

The backend is deployed on AWS EC2.

**Base URL:**
```
http://<ec2-public-ip>:8080/api/v1
```

**Step 1 — Verify the server is live:**
```
GET http://<ec2-public-ip>:8080/health
→ { "status": "Server is running" }
```

**Step 2 — Register:**
```
POST http://<ec2-public-ip>:8080/api/v1/auth/register
{
  "name": "Your Name",
  "email": "you@example.com",
  "password": "yourpassword",
  "p_no": "03001234567",
  "role": "doctor"
}
```

**Step 3 — Login and copy token:**
```
POST http://<ec2-public-ip>:8080/api/v1/auth/login
{
  "email": "you@example.com",
  "password": "yourpassword",
  "role": "doctor"
}
```

**Step 4 — Attach token to all requests:**

In Postman -> Authorization tab -> Bearer Token -> paste token.

---

## Project Structure

```
Group_14_HospitalManagement/
├── backend/
│   ├── main.go
│   ├── .env
│   ├── .env.example
│   ├── go.mod / go.sum
│   ├── db/db.go                    
│   ├── middleware/auth.go           
│   ├── models/models.go            
│   ├── controllers/
│   │   ├── auth.go
│   │   ├── admin.go               
│   │   ├── patient.go              
│   │   ├── donor.go                
│   │   ├── appointment.go          
│   │   ├── prescription.go          
│   │   ├── blood_donation.go      
│   │   ├── blood_request.go        
│   │   ├── blood_manager.go        
│   │   ├── doctor.go               
│   │   ├── pharmacist.go            
│   │   └── Receptionist.go         
│   └── routes/routes.go             
├── frontend/
│   ├── src/
│   │   ├── api/                     
│   │   ├── components/             
│   │   ├── context/AuthContext.jsx  
│   │   ├── pages/                  
│   │   └── App.jsx                
│   ├── .env
│   └── .env.example
├── database/
│   ├── schema.sql                  
│   ├── seed.sql                    
│   └── performance.sql           
├── docs/
│   ├── swagger.yaml                
│   ├── ACID_Documentation.pdf
│   ├── Backend_Explanation.pdf
│   ├── Schema_Documentation.pdf
│   └── ER_Diagram.pdf
├── media/
│   └── RollBack_ScreenShots.pdf    
└── README.md
```

---

## Common Errors & Fixes

| Error | Fix |
|-------|-----|
| `Error loading .env file` | Make sure `.env` exists inside `backend/` folder |
| `Unable to connect to database` | PostgreSQL not running, or wrong password/DB name in `.env` |
| `relation "patient" does not exist` | Run `schema.sql` |
| `bind: address already in use` | Port 8080 taken — kill process or change `PORT` in `.env` |
| `go: module not found` | Run `go mod tidy` from inside `backend/` |
| `401 Unauthorized` | Token missing, expired, or wrong role -> re-login |
| `Role not found` | `AuthMiddleware` didn't run before `RequireRole` -> check routes.go grouping |
| `npm: not recognized` | Node.js not installed -> download from nodejs.org |
| Frontend blank page | Check `VITE_API_URL` in `frontend/.env` matches backend address |

---

## Known Issues & Limitations

- ER Shift setup (inserting doctors into `er_shift_doctor`) must be done directly via `psql` -> there is no API route for ER shift management.
- OT room records must be manually seeded in the `ot` table before OT appointment booking works.
- Medicine edit and delete are not exposed via any API route -> intentional design decision. Pharmacist manages stock via add/restock; admin cannot delete medicine records to preserve prescription history integrity.
- The frontend does not yet implement websocket real-time updates -> page refresh required to see new data from other users.

---

## Stopping the Server

```bash
# Backend
Ctrl + C

# Frontend
Ctrl + C
```
