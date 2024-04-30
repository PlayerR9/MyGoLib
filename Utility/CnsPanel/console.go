// Package CnsPanel provides a structure and functions for handling
// console command flags.
package CnsPanel

import (
	"slices"
	"strings"

	fs "github.com/PlayerR9/MyGoLib/Formatting/FString"
	ers "github.com/PlayerR9/MyGoLib/Units/Errors"
)

// ConsolePanel represents a command line console.
type ConsolePanel struct {
	// Name of the executable.
	executableName string

	// Description of the executable.
	description [][]string

	// Map of commands accepted by the console.
	commands []ConsoleCommandInfo
}

// ParseArgs parses the provided command line arguments
// and executes the corresponding command.
//
// Panics with an error of type *ErrInvalidParameter if no
// arguments are provided, or with an error of type *ErrCallFailed
// if the ParseArgs function fails.
//
// Parameters:
//   - args: The command line arguments to parse. Without the
//     executable name.
//
// Returns:
//   - parsedCommand: The parsed command.
//   - error: An error, if any.
func (cns *ConsolePanel) ParseArgs(args []string) (parsedCommand, error) {
	var pc parsedCommand

	// Check if any arguments were provided
	if len(args) == 0 {
		return pc, ers.NewErrInvalidParameter(
			"args",
			ers.NewErrEmptySlice(),
		)
	}

	// Get the command from the command map
	index := slices.IndexFunc(cns.commands, func(command ConsoleCommandInfo) bool {
		return command.name == args[0]
	})
	if index == -1 {
		return pc, NewErrCommandNotFound(args[0])
	}

	command := cns.commands[index]

	pc.command = command.name
	pc.args = make(Arguments)
	pc.callback = command.callback

	if len(args) == 0 {
		return pc, nil
	}

	// Create a map to store the flags
	var err error

	// Parse the flags if provided
	pc.args, err = parseConsoleFlags(args[1:], command.flags)
	if err != nil {
		return pc, NewErrParsingFlags(command.name, err)
	}

	return pc, nil
}

// parseConsoleFlags parses the provided arguments into console flags.
//
// Parameters:
//   - args: The arguments to parse.
//   - flags: The console flags to parse the arguments into.
//
// Returns:
//   - Arguments: A map of the parsed console flags.
//   - error: An error, if any.
func parseConsoleFlags(args []string, flags []ConsoleFlagInfo) (Arguments, error) {
	// Create a map to store the console flags for quick lookup
	consoleFlagMap := make(map[string]ConsoleFlagInfo)
	for _, consoleFlag := range flags {
		consoleFlagMap[consoleFlag.name] = consoleFlag
	}

	// Create a map to store the parsed results
	parsedResults := make(Arguments)
	currentArgIndex := 0

	for currentArgIndex < len(args) {
		// Get the console flag name from the current argument
		consoleFlagName := args[currentArgIndex]

		// Check if the console flag exists in the map
		consoleFlag, exists := consoleFlagMap[consoleFlagName]
		if !exists {
			return nil, NewErrUnknownFlag()
		}

		// Check if there are enough arguments for the console flag
		if len(consoleFlag.args)+currentArgIndex >= len(args) {
			return nil, NewErrFewArguments()
		}

		// Move to the next argument
		currentArgIndex++

		// Create a temporary slice to store the arguments for the console flag
		tempArgs := make([]string, len(args[currentArgIndex:]))
		copy(tempArgs, args[currentArgIndex:])

		// Call the callback function for the console flag with the arguments
		parsedFlag, err := consoleFlag.callback(tempArgs...)
		if err != nil {
			return nil, err
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
// Parameters:
//   - indentLevel: The level of indentation to use for the CMLine.
//
// Returns:
//   - []string: A slice of strings representing the CMLine.
func (cns *ConsolePanel) FString(indentLevel int) []string {
	indentCfig := fs.NewIndentConfig(fs.DefaultIndentation, 0, true, false)
	indent := indentCfig.String()

	results := make([]string, 0)

	var builder strings.Builder

	// Print the usage information

	builder.WriteString(indent)
	builder.WriteString("Usage: ")
	builder.WriteString(cns.executableName)
	builder.WriteString(" <command> [flags]")

	results = append(results, builder.String())

	// Add the description
	builder.Reset()

	builder.WriteString(indent)
	builder.WriteString("Description:")

	if len(cns.description) == 0 {
		builder.WriteString(" [No description provided]")

		results = append(results, builder.String())
	} else {
		results = append(results, builder.String())

		for _, line := range cns.description {
			builder.Reset()

			// FIXME: Pretty print the description
			builder.WriteString(indent)
			builder.WriteString(indent)

			builder.WriteString(strings.Join(line, " "))

			results = append(results, builder.String())
		}
	}

	// Add the list of commands
	builder.Reset()

	builder.WriteString(indent)
	builder.WriteString("Commands:")

	if len(cns.commands) == 0 {
		builder.WriteString(" None")

		results = append(results, builder.String())
	} else {
		results = append(results, builder.String())

		for _, command := range cns.commands {
			results = append(results, command.FString(indentLevel+1)...)
		}
	}

	return results
}
