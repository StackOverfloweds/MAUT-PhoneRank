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
		"minRearCam": math.MaxFloat64, "maxRearCam": 0,
		"minFrontCam": math.MaxFloat64, "maxFrontCam": 0,
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

		// Rear Camera
		if s.Camera.PrimaryCameraRear < minMax["minRearCam"] {
			minMax["minRearCam"] = s.Camera.PrimaryCameraRear
		}
		if s.Camera.PrimaryCameraRear > minMax["maxRearCam"] {
			minMax["maxRearCam"] = s.Camera.PrimaryCameraRear
		}

		// Front Camera
		if s.Camera.PrimaryCameraFront < minMax["minFrontCam"] {
			minMax["minFrontCam"] = s.Camera.PrimaryCameraFront
		}
		if s.Camera.PrimaryCameraFront > minMax["maxFrontCam"] {
			minMax["maxFrontCam"] = s.Camera.PrimaryCameraFront
		}
	}

	return minMax
}
