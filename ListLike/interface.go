package ListLike

import (
	"fmt"

	intf "github.com/PlayerR9/MyGoLib/Units/Common"
	itr "github.com/PlayerR9/MyGoLib/Units/Iterators"
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

	itr.Iterable[T]

	intf.Copier
	fmt.GoStringer
}
