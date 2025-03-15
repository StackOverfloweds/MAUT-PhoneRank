package user

import (
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
