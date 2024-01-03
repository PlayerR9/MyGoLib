package Header

import (
	"fmt"
	"strings"

	"github.com/PlayerR9/MyGoLib/CustomData/Counters"

	bf "github.com/PlayerR9/MyGoLib/CustomData/Rerouting"
	mb "github.com/PlayerR9/MyGoLib/Formatter/FScreen/MessageBox"
)

type HeaderMessageType int

const (
	UpdateCurrentProcess HeaderMessageType = iota
	SetCounter
	DesetCounter
	IncrementCounter
	ReduceCounter
)

func (enum HeaderMessageType) String() string {
	return [...]string{
		"UpdateCurrentProcess",
		"SetCounter",
		"DesetCounter",
		"IncrementCounter",
		"ReduceCounter",
	}[enum]
}

type HeaderMessage struct {
	messageType HeaderMessageType
	data        interface{}
}

func (hm HeaderMessage) Channel() bf.SendChannel {
	return bf.NewSendChannel(make(chan<- bf.Messager), 1)
}

func (hm HeaderMessage) ParseInexistentEntryPoint() bf.Messager {
	panic("implement me")
}

func (cm HeaderMessage) GetType() HeaderMessageType {
	return cm.messageType
}

func NewSetCurrentProcessMessage(title string) HeaderMessage {
	title = strings.TrimSpace(title)
	if title == "" {
		return HeaderMessage{
			messageType: UpdateCurrentProcess,
			data: mb.NewTextMessage(mb.ErrorText,
				"Trying to set current process with empty title;",
				"skipping command",
			),
		}
	}

	return HeaderMessage{
		messageType: UpdateCurrentProcess,
		data:        title,
	}
}

func NewDesetCurrentProcessMessage() HeaderMessage {
	title := ""

	return HeaderMessage{
		messageType: UpdateCurrentProcess,
		data:        title,
	}
}

func (hm *HeaderMessage) GetTitle() string {
	return hm.data.(string)
}

type counterData struct {
	label   string
	counter *Counters.UpCounter
	ifError mb.TextMessage
}

func (hm *HeaderMessage) GetCounter() *Counters.UpCounter {
	return hm.data.(counterData).counter
}

func (hm *HeaderMessage) GetLabel() string {
	return hm.data.(counterData).label
}

func (hm *HeaderMessage) GetIfError() mb.TextMessage {
	return hm.data.(counterData).ifError
}

func NewSetCounterMessage(label string, totalElements int) HeaderMessage {
	label = strings.TrimSpace(label)
	if label == "" {
		return HeaderMessage{
			messageType: SetCounter,
			data: mb.NewTextMessage(mb.ErrorText,
				"Trying to set counter with empty label;",
				"skipping command",
			),
		}
	} else if totalElements < 0 {
		return HeaderMessage{
			messageType: SetCounter,
			data: mb.NewTextMessage(mb.ErrorText,
				"Trying to set counter with negative 'totalElements';",
				"skipping command",
			),
		}
	}

	counter := Counters.NewUpCounter(label, totalElements)

	return HeaderMessage{
		messageType: SetCounter,
		data: counterData{
			label:   label,
			counter: &counter,
			ifError: mb.NewTextMessage(mb.WarningText,
				fmt.Sprintf("Counter \"%s\" already exists;", label),
				"cannot set it again",
			),
		},
	}
}

func NewDesetCounterMessage(label string) HeaderMessage {
	label = strings.TrimSpace(label)
	if label == "" {
		return HeaderMessage{
			messageType: DesetCounter,
			data: mb.NewTextMessage(mb.ErrorText,
				"Trying to deset counter with empty label;",
				"skipping command",
			),
		}
	}

	return HeaderMessage{
		messageType: DesetCounter,
		data: counterData{
			label:   label,
			counter: nil,
			ifError: mb.NewTextMessage(mb.WarningText,
				fmt.Sprintf("Counter \"%s\" does not exist;", label),
				"already deset",
			),
		},
	}
}

func NewIncrementCounterMessage(label string) HeaderMessage {
	label = strings.TrimSpace(label)
	if label == "" {
		return HeaderMessage{
			messageType: IncrementCounter,
			data: mb.NewTextMessage(mb.ErrorText,
				"Trying to increment counter with empty label;",
				"skipping command",
			),
		}
	}

	return HeaderMessage{
		messageType: IncrementCounter,
		data: counterData{
			label:   label,
			counter: nil,
			ifError: mb.NewTextMessage(mb.WarningText,
				fmt.Sprintf("Counter \"%s\" does not exist;", label),
				"skipping increment",
			),
		},
	}
}

func NewReduceCounterMessage(label string) HeaderMessage {
	label = strings.TrimSpace(label)
	if label == "" {
		return HeaderMessage{
			messageType: ReduceCounter,
			data: mb.NewTextMessage(mb.ErrorText,
				"Trying to reduce counter with empty label;",
				"skipping command",
			),
		}
	}

	return HeaderMessage{
		messageType: ReduceCounter,
		data: counterData{
			label:   label,
			counter: nil,
			ifError: mb.NewTextMessage(mb.WarningText,
				fmt.Sprintf("Counter \"%s\" does not exist;", label),
				"skipping reduce",
			),
		},
	}
}
