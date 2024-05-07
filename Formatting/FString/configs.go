package FString

import "strings"

// IndentConfig is a type that represents the configuration for indentation.
type IndentConfig struct {
	// IgnoreFirst specifies whether the first line should be indented.
	IgnoreFirst bool

	// Indentation is the string that is used for indentation.
	Indentation string

	// InitialLevel is the initial indentation level.
	InitialLevel int
}

// String is a method of fmt.Stringer interface.
//
// Returns:
//   - string: A string representation of the indentation configuration.
func (c *IndentConfig) String() string {
	return strings.Repeat(c.Indentation, c.InitialLevel)
}

// NewIndentConfig is a function that creates a new indentation configuration.
//
// Parameters:
//   - indentation: The string that is used for indentation.
//   - initialLevel: The initial indentation level.
//   - allowVertical: Whether vertical indentation is allowed.
//   - ignoreFirst: Whether the first line should be indented.
//
// Returns:
//   - *IndentConfig: A pointer to the new indentation configuration.
//
// Default values:
// 	==IndentConfig==
//   - Indentation: DefaultIndentation
//   - InitialLevel: 0
//   - IgnoreFirst: true
func NewIndentConfig(indentation string, initialLevel int, ignoreFirst bool) *IndentConfig {
	config := &IndentConfig{
		Indentation: indentation,
		IgnoreFirst: ignoreFirst,
	}

	if initialLevel < 0 {
		config.InitialLevel = -initialLevel
	} else {
		config.InitialLevel = initialLevel
	}

	return config
}

func (config *IndentConfig) apply(values []string) []string {
	if len(values) == 0 {
		if !config.IgnoreFirst {
			return []string{config.Indentation}
		}

		return []string{""}
	}

	var builder strings.Builder

	result := make([]string, len(values))
	copy(result, values)

	if !config.IgnoreFirst {
		builder.WriteString(config.Indentation)
		builder.WriteString(result[0])

		result[0] = builder.String()
	}

	if len(result) == 1 {
		return result
	}

	for i := 1; i < len(result); i++ {
		builder.Reset()

		builder.WriteString(config.Indentation)
		builder.WriteString(result[i])

		result[i] = builder.String()
	}

	return result
}

// SeparatorConfig is a type that represents the configuration for separators.
type SeparatorConfig struct {
	// Separator is the string that is used as a separator.
	Separator string

	// HasFinalSeparator specifies whether the last element should have a separator.
	HasFinalSeparator bool
}

// String is a method of fmt.Stringer interface.
//
// Returns:
//   - string: A string representation of the separator configuration.
func (c *SeparatorConfig) String() string {
	return c.Separator
}

// NewSeparator is a function that creates a new separator configuration.
//
// Parameters:
//   - separator: The string that is used as a separator.
//   - hasFinalSeparator: Whether the last element should have a separator.
//
// Returns:
//   - *SeparatorConfig: A pointer to the new separator configuration.
//
// Default values:
// 	==SeparatorConfig==
//   - Separator: DefaultSeparator
//   - HasFinalSeparator: false
func NewSeparator(separator string, hasFinalSeparator bool) *SeparatorConfig {
	return &SeparatorConfig{
		Separator:         separator,
		HasFinalSeparator: hasFinalSeparator,
	}
}

func (config *SeparatorConfig) apply(values []string) []string {
	switch len(values) {
	case 0:
		if config.HasFinalSeparator {
			return []string{config.Separator}
		}

		return []string{}
	case 1:
		var builder strings.Builder

		builder.WriteString(values[0])

		if config.HasFinalSeparator {
			builder.WriteString(config.Separator)
		}

		return []string{builder.String()}
	default:
		result := make([]string, len(values))
		copy(result, values)

		var builder strings.Builder

		builder.WriteString(result[0])
		builder.WriteString(config.Separator)

		result[0] = builder.String()

		for i := 1; i < len(result)-1; i++ {
			builder.Reset()

			builder.WriteString(result[i])
			builder.WriteString(config.Separator)
			result[i] = builder.String()
		}

		if config.HasFinalSeparator {
			builder.Reset()

			builder.WriteString(result[len(result)-1])
			builder.WriteString(config.Separator)
			result[len(result)-1] = builder.String()
		}

		return result
	}
}

// DelimiterConfig is a type that represents the configuration for delimiters.
type DelimiterConfig struct {
	// Value is the string that is used as a delimiter.
	Value string

	// Inline specifies whether the delimiter should be inline.
	Inline bool
}

// String is a method of fmt.Stringer interface.
//
// Returns:
//   - string: A string representation of the delimiter configuration.
func (c *DelimiterConfig) String() string {
	return c.Value
}

// NewDelimiterConfig is a function that creates a new delimiter configuration.
//
// Parameters:
//   - value: The string that is used as a delimiter.
//   - inline: Whether the delimiter should be inline.
//
// Returns:
//   - *DelimiterConfig: A pointer to the new delimiter configuration.
//
// Default values:
//   - ==DelimiterConfig==
//   - Value: ""
//   - Inline: true
func NewDelimiterConfig(value string, inline bool) *DelimiterConfig {
	return &DelimiterConfig{
		Value:  value,
		Inline: inline,
	}
}

func (config *DelimiterConfig) applyOnLeft(values []string) []string {
	if len(values) == 0 {
		return []string{config.Value}
	}

	result := make([]string, len(values))
	copy(result, values)

	if config.Inline {
		var builder strings.Builder

		builder.WriteString(config.Value)
		builder.WriteString(values[0])

		result[0] = builder.String()
	} else {
		result = append([]string{config.Value}, result...)
	}

	return result
}

func (config *DelimiterConfig) applyOnRight(values []string) []string {
	if len(values) == 0 {
		return []string{config.Value}
	}

	result := make([]string, len(values))
	copy(result, values)

	if config.Inline {
		var builder strings.Builder

		builder.WriteString(values[len(values)-1])
		builder.WriteString(config.Value)

		result[len(values)-1] = builder.String()
	} else {
		result = append(result, config.Value)
	}

	return result
}
