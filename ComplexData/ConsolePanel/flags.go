// Package CnsPanel provides a structure and functions for handling
// console command flags.
package ConsolePanel

import (
	"errors"

	fsd "github.com/PlayerR9/MyGoLib/FString/Document"
	fsp "github.com/PlayerR9/MyGoLib/FString/Printer"
	slext "github.com/PlayerR9/MyGoLib/Units/Slices"
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
	// args is a slice of Argument representing the arguments accepted by
	// the flag. Order doesn't matter.
	args []*Argument

	// description is the documentation of the flag.
	description *fsd.Document

	// required is a boolean indicating whether the flag is required.
	required bool

	// callback is the function that parses the flag arguments.
	callback FlagCallbackFunc
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
func (cfi *FlagInfo) FString(trav *fsp.Traversor) error {
	// Arguments:
	values := make([]string, 0, len(cfi.args))
	for _, arg := range cfi.args {
		values = append(values, arg.String())
	}

	values = append([]string{"Arguments:"}, values...)

	err := trav.AddJoinedLine(" ", values...)
	if err != nil {
		return err
	}

	// Empty line
	trav.EmptyLine()

	/*
		// Description:
		doc := fsd.NewDocumentPrinter("Description", cfi.description, "[No description provided]")
		err = doc.FString(trav)
		if err != nil {
			return ers.NewErrWhile("FString printing description", err)
		}
	*/

	// Empty line
	trav.EmptyLine()

	// Required:
	sp := fsp.NewSimplePrinter("Required", cfi.required, BoolFString)
	err = sp.FString(trav)
	if err != nil {
		return err
	}

	return nil
}

// NewFlagInfo creates a new FlagInfo with the given name and
// arguments.
//
// Parameters:
//   - isRequired: A boolean indicating whether the flag is required.
//   - callback: The function that parses the flag arguments.
//   - args: A slice of strings representing the arguments accepted by
//     the flag.
//
// Returns:
//   - *FlagInfo: A pointer to the new FlagInfo.
//
// Behaviors:
//   - Any nil arguments are filtered out.
//   - If 'callback' is nil, a default callback is used that returns nil without error.
func NewFlagInfo(isRequired bool, callback FlagCallbackFunc, args ...*Argument) *FlagInfo {
	flag := &FlagInfo{
		description: nil,
		required:    isRequired,
	}

	flag.args = slext.FilterNilValues(args)

	if callback == nil {
		flag.callback = NoFlagCallback
	} else {
		flag.callback = callback
	}

	return flag
}

// IsRequired returns whether a FlagInfo is required.
//
// Returns:
//   - bool: A boolean indicating whether the FlagInfo is required.
func (inf *FlagInfo) IsRequired() bool {
	return inf.required
}

// SetDescription sets the description of a FlagInfo.
//
// Parameters:
//   - description: The description of the FlagInfo.
//
// Returns:
//   - *FlagInfo: A pointer to the FlagInfo. This allows for chaining.
func (cfi *FlagInfo) SetDescription(description *fsd.Document) *FlagInfo {
	cfi.description = description

	return cfi
}

// Parse parses the provided arguments into a map of parsed arguments.
//
// Parameters:
//   - args: The arguments to parse.
//
// Returns:
//   - map[string]any: A map of the parsed arguments.
//   - int: The index of the last unsuccessful parse argument.
//   - bool: A boolean indicating whether the error is ignorable.
//   - error: An error, if any.
func (flag *FlagInfo) Parse(args []string) (*FlagParseResult, error) {
	if len(args) == 0 {
		return NewFlagParseResult(nil, -1), errors.New("no arguments provided")
	}

	if len(args) <= len(flag.args) {
		return NewFlagParseResult(nil, 1), errors.New("not enough arguments provided")
	} else if len(args)+1 > len(flag.args) {
		return NewFlagParseResult(nil, 1), errors.New("too many arguments provided")
	}

	parsedArgs := make(map[string]any) // Map to store the parsed arguments

	i := 1
	for i < len(flag.args) {
		arg := flag.args[i]

		parsedArg, err := arg.Parse(args[i])
		if err != nil {
			return NewFlagParseResult(nil, i), err
		}

		parsedArgs[arg.name] = parsedArg
		i++
	}

	parsed, err := flag.callback(parsedArgs)
	if err != nil {
		return NewFlagParseResult(nil, i), err
	}

	return NewFlagParseResult(parsed, i), nil
}

// GetArguments returns the arguments of a FlagInfo.
//
// Returns:
//   - []*Argument: A slice of pointers to the arguments.
//
// Behaviors:
//   - No nil values are returned.
//   - Modifying the returned slice will affect the FlagInfo.
func (inf *FlagInfo) GetArguments() []*Argument {
	return inf.args
}
