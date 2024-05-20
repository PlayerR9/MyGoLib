package General

import (
	"reflect"

	uc "github.com/PlayerR9/MyGoLib/Units/Common"
	ers "github.com/PlayerR9/MyGoLib/Units/Errors"
)

// Min is a function that takes two parameters, a and b, of any type T
// according to the uc.CompareOf function and returns the smaller of the two values.
//
// Parameters:
//   - a, b: The two values to compare.
//
// Return:
//   - T: The smaller of the two values.
//   - bool: True if the values are comparable. False otherwise
func Min[T any](a, b T) (T, bool) {
	res, ok := uc.CompareOf(a, b)
	if !ok {
		return *new(T), false
	}

	if res < 0 {
		return a, true
	} else {
		return b, true
	}
}

// Max is a function that takes two parameters, a and b, of any type T
// according to the uc.CompareOf function and returns the larger of the two values.
//
// Parameters:
//   - a, b: The two values to compare.
//
// Return:
//   - T: The larger of the two values.
//   - bool: True if the values are comparable. False otherwise
func Max[T any](a, b T) (T, bool) {
	res, ok := uc.CompareOf(a, b)
	if !ok {
		return *new(T), false
	}

	if res > 0 {
		return a, true
	} else {
		return b, true
	}
}

// SplitIntoGroups splits the slice into n groups and returns a 2D slice where
// each inner slice represents a group.
//
// Parameters:
//   - slice: The slice to split.
//   - n: The number of groups to split the slice into.
//
// Return:
//   - [][]T: A 2D slice where each inner slice represents a group.
//   - error: An error of type *ers.ErrInvalidParameter if n is less than or equal to 0.
//
// Example:
//
//	slice := []int{1, 2, 3, 4, 5}
//	n := 3
//	groups := SplitIntoGroups(slice, n)
//	fmt.Println(groups) // Output: [[1 4] [2 5] [3]]
func SplitIntoGroups[T any](slice []T, n int) ([][]T, error) {
	if len(slice) == 0 {
		return [][]T{}, nil
	} else if len(slice) == 1 || n == 1 {
		return [][]T{slice}, nil
	}

	if n <= 0 {
		return nil, ers.NewErrInvalidParameter("n", ers.NewErrGT(0))
	}

	groups := make([][]T, n)

	for index, element := range slice {
		groupNumber := (index + 1) % n

		groups[groupNumber] = append(groups[groupNumber], element)
	}

	return groups, nil
}

// IsNil is a function that checks if a value is nil.
//
// Parameters:
//   - value: The value to check.
//
// Return:
//   - bool: True if the value is nil, false otherwise.
func IsNil[T any](value T) bool {
	switch reflect.ValueOf(value).Kind() {
	case reflect.Ptr, reflect.Interface,
		reflect.Map, reflect.Slice, reflect.Chan:
		if reflect.ValueOf(value).IsNil() {
			return true
		}
	}

	return false
}
