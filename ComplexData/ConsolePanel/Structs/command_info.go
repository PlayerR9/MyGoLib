// Package CnsPanel provides a structure and functions for handling
// console command flags.
package Structs

import (
	"errors"
	"strings"

	com "github.com/PlayerR9/MyGoLib/ComplexData/ConsolePanel/Common"
	fs "github.com/PlayerR9/MyGoLib/Formatting/FString"

	res "github.com/PlayerR9/MyGoLib/ComplexData/ConsolePanel/Results"
	cdd "github.com/PlayerR9/MyGoLib/CustomData/Document"

	sm "github.com/PlayerR9/MyGoLib/CustomData/SortedMap"
)

// CommandInfo represents a console command.
type CommandInfo struct {
	// description is the documentation of the command.
	description *cdd.Document

	// flags is a slice of FlagInfo representing the flags accepted by
	// the command.
	flags *sm.SortedMap[string, *FlagInfo]

	// callback is the function to call when the command is executed.
	callback com.CommandCallbackFunc
}

// FString generates a formatted string representation of a CommandInfo.
//
// Parameters:
//   - indentLevel: The level of indentation to use for the CommandInfo.
//
// Returns:
//   - []string: A slice of strings representing the CommandInfo.
//
// Format:
//
//	Description:
//		// <description>
//
//	Flags:
//		- <flag 1>:
//	   	// <description>
//		- <flag 2>:
//	   	// <description>
//		// ...
func (cci *CommandInfo) FString(indentLevel int) []string {
	indentCfig := fs.NewIndentConfig(fs.DefaultIndentation, indentLevel, false)
	indent := indentCfig.String()

	lines := make([]string, 0)

	// Description:
	if cci.description == nil {
		lines = append(lines, "Description: [No description provided]")
	} else {
		lines = append(lines, "Description:")

		descriptionLines := cci.description.FString(indentLevel + 1)
		lines = append(lines, descriptionLines...)
	}

	// Empty line
	lines = append(lines, "")

	// Flags:
	if cci.flags.Size() == 0 {
		lines = append(lines, "Flags: None")
	} else {
		lines = append(lines, "Flags:")

		var builder strings.Builder

		entries := cci.flags.GetEntries()

		builder.WriteString(indent)
		builder.WriteString("- ")
		builder.WriteString(entries[0].First)
		builder.WriteRune(':')

		lines = append(lines, builder.String())

		flagLines := entries[0].Second.FString(indentLevel + 1)
		lines = append(lines, flagLines...)

		for _, p := range entries[1:] {
			builder.Reset()

			builder.WriteString(indent)
			builder.WriteString("- ")
			builder.WriteString(p.First)
			builder.WriteRune(':')

			lines = append(lines, builder.String())

			flagLines := p.Second.FString(indentLevel + 2)
			lines = append(lines, flagLines...)
		}
	}

	// Add the indentation to each line
	for i := 0; i < len(lines); i++ {
		lines[i] = indent + lines[i]
	}

	return lines
}

// NewCommandInfo creates a new CommandInfo with the
// provided command name and callback function.
//
// Parameters:
//   - description: The description of the command.
//   - callback: The function to call when the command is executed.
//
// Returns:
//   - *CommandInfo: A pointer to the new CommandInfo.
//
// Behaviors:
//   - If callback is nil, NoCommandCallback is used.
func NewCommandInfo(description *cdd.Document, callback com.CommandCallbackFunc) *CommandInfo {
	inf := &CommandInfo{
		description: description,
		flags:       sm.NewSortedMap[string, *FlagInfo](),
	}

	if callback != nil {
		inf.callback = callback
	} else {
		inf.callback = com.NoCommandCallback
	}

	return inf
}

// AddFlag adds a new flag to a CommandInfo.
//
// Parameters:
//   - flag: The flag to add.
//   - info: The FlagInfo for the flag.
//
// Returns:
//   - *CommandInfo: A pointer to the CommandInfo. This allows for chaining.
//
// Behaviors:
//   - If info is nil, the flag is not added.
//   - If the flag already exists, the existing flag is replaced with the new one.
func (ci *CommandInfo) AddFlag(flag string, info *FlagInfo) *CommandInfo {
	if info == nil {
		return ci
	}

	ci.flags.AddEntry(flag, info)

	return ci
}

// ParseArguments parses the provided command line arguments
// and returns a ParsedCommand ready to be executed.
//
// Errors:
//   - *ers.ErrInvalidParameter: No arguments provided.
//   - *ErrCommandNotFound: Command not found.
//   - *ErrParsingFlags: Error parsing flags.
//
// Parameters:
//   - args: The command line arguments to parse. Without the
//     executable name.
//
// Returns:
//   - *ParsedCommand: A pointer to the parsed command.
//   - error: An error, if any.
func (inf *CommandInfo) Parse(args []string) (*res.ParsedCommand, error) {
	if inf.flags.Size() == 0 {
		if len(args) != 0 {
			return nil, errors.New("no flags expected")
		}

		command, err := res.NewParsedCommand(make(com.ArgumentsMap), inf.callback)
		if err != nil {
			panic(err)
		}

		return command, nil
	} else if len(args) == 0 {
		return nil, errors.New("no arguments provided")
	}

	argsMap := make(com.ArgumentsMap)
	todo := inf.flags.Copy().(*sm.SortedMap[string, *FlagInfo])

	for len(args) > 0 {
		if todo.Size() == 0 {
			return nil, errors.New("too many arguments provided")
		}

		flag, err := todo.GetEntry(args[0])
		if err != nil {
			return nil, err
		}

		parsed, err := flag.Parse(args[1:])
		if err != nil {
			return nil, err
		}

		argsMap[args[0]] = parsed.Args
		args = args[parsed.Index:]

		todo.Delete(args[0])
	}

	if todo.Size() > 0 {
		for _, p := range todo.GetEntries() {
			if p.Second.IsRequired() {
				return nil, errors.New("missing required flag")
			}
		}
	}

	command, err := res.NewParsedCommand(argsMap, inf.callback)
	if err != nil {
		panic(err)
	}

	return command, nil
}

// GetDescription returns the description of a CommandInfo.
//
// Returns:
//   - *cdd.Document: The description of the CommandInfo.
func (inf *CommandInfo) GetDescription() *cdd.Document {
	return inf.description
}

// GetFlags returns the flags of a CommandInfo.
//
// Returns:
//   - []*FlagInfo: The flags of the CommandInfo.
//
// Behaviors:
//   - Modifying the returned slice will affect the CommandInfo.
func (inf *CommandInfo) GetFlags() []*FlagInfo {
	return inf.flags.Values()
}

// GetCallback returns the callback function of a CommandInfo.
//
// Returns:
//   - CommandCallbackFunc: The callback function of the CommandInfo.
//
// Behaviors:
//   - Never returns nil.
func (inf *CommandInfo) GetCallback() com.CommandCallbackFunc {
	return inf.callback
}
