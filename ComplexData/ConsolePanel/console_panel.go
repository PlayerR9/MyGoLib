// Package CnsPanel provides a structure and functions for handling
// console command flags.
package ConsolePanel

import (
	"errors"
	"fmt"
	"strings"

	cdd "github.com/PlayerR9/MyGoLib/CustomData/Document"
	fs "github.com/PlayerR9/MyGoLib/Formatting/FString"

	sm "github.com/PlayerR9/MyGoLib/CustomData/SortedMap"

	com "github.com/PlayerR9/MyGoLib/ComplexData/ConsolePanel/Common"
	res "github.com/PlayerR9/MyGoLib/ComplexData/ConsolePanel/Results"
	"github.com/PlayerR9/MyGoLib/ComplexData/ConsolePanel/Structs"
)

// ConsolePanel represents a command line console.
type ConsolePanel struct {
	// ExecutableName is the name of the executable.
	ExecutableName string

	// description is the documentation of the executable.
	description *cdd.Document

	// commandMap is a map of command opcodes to CommandInfo.
	commandMap *sm.SortedMap[string, *Structs.CommandInfo]
}

// FString generates a formatted string representation of a ConsolePanel.
//
// Parameters:
//   - indentLevel: The level of indentation to use for the ConsolePanel.
//
// Returns:
//   - []string: A slice of strings representing the ConsolePanel.
//
// Format:
//
//	Usage: <executable name> <commands> [flags]
//
//	Description:
//		// <description>
//
//	Commands:
//		- <command 1>:
//	   	// <description>
//		- <command 2>:
//	   	// <description>
//		// ...
func (cns *ConsolePanel) FString(indentLevel int) []string {
	lines := make([]string, 0)
	var builder strings.Builder
	indentCfig := fs.NewIndentConfig(fs.DefaultIndentation, indentLevel, false)
	indent := indentCfig.String()

	// Usage:
	builder.WriteString("Usage: ")
	builder.WriteString(cns.ExecutableName)
	builder.WriteString(" <command> [flags]")

	lines = append(lines, builder.String())

	// Empty line
	lines = append(lines, "")

	// Description:
	if cns.description == nil {
		lines = append(lines, "Description: [No description provided]")
	} else {
		lines = append(lines, "Description:")

		descriptionLines := cns.description.FString(indentLevel + 1)
		lines = append(lines, descriptionLines...)
	}

	// Empty line
	lines = append(lines, "")

	// Commands:
	if cns.commandMap.Size() == 0 {
		lines = append(lines, "Commands: None")
	} else {
		lines = append(lines, "Commands:")

		commands := cns.commandMap.GetEntries()

		for _, command := range commands {
			builder.Reset()

			builder.WriteString(indent)
			builder.WriteString("- ")
			builder.WriteString(command.First)
			builder.WriteRune(':')

			lines = append(lines, builder.String())

			commandLines := command.Second.FString(indentLevel + 2)
			lines = append(lines, commandLines...)
		}
	}

	// Add the indentation to each line
	for i := 0; i < len(lines); i++ {
		lines[i] = indent + lines[i]
	}

	return lines
}

// NewConsolePanel creates a new ConsolePanel with the provided executable name.
//
// Parameters:
//   - execName: The name of the executable.
//
// Returns:
//   - *ConsolePanel: A pointer to the created ConsolePanel.
func NewConsolePanel(execName string, description *cdd.Document) *ConsolePanel {
	cp := &ConsolePanel{
		ExecutableName: execName,
		description:    description,
		commandMap:     sm.NewSortedMap[string, *Structs.CommandInfo](),
	}

	// Add the help command
	helpCommandInfo := Structs.NewCommandInfo(
		cdd.NewDocument("Displays help information for the console."),
		func(args com.ArgumentsMap) error {
			lines := cp.FString(0)

			for _, line := range lines {
				fmt.Println(line)
			}

			return nil
		},
	)

	cp.AddCommand(
		"help",
		helpCommandInfo,
	)

	return cp
}

// AddCommand adds the provided command to the ConsolePanel.
//
// Parameters:
//   - opcode: The command opcode.
//   - info: The CommandInfo for the command.
//
// Returns:
//   - *ConsolePanel: A pointer to the ConsolePanel. This allows for chaining.
//
// Behaviors:
//   - If opcode is either an empty string or "help", the command is not added.
//   - If info is nil, the command is not added.
//   - If the opcode already exists, the existing command is replaced with the new one.
func (cp *ConsolePanel) AddCommand(opcode string, info *Structs.CommandInfo) *ConsolePanel {
	if info == nil || opcode == "" || opcode == "help" {
		return cp
	}

	cp.commandMap.AddEntry(opcode, info)

	addSpecificHelp := func(info *Structs.CommandInfo) (*Structs.CommandInfo, error) {
		newInfo := Structs.NewFlagInfo(false, nil, nil).SetDescription(
			cdd.NewDocument(
				fmt.Sprintf("Displays the help information for the %q command.", opcode),
			),
		)

		return info.AddFlag(opcode, newInfo), nil
	}

	err := cp.commandMap.ModifyValueFunc("help", addSpecificHelp)
	if err != nil {
		panic(err)
	}

	return cp
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
func (cns *ConsolePanel) ParseArguments(args []string) (*res.ParsedCommand, error) {
	if len(args) == 0 {
		return nil, errors.New("missing command name")
	}

	command, err := cns.commandMap.GetEntry(args[0])
	if err != nil {
		return nil, err
	}

	args = args[1:] // Remove the command name

	var pc *res.ParsedCommand

	if len(args) == 0 {
		pc, err = res.NewParsedCommand(make(com.ArgumentsMap), command.GetCallback())
		if err != nil {
			panic(err)
		}
	} else {
		pc, err = command.Parse(args)
		if err != nil {
			return nil, err
		}
	}

	return pc, nil
}

// GetCommand returns the CommandInfo for the provided opcode.
//
// Parameters:
//   - opcode: The opcode of the command.
//
// Returns:
//   - *CommandInfo: The CommandInfo for the opcode.
//   - bool: A boolean indicating if the command was found.
func (cns *ConsolePanel) GetCommand(opcode string) (*Structs.CommandInfo, error) {
	return cns.commandMap.GetEntry(opcode)
}
