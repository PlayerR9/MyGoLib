package Display

import "fmt"

// ErrWidthExceeded is an error that occurs when the width of a draw table is exceeded.
type ErrWidthExceeded struct {
	// Amount is the amount by which the width was exceeded.
	Amount int
}

// Error returns the error message: "Width exceeded by <Amount>".
//
// Returns:
//   - string: The error message.
func (e *ErrWidthExceeded) Error() string {
	return fmt.Sprintf("Width exceeded by %d", e.Amount)
}

// NewErrWidthExceeded creates a new ErrWidthExceeded with the given amount.
//
// Parameters:
//   - amount: The amount by which the width was exceeded.
//
// Returns:
//   - *ErrWidthExceeded: The new ErrWidthExceeded.
func NewErrWidthExceeded(amount int) *ErrWidthExceeded {
	return &ErrWidthExceeded{
		Amount: amount,
	}
}

// ErrHeightExceeded is an error that occurs when the height of a draw table is exceeded.
type ErrHeightExceeded struct {
	// Amount is the amount by which the height was exceeded.
	Amount int
}

// Error returns the error message: "Height exceeded by <Amount>".
//
// Returns:
//   - string: The error message.
func (e *ErrHeightExceeded) Error() string {
	return fmt.Sprintf("Height exceeded by %d", e.Amount)
}

// NewErrHeightExceeded creates a new ErrHeightExceeded with the given amount.
//
// Parameters:
//   - amount: The amount by which the height was exceeded.
func NewErrHeightExceeded(amount int) *ErrHeightExceeded {
	return &ErrHeightExceeded{
		Amount: amount,
	}
}
