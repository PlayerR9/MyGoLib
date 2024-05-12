package Tray

// ErrEmptyTray is an error type that represents an empty tray.
type ErrEmptyTray struct{}

// Error returns the error message: "tray is empty".
//
// Returns:
//   - string: The error message.
func (e *ErrEmptyTray) Error() string {
	return "tray is empty"
}

// NewErrEmptyTray creates a new ErrEmptyTray.
//
// Returns:
//   - *ErrEmptyTray: The new ErrEmptyTray.
func NewErrEmptyTray() *ErrEmptyTray {
	return &ErrEmptyTray{}
}
