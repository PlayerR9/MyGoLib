// Package CnsPanel provides a structure and functions for handling
// console command flags.
package ConsolePanel

import (
	"errors"
	"fmt"

	fsd "github.com/PlayerR9/MyGoLib/FString/Document"
	fsp "github.com/PlayerR9/MyGoLib/FString/Printer"

	sm "github.com/PlayerR9/MyGoLib/CustomData/OrderedMap"
	ers "github.com/PlayerR9/MyGoLib/Units/Errors"
)

const (
	// DefaultWidth is the default width of the console.
	DefaultWidth int = 80

	// DefaultHeight is the default height of the console.
	DefaultHeight int = 24
)

// ConsolePanel represents a command line console.
type ConsolePanel struct {
	// ExecutableName is the name of the executable.
	ExecutableName string

	// description is the documentation of the executable.
	description *fsd.Document

	// commandMap is a map of command opcodes to CommandInfo.
	commandMap *sm.OrderedMap[string, *CommandInfo]

	// width and height are the dimensions of the console, respectively.
	width, height int
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
func (cns *ConsolePanel) FString(trav *fsp.Traversor) error {
	// Usage:
	err := trav.AddJoinedLine(" ", "Usage:", cns.ExecutableName, "<command>", "[flags]")
	if err != nil {
		return err
	}

	// Empty line
	trav.EmptyLine()

	/*
		// Description:
		doc := fsd.NewDocumentPrinter("Description", cns.description, "[No description provided]")
		err = doc.FString(trav)
		if err != nil {
			return ers.NewErrWhile("FString printing description", err)
		}
	*/

	// Empty line
	trav.EmptyLine()

	// Commands:
	printer := sm.NewOrderedMapPrinter("Commands", cns.commandMap, "command", "no commands")
	err = printer.FString(trav)
	if err != nil {
		return ers.NewErrWhile("FString printing commands", err)
	}

	return nil
}

// NewConsolePanel creates a new ConsolePanel with the provided executable name.
//
// Parameters:
//   - execName: The name of the executable.
//
// Returns:
//   - *ConsolePanel: A pointer to the created ConsolePanel.
func NewConsolePanel(execName string, description *fsd.Document) *ConsolePanel {
	cp := &ConsolePanel{
		ExecutableName: execName,
		description:    description,
		commandMap:     sm.NewOrderedMap[string, *CommandInfo](),
		width:          DefaultWidth,
		height:         DefaultHeight,
	}

	f := func(args map[string]any) (any, error) {
		if len(args) == 0 {
			printer := fsp.NewPrinter(fsp.DefaultFormatter)

			err := fsp.ApplyFormat(printer, cp)
			if err != nil {
				return nil, fmt.Errorf("error printing console panel: %w", err)
			}

			doc, err := printer.MakeDocument()
			if err != nil {
				return nil, fmt.Errorf("error making document: %w", err)
			}

			return doc, nil
			/*
				runeTable := make([][]rune, 0)

				for _, line := range lines {
					rt, err := line.Runes(cp.width, cp.height)
					if err != nil {
						return nil, err
					}

					runeTable = append(runeTable, rt...)
				}

				return runeTable, nil
			*/
		}

		runeTable := make([][]rune, 0)

		for opcode := range args {
			command, ok := cp.GetCommand(opcode)
			if !ok {
				return nil, NewErrCommandNotFound(opcode)
			}

			printer := fsp.NewPrinter(fsp.DefaultFormatter)

			err := fsp.ApplyFormat(printer, command)
			if err != nil {
				return nil, fmt.Errorf("error printing command %q: %w", opcode, err)
			}

			doc, err := printer.MakeDocument()
			if err != nil {
				return nil, fmt.Errorf("error making document: %w", err)
			}

			return doc, nil

			/*
				mlts := trav.GetLines()

				for _, mlt := range mlts {
					rt, err := mlt.Runes(cp.width, cp.height)
					if err != nil {
						return nil, err
					}

					runeTable = append(runeTable, rt...)
				}
			*/
		}

		return runeTable, nil
	}

	// Add the help command
	helpCommandInfo := NewCommandInfo(
		fsd.NewDocument("Displays help information for the console."),
		f,
	)

	cp.commandMap.AddEntry("help", helpCommandInfo)

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
func (cp *ConsolePanel) AddCommand(opcode string, info *CommandInfo) *ConsolePanel {
	if info == nil || opcode == "" || opcode == "help" {
		return cp
	}

	cp.commandMap.AddEntry(opcode, info)

	addSpecificHelp := func(info *CommandInfo) (*CommandInfo, error) {
		newInfo := NewFlagInfo(false, nil, nil).SetDescription(
			fsd.NewDocument(
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
func (cns *ConsolePanel) ParseArguments(args []string) (*ParsedCommand, error) {
	if len(args) == 0 {
		return nil, errors.New("missing command name")
	}

	command, ok := cns.commandMap.GetEntry(args[0])
	if !ok {
		return nil, fmt.Errorf("command %q not found", args[0])
	}

	args = args[1:] // Remove the command name

	pc, err := command.Parse(args)
	if err != nil {
		return nil, err
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
func (cns *ConsolePanel) GetCommand(opcode string) (*CommandInfo, bool) {
	return cns.commandMap.GetEntry(opcode)
}

// SetDimensions sets the dimensions of the console.
//
// Parameters:
//   - width: The width of the console.
//   - height: The height of the console.
//
// Returns:
//   - error: An error of type *ers.ErrInvalidParameter if width or height is less than
//     or equal to 0.
func (cns *ConsolePanel) SetDimensions(width, height int) error {
	if width <= 0 {
		return ers.NewErrInvalidParameter("width", ers.NewErrGT(0))
	} else if height <= 0 {
		return ers.NewErrInvalidParameter("height", ers.NewErrGT(0))
	}

	cns.width = width
	cns.height = height

	return nil
}
