// git tag v0.1.13

package General

import (
	"fmt"
	"log"
	"os"
)

// GLOBAL VARIABLES
var (
	// DebugMode is a boolean that is used to enable or disable debug mode. When debug mode is enabled, the package will print debug messages.
	// **Note:** Debug mode is disabled by default.
	DebugMode bool = false

	debugger *log.Logger = log.New(os.Stdout, "[General] ", log.LstdFlags) // Debugger
)

// ConsoleFlagInfo is a struct that contains information about a console flag.
//
// Fields:
//   - Name: The name of the flag.
//   - NumArgs: Number of arguments the flag takes.
//   - Description: The description of the flag.
//   - Required: Whether or not the flag is required.
//   - Callback: A function that is called when the flag is parsed.
type ConsoleFlagInfo struct {
	// The name of the flag.
	Name string

	// Number of arguments the flag takes.
	NumArgs int

	// The description of the flag.
	Description string

	// Whether or not the flag is required.
	Required bool

	// A function that is called when the flag is parsed.
	Callback func(args ...string) (interface{}, error)
}

// PressEnterToContinue prints "Press enter to continue..." to the console and waits for the user to press enter.
func PressEnterToContinue() {
	fmt.Println("Press enter to continue...") // Print "Press enter to continue..." to the console
	fmt.Scanln()                              // Wait for the user to press enter
}

// PressEnterToExit prints "Press enter to exit..." to the console and waits for the user to press enter. Then, it exits the program with the given exit code.
//
// Parameters:
//   - exit_code: The exit code to exit the program with.
func PressEnterToExit(exit_code int) {
	fmt.Println("Press enter to exit...") // Print "Press enter to exit..." to the console
	fmt.Scanln()                          // Wait for the user to press enter

	os.Exit(exit_code) // Exit the program with the given exit code
}

// MaxFunc returns the index of the maximum element in a slice of values. Panics if the slice is empty.
// For example, MaxFunc(func(a, b int) int { return a - b }, 1, 2, 3, 4) returns 3.
//
// Parameters:
//   - f: The function to use to compare the values. The function should return < 0 if a < b, 0 if a == b, > 0 if a > b.
//   - values: The slice of values to get the maximum of.
//
// Returns:
//   - int: The index of the first occurrence of the maximum element in the slice of values.
func MaxFunc[T any](f func(T, T) int, values ...T) int {
	if len(values) == 0 {
		// Cannot get max of no values, so panic.
		if DebugMode {
			debugger.Panic("Cannot get max of no values")
		} else {
			panic("Cannot get max of no values")
		}
	}

	max := 0 // The index of the maximum element

	// Find the index of the maximum element
	for index, value := range values[1:] {
		if f(value, values[max]) > 0 {
			max = index
		}
	}

	return max
}

// MinFunc returns the index of the minimum element in a slice of values. Panics if the slice is empty.
// For example, MinFunc(func(a, b int) int { return a - b }, 1, 2, 3, 4) returns 0.
//
// Parameters:
//   - f: The function to use to compare the values. The function should return < 0 if a < b, 0 if a == b, > 0 if a > b.
//   - values: The slice of values to get the minimum of.
//
// Returns:
//   - int: The index of the first occurrence of the minimum element in the slice of values.
func MinFunc[T any](f func(T, T) int, values ...T) int {
	if len(values) == 0 {
		// Cannot get min of no values, so panic.
		if DebugMode {
			debugger.Panic("Cannot get min of no values")
		} else {
			panic("Cannot get min of no values")
		}
	}

	min := 0 // The index of the minimum element

	// Find the index of the minimum element
	for index, value := range values[1:] {
		if f(value, values[min]) < 0 {
			min = index
		}
	}

	return min
}
