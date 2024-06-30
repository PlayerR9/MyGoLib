package FileManager

import "strings"

// ErrFileNotOpen is an error type for when a file was not opened.
type ErrFileNotOpen struct{}

// Error implements the error interface.
//
// Message: "file was not opened"
func (e *ErrFileNotOpen) Error() string {
	return "file was not opened"
}

// NewErrFileNotOpen creates a new ErrFileNotOpen error.
//
// Returns:
//   - *ErrFileNotOpen: A pointer to the newly created error.
func NewErrFileNotOpen() *ErrFileNotOpen {
	e := &ErrFileNotOpen{}
	return e
}

// ErrPathNot is an error type for when a path is not as expected.
type ErrPathNot struct {
	// Path is the path that was not as expected.
	Path string

	// Expected is the expected value of the path.
	Expected string
}

// Error implements the error interface.
//
// Message: "path {Path} is not {Expected}"
func (e *ErrPathNot) Error() string {
	values := []string{
		"path",
		e.Path,
		"is not",
		e.Expected,
	}

	msg := strings.Join(values, " ")

	return msg
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
	e := &ErrPathNot{
		Path:     path,
		Expected: expected,
	}

	return e
}
