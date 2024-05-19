package Display

import (
	ers "github.com/PlayerR9/MyGoLib/Units/Errors"
	"github.com/gdamore/tcell"

	cdt "github.com/PlayerR9/MyGoLib/Units/Table"
)

// DtUniter is an interface that represents the content of a unit in a draw table.
type DtUniter interface {
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
	Runes(width, height int) ([][]rune, error)
}

// DtUnit represents a unit in a draw table. It contains the content of the unit
type DtUnit[T DtUniter] struct {
	// content is the content of the unit.
	content T

	// style is the style of the unit.
	style tcell.Style
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
func (u *DtUnit[T]) Draw(table *DrawTable, x, y int) error {
	if table == nil {
		return ers.NewErrNilParameter("table")
	}

	width, height := table.GetWidth(), table.GetHeight()

	runeTable, err := u.content.Runes(width, height)
	if err != nil {
		return err
	}

	// Fix the boundaries of the rune table
	runeTable, x, y = cdt.FixBoundaries(width, height, runeTable, x, y)
	if len(runeTable) == 0 {
		return nil
	}

	for i, row := range runeTable {
		if len(row) == 0 {
			continue
		}

		sequence := make([]*dtCell, 0, len(row))

		for _, r := range row {
			if r == EmptyRuneCell {
				sequence = append(sequence, nil)
			} else {
				sequence = append(sequence, newDtCell(r, u.style))
			}
		}

		table.WriteHorizontalSequence(x, y+i, sequence)
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
		content: content,
		style:   style,
	}
}
