package user

import (
	"time"

	"github.com/StackOverfloweds/MAUT-PhoneRank/database"
	"github.com/StackOverfloweds/MAUT-PhoneRank/models"
	"github.com/gofiber/fiber/v2"
)

func GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	var profile models.Profile
	if err := database.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Profile not found"})
	}
	return c.JSON(profile)

}

// Update Profile (Protected)
func UpdateProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	var input struct {
		Address     *string `json:"address"`
		PhoneNumber *string `json:"phone_number"`
		Birthdate   *string `json:"birthdate"`
		Gender      *string `json:"gender"`
	}

	// Parse JSON input
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Find Profile
	var profile models.Profile
	if err := database.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Profile not found"})
	}

	// Convert Birthdate String to Time
	if input.Birthdate != nil {
		birthdate, err := time.Parse("2006-01-02", *input.Birthdate)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid birthdate format. Use YYYY-MM-DD."})
		}
		profile.Birthdate = &birthdate
	}

	// Update Fields if Provided
	if input.Address != nil {
		profile.Address = input.Address
	}
	if input.PhoneNumber != nil {
		profile.PhoneNumber = input.PhoneNumber
	}
	if input.Gender != nil {
		profile.Gender = input.Gender
	}

	// Save Updates
	database.DB.Save(&profile)

	return c.JSON(fiber.Map{"message": "Profile updated successfully!", "profile": profile})
}
