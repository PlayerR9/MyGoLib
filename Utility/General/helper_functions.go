package General

import (
	"reflect"
)

// SplitIntoGroups splits the slice into n groups and returns a 2D slice where
// each inner slice represents a group.
//
// Parameters:
//   - slice: The slice to split.
//   - n: The number of groups to split the slice into.
//
// Return:
//   - [][]T: A 2D slice where each inner slice represents a group. Nil if
//     n is less than or equal to 0.
//
// Example:
//
//	slice := []int{1, 2, 3, 4, 5}
//	n := 3
//	groups := SplitIntoGroups(slice, n)
//	fmt.Println(groups) // Output: [[1 4] [2 5] [3]]
func SplitIntoGroups[T any](slice []T, n int) [][]T {
	if len(slice) == 0 {
		return [][]T{}
	} else if len(slice) == 1 || n == 1 {
		return [][]T{slice}
	}

	if n <= 0 {
		return nil
	}

	groups := make([][]T, n)

	for index, element := range slice {
		groupNumber := (index + 1) % n

		groups[groupNumber] = append(groups[groupNumber], element)
	}

	return groups
}

// IsNil is a function that checks if a value is nil.
//
// Parameters:
//   - value: The value to check.
//
// Return:
//   - bool: True if the value is nil, false otherwise.
func IsNil[T any](value T) bool {
	kind := reflect.ValueOf(value).Kind()

	switch kind {
	case reflect.Ptr, reflect.Interface,
		reflect.Map, reflect.Slice, reflect.Chan:
		ok := reflect.ValueOf(value).IsNil()
		if ok {
			return true
		}
	}

	return false
}
