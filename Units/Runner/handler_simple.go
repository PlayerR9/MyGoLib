package Runner

import (
	"sync"

	uc "github.com/PlayerR9/MyGoLib/Units/common"
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

// Start implements the Runner interface.
func (h *HandlerSimple) Start() {
	if h.stopChan != nil {
		return
	}

	h.errChan = make(chan error)
	h.stopChan = make(chan struct{})

	h.wg.Add(1)

	go h.run()
}

// Close implements the Runner interface.
func (h *HandlerSimple) Close() {
	if h.stopChan == nil {
		return
	}

	close(h.stopChan)
	h.stopChan = nil

	h.wg.Wait()
}

// IsClosed implements the Runner interface.
func (h *HandlerSimple) IsClosed() bool {
	return h.errChan == nil
}

// ReceiveErr implements the Runner interface.
func (h *HandlerSimple) ReceiveErr() (error, bool) {
	if h.errChan == nil {
		return nil, false
	}

	err, ok := <-h.errChan
	if !ok {
		return nil, false
	} else {
		return err, true
	}
}

// run is a private method of HandlerSimple that is runned by the Go routine.
//
// Behaviors:
//   - Use ers.ErrNoError to exit the Go routine as nil is used to signal
//     that the function has finished successfully but the Go routine is still running.
func (h *HandlerSimple) run() {
	defer h.wg.Done()

	defer func() {
		r := recover()

		if r != nil {
			h.errChan <- ers.NewErrPanic(r)
		}

		h.clean()
	}()

	for {
		select {
		case <-h.stopChan:
			return
		default:
			err := h.routine()
			if err != nil {
				h.errChan <- err
				return
			}
		}
	}
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
//   - If routine is nil, this function returns nil.
//   - The Go routine is not started automatically.
//   - In routine, use *ers.ErrNoError to exit the Go routine as nil is used to signal
//     that the function has finished successfully but the Go routine is still running.
func NewHandlerSimple(routine uc.MainFunc) *HandlerSimple {
	if routine == nil {
		return nil
	}

	return &HandlerSimple{
		routine: routine,
	}
}

// clean is a private method of HandlerSimple that cleans up the handler.
func (h *HandlerSimple) clean() {
	if h.errChan != nil {
		close(h.errChan)
		h.errChan = nil
	}

	if h.stopChan != nil {
		close(h.stopChan)
		h.stopChan = nil
	}
}
