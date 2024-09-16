// Package utils provides common utility functions.
package utils

import "math"

// Round rounds a float64 to the specified number of decimal places.
func Round(x float64, precision int) float64 {
	pow := math.Pow(10, float64(precision))
	return math.Round(x*pow) / pow
}
