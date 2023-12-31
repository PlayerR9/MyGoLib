// git tag v0.1.18

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
