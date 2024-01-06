package MessageBox

import (
	"strings"
	"sync"

	buffer "github.com/PlayerR9/MyGoLib/CustomData/Buffer"
	"github.com/gdamore/tcell"
)

const (
	Padding      int = 2
	PaddingWidth int = 4 // 2 * Padding
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

func (mb *MessageBox) GenerateDrawTables() ([][]rune, []tcell.Style) {
	border := []rune(strings.Repeat("-", mb.content.width+PaddingWidth))
	border = append([]rune("+"), border...) // Left corner
	border = append(border, '+')            // Right corner

	tables, styles := mb.content.GenerateDrawTables()

	// 1. Add the padding to the left and right of each line
	for i := 0; i < mb.content.height; i++ {
		tables[i] = append([]rune{'|', Space}, tables[i]...)
		tables[i] = append(tables[i], Space, '|')
	}

	// 2. Add a border to the top and bottom of the table
	tables = append([][]rune{border}, tables...)
	styles = append([]tcell.Style{StyleMap[NormalText]}, styles...)
	tables = append(tables, border)
	styles = append(styles, StyleMap[NormalText])

	return tables, styles
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

func (mb *MessageBox) CanSetWidth(width int) bool {
	return mb.content.CanSetWidth(width - PaddingWidth)
}

func (mb *MessageBox) SetWidth(width int) {
	mb.content.SetWidth(width)
}

func (mb *MessageBox) CanSetHeight(height int) bool {
	return mb.content.CanSetHeight(height - Padding)
}

func (mb *MessageBox) SetHeight(height int) {
	mb.content.SetHeight(height - Padding)
}

func (mb *MessageBox) GetCurrentHeight() int {
	return mb.content.GetCurrentHeight() + Padding
}
