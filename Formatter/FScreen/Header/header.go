package Header

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"sync"

	bf "github.com/PlayerR9/MyGoLib/CustomData/Buffer"
	"github.com/PlayerR9/MyGoLib/CustomData/Counters"
	mb "github.com/PlayerR9/MyGoLib/Formatter/FScreen/MessageBox"
	"github.com/gdamore/tcell"
)

type Header struct {
	title          string
	currentProcess string
	counters       []*Counters.UpCounter

	messageChannel bf.Buffer[HeaderMessage]
	ErrorChannel   chan mb.TextMessage
	wg             sync.WaitGroup
	once           sync.Once
}

func NewHeader(title string) (*Header, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return nil, errors.New("title cannot be empty")
	}

	return &Header{
		title:          title,
		currentProcess: "",
		counters:       make([]*Counters.UpCounter, 0),
		messageChannel: bf.NewBuffer[HeaderMessage](),
		ErrorChannel:   make(chan mb.TextMessage),
	}, nil
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

func (h *Header) Run() {
	h.once.Do(func() {
		go h.messageChannel.Run()

		h.wg.Add(1)

		for {
			msg, ok := h.messageChannel.Get()
			if !ok {
				h.wg.Done()

				return
			}

			switch msg.GetType() {
			case UpdateCurrentProcess:
				h.currentProcess = msg.GetTitle()
			case SetCounter:
				counterToSet := msg.GetCounter()

				if slices.ContainsFunc(h.counters, func(c *Counters.UpCounter) bool {
					return c.Equal(*counterToSet)
				}) {
					h.ErrorChannel <- msg.GetIfError()
				} else {
					h.counters = append(h.counters, counterToSet)
				}
			case DesetCounter:
				labelToDeset := msg.GetLabel()

				index := slices.IndexFunc(h.counters, func(c *Counters.UpCounter) bool {
					return c.ContainsLabel(labelToDeset)
				})
				if index != -1 {
					h.counters = append(h.counters[:index], h.counters[index+1:]...)
				} else {
					h.ErrorChannel <- msg.GetIfError()
				}
			case IncrementCounter:
				labelToIncrement := msg.GetLabel()

				index := slices.IndexFunc(h.counters, func(c *Counters.UpCounter) bool {
					return c.ContainsLabel(labelToIncrement)
				})
				if index != -1 {
					h.counters[index].Increment()
				} else {
					h.ErrorChannel <- msg.GetIfError()
				}
			case ReduceCounter:
				labelToReduce := msg.GetLabel()

				index := slices.IndexFunc(h.counters, func(c *Counters.UpCounter) bool {
					return c.ContainsLabel(labelToReduce)
				})
				if index != -1 {
					h.counters[index].Reduce()
				} else {
					h.ErrorChannel <- msg.GetIfError()
				}
			default:
				h.ErrorChannel <- mb.NewTextMessage(mb.FatalText,
					fmt.Sprintf("Unknown message type: %v", msg.GetType()),
				)
			}
		}
	})
}

func (h *Header) Fini() {
	h.messageChannel.Fini()
	close(h.ErrorChannel)

	h.wg.Wait()

	// Close and free resources
	h.counters = nil
	h.ErrorChannel = nil
}

func (h *Header) SendMessages(message HeaderMessage, optionalMessages ...HeaderMessage) {
	h.messageChannel.SendMessages(
		message,
		optionalMessages...,
	)
}
