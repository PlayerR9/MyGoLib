package FScreen

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	bf "github.com/PlayerR9/MyGoLib/CustomData/Rerouting"
	h "github.com/PlayerR9/MyGoLib/Formatter/FScreen/Header"
	mb "github.com/PlayerR9/MyGoLib/Formatter/FScreen/MessageBox"
)

var FScreen *bf.Hub

func Init(title string, width, height int) error {
	if strings.TrimSpace(title) == "" {
		return errors.New("title cannot be empty")
	} else if width <= 5 || height <= 2 {
		return errors.New("width and height must be greater than 5 and 2 respectively")
	}

	FScreen = bf.NewHub()

	header := new(h.Header)
	toSendHeader, receiveHeaderErrors := header.Init(title)

	box := new(mb.MessageBox)
	toSendBox := box.Init(width, height)

	FScreen.AddConnection(receiveHeaderErrors, toSendBox)
	FScreen.AddEntryPoint(toSendHeader)
	FScreen.SetInexistentEntryPoint(mb.GetSendToChannel())

	return nil
}

func (fs *FScreen) routingMessages() {
	defer fs.wg.Done()

	for msg := range fs.receiveFrom {
		correctChannel := msg.Channel()

		switch x := msg.(type) {
		case mb.TextMessage:
			fs.sendToMessageBox <- x
		case h.HeaderMessage:
			fs.sendToHeader <- x
		default:
			fs.sendToMessageBox <- mb.NewTextMessage(mb.FatalText,
				fmt.Sprintf("Unknown message type: %T", x),
			)
		}
	}

	errorWg.Wait()
}

func (fs *FScreen) Cleanup() {
	fs.wg.Wait()

	var finiWg sync.WaitGroup
	defer finiWg.Wait()

	finiWg.Add(3)

	go func() {
		screen.Fini()
		finiWg.Done()
		screen = nil
	}()

	go func() {
		fs.msgBuffer.Cleanup()
		finiWg.Done()
		fs.msgBuffer = nil
	}()

	go func() {
		close(fs.sendToMessageBox)
		messageBox.Fini()
		messageBox.Cleanup()
		finiWg.Done()
		messageBox = nil
	}()
}
