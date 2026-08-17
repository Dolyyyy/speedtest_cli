package model

import "math"

// Round rounds a float64 value to the given precision decimal places
func Round(val float64, precision int) float64 {
	p := math.Pow(10, float64(precision))
	return math.Round(val*p) / p
}
