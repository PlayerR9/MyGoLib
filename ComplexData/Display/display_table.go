package Display

import (
	"strings"

	"github.com/gdamore/tcell"

	ers "github.com/PlayerR9/MyGoLib/Units/Errors"
)

const (
	// EmptyRuneCell is a rune that represents an empty cell.
	EmptyRuneCell rune = '\000'
)

// DisplayTable represents a table of cells that can be drawn to the screen.
type DisplayTable struct {
	// table is a 2D slice of runes. Each element in the slice
	// represents a row in the table, and each element in a row represents
	// a cell in the table.
	table [][]rune

	// width is the width of the table.
	width int

	// height is the height of the table.
	height int
}

// NewDrawTable creates a new DisplayTable with the given width and height.
//
// Parameters:
//   - width: The width of the displayTable.
//   - height: The height of the displayTable.
//
// Returns:
//   - *DisplayTable: The new displayTable.
//   - error: An error of type *ers.ErrInvalidParameter if the width or height
//     are less than 0.
func NewDisplayTable(width, height int) (*DisplayTable, error) {
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

	table := make([][]rune, height)
	for i := 0; i < height; i++ {
		table[i] = make([]rune, width)
	}

	return &DisplayTable{
		table:  table,
		width:  width,
		height: height,
	}, nil
}

// GetWidth returns the width of the displayTable.
//
// Returns:
//   - int: The width of the displayTable.
func (dt *DisplayTable) GetWidth() int {
	return dt.width
}

// GetHeight returns the height of the displayTable.
//
// Returns:
//   - int: The height of the displayTable.
func (dt *DisplayTable) GetHeight() int {
	return dt.height
}

// WriteAt writes a cell to the displayTable at the given coordinates.
//
// Parameters:
//   - x: The x-coordinate of the cell.
//   - y: The y-coordinate of the cell.
//   - cell: The cell to write to the displayTable.
//
// Behaviors:
//   - If the x or y coordinates are out of bounds, the function does nothing.
func (dt *DisplayTable) WriteAt(x, y int, cell rune) {
	if x < 0 || x >= dt.width || y < 0 || y >= dt.height {
		return
	}

	dt.table[y][x] = cell
}

// GetAt returns the cell at the given coordinates in the displayTable.
//
// Parameters:
//   - x: The x-coordinate of the cell.
//   - y: The y-coordinate of the cell.
//
// Returns:
//   - rune: The cell at the given coordinates.
//
// Behaviors:
//   - If the x or y coordinates are out of bounds, the function returns EmptyRuneCell.
func (dt *DisplayTable) GetAt(x, y int) rune {
	if x < 0 || x >= dt.width || y < 0 || y >= dt.height {
		return EmptyRuneCell
	} else {
		return dt.table[y][x]
	}
}

// ClearTable clears the displayTable by setting all cells to nil.
func (dt *DisplayTable) ClearTable() {
	for i := 0; i < dt.height; i++ {
		dt.table[i] = make([]rune, dt.width)
	}
}

// WriteVerticalSequence writes a sequence of cells to the displayTable
// starting at the given coordinates. The sequence is written vertically
// from the starting coordinates.
//
// Parameters:
//   - x: The x-coordinate of the starting cell.
//   - y: The y-coordinate of the starting cell.
//   - sequence: The sequence of cells to write to the displayTable.
//
// Behaviors:
//   - Any value that would cause the sequence to be written outside the
//     bounds of the displayTable is ignored.
func (dt *DisplayTable) WriteVerticalSequence(x, y int, sequence []rune) {
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

// WriteHorizontalSequence writes a sequence of cells to the displayTable
// starting at the given coordinates. The sequence is written horizontally
// from the starting coordinates.
//
// Parameters:
//   - x: The x-coordinate of the starting cell.
//   - y: The y-coordinate of the starting cell.
//   - sequence: The sequence of cells to write to the displayTable.
//
// Behaviors:
//   - Any value that would cause the sequence to be written outside the
//     bounds of the displayTable is ignored.
func (dt *DisplayTable) WriteHorizontalSequence(x, y int, sequence []rune) {
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

// GetLines returns each line of the displayTable as a string.
//
// Returns:
//   - []string: The lines of the displayTable.
//
// Behaviors:
//   - Any nil cells in the displayTable are represented by a space character.
func (dt *DisplayTable) GetLines() []string {
	lines := make([]string, 0, dt.height)
	var builder strings.Builder

	for i := 0; i < dt.height; i++ {
		for _, cell := range dt.table[i] {
			if cell != EmptyRuneCell {
				builder.WriteRune(cell)
			} else {
				builder.WriteRune(' ')
			}
		}

		lines = append(lines, builder.String())
		builder.Reset()
	}

	return lines
}

// WriteLineAt writes a string to the displayTable at the given coordinates.
//
// Parameters:
//   - x: The x-coordinate of the starting cell.
//   - y: The y-coordinate of the starting cell.
//   - line: The string to write to the displayTable.
//   - style: The style of the string.
//   - isHorizontal: A boolean that determines if the string should be written
//     horizontally or vertically.
//
// Behaviors:
//   - This is just a convenience function that converts the string to a sequence
//     of cells and calls WriteHorizontalSequence or WriteVerticalSequence.
func (dt *DisplayTable) WriteLineAt(x, y int, line string, style tcell.Style, isHorizontal bool) {
	sequence := []rune(line)

	if isHorizontal {
		dt.WriteHorizontalSequence(x, y, sequence)
	} else {
		dt.WriteVerticalSequence(x, y, sequence)
	}
}

// IsXInBounds checks if the given x-coordinate is within the bounds of the displayTable.
//
// Parameters:
//   - x: The x-coordinate to check.
//
// Returns:
//   - error: An error of type *ers.ErrOutOfBounds if the x-coordinate is out of bounds.
func (dt *DisplayTable) IsXInBounds(x int) error {
	if x < 0 || x >= dt.width {
		return ers.NewErrOutOfBounds(x, 0, dt.width)
	} else {
		return nil
	}
}

// IsYInBounds checks if the given y-coordinate is within the bounds of the displayTable.
//
// Parameters:
//   - y: The y-coordinate to check.
//
// Returns:
//   - error: An error of type *ers.ErrOutOfBounds if the y-coordinate is out of bounds.
func (dt *DisplayTable) IsYInBounds(y int) error {
	if y < 0 || y >= dt.height {
		return ers.NewErrOutOfBounds(y, 0, dt.height)
	} else {
		return nil
	}
}
