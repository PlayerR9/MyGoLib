package PageInterval

// ErrMinLessThanOne is an error type that represents the condition
// where the minimum value is less than one. It implements the error interface.
type ErrMinLessThanOne struct{}

func (e *ErrMinLessThanOne) Error() string {
	return "the minimum page number must be greater than 0"
}

// ErrMaxLessThanMin is an error type that represents the condition
// where the maximum page number is less than the minimum page number. It implements the error interface.
type ErrMaxLessThanMin struct{}

func (e *ErrMaxLessThanMin) Error() string {
	return "the maximum page number must be greater than the minimum page number"
}

// ErrNoPagesHaveBeenSet is an error type that represents the condition
// where no pages have been set in a book or document. It implements the error interface.
type ErrNoPagesHaveBeenSet struct{}

func (e *ErrNoPagesHaveBeenSet) Error() string {
	return "no pages have been set"
}
