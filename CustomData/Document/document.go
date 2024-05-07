package Document

import (
	"strings"

	ffs "github.com/PlayerR9/MyGoLib/Formatting/FString"
)

// Document is a generic data structure that represents a document.
type Document struct {
	// lines is the lines of the document.
	lines []string
}

// String returns the string representation of the document.
//
// Returns:
//   - string: The string representation of the document.
func (d *Document) String() string {
	return strings.Join(d.lines, "\n")
}

func (d *Document) Tmp() []string {
	return d.lines
}

// FString returns the formatted string representation of the document.
//
// Parameters:
//   - indentLevel: The level of indentation.
//
// Returns:
//   - []string: The formatted string representation of the document.
func (d *Document) FString(indentLevel int) []string {
	indentConfig := ffs.NewIndentConfig(ffs.DefaultIndentation, indentLevel, false)
	indent := indentConfig.String()

	lines := make([]string, 0, len(d.lines))

	for _, line := range d.lines {
		lines = append(lines, indent+line)
	}

	return lines
}

// NewDocument creates a new document.
//
// Parameters:
//   - sentences: The sentences to add to the document.
//
// Returns:
//   - *Document: A pointer to the newly created document.
func NewDocument(sentences ...string) *Document {
	d := &Document{
		lines: make([]string, 0),
	}

	d.AddLine(sentences...)

	return d
}

// AddLine adds sentences to the document separated by a space.
// The line is split by the newline character.
//
// Parameters:
//   - line: The line to add.
//
// Example:
//   - AddLine("Hello,", "world!")
//   - AddLine("This is a sentence.")
func (d *Document) AddLine(sentences ...string) {
	if len(sentences) == 0 {
		return
	}

	var builder strings.Builder

	if strings.HasSuffix(sentences[0], "\n") {
		builder.WriteString(strings.TrimSuffix(sentences[0], "\n"))

		d.lines = append(d.lines, strings.Split(builder.String(), "\n")...)
	} else {
		builder.WriteString(sentences[0])
	}

	for _, sentence := range sentences[1:] {
		if strings.HasSuffix(sentence, "\n") {
			builder.WriteRune(' ')
			builder.WriteString(strings.TrimSuffix(sentence, "\n"))

			d.lines = append(d.lines, strings.Split(builder.String(), "\n")...)
		} else {
			builder.WriteRune(' ')
			builder.WriteString(sentence)
		}
	}

	if builder.Len() != 0 {
		d.lines = append(d.lines, strings.Split(builder.String(), "\n")...)
	}
}
