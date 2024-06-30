package MathExt

import (
	"math"
)

// AVG calculates the average of a slice of float64 elements.
//
// Parameters:
//   - elems: The elements to calculate the average of.
//
// Returns:
//   - float64: The average of the elements.
//   - bool: False if the slice is empty. True otherwise.
func AVG(elems []float64) (float64, bool) {
	if len(elems) == 0 {
		return 0, false
	}

	L := float64(len(elems))

	var sum float64

	for _, elem := range elems {
		sum += elem
	}

	return sum / L, true
}

// SQM calculates the Standard Quadratic Mean of a slice of float64 elements.
//
// Parameters:
//   - elems: The elements to calculate the SQM of.
//
// Returns:
//   - float64: The SQM of the elements.
//   - bool: False if the slice is empty. True otherwise.
func SQM(elems []float64) (float64, bool) {
	if len(elems) == 0 {
		return 0, false
	}

	L := float64(len(elems))

	var average float64

	for _, elem := range elems {
		average += elem
	}

	average /= L

	var sqm float64

	for _, elem := range elems {
		sqm += math.Pow(elem-average, 2)
	}

	res := math.Sqrt(sqm / L)

	return res, true
}
