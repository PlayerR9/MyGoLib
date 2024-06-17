// Package CnsPanel provides a structure and functions for handling
// console command flags.
package pkg

import (
	"errors"

	evalSlc "github.com/PlayerR9/MyGoLib/Evaluations/Slices"
	ffs "github.com/PlayerR9/MyGoLib/Formatting/FString"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
	ue "github.com/PlayerR9/MyGoLib/Units/errors"
	us "github.com/PlayerR9/MyGoLib/Units/slice"
	uts "github.com/PlayerR9/MyGoLib/Utility/Sorting"
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
type CommandCallbackFunc func(args map[string]map[string][]any) (any, error)

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
func NoCommandCallback(args map[string]map[string][]any) (any, error) {
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
	flags []*FlagInfo

	// callback is the function to call when the command is executed.
	callback CommandCallbackFunc
}

// Evaluator implements the Evaluable interface.
func (inf *CommandInfo) Evaluator() evalSlc.LeafEvaluater[string, *resultBranch, int, []*FlagParseResult] {
	return &ciEvaluator{
		flags: inf.flags,
	}
}

// Equals checks if the command is equal to another command.
//
// Two commands are equal iff their names are equal.
//
// Parameters:
//   - other: The other command to compare.
//
// Returns:
//   - bool: true if the commands are equal, false otherwise.
func (ci *CommandInfo) Equals(other uc.Equaler) bool {
	if other == nil {
		return false
	}

	otherC, ok := other.(*CommandInfo)
	if !ok {
		return false
	}

	return ci.name == otherC.name
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
func (cci *CommandInfo) FString(trav *ffs.Traversor, opts ...ffs.Option) error {
	if trav == nil {
		return nil
	}

	// Name:
	err := trav.AddJoinedLine("", "- ", cci.name, ":")
	if err != nil {
		return err
	}

	// Description:
	err = trav.AppendString("Description:")
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

		err = ffs.ApplyForm(
			trav.GetConfig(
				ffs.WithModifiedIndent(1),
			),
			trav,
			NewDescriptionPrinter(cci.description),
		)
		if err != nil {
			return err
		}
	}

	// Empty line
	trav.EmptyLine()

	// Flags:
	if len(cci.flags) == 0 {
		err = trav.AddLine("Flags: [No flags provided]")
		if err != nil {
			return err
		}
	} else {
		err = trav.AddLine("Flags:")
		if err != nil {
			return err
		}

		for at, flag := range cci.flags {
			err := ffs.ApplyForm(
				trav.GetConfig(
					ffs.WithModifiedIndent(1),
				),
				trav,
				flag,
			)
			if err != nil {
				return ue.NewErrAt(at+1, "flag", err)
			}
		}
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
func NewCommandInfo(name string, description []string, fn CommandCallbackFunc, flagInfos []*FlagInfo) (*CommandInfo, error) {
	if name == "" {
		return nil, ue.NewErrInvalidParameter(
			"name",
			ue.NewErrEmpty(name),
		)
	}

	newCommand := &CommandInfo{
		name:        name,
		description: description,
	}

	if fn != nil {
		newCommand.callback = fn
	} else {
		newCommand.callback = NoCommandCallback
	}

	flagInfos = us.FilterNilValues(flagInfos)
	flagInfos = us.UniquefyEquals(flagInfos, false)

	newCommand.flags = flagInfos

	return newCommand, nil
}

// GetOpcode returns the name of a CommandInfo.
//
// Returns:
//   - string: The name of the CommandInfo.
func (inf *CommandInfo) GetOpcode() string {
	return inf.name
}

// HasFlag checks if a CommandInfo has a flag with the provided name.
//
// Parameters:
//   - name: The name of the flag to check.
//
// Returns:
//   - bool: True if the flag exists, false otherwise.
func (inf *CommandInfo) HasFlag(name string) bool {
	if name == "" {
		return false
	}

	for _, flag := range inf.flags {
		if flag.GetName() == name {
			return true
		}
	}

	return false
}

// GetFlag returns the FlagInfo with the provided name.
//
// Parameters:
//   - name: The name of the flag to get.
//
// Returns:
//   - *FlagInfo: The FlagInfo with the provided name. Nil if not found.
func (inf *CommandInfo) GetFlag(name string) *FlagInfo {
	if name == "" {
		return nil
	}

	for _, flag := range inf.flags {
		if flag.GetName() == name {
			return flag
		}
	}

	return nil
}

// ParseArguments parses the provided command line arguments
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
func (inf *CommandInfo) Parse(branches []*resultBranch, args []string) ([]*uc.Pair[*ParsedCommand, error], error) {
	// TODO: Handle the case where pos is not 0.

	uts.StableSort(branches, ResultBranchSortFunc, false)

	checkIfBranchHasRequiredFlags := func(branch *resultBranch) (*resultBranch, error) {
		err := branch.errIfInvalidRequiredFlags(inf.flags)
		return branch, err
	}

	solution, ok := us.EvaluateSimpleHelpers(branches, checkIfBranchHasRequiredFlags)
	if !ok {
		return nil, solution[0].GetData().Second
	} else if len(branches) == 0 {
		// No valid arguments is also a valid solution. Albeit, it is questionable.
		command := newParsedCommand(inf.name, nil, inf.callback)

		return []*uc.Pair[*ParsedCommand, error]{
			uc.NewPair(command, error(ue.NewErrIgnorable(
				errors.New("no valid arguments were found"),
			))),
		}, nil
	}

	branches = us.ExtractResults(solution)

	solution, ok = us.EvaluateSimpleHelpers(branches, errIfAnyError)
	if !ok {
		// No valid arguments is also a valid solution. Albeit, it is questionable.
		command := newParsedCommand(inf.name, nil, inf.callback)

		return []*uc.Pair[*ParsedCommand, error]{
			uc.NewPair(command, error(ue.NewErrIgnorable(
				errors.New("no valid arguments were found"),
			))),
		}, solution[0].GetData().Second
	}

	if len(branches) == 0 {
		// No valid arguments is also a valid solution. Albeit, it is questionable.
		command := newParsedCommand(inf.name, nil, inf.callback)

		return []*uc.Pair[*ParsedCommand, error]{
			uc.NewPair(command, error(ue.NewErrIgnorable(
				errors.New("no valid arguments were found"),
			))),
		}, nil
	}

	var possibleCommands []*uc.Pair[*ParsedCommand, error]

	for _, branch := range branches {
		var reason error

		result, err := branch.getResultMap()
		if err != nil {
			reason = err
		} else if len(branch.argsDone) < len(args) {
			reason = ue.NewErrIgnorable(
				errors.New("extra arguments provided"),
			)
		} else {
			reason = nil
		}

		command := newParsedCommand(inf.name, result, inf.callback)

		possibleCommands = append(possibleCommands, uc.NewPair(command, reason))
	}

	return possibleCommands, nil
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
	return inf.flags
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
