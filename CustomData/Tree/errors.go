package Tree

// ErrMissingRoot is an error that is returned when the root of a tree is missing.
type ErrMissingRoot struct{}

// Error returns the error message: "missing root".
//
// Returns:
//   - string: The error message.
func (e *ErrMissingRoot) Error() string {
	return "missing root"
}

// NewErrMissingRoot creates a new ErrMissingRoot.
//
// Returns:
//   - *ErrMissingRoot: The newly created error.
func NewErrMissingRoot() *ErrMissingRoot {
	return &ErrMissingRoot{}
}
