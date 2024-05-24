package FString

import (
	"strings"

	cb "github.com/PlayerR9/MyGoLib/Formatting/ContentBox"
	ue "github.com/PlayerR9/MyGoLib/Units/Errors"
)

// Traversor is a type that represents a traversor for a formatted string.
type Traversor struct {
	// The indentation configuration of the traversor.
	indent *IndentConfig

	// The indent string of the traversor.
	indentStr string

	// source is the source of the traversor.
	source *FString

	// halfLine is the half-line of the traversor.
	halfLine *cb.MultiLineText
}

// IncreaseIndent increases the indentation level of the traversor.
//
// Parameters:
//   - by: The amount by which to increase the indentation level.
//
// Returns:
//   - *Traversor: A pointer to the new traversor.
//
// Behaviors:
//   - If by is negative, it is converted to a positive value.
//   - If the traversor does not have an indentation configuration, the traversor is returned as is.
func (trav *Traversor) IncreaseIndent(by int) *Traversor {
	if trav.indent == nil {
		return trav
	}

	return &Traversor{
		indent:    trav.indent.Increase(by),
		indentStr: trav.indent.String(),
		source:    trav.source, // The lines are shared.
		halfLine:  nil,         // The half-line is reset.
	}
}

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

// ApplyIndent applies the indentation configuration to the traversor.
//
// Parameters:
//   - str: The string to apply the indentation to.
//
// Returns:
//   - string: The string with the indentation applied.
func (trav *Traversor) ApplyIndent(isFirstLine bool, str string) string {
	if trav.indent == nil ||
		(trav.indent.IgnoreFirst && isFirstLine) {
		return str
	}

	var builder strings.Builder

	builder.WriteString(trav.indentStr)
	builder.WriteString(str)

	return builder.String()
}

// AppendRune appends a rune to the half-line of the traversor.
//
// Parameters:
//   - r: The rune to append.
//
// Behaviors:
//   - If the half-line is nil, then a new half-line is created.
func (trav *Traversor) AppendRune(r rune) {
	if trav.halfLine == nil {
		trav.halfLine = cb.NewMultiLineText()
	}

	trav.halfLine.AppendRune(r)
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
	if str == "" {
		return nil
	}

	if trav.halfLine == nil {
		trav.halfLine = cb.NewMultiLineText()

		str = trav.ApplyIndent(true, str)
	}

	return trav.halfLine.AppendSentence(str)
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
//   - If there are no strings, then nothing is done.
func (trav *Traversor) AppendStrings(strs ...string) error {
	if len(strs) == 0 {
		return nil
	}

	for i, str := range strs {
		err := trav.AppendString(str)
		if err != nil {
			return ue.NewErrAt(i+1, "string", err)
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

// AcceptHalfLine is a function that, if there is any in-progress half-line,
// then said half-line is added to the source; resetting it in the process.
func (trav *Traversor) AcceptHalfLine() {
	if trav.halfLine == nil {
		return
	}

	trav.source.addLine(trav.halfLine)
	trav.halfLine = nil
}

// addLineToBuffer is a private function that adds a line to the source of the traversor.
//
// Parameters:
//   - line: The line to add.
//
// Returns:
//   - error: An error of type *Errors.ErrInvalidRuneAt if there is an invalid rune
//     in the line.
func (trav *Traversor) addLineToBuffer(line string) error {
	if line == "" {
		trav.source.addLine(cb.NewMultiLineText())
	} else {
		mlt := cb.NewMultiLineText()

		err := mlt.AppendSentence(trav.ApplyIndent(true, line))
		if err != nil {
			return err
		}

		trav.source.addLine(mlt)
	}

	return nil
}

// AddLine adds a line to the traversor. If there is any in-progress half-line, then the
// line is appended to the half-line before accepting it. Otherwise, a new line with the
// line is added to the source.
//
// Parameters:
//   - line: The line to add.
//
// Returns:
//   - error: An error of type *Errors.ErrInvalidRuneAt if there is an invalid rune
//     in the line.
//
// Behaviors:
//   - When line is empty, however, if there is any in-progress half-line, then the half-line
//     is accepted as is. Otherwise, an empty line is added to the source.
func (trav *Traversor) AddLine(line string) error {
	if trav.halfLine == nil {
		return trav.addLineToBuffer(line)
	}

	err := trav.halfLine.AppendSentence(line)
	if err != nil {
		return err
	}

	trav.AcceptHalfLine()
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

	at := 0

	if trav.halfLine != nil {
		err := trav.halfLine.AppendSentence(lines[at])
		if err != nil {
			return ue.NewErrAt(at, "line", err)
		}

		trav.AcceptHalfLine()

		at++
	}

	for i := at; i < len(lines); i++ {
		err := trav.addLineToBuffer(lines[i])
		if err != nil {
			return ue.NewErrAt(i+1, "line", err)
		}
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

// EmptyLine adds an empty line to the traversor. This is a more efficient way to do
// the same as AddLine("") or AddLines([]string{""}).
//
// Behaviors:
//   - If the half-line is not empty, then the half-line is added to the source
//     (half-line is reset) and an empty line is added to the source.
func (trav *Traversor) EmptyLine() {
	trav.AcceptHalfLine()
	trav.source.addLine(cb.NewMultiLineText())
}
