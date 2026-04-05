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
| Backend | Go 1.25+, Gin (HTTP framework) |
| Database | PostgreSQL 16 |
| Auth | JWT (golang-jwt/jwt v5), bcrypt password hashing |
| ORM / DB Driver | pgx/v5 with pgxpool (raw SQL) |
| Config | godotenv (.env loading) |
| API Spec | OpenAPI 3.0 (swagger.yaml) |
| Containerization | Docker, Docker Compose |
| Web Server | Nginx (frontend serving + API proxy) |
| Deployment | AWS EC2 (Ubuntu 24.04) |
| PDF Export | jsPDF + jspdf-autotable |
| Email | Gmail SMTP via gomail.v2 |

---

## System Architecture

```
Browser (React + Vite)
        │
        │  HTTP/JSON (axios)
        |
    Nginx (port 80)
        │
        ├── /* → serves React build (dist/)
        └── /api/* → proxies to Go backend
                │
                |
        Go Backend (Gin) :8080
                │
                ├── middleware/auth.go      JWT verification
                ├── middleware/auth.go      RBAC
                ├── controllers/            Raw SQL transactions
                │
                │  pgxpool
                |
        PostgreSQL 16 :5432
                │
                ├── 21 tables
                ├── 13 triggers
                ├── 4 views
                └── 4 indexes
```

**Request Lifecycle:**
1. React sends HTTP request with `Authorization: <token>` header
2. Nginx receives request — static files served directly, `/api/*` proxied to Go backend
3. `AuthMiddleware` parses and validates the JWT; sets `user_id` and `role` in Gin context
4. `RequireRole` checks the role against the allowed list for that route
5. Controller runs pre-checks, then opens a `BEGIN` transaction
6. SQL operations execute and any failure triggers `defer tx.Rollback()`
7. On success, `tx.Commit()` persists all changes atomically
8. JSON response returned to React

---

## Features

- **Docker Containerization** — Full stack (frontend + backend + database) orchestrated via Docker Compose; deployable on any machine with a single command
- **PDF Export** — Multiple data tables (donors, patients, prescriptions, inventory) and reports can be downloaded as a formatted PDF report using jsPDF such as Patient History,Prescrption etc.
- **Real-Time Email Notifications** — Automated emails sent via Gmail SMTP on key events using gomail.v2 on a specifiic event. 
- **Analytical Dashboard** — Visual charts and stats giving role-based overviews of hospital activity using Recharts

## UI Examples

### 1. Login Page
Role-based login with JWT token storage. Each role redirects to its own dashboard after authentication. Required because every protected route checks the JWT role claim. Wrong role gets 401/403.

### 2. Receptionist : Book Appointment
Step 1: Receptionist selects date and time, system queries all doctors with no conflicting pending appointment. Step 2: Receptionist picks a doctor from the available list and enters the patient email. Required to demonstrate the available doctors transaction and prevent double booking.

### 3. Blood Manager : Pending Requests + Fulfill
Lists all pending blood requests with patient name, blood group, and quantity needed. One click opens the fulfillment modal. The system automatically selects the oldest compatible blood unit, locks the row with `SELECT FOR UPDATE`, and commits atomically. Required to show Transaction end-to-end and fulfill the blood-management component of the project.

---

## Running the Project

There are two ways to run this project: **Docker (recommended)** or **Manual Setup**.

---

## Option 1: Docker Setup (Recommended)

Docker runs the entire stack (frontend + backend + database) with a single command. No need to install Go, Node.js, or PostgreSQL separately.

### Prerequisites

| Tool | Version |
|------|---------|
| Docker Desktop | Latest |
| Git | Any |

### Steps

**1. Clone the repository**

```bash
git clone https://github.com/rabia-syed-103/Health-Care-Management-System.git
cd Health-Care-Management-System
```

**2. Create root `.env`**

```bash
# Windows
copy .env.example .env

# Mac / Linux
cp .env.example .env
```

Open `.env` and fill in your values:

```env
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=healthcare
JWT_SECRET=anyrandomsecretkey
EMAIL_USER=youremail@gmail.com
EMAIL_PASS=your_16char_app_password
```

**3. Create backend `.env`**

```bash
# Windows
copy backend\.env.example backend\.env

# Mac / Linux
cp backend/.env.example backend/.env
```

Open `backend/.env` and fill in:

```env
DB_HOST=db              ::must be "db" not "localhost" in Docker
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=healthcare
JWT_SECRET=anyrandomsecretkey
EMAIL_USER=youremail@gmail.com
EMAIL_PASS=your_16char_app_password
PORT=8080
```

> ⚠️ `DB_HOST` must be `db` (the Docker service name), not `localhost`

**4. Build and start everything**

```bash
docker-compose up --build
```

This will:
- Pull PostgreSQL, Nginx, Node, Go images
- Build frontend (npm install + vite build)
- Build backend (Go compile)
- Run schema.sql and seed.sql automatically
- Start all 3 containers

> First build takes 10–20 minutes. Subsequent builds take ~30 seconds.

**5. Open in browser**

```
http://localhost
```

**6. Stop everything**

```bash
docker-compose down        # stop containers
docker-compose down -v     # stop + delete database data
```

### Pre-built Docker Images

To skip building locally, pre-built images are available on Docker Hub:

```bash
# Update docker-compose.yml to use images instead of build:
# backend:  image: rabiasyed/hms-backend
# frontend: image: rabiasyed/hms-frontend

docker-compose pull
docker-compose up -d
```

---

## Option 2 — Manual Setup (Without Docker)

Use this if you prefer running services directly on your machine.

### Prerequisites

| Tool | Version |
|------|---------|
| Go | 1.25+ |
| PostgreSQL | 16+ |
| Node.js | 20+ |
| Git | Any |

### Database Setup

```bash
# Open psql
psql -U postgres

# Inside psql
CREATE DATABASE healthcare;
\q

# Run schema
psql -U postgres -d healthcare -f database/schema.sql

# Run seed data
psql -U postgres -d healthcare -f database/seed.sql
```

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
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=healthcare
JWT_SECRET=anyrandomkey
EMAIL_USER=youremail@gmail.com
EMAIL_PASS=your_16char_app_password
PORT=8080
```

**2. Install Go dependencies**

```bash
cd backend
go mod tidy
```

This installs:
- `github.com/gin-gonic/gin` — HTTP router
- `github.com/jackc/pgx/v5` — PostgreSQL driver with connection pooling
- `github.com/golang-jwt/jwt/v5` — JWT auth
- `github.com/joho/godotenv` — .env loading
- `golang.org/x/crypto` — bcrypt password hashing
- `gopkg.in/gomail.v2` — email sending

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
VITE_API_URL=/api/v1
```

> The `vite.config.js` proxy forwards `/api` calls to `localhost:8080` automatically during development.

**3. Start the frontend dev server**

```bash
npm run dev
```

Open `http://localhost:5173` in your browser.

---

## User Roles

| Role | What They Can Do | What They Cannot Do |
|------|-----------------|-------------------|
| **Admin** | Full CRUD on Doctors, Receptionists, Pharmacists, Blood Managers; view staff activity; PDF exports; all read operations |  |
| **Doctor** | Prescribe medicines; create blood requests; view own appointments; view patient history; download PDF reports | Book appointments; dispense medicine; record blood donations; manage patients or donors |
| **Receptionist** | Full CRUD on Patients; book regular & OT appointments; register ER patients; add ER entries; view all appointments; PDF exports | Manage donors; prescribe; dispense; blood operations |
| **Pharmacist** | View pending prescriptions; dispense medicines; add/restock medicine; check medicine stock; PDF exports | Manage patients or donors; book appointments; blood operations |
| **Blood Manager** | Full CRUD on Donors; record blood donations; fulfill blood requests; view inventory, history, expired blood; PDF exports | Prescribe medicine; book appointments; manage patients |

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
| Cancel appointment | Doctor, Receptionist, Admin | `PUT /appointments/:id/cancel` |
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
| Create blood request | Doctor, Admin | `POST /blood-request/create` |
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

### Admin — Staff Management
| Feature | Endpoint |
|---------|----------|
| CRUD Doctors | `GET/POST /admin/doctors`, `PUT/DELETE /admin/doctors/:email` |
| View Doctor activity | `GET /admin/doctors/:email/detail` |
| CRUD Receptionists | `GET/POST /admin/receptionists`, `PUT/DELETE /admin/receptionists/:email` |
| View Receptionist activity | `GET /admin/receptionists/:email/detail` |
| CRUD Pharmacists | `GET/POST /admin/pharmacists`, `PUT/DELETE /admin/pharmacists/:email` |
| View Pharmacist activity | `GET /admin/pharmacists/:email/detail` |
| CRUD Blood Managers | `GET/POST /admin/blood-managers`, `PUT/DELETE /admin/blood-managers/:email` |
| View Blood Manager activity | `GET /admin/blood-managers/:email/detail` |

---

## Transaction Scenarios

### 1. Prescribe Medicines
**Trigger:** Doctor submits prescription with patient email, doctor email, and a list of medicine names and quantities.\
**Atomic operations:** Verify patient exists -> verify doctor exists -> verify appointment exists -> `BEGIN` -> `INSERT prescription` -> for each medicine: look up by name, check expiry, `INSERT prescription_medicine` -> `COMMIT`\
**Rollback causes:** Medicine not found; medicine expired; DB trigger detects duplicate same-day prescription for same doctor+patient; any INSERT fails.\
**Endpoint:** `POST /api/v1/prescriptions/prescribe`\
**Code:** `controllers/prescription.go -> PrescribeMedicines()`

### 2. Dispense Medicines
**Trigger:** Pharmacist submits prescription ID and pharmacist ID.\
**Atomic operations:** Fetch all medicines on prescription -> pre-check stock and expiry -> `BEGIN` -> for each medicine: `INSERT dispensing` (trigger `trg_dispensing_after_insert` auto-deducts stock and trigger `trg_dispensing_before` re-validates) -> `COMMIT`\
**Rollback causes:** Insufficient stock; expired medicine; trigger violations.\
**Endpoint:** `POST /api/v1/dispensing/dispense`\
**Code:** `controllers/prescription.go -> DispenseMedicines()`

### 3. Book Regular Appointment
**Trigger:** Receptionist submits patient email, doctor ID, date, and time.\
**Atomic operations:** Verify receptionist, patient, doctor -> `BEGIN` -> check doctor has no pending appointment at same date+time -> `INSERT appointment (ot_id=NULL, status='pending')` -> `COMMIT`\
**Rollback causes:** Doctor already has a pending appointment at that slot.\
**Endpoint:** `POST /api/v1/appointments/book`\
**Code:** `controllers/appointment.go -> BookAppointment()`

### 4. Book OT Appointment
**Trigger:** Receptionist submits patient email, doctor ID, date, and time (doctor has requested an OT surgery).\
**Atomic operations:** Verify all parties -> `BEGIN` -> check doctor free -> find OT where `is_available=TRUE` and not assigned at same date+time -> `SELECT FOR UPDATE` on OT (race-condition lock) -> `INSERT appointment (ot_id=found OT, status='pending')` -> `UPDATE ot SET is_available=FALSE` -> `COMMIT`\
**Rollback causes:** Doctor not free; no OT available; concurrent transaction grabbed the OT first.\
**Endpoint:** `POST /api/v1/appointments/book-ot`\
**Code:** `controllers/appointment.go -> BookOTAppointment()`

### 5. Blood Donation
**Trigger:** Blood manager submits donor email and manager ID.\
**Atomic operations:** Verify manager and donor -> pre-check 90-day rule -> `BEGIN` -> `SELECT FOR UPDATE` on donor (prevents duplicate donation at two locations) -> re-check 90 days after lock -> `INSERT blood (unit=0, expiry=today+42 days)` -> `INSERT donation` (triggers `trg_donor_management_after` increment blood unit +1 and update `donor.last_donate=today`) -> `COMMIT`\
**Rollback causes:** Donor donated within last 90 days; concurrent transaction already processed this donor.\
**Endpoint:** `POST /api/v1/blood/donate`\
**Code:** `controllers/blood_donation.go -> RecordBloodDonation()`

### 6. Create Blood Request
**Trigger:** Doctor submits patient email, doctor ID, and quantity needed.\
**Atomic operations:** Verify doctor and patient (by email) -> `BEGIN` -> check patient has no existing pending request -> `INSERT blood_request (status='pending')` -> `COMMIT`\
**Rollback causes:** Patient already has a pending blood request.\
**Endpoint:** `POST /api/v1/blood-request/create`\
**Code:** `controllers/blood_request.go -> CreateBloodRequest()`

### 7. Fulfill Blood Request
**Trigger:** Blood manager submits request ID and manager ID.\
**Atomic operations:** Fetch request (patient blood group + quantity needed) -> `BEGIN` -> `SELECT FOR UPDATE` on blood_request -> re-check still 'pending' -> auto-select compatible non-expired blood ordered by oldest expiry first -> `SELECT FOR UPDATE` on blood row -> `INSERT blood_request_fulfillment` (triggers `trg_blood_fulfillment_before` validates; `trg_blood_fulfillment_after` deducts stock, updates blood status, marks request 'complete') -> `COMMIT`\
**Rollback causes:** Request already fulfilled; no compatible blood available; stock changed between pre-check and lock.\
**Endpoint:** `POST /api/v1/blood-fulfill/fulfill`\
**Code:** `controllers/blood_request.go -> FulfillBloodRequest()`

---

## ACID Compliance

| Property | Implementation |
|----------|---------------|
| **Atomicity** | Every transaction uses `BEGIN` / `COMMIT` / `defer tx.Rollback()` in Go. Any step failure rolls back all changes — no partial writes persist. |
| **Consistency** | PostgreSQL `NOT NULL`, `FOREIGN KEY`, `CHECK` constraints on all tables. 13 triggers enforce domain rules: duplicate prescription check, 90-day donation rule, stock deduction, expiry validation, blood compatibility. Medicine stock cannot go negative (trigger prevents it). |
| **Isolation** | `SELECT FOR UPDATE` row-level locks used in transactions 4, 5, 7 to prevent race conditions when two users target the same OT room, donor record, or blood unit simultaneously. |
| **Durability** | PostgreSQL WAL (Write-Ahead Logging) ensures committed transactions survive crashes. `pgxpool` connection pooling maintains persistent connections for reliability. |

---

## Indexing & Performance

The following indexes are defined in `database/performance.sql`:

| Index | Table | Column(s) | Reason |
|-------|-------|-----------|--------|
| `idx_appointment_doctor_date` | appointment | doctor_id, date, time | Speeds up doctor availability check on every appointment booking |
| `idx_appointment_patient` | appointment | patient_id | Speeds up patient history lookups |
| `idx_prescription_patient` | prescription | patient_id | Speeds up prescription history per patient |
| `idx_blood_group_status` | blood | b_gr, status, expiry_date | Speeds up compatible blood auto-selection in Transaction 7 |

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
| PUT | `/appointments/:id/cancel` | Doctor, Receptionist, Admin | Cancel appointment |
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
| GET | `/admin/doctors/:email/detail` | Admin | View doctor profile + appointments |
| GET | `/admin/receptionists/:email/detail` | Admin | View receptionist profile + activity |
| GET | `/admin/pharmacists/:email/detail` | Admin | View pharmacist profile + dispensing |
| GET | `/admin/blood-managers/:email/detail` | Admin | View blood manager profile + donations |

All protected routes require `Authorization: Bearer <token>` header. Full request/response schemas, status codes, and examples are in `docs/swagger.yaml`.

---

## AWS Deployment

The system is deployed on AWS EC2 (Ubuntu 24.04) using Docker.

> ⚠️ **Important:** AWS assigns a new public IP every time the instance is stopped and started. So if you directly want to run. Use ```http://<Latest EC2 PUBLIC IP>```

**Base URL:**
```
http://<ec2-public-ip>/api/v1
```

**Verify the server is live:**
```
GET http://<ec2-public-ip>/health
→ { "status": "Server is running" }
```

### Deploy on a fresh AWS EC2 instance

```bash
# 1. Install Docker
sudo apt update
sudo apt install -y docker.io docker-compose
sudo systemctl start docker
sudo systemctl enable docker
sudo usermod -aG docker ubuntu
# logout and login again

# 2. Clone repo
git clone https://github.com/rabia-syed-103/Health-Care-Management-System.git
cd Health-Care-Management-System

# 3. Create root .env
nano .env
# fill in your values

# 4. Create backend .env (DB_HOST must be "db")
nano backend/.env

# 5. Run
docker-compose up -d
```

### Update deployment after code changes

```bash
# On your PC: build and push new images
docker build -t username/hms-backend ./backend
docker build -t username/hms-frontend ./frontend
docker push username/hms-backend
docker push username/hms-frontend

# On AWS
git pull
docker-compose pull
docker-compose up -d
```

### AWS Security Group: required open ports

| Port | Purpose |
|------|---------|
| 22 | SSH |
| 80 | Frontend (React via Nginx) |
| 8080 | Backend API (optional: only if accessing directly) |

---

## Project Structure

```
Group_14_HospitalManagement/
├── .env                             ::root env for docker-compose
├── docker-compose.yml               ::orchestrates all 3 services
├── backend/
│   ├── Dockerfile
│   ├── .dockerignore
│   ├── main.go
│   ├── .env
│   ├── .env.example
│   ├── go.mod 
|   ├── go.sum 
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
│   ├── Dockerfile
│   ├── .dockerignore
│   ├── nginx.conf
│   ├── src/
│   │   ├── api/
│   │   ├── components/
│   │   ├── context/AuthContext.jsx
│   │   ├── utils/pdfExport.js      ::PDF export utility
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
| `Error loading .env file` | Make sure `.env` exists inside `backend/` folder. In Docker this is handled automatically via docker-compose environment. |
| `Unable to connect to database` | PostgreSQL not running, or wrong password/DB name in `.env`. In Docker make sure `DB_HOST=db` not `localhost`. |
| `relation "patient" does not exist` | Run `schema.sql` |
| `bind: address already in use` | Port 8080 taken — kill process or change `PORT` in `.env` |
| `go: module not found` | Run `go mod tidy` from inside `backend/` |
| `401 Unauthorized` | Token missing, expired, or wrong role -> re-login |
| `Role not found` | `AuthMiddleware` didn't run before `RequireRole` -> check routes.go grouping |
| `npm: not recognized` | Node.js not installed -> download from nodejs.org |
| Frontend blank page | Check `VITE_API_URL` in `frontend/.env` — use `/api/v1` for Docker, `http://localhost:8080/api/v1` for manual setup |
| `no space left on device` (Docker) | Run `docker system prune -af` to free disk space |
| `TLS handshake timeout` (Docker pull) | Network issue — retry or switch to mobile hotspot |
| `ContainerConfig` error (docker-compose) | Run `docker rm -f hms_backend hms_frontend hms_db` then `docker-compose up -d` |
| AWS IP changed after restart | Stop/start assigns a new IP — use Elastic IP for a fixed address |

---

## Known Issues & Limitations

- ER Shift setup (inserting doctors into `er_shift_doctor`) must be done directly via `psql` — there is no API route for ER shift management.
- OT room records must be manually seeded in the `ot` table before OT appointment booking works.
- Medicine edit and delete are not exposed via any API route, intentional design decision. Pharmacist manages stock via add/restock; admin cannot delete medicine records to preserve prescription history integrity.
- When switching between Docker and manual dev setup, `VITE_API_URL` in `frontend/.env` must be updated accordingly (`/api/v1` for Docker, `http://localhost:8080/api/v1` for manual).
- AWS public IP changes on every stop/start of the EC2 instance unless an Elastic IP is configured.

---

## Stopping the Server

**Docker:**

To stop Containers:
```bash
docker-compose down  
```
To stop and wipe Database:
```
docker-compose down -v    
```

**Manual:**
```bash
# Backend
Ctrl + C

# Frontend
Ctrl + C
```