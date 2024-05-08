package MapExt

// PredicateFilter is a type that defines a map filter function.
//
// Parameters:
//   - T: The type of the elements in the map.
//
// Returns:
//   - bool: True if the element satisfies the filter function, otherwise false.
type PredicateFilter[T comparable, E any] func(T, E) bool

// Intersect returns a PredicateFilter function that checks if an element
// satisfies all the PredicateFilter functions in funcs.
//
// Parameters:
//   - funcs: A map of PredicateFilter functions.
//
// Returns:
//   - PredicateFilter: A PredicateFilter function that checks if a element satisfies
//     all the PredicateFilter functions in funcs.
//
// Behaviors:
//   - It returns false as soon as it finds a function in funcs that the element
//     does not satisfy.
//   - If no functions are provided, it returns a function that always returns true.
func Intersect[T comparable, E any](funcs ...PredicateFilter[T, E]) PredicateFilter[T, E] {
	if len(funcs) == 0 {
		return func(k T, v E) bool { return true }
	}

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
//   - funcs: A map of PredicateFilter functions.
//
// Returns:
//   - PredicateFilter: A PredicateFilter function that checks if a element satisfies
//     at least one of the PredicateFilter functions in funcs.
//
// Behaviors:
//   - It returns true as soon as it finds a function in funcs that the element
//     satisfies.
//   - If no functions are provided, it returns a function that always returns false.
func Union[T comparable, E any](funcs ...PredicateFilter[T, E]) PredicateFilter[T, E] {
	if len(funcs) == 0 {
		return func(k T, v E) bool { return false }
	}

	return func(k T, v E) bool {
		for _, f := range funcs {
			if f(k, v) {
				return true
			}
		}

		return false
	}
}

// MapFilter is a function that iterates over the
//
// Parameters:
//   - S: map of elements.
//   - filter: function that takes an element and returns a bool.
//
// Returns:
//   - []T: map of elements that satisfy the filter function.
//
// Behaviors:
//   - An element is said to satisfy the filter function if the function returns true
//     for that element.
//   - If S is empty, it returns a non-nil empty map.
//   - If filter is nil, it returns S as is
func MapFilter[T comparable, E any](S map[T]E, filter PredicateFilter[T, E]) map[T]E {
	if len(S) == 0 {
		return make(map[T]E)
	} else if filter == nil {
		return S
	}

	result := make(map[T]E)

	for k, v := range S {
		if filter(k, v) {
			result[k] = v
		}
	}

	return result
}
