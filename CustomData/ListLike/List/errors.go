// Package errors provides error handling for list operations.
// It includes error types for empty and full list conditions.
package ListLike

import "fmt"

// ErrEmptyList represents an error when attempting to perform
// a list operation on an empty list.
type ErrEmptyList struct {
	// list is the type of list that caused the error.
	list string
}

// NewErrEmptyList creates a new ErrEmptyList.
//
// Parameters:
//
//   - list: the type of list that caused the error.
//
// Returns:
//
//   - *ErrEmptyList: a pointer to the new ErrEmptyList.
func NewErrEmptyList[T any](list Lister[T]) *ErrEmptyList {
	return &ErrEmptyList{list: fmt.Sprintf("%T", list)}
}

// Error implements the error interface for the ErrEmptyList type.
//
// Returns:
//
//   - string: a string representation of the error.
func (e *ErrEmptyList) Error() string {
	return fmt.Sprintf("list (%v) is empty", e.list)
}

// ErrFullList represents an error when attempting to prepend an element
// into a full list.
type ErrFullList struct {
	// list is the type of list that caused the error.
	list string
}

// NewErrFullList creates a new ErrFullList.
//
// Parameters:
//
//   - list: the type of list that caused the error.
//
// Returns:
//
//   - *ErrFullList: a pointer to the new ErrFullList.
func NewErrFullList[T any](list Lister[T]) *ErrFullList {
	return &ErrFullList{list: fmt.Sprintf("%T", list)}
}

// Error implements the error interface for the ErrFullList type.
//
// Returns:
//
//   - string: a string representation of the error.
func (e *ErrFullList) Error() string {
	return fmt.Sprintf("list (%v) is full", e.list)
}
