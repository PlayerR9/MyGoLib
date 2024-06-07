package Slices

// ErrLastNotFound is an error type for when the last element is not found.
type ErrLastNotFound struct{}

// Error implements the error interface.
//
// It returns the message: "last element not found".
func (e *ErrLastNotFound) Error() string {
	return "last element not found"
}

// NewErrLastNotFound creates a new ErrLastNotFound.
//
// Returns:
//   - *ErrLastNotFound: The new ErrLastNotFound.
func NewErrLastNotFound() *ErrLastNotFound {
	return &ErrLastNotFound{}
}
