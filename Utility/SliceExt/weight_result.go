package SliceExt

import (
	cdp "github.com/PlayerR9/MyGoLib/CustomData/Pair"
)

// WeightFunc is a type that defines a function that assigns a weight to an element.
//
// Parameters:
//   - elem: The element to assign a weight to.
//
// Returns:
//   - float64: The weight of the element.
//   - bool: True if the weight is valid, otherwise false.
type WeightFunc[T any] func(elem T) (float64, bool)

// WeightResult is a type that represents an element with its corresponding weight.
type WeightResult[T any] cdp.Pair[T, float64]

// NewWeightResult creates a new WeightResult with the given weight and element.
//
// Parameters:
//   - elem: The element.
//   - weight: The weight of the element.
//
// Returns:
//   - WeightResult[T]: The new WeightResult.
func NewWeightResult[T any](elem T, weight float64) WeightResult[T] {
	return WeightResult[T]{
		Second: weight,
		First:  elem,
	}
}

// ApplyWeightFunc is a function that iterates over the slice and applies the weight
// function to each element. The returned slice contains the elements with their
// corresponding weights but only if the weight is valid.
//
// If S is empty, the function returns an empty slice.
//
// Parameters:
//   - S: slice of elements.
//   - f: the weight function.
//
// Returns:
//   - []WeightResult[T]: slice of elements with their corresponding weights.
func ApplyWeightFunc[T any](S []T, f WeightFunc[T]) []WeightResult[T] {
	if len(S) == 0 {
		return []WeightResult[T]{}
	}

	trimmed := make([]WeightResult[T], 0, len(S))

	for _, e := range S {
		weight, ok := f(e)
		if !ok {
			continue
		}

		trimmed = append(trimmed, WeightResult[T]{
			First:  e,
			Second: weight,
		})
	}

	return trimmed
}

// FilterByPositiveWeight is a function that iterates over weight results and
// returns the elements with the maximum weight. If multiple elements have the
// same maximum weight, they are all returned.
//
// If S is empty, the function returns an empty slice.
//
// Parameters:
//   - S: slice of weight results.
//
// Returns:
//   - []T: slice of elements with the maximum weight.
func FilterByPositiveWeight[T any](S []WeightResult[T]) []T {
	if len(S) == 0 {
		return []T{}
	}

	solution := []T{S[0].First}
	maxWeight := S[0].Second

	for _, e := range S[1:] {
		if e.Second > maxWeight {
			maxWeight = e.Second
			solution = []T{e.First}
		} else if e.Second == maxWeight {
			solution = append(solution, e.First)
		}
	}

	return solution
}

// FilterByNegativeWeight is a function that iterates over weight results and
// returns the elements with the minimum weight. If multiple elements have the
// same minimum weight, they are all returned.
//
// If S is empty, the function returns an empty slice.
//
// Parameters:
//   - S: slice of weight results.
//
// Returns:
//   - []T: slice of elements with the minimum weight.
func FilterByNegativeWeight[T any](S []WeightResult[T]) []T {
	if len(S) == 0 {
		return []T{}
	}

	solution := []T{S[0].First}
	minWeight := S[0].Second

	for _, e := range S[1:] {
		if e.Second < minWeight {
			minWeight = e.Second
			solution = []T{e.First}
		} else if e.Second == minWeight {
			solution = append(solution, e.First)
		}
	}

	return solution
}
