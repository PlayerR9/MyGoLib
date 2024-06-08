package FString

import (
	"fmt"

	ue "github.com/PlayerR9/MyGoLib/Units/errors"
)

/////////////////////////////////////////////////

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
	DefaultFormatter FormatConfig = NewFormatter(
		DefaultIndentationConfig,
		DefaultSeparatorConfig,
	)
)

// FormatConfig is a type that represents a configuration for formatting.
// [Indentation] [Left Delimiter] [Right Delimiter] [Separator]
type FormatConfig [4]Configer

const (
	// ConfInd_Idx is the index for the indentation configuration.
	ConfInd_Idx = iota

	// ConfDelL_Idx is the index for the left delimiter configuration.
	ConfDelL_Idx

	// ConfDelR_Idx is the index for the right delimiter configuration.
	ConfDelR_Idx

	// ConfSep_Idx is the index for the separator configuration.
	ConfSep_Idx
)

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
func NewFormatter(options ...Configer) (form FormatConfig) {
	if len(options) == 0 {
		return
	}

	for _, opt := range options {
		switch opt := opt.(type) {
		case *IndentConfig:
			form[ConfInd_Idx] = opt
		case *DelimiterConfig:
			if opt.left {
				form[ConfDelL_Idx] = opt
			} else {
				form[ConfDelR_Idx] = opt
			}
		case *SeparatorConfig:
			form[ConfSep_Idx] = opt
		default:
			panic(fmt.Errorf("invalid configuration type: %T", opt))
		}
	}

	return
}

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
func ApplyForm[T FStringer](form FormatConfig, trav *Traversor, elem T) error {
	if trav == nil {
		// Do nothing if the traversor is nil.
		return nil
	}

	err := elem.FString(newTraversor(form, trav.source))
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
func ApplyFormMany[T FStringer](form FormatConfig, trav *Traversor, elems []T) error {
	if trav == nil || len(elems) == 0 {
		// Do nothing if the traversor is nil or if there are no elements.
		return nil
	}

	for i, elem := range elems {
		err := elem.FString(newTraversor(form, trav.source))
		if err != nil {
			return ue.NewErrAt(i+1, "FStringer element", err)
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
func ApplyFormFunc[T FStringer](form FormatConfig, trav *Traversor, elem T, f FStringFunc[T]) error {
	if trav == nil {
		// Do nothing if the traversor is nil.
		return nil
	}

	err := f(newTraversor(form, trav.source), elem)
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
func ApplyFormManyFunc[T FStringer](form FormatConfig, trav *Traversor, elems []T, f FStringFunc[T]) error {
	if trav == nil || len(elems) == 0 {
		// Do nothing if the traversor is nil or if there are no elements.
		return nil
	}

	for i, elem := range elems {
		err := f(newTraversor(form, trav.source), elem)
		if err != nil {
			return ue.NewErrAt(i+1, "FStringer element", err)
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
//   - FormatConfig: A pointer to the new formatter.
func MergeForm(form1, form2 FormatConfig) FormatConfig {
	var form FormatConfig

	for i := 0; i < 4; i++ {
		if form1[i] != nil {
			form[i] = form1[i]
		} else {
			form[i] = form2[i]
		}
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
