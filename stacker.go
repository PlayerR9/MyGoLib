package MyGoLib

import (
	"fmt"

	uc "github.com/PlayerR9/MyGoLib/Units/common"
)

// Stacker is an interface that defines methods for a stack data structure.
type Stacker[T any] interface {
	// Push is a method that adds a value of type T to the end of the stack.
	//
	// Parameters:
	//   - value: The value of type T to add to the stack.
	//
	// Returns:
	//   - bool: True if the value was successfully added to the stack, false otherwise.
	Push(value T) bool

	// PushMany is a method that adds multiple values of type T to the end of the stack.
	//
	// Parameters:
	//   - values: The values of type T to add to the stack.
	//
	// Returns:
	//   - int: The number of values that were successfully added to the stack.
	PushMany(values []T) int

	// Pop is a method that pops an element from the stack and returns it.
	//
	// Returns:
	//   - T: The value of type T that was popped.
	//   - bool: True if the value was successfully popped, false otherwise.
	Pop() (T, bool)

	// Peek is a method that returns the value at the front of the stack without removing
	// it.
	//
	// Returns:
	//   - T: The value of type T at the front of the stack.
	//   - bool: True if the value was successfully peeked, false otherwise.
	Peek() (T, bool)

	// IsEmpty is a method that checks whether the list is empty.
	//
	// Returns:
	//
	//   - bool: True if the list is empty, false otherwise.
	IsEmpty() bool

	// Size method returns the number of elements currently in the list.
	//
	// Returns:
	//
	//   - int: The number of elements in the list.
	Size() int

	// Clear method is used to remove all elements from the list, making it empty.
	Clear()

	// Capacity is a method that returns the maximum number of elements that the list can hold.
	//
	// Returns:
	//
	//   - int: The maximum number of elements that the list can hold. -1 if there is no limit.
	Capacity() int

	// IsFull is a method that checks whether the list is full.
	//
	// Returns:
	//
	//   - bool: True if the list is full, false otherwise.
	IsFull() bool

	uc.Slicer[T]
	uc.Copier
	fmt.GoStringer
}
