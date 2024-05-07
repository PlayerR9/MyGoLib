package FString

import (
	"strings"
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
