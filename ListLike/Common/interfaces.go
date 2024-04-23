package Common

import (
	"fmt"

	itf "github.com/PlayerR9/MyGoLib/CustomData/Iterators"
	itff "github.com/PlayerR9/MyGoLib/Units/Interfaces"
)

// ListLike is an interface that defines methods for a list data structure.
type ListLike[T any] interface {
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

	itf.Iterable[T]

	itff.Slicer[T]

	itff.Copier

	fmt.Stringer
}

// Lister is an interface that defines methods for a list data structure.
type Lister[T any] interface {
	// Append is a method that adds a value of type T to the end of the list.
	//
	// Parameters:
	//
	//   - value: The value of type T to add to the list.
	//
	// Returns:
	//
	//   - error: An error if the list is full.
	Append(value T) error

	// DeleteFirst is a method that deletes an element from the front of the list and
	// returns it. If the list is empty, it will panic.
	//
	// Returns:
	//
	//   - T: The value of type T that was deleted.
	DeleteFirst() T

	// PeekFirst is a method that returns the value at the front of the list without
	// removing it. If the list is empty, it will panic.
	//
	// Returns:
	//
	//   - T: The value of type T at the front of the list.
	PeekFirst() T

	// Prepend is a method that adds a value of type T to the end of the list.
	//
	// Parameters:
	//
	//   - value: The value of type T to add to the list.
	//
	// Returns:
	//
	//   - error: An error if the list is full.
	Prepend(value T) error

	// DeleteLast is a method that deletes an element from the end of the list and
	// returns it. If the list is empty, it will panic.
	//
	// Returns:
	//
	//   - T: The value of type T that was deleted.
	DeleteLast() T

	// PeekLast is a method that returns the value at the end of the list without
	// removing it. If the list is empty, it will panic.
	//
	// Returns:
	//
	//   - T: The value of type T at the end of the list.
	PeekLast() T

	ListLike[T]
}

// Queuer is an interface that defines methods for a queue data structure.
type Queuer[T any] interface {
	// Enqueue is a method that adds a value of type T to the end of the queue.
	// If the queue is full, it will panic.
	//
	// Parameters:
	//
	//   - value: The value of type T to add to the queue.
	Enqueue(value T)

	// Dequeue is a method that dequeues an element from the queue and returns it.
	// If the queue is empty, it will panic.
	//
	// Returns:
	//
	//   - T: The value of type T that was dequeued.
	Dequeue() T

	// Peek is a method that returns the value at the front of the queue without
	// removing it.
	// If the queue is empty, it will panic.
	//
	// Returns:
	//
	//   - T: The value of type T at the front of the queue.
	Peek() T

	ListLike[T]
}

// Stacker is an interface that defines methods for a stack data structure.
type Stacker[T any] interface {
	// Push is a method that adds a value of type T to the end of the stack.
	// If the stack is full, it will panic.
	//
	// Parameters:
	//
	//   - value: The value of type T to add to the stack.
	Push(value T)

	// Pop is a method that pops an element from the stack and returns it.
	// If the stack is empty, it will panic.
	//
	// Returns:
	//
	//   - T: The value of type T that was popped.
	Pop() T

	// Peek is a method that returns the value at the front of the stack without removing
	// it.
	// If the stack is empty, it will panic.
	//
	// Returns:
	//
	//   - T: The value of type T at the front of the stack.
	Peek() T

	ListLike[T]
}

// Lister is an interface that defines methods for a list data structure.
type SafeLister[T any] interface {
	// Append is a method that adds a value of type T to the end of the list.
	//
	// Parameters:
	//
	//   - value: The value of type T to add to the list.
	//
	// Returns:
	//
	//   - error: An error if the list is full.
	Append(value T) error

	// DeleteFirst is a method that deletes an element from the front of the list and
	// returns it. If the list is empty, it will panic.
	//
	// Returns:
	//
	//   - T: The value of type T that was deleted.
	DeleteFirst() (T, error)

	// PeekFirst is a method that returns the value at the front of the list without
	// removing it. If the list is empty, it will panic.
	//
	// Returns:
	//
	//   - T: The value of type T at the front of the list.
	PeekFirst() (T, error)

	// Prepend is a method that adds a value of type T to the end of the list.
	//
	// Parameters:
	//
	//   - value: The value of type T to add to the list.
	//
	// Returns:
	//
	//   - error: An error if the list is full.
	Prepend(value T) error

	// DeleteLast is a method that deletes an element from the end of the list and
	// returns it. If the list is empty, it will panic.
	//
	// Returns:
	//
	//   - T: The value of type T that was deleted.
	DeleteLast() (T, error)

	// PeekLast is a method that returns the value at the end of the list without
	// removing it. If the list is empty, it will panic.
	//
	// Returns:
	//
	//   - T: The value of type T at the end of the list.
	PeekLast() (T, error)

	ListLike[T]
}

// Queuer is an interface that defines methods for a queue data structure.
type SafeQueuer[T any] interface {
	// Enqueue is a method that adds a value of type T to the end of the queue.
	// If the queue is full, it will panic.
	//
	// Parameters:
	//
	//   - value: The value of type T to add to the queue.
	Enqueue(value T) error

	// Dequeue is a method that dequeues an element from the queue and returns it.
	// If the queue is empty, it will panic.
	//
	// Returns:
	//
	//   - T: The value of type T that was dequeued.
	Dequeue() (T, error)

	// Peek is a method that returns the value at the front of the queue without
	// removing it.
	// If the queue is empty, it will panic.
	//
	// Returns:
	//
	//   - T: The value of type T at the front of the queue.
	Peek() (T, error)

	ListLike[T]
}
