package ListLike

import (
	List "github.com/PlayerR9/MyGoLib/CustomData/ListLike/List"
	Queue "github.com/PlayerR9/MyGoLib/CustomData/ListLike/Queue"
	Stack "github.com/PlayerR9/MyGoLib/CustomData/ListLike/Stack"
)

type ListLike[T any] interface {
	List.Lister[T]
	Queue.Queuer[T]
	Stack.Stacker[T]
}
