package MathExt

// ErrInvalidBase is an error that is returned when the base is less than 1.
type ErrInvalidBase struct{}

// Error is a method of ErrInvalidBase that returns message: "base cannot be less than 1".
//
// Returns:
//   - string: The error message.
func (e *ErrInvalidBase) Error() string {
	return "base cannot be less than 1"
}

// NewErrInvalidBase creates a new ErrInvalidBase error.
//
// Returns:
//   - *ErrInvalidBase: The new ErrInvalidBase error.
func NewErrInvalidBase() *ErrInvalidBase {
	return &ErrInvalidBase{}
}

// ErrSubtractionUnderflow is an error that is returned when a subtraction
// operation results in a negative number.
type ErrSubtractionUnderflow struct{}

// Error is a method of ErrSubtractionUnderflow that returns the message:
// "subtraction underflow".
//
// Returns:
//   - string: The error message.
func (e *ErrSubtractionUnderflow) Error() string {
	return "subtraction underflow"
}

// NewErrSubtractionUnderflow creates a new ErrSubtractionUnderflow error.
//
// Returns:
//   - *ErrSubtractionUnderflow: The new ErrSubtractionUnderflow error.
func NewErrSubtractionUnderflow() *ErrSubtractionUnderflow {
	return &ErrSubtractionUnderflow{}
}
