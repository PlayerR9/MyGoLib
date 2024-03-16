// Package CMLine provides a structure and functions for handling
// console command flags.
package CMLine

import (
	"fmt"
	"strings"
)

// ConsoleFlagInfo represents a flag for a console command.
type ConsoleFlagInfo struct {
	// Name of the flag.
	name string

	// Slice of strings representing the arguments accepted by
	// the flag.
	args []string

	// Brief explanation of what the flag does.
	description []string

	// Boolean indicating whether the flag is required.
	required bool

	// Function invoked when the flag is used.
	callback func(...string) (any, error)
}

// FlagInfoOption is a function type that modifies ConsoleFlagInfo.
//
// Parameters:
//
//   - ConsoleFlagInfo: The ConsoleFlagInfo to modify.
type FlagInfoOption func(*ConsoleFlagInfo)

// WithArgs is a FlagInfoOption that sets the arguments for a
// ConsoleFlagInfo.
// It trims the space from each argument and ignores empty
// arguments.
//
// Parameters:
//
//   - args: The arguments to set.
//
// Returns:
//
//   - FlagInfoOption: A FlagInfoOption that sets the arguments for a
//     ConsoleFlagInfo.
func WithArgs(args ...string) FlagInfoOption {
	return func(flag *ConsoleFlagInfo) {
		flag.args = make([]string, 0, len(args))
		for _, arg := range args {
			arg = strings.TrimSpace(arg)
			if arg != "" {
				flag.args = append(flag.args, arg)
			}
		}
		flag.args = flag.args[:len(flag.args)]
	}
}

// WithFlagDescription is a FlagInfoOption that sets the description
// for a ConsoleFlagInfo.
// It splits each line of the description by newline characters.
//
// Parameters:
//
//   - description: The description to set.
//
// Returns:
//
//   - FlagInfoOption: A FlagInfoOption that sets the description for a
//     ConsoleFlagInfo.
func WithFlagDescription(description ...string) FlagInfoOption {
	return func(flag *ConsoleFlagInfo) {
		for _, line := range description {
			fields := strings.Split(line, "\n")
			flag.description = append(flag.description, fields...)
		}
	}
}

// WithRequired is a FlagInfoOption that sets whether a
// ConsoleFlagInfo is required.
//
// Parameters:
//
//   - required: Whether the flag is required.
//
// Returns:
//
//   - FlagInfoOption: A FlagInfoOption that sets whether a
//     ConsoleFlagInfo is required.
func WithRequired(required bool) FlagInfoOption {
	return func(flag *ConsoleFlagInfo) {
		flag.required = required
	}
}

// FString generates a formatted string representation of a
// ConsoleFlagInfo, including the flag name, arguments, description,
// and whether it is required.
//
// Parameters:
//
//   - indentLevel: The level of indentation to use. Sign is ignored.
//
// Returns:
//
//   - string: A string representing the ConsoleFlagInfo.
func (cfi *ConsoleFlagInfo) FString(indentLevel int) string {
	if indentLevel < 0 {
		indentLevel *= -1
	}

	var builder strings.Builder

	indentation := strings.Repeat("\t", indentLevel)

	// Add the flag name
	fmt.Fprintf(&builder, "%s%s", indentation, cfi.name)

	// Add the arguments
	for _, arg := range cfi.args {
		fmt.Fprintf(&builder, " <%s>", arg)
	}

	// Add the description
	fmt.Fprintf(&builder, "\n%sDescription:\n%s\t%s", indentation, indentation, cfi.description)

	// Add the required information
	if len(cfi.args) != 0 {
		fmt.Fprintf(&builder, "\n%sRequired: ", indentation)

		if cfi.required {
			fmt.Fprint(&builder, "Yes")
		} else {
			fmt.Fprint(&builder, "No")
		}
	}

	return builder.String()
}
