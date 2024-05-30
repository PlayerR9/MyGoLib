// Package errors provides a custom error type for out-of-bound errors.
package Errors

import (
	"fmt"
	"reflect"
	"strings"

	uc "github.com/PlayerR9/MyGoLib/Units/Common"
)

// ErrNoError represents an error when no error occurs.
type ErrNoError struct {
	// Err is the reason for the no error error.
	Err error
}

// Error is a method of the error interface.
//
// Returns:
//   - string: The error message of the no error error
//     (no mention of being a no error error).
func (e *ErrNoError) Error() string {
	if e.Err == nil {
		return "no error"
	} else {
		return e.Err.Error()
	}
}

// Unwrap is a method of the errors interface.
//
// Returns:
//   - error: The reason for the no error error.
func (e *ErrNoError) Unwrap() error {
	return e.Err
}

// ChangeReason changes the reason of the no error error.
//
// Parameters:
//   - reason: The new reason for the no error error.
func (e *ErrNoError) ChangeReason(reason error) {
	e.Err = reason
}

// NewErrNoError creates a new ErrNoError error.
//
// Parameters:
//   - err: The reason for the no error error.
//
// Returns:
//   - *ErrNoError: A pointer to the newly created ErrNoError.
func NewErrNoError(err error) *ErrNoError {
	return &ErrNoError{Err: err}
}

// ErrInvalidParameter represents an error when a parameter is invalid.
type ErrInvalidParameter struct {
	// Parameter is the name of the Parameter.
	Parameter string

	// Reason is the Reason for the invalidity of the parameter.
	Reason error
}

// Error is a method of the error interface.
//
// If the reason is not provided (nil), no reason is included in the error message.
//
// Returns:
//   - string: The error message.
func (e *ErrInvalidParameter) Error() string {
	if e.Reason == nil {
		return fmt.Sprintf("parameter (%s) is invalid", e.Parameter)
	} else {
		return fmt.Sprintf("parameter (%s) is invalid: %s", e.Parameter, e.Reason.Error())
	}
}

// Unwrap returns the reason for the invalidity of the parameter.
// It is used for error unwrapping.
//
// Returns:
//   - error: The reason for the invalidity of the parameter.
func (e *ErrInvalidParameter) Unwrap() error {
	return e.Reason
}

// ChangeReason changes the reason for the invalidity of the parameter.
//
// Parameters:
//   - reason: The new reason for the invalidity of the parameter.
func (e *ErrInvalidParameter) ChangeReason(reason error) {
	e.Reason = reason
}

// NewErrInvalidParameter creates a new ErrInvalidParameter error.
// If the reason is not provided (nil), the reason is set to
// "parameter is invalid" by default.
//
// Parameters:
//   - parameter: The name of the parameter.
//   - reason: The reason for the invalidity.
//
// Returns:
//   - error: A pointer to the newly created ErrInvalidParameter.
func NewErrInvalidParameter(parameter string, reason error) *ErrInvalidParameter {
	return &ErrInvalidParameter{
		Parameter: parameter,
		Reason:    reason,
	}
}

// ErrNilParameter represents an error when a parameter is nil.
// This is a shorthand for NewErrInvalidParameter(parameter, NewErrNilValue()).
//
// Parameters:
//   - parameter: The name of the parameter.
//
// Returns:
//   - *ErrInvalidParameter: A pointer to the newly created ErrInvalidParameter.
func NewErrNilParameter(parameter string) *ErrInvalidParameter {
	return &ErrInvalidParameter{
		Parameter: parameter,
		Reason:    NewErrNilValue(),
	}
}

// ErrInvalidCall represents an error that occurs when a function
// is not called correctly.
type ErrInvalidCall struct {
	// FnName is the name of the function.
	FnName string

	// Signature is the Signature of the function.
	Signature reflect.Type

	// Reason is the Reason for the failure.
	Reason error
}

// Error is a method of the error interface.
//
// Returns:
//
//   - string: The error message.
func (e *ErrInvalidCall) Error() string {
	var builder strings.Builder

	builder.WriteString("call")
	builder.WriteRune(' ')
	builder.WriteString("to")
	builder.WriteRune(' ')
	builder.WriteString(e.FnName)
	builder.WriteString(e.Signature.String())
	builder.WriteRune(' ')
	builder.WriteString("failed")

	if e.Reason != nil {
		builder.WriteRune(':')
		builder.WriteRune(' ')
		builder.WriteString(e.Reason.Error())
	}

	return builder.String()
}

// Unwrap returns the underlying error that caused the ErrInvalidCall.
// It is used for error unwrapping.
//
// Returns:
//   - error: The reason for the failure.
func (e *ErrInvalidCall) Unwrap() error {
	return e.Reason
}

// ChangeReason changes the reason for the failure of the function call.
//
// Parameters:
//   - reason: The new reason for the failure.
func (e *ErrInvalidCall) ChangeReason(reason error) {
	e.Reason = reason
}

// NewErrInvalidCall creates a new ErrInvalidCall. If the reason is not provided (nil),
// the reason is set to "an error occurred while calling the function" by default.
//
// Parameters:
//   - functionName: The name of the function.
//   - function: The function that failed.
//   - reason: The reason for the failure.
//
// Returns:
//   - *ErrInvalidCall: A pointer to the new ErrInvalidCall.
func NewErrInvalidCall(functionName string, function any, reason error) *ErrInvalidCall {
	return &ErrInvalidCall{
		FnName:    functionName,
		Signature: reflect.ValueOf(function).Type(),
		Reason:    reason,
	}
}

// ErrIgnorable represents an error that can be ignored. Useful for indicating
// that an error is ignorable.
type ErrIgnorable struct {
	// Err is the error that can be ignored.
	Err error
}

// Error is a method of the error interface.
// It does not mention that the error is ignorable.
//
// Returns:
//
//   - string: The error message of the ignorable error (no mention of being ignorable).
func (e *ErrIgnorable) Error() string {
	if e.Err == nil {
		return "ignorable error"
	} else {
		return e.Err.Error()
	}
}

// Unwrap returns the error that can be ignored.
// It is used for error unwrapping.
//
// Returns:
//   - error: The error that can be ignored.
func (e *ErrIgnorable) Unwrap() error {
	return e.Err
}

// ChangeReason changes the reason for the ignorable error.
//
// Parameters:
//   - reason: The new reason for the ignorable error.
func (e *ErrIgnorable) ChangeReason(reason error) {
	e.Err = reason
}

// NewErrIgnorable creates a new ErrIgnorable error.
//
// Parameters:
//   - err: The error that can be ignored.
//
// Returns:
//   - *ErrIgnorable: A pointer to the newly created ErrIgnorable.
func NewErrIgnorable(err error) *ErrIgnorable {
	return &ErrIgnorable{
		Err: err,
	}
}

// ErrInvalidRune represents an error when an invalid rune is encountered.
type ErrInvalidRune struct {
	// Reason is the reason for the invalidity of the rune.
	Reason error
}

// Error is a method of the error interface that returns the error message.
//
// If the reason is not provided (nil), no reason is included in the error message.
//
// Returns:
//   - string: The error message.
func (e *ErrInvalidRune) Error() string {
	if e.Reason == nil {
		return "rune is invalid"
	} else {
		return fmt.Sprintf("invalid rune: %s", e.Reason.Error())
	}
}

// Unwrap returns the reason for the invalidity of the rune.
// It is used for error unwrapping.
//
// Returns:
//   - error: The reason for the invalidity of the rune.
func (e *ErrInvalidRune) Unwrap() error {
	return e.Reason
}

// ChangeReason changes the reason for the invalidity of the rune.
//
// Parameters:
//   - reason: The new reason for the invalidity of the rune.
func (e *ErrInvalidRune) ChangeReason(reason error) {
	e.Reason = reason
}

// NewErrInvalidRune creates a new ErrInvalidRuneAt error.
//
// Parameters:
//   - reason: The reason for the invalidity of the rune.
//
// Returns:
//   - *ErrInvalidRune: A pointer to the newly created ErrInvalidRune.
func NewErrInvalidRune(reason error) *ErrInvalidRune {
	return &ErrInvalidRune{
		Reason: reason,
	}
}

// ErrAt represents an error that occurs at a specific index.
type ErrAt struct {
	// Index is the index where the error occurred.
	Index int

	// Name is the name of the index.
	Name string

	// Reason is the reason for the error.
	Reason error
}

// Error is a method of the error interface that returns the error message.
//
// If the reason is not provided (nil), the error message is
// "at index %d: something went wrong".
//
// Returns:
//   - string: The error message.
func (e *ErrAt) Error() string {
	var builder strings.Builder

	if e.Reason == nil {
		builder.WriteString("something went wrong at the ")
	}

	var name string

	if e.Name != "" {
		name = e.Name
	} else {
		name = "index"
	}

	builder.WriteString(uc.GetOrdinalSuffix(e.Index))
	builder.WriteRune(' ')
	builder.WriteString(name)

	if e.Reason != nil {
		builder.WriteString(" is invalid: ")
		builder.WriteString(e.Reason.Error())
	}

	return builder.String()
}

// Unwrap returns the reason for the error.
// It is used for error unwrapping.
//
// Returns:
//   - error: The reason for the error.
func (e *ErrAt) Unwrap() error {
	return e.Reason
}

// ChangeReason changes the reason for the error.
//
// Parameters:
//   - reason: The new reason for the error.
func (e *ErrAt) ChangeReason(reason error) {
	e.Reason = reason
}

// NewErrAt creates a new ErrAt error.
//
// Parameters:
//   - index: The index where the error occurred.
//   - name: The name of the index.
//   - reason: The reason for the error.
//
// Returns:
//   - *ErrAt: A pointer to the newly created ErrAt.
func NewErrAt(index int, name string, reason error) *ErrAt {
	return &ErrAt{
		Index:  index,
		Name:   name,
		Reason: reason,
	}
}

// ErrAfter is an error that is returned when something goes wrong after a certain
// element in a stream of data.
type ErrAfter struct {
	// After is the element that was processed before the error occurred.
	After string

	// Reason is the reason for the error.
	Reason error
}

// Error is a method of the error interface.
//
// Returns:
//   - string: The error message.
func (e *ErrAfter) Error() string {
	if e.Reason == nil {
		return fmt.Sprintf("something went wrong after %s", e.After)
	} else {
		return fmt.Sprintf("after %s: %s", e.After, e.Reason.Error())
	}
}

// Unwrap returns the reason for the error.
// It is used for error unwrapping.
//
// Returns:
//   - error: The reason for the error.
func (e *ErrAfter) Unwrap() error {
	return e.Reason
}

// ChangeReason changes the reason for the error.
//
// Parameters:
//   - reason: The new reason for the error.
func (e *ErrAfter) ChangeReason(reason error) {
	e.Reason = reason
}

// NewErrAfter creates a new ErrAfter error.
//
// Parameters:
//   - after: The element that was processed before the error occurred.
//   - reason: The reason for the error.
//
// Returns:
//   - *ErrAfter: A pointer to the new ErrAfter error.
func NewErrAfter(after string, reason error) *ErrAfter {
	return &ErrAfter{
		After:  after,
		Reason: reason,
	}
}

// ErrBefore is an error that is returned when something goes wrong before
// a certain element in a stream of data.
type ErrBefore struct {
	// Before is the element that was processed before the error occurred.
	Before string

	// Reason is the reason for the error.
	Reason error
}

// Error is a method of the error interface.
//
// Returns:
//   - string: The error message.
func (e *ErrBefore) Error() string {
	if e.Reason == nil {
		return fmt.Sprintf("something went wrong before %s", e.Before)
	} else {
		return fmt.Sprintf("before %s: %s", e.Before, e.Reason.Error())
	}
}

// Unwrap returns the reason for the error.
// It is used for error unwrapping.
//
// Returns:
//   - error: The reason for the error.
func (e *ErrBefore) Unwrap() error {
	return e.Reason
}

// ChangeReason changes the reason for the error.
//
// Parameters:
//   - reason: The new reason for the error.
func (e *ErrBefore) ChangeReason(reason error) {
	e.Reason = reason
}

// NewErrBefore creates a new ErrBefore error.
//
// Parameters:
//   - before: The element that was processed before the error occurred.
//   - reason: The reason for the error.
//
// Returns:
//   - *ErrBefore: A pointer to the new ErrBefore error.
func NewErrBefore(before string, reason error) *ErrBefore {
	return &ErrBefore{
		Before: before,
		Reason: reason,
	}
}

// ErrInvalidUsage represents an error that occurs when a function is used incorrectly.
type ErrInvalidUsage struct {
	// Reason is the reason for the invalid usage.
	Reason error

	// Usage is the usage of the function.
	Usage string
}

// Error is a method of the error interface.
//
// Returns:
//   - string: The error message.
func (e *ErrInvalidUsage) Error() string {
	var reason string

	if e.Reason == nil {
		reason = "invalid usage"
	} else {
		reason = e.Reason.Error()
	}

	if e.Usage == "" {
		return reason
	}

	return fmt.Sprintf("%s. %s", reason, e.Usage)
}

// Unwrap returns the reason for the invalid usage.
// It is used for error unwrapping.
//
// Returns:
//   - error: The reason for the invalid usage.
func (e *ErrInvalidUsage) Unwrap() error {
	return e.Reason
}

// ChangeReason changes the reason for the invalid usage.
//
// Parameters:
//   - reason: The new reason for the invalid usage.
func (e *ErrInvalidUsage) ChangeReason(reason error) {
	e.Reason = reason
}

// NewErrInvalidUsage creates a new ErrInvalidUsage error.
//
// Parameters:
//   - reason: The reason for the invalid usage.
//   - usage: The usage of the function.
//
// Returns:
//   - *ErrInvalidUsage: A pointer to the new ErrInvalidUsage error.
func NewErrInvalidUsage(reason error, usage string) *ErrInvalidUsage {
	return &ErrInvalidUsage{
		Reason: reason,
		Usage:  usage,
	}
}

// ErrUnexpectedError represents an error that occurs unexpectedly.
type ErrUnexpectedError struct {
	// Reason is the reason for the unexpected error.
	Reason error
}

// Error is a method of the error interface.
//
// Returns:
//   - string: The error message.
func (e *ErrUnexpectedError) Error() string {
	if e.Reason == nil {
		return "unexpected error"
	} else {
		return fmt.Sprintf("unexpected error: %s", e.Reason.Error())
	}
}

// Unwrap returns the reason for the unexpected error.
// It is used for error unwrapping.
//
// Returns:
//   - error: The reason for the unexpected error.
func (e *ErrUnexpectedError) Unwrap() error {
	return e.Reason
}

// ChangeReason changes the reason for the unexpected error.
//
// Parameters:
//   - reason: The new reason for the unexpected error.
func (e *ErrUnexpectedError) ChangeReason(reason error) {
	e.Reason = reason
}

// NewErrUnexpectedError creates a new ErrUnexpectedError error.
//
// Parameters:
//   - reason: The reason for the unexpected error.
//
// Returns:
//   - *ErrUnexpectedError: A pointer to the new ErrUnexpectedError error.
func NewErrUnexpectedError(reason error) *ErrUnexpectedError {
	return &ErrUnexpectedError{
		Reason: reason,
	}
}

// ErrVariableError represents an error that occurs when a variable is invalid.
type ErrVariableError struct {
	// Variable is the name of the variable that caused the error.
	Variable string

	// Reason is the reason for the variable error.
	Reason error
}

// Error returns the error message: "variable (<variable>) error"
// or "variable (<variable>) error: <reason>" if the reason is provided.
//
// Returns:
//   - string: The error message.
func (e *ErrVariableError) Error() string {
	if e.Reason == nil {
		return fmt.Sprintf("variable (%s) error", e.Variable)
	} else {
		return fmt.Sprintf("variable (%s) error: %s", e.Variable, e.Reason.Error())
	}
}

// Unwrap returns the reason for the variable error.
// It is used for error unwrapping.
//
// Returns:
//   - error: The reason for the variable error.
func (e *ErrVariableError) Unwrap() error {
	return e.Reason
}

// ChangeReason changes the reason for the variable error.
//
// Parameters:
//   - reason: The new reason for the variable error.
func (e *ErrVariableError) ChangeReason(reason error) {
	e.Reason = reason
}

// NewErrVariableError creates a new ErrVariableError error.
//
// Parameters:
//   - variable: The name of the variable that caused the error.
//   - reason: The reason for the variable error.
//
// Returns:
//   - *ErrVariableError: A pointer to the new ErrVariableError error.
func NewErrVariableError(variable string, reason error) *ErrVariableError {
	return &ErrVariableError{
		Variable: variable,
		Reason:   reason,
	}
}
