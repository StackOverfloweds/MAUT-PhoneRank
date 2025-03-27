package maut

import (
	"math"

	"github.com/StackOverfloweds/MAUT-PhoneRank/models"
)

// GetMinMaxValues finds the minimum and maximum values for normalization.
func GetMinMaxValues(smartphones []models.Smartphone) map[string]float64 {
	minMax := map[string]float64{
		"minProcessor": math.MaxFloat64, "maxProcessor": 0,
		"minRAM": math.MaxFloat64, "maxRAM": 0,
		"minPrice": math.MaxFloat64, "maxPrice": 0,
		"minDisplay": math.MaxFloat64, "maxDisplay": 0,
	}

	for _, s := range smartphones {
		// Processor
		if s.Processor.Speed < minMax["minProcessor"] {
			minMax["minProcessor"] = s.Processor.Speed
		}
		if s.Processor.Speed > minMax["maxProcessor"] {
			minMax["maxProcessor"] = s.Processor.Speed
		}

		// RAM
		if float64(s.RAMCapacity) < minMax["minRAM"] {
			minMax["minRAM"] = float64(s.RAMCapacity)
		}
		if float64(s.RAMCapacity) > minMax["maxRAM"] {
			minMax["maxRAM"] = float64(s.RAMCapacity)
		}

		// Price
		if s.Price < minMax["minPrice"] {
			minMax["minPrice"] = s.Price
		}
		if s.Price > minMax["maxPrice"] {
			minMax["maxPrice"] = s.Price
		}

		// Display
		if float64(s.Display.RefreshRate) < minMax["minDisplay"] {
			minMax["minDisplay"] = float64(s.Display.RefreshRate)
		}
		if float64(s.Display.RefreshRate) > minMax["maxDisplay"] {
			minMax["maxDisplay"] = float64(s.Display.RefreshRate)
		}
	}

	return minMax
}
