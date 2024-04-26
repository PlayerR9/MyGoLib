package SliceExt

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
type WeightResult[T any] struct {
	// Elem is the element.
	Elem T

	// Weight is the weight of the element.
	Weight float64
}

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
		Weight: weight,
		Elem:   elem,
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
			Elem:   e,
			Weight: weight,
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
// 	 - S: slice of weight results.
//
// Returns:
//   - []T: slice of elements with the maximum weight.
func FilterByPositiveWeight[T any](S []WeightResult[T]) []T {
	if len(S) == 0 {
		return []T{}
	}

	solution := []T{S[0].Elem}
	maxWeight := S[0].Weight

	for _, e := range S[1:] {
		if e.Weight > maxWeight {
			maxWeight = e.Weight
			solution = []T{e.Elem}
		} else if e.Weight == maxWeight {
			solution = append(solution, e.Elem)
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
// 	 - S: slice of weight results.
//
// Returns:
//   - []T: slice of elements with the minimum weight.
func FilterByNegativeWeight[T any](S []WeightResult[T]) []T {
	if len(S) == 0 {
		return []T{}
	}

	solution := []T{S[0].Elem}
	minWeight := S[0].Weight

	for _, e := range S[1:] {
		if e.Weight < minWeight {
			minWeight = e.Weight
			solution = []T{e.Elem}
		} else if e.Weight == minWeight {
			solution = append(solution, e.Elem)
		}
	}

	return solution
}
