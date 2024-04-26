package MathExt

import (
	"errors"
	"math"
)

func AVG(elems []float64) (float64, error) {
	if len(elems) == 0 {
		return 0, errors.New("cannot calculate the average of an empty slice")
	}

	L := float64(len(elems))

	var sum float64

	for _, elem := range elems {
		sum += elem
	}

	return sum / L, nil
}

func SQM(elems []float64) (float64, error) {
	if len(elems) == 0 {
		return 0, errors.New("cannot calculate the SQM of an empty slice")
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
