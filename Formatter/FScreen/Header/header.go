package Header

import (
	"fmt"
	"slices"
	"strings"
	"sync"

	"github.com/PlayerR9/MyGoLib/CustomData/Counters"
	bf "github.com/PlayerR9/MyGoLib/CustomData/Rerouting"
	mb "github.com/PlayerR9/MyGoLib/Formatter/FScreen/MessageBox"
	"github.com/gdamore/tcell"
)

type Header struct {
	title          string
	currentProcess string
	counters       []*Counters.UpCounter

	msgBuffer   *bf.Buffer[bf.Messager]
	receiveFrom <-chan bf.Messager

	receiveErrors chan mb.TextMessage
	wg            sync.WaitGroup
	once          sync.Once
}

func (h *Header) Init(title string) (bf.SendChannel, <-chan mb.TextMessage) {
	var sendTo chan<- bf.Messager

	h.once.Do(func() {
		title = strings.TrimSpace(title)
		if title == "" {
			panic("title cannot be empty")
		}

		h.title = title
		h.counters = make([]*Counters.UpCounter, 0)
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

	return bf.NewSendChannel(sendTo, 1), h.receiveErrors
}

func (h *Header) executeCommand(msg bf.Messager) {
	command, ok := msg.(HeaderMessage)
	if !ok {
		h.receiveErrors <- mb.NewTextMessage(mb.FatalText,
			fmt.Sprintf("Unknown message type: %T", msg),
		)
		return
	}

	switch command.GetType() {
	case UpdateCurrentProcess:
		h.currentProcess = command.GetTitle()
	case SetCounter:
		counterToSet := command.GetCounter()

		if slices.ContainsFunc(h.counters, func(c *Counters.UpCounter) bool {
			return c.Equal(*counterToSet)
		}) {
			h.receiveErrors <- command.GetIfError()
		} else {
			h.counters = append(h.counters, counterToSet)
		}
	case DesetCounter:
		labelToDeset := command.GetLabel()

		index := slices.IndexFunc(h.counters, func(c *Counters.UpCounter) bool {
			return c.ContainsLabel(labelToDeset)
		})
		if index != -1 {
			h.counters = append(h.counters[:index], h.counters[index+1:]...)
		} else {
			h.receiveErrors <- command.GetIfError()
		}
	case IncrementCounter:
		labelToIncrement := command.GetLabel()

		index := slices.IndexFunc(h.counters, func(c *Counters.UpCounter) bool {
			return c.ContainsLabel(labelToIncrement)
		})
		if index != -1 {
			h.counters[index].Increment()
		} else {
			h.receiveErrors <- command.GetIfError()
		}
	case ReduceCounter:
		labelToReduce := command.GetLabel()

		index := slices.IndexFunc(h.counters, func(c *Counters.UpCounter) bool {
			return c.ContainsLabel(labelToReduce)
		})
		if index != -1 {
			h.counters[index].Reduce()
		} else {
			h.receiveErrors <- command.GetIfError()
		}
	default:
		h.receiveErrors <- mb.NewTextMessage(mb.FatalText,
			fmt.Sprintf("Unknown message type: %v", command.GetType()),
		)
	}
}

func (h *Header) Wait() {
	h.wg.Wait()
}

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

func (h *Header) SetScreen(y, width int, style tcell.Style, screen tcell.Screen) (int, tcell.Screen) {
	// Print the title
	startPos := (width - len(h.title)) / 2
	for x, char := range h.title {
		screen.SetContent(startPos+x, y, char, nil, style)
	}

	// Print the current process (if interface{})
	if h.currentProcess != "" {
		y += 2

		for x, char := range h.currentProcess {
			screen.SetContent(x, y, char, nil, style)
		}
	}

	// Print all the counters (if interface{})
	if len(h.counters) != 0 {
		y++

		for _, counter := range h.counters {
			y++

			for x, char := range counter.String() {
				screen.SetContent(x, y, char, nil, style)
			}
		}
	}

	return y, screen
}
