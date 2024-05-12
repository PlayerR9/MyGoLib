package Runner

import (
	"sync"

	ers "github.com/PlayerR9/MyGoLib/Units/Errors"
)

// ExecuteMsgFunc is a function type that executes a message.
//
// Parameters:
//   - msg: The message to execute.
//
// Returns:
//   - error: An error if the message cannot be executed.
type ExecuteMsgFunc[T any] func(msg T) error

// HandlerSend is a struct that represents a Go routine handler.
// It is used to handle the result of a Go routine.
type HandlerSend[T any] struct {
	// wg is a WaitGroup that is used to wait for the Go routine to finish.
	wg sync.WaitGroup

	// errChan is the error status of the Go routine.
	errChan chan error

	// routine is the Go routine that is run by the handler.
	routine ExecuteMsgFunc[T]

	// sendChan is the channel to send messages to the Go routine.
	sendChan chan T
}

// run is a private method of HandlerSend that is runned by the Go routine.
//
// Behaviors:
//   - Use ers.ErrNoError to exit the Go routine as nil is used to signal
//     that the function has finished successfully but the Go routine is still running.
func (h *HandlerSend[T]) run() {
	defer h.wg.Done()

	defer func() {
		if r := recover(); r != nil {
			h.errChan <- ers.NewErrPanic(r)
		}
	}()

	for msg := range h.sendChan {
		err := h.routine(msg)
		if err == nil {
			continue
		}

		if ers.As[*ers.ErrNoError](err) {
			return
		}

		h.errChan <- err
	}
}

// Start is a method of HandlerSend that starts the Go routine.
//
// Returns:
//   - error: An error  if the Go routine cannot be started.
//
// Errors:
//   - *ErrAlreadyRunning: The Go routine is already running.
//   - *ers.ErrInvalidType: The status of the Go routine is invalid.
//
// Behaviors:
//   - If the Go routine is stopped, it will be started.
//   - If the Go routine is running, the error *ErrAlreadyRunning is returned.
//   - If the Go routine is not initialized, the error *ErrNilValue is returned.
func (h *HandlerSend[T]) Start() error {
	if h.sendChan != nil {
		return NewErrAlreadyRunning()
	}

	h.errChan = make(chan error)
	h.sendChan = make(chan T)

	h.wg.Add(1)

	go h.run()

	return nil
}

// Wait is a method of HandlerSend that waits for the Go routine to finish.
//
// Behaviors:
//   - If the Go routine is not running, this method does nothing.
func (h *HandlerSend[T]) Wait() {
	if h.sendChan == nil {
		return
	}

	h.wg.Wait()
}

// Close closes the runner.
//
// Behaviors:
//   - If the Go routine is not running, this method does nothing.
func (h *HandlerSend[T]) Close() {
	if h.sendChan == nil {
		return
	}

	close(h.sendChan)

	h.wg.Wait()

	// Clean up
	h.sendChan = nil

	close(h.errChan)
	h.errChan = nil
}

// IsRunning returns whether the runner is running.
//
// Returns:
//   - bool: True if the runner is running, false otherwise.
func (h *HandlerSend[T]) IsRunning() bool {
	return h.sendChan != nil
}

// NewHandlerSend creates a new HandlerSend.
//
// Parameters:
//   - routine: The Go routine to run.
//
// Returns:
//   - *HandlerSend: A pointer to the HandlerSend that handles the result of the Go routine.
//
// Behaviors:
//   - The Go routine is not started automatically.
//   - In routine, use *ers.ErrNoError to exit the Go routine as nil is used to signal
//     that the function has finished successfully but the Go routine is still running.
func NewHandlerSend[T any](routine ExecuteMsgFunc[T]) *HandlerSend[T] {
	return &HandlerSend[T]{
		routine: routine,
	}
}

// GetErrChannel is a method of HandlerSend that returns the error status of the Go routine.
//
// Returns:
//   - error: The error status of the Go routine.
func (h *HandlerSend[T]) GetErrChannel() <-chan error {
	return h.errChan
}

// Send is a method of HandlerSend that sends a message to the Go routine.
// If the Go routine is not running, false is returned.
//
// Parameters:
//   - msg: The message to send.
//
// Returns:
//   - bool: True if the message is sent, false otherwise.
func (h *HandlerSend[T]) Send(msg T) bool {
	if h.sendChan == nil {
		return false
	}

	h.sendChan <- msg

	return true
}
