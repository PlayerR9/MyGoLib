package MessageBox

import (
	"sync"

	buffer "github.com/PlayerR9/MyGoLib/CustomData/Buffer"
	"github.com/gdamore/tcell"
)

const (
	Padding      = 2
	PaddingWidth = 4 // 2 * Padding
)

type MessageBox struct {
	content *ContentBox

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
		mb.content, _ = NewContentBox(width, height)

		sendTo, mb.receiveFrom = mb.msgBuffer.Init(1)

		mb.notEmpty = *sync.NewCond(&sync.Mutex{})

		mb.wg.Add(1)

		go mb.executeCommands()
	})

	return sendTo, nil
}

func (mb *MessageBox) SetSize(width, height int) error {
	err := mb.content.ResizeWidth(width)
	if err != nil {
		return &ErrWidthTooSmall{}
	}

	err = mb.content.ResizeHeight(height)
	if err != nil {
		return &ErrHeightTooSmall{}
	}

	return nil
}

func (mb *MessageBox) drawHorizontalBorderAt(y int, style tcell.Style, screen tcell.Screen) tcell.Screen {
	screen.SetContent(0, y, '+', nil, style) // Left corner

	for x := 1; x < mb.content.width+1+PaddingWidth; x++ {
		screen.SetContent(x, y, '-', nil, style)
	}

	screen.SetContent(mb.content.width+1+PaddingWidth, y, '+', nil, style) // Right corner

	return screen
}

func (mb *MessageBox) Draw(y int, screen tcell.Screen) (int, tcell.Screen) {
	style := StyleMap[NormalText]

	screen = mb.drawHorizontalBorderAt(y, style, screen) // Top border

	for i := 0; i < mb.content.firstEmptyLine.Get(); i++ {
		y++
		screen.SetContent(0, y, '|', nil, style) // Left border

		for j, cell := range mb.content.table[i] {
			screen.SetContent(Padding+j, y, cell, nil, mb.content.styles[i])
		}

		screen.SetContent(mb.content.width+PaddingWidth, y, '|', nil, style) // Right border
	}

	for i := mb.content.firstEmptyLine.Get(); i < mb.content.height; i++ {
		y++
		screen.SetContent(0, y, '|', nil, style) // Left border

		for j, cell := range mb.content.table[i] {
			screen.SetContent(Padding+j, y, cell, nil, mb.content.styles[i])
		}

		screen.SetContent(mb.content.width+PaddingWidth, y, '|', nil, style) // Right border
	}

	y++
	screen = mb.drawHorizontalBorderAt(y, style, screen) // Bottom border

	return y, screen
}

////////////////// OLD CODE //////////////////////

// REMEMBER TO INITIALIZE THE MESSAGEBOX WITH THE PADDING

func (mb *MessageBox) Close() {

}

// Clear interface{} information to prevent a deadlock and release the memory
// FIXME: Check if this works
func (mb *MessageBox) Fini() {
	// Wake up the message box if it is waiting for a message
	// this will prevent a deadlock
	mb.notEmpty.Broadcast()

	// BAD: This shouldn't be done in the first place
}

func (mb *MessageBox) Wait() {
	mb.wg.Wait()
}

func (mb *MessageBox) Cleanup() {
	mb.wg.Wait()

	mb.content.Cleanup()
	mb.notEmpty.L = nil
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
		for !mb.content.HasEmptyLine() {
			mb.notEmpty.Wait()
		}

		// Prevent infinite wait time when the message box is closed
		if mb.content.firstEmptyLine.Get() == -1 {
			break
		}

		mb.notEmpty.L.Unlock()

		// Enqueue the message into the message box
		switch msg.GetType() {
		case BreakLine:
			mb.content.EnqueueLineSeparator(Space)
		case SeparatorLine:
			mb.content.EnqueueLineSeparator(Separator)
		default:
			mb.content.EnqueueContents(msg.GetContents(), style)
		}

		if mb.content.CanShiftUp() {
			mb.content.ShiftUp()
		}
	}
}
