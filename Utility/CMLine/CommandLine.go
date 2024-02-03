// Package CMLine provides a structure and functions for handling
// console command flags.
package CMLine

import (
	"errors"
	"fmt"
	"strings"

	ers "github.com/PlayerR9/MyGoLib/Utility/Errors"
)

// ConsoleFlagInfo represents a flag for a console command. It
// contains the name, arguments, and description of the flag,
// as well as a boolean indicating whether the flag is required.
// The Callback function is invoked when the flag is used, with
// the provided arguments.
// It returns a result and an error, if any.
type ConsoleFlagInfo struct {
	// Name of the flag.
	name string

	// Slice of strings representing the arguments accepted by
	// the flag.
	args []string

	// Brief explanation of what the flag does.
	description []string

	// Boolean indicating whether the flag is required.
	required bool

	// Function invoked when the flag is used.
	callback func(...string) (any, error)
}

// FlagInfoOption is a function type that modifies ConsoleFlagInfo.
type FlagInfoOption func(*ConsoleFlagInfo)

// WithArgs is a FlagInfoOption that sets the arguments for a
// ConsoleFlagInfo.
// It trims the space from each argument and ignores empty
// arguments.
//
// Parameters:
//
//   - args: The arguments to set.
//
// Returns:
//
//   - A FlagInfoOption that sets the arguments for a
//
// ConsoleFlagInfo.
func WithArgs(args ...string) FlagInfoOption {
	return func(flag *ConsoleFlagInfo) {
		flag.args = make([]string, 0, len(args))
		for _, arg := range args {
			arg = strings.TrimSpace(arg)
			if arg != "" {
				flag.args = append(flag.args, arg)
			}
		}
		flag.args = flag.args[:len(flag.args)]
	}
}

// WithFlagDescription is a FlagInfoOption that sets the description
// for a ConsoleFlagInfo.
// It splits each line of the description by newline characters.
//
// Parameters:
//
//   - description: The description to set.
//
// Returns:
//
//   - A FlagInfoOption that sets the description for a
//     ConsoleFlagInfo.
func WithFlagDescription(description ...string) FlagInfoOption {
	return func(flag *ConsoleFlagInfo) {
		for _, line := range description {
			fields := strings.Split(line, "\n")
			flag.description = append(flag.description, fields...)
		}
	}
}

// WithRequired is a FlagInfoOption that sets whether a
// ConsoleFlagInfo is required.
//
// Parameters:
//
//   - required: Whether the flag is required.
//
// Returns:
//
//   - A FlagInfoOption that sets whether a ConsoleFlagInfo
//     is required.
func WithRequired(required bool) FlagInfoOption {
	return func(flag *ConsoleFlagInfo) {
		flag.required = required
	}
}

// FString generates a formatted string representation of a
// ConsoleFlagInfo.
// It includes the flag name, arguments, description, and whether
// it is required.
//
// Parameters:
//
//   - indentLevel: The level of indentation to use.
//
// Returns:
//
//   - A string representing the ConsoleFlagInfo.
func (cfi *ConsoleFlagInfo) FString(indentLevel int) string {
	var builder strings.Builder

	indentation := strings.Repeat("\t", indentLevel)

	// Add the flag name
	fmt.Fprintf(&builder, "%s%s", indentation, cfi.name)

	// Add the arguments
	for _, arg := range cfi.args {
		fmt.Fprintf(&builder, " <%s>", arg)
	}

	// Add the description
	fmt.Fprintf(&builder, "\n%sDescription:\n%s\t%s", indentation, indentation, cfi.description)

	// Add the required information
	if len(cfi.args) != 0 {
		fmt.Fprintf(&builder, "\n%sRequired: ", indentation)

		if cfi.required {
			fmt.Fprint(&builder, "Yes")
		} else {
			fmt.Fprint(&builder, "No")
		}
	}

	return builder.String()
}

// ConsoleCommandInfo represents a console command. It contains
// the name, description, and flags of the command. The Callback
// function is invoked when the command is executed, with the
// provided arguments. It returns a result and an error, if any.
type ConsoleCommandInfo struct {
	// Name of the command.
	name string

	// Brief explanation of what the command does.
	description []string

	// Slice of ConsoleFlagInfo representing the flags accepted by
	// the command.
	flags []*ConsoleFlagInfo

	// Function invoked when the command is executed.
	callback func(map[string]any) (any, error)
}

// CommandInfoOption is a function type that modifies
// ConsoleCommandInfo.
type CommandInfoOption func(*ConsoleCommandInfo) error

// WithFlag is a CommandInfoOption that adds a new flag to a
// ConsoleCommandInfo.
// It creates a new ConsoleFlagInfo with the provided name and
// callback, and applies the provided options to it. If the
// flag name is empty or the callback is nil, it sets the
// errReason of the ConsoleCommandInfo.
//
// Parameters:
//
//   - name: The name of the flag.
//   - callback: The function to call when the flag is used.
//   - options: The options to apply to the flag.
//
// Returns:
//
//   - A CommandInfoOption that adds a new flag to a
//     ConsoleCommandInfo.
func WithFlag(name string, callback func(...string) (any, error), options ...FlagInfoOption) CommandInfoOption {
	return func(command *ConsoleCommandInfo) error {
		newFlag := &ConsoleFlagInfo{
			name:        name,
			args:        make([]string, 0),
			description: make([]string, 0),
			required:    false,
			callback:    callback,
		}

		name = strings.TrimSpace(name)
		if name == "" {
			return fmt.Errorf("could not create flag: %v",
				ers.NewErrInvalidParameter(
					"name", errors.New("flag name cannot be empty")),
			)
		}

		if callback == nil {
			return fmt.Errorf("could not create flag %s: %v", name,
				ers.NewErrInvalidParameter(
					"callback", errors.New("flag callback cannot be nil")),
			)
		}

		for _, option := range options {
			option(newFlag)
		}

		command.flags = append(command.flags, newFlag)

		return nil
	}
}

// WithCallback is a CommandInfoOption that sets the callback for
// a ConsoleCommandInfo.
// If the provided callback is nil, it sets the errReason of the
// ConsoleCommandInfo.
//
// Parameters:
//
//   - callback: The function to call when the command is used.
//
// Returns:
//
//   - A CommandInfoOption that sets the callback for a
//     ConsoleCommandInfo.
func WithCallback(callback func(map[string]any) (any, error)) CommandInfoOption {
	return func(command *ConsoleCommandInfo) error {
		if callback == nil {
			return ers.NewErrInvalidParameter("callback").
				WithReason(errors.New("callback cannot be nil"))
		}

		command.callback = callback

		return nil
	}
}

// WithCommandDescription is a CommandInfoOption that sets the
// description for a ConsoleCommandInfo.
// It splits each line of the description by newline characters.
//
// Parameters:
//
//   - description: The description to set.
//
// Returns:
//
//   - A CommandInfoOption that sets the description for a
//     ConsoleCommandInfo.
func WithCommandDescription(description ...string) CommandInfoOption {
	return func(command *ConsoleCommandInfo) error {
		for _, line := range description {
			fields := strings.Split(line, "\n")
			command.description = append(command.description, fields...)
		}

		return nil
	}
}

// FString generates a formatted string representation of a ConsoleCommandInfo.
// It includes the command name, description, usage information for each flag,
// and the list of flags and their details.
//
// Parameters:
//
//   - indentLevel: The level of indentation to use.
//
// Returns:
//
//   - A string representing the ConsoleCommandInfo.
func (cci *ConsoleCommandInfo) FString(indentLevel int) string {
	indentation := strings.Repeat("\t", indentLevel)
	var builder strings.Builder

	// Add the command name
	fmt.Fprintf(&builder, "%sCommand: %s\n", indentation, cci.name)

	// Add the command description
	if len(cci.description) == 0 {
		fmt.Fprintf(&builder, "%sDescription: [No description provided]\n", indentation)
	} else {
		fmt.Fprintf(&builder, "%sDescription:\n", indentation)
		for _, line := range cci.description {
			fmt.Fprintf(&builder, "%s\t%s\n", indentation, line)
		}
	}

	// Add the usage information for each flag
	for _, flag := range cci.flags {
		fmt.Fprintf(&builder, "%sUsage: %s", indentation, cci.name)

		if flag.required {
			fmt.Fprintf(&builder, " %s", flag.name)
		} else {
			fmt.Fprintf(&builder, " [%s]", flag.name)
		}

		for _, arg := range flag.args {
			fmt.Fprintf(&builder, " <%s>", arg)
		}

		builder.WriteString("\n")
	}

	// Add the flag information
	if len(cci.flags) == 0 {
		fmt.Fprintf(&builder, "%sFlags: None\n", indentation)
	} else {
		fmt.Fprintf(&builder, "%sFlags:\n", indentation)

		for _, flag := range cci.flags {
			fmt.Fprintf(&builder, "%s\n", flag.FString(indentLevel+1))
		}
	}

	return builder.String()
}

// CMLine represents a command line interface. It contains the
// name of the executable and a map of commands that the
// interface accepts.
type CMLine struct {
	// Name of the executable.
	executableName string

	// Description of the executable.
	description []string

	// Map of commands accepted by the interface.
	commands map[string]*ConsoleCommandInfo
}

// CMLineOption is a function type that modifies CMLine.
type CMLineOption func(*CMLine) error

// WithExecutableName is a CMLineOption that sets the executable
// name for a CMLine.
// It trims the space from the name. If the name is empty, it sets
// the errReason of the CMLine.
//
// Parameters:
//
//   - name: The name of the executable.
//
// Returns:
//
//   - A CMLineOption that sets the executable name for a CMLine.
func WithExecutableName(name string) CMLineOption {
	return func(cm *CMLine) error {
		name = strings.TrimSpace(name)
		if name == "" {
			return ers.NewErrInvalidParameter("name").
				WithReason(errors.New("executable name cannot be empty"))
		}

		cm.executableName = name

		return nil
	}
}

// WithCommand is a CMLineOption that adds a new command to a CMLine.
// It trims the space from the name. If the name is empty, it sets
// the errReason of the CMLine.
//
// Parameters:
//
//   - name: The name of the command.
//   - options: The options to apply to the command.
//
// Returns:
//
//   - A CMLineOption that adds a new command to a CMLine.
func WithCommand(name string, options ...CommandInfoOption) CMLineOption {
	return func(cm *CMLine) error {
		name = strings.TrimSpace(name)
		if name == "" {
			return ers.NewErrInvalidParameter("name").
				WithReason(errors.New("name cannot be empty"))
		}

		newCommand := &ConsoleCommandInfo{
			name:        name,
			description: make([]string, 0),
			flags:       make([]*ConsoleFlagInfo, 0),
			callback:    nil,
		}

		for i, option := range options {
			err := option(newCommand)
			if err != nil {
				return fmt.Errorf("invalid option %d for command %s: %v", i, name, err)
			}
		}

		cm.commands[name] = newCommand

		return nil
	}
}

// WithDescription is a CMLineOption that sets the description for
// a CMLine.
// It splits each line of the description by newline characters.
//
// Parameters:
//
//   - description: The description to set.
//
// Returns:
//
//   - A CMLineOption that sets the description for a CMLine.
func WithDescription(description ...string) CMLineOption {
	return func(cm *CMLine) error {
		for _, line := range description {
			fields := strings.Split(line, "\n")
			cm.description = append(cm.description, fields...)
		}

		return nil
	}
}

// NewCMLine creates a new CMLine with the given options. If any
// errors occur while creating the commands, it sets the errReason
// of the CMLine.
//
// Parameters:
//
//   - options: The options to apply to the CMLine.
//
// Returns:
//
//   - A pointer to the created CMLine.
func NewCMLine(options ...CMLineOption) (*CMLine, error) {
	cml := &CMLine{
		commands: make(map[string]*ConsoleCommandInfo),
	}

	for i, option := range options {
		err := option(cml)
		if err != nil {
			return nil, fmt.Errorf("invalid option %d: %v", i, err)
		}
	}

	return cml, nil
}

// ParseCommandLine parses the provided command line arguments
// and executes the corresponding command.
// It returns the name of the executed command, the result of
// the command, and an error, if any.
//
// Parameters:
//
//   - args: The command line arguments to parse.
//
// Returns:
//
//   - The name of the executed command.
//   - The result of the command.
//   - An error, if any.
func (cml *CMLine) ParseCommandLine(args []string) (string, any, error) {
	// Check if any arguments were provided
	if len(args) == 0 {
		return "", nil, ers.NewErrInvalidParameter("args").
			WithReason(errors.New("no arguments provided"))
	}

	// Get the command from the command map
	command, exists := cml.commands[args[0]]
	if !exists {
		return "", nil, fmt.Errorf("command '%s' not found", args[0])
	}

	// Create a map to store the flags
	commandFlags := make(map[string]any)

	if len(args) > 1 {
		// Parse the flags if provided
		var err error

		commandFlags, err = parseConsoleFlags(args[1:], command.flags)
		if err != nil {
			return "", nil, err
		}
	}

	// Check if the command has a callback function
	if command.callback == nil {
		return command.name, nil, nil
	}

	// Call the callback function with the flags
	commandSolution, err := command.callback(commandFlags)
	if err != nil {
		return "", nil, fmt.Errorf("failed to execute command '%s': reason=%v", command.name, err)
	}

	return command.name, commandSolution, nil
}

// parseConsoleFlags parses the provided arguments into console flags.
// It returns a map of the parsed flags and an error, if any.
//
// Parameters:
//
//   - args: The arguments to parse.
//   - flags: The console flags to parse the arguments into.
//
// Returns:
//
//   - A map of the parsed flags.
//   - An error, if any.
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
//   - A string representing the CMLine.
func (cml *CMLine) FString() string {
	var builder strings.Builder

	fmt.Fprintf(&builder, "Usage: %s <command> [flags]\n", cml.executableName)

	if len(cml.description) == 0 {
		fmt.Fprintf(&builder, "Description: [No description provided]\n")
	} else {
		fmt.Fprintf(&builder, "Description:\n")
		for _, line := range cml.description {
			fmt.Fprintf(&builder, "\t%s\n", line)
		}
	}

	if len(cml.commands) == 0 {
		fmt.Fprint(&builder, "Commands: None\n")
	} else {
		fmt.Fprint(&builder, "Commands:\n")

		for _, command := range cml.commands {
			fmt.Fprintf(&builder, "%s\n", command.FString(1))
		}
	}

	return builder.String()
}
