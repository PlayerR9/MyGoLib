package Display

import (
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

// NewDtCell creates a new DtCell with the given content and style.
//
// Parameters:
//   - content: The content of the cell.
//   - style: The style of the cell.
//
// Returns:
//   - *DtCell: The new DtCell.
func NewDtCell(content rune, style tcell.Style) *DtCell {
	return &DtCell{
		Content: content,
		Style:   style,
	}
}
