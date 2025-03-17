package routes

import (
	auth "github.com/StackOverfloweds/MAUT-PhoneRank/controllers/Auth"
	Smartphone "github.com/StackOverfloweds/MAUT-PhoneRank/controllers/Smartphone"
	user "github.com/StackOverfloweds/MAUT-PhoneRank/controllers/User"
	"github.com/StackOverfloweds/MAUT-PhoneRank/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	// routes for authentication
	authRoutes := app.Group("/auth")
	authRoutes.Post("/register", auth.Register)
	authRoutes.Post("/login", auth.Login)
	authRoutes.Post("/verify-otp", auth.VerifyOTP)
	authRoutes.Post("/logout", auth.Logout)

	//make routes for profile
	userProf := app.Group("/user", middleware.JWTMiddleware())
	userProf.Put("/profile", user.UpdateProfile)

	// routes for smartphone
	smartphoneRoutes := app.Group("/smartphone")
	smartphoneRoutes.Get("/:id", Smartphone.GetSmartphoneDetail)
	smartphoneRoutes.Post("/search", Smartphone.SearchSmartphone)
	smartphoneRoutes.Post("/maut", Smartphone.CalculateMAUT)
}
