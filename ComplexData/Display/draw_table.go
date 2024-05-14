package Display

import (
	"fmt"

	"github.com/gdamore/tcell"
)

// dtCell represents a cell in a data table. It contains a rune which
// represents the content of the cell, and a tcell.Style which represents
// the style of the cell.
type DtCell struct {
	// Content is the rune that represents the content of the cell.
	Content rune

	// Style is the tcell.Style that represents the style of the cell.
	Style tcell.Style
}

func (cell *DtCell) String() string {
	return string(cell.Content)
}

func NewDtCell(content rune, style tcell.Style) *DtCell {
	return &DtCell{
		Content: content,
		Style:   style,
	}
}

type DrawTable struct {
	tableHeight int // Height of the table
	*p.PtTable[*DtCell]
}

func (dt *DrawTable) Init(dt p.Partitionable[*DtCell]) *DrawTable {
	dt.tableHeight = dt.GetHeight()
	dt.PtTable = new(p.PtTable[*DtCell]).Init(dt)
	dt.untriggeredEvent = nil

	return dt
}

/// FIX THIS

func (dt *DrawTable) HasEvent() bool {
	return dt.untriggeredEvent != nil
}

func (dt *DrawTable) SetEvent(ev *DtEvent) {
	dt.untriggeredEvent = ev
}

func (dt *DrawTable) WriteLineAt(x, y int, contents string, style tcell.Style) error {
	if len(contents) == 0 {
		return nil // Nothing to do
	}

	width, height := dt.GetWidth(), dt.GetHeight()

	if x < 0 || x >= width {
		return ers.NewErrInvalidParameter("x").
			WithReason(ers.NewErrOutOfBound(0, width, x))
	}

	if y < 0 || y >= height {
		return ers.NewErrInvalidParameter("y").
			WithReason(ers.NewErrOutOfBound(0, height, y))
	}

	if x+len(contents) > width {
		return ers.NewErrInvalidParameter("contents").
			WithReason(fmt.Errorf("length (%d) exceeds table's width by %d", len(contents), x+len(contents)-width))
	}

	sequence := make([]*DtCell, 0, len(contents))
	for _, content := range contents {
		sequence = append(sequence, &DtCell{Content: content, Style: style})
	}

	dt.WriteSequence(x, y, sequence, false)

	return nil
}

func (dt *drawTable) TriggerEvent() error {
	if dt.untriggeredEvent == nil {
		return fmt.Errorf("no event to trigger")
	}

	if dt.Parent == nil {
		return fmt.Errorf("cannot trigger event: table doesn't support triggering events")
	}

	dt.FirstParent.enqueueEvent(dt.untriggeredEvent)
	dt.untriggeredEvent = nil

	return nil
}
