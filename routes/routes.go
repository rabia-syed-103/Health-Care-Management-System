package routes

import (
	"hospital-backend/controllers"
	"hospital-backend/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {

	api := r.Group("/api/v1")

	// ─── AUTH ROUTES (no auth needed) ────
	auth := api.Group("/auth")
	{
		auth.POST("/register", controllers.Register)
		auth.POST("/login", controllers.Login)
	}

	// ─── PROTECTED ROUTES ────────────────
	// AuthMiddleware runs first for everything below
	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		// ─── PRESCRIPTION ROUTES ─────────
		prescriptions := protected.Group("/prescriptions")
		prescriptions.Use(middleware.RequireRole("doctor", "admin"))
		{
			prescriptions.POST("/prescribe", controllers.PrescribeMedicines)
		}

		// ─── DISPENSING ROUTES ───────────
		dispensing := protected.Group("/dispensing")
		dispensing.Use(middleware.RequireRole("pharmacist", "admin"))
		{
			dispensing.POST("/dispense", controllers.DispenseMedicines)
		}

		// ─── APPOINTMENT ROUTES ──────────
		appointments := protected.Group("/appointments")
		appointments.Use(middleware.RequireRole("receptionist", "admin"))
		{
			appointments.POST("/book", controllers.BookAppointment)
			appointments.POST("/book-ot", controllers.BookOTAppointment)
		}
		blood := protected.Group("/blood")
		blood.Use(middleware.RequireRole("blood_manager", "admin"))
		{
			blood.POST("/donate", controllers.RecordBloodDonation)
		}
		// Blood Request — doctor creates, blood_manager fulfills
		bloodRequest := protected.Group("/blood-request")
		bloodRequest.Use(middleware.RequireRole("doctor", "admin"))
		{
			bloodRequest.POST("/create", controllers.CreateBloodRequest)
		}

		bloodFulfill := protected.Group("/blood-fulfill")
		bloodFulfill.Use(middleware.RequireRole("blood_manager", "admin"))
		{
			bloodFulfill.POST("/fulfill", controllers.FulfillBloodRequest)
		}
		patients := protected.Group("/patients")
		patients.Use(middleware.RequireRole("receptionist", "admin"))
		{
			patients.POST("/", controllers.AddPatient)            // Add
			patients.GET("/", controllers.GetAllPatients)         // View all
			patients.GET("/:email", controllers.GetPatient)       // View one by email
			patients.PUT("/:email", controllers.EditPatient)      // Edit by email
			patients.DELETE("/:email", controllers.DeletePatient) // Delete by email
		}

		// DONOR ROUTES — blood_manager + admin
		donors := protected.Group("/donors")
		donors.Use(middleware.RequireRole("blood_manager", "admin"))
		{
			donors.POST("/", controllers.AddDonor)            // Add
			donors.GET("/", controllers.GetAllDonors)         // View all
			donors.GET("/:email", controllers.GetDonor)       // View one by email
			donors.PUT("/:email", controllers.EditDonor)      // Edit by email
			donors.DELETE("/:email", controllers.DeleteDonor) // Delete by email
		}
		receptionist := protected.Group("/receptionist")
		receptionist.Use(middleware.RequireRole("receptionist", "admin"))
		{
			// Appointments
			receptionist.GET("/appointments", controllers.GetAllAppointments) // View all appointments

			// ER Patients
			receptionist.GET("/er-patients", controllers.GetAllERPatients) // View all ER patients
			receptionist.POST("/er-patients", controllers.AddERPatient)    // Register new ER patient

			// ER Patient Entries
			receptionist.GET("/er-entries", controllers.GetAllERPatientEntries) // View all ER entries
			receptionist.POST("/er-entries", controllers.AddERPatientEntry)     // Add ER patient entry
		}
		doctor := protected.Group("/doctor")
		doctor.Use(middleware.RequireRole("doctor", "admin"))
		{
			// View own appointments (pulled from JWT token — no ID needed in URL)
			doctor.GET("/my-appointments", controllers.GetMyAppointments)

			// View full history of a patient by email
			doctor.GET("/patient-history/:email", controllers.GetPatientHistory)
		}
		pharmacist := protected.Group("/pharmacist")
		pharmacist.Use(middleware.RequireRole("pharmacist", "admin"))
		{
			pharmacist.GET("/pending-prescriptions", controllers.GetPendingPrescriptions) // View pending
			pharmacist.GET("/medicine-stock", controllers.GetMedicineStock)               // Check stock
			pharmacist.POST("/medicine", controllers.AddMedicine)                         // Add or restock
		}
		bloodManager := protected.Group("/blood-manager")
		bloodManager.Use(middleware.RequireRole("blood_manager", "admin"))
		{
			bloodManager.GET("/donations", controllers.GetDonationHistory)             // All donation history
			bloodManager.GET("/inventory", controllers.GetBloodInventory)              // Active blood inventory
			bloodManager.GET("/pending-requests", controllers.GetPendingBloodRequests) // Pending blood requests
			bloodManager.GET("/expired-blood", controllers.GetExpiredBlood)            // Expired blood units
		}
	}

}
