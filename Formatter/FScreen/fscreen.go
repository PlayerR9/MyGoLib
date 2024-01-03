package FScreen

import (
	"fmt"
	"strings"
	"sync"

	bf "github.com/PlayerR9/MyGoLib/CustomData/Buffer"
	h "github.com/PlayerR9/MyGoLib/Formatter/FScreen/Header"
	mb "github.com/PlayerR9/MyGoLib/Formatter/FScreen/MessageBox"
	"github.com/gdamore/tcell"
)

const (
	Padding      = 2
	PaddingWidth = 4 // 2 * Padding
)

var (
	messageBox *mb.MessageBox = nil
	screen     tcell.Screen   = nil
)

type FScreen struct {
	h *h.Header
	d *Display

	msgBuffer           *bf.Buffer[interface{}]
	sendToHeader        chan<- h.HeaderMessage
	sendToMessageBox    chan<- mb.TextMessage
	receiveFrom         <-chan interface{}
	receiveHeaderErrors <-chan mb.TextMessage

	sendToDisplay        chan<- DisplayOpcode
	receiveDisplayErrors <-chan mb.TextMessage

	wg   sync.WaitGroup
	once sync.Once
}

/*
REMEMBER TO INITIALIZE THE SCREEN IN MAIN PROGRAM

screen, err := tcell.NewScreen()
	if err != nil {
		return err
	}

	err = screen.Init()
	if err != nil {
		return err
	}
	screen.Clear()
*/

func (fs *FScreen) Init(title string) chan<- interface{} {
	if strings.TrimSpace(title) == "" {
		panic("title cannot be empty")
	}

	width, height := screen.Size()

	if width < 2*Padding || height < Padding {
		panic(fmt.Sprintf("screen too small: %dx%d", width, height))
	}

	var sendTo chan<- interface{}

	fs.once.Do(func() {
		fs.sendToMessageBox = messageBox.Init(width-2*Padding, height-Padding)
		fs.sendToHeader, fs.receiveHeaderErrors = fs.h.Init(title)
		sendTo, fs.receiveFrom = fs.msgBuffer.Init(1)
		fs.sendToDisplay, fs.receiveDisplayErrors = fs.d.Init(messageBox, fs.h)

		fs.wg.Add(1)

		go fs.routingMessages()
	})

	return sendTo
}

func (fs *FScreen) routingMessages() {
	defer fs.wg.Done()

	// Start the error channel listeners
	var errorWg sync.WaitGroup

	errorWg.Add(2)

	go func() {
		for msg := range fs.receiveHeaderErrors {
			fs.sendToMessageBox <- msg
		}

		errorWg.Done()
	}()

	go func() {
		for msg := range fs.receiveDisplayErrors {
			fs.sendToMessageBox <- msg
		}

		errorWg.Done()
	}()

	for msg := range fs.receiveFrom {
		switch x := msg.(type) {
		case mb.TextMessage:
			fs.sendToMessageBox <- x
		case h.HeaderMessage:
			fs.sendToHeader <- x
		default:
			fs.sendToMessageBox <- mb.NewTextMessage(mb.FatalText,
				fmt.Sprintf("Unknown message type: %T", x),
			)
		}
	}

	errorWg.Wait()
}

func (fs *FScreen) Cleanup() {
	fs.wg.Wait()

	var finiWg sync.WaitGroup
	defer finiWg.Wait()

	finiWg.Add(4)

	go func() {
		fs.d.Cleanup()
		finiWg.Done()
		fs.d = nil
	}()

	go func() {
		screen.Fini()
		finiWg.Done()
		screen = nil
	}()

	go func() {
		fs.msgBuffer.Cleanup()
		finiWg.Done()
		fs.msgBuffer = nil
	}()

	go func() {
		close(fs.sendToMessageBox)
		messageBox.Fini()
		messageBox.Cleanup()
		finiWg.Done()
		messageBox = nil
	}()
}

type DisplayOpcode int

const (
	DOMustPrint DisplayOpcode = iota
	DOShouldPrint
)

func (enum DisplayOpcode) String() string {
	return [...]string{
		"MustPrint",
		"ShouldPrint",
	}[enum]
}

type Display struct {
	box *mb.MessageBox
	h   *h.Header

	screen tcell.Screen
	style  tcell.Style

	msgBuffer    *bf.Buffer[DisplayOpcode]
	sendToHeader chan<- h.HeaderMessage
	receiveFrom  <-chan DisplayOpcode

	receiveErrors chan mb.TextMessage
	wg            sync.WaitGroup
	once          sync.Once
}

func (d *Display) Init(box *mb.MessageBox, h *h.Header) (chan<- DisplayOpcode, <-chan mb.TextMessage) {
	var sendTo chan<- DisplayOpcode

	d.once.Do(func() {
		d.box = box
		d.h = h

		sendTo, d.receiveFrom = d.msgBuffer.Init(1)
		d.receiveErrors = make(chan mb.TextMessage, 1)
		d.style = box.GetDefaultStyle()

		d.wg.Add(1)

		go d.executeOpcodes()
	})

	return sendTo, d.receiveErrors
}

func (d *Display) executeOpcodes() {
	printScreen := func() {
		y := 0 // y is the current line

		d.screen.Clear()

		width, height := d.screen.Size()
		d.box.ResizeWidth(width - 2*Padding)
		d.box.ResizeHeight(height - Padding)

		y, d.screen = d.h.SetScreen(y, width, d.style, d.screen)
		y += 2

		_, d.screen = d.box.SetScreen(y, d.screen)

		d.screen.Show()
	}

	defer d.wg.Done()
	shouldHavePrinted := false

	for opcode := range d.receiveFrom {
		switch opcode {
		case DOMustPrint:
			printScreen()

			shouldHavePrinted = false
		case DOShouldPrint:
			shouldHavePrinted = true
		default:
			d.receiveErrors <- mb.NewTextMessage(mb.FatalText,
				fmt.Sprintf("Unknown opcode: %T", opcode),
			)
		}
	}

	if shouldHavePrinted {
		printScreen()
	}

	close(d.receiveErrors)
}

func (d *Display) Wait() {
	d.wg.Wait()
}

func (d *Display) Cleanup() {
	d.wg.Wait()

	var finiWg sync.WaitGroup
	defer finiWg.Wait()

	finiWg.Add(4)

	go func() {
		d.box.Fini()
		finiWg.Done()
		d.box = nil
	}()

	go func() {
		close(d.sendToHeader)
		d.h.Cleanup()
		finiWg.Done()
		d.h = nil
	}()

	go func() {
		d.screen.Fini()
		finiWg.Done()
		d.screen = nil
	}()

	go func() {
		d.msgBuffer.Cleanup()
		finiWg.Done()
		d.msgBuffer = nil
	}()

	d.receiveErrors = nil
}
