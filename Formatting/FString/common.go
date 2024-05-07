package FString

import (
	"strings"
)

const (
	// DefaultIndentation is the default indentation string.
	DefaultIndentation string = "   "

	// DefaultSeparator is the default separator string.
	DefaultSeparator string = ", "
)

var (
	// ArrayDefault is the default options for an array.
	// [1, 2, 3]
	ArrayDefault *Formatter = func() *Formatter {
		var builder Builder

		builder.SetDelimiterLeft(NewDelimiterConfig("[", false))
		builder.SetDelimiterRight(NewDelimiterConfig("]", false))
		builder.SetSeparator(NewSeparator(DefaultSeparator, false))

		return builder.Build()
	}()
)

// FStringer is an interface that defines the behavior of a type that can be
// converted to a string representation.
type FStringer interface {
	// FString returns a string representation of the object.
	//
	// Parameters:
	//   - int: The current indentation level.
	//
	// Returns:
	//   - []string: A slice of strings that represent the object.
	FString() []string
}

func FString(form *Formatter, values []string) string {
	if form == nil {
		var builder Builder

		form = builder.Build()
	}

	return strings.Join(form.Apply(values), "\n")
}
