package Helpers

import (
	slext "github.com/PlayerR9/MyGoLib/Utility/SliceExt"
)

// Helper is an interface that represents a helper.
type Helper interface {
}

// ResultEvaluator is a generic type that represents the result of a function evaluation.
type ResultEvaluator[H Helper, T any] struct {
	// hs is the slice of results of the function evaluation.
	hs []H
}

func (rs *ResultEvaluator[H, T]) SuccessOrFail() ([]H, bool) {
	// 1. Sort the slice of results by weight.
	slext.FilterByPositiveWeight()
}

// FilterSuccessOrFail filters the slice by returning only the successful
// results. However, if no successful results are found, it returns the
// original slice. The boolean indicates whether the slice was filtered
// or not.
//
// Parameters:
//   - batch: The slice to filter.
//
// Returns:
//   - []WeightedHelper: The filtered slice.
//   - bool: True if the slice was filtered, false otherwise.
func FilterSuccessOrFail[T any](batch []WeightedHelper[T]) ([]WeightedHelper[T], bool) {
	return slext.SFSeparateEarly(batch, FilterNilWeightedHelper)
}
