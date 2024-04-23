package DoubleLL

// ErrNoElementsHaveBeenPopped represents an error where no elements have been popped.
type ErrNoElementsHaveBeenPopped struct{}

// Error is a method of the error interface.
//
// Returns:
//
//   - string: The error message.
func (e *ErrNoElementsHaveBeenPopped) Error() string {
	return "no elements have been popped"
}

// NewErrNoElementsHaveBeenPopped creates a new instance of ErrNoElementsHaveBeenPopped.
//
// Returns:
//
//   - *ErrNoElementsHaveBeenPopped: A pointer to the newly created error.
func NewErrNoElementsHaveBeenPopped() *ErrNoElementsHaveBeenPopped {
	return &ErrNoElementsHaveBeenPopped{}
}
