package FileManager

import "strings"

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

// ErrPathNot is an error type for when a path is not as expected.
type ErrPathNot struct {
	// Path is the path that was not as expected.
	Path string

	// Expected is the expected value of the path.
	Expected string
}

// Error returns the error message: "path <path> is not <expected>".
//
// Returns:
//   - string: The error message.
func (e *ErrPathNot) Error() string {
	var builder strings.Builder

	builder.WriteString("path ")
	builder.WriteString(e.Path)
	builder.WriteString(" is not ")
	builder.WriteString(e.Expected)

	return builder.String()
}

// NewErrPathNot creates a new ErrPathNot error.
//
// Parameters:
//   - path: A string representing the path that was not as expected.
//   - expected: A string representing the expected value of the path.
//
// Returns:
//   - *ErrPathNot: A pointer to the newly created error.
func NewErrPathNot(path, expected string) *ErrPathNot {
	return &ErrPathNot{
		Path:     path,
		Expected: expected,
	}
}
