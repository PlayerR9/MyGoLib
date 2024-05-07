package FString

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

// SetIndentation is a function that sets the indentation configuration of the builder.
//
// Ignore this function if you don't want to indent at all. If you want the rules of
// vertical indentation to be applied, without any indentation, use "" as the indentation
// string.
//
// Parameters:
//   - config: The indentation configuration.
func (b *Builder) SetIndentation(config *IndentConfig) {
	b.indent = config
}

// SetDelimiterLeft is a function that sets the left delimiter configuration of the builder.
//
// Parameters:
//   - delimiter: The left delimiter configuration.
func (b *Builder) SetDelimiterLeft(delimiter *DelimiterConfig) {
	b.delimiterLeft = delimiter
}

// SetDelimiterRight is a function that sets the right delimiter configuration of the builder.
//
// Parameters:
//   - delimiter: The right delimiter configuration.
func (b *Builder) SetDelimiterRight(delimiter *DelimiterConfig) {
	b.delimiterRight = delimiter
}

// SetSeparator is a function that sets the separator configuration of the builder.
//
// Parameters:
//   - config: The separator configuration.
func (b *Builder) SetSeparator(config *SeparatorConfig) {
	b.separator = config
}

// Build is a method of the Builder type that creates a formatter with the
// configuration of the builder.
//
// Returns:
//   - *Formatter: A pointer to the newly created formatter.
//
// Information:
//   - Options that are not specified will be set to their default values:
//   - ==IndentConfig==
//   - Nil (no indentation is used by default)
//   - ==SeparatorConfig==
//   - Separator: DefaultSeparator
//   - HasFinalSeparator: false
//   - ==DelimiterConfig (Left and Right)==
//   - Nil (no delimiters are used by default)
func (b *Builder) Build() *Formatter {
	var separatorConfig *SeparatorConfig

	if b.separator == nil {
		separatorConfig = NewSeparator(DefaultSeparator, false)
	} else {
		separatorConfig = b.separator
	}

	return &Formatter{
		indent:         b.indent,
		delimiterLeft:  b.delimiterLeft,
		delimiterRight: b.delimiterRight,
		separator:      separatorConfig,
	}
}
