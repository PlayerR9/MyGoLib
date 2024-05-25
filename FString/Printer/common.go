package Printer

import (
	"fmt"

	uc "github.com/PlayerR9/MyGoLib/Units/Common"
)

// FStringer is an interface that defines the behavior of a type that can be
// converted to a string representation.
type FStringer interface {
	// FString returns a string representation of the object.
	//
	// Parameters:
	//   - trav: The traversor to use for printing.
	//
	// Returns:
	//   - error: An error if there was a problem generating the string.
	FString(trav *Traversor) error
}

// FStringFunc is a function that generates a formatted string representation of an object.
//
// Parameters:
//   - trav: The traversor to use for printing.
//   - elem: The element to print.
//
// Returns:
//   - error: An error if there was a problem generating the string.
type FStringFunc[T any] func(trav *Traversor, elem T) error

var (
	// ArrayLikeFormat is the default options for an array-like object.
	// [1, 2, 3]
	ArrayLikeFormat *Formatter = NewFormatter(
		nil,
		NewDelimiterConfig("[", false),
		NewDelimiterConfig("]", false),
		NewSeparator(DefaultSeparator, false),
	)
)

// AdvancedFields extracts the fields in a similar way to strings.Fields, but with more
// advanced formatting.
//
// Parameters:
//   - str: The string to extract the fields from.
//
// Returns:
//   - []string: The fields in the string.
//   - error: An error if the extraction fails.
//
// Behaviors:
//   - The function uses a Printer to extract the fields.
//   - The function uses the default options for the Printer.
//   - The function returns the fields in the order they appear in the string.
func AdvancedFields(str string) ([]string, error) {
	printer := NewPrinter(nil)

	err := ApplyFormatFunc(printer, str, func(trav *Traversor, str string) error {
		err := trav.AppendString(str)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error applying format: %w", err)
	}

	rawPages, err := printer.GetRaw()
	if err != nil {
		return nil, fmt.Errorf("error getting raw pages: %w", err)
	}

	var fields [][]string

	for _, rp := range rawPages {
		for _, line := range rp {
			fields = append(fields, line.GetRawContent()...)
		}
	}

	res := make([]string, 0)

	for _, field := range fields {
		res = append(res, field...)
	}

	return res, nil
}

//////////////////////////////////////////////////////////////

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
