package tests

import (
	"log"

	auth "github.com/StackOverfloweds/MAUT-PhoneRank/controllers/Auth"
	"github.com/StackOverfloweds/MAUT-PhoneRank/database"
	"github.com/StackOverfloweds/MAUT-PhoneRank/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// SetupTestDB initializes a test database using SQLite in-memory.
func SetupTestDB() {
	var err error
	database.DB, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{}) // In-memory DB for testing
	if err != nil {
		log.Fatalf("Failed to initialize test database: %v", err)
	}

	// Auto-migrate necessary models
	database.DB.AutoMigrate(&models.User{}, &models.Profile{})
}

// SetupApp initializes the Fiber app with authentication routes for testing.
func SetupApp() *fiber.App {
	app := fiber.New()
	app.Post("/auth/login", auth.Login)
	app.Post("/auth/verify-otp", auth.VerifyOTP)
	app.Post("/auth/register", auth.Register)
	app.Post("/auth/logout", auth.Logout)
	return app
}
