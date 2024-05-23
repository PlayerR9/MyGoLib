package FString

import (
	"strings"

	cb "github.com/PlayerR9/MyGoLib/Formatting/ContentBox"
)

const (
	// Hellip is the ellipsis character.
	Hellip string = "..."

	// HellipLen is the length of the ellipsis character.
	HellipLen int = len(Hellip)

	// MarginLeft is the left margin of the content box.
	MarginLeft int = 1
)

// FString is a type that represents a formatted string.
type FString struct {
	// lines is the lines of the formatted string.
	lines []*cb.MultiLineText
}

// String returns the string representation of the traversor.
//
// Returns:
//   - string: The string representation of the traversor.
func (fs *FString) String() string {
	values := make([]string, 0, len(fs.lines))

	for _, line := range fs.lines {
		values = append(values, line.GetLines()...)
	}

	return strings.Join(values, "\n")
}

// NewFString creates a new formatted string.
//
// Returns:
//   - *FString: A pointer to the newly created formatted string.
func NewFString() *FString {
	return &FString{
		lines: make([]*cb.MultiLineText, 0),
	}
}

// Traversor creates a traversor for the formatted string.
//
// Parameters:
//   - indent: The indentation configuration of the traversor.
//
// Returns:
//   - *Traversor: A pointer to the newly created traversor.
//
// Behaviors:
//   - If the indentation configuration is nil, the default indentation configuration is used.
func (fs *FString) Traversor(indent *IndentConfig) *Traversor {
	if indent == nil {
		indent = NewIndentConfig(DefaultIndentation, 0, true)
	}

	return &Traversor{
		indent: indent,
		source: fs,
		buffer: make([]*cb.MultiLineText, 0),
	}
}

func (fs *FString) addLine(mlt *cb.MultiLineText) {
	if mlt == nil {
		return
	}

	fs.lines = append(fs.lines, mlt)
}

func (fs *FString) GetLines() []*cb.MultiLineText {
	return fs.lines
}

/*
func (fs *FString) Boxed(width, height int) ([]string, error) {
	fs.fix()

	all_fields := fs.getAllFields()

	fss := make([]*FString, 0, len(all_fields))

	for _, fields := range all_fields {
		fs := &FString{
			lines: fields,
		}

		fss = append(fss, fs)
	}

	lines := make([]string, 0)

	for _, fs := range fss {
		ts, err := fs.generateContentBox(width, height)
		if err != nil {
			return nil, err
		}

		leftLimit, ok := ts.GetFurthestRightEdge()
		if !ok {
			panic("could not get furthest right edge")
		}

		for _, line := range ts.GetLines() {
			fitted, err := sext.FitString(line.String(), leftLimit)
			if err != nil {
				return nil, err
			}

			lines = append(lines, fitted)
		}
	}

	return lines, nil
}


func (fs *FString) fix() {
	// 1. Fix newline boundaries
	newLines := make([]string, 0)

	for _, line := range fs.lines {
		newFields := strings.Split(line, "\n")

		newLines = append(newLines, newFields...)
	}

	fs.lines = newLines
}

// Must call Fix() before calling this function.
func (fs *FString) getAllFields() [][]string {
	// TO DO: Handle special WHITESPACE characters

	fieldList := make([][]string, 0)

	for _, content := range fs.lines {
		fields := strings.Fields(content)

		if len(fields) != 0 {
			fieldList = append(fieldList, fields)
		}
	}

	return fieldList
}
*/
