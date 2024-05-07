// Package CnsPanel provides a structure and functions for handling
// console command flags.
package CnsPanel

import (
	"slices"
	"strings"

	fs "github.com/PlayerR9/MyGoLib/Formatting/FString"
)

// CommandInfo represents a console command.
type CommandInfo struct {
	// Name of the command.
	name string

	// Brief explanation of what the command does.
	description *Description

	// Slice of FlagInfo representing the flags accepted by
	// the command.
	flags []*FlagInfo

	// Function invoked when the command is executed.
	callback CommandCallbackFunc
}

// FString generates a formatted string representation of a CommandInfo.
// It includes the command name, description, usage information for each flag,
// and the list of flags and their details.
//
// Parameters:
//   - indentLevel: The level of indentation to use. Sign is ignored.
//
// Returns:
//   - []string: A slice of strings representing the CommandInfo.
func (cci *CommandInfo) FString(indentLevel int) []string {
	indentCfig := fs.NewIndentConfig(fs.DefaultIndentation, indentLevel, false)
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

	if cci.description == nil {
		builder.WriteString("Description: [No description provided]")
		results = append(results, builder.String())
	} else {
		builder.WriteString("Description:")
		results = append(results, builder.String())

		lines := cci.description.FString(indentLevel + 1)
		results = append(results, lines...)
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

// NewCommandInfo creates a new CommandInfo with the
// provided name and callback function.
//
// Parameters:
//   - commandName: The name of the command.
//   - callback: The function to call when the command is used.
//
// Returns:
//   - *CommandInfo: A pointer to the new CommandInfo.
func NewCommandInfo(commandName string, callback CommandCallbackFunc) *CommandInfo {
	return &CommandInfo{
		name:        commandName,
		description: nil,
		flags:       make([]*FlagInfo, 0),
		callback:    callback,
	}
}

// AddFlags is a method of CommandInfo that adds a new flag to a
// CommandInfo.
//
// Nil flags are ignored.
//
// Parameters:
//   - flags: The flags to add.
func (ci *CommandInfo) AddFlags(flags ...*FlagInfo) {
	if len(flags) == 0 {
		return
	}

	for _, flag := range flags {
		if flag == nil {
			continue
		}

		index := slices.IndexFunc(ci.flags, func(f *FlagInfo) bool {
			return f.Equals(flag)
		})
		if index != -1 {
			ci.flags[index] = flag
		} else {
			ci.flags = append(ci.flags, flag)
		}
	}
}

// SetDescription sets the description for a CommandInfo.
// It splits each line of the description by newline characters.
//
// Parameters:
//   - description: The description to set.
func (ci *CommandInfo) SetDescription(description *Description) {
	ci.description = description
}
