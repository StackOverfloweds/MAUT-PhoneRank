package auth

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

func Logout(c *fiber.Ctx) error {
	autHeader := c.Get("Authorization")

	if autHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	//extract token from "bearer token"
	tokenPart := strings.Split(autHeader, "")
	if len(tokenPart) != 2 || tokenPart[0] != "Bearer" {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid Authorization header format"})
	}

	return c.JSON(fiber.Map{"message": "Logout successful. Please remove the token on the client side."})
}
