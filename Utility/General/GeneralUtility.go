// git tag v0.1.33

package General

import (
	"fmt"
	"os"
	"reflect"
)

func WaitForUserConfirmation() {
	fmt.Println("Press enter to proceed...")
	fmt.Scanln()
}

// ExitFromProgram is a utility function that handles program termination in
// case of an error.
// It prints the error message if the error is not nil, prompts the user to
// press enter to exit, and then terminates the program with a status code of 1.
//
// Parameters:
//
//   - err: The error that caused the program termination. If err is nil, no error
//     message is printed.
//
// Note: This function does not return. The program is terminated by calling
// os.Exit(1).
func ExitFromProgram(err error) {
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Press enter to exit...")
	fmt.Scanln()

	os.Exit(1)
}

func FindMaximumValue[T any](comparisonFunc func(T, T) int, inputValues ...T) T {
	if len(inputValues) == 0 {
		panic("Cannot find maximum value in an empty set")
	}

	maxValue := inputValues[0]
	for _, currentValue := range inputValues[1:] {
		if comparisonFunc(currentValue, maxValue) > 0 {
			maxValue = currentValue
		}
	}

	return maxValue
}

func FindMinimumValue[T any](comparisonFunc func(T, T) int, inputValues ...T) T {
	if len(inputValues) == 0 {
		panic("Cannot find minimum value in an empty set")
	}

	minValue := inputValues[0]
	for _, currentValue := range inputValues[1:] {
		if comparisonFunc(currentValue, minValue) < 0 {
			minValue = currentValue
		}
	}

	return minValue
}

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

// MaxInt takes two integers as input and returns the larger value.
// If the values are equal, it returns the same value.
//
// Example:
//
//	a := 5
//	b := 7
//	max := MaxInt(a, b)
//	fmt.Println(max) // Output: 7
func MaxInt(a, b int) int {
	if a > b {
		return a
	}

	return b
}

// MinInt takes two integers as input and returns the smaller value.
// If the values are equal, it returns the same value.
//
// Example:
//
//	a := 5
//	b := 7
//	min := MinInt(a, b)
//	fmt.Println(min) // Output: 5
func MinInt(a, b int) int {
	if a < b {
		return a
	}

	return b
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

	if n <= 0 {
		panic("The number of groups must be positive and non-zero")
	}

	groups := make([][]T, n)

	for index, element := range slice {
		groupNumber := (index + 1) % n

		groups[groupNumber] = append(groups[groupNumber], element)
	}

	return groups
}
