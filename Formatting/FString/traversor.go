package FString

import (
	"strings"

	cb "github.com/PlayerR9/MyGoLib/Formatting/ContentBox"
)

// Traversor is a type that represents a traversor for a formatted string.
type Traversor struct {
	// The indentation configuration of the traversor.
	indent *IndentConfig

	// source is the source of the traversor.
	source *FString

	// buffer is the buffer of the traversor.
	buffer []*cb.MultiLineText

	// halfLine is the half line of the traversor.
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
		indent:   trav.indent.Increase(by),
		source:   trav.source,                  // The lines are shared.
		buffer:   make([]*cb.MultiLineText, 0), // The buffer is reset.
		halfLine: nil,                          // The half line is reset.
	}
}

// GetIndent returns the indentation string of the traversor.
//
// Returns:
//   - string: The indentation string of the traversor.
func (trav *Traversor) GetIndent() string {
	if trav.indent == nil {
		return ""
	}

	return strings.Repeat(trav.indent.Indentation, trav.indent.InitialLevel)
}

// AddLines adds lines to the traversor.
//
// Parameters:
//   - lines: The lines to add.
//
// Behaviors:
//   - If the half line is not empty, then the first line is added to the
//     half line, the new half line is added to the buffer (half line is reset),
//     and the rest of the lines are added to the buffer.
//   - If the half line is empty, then the lines are added to the buffer. However,
//     if the lines are empty, then an empty line is added to the buffer.
func (trav *Traversor) AddLines(lines ...string) {
	if len(lines) == 0 {
		if trav.halfLine == nil {
			// Add an empty line if the half line is empty.
			trav.buffer = append(trav.buffer, cb.NewMultiLineText())
		} else {
			// Accept the half line.
			trav.buffer = append(trav.buffer, trav.halfLine)
			trav.halfLine = nil
		}

		return
	}

	if trav.halfLine != nil {
		err := trav.halfLine.AppendSentence(lines[0])
		if err != nil {
			panic(err)
		}

		trav.buffer = append(trav.buffer, trav.halfLine)
		trav.halfLine = nil

		lines = lines[1:]

		if len(lines) == 0 {
			return
		}
	}

	for _, line := range lines {
		mlt := cb.NewMultiLineText()

		err := mlt.AppendSentence(line)
		if err != nil {
			panic(err)
		}

		trav.buffer = append(trav.buffer, mlt)
	}
}

// AppendString appends a string to the half line of the traversor.
//
// Parameters:
//   - str: The string to append.
func (trav *Traversor) AppendString(str string) {
	if trav.halfLine == nil {
		trav.halfLine = cb.NewMultiLineText()
	}

	err := trav.halfLine.AppendSentence(str)
	if err != nil {
		panic(err)
	}
}

// AppendStrings appends strings to the half line of the traversor with a separator
// between each string.
//
// Parameters:
//   - separator: The separator between each string.
//   - strs: The strings to append.
//
// Behaviors:
//   - If there are no strings, then nothing is appended.
func (trav *Traversor) AppendStrings(separator string, strs ...string) {
	if len(strs) == 0 {
		return
	}

	if trav.halfLine == nil {
		trav.halfLine = cb.NewMultiLineText()
	}

	err := trav.halfLine.AppendSentence(strings.Join(strs, separator))
	if err != nil {
		panic(err)
	}
}

// Apply adds the buffer to the lines of the traversor.
//
// Behaviors:
//   - If the half line is not empty, then the half line is added to the buffer
//     (half line is reset) and the buffer is added to the lines of the traversor.
//   - Each line in the buffer is indented by the indentation configuration.
func (trav *Traversor) Apply() {
	if trav.halfLine != nil {
		trav.buffer = append(trav.buffer, trav.halfLine)
		trav.halfLine = nil
	}

	indent := trav.GetIndent()

	for _, line := range trav.buffer {
		lines := line.GetLines()
		mlt := cb.NewMultiLineText()

		for _, line := range lines {
			err := mlt.AppendSentence(indent + line)
			if err != nil {
				panic(err)
			}
		}

		trav.source.addLine(mlt)
	}

	trav.buffer = trav.buffer[:0]
}

// EmptyLine adds an empty line to the traversor.
//
// Behaviors:
//   - If the half line is not empty, then the half line is added to the buffer
//     (half line is reset) and an empty line is added to the buffer.
func (trav *Traversor) EmptyLine() {
	if trav.halfLine != nil {
		trav.buffer = append(trav.buffer, trav.halfLine)
		trav.halfLine = nil
	}

	trav.buffer = append(trav.buffer, cb.NewMultiLineText())
}
