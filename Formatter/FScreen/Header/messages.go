package Header

import (
	"fmt"
	"strings"

	"github.com/PlayerR9/MyGoLib/CustomData/Counters"

	mb "github.com/PlayerR9/MyGoLib/Formatter/FScreen/MessageBox"
)

// HeaderMessageType represents the type of a message header.
// It is an integer type and can take one of the following constants:
// - UpdateCurrentProcess
// - SetCounter
// - DesetCounter
// - IncrementCounter
// - ReduceCounter
type HeaderMessageType int

// These are the constants for HeaderMessageType.
// UpdateCurrentProcess is used when the current process needs to be updated.
// SetCounter is used when the counter needs to be set to a specific value.
// UnsetCounter is used when the counter needs to be unset or reset.
// IncrementCounter is used when the counter needs to be incremented.
// ReduceCounter is used when the counter needs to be reduced.
const (
	UpdateCurrentProcess HeaderMessageType = iota
	SetCounter
	UnsetCounter
	IncrementCounter
	ReduceCounter
)

// String returns the string representation of the HeaderMessageType.
// It returns one of the following values:
// - "UpdateCurrentProcess"
// - "SetCounter"
// - "UnsetCounter"
// - "IncrementCounter"
// - "ReduceCounter"
// This method is useful when the string representation of the enum value is
// needed, for example, in debugging or logging.
func (enum HeaderMessageType) String() string {
	return [...]string{
		"UpdateCurrentProcess",
		"SetCounter",
		"UnsetCounter",
		"IncrementCounter",
		"ReduceCounter",
	}[enum]
}

// HeaderMessage represents a message that includes a specific header
// type and associated data.
// It is a struct that contains two fields: MessageType and Data.
type HeaderMessage struct {
	// The type of the message.
	MessageType HeaderMessageType

	// The data associated with the message.
	Data interface{}
}

// NewUpdateCurrentProcessMessage is a function that creates a new
// HeaderMessage for setting the current process.
// It takes a string argument, title, which represents the title
// of the process to be set.
// The function first trims any leading or trailing whitespace
// from the title.
// If the trimmed title is empty, the function returns a HeaderMessage
// with messageType set to UpdateCurrentProcess and data set to an
// error message.
// If the trimmed title is not empty, the function returns a
// HeaderMessage with messageType set to UpdateCurrentProcess and
// data set to the provided title.
// The returned HeaderMessage can be used to set the current process
// in the application.
func NewUpdateCurrentProcessMessage(title string) HeaderMessage {
	title = strings.TrimSpace(title)
	if title == "" {
		return HeaderMessage{
			MessageType: UpdateCurrentProcess,
			Data: mb.NewTextMessage(mb.ErrorText,
				"Trying to set current process with empty title;",
				"skipping command",
			),
		}
	}

	return HeaderMessage{
		MessageType: UpdateCurrentProcess,
		Data:        title,
	}
}

// NewEmptyCurrentProcessMessage is a function that creates a new
// HeaderMessage for unsetting the current process.
// It does not take any arguments.
// The function returns a HeaderMessage with messageType set to
// UpdateCurrentProcess and data set to an empty string.
// The returned HeaderMessage can be used to unset the current
// process in the application.
func NewEmptyCurrentProcessMessage() HeaderMessage {
	title := ""

	return HeaderMessage{
		MessageType: UpdateCurrentProcess,
		Data:        title,
	}
}

// GetTitle attempts to retrieve the title from the HeaderMessage.
// It returns the title as a string if it exists and is a string.
// If the title does not exist or is not a string, it returns an error.
func (hm *HeaderMessage) GetTitle() (string, error) {
	title, ok := hm.Data.(string)
	if !ok {
		return "", &ErrDataNotString{}
	}
	return title, nil
}

// CounterData represents a structure that holds information about a counter.
// It includes a label for the counter, the counter object itself, and a
// message to be displayed in case of an error.
type CounterData struct {
	// label is the description of the counter data.
	label string

	// counter is the actual counter object.
	counter *Counters.UpCounter

	// ifError is the error message to be displayed if there is an error.
	ifError mb.TextMessage
}

// GetCounter attempts to retrieve the counter from the HeaderMessage.
// It returns the counter as a *Counters.UpCounter if it exists and is
// of type counterData.
// If the counter does not exist or is not of type counterData, it
// returns an error of type ErrDataNotCounter.
func (hm *HeaderMessage) GetCounter() (*Counters.UpCounter, error) {
	counter, ok := hm.Data.(CounterData)
	if !ok {
		return nil, &ErrDataNotCounter{}
	}

	return counter.counter, nil
}

// GetLabel attempts to retrieve the label from the HeaderMessage.
// It returns the label as a string if it exists and is of type counterData.
// If the label does not exist or is not of type counterData, it
// returns an error of type ErrDataNotCounter.
func (hm *HeaderMessage) GetLabel() (string, error) {
	counter, ok := hm.Data.(CounterData)
	if !ok {
		return "", &ErrDataNotCounter{}
	}

	return counter.label, nil
}

// GetIfError attempts to retrieve the ifError field from the HeaderMessage.
// It returns the ifError as a mb.TextMessage if it exists and is of type counterData.
// If the ifError does not exist or is not of type counterData, it returns an
// error of type ErrDataNotCounter.
func (hm *HeaderMessage) GetIfError() (mb.TextMessage, error) {
	counter, ok := hm.Data.(CounterData)
	if !ok {
		return mb.TextMessage{}, &ErrDataNotCounter{}
	}

	return counter.ifError, nil
}

// NewSetCounterMessage creates a new HeaderMessage of type SetCounter.
// It takes a label string and totalElements int as parameters.
// The label is trimmed of leading and trailing whitespace.
// If the label is empty, it returns a HeaderMessage with an error
// message indicating an empty label.
// If the totalElements is negative, it returns a HeaderMessage with an
// error message indicating a negative totalElements.
// Otherwise, it creates a new UpCounter with the given label and totalElements,
// and returns a HeaderMessage with the counter data.
func NewSetCounterMessage(label string, totalElements int) HeaderMessage {
	label = strings.TrimSpace(label)
	if label == "" {
		return HeaderMessage{
			MessageType: SetCounter,
			Data: mb.NewTextMessage(mb.WarningText,
				"Trying to set counter with empty label;",
				"skipping command",
			),
		}
	} else if totalElements < 0 {
		return HeaderMessage{
			MessageType: SetCounter,
			Data: mb.NewTextMessage(mb.WarningText,
				"Trying to set counter with negative 'totalElements';",
				"skipping command",
			),
		}
	}

	counter := Counters.NewUpCounter(label, totalElements)

	return HeaderMessage{
		MessageType: SetCounter,
		Data: CounterData{
			label:   label,
			counter: &counter,
			ifError: mb.NewTextMessage(mb.WarningText,
				fmt.Sprintf("Counter \"%s\" already exists;", label),
				"cannot set it again",
			),
		},
	}
}

// NewUnsetCounterMessage creates a new HeaderMessage of type UnsetCounter.
// It takes a label string as a parameter.
// The label is trimmed of leading and trailing whitespace.
// If the label is empty, it returns a HeaderMessage with an error message
// indicating an empty label.
// Otherwise, it returns a HeaderMessage with the provided label and a warning
// message if the counter does not exist.
func NewUnsetCounterMessage(label string) HeaderMessage {
	label = strings.TrimSpace(label)
	if label == "" {
		return HeaderMessage{
			MessageType: UnsetCounter,
			Data: mb.NewTextMessage(mb.WarningText,
				"Trying to unset counter with empty label;",
				"skipping command",
			),
		}
	}

	return HeaderMessage{
		MessageType: UnsetCounter,
		Data: CounterData{
			label:   label,
			counter: nil,
			ifError: mb.NewTextMessage(mb.WarningText,
				fmt.Sprintf("Counter \"%s\" does not exist;", label),
				"already unset",
			),
		},
	}
}

// NewIncrementCounterMessage creates a new HeaderMessage of type IncrementCounter.
// It takes a label string as a parameter.
// The label is trimmed of leading and trailing whitespace.
// If the label is empty, it returns a HeaderMessage with an error message indicating
// an empty label.
// Otherwise, it returns a HeaderMessage with the provided label and a warning message
// if the counter does not exist.
func NewIncrementCounterMessage(label string) HeaderMessage {
	label = strings.TrimSpace(label)
	if label == "" {
		return HeaderMessage{
			MessageType: IncrementCounter,
			Data: mb.NewTextMessage(mb.WarningText,
				"Trying to increment counter with empty label;",
				"skipping command",
			),
		}
	}

	return HeaderMessage{
		MessageType: IncrementCounter,
		Data: CounterData{
			label:   label,
			counter: nil,
			ifError: mb.NewTextMessage(mb.WarningText,
				fmt.Sprintf("Counter \"%s\" does not exist;", label),
				"skipping increment",
			),
		},
	}
}

// NewReduceCounterMessage creates a new HeaderMessage of type ReduceCounter.
// It takes a label string as a parameter.
// The label is trimmed of leading and trailing whitespace.
// If the label is empty, it returns a HeaderMessage with an error message
// indicating an empty label.
// Otherwise, it returns a HeaderMessage with the provided label and a warning
// message if the counter does not exist.
func NewReduceCounterMessage(label string) HeaderMessage {
	label = strings.TrimSpace(label)
	if label == "" {
		return HeaderMessage{
			MessageType: ReduceCounter,
			Data: mb.NewTextMessage(mb.WarningText,
				"Trying to reduce counter with empty label;",
				"skipping command",
			),
		}
	}

	return HeaderMessage{
		MessageType: ReduceCounter,
		Data: CounterData{
			label:   label,
			counter: nil,
			ifError: mb.NewTextMessage(mb.WarningText,
				fmt.Sprintf("Counter \"%s\" does not exist;", label),
				"skipping reduce",
			),
		},
	}
}
