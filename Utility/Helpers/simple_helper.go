// Package Helpers provides a set of helper functions and types that
// can be used for automatic error handling and result evaluation.
//
// However, this is still Work In Progress and is not yet fully
// implemented.
package Helpers

import up "github.com/PlayerR9/MyGoLib/Units/Pair"

// SimpleHelper is a type that represents the result of a function evaluation
// that can either be successful or a failure.
type SimpleHelper[O any] struct {
	// result is the result of the function evaluation.
	result O

	// reason is the error that occurred during the function evaluation.
	reason error
}

// GetData returns the result of the function evaluation.
//
// Returns:
//   - *up.Pair[O, error]: The result of the function evaluation.
func (h *SimpleHelper[O]) GetData() *up.Pair[O, error] {
	return up.NewPair(h.result, h.reason)
}

// GetWeight returns the weight of the element.
//
// Returns:
//   - float64: The weight of the element.
func (h *SimpleHelper[O]) GetWeight() float64 {
	return 0.0
}

// NewSimpleHelper creates a new SimpleHelper with the given result and reason.
//
// Parameters:
//   - result: The result of the function evaluation.
//   - reason: The error that occurred during the function evaluation.
//
// Returns:
//   - SimpleHelper: The new SimpleHelper.
func NewSimpleHelper[O any](result O, reason error) *SimpleHelper[O] {
	return &SimpleHelper[O]{
		result: result,
		reason: reason,
	}
}
