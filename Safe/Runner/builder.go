package Runner

import (
	uc "github.com/PlayerR9/MyGoLib/Units/Common"
)

// BatchBuilder is a struct that represents a Go routine batch builder.
// It is used to build a batch of Go routines and their handlers.
type BatchBuilder struct {
	// routines is a slice of Go routines.
	routines []uc.MainFunc

	// ids is a slice of identifiers for the Go routines.
	ids []string
}

// Add is a method of BatchBuilder that adds a Go routine to the batch.
//
// Parameters:
//   - identifier: The identifier of the Go routine.
//   - routine: The Go routine to add to the batch.
func (b *BatchBuilder) Add(identifier string, routine uc.MainFunc) {
	b.routines = append(b.routines, routine)
	b.ids = append(b.ids, identifier)
}

// Build is a method of BatchBuilder that builds the batch of Go routines
// and their handlers.
//
// Returns:
//   - []*HandlerSimple: A slice of pointers to the HandlerSimple instances that handle the Go routines.
//
// All Go routines are automatically run when this method is called.
// Finally, the batch is cleared after the Go routines are built.
func (b *BatchBuilder) Build() map[string]*HandlerSimple {
	if b.routines == nil {
		return nil
	}

	handlers := make(map[string]*HandlerSimple)

	for i, routine := range b.routines {
		h := NewHandlerSimple(routine)

		h.Start()

		handlers[b.ids[i]] = h
	}

	b.routines = nil
	b.ids = nil

	return handlers
}

// Clear is a method of BatchBuilder that clears the batch of Go routines.
func (b *BatchBuilder) Clear() {
	b.routines = nil
	b.ids = nil
}
