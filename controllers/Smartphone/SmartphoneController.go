package smartphone

import (
	"fmt"

	"github.com/StackOverfloweds/MAUT-PhoneRank/database"
	maut "github.com/StackOverfloweds/MAUT-PhoneRank/helpers/MAUT"
	"github.com/StackOverfloweds/MAUT-PhoneRank/models"
	"github.com/gofiber/fiber/v2"
)

/*
SearchSmartphone - Filters `getSmartphone` based on user criteria.
It correctly references the related `Brand`, `Processor`, etc.

Request Body (JSON):

	{
	  "brand": "Samsung",
	  "min_price": 4000000,
	  "max_price": 8000000,
	  "min_ram": 6
	}

Returns:
  - A JSON list of smartphones matching the filters.
*/
func SearchSmartphone(c *fiber.Ctx) error {
	// Define input struct
	type SearchRequest struct {
		Brand    string  `json:"brand,omitempty"`
		MinPrice float64 `json:"min_price,omitempty"`
		MaxPrice float64 `json:"max_price,omitempty"`
		MinRAM   int     `json:"min_ram,omitempty"`
	}

	var req SearchRequest

	// Parse JSON request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON input"})
	}

	// Initialize GORM query
	query := database.DB.
		Preload("Brand").
		Preload("Processor").
		Preload("Battery").
		Preload("Display").
		Preload("Camera")

	// Apply filters based on the request
	if req.Brand != "" {
		query = query.Joins("JOIN brands ON brands.id = smartphones.brand_id").
			Where("LOWER(brands.name) = LOWER(?)", req.Brand)
	}
	if req.MinPrice > 0 {
		query = query.Where("price >= ?", req.MinPrice)
	}
	if req.MaxPrice > 0 {
		query = query.Where("price <= ?", req.MaxPrice)
	}
	if req.MinRAM > 0 {
		query = query.Where("ram_capacity >= ?", req.MinRAM)
	}

	// Log the query for debugging
	fmt.Println("Executing query:", query.Debug().Statement.SQL.String())

	// Execute the query
	var smartphones []models.Smartphone
	if err := query.Find(&smartphones).Error; err != nil {
		fmt.Println("Database error:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch smartphones"})
	}

	// Return the filtered smartphones
	return c.JSON(smartphones)
}

/*
CalculateMAUT - Processes the MAUT (Multi-Attribute Utility Theory) calculation for smartphones.
This function determines the best smartphone based on predefined criteria weights.

Process:
  - Retrieves all smartphone data from the database.
  - Calculates the minimum and maximum values for each attribute.
  - Normalizes the data using MAUT normalization.
  - Computes the utility score for each smartphone.
  - Sorts smartphones from highest to lowest utility score.

Returns:
  - A JSON list of smartphones ranked by their MAUT utility score.
*/
func CalculateMAUT(c *fiber.Ctx) error {
	var smartphones []models.Smartphone
	database.DB.Preload("Processor").Preload("Camera").Find(&smartphones)

	minMaxValues := maut.GetMinMaxValues(smartphones)

	weights := map[string]float64{
		"price":  0.3,
		"ram":    0.4,
		"camera": 0.3,
	}

	maut.CalculateUtility(smartphones, minMaxValues, weights)

	maut.SortSmartphonesByScore(smartphones)

	return c.Status(200).JSON(smartphones)
}

/*
GetSmartphoneDetail - Retrieves detailed information about a specific smartphone.
This function fetches the smartphone details along with related data.
If the smartphone is not found, it returns a 404 error.
*/
func GetSmartphoneDetail(c *fiber.Ctx) error {
	smartphoneID := c.Params("id")
	if smartphoneID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Smartphone ID is required"})
	}

	var smartphone models.Smartphone

	if err := database.DB.
		Preload("Brand").
		Preload("Processor").
		Preload("Battery").
		Preload("Display").
		Preload("Camera").
		First(&smartphone, "id = ?", smartphoneID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Smartphone not found"})
	}

	return c.JSON(smartphone)
}
