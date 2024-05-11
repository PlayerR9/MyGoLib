package Keyboard

import (
	"fmt"
	_ "image/png"
	"sync"

	sfb "github.com/PlayerR9/MyGoLib/Safe/Buffer"
	"github.com/eiannone/keyboard"
)

// Keyboard handles keyboard input using the eiannone/keyboard package.
type Keyboard struct {
	// buffer is a safe buffer for keyboard.Key values.
	buffer *sfb.Buffer[keyboard.Key]

	// sendKeyChan is the send channel for keyboard.Key values.
	sendKeyChan chan<- keyboard.Key

	// receiveKeyChan is the receive channel for keyboard.Key values.
	receiveKeyChan <-chan keyboard.Key

	// errChan is the error channel for the Keyboard.
	errChan chan error

	// closeChan is the close channel for the Keyboard.
	closeChan chan struct{}

	// wg is the wait group for the Keyboard.
	wg sync.WaitGroup
}

// NewKeyboard creates a new Keyboard.
//
// Returns:
//   - *Keyboard: The new Keyboard.
func NewKeyboard() *Keyboard {
	k := &Keyboard{
		errChan:   make(chan error),
		closeChan: make(chan struct{}),
	}

	buffer, err := sfb.NewBuffer[keyboard.Key](1)
	if err != nil {
		panic(err)
	}

	k.buffer = buffer

	k.sendKeyChan = buffer.GetSendChannel()
	k.receiveKeyChan = buffer.GetReceiveChannel()

	return k
}

// GetErrorChannel returns the error channel for the Keyboard.
//
// Returns:
//   - <-chan error: The error channel.
func (k *Keyboard) GetErrorChannel() <-chan error {
	return k.errChan
}

// GetKeyChannel returns the key channel for the Keyboard.
//
// Returns:
//   - <-chan keyboard.Key: The key channel.
func (k *Keyboard) GetKeyChannel() <-chan keyboard.Key {
	return k.receiveKeyChan
}

// Close closes the Keyboard.
//
// Returns:
//   - error: An error if the Keyboard could not be closed.
func (k *Keyboard) Close() error {
	if k.sendKeyChan == nil {
		return fmt.Errorf("keyboard already closed")
	}

	close(k.sendKeyChan)

	k.buffer.Wait()

	close(k.closeChan)

	k.wg.Wait()

	close(k.errChan)

	// Clean up

	k.buffer = nil
	k.sendKeyChan = nil
	k.receiveKeyChan = nil

	err := keyboard.Close()
	if err != nil {
		return err
	}

	return nil
}

// Start starts the Keyboard.
//
// Returns:
//   - error: An error if the Keyboard could not be started.
func (k *Keyboard) Start() error {
	err := keyboard.Open()
	if err != nil {
		return err
	}

	k.wg.Add(1)

	go k.keyListener()

	return nil
}

// Wait waits for the Keyboard to finish.
// It may cause a deadlock if the Keyboard is not closed.
func (k *Keyboard) Wait() {
	k.wg.Wait()
}

// keyListener is an helper function that listens for keyboard input.
func (k *Keyboard) keyListener() {
	defer k.wg.Done()

	for {
		select {
		case <-k.closeChan:
			return
		default:
			_, key, err := keyboard.GetKey()
			if err != nil {
				k.errChan <- err
			} else {
				k.sendKeyChan <- key
			}
		}
	}
}
