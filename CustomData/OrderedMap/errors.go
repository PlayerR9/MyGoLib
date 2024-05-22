package OrderedMap

// ErrKeyNotFound is an error type that represents a key not found error.
type ErrKeyNotFound struct{}

// Error returns the error message: "key %q not found".
//
// Returns:
//   - string: The error message.
func (e *ErrKeyNotFound) Error() string {
	return "key not found"
}

// NewErrKeyNotFound creates a new ErrKeyNotFound.
//
// Returns:
//   - *ErrKeyNotFound: A pointer to the new ErrKeyNotFound.
func NewErrKeyNotFound() *ErrKeyNotFound {
	return &ErrKeyNotFound{}
}
