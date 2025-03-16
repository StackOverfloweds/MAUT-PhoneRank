package auth

import (
	"github.com/StackOverfloweds/MAUT-PhoneRank/database"
	"github.com/StackOverfloweds/MAUT-PhoneRank/models"
	"github.com/gofiber/fiber/v2"
)

/*
Register - Handles user registration.
This function accepts a JSON request containing username, full name, and phone number.
It checks if the username or phone number is already taken.
If the user is new, it creates both a User and a Profile in the database.
The role is automatically set to "user" (cannot be changed during registration).
*/
func Register(c *fiber.Ctx) error {
	var input struct {
		Username string `json:"username"`
		FullName string `json:"full_name"`
		Phone    string `json:"phone"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	var existingUser models.User
	if database.DB.Where("username = ? OR phone = ?", input.Username, input.Phone).First(&existingUser).RowsAffected > 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Username or phone number already taken"})
	}

	user := models.User{
		Username: input.Username,
		Phone:    input.Phone,
		Role:     "user",
	}

	if err := database.DB.Create(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create user"})
	}

	profile := models.Profile{
		UserID:   user.ID,
		FullName: input.FullName,
	}

	if err := database.DB.Create(&profile).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create profile"})
	}

	return c.JSON(fiber.Map{
		"message": "User created successfully",
	})
}
