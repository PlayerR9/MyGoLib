// Package CnsPanel provides a structure and functions for handling
// console command flags.
package pkg

import (
	"fmt"

	evalSlc "github.com/PlayerR9/MyGoLib/Evaluations/Slices"
	ffs "github.com/PlayerR9/MyGoLib/Formatting/FString"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
	us "github.com/PlayerR9/MyGoLib/Units/slice"
	uts "github.com/PlayerR9/MyGoLib/Utility/Sorting"
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
type FlagCallbackFunc func(argMap map[string][]any) (map[string][]any, error)

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
func NoFlagCallback(argMap map[string][]any) (map[string][]any, error) {
	return argMap, nil
}

// FlagInfo represents a flag for a console command.
type FlagInfo struct {
	// name is the name of the flag.
	name string

	// argList is a slice of Argument representing the arguments accepted by
	// the flag. Order doesn't matter.
	argList []*ArgInfo

	// description is the documentation of the flag.
	description []string

	// required is a boolean indicating whether the flag is required.
	required bool

	// callback is the function that parses the flag arguments.
	callback FlagCallbackFunc
}

// Evaluator implements the Evaluable interface.
func (inf *FlagInfo) Evaluator() evalSlc.LeafEvaluater[string, *FlagParseResult, *ArgInfo, []*resultArg] {
	return &flgEvaluator{
		argList:      inf.argList,
		startIndices: make([]int, 0),
		args:         make([]string, 0),
	}
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
func (cfi *FlagInfo) FString(trav *ffs.Traversor, opts ...ffs.Option) error {
	// Name:
	err := trav.AddJoinedLine(" ", "Flag:", cfi.name)
	if err != nil {
		return err
	}

	// Arguments:
	values := make([]string, 0, len(cfi.argList))
	for _, arg := range cfi.argList {
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

		err = ffs.ApplyForm(
			trav.GetConfig(
				ffs.WithModifiedIndent(1),
			),
			trav,
			NewDescriptionPrinter(cfi.description),
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
		return nil, uc.NewErrInvalidParameter(
			"name",
			uc.NewErrEmpty(name),
		)
	}

	newFlag := &FlagInfo{
		name:        name,
		description: doc,
		required:    isRequired,
	}

	argInfs = us.FilterNilValues(argInfs)
	newFlag.argList = us.UniquefyEquals(argInfs, false)

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
func (flag *FlagInfo) Parse(branches []*FlagParseResult, args []string) ([]*FlagParseResult, error) {
	if len(branches) == 0 {
		return nil, fmt.Errorf("no arguments provided")
	}

	solutions, ok := us.EvaluateWeightHelpers(
		branches,
		func(b *FlagParseResult) (*FlagParseResult, error) {
			parsed, err := flag.callback(b.GetResult())
			if err != nil {
				return nil, err
			}

			return &FlagParseResult{
				argMap:        parsed,
				argumentsDone: b.GetArgumentsDone(),
			}, nil
		},
		func(b *FlagParseResult) (float64, bool) {
			result := b.GetResult()

			return float64(len(result)), true
		},
		true,
	)
	if !ok {
		return nil, uc.NewErrPossibleError(fmt.Errorf("no valid arguments"), solutions[0].GetData().Second)
	}

	actualSolutions := us.ExtractResults(solutions)

	uts.StableSort(actualSolutions, FlagParseResultSortFunc, false)

	return actualSolutions, nil
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
	return inf.argList
}

// GetName returns the name of a FlagInfo.
//
// Returns:
//   - string: The name of the FlagInfo.
func (inf *FlagInfo) GetName() string {
	return inf.name
}
