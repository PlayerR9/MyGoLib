// Package Helpers provides a set of helper functions and types that
// can be used for automatic error handling and result evaluation.
//
// However, this is still Work In Progress and is not yet fully
// implemented.
package Helpers

// SimpleHelper is a type that represents the result of a function evaluation
// that can either be successful or a failure.
type SimpleHelper[T any] struct {
	// result is the result of the function evaluation.
	result T

	// reason is the error that occurred during the function evaluation.
	reason error
}

// NewSimpleHelper creates a new SimpleHelper with the given result and reason.
//
// Parameters:
//   - result: The result of the function evaluation.
//   - reason: The error that occurred during the function evaluation.
//
// Returns:
//   - SimpleHelper: The new SimpleHelper.
func NewSimpleHelper[T any](result T, reason error) *SimpleHelper[T] {
	return &SimpleHelper[T]{
		result: result,
		reason: reason,
	}
}

/*
// IsSuccess checks if the SimpleHelper is successful.
//
// Returns:
//   - bool: True if the SimpleHelper is successful, false otherwise.
func (h *SimpleHelper[T]) IsSuccess() bool {
	return h.Reason == nil
}

// DoIfSuccess executes a function if the SimpleHelper is successful.
//
// Parameters:
//   - f: The function to execute.
func (h *SimpleHelper[T]) DoIfSuccess(f uc.DoFunc[T]) {
	if h.Reason == nil {
		f(h.Result)
	}
}

// DoIfSuccess executes a function if the SimpleHelper is successful.
//
// Parameters:
//   - f: The function to execute.
func (h *SimpleHelper[T]) DoIfFailure(f uc.DoFunc[T]) {
	if h.Reason != nil {
		f(h.Result)
	}
}


// EvaluateFunc evaluates a function and returns the result as an HResult.
//
// Parameters:
//   - elem: The element to evaluate.
//   - f: The function to evaluate.
//
// Returns:
//   - HResult: The result of the function evaluation.
func EvaluateFunc[E, R any](elem E, f uc.EvalOneFunc[E, R]) SimpleHelper[R] {
	res, err := f(elem)

	return SimpleHelper[R]{
		Result: res,
		Reason: err,
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
func EvaluateMany[E, R any](elem E, f uc.EvalManyFunc[E, R]) []SimpleHelper[R] {
	res, err := f(elem)

	if len(res) == 0 {
		h := NewHResult(*new(R), err, 0.0) // TODO: Implement weight calculation

		return []SimpleHelper[R]{h}
	}

	results := make([]SimpleHelper[R], len(res))

	for i, r := range res {
		results[i] = NewHResult(r, err, 0.0) // TODO: Implement weight calculation
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
func FilterOut[T any](batch []SimpleHelper[T]) []SimpleHelper[T] {
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
func FilterSuccessOrFail[T any](batch []SimpleHelper[T]) ([]SimpleHelper[T], bool) {
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
func ApplyFunc[E, R any](elems []E, f uc.EvalOneFunc[E, R]) []SimpleHelper[R] {
	results := make([]SimpleHelper[R], len(elems))

	for i, elem := range elems {
		res, err := f(elem)

		results[i] = NewHResult(res, err, 0.0) // TODO: Implement weight calculation
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
func ExtractResults[T any](batch []SimpleHelper[T]) []T {
	results := make([]T, len(batch))

	for i, res := range batch {
		results[i] = res.Result
	}

	return results
}
*/
