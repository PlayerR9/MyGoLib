package Display

// TableDrawer is an interface that represents an object that can draw a table.
type TableDrawer interface {
	// Draw draws the table to the screen.
	//
	// Parameters:
	//   - dt: The DrawTable to draw.
	//
	// Returns:
	//   - error: An error if the table could not be drawn.
	Draw(dt *DrawTable) error
}
