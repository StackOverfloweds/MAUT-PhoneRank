package auth

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

/*
Logout - Handles user logout by validating and removing the token.
This function checks if the Authorization header contains a valid Bearer token.
If valid, the server instructs the client to remove the token.
Note: The token itself is not stored server-side, so invalidation is handled on the client.
*/
func Logout(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")

	// Check if the Authorization header is missing
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Extract token from "Bearer token" format
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid Authorization header format"})
	}

	// Instruct the client to remove the token (no server-side storage of tokens)
	return c.JSON(fiber.Map{"message": "Logout successful. Please remove the token on the client side."})
}
