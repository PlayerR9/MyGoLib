package Sorting

import "fmt"

// ErrNotComparable is an error type that is returned when two values are not
// comparable.
type ErrNotComparable[A, B any] struct {
	// First is the first value that is not comparable.
	First A

	// Second is the second value that is not comparable.
	Second B
}

// Error returns the error message: "values <First> and <Second> are not comparable".
//
// Returns:
//   - string: The error message.
func (e *ErrNotComparable[A, B]) Error() string {
	return fmt.Sprintf("values %T and %T are not comparable", e.First, e.Second)
}

// NewErrNotComparable creates a new ErrNotComparable error with the provided values.
//
// Parameters:
//   - first: The first value that is not comparable.
//   - second: The second value that is not comparable.
//
// Returns:
//   - *ErrNotComparable: A pointer to the new error.
func NewErrNotComparable[A, B any](first A, second B) *ErrNotComparable[A, B] {
	return &ErrNotComparable[A, B]{First: first, Second: second}
}
