package FScreen

import (
	"fmt"
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
)

type FScreen struct {
	h *h.Header

	screen tcell.Screen

	messageChannel bf.Buffer[interface{}]

	wg   sync.WaitGroup
	once sync.Once
}

func NewFScreen(title string) (*FScreen, error) {
	fscreen := &FScreen{
		messageChannel: bf.NewBuffer[interface{}](),
	}

	if header, err := h.NewHeader(title); err != nil {
		return nil, fmt.Errorf("could not create header: %w", err)
	} else {
		fscreen.h = header
	}

	screen, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}

	err = screen.Init()
	if err != nil {
		return nil, err
	}
	screen.Clear()

	width, height := screen.Size()
	messageBox = mb.NewMessageBox(width-2*Padding, height-Padding)

	fscreen.screen = screen

	return fscreen, nil
}

func (fs *FScreen) Run() error {
	fs.once.Do(func() {
		fs.wg.Add(1)

		// Start the header, message box, and display

		display := NewDisplay(messageBox, fs.h)

		go fs.h.Run()
		go messageBox.Run()
		go display.Run()

		// Start the error channel listeners
		var errorWg sync.WaitGroup

		errorWg.Add(2)

		go func() {
			for msg := range fs.h.ErrorChannel {
				messageBox.SendMessages(msg)
			}

			errorWg.Done()
		}()

		go func() {
			for msg := range display.ErrorChannel {
				messageBox.SendMessages(msg)
			}

			errorWg.Done()
		}()

		for {
			msg, ok := fs.messageChannel.Get()
			if !ok {
				break
			}

			switch x := msg.(type) {
			case mb.TextMessage:
				messageBox.SendMessages(x)
			case h.HeaderMessage:
				fs.h.SendMessages(x)
			default:
				messageBox.SendMessages(
					mb.NewTextMessage(mb.FatalText,
						fmt.Sprintf("Unknown message type: %T", x),
					),
				)
			}
		}

		// Release some resources

		var finiWg sync.WaitGroup

		finiWg.Add(2)

		go func() {
			display.Fini()
			finiWg.Done()
			display = nil
		}()

		go func() {
			fs.screen.Fini()
			finiWg.Done()
			fs.screen = nil
		}()

		finiWg.Wait()

		errorWg.Wait()

		fs.wg.Done()
	})

	return nil
}

func (fs *FScreen) Fini() {
	fs.messageChannel.Fini()
	close(fs.h.ErrorChannel)

	fs.wg.Wait()

	// Close and free resources
	var finiWg sync.WaitGroup

	finiWg.Add(2)

	go func() {
		messageBox.Fini()
		finiWg.Done()
		messageBox = nil
	}()

	go func() {
		fs.h.Fini()
		finiWg.Done()
		fs.h = nil
	}()

	finiWg.Wait()
}

func (fs *FScreen) SendMessages(message interface{}, optionalMessages ...interface{}) {
	fs.messageChannel.SendMessages(
		message,
		optionalMessages...,
	)
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

	messageChannel bf.Buffer[DisplayOpcode]
	ErrorChannel   chan mb.TextMessage
	wg             sync.WaitGroup
	once           sync.Once
}

func NewDisplay(box *mb.MessageBox, h *h.Header) *Display {
	return &Display{
		box: box,
		h:   h,

		screen: nil,
		style:  box.GetDefaultStyle(),

		messageChannel: bf.NewBuffer[DisplayOpcode](),
		ErrorChannel:   make(chan mb.TextMessage),
	}
}

func (d *Display) Run() {
	d.once.Do(func() {
		go d.messageChannel.Run()

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

		d.wg.Add(1)

		shouldHavePrinted := false

		for {
			opcode, ok := d.messageChannel.Get()
			if !ok {
				if shouldHavePrinted {
					printScreen()
				}

				d.wg.Done()

				break
			}

			switch opcode {
			case DOMustPrint:
				printScreen()

				shouldHavePrinted = false
			case DOShouldPrint:
				shouldHavePrinted = true
			default:
				d.ErrorChannel <- mb.NewTextMessage(mb.FatalText,
					fmt.Sprintf("Unknown opcode: %T", opcode),
				)
			}
		}
	})
}

func (d *Display) Fini() {
	d.messageChannel.Fini()
	close(d.ErrorChannel)

	d.wg.Wait()

	// Close and free resources
	var finiWg sync.WaitGroup

	finiWg.Add(3)

	go func() {
		d.box.Fini()
		finiWg.Done()
		d.box = nil
	}()

	go func() {
		d.h.Fini()
		finiWg.Done()
		d.h = nil
	}()

	go func() {
		d.screen.Fini()
		finiWg.Done()
		d.screen = nil
	}()

	d.ErrorChannel = nil

	finiWg.Wait()
}
