// Package CnsPanel provides a structure and functions for handling
// console command flags.
package CmdLineParser

import (
	"fmt"
	"slices"

	"github.com/PlayerR9/MyGoLib/ComplexData/CmdLineParser/pkg"
	evalSlc "github.com/PlayerR9/MyGoLib/Evaluations/Slices"
	ffs "github.com/PlayerR9/MyGoLib/Formatting/FString"
	ue "github.com/PlayerR9/MyGoLib/Units/Errors"
)

const (
	// DefaultWidth is the default width of the console.
	DefaultWidth int = 80

	// DefaultHeight is the default height of the console.
	DefaultHeight int = 24
)

// CmdLineParser represents a command line console.
type CmdLineParser struct {
	// ExecutableName is the name of the executable.
	ExecutableName string

	// description is the documentation of the executable.
	description []string

	// commandList is a map of command opcodes to CommandInfo.
	commandList []*pkg.CommandInfo

	// width and height are the dimensions of the console, respectively.
	width, height int
}

// FString generates a formatted string representation of a CmdLineParser.
//
// Parameters:
//   - indentLevel: The level of indentation to use for the CmdLineParser.
//
// Returns:
//   - []string: A slice of strings representing the CmdLineParser.
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
func (cns *CmdLineParser) FString(trav *ffs.Traversor) error {
	// Usage:
	err := trav.AddJoinedLine(" ", "Usage:", cns.ExecutableName, "<command>", "[flags]")
	if err != nil {
		return err
	}

	// Empty line
	trav.EmptyLine()

	err = trav.AppendString("Description:")
	if err != nil {
		return err
	}

	if len(cns.description) == 0 {
		err = trav.AppendRune(' ')
		if err != nil {
			return err
		}

		err := trav.AppendString("[No description provided]")
		if err != nil {
			return err
		}

		trav.AcceptLine()
	} else {
		trav.AcceptLine()

		err = ffs.ApplyForm(
			trav.GetConfig(
				ffs.WithIncreasedIndent(),
			),
			trav,
			pkg.NewDescriptionPrinter(cns.description),
		)
		if err != nil {
			return err
		}
	}

	/*
		// Description:
		doc := fsd.NewDocumentPrinter("Description", cns.description, )
		err = doc.FString(trav)
		if err != nil {
			return ue.NewErrWhile("FString printing description", err)
		}
	*/

	// Empty line
	trav.EmptyLine()

	// Commands:
	if len(cns.commandList) == 0 {
		err = trav.AddLine("Commands: No commands provided.")
		if err != nil {
			return err
		}
	} else {
		err = trav.AddLine("Commands:")
		if err != nil {
			return err
		}

		for at, command := range cns.commandList {
			err := ffs.ApplyForm(
				trav.GetConfig(
					ffs.WithIncreasedIndent(),
				),
				trav,
				command,
			)
			if err != nil {
				return ue.NewErrAt(at+1, "command", err)
			}
		}
	}

	return nil
}

// NewCmdLineParser creates a new CmdLineParser with the provided executable name.
//
// Parameters:
//   - execName: The name of the executable.
//
// Returns:
//   - *CmdLineParser: A pointer to the created CmdLineParser.
//
// addCommand adds the provided command to the CmdLineParser.
//
// Parameters:
//   - opcode: The command opcode.
//   - info: The CommandInfo for the command.
//
// Returns:
//   - *CmdLineParser: A pointer to the CmdLineParser. This allows for chaining.
//
// Behaviors:
//   - If opcode is either an empty string or "help", the command is not added.
//   - If info is nil, the command is not added.
//   - If the opcode already exists, the existing command is replaced with the new one.
func NewCmdLineParser(execName string, description []string, commandBuilder *CmdBuilder) (*CmdLineParser, error) {
	cp := &CmdLineParser{
		ExecutableName: execName,
		description:    description,
		width:          DefaultWidth,
		height:         DefaultHeight,
		commandList:    make([]*pkg.CommandInfo, 0),
	}

	f := func(args map[string]any) (any, error) {
		if len(args) == 0 {
			doc, err := ffs.SprintFString(ffs.DefaultFormatter, cp)
			if err != nil {
				return nil, fmt.Errorf("error printing console panel: %w", err)
			}

			return ffs.Stringfy(doc), nil
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

			doc, err := ffs.SprintFString(ffs.DefaultFormatter, command)
			if err != nil {
				return nil, fmt.Errorf("error printing command %q: %w", opcode, err)
			}

			return ffs.Stringfy(doc), nil

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
	doc, err := ffs.Sprintln(ffs.DefaultFormatter, "Displays help information for the console.")
	if err != nil {
		return nil, fmt.Errorf("error generating help doc: %w", err)
	}

	helpCommandInfo, err := pkg.NewCommandInfo(
		HelpOpcode,
		ffs.Stringfy(doc),
		f,
		make([]*pkg.FlagInfo, 0),
	)
	if err != nil {
		return nil, fmt.Errorf("error making help command: %w", err)
	}

	cp.commandList = append(cp.commandList, helpCommandInfo)

	if commandBuilder == nil {
		return cp, nil
	}

	commandList, err := commandBuilder.Build()
	if err != nil {
		return nil, fmt.Errorf("error building command: %w", err)
	}

	cp.commandList = commandList

	helpCommand, err := generateHelpCommand(cp)
	if err != nil {
		return nil, fmt.Errorf("error making help command: %w", err)
	}

	cp.commandList = append(cp.commandList, helpCommand)

	return cp, nil
}

// GetCommandByOpcode returns the CommandInfo for the provided opcode.
//
// Parameters:
//   - opcode: The opcode of the command.
//
// Returns:
//   - *CommandInfo: The CommandInfo for the opcode. Nil if not found.
func (cns *CmdLineParser) GetCommandByOpcode(opcode string) *pkg.CommandInfo {
	for _, command := range cns.commandList {
		if command.GetOpcode() == opcode {
			return command
		}
	}

	return nil
}

// Parse parses the provided command line arguments
// and returns a ParsedCommand ready to be executed.
//
// Errors:
//   - *ue.ErrInvalidParameter: No arguments provided.
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
func (cns *CmdLineParser) Parse(args []string) (*pkg.ParsedCommand, error) {
	if len(args) == 0 {
		return nil, ue.NewErrInvalidParameter("args", ue.NewErrEmpty(args))
	}

	command := cns.GetCommandByOpcode(args[0])
	if command == nil {
		return nil, NewErrCommandNotFound(args[0])
	}

	args = args[1:] // Remove the command name

	branches, err := evalSlc.Evaluate(command, args)
	if err != nil {
		return nil, err
	}

	pc, err := command.Parse(branches, args)
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
func (cns *CmdLineParser) GetCommand(opcode string) (*pkg.CommandInfo, bool) {
	index := slices.IndexFunc(cns.commandList, func(c *pkg.CommandInfo) bool {
		return c.GetOpcode() == opcode
	})
	if index == -1 {
		return nil, false
	}

	return cns.commandList[index], true
}

// SetDimensions sets the dimensions of the console.
//
// Parameters:
//   - width: The width of the console.
//   - height: The height of the console.
//
// Returns:
//   - error: An error of type *ue.ErrInvalidParameter if width or height is less than
//     or equal to 0.
func (cns *CmdLineParser) SetDimensions(width, height int) error {
	if width <= 0 {
		return ue.NewErrInvalidParameter("width", ue.NewErrGT(0))
	} else if height <= 0 {
		return ue.NewErrInvalidParameter("height", ue.NewErrGT(0))
	}

	cns.width = width
	cns.height = height

	return nil
}
