// Package Helpers provides a set of helper functions and types that
// can be used for automatic error handling and result evaluation.
//
// However, this is still Work In Progress and is not yet fully
// implemented.
package Helpers

import (
	uf "github.com/PlayerR9/MyGoLib/Units/Functions"
	slext "github.com/PlayerR9/MyGoLib/Utility/SliceExt"
)

// HResult is a generic type that represents the result of a function
// evaluation.
type HResult[T any] struct {
	// Result is the result of the function evaluation.
	Result T

	// Reason is the error that occurred during the function evaluation.
	Reason error
}

// NewHResult creates a new HResult with the given result and reason.
//
// Parameters:
//   - result: The result of the function evaluation.
//   - reason: The error that occurred during the function evaluation.
//
// Returns:
//   - HResult: The new HResult.
func NewHResult[T any](result T, reason error) HResult[T] {
	return HResult[T]{
		Result: result,
		Reason: reason,
	}
}

// IsSuccess checks if the HResult is successful.
//
// Returns:
//   - bool: True if the HResult is successful, false otherwise.
func (hr HResult[T]) IsSuccess() bool {
	return hr.Reason == nil
}

// DoIfSuccess executes a function if the HResult is successful.
//
// Parameters:
//   - f: The function to execute.
func (hr HResult[T]) DoIfSuccess(f uf.DoFunc[T]) {
	if hr.Reason != nil {
		return
	}

	f(hr.Result)
}

// DoIfFailure executes a function if the HResult is a failure.
//
// Parameters:
//   - f: The function to execute.
func (hr HResult[T]) DoIfFailure(f uf.DualDoFunc[T, error]) {
	if hr.Reason == nil {
		return
	}

	f(hr.Result, hr.Reason)
}

// EvaluateFunc evaluates a function and returns the result as an HResult.
//
// Parameters:
//   - f: The function to evaluate.
//
// Returns:
//   - HResult: The result of the function evaluation.
func EvaluateFunc[T any](f EvalFunc[T]) HResult[T] {
	res, err := f()

	return HResult[T]{
		Result: res,
		Reason: err,
	}
}

// EvaluateMany evaluates a function that returns multiple values and
// returns the results as an array of HResults.
//
// Parameters:
//   - f: The function to evaluate.
//
// Returns:
//   - []HResult: The results of the function evaluation.
func EvaluateMany[T any](f EvalManyFunc[T]) []HResult[T] {
	res, err := f()

	if len(res) == 0 {
		return []HResult[T]{NewHResult(*new(T), err)}
	}

	results := make([]HResult[T], len(res))

	for i, r := range res {
		results[i] = NewHResult(r, err)
	}

	return results
}

// FilterOut removes all elements from the slice that have an error.
//
// Parameters:
//   - batch: The slice to filter.
//
// Returns:
//   - []HResult: The filtered slice.
func FilterOut[T any](batch []HResult[T]) []HResult[T] {
	batch = slext.SliceFilter(batch, FilterNilHResult)

	return batch
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
//   - []HResult: The filtered slice.
//   - bool: True if the slice was filtered, false otherwise.
func FilterSuccessOrFail[T any](batch []HResult[T]) ([]HResult[T], bool) {
	return slext.SFSeparateEarly(batch, FilterNilHResult)
}
