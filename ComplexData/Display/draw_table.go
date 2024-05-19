package Display

import (
	"strings"

	"github.com/gdamore/tcell"

	cdt "github.com/PlayerR9/MyGoLib/Units/Table"
)

// DrawTable represents a table of cells that can be drawn to the screen.
type DrawTable struct {
	*cdt.Table[*dtCell]
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
	table, err := cdt.NewTable[*dtCell](width, height)
	if err != nil {
		return nil, err
	}

	return &DrawTable{table}, nil
}

// GetLines returns each line of the drawTable as a string.
//
// Returns:
//   - []string: The lines of the drawTable.
//
// Behaviors:
//   - Any nil cells in the drawTable are represented by a space character.
func (dt *DrawTable) GetLines() []string {
	table := dt.GetFullTable()

	lines := make([]string, 0, len(table))
	var builder strings.Builder

	for _, row := range table {
		for _, cell := range row {
			if cell != nil {
				builder.WriteRune(cell.content)
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

	sequence := make([]*dtCell, 0, len(runes))

	for _, r := range runes {
		sequence = append(sequence, newDtCell(r, style))
	}

	if isHorizontal {
		dt.WriteHorizontalSequence(x, y, sequence)
	} else {
		dt.WriteVerticalSequence(x, y, sequence)
	}
}
