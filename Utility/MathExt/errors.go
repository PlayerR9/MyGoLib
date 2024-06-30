package MathExt

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
	e := &ErrSubtractionUnderflow{}
	return e
}
