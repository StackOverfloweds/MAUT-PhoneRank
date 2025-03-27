package maut

import (
	"fmt"
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
		normalizedDisplay := NormalizeValue(float64(s.Display.RefreshRate), minMax["minDisplay"], minMax["maxDisplay"])

		// Compute utility score
		score := (weights["processor"] * normalizedProcessor) +
			(weights["ram"] * normalizedRAM) +
			(weights["price"] * normalizedPrice) +
			(weights["display"] * normalizedDisplay)

		// Store the MAUT score (rounded to 3 decimals)
		utilityScores[s.ID] = math.Round(score*1000) / 1000

	}
	fmt.Println("Utility Score : ", utilityScores)

	return utilityScores
}
