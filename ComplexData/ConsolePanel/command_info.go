// Package CnsPanel provides a structure and functions for handling
// console command flags.
package ConsolePanel

import (
	"errors"
	"fmt"

	fss "github.com/PlayerR9/MyGoLib/Formatting/FString"

	sm "github.com/PlayerR9/MyGoLib/CustomData/OrderedMap"

	ers "github.com/PlayerR9/MyGoLib/Units/Errors"
)

// CommandCallbackFunc is a function type that represents a callback
// function for a console command.
//
// Parameters:
//   - args: A map of string keys to any values representing the
//     arguments passed to the command.
//
// Returns:
//   - error: An error if the command fails.
//   - any: The result of the command. (if any)
type CommandCallbackFunc func(args map[string]any) (any, error)

// NoCommandCallback is a default callback function for a console command
// when no callback is provided.
//
// Parameters:
//   - args: A map of string keys to any values representing the
//     arguments passed to the command.
//
// Returns:
//   - error: nil
//   - any: nil
func NoCommandCallback(args map[string]any) (any, error) {
	return nil, nil
}

// CommandInfo represents a console command.
type CommandInfo struct {
	// name is the name of the command.
	name string

	// description is the documentation of the command.
	description []string

	// flags is a slice of FlagInfo representing the flags accepted by
	// the command.
	flags *sm.OrderedMap[string, *FlagInfo]

	// callback is the function to call when the command is executed.
	callback CommandCallbackFunc
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
func (cci *CommandInfo) FString(trav *fss.Traversor) error {
	if trav == nil {
		return nil
	}

	// Description:
	err := trav.AppendString("Description:")
	if err != nil {
		return err
	}

	if len(cci.description) == 0 {
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

		err = fss.ApplyForm(
			trav.GetConfig(
				fss.WithIncreasedIndent(),
			),
			trav,
			&descriptionPrinter{cci.description},
		)
		if err != nil {
			return err
		}
	}

	// Empty line
	trav.EmptyLine()

	// Flags:
	printer := sm.NewOrderedMapPrinter("Flags", cci.flags, "flag", "no flags")
	err = printer.FString(trav)
	if err != nil {
		return ers.NewErrWhile("FString printing flags", err)
	}

	return nil
}

// NewCommandInfo creates a new CommandInfo with the
// provided command name and callback function.
//
// Parameters:
//   - name: The name of the command.
//   - description: The description of the command.
//   - callback: The function to call when the command is executed.
//
// Returns:
//   - *CommandInfo: A pointer to the new CommandInfo.
//
// Behaviors:
//   - If callback is nil, NoCommandCallback is used.
func NewCommandInfo(name string, description []string, callback CommandCallbackFunc) *CommandInfo {
	inf := &CommandInfo{
		name:        name,
		description: description,
		flags:       sm.NewOrderedMap[string, *FlagInfo](),
	}

	if callback != nil {
		inf.callback = callback
	} else {
		inf.callback = NoCommandCallback
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
func (inf *CommandInfo) Parse(args []string) (*ParsedCommand, error) {
	if inf.flags.Size() == 0 {
		if err := ers.NewErrEmpty(args); err == nil {
			return nil, errors.New("no flags expected")
		}

		command := newParsedCommand(inf.name, make(map[string]any), inf.callback)

		return command, nil
	}

	argsMap := make(map[string]any)
	todo := inf.flags.Copy().(*sm.OrderedMap[string, *FlagInfo])

	for len(args) > 0 {
		if todo.Size() == 0 {
			return nil, errors.New("too many arguments provided")
		}

		flag, ok := todo.GetEntry(args[0])
		if !ok {
			return nil, fmt.Errorf("invalid flag %q", args[0])
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
		iter := todo.Iterator()

		for {
			entry, err := iter.Consume()
			if err != nil {
				break
			}

			if entry.Second.IsRequired() {
				return nil, errors.New("missing required flag")
			}
		}
	}

	command := newParsedCommand(inf.name, argsMap, inf.callback)

	return command, nil
}

// GetDescription returns the description of a CommandInfo.
//
// Returns:
//   - *fsd.Document: The description of the CommandInfo.
func (inf *CommandInfo) GetDescription() []string {
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
func (inf *CommandInfo) GetCallback() CommandCallbackFunc {
	return inf.callback
}
