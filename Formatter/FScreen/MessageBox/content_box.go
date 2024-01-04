package MessageBox

import (
	"fmt"
	"math"
	"slices"
	"strings"
	"sync"

	rws "github.com/PlayerR9/MyGoLib/CustomData/RWSafe"
	sext "github.com/PlayerR9/MyGoLib/Utility/StrExt"
	"github.com/gdamore/tcell"
)

// WriteStringAt writes a string at the specified coordinates in the ContentBox.
// It starts writing the string at the given x and y coordinates and continues
// horizontally until the string is fully written, until it reaches the width of the ContentBox,
// or until an error occurs.
//
// If the x coordinate plus the length of the ellipsis is greater than or equal to the width of the ContentBox,
// the function will return an ErrXOutOfBounds error.
//
// If the y coordinate is greater than or equal to the height of the ContentBox
// (i.e., the length of the table), the function will return an ErrYOutOfBounds error.
//
// If the length of the text plus the x coordinate is greater than the width of the ContentBox,
// the function will attempt to replace the suffix of the text with an ellipsis (Hellip).
// If this replacement is successful and the text can now fit within the width of the ContentBox,
// the function will write the text and return. If the text still cannot fit within the width of the ContentBox,
// the function will return an ErrTextTooLong error.
//
// Parameters:
//
//   - x: The x coordinate where the writing should start.
//   - y: The y coordinate where the writing should start.
//   - text: The string to be written.
//
// Returns:
//
//   - The x coordinate of the first empty cell after the last character written.
//   - An error if the x coordinate is out of bounds, the y coordinate is out of bounds, or the text is too long.
func (cb *ContentBox) WriteStringAt(x, y int, text string) (int, error) {
	if x+HellipLen >= cb.width {
		return x, &ErrXOutOfBounds{}
	} else if y >= len(cb.table) {
		return x, &ErrYOutOfBounds{}
	}

	if x+len(text) >= cb.width {
		var err error

		text, err = sext.ReplaceSuffix(text[:cb.width-x-HellipLen-1], Hellip)
		if err != nil {
			return x, fmt.Errorf("unable to replace suffix of text '%s' with ellipsis: %w", text, err)
		} else if x+len(text) >= cb.width {
			return x, &ErrTextTooLong{}
		}
	}

	row := cb.table[y]

	for _, char := range text {
		row[x] = char
		x++
	}

	return x, nil
}

// GetLastWriteableFieldIndex determines the last field from the provided slice of strings
// that can be written in full on the current line of the ContentBox, starting at the
// specified x coordinate. The function also checks if the next field in the slice can
// be written with an ellipsis at the end, if it cannot be written in full.
//
// The function operates in two steps:
//  1. It iterates over the fields and updates the x coordinate after each field by the length of
//     the field plus the value of FieldSpacing.
//  2. If the x coordinate plus the length of the ellipsis is greater than or equal to the width of the ContentBox,
//     it returns the index of the current field.
//
// If no fields can be written, the function returns the length of the fields slice minus 1.
//
// Parameters:
//
//   - x: The x coordinate where the writing should start.
//   - fields: The slice of strings to be written.
//
// Returns:
//
//   - The index of the last field that can be written in full or with an ellipsis at the end.
func (cb *ContentBox) GetLastWriteableFieldIndex(x int, fields []string) int {
	if index := slices.IndexFunc(fields, func(s string) bool {
		if x+len(s) >= cb.width || x+len(s)+HellipLen >= cb.width {
			return true
		}

		x += len(s) + FieldSpacing
		return false
	}); index < 0 {
		return len(fields) - 1
	} else {
		return index - 1
	}
}

////////// NON-COMMENTED /////

const (
	FieldSpacing = 2
)

type ContentBox struct {
	table  [][]rune
	styles []tcell.Style

	firstEmptyLine *rws.RWSafe[int] // -1 indicates the message box is closed

	width, height int

	startedShifting bool

	notEmpty sync.Cond
}

func NewContentBox(width, height int) (*ContentBox, error) {
	if width < 5 {
		return nil, &ErrWidthTooSmall{}
	} else if height < 2 {
		return nil, &ErrHeightTooSmall{}
	}

	table := make([][]rune, height)
	for i := 0; i < height; i++ {
		table[i] = make([]rune, width)
	}

	return &ContentBox{
		table:           table,
		styles:          make([]tcell.Style, height),
		firstEmptyLine:  rws.NewRWSafe(0),
		width:           width,
		height:          height,
		startedShifting: false,
		notEmpty:        *sync.NewCond(&sync.Mutex{}),
	}, nil
}

func (cb *ContentBox) EnqueueLineSeparator(char rune) {
	emptyLine := cb.firstEmptyLine.Get()

	if emptyLine == cb.height {
		// Write the separator outside the message box

		cb.table = append(cb.table, make([]rune, cb.width))
		for i := 0; i < cb.width; i++ {
			cb.table[emptyLine][i] = char
		}

		cb.styles = append(cb.styles, StyleMap[NormalText])
	} else {
		// Write the separator inside the message box
		for i := 0; i < cb.width; i++ {
			cb.table[emptyLine][i] = char
		}

		cb.firstEmptyLine.Set(emptyLine + 1)
	}
}

func (cb *ContentBox) EnqueueContents(contents []string, style tcell.Style) {
	remaining := cb.enqueueContents(contents, style, false)

	for len(remaining) > 0 {
		remaining = cb.enqueueContents(remaining, style, true)
	}
}

func (cb *ContentBox) CanShiftUp() bool {
	if cb.startedShifting && cb.firstEmptyLine.Get() == 0 {
		cb.startedShifting = false
		return false
	}

	return (cb.startedShifting && cb.firstEmptyLine.Get() != 0) || cb.firstEmptyLine.Get() >= cb.height/2
}

func (cb *ContentBox) ShiftUp() {
	emptyLine := cb.firstEmptyLine.Get()
	cb.startedShifting = emptyLine != 0

	if emptyLine == 0 {
		return // Nothing to do
	}

	// 1. Find the amount of lines occupied by the first message
	size := 1

	for size < emptyLine && cb.table[size][0] == Space {
		size++
	}

	// 2. Shift inbounds values up by size lines
	upperLimit := int(math.Min(float64(cb.height), float64(emptyLine)) - float64(size))

	for i := 0; i < upperLimit; i++ {
		cb.table[i] = cb.table[i+size]
		cb.styles[i] = cb.styles[i+size]
	}

	emptyLine -= size

	if len(cb.table) > cb.height {
		// 3. Shift out of bounds values up by size lines and reduce the table size (if necessary)
		for i := upperLimit; i < emptyLine; i++ {
			cb.table[i] = cb.table[i+size]
			cb.styles[i] = cb.styles[i+size]
		}

		// 4. Reduce the table size to the new first upper limit
		upperLimit = int(math.Max(float64(emptyLine), float64(cb.height)))
		cb.table = cb.table[:upperLimit]
		cb.styles = cb.styles[:upperLimit]
	}
}

func (cb *ContentBox) ResizeWidth(width int) error {
	if width == cb.width {
		return nil // Nothing to do
	}

	newTable := make([][]rune, cb.height)
	for i := 0; i < cb.height; i++ {
		newTable[i] = make([]rune, width)
	}

	newStyles := make([]tcell.Style, cb.height)

	if width > cb.width {
		// Grow: Copy the old table and the old styles
		// but fill the new space with spaces

		for i, row := range cb.table {
			copy(newTable[i], row)

			for j := cb.width; j < width; j++ {
				newTable[i][j] = Space
			}
		}

		copy(newStyles, cb.styles)
	} else {
		// Shrink: Copy the old table and the old styles
		// but discard the extra characters (replacing them with hellipses)

		// NOTE: This way of doing it is really inefficient but
		// it is the simplest way to do it as of now
		// improve it later

		for i, row := range cb.table {
			line := string(row[:width])

			// enqueueContents will take care of the hellipses
			cb.enqueueContents([]string{line}, cb.styles[i], row[0] == Space)

			newStyles[i] = cb.styles[i]
		}
	}

	cb.table = newTable
	cb.styles = newStyles

	return nil
}

func (cb *ContentBox) ResizeHeight(height int) error {
	if height <= cb.height {
		// Nothing to do

		// Shrink: Just make the extra lines
		// out of bounds
		return nil
	}

	// Grow: Fill the new lines with spaces

	for i := cb.height; i < height; i++ {
		cb.table[i] = make([]rune, cb.width)
		cb.styles = append(cb.styles, StyleMap[NormalText])
	}

	cb.height = height

	return nil
}

func (cb *ContentBox) Clear() {
	cb.firstEmptyLine.Set(0)

	cb.table = cb.table[:cb.height]
	cb.styles = cb.styles[:cb.height]
}

func (cb *ContentBox) HasEmptyLine() bool {
	return cb.firstEmptyLine.Get() <= cb.height
}

func (cb *ContentBox) enqueueContentsCanTryLine(x, y int, fields []string) (int, int, []string) {
	lastIndex := cb.GetLastWriteableFieldIndex(IndentLevel, fields)
	if lastIndex == -1 {
		panic("No fields can be written. This should never happen")
	}

	if lastIndex == len(fields)-1 {
		x, _ := cb.WriteStringAt(x, y, fields[0])

		for _, field := range fields[1:] {
			x, _ = cb.WriteStringAt(x+FieldSpacing, y, field)
		}

		return x, y, nil
	} else {
		x, _ := cb.WriteStringAt(x, y, fields[0])

		for _, field := range fields[1:lastIndex] {
			x, _ = cb.WriteStringAt(x+FieldSpacing, y, field)
		}

		return IndentLevel, y + 1, fields[lastIndex+1:]
	}
}

func (cb *ContentBox) enqueueContentsTryLine(x, y int, fields []string) (int, int) {
	lastIndex := cb.GetLastWriteableFieldIndex(IndentLevel, fields)
	if lastIndex == -1 {
		panic("No fields can be written. This should never happen")
	}

	if lastIndex == len(fields)-1 {
		x, _ := cb.WriteStringAt(x, y, fields[0])

		for _, field := range fields[1:] {
			x, _ = cb.WriteStringAt(x+FieldSpacing, y, field)
		}

		return x, y
	} else {
		x, _ := cb.WriteStringAt(x, y, fields[0])

		for _, field := range fields[1:lastIndex] {
			x, _ = cb.WriteStringAt(x+FieldSpacing, y, field)
		}

		cb.WriteStringAt(x, y, Hellip)

		return IndentLevel, y + 1
	}
}

// WARNING: This function doesn't do any parameter checks
func (cb *ContentBox) enqueueContents(contents []string, style tcell.Style, isIndented bool) []string {
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

	x, y := cb.enqueueContentsTryLine(0, cb.firstEmptyLine.Get(), fieldMatrix[0])
	index := 1

	for index < len(fieldMatrix) {
		var fields []string

		if y == cb.firstEmptyLine.Get() {
			var remainder []string

			for index++; index < len(fieldMatrix); index++ {
				x, y, remainder = cb.enqueueContentsCanTryLine(x+FieldSpacing, y, fieldMatrix[index])
				if len(remainder) != 0 {
					break
				}
			}

			fields = remainder
		} else {
			cb.firstEmptyLine.Set(y)

			// Increase the table size if necessary
			if cb.firstEmptyLine.Get() >= len(cb.table) {
				cb.table = append(cb.table, make([]rune, cb.width))
				cb.styles = append(cb.styles, StyleMap[NormalText])
			}

			// Go to the next non-empty line
			index++

			if index >= len(fieldMatrix) {
				break
			}

			fields = fieldMatrix[index]
		}

		x, y = cb.enqueueContentsTryLine(IndentLevel, y+1, fields)
	}

	return nil
}

func (cb *ContentBox) drawHorizontalBorderAt(y int, style tcell.Style, screen tcell.Screen) tcell.Screen {
	screen.SetContent(0, y, '+', nil, style) // Left corner

	for x := 1; x < cb.width+1+PaddingWidth; x++ {
		screen.SetContent(x, y, '-', nil, style)
	}

	screen.SetContent(cb.width+1+PaddingWidth, y, '+', nil, style) // Right corner

	return screen
}

func (cb *ContentBox) Draw(y int, screen tcell.Screen) (int, tcell.Screen) {
	style := StyleMap[NormalText]

	screen = cb.drawHorizontalBorderAt(y, style, screen) // Top border

	for i := 0; i < cb.firstEmptyLine.Get(); i++ {
		y++
		screen.SetContent(0, y, '|', nil, style) // Left border

		for j, cell := range cb.table[i] {
			screen.SetContent(Padding+j, y, cell, nil, cb.styles[i])
		}

		screen.SetContent(cb.width+PaddingWidth, y, '|', nil, style) // Right border
	}

	for i := cb.firstEmptyLine.Get(); i < cb.height; i++ {
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

////////////////// OLD CODE //////////////////////

// REMEMBER TO INITIALIZE THE MESSAGEBOX WITH THE PADDING

func (cb *ContentBox) Close() {

}

// Clear interface{} information to prevent a deadlock and release the memory
// FIXME: Check if this works
func (cb *ContentBox) Fini() {
	// Wake up the message box if it is waiting for a message
	// this will prevent a deadlock
	cb.firstEmptyLine.Set(-1)
	cb.notEmpty.Broadcast()

	// BAD: This shouldn't be done in the first place
}

func (cb *ContentBox) Cleanup() {
	cb.table = nil
	cb.styles = nil
	cb.notEmpty.L = nil
	cb.firstEmptyLine = nil
}

////////////////////////////////

///////////////////////////////

const (
	Hellip    = "..."
	HellipLen = 3
)

var StyleMap map[TextMessageType]tcell.Style = map[TextMessageType]tcell.Style{
	NormalText:  tcell.StyleDefault.Foreground(tcell.ColorBlack).Background(tcell.ColorGhostWhite),
	DebugText:   tcell.StyleDefault.Bold(true).Foreground(tcell.ColorSlateGray).Background(tcell.ColorGhostWhite),
	WarningText: tcell.StyleDefault.Foreground(tcell.ColorDarkOrange).Background(tcell.ColorGhostWhite),
	ErrorText:   tcell.StyleDefault.Foreground(tcell.ColorFireBrick).Background(tcell.ColorGhostWhite),
	FatalText:   tcell.StyleDefault.Bold(true).Foreground(tcell.ColorDarkRed).Background(tcell.ColorGhostWhite),
	SuccessText: tcell.StyleDefault.Foreground(tcell.ColorDarkGreen).Background(tcell.ColorGhostWhite),
}

const (
	IndentLevel int = 2
)

func (cb *ContentBox) GetHeight() int {
	return cb.height
}

// Returns the first empty index after the last character
// WARNING: This function doesn't do interface{} checks
func (cb *ContentBox) WriteFieldsAt(x, y int, fields []string) (int, int) {
	if len(fields) == 1 {
		if x+len(fields[0]) < cb.width {
			x, _ := cb.WriteStringAt(x, y, fields[0])
			return x, y
		}

		x, _ = cb.WriteStringAt(x, y, fields[0][:cb.width-x-HellipLen])
		x, _ = cb.WriteStringAt(x, y, Hellip)
		return x, y
	} else if x+len(fields[0]) >= cb.width {
		x, _ = cb.WriteStringAt(x, y, fields[0][:cb.width-x-HellipLen])
		cb.WriteStringAt(x, y, Hellip)

		return IndentLevel, y + 1
	}

	x, _ = cb.WriteStringAt(x, y, fields[0])

	index := -1
	for i, field := range fields[1:] {
		if len(field)+x+2 >= cb.width {
			index = i
			break
		}

		x, _ = cb.WriteStringAt(x+2, y, field)
	}

	if index != -1 {
		cb.WriteStringAt(x, y, Hellip)

		return IndentLevel, y + 1
	}

	return x, y
}
