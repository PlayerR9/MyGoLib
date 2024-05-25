package Printer

import (
	"strings"

	ue "github.com/PlayerR9/MyGoLib/Units/Errors"
)

// Traversor is a type that represents a traversor for a formatted string.
type Traversor struct {
	// The indentation configuration of the builder.
	indentConfig *IndentConfig

	// The left delimiter configuration of the builder.
	delLConfig *DelimiterConfig

	// The right delimiter configuration of the builder.
	delRConfig *DelimiterConfig

	// The sepConfig configuration of the builder.
	sepConfig *SeparatorConfig

	// indentation is the string that is used for indentation
	// on the left side of the traversor.
	indentation string

	// source is the source of the traversor.
	source *Printer
}

// newTraversor creates a new traversor.
//
// Parameters:
//   - config: The configuration of the traversor.
//   - source: The source of the traversor.
//
// Returns:
//   - *Traversor: The new traversor.
func newTraversor(config *Formatter, source *Printer) *Traversor {
	trav := &Traversor{
		indentConfig: config.indent,
		delLConfig:   config.delimiterLeft,
		delRConfig:   config.delimiterRight,
		sepConfig:    config.separator,
		source:       source, // shared source
	}

	if trav.indentConfig == nil {
		trav.indentation = ""
	} else {
		trav.indentation = strings.Repeat(trav.indentConfig.str, trav.indentConfig.level)
	}

	return trav
}

// AppendNBSP appends a non-breaking space to the half-line of the traversor.
//
// Behaviors:
//   - This is equivalent to calling AppendRune('\u00A0') but much more efficient.
func (trav *Traversor) AppendNBSP() {
	if trav.source.isFirstOfLine() && trav.indentConfig != nil {
		trav.source.writeIndent(trav.indentConfig.str)
	}

	trav.source.appendNBSP()
}

// AppendRune appends a rune to the half-line of the traversor.
//
// Parameters:
//   - r: The rune to append.
//
// Returns:
//   - error: An error if the rune could not be appended.
//
// Errors:
//   - *Errors.ErrInvalidRune: If the rune is invalid.
//
// Behaviors:
//   - If the half-line is nil, then a new half-line is created.
func (trav *Traversor) AppendRune(r rune) error {
	if trav.source.isFirstOfLine() && trav.indentConfig != nil {
		trav.source.writeIndent(trav.indentConfig.str)
	}

	err := trav.source.appendRune(trav.indentConfig.str, r)
	if err != nil {
		return err
	}

	return nil
}

// AppendString appends a string to the half-line of the traversor.
//
// Parameters:
//   - str: The string to append.
//
// Returns:
//   - error: An error of type *Errors.ErrInvalidRuneAt if there is an invalid rune
//     in the string.
//
// Behaviors:
//   - IF str is empty: nothing is done.
func (trav *Traversor) AppendString(str string) error {
	if trav.source.isFirstOfLine() && trav.indentConfig != nil {
		trav.source.writeIndent(trav.indentConfig.str)
	}

	err := trav.source.appendString(trav.indentConfig.str, str)
	if err != nil {
		return err
	}

	return nil
}

// AppendStrings appends multiple strings to the half-line of the traversor.
//
// Parameters:
//   - strs: The strings to append.
//
// Returns:
//   - error: An error of type *Errors.ErrAt if there is an error appending a string.
//
// Behaviors:
//   - This is equivalent to calling AppendString for each string in strs but more efficient.
func (trav *Traversor) AppendStrings(strs []string) error {
	if len(strs) == 0 {
		return nil
	}

	if trav.source.isFirstOfLine() && trav.indentConfig != nil {
		trav.source.writeIndent(trav.indentConfig.str)
	}

	err := trav.source.appendString(trav.indentConfig.str, strs[0])
	if err != nil {
		return ue.NewErrAt(0, "string", err)
	}

	for i := 1; i < len(strs); i++ {
		err := trav.source.appendString(trav.indentConfig.str, strs[i])
		if err != nil {
			return ue.NewErrAt(i, "string", err)
		}
	}

	return nil
}

// AppendJoinedString appends a joined string to the half-line of the traversor.
//
// Parameters:
//   - sep: The separator to use.
//   - fields: The fields to join.
//
// Returns:
//   - error: An error of type *Errors.ErrInvalidRuneAt if there is an invalid rune
//     in the string.
//
// Behaviors:
//   - This is equivalent to calling AppendString(strings.Join(fields, sep)).
func (trav *Traversor) AppendJoinedString(sep string, fields ...string) error {
	return trav.AppendString(strings.Join(fields, sep))
}

// AcceptWord is a function that, if there is any in-progress word, then said word is added
// to the source.
func (trav *Traversor) AcceptWord() {
	trav.source.acceptWord()
}

// AcceptLine is a function that accepts the current line of the traversor.
//
// Behaviors:
//   - This also accepts the current word if any.
func (trav *Traversor) AcceptLine() {
	trav.source.acceptLine()
}

// AddLine adds a line to the traversor. If there is any in-progress line, then the line is
// appended to the line before accepting it. Otherwise, a new line with the line is added to
// the source.
//
// Parameters:
//   - line: The line to add.
//
// Returns:
//   - error: An error of type *Errors.ErrAt if there is an error adding the line.
//
// Behaviors:
//   - If line is empty, then an empty line is added to the source.
func (trav *Traversor) AddLine(line string) error {
	trav.source.acceptLine() // Accept the current line if any.

	if trav.indentConfig != nil {
		trav.source.writeIndent(trav.indentConfig.str)
	}

	err := trav.source.appendString(trav.indentConfig.str, line)
	if err != nil {
		return err
	}

	trav.source.acceptLine() // Accept the line.

	return nil
}

// AddLines adds multiple lines to the traversor in a more efficient way than
// adding each line individually.
//
// Parameters:
//   - lines: The lines to add.
//
// Returns:
//   - error: An error of type *Errors.ErrAt if there is an error adding a line.
//
// Behaviors:
//   - If there are no lines, then nothing is done.
func (trav *Traversor) AddLines(lines []string) error {
	if len(lines) == 0 {
		return nil
	}

	trav.source.acceptLine() // Accept the current line if any.

	for i := 0; i < len(lines); i++ {
		if trav.indentConfig != nil {
			trav.source.writeIndent(trav.indentConfig.str)
		}

		err := trav.source.appendString(trav.indentConfig.str, lines[0])
		if err != nil {
			return ue.NewErrAt(i, "line", err)
		}

		trav.source.acceptLine() // Accept the line.
	}

	return nil
}

// AddJoinedLine adds a joined line to the traversor. This is a more efficient way to do
// the same as AddLine(strings.Join(fields, sep)).
//
// Parameters:
//   - sep: The separator to use.
//   - fields: The fields to join.
//
// Returns:
//   - error: An error of type *Errors.ErrInvalidRuneAt if there is an invalid rune
//     in the line.
//
// Behaviors:
//   - If fields is empty, then nothing is done.
func (trav *Traversor) AddJoinedLine(sep string, fields ...string) error {
	return trav.AddLine(strings.Join(fields, sep))
}

// EmptyLine adds an empty line to the traversor. This is a more efficient way to do
// the same as AddLine("") or AddLines([]string{""}).
//
// Behaviors:
//   - If the half-line is not empty, then the half-line is added to the source
//     (half-line is reset) and an empty line is added to the source.
func (trav *Traversor) EmptyLine() {
	trav.source.acceptLine() // Accept the current line if any.

	if trav.indentConfig != nil {
		trav.source.writeIndent(trav.indentConfig.str)
	}

	trav.source.acceptLine() // Accept the line.
}

// ConfigOption is a type that represents a configuration option for a formatter.
type ConfigOption func(*Formatter)

// WithIncreasedIndent is a function that increases the indentation level of the formatter
// by one.
//
// Returns:
//   - ConfigOption: The configuration option.
func WithIncreasedIndent() ConfigOption {
	return func(f *Formatter) {
		if f.indent == nil {
			return
		}

		f.indent.level++
	}
}

// WithDecreasedIndent is a function that decreases the indentation level of the formatter
// by one.
//
// Returns:
//   - ConfigOption: The configuration option.
//
// Behaviors:
//   - If the indentation level is already 0, it is not decreased.
func WithDecreasedIndent() ConfigOption {
	return func(f *Formatter) {
		if f.indent == nil || f.indent.level == 0 {
			return
		}

		f.indent.level--
	}
}

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
		return func(f *Formatter) {}
	} else {
		return func(f *Formatter) {
			if f.indent == nil {
				return
			}

			f.indent.level += by
			if f.indent.level < 0 {
				f.indent.level = 0
			}
		}
	}
}

// GetConfig is a method that returns a copy of the configuration of the traversor.
//
// Parameters:
//   - options: The options to apply to the configuration.
//
// Returns:
//   - *Formatter: A copy of the configuration of the traversor.
func (trav *Traversor) GetConfig(options ...ConfigOption) *Formatter {
	configCopy := &Formatter{
		indent:         trav.indentConfig.Copy().(*IndentConfig),
		delimiterLeft:  trav.delLConfig.Copy().(*DelimiterConfig),
		delimiterRight: trav.delRConfig.Copy().(*DelimiterConfig),
		separator:      trav.sepConfig.Copy().(*SeparatorConfig),
	}

	for _, option := range options {
		option(configCopy)
	}

	return configCopy
}

//////////////////////////////////////////////////////////////

/*
// GetIndent returns the indentation string of the traversor.
//
// Returns:
//   - string: The indentation string of the traversor.
func (trav *Traversor) GetIndent() string {
	if trav.indent == nil {
		return ""
	} else {
		return trav.indentStr
	}
}

// ApplyIndent applies the indentation configuration to a specified string.
//
// Parameters:
//   - str: The string to apply the indentation to.
//
// Returns:
//   - string: The string with the indentation applied.
func (trav *Traversor) ApplyIndent(isFirstLine bool, str string) string {
	if trav.indent == nil || !trav.source.isFirstOfLine() {
		return str
	}

	var builder strings.Builder

	builder.WriteString(trav.indentStr)
	builder.WriteString(str)

	return builder.String()
}
*/

/*
// AddMultiline adds a multiline to the traversor. But first, it accepts any in-progress
// half-line.
//
// Parameters:
//   - mlt: The multiline to add.
//
// Behaviors:
//   - If the multiline is nil, then nothing is done.
func (trav *Traversor) AddMultiline(mlt *cb.MultiLineText) {
	if mlt == nil {
		return
	}

	trav.AcceptHalfLine()
	trav.source.addLine(mlt)
}
*/
