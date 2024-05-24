package FString

import (
	uc "github.com/PlayerR9/MyGoLib/Units/Common"
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
	//   - error: An error if there was a problem generating the string.
	FString(trav *Traversor) error
}

/*
func FString(form *Formatter, values []string) string {
	if form == nil {
		var builder Builder

		form = builder.Build()
	}

	return strings.Join(form.Apply(values), "\n")
}
*/

// SimplePrinter is a simple printer that prints a value with a name.
type SimplePrinter[T comparable] struct {
	// name is the name of the value.
	name string

	// value is the value to print.
	value T

	// fn is the function to use to convert the value to a string.
	fn func(T) (string, error)
}

// FString generates a formatted string representation of a SimplePrinter.
//
// Format:
//
//	<name>: <value>
//
// Parameters:
//   - trav: The traversor to use for printing.
//
// Returns:
//   - error: An error if the printing fails.
func (sp *SimplePrinter[T]) FString(trav *Traversor) error {
	str, err := sp.fn(sp.value)
	if err != nil {
		return err
	}

	err = trav.AddJoinedLine("", sp.name, ": ", str)
	if err != nil {
		return err
	}

	return nil
}

// NewSimplePrinter creates a new SimplePrinter with the provided name and value.
//
// Parameters:
//   - name: The name of the value.
//   - value: The value to print.
//   - fn: The function to use to convert the value to a string.
//
// Returns:
//   - *SimplePrinter: The new SimplePrinter.
//
// Behaviors:
//   - If the function is nil, the function uses uc.StringOf to convert the value to a string.
func NewSimplePrinter[T comparable](name string, value T, fn func(T) (string, error)) *SimplePrinter[T] {
	if fn == nil {
		fn = func(v T) (string, error) {
			return uc.StringOf(v), nil
		}
	}

	return &SimplePrinter[T]{
		name:  name,
		value: value,
		fn:    fn,
	}
}
