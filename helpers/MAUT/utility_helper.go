package maut

import (
	"math"

	"github.com/StackOverfloweds/MAUT-PhoneRank/models"
)

/*
CalculateUtility - Computes the MAUT score for each smartphone
and returns a map of smartphone ID -> calculated utility score.
*/
func CalculateUtility(smartphones []models.Smartphone, minMax map[string]float64, weights map[string]float64) map[string]float64 {
	// Validate and normalize weights
	totalWeight := 0.0
	for _, weight := range weights {
		totalWeight += weight
	}
	if totalWeight != 1.0 {
		for key := range weights {
			weights[key] /= totalWeight
		}
	}

	// Create map to store utility scores
	utilityScores := make(map[string]float64)

	// Calculate utility for each smartphone
	for _, s := range smartphones {
		// Normalize values
		normalizedProcessor := NormalizeValue(s.Processor.Speed, minMax["minProcessor"], minMax["maxProcessor"])
		normalizedRAM := NormalizeValue(float64(s.RAMCapacity), minMax["minRAM"], minMax["maxRAM"])
		normalizedPrice := NormalizeValue(minMax["maxPrice"]-s.Price, minMax["minPrice"], minMax["maxPrice"]) // Lower price is better
		normalizedRearCam := NormalizeValue(s.Camera.PrimaryCameraRear, minMax["minRearCam"], minMax["maxRearCam"])
		normalizedFrontCam := NormalizeValue(s.Camera.PrimaryCameraFront, minMax["minFrontCam"], minMax["maxFrontCam"])

		// Compute utility score
		score := (weights["processor"] * normalizedProcessor) +
			(weights["ram"] * normalizedRAM) +
			(weights["price"] * normalizedPrice) +
			(weights["rear_cam"] * normalizedRearCam) +
			(weights["front_cam"] * normalizedFrontCam)

		// Store the MAUT score (rounded to 3 decimals)
		utilityScores[s.ID] = math.Round(score*1000) / 1000
	}

	return utilityScores
}
