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

// Unwrapper is an interface that defines a method to unwrap an error.
type Unwrapper interface {
	// Unwrap returns the error that this error wraps.
	//
	// Returns:
	//   - error: The error that this error wraps.
	Unwrap() error

	// ChangeReason changes the reason of the error.
	//
	// Parameters:
	//   - reason: The new reason of the error.
	ChangeReason(reason error)
}

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

// IsNoError checks if an error is a no error error or if it is nil.
//
// Parameters:
//   - err: The error to check.
//
// Returns:
//   - bool: True if the error is a no error error or if it is nil, otherwise false.
func IsNoError(err error) bool {
	if err == nil {
		return true
	}

	var errNoError *ErrNoError

	return errors.As(err, &errNoError)
}

// IsErrIgnorable checks if an error is an *ErrIgnorable or *ErrInvalidParameter error.
// If the error is nil, the function returns false.
//
// Parameters:
//   - err: The error to check.
//
// Returns:
//   - bool: True if the error is an *ErrIgnorable or *ErrInvalidParameter error,
//     otherwise false.
func IsErrIgnorable(err error) bool {
	if err == nil {
		return false
	}

	var ignorable *ErrIgnorable

	if errors.As(err, &ignorable) {
		return true
	}

	var invalid *ErrInvalidParameter

	return errors.As(err, &invalid)
}

// LimitErrorMsg limits the error message to a certain number of unwraps.
// It returns the top level error for allowing to print the error message
// with the limit of unwraps applied.
//
// If the error is nil or the limit is less than 0, the function does nothing.
//
// Parameters:
//   - err: The error to limit.
//   - limit: The limit of unwraps.
//
// Returns:
//   - error: The top level error with the limit of unwraps applied.
func LimitErrorMsg(err error, limit int) error {
	if err == nil || limit < 0 {
		return err
	}

	currErr := err

	for i := 0; i < limit; i++ {
		target, ok := currErr.(Unwrapper)
		if !ok {
			return err
		}

		reason := target.Unwrap()
		if reason == nil {
			return err
		}

		currErr = reason
	}

	// Limit reached
	target, ok := currErr.(Unwrapper)
	if !ok {
		return err
	}

	target.ChangeReason(nil)

	return err
}
