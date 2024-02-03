package ListLike

import "fmt"

// ErrEmptyList is a struct that represents an error when attempting to perform a list
// operation on an empty list.
// It has a single field, operation, of type ListOperationType, which indicates the type
// of operation that caused the error.
type ErrEmptyList struct {
	list string
}

// NewErrEmptyList creates a new ErrEmptyList.
// It takes the following parameter:
//
//   - operation is the type of operation that caused the error.
//
// The function returns the following:
//
//   - A pointer to the new ErrEmptyList.
func NewErrEmptyList[T any](list Lister[T]) *ErrEmptyList {
	return &ErrEmptyList{list: fmt.Sprintf("%T", list)}
}

// Error is a method of the ErrEmptyList type that implements the error interface. It
// returns a string representation of the error.
// The method constructs the error message by concatenating the string "could not ", the
// string representation of the operation that caused the error,
// and the string ": list is empty". This provides a clear and descriptive error message
// when attempting to perform a list operation on an empty list.
func (e *ErrEmptyList) Error() string {
	return fmt.Sprintf("list (%v) is empty", e.list)
}

// ErrFullList is a struct that represents an error when attempting to prepend an element
// into a full list.
// It does not have any fields as the error condition is solely based on the state of the
// list being full.
type ErrFullList struct {
	list string
}

// NewErrFullList creates a new ErrFullList.
// It takes the following parameter:
//
//   - operation is the type of operation that caused the error.
//
// The function returns the following:
//
//   - A pointer to the new ErrFullList.
func NewErrFullList[T any](list Lister[T]) *ErrFullList {
	return &ErrFullList{list: fmt.Sprintf("%T", list)}
}

// Error is a method of the ErrFullList type that implements the error interface. It
// returns a string representation of the error.
// The method returns the string "could not prepend: list is full", providing a clear and
// descriptive error message when attempting to prepend an element into a full list.
func (e *ErrFullList) Error() string {
	return fmt.Sprintf("list (%v) is full", e.list)
}
