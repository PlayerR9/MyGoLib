package CnsPanel

import (
	"strings"

	fs "github.com/PlayerR9/MyGoLib/Formatting/FString"
)

// Description represents a description for a flag, command,
// or any other object that requires a description.
type Description struct {
	// lines are the lines of the description.
	lines []string
}

// String returns the description as a string.
//
// Returns:
//   - string: The description as a string.
func (d *Description) String() string {
	return strings.Join(d.lines, "\n")
}

// FString generates a formatted string representation of a Description.
//
// Parameters:
//   - indentLevel: The level of indentation to use. Sign is ignored.
//
// Returns:
//   - []string: A slice of strings representing the Description.
func (d *Description) FString(indentLevel int) []string {
	indentCfig := fs.NewIndentConfig(fs.DefaultIndentation, indentLevel, true, false)
	indent := indentCfig.String()

	var builder strings.Builder

	results := make([]string, 0, len(d.lines))

	for _, line := range d.lines {
		builder.Reset()

		builder.WriteString(indent)
		builder.WriteString(line)

		results = append(results, builder.String())
	}

	return results
}

// NewDescription creates a new Description.
//
// Returns:
//   - *Description: A pointer to the new Description.
func NewDescription() *Description {
	return &Description{
		lines: make([]string, 0),
	}
}

// AppendFullLine appends a full line to the description.
//
// Parameters:
//   - line: The line to append.
func (d *Description) AppendFullLine(line string) {
	fields := strings.Split(line, "\n")

	d.lines = append(d.lines, fields...)
}

// Append appends multiple elements to the description in a single line.
//
// Parameters:
//   - elem: The elements of a single line.
func (d *Description) Append(elem ...string) {
	d.lines = append(d.lines, strings.Join(elem, " "))
}

// AddNewline adds a newline to the description.
func (d *Description) AddNewline() {
	d.lines = append(d.lines, "")
}
