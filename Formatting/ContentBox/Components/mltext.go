package Components

import (
	cdd "github.com/PlayerR9/MyGoLib/ComplexData/Display"
	sx "github.com/PlayerR9/MyGoLib/Formatting/ContentBox/String"

	sext "github.com/PlayerR9/MyGoLib/Utility/StringExt"

	ers "github.com/PlayerR9/MyGoLib/Units/Errors"

	"github.com/gdamore/tcell"
)

type MultiLineText struct {
	lines    [][]*sx.String
	style    tcell.Style
	splitStr string // empty string means no split character
}

func (mlt *MultiLineText) Draw(table *cdd.DrawTable, x, y int) error {
	if table == nil {
		return ers.NewErrNilParameter("table")
	}

	if mlt.IsEmpty() {
		return nil // Nothing to draw
	}

	cb := NewContentBox(mlt.lines, mlt.style, mlt.splitStr)

	return cb.Draw(table, x, y)
}

func (mlt *MultiLineText) ForceDraw(table *cdd.DrawTable, x, y int) error {
	if table == nil {
		return ers.NewErrNilParameter("table")
	}

	if mlt.IsEmpty() {
		return nil // Nothing to draw
	}

	cb := NewContentBox(mlt.lines, mlt.style, mlt.splitStr)

	return cb.ForceDraw(table, x, y)
}

func NewMultiLineText(style tcell.Style, splitString string) *MultiLineText {
	mlt := &MultiLineText{
		lines: make([][]*sx.String, 0),
		style: style,
	}

	return mlt
}

func (mlt *MultiLineText) AppendSentence(sentence string) error {
	lines, err := sext.AdvancedFieldsSplitter(sentence, IndentLevel)
	if err != nil {
		return err
	}

	for _, line := range lines {
		newWords := make([]*sx.String, 0, len(line))

		for _, words := range line {
			newWords = append(newWords, sx.NewString(words, mlt.style))
		}

		mlt.lines = append(mlt.lines, newWords)
	}

	return nil
}

func (mlt *MultiLineText) IsEmpty() bool {
	return len(mlt.lines) == 0
}

func (mlt *MultiLineText) GetLines() []*sx.String {
	if len(mlt.lines) == 0 {
		return nil
	}

	lines := make([]*sx.String, 0, len(mlt.lines))

	for _, line := range mlt.lines {
		lines = append(lines, sx.Join(line, mlt.splitStr))
	}

	return lines
}
