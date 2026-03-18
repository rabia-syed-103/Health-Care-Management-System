package routes

import (
	"hospital-management/controllers"
	"hospital-management/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {

	api := r.Group("/api/v1")

	// AUTH ROUTES
	auth := api.Group("/auth")
	{
		auth.POST("/register", controllers.Register)
		auth.POST("/login", controllers.Login)
	}

	//  PROTECTED ROUTES

	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		//  PRESCRIPTION
		prescriptions := protected.Group("/prescriptions")
		prescriptions.Use(middleware.RequireRole("doctor", "admin"))
		{
			prescriptions.POST("/prescribe", controllers.PrescribeMedicines)
		}

		// DISPENSING
		dispensing := protected.Group("/dispensing")
		dispensing.Use(middleware.RequireRole("pharmacist", "admin"))
		{
			dispensing.POST("/dispense", controllers.DispenseMedicines)
		}

		//  APPOINTMENT
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
		// Blood Request
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
			patients.POST("/", controllers.AddPatient)
			patients.GET("/", controllers.GetAllPatients)
			patients.GET("/:email", controllers.GetPatient)
			patients.PUT("/:email", controllers.EditPatient)
			patients.DELETE("/:email", controllers.DeletePatient)
		}

		// DONOR
		donors := protected.Group("/donors")
		donors.Use(middleware.RequireRole("blood_manager", "admin"))
		{
			donors.POST("/", controllers.AddDonor)
			donors.GET("/", controllers.GetAllDonors)
			donors.GET("/:email", controllers.GetDonor)
			donors.PUT("/:email", controllers.EditDonor)
			donors.DELETE("/:email", controllers.DeleteDonor)
		}
		receptionist := protected.Group("/receptionist")
		receptionist.Use(middleware.RequireRole("receptionist", "admin"))
		{
			receptionist.GET("/appointments", controllers.GetAllAppointments)
			appointments.POST("/available-doctors", controllers.GetAvailableDoctors)
			receptionist.GET("/er-patients", controllers.GetAllERPatients)
			receptionist.POST("/er-patients", controllers.AddERPatient)

			receptionist.GET("/er-entries", controllers.GetAllERPatientEntries)
			receptionist.POST("/er-entries", controllers.AddERPatientEntry)
		}
		doctor := protected.Group("/doctor")
		doctor.Use(middleware.RequireRole("doctor", "admin"))
		{

			doctor.GET("/my-appointments", controllers.GetMyAppointments)

			doctor.GET("/patient-history/:email", controllers.GetPatientHistory)
		}
		pharmacist := protected.Group("/pharmacist")
		pharmacist.Use(middleware.RequireRole("pharmacist", "admin"))
		{
			pharmacist.GET("/pending-prescriptions", controllers.GetPendingPrescriptions)
			pharmacist.GET("/medicine-stock", controllers.GetMedicineStock)
			pharmacist.POST("/medicine", controllers.AddMedicine)
			pharmacist.GET("/medicine/:name", controllers.GetMedicineByName)
		}
		bloodManager := protected.Group("/blood-manager")
		bloodManager.Use(middleware.RequireRole("blood_manager", "admin"))
		{
			bloodManager.GET("/donations", controllers.GetDonationHistory)
			bloodManager.GET("/inventory", controllers.GetBloodInventory)
			bloodManager.GET("/pending-requests", controllers.GetPendingBloodRequests)
			bloodManager.GET("/expired-blood", controllers.GetExpiredBlood)
		}

		// ADMIN ONLY
		admin := protected.Group("/admin")
		admin.Use(middleware.RequireRole("admin"))
		{
			// Doctors
			admin.GET("/doctors", controllers.AdminGetAllDoctors)
			admin.POST("/doctors", controllers.AdminAddDoctor)
			admin.PUT("/doctors/:email", controllers.AdminEditDoctor)
			admin.DELETE("/doctors/:email", controllers.AdminDeleteDoctor)

			// Receptionists
			admin.GET("/receptionists", controllers.AdminGetAllReceptionists)
			admin.POST("/receptionists", controllers.AdminAddReceptionist)
			admin.PUT("/receptionists/:email", controllers.AdminEditReceptionist)
			admin.DELETE("/receptionists/:email", controllers.AdminDeleteReceptionist)

			// Pharmacists
			admin.GET("/pharmacists", controllers.AdminGetAllPharmacists)
			admin.POST("/pharmacists", controllers.AdminAddPharmacist)
			admin.PUT("/pharmacists/:email", controllers.AdminEditPharmacist)
			admin.DELETE("/pharmacists/:email", controllers.AdminDeletePharmacist)

			// Blood Managers
			admin.GET("/blood-managers", controllers.AdminGetAllBloodManagers)
			admin.POST("/blood-managers", controllers.AdminAddBloodManager)
			admin.PUT("/blood-managers/:email", controllers.AdminEditBloodManager)
			admin.DELETE("/blood-managers/:email", controllers.AdminDeleteBloodManager)

		}
	}

}
