package maut

import (
	"math"

	"github.com/StackOverfloweds/MAUT-PhoneRank/models"
)

// GetMinMaxValues finds the minimum and maximum values ​​for normalization.
func GetMinMaxValues(smartphones []models.Smartphone) map[string]float64 {
	minMax := map[string]float64{
		"minPrice": math.MaxFloat64, "maxPrice": 0,
		"minRAM": math.MaxFloat64, "maxRAM": 0,
		"minCamera": math.MaxFloat64, "maxCamera": 0,
	}

	for _, s := range smartphones {
		// Harga
		if s.Price < minMax["minPrice"] {
			minMax["minPrice"] = s.Price
		}
		if s.Price > minMax["maxPrice"] {
			minMax["maxPrice"] = s.Price
		}

		// RAM
		if float64(s.RAMCapacity) < minMax["minRAM"] {
			minMax["minRAM"] = float64(s.RAMCapacity)
		}
		if float64(s.RAMCapacity) > minMax["maxRAM"] {
			minMax["maxRAM"] = float64(s.RAMCapacity)
		}

		// Kamera belakang
		if s.Camera.PrimaryCameraRear < minMax["minCamera"] {
			minMax["minCamera"] = s.Camera.PrimaryCameraRear
		}
		if s.Camera.PrimaryCameraRear > minMax["maxCamera"] {
			minMax["maxCamera"] = s.Camera.PrimaryCameraRear
		}
	}

	return minMax
}
