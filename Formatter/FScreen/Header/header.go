// git tag v0.1.48

package Header

import (
	"fmt"
	"slices"
	"strings"
	"sync"

	bf "github.com/PlayerR9/MyGoLib/CustomData/Buffer"
	"github.com/PlayerR9/MyGoLib/CustomData/Counters"
	mb "github.com/PlayerR9/MyGoLib/Formatter/FScreen/MessageBox"
	sext "github.com/PlayerR9/MyGoLib/Utility/StrExt"
	"github.com/gdamore/tcell"
)

const (
	Hellip string = "..."
)

// style is a variable of type tcell.Style. It represents the style used
// for normal text in the header.
// The style is fetched from the StyleMap of MessageBox using the key
// MessageBox.NormalText.
var style tcell.Style = mb.StyleMap[mb.NormalText]

// Header represents the header of a process or application.
// It contains information about the title, current process, counters, message buffer,
// channels for receiving messages and errors, synchronization primitives, and the width
// of the header.
type Header struct {
	// title represents the title of the header.
	title string

	// currentProcess represents the current process of the header.
	currentProcess string

	// counters represents a slice of up counters.
	counters []*Counters.UpCounter

	// msgBuffer represents the message buffer for the header.
	msgBuffer *bf.Buffer[HeaderMessage]

	// receiveFrom represents the channel to receive header messages from.
	receiveFrom <-chan HeaderMessage

	// receiveErrors represents the channel to receive text messages for errors.
	receiveErrors chan mb.TextMessage

	// wg represents the wait group for synchronization.
	wg sync.WaitGroup

	// once ensures that certain operations are performed only once.
	once sync.Once

	// width represents the width of the header.
	width int
}

// GetReceiveErrorsFromChannel is a method on the Header struct.
// It returns the channel that the Header struct uses to receive TextMessage errors.
// This method allows external code to read from the receiveErrors channel
// without directly accessing the struct field.
func (h *Header) GetReceiveErrorsFromChannel() <-chan mb.TextMessage {
	return h.receiveErrors
}

// Init is a method on the Header struct. It initializes the Header with the
// given title and returns a channel to send HeaderMessage and an error.
// The title is trimmed and checked for emptiness. If the title is empty, it
// returns an ErrEmptyTitle error.
// The returned channel can be used to send HeaderMessage to the Header.
// The Header executes commands received from the channel in a separate goroutine.
// Any errors encountered during execution are sent to the receiveErrors channel.
// The method uses the sync.Once primitive to ensure that the initialization
// only happens once.
func (h *Header) Init(title string) (chan<- HeaderMessage, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return nil, &ErrEmptyTitle{}
	}

	var sendTo chan<- HeaderMessage

	h.once.Do(func() {
		h.title = title
		h.counters = make([]*Counters.UpCounter, 0)
		h.msgBuffer = new(bf.Buffer[HeaderMessage])
		sendTo, h.receiveFrom = h.msgBuffer.Init(1)
		h.receiveErrors = make(chan mb.TextMessage, 1)

		h.wg.Add(1)

		go func() {
			defer h.wg.Done()

			for msg := range h.receiveFrom {
				h.executeCommand(msg)
			}

			close(h.receiveErrors)
		}()
	})

	return sendTo, nil
}

// executeCommand is a method on the Header struct. It executes the given HeaderMessage
// and performs the corresponding actions based on the message type.
// The method supports the following message types: UpdateCurrentProcess, SetCounter,
// DesetCounter, IncrementCounter, and ReduceCounter.
// For UpdateCurrentProcess, it updates the current process with the title from the message.
// For SetCounter, it checks if the counter already exists. If it does, it sends an error
// message to the receiveErrors channel. Otherwise, it adds the counter to the counters slice.
// For DesetCounter, it removes the counter with the label from the message from the
// counters slice. If the counter does not exist, it sends an error message to the
// receiveErrors channel.
// For IncrementCounter, it increments the counter with the label from the message. If
// the counter does not exist, it sends an error message to the receiveErrors channel.
// For ReduceCounter, it reduces the counter with the label from the message. If the counter
// does not exist, it sends an error message to the receiveErrors channel.
// For unknown message types, it sends an error message to the receiveErrors channel.
// The method uses the slices package to find counters in the counters slice.
func (h *Header) executeCommand(msg HeaderMessage) {
	switch msg.MessageType {
	case UpdateCurrentProcess:
		title, err := msg.GetTitle()
		if err != nil {
			h.receiveErrors <- mb.NewTextMessage(mb.FatalText,
				fmt.Sprintf("Error getting title: %v", err),
			)

			return
		}

		h.currentProcess = title
	case SetCounter:
		counterToSet, err := msg.GetCounter()
		if err != nil {
			h.receiveErrors <- mb.NewTextMessage(mb.FatalText,
				fmt.Sprintf("Error getting counter: %v", err),
			)

			return
		}

		if slices.ContainsFunc(h.counters, func(c *Counters.UpCounter) bool {
			return c.Equal(*counterToSet)
		}) {
			ifErr, err := msg.GetIfError()
			if err != nil {
				h.receiveErrors <- mb.NewTextMessage(mb.FatalText,
					fmt.Sprintf("Error getting ifError: %v", err),
				)

				return
			}

			h.receiveErrors <- ifErr
		} else {
			h.counters = append(h.counters, counterToSet)
		}
	case UnsetCounter:
		labelToDeset, err := msg.GetLabel()
		if err != nil {
			h.receiveErrors <- mb.NewTextMessage(mb.FatalText,
				fmt.Sprintf("Error getting label: %v", err),
			)

			return
		}

		index := slices.IndexFunc(h.counters, func(c *Counters.UpCounter) bool {
			return c.ContainsLabel(labelToDeset)
		})
		if index != -1 {
			h.counters = append(h.counters[:index], h.counters[index+1:]...)
		} else {
			ifErr, err := msg.GetIfError()
			if err != nil {
				h.receiveErrors <- mb.NewTextMessage(mb.FatalText,
					fmt.Sprintf("Error getting ifError: %v", err),
				)

				return
			}

			h.receiveErrors <- ifErr
		}
	case IncrementCounter:
		labelToIncrement, err := msg.GetLabel()
		if err != nil {
			h.receiveErrors <- mb.NewTextMessage(mb.FatalText,
				fmt.Sprintf("Error getting label: %v", err),
			)

			return
		}

		index := slices.IndexFunc(h.counters, func(c *Counters.UpCounter) bool {
			return c.ContainsLabel(labelToIncrement)
		})
		if index != -1 {
			h.counters[index].Increment()
		} else {
			ifErr, err := msg.GetIfError()
			if err != nil {
				h.receiveErrors <- mb.NewTextMessage(mb.FatalText,
					fmt.Sprintf("Error getting ifError: %v", err),
				)

				return
			}

			h.receiveErrors <- ifErr
		}
	case ReduceCounter:
		labelToReduce, err := msg.GetLabel()
		if err != nil {
			h.receiveErrors <- mb.NewTextMessage(mb.FatalText,
				fmt.Sprintf("Error getting label: %v", err),
			)

			return
		}

		index := slices.IndexFunc(h.counters, func(c *Counters.UpCounter) bool {
			return c.ContainsLabel(labelToReduce)
		})
		if index != -1 {
			h.counters[index].Reduce()
		} else {
			ifErr, err := msg.GetIfError()
			if err != nil {
				h.receiveErrors <- mb.NewTextMessage(mb.FatalText,
					fmt.Sprintf("Error getting ifError: %v", err),
				)

				return
			}

			h.receiveErrors <- ifErr
		}
	default:
		h.receiveErrors <- mb.NewTextMessage(mb.FatalText,
			fmt.Sprintf("Unknown message type: %v", msg.MessageType),
		)
	}
}

// Wait is a method on the Header struct. It blocks until all goroutines
// associated with the Header have completed.
// This method is typically used to ensure that all operations started by
// the Header (such as processing HeaderMessages) have completed before
// proceeding.
// It uses the sync.WaitGroup that is part of the Header struct to track
// the completion of goroutines.
func (h *Header) Wait() {
	h.wg.Wait()
}

// Cleanup is a method on the Header struct. It cleans up the resources
// used by the Header.
// The method waits for all goroutines to finish using the WaitGroup from
// the Header struct.
// It then cleans up the message buffer in a separate goroutine and waits
// for it to finish using a local WaitGroup.
// After the message buffer is cleaned up, it sets the message buffer to nil.
// Finally, it sets the counters and receiveErrors slices to nil, effectively
// releasing the resources.
func (h *Header) Cleanup() {
	h.wg.Wait()

	var finiWg sync.WaitGroup
	defer finiWg.Wait()

	finiWg.Add(1)

	go func() {
		h.msgBuffer.Cleanup()
		finiWg.Done()
		h.msgBuffer = nil
	}()

	h.counters = nil
	h.receiveErrors = nil
}

func (h *Header) trimStringIfTooLarge(str string) string {
	if len(str) > h.width {
		newStr, _ := sext.ReplaceSuffix(str[:h.width], Hellip)

		return newStr
	}

	return str
}

// Draw is a method on the Header struct. It prints the header on the screen
// at the specified y position.
// The method takes the y position and the screen as input parameters.
// It first calculates the start position for the title to center it on the
// screen.
// It then prints the title character by character on the screen at the calculated
// position.
// If there is a current process, it prints it on the screen below the title.
// If there are any counters, it prints them on the screen below the current
// process (or the title if there is no current process).
// The method returns the updated y position and the modified screen.
//
// Parameters:
//   - y: the y position where the header should be printed
//   - screen: the screen where the header should be printed
//
// Returns:
//   - int: the updated y position
//   - tcell.Screen: the modified screen
func (h *Header) GenerateDrawTables() ([][]rune, []tcell.Style) {
	emptyLine := []rune(strings.Repeat(" ", h.width))
	y := 0 // The current y position

	table := make([][]rune, 0)
	styles := make([]tcell.Style, 0)

	// 1. Print the title
	var title string

	if len(h.title)+8 > h.width {
		title, _ = sext.ReplaceSuffix(h.title, Hellip)
	} else {
		title = h.title
	}

	title = "*** " + title + " ***"

	// Print the title
	table = append(table, emptyLine)
	styles = append(styles, style)

	startPos := (h.width - len(title)) / 2
	copy(table[y][startPos:], []rune(title))

	// 2. Empty line
	table = append(table, emptyLine)
	styles = append(styles, style)
	y++

	// 3. Print the current process (if any)
	if h.currentProcess != "" {
		table = append(table, emptyLine)
		styles = append(styles, style)
		y++

		copy(table[y][0:], []rune(h.trimStringIfTooLarge(h.currentProcess)))
	}

	// 4. Empty line
	table = append(table, emptyLine)
	styles = append(styles, style)
	y++

	// 5. Print all the counters (if any)
	if len(h.counters) != 0 {
		for _, counter := range h.counters {
			table = append(table, emptyLine)
			styles = append(styles, style)
			y++

			copy(table[y][0:], []rune(h.trimStringIfTooLarge(counter.String())))
		}
	}

	return table, styles
}

// SetWidth is a method on the Header struct. It sets the width of the Header.
// The method always returns nil, indicating that it never fails.
//
// Parameters:
//   - width: the new width of the Header
//   - height: the new height of the Header (currently ignored)
//
// Returns:
//   - error: always nil
func (h *Header) CanSetWidth(width int) bool {
	return width >= 9
}

// SetWidth is a method on the Header struct. It sets the width of the Header.
// The method always returns nil, indicating that it never fails.
//
// Parameters:
//   - width: the new width of the Header
//   - height: the new height of the Header (currently ignored)
//
// Returns:
//   - error: always nil
func (h *Header) SetWidth(width int) {
	h.width = width
}

// CanSetSize is a method on the Header struct. It sets the width and height of the Header.
// The method currently only sets the width of the Header and ignores the height parameter.
// The method always returns nil, indicating that it never fails.
//
// Parameters:
//   - width: the new width of the Header
//   - height: the new height of the Header (currently ignored)
//
// Returns:
//   - error: always nil
func (h *Header) CanSetHeight(height int) bool {
	return height >= 1
}

// SetSize is a method on the Header struct. It sets the width and height of the Header.
// The method currently only sets the width of the Header and ignores the height parameter.
// The method always returns nil, indicating that it never fails.
//
// Parameters:
//   - width: the new width of the Header
//   - height: the new height of the Header (currently ignored)
//
// Returns:
//   - error: always nil
func (h *Header) SetHeight(height int) {
	// Do nothing
}

func (h *Header) GetCurrentHeight() int {
	y := 3 + len(h.counters) // The current y position

	if h.currentProcess != "" {
		return y + 1
	} else {
		return y
	}
}
