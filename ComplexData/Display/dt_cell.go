package Display

import (
	"github.com/gdamore/tcell"
)

// dtCell represents a cell in a data table. It contains a rune which
// represents the content of the cell, and a tcell.Style which represents
// the style of the cell.
type dtCell struct {
	// content is the rune that represents the content of the cell.
	content rune

	// style is the tcell.style that represents the style of the cell.
	style tcell.Style
}

// newDtCell creates a new DtCell with the given content and style.
//
// Parameters:
//   - content: The content of the cell.
//   - style: The style of the cell.
//
// Returns:
//   - *DtCell: The new DtCell.
func newDtCell(content rune, style tcell.Style) *dtCell {
	return &dtCell{
		content: content,
		style:   style,
	}
}
