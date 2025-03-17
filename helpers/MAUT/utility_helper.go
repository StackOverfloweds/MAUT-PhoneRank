package maut

import (
	"math"

	"github.com/StackOverfloweds/MAUT-PhoneRank/models"
)

// CalculateUtility calculates the MAUT score for each smartphone.
func CalculateUtility(smartphones []models.Smartphone, minMax map[string]float64, weights map[string]float64) {
	for i := range smartphones {
		s := &smartphones[i]

		normalizedPrice := NormalizeValue(s.Price, minMax["minPrice"], minMax["maxPrice"])
		normalizedRAM := NormalizeValue(float64(s.RAMCapacity), minMax["minRAM"], minMax["maxRAM"])
		normalizedCamera := NormalizeValue(s.Camera.PrimaryCameraRear, minMax["minCamera"], minMax["maxCamera"])

		// Perhitungan skor akhir
		score := (weights["price"] * normalizedPrice) +
			(weights["ram"] * normalizedRAM) +
			(weights["camera"] * normalizedCamera)

		s.AvgRating = math.Round(score*1000) / 1000 // Dibulatkan ke 3 desimal
	}
}
