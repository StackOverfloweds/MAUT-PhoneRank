package main

import (
	"log"

	"github.com/StackOverfloweds/MAUT-PhoneRank/database"
	"github.com/StackOverfloweds/MAUT-PhoneRank/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	database.ConnectDB()

	app := fiber.New()

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("db", database.DB)
		return c.Next()
	})

	routes.SetupRoutes(app)

	log.Fatal(app.Listen(":3000"))
}
