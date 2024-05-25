package Printer

import "strings"

type ErrFinalization struct {
	// Reason is the reason for the error.
	Reason error
}

// Error returns the error message: "could not finalize the last page: <reason>".
//
// Returns:
//   - string: The error message.
//
// Behaviors:
//   - If the reason is nil, the error message is: "last page could not be finalized" instead.
func (e *ErrFinalization) Error() string {
	if e.Reason == nil {
		return "last page could not be finalized"
	} else {
		var builder strings.Builder

		builder.WriteString("could not finalize the last page: ")
		builder.WriteString(e.Reason.Error())

		return builder.String()
	}
}

// Unwrap returns the reason for the error.
//
// Returns:
//   - error: The reason for the error.
func (e *ErrFinalization) Unwrap() error {
	return e.Reason
}

// NewErrFinalization creates a new ErrFinalization with the provided reason.
//
// Parameters:
//   - reason: The reason for the error.
//
// Returns:
//   - *ErrFinalization: The new error.
func NewErrFinalization(reason error) *ErrFinalization {
	return &ErrFinalization{
		Reason: reason,
	}
}

// ErrInvalidPage is an error that occurs when a page is invalid.
type ErrInvalidPage struct {
	// Reason is the reason for the error.
	Reason error
}

// Error returns the error message: "invalid page: <reason>".
//
// Returns:
//   - string: The error message.
//
// Behaviors:
//   - If the reason is nil, the error message is: "page is invalid" instead.
func (e *ErrInvalidPage) Error() string {
	if e.Reason == nil {
		return "page is invalid"
	} else {
		var builder strings.Builder

		builder.WriteString("invalid page: ")
		builder.WriteString(e.Reason.Error())

		return builder.String()
	}
}

// Unwrap returns the reason for the error.
//
// Returns:
//   - error: The reason for the error.
func (e *ErrInvalidPage) Unwrap() error {
	return e.Reason
}

// NewErrInvalidPage creates a new ErrInvalidPage with the provided reason.
//
// Parameters:
//   - reason: The reason for the error.
//
// Returns:
//   - *ErrInvalidPage: The new error.
func NewErrInvalidPage(reason error) *ErrInvalidPage {
	return &ErrInvalidPage{
		Reason: reason,
	}
}
