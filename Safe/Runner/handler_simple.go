package Runner

import (
	"sync"

	uc "github.com/PlayerR9/MyGoLib/Units/Common"
	ers "github.com/PlayerR9/MyGoLib/Units/errors"
)

// HandlerSimple is a struct that represents a Go routine handler.
// It is used to handle the result of a Go routine.
type HandlerSimple struct {
	// wg is a WaitGroup that is used to wait for the Go routine to finish.
	wg sync.WaitGroup

	// errChan is the error status of the Go routine.
	errChan chan error

	// routine is the Go routine that is run by the handler.
	routine uc.MainFunc

	// stopChan is a channel that is used to stop the Go routine.
	stopChan chan struct{}
}

// run is a private method of HandlerSimple that is runned by the Go routine.
//
// Behaviors:
//   - Use ers.ErrNoError to exit the Go routine as nil is used to signal
//     that the function has finished successfully but the Go routine is still running.
func (h *HandlerSimple) run() {
	defer h.wg.Done()

	defer func() {
		if r := recover(); r != nil {
			h.errChan <- ers.NewErrPanic(r)
		}
	}()

	for {
		select {
		case <-h.stopChan:
			return
		default:
			err := h.routine()
			if err != nil {
				if ers.As[*ers.ErrNoError](err) {
					return
				}

				h.errChan <- err
			}
		}
	}
}

// Start is a method of HandlerSimple that starts the Go routine.
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
func (h *HandlerSimple) Start() error {
	if h.stopChan != nil {
		return NewErrAlreadyRunning()
	}

	h.errChan = make(chan error)
	h.stopChan = make(chan struct{})

	h.wg.Add(1)

	go h.run()

	return nil
}

// Wait is a method of HandlerSimple that waits for the Go routine to finish.
//
// Behaviors:
//   - If the Go routine is not running, this method does nothing.
func (h *HandlerSimple) Wait() {
	if h.stopChan == nil {
		return
	}

	h.wg.Wait()
}

// Close closes the runner.
//
// Behaviors:
//   - If the Go routine is not running, this method does nothing.
func (h *HandlerSimple) Close() {
	if h.stopChan == nil {
		return
	}

	close(h.stopChan)

	h.wg.Wait()

	// Clean up
	close(h.errChan)
	h.errChan = nil

	h.stopChan = nil
}

// IsRunning returns whether the runner is running.
//
// Returns:
//   - bool: True if the runner is running, false otherwise.
func (h *HandlerSimple) IsRunning() bool {
	return h.stopChan != nil
}

// NewHandlerSimple creates a new HandlerSimple.
//
// Parameters:
//   - routine: The Go routine to run.
//
// Returns:
//   - *HandlerSimple: A pointer to the HandlerSimple that handles the result of the Go routine.
//
// Behaviors:
//   - The Go routine is not started automatically.
//   - In routine, use *ers.ErrNoError to exit the Go routine as nil is used to signal
//     that the function has finished successfully but the Go routine is still running.
func NewHandlerSimple(routine uc.MainFunc) *HandlerSimple {
	return &HandlerSimple{
		routine: routine,
	}
}

// GetErrChannel is a method of HandlerSimple that returns the error status of the Go routine.
//
// Returns:
//   - error: The error status of the Go routine.
func (h *HandlerSimple) GetErrChannel() <-chan error {
	return h.errChan
}
