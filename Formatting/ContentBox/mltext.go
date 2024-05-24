package ContentBox

import (
	"strings"

	sext "github.com/PlayerR9/MyGoLib/Utility/StringExt"
)

// MultiLineText is a unit that represents a multi-line text.
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

// AppendRune appends a rune to the multi-line text.
//
// Parameters:
//   - r: The rune to append.
//
// Behaviors:
//   - If the rune is a newline, then a new line is created.
func (mlt *MultiLineText) AppendRune(r rune) {
	if r == '\n' {
		mlt.lines = append(mlt.lines, []string{})
		return
	}

	if len(mlt.lines) == 0 {
		mlt.lines = append(mlt.lines, []string{})
	}

	mlt.lines[len(mlt.lines)-1] = append(mlt.lines[len(mlt.lines)-1], string(r))

	return
}

// AppendSentence appends a sentence to the multi-line text.
//
// Parameters:
//   - sentence: The sentence to append.
//
// Returns:
//   - error: An error of type *Errors.ErrInvalidRuneAt if there is an invalid rune.
//
// Behaviors:
//   - If the sentence is empty, then nothing is appended.
func (mlt *MultiLineText) AppendSentence(sentence string) error {
	if sentence == "" {
		return nil
	}

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
