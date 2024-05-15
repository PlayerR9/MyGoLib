package Components

import (
	"strings"

	cdd "github.com/PlayerR9/MyGoLib/ComplexData/Display"

	sext "github.com/PlayerR9/MyGoLib/Utility/StringExt"

	ers "github.com/PlayerR9/MyGoLib/Units/Errors"

	"github.com/gdamore/tcell"
	"github.com/markphelps/optional"
)

type MultiLineText struct {
	lines     [][]string
	style     tcell.Style
	splitChar optional.Rune
}

func (mlt *MultiLineText) Draw(table *cdd.DrawTable) error {
	if table == nil {
		return ers.NewErrNilParameter("table")
	}

	if mlt.IsEmpty() {
		return nil // Nothing to draw
	}

	cb := NewContentBox(mlt.lines)

	char, err := mlt.splitChar.Get()
	if err == nil {
		cb.SetSeparator(char)
	}

	return cb.Draw(table)
}

// Display-related
func (mlt *MultiLineText) SetStyle(style tcell.Style) {
	mlt.style = style
}

type MultiLineTextOption func(*MultiLineText)

func WithSplitChar(splitChar rune) MultiLineTextOption {
	return func(mlt *MultiLineText) {
		mlt.splitChar = optional.NewRune(splitChar)
	}
}

func NewMultiLineText(options ...MultiLineTextOption) *MultiLineText {
	mlt := &MultiLineText{
		lines: make([][]string, 0),
		style: tcell.StyleDefault,
	}

	for _, option := range options {
		option(mlt)
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
	lines := make([]string, 0, len(mlt.lines))

	if mlt.splitChar.Present() {
		for _, line := range mlt.lines {
			lines = append(lines, strings.Join(line, string(mlt.splitChar.MustGet())))
		}
	} else {
		for _, line := range mlt.lines {
			lines = append(lines, strings.Join(line, ""))
		}
	}

	return lines
}
