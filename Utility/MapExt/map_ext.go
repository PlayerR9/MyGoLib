package MapExt

// PredicateFilter is a type that defines a slice filter function.
//
// Parameters:
//   - T: The type of the elements in the slice.
//
// Returns:
//   - bool: True if the element satisfies the filter function, otherwise false.
type PredicateFilter[T comparable, E any] func(T, E) bool

// Intersect returns a PredicateFilter function that checks if an element
// satisfies all the PredicateFilter functions in funcs.
// It returns false as soon as it finds a function in funcs that the element
// does not satisfy.
//
// Parameters:
//   - funcs: A slice of PredicateFilter functions.
//
// Returns:
//   - PredicateFilter: A PredicateFilter function that checks if a element satisfies
//     all the PredicateFilter functions in funcs.
func Intersect[T comparable, E any](funcs ...PredicateFilter[T, E]) PredicateFilter[T, E] {
	return func(k T, v E) bool {
		for _, f := range funcs {
			if !f(k, v) {
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
//   - funcs: A slice of PredicateFilter functions.
//
// Returns:
//   - PredicateFilter: A PredicateFilter function that checks if a element satisfies
//     at least one of the PredicateFilter functions in funcs.
func Union[T comparable, E any](funcs ...PredicateFilter[T, E]) PredicateFilter[T, E] {
	return func(k T, v E) bool {
		for _, f := range funcs {
			if f(k, v) {
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
//   - S: slice of elements.
//   - filter: function that takes an element and returns a bool.
//
// Returns:
//   - []T: slice of elements that satisfy the filter function.
func MapFilter[T comparable, E any](S map[T]E, filter PredicateFilter[T, E]) map[T]E {
	result := make(map[T]E)

	for k, v := range S {
		if filter(k, v) {
			result[k] = v
		}
	}

	return result
}
