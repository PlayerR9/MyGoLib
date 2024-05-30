// Package CnsPanel provides a structure and functions for handling
// console command flags.
package pkg

import (
	"fmt"

	fss "github.com/PlayerR9/MyGoLib/Formatting/FString"
	uc "github.com/PlayerR9/MyGoLib/Units/Common"
	ue "github.com/PlayerR9/MyGoLib/Units/Errors"
	us "github.com/PlayerR9/MyGoLib/Units/Slices"
)

// FlagCallbackFunc is a function type that represents a callback
// function for a console command flag.
//
// Parameters:
//   - argMap: A map of string keys to any values representing the
//     arguments passed to the flag.
//
// Returns:
//   - map[string]any: A map of string keys to any values representing
//     the parsed arguments.
//   - error: An error if the flag fails.
type FlagCallbackFunc func(argMap map[string]any) (map[string]any, error)

// NoFlagCallback is a default callback function for a console command flag
// when no callback is provided.
//
// Parameters:
//   - args: A slice of strings representing the arguments passed to
//     the flag.
//
// Returns:
//   - map[string]any: A map of string keys to any values representing
//     the parsed arguments.
//   - error: nil
func NoFlagCallback(argMap map[string]any) (map[string]any, error) {
	return argMap, nil
}

// FlagInfo represents a flag for a console command.
type FlagInfo struct {
	// name is the name of the flag.
	name string

	// args is a slice of Argument representing the arguments accepted by
	// the flag. Order doesn't matter.
	args []*ArgInfo

	// description is the documentation of the flag.
	description []string

	// required is a boolean indicating whether the flag is required.
	required bool

	// callback is the function that parses the flag arguments.
	callback FlagCallbackFunc
}

// Equals checks if the flag is equal to another flag.
//
// Two flags are equal iff their names are equal.
//
// Parameters:
//   - other: The other flag to compare.
//
// Returns:
//   - bool: true if the flags are equal, false otherwise.
func (inf *FlagInfo) Equals(other uc.Equaler) bool {
	if other == nil {
		return false
	}

	otherFlag, ok := other.(*FlagInfo)
	if !ok {
		return false
	}

	return inf.name == otherFlag.name
}

// FString generates a formatted string representation of a FlagInfo.
//
// Parameters:
//   - indentLevel: The level of indentation to use for the FlagInfo.
//
// Returns:
//   - []string: A slice of strings representing the FlagInfo.
//
// Format:
//
//	Arguments: <arg1> <arg2> ...
//
//	Description:
//		// <description>
//
//	Required: <Yes/No>
func (cfi *FlagInfo) FString(trav *fss.Traversor) error {
	// Name:
	err := trav.AddJoinedLine(" ", "Flag:", cfi.name)
	if err != nil {
		return err
	}

	// Arguments:
	values := make([]string, 0, len(cfi.args))
	for _, arg := range cfi.args {
		values = append(values, arg.String())
	}

	values = append([]string{"Arguments:"}, values...)

	err = trav.AddJoinedLine(" ", values...)
	if err != nil {
		return err
	}

	// Empty line
	trav.EmptyLine()

	err = trav.AppendString("Description:")
	if err != nil {
		return err
	}

	if len(cfi.description) == 0 {
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
			&DescriptionPrinter{cfi.description},
		)
		if err != nil {
			return err
		}
	}

	// Empty line
	trav.EmptyLine()

	// Required:
	err = trav.AppendString("Required: ")
	if err != nil {
		return err
	}

	if cfi.required {
		err = trav.AppendString("Yes")
	} else {
		err = trav.AppendString("No")
	}

	if err != nil {
		return err
	}

	return nil
}

func NewFlagInfo(name string, doc []string, isRequired bool, fn FlagCallbackFunc, argInfs []*ArgInfo) (*FlagInfo, error) {
	if name == "" {
		return nil, ue.NewErrInvalidParameter(
			"name",
			ue.NewErrEmpty(name),
		)
	}

	newFlag := &FlagInfo{
		name:        name,
		description: doc,
		required:    isRequired,
	}

	argInfs = us.FilterNilValues(argInfs)
	newFlag.args = us.UniquefyEquals(argInfs, false)

	if fn != nil {
		newFlag.callback = fn
	} else {
		newFlag.callback = NoFlagCallback
	}

	return newFlag, nil
}

// IsRequired returns whether a FlagInfo is required.
//
// Returns:
//   - bool: A boolean indicating whether the FlagInfo is required.
func (inf *FlagInfo) IsRequired() bool {
	return inf.required
}

// Parse parses the provided arguments into a map of parsed arguments.
//
// Parameters:
//   - args: The arguments to parse.
//
// Returns:
//   - map[string]any: A map of the parsed arguments.
//   - error: An error if the arguments are invalid.
func (flag *FlagInfo) Parse(args []string) (map[string]any, error) {
	if len(args) < len(flag.args) {
		return nil, fmt.Errorf("missing argument %q", flag.args[len(args)].GetName())
	}

	parsedArgs := make(map[string]any) // Map to store the parsed arguments

	for i, arg := range flag.args {
		res, err := arg.Parse(args[i])
		if err != nil {
			return nil, ue.NewErrAt(i+1, "argument", err)
		}

		parsedArgs[arg.GetName()] = res
	}

	parsed, err := flag.callback(parsedArgs)
	if err != nil {
		return nil, fmt.Errorf("invalid arguments: %w", err)
	}

	var reason error

	if len(args) > len(flag.args) {
		return nil, ue.NewErrIgnorable(
			fmt.Errorf("argument %q is an extra argument", args[len(flag.args)]),
		)
	} else {
		reason = nil
	}

	return parsed, reason
}

// GetArguments returns the arguments of a FlagInfo.
//
// Returns:
//   - []*Argument: A slice of pointers to the arguments.
//
// Behaviors:
//   - No nil values are returned.
//   - Modifying the returned slice will affect the FlagInfo.
func (inf *FlagInfo) GetArguments() []*ArgInfo {
	return inf.args
}

// GetName returns the name of a FlagInfo.
//
// Returns:
//   - string: The name of the FlagInfo.
func (inf *FlagInfo) GetName() string {
	return inf.name
}
