package maut

// NormalizeValue calculates normalization based on min & max values.
func NormalizeValue(value, min, max float64) float64 {
	if max-min == 0 {
		return 0
	}
	return (value - min) / (max - min)
}
