package General

// ExitCode is a custom type that represents the exit code of a program.
type ExitCode int8

const (
	// Success indicates that the program has finished successfully.
	Success ExitCode = iota

	// Panic indicates that a panic occurred during the execution of the program.
	Panic

	// SetupFailed indicates that the program could not be set up.
	SetupFailed

	// Error indicates that an error occurred during the execution of the program.
	Error
)

// String is a method of fmt.Stringer interface.
//
// Return:
//
//   - string: A string representation of the exit code.
func (ec ExitCode) String() string {
	return [...]string{
		"Program has finished successfully",
		"Panic occurred",
		"Cound not set up the program",
		"An error occurred",
	}[ec]
}
