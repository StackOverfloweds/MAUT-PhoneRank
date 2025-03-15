package routes

import (
	"github.com/StackOverfloweds/MAUT-PhoneRank/controllers"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Get("/smartphone", controllers.GetSmartphone)
	api.Post("/smarthone", controllers.CreateSmartphone)
}
