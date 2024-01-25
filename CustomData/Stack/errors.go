package Stack

import "fmt"

// StackOperationType is an integer type that represents the type of operation performed
// on a stack.
// It is used in error handling to specify the operation that caused an error.
type StackOperationType int

const (
	// Pop represents a pop operation, which removes an element from the stack.
	Pop StackOperationType = iota

	// Peek represents a peek operation, which retrieves the element at the front of the
	// stack without removing it.
	Peek
)

// String is a method of the StackOperationType type. It returns a string representation
// of the stack operation type.
//
// The method uses an array of strings where the index corresponds to the integer value
// of the StackOperationType.
// The string at the corresponding index is returned as the string representation of the
// StackOperationType.
//
// This method is typically used for error messages and logging.
func (qot StackOperationType) String() string {
	return [...]string{
		"pop",
		"peek",
	}[qot]
}

// ErrEmptyStack is a struct that represents an error when attempting to perform a stack
// operation on an empty stack.
// It has a single field, operation, of type StackOperationType, which indicates the type
// of operation that caused the error.
type ErrEmptyStack struct {
	operation StackOperationType
}

// NewErrEmptyStack creates a new ErrEmptyStack.
// It takes the following parameter:
//
//   - operation is the type of operation that caused the error.
//
// The function returns the following:
//
//   - A pointer to the new ErrEmptyStack.
func NewErrEmptyStack(operation StackOperationType) *ErrEmptyStack {
	return &ErrEmptyStack{operation: operation}
}

// Error is a method of the ErrEmptyStack type that implements the error interface. It
// returns a string representation of the error.
// The method constructs the error message by concatenating the string "could not ", the
// string representation of the operation that caused the error,
// and the string ": stack is empty". This provides a clear and descriptive error message
// when attempting to perform a stack operation on an empty stack.
func (e *ErrEmptyStack) Error() string {
	return fmt.Sprintf("could not %v: stack is empty", e.operation)
}

// ErrFullStack is a struct that represents an error when attempting to push an element
// into a full stack.
// It does not have any fields as the error condition is solely based on the state of the
// stack being full.
type ErrFullStack struct{}

// Error is a method of the ErrFullStack type that implements the error interface. It
// returns a string representation of the error.
// The method returns the string "could not push: stack is full", providing a clear and
// descriptive error message when attempting to push an element into a full stack.
func (e *ErrFullStack) Error() string {
	return "could not push: stack is full"
}

// ErrNegativeCapacity is a struct that represents an error when a negative capacity is
// provided for a stack.
// It does not have any fields as the error condition is solely based on the provided
// capacity being negative.
type ErrNegativeCapacity struct{}

// Error is a method of the ErrNegativeCapacity type that implements the error interface.
// It returns a string representation of the error.
// The method returns the string "capacity of a stack cannot be negative", providing a
// clear and descriptive error message when a negative capacity is provided for a stack.
func (e *ErrNegativeCapacity) Error() string {
	return "capacity of a stack cannot be negative"
}

// ErrTooManyValues is a struct that represents an error when too many values are
// provided for initializing a stack.
// It does not have any fields as the error condition is solely based on the number of
// provided values exceeding the capacity of the stack.
type ErrTooManyValues struct{}

// Error is a method of the ErrTooManyValues type that implements the error interface.
// It returns a string representation of the error.
// The method returns the string "could not initialize stack: too many values", providing
// a clear and descriptive error message when too many values are provided for initializing
// a stack.
func (e *ErrTooManyValues) Error() string {
	return "could not initialize stack: too many values"
}

// ErrOutOfBoundsIterator is a struct that represents an error when an iterator goes
// out of bounds.
// It does not have any fields as the error condition is solely based on the iterator
// exceeding the bounds of the data structure it is iterating over.
type ErrOutOfBoundsIterator struct{}

// Error is a method of the ErrOutOfBoundsIterator type that implements the error
// interface. It returns a string representation of the error.
// The method returns the string "iterator out of bounds", providing a clear and
// descriptive error message when an iterator goes out of bounds.
func (e *ErrOutOfBoundsIterator) Error() string {
	return "iterator out of bounds"
}
