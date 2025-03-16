package user

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	user "github.com/StackOverfloweds/MAUT-PhoneRank/controllers/User"
	"github.com/StackOverfloweds/MAUT-PhoneRank/database"
	"github.com/StackOverfloweds/MAUT-PhoneRank/models"
	"github.com/StackOverfloweds/MAUT-PhoneRank/tests"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestGetProfile(t *testing.T) {
	tests.SetupTestDB()

	users := models.User{
		Username: "testuser",
		Phone:    "081234567890",
		Role:     "user",
	}
	database.DB.Create(&users)

	profile := models.Profile{
		UserID:   users.ID,
		FullName: "Test User",
	}
	database.DB.Create(&profile)

	app := fiber.New()
	app.Get("/profile", func(c *fiber.Ctx) error {
		c.Locals("user_id", users.ID)
		return user.GetProfile(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/profile", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response models.Profile
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, profile.FullName, response.FullName)
}

func TestUpdateProfile(t *testing.T) {
	tests.SetupTestDB()

	users := models.User{
		Username: "testuser",
		Phone:    "081234567890",
		Role:     "user",
	}
	database.DB.Create(&users)

	profile := models.Profile{
		UserID:   users.ID,
		FullName: "Test User",
	}
	database.DB.Create(&profile)

	app := fiber.New()
	app.Put("/profile", func(c *fiber.Ctx) error {
		c.Locals("user_id", users.ID)
		return user.UpdateProfile(c)
	})

	body := `{"address": "New Address", "phone_number": "081234567891"}`
	req := httptest.NewRequest(http.MethodPut, "/profile", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var updatedProfile models.Profile
	database.DB.Where("user_id = ?", users.ID).First(&updatedProfile)
	assert.Equal(t, "New Address", *updatedProfile.Address)
	assert.Equal(t, "081234567891", *updatedProfile.BackupPhone)
}
