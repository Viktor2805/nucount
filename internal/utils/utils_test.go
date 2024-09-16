package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestRound tests the Round function with various inputs.
func TestRound(t *testing.T) {
	tests := []struct {
		input     float64
		precision int
		expected  float64
	}{
		{124.480341, 1, 124.5},
		{124.480341, 2, 124.48},
		{124.480341, 3, 124.480},
		{124.480341, 4, 124.4803},
		{124.480341, 5, 124.48034},
		{124.480341, 6, 124.480341},
		{1.555, 1, 1.6},
		{1.555, 2, 1.56},
		{1.555, 3, 1.555},
		{1.555, 0, 2},
		{1.554, 0, 2},
	}

	for _, test := range tests {
		result := Round(test.input, test.precision)
		assert.Equal(t, test.expected, result, "Round(%f, %d) = %f; want %f", test.input, test.precision, result, test.expected)
	}

}
