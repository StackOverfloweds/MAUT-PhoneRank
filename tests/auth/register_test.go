package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/StackOverfloweds/MAUT-PhoneRank/database"
	"github.com/StackOverfloweds/MAUT-PhoneRank/models"
	"github.com/StackOverfloweds/MAUT-PhoneRank/tests" // ✅ Import shared test helpers
	"github.com/stretchr/testify/assert"
)

func TestRegister_Success(t *testing.T) {
	tests.SetupTestDB() // ✅ Call shared test helper

	app := tests.SetupApp()
	body := `{
		"username": "testuser",
		"full_name": "Test User",
		"phone": "082253739918"
	}`

	req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Parse response body
	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)

	// Ensure user was created
	var user models.User
	database.DB.Where("username = ?", "testuser").First(&user)
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, "082253739918", user.Phone)
	assert.Equal(t, "user", user.Role)

	// Ensure profile was created
	var profile models.Profile
	database.DB.Where("user_id = ?", user.ID).First(&profile)
	assert.Equal(t, "Test User", profile.FullName)

	// Clean up
	database.DB.Where("phone = ?", "082253739918").Delete(&models.User{})
	database.DB.Where("full_name = ?", "Test User").Delete(&models.Profile{})
}

/*
TestRegister_ExistingUser - Tests registration with a duplicate username or phone number.
*/
func TestRegister_ExistingUser(t *testing.T) {
	tests.SetupTestDB() // ✅ Call shared test helper

	// Create an existing user
	existingUser := models.User{
		Username: "existinguser",
		Phone:    "082253739919",
		Role:     "user",
	}
	database.DB.Create(&existingUser)

	app := tests.SetupApp()
	body := `{
		"username": "existinguser",
		"full_name": "Existing User",
		"phone": "082253739919"
	}`

	req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Parse response body
	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)

	assert.Equal(t, "Username or phone number already taken", response["error"])

	// Clean up
	database.DB.Where("phone = ?", "082253739919").Delete(&models.User{})
}

/*
TestRegister_InvalidRequestBody - Tests registration with an invalid request body.
*/
func TestRegister_InvalidRequestBody(t *testing.T) {
	tests.SetupTestDB() // ✅ Call shared test helper

	app := tests.SetupApp()
	body := `{"username": 123, "full_name": true, "phone": []}` // Invalid types

	req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Parse response body
	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)

	assert.Equal(t, "Invalid request body", response["error"])
}
