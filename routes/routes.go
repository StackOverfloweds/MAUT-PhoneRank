package routes

import (
	auth "github.com/StackOverfloweds/MAUT-PhoneRank/controllers/Auth"
	user "github.com/StackOverfloweds/MAUT-PhoneRank/controllers/User"
	"github.com/StackOverfloweds/MAUT-PhoneRank/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	// Public Routes
	app.Post("/register", auth.Register)
	app.Post("/login", auth.Login)
	app.Post("/logout", auth.Logout)

	// Protected Routes (Require JWT)
	api := app.Group("/api", middleware.JWTMiddleware())
	api.Put("/profile", user.UpdateProfile)
}
