// Package Helpers provides a set of helper functions and types that
// can be used for automatic error handling and result evaluation.
//
// However, this is still Work In Progress and is not yet fully
// implemented.
package Helpers

import (
	uc "github.com/PlayerR9/MyGoLib/Units/Common"
	cdp "github.com/PlayerR9/MyGoLib/Units/Pair"
	slext "github.com/PlayerR9/MyGoLib/Utility/SliceExt"
)

// HResult is a generic type that represents the result of a function
// evaluation.
type HResult[T any] cdp.Pair[T, error]

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
		First:  result,
		Second: reason,
	}
}

// IsSuccess checks if the HResult is successful.
//
// Returns:
//   - bool: True if the HResult is successful, false otherwise.
func (hr HResult[T]) IsSuccess() bool {
	return hr.Second == nil
}

// DoIfSuccess executes a function if the HResult is successful.
//
// Parameters:
//   - f: The function to execute.
func (hr HResult[T]) DoIfSuccess(f uc.DoFunc[T]) {
	if hr.Second != nil {
		return
	}

	f(hr.First)
}

// DoIfFailure executes a function if the HResult is a failure.
//
// Parameters:
//   - f: The function to execute.
func (hr HResult[T]) DoIfFailure(f uc.DualDoFunc[T, error]) {
	if hr.Second == nil {
		return
	}

	f(hr.First, hr.Second)
}

// EvaluateFunc evaluates a function and returns the result as an HResult.
//
// Parameters:
//   - elem: The element to evaluate.
//   - f: The function to evaluate.
//
// Returns:
//   - HResult: The result of the function evaluation.
func EvaluateFunc[E, R any](elem E, f uc.EvalOneFunc[E, R]) HResult[R] {
	res, err := f(elem)

	return HResult[R]{
		First:  res,
		Second: err,
	}
}

// EvaluateMany evaluates a function that returns multiple values and
// returns the results as an array of HResults.
//
// Parameters:
//   - elem: The element to evaluate.
//   - f: The function to evaluate.
//
// Returns:
//   - []HResult: The results of the function evaluation.
func EvaluateMany[E, R any](elem E, f uc.EvalManyFunc[E, R]) []HResult[R] {
	res, err := f(elem)

	if len(res) == 0 {
		return []HResult[R]{
			{
				First:  *new(R),
				Second: err,
			},
		}
	}

	results := make([]HResult[R], len(res))

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

// ApplyFunc applies a function to each element in the slice and returns
// the results as an array of HResults.
//
// Parameters:
//   - elems: The slice to apply the function to.
//   - f: The function to apply.
//
// Returns:
//   - []HResult: The results of the function application.
func ApplyFunc[E, R any](elems []E, f uc.EvalOneFunc[E, R]) []HResult[R] {
	results := make([]HResult[R], len(elems))

	for i, elem := range elems {
		res, err := f(elem)

		results[i] = HResult[R]{
			First:  res,
			Second: err,
		}
	}

	return results
}

// ExtractResults extracts the results from the HResults.
//
// Parameters:
//   - batch: The slice of HResults.
//
// Returns:
//   - []T: The extracted results.
func ExtractResults[T any](batch []HResult[T]) []T {
	results := make([]T, len(batch))

	for i, res := range batch {
		results[i] = res.First
	}

	return results
}
