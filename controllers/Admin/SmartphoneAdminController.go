package admin

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/StackOverfloweds/MAUT-PhoneRank/database"
	"github.com/StackOverfloweds/MAUT-PhoneRank/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

/*
*
Create for smartphone
*
*/
func CreateSmartphone(c *fiber.Ctx) error {
	var raw map[string]interface{}
	if err := c.BodyParser(&raw); err != nil {
		log.Println("Failed to parse JSON:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid JSON format",
			"err":   err.Error(),
		})
	}

	// Begin the transaction
	tx := database.DB.Begin() // Assuming you have a `database.DB` for your GORM instance

	// Ensure to rollback if any error occurs
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var brand models.Brand

	// Check if the 'brand' field exists and extract the name from raw data
	if brandMap, ok := raw["brand"].(map[string]interface{}); ok {
		if name, ok := brandMap["name"].(string); ok {
			// Cek apakah brand sudah ada di database
			err := tx.Where("name = ?", name).First(&brand).Error
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					// Brand tidak ditemukan, maka buat yang baru
					brand = models.Brand{Name: name}
					if err := tx.Create(&brand).Error; err != nil {
						tx.Rollback()
						log.Println("Failed to create brand:", err)
						return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
							"error": "Failed to create brand",
						})
					}
				} else {
					// Jika error lain (misalnya koneksi), rollback dan kembalikan error
					tx.Rollback()
					log.Println("Error while checking for existing brand:", err)
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
						"error": "Error while checking for existing brand",
					})
				}
			} else {
				// Brand ditemukan, tidak perlu create lagi
				log.Println("Brand already exists, using existing one:", brand.Name)
			}
		}
	}

	// Display
	var display models.Display
	if displayMap, ok := raw["display"].(map[string]interface{}); ok {
		if v, ok := displayMap["screen_size"].(float64); ok {
			display.ScreenSize = v
		}
		if v, ok := displayMap["refresh_rate"].(float64); ok {
			display.RefreshRate = int(v)
		}
		if v, ok := displayMap["resolution_width"].(float64); ok {
			display.ResolutionWidth = int(v)
		}
		if v, ok := displayMap["resolution_height"].(float64); ok {
			display.ResolutionHeight = int(v)
		}
	}

	// Save Display to DB
	if err := tx.Create(&display).Error; err != nil {
		tx.Rollback()
		log.Println("Failed to create display:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create display",
		})
	}

	// Battery
	var battery models.Battery
	if batteryMap, ok := raw["battery"].(map[string]interface{}); ok {
		if v, ok := batteryMap["capacity"].(float64); ok {
			battery.Capacity = int(v)
		}
		if v, ok := batteryMap["fast_charging_available"].(bool); ok {
			battery.FastChargingAvailable = v
		}
		if v, ok := batteryMap["fast_charging"].(float64); ok {
			battery.FastCharging = int(v)
		}
	}

	// Save Battery to DB
	if err := tx.Create(&battery).Error; err != nil {
		tx.Rollback()
		log.Println("Failed to create battery:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create battery",
		})
	}

	// Camera
	var camera models.Camera
	if cameraMap, ok := raw["camera"].(map[string]interface{}); ok {
		if v, ok := cameraMap["num_rear_cameras"].(float64); ok {
			camera.NumRearCameras = int(v)
		}
		if v, ok := cameraMap["primary_camera_rear"].(float64); ok {
			camera.PrimaryCameraRear = v
		}
		if v, ok := cameraMap["primary_camera_front"].(float64); ok {
			camera.PrimaryCameraFront = v
		}
	}

	// Save Camera to DB
	if err := tx.Create(&camera).Error; err != nil {
		tx.Rollback()
		log.Println("Failed to create camera:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create camera",
		})
	}

	// Processor
	var processor models.Processor
	if processorMap, ok := raw["processor"].(map[string]interface{}); ok {
		if v, ok := processorMap["brand"].(string); ok {
			processor.Brand = v
		}
		if v, ok := processorMap["model"].(string); ok {
			processor.Model = v
		}
		if v, ok := processorMap["num_cores"].(float64); ok {
			processor.NumCores = int(v)
		}
		if v, ok := processorMap["speed"].(float64); ok {
			processor.Speed = v
		}
	}

	// Save Processor to DB
	if err := tx.Create(&processor).Error; err != nil {
		tx.Rollback()
		log.Println("Failed to create processor:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create processor",
		})
	}

	// Smartphone
	var smartphone models.Smartphone
	if model, ok := raw["model"].(string); ok {
		smartphone.Model = model
	}
	if price, ok := raw["price"].(float64); ok {
		smartphone.Price = price
	}
	if avgRating, ok := raw["avg_rating"].(float64); ok {
		smartphone.AvgRating = avgRating
	}
	if is5G, ok := raw["is_5G"].(bool); ok {
		smartphone.Is5G = is5G
	}
	if ram, ok := raw["ram_capacity"].(float64); ok {
		smartphone.RAMCapacity = int(ram)
	}
	if memory, ok := raw["internal_memory"].(float64); ok {
		smartphone.InternalMemory = int(memory)
	}
	if os, ok := raw["os"].(string); ok {
		smartphone.OS = os
	}
	if extendedMemory, ok := raw["extended_memory_available"].(bool); ok {
		smartphone.ExtendedMemoryAvailable = extendedMemory
	}

	// Assign foreign keys for relations
	smartphone.BrandID = brand.ID
	smartphone.DisplayID = display.ID
	smartphone.BatteryID = battery.ID
	smartphone.CameraID = camera.ID
	smartphone.ProcessorID = processor.ID

	// Save Smartphone to DB
	if err := tx.Create(&smartphone).Error; err != nil {
		tx.Rollback()
		log.Println("Failed to create smartphone:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create smartphone",
		})
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		log.Println("Failed to commit transaction:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to commit transaction",
		})
	}

	// Response
	return c.Status(200).JSON(fiber.Map{
		"message":    "Smartphone data created successfully",
		"smartphone": smartphone,
		"brand":      brand,
		"display":    display,
		"battery":    battery,
		"camera":     camera,
		"processor":  processor,
	})
}

func UpdateSmartphone(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Println("Updating smartphone with ID:", id)

	var raw map[string]interface{}
	if err := c.BodyParser(&raw); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid JSON"})
	}

	tx := database.DB.Begin()

	var smartphone models.Smartphone
	if err := tx.Preload("Display").Preload("Battery").Preload("Camera").Preload("Processor").Where("id = ?", id).First(&smartphone).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Smartphone not found"})
	}

	// Helper untuk parse string -> float64
	parseFloat := func(val interface{}) (float64, error) {
		switch v := val.(type) {
		case float64:
			return v, nil
		case string:
			return strconv.ParseFloat(v, 64)
		default:
			return 0, fmt.Errorf("invalid type")
		}
	}

	// Helper untuk parse string -> int
	parseInt := func(val interface{}) (int, error) {
		f, err := parseFloat(val)
		return int(f), err
	}

	// Update brand
	if brandMap, ok := raw["brand"].(map[string]interface{}); ok {
		if name, ok := brandMap["name"].(string); ok {
			var brand models.Brand
			if err := tx.Where("name = ?", name).First(&brand).Error; err != nil {
				brand = models.Brand{Name: name}
				if err := tx.Create(&brand).Error; err != nil {
					tx.Rollback()
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create brand"})
				}
			}
			smartphone.BrandID = brand.ID
		}
	}

	// Update Display
	if displayMap, ok := raw["display"].(map[string]interface{}); ok {
		screenSize, _ := parseFloat(displayMap["screen_size"])
		refreshRate, _ := parseInt(displayMap["refresh_rate"])
		resW, _ := parseInt(displayMap["resolution_width"])
		resH, _ := parseInt(displayMap["resolution_height"])

		tx.Model(&smartphone.Display).Updates(models.Display{
			ScreenSize:       screenSize,
			RefreshRate:      refreshRate,
			ResolutionWidth:  resW,
			ResolutionHeight: resH,
		})
	}

	// Update Battery
	if batteryMap, ok := raw["battery"].(map[string]interface{}); ok {
		capacity, _ := parseInt(batteryMap["capacity"])
		fastCharging, _ := parseInt(batteryMap["fast_charging"])
		fastAvailable, _ := batteryMap["fast_charging_available"].(bool)

		tx.Model(&smartphone.Battery).Updates(models.Battery{
			Capacity:              capacity,
			FastChargingAvailable: fastAvailable,
			FastCharging:          fastCharging,
		})
	}

	// Update Camera
	if cameraMap, ok := raw["camera"].(map[string]interface{}); ok {
		numRear, _ := parseInt(cameraMap["num_rear_cameras"])
		rear, _ := parseFloat(cameraMap["primary_camera_rear"])
		front, _ := parseFloat(cameraMap["primary_camera_front"])

		tx.Model(&smartphone.Camera).Updates(models.Camera{
			NumRearCameras:     numRear,
			PrimaryCameraRear:  rear,
			PrimaryCameraFront: front,
		})
	}

	// Update Processor
	if processorMap, ok := raw["processor"].(map[string]interface{}); ok {
		numCores, _ := parseInt(processorMap["num_cores"])
		speed, _ := parseFloat(processorMap["speed"])
		brand, _ := processorMap["brand"].(string)
		model, _ := processorMap["model"].(string)

		tx.Model(&smartphone.Processor).Updates(models.Processor{
			Brand:    brand,
			Model:    model,
			NumCores: numCores,
			Speed:    speed,
		})
	}

	// Update Smartphone
	modelStr, _ := raw["model"].(string)
	osStr, _ := raw["os"].(string)
	is5G, _ := raw["is_5G"].(bool)
	extMemAvail, _ := raw["extended_memory_available"].(bool)

	price, _ := parseFloat(raw["price"])
	avgRating, _ := parseFloat(raw["avg_rating"])
	ram, _ := parseInt(raw["ram_capacity"])
	internalMem, _ := parseInt(raw["internal_memory"])

	tx.Model(&smartphone).Updates(models.Smartphone{
		Model:                   modelStr,
		Price:                   price,
		AvgRating:               avgRating,
		Is5G:                    is5G,
		RAMCapacity:             ram,
		InternalMemory:          internalMem,
		OS:                      osStr,
		ExtendedMemoryAvailable: extMemAvail,
	})

	if err := tx.Commit().Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update smartphone"})
	}

	return c.JSON(fiber.Map{"message": "Smartphone updated successfully"})
}

func DeleteSmartphone(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID tidak boleh kosong",
		})
	}

	tx := database.DB.Begin()
	if tx.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memulai transaksi"})
	}

	var smartphone models.Smartphone
	if err := tx.First(&smartphone, "id = ?", id).Error; err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Smartphone tidak ditemukan"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mencari smartphone"})
	}

	// Hapus smartphone utama dalam transaksi
	if err := tx.Delete(&smartphone).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menghapus smartphone"})
	}

	// Commit dulu agar constraint dilepas
	if err := tx.Commit().Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menyelesaikan transaksi"})
	}

	// Baru setelah itu, hapus komponen terkait di luar transaksi
	db := database.DB // gunakan db normal
	db.Delete(&models.Display{}, "id = ?", smartphone.DisplayID)
	db.Delete(&models.Battery{}, "id = ?", smartphone.BatteryID)
	db.Delete(&models.Camera{}, "id = ?", smartphone.CameraID)
	db.Delete(&models.Processor{}, "id = ?", smartphone.ProcessorID)

	return c.JSON(fiber.Map{"message": "Smartphone berhasil dihapus"})
}
