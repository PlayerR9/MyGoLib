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
//   - Args: The argument name of the flag.
//   - Description: The description of the flag.
//   - Required: Whether or not the flag is required.
//   - Callback: A function that is called when the flag is parsed.
type ConsoleFlagInfo struct {
	// The name of the flag.
	Name string

	// The argument name of the flag.
	Args []string

	// The description of the flag.
	Description string

	// Whether or not the flag is required.
	Required bool

	// A function that is called when the flag is parsed.
	Callback func(args ...string) (interface{}, error)
}

func (cfi ConsoleFlagInfo) ToString() string {
	str := ""

	str += "Name: " + cfi.Name + "\n"
	str += "Args:"

	for _, arg := range cfi.Args {
		str += " " + arg
	}

	str += "\n"

	str += "Description: " + cfi.Description + "\n"

	str += "Required: "

	if cfi.Required {
		str += "Yes"
	} else {
		str += "No"
	}

	return str
}

func UsageToString(executable_name string, flags []ConsoleFlagInfo) string {
	str := ""

	str += "Usage: " + executable_name

	for _, flag := range flags {
		str += " "

		if !flag.Required {
			str += "["
		}

		str += flag.Name

		for i, arg := range flag.Args {
			if i != 0 {
				str += " "
			}

			str += "<" + arg + ">"
		}

		if !flag.Required {
			str += "]"
		}
	}

	return str
}

func HelpToString(executable_name string, flags []ConsoleFlagInfo) string {
	str := UsageToString(executable_name, flags) + "\n\n" + "Flags:\n"

	for i, flag := range flags {
		if i != 0 {
			str += "\n\n"
		}

		str += flag.ToString()
	}

	return str
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
