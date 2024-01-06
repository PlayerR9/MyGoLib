// git tag v0.1.45

package FScreen

import (
	"fmt"
	"sync"

	buffer "github.com/PlayerR9/MyGoLib/CustomData/Buffer"
	h "github.com/PlayerR9/MyGoLib/Formatter/FScreen/Header"
	mb "github.com/PlayerR9/MyGoLib/Formatter/FScreen/MessageBox"
)

var (
	header     *h.Header      = nil
	messageBox *mb.MessageBox = nil
)

type FScreen struct {
	msgBuffer   *buffer.Buffer[interface{}]
	receiveFrom <-chan interface{}
	once        sync.Once
	wg          sync.WaitGroup

	sendToMessageBox        chan<- mb.TextMessage
	receiveMessageBoxErrors <-chan mb.TextMessage

	sendToHeader        chan<- h.HeaderMessage
	receiveHeaderErrors <-chan mb.TextMessage
}

func (fs *FScreen) Init(title string, width, height int) (chan<- interface{}, error) {
	var err error

	messageBox = new(mb.MessageBox)
	fs.sendToMessageBox, err = messageBox.Init(width, height)
	if err != nil {
		return nil, err
	}

	header = new(h.Header)
	fs.sendToHeader, err = header.Init(title)
	if err != nil {
		return nil, err
	}

	var sendTo chan<- interface{}

	fs.once.Do(func() {
		fs.receiveMessageBoxErrors = messageBox.GetReceiveErrorsFromChannel()
		fs.receiveHeaderErrors = header.GetReceiveErrorsFromChannel()
		fs.msgBuffer = new(buffer.Buffer[interface{}])
		sendTo, fs.receiveFrom = fs.msgBuffer.Init(1)

		fs.wg.Add(1)

		go fs.routingMessages()
	})

	return sendTo, nil
}

func (fs *FScreen) routingMessages() {
	defer fs.wg.Done()

	// Start the rerouting channel listeners
	var reroutingWg sync.WaitGroup

	reroutingWg.Add(2)

	go func() {
		defer reroutingWg.Done()

		for msg := range fs.receiveHeaderErrors {
			fs.sendToMessageBox <- msg
		}
	}()

	go func() {
		defer reroutingWg.Done()

		for msg := range fs.receiveMessageBoxErrors {
			fs.sendToMessageBox <- msg
		}
	}()

	for msg := range fs.receiveFrom {
		switch x := msg.(type) {
		case mb.TextMessage:
			fs.sendToMessageBox <- x
		case h.HeaderMessage:
			switch t := x.Data.(type) {
			case mb.TextMessage:
				fs.sendToMessageBox <- t
			case string, h.CounterData:
				fs.sendToHeader <- x
			default:
				fs.sendToMessageBox <- mb.NewTextMessage(mb.ErrorText,
					fmt.Sprintf("Unknown header message type: %T", x),
				)
			}
		default:
			fs.sendToMessageBox <- mb.NewTextMessage(mb.FatalText,
				fmt.Sprintf("Unknown message type: %T", x),
			)
		}
	}

	reroutingWg.Wait()
}

func (fs *FScreen) Cleanup() {
	fs.wg.Wait()

	var cleanWg sync.WaitGroup
	defer cleanWg.Wait()

	cleanWg.Add(3)

	go func() {
		fs.msgBuffer.Cleanup()
		cleanWg.Done()
		fs.msgBuffer = nil
		fs.receiveFrom = nil
	}()

	go func() {
		close(fs.sendToMessageBox)
		messageBox.Cleanup()
		cleanWg.Done()
		messageBox = nil
	}()

	go func() {
		close(fs.sendToHeader)
		header.Cleanup()
		cleanWg.Done()
		header = nil
		fs.sendToHeader = nil
		fs.receiveHeaderErrors = nil
	}()
}

func (fs *FScreen) Wait() {
	fs.wg.Wait()
}
