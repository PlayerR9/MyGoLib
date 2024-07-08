package FString

// Settinger is an interface that represents the settings for the formatting
// functions.
type Settinger interface{}

// Option is a function that sets the settings for the formatting functions.
//
// Parameters:
//   - Settinger: The settings to set.
type Option func(Settinger)

// ConfigOption is a type that represents a configuration option for a formatter.
type ConfigOption func(*FormatConfig)

// WithModifiedIndent is a function that modifies the indentation level of the formatter
// by a specified amount relative to the current indentation level.
//
// Parameters:
//   - by: The amount by which to modify the indentation level.
//
// Returns:
//   - ConfigOption: The configuration option.
//
// Behaviors:
//   - Negative values will decrease the indentation level while positive values will
//     increase it. If the value is 0, then nothing is done and when the indentation level
//     is 0, it is not decreased.
func WithModifiedIndent(by int) ConfigOption {
	if by == 0 {
		return func(f *FormatConfig) {}
	} else {
		return func(f *FormatConfig) {
			config := f.indentation
			if config == nil {
				return
			}

			config.level += by
			if config.level < 0 {
				config.level = 0
			}
		}
	}
}

// WithLeftDelimiter is a function that modifies the left delimiter of the formatter.
//
// Parameters:
//   - str: The string to use as the left delimiter.
//
// Returns:
//   - ConfigOption: The configuration option.
//
// Behaviors:
//   - If str is empty, then the left delimiter is removed.
func WithLeftDelimiter(str string) ConfigOption {
	if str == "" {
		return func(f *FormatConfig) {
			f.delimiterLeft = nil
		}
	} else {
		newConfig := &DelimiterConfig{
			str:  str,
			left: true,
		}

		return func(f *FormatConfig) {
			f.delimiterLeft = newConfig
		}
	}
}

// WithRightDelimiter is a function that modifies the right delimiter of the formatter.
//
// Parameters:
//   - str: The string to use as the right delimiter.
//
// Returns:
//   - ConfigOption: The configuration option.
//
// Behaviors:
//   - If str is empty, then the right delimiter is removed.
func WithRightDelimiter(str string) ConfigOption {
	if str == "" {
		return func(f *FormatConfig) {
			f.delimiterRight = nil
		}
	} else {
		newConfig := &DelimiterConfig{
			str:  str,
			left: false,
		}

		return func(f *FormatConfig) {
			f.delimiterRight = newConfig
		}
	}
}
