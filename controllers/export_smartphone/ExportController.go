package exportsmartphone

import (
	"encoding/json"
	"log"
	"os"

	"github.com/StackOverfloweds/MAUT-PhoneRank/database"
	"github.com/StackOverfloweds/MAUT-PhoneRank/models"
	"github.com/gofiber/fiber/v2"
)

// ExportJSONSmartphone - Mengekspor data smartphone ke file JSON
func ExportJSONSmartphone(c *fiber.Ctx) error {
	var smartphones []models.Smartphone

	// Query database dengan preload relasi
	if err := database.DB.
		Preload("Brand").
		Preload("Processor").
		Preload("Battery").
		Preload("Display").
		Preload("Camera").
		Find(&smartphones).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch smartphones"})
	}

	if len(smartphones) == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "No smartphones found"})
	}

	// Array untuk menyimpan data JSON
	var smartphoneList []fiber.Map

	for _, phone := range smartphones {
		smartphoneList = append(smartphoneList, fiber.Map{
			"brand": phone.Brand.Name,
			"model": phone.Model,
		})
	}

	// Simpan ke file JSON
	file, err := os.Create("smartphones.json")
	if err != nil {
		log.Println("Error creating file:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create JSON file"})
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Format JSON agar lebih rapi
	if err := encoder.Encode(smartphoneList); err != nil {
		log.Println("Error writing JSON:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to write JSON file"})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Smartphones exported successfully",
		"file":    "smartphones.json",
	})
}
