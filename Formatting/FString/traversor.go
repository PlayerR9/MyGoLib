package FString

import "strings"

// Traversor is a type that represents a traversor for a formatted string.
type Traversor struct {
	// The indentation configuration of the traversor.
	indent *IndentConfig

	// lines is the lines of the traversor.
	lines *[]string

	// buffer is the buffer of the traversor.
	buffer []string

	// halfLine is the half line of the traversor.
	halfLine strings.Builder
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
//		- If by is negative, it is converted to a positive value.
//    - If the traversor does not have an indentation configuration, the traversor is returned as is.
func (trav *Traversor) IncreaseIndent(by int) *Traversor {
	if trav.indent == nil {
		return trav
	}

	return &Traversor{
		indent: trav.indent.Increase(by),
		lines:  trav.lines,        // The lines are shared.
		buffer: make([]string, 0), // The buffer is reset.
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
// 	- If the half line is not empty, then the first line is added to the
//    half line, the new half line is added to the buffer (half line is reset),
//    and the rest of the lines are added to the buffer.
//    - If the half line is empty, then the lines are added to the buffer. However,
//      if the lines are empty, then an empty line is added to the buffer.
func (trav *Traversor) AddLines(lines ...string) {
	if len(lines) == 0 {
		if trav.halfLine.Len() == 0 {
			// Add an empty line if the half line is empty.
			trav.buffer = append(trav.buffer, "")
		} else {
			// Accept the half line.
			trav.buffer = append(trav.buffer, trav.halfLine.String())
			trav.halfLine.Reset()
		}

		return
	}

	if trav.halfLine.Len() > 0 {
		trav.halfLine.WriteString(lines[0])

		trav.buffer = append(trav.buffer, trav.halfLine.String())
		trav.halfLine.Reset()

		lines = lines[1:]

		if len(lines) == 0 {
			return
		}
	}

	trav.buffer = append(trav.buffer, lines...)
}

// AppendString appends a string to the half line of the traversor.
//
// Parameters:
//   - str: The string to append.
func (trav *Traversor) AppendString(str string) {
	trav.halfLine.WriteString(str)
}

// AppendStrings appends strings to the half line of the traversor with a separator
// between each string.
//
// Parameters:
//   - separator: The separator between each string.
//   - strs: The strings to append.
//
// Behaviors:
//		- If there are no strings, then nothing is appended.
func (trav *Traversor) AppendStrings(separator string, strs ...string) {
	if len(strs) == 0 {
		return
	}

	trav.halfLine.WriteString(strings.Join(strs, separator))
}

// Apply adds the buffer to the lines of the traversor.
//
// Behaviors:
//		- If the half line is not empty, then the half line is added to the buffer
//    (half line is reset) and the buffer is added to the lines of the traversor.
//    - Each line in the buffer is indented by the indentation configuration.
func (trav *Traversor) Apply() {
	if trav.halfLine.Len() > 0 {
		trav.buffer = append(trav.buffer, trav.halfLine.String())
		trav.halfLine.Reset()
	}

	indent := trav.GetIndent()

	for _, line := range trav.buffer {
		*trav.lines = append(*trav.lines, indent+line)
	}

	trav.buffer = trav.buffer[:0]
}

// EmptyLine adds an empty line to the traversor.
//
// Behaviors:
//		- If the half line is not empty, then the half line is added to the buffer
//    (half line is reset) and an empty line is added to the buffer.
func (trav *Traversor) EmptyLine() {
	if trav.halfLine.Len() > 0 {
		trav.buffer = append(trav.buffer, trav.halfLine.String())
		trav.halfLine.Reset()
	}

	trav.buffer = append(trav.buffer, "")
}
