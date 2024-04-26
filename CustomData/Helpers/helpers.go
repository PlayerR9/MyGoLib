// Package Helpers provides a set of helper functions and types that
// can be used for automatic error handling and result evaluation.
//
// However, this is still Work In Progress and is not yet fully
// implemented.
package Helpers

// EvalFunc is a generic type that represents a function that returns
// a value and an error.
//
// Returns:
//   - T: The result of the function evaluation.
//   - error: The error that occurred during the function evaluation.
type EvalFunc[T any] func() (T, error)

// EvalManyFunc is a generic type that represents a function that returns
// multiple values and an error.
//
// Returns:
//   - []T: The result of the function evaluation.
//   - error: The error that occurred during the function evaluation.
type EvalManyFunc[T any] func() ([]T, error)

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
		return []HResult[T]{NewHResult[T](*new(T), err)}
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
func FilterOut[T any](batch []HResult[T]) {
	// 1. Remove all elements from the slice that have an error
	for i := 0; i < len(batch); i++ {
		if batch[i].Reason != nil {
			batch = append(batch[:i], batch[i+1:]...)
			i--
		}
	}
}
