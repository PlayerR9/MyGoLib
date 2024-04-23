package SliceExt

import (
	"slices"

	intf "github.com/PlayerR9/MyGoLib/Units/Interfaces"
	util "github.com/PlayerR9/MyGoLib/util"
)

// FilterByPositiveWeight is a function that iterates over the slice and applies
// the weight function to each element. The returned slice contains the elements
// with the maximum weight. If multiple elements have the same maximum weight,
// they are all returned.
//
// Parameters:
//
//   - S: slice of elements.
//   - weightFunc: function that takes an element and returns an integer.
//
// Returns:
//
//   - []T: slice of elements with the maximum weight.
func FilterByPositiveWeight[T any](S []T, weightFunc func(T) (int, bool)) []T {
	if len(S) == 0 {
		return []T{}
	}

	var solution []T
	var maxWeight int

	index := -1

	for i, e := range S {
		currentValue, ok := weightFunc(e)
		if ok {
			solution = []T{e}
			maxWeight = currentValue
			index = i
			break
		}
	}

	if index == -1 {
		return []T{}
	}

	for _, e := range S[index+1:] {
		currentValue, ok := weightFunc(e)
		if !ok {
			continue
		}

		if currentValue > maxWeight {
			maxWeight = currentValue

			// Clear the solution and add the new element
			for i := range solution {
				solution[i] = *new(T)
			}

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
//   - S: slice of elements.
//   - weightFunc: function that takes an element and returns an integer.
//
// Returns:
//
//   - []T: slice of elements with the minimum weight.
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

// PredicateFilter is a type that defines a slice filter function.
//
// Parameters:
//
//   - T: The type of the elements in the slice.
//
// Returns:
//
//   - bool: True if the element satisfies the filter function, otherwise false.
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
//   - PredicateFilter: A PredicateFilter function that checks if a element satisfies
//     all the PredicateFilter functions in funcs.
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
//   - PredicateFilter: A PredicateFilter function that checks if a element satisfies
//     at least one of the PredicateFilter functions in funcs.
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
//   - S: slice of elements.
//   - filter: function that takes an element and returns a bool.
//
// Returns:
//
//   - []T: slice of elements that satisfy the filter function.
func SliceFilter[T any](S []T, filter PredicateFilter[T]) []T {
	solution := make([]T, len(S))
	copy(solution, S)

	pos := 0
	for _, item := range S {
		if filter(item) {
			solution[pos] = item
			pos++
		}
	}

	return solution[:pos]
}

// SFSeparate is a function that iterates over the slice and applies the filter
// function to each element. The returned slices contain the elements that
// satisfy and do not satisfy the filter function.
//
// Parameters:
//
//   - S: slice of elements.
//   - filter: function that takes an element and returns a bool.
//
// Returns:
//
//   - []T: slice of elements that satisfy the filter function.
//   - []T: slice of elements that do not satisfy the filter function.
func SFSeparate[T any](S []T, filter PredicateFilter[T]) ([]T, []T) {
	success := make([]T, 0)
	failed := make([]T, 0)

	for _, item := range S {
		if filter(item) {
			success = append(success, item)
		} else {
			failed = append(failed, item)
		}
	}

	return success, failed
}

// SFSeparateEarly is a variant of SFSeparate that returns all successful elements.
// If there are none, it returns the original slice and false.
//
// Parameters:
//
//   - S: slice of elements.
//   - filter: function that takes an element and returns a bool.
//
// Returns:
//
//   - []T: slice of elements that satisfy the filter function or the original slice.
//   - bool: true if there are successful elements, otherwise false.
func SFSeparateEarly[T any](S []T, filter PredicateFilter[T]) ([]T, bool) {
	index := slices.IndexFunc(S, filter)
	if index == -1 {
		return S, false
	}

	success := make([]T, len(S)-index)
	copy(success, S[index:])

	pos := 0
	for _, item := range S[index:] {
		if filter(item) {
			success[pos] = item
			pos++
		}
	}

	return success[:pos], true
}

// RemoveDuplicates removes duplicate elements from the slice.
//
// Parameters:
//
//   - S: slice of elements.
//
// Returns:
//
//   - []T: slice of elements with duplicates removed.
func RemoveDuplicates[T intf.Comparable](S []T) []T {
	if len(S) < 2 {
		return S
	}

	seen := make(map[T]bool)

	unique := make([]T, 0, len(S))

	for _, e := range S {
		if _, ok := seen[e]; !ok {
			seen[e] = true
			unique = append(unique, e)
		}
	}

	return unique
}

// RemoveDuplicatesFunc removes duplicate elements from the slice.
//
// Parameters:
//
//   - S: slice of elements.
//   - equals: function that takes two elements and returns a bool.
//
// Returns:
//
//   - []T: slice of elements with duplicates removed.
func RemoveDuplicatesFunc[T interface {
	Equals(T) bool
}](S []T) []T {
	if len(S) < 2 {
		return S
	}

	unique := make([]T, 0, len(S))

	for _, e := range S {
		found := false

		for _, u := range unique {
			if u.Equals(e) {
				found = true
				break
			}
		}

		if !found {
			unique = append(unique, e)
		}
	}

	return unique
}

// IndexOfDuplicate returns the index of the first duplicate element in the slice.
// If there are no duplicates, it returns -1.
//
// Parameters:
//
//   - S: slice of elements.
//
// Returns:
//
//   - int: index of the first duplicate element or -1 if there are no duplicates.
func IndexOfDuplicate[T intf.Comparable](S []T) int {
	if len(S) < 2 {
		return -1
	}

	seen := make(map[T]bool)

	for i, e := range S {
		if _, ok := seen[e]; ok {
			return i
		}

		seen[e] = true
	}

	return -1
}

// IndexOfDuplicateFunc returns the index of the first duplicate element in the slice.
// If there are no duplicates, it returns -1.
//
// Parameters:
//
//   - S: slice of elements.
//   - equals: function that takes two elements and returns a bool.
//
// Returns:
//
//   - int: index of the first duplicate element or -1 if there are no duplicates.
func IndexOfDuplicateFunc[T any](S []T, equals func(T, T) bool) int {
	if len(S) < 2 {
		return -1
	}

	for i, e := range S {
		for j := i + 1; j < len(S); j++ {
			if equals(e, S[j]) {
				return i
			}
		}
	}

	return -1
}

// FindSubBytesFrom finds the first occurrence of a subslice in a byte
// slice starting from a given index.
//
// If the subslice is not found, the function returns -1.
//
// Parameters:
//
//   - s: The byte slice to search in.
//   - subS: The byte slice to search for.
//   - at: The index to start searching from.
//
// Returns:
//
//   - int: The index of the first occurrence of the subslice.
func FindSubsliceFrom[T intf.Comparable](s []T, subS []T, at int) int {
	if len(subS) == 0 || len(s) == 0 || at+len(subS) > len(s) {
		return -1
	}

	if at < 0 {
		at = 0
	}

	possibleStarts := make([]int, 0)

	// Find all possible starting points.
	for i := at; i < len(s)-len(subS); i++ {
		if s[i] == subS[0] {
			possibleStarts = append(possibleStarts, i)
		}
	}

	// Check only the possible starting points that have enough space
	// to contain the subslice in full.
	top := 0

	for i := 0; i < len(possibleStarts)-1; i++ {
		if possibleStarts[i+1]-possibleStarts[i] >= len(subS) {
			possibleStarts[top] = possibleStarts[i]
			top++
		}
	}

	possibleStarts = possibleStarts[:top]

	// Check if the subslice is present at any of the possible starting points
	for _, start := range possibleStarts {
		found := true

		for j := 0; j < len(subS); j++ {
			if s[start+j] != subS[j] {
				found = false
				break
			}
		}

		if found {
			return start
		}
	}

	return -1
}

// FindSubsliceFromFunc finds the first occurrence of a subslice in a byte
// slice starting from a given index using a custom comparison function.
//
// If the subslice is not found, the function returns -1.
//
// Parameters:
//
//   - s: The byte slice to search in.
//   - subS: The byte slice to search for.
//   - at: The index to start searching from.
//   - equals: The comparison function.
//
// Returns:
//
//   - int: The index of the first occurrence of the subslice.
func FindSubsliceFromFunc[T util.Equaler[T]](s []T, subS []T, at int) int {
	if len(subS) == 0 || len(s) == 0 || at+len(subS) > len(s) {
		return -1
	}

	if at < 0 {
		at = 0
	}

	possibleStarts := make([]int, 0)

	// Find all possible starting points.
	for i := at; i < len(s)-len(subS); i++ {
		if s[i].Equals(subS[0]) {
			possibleStarts = append(possibleStarts, i)
		}
	}

	// Check only the possible starting points that have enough space
	// to contain the subslice in full.
	top := 0

	for i := 0; i < len(possibleStarts)-1; i++ {
		if possibleStarts[i+1]-possibleStarts[i] >= len(subS) {
			possibleStarts[top] = possibleStarts[i]
			top++
		}
	}

	possibleStarts = possibleStarts[:top]

	// Check if the subslice is present at any of the possible starting points
	for _, start := range possibleStarts {
		found := true

		for j := 0; j < len(subS); j++ {
			if !s[start+j].Equals(subS[j]) {
				found = false
				break
			}
		}

		if found {
			return start
		}
	}

	return -1
}
