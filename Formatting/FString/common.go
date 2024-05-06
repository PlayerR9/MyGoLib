package FString

import "strings"

const (
	// DefaultIndentation is the default indentation string.
	DefaultIndentation string = "   "

	// DefaultSeparator is the default separator string.
	DefaultSeparator string = ", "
)

// Options is a type that represents the options that can be passed to the builder.
type Options []BuildOption

var (
	// ArrayDefault is the default options for an array.
	// [1, 2, 3]
	ArrayDefault Options = Options{
		WithIndentation(NewIndentConfig(DefaultIndentation, 0, false, true)),
		WithDelimiterLeft(NewDelimiterConfig("[", false)),
		WithDelimiterRight(NewDelimiterConfig("]", false)),
		WithSeparator(NewSeparator(DefaultSeparator, false)),
	}
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
	FString(int) []string
}

// FString is a function that returns a string representation of an object that
// implements the FStringer interface.
//
// It joins the strings returned by the FString method of the object using a newline
// character with no indentation at the beginning.
//
// Parameters:
//   - obj: The object that implements the FStringer interface.
//
// Returns:
//   - string: A string representation of the object.
func FString(obj FStringer) string {
	if obj == nil {
		return ""
	}

	return strings.Join(obj.FString(0), "\n")
}
