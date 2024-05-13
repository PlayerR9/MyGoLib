package FString

import (
	"errors"
	"fmt"
	"strings"

	splt "github.com/PlayerR9/MyGoLib/Formatting/Strings"
	sext "github.com/PlayerR9/MyGoLib/Utility/StringExt"
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
	lines []string
}

// String returns the string representation of the traversor.
//
// Returns:
//   - string: The string representation of the traversor.
func (fs *FString) String() string {
	return strings.Join(fs.lines, "\n")
}

// NewFString creates a new formatted string.
//
// Returns:
//   - *FString: A pointer to the newly created formatted string.
func NewFString() *FString {
	return &FString{
		lines: make([]string, 0),
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
		lines:  &fs.lines,
		buffer: make([]string, 0),
	}
}

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
		for index := 0; index < len(fs.lines); index++ {
			if fs.lines[index] == "" {
				lines = append(lines, strings.Repeat(" ", width))
			} else {
				ts, err := fs.generateContentBox(width, height, index)
				if err != nil {
					return nil, err
				}

				leftLimit := ts.GetFurthestRightEdge()

				for _, line := range ts.Lines {
					fitted, err := sext.FitString(line.String(), leftLimit)
					if err != nil {
						return nil, err
					}

					lines = append(lines, fitted)
				}
			}
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

	for _, content := range fs.lines {
		fields := strings.Fields(content)
	}

	fieldList := make([][]string, 0)

	for _, content := range fs.lines {
		fields := strings.Fields(content)

		if len(fields) != 0 {
			fieldList = append(fieldList, fields)
		}
	}

	return fieldList
}

func (fs *FString) insertField(fields []string, width, height int) (bool, error) {
	ts, err := splt.NewTextSplit(width, height)
	if err != nil {
		return false, err
	}

	// 1. Calculate the number of lines necessary to fit the fields in the width.
	numberOfLines, err := splt.CalculateNumberOfLines(fields, width)
	if err != nil {
		// 2. If it is not possible to fit the entire field in one line, then
		// we must trim it and add an ellipsis character.
		line, err := sext.ReplaceSuffix(strings.Join(fields, " ")[:width], Hellip)
		if err != nil {
			return false, errors.New("line is too short to fit the ellipsis character")
		}

		// 3. Insert the line in the text split.
		ok := ts.InsertWord(line)
		if !ok {
			panic("could not insert word")
		}

		return true, nil
	}

	if numberOfLines > 1 {
		splt := splt.NewSpltLine(fields[0])

		// Keep adding words to the line until it is not possible to fit the next word.
		possibleNewLine := false

		for _, field := range fields[1:] {
			if splt.Len+1+len(field)+HellipLen > width {
				splt.InsertWord(field + Hellip)

				possibleNewLine = true
				break
			}

			splt.InsertWord(field)
		}

		// 4. Insert the line in the text split.
		ok := ts.InsertWord(splt.Line[0])
		if !ok {
			panic("could not insert word")
		}

		return possibleNewLine, nil
	}

	// Otherwise, use the text splitter to split the fields in equal-sized lines.

	tsr, err := splt.NewTextSplitter(width)
	if err != nil {
		return false, err
	}

	err = tsr.SetHeight(numberOfLines)
	if err != nil {
		return false, err
	}

	err = tsr.SplitInEqualSizedLines(fields)
	if err != nil {
		return false, err
	}

	halfTs, err := tsr.GetSolution()
	if err != nil {
		return false, err
	}

	ok := ts.InsertWord(halfTs.GetFirstLine().String())
	if !ok {
		panic("could not insert word")
	}

	return false, nil
}

func (fs *FString) generateContentBox(width, height int, index int) (*splt.TextSplit, error) {
	fieldList := fs.lines

	// Try to optimize lines

	ts, err := splt.NewTextSplit(width-MarginLeft, height)
	if err != nil {
		return nil, fmt.Errorf("width cannot be less than %d", MarginLeft)
	}

	fields := fieldList[index]

	possibleNewLine, err := fs.insertField(fields, width-MarginLeft, height)
	if err != nil {
		return nil, err
	}

	for _, fields := range fieldList[index+1:] {
		if possibleNewLine {
			for len(fields) > 0 {
				if !ts.CanInsertWord(fields[0], len(ts.Lines)-1) {
					break
				}

				ts.InsertWordAt(fields[0], len(ts.Lines)-1)
				fields = fields[1:]
			}
		}

		if len(fields) == 0 {
			continue
		}

		possibleNewLine, err = fs.insertField(fields, width-MarginLeft, height)
		if err != nil {
			return nil, err
		}
	}

	// 3. Final last optimization

	text := make([]string, 0, len(ts.Lines))
	for _, line := range ts.Lines {
		text = append(text, line.String())
	}

	tsr, err := splt.NewTextSplitter(width)
	if err != nil {
		return nil, err
	}

	err = tsr.SplitInEqualSizedLines(text)
	if err != nil {
		return nil, err
	}

	optimizedTs, err := tsr.GetSolution()
	if err != nil {
		return nil, err
	}

	return optimizedTs, nil
}
