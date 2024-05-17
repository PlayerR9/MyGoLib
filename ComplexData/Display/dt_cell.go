package Display

import (
	"github.com/gdamore/tcell"

	ers "github.com/PlayerR9/MyGoLib/Units/Errors"
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

// Draw is a method of cdd.TableDrawer that draws the cell to the table at the given x and y
// coordinates.
//
// Parameters:
//   - table: The table to draw the cell to.
//   - x: The x coordinate to draw the cell at.
//   - y: The y coordinate to draw the cell at.
//
// Returns:
//   - error: An error of type *ers.ErrInvalidParameter if the table is nil, or if the x or y
//     coordinates are out of bounds.
func (c *DtCell) Draw(table *DrawTable, x, y int) error {
	if table == nil {
		return ers.NewErrNilParameter("table")
	} else {
		return table.WriteAt(x, y, c)
	}
}

// ForceDraw is a method of cdd.TableDrawer that draws the cell to the table at the given x
// and y coordinates.
//
// Parameters:
//   - table: The table to draw the cell to.
//   - x: The x coordinate to draw the cell at.
//   - y: The y coordinate to draw the cell at.
//
// Returns:
//   - error: An error of type *ers.ErrInvalidParameter if the table is nil.
//
// Behaviors:
//   - If the x or y coordinates are out of bounds, the cell will still not be drawn.
func (c *DtCell) ForceDraw(table *DrawTable, x, y int) error {
	if table == nil {
		return ers.NewErrNilParameter("table")
	} else {
		table.ForceWriteAt(x, y, c)
		return nil
	}
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
