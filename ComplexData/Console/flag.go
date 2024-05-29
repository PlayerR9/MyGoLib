package Console

import (
	"strings"

	ffs "github.com/PlayerR9/MyGoLib/Formatting/FString"
)

// FlagParsingFunc is a function that parses a string into a value.
//
// Parameters:
//   - str: The string to parse.
//
// Returns:
//   - any: The parsed value.
//   - error: An error if the parsing fails.
type FlagParsingFunc func(str string) (any, error)

// Flag represents a flag for a command.
type Flag struct {
	// name is the name of the flag.
	name string

	// parsingFunc is the function used to parse the flag.
	parsingFunc FlagParsingFunc
}

// String implements the fmt.Stringer interface.
func (f *Flag) String() string {
	var builder strings.Builder

	builder.WriteRune('<')
	builder.WriteString(f.name)
	builder.WriteRune('>')

	return builder.String()
}

// FString implements the FStringer interface.
func (f *Flag) FString(trav *ffs.Traversor) error {
	if trav == nil {
		return nil
	}

	err := trav.AppendRune('<')
	if err != nil {
		return err
	}

	err = trav.AppendString(f.name)
	if err != nil {
		return err
	}

	err = trav.AppendRune('>')
	if err != nil {
		return err
	}

	return nil
}

// NewFlag creates a new flag with the specified name and parsing function.
//
// Parameters:
//   - name: The name of the flag.
//   - parsingFunc: The function used to parse the flag.
//
// Returns:
//   - *Flag: The new flag.
//
// Behaviors:
//   - If the parsing function is nil, the flag uses the default parsing function, which returns the string as is.
func NewFlag(name string, parsingFunc FlagParsingFunc) *Flag {
	flag := &Flag{
		name: name,
	}

	if parsingFunc != nil {
		flag.parsingFunc = parsingFunc
	} else {
		flag.parsingFunc = func(str string) (any, error) {
			return str, nil
		}
	}

	return flag
}
