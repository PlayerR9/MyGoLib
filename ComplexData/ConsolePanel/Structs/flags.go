// Package CnsPanel provides a structure and functions for handling
// console command flags.
package Structs

import (
	"errors"
	"strings"

	cdd "github.com/PlayerR9/MyGoLib/CustomData/Document"
	fs "github.com/PlayerR9/MyGoLib/Formatting/FString"
	slext "github.com/PlayerR9/MyGoLib/Utility/SliceExt"

	com "github.com/PlayerR9/MyGoLib/ComplexData/ConsolePanel/Common"
	res "github.com/PlayerR9/MyGoLib/ComplexData/ConsolePanel/Results"
)

// FlagInfo represents a flag for a console command.
type FlagInfo struct {
	// args is a slice of Argument representing the arguments accepted by
	// the flag. Order doesn't matter.
	args []*Argument

	// description is the documentation of the flag.
	description *cdd.Document

	// required is a boolean indicating whether the flag is required.
	required bool

	// callback is the function that parses the flag arguments.
	callback com.FlagCallbackFunc
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
func (cfi *FlagInfo) FString(indentLevel int) []string {
	lines := make([]string, 0)
	var builder strings.Builder

	// Arguments:
	values := make([]string, 0, len(cfi.args))
	for _, arg := range cfi.args {
		values = append(values, arg.String())
	}

	builder.WriteString("Arguments: ")
	builder.WriteString(strings.Join(values, " "))

	lines = append(lines, builder.String())

	// Empty line
	lines = append(lines, "")

	// Description:
	if cfi.description == nil {
		lines = append(lines, "Description: [No description provided]")
	} else {
		lines = append(lines, "Description:")

		descriptionLines := cfi.description.FString(indentLevel + 1)
		lines = append(lines, descriptionLines...)
	}

	// Empty line
	lines = append(lines, "")

	// Required:
	if cfi.required {
		lines = append(lines, "Required: Yes")
	} else {
		lines = append(lines, "Required: No")
	}

	// Add the indentation to each line
	indentCfig := fs.NewIndentConfig(fs.DefaultIndentation, indentLevel, false)
	indent := indentCfig.String()

	for i := 0; i < len(lines); i++ {
		lines[i] = indent + lines[i]
	}

	return lines
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
func NewFlagInfo(isRequired bool, callback com.FlagCallbackFunc, args ...*Argument) *FlagInfo {
	flag := &FlagInfo{
		description: nil,
		required:    isRequired,
	}

	flag.args = slext.FilterNilValues(args)

	if callback == nil {
		flag.callback = com.NoFlagCallback
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
func (cfi *FlagInfo) SetDescription(description *cdd.Document) *FlagInfo {
	cfi.description = description

	return cfi
}

// Parse parses the provided arguments into a map of parsed arguments.
//
// Parameters:
//   - args: The arguments to parse.
//
// Returns:
//   - ArgumentsMap: A map of the parsed arguments.
//   - int: The index of the last unsuccessful parse argument.
//   - bool: A boolean indicating whether the error is ignorable.
//   - error: An error, if any.
func (flag *FlagInfo) Parse(args []string) (*res.FlagParseResult, error) {
	if len(args) == 0 {
		return res.NewFlagParseResult(nil, -1), errors.New("no arguments provided")
	}

	if len(args) <= len(flag.args) {
		return res.NewFlagParseResult(nil, 1), errors.New("not enough arguments provided")
	} else if len(args)+1 > len(flag.args) {
		return res.NewFlagParseResult(nil, 1), errors.New("too many arguments provided")
	}

	parsedArgs := make(com.ArgumentsMap) // Map to store the parsed arguments

	i := 1
	for i < len(flag.args) {
		arg := flag.args[i]

		parsedArg, err := arg.Parse(args[i])
		if err != nil {
			return res.NewFlagParseResult(nil, i), err
		}

		parsedArgs[arg.name] = parsedArg
		i++
	}

	parsed, err := flag.callback(parsedArgs)
	if err != nil {
		return res.NewFlagParseResult(nil, i), err
	}

	return res.NewFlagParseResult(parsed, i), nil
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
