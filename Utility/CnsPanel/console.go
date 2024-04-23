// Package CnsPanel provides a structure and functions for handling
// console command flags.
package CnsPanel

import (
	"errors"
	"fmt"

	fs "github.com/PlayerR9/MyGoLib/Formatting/FString"
	ers "github.com/PlayerR9/MyGoLib/Units/Errors"
)

// ConsoleBuilder represents a builder for a command line interface.
type ConsoleBuilder struct {
	// Name of the executable.
	execName string

	// Description of the executable.
	description [][]string

	// Map of commands accepted by the interface.
	commands map[string]*ConsoleCommandInfo
}

// SetExecutableName is a method of ConsoleBuilder that sets the
// executable name for a ConsoleBuilder.
//
// Parameters:
//
//   - name: The name of the executable.
func (b *ConsoleBuilder) SetExecutableName(name string) {
	b.execName = name
}

// AppendParagraph is a method of ConsoleBuilder that appends a
// paragraph to the description of a ConsoleBuilder.
//
// Parameters:
//
//   - contents: The contents of the paragraph to append.
func (b *ConsoleBuilder) AppendParagraph(contents ...string) {
	if b.description == nil {
		b.description = [][]string{contents}
	} else {
		b.description = append(b.description, contents)
	}
}

// AddCommand is a method of ConsoleBuilder that adds a new command
// to a ConsoleBuilder.
//
// Parameters:
//
//   - commandName: The name of the command.
//   - options: The options to apply to the command.
//
// Returns:
//
//   - error: An error if the command cannot be added.
func (b *ConsoleBuilder) AddCommand(commandName string, options ...CommandInfoOption) error {
	if b.commands == nil {
		b.commands = make(map[string]*ConsoleCommandInfo)
	} else if _, exists := b.commands[commandName]; exists {
		return fmt.Errorf("command '%s' already exists", commandName)
	}

	newCommand := &ConsoleCommandInfo{
		name:        commandName,
		description: make([]string, 0),
		flags:       make([]*ConsoleFlagInfo, 0),
		callback:    nil,
	}

	for i, option := range options {
		err := option(newCommand)
		if err != nil {
			return fmt.Errorf("invalid option %d for command %s: %v", i, commandName, err)
		}
	}

	b.commands[commandName] = newCommand

	return nil
}

// Build is a method of ConsoleBuilder that builds a CMLine from a
// ConsoleBuilder.
//
// Returns:
//
//   - *CMLine: A CMLine built from the ConsoleBuilder.
func (b *ConsoleBuilder) Build() *ConsolePanel {
	cm := &ConsolePanel{executableName: b.execName}

	if b.commands == nil {
		cm.commands = make(map[string]*ConsoleCommandInfo)
		cm.description = make([][]string, 0)
	} else {
		cm.commands = b.commands
		cm.description = b.description
	}

	// Clear the ConsoleBuilder
	for i := range b.description {
		b.description[i] = nil
	}

	b.description = nil

	b.commands = nil

	return cm
}

// Reset is a method of ConsoleBuilder that resets a ConsoleBuilder.
func (b *ConsoleBuilder) Reset() {
	b.execName = ""

	for i := range b.description {
		b.description[i] = nil
	}

	b.description = nil
	b.commands = nil
}

// ConsolePanel represents a command line interface.
type ConsolePanel struct {
	// Name of the executable.
	executableName string

	// Description of the executable.
	description [][]string

	// Map of commands accepted by the interface.
	commands map[string]*ConsoleCommandInfo
}

// ParseArgs parses the provided command line arguments
// and executes the corresponding command.
//
// Panics with an error of type *ErrInvalidParameter if no
// arguments are provided, or with an error of type *ErrCallFailed
// if the ParseArgs function fails.
//
// Parameters:
//
//   - args: The command line arguments to parse. Without the
//     executable name.
//
// Returns:
//
//   - *parsedCommand: The parsed command.
//   - error: An error, if any.
func (cns *ConsolePanel) ParseArgs(args []string) (*parsedCommand, error) {
	// Check if any arguments were provided
	if len(args) == 0 {
		return nil, ers.NewErrInvalidParameter(
			"args",
			errors.New("no arguments provided"),
		)
	}

	// Get the command from the command map
	command, exists := cns.commands[args[0]]
	if !exists {
		return nil, fmt.Errorf("command '%s' not found", args[0])
	}

	pc := &parsedCommand{
		command:  command.name,
		args:     make(map[string]any),
		callback: command.callback,
	}

	// Create a map to store the flags
	var err error

	if len(args) > 1 {
		// Parse the flags if provided
		pc.args, err = parseConsoleFlags(args[1:], command.flags)
		if err != nil {
			return nil, err
		}
	}

	return pc, nil
}

// parseConsoleFlags parses the provided arguments into console flags.
//
// Parameters:
//
//   - args: The arguments to parse.
//   - flags: The console flags to parse the arguments into.
//
// Returns:
//
//   - map[string]any: A map of the parsed flags.
//   - error: An error, if any.
func parseConsoleFlags(args []string, flags []*ConsoleFlagInfo) (map[string]any, error) {
	// Create a map to store the console flags for quick lookup
	consoleFlagMap := make(map[string]*ConsoleFlagInfo)
	for _, consoleFlag := range flags {
		consoleFlagMap[consoleFlag.name] = consoleFlag
	}

	// Create a map to store the parsed results
	parsedResults := make(map[string]any)
	currentArgIndex := 0

	for currentArgIndex < len(args) {
		// Get the console flag name from the current argument
		consoleFlagName := args[currentArgIndex]

		// Check if the console flag exists in the map
		consoleFlag, exists := consoleFlagMap[consoleFlagName]
		if !exists {
			return nil, fmt.Errorf("unknown flag '%s' provided", consoleFlagName)
		}

		// Check if there are enough arguments for the console flag
		if len(consoleFlag.args)+currentArgIndex >= len(args) {
			return nil, fmt.Errorf("flag '%s' requires more arguments", consoleFlag.name)
		}

		// Move to the next argument
		currentArgIndex++

		// Create a temporary slice to store the arguments for the console flag
		tempArgs := make([]string, len(args[currentArgIndex:]))
		copy(tempArgs, args[currentArgIndex:])

		// Call the callback function for the console flag with the arguments
		parsedFlag, err := consoleFlag.callback(tempArgs...)
		if err != nil {
			return nil, fmt.Errorf("failed to parse flag '%s': reason=%v", consoleFlag.name, err)
		}

		// Store the result of the callback function in the parsed results map
		// amd move to the next argument
		parsedResults[consoleFlag.name] = parsedFlag
		currentArgIndex += len(consoleFlag.args)
	}

	return parsedResults, nil
}

// FString generates a formatted string representation of a CMLine.
// It includes the usage information, and the list of commands and
// their details.
//
// Returns:
//
//   - []string: A slice of strings representing the CMLine.
func (cns *ConsolePanel) FString(indentLevel int) []string {
	indentCfig := fs.NewIndentConfig(fs.DefaultIndentation, 0, true, false)
	indent := indentCfig.String()

	results := make([]string, 0)

	results = append(results, fmt.Sprintf("%sUsage: %s <command> [flags]", indent, cns.executableName))

	if len(cns.description) == 0 {
		results = append(results, fmt.Sprintf("%sDescription: [No description provided]", indent))
	} else {
		results = append(results, fmt.Sprintf("%sDescription:", indent))

		for _, line := range cns.description {
			// FIXME: Pretty print the description
			results = append(results, fmt.Sprintf("%s\t%s", indent, line))
		}
	}

	if len(cns.commands) == 0 {
		results = append(results, fmt.Sprintf("%sCommands: None", indent))
	} else {
		results = append(results, fmt.Sprintf("%sCommands:", indent))

		for _, command := range cns.commands {
			results = append(results, command.FString(indentLevel+1)...)
		}
	}

	return results
}
