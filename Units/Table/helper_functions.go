package Table

// fixVerticalBoundaries is a helper function that fixes the vertical boundaries
// of a table.
//
// Parameters:
//   - maxHeight: The maximum height of the table.
//   - elems: The elements of the table.
//   - y: The y-coordinate to fix the boundaries at.
//
// Returns:
//   - [][]T: The elements of the table with the boundaries fixed.
//   - int: The y-coordinate with the boundaries fixed.
func fixVerticalBoundaries[T any](maxHeight int, elems [][]T, y int) ([][]T, int) {
	if y < 0 {
		return elems[-y:], 0
	} else if y >= maxHeight {
		return elems[:maxHeight-y], maxHeight
	}

	totalHeight := len(elems) + y

	if totalHeight <= maxHeight {
		return elems, y
	} else {
		return elems[:maxHeight-y], maxHeight
	}
}

// fixHorizontalBoundaries is a helper function that fixes the horizontal boundaries
// of a table.
//
// Parameters:
//   - maxWidth: The maximum width of the table.
//   - elems: The elements of the table.
//   - x: The x-coordinate to fix the boundaries at.
//
// Returns:
//   - [][]T: The elements of the table with the boundaries fixed.
//   - int: The x-coordinate with the boundaries fixed.
func fixHorizontalBoundaries[T any](maxWidth int, elems [][]T, x int) ([][]T, int) {
	if x < 0 {
		for i, row := range elems {
			if -x >= len(row) {
				elems[i] = nil
			} else {
				elems[i] = row[-x:]
			}
		}

		return elems, 0
	} else if x >= maxWidth {
		for i, row := range elems {
			if x >= len(row) {
				elems[i] = nil
			} else {
				elems[i] = row[:maxWidth-x]
			}
		}

		return elems, maxWidth
	}

	for i, row := range elems {
		totalWidth := len(row) + x

		if totalWidth < maxWidth {
			continue
		}

		elems[i] = row[:maxWidth-x]
	}

	return elems, x
}

// FixBoundaries is a function that fixes the boundaries of a table of elements based
// on the maximum width and height of the table.
//
// Parameters:
//   - maxWidth: The maximum width of the table.
//   - maxHeight: The maximum height of the table.
//   - elems: The elements of the table.
//   - x: The x-coordinate to fix the boundaries at.
//   - y: The y-coordinate to fix the boundaries at.
//
// Returns:
//   - [][]T: The elements of the table with the boundaries fixed.
//   - int: The x-coordinate with the boundaries fixed.
//   - int: The y-coordinate with the boundaries fixed.
//
// Behaviors:
//   - If maxWidth is less than 0, it is set to 0.
//   - If maxHeight is less than 0, it is set to 0.
//   - If elems is empty, nil is returned.
func FixBoundaries[T any](maxWidth, maxHeight int, elems [][]T, x, y int) ([][]T, int, int) {
	if maxWidth < 0 {
		maxWidth = 0
	}

	if maxHeight < 0 {
		maxHeight = 0
	}

	if len(elems) == 0 {
		return nil, 0, 0
	}

	elems, y = fixVerticalBoundaries(maxHeight, elems, y)
	elems, x = fixHorizontalBoundaries(maxWidth, elems, x)

	return elems, x, y
}
