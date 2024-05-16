package Display

import (
	ers "github.com/PlayerR9/MyGoLib/Units/Errors"
	"github.com/gdamore/tcell"
)

// DtUniter is an interface that represents the content of a unit in a draw table.
type DtUniter interface {
	// GetRunes returns the content of the unit as a 2D slice of runes.
	//
	// Returns:
	//   - [][]rune: The content of the unit as a 2D slice of runes.
	GetRunes() [][]rune
}

// DtUnit represents a unit in a draw table. It contains the content of the unit
type DtUnit[T DtUniter] struct {
	// Content is the content of the unit.
	Content T

	// Style is the style of the unit.
	Style tcell.Style
}

// Draw is a method of cdd.TableDrawer that draws the unit to the table at the given x and y
// coordinates.
//
// Parameters:
//   - table: The table to draw the unit to.
//   - x: The x coordinate to draw the unit at.
//   - y: The y coordinate to draw the unit at.
//
// Returns:
//   - error: An error if the table cannot be drawn to.
//
// Errors:
//   - *ers.ErrInvalidParameter: If the table is nil, or if the x or y coordinates are out of
//     bounds.
//   - *ErrHeightExceeded: If the unit would be drawn below the table.
//   - *ErrWidthExceeded: If the unit would be drawn to the right of the table.
func (u *DtUnit[T]) Draw(table *DrawTable, x, y int) error {
	if table == nil {
		return ers.NewErrNilParameter("table")
	}

	if x < 0 || x >= table.width {
		return ers.NewErrInvalidParameter(
			"x",
			ers.NewErrOutOfBounds(x, 0, table.width),
		)
	} else if y < 0 || y >= table.height {
		return ers.NewErrInvalidParameter(
			"y",
			ers.NewErrOutOfBounds(y, 0, table.height),
		)
	}

	runes := u.Content.GetRunes()

	for i, row := range runes {
		if y+i >= table.height {
			return NewErrHeightExceeded(y + i - table.height)
		}

		if len(row) == 0 {
			continue
		}

		for j, r := range row {
			if x+j >= table.width {
				return NewErrWidthExceeded(x + j - table.width)
			}

			if r != '\000' {
				table.table[y+i][x+j] = NewDtCell(r, u.Style)
			}
		}
	}

	return nil
}

// Draw is a method of cdd.TableDrawer that draws the unit to the table at the given x and y
// coordinates.
//
// Parameters:
//   - table: The table to draw the unit to.
//   - x: The x coordinate to draw the unit at.
//   - y: The y coordinate to draw the unit at.
//
// Returns:
//   - error: An error of type *ers.ErrInvalidParameter if the table is nil.
//
// Behaviors:
//   - Any value that would be drawn outside of the table is not drawn.
func (u *DtUnit[T]) ForceDraw(table *DrawTable, x, y int) error {
	if table == nil {
		return ers.NewErrNilParameter("table")
	}

	runes := u.Content.GetRunes()

	runes, x, y = FixBoundaries(table.width, table.height, runes, x, y)
	if len(runes) == 0 {
		return nil
	}

	for i, row := range runes {
		if len(row) == 0 {
			continue
		}

		for j, r := range row {
			if x+j >= table.width {
				break
			}

			if r != '\000' {
				table.table[y+i][x+j] = NewDtCell(r, u.Style)
			}
		}
	}

	return nil
}

// NewDtUnit creates a new DtUnit with the given content and style.
//
// Parameters:
//   - content: The content of the unit.
//   - style: The style of the unit.
//
// Returns:
//   - *DtUnit: The new DtUnit.
func NewDtUnit[T DtUniter](content T, style tcell.Style) *DtUnit[T] {
	return &DtUnit[T]{
		Content: content,
		Style:   style,
	}
}
