package Slices

import (
	uc "github.com/PlayerR9/MyGoLib/Units/Common"
)

// Find returns the index of the first occurrence of an element in the slice.
//
// Parameters:
//   - S: slice of elements.
//   - elem: element to find.
//
// Returns:
//   - int: index of the first occurrence of the element or -1 if not found.
func Find[T uc.Equaler[T]](S []T, elem T) int {
	for i, e := range S {
		if e.Equals(elem) {
			return i
		}
	}

	return -1
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
func RemoveDuplicates[T any](S []T) []T {
	if len(S) < 2 {
		return S
	}

	indexMap := make(map[int]T)
	unique := make([]T, 0)

	for i, elem := range S {
		found := false

		f := func(other T) bool {

		}

		for _, value := range indexMap {
			if uc.EqualOf(value, elem) {
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

// Difference returns the elements that are in S1 but not in S2.
//
// Parameters:
//   - S1: The first slice of elements.
//   - S2: The second slice of elements.
func Difference[T comparable](S1, S2 []T) []T {
	if len(S1) == 0 {
		return S2
	} else if len(S2) == 0 {
		return S1
	}

	seen := make(map[T]bool)

	for _, e := range S2 {
		seen[e] = true
	}

	diff := make([]T, 0)

	for _, e := range S1 {
		if _, ok := seen[e]; !ok {
			diff = append(diff, e)
		}
	}

	return diff
}
