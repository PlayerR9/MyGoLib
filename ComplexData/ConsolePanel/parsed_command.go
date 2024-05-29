// Package CnsPanel provides a structure and functions for handling
// console command flags.
package ConsolePanel

// ParsedCommand represents a parsed console command.
type ParsedCommand struct {
	// name is the name of the command.
	name string

	// args is the arguments passed to the command.
	args map[string]any

	// callback is the function to call when the command is used.
	callback CommandCallbackFunc
}

// newParsedCommand creates a new ParsedCommand with the provided name, arguments,
// and callback function.
//
// Parameters:
//   - name: The name of the command.
//   - args: A map of string keys to any values representing the arguments passed
//     to the command.
//   - callbackFunc: The function to call when the command is used.
//
// Returns:
//   - *ParsedCommand: A pointer to the new ParsedCommand.
func newParsedCommand(name string, args map[string]any, callbackFunc CommandCallbackFunc) *ParsedCommand {
	return &ParsedCommand{
		name:     name,
		args:     args,
		callback: callbackFunc,
	}
}

// GetName returns the name of the command.
//
// Returns:
//   - string: The name of the command.
func (pc *ParsedCommand) GetName() string {
	return pc.name
}

// Execute executes the callback function for the parsed command.
//
// Returns:
//   - any: The result of the callback function.
//   - error: An error if the callback function fails.
func (pc *ParsedCommand) Execute() (any, error) {
	return pc.callback(pc.args)
}
