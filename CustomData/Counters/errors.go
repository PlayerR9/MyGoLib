package Counters

import "fmt"

// ErrInvalidParameter represents an error that occurs when an
// invalid parameter is provided.
// It includes the name of the parameter and the reason for the
// error.
type ErrInvalidParameter struct {
	// parameter is the name of the parameter that caused the error.
	parameter string

	// reason is the specific error that occurred.
	reason error
}

// Parameter sets the name of the parameter that caused the error.
// It follows the builder pattern, returning the ErrInvalidParameter
// itself for chaining.
//
// The function takes the following parameter:
//
// 	- p is the name of the parameter that caused the error.
//
// The function returns the following:
//
// 	- A pointer to the ErrInvalidParameter itself.
func (e *ErrInvalidParameter) Parameter(p string) *ErrInvalidParameter {
	e.parameter = p
	return e
}

// Reason sets the specific error that occurred.
// It follows the builder pattern, returning the ErrInvalidParameter
// itself for chaining.
//
// The function takes the following parameter:
//
// 	- r is the specific error that occurred.
//
// The function returns the following:
//
// 	- A pointer to the ErrInvalidParameter itself.
func (e *ErrInvalidParameter) Reason(r error) *ErrInvalidParameter {
	e.reason = r
	return e
}

// Error returns a string representation of the ErrInvalidParameter.
// It includes the name of the parameter and the reason for the error.
//
// The function returns the following:
//
// 	- A string representing the ErrInvalidParameter.
func (e *ErrInvalidParameter) Error() string {
	return fmt.Sprintf("invalid parameter %s: reason=%v", e.parameter, e.reason)
}

// Unwrap returns the specific error that occurred, which is wrapped
// in the ErrInvalidParameter.
// This allows the use of errors.Is and errors.As to check for specific
// errors.
//
// The function returns the following:
//
// 	- The specific error that occurred.
func (e *ErrInvalidParameter) Unwrap() error {
	return e.reason
}

// ErrCannotAdvanceCounter represents an error that occurs when a
// counter cannot be advanced.
// It includes the counter that caused the error and the reason
// for the error.
type ErrCannotAdvanceCounter struct {
	// counter is the counter that caused the error.
	counter Counter

	// reason is the specific error that occurred.
	reason error
}

// Counter sets the counter that caused the error.
// It follows the builder pattern, returning the ErrCannotAdvanceCounter
// itself for chaining.
//
// The function takes the following parameter:
//
// 	- c is the counter that caused the error.
//
// The function returns the following:
//
// 	- A pointer to the ErrCannotAdvanceCounter itself.
func (e *ErrCannotAdvanceCounter) Counter(c Counter) *ErrCannotAdvanceCounter {
	e.counter = c
	return e
}

// Reason sets the specific error that occurred.
// It follows the builder pattern, returning the ErrCannotAdvanceCounter
// itself for chaining.
//
// The function takes the following parameter:
//
// 	- r is the specific error that occurred.
//
// The function returns the following:
//
// 	- A pointer to the ErrCannotAdvanceCounter itself.
func (e *ErrCannotAdvanceCounter) Reason(r error) *ErrCannotAdvanceCounter {
	e.reason = r
	return e
}

// Error returns a string representation of the ErrCannotAdvanceCounter.
// It includes the type of the counter and the reason for the error.
//
// The function returns the following:
//
// 	- A string representing the ErrCannotAdvanceCounter.
func (e *ErrCannotAdvanceCounter) Error() string {
	return fmt.Sprintf("cannot advance %T: reason=%v", e.counter, e.reason)
}

// Unwrap returns the specific error that occurred, which is wrapped in
// the ErrCannotAdvanceCounter.
// This allows the use of errors.Is and errors.As to check for specific
// errors.
//
// The function returns the following:
//
// 	- The specific error that occurred.
func (e *ErrCannotAdvanceCounter) Unwrap() error {
	return e.reason
}

// ErrCannotRetreatCounter represents an error that occurs when a
// counter cannot be retreated.
// It includes the counter that caused the error and the reason
// for the error.
type ErrCannotRetreatCounter struct {
	// counter is the counter that caused the error.
	counter Counter

	// reason is the specific error that occurred.
	reason error
}

// Counter sets the counter that caused the error.
// It follows the builder pattern, returning the ErrCannotRetreatCounter
// itself for chaining.
//
// The function takes the following parameter:
//
// 	- c is the counter that caused the error.
//
// The function returns the following:
//
// 	- A pointer to the ErrCannotRetreatCounter itself.
func (e *ErrCannotRetreatCounter) Counter(c Counter) *ErrCannotRetreatCounter {
	e.counter = c
	return e
}

// Reason sets the specific error that occurred.
// It follows the builder pattern, returning the ErrCannotRetreatCounter
// itself for chaining.
//
// The function takes the following parameter:
//
// 	- r is the specific error that occurred.
//
// The function returns the following:
//
// 	- A pointer to the ErrCannotRetreatCounter itself.
func (e *ErrCannotRetreatCounter) Reason(r error) *ErrCannotRetreatCounter {
	e.reason = r
	return e
}

// Error returns a string representation of the ErrCannotRetreatCounter.
// It includes the type of the counter and the reason for the error.
//
// The function returns the following:
//
// 	- A string representing the ErrCannotRetreatCounter.
func (e *ErrCannotRetreatCounter) Error() string {
	return fmt.Sprintf("cannot retreat %T: reason=%v", e.counter, e.reason)
}

// Unwrap returns the specific error that occurred, which is wrapped
// in the ErrCannotRetreatCounter.
// This allows the use of errors.Is and errors.As to check for specific
// errors.
//
// The function returns the following:
//
// 	- The specific error that occurred.
func (e *ErrCannotRetreatCounter) Unwrap() error {
	return e.reason
}
