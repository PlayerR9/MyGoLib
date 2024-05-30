package pkg

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	uc "github.com/PlayerR9/MyGoLib/Units/Common"
	ue "github.com/PlayerR9/MyGoLib/Units/Errors"
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

// ArgInfo represents an argument of a flag.
type ArgInfo struct {
	// name of the argument.
	name string

	// qty is the number of the values of the argument.
	qty [2]int

	// parserFunc is the function that parses the argument.
	parserFunc ArgumentParserFunc
}

// Equals checks if the argument is equal to another argument.
//
// Two arguments are equal iff their names are equal.
//
// Parameters:
//   - other: The other argument to compare.
//
// Returns:
//   - bool: true if the arguments are equal, false otherwise.
func (a *ArgInfo) Equals(other uc.Equaler) bool {
	if other == nil {
		return false
	}

	otherA, ok := other.(*ArgInfo)
	if !ok {
		return false
	}

	return a.name == otherA.name
}

// String returns the string: <name>
//
// Returns:
//   - string: The string representation of the argument.
func (a *ArgInfo) String() string {
	var builder strings.Builder

	builder.WriteString(a.name)

	if a.qty[0] != 0 {
		builder.WriteRune(':')
		builder.WriteString(strconv.Itoa(a.qty[0]))
	}

	if a.qty[1] != -1 && a.qty[0] != a.qty[1] {
		if a.qty[0] != 0 {
			builder.WriteRune('-')
		} else {
			builder.WriteString(":")
		}

		builder.WriteString(strconv.Itoa(a.qty[1]))
	}

	var str strings.Builder

	if a.qty[1] == -1 {
		if a.qty[0] == 0 {
			str.WriteRune('(')
			builder.WriteRune(')')
		} else {
			str.WriteRune('{')
			builder.WriteRune('}')
		}
	} else {
		if a.qty[0] == 0 {
			str.WriteRune('[')
			builder.WriteRune(']')
		} else {
			str.WriteRune('<')
			builder.WriteRune('>')
		}
	}

	str.WriteString(builder.String())

	return str.String()
}

func checkArgumentFormat(format string) ([2]int, error) {
	var qty [2]int

	fields := strings.Split(format, "-")

	if len(fields) > 2 {
		return qty, fmt.Errorf("expected 2 fields, got %d", len(fields))
	}

	var min int

	if fields[0] == "" {
		min = 0
	} else {
		var err error

		min, err = strconv.Atoi(fields[0])
		if err != nil {
			return qty, err
		}
	}

	qty[0] = min

	var max int

	if fields[1] == "" {
		max = -1
	} else {
		var err error

		max, err = strconv.Atoi(fields[1])
		if err != nil {
			return qty, err
		}

		if max < min {
			return qty, errors.New("max is less than min")
		}
	}

	qty[1] = max

	return qty, nil
}

func NewArgument(format string, fn ArgumentParserFunc) (*ArgInfo, error) {
	if format == "" {
		return nil, ue.NewErrInvalidParameter(
			"format",
			ue.NewErrEmpty(format),
		)
	}

	fields := strings.Fields(format)

	newArg := &ArgInfo{
		name: fields[0],
	}

	if fn != nil {
		newArg.parserFunc = fn
	} else {
		newArg.parserFunc = NoArgumentParser
	}

	if len(fields) > 2 {
		return nil, NewErrFormatError(errors.New("expected 1 or 2 fields"))
	}

	if len(fields) == 1 {
		newArg.qty = [2]int{1, 1}
	} else {
		qty, err := checkArgumentFormat(fields[1])
		if err != nil {
			return nil, NewErrFormatError(err)
		}

		newArg.qty = qty
	}

	return newArg, nil
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
func (a *ArgInfo) Parse(s string) (any, error) {
	parsed, err := a.parserFunc(s)
	if err != nil {
		return nil, err
	}

	return parsed, nil
}

// GetName returns the name of the argument.
//
// Returns:
//   - string: The name of the argument.
func (a *ArgInfo) GetName() string {
	return a.name
}
