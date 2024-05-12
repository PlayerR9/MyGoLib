// Package CnsPanel provides a structure and functions for handling
// console command flags.
package CnsPanel

import (
	"slices"
	"strings"

	cdd "github.com/PlayerR9/MyGoLib/CustomData/Document"
	fs "github.com/PlayerR9/MyGoLib/Formatting/FString"
)

// FlagInfo represents a flag for a console command.
type FlagInfo struct {
	// Name of the flag.
	name string

	// Slice of strings representing the arguments accepted by
	// the flag. Order doesn't matter.
	args []string

	// Brief explanation of what the flag does.
	description *cdd.Document

	// Boolean indicating whether the flag is required.
	required bool

	// Function invoked when the flag is used.
	callback FlagCallbackFunc
}

// Equals compares two FlagInfo objects to determine if they are
// equal. Two FlagInfo objects are considered equal if they have
// the same name and arguments. The latter of which can be in any order.
//
// Parameters:
//   - other: The other FlagInfo to compare.
//
// Returns:
//   - bool: A boolean indicating whether the two FlagInfo objects
//     are equal.
func (fi *FlagInfo) Equals(other *FlagInfo) bool {
	if fi == nil || other == nil {
		return false
	}

	if fi.name != other.name {
		return false
	}

	if len(fi.args) != len(other.args) {
		return false
	}

	for _, arg := range fi.args {
		if !slices.Contains(other.args, arg) {
			return false
		}
	}

	return true
}

// FString generates a formatted string representation of a
// FlagInfo, including the flag name, arguments, description,
// and whether it is required.
//
// Parameters:
//   - indentLevel: The level of indentation to use. Sign is ignored.
//
// Returns:
//   - []string: A slice of strings representing the FlagInfo.
func (cfi *FlagInfo) FString(indentLevel int) []string {
	indentCfig := fs.NewIndentConfig(fs.DefaultIndentation, indentLevel, false)
	indent := indentCfig.String()

	results := make([]string, 0)

	var builder strings.Builder

	builder.WriteString(indent)
	builder.WriteString(cfi.name)

	// Add the arguments
	for _, arg := range cfi.args {
		builder.WriteRune(' ')
		builder.WriteRune('<')
		builder.WriteString(arg)
		builder.WriteRune('>')
	}

	results = append(results, builder.String())

	// Add the description
	builder.Reset()
	builder.WriteString(indent)

	if cfi.description == nil {
		builder.WriteString("Description: [No description provided]")

		results = append(results, builder.String())
	} else {
		builder.WriteString("Description:")
		results = append(results, builder.String())

		lines := cfi.description.FString(indentLevel + 1)
		results = append(results, lines...)
	}

	// Add the required information
	if len(cfi.args) != 0 {
		builder.Reset()
		builder.WriteString(indent)
		builder.WriteString("Required:")

		if cfi.required {
			builder.WriteString(" Yes")
		} else {
			builder.WriteString(" No")
		}

		results = append(results, builder.String())
	}

	return results
}

// NewFlagInfo creates a new FlagInfo with the given name and
// arguments.
//
// Parameters:
//   - name: The name of the flag.
//   - isRequired: A boolean indicating whether the flag is required.
//   - args: A slice of strings representing the arguments accepted by
//     the flag.
//
// Returns:
//   - *FlagInfo: A pointer to the new FlagInfo.
func NewFlagInfo(name string, isRequired bool, callback FlagCallbackFunc, args ...string) *FlagInfo {
	flag := &FlagInfo{
		name:        name,
		args:        make([]string, 0),
		description: nil,
		required:    isRequired,
	}

	flag.args = make([]string, 0, len(args))
	for _, arg := range args {
		arg = strings.TrimSpace(arg)
		if arg != "" {
			flag.args = append(flag.args, arg)
		}
	}
	flag.args = flag.args[:len(flag.args)]

	if callback == nil {
		flag.callback = func(...string) (any, error) {
			return nil, nil
		}
	} else {
		flag.callback = callback
	}

	return flag
}

// SetDescription sets the description of a FlagInfo.
//
// Parameters:
//   - description: The description of the FlagInfo.
func (cfi *FlagInfo) SetDescription(description *cdd.Document) {
	cfi.description = description
}
