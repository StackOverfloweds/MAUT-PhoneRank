package middleware

import (
	"github.com/gofiber/fiber/v2"
)

func CheckAdmin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role")
		if role != "admin" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Access denied. admin only",
			})
		}
		return c.Next()
	}
}
