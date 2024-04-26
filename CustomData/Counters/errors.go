package Counters

// ErrCurrentCountBelowZero represents an error where the current count
// is already at or below zero.
type ErrCurrentCountBelowZero struct{}

// Error is a method of the error interface that returns the error message.
//
// Returns:
//   - string: The error message.
func (e *ErrCurrentCountBelowZero) Error() string {
	return "current count is already at or below zero"
}

// NewErrCurrentCountBelowZero creates a new ErrCurrentCountBelowZero error.
//
// Returns:
//   - *ErrCurrentCountBelowZero: A pointer to the new error.
func NewErrCurrentCountBelowZero() *ErrCurrentCountBelowZero {
	return &ErrCurrentCountBelowZero{}
}

// ErrCurrentCountAboveUpperLimit represents an error where the current count
// is already at or beyond the upper limit.
type ErrCurrentCountAboveUpperLimit struct{}

// Error is a method of the error interface.
//
// Returns:
//   - string: The error message.
func (e *ErrCurrentCountAboveUpperLimit) Error() string {
	return "current count is already at or beyond the upper limit"
}

// NewErrCurrentCountAboveUpperLimit creates a new
// ErrCurrentCountAboveUpperLimit error.
//
// Returns:
//   - *ErrCurrentCountAboveUpperLimit: A pointer to the new error.
func NewErrCurrentCountAboveUpperLimit() *ErrCurrentCountAboveUpperLimit {
	return &ErrCurrentCountAboveUpperLimit{}
}
