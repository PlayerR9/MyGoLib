package General

import (
	"reflect"

	itf "github.com/PlayerR9/MyGoLib/Units/Common"
	ers "github.com/PlayerR9/MyGoLib/Units/Errors"
)

// Min is a function that takes two parameters, a and b, of any type T that
// implements the Comparable interface and returns the smaller of the two values.
//
// Parameters:
//   - a, b: The two values to compare.
//
// Return:
//   - T: The smaller of the two values.
func Min[T itf.Comparable](a, b T) T {
	if a < b {
		return a
	} else {
		return b
	}
}

// Max is a function that takes two parameters, a and b, of any type T that
// implements the Comparable interface and returns the larger of the two values.
//
// Parameters:
//   - a, b: The two values to compare.
//
// Return:
//   - T: The larger of the two values.
func Max[T itf.Comparable](a, b T) T {
	if a < b {
		return b
	} else {
		return a
	}
}

// DeepCopy is a function that performs a deep copy of a given value.
//
// Parameters:
//   - value: The value to copy.
//
// Return:
//   - any: A deep copy of the input value.
func DeepCopy(value any) any {
	typ := reflect.TypeOf(value)

	if typ.Kind() != reflect.Ptr {
		// Perform a shallow copy if the value is not a pointer
		return value
	}

	// Create a new instance of the underlying type
	newValue := reflect.New(typ.Elem()).Interface()

	// Use reflection to copy the value
	reflect.ValueOf(newValue).Elem().Set(reflect.ValueOf(value).Elem())

	return newValue
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

	if n < 0 {
		return nil, ers.NewErrInvalidParameter(
			"n",
			ers.NewErrGTE(0),
		)
	} else if n == 0 {
		return nil, ers.NewErrInvalidParameter(
			"n",
			ers.NewErrInvalidValue(0),
		)
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
