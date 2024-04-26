package Stream

// ErrNoMoreItems is an error that indicates that
// there are no more items in the stream.
type ErrNoMoreItems struct{}

// Error is a method of errors that returns the error message.
//
// Returns:
//  	- string: The error message.
func (e *ErrNoMoreItems) Error() string {
	return "no more items in the stream"
}

// NewErrNoMoreItems creates a new ErrNoMoreItems error.
//
// Returns:
//  	- *ErrNoMoreItems: A pointer to the new error.
func NewErrNoMoreItems() *ErrNoMoreItems {
	return &ErrNoMoreItems{}
}
