package Errors

import (
	"sync"
)

// GRHandler is a struct that represents a Go routine handler.
// It is used to handle the result of a Go routine.
type GRHandler struct {
	// wg is a WaitGroup that is used to wait for the Go routine to finish.
	wg sync.WaitGroup

	// id is the identifier of the Go routine.
	id string

	// errStatus is the error status of the Go routine.
	errStatus error
}

// GoRun is a function that runs a Go routine and returns a GRHandler
// to handle the result of the Go routine.
//
// Parameters:
//
//   - id: The identifier of the Go routine.
//   - routine: The Go routine to run.
//
// Returns:
//
//   - *GRHandler: A pointer to the GRHandler that handles the result of the Go routine.
//
// The Go routine is automatically run when this function is called.
func GoRun(id string, routine func() error) *GRHandler {
	h := &GRHandler{
		id:        id,
		errStatus: nil,
	}

	h.wg.Add(1)

	go func() {
		defer h.wg.Done()

		defer func() {
			if r := recover(); r != nil {
				h.errStatus = &ErrPanic{value: r}
			}
		}()

		h.errStatus = routine()
	}()

	return h
}

// Wait is a method of GRHandler that waits for the Go routine to finish
// and returns the error status of the Go routine.
//
// Returns:
//
//   - error: The error status of the Go routine.
func (h *GRHandler) Wait() error {
	h.wg.Wait()

	return h.errStatus
}

// Identifier is a method of GRHandler that returns the identifier of the Go routine.
//
// Returns:
//
//   - string: The identifier of the Go routine.
func (h *GRHandler) Identifier() string {
	return h.id
}

// BatchBuilder is a struct that represents a Go routine batch builder.
// It is used to build a batch of Go routines and their handlers.
type BatchBuilder struct {
	// routines is a slice of Go routines.
	routines []func() error

	// ids is a slice of identifiers for the Go routines.
	ids []string
}

// Add is a method of BatchBuilder that adds a Go routine to the batch.
//
// Parameters:
//
//   - identifier: The identifier of the Go routine.
//   - routine: The Go routine to add to the batch.
func (b *BatchBuilder) Add(identifier string, routine func() error) {
	b.routines = append(b.routines, routine)
	b.ids = append(b.ids, identifier)
}

// Build is a method of BatchBuilder that builds the batch of Go routines
// and their handlers.
//
// Returns:
//
//   - []*GRHandler: A slice of pointers to the GRHandler instances that handle the Go routines.
//
// All Go routines are automatically run when this method is called.
// Finally, the batch is cleared after the Go routines are built.
func (b *BatchBuilder) Build() []*GRHandler {
	if b.routines == nil {
		return nil
	}

	pairings := make([]*GRHandler, 0, len(b.routines))

	for i, routine := range b.routines {
		pairings = append(pairings, GoRun(b.ids[i], routine))
	}

	b.routines = nil
	b.ids = nil

	return pairings
}

// Clear is a method of BatchBuilder that clears the batch of Go routines.
func (b *BatchBuilder) Clear() {
	b.routines = nil
	b.ids = nil
}

// WaitAll is a function that waits for all Go routines in the batch to finish
// and returns a slice of errors that represent the error statuses of the Go routines.
//
// Parameters:
//
//   - batch: A slice of pointers to the GRHandler instances that handle the Go routines.
//
// Returns:
//
//   - []error: A slice of errors that represent the error statuses of the Go routines.
func WaitAll(batch []*GRHandler) []error {
	errs := make([]error, 0, len(batch))

	for _, pair := range batch {
		if pair == nil {
			continue
		}

		errs = append(errs, pair.Wait())
	}

	return errs
}
