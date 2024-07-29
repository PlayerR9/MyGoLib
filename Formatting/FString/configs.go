package FString

import (
	"strings"

	uc "github.com/PlayerR9/lib_units/common"
)

const (
	// DefaultIndentation is the default indentation string.
	DefaultIndentation string = "   "

	// DefaultSeparator is the default separator string.
	DefaultSeparator string = ", "
)

var (
	// DefaultIndentationConfig is the default indentation configuration.
	DefaultIndentationConfig *IndentConfig

	// DefaultSeparatorConfig is the default separator configuration.
	DefaultSeparatorConfig *SeparatorConfig

	// DefaultFormatterConfig is the default formatter configuration.
	DefaultFormatterConfig *FormatterConfig
)

func init() {
	DefaultIndentationConfig = NewIndentConfig(DefaultIndentation, 0)
	DefaultSeparatorConfig = NewSeparator(DefaultSeparator, false)
	DefaultFormatterConfig = NewFormatterConfig(3, 1)
}

// IndentConfig is a type that represents the configuration for indentation.
type IndentConfig struct {
	// str is the string that is used for indentation.
	str string

	// InitialLevel is the current indentation level.
	level int
}

// Copy implements the uc.Copier interface.
func (c *IndentConfig) Copy() uc.Copier {
	return &IndentConfig{
		str:   c.str,
		level: c.level,
	}
}

// NewIndentConfig is a function that creates a new indentation configuration.
//
// Parameters:
//   - indentation: The string that is used for indentation.
//   - initialLevel: The initial indentation level.
//
// Returns:
//   - *IndentConfig: A pointer to the new indentation configuration.
//
// Default values:
//
//		==IndentConfig==
//	  - Indentation: DefaultIndentation
//	  - InitialLevel: 0
//
// Behaviors:
//   - If initialLevel is negative, it is set to 0.
//   - If indentation is empty, it is set to DefaultIndentation.
func NewIndentConfig(str string, initialLevel int) *IndentConfig {
	if initialLevel < 0 {
		initialLevel = 0
	}

	if str == "" {
		str = DefaultIndentation
	}

	config := &IndentConfig{
		str:   str,
		level: initialLevel,
	}

	return config
}

// GetIndentation is a method that returns the applied indentation.
//
// Returns:
//   - string: The applied indentation.
func (c *IndentConfig) GetIndentation() string {
	return strings.Repeat(c.str, c.level)
}

// GetIndentStr is a method that returns the indentation string.
//
// Returns:
//   - string: The indentation string.
func (c *IndentConfig) GetIndentStr() string {
	return c.str
}

// SeparatorConfig is a type that represents the configuration for separators.
type SeparatorConfig struct {
	// str is the string that is used as a separator.
	str string

	// includeFinal specifies whether the last element should have a separator.
	includeFinal bool
}

// Copy implements the uc.Copier interface.
func (c *SeparatorConfig) Copy() uc.Copier {
	return &SeparatorConfig{
		str:          c.str,
		includeFinal: c.includeFinal,
	}
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
//
//		==SeparatorConfig==
//	  - Separator: DefaultSeparator
//	  - HasFinalSeparator: false
//
// Behaviors:
//   - If separator is empty, it is set to DefaultSeparator.
func NewSeparator(sep string, includeFinal bool) *SeparatorConfig {
	if sep == "" {
		sep = DefaultSeparator
	}

	return &SeparatorConfig{
		str:          sep,
		includeFinal: includeFinal,
	}
}

// DelimiterConfig is a type that represents the configuration for delimiters.
type DelimiterConfig struct {
	// str is the string that is used as a delimiter.
	str string

	// isInline specifies whether the delimiter should be inline.
	isInline bool

	// left specifies whether the delimiter is on the left side.
	left bool
}

// Copy implements the uc.Copier interface.
func (c *DelimiterConfig) Copy() uc.Copier {
	return &DelimiterConfig{
		str:      c.str,
		isInline: c.isInline,
		left:     c.left,
	}
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
func NewDelimiterConfig(str string, isInline, left bool) *DelimiterConfig {
	return &DelimiterConfig{
		str:      str,
		isInline: isInline,
		left:     left,
	}
}

// FormatterConfig is a type that represents the configuration for formatting.
type FormatterConfig struct {
	// tabSize is the size of the tab.
	tabSize int

	// spacingSize is the size of the spacing.
	spacingSize int
}

// Copy implements the uc.Copier interface.
func (c *FormatterConfig) Copy() uc.Copier {
	return &FormatterConfig{
		tabSize:     c.tabSize,
		spacingSize: c.spacingSize,
	}
}

// NewFormatterConfig is a function that creates a new formatter configuration.
//
// Parameters:
//   - tabSize: The size of the tab.
//   - spacingSize: The size of the spacing.
//
// Returns:
//   - *FormatterConfig: A pointer to the new formatter configuration.
//
// Default values:
//
//		==FormatterConfig==
//	  - TabSize: 3
//	  - SpacingSize: 1
//
// Behaviors:
//   - If tabSize is less than 1, it is set to 3.
//   - If spacingSize is less than 1, it is set to 1.
func NewFormatterConfig(tabSize, spacingSize int) *FormatterConfig {
	if tabSize < 1 {
		tabSize = 3
	}

	if spacingSize < 1 {
		spacingSize = 1
	}

	return &FormatterConfig{
		tabSize:     tabSize,
		spacingSize: spacingSize,
	}
}

//////////////////////////////////////////////////////////////

/*



func (config *IndentConfig) apply(values []string) []string {
	if len(values) == 0 {
		return []string{config.Indentation}
	}

	var builder strings.Builder

	result := make([]string, len(values))
	copy(result, values)

	for i := 0; i < len(result); i++ {
		builder.Reset()

		builder.WriteString(config.Indentation)
		builder.WriteString(result[i])

		result[i] = builder.String()
	}

	return result
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
*/
