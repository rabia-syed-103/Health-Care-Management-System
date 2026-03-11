# Hospital Management System 

# Setup & Run Guide

## Prerequisites

Make sure these are installed before starting:

| Tool | Version | Download |
|---|---|---|
| Go | 1.21+ | https://go.dev/dl |
| PostgreSQL | 15+ | https://www.postgresql.org/download |
| Git | Any | https://git-scm.com |
| Postman | Any | https://www.postman.com/downloads |


## Step 1: Clone the Repository

```bash
git clone https://github.com/rabia-syed-103/Health-Care-Management-System.git
cd Health-Care-Management-System
```


## Step 2: Create the Database

Open your terminal and run:

```bash
psql -U postgres
```

Then inside psql:

```sql
CREATE DATABASE hospital_db;
\q
```



## Step 3: Run the Schema

This creates all 21 tables, 5 triggers, 4 views, and 4 indexes:

```bash
psql -U postgres -d hospital_db -f schema.sql
```

You should see all commands running without error



## Step 4: Run the Seed Data (Optional)

```bash
psql -U postgres -d hospital_db -f seed.sql
```

## Step 5: Change the .env File

In the root of the project, create a file named `.env`:

```bash
# Windows
copy .env.example .env

# Mac / Linux
cp .env.example .env
```

Then open `.env` and fill in your values:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_postgres_password
DB_NAME=hospital_db
JWT_SECRET=anyrandomsecretkey123
PORT=8080
```

> **Note:** `JWT_SECRET` can be any string.

## Step 6: Install Go Dependencies

```bash
go mod tidy
```

It will downloads all required packages:
- `github.com/gin-gonic/gin`
- `github.com/jackc/pgx/v5`
- `github.com/golang-jwt/jwt/v5`
- `github.com/joho/godotenv`
- `golang.org/x/crypto`

> **Note:** `go mod tidy` might take a couple of minutes depending on your internet speed.


## Step 7: Run the Server

```bash
go run main.go
```

You should see:

```
Connected to PostgreSQL successfully!
[GIN-debug] Listening and serving HTTP on :8080
```

## Step 8: Verify Server is Running

Open Postman or your browser and paste:

```
GET http://localhost:8080/health
```

Expected response:

```json
{
  "status": "Server is running"
}
```
### Step 9: Register and Login

No credentials are needed to register. Just paste the register endpoint with your role and you will get a JWT token back.

**Register:**
```
POST http://localhost:8080/api/v1/auth/register
```
```json
{
  "name": "Your Name",
  "email": "you@example.com",
  "password": "yourpassword",
  "p_no": "03001234567",
  "role": "doctor"
}
```

**Login:**
```
POST http://localhost:8080/api/v1/auth/login
```
```json
{
  "email": "you@example.com",
  "password": "yourpassword",
  "role": "doctor"
}
```
Copy the `token` from the response.


### Step 10: Attach Token to All Requests

In Postman, for every request after login:

1. Go to the **Authorization** tab
2. Select **Bearer Token**
3. Paste your token

Every protected endpoint requires this. You will get `401 Unauthorized` without it.


### Step 11: Send Requests

Use the base URL format:

```
http://localhost:8080/api/v1/<endpoint>
```

**Example:**
```
GET http://localhost:8080/api/v1/patients/
```

Refer to `swagger.yaml` or the `Backend_Explanation.pdf` for the full list of endpoints and request bodies.

## Project Structure

```
hospital-backend/
├── main.go
├── .env
├── .env.example
├── schema.sql
├── seed.sql
├── swagger.yaml
├── README.md
├── db/
│   └── db.go
├── middleware/
│   └── auth.go
├── models/
│   └── models.go
├── controllers/
│   ├── auth.go
│   ├── admin.go
│   ├── patient.go
│   ├── donor.go
│   ├── receptionist.go
│   ├── doctor.go
│   ├── pharmacist.go
│   ├── blood_manager.go
│   ├── blood_donation.go
│   ├── blood_request.go
│   ├── appointment.go
│   └── prescription.go
└── routes/
    └── routes.go
```

## Accessing the Deployed API

The backend is already deployed and running on AWS EC2. You do not need to install anything or set up a database. Just open Postman and start making requests.


### Base URL

All requests go to:

```
http://<ec2-public-ip>:8080/api/v1
```

> The instructor/team member will provide the exact public IP.


### Step 1: Verify the Server is Live

Open Postman and send:

```
GET http://<ec2-public-ip>:8080/health
```

```json
{
  "status": "Server is running"
}
```

If this works, the server is up and you can test all endpoints.



### Step 2: Register and Login

No credentials are needed to register. Just paste the register endpoint with your role and you will get a JWT token back.

**Register:**
```
POST http://<ec2-public-ip>:8080/api/v1/auth/register
```
```json
{
  "name": "Your Name",
  "email": "you@example.com",
  "password": "yourpassword",
  "p_no": "03001234567",
  "role": "doctor"
}
```

**Login:**
```
POST http://<ec2-public-ip>:8080/api/v1/auth/login
```
```json
{
  "email": "you@example.com",
  "password": "yourpassword",
  "role": "doctor"
}
```
Copy the `token` from the response.


### Step 3: Attach Token to All Requests

In Postman, for every request after login:

1. Go to the **Authorization** tab
2. Select **Bearer Token**
3. Paste your token

Every protected endpoint requires this. You will get `401 Unauthorized` without it.


### Step 4: Send Requests

Use the base URL format:

```
http://<ec2-public-ip>:8080/api/v1/<endpoint>
```

**Example:**
```
GET http://<ec2-public-ip>:8080/api/v1/patients/
```

Refer to `swagger.yaml` or the `Backend_Explanation.pdf` for the full list of endpoints and request bodies.

## Common Errors & Fixes

### `Error loading .env file`
The `.env` file is missing or not in the root folder. Make sure it exists at `hospital-backend/.env`.

### `Unable to connect to database`
- PostgreSQL is not running — start it from Services (Windows) or `brew services start postgresql` (Mac)
- Wrong password in `.env`
- Wrong DB name — make sure `hospital_db` was created in Step 2

### `Database ping failed`
PostgreSQL is running but rejecting the connection. Check `DB_USER` and `DB_PASSWORD` in `.env`.

### `relation "patient" does not exist`
Schema was not run. Go back to Step 3 and run `schema.sql`.

### `bind: address already in use`
Port 8080 is taken. Either kill the other process or change `PORT=8081` in `.env`.

### `go: module not found`
Run `go mod tidy` again from inside the `hospital-backend/` folder.



## Stopping the Server

Press `Ctrl + C` in the terminal where the server is running.
