package Runner

// ErrAlreadyRunning is an error type that represents an error where
// a process is already running.
type ErrAlreadyRunning struct{}

// Error returns the error message: "the process is already running".
//
// Returns:
//   - string: The error message.
func (e *ErrAlreadyRunning) Error() string {
	return "the process is already running"
}

// NewErrAlreadyRunning creates a new ErrAlreadyRunning error.
//
// Returns:
//   - *ErrAlreadyRunning: The new error.
func NewErrAlreadyRunning() *ErrAlreadyRunning {
	return &ErrAlreadyRunning{}
}
