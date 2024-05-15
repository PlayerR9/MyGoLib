package Display

import (
	"strings"

	"github.com/gdamore/tcell"

	ers "github.com/PlayerR9/MyGoLib/Units/Errors"
)

// DrawTable represents a table of cells that can be drawn to the screen.
type DrawTable struct {
	// table is a 2D slice of DtCell pointers. Each element in the slice
	// represents a row in the table, and each element in a row represents
	// a cell in the table.
	table [][]*DtCell

	// width is the width of the table.
	width int

	// height is the height of the table.
	height int
}

// NewDrawTable creates a new DrawTable with the given width and height.
//
// Parameters:
//   - width: The width of the drawTable.
//   - height: The height of the drawTable.
//
// Returns:
//   - *DrawTable: The new drawTable.
//   - error: An error of type *ers.ErrInvalidParameter if the width or height
//     are less than 0.
func NewDrawTable(width, height int) (*DrawTable, error) {
	if width < 0 {
		return nil, ers.NewErrInvalidParameter(
			"width",
			ers.NewErrGTE(0),
		)
	} else if height < 0 {
		return nil, ers.NewErrInvalidParameter(
			"height",
			ers.NewErrGTE(0),
		)
	}

	table := make([][]*DtCell, height)
	for i := 0; i < height; i++ {
		table[i] = make([]*DtCell, width)
	}

	return &DrawTable{
		table:  table,
		width:  width,
		height: height,
	}, nil
}

// GetWidth returns the width of the drawTable.
//
// Returns:
//   - int: The width of the drawTable.
func (dt *DrawTable) GetWidth() int {
	return dt.width
}

// GetHeight returns the height of the drawTable.
//
// Returns:
//   - int: The height of the drawTable.
func (dt *DrawTable) GetHeight() int {
	return dt.height
}

// WriteAt writes a cell to the drawTable at the given coordinates.
//
// Parameters:
//   - x: The x-coordinate of the cell.
//   - y: The y-coordinate of the cell.
//   - cell: The cell to write to the drawTable.
//
// Behaviors:
//   - If the x or y coordinates are out of bounds, the function does nothing.
func (dt *DrawTable) WriteAt(x, y int, cell *DtCell) {
	if x < 0 || x >= dt.width || y < 0 || y >= dt.height {
		return
	}

	dt.table[y][x] = cell
}

// GetAt returns the cell at the given coordinates in the drawTable.
//
// Parameters:
//   - x: The x-coordinate of the cell.
//   - y: The y-coordinate of the cell.
//
// Returns:
//   - *DtCell: The cell at the given coordinates.
//
// Behaviors:
//   - If the x or y coordinates are out of bounds, the function returns nil.
func (dt *DrawTable) GetAt(x, y int) *DtCell {
	if x < 0 || x >= dt.width || y < 0 || y >= dt.height {
		return nil
	} else {
		return dt.table[y][x]
	}
}

// ClearTable clears the drawTable by setting all cells to nil.
func (dt *DrawTable) ClearTable() {
	for i := 0; i < dt.height; i++ {
		dt.table[i] = make([]*DtCell, dt.width)
	}
}

// WriteVerticalSequence writes a sequence of cells to the drawTable
// starting at the given coordinates. The sequence is written vertically
// from the starting coordinates.
//
// Parameters:
//   - x: The x-coordinate of the starting cell.
//   - y: The y-coordinate of the starting cell.
//   - sequence: The sequence of cells to write to the drawTable.
//
// Behaviors:
//   - Any value that would cause the sequence to be written outside the
//     bounds of the drawTable is ignored.
func (dt *DrawTable) WriteVerticalSequence(x, y int, sequence []*DtCell) {
	if len(sequence) == 0 || x < 0 || x >= dt.width || y >= dt.height {
		return
	}

	if y < 0 {
		sequence = sequence[-y:]
		y = 0
	} else if y+len(sequence) > dt.height {
		sequence = sequence[:dt.height-y]
	}

	for i, cell := range sequence {
		dt.table[y+i][x] = cell
	}
}

// WriteHorizontalSequence writes a sequence of cells to the drawTable
// starting at the given coordinates. The sequence is written horizontally
// from the starting coordinates.
//
// Parameters:
//   - x: The x-coordinate of the starting cell.
//   - y: The y-coordinate of the starting cell.
//   - sequence: The sequence of cells to write to the drawTable.
//
// Behaviors:
//   - Any value that would cause the sequence to be written outside the
//     bounds of the drawTable is ignored.
func (dt *DrawTable) WriteHorizontalSequence(x, y int, sequence []*DtCell) {
	if len(sequence) == 0 || y < 0 || y >= dt.height || x >= dt.width {
		return
	}

	if x < 0 {
		sequence = sequence[-x:]
		x = 0
	} else if x+len(sequence) > dt.width {
		sequence = sequence[:dt.width-x]
	}

	copy(dt.table[y][x:], sequence)
}

// GetLines returns each line of the drawTable as a string.
//
// Returns:
//   - []string: The lines of the drawTable.
//
// Behaviors:
//   - Any nil cells in the drawTable are represented by a space character.
func (dt *DrawTable) GetLines() []string {
	lines := make([]string, 0, dt.height)
	var builder strings.Builder

	for i := 0; i < dt.height; i++ {
		for _, cell := range dt.table[i] {
			if cell != nil {
				builder.WriteRune(cell.Content)
			} else {
				builder.WriteRune(' ')
			}
		}

		lines = append(lines, builder.String())
		builder.Reset()
	}

	return lines
}

// WriteLineAt writes a string to the drawTable at the given coordinates.
//
// Parameters:
//   - x: The x-coordinate of the starting cell.
//   - y: The y-coordinate of the starting cell.
//   - line: The string to write to the drawTable.
//   - style: The style of the string.
//   - isHorizontal: A boolean that determines if the string should be written
//     horizontally or vertically.
//
// Behaviors:
//   - This is just a convenience function that converts the string to a sequence
//     of cells and calls WriteHorizontalSequence or WriteVerticalSequence.
func (dt *DrawTable) WriteLineAt(x, y int, line string, style tcell.Style, isHorizontal bool) {
	runes := []rune(line)

	sequence := make([]*DtCell, 0, len(runes))

	for _, r := range runes {
		sequence = append(sequence, NewDtCell(r, style))
	}

	if isHorizontal {
		dt.WriteHorizontalSequence(x, y, sequence)
	} else {
		dt.WriteVerticalSequence(x, y, sequence)
	}
}
