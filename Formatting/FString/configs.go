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

	// AllowVertical specifies whether vertical indentation is allowed.
	AllowVertical bool
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
func NewIndentConfig(indentation string, initialLevel int, allowVertical bool, ignoreFirst bool) *IndentConfig {
	config := &IndentConfig{
		Indentation:   indentation,
		AllowVertical: allowVertical,
		IgnoreFirst:   ignoreFirst,
	}

	if initialLevel < 0 {
		config.InitialLevel = -initialLevel
	} else {
		config.InitialLevel = initialLevel
	}

	return config
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
func NewSeparator(separator string, hasFinalSeparator bool) *SeparatorConfig {
	return &SeparatorConfig{
		Separator:         separator,
		HasFinalSeparator: hasFinalSeparator,
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
func NewDelimiterConfig(value string, inline bool) *DelimiterConfig {
	return &DelimiterConfig{
		Value:  value,
		Inline: inline,
	}
}
