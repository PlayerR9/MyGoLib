package Common

// ObserverFunc is a function type that takes a value of type T and returns an error.
// It is used to observe the values of the nodes during a traversal.
type ObserverFunc[T any] func(T) error

// DoFunc is a generic type that represents a function that takes a value
// and does something with it.
//
// Parameters:
//   - T: The type of the value.
type DoFunc[T any] func(T)

// DualDoFunc is a generic type that represents a function that takes two
// values and does something with them.
//
// Parameters:
//   - T: The type of the first value.
//   - U: The type of the second value.
type DualDoFunc[T any, U any] func(T, U)

// EvalOneFunc is a function that evaluates one element.
//
// Parameters:
//   - elem: The element to evaluate.
//
// Returns:
//   - R: The result of the evaluation.
//   - error: An error if the evaluation failed.
type EvalOneFunc[E, R any] func(elem E) (R, error)

// EvalManyFunc is a function that evaluates many elements.
//
// Parameters:
//   - elem: The element to evaluate.
//
// Returns:
//   - []R: The results of the evaluation.
//   - error: An error if the evaluation failed.
type EvalManyFunc[E, R any] func(elem E) ([]R, error)
