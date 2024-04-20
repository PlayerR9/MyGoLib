package FString

import (
	"fmt"
	"strings"
)

// BuildOption is a type that represents a function that modifies the builder.
//
// Parameters:
//
//   - *Builder: The builder to modify.
type BuildOption func(*Builder)

// WithIndentation is a function that sets the indentation configuration of the builder.
//
// Ignore this function if you don't want to indent at all. If you want the rules of
// vertical indentation to be applied, without any indentation, use "" as the indentation
// string.
//
// Parameters:
//
//   - config: The indentation configuration.
//
// Returns:
//
//   - BuildOption: A function that modifies the builder.
func WithIndentation(config *IndentConfig) BuildOption {
	return func(b *Builder) {
		b.indent = config
	}
}

// WithDelimiterLeft is a function that sets the left delimiter configuration of the builder.
//
// Parameters:
//
//   - delimiter: The left delimiter configuration.
//
// Returns:
//
//   - BuildOption: A function that modifies the builder.
func WithDelimiterLeft(delimiter *DelimiterConfig) BuildOption {
	return func(b *Builder) {
		b.delimiterLeft = delimiter
	}
}

// WithDelimiterRight is a function that sets the right delimiter configuration of the builder.
//
// Parameters:
//
//   - delimiter: The right delimiter configuration.
//
// Returns:
//
//   - BuildOption: A function that modifies the builder.
func WithDelimiterRight(delimiter *DelimiterConfig) BuildOption {
	return func(b *Builder) {
		b.delimiterRight = delimiter
	}
}

// WithSeparator is a function that sets the separator configuration of the builder.
//
// Parameters:
//
//   - config: The separator configuration.
//
// Returns:
//
//   - BuildOption: A function that modifies the builder.
func WithSeparator(config *SeparatorConfig) BuildOption {
	return func(b *Builder) {
		b.separator = config
	}
}

// Builder is a type that represents a builder for creating formatted strings.
type Builder struct {
	// The indentation configuration of the builder.
	indent *IndentConfig

	// The left delimiter configuration of the builder.
	delimiterLeft *DelimiterConfig

	// The right delimiter configuration of the builder.
	delimiterRight *DelimiterConfig

	// The separator configuration of the builder.
	separator *SeparatorConfig
}

// NewBuilder is a function that creates a new builder with the given options.
//
// Parameters:
//
//   - options: The options to set on the builder.
//
// Returns:
//
//   - *Builder: A pointer to the new builder.
//
// Information:
//
//   - Options that are not specified will be set to their default values:
//   - ==IndentConfig==
//   - Indentation: DefaultIndentation
//   - InitialLevel: 0
//   - AllowVertical: false
//   - ==SeparatorConfig==
//   - Separator: DefaultSeparator
//   - HasFinalSeparator: false
//   - ==DelimiterConfig==
//   - Value: ""
//   - Inline: true
func NewBuilder(options ...BuildOption) *Builder {
	b := &Builder{}

	for _, option := range options {
		option(b)
	}

	if b.indent == nil {
		b.indent = &IndentConfig{
			Indentation:   DefaultIndentation,
			InitialLevel:  0,
			AllowVertical: false,
		}
	}

	if b.separator == nil {
		b.separator = &SeparatorConfig{
			Separator:         "",
			HasFinalSeparator: false,
		}
	}

	if b.delimiterLeft == nil {
		b.delimiterLeft = &DelimiterConfig{
			Value:  "",
			Inline: true,
		}
	}

	if b.delimiterRight == nil {
		b.delimiterRight = &DelimiterConfig{
			Value:  "",
			Inline: true,
		}
	}

	return b
}

// Build is a method of the Builder type that creates a formatted string from the given values.
//
// Parameters:
//
//   - values: The values to format.
//
// Returns:
//
//   - string: The formatted string.
func (b *Builder) Build(values []string) string {
	// 1. Add the separator between each value.
	vals := make([]string, len(values))
	copy(vals, values)

	if b.separator.HasFinalSeparator {
		for i := 0; i < len(vals); i++ {
			vals[i] = fmt.Sprintf("%s%s", vals[i], b.separator.Separator)
		}
	} else {
		if len(values) > 0 {
			for i := 0; i < len(vals)-1; i++ {
				vals[i] = fmt.Sprintf("%s%s", vals[i], b.separator.Separator)
			}
		}
	}

	if b.indent.AllowVertical {
		indent := b.indent.String()

		if len(vals) == 0 {
			if b.indent.IgnoreFirst {
				return fmt.Sprintf("%s%s", b.delimiterLeft, b.delimiterRight)
			} else {
				return fmt.Sprintf("%s%s%s", indent, b.delimiterLeft, b.delimiterRight)
			}
		}

		var builder strings.Builder

		if b.indent.IgnoreFirst {
			fmt.Fprintf(&builder, "%s\n", b.delimiterLeft)
		} else {
			fmt.Fprintf(&builder, "%s%s\n", indent, b.delimiterLeft)
		}

		for _, value := range vals {
			fmt.Fprintf(&builder, "%s%s%s\n", indent, b.indent.Indentation, value)
		}

		fmt.Fprintf(&builder, "%s%s", indent, b.delimiterRight)

		return builder.String()
	} else {
		indent := b.indent.String()

		if len(vals) == 0 {
			if b.indent.IgnoreFirst {
				return fmt.Sprintf("%s%s", b.delimiterLeft, b.delimiterRight)
			} else {
				return fmt.Sprintf("%s%s%s", indent, b.delimiterLeft, b.delimiterRight)
			}
		}

		// Add the delimiter.
		if b.delimiterLeft.Value != "" {
			if !b.indent.IgnoreFirst {
				vals = append(
					[]string{
						fmt.Sprintf("%s%s", indent, b.delimiterLeft.Value),
					}, vals...,
				)
			} else {
				vals = append([]string{b.delimiterLeft.Value}, vals...)
			}
		}

		if b.delimiterRight.Value != "" {
			vals = append(vals, b.delimiterRight.Value)
		}

		return strings.Join(vals, "")
	}
}

// Format is a method of the Builder type that creates a formatted string from the given values.
//
// Parameters:
//
//   - values: The values to format.
//
// Returns:
//
//   - []string: The formatted string.
func (b *Builder) Format(values []string) []string {
	// 1. Add the separator between each value.
	vals := make([]string, len(values))
	copy(vals, values)

	if b.separator.HasFinalSeparator {
		for i := 0; i < len(vals); i++ {
			vals[i] = fmt.Sprintf("%s%s", vals[i], b.separator.Separator)
		}
	} else {
		if len(values) > 0 {
			for i := 0; i < len(vals)-1; i++ {
				vals[i] = fmt.Sprintf("%s%s", vals[i], b.separator.Separator)
			}
		}
	}

	if b.indent.AllowVertical {
		indent := b.indent.String()

		if len(vals) == 0 {
			if b.indent.IgnoreFirst {
				return []string{fmt.Sprintf("%s%s", b.delimiterLeft, b.delimiterRight)}
			} else {
				return []string{fmt.Sprintf("%s%s%s", indent, b.delimiterLeft, b.delimiterRight)}
			}
		}

		result := make([]string, 0, len(vals))

		if b.indent.IgnoreFirst {
			result = append(result, b.delimiterLeft.String())
		} else {
			result = append(result, fmt.Sprintf("%s%s", indent, b.delimiterLeft))
		}

		for _, value := range vals {
			result = append(result, fmt.Sprintf("%s%s%s", indent, b.indent.Indentation, value))
		}

		result = append(result, fmt.Sprintf("%s%s", indent, b.delimiterRight))

		return result
	} else {
		indent := b.indent.String()

		if len(vals) == 0 {
			if b.indent.IgnoreFirst {
				return []string{fmt.Sprintf("%s%s", b.delimiterLeft, b.delimiterRight)}
			} else {
				return []string{fmt.Sprintf("%s%s%s", indent, b.delimiterLeft, b.delimiterRight)}
			}
		}

		// Add the delimiter.
		if b.delimiterLeft.Value != "" {
			if !b.indent.IgnoreFirst {
				vals = append(
					[]string{
						fmt.Sprintf("%s%s", indent, b.delimiterLeft.Value),
					}, vals...,
				)
			} else {
				vals = append([]string{b.delimiterLeft.Value}, vals...)
			}
		}

		if b.delimiterRight.Value != "" {
			vals = append(vals, b.delimiterRight.Value)
		}

		return vals
	}
}
