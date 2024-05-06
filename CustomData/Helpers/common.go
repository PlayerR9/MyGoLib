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
