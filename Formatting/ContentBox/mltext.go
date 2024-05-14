package ContentBox

import (
	"errors"
	"fmt"
	"strings"

	fs "github.com/PlayerR9/MyGoLib/Formatting/Strings"
	sext "github.com/PlayerR9/MyGoLib/Utility/StringExt"
)

type MultiLineText struct {
	text []string
}

func (mlt *MultiLineText) String() string {
	return strings.Join(mlt.text, " ")
}

func NewMultiLineText(text ...string) *MultiLineText {
	return &MultiLineText{
		text: text,
	}
}

func (mlt *MultiLineText) AppendString(text string) {
	mlt.text = append(mlt.text, text)
}

func (mlt *MultiLineText) Draw(width, height int) ([]string, error) {
	ts, err := mlt.GenerateContentBox(width, height)
	if err != nil {
		return nil, fmt.Errorf("could not generate content box: %s", err.Error())
	}

	leftLimit, ok := ts.GetFurthestRightEdge()
	if !ok {
		leftLimit = width
	}

	newLines := make([]string, 0)

	lines := ts.GetLines()

	// First line
	line := lines[0].String()

	var builder strings.Builder

	builder.WriteString(line)
	builder.WriteString(strings.Repeat(Space, leftLimit-len(line)))

	newLines = append(newLines, builder.String())
	builder.Reset()

	// Rest of the lines
	for _, line := range lines[1:] {
		builder.WriteString(strings.Repeat(Space, IndentLevel))

		wl := line.String()

		builder.WriteString(wl)
		builder.WriteString(strings.Repeat(Space, leftLimit-len(wl)))

		newLines = append(newLines, builder.String())
		builder.Reset()
	}

	return newLines, nil
}

func FormatAndSplitText(ts *fs.TextSplit, fields []string, hasIndent bool, width, maxHeight int) ([]string, bool, error) {
	if hasIndent {
		width -= IndentLevel
	}

	numberOfLines, err := fs.CalculateNumberOfLines(fields, width)
	if err != nil {
		line, err := sext.ReplaceSuffix(strings.Join(fields, " ")[:width], Hellip)
		if err != nil {
			return nil, false, fmt.Errorf("could not replace suffix: %s", err.Error())
		}

		ok := ts.InsertWord(line)
		if !ok {
			return nil, false, errors.New("could not insert word")
		}

		return nil, !hasIndent, nil
	}

	if numberOfLines == 1 {
		halfTs, err := fs.SplitInEqualSizedLines(fields, width, numberOfLines)
		if err != nil {
			return nil, false, fmt.Errorf("could not split text: %s", err.Error())
		}

		ok := ts.InsertWord(halfTs.GetFirstLine().String())
		if !ok {
			return nil, false, errors.New("could not insert word")
		}

		return nil, false, nil
	}

	for len(fields) > 0 {
		ok := ts.InsertWord(fields[0])
		if !ok {
			return fields, true, nil
		}

		fields = fields[1:]
	}

	return nil, false, nil
}

// EnqueueContents writes a slice of strings (contents) into the ContentBox
// with a specified style. Each string in the contents is split into fields
// (words) and written into the ContentBox.
// The fields are written with spacing defined by FieldSpacing.
//
// Parameters:
//
//   - contents: The slice of strings to be written.
//   - style: The style to be applied to the strings.
//
// Returns:
//
//   - An error if writing any of the strings fails, otherwise nil.
func (mlt *MultiLineText) GenerateContentBox(width int, height int) (*fs.TextSplit, error) {
	// 1. Create a list of non-empty fields
	fieldList := make([][]string, 0)

	for _, content := range mlt.text {
		fields := strings.Fields(content)

		if len(fields) != 0 {
			fieldList = append(fieldList, fields)
		}
	}
	if len(fieldList) == 0 {
		return fs.NewTextSplit(width, height)
	}

	// 2. Try to optimize lines
	ts, err := fs.NewTextSplit(width-IndentLevel, height)
	if err != nil {
		return nil, fmt.Errorf("could not create text split: %s", err.Error())
	}

	fields := fieldList[0]

	// FIXME: Discard any remaining fields
	_, possibleNewLine, err := FormatAndSplitText(ts, fields, false, width, height)
	if err != nil {
		return nil, fmt.Errorf("could not insert field: %s", err.Error())
	}

	for _, fields := range fieldList[1:] {
		if possibleNewLine {
			for len(fields) > 0 {
				ok := ts.InsertWord(fields[0])
				if !ok {
					break
				}

				fields = fields[1:]
			}
		}

		if len(fields) == 0 {
			continue
		}

		// FIXME: Discard any remaining fields
		_, possibleNewLine, err = FormatAndSplitText(ts, fields, true, width, height)
		if err != nil {
			return nil, fmt.Errorf("could not insert field: %s", err.Error())
		}
	}

	// 3. Final last optimization

	text := make([]string, 0)

	for _, line := range ts.GetLines() {
		text = append(text, line.String())
	}

	optimizedTs, err := fs.SplitInEqualSizedLines(text, width, -1)
	if err != nil {
		return ts, fmt.Errorf("could not split text: %s", err.Error())
	}

	return optimizedTs, nil
}
