package ListLike

import (
	"fmt"

	itf "github.com/PlayerR9/MyGoLib/Interfaces"
	"github.com/markphelps/optional"
)

type ListLike[T any] interface {
	// The IsEmpty method checks if the list is empty and returns a boolean value
	// indicating whether it is empty or not.
	IsEmpty() bool

	// The Size method returns the number of elements currently in the list.
	Size() int

	// The Capacity method returns the maximum number of elements that the list can hold.
	Capacity() optional.Int

	// The Clear method is used to remove all elements from the list, making it empty.
	Clear()

	// The IsFull method checks if the list is full, meaning it has reached its maximum
	// capacity and cannot accept any more elements.
	IsFull() bool

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
