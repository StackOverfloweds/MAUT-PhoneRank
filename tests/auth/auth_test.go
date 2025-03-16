package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/StackOverfloweds/MAUT-PhoneRank/database"
	"github.com/StackOverfloweds/MAUT-PhoneRank/helpers"
	"github.com/StackOverfloweds/MAUT-PhoneRank/models"
	"github.com/StackOverfloweds/MAUT-PhoneRank/tests" // ✅ Import shared test helpers
	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	tests.SetupTestDB() // ✅ Ensure database is initialized

	// Create a test user before testing login
	database.DB.Create(&models.User{
		Username: "testuser",
		Phone:    "082253739918",
		Role:     "user",
	})

	app := tests.SetupApp()
	body := `{"phone_number": "082253739918"}`

	req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Clean up test data
	database.DB.Where("phone = ?", "082253739918").Delete(&models.User{})
}

/*
TestLogin_UnregisteredUser - Tests login attempt with an unregistered phone number.
*/
func TestLogin_UnregisteredUser(t *testing.T) {
	app := tests.SetupApp()
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
	tests.SetupTestDB() // ✅ Ensure database is initialized

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

	app := tests.SetupApp() // ✅ Use tests.SetupApp() instead of setupApp()
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
