package FString

import "strings"

type IndentConfig struct {
	IgnoreFirst   bool
	Indentation   string
	InitialLevel  int
	AllowVertical bool
}

func (c *IndentConfig) String() string {
	return strings.Repeat(c.Indentation, c.InitialLevel)
}

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

type SeparatorConfig struct {
	Separator         string
	HasFinalSeparator bool
}

func (c *SeparatorConfig) String() string {
	return c.Separator
}

func NewSeparator(separator string, hasFinalSeparator bool) *SeparatorConfig {
	return &SeparatorConfig{
		Separator:         separator,
		HasFinalSeparator: hasFinalSeparator,
	}
}

type DelimiterConfig struct {
	Value  string
	Inline bool
}

func (c *DelimiterConfig) String() string {
	return c.Value
}

func NewDelimiterConfig(value string, inline bool) *DelimiterConfig {
	return &DelimiterConfig{
		Value:  value,
		Inline: inline,
	}
}
