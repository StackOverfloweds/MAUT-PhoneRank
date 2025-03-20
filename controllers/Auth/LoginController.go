package auth

import (
	"time"

	"github.com/StackOverfloweds/MAUT-PhoneRank/database"
	"github.com/StackOverfloweds/MAUT-PhoneRank/helpers"
	jwts "github.com/StackOverfloweds/MAUT-PhoneRank/helpers/JWTs"
	"github.com/StackOverfloweds/MAUT-PhoneRank/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

/*
Login - Requests an OTP for authentication
This function checks if the phone number is registered.
If registered, it generates and sends an OTP via Fonnte API.
If the phone number is not found, it returns a 404 error.
*/
func Login(c *fiber.Ctx) error {
	var input struct {
		PhoneNumber string `json:"phone_number"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	var user models.User
	if database.DB.Where("phone = ?", input.PhoneNumber).First(&user).RowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Phone number is not registered"})
	}

	otp := helpers.GenerateOTP()
	if err := helpers.SendOTP(input.PhoneNumber, otp); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "OTP sent successfully!"})
}

/*
VerifyOTP - Verifies the OTP entered by the user
It checks if the OTP is valid and belongs to a registered phone number.
If valid, the OTP is deleted and a JWT token is generated.
If invalid or expired, an error is returned.
*/
func VerifyOTP(c *fiber.Ctx) error {
	var input struct {
		OTP string `json:"otp"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Find phone number associated with OTP
	phoneNumber, err := helpers.FindPhoneByOTP(input.OTP)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid or expired OTP"})
	}

	helpers.DeleteOTP(phoneNumber)

	var user models.User
	result := database.DB.Where("phone = ?", phoneNumber).First(&user)
	if result.RowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(), // Expire in 24 hours
	})

	tokenString, err := token.SignedString(jwts.GetJWTSecret())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	return c.JSON(fiber.Map{
		"message": "OTP verification successful!",
		"token":   tokenString,
		"role":    user.Role,
	})
}
