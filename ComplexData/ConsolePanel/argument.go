package ConsolePanel

import (
	"strings"
)

// ArgumentParserFunc is a function type that represents a function
// that parses a string argument.
//
// Parameters:
//   - string: The string to parse.
//
// Returns:
//   - any: The parsed value.
type ArgumentParserFunc func(string) (any, error)

// NoArgumentParser is a default argument parser function that returns
// the string as is.
//
// Parameters:
//   - string: The string to parse.
//
// Returns:
//   - any: The string as is.
//   - error: nil
func NoArgumentParser(s string) (any, error) {
	return s, nil
}

// Argument represents an argument of a flag.
type Argument struct {
	// name of the argument.
	name string

	// parserFunc is the function that parses the argument.
	parserFunc ArgumentParserFunc
}

// String returns the string: <name>
//
// Returns:
//   - string: The string representation of the argument.
func (a *Argument) String() string {
	var builder strings.Builder

	builder.WriteRune('<')
	builder.WriteString(a.name)
	builder.WriteRune('>')

	return builder.String()
}

// NewArgument creates a new argument.
//
// Parameters:
//   - name: The name of the argument.
//   - argumentParserFunc: The function that parses the argument.
//
// Returns:
//   - *Argument: A pointer to the newly created argument.
//
// Behaviors:
//   - If argumentParserFunc is nil, the default NoArgumentParser is used.
func NewArgument(name string, argumentParserFunc ArgumentParserFunc) *Argument {
	arg := &Argument{
		name: name,
	}

	if argumentParserFunc != nil {
		arg.parserFunc = argumentParserFunc
	} else {
		arg.parserFunc = NoArgumentParser
	}

	return arg
}

// Parse parses a string into the argument.
//
// Parameters:
//   - s: The string to parse.
//
// Returns:
//   - any: The parsed value.
//   - error: An error if the parsing fails.
//
// Errors:
//   - The error returned is the error from the parser function.
func (a *Argument) Parse(s string) (any, error) {
	parsed, err := a.parserFunc(s)
	if err != nil {
		return nil, err
	}

	return parsed, nil
}
