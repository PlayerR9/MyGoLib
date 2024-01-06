// git tag v0.1.47

package Display

import (
	"math"
	"sync"
	"time"

	"github.com/gdamore/tcell"
)

// Displayable defines the behavior of an object that can be displayed
// on a screen.
//
// Any type that implements these methods is said to satisfy the
// Displayable interface.
// This means that the type can be used wherever a Displayable is expected.
type Displayable interface {
	// Checks if the width of the object can be set to a given value.
	// Returns a boolean.
	CanSetWidth(width int) bool

	// Sets the width of the object to a given value.
	SetWidth(width int)

	// Checks if the height of the object can be set to a given value.
	// Returns a boolean.
	CanSetHeight(height int) bool

	// Sets the height of the object to a given value.
	SetHeight(height int)

	// Draws the object on a given screen and returns the updated screen.
	GenerateDrawTables() ([][]rune, []tcell.Style)
}

// screen is a global variable of type tcell.Screen.
// It is initially set to nil.
//
// The screen variable represents the screen that is currently being used
// for display.
// It is used throughout the display.go file to draw components, handle user
// input, and manage the display's state.
//
// The screen is initialized in the Init method of the Display struct, and it
// is finalized and set to nil in the Cleanup method.
var screen tcell.Screen = nil

// Display represents a display screen with a frame rate.
// It also includes synchronization primitives for managing the display's
// state.
type Display struct {
	// The duration between frames on the display. It is of type time.Duration.
	frameRate time.Duration

	// A boolean that indicates whether a close signal has been received.
	closeReceived bool

	// A sync.WaitGroup that allows the program to wait for all goroutines
	// related to this display to finish before continuing.
	wg sync.WaitGroup

	// A sync.Once that ensures that a certain operation (usually cleanup
	// or initialization) is performed only once.
	once sync.Once
}

// Init initializes the Display with a given frame rate.
//
// Parameters:
//
//   - frameRate: The frame rate for the Display, specified as a float64.
//   - defaultStyle: The default style for the screen. It must be of type
//     tcell.Style.
//
// Returns:
//   - An error if any occurred during initialization.
//
// The frame rate is converted to a time.Duration by dividing 1000 by the
// frame rate, rounding the result, and multiplying by time.Millisecond.
// This gives the duration between frames for the Display.
func (d *Display) Init(frameRate float64, defaultStyle tcell.Style) error {
	var err error

	// Initialize the screen
	screen, err = tcell.NewScreen()
	if err != nil {
		return err
	}

	// Set the default style for the screen
	screen.SetStyle(defaultStyle)

	// Initialize the screen
	err = screen.Init()
	if err != nil {
		return err
	}
	screen.Clear() // Clear the screen to avoid artifacts

	d.frameRate = time.Duration(math.Round(1000/frameRate)) * time.Millisecond
	d.closeReceived = false

	return nil
}

// Start begins the main loop of the Display, which clears the screen,
// sets the sizes of the element to display, draws the element, shows
// the screen, and then waits for the duration of the frame rate.
// This loop continues until the Display.Close() method is called or
// the program terminates unexpectedly.
//
// Parameters:
//
//   - elementToDisplay: The element to display on the screen. It must
//     satisfy the Displayable interface.
//
// The method uses a sync.Once to ensure that the main loop is only
// started once.
//
// If an error occurs while setting the size of the element to display,
// the method panics with an appropriate message.
//
// After the element has been drawn and shown, the method waits for the
// duration of the frame rate before starting the next iteration of the
// loop.
//
// When the loop ends, the method sets the closeReceived field to true,
// clears the screen, finalizes it, and sets it to nil.
func (d *Display) Start(elementToDisplay Displayable) {
	d.once.Do(func() {
		previousWidth, previousHeight := -1, -1

		// Cleanup procedure
		defer func() {
			d.closeReceived = true
			screen.Clear() // Clear the screen
			screen.Fini()
			screen = nil
		}()

		for !d.closeReceived {
			screen.Clear()

			width, height := screen.Size()

			// Set the size of the element to display
			if elementToDisplay.CanSetWidth(width) {
				previousWidth = width
			} else if previousWidth == -1 {
				panic("Cannot set width of element to display")
			}

			elementToDisplay.SetWidth(previousWidth)

			if elementToDisplay.CanSetHeight(height) {
				previousHeight = height
			} else if previousHeight == -1 {
				panic("Cannot set height of element to display")
			}

			elementToDisplay.SetHeight(previousHeight)

			// Draw the element to display
			tables, styles := elementToDisplay.GenerateDrawTables()

			for i := 0; i < len(tables); i++ {
				for j := 0; j < len(tables[i]); j++ {
					screen.SetContent(j, i, tables[i][j], nil, styles[i])
				}
			}

			screen.Show()

			// Wait for some time
			time.Sleep(d.frameRate)
		}
	})
}

// Stop signals to the main loop of the Display to stop and waits
// for all goroutines related to the Display to finish.
//
// The method should be called when the Display is no longer needed,
// to ensure that all resources are properly released and all
// goroutines have finished.
func (d *Display) Stop() {
	d.closeReceived = true
	d.wg.Wait()
}
