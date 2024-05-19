package ContentBox

import (
	"strings"

	sext "github.com/PlayerR9/MyGoLib/Utility/StringExt"
)

type MultiLineText struct {
	lines [][]string
}

// Runes returns the content of the unit as a 2D slice of runes
// given the size of the table.
//
// Parameters:
//   - width: The width of the table.
//   - height: The height of the table.
//
// Returns:
//   - [][]rune: The content of the unit as a 2D slice of runes.
//   - error: An error if the content could not be converted to runes.
//
// Behaviors:
//   - Always assume that the width and height are greater than 0. No need to check for
//     this.
//   - Errors are only for critical issues, such as the content not being able to be
//     converted to runes. However, out of bounds or other issues should not error.
//     Instead, the content should be drawn as much as possible before unable to be
//     drawn.
func (mlt *MultiLineText) Runes(width, height int) ([][]rune, error) {
	if mlt.IsEmpty() {
		return nil, nil // Nothing to draw
	}

	cb := NewContentBox(mlt.lines)

	return cb.Runes(width, height)
}

func NewMultiLineText() *MultiLineText {
	mlt := &MultiLineText{
		lines: make([][]string, 0),
	}

	return mlt
}

func (mlt *MultiLineText) AppendSentence(sentence string) error {
	lines, err := sext.AdvancedFieldsSplitter(sentence, IndentLevel)
	if err != nil {
		return err
	}

	mlt.lines = append(mlt.lines, lines...)

	return nil
}

func (mlt *MultiLineText) IsEmpty() bool {
	return len(mlt.lines) == 0
}

func (mlt *MultiLineText) GetLines() []string {
	if len(mlt.lines) == 0 {
		return nil
	}

	lines := make([]string, 0, len(mlt.lines))

	for _, line := range mlt.lines {
		lines = append(lines, strings.Join(line, ""))
	}

	return lines
}
