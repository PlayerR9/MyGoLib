package MathExt

type ErrInvalidBase struct{}

func NewErrInvalidBase() error {
	return &ErrInvalidBase{}
}

func (e *ErrInvalidBase) Error() string {
	return "base cannot be less than 1"
}
