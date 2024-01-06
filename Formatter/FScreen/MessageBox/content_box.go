// git tag v0.1.46

package MessageBox

import (
	"fmt"
	"math"
	"slices"
	"strings"

	sext "github.com/PlayerR9/MyGoLib/Utility/StrExt"
	"github.com/gdamore/tcell"
)

const (
	// FieldSpacing defines the number of spaces between each field (word)
	// when they are written into the ContentBox.
	// It is set to 1, meaning there will be one spaces between each field.
	FieldSpacing int = 1

	// Hellip defines the string to be used as an ellipsis when the content
	// of the ContentBox is truncated.
	// It is set to "...", which is the standard representation of an ellipsis
	// in text.
	Hellip string = "..."

	// HellipLen defines the length of the Hellip string.
	// It is set to 3, which is the number of characters in the Hellip string.
	HellipLen int = len(Hellip)

	// IndentLevel defines the number of spaces used for indentation when
	// writing content into the ContentBox.
	// It is set to 2, meaning there will be two spaces at the start of each
	// new line of content.
	IndentLevel int = 2
)

// ContentBox represents a box in the terminal where content can be written.
// It maintains a table of runes representing the content, and a corresponding
// array of styles for each rune.
// It also keeps track of the first empty line in the box, the width and
// height of the box, and a flag indicating whether the box has started shifting
// content.
type ContentBox struct {
	// A 2D slice of runes representing the content of the box.
	table [][]rune

	// A slice of tcell.Style objects corresponding to the styles of the runes in the table.
	styles []tcell.Style

	// A pointer to an integer representing the first empty line in the box.
	// A value of -1 indicates the box is closed.
	firstEmptyLine int

	// The dimensions of the box.
	width, height int

	// A slice of runes representing an empty line.
	// This is used to simplify the process of clearing the table.
	emptyLine []rune
}

// NewContentBox creates a new ContentBox with the specified width and height.
// It returns an error if the width is less than 5 or the height is less than 2.
//
// Parameters:
//
//   - width: The width of the ContentBox.
//   - height: The height of the ContentBox.
//
// Returns:
//
//   - The newly created ContentBox.
//   - An error if the width is less than 5 or the height is less than 2, otherwise nil.
func NewContentBox(width, height int) (ContentBox, error) {
	var cb ContentBox

	if width < 5 {
		return cb, &ErrWidthTooSmall{}
	} else if height < 2 {
		return cb, &ErrHeightTooSmall{}
	}

	// Generate an empty line
	cb.emptyLine = []rune(strings.Repeat(string(Space), width))

	// Fill the table with spaces
	cb.table = make([][]rune, height)
	cb.styles = make([]tcell.Style, height)

	for i := 0; i < height; i++ {
		cb.table[i] = make([]rune, width)
		copy(cb.table[i], cb.emptyLine)
		cb.styles[i] = StyleMap[NormalText]
	}

	cb.firstEmptyLine = 0
	cb.width = width
	cb.height = height

	return cb, nil
}

// Clear resets the ContentBox to its initial state.
// It sets the first empty line to 0, trims the table and styles to the height of the box,
// and sets the startedShifting flag to false.
//
// This function is used to clear the content of the box and prepare it for new content.
func (cb *ContentBox) Clear() {
	cb.firstEmptyLine = 0

	cb.table = make([][]rune, cb.height)
	cb.styles = make([]tcell.Style, cb.height)

	for i := 0; i < cb.height; i++ {
		cb.table[i] = make([]rune, cb.width)
		copy(cb.table[i], cb.emptyLine)

		cb.styles[i] = StyleMap[NormalText]
	}
}

// ResizeWidth changes the width of the ContentBox to the specified width.
// If the new width is greater than the current width, the function grows the
// box by filling the new space with spaces.
// If the new width is less than the current width, the function shrinks the
// box by discarding the extra characters and replacing them with hellipses.
//
// Parameters:
//
//   - width: The new width of the ContentBox.
func (cb *ContentBox) ResizeWidth(width int) {
	defer func() { cb.width = width }()
	if width == cb.width {
		return // Nothing to do
	}

	newTable := make([][]rune, cb.height)
	for i := 0; i < cb.height; i++ {
		newTable[i] = make([]rune, width)
		copy(newTable[i], cb.table[i])
	}

	if width > cb.width {
		// Grow: Copy the old table and the old styles
		// but fill the new space with spaces
		// NOTE: These loops are inverted
		for j := cb.width; j < width; j++ {
			for i := 0; i < cb.height; i++ {
				newTable[i][j] = Space
			}
		}
	} else {
		// Shrink: Copy the old table and the old styles
		// but discard the extra characters (replacing them with hellipses)
		for i := 0; i < cb.height; i++ {
			if newTable[i][width-1] != Space {
				copy(newTable[i][width-HellipLen:], []rune(Hellip))
			}
		}
	}

	cb.table = newTable
	cb.emptyLine = []rune(strings.Repeat(string(Space), width))
}

// ResizeHeight changes the height of the ContentBox to the specified height.
// If the new height is greater than the current height, the function grows
// the box by increasing the table size.
// If the new height is less than or equal to the current height, the function
// shrinks the box by making the extra lines out of bounds.
//
// Parameters:
//
//   - height: The new height of the ContentBox.
func (cb *ContentBox) ResizeHeight(height int) {
	defer func() { cb.height = height }()

	upperLimit := int(math.Max(float64(len(cb.table)), float64(cb.height)))

	if height <= upperLimit {
		return
	}

	for i := upperLimit; i < height; i++ {
		cb.table = append(cb.table, cb.emptyLine)
		cb.styles = append(cb.styles, StyleMap[NormalText])
	}
}

// HasEmptyLine checks if the ContentBox has an empty line where new content
// can be written.
// The function returns true if the first empty line is less than or equal
// to the height of the box, meaning there is an empty line.
// Otherwise, it returns false, meaning the box is full.
//
// Returns:
//
//   - A boolean indicating whether the ContentBox has an empty line.
func (cb *ContentBox) HasEmptyLine() bool {
	return cb.firstEmptyLine <= cb.height
}

// WriteStringAt writes a string at the specified coordinates in the
// ContentBox.
// It starts writing the string at the given x and y coordinates and
// continues horizontally until the string is fully written, until it
// reaches the width of the ContentBox, or until an error occurs.
//
// If the x coordinate plus the length of the ellipsis is greater than or
// equal to the width of the ContentBox, the function will return an
// ErrXOutOfBounds error.
//
// If the y coordinate is greater than or equal to the height of the
// ContentBox (i.e., the length of the table), the function will return
// an ErrYOutOfBounds error.
//
// If the length of the text plus the x coordinate is greater than the
// width of the ContentBox, the function will attempt to replace the suffix
// of the text with an ellipsis (Hellip). If this replacement is successful
// and the text can now fit within the width of the ContentBox, the
// function will write the text and return. If the text still cannot fit
// within the width of the ContentBox, the function will return an
// ErrTextTooLong error.
//
// Parameters:
//
//   - x: The x coordinate where the writing should start.
//   - y: The y coordinate where the writing should start.
//   - text: The string to be written.
//
// Returns:
//
//   - The x coordinate of the first empty cell after the last character
//     written.
//   - An error if the x coordinate is out of bounds, the y coordinate is
//     out of bounds, or the text is too long.
func (cb *ContentBox) WriteStringAt(x, y int, text string) (int, error) {
	if x+HellipLen >= cb.width {
		return x, &ErrXOutOfBounds{}
	} else if y >= len(cb.table) {
		return x, &ErrYOutOfBounds{}
	}

	if x+len(text) >= cb.width {
		var err error

		text, err := sext.ReplaceSuffix(text[:cb.width-x-HellipLen-1], Hellip)
		if err != nil {
			return x, fmt.Errorf("unable to replace suffix of text '%s' with ellipsis: %w", text, err)
		}
	}

	copy(cb.table[y][x:x+len(text)], []rune(text))

	return x + len(text), nil
}

// GetLastWriteableFieldIndex determines the last field from the provided
// slice of strings that can be written in full on the current line of the
// ContentBox, starting at the specified x coordinate. The function also
// checks if the next field in the slice can be written with an ellipsis
// at the end, if it cannot be written in full.
//
// If no fields can be written, the function returns the length of the
// fields slice minus 1.
//
// Parameters:
//
//   - x: The x coordinate where the writing should start.
//   - fields: The slice of strings to be written.
//
// Returns:
//
//   - The index of the last field that can be written in full or with an
//     ellipsis at the end.
func (cb *ContentBox) GetLastWriteableFieldIndex(x int, fields []string) int {
	if index := slices.IndexFunc(fields, func(s string) bool {
		length := len(s) + x
		if length >= cb.width || length+HellipLen >= cb.width {
			return true
		}

		x = length + FieldSpacing
		return false
	}); index < 0 {
		return len(fields) - 1
	} else {
		return index - 1
	}
}

// writeFieldsWithSpacing writes a slice of strings (fields) into the
// ContentBox starting at the specified x and y coordinates. Each string
// is written with a spacing defined by FieldSpacing.
//
// If an error occurs while writing a string, the function stops and
// returns the current x coordinate and the error.
//
// Parameters:
//
//   - x: The x coordinate where the writing should start.
//   - y: The y coordinate where the writing should start.
//   - fields: The slice of strings to be written.
//
// Returns:
//
//   - The x coordinate after the last string has been written.
//   - An error if writing any of the strings fails, otherwise nil.
func (cb *ContentBox) writeFieldsWithSpacing(x, y int, fields []string) (int, error) {
	var err error

	x, err = cb.WriteStringAt(x, y, fields[0])
	if err != nil {
		return x, fmt.Errorf("unable to write string '%s' at x = %d, y = %d: %w", fields[0], x, y, err)
	}

	for _, field := range fields[1:] {
		x, err = cb.WriteStringAt(x+FieldSpacing, y, field)
		if err != nil {
			return x, fmt.Errorf("unable to write string '%s' at x = %d, y = %d: %w", field, x, y, err)
		}
	}

	return x, nil
}

// handleFirstLine writes the first line of fields into the ContentBox
// starting at the specified x and y coordinates. The fields are written
// with spacing defined by FieldSpacing.
//
// Parameters:
//
//   - x: The x coordinate where the writing should start.
//   - y: The y coordinate where the writing should start.
//   - fields: The slice of strings to be written.
//
// Returns:
//
//   - The x coordinate after the last string has been written.
//   - The y coordinate after the last string has been written.
//   - An error if writing any of the strings fails, otherwise nil.
func (cb *ContentBox) handleFirstLine(x, y int, fields []string) (int, int, error) {
	const DebugMode bool = false

	var err error

	lastIndex := cb.GetLastWriteableFieldIndex(x, fields)
	if lastIndex == -1 {
		return x, y, fmt.Errorf("no fields can be written. This should never happen")
	}
	lastIndex++

	// DEBUG: Print the last index
	if DebugMode {
		fmt.Printf("Last index (first line): %d\n", lastIndex)
	}

	x, err = cb.writeFieldsWithSpacing(x, y, fields[:lastIndex])
	if err != nil {
		return x, y, err
	} else if lastIndex >= len(fields) {
		return x, y, nil
	}

	_, err = cb.WriteStringAt(x, y, Hellip)
	if err != nil {
		return x, y, fmt.Errorf("unable to write string '%s' at x = %d, y = %d: %w", Hellip, x, y, err)
	} else {
		return IndentLevel, y + 1, nil
	}
}

// tryHandlingConsecutiveLines attempts to write as many fields as possible
// from the fieldMatrix into the ContentBox starting at the specified x and
// y coordinates. The fields are written with spacing defined by FieldSpacing.
//
// Parameters:
//
//   - x: The x coordinate where the writing should start.
//   - y: The y coordinate where the writing should start.
//   - index: The index of the fieldMatrix where the writing should start.
//   - fieldMatrix: The 2D slice of strings to be written.
//
// Returns:
//
//   - The x coordinate after the last string has been written.
//   - The y coordinate after the last string has been written.
//   - The index of the fieldMatrix after the last string has been written.
//   - The index of the field in the fieldMatrix after the last string has
//     been written.
//   - An error if writing any of the strings fails, otherwise nil.
func (cb *ContentBox) tryHandlingConsecutiveLines(x, y, index int, fieldMatrix [][]string) (int, int, int, int, error) {
	const DebugMode bool = false

	var err error

	if y > cb.firstEmptyLine {
		// WARNING: This might be wrong
		// if errors occur here, check this

		// Increase the table size (if necessary)
		if y >= len(cb.table) {
			for i := len(cb.table); i < y; i++ {
				cb.table = append(cb.table, cb.emptyLine)
				cb.styles = append(cb.styles, StyleMap[NormalText])
			}
		}

		return x, y, index, 0, nil
	}

	// Try to append to the current line as many fields as possible
	// until the line is full (ignore ellipsing)
	for ; index < len(fieldMatrix); index++ {
		x += FieldSpacing

		lastIndex := cb.GetLastWriteableFieldIndex(x, fieldMatrix[index])
		if lastIndex == -1 {
			return x, y, index, 0, fmt.Errorf("no fields can be written. This should never happen")
		}
		lastIndex++

		// DEBUG: Print the last index
		if DebugMode {
			fmt.Printf("Last index (consecutive lines): %d\n", lastIndex)
		}

		x, err = cb.writeFieldsWithSpacing(x, y, fieldMatrix[index][:lastIndex])
		if err != nil {
			return x, y, index, 0, err
		}

		if lastIndex >= len(fieldMatrix[index]) {
			return IndentLevel, y + 1, index, lastIndex + 1, nil
		}
	}

	return x, y, index, 0, nil
}

// EnqueueContents writes a slice of strings (contents) into the ContentBox
// with a specified style. Each string in the contents is split into fields
// (words) and written into the ContentBox.
// The fields are written with spacing defined by FieldSpacing.
//
// Parameters:
//
//   - contents: The slice of strings to be written.
//   - style: The style to be applied to the strings.
//
// Returns:
//
//   - An error if writing any of the strings fails, otherwise nil.
func (cb *ContentBox) EnqueueContents(contents []string, style tcell.Style) error {
	const DebugMode bool = false

	// Create the 2D field matrix
	fieldMatrix := make([][]string, 0)

	for _, content := range contents {
		fields := strings.Fields(content)

		if len(fields) != 0 {
			fieldMatrix = append(fieldMatrix, fields)
		}
	}
	if len(fieldMatrix) == 0 {
		return nil // Skip empty lines
	}

	// DEBUG: Print the field matrix
	if DebugMode {
		fmt.Println("Field matrix:")

		for _, fields := range fieldMatrix {
			fmt.Println(fields)
		}
	}

	originalY := cb.firstEmptyLine

	// Handle first line
	x, y, err := cb.handleFirstLine(0, originalY, fieldMatrix[0])
	if err != nil {
		return err
	}

	// Handle the rest of the lines
	startNextLine, index := 0, 1

	for index < len(fieldMatrix) {
		_, y, index, startNextLine, err = cb.tryHandlingConsecutiveLines(x, y, index, fieldMatrix)
		if err != nil {
			return err
		} else if index >= len(fieldMatrix) {
			break
		} else if startNextLine >= len(fieldMatrix[index]) {
			index++
			continue
		}

		// Write to the next line
		x, y, err = cb.handleFirstLine(IndentLevel, y+1, fieldMatrix[index][startNextLine:])
		if err != nil {
			return err
		}
	}

	// Set the style of the lines that were just written
	for i := originalY; i < cb.firstEmptyLine; i++ {
		cb.styles[i] = style
	}

	if y == cb.firstEmptyLine {
		cb.firstEmptyLine++
	} else {
		cb.firstEmptyLine = y
	}

	return nil
}

// EnqueueLineSeparator writes a line of the specified rune (char) into the
// ContentBox as a separator.
// If the first empty line is at the height of the box, the separator is
// written outside the box.
// Otherwise, the separator is written inside the box at the first empty line.
//
// Parameters:
//
//   - char: The rune to be used as the line separator.
func (cb *ContentBox) EnqueueLineSeparator(char rune) {
	emptyLine := cb.firstEmptyLine // Reduce the risk of concurrent writes
	defer func() { cb.firstEmptyLine = emptyLine }()

	if emptyLine == cb.height {
		// Write the separator outside the message box

		cb.table = append(cb.table, make([]rune, cb.width))
		for i := 0; i < cb.width; i++ {
			cb.table[emptyLine][i] = char
		}

		cb.styles = append(cb.styles, StyleMap[NormalText])
	} else {
		// Write the separator inside the message box
		cb.table[emptyLine] = []rune(strings.Repeat(string(char), cb.width))

		emptyLine++
	}
}

// CanShiftUp checks if the ContentBox can shift its content up.
// The function returns true if the ContentBox has started shifting and the first
// empty line is not at the top, or if the first empty line is at or beyond the
// halfway point of the box's height.
// If the ContentBox has started shifting and the first empty line is at the top,
// it stops shifting and returns false.
//
// Returns:
//
//   - A boolean indicating whether the ContentBox can shift its content up.
func (cb *ContentBox) CanShiftUp() bool {
	return cb.firstEmptyLine >= cb.height
}

// ShiftUp shifts the content of the ContentBox up by a certain number of lines.
// The number of lines to shift up is determined by the number of lines occupied by the
// first message in the box.
// If the first empty line is at the top of the box, the function does nothing.
//
// This function is used to make room for new messages when the box is full.
func (cb *ContentBox) ShiftUp() {
	emptyLine := cb.firstEmptyLine
	defer func() { cb.firstEmptyLine = emptyLine }()

	if emptyLine == 0 {
		return // Nothing to do
	}

	// 1. Find the amount of lines occupied by the first message
	size := 1

	for size < emptyLine && cb.table[size][IndentLevel] == Space {
		size++
	}

	// 2. Shift inbounds values up by size lines
	upperLimit := int(math.Min(float64(cb.height), float64(emptyLine)) - float64(size))

	for i := 0; i < upperLimit; i++ {
		copy(cb.table[i], cb.table[i+size])
		cb.styles[i] = cb.styles[i+size]
	}

	for i := upperLimit; i < emptyLine; i++ {
		copy(cb.table[i], cb.emptyLine)
		cb.styles[i] = StyleMap[NormalText]
	}

	emptyLine -= size

	if len(cb.table) > cb.height {
		// 3. Shift out of bounds values up by size lines and reduce the table size (if necessary)
		for i := upperLimit; i < emptyLine; i++ {
			copy(cb.table[i], cb.table[i+size])
			cb.styles[i] = cb.styles[i+size]
		}

		// 4. Reduce the table size to the new first upper limit
		upperLimit = int(math.Max(float64(emptyLine), float64(cb.height)))
		cb.table = cb.table[:upperLimit]
		cb.styles = cb.styles[:upperLimit]
	}
}

// drawHorizontalBorderAt draws a horizontal border at the specified
// y-coordinate on the screen with the specified style.
// The border is drawn from the left corner to the right corner of the
// ContentBox, with '+' characters at the corners and '-' characters in between.
//
// Parameters:
//
//   - y: The y-coordinate where the border is to be drawn.
//   - style: The style to be used for the border.
//   - screen: The screen where the border is to be drawn.
//
// Returns:
//
//   - The screen after the border has been drawn.
func (cb *ContentBox) drawHorizontalBorderAt(y int, style tcell.Style, screen tcell.Screen) tcell.Screen {
	screen.SetContent(0, y, '+', nil, style) // Left corner

	for x := 1; x < cb.width+1+PaddingWidth; x++ {
		screen.SetContent(x, y, '-', nil, style)
	}

	screen.SetContent(cb.width+1+PaddingWidth, y, '+', nil, style) // Right corner

	return screen
}

// Draw renders the ContentBox on the screen at the specified y-coordinate.
// The function draws the top and bottom borders, and the content of the box
// between the borders.
// The content is drawn line by line, with '|' characters at the left and
// right borders, and the content of the box in between.
// The function returns the y-coordinate of the last line drawn and the
// screen after the box has been drawn.
//
// Parameters:
//
//   - y: The y-coordinate where the box is to be drawn.
//   - screen: The screen where the box is to be drawn.
//
// Returns:
//
//   - The y-coordinate of the last line drawn.
//   - The screen after the box has been drawn.
func (cb *ContentBox) Draw(y int, screen tcell.Screen) (int, tcell.Screen) {
	style := StyleMap[NormalText]

	screen = cb.drawHorizontalBorderAt(y, style, screen) // Top border

	for i := 0; i < cb.firstEmptyLine; i++ {
		y++
		screen.SetContent(0, y, '|', nil, style) // Left border

		for j, cell := range cb.table[i] {
			screen.SetContent(Padding+j, y, cell, nil, cb.styles[i])
		}

		screen.SetContent(cb.width+PaddingWidth, y, '|', nil, style) // Right border
	}

	for i := cb.firstEmptyLine; i < cb.height; i++ {
		y++
		screen.SetContent(0, y, '|', nil, style) // Left border

		for j, cell := range cb.table[i] {
			screen.SetContent(Padding+j, y, cell, nil, cb.styles[i])
		}

		screen.SetContent(cb.width+PaddingWidth, y, '|', nil, style) // Right border
	}

	y++
	screen = cb.drawHorizontalBorderAt(y, style, screen) // Bottom border

	return y, screen
}
