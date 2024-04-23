package DoubleLL

// ErrNoElementsHaveBeenPopped represents an error where no elements have been popped.
type ErrNoElementsHaveBeenPopped struct{}

// Error is a method of the error interface.
//
// Returns:
//
//   - string: The error message.
func (e ErrNoElementsHaveBeenPopped) Error() string {
	return "no elements have been popped"
}

// NewErrNoElementsHaveBeenPopped creates a new instance of ErrNoElementsHaveBeenPopped.
//
// Returns:
//
//   - error: An error of type *ErrNoElementsHaveBeenPopped.
func NewErrNoElementsHaveBeenPopped() error {
	return &ErrNoElementsHaveBeenPopped{}
}
