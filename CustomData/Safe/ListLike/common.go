package ListLike

import (
	"fmt"

	itf "github.com/PlayerR9/MyGoLib/Interfaces"
)

type ListLike[T any] interface {
	// The IsEmpty method checks if the list is empty and returns a boolean value
	// indicating whether it is empty or not.
	IsEmpty() bool

	// The Size method returns the number of elements currently in the list.
	Size() int

	// The Clear method is used to remove all elements from the list, making it empty.
	Clear()

	// The String method returns a string representation of the list.
	// It is useful for debugging and logging purposes.
	fmt.Stringer

	// CutNilValues is a method that removes all nil values from the list.
	// It is useful for cleaning up the list and removing any empty or nil elements.
	CutNilValues()

	// The itf.Iterable interface is used to provide an iterator for the list.
	itf.Iterable[T]

	// The itf.Slicer interface is used to provide a slicer for the list.
	itf.Slicer[T]

	// The itf.Copier interface is used to provide a method for copying the list.
	itf.Copier
}

// queueLinkedNode represents a node in a linked list.
type queueLinkedNode[T any] struct {
	// value is the value stored in the node.
	value T

	// next is a pointer to the next queueLinkedNode in the list.
	next *queueLinkedNode[T]
}

// listLinkedNode represents a node in a linked list. It holds a value of a
// generic type and a reference to the next and previous nodes in the list.
type listLinkedNode[T any] struct {
	// The value stored in the node.
	value T

	// A reference to the previous and next nodes in the list, respectively.
	prev, next *listLinkedNode[T]
}
