package FString

import (
	"strings"

	uc "github.com/PlayerR9/MyGoLib/Units/Common"
)

// Stringify converts a formatted string to a string.
//
// Parameters:
//   - doc: The formatted string.
//
// Returns:
//   - [][]string: The stringified formatted string.
func Stringfy(doc [][][][]string) []string {
	var pages []string

	for _, page := range doc {
		var sections []string

		for _, section := range page {
			var lines []string

			for _, line := range section {
				lines = append(lines, strings.Join(line, " "))
			}

			sections = append(sections, strings.Join(lines, "\n"))
		}

		pages = append(pages, strings.Join(sections, "\n"))
	}

	return pages
}

/////////////////////////////////////////////////

// FStringer is an interface that defines the behavior of a type that can be
// converted to a string representation.
type FStringer interface {
	// FString returns a string representation of the object.
	//
	// Parameters:
	//   - trav: The traversor to use for printing.
	//   - opts: The options to use for printing.
	//
	// Returns:
	//   - error: An error if there was a problem generating the string.
	FString(trav *Traversor, opts ...Option) error
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
	ArrayLikeFormat FormatConfig = NewFormatter(
		NewDelimiterConfig("[", false, true),
		NewDelimiterConfig("]", false, false),
		NewSeparator(DefaultSeparator, false),
	)
)

// FStringArray generates a formatted string representation of an array-like object.
//
// Parameters:
//   - format: The format to use for printing.
//   - values: The values to print.
//
// Returns:
//   - string: The formatted string.
//   - error: An error if the printing fails.
func FStringArray(format FormatConfig, values []string) (string, error) {
	doc, err := Sprint(format, values...)
	if err != nil {
		return "", err
	}

	return strings.Join(Stringfy(doc), "\f"), nil
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
