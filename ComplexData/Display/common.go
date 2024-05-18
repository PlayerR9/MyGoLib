package Display

// TableDrawer is an interface that represents an object that can draw a table.
type TableDrawer interface {
	// Draw draws the table to the screen at the given x and y coordinates.
	//
	// Parameters:
	//   - dt: The DrawTable to draw.
	//   - x: The x-coordinate to draw the table at.
	//   - y: The y-coordinate to draw the table at.
	//
	// Returns:
	//   - error: An error if the table could not be drawn.
	Draw(dt *DrawTable, x, y int) error

	// ForceDraw draws the table to the screen at the given x and y coordinates.
	// This function should draw the table even if values are out of bounds.
	//
	// Parameters:
	//   - dt: The DrawTable to draw.
	//   - x: The x-coordinate to draw the table at.
	//   - y: The y-coordinate to draw the table at.
	//
	// Returns:
	//   - error: An error if the table could not be drawn.
	//
	// Behaviors:
	//   - Values that cannot be drawn should be ignored.
	//   - The error should only be used for critical errors that prevent the
	//     table from being drawn. (i.e., error is returned when a panic would
	//     occur if the table was drawn.)
	ForceDraw(dt *DrawTable, x, y int) error
}
