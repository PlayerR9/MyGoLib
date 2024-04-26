package Errors

import (
	"errors"
	"io/fs"
	"log"
	"os"
)

const (
	// DefaultLFPermissions is the default file permissions for creating a log file.
	DefaultLFPermissions fs.FileMode = 0666

	// DefaultLFExtension is the default file extension for creating a log file.
	DefaultLFExtension string = ".log.md"

	// DefaultLFFlags is the default flags for creating a log file.
	DefaultLFFlags int = os.O_CREATE | os.O_WRONLY | os.O_APPEND

	// DefaultLoggerFlags is the default flags for creating a logger.
	DefaultLoggerFlags int = log.LstdFlags | log.Llongfile
)

// As is function that checks if an error is of type T.
//
// If the error is nil, the function returns false.
//
// Parameters:
//   - err: The error to check.
//
// Returns:
//   - bool: true if the error is of type T, false otherwise.
func As[T any](err error) bool {
	if err == nil {
		return false
	}

	var target T

	return errors.As(err, &target)
}
