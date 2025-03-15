package auth

import (
	"time"

	"github.com/StackOverfloweds/MAUT-PhoneRank/database"
	"github.com/StackOverfloweds/MAUT-PhoneRank/helpers"
	"github.com/StackOverfloweds/MAUT-PhoneRank/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *fiber.Ctx) error {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	// find user by usernmae
	var user models.User
	if err := database.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid username or password"})
	}

	// compare password

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid username or password"})
	}

	//generate jwt
	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(helpers.GetJWTSecret())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	// Return JWT Token
	return c.JSON(fiber.Map{"token": tokenString})
}
