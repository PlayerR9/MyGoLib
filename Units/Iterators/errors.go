package Iterators

// ErrNotInitialized is an error type that is returned when an iterator
// is not initialized.
type ErrNotInitialized struct{}

// Error implements the error interface.
//
// Message: "iterator is not initialized"
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

// Error implements the error interface.
//
// Message: "iterator is exhausted"
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
