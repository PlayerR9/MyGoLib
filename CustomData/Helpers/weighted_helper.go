// Package Helpers provides a set of helper functions and types that
// can be used for automatic error handling and result evaluation.
//
// However, this is still Work In Progress and is not yet fully
// implemented.
package Helpers

// WeightedHelper is a generic type that represents the result of a function
// evaluation.
type WeightedHelper[T any] struct {
	// result is the result of the function evaluation.
	result T

	// reason is the error that occurred during the function evaluation.
	reason error

	// weight is the weight of the result (i.e., the probability of the result being correct)
	// or the most likely error (if the result is an error).
	weight float64
}

// GetData returns the result of the function evaluation.
//
// Returns:
//   - T: The result of the function evaluation.
func (h *WeightedHelper[T]) GetData() T {
	return h.result
}

// GetReason returns the error that occurred during the function evaluation.
//
// Returns:
//   - error: The error that occurred during the function evaluation.
func (h *WeightedHelper[T]) GetWeight() float64 {
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
func NewWeightedHelper[T any](result T, reason error, weight float64) *WeightedHelper[T] {
	return &WeightedHelper[T]{
		result: result,
		reason: reason,
		weight: weight,
	}
}

/*


// IsSuccess checks if the WeightedHelper is successful.
//
// Returns:
//   - bool: True if the WeightedHelper is successful, false otherwise.
func (hr WeightedHelper[T]) IsSuccess() bool {
	return hr.Reason == nil
}

// DoIfSuccess executes a function if the WeightedHelper is successful.
//
// Parameters:
//   - f: The function to execute.
func (hr WeightedHelper[T]) DoIfSuccess(f uc.DoFunc[T]) {
	if hr.Reason != nil {
		return
	}

	f(hr.Result)
}

// DoIfFailure executes a function if the WeightedHelper is a failure.
//
// Parameters:
//   - f: The function to execute.
func (hr WeightedHelper[T]) DoIfFailure(f uc.DualDoFunc[T, error]) {
	if hr.Reason == nil {
		return
	}

	f(hr.Result, hr.Reason)
}

// EvaluateFunc evaluates a function and returns the result as an WeightedHelper.
//
// Parameters:
//   - elem: The element to evaluate.
//   - f: The function to evaluate.
//
// Returns:
//   - WeightedHelper: The result of the function evaluation.
func EvaluateFunc[E, R any](elem E, f uc.EvalOneFunc[E, R]) WeightedHelper[R] {
	res, err := f(elem)

	return WeightedHelper[R]{
		Result: res,
		Reason: err,
		Weight: 0.0, // TODO: Implement weight calculation
	}
}

// EvaluateMany evaluates a function that returns multiple values and
// returns the results as an array of WeightedHelpers.
//
// Parameters:
//   - elem: The element to evaluate.
//   - f: The function to evaluate.
//
// Returns:
//   - []WeightedHelper: The results of the function evaluation.
func EvaluateMany[E, R any](elem E, f uc.EvalManyFunc[E, R]) []WeightedHelper[R] {
	res, err := f(elem)

	if len(res) == 0 {
		h := NewWeightedHelper(*new(R), err, 0.0) // TODO: Implement weight calculation

		return []WeightedHelper[R]{h}
	}

	results := make([]WeightedHelper[R], len(res))

	for i, r := range res {
		results[i] = NewWeightedHelper(r, err, 0.0) // TODO: Implement weight calculation
	}

	return results
}

// FilterOut removes all elements from the slice that have an error.
//
// Parameters:
//   - batch: The slice to filter.
//
// Returns:
//   - []WeightedHelper: The filtered slice.
func FilterOut[T any](batch []WeightedHelper[T]) []WeightedHelper[T] {
	batch = slext.SliceFilter(batch, FilterNilWeightedHelper)

	return batch
}



// ApplyFunc applies a function to each element in the slice and returns
// the results as an array of WeightedHelpers.
//
// Parameters:
//   - elems: The slice to apply the function to.
//   - f: The function to apply.
//
// Returns:
//   - []WeightedHelper: The results of the function application.
func ApplyFunc[E, R any](elems []E, f uc.EvalOneFunc[E, R]) []WeightedHelper[R] {
	results := make([]WeightedHelper[R], len(elems))

	for i, elem := range elems {
		res, err := f(elem)

		results[i] = NewWeightedHelper(res, err, 0.0) // TODO: Implement weight calculation
	}

	return results
}

// ExtractResults extracts the results from the WeightedHelpers.
//
// Parameters:
//   - batch: The slice of WeightedHelpers.
//
// Returns:
//   - []T: The extracted results.
func ExtractResults[T any](batch []WeightedHelper[T]) []T {
	results := make([]T, len(batch))

	for i, res := range batch {
		results[i] = res.Result
	}

	return results
}
*/
