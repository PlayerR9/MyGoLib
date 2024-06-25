// Package errors provides a custom error type for out-of-bound errors.
package common

import (
	"reflect"
	"strconv"
	"strings"
)

// ErrWhile represents an error that occurs while performing an operation.
type ErrWhile struct {
	// Operation is the operation that was being performed.
	Operation string

	// Reason is the reason for the error.
	Reason error
}

// Error implements the Unwrapper interface.
//
// Message: "error while {operation}: {reason}"
//
// However, if the reason is nil, the message is "an error occurred while {operation}"
// instead.
func (e *ErrWhile) Error() string {
	var builder strings.Builder

	if e.Reason == nil {
		builder.WriteString("an error occurred while ")
		builder.WriteString(e.Operation)
	} else {
		builder.WriteString("error while ")
		builder.WriteString(e.Operation)
		builder.WriteString(": ")
		builder.WriteString(e.Reason.Error())
	}

	return builder.String()
}

// Unwrap implements the Unwrapper interface.
func (e *ErrWhile) Unwrap() error {
	return e.Reason
}

// ChangeReason implements the Unwrapper interface.
func (e *ErrWhile) ChangeReason(reason error) {
	e.Reason = reason
}

// NewErrWhile creates a new ErrWhile error.
//
// Parameters:
//   - operation: The operation that was being performed.
//   - reason: The reason for the error.
//
// Returns:
//   - *ErrWhile: A pointer to the newly created ErrWhile.
func NewErrWhile(operation string, reason error) *ErrWhile {
	return &ErrWhile{
		Operation: operation,
		Reason:    reason,
	}
}

// ErrWhileAt represents an error that occurs while performing an operation at a specific index.
type ErrWhileAt struct {
	// Index is the index where the error occurred.
	Index int

	// Element is the element where the index is pointing to.
	Element string

	// Operation is the operation that was being performed.
	Operation string

	// Reason is the reason for the error.
	Reason error
}

// Error implements the Unwrapper interface.
//
// Message: "an error occurred while {operation} at index {index} {element}: {reason}"
//
// However, if the reason is nil, the message is "an error occurred while {operation}
// at index {index} {element}" instead.
func (e *ErrWhileAt) Error() string {
	var builder strings.Builder

	if e.Reason == nil {
		builder.WriteString("an error occurred ")
	}

	builder.WriteString("while ")
	builder.WriteString(e.Operation)
	builder.WriteRune(' ')
	builder.WriteString(GetOrdinalSuffix(e.Index))
	builder.WriteRune(' ')
	builder.WriteString(e.Element)

	if e.Reason != nil {
		builder.WriteString(": ")
		builder.WriteString(e.Reason.Error())
	}

	return builder.String()
}

// Unwrap implements the Unwrapper interface.
func (e *ErrWhileAt) Unwrap() error {
	return e.Reason
}

// ChangeReason implements the Unwrapper interface.
func (e *ErrWhileAt) ChangeReason(reason error) {
	e.Reason = reason
}

// NewErrWhileAt creates a new ErrWhileAt error.
//
// Parameters:
//   - operation: The operation that was being performed.
//   - index: The index where the error occurred.
//   - elem: The element where the index is pointing to.
//   - reason: The reason for the error.
//
// Returns:
//   - *ErrWhileAt: A pointer to the newly created ErrWhileAt.
func NewErrWhileAt(operation string, index int, elem string, reason error) *ErrWhileAt {
	return &ErrWhileAt{
		Index:     index,
		Operation: operation,
		Element:   elem,
		Reason:    reason,
	}
}

// ErrNoError represents an error when no error occurs.
type ErrNoError struct {
	// Err is the reason for the no error error.
	Err error
}

// Error implements the Unwrapper interface.
//
// Message: "no error" if the reason is nil, otherwise the error message.
func (e *ErrNoError) Error() string {
	var msg string

	if e.Err == nil {
		msg = "no error"
	} else {
		msg = e.Err.Error()
	}

	return msg
}

// Unwrap implements the Unwrapper interface.
func (e *ErrNoError) Unwrap() error {
	return e.Err
}

// ChangeReason implements the Unwrapper interface.
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
	return &ErrNoError{
		Err: err,
	}
}

// ErrInvalidParameter represents an error when a parameter is invalid.
type ErrInvalidParameter struct {
	// Parameter is the name of the Parameter.
	Parameter string

	// Reason is the Reason for the invalidity of the parameter.
	Reason error
}

// Error implements the Unwrapper interface.
//
// Message: "parameter ({parameter}) is invalid: {reason}".
//
// However, if the reason is nil, the message is "parameter ({parameter})
// is invalid" instead.
func (e *ErrInvalidParameter) Error() string {
	var builder strings.Builder

	builder.WriteString("parameter (")
	builder.WriteString(strconv.Quote(e.Parameter))
	builder.WriteString(") is invalid")

	if e.Reason != nil {
		builder.WriteString(": ")
		builder.WriteString(e.Reason.Error())
	}

	return builder.String()
}

// Unwrap implements the Unwrapper interface.
func (e *ErrInvalidParameter) Unwrap() error {
	return e.Reason
}

// ChangeReason implements the Unwrapper interface.
func (e *ErrInvalidParameter) ChangeReason(reason error) {
	e.Reason = reason
}

// NewErrInvalidParameter creates a new ErrInvalidParameter error.
//
// Parameters:
//   - parameter: The name of the parameter.
//   - reason: The reason for the invalidity.
//
// Returns:
//   - *ErrInvalidParameter: A pointer to the newly created ErrInvalidParameter.
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

// Error implements the Unwrapper interface.
//
// Message: "call to {function}({signature}) failed: {reason}".
//
// However, if the reason is nil, the message is "call to {function}({signature})
// failed" instead.
func (e *ErrInvalidCall) Error() string {
	var builder strings.Builder

	builder.WriteString("call to ")
	builder.WriteString(e.FnName)
	builder.WriteString(e.Signature.String())
	builder.WriteString(" failed")

	if e.Reason != nil {
		builder.WriteString(": ")
		builder.WriteString(e.Reason.Error())
	}

	return builder.String()
}

// Unwrap implements the Unwrapper interface.
func (e *ErrInvalidCall) Unwrap() error {
	return e.Reason
}

// ChangeReason implements the Unwrapper interface.
func (e *ErrInvalidCall) ChangeReason(reason error) {
	e.Reason = reason
}

// NewErrInvalidCall creates a new ErrInvalidCall.
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

// Error implements the Unwrapper interface.
//
// Message: "ignorable error" if the reason is nil, otherwise the error message.
func (e *ErrIgnorable) Error() string {
	var msg string

	if e.Err == nil {
		msg = "an ignorable error occurred"
	} else {
		msg = e.Err.Error()
	}

	return msg
}

// Unwrap implements the Unwrapper interface.
func (e *ErrIgnorable) Unwrap() error {
	return e.Err
}

// ChangeReason implements the Unwrapper interface.
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

// Error implements the Unwrapper interface.
//
// Message: "invalid rune: {reason}".
//
// However, if the reason is nil, the message is "rune is invalid" instead.
func (e *ErrInvalidRune) Error() string {
	if e.Reason == nil {
		return "rune is invalid"
	}

	var builder strings.Builder

	builder.WriteString("invalid rune: ")
	builder.WriteString(e.Reason.Error())

	return builder.String()
}

// Unwrap implements the Unwrapper interface.
func (e *ErrInvalidRune) Unwrap() error {
	return e.Reason
}

// ChangeReason implements the Unwrapper interface.
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

// Error implements the Unwrapper interface.
//
// Message: "something went wrong at the {index} {name}: {reason}".
//
// However, if the reason is nil, the message is "something went wrong at the {index}
// {name}" instead.
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

	builder.WriteString(GetOrdinalSuffix(e.Index))
	builder.WriteRune(' ')
	builder.WriteString(name)

	if e.Reason != nil {
		builder.WriteString(" is invalid: ")
		builder.WriteString(e.Reason.Error())
	}

	return builder.String()
}

// Unwrap implements the Unwrapper interface.
func (e *ErrAt) Unwrap() error {
	return e.Reason
}

// ChangeReason implements the Unwrapper interface.
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

// Error implements the Unwrapper interface.
//
// Message: "after {after}: {reason}".
//
// However, if the reason is nil, the message is "something went wrong after {after}"
// instead.
func (e *ErrAfter) Error() string {
	var builder strings.Builder

	if e.Reason == nil {
		builder.WriteString("something went wrong after ")
		builder.WriteString(e.After)
	} else {
		builder.WriteString("after ")
		builder.WriteString(e.After)
		builder.WriteString(": ")
		builder.WriteString(e.Reason.Error())
	}

	return builder.String()
}

// Unwrap implements the Unwrapper interface.
func (e *ErrAfter) Unwrap() error {
	return e.Reason
}

// ChangeReason implements the Unwrapper interface.
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

// Error implements the Unwrapper interface.
//
// Message: "before {before}: {reason}".
//
// However, if the reason is nil, the message is "something went wrong before {before}"
// instead.
func (e *ErrBefore) Error() string {
	var builder strings.Builder

	if e.Reason == nil {
		builder.WriteString("something went wrong before ")
		builder.WriteString(e.Before)
	} else {
		builder.WriteString("before ")
		builder.WriteString(e.Before)
		builder.WriteString(": ")
		builder.WriteString(e.Reason.Error())
	}

	return builder.String()
}

// Unwrap implements the Unwrapper interface.
func (e *ErrBefore) Unwrap() error {
	return e.Reason
}

// ChangeReason implements the Unwrapper interface.
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

// Error is a method of the Unwrapper interface.
//
// Message: "{reason}. {usage}".
//
// However, if the reason is nil, the message is "invalid usage. {usage}" instead.
//
// If the usage is empty, no usage is added to the message.
func (e *ErrInvalidUsage) Error() string {
	var builder strings.Builder

	if e.Reason == nil {
		builder.WriteString("invalid usage")
	} else {
		builder.WriteString(e.Reason.Error())
	}

	if e.Usage == "" {
		builder.WriteString(". ")
		builder.WriteString(e.Usage)
	}

	return builder.String()
}

// Unwrap implements the Unwrapper interface.
func (e *ErrInvalidUsage) Unwrap() error {
	return e.Reason
}

// ChangeReason implements the Unwrapper interface.
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

// Error implements the Unwrapper interface.
//
// Message: "unexpected error: {reason}".
//
// However, if the reason is nil, the message is "unexpected error" instead.
func (e *ErrUnexpectedError) Error() string {
	if e.Reason == nil {
		return "unexpected error"
	}

	var builder strings.Builder

	builder.WriteString("unexpected error: ")
	builder.WriteString(e.Reason.Error())

	return builder.String()
}

// Unwrap implements the Unwrapper interface.
func (e *ErrUnexpectedError) Unwrap() error {
	return e.Reason
}

// ChangeReason implements the Unwrapper interface.
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

// Error implements the Unwrapper interface.
//
// Message: "variable ({variable}) error: {reason}".
//
// However, if the reason is nil, the message is "variable ({variable}) error" instead.
func (e *ErrVariableError) Error() string {
	var builder strings.Builder

	builder.WriteString("variable (")
	builder.WriteString(strconv.Quote(e.Variable))
	builder.WriteString(") error")

	if e.Reason != nil {
		builder.WriteString(": ")
		builder.WriteString(e.Reason.Error())
	}

	return builder.String()
}

// Unwrap implements the Unwrapper interface.
func (e *ErrVariableError) Unwrap() error {
	return e.Reason
}

// ChangeReason implements the Unwrapper interface.
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

// ErrPossibleError represents an error that occurs when a possible error is encountered.
type ErrPossibleError struct {
	// Reason is the reason for the possible error.
	Reason error

	// Possible is the possible error.
	Possible error
}

// Error implements the Unwrapper interface.
//
// Message: "{reason}. It is possible that {possible}".
//
// However, if the reason is nil, the message is "no error occurred. It is possible
// that {possible}" instead.
//
// If the possible error is nil, the "It is possible that {possible}" part is omitted.
func (e *ErrPossibleError) Error() string {
	var builder strings.Builder

	if e.Reason == nil {
		builder.WriteString("no error occurred")
	} else {
		builder.WriteString(e.Reason.Error())
	}

	if e.Possible != nil {
		builder.WriteString(". It is possible that ")
		builder.WriteString(e.Possible.Error())
	}

	return builder.String()
}

// Unwrap implements the Unwrapper interface.
func (e *ErrPossibleError) Unwrap() error {
	return e.Reason
}

// ChangeReason implements the Unwrapper interface.
func (e *ErrPossibleError) ChangeReason(reason error) {
	e.Reason = reason
}

// NewErrPossibleError creates a new ErrPossibleError error.
//
// Parameters:
//   - reason: The reason for the possible error.
//   - possible: The possible error.
//
// Returns:
//   - *ErrPossibleError: A pointer to the new ErrPossibleError error.
func NewErrPossibleError(reason error, possible error) *ErrPossibleError {
	return &ErrPossibleError{
		Reason:   reason,
		Possible: possible,
	}
}
