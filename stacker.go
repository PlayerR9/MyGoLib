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

	// CutNilValues is a method that removes all nil values from the list.
	// It is useful for cleaning up the list and removing any empty or nil elements.
	CutNilValues()

	uc.Slicer[T]
	uc.Copier
	fmt.GoStringer
}
