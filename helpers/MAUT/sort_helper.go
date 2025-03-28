package maut

import (
	"sort"

	"github.com/StackOverfloweds/MAUT-PhoneRank/models"
)

// SortSmartphonesByScore sorts smartphones by MAUT score (descending order).
func SortSmartphonesByScore(smartphones []models.Smartphone) {
	sort.Slice(smartphones, func(i, j int) bool {
		return smartphones[i].AvgRating > smartphones[j].AvgRating
	})
}
