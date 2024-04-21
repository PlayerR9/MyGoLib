// Package CnsPanel provides a structure and functions for handling
// console command flags.
package CnsPanel

import (
	"errors"
	"fmt"
	"strings"

	fs "github.com/PlayerR9/MyGoLib/Formatting/FString"
	ers "github.com/PlayerR9/MyGoLib/Utility/Errors"
)

// ConsoleCommandInfo represents a console command.
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
//
// Parameters:
//
//   - command: The ConsoleCommandInfo to modify.
//
// Returns:
//
//   - error: An error if the modification fails.
type CommandInfoOption func(*ConsoleCommandInfo) error

// WithFlag is a CommandInfoOption that adds a new flag to a
// ConsoleCommandInfo.
// It creates a new ConsoleFlagInfo with the provided name and
// callback, and applies the provided options to it. If the
// flag name is empty or the callback is nil, it returns an error
// of type *ErrInvalidParameter.
//
// Parameters:
//
//   - name: The name of the flag.
//   - callback: The function to call when the flag is used.
//   - options: The options to apply to the flag.
//
// Returns:
//
//   - CommandInfoOption: A CommandInfoOption that adds a new flag to a
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
			return ers.NewErrInvalidParameter("name").Wrap(
				errors.New("flag name cannot be empty"),
			)
		}

		if callback == nil {
			return ers.NewErrNilParameter("callback")
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
// If the provided callback is nil, it returns an error of type
// *ErrInvalidParameter.
//
// Parameters:
//
//   - callback: The function to call when the command is used.
//
// Returns:
//
//   - CommandInfoOption: A CommandInfoOption that sets the callback for
//     a ConsoleCommandInfo.
func WithCallback(callback func(map[string]any) (any, error)) CommandInfoOption {
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
//
//   - description: The description to set.
//
// Returns:
//
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

// FString generates a formatted string representation of a ConsoleCommandInfo.
// It includes the command name, description, usage information for each flag,
// and the list of flags and their details.
//
// Parameters:
//
//   - indentLevel: The level of indentation to use. Sign is ignored.
//
// Returns:
//
//   - []string: A slice of strings representing the ConsoleCommandInfo.
func (cci *ConsoleCommandInfo) FString(indentLevel int) []string {
	indentCfig := fs.NewIndentConfig(fs.DefaultIndentation, 0, true, false)
	indent := indentCfig.String()

	results := make([]string, 0)

	// Add the command name
	results = append(results, fmt.Sprintf("%sCommand: %s", indent, cci.name))

	// Add the command description
	if len(cci.description) == 0 {
		results = append(results, fmt.Sprintf("%sDescription: [No description provided]", indent))
	} else {
		results = append(results, fmt.Sprintf("%sDescription:", indent))

		for _, line := range cci.description {
			results = append(results, fmt.Sprintf("%s\t%s", indent, line))
		}
	}

	// Add the usage information for each flag
	var builder strings.Builder

	for _, flag := range cci.flags {
		fmt.Fprintf(&builder, "%sUsage: %s", indent, cci.name)

		if flag.required {
			fmt.Fprintf(&builder, " %s", flag.name)
		} else {
			fmt.Fprintf(&builder, " [%s]", flag.name)
		}

		for _, arg := range flag.args {
			fmt.Fprintf(&builder, " <%s>", arg)
		}

		results = append(results, builder.String())
		builder.Reset()
	}

	// Add the flag information
	if len(cci.flags) == 0 {
		results = append(results, fmt.Sprintf("%sFlags: None", indent))
	} else {
		results = append(results, fmt.Sprintf("%sFlags:", indent))

		for _, flag := range cci.flags {
			results = append(results, flag.FString(indentLevel+1)...)
		}
	}

	return results
}

type parsedCommand struct {
	command  string
	args     map[string]any
	callback func(map[string]any) (any, error)
}

func (pc *parsedCommand) Command() string {
	return pc.command
}

// ParseCommandLine parses the provided command line arguments
// and executes the corresponding command.
//
// Panics with an error of type *ErrInvalidParameter if no
// arguments are provided, or with an error of type *ErrCallFailed
// if the ParseCommandLine function fails.
//
// Parameters:
//
//   - args: The command line arguments to parse.
//
// Returns:
//
//   - string: The name of the executed command.
//   - any: The result of the command.
func (pc *parsedCommand) Run() (any, error) {
	if pc.callback == nil {
		return nil, nil
	}

	res, err := pc.callback(pc.args)
	if err != nil {
		return nil, err
	}

	return res, nil
}
