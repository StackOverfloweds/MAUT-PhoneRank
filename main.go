package main

import (
	"log"

	"github.com/StackOverfloweds/MAUT-PhoneRank/database"
	"github.com/StackOverfloweds/MAUT-PhoneRank/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No. Env file Found")
	}

	// setup db con
	database.ConnectDB()

	app := fiber.New()

	routes.SetupRoutes(app)

	// start server
	log.Fatal(app.Listen(":3000"))

}
