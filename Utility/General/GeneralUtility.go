// git tag v0.1.18

package General

import (
	"fmt"
	"log"
	"os"
	"strings"
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

// ToString returns a string representation of the ConsoleFlagInfo struct.
//
// Returns:
//   - string: A string representation of the ConsoleFlagInfo struct.
func (cfi ConsoleFlagInfo) ToString(indent_level int) (str string) {
	var indentation string

	for i := 0; i < indent_level; i++ {
		indentation += "\t"
	}

	str += fmt.Sprintf("%sName: %s\n%sArgs:", indentation, cfi.Name, indentation)

	if len(cfi.Args) == 0 {
		str += " None"
	} else {
		for _, arg := range cfi.Args {
			str += " <" + arg + ">"
		}
	}

	str += fmt.Sprintf("\n%sDescription: %s", indentation, cfi.Description)

	if len(cfi.Args) != 0 {
		str += fmt.Sprintf("\n%sRequired: ", indentation)

		if cfi.Required {
			str += "Yes"
		} else {
			str += "No"
		}
	}

	return
}

// UsageToString returns a string representation of the usage of a command.
//
// Parameters:
//   - executable_name: The name of the executable.
//   - command: The name of the command.
//   - flags: The flags of the command.
//
// Returns:
//   - string: A string representation of the usage of a command.
func UsageToString(executable_name, command string, flags []ConsoleFlagInfo, indent_level int) (str string) {
	var indentation string

	for i := 0; i < indent_level; i++ {
		indentation += "\t"
	}

	str += fmt.Sprintf("%sUsage: %s %s", indentation, executable_name, command)

	for _, flag := range flags {
		if !flag.Required {
			str += " ["
		} else {
			str += " "
		}

		str += flag.Name

		for _, arg := range flag.Args {
			str += fmt.Sprintf(" <%s>", arg)
		}

		if !flag.Required {
			str += "]"
		}
	}

	return str
}

// HelpToString returns a string representation of the help of a command.
//
// Parameters:
//   - executable_name: The name of the executable.
//   - flags: The flags of the command.
//
// Returns:
//   - string: A string representation of the help of a command.
func HelpToString(executable_name string, flags map[string][]ConsoleFlagInfo) (str string) {
	for command, flag_set := range flags {
		str += UsageToString(executable_name, command, flag_set, 0)

		if len(flag_set) == 0 {
			str += "\n\n"
			continue
		}

		str += fmt.Sprintf("\nFlags:\n%s", flag_set[0].ToString(1))

		for _, flag := range flag_set[1:] {
			str += fmt.Sprintf("\n%s", flag.ToString(1))
		}

		str += "\n\n"
	}

	return strings.TrimSuffix(str, "\n\n")
}

// ParseConsoleFlags parses the flags of a command.
//
// Parameters:
//   - args: The arguments of the command.
//   - flags: The flags of the command.
//
// Returns:
//   - string: The name of the command.
//   - map[string]interface{}: A map of the flags and their values.
//   - error: An error if one occurred.
func ParseConsoleFlags(args []string, flags map[string][]ConsoleFlagInfo) (string, map[string]interface{}, error) {
	// Check if the command is present
	if len(args) == 0 {
		return "", nil, fmt.Errorf("no command specified")
	}

	command := args[0]

	if _, ok := flags[command]; !ok {
		return "", nil, fmt.Errorf("command %s not found", command)
	}

	// Parse flags
	flas_set := flags[command]

	results := make(map[string]interface{})

	// Check if enough arguments are present
	var min int = 1
	var max int = 1

	for _, f := range flas_set {
		if f.Required {
			min += len(f.Args) + 1
		}

		max += len(f.Args) + 1
	}

	if len(args) < min {
		return "", nil, fmt.Errorf("not enough arguments for command %s; expected at least %d, got %d", command, min, len(args))
	} else if len(args) >= max {
		return "", nil, fmt.Errorf("too many arguments for command %s; expected at most %d, got %d", command, max, len(args))
	}

	// Parse flags
	arg_index := 1

	for _, f := range flas_set {
		if arg_index >= len(args) {
			break
		}

		if f.Name != args[arg_index] {
			if f.Required {
				return "", nil, fmt.Errorf("required flag %s not present for command %s", f.Name, command)
			}

			continue
		}

		if len(f.Args)+arg_index >= len(args) {
			return "", nil, fmt.Errorf("flag %s present but not enough arguments specified for command %s", f.Name, command)
		}

		arg_index++

		args_tmp := make([]string, 0)

		for i := 0; i < len(f.Args); i++ {
			args_tmp = append(args_tmp, args[arg_index+i])
		}

		// Call callback function for flag
		inf_tmp, err := f.Callback(args_tmp...)
		if err != nil {
			return "", nil, fmt.Errorf("invalid argument for flag %s of command %s: %v", f.Name, command, err)
		}

		// Set result
		results[f.Name] = inf_tmp

		arg_index += len(f.Args)
	}

	return command, results, nil
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
