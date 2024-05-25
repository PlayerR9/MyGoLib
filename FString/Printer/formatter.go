package Printer

import ue "github.com/PlayerR9/MyGoLib/Units/Errors"

var (
	// DefaultFormatter is the default formatter.
	//
	// ==IndentConfig==
	//   - DefaultIndentationConfig
	//
	// ==SeparatorConfig==
	//   - DefaultSeparatorConfig
	//
	// ==DelimiterConfig (Left and Right)==
	//   - Nil (no delimiters are used by default)
	DefaultFormatter *Formatter = NewFormatter(
		DefaultIndentationConfig,
		nil,
		nil,
		DefaultSeparatorConfig,
	)
)

// Formatter is a type that represents a builder for creating formatted strings.
type Formatter struct {
	// The indentation configuration of the builder.
	indent *IndentConfig

	// The left delimiter configuration of the builder.
	delimiterLeft *DelimiterConfig

	// The right delimiter configuration of the builder.
	delimiterRight *DelimiterConfig

	// The separator configuration of the builder.
	separator *SeparatorConfig
}

// Copy is a method that creates a copy of the formatter.
//
// Returns:
//   - *Formatter: A pointer to the new formatter.
func (f *Formatter) Copy() *Formatter {
	return &Formatter{
		indent:         f.indent.Copy().(*IndentConfig),
		delimiterLeft:  f.delimiterLeft.Copy().(*DelimiterConfig),
		delimiterRight: f.delimiterRight.Copy().(*DelimiterConfig),
		separator:      f.separator.Copy().(*SeparatorConfig),
	}
}

// NewFormatter creates a new formatter with the given configuration.
//
// Nil parameters will be set to their default values.
//
// Parameters:
//   - indent: The indentation configuration of the builder.
//   - delimiterLeft: The left delimiter configuration of the builder.
//   - delimiterRight: The right delimiter configuration of the builder.
//   - separator: The separator configuration of the builder.
//
// Returns:
//   - *Formatter: A pointer to the newly created formatter.
//
// ==IndentConfig==
//   - Nil (no indentation is used by default)
//
// ==SeparatorConfig==
//   - Nil (no separator is used by default)
//
// ==DelimiterConfig (Left and Right)==
//   - Nil (no delimiters are used by default)
func NewFormatter(indent *IndentConfig, delimiterLeft *DelimiterConfig, delimiterRight *DelimiterConfig, separator *SeparatorConfig) *Formatter {
	f := &Formatter{
		indent:         indent,
		delimiterLeft:  delimiterLeft,
		delimiterRight: delimiterRight,
		separator:      separator,
	}

	return f
}

// Apply is a function that applies the format to an element.
//
// Parameters:
//   - trav: The traversor to use for formatting.
//   - elem: The element to format.
//
// Returns:
//   - error: An error if the formatting fails.
//
// Behaviors:
//   - If the traversor is nil, the function does nothing.
func (f *Formatter) Apply(trav *Traversor, elem FStringer) error {
	if trav == nil {
		// Do nothing if the traversor is nil.
		return nil
	}

	err := elem.FString(newTraversor(f, trav.source))
	if err != nil {
		return err
	}

	return nil
}

// ApplyMany is a function that applies the format to multiple elements at once.
//
// Parameters:
//   - trav: The traversor to use for formatting.
//   - elems: The elements to format.
//
// Returns:
//   - error: An error if type Errors.ErrAt if the formatting fails on
//     a specific element.
//
// Behaviors:
//   - If the traversor is nil, the function does nothing.
func (f *Formatter) ApplyMany(trav *Traversor, elems []FStringer) error {
	if trav == nil || len(elems) == 0 {
		// Do nothing if the traversor is nil or if there are no elements.
		return nil
	}

	for i, elem := range elems {
		err := elem.FString(newTraversor(f, trav.source))
		if err != nil {
			return ue.NewErrAt(i+1, "FStringer element", err)
		}
	}

	return nil
}

// ApplyMany is a function that applies the format to multiple elements at once.
//
// Parameters:
//   - trav: The traversor to use for formatting.
//   - elems: The elements to format.
//
// Returns:
//   - error: An error if type Errors.ErrAt if the formatting fails on
//     a specific element.
//
// Behaviors:
//   - If the traversor is nil, the function does nothing.
func ApplyMany[T FStringer](f *Formatter, trav *Traversor, elems []T) error {
	if trav == nil || len(elems) == 0 {
		// Do nothing if the traversor is nil or if there are no elements.
		return nil
	}

	var form *Formatter

	if f == nil {
		form = NewFormatter(nil, nil, nil, nil)
	} else {
		form = f
	}

	for i, elem := range elems {
		err := elem.FString(newTraversor(form, trav.source))
		if err != nil {
			return ue.NewErrAt(i+1, "FStringer element", err)
		}
	}

	return nil
}

// Merge is a function that merges the given formatter with the current one.
//
// Parameters:
//   - other: The other formatter to merge with.
//
// Returns:
//   - *Formatter: A pointer to the new formatter.
//
// Behaviors:
//   - If the current formatter has a value set for a configuration, it will
// 	be used. Otherwise, the value from the other formatter will be used.
func (f *Formatter) Merge(other *Formatter) *Formatter {
	form := &Formatter{}

	if f.indent != nil {
		form.indent = f.indent
	} else {
		form.indent = other.indent
	}

	if f.delimiterLeft != nil {
		form.delimiterLeft = f.delimiterLeft
	} else {
		form.delimiterLeft = other.delimiterLeft
	}

	if f.delimiterRight != nil {
		form.delimiterRight = f.delimiterRight
	} else {
		form.delimiterRight = other.delimiterRight
	}

	if f.separator != nil {
		form.separator = f.separator
	} else {
		form.separator = other.separator
	}

	return form
}

//////////////////////////////////////////////////////////////

/*
// Apply is a method of the Formatter type that creates a formatted string from the given values.
//
// Parameters:
//   - values: The values to format.
//
// Returns:
//   - []string: The formatted string.
func (form *Formatter) Apply(values []string) []string {
	// 1. Add the separator between each value.
	if form.separator != nil {
		values = form.separator.apply(values)
	}

	// 2. Add the left delimiter (if any).
	if form.delimiterLeft != nil {
		values = form.delimiterLeft.applyOnLeft(values)
	}

	// 3. Add the right delimiter (if any).
	if form.delimiterRight != nil {
		values = form.delimiterRight.applyOnRight(values)
	}

	// 4. Apply indentation to all the values.
	if form.indent != nil {
		values = form.indent.apply(values)
	} else {
		values = []string{strings.Join(values, "")}
	}

	return values
}

// ApplyString is a method of the Formatter type that works like Apply but returns a single string.
//
// Parameters:
//   - values: The values to format.
//
// Returns:
//   - string: The formatted string.
func (form *Formatter) ApplyString(values []string) string {
	return strings.Join(form.Apply(values), "\n")
}
*/
