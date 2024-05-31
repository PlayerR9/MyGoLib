// Package CnsPanel provides a structure and functions for handling
// console command flags.
package pkg

// ParsedCommand represents a parsed console command.
type ParsedCommand struct {
	// name is the name of the command.
	name string

	// args is the arguments passed to the command.
	args map[string]map[string][]any

	// callback is the function to call when the command is used.
	callback CommandCallbackFunc
}

// newParsedCommand creates a new ParsedCommand with the provided name, arguments,
// and callback function.
//
// Parameters:
//   - name: The name of the command.
//   - args: The arguments passed to the command.
//   - callbackFunc: The function to call when the command is used.
//
// Returns:
//   - *ParsedCommand: A pointer to the new ParsedCommand.
func newParsedCommand(name string, result map[string]*FlagParseResult, callbackFunc CommandCallbackFunc) *ParsedCommand {
	pc := &ParsedCommand{
		name:     name,
		args:     make(map[string]map[string][]any),
		callback: callbackFunc,
	}

	if result == nil {
		return pc
	}

	for k, v := range result {
		pc.args[k] = v.GetResult()
	}

	return pc
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
