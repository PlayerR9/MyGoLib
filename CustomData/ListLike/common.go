package ListLike

import (
	List "github.com/PlayerR9/MyGoLib/CustomData/ListLike/List"
	Queue "github.com/PlayerR9/MyGoLib/CustomData/ListLike/Queue"
	Stack "github.com/PlayerR9/MyGoLib/CustomData/ListLike/Stack"
)

type ListLike[T any] interface {
	// WithCapacity is a special function that modifies an existing list-like data
	// structure to have a specific capacity. Panics if the list already has a capacity
	// set or if the new capacity is less than the current size of the list-like data
	// structure.
	//
	// As a result, it is recommended to use this function only when creating a new
	// list-like data structure.
	WithCapacity(int) ListLike[T]

	List.Lister[T]
	Queue.Queuer[T]
	Stack.Stacker[T]
}
