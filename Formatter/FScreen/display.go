// git tag v0.1.39

package FScreen

import (
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/gdamore/tcell"
)

type Displayer interface {
	SetSize(width, height int) error
	Draw(int, tcell.Screen) (int, tcell.Screen)
}

type Display struct {
	elementsToShow []Displayer
	frameRate      time.Duration

	screen        tcell.Screen
	style         tcell.Style
	closeReceived bool

	wg   sync.WaitGroup
	once sync.Once
}

func NewDisplay(frameRate float64, elementToDisplay []Displayer) (Display, error) {
	// Initialize the screen
	screen, err := tcell.NewScreen()
	if err != nil {
		return Display{}, err
	}

	err = screen.Init()
	if err != nil {
		return Display{}, err
	}
	screen.Clear()

	return Display{
		elementsToShow: elementToDisplay,
		frameRate:      time.Duration(math.Round(1000/frameRate)) * time.Millisecond,
		screen:         screen,
		style:          tcell.StyleDefault,
		closeReceived:  false,
	}, nil
}

func (d *Display) Start() {
	d.once.Do(func() {
		d.wg.Add(1)
		defer d.wg.Done()

		for !d.closeReceived {
			d.screen.Clear()

			width, height := d.screen.Size()
			for _, element := range d.elementsToShow {
				err := element.SetSize(width, height)
				if err != nil {
					// TO DO: Handle error
					panic(fmt.Errorf("error setting size for element: %v", err))
				}
			}

			y := 0 // y is the current line
			var offset int

			for _, element := range d.elementsToShow {
				offset, d.screen = element.Draw(y, d.screen)
				y += offset + 2
			}

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

	d.elementsToShow = nil
	d.screen = nil
}
