package MathExt

import (
	"math"

	ers "github.com/PlayerR9/MyGoLib/Units/Errors"
)

// AVG calculates the average of a slice of float64 elements.
//
// Parameters:
//   - elems: The elements to calculate the average of.
//
// Returns:
//   - float64: The average of the elements.
//   - error: An error of type *ErrInvalidParameter if the slice is empty.
func AVG(elems []float64) (float64, error) {
	if len(elems) == 0 {
		return 0, ers.NewErrInvalidParameter("elems", ers.NewErrEmpty(elems))
	}

	L := float64(len(elems))

	var sum float64

	for _, elem := range elems {
		sum += elem
	}

	return sum / L, nil
}

// SQM calculates the Standard Quadratic Mean of a slice of float64 elements.
//
// Parameters:
//   - elems: The elements to calculate the SQM of.
//
// Returns:
//   - float64: The SQM of the elements.
//   - error: An error of type *ErrInvalidParameter if the slice is empty.
func SQM(elems []float64) (float64, error) {
	if len(elems) == 0 {
		return 0, ers.NewErrInvalidParameter("elems", ers.NewErrEmpty(elems))
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

	return math.Sqrt(sqm / L), nil
}
