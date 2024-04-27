package SliceExt

import (
	"slices"

	intf "github.com/PlayerR9/MyGoLib/Units/Interfaces"
)

// PredicateFilter is a type that defines a slice filter function.
//
// Parameters:
//   - T: The type of the elements in the slice.
//
// Returns:
//   - bool: True if the element satisfies the filter function, otherwise false.
type PredicateFilter[T any] func(T) bool

// Intersect returns a PredicateFilter function that checks if an element
// satisfies all the PredicateFilter functions in funcs.
// It returns false as soon as it finds a function in funcs that the element
// does not satisfy.
//
// If no filter functions are provided, then all elements are considered to satisfy
// the filter function.
//
// Parameters:
//   - funcs: A slice of PredicateFilter functions.
//
// Returns:
//   - PredicateFilter: A PredicateFilter function that checks if a element satisfies
//     all the PredicateFilter functions in funcs.
func Intersect[T any](funcs ...PredicateFilter[T]) PredicateFilter[T] {
	return func(e T) bool {
		if len(funcs) == 0 {
			return true
		}

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
// If no filter functions are provided, then all elements are considered to satisfy
// the filter function.
//
// Parameters:
//   - funcs: A slice of PredicateFilter functions.
//
// Returns:
//   - PredicateFilter: A PredicateFilter function that checks if a element satisfies
//     at least one of the PredicateFilter functions in funcs.
func Union[T any](funcs ...PredicateFilter[T]) PredicateFilter[T] {
	return func(e T) bool {
		if len(funcs) == 0 {
			return true
		}

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
// If S is empty, the function returns an empty slice.
//
// Parameters:
//   - S: slice of elements.
//   - filter: function that takes an element and returns a bool.
//
// Returns:
//   - []T: slice of elements that satisfy the filter function.
func SliceFilter[T any](S []T, filter PredicateFilter[T]) []T {
	if len(S) == 0 {
		return []T{}
	}

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
// If S is empty, the function returns two empty slices.
//
// Parameters:
//   - S: slice of elements.
//   - filter: function that takes an element and returns a bool.
//
// Returns:
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
// If S is empty, the function returns an empty slice and true.
//
// Parameters:
//   - S: slice of elements.
//   - filter: function that takes an element and returns a bool.
//
// Returns:
//   - []T: slice of elements that satisfy the filter function or the original slice.
//   - bool: true if there are successful elements, otherwise false.
func SFSeparateEarly[T any](S []T, filter PredicateFilter[T]) ([]T, bool) {
	if len(S) == 0 {
		return []T{}, true
	}

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
// If there are less than two elements in the slice, the function
// returns the original slice.
//
// Parameters:
//   - S: slice of elements.
//
// Returns:
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
// If there are less than two elements in the slice, the function
// returns the original slice.
//
// Parameters:
//   - S: slice of elements.
//   - equals: function that takes two elements and returns a bool.
//
// Returns:
//   - []T: slice of elements with duplicates removed.
func RemoveDuplicatesFunc[T intf.Equaler[T]](S []T) []T {
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
// If there are less than two elements in the slice, the function returns -1.
//
// Parameters:
//   - S: slice of elements.
//
// Returns:
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
// If there are less than two elements in the slice, the function returns -1.
//
// Parameters:
//   - S: slice of elements.
//   - equals: function that takes two elements and returns a bool.
//
// Returns:
//   - int: index of the first duplicate element or -1 if there are no duplicates.
func IndexOfDuplicateFunc[T intf.Equaler[T]](S []T) int {
	if len(S) < 2 {
		return -1
	}

	for i, e := range S {
		for j := i + 1; j < len(S); j++ {
			if e.Equals(S[j]) {
				return i
			}
		}
	}

	return -1
}

// FindSubBytesFrom finds the first occurrence of a subslice in a byte
// slice starting from a given index. If at is negative, it is set to 0.
//
// If S or subS is empty, the function returns -1.
// If the subslice is not found, the function returns -1.
//
// Parameters:
//   - S: The byte slice to search in.
//   - subS: The byte slice to search for.
//   - at: The index to start searching from.
//
// Returns:
//   - int: The index of the first occurrence of the subslice.
func FindSubsliceFrom[T intf.Comparable](S []T, subS []T, at int) int {
	if len(subS) == 0 || len(S) == 0 || at+len(subS) > len(S) {
		return -1
	}

	if at < 0 {
		at = 0
	}

	possibleStarts := make([]int, 0)

	// Find all possible starting points.
	for i := at; i < len(S)-len(subS); i++ {
		if S[i] == subS[0] {
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
			if S[start+j] != subS[j] {
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
// If at is negative, it is set to 0.
//
// If S or subS is empty, the function returns -1.
// If the subslice is not found, the function returns -1.
//
// Parameters:
//   - S: The byte slice to search in.
//   - subS: The byte slice to search for.
//   - at: The index to start searching from.
//   - equals: The comparison function.
//
// Returns:
//   - int: The index of the first occurrence of the subslice.
func FindSubsliceFromFunc[T intf.Equaler[T]](S []T, subS []T, at int) int {
	if len(subS) == 0 || len(S) == 0 || at+len(subS) > len(S) {
		return -1
	}

	if at < 0 {
		at = 0
	}

	possibleStarts := make([]int, 0)

	// Find all possible starting points.
	for i := at; i < len(S)-len(subS); i++ {
		if S[i].Equals(subS[0]) {
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
			if !S[start+j].Equals(subS[j]) {
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
