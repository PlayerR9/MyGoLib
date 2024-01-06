// git tag v0.1.42

package FScreen

import (
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/gdamore/tcell"
)

type Display struct {
	frameRate time.Duration

	screen        tcell.Screen
	style         tcell.Style
	closeReceived bool

	wg   sync.WaitGroup
	once sync.Once
}

func NewDisplay(frameRate float64) (Display, error) {
	// Initialize the screen
	screen, err := tcell.NewScreen()
	if err != nil {
		return Display{}, err
	}

	screen.SetStyle(DefaultStyle)

	err = screen.Init()
	if err != nil {
		return Display{}, err
	}
	screen.Clear()

	return Display{
		frameRate:     time.Duration(math.Round(1000/frameRate)) * time.Millisecond,
		screen:        screen,
		style:         tcell.StyleDefault,
		closeReceived: false,
	}, nil
}

func (d *Display) Start() {
	d.once.Do(func() {
		d.wg.Add(1)
		defer d.wg.Done()

		for !d.closeReceived {
			d.screen.Clear()

			width, height := d.screen.Size()

			// Set the size of the header
			err := header.SetSize(width, height)
			if err != nil {
				// TO DO: Handle error
				panic(fmt.Errorf("error setting size for header: %v", err))
			}

			// Set the size of the message box
			err = messageBox.SetSize(width, height)
			if err != nil {
				// TO DO: Handle error
				panic(fmt.Errorf("error setting size for message box: %v", err))
			}

			y := 0 // y is the current line
			var offset int

			// Draw the header
			offset, d.screen = header.Draw(y, d.screen)
			y += offset + 2

			// Draw the message box
			offset, d.screen = messageBox.Draw(y, d.screen)
			y += offset + 2

			d.screen.Show()

			// Wait for some time
			time.Sleep(d.frameRate)
		}
	})
}

func (d *Display) Stop() {
	d.closeReceived = true
}

func (d *Display) Wait() {
	d.wg.Wait()
}

func (d *Display) Cleanup() {
	d.wg.Wait()

	var finiWg sync.WaitGroup
	defer finiWg.Wait()

	finiWg.Add(1)

	go func() {
		d.screen.Fini()
		finiWg.Done()
		d.screen = nil
	}()

	d.screen = nil
}
