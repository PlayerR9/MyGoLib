package General

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"strings"

	ers "github.com/PlayerR9/MyGoLib/Utility/Errors"
)

// Comparable is an interface that defines the behavior of a type that can be
// compared with other values of the same type using the < and > operators.
// The interface is implemented by the built-in types int, int8, int16, int32,
// int64, uint, uint8, uint16, uint32, uint64, float32, float64, and string.
type Comparable interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64 | ~string
}

// Min is a function that takes two parameters, a and b, of any type T that
// implements the Comparable interface and returns the smaller of the two values.
//
// Parameters:
//
//   - a, b: The two values to compare.
//
// Return:
//
//   - T: The smaller of the two values.
func Min[T Comparable](a, b T) T {
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
//
//   - a, b: The two values to compare.
//
// Return:
//
//   - T: The larger of the two values.
func Max[T Comparable](a, b T) T {
	if a < b {
		return b
	} else {
		return a
	}
}

// DeepCopy is a function that performs a deep copy of a given value.
//
// Parameters:
//
//   - value: The value to copy.
//
// Return:
//
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
// Panics with error of type *ers.ErrInvalidParameter if the number of groups
// is less than or equal to 0.
//
// Parameters:
//
//   - slice: The slice to split.
//   - n: The number of groups to split the slice into.
//
// Return:
//
//   - [][]T: A 2D slice where each inner slice represents a group.
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
		return nil, ers.NewErrInvalidParameter("n").
			Wrap(fmt.Errorf("negative group number (%d) are not allowed", n))
	} else if n == 0 {
		return nil, ers.NewErrInvalidParameter("n").
			Wrap(errors.New("cannot split into 0 groups"))
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
//
//   - value: The value to check.
//
// Return:
//
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

// RunInPowerShell is a function that returns a function that runs a program in
// a new PowerShell process.
//
// Upon calling the returned function, a new PowerShell process is started with
// the specified program and arguments. The function returns an error if the
// process cannot be started.
//
// Parameters:
//
//   - program: The path to the program to run.
//   - args: The arguments to pass to the program.
//
// Return:
//
//   - func() error: A function that runs the program in a new PowerShell process.
func RunInPowerShell(program string, args ...string) func() error {
	var builder strings.Builder

	fmt.Fprintf(&builder, "'-NoExit', '%s'", program)

	for _, arg := range args {
		fmt.Fprintf(&builder, ", '%s'", arg)
	}

	cmd := exec.Command(
		"powershell", "-Command", "Start-Process", "powershell", "-ArgumentList",
		builder.String(),
	)

	return cmd.Run
}

// ExitCode is a custom type that represents the exit code of a program.
type ExitCode int

const (
	// Success indicates that the program has finished successfully.
	Success ExitCode = iota

	// Panic indicates that a panic occurred during the execution of the program.
	Panic

	// SetupFailed indicates that the program could not be set up.
	SetupFailed

	// Error indicates that an error occurred during the execution of the program.
	Error
)

// String returns a string representation of the exit code.
//
// Return:
//
//   - string: A string representation of the exit code.
func (ec ExitCode) String() string {
	return [...]string{
		"Program has finished successfully",
		"Panic occurred",
		"Cound not set up the program",
		"An error occurred",
	}[ec]
}

// RecoverFromPanic is a function that recovers from a panic and logs the error
// message to a log file.
//
// If a panic occurs during the execution of the program, this function will
// recover from the panic, log the error message to the specified log file, and
// exit the program with the Panic exit code.
//
// Parameters:
//
//   - logger: The logger to use for logging the error message. If nil, the
//     error message is logged to the console.
func RecoverFromPanic(logger *ers.FileLogger) {
	if r := recover(); r != nil {
		if logger != nil {
			logger.Printf("%v: %v\n", Panic, r)

			fmt.Println("An unexpected error occurred. For more information, see the log file:", logger.GetFileName())
		} else {
			fmt.Printf("%v: %v\n", Panic, r)
		}
		fmt.Println()

		fmt.Println("Press the enter key to exit...")
		fmt.Scanln()
	}

	os.Exit(int(Panic))
}

// FinalizeResult is a function that finalizes the result of a program and logs
// the error message to a log file.
//
// If an error occurs during the execution of the program, this function will
// log the error message to the specified log file and exit the program with the
// Error exit code. If no error occurs, the function will log a success message
// to the log file and exit the program with the Success exit code.
//
// Parameters:
//
//   - logger: The logger to use for logging the error message. If nil, the
//     error message is logged to the console.
//   - result: The result of the program. If nil, the program is considered to
//     have finished successfully.
//
// Return:
//
//   - int: The exit code of the program.
//
// Example:
//
// var logger *ers.FileLogger
//
//	func main() {
//		// Set up the logger
//		var err error
//		logger, err = ers.NewFileLogger("log.txt")
//		if err != nil {
//			os.Exit(FinalizeResult(logger, err, true))
//		}
//		defer logger.Close()
//
//		defer RecoverFromPanic(logger) // handle panics gracefully
//
//		os.Exit(FinalizeResult(logger, mainBody(), false)) // handle errors gracefully
//	}
//
//	func mainBody() error {
//		// Perform the main logic of the program
//	}
func FinalizeResult(logger *ers.FileLogger, result error, isSetup bool) int {
	var errType ExitCode

	if isSetup {
		errType = SetupFailed
	} else {
		errType = Error
	}

	if result != nil {
		if logger != nil {
			logger.Printf("%v: %v\n", errType, result)

			if isSetup {
				fmt.Print("Could not set up the program.")
			} else {
				fmt.Print("An error occurred.")
			}

			fmt.Println(" For more information, see the log file: ", logger.GetFileName())
		} else {
			fmt.Printf("%v: %v\n", errType, result)
		}
	} else {
		fmt.Printf("%v\n", Success)
	}

	fmt.Println()

	fmt.Println("Press the enter key to exit...")
	fmt.Scanln()

	if result != nil {
		if isSetup {
			return int(SetupFailed)
		} else {
			return int(Error)
		}
	} else {
		return int(Success)
	}
}
