package ContentBox

import (
	"fmt"
	"strings"

	fs "github.com/PlayerR9/MyGoLib/Formatting/Strings"

	sext "github.com/PlayerR9/MyGoLib/Utility/StringExt"
)

const (
	// Asterisks are the characters used emphasize the title.
	Asterisks string = "***"

	// AsterisksLen is the length of the asterisks.
	AsterisksLen int = len(Asterisks)

	// TitleMinWidth is the minimum width of the title.
	TitleMinWidth int = 2 * (AsterisksLen + 1) // +1 for the space between the asterisks and the title
)

// TitleBox is a type that represents a title box.
type TitleBox struct {
	// content is the content of the title box.
	content string
}

// NewTitleBox creates a new title box with the given content.
//
// Parameters:
//   - content: The content of the title box.
//
// Returns:
//   - *TitleBox: A pointer to the newly created title box.
func NewTitleBox(content string) *TitleBox {
	return &TitleBox{
		content: content,
	}
}

// writeTitle is a helper method that writes the title to the title box.
//
// Parameters:
//   - width: The width of the title box.
//
// Returns:
//   - []string: The lines of the title.
//   - error: An error if it occurs during the writing.
func (tb *TitleBox) writeTitle(width int) ([]string, error) {
	// FIXME: Find a better function instead of strings.Fields
	contents := strings.Fields(tb.content)

	numberOfLines, err := fs.CalculateNumberOfLines(contents, width-TitleMinWidth)
	if err != nil {
		return nil, fmt.Errorf("could not calculate number of lines: %s", err.Error())
	}

	ts, err := fs.SplitInEqualSizedLines(contents, width-TitleMinWidth, numberOfLines)
	if err != nil {
		return nil, fmt.Errorf("could not split text in equal sized lines: %s", err.Error())
	}

	lines := make([]string, 0, ts.GetHeight())
	var builder strings.Builder

	for _, line := range ts.GetLines() {
		startPos := (width - (line.Length() + TitleMinWidth)) / 2

		// Add white spaces to the left
		builder.WriteString(strings.Repeat(" ", startPos))

		// Add the title
		builder.WriteString(Asterisks)
		builder.WriteRune(' ')
		builder.WriteString(line.String())
		builder.WriteRune(' ')
		builder.WriteString(Asterisks)

		// Add white spaces to the right
		builder.WriteString(strings.Repeat(" ", width-(line.Length()+TitleMinWidth)-startPos))

		lines = append(lines, builder.String())
		builder.Reset()
	}

	return lines, nil
}

// DrawTitle is a method of TitleBox that draws the title.
//
// Parameters:
//   - width: The width of the title box.
//
// Returns:
//   - []string: The lines of the title.
//   - error: An error if it occurs during the drawing.
func (tb *TitleBox) DrawTitle(width int) ([]string, error) {
	// 1. Generate the lines
	lines, err := tb.writeTitle(width)
	if err == nil {
		return lines, nil
	}

	// Try to truncate the title
	title, err := sext.ReplaceSuffix(tb.content[:width-TitleMinWidth], Hellip)
	if err != nil {
		return nil, fmt.Errorf("could not draw title: %s", err.Error())
	}

	var builder strings.Builder

	builder.WriteString(Asterisks)
	builder.WriteRune(' ')
	builder.WriteString(title)
	builder.WriteRune(' ')
	builder.WriteString(Asterisks)

	return []string{builder.String()}, nil
}
