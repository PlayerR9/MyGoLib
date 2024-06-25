package Table

import (
	uc "github.com/PlayerR9/MyGoLib/Units/common"
)

// Table[T] represents a table of cells that can be drawn to the screen.
type Table[T any] struct {
	// table is a 2D slice of elements of pointer type T. Each element in the slice
	// represents a row in the table, and each element in a row represents
	// a cell in the table.
	table [][]T

	// width is the width of the table.
	width int

	// height is the height of the table.
	height int
}

// Iterator returns an iterator for the table. It iterates over the rows and columns
// as if it were an array of elements of type T.
//
// Returns:
//   - uc.Iterater[T]: An iterator for the table.
//
// Behaviors:
//   - For efficiency, the iterator is procedural. It is not a generator.
func (t *Table[T]) Iterator() uc.Iterater[T] {
	iter := uc.NewDynamicIterator(
		uc.NewSimpleIterator(t.table),
		func(row []T) uc.Iterater[T] {
			return uc.NewSimpleIterator(row)
		},
	)

	return iter
}

// NewDrawTable creates a new Table[T] with the given width and height.
//
// Parameters:
//   - width: The width of the table.
//   - height: The height of the table.
//
// Returns:
//   - *Table[T]: The new table.
//
// Behaviors:
//   - If the width or height is negative, the absolute value is used.
func NewTable[T any](width, height int) *Table[T] {
	if width < 0 {
		width = -width
	}

	if height < 0 {
		height = -height
	}

	table := make([][]T, height)
	for i := 0; i < height; i++ {
		table[i] = make([]T, width)
	}

	return &Table[T]{
		table:  table,
		width:  width,
		height: height,
	}
}

// GetWidth returns the width of the table.
//
// Returns:
//   - int: The width of the table.
func (dt *Table[T]) GetWidth() int {
	return dt.width
}

// GetHeight returns the height of the table.
//
// Returns:
//   - int: The height of the table.
func (dt *Table[T]) GetHeight() int {
	return dt.height
}

// WriteAt writes a cell to the table at the given coordinates.
//
// Parameters:
//   - x: The x-coordinate of the cell.
//   - y: The y-coordinate of the cell.
//   - cell: The cell to write to the table.
//
// Behaviors:
//   - If the x or y coordinates are out of bounds, the function does nothing.
func (dt *Table[T]) WriteAt(x, y int, cell T) {
	if x < 0 || x >= dt.width || y < 0 || y >= dt.height {
		return
	}

	dt.table[y][x] = cell
}

// GetAt returns the cell at the given coordinates in the table.
//
// Parameters:
//   - x: The x-coordinate of the cell.
//   - y: The y-coordinate of the cell.
//
// Returns:
//   - T: The cell at the given coordinates.
//
// Behaviors:
//   - If the x or y coordinates are out of bounds, the function returns
//     the zero value of type T.
func (dt *Table[T]) GetAt(x, y int) T {
	if x < 0 || x >= dt.width || y < 0 || y >= dt.height {
		return *new(T)
	} else {
		return dt.table[y][x]
	}
}

// ClearTable clears the table by setting all cells to their zero value.
func (dt *Table[T]) ClearTable() {
	for i := 0; i < dt.height; i++ {
		dt.table[i] = make([]T, dt.width)
	}
}

// WriteVerticalSequence writes a sequence of cells to the table
// starting at the given coordinates. The sequence is written vertically
// from the starting coordinates.
//
// Parameters:
//   - x: The x-coordinate of the starting cell.
//   - y: The y-coordinate of the starting cell.
//   - sequence: The sequence of cells to write to the table.
//
// Behaviors:
//   - Any value that would cause the sequence to be written outside the
//     bounds of the table is ignored.
//   - The y-coordinate is updated to the next available cell after the
//     sequence is written.
func (dt *Table[T]) WriteVerticalSequence(x, y *int, sequence []T) {
	actualX, actualY := *x, *y

	if len(sequence) == 0 || actualX < 0 || actualX >= dt.width || actualY >= dt.height {
		return
	}

	if actualY < 0 {
		sequence = sequence[-actualY:]

		*y = 0
	} else if actualY+len(sequence) > dt.height {
		sequence = sequence[:dt.height-actualY]
	}

	for i, cell := range sequence {
		dt.table[actualY+i][actualX] = cell
	}

	*y += len(sequence)
}

// WriteHorizontalSequence writes a sequence of cells to the table
// starting at the given coordinates. The sequence is written horizontally
// from the starting coordinates.
//
// Parameters:
//   - x: The x-coordinate of the starting cell.
//   - y: The y-coordinate of the starting cell.
//   - sequence: The sequence of cells to write to the table.
//
// Behaviors:
//   - Any value that would cause the sequence to be written outside the
//     bounds of the table is ignored.
//   - The x-coordinate is updated to the next available cell after the
//     sequence is written.
func (dt *Table[T]) WriteHorizontalSequence(x, y *int, sequence []T) {
	actualX, actualY := *x, *y

	if len(sequence) == 0 || actualY < 0 || actualY >= dt.height || actualX >= dt.width {
		return
	}

	if actualX < 0 {
		sequence = sequence[-actualX:]

		*x = 0
	} else if actualX+len(sequence) > dt.width {
		sequence = sequence[:dt.width-actualX]
	}

	copy(dt.table[actualY][actualX:], sequence)

	*x = actualX + len(sequence)
}

// GetFullTable returns the full table as a 2D slice of elements of type T.
//
// Returns:
//   - [][]T: The full table.
func (dt *Table[T]) GetFullTable() [][]T {
	return dt.table
}

// IsXInBounds checks if the given x-coordinate is within the bounds of the table.
//
// Parameters:
//   - x: The x-coordinate to check.
//
// Returns:
//   - error: An error of type *uc.ErrOutOfBounds if the x-coordinate is out of bounds.
func (dt *Table[T]) IsXInBounds(x int) error {
	if x < 0 || x >= dt.width {
		return uc.NewErrOutOfBounds(x, 0, dt.width)
	} else {
		return nil
	}
}

// IsYInBounds checks if the given y-coordinate is within the bounds of the table.
//
// Parameters:
//   - y: The y-coordinate to check.
//
// Returns:
//   - error: An error of type *uc.ErrOutOfBounds if the y-coordinate is out of bounds.
func (dt *Table[T]) IsYInBounds(y int) error {
	if y < 0 || y >= dt.height {
		return uc.NewErrOutOfBounds(y, 0, dt.height)
	} else {
		return nil
	}
}

// WriteTableAt writes the given table to the table at the given coordinates.
//
// Parameters:
//   - x: The x-coordinate to write the table at.
//   - y: The y-coordinate to write the table at.
//   - table: The table to write to the table.
//
// Behaviors:
//   - If the table is nil, the function does nothing.
//   - Any value that would cause the table to be written outside the bounds
//     of the table is ignored.
//   - The x and y coordinates are updated to the next available cell after
//     the table is written.
func (dt *Table[T]) WriteTableAt(table *Table[T], x, y *int) {
	if table == nil {
		return
	}

	offsetX, offsetY := 0, 0
	X, Y := *x, *y

	for offsetY < table.height && Y+offsetY < dt.height {
		offsetX = 0

		for offsetX < table.width && X+offsetX < dt.width {
			dt.table[Y+offsetY][X+offsetX] = table.table[offsetY][offsetX]
			offsetX++
		}

		offsetY++
	}

	*x += offsetX
	*y += offsetY
}
