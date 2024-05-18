package Display

import (
	"strings"

	"github.com/gdamore/tcell"

	cdt "github.com/PlayerR9/MyGoLib/CustomData/Table"
)

const (
	// EmptyRuneCell is a rune that represents an empty cell.
	EmptyRuneCell rune = '\000'
)

// DisplayTable represents a table of cells that can be drawn to the screen.
type DisplayTable struct {
	// table is the table of cells in the displayTable. Each cell is a rune.
	*cdt.Table[rune]
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
	table, err := cdt.NewTable[rune](height, width)
	if err != nil {
		return nil, err
	}

	return &DisplayTable{
		table,
	}, nil
}

// GetLines returns each line of the displayTable as a string.
//
// Returns:
//   - []string: The lines of the displayTable.
//
// Behaviors:
//   - Any nil cells in the displayTable are represented by a space character.
func (dt *DisplayTable) GetLines() []string {
	table := dt.GetFullTable()

	lines := make([]string, 0, len(table))
	var builder strings.Builder

	for _, row := range table {
		for _, cell := range row {
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
