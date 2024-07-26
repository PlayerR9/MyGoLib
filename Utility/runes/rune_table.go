package runes

import (
	"strings"

	uc "github.com/PlayerR9/MyGoLib/Units/common"
)

// RuneTable is a table of runes.
type RuneTable struct {
	// table is the table of runes.
	table [][]rune
}

// String implements the fmt.Stringer interface.
func (rt *RuneTable) String() string {
	var builder strings.Builder

	for _, row := range rt.table {
		builder.WriteString(string(row))
		builder.WriteRune('\n')
	}

	return builder.String()
}

// NewRuneTable creates a new RuneTable with the given lines.
//
// Parameters:
//   - lines: The lines to add to the table.
//
// Returns:
//   - *RuneTable: The new RuneTable.
//   - error: An error if any.
func NewRuneTable(lines []string) (*RuneTable, error) {
	table := make([][]rune, 0, len(lines))

	for i, line := range lines {
		row, err := StringToUtf8(line)
		if err != nil {
			return nil, uc.NewErrAt(i+1, "line", err)
		}

		table = append(table, row)
	}

	rt := &RuneTable{
		table: table,
	}

	return rt, nil
}

// GetRightMostEdge gets the right most edge of the content.
//
// Parameters:
//   - content: The content.
//
// Returns:
//   - int: The right most edge.
func (rt *RuneTable) GetRightMostEdge() int {
	var longest_line int

	for _, row := range rt.table {
		if len(row) > longest_line {
			longest_line = len(row)
		}
	}

	return longest_line
}

// AlignRightEdge aligns the right edge of the table.
//
// Returns:
//   - int: The right most edge.
func (rt *RuneTable) AlignRightEdge() int {
	edge := rt.GetRightMostEdge()

	for i := 0; i < len(rt.table); i++ {
		curr_row := rt.table[i]

		padding := edge - len(curr_row)

		padding_right := make([]rune, 0, padding)
		for i := 0; i < padding; i++ {
			padding_right = append(padding_right, ' ')
		}

		rt.table[i] = append(curr_row, padding_right...)
	}

	return edge
}

// PrependTopRow prepends a row to the top of the table.
//
// Parameters:
//   - row: The row to prepend.
func (rt *RuneTable) PrependTopRow(row []rune) {
	rt.table = append([][]rune{row}, rt.table...)
}

// AppendBottomRow appends a row to the bottom of the table.
//
// Parameters:
//   - row: The row to append.
func (rt *RuneTable) AppendBottomRow(row []rune) {
	rt.table = append(rt.table, row)
}

// PrefixEachRow prefixes each row with the given prefix.
//
// Parameters:
//   - prefix: The prefix to add to each row.
func (rt *RuneTable) PrefixEachRow(prefix []rune) {
	for i := 0; i < len(rt.table); i++ {
		new_row := append(prefix, rt.table[i]...)
		rt.table[i] = new_row
	}
}

// SuffixEachRow suffixes each row with the given suffix.
//
// Parameters:
//   - suffix: The suffix to add to each row.
func (rt *RuneTable) SuffixEachRow(suffix []rune) {
	for i := 0; i < len(rt.table); i++ {
		new_row := append(rt.table[i], suffix...)
		rt.table[i] = new_row
	}
}
