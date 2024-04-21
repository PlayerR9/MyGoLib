package Counters

// ErrCurrentCountBelowZero represents an error where the current count
// is already at or below zero.
type ErrCurrentCountBelowZero struct{}

// Error is a method of the error interface.
//
// Returns:
//
//   - string: The error message.
func (e ErrCurrentCountBelowZero) Error() string {
	return "current count is already at or below zero"
}

// NewErrCurrentCountBelowZero creates a new instance of ErrCurrentCountBelowZero.
//
// Returns:
//
//   - error: An error of type *ErrCurrentCountBelowZero.
func NewErrCurrentCountBelowZero() error {
	return &ErrCurrentCountBelowZero{}
}

// ErrCurrentCountAboveUpperLimit represents an error where the current count
// is already at or beyond the upper limit.
type ErrCurrentCountAboveUpperLimit struct{}

// Error is a method of the error interface.
//
// Returns:
//
//   - string: The error message.
func (e ErrCurrentCountAboveUpperLimit) Error() string {
	return "current count is already at or beyond the upper limit"
}

// NewErrCurrentCountAboveUpperLimit creates a new instance of ErrCurrentCountAboveUpperLimit.
//
// Returns:
//
//   - error: An error of type *ErrCurrentCountAboveUpperLimit.
func NewErrCurrentCountAboveUpperLimit() error {
	return &ErrCurrentCountAboveUpperLimit{}
}
