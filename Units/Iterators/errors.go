package Iterators

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
