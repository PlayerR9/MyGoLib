package MathExt

// ErrInvalidBase is an error that is returned when the base is less than 1.
type ErrInvalidBase struct{}

// NewErrInvalidBase creates a new ErrInvalidBase error.
//
// Returns:
//
//   - error: A new ErrInvalidBase error.
func NewErrInvalidBase() error {
	return &ErrInvalidBase{}
}

// Error is a method of ErrInvalidBase that returns the error message.
//
// Returns:
//
//   - string: The error message.
func (e *ErrInvalidBase) Error() string {
	return "base cannot be less than 1"
}
