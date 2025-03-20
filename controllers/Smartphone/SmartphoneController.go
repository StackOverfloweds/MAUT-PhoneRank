package smartphone

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/StackOverfloweds/MAUT-PhoneRank/database"
	img "github.com/StackOverfloweds/MAUT-PhoneRank/helpers/Image"
	maut "github.com/StackOverfloweds/MAUT-PhoneRank/helpers/MAUT"
	"github.com/StackOverfloweds/MAUT-PhoneRank/models"
	"github.com/gofiber/fiber/v2"
)

/*
SearchSmartphoneWithMAUT - Filters smartphones based on user input,
then applies MAUT ranking to determine the best choices.

Request Body (JSON):

	{
	  "brand": "Samsung",
	  "min_price": 4000000,
	  "max_price": 8000000,
	  "min_ram": 6
	}

Returns:
  - A ranked list of smartphones based on MAUT analysis.
*/
func SearchSmartphoneWithMAUT(c *fiber.Ctx) error {
	type SearchRequest struct {
		Brand    string  `json:"brand,omitempty"`
		MinPrice float64 `json:"min_price,omitempty"`
		MaxPrice float64 `json:"max_price,omitempty"`
		MinRAM   int     `json:"min_ram,omitempty"`
	}

	var req SearchRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON input"})
	}

	query := database.DB.
		Preload("Brand").
		Preload("Processor").
		Preload("Battery").
		Preload("Display").
		Preload("Camera")

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

	var smartphones []models.Smartphone
	if err := query.Find(&smartphones).Error; err != nil {
		fmt.Println("Database error:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch smartphones"})
	}

	if len(smartphones) == 0 {
		fmt.Println("⚠️ No smartphones found matching criteria")
		return c.Status(404).JSON(fiber.Map{"error": "No smartphones found matching criteria"})
	}

	/*
		Calculation using MAUT
	*/

	minMaxValues := maut.GetMinMaxValues(smartphones)

	weights := map[string]float64{
		"processor": 0.25,
		"ram":       0.25,
		"price":     0.2,
		"rear_cam":  0.15,
		"front_cam": 0.15,
	}

	totalWeight := 0.0
	for _, weight := range weights {
		totalWeight += weight
	}
	if totalWeight != 1.0 {
		for key := range weights {
			weights[key] /= totalWeight
		}
	}

	maut.CalculateUtility(smartphones, minMaxValues, weights)

	maut.SortSmartphonesByScore(smartphones)

	var smartphoneDetails []fiber.Map

	for _, phone := range smartphones {
		imageURL, err := img.SearchSmartphoneImage(phone.Brand.Name, phone.Model)
		if err != nil {
			imageURL = ""
		}

		// Append each smartphone's data with image
		smartphoneDetails = append(smartphoneDetails, fiber.Map{
			"smartphone": phone,
			"image_url":  imageURL,
		})
	}

	fmt.Println("MAUT ranking completed")

	return c.Status(200).JSON(smartphoneDetails)
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

	imageURL, err := img.SearchSmartphoneImage(smartphone.Brand.Name, smartphone.Model)

	if err != nil || imageURL == "" || strings.Contains(imageURL, "No+Image") {
		imageURL = ""
	}

	return c.JSON(fiber.Map{
		"smartphone": smartphone,
		"image_url":  imageURL,
	})
}

func GetAllSmartphones(c *fiber.Ctx) error {
	var smartphones []models.Smartphone

	var totalCount int64
	if err := database.DB.Model(&models.Smartphone{}).Count(&totalCount).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to count smartphones"})
	}

	limit, err := strconv.Atoi(c.Query("limit", "30"))
	if err != nil || limit < 1 {
		limit = 30
	}

	if limit > int(totalCount) {
		limit = int(totalCount)
	}

	offset := 0
	if totalCount > int64(limit) {
		rand.Seed(time.Now().UnixNano()) // Seed random
		offset = rand.Intn(int(totalCount) - limit + 1)
	}

	if err := database.DB.
		Preload("Brand").
		Preload("Processor").
		Preload("Battery").
		Preload("Display").
		Preload("Camera").
		Offset(offset).
		Limit(limit).
		Find(&smartphones).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch smartphones"})
	}

	if len(smartphones) == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "No smartphones found"})
	}

	type result struct {
		Phone    models.Smartphone
		ImageURL string
	}

	results := make([]result, len(smartphones))
	var wg sync.WaitGroup
	wg.Add(len(smartphones))

	for i, phone := range smartphones {
		go func(i int, phone models.Smartphone) {
			defer wg.Done()
			imageURL, err := img.SearchSmartphoneImage(phone.Brand.Name, phone.Model)
			if err != nil {
				imageURL = ""
			}
			results[i] = result{Phone: phone, ImageURL: imageURL}
		}(i, phone)
	}

	wg.Wait()

	var smartphoneDetails []fiber.Map
	for _, res := range results {
		if res.ImageURL != "" {
			smartphoneDetails = append(smartphoneDetails, fiber.Map{
				"smartphone": res.Phone,
				"image_url":  res.ImageURL,
			})
		}
	}
	// If no smartphones have valid ImageURLs, return an error
	if len(smartphoneDetails) == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "No smartphones with valid images found"})
	}

	return c.Status(200).JSON(smartphoneDetails)
}
