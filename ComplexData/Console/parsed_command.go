package Console

// ParsedCommand represents a parsed command.
type ParsedCommand struct {
	// command is the name of the command
	command string

	// flagMap is a map of flags and their values
	flagMap map[string]any

	// fn is the function to execute
	fn ConsoleFunc
}

// NewParsedCommand creates a new ParsedCommand.
//
// Parameters:
//   - command: The name of the command.
//   - flagMap: A map of flags and their values.
//   - fn: The function to execute.
//
// Returns:
//   - *ParsedCommand: The new ParsedCommand.
func NewParsedCommand(command string, flagMap map[string]any, fn ConsoleFunc) *ParsedCommand {
	return &ParsedCommand{
		command: command,
		flagMap: flagMap,
		fn:      fn,
	}
}

// Execute executes the command with the given flags.
//
// Returns:
//   - error: An error if the command failed to execute.
func (pc *ParsedCommand) Execute() (any, error) {
	return pc.fn(pc.flagMap)
}

// GetCommand returns the name of the command.
//
// Returns:
//   - string: The name of the command.
func (pc *ParsedCommand) GetCommand() string {
	return pc.command
}
