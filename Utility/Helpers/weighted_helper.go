// Package Helpers provides a set of helper functions and types that
// can be used for automatic error handling and result evaluation.
//
// However, this is still Work In Progress and is not yet fully
// implemented.
package Helpers

import (
	up "github.com/PlayerR9/MyGoLib/Units/Pair"
)

// WeightedHelper is a generic type that represents the result of a function
// evaluation.
type WeightedHelper[O any] struct {
	// result is the result of the function evaluation.
	result O

	// reason is the error that occurred during the function evaluation.
	reason error

	// weight is the weight of the result (i.e., the probability of the result being correct)
	// or the most likely error (if the result is an error).
	weight float64
}

// GetData returns the result of the function evaluation.
//
// Returns:
//   - O: The result of the function evaluation.
func (h *WeightedHelper[O]) GetData() *up.Pair[O, error] {
	return up.NewPair(h.result, h.reason)
}

// GetReason returns the error that occurred during the function evaluation.
//
// Returns:
//   - error: The error that occurred during the function evaluation.
func (h *WeightedHelper[O]) GetWeight() float64 {
	return h.weight
}

// NewWeightedHelper creates a new WeightedHelper with the given result, reason, and weight.
//
// Parameters:
//   - result: The result of the function evaluation.
//   - reason: The error that occurred during the function evaluation.
//   - weight: The weight of the result. The higher the weight, the more likely the result
//     is correct.
//
// Returns:
//   - WeightedHelper: The new WeightedHelper.
func NewWeightedHelper[O any](result O, reason error, weight float64) *WeightedHelper[O] {
	return &WeightedHelper[O]{
		result: result,
		reason: reason,
		weight: weight,
	}
}
