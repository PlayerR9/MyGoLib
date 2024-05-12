package CnsPanel

import "fmt"

// ErrCommandNotFound represents an error where a command is not found.
type ErrCommandNotFound struct {
	// The command that was not found.
	Command string
}

// Error returns the error message for an ErrCommandNotFound.
//
// Returns:
//   - string: The error message.
func (e *ErrCommandNotFound) Error() string {
	return fmt.Sprintf("command %q not found", e.Command)
}

// NewErrCommandNotFound creates a new ErrCommandNotFound with the provided command.
//
// Parameters:
//   - command: The command that was not found.
//
// Returns:
//   - *ErrCommandNotFound: A pointer to the new ErrCommandNotFound.
func NewErrCommandNotFound(command string) *ErrCommandNotFound {
	return &ErrCommandNotFound{Command: command}
}

// ErrParsingFlags represents an error where flags could not be parsed.
type ErrParsingFlags struct {
	// Name is the flag that could not be parsed.
	Name string

	// Reason is the reason the flag could not be parsed.
	Reason error
}

// Error returns the error message for an ErrParsingFlags.
//
// Returns:
//   - string: The error message.
func (e *ErrParsingFlags) Error() string {
	if e.Reason == nil {
		return fmt.Sprintf("flag %q could not be parsed", e.Name)
	} else {
		return fmt.Sprintf("invalid arguments for flag %q: %s", e.Name, e.Reason.Error())
	}
}

// Unwrap returns the reason the flags could not be parsed.
//
// Returns:
//   - error: The reason the flags could not be parsed.
func (e *ErrParsingFlags) Unwrap() error {
	return e.Reason
}

// NewErrParsingFlags creates a new ErrParsingFlags with the provided name and reason.
//
// Parameters:
//   - name: The flag that could not be parsed.
//   - reason: The reason the flag could not be parsed.
//
// Returns:
//   - *ErrParsingFlags: A pointer to the new ErrParsingFlags.
func NewErrParsingFlags(name string, reason error) *ErrParsingFlags {
	return &ErrParsingFlags{Name: name, Reason: reason}
}

// ErrUnknownFlag represents an error where an unknown flag is provided.
type ErrUnknownFlag struct{}

// Error returns the error message for an ErrUnknownFlag.
//
// Returns:
//   - string: The error message.
func (e *ErrUnknownFlag) Error() string {
	return "unknown flag"
}

// NewErrUnknownFlag creates a new ErrUnknownFlag.
//
// Returns:
//   - *ErrUnknownFlag: A pointer to the new ErrUnknownFlag.
func NewErrUnknownFlag() *ErrUnknownFlag {
	return &ErrUnknownFlag{}
}

// ErrFewArguments represents an error where not enough arguments are provided.
type ErrFewArguments struct{}

// Error returns the error message for an ErrFewArguments.
//
// Returns:
//   - string: The error message.
func (e *ErrFewArguments) Error() string {
	return "not enough arguments were provided"
}

// NewErrFewArguments creates a new ErrFewArguments.
//
// Returns:
//   - *ErrFewArguments: A pointer to the new ErrFewArguments.
func NewErrFewArguments() *ErrFewArguments {
	return &ErrFewArguments{}
}
