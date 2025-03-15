package database

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/StackOverfloweds/MAUT-PhoneRank/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ImportSmartphones reads data from CSV and inserts it into the database
func ImportSmartphones(db *gorm.DB) {
	// Open CSV file
	file, err := os.Open("data/smartphones_converted.csv")
	if err != nil {
		log.Fatal("❌ Failed to open CSV file:", err)
	}
	defer file.Close()

	// Read CSV content
	reader := csv.NewReader(file)
	reader.Comma = ';' // Use semicolon as delimiter
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal("❌ Failed to read CSV file:", err)
	}

	// Skip header row
	for i, row := range records {
		if i == 0 {
			continue
		}

		// Extract values from CSV
		brandName := row[0]
		model := row[1]
		price, _ := strconv.ParseFloat(row[2], 64)
		avgRating, _ := strconv.ParseFloat(row[3], 64)
		is5G := row[4] == "1"
		processorBrand := row[5]
		numCores, _ := strconv.Atoi(row[6])
		processorSpeed, _ := strconv.ParseFloat(row[7], 64)
		batteryCapacity, _ := strconv.Atoi(row[8])
		fastChargingAvailable := row[9] == "1"
		fastCharging, _ := strconv.Atoi(row[10])
		ramCapacity, _ := strconv.Atoi(row[11])
		internalMemory, _ := strconv.Atoi(row[12])
		screenSize, _ := strconv.ParseFloat(row[13], 64)
		refreshRate, _ := strconv.Atoi(row[14])
		numRearCameras, _ := strconv.Atoi(row[15])
		osType := row[16]
		primaryCameraRear, _ := strconv.ParseFloat(row[17], 64)
		primaryCameraFront, _ := strconv.ParseFloat(row[18], 64)
		extendedMemoryAvailable := row[19] == "1"
		resolutionHeight, _ := strconv.Atoi(row[20])
		resolutionWidth, _ := strconv.Atoi(row[21])

		// 1️⃣ Insert Brand (if not exists)
		var brand models.Brand
		db.FirstOrCreate(&brand, models.Brand{Name: brandName})

		// 2️⃣ Insert Processor (if not exists)
		var processor models.Processor
		db.FirstOrCreate(&processor, models.Processor{
			Brand:    processorBrand,
			Model:    fmt.Sprintf("%s %d-Core", processorBrand, numCores),
			NumCores: numCores,
			Speed:    processorSpeed,
		})

		// 3️⃣ Insert Battery (if not exists)
		var battery models.Battery
		db.FirstOrCreate(&battery, models.Battery{
			Capacity:              batteryCapacity,
			FastChargingAvailable: fastChargingAvailable,
			FastCharging:          fastCharging,
		})

		// 4️⃣ Insert Display (if not exists)
		var display models.Display
		db.FirstOrCreate(&display, models.Display{
			ScreenSize:       screenSize,
			RefreshRate:      refreshRate,
			ResolutionWidth:  resolutionWidth,
			ResolutionHeight: resolutionHeight,
		})

		// 5️⃣ Insert Camera (if not exists)
		var camera models.Camera
		db.FirstOrCreate(&camera, models.Camera{
			NumRearCameras:     numRearCameras,
			PrimaryCameraRear:  primaryCameraRear,
			PrimaryCameraFront: primaryCameraFront,
		})

		// 6️⃣ Check if Smartphone Already Exists (Avoid Duplicates)
		var existingSmartphone models.Smartphone
		result := db.Where("brand_id = ? AND model = ?", brand.ID, model).First(&existingSmartphone)
		if result.RowsAffected > 0 {
			continue // Skip inserting this smartphone
		}

		// 7️⃣ Insert Smartphone (only if not exists)
		smartphone := models.Smartphone{
			ID:                      uuid.New().String(),
			BrandID:                 brand.ID,
			Model:                   model,
			Price:                   price,
			AvgRating:               avgRating,
			Is5G:                    is5G,
			ProcessorID:             processor.ID,
			BatteryID:               battery.ID,
			DisplayID:               display.ID,
			CameraID:                camera.ID,
			RAMCapacity:             ramCapacity,
			InternalMemory:          internalMemory,
			OS:                      osType,
			ExtendedMemoryAvailable: extendedMemoryAvailable,
		}

		db.Create(&smartphone)
	}

	fmt.Println("🚀 CSV Data Successfully Imported!")
}
