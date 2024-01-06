// git tag v0.1.47

package FScreen

import (
	"fmt"
	"strings"
	"sync"

	buffer "github.com/PlayerR9/MyGoLib/CustomData/Buffer"
	d "github.com/PlayerR9/MyGoLib/CustomData/Display"
	h "github.com/PlayerR9/MyGoLib/Formatter/FScreen/Header"
	mb "github.com/PlayerR9/MyGoLib/Formatter/FScreen/MessageBox"
	"github.com/gdamore/tcell"
)

// Default values for the FScreen
const (
	DefaultScreenWidth  int = 80
	DefaultScreenHeight int = 24
)

var (
	DefaultStyle tcell.Style = mb.StyleMap[mb.NormalText]
)

// Global variables for the FScreen
var (
	header     *h.Header      = nil
	messageBox *mb.MessageBox = nil
	display    *d.Display     = nil
)

type FScreen struct {
	msgBuffer   *buffer.Buffer[interface{}]
	receiveFrom <-chan interface{}

	width int

	sendToHeader        chan<- h.HeaderMessage
	receiveHeaderErrors <-chan mb.TextMessage

	sendToMessageBox        chan<- mb.TextMessage
	receiveMessageBoxErrors <-chan mb.TextMessage

	once sync.Once
	wg   sync.WaitGroup
}

func (fs *FScreen) Init(title string, frameRate float64) (chan<- interface{}, error) {
	var err error
	var sendTo chan<- interface{}

	// Initialize the header
	header = new(h.Header)
	fs.sendToHeader, err = header.Init(title)
	if err != nil {
		return nil, err
	}

	// Initialize the message box
	messageBox = new(mb.MessageBox)
	fs.sendToMessageBox, err = messageBox.Init(DefaultScreenWidth, DefaultScreenHeight)
	if err != nil {
		return nil, err
	}

	// Initialize the display
	display = new(d.Display)
	err = display.Init(frameRate, DefaultStyle)
	if err != nil {
		return nil, err
	}

	fs.once.Do(func() {
		// Initialize the FScreen
		fs.msgBuffer = new(buffer.Buffer[interface{}])
		sendTo, fs.receiveFrom = fs.msgBuffer.Init(1)

		// Initialize the channels
		fs.receiveHeaderErrors = header.GetReceiveErrorsFromChannel()
		fs.receiveMessageBoxErrors = messageBox.GetReceiveErrorsFromChannel()

		fs.wg.Add(1)

		go fs.routingMessages()
	})

	return sendTo, nil
}

func (fs *FScreen) routingMessages() {
	defer fs.wg.Done()

	// Start the rerouting channel listeners
	var reroutingWg sync.WaitGroup

	reroutingWg.Add(2)

	go func() {
		defer reroutingWg.Done()

		for msg := range fs.receiveHeaderErrors {
			fs.sendToMessageBox <- msg
		}
	}()

	go func() {
		defer reroutingWg.Done()

		for msg := range fs.receiveMessageBoxErrors {
			fs.sendToMessageBox <- msg
		}
	}()

	for msg := range fs.receiveFrom {
		switch x := msg.(type) {
		case mb.TextMessage:
			fs.sendToMessageBox <- x
		case h.HeaderMessage:
			switch t := x.Data.(type) {
			case mb.TextMessage:
				fs.sendToMessageBox <- t
			case string, h.CounterData:
				fs.sendToHeader <- x
			default:
				fs.sendToMessageBox <- mb.NewTextMessage(mb.ErrorText,
					fmt.Sprintf("Unknown header message type: %T", x),
				)
			}
		default:
			fs.sendToMessageBox <- mb.NewTextMessage(mb.FatalText,
				fmt.Sprintf("Unknown message type: %T", x),
			)
		}
	}

	reroutingWg.Wait()
}

func (fs *FScreen) Cleanup() {
	fs.wg.Wait()

	var cleanWg sync.WaitGroup
	defer cleanWg.Wait()

	cleanWg.Add(3)

	go func() {
		fs.msgBuffer.Cleanup()
		cleanWg.Done()
		fs.msgBuffer = nil
		fs.receiveFrom = nil
	}()

	go func() {
		close(fs.sendToMessageBox)
		messageBox.Cleanup()
		cleanWg.Done()
		messageBox = nil
	}()

	go func() {
		close(fs.sendToHeader)
		header.Cleanup()
		cleanWg.Done()
		header = nil
		fs.sendToHeader = nil
		fs.receiveHeaderErrors = nil
	}()
}

func (fs *FScreen) Wait() {
	fs.wg.Wait()
}

func (fs *FScreen) CanSetWidth(width int) bool {
	ok := header.CanSetWidth(width)
	if !ok {
		return false
	}

	return messageBox.CanSetWidth(width)
}

func (fs *FScreen) SetWidth(width int) {
	fs.width = width

	header.SetWidth(width)
	messageBox.SetWidth(width)
}

func (fs *FScreen) CanSetHeight(height int) bool {
	return header.GetCurrentHeight()+2+messageBox.GetCurrentHeight() >= height
}

func (fs *FScreen) SetHeight(height int) {
	n := header.GetCurrentHeight()

	header.SetHeight(n)
	messageBox.SetHeight(height - n - 1)
}

func (fs *FScreen) GenerateDrawTables() ([][]rune, []tcell.Style) {
	emptyLine := []rune(strings.Repeat(" ", fs.width))

	headerTable, headerStyle := header.GenerateDrawTables()
	messageBoxTable, messageBoxStyle := messageBox.GenerateDrawTables()

	headerTable = append(headerTable, emptyLine)
	headerStyle = append(headerStyle, mb.StyleMap[mb.NormalText])

	for i, row := range messageBoxTable {
		headerTable = append(headerTable, row)
		headerStyle = append(headerStyle, messageBoxStyle[i])
	}

	return headerTable, headerStyle
}
