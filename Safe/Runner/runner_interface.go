package Runner

// Runner is an interface that defines the behavior of a type that can be started,
// stopped, and waited for.
type Runner interface {
	// Start starts the runner.
	//
	// Returns:
	//   - error: An error of type *ErrAlreadyRunning if the runner is already running.
	//   or any other error if the runner could not be started.
	Start() error

	// Wait waits for the runner to finish.
	Wait()

	// Close closes the runner.
	Close()

	// IsRunning returns true if the runner is running, false otherwise.
	IsRunning() bool
}
