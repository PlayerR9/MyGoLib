package Header

// ErrEmptyTitle represents an error that is thrown when a title is expected but not provided.
// It is a struct type that implements the error interface by providing an Error method.
type ErrEmptyTitle struct{}

func (e ErrEmptyTitle) Error() string {
	return "title cannot be empty"
}

// ErrDataNotString is a custom error type that represents an
// error when the data is not a string.
type ErrDataNotString struct{}

func (e ErrDataNotString) Error() string {
	return "data is not a string"
}

// ErrDataNotCounter is a custom error type that represents an
// error when the data is not a counter.
type ErrDataNotCounter struct{}

func (e ErrDataNotCounter) Error() string {
	return "data is not a counter"
}
