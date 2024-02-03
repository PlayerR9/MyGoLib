package General

import (
	"errors"
	"fmt"
	"reflect"

	ers "github.com/PlayerR9/MyGoLib/Utility/Errors"
)

// DeepCopy is a function that performs a deep copy of a given value.
// It takes a parameter, value, of any type and returns a new value that is a
// deep copy of the input value.
//
// The function first gets the type of the input value using the reflect.TypeOf
// function.
//
// If the kind of the type is not a pointer, the function performs a shallow
// copy of the value by simply returning the input value.
// This is because non-pointer values in Go are passed by value, so a new
// copy is created when the value is passed to the function.
//
// If the kind of the type is a pointer, the function creates a new instance
// of the underlying type of the pointer using the reflect.New function.
// The new instance is then set to the value pointed to by the input pointer
// using the reflect.ValueOf function and the Set method of the reflect.Value type.
//
// The function then returns the new value, which is a deep copy of the input value.
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

// SplitIntoGroups takes a slice of any type and an integer n as input.
// It splits the slice into n groups and returns a 2D slice where each
// inner slice represents a group.
// The elements are distributed based on their index in the original slice.
// If there are more elements than groups, the remaining elements are
// distributed to the groups in a round-robin fashion.
// If the number of groups (n) is less than or equal to 0, it panics with
// the message "The number of groups must be positive and non-zero".
// If the length of the slice is less than 2 or n is 1, it returns the
// original slice wrapped in a 2D slice.
//
// Example:
//
//	slice := []int{1, 2, 3, 4, 5, 6}
//	n := 3
//	groups := SplitIntoGroups(slice, n)
//	fmt.Println(groups) // Output: [[1 4] [2 5] [3 6]]
func SplitIntoGroups[T any](slice []T, n int) [][]T {
	if len(slice) == 0 {
		return [][]T{}
	} else if len(slice) == 1 || n == 1 {
		return [][]T{slice}
	}

	if n < 0 {
		panic(ers.NewErrInvalidParameter(
			"n", fmt.Errorf("negative group number (%d) are not allowed", n),
		))
	} else if n == 0 {
		panic(ers.NewErrInvalidParameter(
			"n", errors.New("cannot split into 0 groups"),
		))
	}

	groups := make([][]T, n)

	for index, element := range slice {
		groupNumber := (index + 1) % n

		groups[groupNumber] = append(groups[groupNumber], element)
	}

	return groups
}

func TransformToPointer[T any](value T) *T {
	return &value
}
