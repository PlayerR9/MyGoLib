package FileManager

// ErrFileNotOpen is an error type for when a file was not opened.
type ErrFileNotOpen struct{}

// Error returns the error message: "file was not opened".
func (e *ErrFileNotOpen) Error() string {
	return "file was not opened"
}

// NewErrFileNotOpen creates a new ErrFileNotOpen error.
//
// Returns:
//   - *ErrFileNotOpen: A pointer to the newly created error.
func NewErrFileNotOpen() *ErrFileNotOpen {
	return &ErrFileNotOpen{}
}
