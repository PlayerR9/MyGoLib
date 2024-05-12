package SliceExt

import (
	uc "github.com/PlayerR9/MyGoLib/Units/Common"
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
//
// Parameters:
//   - funcs: A slice of PredicateFilter functions.
//
// Returns:
//   - PredicateFilter: A PredicateFilter function that checks if a element satisfies
//     all the PredicateFilter functions in funcs.
//
// Behavior:
//   - If no filter functions are provided, then all elements are considered to satisfy
//     the filter function.
//   - It returns false as soon as it finds a function in funcs that the element
//     does not satisfy.
func Intersect[T any](funcs ...PredicateFilter[T]) PredicateFilter[T] {
	if len(funcs) == 0 {
		return func(elem T) bool { return true }
	}

	return func(elem T) bool {
		for _, f := range funcs {
			if !f(elem) {
				return false
			}
		}

		return true
	}
}

// ParallelIntersect returns a PredicateFilter function that checks if an element
// satisfies all the PredicateFilter functions in funcs concurrently.
//
// Parameters:
//   - funcs: A slice of PredicateFilter functions.
//
// Returns:
//   - PredicateFilter: A PredicateFilter function that checks if a element satisfies
//     all the PredicateFilter functions in funcs.
//
// Behavior:
//   - If no filter functions are provided, then all elements are considered to satisfy
//     the filter function.
//   - It returns false as soon as it finds a function in funcs that the element
//     does not satisfy.
func ParallelIntersect[T any](funcs ...PredicateFilter[T]) PredicateFilter[T] {
	if len(funcs) == 0 {
		return func(elem T) bool { return true }
	}

	return func(elem T) bool {
		resultChan := make(chan bool, len(funcs))

		for _, f := range funcs {
			go func(f PredicateFilter[T]) {
				resultChan <- f(elem)
			}(f)
		}

		for range funcs {
			if !<-resultChan {
				return false
			}
		}

		return true
	}
}

// Union returns a PredicateFilter function that checks if an element
// satisfies at least one of the PredicateFilter functions in funcs.
//
// Parameters:
//   - funcs: A slice of PredicateFilter functions.
//
// Returns:
//   - PredicateFilter: A PredicateFilter function that checks if a element satisfies
//     at least one of the PredicateFilter functions in funcs.
//
// Behavior:
//   - If no filter functions are provided, then no elements are considered to satisfy
//     the filter function.
//   - It returns true as soon as it finds a function in funcs that the element
//     satisfies.
func Union[T any](funcs ...PredicateFilter[T]) PredicateFilter[T] {
	if len(funcs) == 0 {
		return func(elem T) bool { return false }
	}

	return func(elem T) bool {
		for _, f := range funcs {
			if f(elem) {
				return true
			}
		}

		return false
	}
}

// ParallelUnion returns a PredicateFilter function that checks if an element
// satisfies at least one of the PredicateFilter functions in funcs concurrently.
//
// Parameters:
//   - funcs: A slice of PredicateFilter functions.
//
// Returns:
//   - PredicateFilter: A PredicateFilter function that checks if a element satisfies
//     at least one of the PredicateFilter functions in funcs.
//
// Behavior:
//   - If no filter functions are provided, then no elements are considered to satisfy
//     the filter function.
//   - It returns true as soon as it finds a function in funcs that the element
//     satisfies.
func ParallelUnion[T any](funcs ...PredicateFilter[T]) PredicateFilter[T] {
	if len(funcs) == 0 {
		return func(elem T) bool { return false }
	}

	return func(elem T) bool {
		resultChan := make(chan bool, len(funcs))

		for _, f := range funcs {
			go func(f PredicateFilter[T]) {
				resultChan <- f(elem)
			}(f)
		}

		for range funcs {
			if <-resultChan {
				return true
			}
		}

		return false
	}
}

// SliceFilter is a function that iterates over the slice and applies the filter
// function to each element.
//
// Parameters:
//   - S: slice of elements.
//   - filter: function that takes an element and returns a bool.
//
// Returns:
//   - []T: slice of elements that satisfy the filter function.
//
// Behavior:
//   - If S is empty, the function returns a nil slice.
//   - If S has only one element and it satisfies the filter function, the function
//     returns a slice with that element. Otherwise, it returns a nil slice.
//   - An element is said to satisfy the filter function if the function returns true
//     when applied to the element.
func SliceFilter[T any](S []T, filter PredicateFilter[T]) []T {
	solution := make([]T, 0)

	for _, item := range S {
		if filter(item) {
			solution = append(solution, item)
		}
	}

	return solution
}

// FilterNilValues is a function that iterates over the slice and removes the
// nil elements.
//
// Parameters:
//   - S: slice of elements.
//
// Returns:
//   - []*T: slice of elements that satisfy the filter function.
//
// Behavior:
//   - If S is empty, the function returns a nil slice.
func FilterNilValues[T any](S []*T) []*T {
	solution := make([]*T, 0)

	for _, item := range S {
		if item != nil {
			solution = append(solution, item)
		}
	}

	return solution
}

// SFSeparate is a function that iterates over the slice and applies the filter
// function to each element. The returned slices contain the elements that
// satisfy and do not satisfy the filter function.
//
// Parameters:
//   - S: slice of elements.
//   - filter: function that takes an element and returns a bool.
//
// Returns:
//   - []T: slice of elements that satisfy the filter function.
//   - []T: slice of elements that do not satisfy the filter function.
//
// Behavior:
//   - If S is empty, the function returns two empty slices.
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
//   - S: slice of elements.
//   - filter: function that takes an element and returns a bool.
//
// Returns:
//   - []T: slice of elements that satisfy the filter function or the original slice.
//   - bool: true if there are successful elements, otherwise false.
//
// Behavior:
//   - If S is empty, the function returns an empty slice and true.
func SFSeparateEarly[T any](S []T, filter PredicateFilter[T]) ([]T, bool) {
	if len(S) == 0 {
		return []T{}, true
	}

	success := make([]T, 0)
	for _, item := range S {
		if filter(item) {
			success = append(success, item)
		}
	}

	if len(success) == 0 {
		return S, false
	}

	return success, true
}

// RemoveDuplicates removes duplicate elements from the slice.
//
// Parameters:
//   - S: slice of elements.
//
// Returns:
//   - []T: slice of elements with duplicates removed.
//
// Behavior:
//   - The function preserves the order of the elements in the slice.
//   - If there are multiple duplicates of an element, only the first
//     occurrence is kept.
//   - If there are less than two elements in the slice, the function
//     returns the original slice.
func RemoveDuplicates[T comparable](S []T) []T {
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
//   - S: slice of elements.
//   - equals: function that takes two elements and returns a bool.
//
// Returns:
//   - []T: slice of elements with duplicates removed.
//
// Behavior:
//   - The function preserves the order of the elements in the slice.
//   - If there are multiple duplicates of an element, only the first
//     occurrence is kept.
//   - If there are less than two elements in the slice, the function
//     returns the original slice.
func RemoveDuplicatesFunc[T uc.Equaler[T]](S []T) []T {
	if len(S) < 2 {
		return S
	}

	indexMap := make(map[int]T)
	unique := make([]T, 0)

	for i, elem := range S {
		found := false

		for _, value := range indexMap {
			if value.Equals(elem) {
				found = true
				break
			}
		}

		if !found {
			unique = append(unique, elem)
			indexMap[i] = elem
		}
	}

	return unique
}

// IndexOfDuplicate returns the index of the first duplicate element in the slice.
//
// Parameters:
//   - S: slice of elements.
//
// Returns:
//   - int: index of the first duplicate element or -1 if there are no duplicates.
func IndexOfDuplicate[T comparable](S []T) int {
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
//
// Parameters:
//   - S: slice of elements.
//
// Returns:
//   - int: index of the first duplicate element or -1 if there are no duplicates.
func IndexOfDuplicateFunc[T uc.Equaler[T]](S []T) int {
	if len(S) < 2 {
		return -1
	}

	indexMap := make(map[int]T)

	for i, elem := range S {
		found := false

		for _, value := range indexMap {
			if value.Equals(elem) {
				found = true
				break
			}
		}

		if found {
			return i
		}

		indexMap[i] = elem
	}

	return -1
}

// computeLPSArray is a helper function that computes the Longest Prefix
// Suffix (LPS) array for the Knuth-Morris-Pratt algorithm.
//
// Parameters:
//   - subS: The subslice to compute the LPS array for.
//   - lps: The LPS array to store the results in.
//
// Behavior:
//   - The function modifies the lps array in place.
//   - The lps array is initialized with zeros.
//   - The lps array is used to store the length of the longest prefix
//     that is also a suffix for each index in the subslice.
//   - The first element of the lps array is always 0.
func computeLPSArray[T comparable](subS []T, lps []int) {
	length := 0
	i := 1
	lps[0] = 0 // lps[0] is always 0

	// the loop calculates lps[i] for i = 1 to len(subS)-1
	for i < len(subS) {
		if subS[i] == subS[length] {
			length++
			lps[i] = length
			i++
		} else {
			if length != 0 {
				length = lps[length-1]
			} else {
				lps[i] = 0
				i++
			}
		}
	}
}

// FindSubBytesFrom finds the first occurrence of a subslice in a byte
// slice starting from a given index.
//
// Parameters:
//   - S: The byte slice to search in.
//   - subS: The byte slice to search for.
//   - at: The index to start searching from.
//
// Returns:
//   - int: The index of the first occurrence of the subslice.
//
// Behavior:
//   - The function uses the Knuth-Morris-Pratt algorithm to find the subslice.
//   - If S or subS is empty, the function returns -1.
//   - If the subslice is not found, the function returns -1.
//   - If at is negative, it is set to 0.
func FindSubsliceFrom[T comparable](S []T, subS []T, at int) int {
	if len(subS) == 0 || len(S) == 0 || at+len(subS) > len(S) {
		return -1
	}

	if at < 0 {
		at = 0
	}

	lps := make([]int, len(subS))
	computeLPSArray(subS, lps)

	i := at
	j := 0
	for i < len(S) {
		if S[i] == subS[j] {
			i++
			j++
		}

		if j == len(subS) {
			return i - j
		} else if i < len(S) && S[i] != subS[j] {
			if j != 0 {
				j = lps[j-1]
			} else {
				i = i + 1
			}
		}
	}

	return -1
}

// computeLPSArrayFunc is a helper function that computes the Longest Prefix
// Suffix (LPS) array for the Knuth-Morris-Pratt algorithm using a custom
// comparison function.
//
// Parameters:
//   - subS: The subslice to compute the LPS array for.
//   - lps: The LPS array to store the results in.
//
// Behavior:
//   - The function modifies the lps array in place.
//   - The lps array is initialized with zeros.
//   - The lps array is used to store the length of the longest prefix
//     that is also a suffix for each index in the subslice.
//   - The first element of the lps array is always 0.
func computeLPSArrayFunc[T uc.Equaler[T]](subS []T, lps []int) {
	length := 0
	i := 1
	lps[0] = 0 // lps[0] is always 0

	// the loop calculates lps[i] for i = 1 to len(subS)-1
	for i < len(subS) {
		if subS[i].Equals(subS[length]) {
			length++
			lps[i] = length
			i++
		} else {
			if length != 0 {
				length = lps[length-1]
			} else {
				lps[i] = 0
				i++
			}
		}
	}
}

// FindSubsliceFromFunc finds the first occurrence of a subslice in a byte
// slice starting from a given index using a custom comparison function.
//
// Parameters:
//   - S: The byte slice to search in.
//   - subS: The byte slice to search for.
//   - at: The index to start searching from.
//
// Returns:
//   - int: The index of the first occurrence of the subslice.
//
// Behavior:
//   - The function uses the Knuth-Morris-Pratt algorithm to find the subslice.
//   - If S or subS is empty, the function returns -1.
//   - If the subslice is not found, the function returns -1.
//   - If at is negative, it is set to 0.
func FindSubsliceFromFunc[T uc.Equaler[T]](S []T, subS []T, at int) int {
	if len(subS) == 0 || len(S) == 0 || at+len(subS) > len(S) {
		return -1
	}

	if at < 0 {
		at = 0
	}

	lps := make([]int, len(subS))
	computeLPSArrayFunc(subS, lps)

	i := at
	j := 0
	for i < len(S) {
		if S[i].Equals(subS[j]) {
			i++
			j++
		}

		if j == len(subS) {
			return i - j
		} else if i < len(S) && !S[i].Equals(subS[j]) {
			if j != 0 {
				j = lps[j-1]
			} else {
				i = i + 1
			}
		}
	}

	return -1
}
