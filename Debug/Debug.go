package Debug

import (
	"fmt"
	"log"
	"os"
	"runtime"

	sext "github.com/PlayerR9/MyGoLib/StrExt"
)

const (
	// Log levels
	DEBUG = iota
	INFO
	WARN
	ERROR
	FATAL
	TEST
	OTHER
)

// GetLogLevelPrefix is a function that is used to get the log level prefix. Panics if the log level is unknown.
//
// Parameters:
//   - log_level: The log level. It can be one of the following: DEBUG, INFO, WARN, ERROR, FATAL, OTHER.
//
// Returns:
//   - string: The log level prefix.
func GetLogLevelPrefix(log_level int) string {
	switch log_level {
	case DEBUG:
		return "[DEBUG]"
	case INFO:
		return "[INFO]"
	case WARN:
		return "[WARNING]"
	case ERROR:
		return "[ERROR]"
	case FATAL:
		return "[FATAL]"
	case TEST:
		return "[TEST]"
	case OTHER:
		return "[OTHER]"
	default:
		panic("Unknown log level.")
	}
}

// Debugger is a struct that is used to print debug information about the test. It can be active or not, and the information that it prints
// can be modified
type Debugger struct {
	is_active bool        // If true, the debugger will print debug information during the test.
	logger    *log.Logger // The logger used to print debug information. It is initialized with the default logger. If you want to use a custom logger, you can set it with the SetDebugMode function.
}

// NewDebugger is a function that is used to create a new debugger.
//
// Parameters:
//   - is_active: If true, the debugger will print debug information during the test. If false, the debugger will not print debug information during the test.
//   - logger: The logger used to print debug information.
//
// Returns:
//   - Debugger: The debugger.
//
// Information:
//   - If the logger is nil, it will be initialized with the default logger.
func NewDebugger(is_active bool, logger *log.Logger) Debugger {
	if logger == nil {
		return Debugger{is_active, log.New(os.Stdout, "", log.LstdFlags)}
	}

	return Debugger{is_active, logger}
}

// Activate is a function that activates the debugger
func (d *Debugger) Activate() {
	d.is_active = true
}

// Deactivate is a function that deactivates the debugger
func (d *Debugger) Deactivate() {
	d.is_active = false
}

// IsActive is a function that returns true if the debugger is active, false otherwise
func (d Debugger) IsActive() bool {
	return d.is_active
}

// SetLogger is a function that is used to set the logger used to print debug information.
//
// Parameters:
//   - logger: The logger used to print debug information.
//
// Information:
//   - If the logger is nil, it will be initialized with the default logger.
func (d *Debugger) SetLogger(logger *log.Logger) {
	if logger == nil {
		d.logger = log.New(os.Stdout, "", log.LstdFlags)
	} else {
		d.logger = logger
	}
}

// Println is a function that is used to print debug information about the test. It uses the logger Println() function.
//
// Parameters:
//   - log_level: The log level. It can be one of the following: INFO, WARNING, ERROR, FATAL, DEBUG, OTHER.
//   - v: The debug information.
//
// Information:
//   - If the debugger is not active, no debug information is printed.
//   - If log_level is FATAl, the test will be stopped.
func (d Debugger) Println(log_level int, v ...interface{}) {
	// No debug information is printed if the debugger is not active
	if !d.is_active {
		return
	}

	// Get information about the caller
	_, fname, line_counter, _ := runtime.Caller(1)

	if log_level < DEBUG || log_level > OTHER {
		panic("Unknown log level.")
	}

	// Print the debug header
	d.logger.Printf("%s: at the %s line of file %s:\n", GetLogLevelPrefix(log_level), sext.PrintOrdinalNumber(line_counter), fname)

	// Print the debug information
	writer := d.logger.Writer()

	fmt.Fprintln(writer, v...)

	// If the log level is FATAL, the test will be stopped
	if log_level == FATAL {
		d.logger.Fatalln("Test stopped.")
	}
}

// Print is a function that is used to print debug information about the test. It uses the logger Print() function.
//
// Parameters:
//   - log_level: The log level. It can be one of the following: INFO, WARNING, ERROR, FATAL, DEBUG, OTHER.
//   - v: The debug information.
//
// Information:
//   - If the debugger is not active, no debug information is printed.
//   - If log_level is FATAl, the test will be stopped.
func (d Debugger) Print(log_level int, v ...interface{}) {
	// No debug information is printed if the debugger is not active
	if !d.is_active {
		return
	}

	// Get information about the caller
	_, fname, line_counter, _ := runtime.Caller(1)

	if log_level < DEBUG || log_level > OTHER {
		panic("Unknown log level.")
	}

	// Print the debug header
	d.logger.Printf("%s: at the %s line of file %s:\n", GetLogLevelPrefix(log_level), sext.PrintOrdinalNumber(line_counter), fname)

	// Print the debug information
	writer := d.logger.Writer()

	fmt.Fprint(writer, v...)

	// If the log level is FATAL, the test will be stopped
	if log_level == FATAL {
		d.logger.Fatalln("Test stopped.")
	}
}

// Printf is a function that is used to print debug information about the test. It uses the logger Printf() function.
//
// Parameters:
//   - log_level: The log level. It can be one of the following: INFO, WARNING, ERROR, FATAL, DEBUG, OTHER.
//   - format: The format string.
//   - v: The debug information.
//
// Information:
//   - If the debugger is not active, no debug information is printed.
//   - If log_level is FATAl, the test will be stopped.
func (d Debugger) Printf(log_level int, format string, v ...interface{}) {
	// No debug information is printed if the debugger is not active
	if !d.is_active {
		return
	}

	// Get information about the caller
	_, fname, line_counter, _ := runtime.Caller(1)

	if log_level < DEBUG || log_level > OTHER {
		panic("Unknown log level.")
	}

	// Print the debug header
	d.logger.Printf("%s: at the %s line of file %s:\n", GetLogLevelPrefix(log_level), sext.PrintOrdinalNumber(line_counter), fname)

	// Print the debug information
	writer := d.logger.Writer()

	fmt.Fprintf(writer, format, v...)

	// If the log level is FATAL, the test will be stopped
	if log_level == FATAL {
		d.logger.Fatalln("Test stopped.")
	}
}

// TestCondition is a function that is used to test a condition for a set of elements.
//
// Parameters:
//   - debugger: The debugger used to print debug information.
//   - condition_to_test: The condition to test. It returns the log level and an error. If the error is not nil, the test failed.
//   - elements: The elements to test.
//   - printElement: A function that is used to print the element.
//
// Returns:
//   - bool: True if the test was successful, false otherwise.
//
// Information:
//   - If the debugger is not active, no debug information is printed.
//   - If log_level is FATAl, the test will be stopped.
func TestCondition[T any](debugger Debugger, condition_to_test func(T) (int, error), printElement func(T) string, elements ...T) bool {
	// No debug information is printed if the debugger is not active
	if !debugger.is_active {
		for _, e := range elements {
			_, err := condition_to_test(e)

			if err != nil {
				return false
			}
		}

		return true
	}

	// Otherwise

	// Initialize the variables
	sucessful_test := true

	writer := debugger.logger.Writer()

	// Get information about the caller
	_, fname, line_counter, _ := runtime.Caller(1)

	// Print the debug header
	debugger.logger.Printf("%s: at the %s line of file %s:\n", GetLogLevelPrefix(TEST), sext.PrintOrdinalNumber(line_counter), fname)

	for i, e := range elements {
		// Print the test information
		writer.Write([]byte(fmt.Sprintf("Testing the %s element: %s\n", sext.PrintOrdinalNumber(i+1), printElement(e))))

		// Test the condition
		log_level, err := condition_to_test(e)

		if log_level < DEBUG || log_level > OTHER {
			panic("Unknown log level.")
		}

		if err != nil && sucessful_test {
			sucessful_test = false
		}

		// Print the test result
		if err != nil {
			writer.Write([]byte(fmt.Sprintf("%s: The %s element does not satisfy the condition: %s\n", sext.PrintOrdinalNumber(i+1), GetLogLevelPrefix(log_level), err)))
		} else {
			writer.Write([]byte(fmt.Sprintf("%s: The %s element satisfies the condition.\n", sext.PrintOrdinalNumber(i+1), GetLogLevelPrefix(log_level))))
		}

		// If the log level is FATAL, the test will be stopped
		if log_level == FATAL {
			debugger.logger.Fatalln("Test stopped.")
		}
	}

	// Inform that the test is finished
	debugger.logger.Printf("%s: Test finished.\n", GetLogLevelPrefix(TEST))

	return sucessful_test
}
