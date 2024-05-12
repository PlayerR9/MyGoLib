package General

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	uc "github.com/PlayerR9/MyGoLib/Units/Common"
)

// RunInPowerShell is a function that returns a function that runs a program in
// a new PowerShell process.
//
// Upon calling the returned function, a new PowerShell process is started with
// the specified program and arguments. The function returns an error if the
// process cannot be started.
//
// Parameters:
//   - program: The path to the program to run.
//   - args: The arguments to pass to the program.
//
// Return:
//   - MainFunc: A function that runs the program in a new PowerShell process.
func RunInPowerShell(program string, args ...string) uc.MainFunc {
	var builder strings.Builder

	builder.WriteString("'-NoExit', '")
	builder.WriteString(program)
	builder.WriteString("'")

	for _, arg := range args {
		builder.WriteRune(',')
		builder.WriteRune(' ')
		builder.WriteRune('\'')
		builder.WriteString(arg)
		builder.WriteRune('\'')
	}

	cmd := exec.Command(
		"powershell", "-Command", "Start-Process", "powershell", "-ArgumentList",
		builder.String(),
	)

	return cmd.Run
}

// RecoverFromPanic is a function that recovers from a panic and logs the error
// message to a log file.
//
// If a panic occurs during the execution of the program, this function will
// recover from the panic, log the error message to the specified log file, and
// exit the program with the Panic exit code.
//
// Parameters:
//   - logger: The logger to use for logging the error message. If nil, the
//     error message is logged to the console.
func RecoverFromPanic(logger *log.Logger) {
	if r := recover(); r != nil {
		if logger != nil {
			logger.Printf("%s: %v\n", Panic.String(), r)

			fmt.Println("An unexpected error occurred. For more information, see the log")
		} else {
			fmt.Printf("%s: %v\n", Panic.String(), r)
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
//   - logger: The logger to use for logging the error message. If nil, the
//     error message is logged to the console.
//   - result: The result of the program. If nil, the program is considered to
//     have finished successfully.
//
// Return:
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
func FinalizeResult(logger *log.Logger, result error, isSetup bool) int {
	var errType ExitCode

	if isSetup {
		errType = SetupFailed
	} else {
		errType = Error
	}

	if result != nil {
		if logger != nil {
			logger.Printf("%s: %s\n", errType.String(), result.Error())

			if isSetup {
				fmt.Print("Could not set up the program.")
			} else {
				fmt.Print("An error occurred.")
			}

			fmt.Println(" For more information, see the log")
		} else {
			fmt.Printf("%s: %s\n", errType.String(), result.Error())
		}
	} else {
		fmt.Printf("%s\n", Success.String())
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
