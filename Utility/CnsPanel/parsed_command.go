// Package CnsPanel provides a structure and functions for handling
// console command flags.
package CnsPanel

// ParsedCommand represents a parsed console command.
type ParsedCommand struct {
	// command is the name of the command.
	command string

	// args is the arguments passed to the command.
	args Arguments

	// callback is the function to call when the command is used.
	callback CommandCallbackFunc
}

// Command returns the name of the command.
//
// Returns:
//   - string: The name of the command.
func (pc *ParsedCommand) Command() string {
	return pc.command
}

// Run executes the callback function for the parsed command.
//
// Returns:
//   - error: An error, if any.
func (pc *ParsedCommand) Run() error {
	if pc.callback == nil {
		return nil
	}

	return pc.callback(pc.args)
}
