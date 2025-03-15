package auth

import (
	"github.com/StackOverfloweds/MAUT-PhoneRank/database"
	"github.com/StackOverfloweds/MAUT-PhoneRank/models"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *fiber.Ctx) error {
	var input struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
		FullName string `json:"full_name"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"Error": "Invalid request body"})
	}

	//  check if the username already exists
	var existingUser models.User
	if database.DB.Where("username = ? OR email = ?", input.Username, input.Email).First(&existingUser).RowsAffected > 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Username or email already taken"})
	}

	hashPass, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to hash password"})
	}

	//create user
	user := models.User{
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: string(hashPass),
		Role:         "user",
	}

	if err := database.DB.Create(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create user"})
	}

	profile := models.Profile{
		UserID:   user.ID,
		FullName: input.FullName,
	}
	// Save Profile to Database
	if err := database.DB.Create(&profile).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create profile"})
	}

	return c.JSON(fiber.Map{
		"message": "User created successfully",
	})

}
