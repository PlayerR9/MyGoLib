package Sorting

import (
	"slices"
)

// SortFunc is a function type that defines a comparison function for elements.
//
// Parameters:
//   - e1: The first element to compare.
//   - e2: The second element to compare.
//
// Returns:
//   - int: A negative value if the first element is less than the second element,
//     0 if they are equal, and a positive value if the first element is greater
//     than the second element.
type SortFunc[T any] func(e1, e2 T) int

// StableSort is a function that sorts a slice of elements in ascending order while preserving
// the order of equal elements.
//
// Parameters:
//   - slice: The slice to sort.
//   - sf: A function that compares two elements.
//   - isAsc: A flag indicating if the sort is in ascending order.
//
// Behaviors:
//   - This function uses the Merge Sort algorithm to sort the slice.
//   - The elements in the slice must implement the Comparer interface.
func StableSort[T any](S []T, sf SortFunc[T], isAsc bool) {
	if len(S) < 2 || sf == nil {
		return
	}

	if !isAsc {
		sf = func(a, b T) int {
			return -sf(a, b)
		}
	}

	slices.SortStableFunc(S, sf)
}

// Sort is a function that sorts a slice of elements in ascending order.
//
// Parameters:
//   - slice: The slice to sort.
//   - sf: A function that compares two elements.
//   - isAsc: A flag indicating if the sort is in ascending order.
//
// Behaviors:
//   - This function uses the Quick Sort algorithm to sort the slice.
//   - The elements in the slice must implement the Comparer interface.
func Sort[T any](S []T, sf SortFunc[T], isAsc bool) {
	if len(S) < 2 || sf == nil {
		return
	}

	if !isAsc {
		sf = func(a, b T) int {
			return -sf(a, b)
		}
	}

	sortQuick(S, 0, len(S)-1, sf)
}

// sortQuick is a helper function that sorts a slice of elements using the Quick Sort algorithm.
//
// Parameters:
//   - S: The slice to sort.
//   - l: The left index of the slice.
//   - r: The right index of the slice.
//   - sf: A function that compares two elements.
func sortQuick[T any](S []T, l, r int, sf SortFunc[T]) {
	if l >= r {
		return
	}

	p := partition(S, l, r, sf)
	sortQuick(S, l, p-1, sf)
	sortQuick(S, p+1, r, sf)
}

// partition is a helper function that partitions a slice of elements for the Quick Sort algorithm.
//
// Parameters:
//   - S: The slice to partition.
//   - l: The left index of the slice.
//   - r: The right index of the slice.
//
// Returns:
//   - int: The index of the pivot element.
func partition[T any](S []T, l, r int, sf SortFunc[T]) int {
	pivot := S[r]
	i := l

	for j := l; j < r; j++ {
		res := sf(S[j], pivot)

		if res < 0 {
			S[i], S[j] = S[j], S[i]
			i++
		}
	}

	S[i], S[r] = S[r], S[i]

	return i
}
