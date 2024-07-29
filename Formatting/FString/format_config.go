package FString

import (
	"fmt"

	uc "github.com/PlayerR9/lib_units/common"
)

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
	DefaultFormatter *FormatConfig
)

func init() {
	DefaultFormatter = NewFormatter(
		DefaultIndentationConfig,
		DefaultSeparatorConfig,
	)
}

// FormatterConfig is a type that represents a configuration for the
// general formatter.
type FormatConfig struct {
	// indentation is the configuration for the indentation.
	indentation *IndentConfig

	// delimiterLeft is the configuration for the left delimiter.
	delimiterLeft *DelimiterConfig

	// delimiterRight is the configuration for the right delimiter.
	delimiterRight *DelimiterConfig

	// separator is the configuration for the separator.
	separator *SeparatorConfig

	// format is the configuration for the general formatter.
	format *FormatterConfig
}

// Copy implements the uc.Copier interface.
func (form *FormatConfig) Copy() uc.Copier {
	formCopy := new(FormatConfig)

	if form.indentation != nil {
		configCopy := form.indentation.Copy().(*IndentConfig)
		formCopy.indentation = configCopy
	}

	if form.delimiterLeft != nil {
		configCopy := form.delimiterLeft.Copy().(*DelimiterConfig)
		formCopy.delimiterLeft = configCopy
	}

	if form.delimiterRight != nil {
		configCopy := form.delimiterRight.Copy().(*DelimiterConfig)
		formCopy.delimiterRight = configCopy
	}

	if form.separator != nil {
		configCopy := form.separator.Copy().(*SeparatorConfig)
		formCopy.separator = configCopy
	}

	if form.format != nil {
		configCopy := form.format.Copy().(*FormatterConfig)
		formCopy.format = configCopy
	}

	return formCopy
}

// NewFormatter is a function that creates a new formatter with the given configuration.
//
// Parameters:
//   - options: The configuration for the formatter.
//
// Returns:
//   - form: A pointer to the new formatter.
//
// Behaviors:
//   - The function panics if an invalid configuration type is given. (i.e., not IndentConfig,
//     DelimiterConfig, or SeparatorConfig)
//   - If no formatter configuration is given, the default formatter configuration is used.
func NewFormatter(options ...uc.Copier) *FormatConfig {
	form := new(FormatConfig)

	for _, opt := range options {
		switch opt := opt.(type) {
		case *IndentConfig:
			form.indentation = opt
		case *DelimiterConfig:
			if opt.left {
				form.delimiterLeft = opt
			} else {
				form.delimiterRight = opt
			}
		case *SeparatorConfig:
			form.separator = opt
		case *FormatterConfig:
			form.format = opt
		default:
			panic(fmt.Errorf("invalid configuration type: %T", opt))
		}
	}

	if form.format == nil {
		form.format = DefaultFormatterConfig
	}

	return form
}

// GetTabSize is a function that returns the tab size of the formatter.
//
// Returns:
//   - int: The tab size.
func (form *FormatConfig) GetTabSize() int {
	size := form.format.tabSize
	return size
}

// GetIndentationSize is a function that returns the indentation size of the formatter.
//
// Returns:
//   - int: The indentation size.
func (form *FormatConfig) GetSpacingSize() int {
	size := form.format.spacingSize
	return size
}

/////////////////////////////////////////////////

// ApplyForm is a function that applies the format to an element.
//
// Parameters:
//   - form: The formatter to use for formatting.
//   - trav: The traversor to use for formatting.
//   - elem: The element to format.
//
// Returns:
//   - error: An error if the formatting fails.
//
// Behaviors:
//   - If the traversor is nil, the function does nothing.
func ApplyForm[T FStringer](form *FormatConfig, trav *Traversor, elem T) error {
	if trav == nil {
		// Do nothing if the traversor is nil.
		return nil
	}

	otherTrav := newTraversor(form, trav.source)

	err := elem.FString(otherTrav)
	if err != nil {
		return err
	}

	return nil
}

// ApplyFormMany is a function that applies the format to multiple elements at once.
//
// Parameters:
//   - form: The formatter to use for formatting.
//   - trav: The traversor to use for formatting.
//   - elems: The elements to format.
//
// Returns:
//   - error: An error if type Errors.ErrAt if the formatting fails on
//     a specific element.
//
// Behaviors:
//   - If the traversor is nil, the function does nothing.
func ApplyFormMany[T FStringer](form *FormatConfig, trav *Traversor, elems []T) error {
	if trav == nil || len(elems) == 0 {
		// Do nothing if the traversor is nil or if there are no elements.
		return nil
	}

	otherTrav := newTraversor(form, trav.source)

	for i, elem := range elems {
		err := elem.FString(otherTrav)
		if err != nil {
			return uc.NewErrAt(i+1, "FStringer element", err)
		}
	}

	return nil
}

// ApplyFormFunc is a function that applies the format to an element.
//
// Parameters:
//   - form: The formatter to use for formatting.
//   - trav: The traversor to use for formatting.
//   - elem: The element to format.
//
// Returns:
//   - error: An error if the formatting fails.
//
// Behaviors:
//   - If the traversor is nil, the function does nothing.
func ApplyFormFunc[T any](form *FormatConfig, trav *Traversor, elem T, f FStringFunc[T]) error {
	if trav == nil {
		// Do nothing if the traversor is nil.
		return nil
	}

	otherTrav := newTraversor(form, trav.source)

	err := f(otherTrav, elem)
	if err != nil {
		return err
	}

	return nil
}

// ApplyFormManyFunc is a function that applies the format to multiple elements at once.
//
// Parameters:
//   - form: The formatter to use for formatting.
//   - trav: The traversor to use for formatting.
//   - elems: The elements to format.
//
// Returns:
//   - error: An error if type Errors.ErrAt if the formatting fails on
//     a specific element.
//
// Behaviors:
//   - If the traversor is nil, the function does nothing.
func ApplyFormManyFunc[T any](form *FormatConfig, trav *Traversor, elems []T, f FStringFunc[T]) error {
	if trav == nil || len(elems) == 0 {
		// Do nothing if the traversor is nil or if there are no elements.
		return nil
	}

	otherTrav := newTraversor(form, trav.source)

	for i, elem := range elems {
		err := f(otherTrav, elem)
		if err != nil {
			return uc.NewErrAt(i+1, "FStringer element", err)
		}
	}

	return nil
}

// MergeForm is a function that merges the given formatter with the current one;
// prioritizing the values of the first formatter.
//
// Parameters:
//   - form1: The first formatter.
//   - form2: The second formatter.
//
// Returns:
//   - *FormatConfig: A pointer to the new formatter.
func MergeForm(form1, form2 *FormatConfig) *FormatConfig {
	if form1 == nil {
		var res *FormatConfig

		if form2 == nil {
			res = DefaultFormatter.Copy().(*FormatConfig)
		} else {
			res = form2.Copy().(*FormatConfig)
		}

		return res
	} else if form2 == nil {
		res := form1.Copy().(*FormatConfig)
		return res
	}

	res := new(FormatConfig)

	if form1.indentation != nil {
		res.indentation = form1.indentation
	} else {
		res.indentation = form2.indentation
	}

	if form1.delimiterLeft != nil {
		res.delimiterLeft = form1.delimiterLeft
	} else {
		res.delimiterLeft = form2.delimiterLeft
	}

	if form1.delimiterRight != nil {
		res.delimiterRight = form1.delimiterRight
	} else {
		res.delimiterRight = form2.delimiterRight
	}

	if form1.separator != nil {
		res.separator = form1.separator
	} else {
		res.separator = form2.separator
	}

	if form1.format != nil {
		res.format = form1.format
	} else {
		res.format = form2.format
	}

	return res
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
func (form FormatConfig) Apply(values []string) []string {
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
func (form FormatConfig) ApplyString(values []string) string {
	return strings.Join(form.Apply(values), "\n")
}
*/
