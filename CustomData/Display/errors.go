package Display

type ErrWidthTooSmall struct{}

func (e ErrWidthTooSmall) Error() string {
	return "width must be at least 5"
}

type ErrHeightTooSmall struct{}

func (e ErrHeightTooSmall) Error() string {
	return "height must be at least 2"
}

type ErrXOutOfBounds struct{}

func (e ErrXOutOfBounds) Error() string {
	return "x is out of bounds"
}

type ErrYOutOfBounds struct{}

func (e ErrYOutOfBounds) Error() string {
	return "y is out of bounds"
}

type ErrTextTooLong struct{}

func (e ErrTextTooLong) Error() string {
	return "text is too long"
}
