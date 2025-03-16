package user

import (
	"time"

	"github.com/StackOverfloweds/MAUT-PhoneRank/database"
	"github.com/StackOverfloweds/MAUT-PhoneRank/models"
	"github.com/gofiber/fiber/v2"
)

/*
GetProfile - Retrieves the user profile.
This function fetches the profile data of the currently authenticated user.
The user ID is extracted from the request context (`c.Locals("user_id")`).
If the profile is not found, it returns a 404 error.
*/
func GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	var profile models.Profile
	if err := database.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Profile not found"})
	}

	return c.JSON(profile)
}

/*
UpdateProfile - Updates the user's profile.
This function allows authenticated users to update their profile information,
including address, backup phone, birthdate, and gender.
The user ID is extracted from the request context (`c.Locals("user_id")`).
If the profile is found, the provided fields are updated and saved to the database.
*/
func UpdateProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	var input struct {
		Address     *string `json:"address"`
		BackupPhone *string `json:"phone_number"`
		Birthdate   *string `json:"birthdate"`
		Gender      *string `json:"gender"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	var profile models.Profile
	if err := database.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Profile not found"})
	}

	if input.Birthdate != nil {
		birthdate, err := time.Parse("2006-01-02", *input.Birthdate)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid birthdate format. Use YYYY-MM-DD."})
		}
		profile.Birthdate = &birthdate
	}

	if input.Address != nil {
		profile.Address = input.Address
	}
	if input.BackupPhone != nil {
		profile.BackupPhone = input.BackupPhone
	}
	if input.Gender != nil {
		profile.Gender = input.Gender
	}

	database.DB.Save(&profile)

	return c.JSON(fiber.Map{"message": "Profile updated successfully!", "profile": profile})
}
