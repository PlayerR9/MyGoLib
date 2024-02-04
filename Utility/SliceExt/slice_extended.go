package SliceExt

// FilterByPositiveWeight is a function that iterates over the slice and applies
// the weight function to each element. The returned slice contains the elements
// with the maximum weight. If multiple elements have the same maximum weight,
// they are all returned.
//
// Parameters:
//
// 	- S: slice of elements.
// 	- weightFunc: function that takes an element and returns an integer.
//
// Returns:
//
// 	- slice of elements with the maximum weight.
func FilterByPositiveWeight[T any](S []T, weightFunc func(T) int) []T {
	if len(S) == 0 {
		return []T{}
	}

	solution := []T{S[0]}

	maxWeight := weightFunc(S[0])

	for _, e := range S[1:] {
		currentValue := weightFunc(e)

		if currentValue > maxWeight {
			maxWeight = currentValue
			solution = []T{e}
		} else if currentValue == maxWeight {
			solution = append(solution, e)
		}
	}

	return solution
}

// FilterByNegativeWeight is a function that iterates over the slice and applies
// the weight function to each element. The returned slice contains the elements
// with the minimum weight. If multiple elements have the same minimum weight,
// they are all returned.
//
// Parameters:
//
// 	- S: slice of elements.
// 	- weightFunc: function that takes an element and returns an integer.
//
// Returns:
//
// 	- slice of elements with the minimum weight.
func FilterByNegativeWeight[T any](S []T, weightFunc func(T) int) []T {
	if len(S) == 0 {
		return []T{}
	}

	solution := []T{S[0]}

	minWeight := weightFunc(S[0])

	for _, e := range S[1:] {
		currentValue := weightFunc(e)

		if currentValue < minWeight {
			minWeight = currentValue
			solution = []T{e}
		} else if currentValue == minWeight {
			solution = append(solution, e)
		}
	}

	return solution
}

// PredicateFilter is a type that defines a function that takes an element
// and returns a bool.
type PredicateFilter[T any] func(T) bool

// Intersect returns a PredicateFilter function that checks if an element
// satisfies all the PredicateFilter functions in funcs.
// It returns false as soon as it finds a function in funcs that the element
// does not satisfy.
//
// Parameters:
//
//   - funcs: A slice of PredicateFilter functions.
//
// Returns:
//
//   - A PredicateFilter function that checks if a element satisfies all the
//     PredicateFilter functions in funcs.
func Intersect[T any](funcs ...PredicateFilter[T]) PredicateFilter[T] {
	return func(e T) bool {
		for _, f := range funcs {
			if !f(e) {
				return false
			}
		}

		return true
	}
}

// Union returns a PredicateFilter function that checks if an element
// satisfies at least one of the PredicateFilter functions in funcs.
// It returns true as soon as it finds a function in funcs that the element
// satisfies.
//
// Parameters:
//
//   - funcs: A slice of PredicateFilter functions.
//
// Returns:
//
//   - A PredicateFilter function that checks if a element satisfies at least
//     one of the PredicateFilter functions in funcs.
func Union[T any](funcs ...PredicateFilter[T]) PredicateFilter[T] {
	return func(e T) bool {
		for _, f := range funcs {
			if f(e) {
				return true
			}
		}

		return false
	}
}

// SliceFilter is a function that iterates over the slice and applies the filter
// function to each element. The returned slice contains the elements that
// satisfy the filter function.
//
// Parameters:
//
// 	- S: slice of elements.
// 	- filter: function that takes an element and returns a bool.
//
// Returns:
//
// 	- slice of elements that satisfy the filter function.
func SliceFilter[T any](S []T, filter PredicateFilter[T]) []T {
	solution := make([]T, len(S))
	copy(solution, S)

	pos := 0
	for _, item := range solution {
		if filter(item) {
			solution[pos] = item
			pos++
		}
	}

	return solution[:pos]
}
