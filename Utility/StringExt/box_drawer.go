package StringExt

import (
	uc "github.com/PlayerR9/MyGoLib/Units/common"
)

var (
	// DefaultBoxStyle is the default box style.
	DefaultBoxStyle *BoxStyle
)

func init() {
	DefaultBoxStyle = &BoxStyle{
		LineType: BtNormal,
		IsHeavy:  false,
		Padding:  [4]int{1, 1, 1, 1},
	}
}

// BoxBorderType is the type of the box border.
type BoxBorderType int

const (
	// BtNormal is the normal box border type.
	BtNormal BoxBorderType = iota

	// BtTriple is the triple box border type.
	BtTriple

	// BtQuadruple is the quadruple box border type.
	BtQuadruple

	// BtDouble is the double box border type.
	BtDouble

	// BtRounded is like BtNormal but with rounded corners.
	BtRounded
)

// BoxStyle is the style of the box.
type BoxStyle struct {
	// LineType is the type of the line.
	LineType BoxBorderType

	// IsHeavy is whether the line is heavy or not.
	// Only applicable to BtNormal, BtTriple, and BtQuadruple.
	IsHeavy bool

	// Padding is the padding of the box.
	// [Top, Right, Bottom, Left]
	Padding [4]int
}

// NewBoxStyle creates a new box style.
//
// Parameters:
//   - lineType: The line type.
//   - isHeavy: Whether the line is heavy or not.
//   - padding: The padding of the box. [Top, Right, Bottom, Left]
//
// Returns:
//   - *BoxStyle: The new box style.
//
// Behaviors:
//   - If the padding is negative, it will be set to 0.
func NewBoxStyle(lineType BoxBorderType, isHeavy bool, padding [4]int) *BoxStyle {
	for i := 0; i < 4; i++ {
		if padding[i] < 0 {
			padding[i] = 0
		}
	}

	bs := &BoxStyle{
		LineType: lineType,
		IsHeavy:  isHeavy,
		Padding:  padding,
	}

	return bs
}

// GetCorners gets the corners of the box.
//
// Returns:
//   - [4]rune: The corners. [TopLeft, TopRight, BottomLeft, BottomRight]
func (bs *BoxStyle) GetCorners() [4]rune {
	var corners [4]rune

	if bs.IsHeavy {
		corners = [4]rune{'┏', '┓', '┗', '┛'}
	} else {
		corners = [4]rune{'┌', '┐', '└', '┘'}
	}

	return corners
}

// GetTopBorder gets the top border of the box.
//
// It also applies to the bottom border as they are the same.
//
// Returns:
//   - string: The top border.
func (bs *BoxStyle) GetTopBorder() rune {
	var tb_border rune

	switch bs.LineType {
	case BtNormal:
		if bs.IsHeavy {
			tb_border = '━'
		} else {
			tb_border = '─'
		}
	case BtTriple:
		if bs.IsHeavy {
			tb_border = '┅'
		} else {
			tb_border = '┄'
		}
	case BtQuadruple:
		if bs.IsHeavy {
			tb_border = '┉'
		} else {
			tb_border = '┅'
		}
	case BtDouble:
		tb_border = '═'
	case BtRounded:
		tb_border = '─'
	}

	return tb_border
}

// GetSideBorder gets the side border of the box.
//
// It also applies to the left border as they are the same.
//
// Returns:
//   - string: The side border.
func (bs *BoxStyle) GetSideBorder() rune {
	var side_border rune

	switch bs.LineType {
	case BtNormal:
		if bs.IsHeavy {
			side_border = '┃'
		} else {
			side_border = '│'
		}
	case BtTriple:
		if bs.IsHeavy {
			side_border = '┇'
		} else {
			side_border = '┆'
		}
	case BtQuadruple:
		if bs.IsHeavy {
			side_border = '┋'
		} else {
			side_border = '┆'
		}
	case BtDouble:
		side_border = '║'
	case BtRounded:
		side_border = '│'
	}

	return side_border
}

// makeEmptyRow is a helper function to make an empty row.
//
// Parameters:
//   - width: The width of the row.
//   - side_border: The side border of the row.
//
// Returns:
//   - []rune: The empty row.
//
// Assertions:
//   - width >= 0
func makeEmptyRow(width int, side_border rune) []rune {
	uc.AssertParam("width", width >= 0, uc.NewErrGTE(0))

	empty_row := make([]rune, 0, width+2)
	for i := 1; i < width; i++ {
		empty_row = append(empty_row, ' ')
	}

	empty_row = append([]rune{side_border}, empty_row...)
	empty_row = append(empty_row, side_border)

	return empty_row
}

// makeSidePadding is a helper function to make side padding.
//
// Parameters:
//   - width: The width of the padding.
//
// Returns:
//   - []rune: The side padding.
func makeSidePadding(width int) []rune {
	side_padding := make([]rune, 0, width)
	for i := 0; i < width; i++ {
		side_padding = append(side_padding, ' ')
	}

	return side_padding
}

// makeTBBorder is a helper function to make a top or bottom border.
//
// Parameters:
//   - width: The width of the border.
//   - border: The border character.
//   - left_corner: The left corner character.
//   - right_corner: The right corner character.
//
// Returns:
//   - []rune: The top or bottom border.
//
// Assertions:
//   - width >= 0
func makeTBBorder(width int, border, left_corner, right_corner rune) []rune {
	uc.AssertParam("width", width >= 0, uc.NewErrGTE(0))

	row := make([]rune, 0, width+2)

	row = append(row, left_corner)
	for i := 0; i < width; i++ {
		row = append(row, border)
	}

	row = append(row, right_corner)

	return row
}

// DrawBox draws a box around the content.
//
// Format: If the content is ["Hello", "World"], the box will be:
//
//	┏━━━━━━━┓
//	┃ Hello ┃
//	┃ World ┃
//	┗━━━━━━━┛
//
// Parameters:
//   - content: The content.
//
// Returns:
//   - string: The content in a box.
//
// Behaviors:
//   - If the box style is nil, the default box style will be used.
func (bs *BoxStyle) ApplyStrings(content []string) (*RuneTable, error) {
	for i := 0; i < 4; i++ {
		if bs.Padding[i] < 0 {
			bs.Padding[i] = 0
		}
	}

	side_border := bs.GetSideBorder()
	left_padding := makeSidePadding(bs.Padding[3])
	right_padding := makeSidePadding(bs.Padding[1])
	tbb_char := bs.GetTopBorder()
	corners := bs.GetCorners()
	prefix := append([]rune{side_border}, left_padding...)
	suffix := append(right_padding, side_border)

	table, err := NewRuneTable(content)
	if err != nil {
		return nil, err
	}

	right_edge := table.AlignRightEdge()

	total_width := right_edge + bs.Padding[1] + bs.Padding[3]
	empty_row := makeEmptyRow(total_width, side_border)

	top_border := makeTBBorder(total_width, tbb_char, corners[0], corners[1])
	bottom_border := makeTBBorder(total_width, tbb_char, corners[2], corners[3])

	for i := 0; i < bs.Padding[0]; i++ {
		table.PrependTopRow(empty_row)
	}

	for i := 0; i < bs.Padding[2]; i++ {
		table.AppendBottomRow(empty_row)
	}
	table.PrefixEachRow(prefix)
	table.SuffixEachRow(suffix)
	table.PrependTopRow(top_border)
	table.AppendBottomRow(bottom_border)

	return table, nil
}
