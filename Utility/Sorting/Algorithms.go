package Sorting

// CompareFunc is a function that compares two elements of a slice.
//
// Parameters:
//   - T: The type of the elements to compare.
//
// Returns:
//   - int: < 0 if e1 < e2, 0 if e1 == e2, > 0 if e1 > e2.
type CompareFunc[T any] func(T, T) int

// SelectionSort sorts a slice of elements using the selection sort algorithm.
//
// Parameters:
//   - elements: The slice of elements to sort.
//   - compare: The function to use to compare two elements of the slice.
//
// Information:
//  - In place: Yes.
func SelectionSort[T any](elements []T, compare CompareFunc[T]) {
	for k := 0; k < len(elements); k++ {
		min := k

		for j := k + 1; j < len(elements); j++ {
			if compare(elements[j], elements[min]) < 0 {
				min = j
			}
		}

		elements[k], elements[min] = elements[min], elements[k]
	}
}

// InsertionSort sorts a slice of elements using the insertion sort algorithm.
//
// Parameters:
//   - elements: The slice of elements to sort.
//   - compare: The function to use to compare two elements of the slice.
//
// Information:
//  - In place: Yes.
func InsertionSort[T any](elements []T, compare CompareFunc[T]) {
	for k := 1; k < len(elements); k++ {
		x := elements[k]
		j := k - 1

		for j >= 0 && compare(elements[j], x) > 0 {
			elements[j+1] = elements[j]
			j--
		}

		elements[j+1] = x
	}
}

// BubbleSort sorts a slice of elements using the bubble sort algorithm.
//
// Parameters:
//   - elements: The slice of elements to sort.
//   - compare: The function to use to compare two elements of the slice.
//
// Information:
//  - In place: Yes.
func BubbleSort[T any](elements []T, compare CompareFunc[T]) {
	swapped := true

	for {
		swapped = false

		for j := 1; j < len(elements)-1; j++ {
			if compare(elements[j], elements[j-1]) > 0 {
				elements[j], elements[j-1] = elements[j-1], elements[j]
				swapped = true
			}
		}

		if !swapped {
			break
		}
	}
}

// merge is a helper function for MergeSort. It merges two sorted slices into a single sorted slice.
//
// Parameters:
//   - left: The left slice to merge.
//   - right: The right slice to merge.
//   - compare: The function to use to compare two elements of the slice.
//
// Returns:
//   - []T: The merged slice.
func merge[T any](left, right []T, compare CompareFunc[T]) []T {
	result := make([]T, len(left)+len(right))
	i1, i2, k := 0, 0, 0

	for i1 < len(left) && i2 < len(right) {
		if compare(left[i1], right[i2]) <= 0 {
			result[k] = left[i1]
			i1++
		} else {
			result[k] = right[i2]
			i2++
		}

		k++
	}

	if i1 < len(left) {
		copy(result[k:], left[i1:])
	} else {
		copy(result[k:], right[i2:])
	}

	return result
}

// MergeSort sorts a slice of elements using the merge sort algorithm.
//
// Parameters:
//   - elements: The slice of elements to sort.
//   - compare: The function to use to compare two elements of the slice.
//
// Information:
//  - In place: Yes
func MergeSort[T any](elements []T, compare CompareFunc[T]) []T {
	if len(elements) > 1 {
		return elements
	}

	middle := len(elements) / 2

	left := MergeSort(elements[:middle], compare)
	right := MergeSort(elements[middle:], compare)

	MergeSort(left, compare)
	MergeSort(right, compare)

	return merge(left, right, compare)
}
