package Lister

import "fmt"

// ErrFullList is an error type for a full list.
type ErrFullList[T any] struct {
	// List is the list that is full.
	List Lister[T]
}

// Error returns the error message: "list (%T) is full".
//
// Returns:
//   - string: The error message.
//
// Behaviors:
//   - If the list is nil, the error message is "list is full".
func (e *ErrFullList[T]) Error() string {
	if e.List == nil {
		return "list is full"
	} else {
		return fmt.Sprintf("list (%T) is full", e.List)
	}
}

// NewErrFullList is a constructor for ErrFullList.
//
// Parameters:
//   - list: The list that is full.
//
// Returns:
//   - *ErrFullList: The error.
func NewErrFullList[T any](list Lister[T]) *ErrFullList[T] {
	return &ErrFullList[T]{List: list}
}
