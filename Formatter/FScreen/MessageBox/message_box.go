package MessageBox

import (
	"strings"
	"sync"

	buffer "github.com/PlayerR9/MyGoLib/CustomData/Buffer"
	rws "github.com/PlayerR9/MyGoLib/CustomData/RWSafe"
	"github.com/gdamore/tcell"
)

func (mb *MessageBox) SetSize(width, height int) error {
	return nil
}

func (mb *MessageBox) Draw(y int, screen tcell.Screen) (int, tcell.Screen) {
	return y, nil
}

type ErrWidthTooSmall struct{}

func (e ErrWidthTooSmall) Error() string {
	return "width must be at least 5"
}

type ErrHeightTooSmall struct{}

func (e ErrHeightTooSmall) Error() string {
	return "height must be at least 2"
}

const (
	Padding      = 2
	PaddingWidth = 4 // 2 * Padding
)

// REMEMBER TO INITIALIZE THE MESSAGEBOX WITH THE PADDING

type MessageBox struct {
	table  [][]rune
	styles []tcell.Style

	firstEmptyLine *rws.RWSafe[int] // -1 indicates the message box is closed

	width, height int

	isShifting bool

	msgBuffer   *buffer.Buffer[TextMessage]
	receiveFrom <-chan TextMessage

	notEmpty sync.Cond

	wg   sync.WaitGroup
	once sync.Once
}

func (mb *MessageBox) Init(width, height int) (chan<- TextMessage, error) {
	if width < 5 {
		return nil, &ErrWidthTooSmall{}
	} else if height < 2 {
		return nil, &ErrHeightTooSmall{}
	}

	var sendTo chan<- TextMessage

	mb.once.Do(func() {
		mb.table = make([][]rune, height)
		for i := 0; i < height; i++ {
			mb.table[i] = make([]rune, width)
		}

		mb.styles = make([]tcell.Style, height)

		mb.width = width
		mb.height = height

		mb.firstEmptyLine = rws.NewRWSafe(0)
		mb.isShifting = false

		sendTo, mb.receiveFrom = mb.msgBuffer.Init(1)

		mb.notEmpty = *sync.NewCond(&sync.Mutex{})

		mb.wg.Add(1)

		go mb.executeCommands()
	})

	return sendTo, nil
}

func (mb *MessageBox) GetDefaultStyle() tcell.Style {
	return StyleMap[NormalText]
}

func (mb *MessageBox) HasEmptyLine() bool {
	return mb.firstEmptyLine.Get() < mb.height
}

func (mb *MessageBox) Close() {

}

// Clear interface{} information to prevent a deadlock and release the memory
// FIXME: Check if this works
func (mb *MessageBox) Fini() {
	// Wake up the message box if it is waiting for a message
	// this will prevent a deadlock
	mb.firstEmptyLine.Set(-1)
	mb.notEmpty.Broadcast()

	// BAD: This shouldn't be done in the first place
}

func (mb *MessageBox) Wait() {
	mb.wg.Wait()
}

func (mb *MessageBox) Cleanup() {
	mb.wg.Wait()

	mb.table = nil
	mb.styles = nil
	mb.notEmpty.L = nil
	mb.firstEmptyLine = nil
}

////////////////////////////////

func (mb *MessageBox) executeCommands() {
	defer mb.wg.Done()

	for msg := range mb.receiveFrom {
		if msg.IsEmpty() {
			continue // Skip empty messages
		}

		// Get the style
		var style tcell.Style

		if val, ok := StyleMap[msg.GetType()]; ok {
			style = val
		} else {
			style = StyleMap[NormalText]
		}

		// Wait for an empty line
		mb.notEmpty.L.Lock()
		for !mb.HasEmptyLine() {
			mb.notEmpty.Wait()
		}

		// Prevent infinite wait time when the message box is closed
		if mb.firstEmptyLine.Get() == -1 {
			break
		}

		mb.notEmpty.L.Unlock()

		// Enqueue the message into the message box
		switch msg.GetType() {
		case BreakLine, SeparatorLine:
			var char rune

			if msg.GetType() == BreakLine {
				char = Space
			} else {
				char = Separator
			}

			for i := 0; i < mb.width; i++ {
				mb.table[mb.firstEmptyLine.Get()][i] = char
			}

			mb.firstEmptyLine.Set(mb.firstEmptyLine.Get() + 1)
		default:
			remaining := mb.enqueueContents(msg.GetContents(), style, false)

			for len(remaining) > 0 {
				if mb.HasEmptyLine() {
					remaining = mb.enqueueContents(remaining, style, true)
				} else {
					remaining = mb.outOfBoundsEnqueueContents(remaining, style, true)
				}
			}
		}
	}
}

// WARNING: This function doesn't do interface{} checks
// FIXME: This function doesn't work correctly
func (mb *MessageBox) enqueueContents(contents []string, style tcell.Style, isIndented bool) []string {
	fields := strings.Fields(contents[0])
	if len(fields) == 0 {
		return nil
	}

	x := 0
	if isIndented {
		x = IndentLevel
	}

	var y int
	x, y = mb.WriteFieldsAt(x, mb.firstEmptyLine.Get(), fields) // ASSUMPTION: This works correctly
	mb.styles[y] = style
	if y != mb.firstEmptyLine.Get() {
		mb.firstEmptyLine.Set(y)
		return contents[1:]
	}

	// Try to see if the next content fits on the same line
	for index, content := range contents[1:] {
		fields = strings.Fields(content)
		if len(fields) == 0 {
			continue // Skip empty lines
		}

		lastValidField := mb.CanWriteFieldsAt(x, fields)
		if lastValidField == -1 {
			mb.firstEmptyLine.Set(mb.firstEmptyLine.Get() + 1)
			return contents[index:]
		}

		for _, field := range fields[:lastValidField+1] {
			x = mb.WriteAt(x+2, y, field)
		}

		if lastValidField != len(fields)-1 {
			mb.firstEmptyLine.Set(mb.firstEmptyLine.Get() + 1)
			return append([]string{strings.Join(fields[lastValidField+1:], " ")}, contents[index+1:]...)
		}
	}

	return nil
}

// WARNING: This function doesn't do interface{} checks
// FIXME: This function doesn't work correctly
func (mb *MessageBox) outOfBoundsEnqueueContents(contents []string, style tcell.Style, isIndented bool) []string {
	fields := strings.Fields(contents[0])
	if len(fields) == 0 {
		return nil
	}

	x := 0
	if isIndented {
		x = IndentLevel
	}

	var y int
	x, y = mb.WriteFieldsAt(x, mb.firstEmptyLine.Get(), fields) // ASSUMPTION: This works correctly
	mb.styles[y] = style
	if y != mb.firstEmptyLine.Get() {
		mb.firstEmptyLine.Set(y)
		return contents[1:]
	}

	// Try to see if the next content fits on the same line
	for index, content := range contents[1:] {
		fields = strings.Fields(content)
		if len(fields) == 0 {
			continue // Skip empty lines
		}

		lastValidField := mb.CanWriteFieldsAt(x, fields)
		if lastValidField == -1 {
			mb.firstEmptyLine.Set(mb.firstEmptyLine.Get() + 1)
			return contents[index:]
		}

		for _, field := range fields[:lastValidField+1] {
			x = mb.WriteAt(x+2, y, field)
		}

		if lastValidField != len(fields)-1 {
			mb.firstEmptyLine.Set(mb.firstEmptyLine.Get() + 1)
			return append([]string{strings.Join(fields[lastValidField+1:], " ")}, contents[index+1:]...)
		}
	}

	return nil
}

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

func (mb *MessageBox) GetHeight() int {
	return mb.height
}

func (mb *MessageBox) ResizeWidth(width int) {
	if width == mb.width {
		return
	}

	newTable := make([][]rune, mb.height)
	for i := 0; i < mb.height; i++ {
		newTable[i] = make([]rune, width)
	}

	newStyles := make([]tcell.Style, mb.height)

	if width > mb.width {
		// Grow: Copy the old table and the old styles
		// but fill the new space with spaces

		for i, row := range mb.table {
			copy(newTable[i], row)

			for j := mb.width; j < width; j++ {
				newTable[i][j] = Space
			}
		}

		copy(newStyles, mb.styles)
	} else {
		// Shrink: Copy the old table and the old styles
		// but discard the extra characters (replacing them with hellipses)
		for i, row := range newTable {
			for j := range row {
				newTable[i][j] = mb.table[i][j]
			}

			if mb.table[i][width-1] == Space {
				continue
			}

			for j, char := range Hellip {
				newTable[i][width-HellipLen+j] = char
			}
		}

		for i := range newStyles {
			newStyles[i] = mb.styles[i]
		}
	}

	mb.table = newTable
	mb.styles = newStyles
	mb.width = width
}

func (mb *MessageBox) ResizeHeight(height int) {
	if height == mb.height {
		return
	}

	newTable := make([][]rune, height)
	for i := 0; i < height; i++ {
		newTable[i] = make([]rune, mb.width)
	}

	newStyles := make([]tcell.Style, height)

	if height > mb.height {
		// Grow: Copy the old table and the old styles
		// but fill the new lines with spaces

		for i, row := range mb.table {
			copy(newTable[i], row)
		}

		for i := mb.height; i < height; i++ {
			for j := 0; j < mb.width; j++ {
				newTable[i][j] = Space
			}
		}

		copy(newStyles, mb.styles)
	} else {
		// Shrink: Copy the old table and the old styles
		// but discard the older lines
		for i := height - 1; i >= 0; i-- {
			copy(newTable[i], mb.table[i])
			newStyles[i] = mb.styles[i]
		}

		mb.firstEmptyLine.Set(mb.firstEmptyLine.Get() - (mb.height - height))
	}

	mb.table = newTable
	mb.styles = newStyles
	mb.height = height
}

// Returns the first empty index after the last character
// WARNING: This function doesn't do interface{} checks
func (mb *MessageBox) WriteAt(x, y int, word string) int {
	row := mb.table[y]

	for _, char := range word {
		row[x] = char
		x++
	}

	return x
}

// WARNING: This function doesn't do interface{} checks
func (mb *MessageBox) CanWriteAt(x int, word string) bool {
	return x+len(word) <= mb.width
}

// Returns the first empty index after the last character
// WARNING: This function doesn't do interface{} checks
func (mb *MessageBox) WriteFieldsAt(x, y int, fields []string) (int, int) {
	if len(fields) == 1 {
		if mb.CanWriteAt(x, fields[0]) {
			return mb.WriteAt(x, y, fields[0]), y
		}

		x = mb.WriteAt(x, y, fields[0][:mb.width-x-HellipLen])
		return mb.WriteAt(x, y, Hellip), y
	} else if !mb.CanWriteAt(x, fields[0]) {
		x = mb.WriteAt(x, y, fields[0][:mb.width-x-HellipLen])
		mb.WriteAt(x, y, Hellip)

		return IndentLevel, y + 1
	}

	x = mb.WriteAt(x, y, fields[0])

	index := -1
	for i, field := range fields[1:] {
		if !mb.CanWriteAt(x+2, field) {
			index = i
			break
		}

		x = mb.WriteAt(x+2, y, field)
	}

	if index != -1 {
		mb.WriteAt(x, y, Hellip)

		return IndentLevel, y + 1
	}

	return x, y
}

// -1 = no fields, >= 0 index of the last field that fits
// WARNING: This function doesn't do interface{} checks
func (mb *MessageBox) CanWriteFieldsAt(x int, fields []string) int {
	if !mb.CanWriteAt(x, fields[0]) {
		return -1
	}

	x += len(fields[0])

	for i, field := range fields[1:] {
		if !mb.CanWriteAt(x+2, field) {
			return i
		}

		x += len(field) + 2
	}

	return len(fields) - 1
}

func (mb *MessageBox) CanShiftUp() bool {
	if mb.isShifting {
		if mb.firstEmptyLine.Get() == 0 {
			mb.isShifting = false
			return false
		}
	} else if mb.firstEmptyLine.Get() < mb.height/2 {
		return false
	}

	return true
}

func (mb *MessageBox) ShiftUp() {
	size := 1

	for size < mb.firstEmptyLine.Get() && mb.table[mb.firstEmptyLine.Get()-size][0] == Space {
		size++
	}

	for i := 0; i < mb.height-size; i++ {
		mb.table[i] = mb.table[i+size]
		mb.styles[i] = mb.styles[i+size]
	}

	mb.firstEmptyLine.Set(mb.firstEmptyLine.Get() - size)
	mb.isShifting = true
}

func (mb *MessageBox) Clear() {
	for i := 0; i < mb.height; i++ {
		for j := 0; j < mb.width; j++ {
			mb.table[i][j] = Space
		}
	}
}

func (mb *MessageBox) DrawHorizontalBorderAt(y int, style tcell.Style, screen tcell.Screen) tcell.Screen {
	screen.SetContent(0, y, '+', nil, style) // Left corner

	for x := 1; x < mb.width+1+PaddingWidth; x++ {
		screen.SetContent(x, y, '-', nil, style)
	}

	screen.SetContent(mb.width+1+PaddingWidth, y, '+', nil, style) // Right corner

	return screen
}

func (mb *MessageBox) SetScreen(y int, screen tcell.Screen) (int, tcell.Screen) {
	style := StyleMap[NormalText]

	screen = mb.DrawHorizontalBorderAt(y, style, screen) // Top border

	for y++; y < mb.height; y++ {
		screen.SetContent(0, y, '|', nil, style) // Left border

		if y < mb.firstEmptyLine.Get() {
			for i, cell := range mb.table[y] {
				screen.SetContent(Padding+i, y, cell, nil, mb.styles[y])
			}
		} else {
			for i := 0; i < mb.width; i++ {
				screen.SetContent(Padding+i, y, Space, nil, style)
			}
		}

		screen.SetContent(mb.width+PaddingWidth, y, '|', nil, style) // Right border
	}

	return y, mb.DrawHorizontalBorderAt(y, style, screen) // Bottom border
}
