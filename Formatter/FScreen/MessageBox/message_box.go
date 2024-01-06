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

var StyleMap map[TextMessageType]tcell.Style = map[TextMessageType]tcell.Style{
	NormalText:  tcell.StyleDefault.Foreground(tcell.ColorBlack).Background(tcell.ColorGhostWhite),
	DebugText:   tcell.StyleDefault.Bold(true).Foreground(tcell.ColorSlateGray).Background(tcell.ColorGhostWhite),
	WarningText: tcell.StyleDefault.Foreground(tcell.ColorDarkOrange).Background(tcell.ColorGhostWhite),
	ErrorText:   tcell.StyleDefault.Foreground(tcell.ColorFireBrick).Background(tcell.ColorGhostWhite),
	FatalText:   tcell.StyleDefault.Bold(true).Foreground(tcell.ColorDarkRed).Background(tcell.ColorGhostWhite),
	SuccessText: tcell.StyleDefault.Foreground(tcell.ColorDarkGreen).Background(tcell.ColorGhostWhite),
}

type MessageBox struct {
	content ContentBox

	msgBuffer     *buffer.Buffer[TextMessage]
	receiveFrom   <-chan TextMessage
	receiveErrors chan TextMessage

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

		mb.msgBuffer = new(buffer.Buffer[TextMessage])
		sendTo, mb.receiveFrom = mb.msgBuffer.Init(1)
		mb.receiveErrors = make(chan TextMessage, 1)

		mb.wg.Add(1)

		go mb.executeCommands()
	})

	return sendTo, nil
}

func (mb *MessageBox) GetReceiveErrorsFromChannel() <-chan TextMessage {
	return mb.receiveErrors
}

func (mb *MessageBox) SetSize(width, height int) error {
	if width-PaddingWidth < 5 {
		return &ErrWidthTooSmall{}
	} else if height-Padding < 2 {
		return &ErrHeightTooSmall{}
	}

	mb.content.ResizeWidth(width)
	mb.content.ResizeHeight(height)

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

	for i := 0; i < mb.content.firstEmptyLine; i++ {
		y++
		screen.SetContent(0, y, '|', nil, style) // Left border

		for j, cell := range mb.content.table[i] {
			screen.SetContent(Padding+j, y, cell, nil, mb.content.styles[i])
		}

		screen.SetContent(mb.content.width+PaddingWidth, y, '|', nil, style) // Right border
	}

	for i := mb.content.firstEmptyLine; i < mb.content.height; i++ {
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

func (mb *MessageBox) Wait() {
	mb.wg.Wait()
}

func (mb *MessageBox) Cleanup() {
	mb.wg.Wait()

	var finiWg sync.WaitGroup

	finiWg.Add(1)

	go func() {
		mb.msgBuffer.Cleanup()
		finiWg.Done()
	}()

	finiWg.Wait()

	close(mb.receiveErrors)
}

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

		if mb.content.CanShiftUp() {
			mb.content.ShiftUp()
		}

		// Enqueue the message into the message box
		switch msg.GetType() {
		case BreakLine:
			mb.content.EnqueueLineSeparator(Space)
		case SeparatorLine:
			mb.content.EnqueueLineSeparator(Separator)
		default:
			err := mb.content.EnqueueContents(msg.GetContents(), style)
			if err != nil {
				mb.receiveErrors <- NewTextMessage(ErrorText,
					"Could not enqueue the message into the message box:",
					err.Error(),
				)
			}
		}
	}
}

////////////////// OLD CODE //////////////////////

// Clear interface{} information to prevent a deadlock and release the memory
// FIXME: Check if this works
func (mb *MessageBox) Fini() {

	// BAD: This shouldn't be done in the first place
}

////////////////////////////////
