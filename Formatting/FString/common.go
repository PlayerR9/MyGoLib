package FString

import (
	"strings"

	uc "github.com/PlayerR9/lib_units/common"
)

// Stringify converts a formatted string to a string.
//
// Parameters:
//   - doc: The formatted string.
//   - spacing: The spacing to use for the string.
//
// Returns:
//   - [][]string: The stringified formatted string.
//
// Behaviors:
//   - If the spacing is less than or equal to 0, the spacing is set to 1.
func Stringfy(doc [][][][]string, spacing int) []string {
	if spacing <= 0 {
		spacing = 1
	}

	space := strings.Repeat(" ", spacing)

	var pages1 [][][]string

	for _, page := range doc {
		var sections [][]string

		for _, section := range page {
			var lines []string

			for _, line := range section {
				joinedStr := strings.Join(line, space)
				lines = append(lines, joinedStr)
			}

			sections = append(sections, lines)
		}

		pages1 = append(pages1, sections)
	}

	var pages2 [][]string

	for _, page := range pages1 {
		var sections []string

		for _, section := range page {
			joinedStr := strings.Join(section, "\n")
			sections = append(sections, joinedStr)
		}

		pages2 = append(pages2, sections)
	}

	var pages []string

	for _, page := range pages2 {
		joinedStr := strings.Join(page, "\n\n")
		pages = append(pages, joinedStr)
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
	ArrayLikeFormat *FormatConfig
)

func init() {
	ArrayLikeFormat = NewFormatter(
		NewDelimiterConfig("[", false, true),
		NewDelimiterConfig("]", false, false),
		NewSeparator(DefaultSeparator, false),
	)
}

// FStringArray generates a formatted string representation of an array-like object.
//
// Parameters:
//   - format: The format to use for printing.
//   - values: The values to print.
//
// Returns:
//   - string: The formatted string.
//   - error: An error if the printing fails.
//
// Behaviors:
//   - If the format is nil, the function uses ArrayLikeFormat.
func FStringArray(format *FormatConfig, values []string) (string, error) {
	if format == nil {
		format = ArrayLikeFormat.Copy().(*FormatConfig)
	}

	doc, err := Sprint(format, values...)
	if err != nil {
		return "", err
	}

	pages := Stringfy(doc, 1)
	joinedStr := strings.Join(pages, "\f")

	return joinedStr, nil
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

// ApplyTravFunc applies a function to the printer. Useful for when you want to apply a function
// that does not implement the FStringer interface.
//
// Parameters:
//   - trav: The traversor to use.
//   - elem: The element to apply the function to.
//   - f: The function to apply.
//
// Returns:
//   - error: An error if the function fails.
//
// Errors:
//   - *ErrFinalization: If the finalization of the page fails.
//   - any error returned by the function.
func ApplyTravFunc[T any](trav *Traversor, elem T, f FStringFunc[T]) error {
	err := f(trav, elem)
	if err != nil {
		return err
	}

	return nil
}

// ApplyTravFuncMany applies a function to the printer. Useful for when you want to apply a function
// that does not implement the FStringer interface.
//
// Parameters:
//   - trav: The traversor to use.
//   - f: The function to apply.
//   - elems: The elements to apply the function to.
//
// Returns:
//   - error: An error if the function fails.
//
// Errors:
//   - *ErrFinalization: If the finalization of the page fails.
//   - *Errors.ErrAt: If an error occurs on a specific element.
//   - any error returned by the function.
func ApplyTravFuncMany[T any](trav *Traversor, f FStringFunc[T], elems []T) error {
	if len(elems) == 0 {
		return nil
	}

	for i, elem := range elems {
		err := f(trav, elem)
		if err != nil {
			return uc.NewErrAt(i+1, "element", err)
		}
	}

	return nil
}
