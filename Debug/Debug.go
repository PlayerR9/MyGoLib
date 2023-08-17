package Debug




Your idea of using a map of custom log levels with associated information is a good approach to allow users to define and manage their own log levels. It provides
flexibility and control over log level behavior. However, managing this map and providing efficient access to log levels can be enhanced further.

Here's an implementation using your suggested approach with a few modifications to streamline the management of custom log levels:

```go
package mylogger

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"sync"
)

type LogLevel int

type LogLevelBehavior struct {
	Prefix string
	Fatal  bool
}

var (
	defaultLogLevels = map[LogLevel]LogLevelBehavior{
		DEBUG: {Prefix: "DEBUG:", Fatal: false},
		WARN:  {Prefix: "WARN:", Fatal: false},
		// Add more default log levels as needed
	}

	customLogLevels = make(map[int]LogLevelBehavior)
	logLevelMutex   sync.RWMutex
)

type CustomLogger struct {
	logger *log.Logger
	level  LogLevel
}

type dynamicPrefixWriter struct {
	writer io.Writer
	prefix string
}

func (w dynamicPrefixWriter) Write(p []byte) (n int, err error) {
	prefixedMessage := fmt.Sprintf("%s %s", w.prefix, p)
	return w.writer.Write([]byte(prefixedMessage))
}

func NewCustomLogger(output io.Writer, level LogLevel) *CustomLogger {
	prefix := getLogLevelPrefix(level)
	logger := log.New(&dynamicPrefixWriter{writer: output, prefix: prefix}, "", log.Lshortfile)

	return &CustomLogger{
		logger: logger,
		level:  level,
	}
}

func getLogLevelPrefix(level LogLevel) string {
	logLevelMutex.RLock()
	defer logLevelMutex.RUnlock()

	behavior, exists := defaultLogLevels[level]
	if exists {
		return behavior.Prefix
	}
	customBehavior, customExists := customLogLevels[int(level)]
	if customExists {
		return customBehavior.Prefix
	}
	return "INFO:" // Default prefix for unknown levels
}

func (c *CustomLogger) logMessage(level LogLevel, message string) {
	logLevelMutex.RLock()
	defer logLevelMutex.RUnlock()

	behavior, exists := defaultLogLevels[level]
	if exists {
		if behavior.Fatal {
			log.Fatalf("%s %s", behavior.Prefix, message)
		} else {
			c.logger.Output(2, message)
		}
	} else {
		customBehavior, customExists := customLogLevels[int(level)]
		if customExists {
			if customBehavior.Fatal {
				log.Fatalf("%s %s", customBehavior.Prefix, message)
			} else {
				c.logger.Output(2, message)
			}
		} else {
			c.logger.Output(2, message)
		}
	}
}

func (c *CustomLogger) Log(level LogLevel, message string) {
	if c.level <= level {
		c.logMessage(level, message)
	}
}

// ... Add methods to manage custom log levels like AddLogLevel, GetLogLevel, DeleteLogLevel
```

In this implementation:

1. The default log levels are stored in `defaultLogLevels` map.
2. Custom log levels are stored in the `customLogLevels` map, which is guarded by a `sync.RWMutex` to ensure thread-safe access.
3. The `getLogLevelPrefix` function retrieves prefixes based on the provided log level. It first checks if the level exists in the default or custom log levels.
4. The `logMessage` function handles logging messages with the appropriate prefix and behavior based on the log level.

You can then implement methods to manage custom log levels such as `AddLogLevel`, `GetLogLevel`, and `DeleteLogLevel` to allow users to customize their logging behavior.

This approach reduces the complexity of checking and managing log levels, allowing for efficient customization while still utilizing the default log levels when needed.



type LogLevel int

const (
	DEBUG LogLevel = iota
	WARN
	// Add more log levels as needed
)

type CustomLogger struct {
	logger *log.Logger
	level  LogLevel
}

type dynamicPrefixWriter struct {
	writer io.Writer
	prefix string
}

func (w dynamicPrefixWriter) Write(p []byte) (n int, err error) {
	prefixedMessage := fmt.Sprintf("%s %s", w.prefix, p)
	return w.writer.Write([]byte(prefixedMessage))
}

func NewCustomLogger(output io.Writer, level LogLevel) *CustomLogger {
	prefix := getLogLevelPrefix(level)
	logger := log.New(&dynamicPrefixWriter{writer: output, prefix: prefix}, "", log.Lshortfile)

	return &CustomLogger{
		logger: logger,
		level:  level,
	}
}

func getLogLevelPrefix(level LogLevel) string {
	switch level {
	case DEBUG:
		return "DEBUG:"
	case WARN:
		return "WARN:"
	// Add more cases for other log levels
	default:
		return "INFO:" // Default prefix for unknown levels
	}
}

func (c *CustomLogger) Debug(message string) {
	if c.level <= DEBUG {
		c.logger.Output(2, message)
	}
}

func (c *CustomLogger) Warn(message string) {
	if c.level <= WARN {
		c.logger.Output(2, message)
	}
}

// Add more log methods for other levels as needed
```

In this example, the `dynamicPrefixWriter` type wraps an `io.Writer` and dynamically adds the appropriate prefix to the log entries. The `NewCustomLogger()` function constructs a custom logger using the `dynamicPrefixWriter`. The `getLogLevelPrefix()` function translates the log level into the corresponding prefix.

With this approach, you'll be able to achieve dynamic prefixes based on log levels without creating a separate logger for each level, thus avoiding excessive memory usage.

Here's how you would use the custom logger in your application:

```go
package main

import (
	"yourpackage/mylogger"
	"os"
)

func main() {
	logger := mylogger.NewCustomLogger(os.Stdout, mylogger.DEBUG)
	logger.Debug("This is a debug message.")
	logger.Warn("This is a warning message.")
}
```

Remember to adjust the `getLogLevelPrefix()` function to include additional log levels as needed.




package Debug

import "log"

// The log levels
const (
	// The default log level. It is used to print information about the test.
	INFO = iota

	// It is used to print warnings about the test. For example, if the test is not running as expected.
	WARNING

	// It is used to print errors about the test. For example, if the test is not running as expected.
	ERROR

	// It is used to print fatal errors about the test. For example, if the test is not running as expected. If a fatal error is printed, the test will be stopped.
	FATAL

	// It is used to print debug information about the test. When using this log level, the debugger will print debug information about the test regardless of the debug mode.
	DEBUG

	// A generic log level. It is used to print information about the test that is not related to the other log levels.
	OTHER
)

// The default logger
var defaultLogger log.Logger = *log.New(log.Writer(), "DEBUG: ", log.Lshortfile)

func print_log_level(log_level int) string {
	switch log_level {
	case INFO:
		return "[INFO]: "
	case WARNING:
		return "[WARNING]: "
	case ERROR:
		return "[ERROR]: "
	case FATAL:
		return "[FATAL]: "
	case DEBUG:
		return "[DEBUG]: "
	}

	return "[OTHER]: "
}

// Debugger is a struct that is used to print debug information about the test.
// It has a mode that can be enabled or disabled. If the mode is enabled, the debugger will print debug information about the test. If the mode is disabled,
// the debugger will not print debug information about the test. The mode can be enabled or disabled with the SetDebugMode and UnsetDebugMode functions.
// The debugger has a logger that is used to print debug information. It is initialized with the default logger. If you want to use a custom logger, you can set it
// with the SetDebugMode function. The debugger has a Println function that is used to print debug information. It has a log level that is used to print debug information
// about the test. The log level can be one of the following: INFO, WARNING, ERROR, FATAL, DEBUG, OTHER. If the log level is DEBUG, the debugger will print debug information
// about the test regardless of the debug mode. If the log level is INFO, WARNING, ERROR, FATAL, OTHER, the debugger will print debug information about the test only if the debug
// mode is enabled. If the log level is FATAL, the debugger will print debug information about the test only if the debug mode is enabled and the test will be stopped.
type Debugger struct {
	mode   bool        // If true, the debugger will print debug information during the test.
	logger *log.Logger // The logger used to print debug information. It is initialized with the default logger. If you want to use a custom logger, you can set it with the SetDebugMode function.
}

// NewDebugger is a function that is used to create a new debugger.
//
// Parameters:
//  - logger: The logger used to print debug information.
//
// Returns:
//  - Debugger: The debugger.
//
// Information:
//  - If the logger is nil, it will be initialized with the default logger.
//  - The debugger has a mode that can be enabled or disabled. By default, the mode is disabled. It can be enabled or disabled with the SetDebugMode and UnsetDebugMode functions.
func NewDebugger(logger *log.Logger) Debugger {
	if logger == nil {
		logger = &defaultLogger
	}

	return Debugger{false, logger}
}

// SetDebugMode is a function that is used to enable the debug mode. It will print a message to inform the user that the debug mode is being enabled
// or that the debug mode is already enabled.
func (d *Debugger) SetDebugMode() {
	if d.mode {
		d.logger.Println("Debug mode is already enabled.")
	} else {
		d.logger.Println("Debug mode is being enabled.")
	}

	d.mode = true
}

// UnsetDebugMode is a function that is used to disable the debug mode. It will print a message to inform the user that the debug mode is being disabled
// or that the debug mode is already disabled.
func (d *Debugger) UnsetDebugMode() {
	if !d.mode {
		d.logger.Println("Debug mode is already disabled.")
	} else {
		d.logger.Println("Debug mode is being disabled.")
	}

	d.mode = false
}

// GetDebugMode is a function that is used to get the debug mode.
//
// Returns:
//  - bool: The debug mode.
func (d Debugger) GetDebugMode() bool {
	return d.mode
}

// SetLogger is a function that is used to set the logger used to print debug information.
//
// Parameters:
//  - logger: The logger used to print debug information.
//
// Information:
//  - If the logger is nil, it will be initialized with the default logger.
func (d *Debugger) SetLogger(logger *log.Logger) {
	if logger == nil {
		logger = &defaultLogger
	}

	d.logger = logger
}

// GetLogger is a function that is used to get the logger used to print debug information.
//
// Returns:
//  - *log.Logger: The logger used to print debug information.
func (d Debugger) GetLogger() *log.Logger {
	return d.logger
}

// Println is a function that is used to print debug information about the test.
//
// Parameters:
//  - log_level: The log level. It can be one of the following: INFO, WARNING, ERROR, FATAL, DEBUG, OTHER.
//  - v: The debug information.
//
// Information:
//  - If the log level is DEBUG, the debugger will print debug information about the test regardless of the debug mode. If the log level is INFO, WARNING, ERROR, FATAL, OTHER,
//    the debugger will print debug information about the test only if the debug mode is enabled. If the log level is FATAL, the debugger will print debug information about the
//    test only if the debug mode is enabled and the test will be stopped.
func (d Debugger) Println(log_level int, v ...interface{}) {
	if log_level == DEBUG {
		d.logger.Print("[DEBUG]: ")
		d.logger.Println(v...)
		return
	}

	if d.mode {
		d.logger.Print(print_log_level(log_level))

		if log_level == FATAL {
			d.logger.Fatalln(v...)
		} else {
			d.logger.Println(v...)
		}
	}
}

// TestCondition is a function that is used to test if the elements satisfy the condition. It will print debug information about the test.
//
// Parameters:
//  - debugger: The debugger used to print debug information.
//  - condition: The condition to test.
//  - elements: The elements to test.
//
// Information:
//  - If the condition returns an error, the debugger will print debug information about the error.
//  - If the condition returns FATAL, the debugger will print debug information about the error and the test will be stopped.
//  - If the condition returns INFO, WARNING, ERROR, OTHER, the debugger will print debug information about the error only if the debug mode is enabled.
func TestCondition[T any](debugger Debugger, condition func(T) (int, error), elements ...T) {
	for i, e := range elements {
		code, err := condition(e)

		debugger.logger.Print(print_log_level(code))

		if err != nil {
			debugger.logger.Println("Element", i, "does not satisfy the condition: ", err)
		} else {
			debugger.logger.Println("Element", i, "satisfies the condition.")
		}

		if code == FATAL {
			debugger.logger.Fatalln("Test stopped.")
		}
	}

	debugger.logger.Println("Test finished.")
}
