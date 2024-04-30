// Package CnsPanel provides a structure and functions for handling
// console command flags.
package CnsPanel

import (
	"strings"

	fs "github.com/PlayerR9/MyGoLib/Formatting/FString"
	ers "github.com/PlayerR9/MyGoLib/Units/Errors"
)

// CommandInfoOption is a function type that modifies
// ConsoleCommandInfo.
//
// Parameters:
//   - command: The ConsoleCommandInfo to modify.
//
// Returns:
//   - error: An error if the modification fails.
type CommandInfoOption func(*ConsoleCommandInfo) error

// WithFlag is a CommandInfoOption that adds a new flag to a
// ConsoleCommandInfo.
// It creates a new ConsoleFlagInfo with the provided name and
// callback, and applies the provided options to it. If the
// flag name is empty or the callback is nil, it returns an error
// of type *ers.ErrInvalidParameter.
//
// Parameters:
//   - name: The name of the flag.
//   - callback: The function to call when the flag is used.
//   - options: The options to apply to the flag.
//
// Returns:
//   - CommandInfoOption: A CommandInfoOption that adds a new flag to a
//     ConsoleCommandInfo.
func WithFlag(name string, callback FlagCallbackFunc, options ...FlagInfoOption) CommandInfoOption {
	return func(command *ConsoleCommandInfo) error {
		newFlag := ConsoleFlagInfo{
			name:        name,
			args:        make([]string, 0),
			description: make([]string, 0),
			required:    false,
			callback:    callback,
		}

		name = strings.TrimSpace(name)
		if name == "" {
			return ers.NewErrInvalidParameter(
				"name",
				ers.NewErrEmptyString(),
			)
		}

		if callback == nil {
			return ers.NewErrNilParameter("callback")
		}

		for _, option := range options {
			option(&newFlag)
		}

		command.flags = append(command.flags, newFlag)

		return nil
	}
}

// WithCallback is a CommandInfoOption that sets the callback for
// a ConsoleCommandInfo.
// If the provided callback is nil, it returns an error of type
// *ers.ErrInvalidParameter.
//
// Parameters:
//   - callback: The function to call when the command is used.
//
// Returns:
//   - CommandInfoOption: A CommandInfoOption that sets the callback for
//     a ConsoleCommandInfo.
func WithCallback(callback CommandCallbackFunc) CommandInfoOption {
	return func(command *ConsoleCommandInfo) error {
		if callback == nil {
			return ers.NewErrNilParameter("callback")
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
//   - description: The description to set.
//
// Returns:
//   - CommandInfoOption: A CommandInfoOption that sets the description for
//     a ConsoleCommandInfo.
func WithCommandDescription(description ...string) CommandInfoOption {
	return func(command *ConsoleCommandInfo) error {
		for _, line := range description {
			fields := strings.Split(line, "\n")
			command.description = append(command.description, fields...)
		}

		return nil
	}
}

// ConsoleCommandInfo represents a console command.
type ConsoleCommandInfo struct {
	// Name of the command.
	name string

	// Brief explanation of what the command does.
	description []string

	// Slice of ConsoleFlagInfo representing the flags accepted by
	// the command.
	flags []ConsoleFlagInfo

	// Function invoked when the command is executed.
	callback CommandCallbackFunc
}

// FString generates a formatted string representation of a ConsoleCommandInfo.
// It includes the command name, description, usage information for each flag,
// and the list of flags and their details.
//
// Parameters:
//   - indentLevel: The level of indentation to use. Sign is ignored.
//
// Returns:
//   - []string: A slice of strings representing the ConsoleCommandInfo.
func (cci *ConsoleCommandInfo) FString(indentLevel int) []string {
	indentCfig := fs.NewIndentConfig(fs.DefaultIndentation, 0, true, false)
	indent := indentCfig.String()

	results := make([]string, 0)
	var builder strings.Builder

	// Add the command name
	builder.WriteString(indent)
	builder.WriteString("Command: ")
	builder.WriteString(cci.name)

	results = append(results, builder.String())

	// Add the command description
	builder.Reset()

	builder.WriteString(indent)
	builder.WriteString("Description:")

	if len(cci.description) == 0 {
		builder.WriteString(" [No description provided]")
		results = append(results, builder.String())
	} else {
		results = append(results, builder.String())

		for _, line := range cci.description {
			builder.Reset()

			builder.WriteString(indent)
			builder.WriteString(indent)
			builder.WriteString(line)

			results = append(results, builder.String())
		}
	}

	// Add the usage information for each flag
	for _, flag := range cci.flags {
		builder.Reset()

		builder.WriteString(indent)
		builder.WriteString("Usage: ")
		builder.WriteString(cci.name)

		builder.WriteRune(' ')

		if flag.required {
			builder.WriteString(flag.name)
		} else {
			builder.WriteRune('[')
			builder.WriteString(flag.name)
			builder.WriteRune(']')
		}

		for _, arg := range flag.args {
			builder.WriteRune(' ')
			builder.WriteRune('<')
			builder.WriteString(arg)
			builder.WriteRune('>')
		}

		results = append(results, builder.String())
	}

	// Add the flag information
	builder.Reset()

	builder.WriteString(indent)
	builder.WriteString("Flags:")

	if len(cci.flags) == 0 {
		builder.WriteString(" None")
		results = append(results, builder.String())
	} else {
		results = append(results, builder.String())

		for _, flag := range cci.flags {
			results = append(results, flag.FString(indentLevel+1)...)
		}
	}

	return results
}

// NewConsoleCommandInfo creates a new ConsoleCommandInfo with the
// provided name and options.
//
// Parameters:
//   - commandName: The name of the command.
//   - options: The options to apply to the command.
//
// Returns:
//   - ConsoleCommandInfo: The new ConsoleCommandInfo.
func NewConsoleCommandInfo(commandName string, options ...CommandInfoOption) (ConsoleCommandInfo, error) {
	newCommand := ConsoleCommandInfo{
		name:        commandName,
		description: make([]string, 0),
		flags:       make([]ConsoleFlagInfo, 0),
		callback:    nil,
	}

	for i, option := range options {
		err := option(&newCommand)
		if err != nil {
			return newCommand, ers.NewErrAt(i, err)
		}
	}

	return newCommand, nil
}
