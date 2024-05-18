package Console

type ParsedCommand struct {
	// command is the name of the command
	command string

	// flagMap is a map of flags and their values
	flagMap map[string]any

	// fn is the function to execute
	fn ConsoleFunc
}

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
func (pc *ParsedCommand) GetCommand() string {
	return pc.command
}
