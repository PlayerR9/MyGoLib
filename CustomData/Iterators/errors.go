package Iterators

// ErrNotInitialized is an error type that is returned when an iterator
// is not initialized.
type ErrNotInitialized struct{}

// Error is a method of the error interface that returns the error message.
//
// Returns:
//   - string: The error message.
func (e *ErrNotInitialized) Error() string {
	return "iterator is not initialized"
}

// NewErrNotInitialized creates a new ErrNotInitialized error.
//
// Returns:
//   - *ErrNotInitialized: A pointer to the new error.
func NewErrNotInitialized() *ErrNotInitialized {
	return &ErrNotInitialized{}
}

// ErrExhaustedIter is an error type that is returned when an iterator
// is exhausted (i.e., there are no more elements to consume).
type ErrExhaustedIter struct{}

// Error is a method of the error interface that returns the error message.
//
// Returns:
//   - string: The error message.
func (e *ErrExhaustedIter) Error() string {
	return "iterator is exhausted"
}

// NewErrExhaustedIter creates a new ErrExhaustedIter error.
//
// Returns:
//   - *ErrExhaustedIter: A pointer to the new error.
func NewErrExhaustedIter() *ErrExhaustedIter {
	return &ErrExhaustedIter{}
}
