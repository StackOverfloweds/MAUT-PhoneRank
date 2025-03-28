package brand

import (
	"github.com/StackOverfloweds/MAUT-PhoneRank/database"
	"github.com/gofiber/fiber/v2"
)

/*
GetAllBrand - Retrieves all brand names from the brands table.

Returns:
  - JSON list of brand names.
  - 500 if there is a database error.
  - 404 if no brands are found.
*/
func GetAllBrand(c *fiber.Ctx) error {
	var brandNames []string

	// Fetch all brand names from the brands table
	err := database.DB.Table("brands").Pluck("name", &brandNames).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to retrieve brand names", "details": err.Error()})
	}

	// If no brands found, return 404
	if len(brandNames) == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "No brands found"})
	}

	return c.JSON(fiber.Map{"brands": brandNames})
}
