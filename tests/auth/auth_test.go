package auth

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	auth "github.com/StackOverfloweds/MAUT-PhoneRank/controllers/Auth"
	"github.com/StackOverfloweds/MAUT-PhoneRank/database"
	"github.com/StackOverfloweds/MAUT-PhoneRank/helpers"
	"github.com/StackOverfloweds/MAUT-PhoneRank/models"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Setup test database (using SQLite for in-memory testing)
func setupTestDB() {
	var err error
	database.DB, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{}) // Use in-memory DB for tests
	if err != nil {
		log.Fatalf("Failed to initialize test database: %v", err)
	}

	// Auto-migrate models
	database.DB.AutoMigrate(&models.User{}, &models.Profile{})
}

// Setup Fiber app for testing
func setupApp() *fiber.App {
	app := fiber.New()
	app.Post("/auth/login", auth.Login)
	app.Post("/auth/verify-otp", auth.VerifyOTP)
	return app
}

/*
TestLogin - Tests the Login function.
It checks whether an OTP is successfully sent for a registered user.
*/
func TestLogin(t *testing.T) {
	setupTestDB() // Ensure database is initialized

	// Create a test user
	database.DB.Create(&models.User{
		Username: "testuser",
		Phone:    "082253739918",
		Role:     "user",
	})

	app := setupApp()
	body := `{"phone_number": "082253739918"}`

	req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Clean up mock data
	database.DB.Where("phone = ?", "082253739918").Delete(&models.User{})
}

/*
TestLogin_UnregisteredUser - Tests login attempt with an unregistered phone number.
*/
func TestLogin_UnregisteredUser(t *testing.T) {
	app := setupApp()
	body := `{"phone_number": "000000000000"}`

	req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

/*
TestVerifyOTP - Tests OTP verification process.
This test simulates sending an OTP, verifying it, and checking for a valid token response.
*/
func TestVerifyOTP(t *testing.T) {
	setupTestDB() // Ensure database is initialized

	// Create a test user
	user := models.User{
		Username: "testuser",
		Phone:    "082253739918",
		Role:     "user",
	}
	database.DB.Create(&user)

	// Generate and store OTP manually
	otp := helpers.GenerateOTP()
	helpers.SendOTP(user.Phone, otp)

	app := setupApp()
	body := `{"otp": "` + otp + `"}`

	req := httptest.NewRequest(http.MethodPost, "/auth/verify-otp", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Parse response body
	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)

	// Check if token exists
	assert.NotEmpty(t, response["token"])
	assert.Equal(t, "user", response["role"])

	// Clean up mock data
	database.DB.Where("phone = ?", "082253739918").Delete(&models.User{})
}
