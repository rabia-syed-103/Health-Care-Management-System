package routes

import (
	"hospital-backend/controllers"
	"hospital-backend/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {

	api := r.Group("/api")

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
	}

}
