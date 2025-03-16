package routes

import (
	auth "github.com/StackOverfloweds/MAUT-PhoneRank/controllers/Auth"
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

	// Protected Routes (Require JWT)
	api := app.Group("/api", middleware.JWTMiddleware())
	api.Put("/profile", user.UpdateProfile)
}
