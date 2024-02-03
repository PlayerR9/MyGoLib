package ListLike

import (
	"fmt"
	"strings"
	"sync"

	ers "github.com/PlayerR9/MyGoLib/Utility/Errors"
)

// LimitedSafeList is a generic type in Go that represents a thread-safe list data
// structure implemented using a linked list.
type LimitedSafeList[T any] struct {
	// front and back are pointers to the first and last nodes in the safe list,
	// respectively.
	front, back *linkedNode[T]

	// frontMutex and backMutex are sync.RWMutexes, which are used to ensure that
	// concurrent reads and writes to the front and back nodes are thread-safe.
	frontMutex, backMutex sync.RWMutex

	// size is the current number of elements in the list.
	size int

	// capacity is the maximum number of elements that the list can hold.
	capacity int
}

// NewLimitedSafeList is a function that creates and returns a new instance of a LimitedSafeList.
// It takes a variadic parameter of type T, which represents the initial values to be
// stored in the list.
//
// If no initial values are provided, the function simply returns a new LimitedSafeList with
// all its fields set to their zero values.
//
// If initial values are provided, the function creates a new LimitedSafeList and initializes
// its size. It then creates a linked list of linkedNodes from the initial values, with
// each node holding one value, and sets the front and back pointers of the list.
// The new LimitedSafeList is then returned.
func NewLimitedSafeList[T any](capacity int, values ...*T) *LimitedSafeList[T] {
	if capacity <= 0 {
		panic(ers.NewErrInvalidParameter(
			"capacity", fmt.Errorf("negative capacity (%d) is not allowed", capacity),
		))
	} else if len(values) > capacity {
		panic(ers.NewErrInvalidParameter(
			"values", fmt.Errorf("number of values (%d) exceeds the provided capacity (%d)",
				len(values), capacity),
		))
	}

	if len(values) == 0 {
		return new(LimitedSafeList[T])
	}

	list := new(LimitedSafeList[T])
	list.size = len(values)

	// First node
	node := &linkedNode[T]{value: values[0]}

	list.front = node
	list.back = node

	// Subsequent nodes
	for _, element := range values[1:] {
		node = &linkedNode[T]{
			value: element,
			prev:  list.back,
		}

		list.back.next = node
		list.back = node
	}

	return list
}

func (list *LimitedSafeList[T]) Append(value *T) {
	list.backMutex.Lock()
	defer list.backMutex.Unlock()

	list.frontMutex.RLock()
	if list.size >= list.capacity {
		list.frontMutex.RUnlock()

		panic(ers.NewErrOperationFailed(
			"append element", NewErrFullList(list),
		))
	} else {
		list.frontMutex.RUnlock()
	}

	node := &linkedNode[T]{value: value}

	if list.back != nil {
		list.back.next = node
		node.prev = list.back
	} else {
		// The list is empty
		list.frontMutex.Lock()
		list.front = node
		list.frontMutex.Unlock()
	}

	list.back = node

	list.size++
}

func (list *LimitedSafeList[T]) DeleteFirst() *T {
	list.frontMutex.Lock()
	defer list.frontMutex.Unlock()

	list.backMutex.Lock()
	if list.front == nil {
		list.backMutex.Unlock()

		panic(ers.NewErrOperationFailed(
			"delete first element", NewErrEmptyList(list),
		))
	} else {
		list.backMutex.Unlock()
	}

	var value *T

	value, list.front = list.front.value, list.front.next
	if list.front == nil {
		// The list has only one element
		list.backMutex.Lock()
		list.back = nil
		list.backMutex.Unlock()
	} else {
		list.front.prev = nil
	}

	list.size--

	return value
}

func (list *LimitedSafeList[T]) PeekFirst() *T {
	list.frontMutex.RLock()
	defer list.frontMutex.RUnlock()

	if list.front == nil {
		panic(ers.NewErrOperationFailed(
			"peek first element", NewErrEmptyList(list),
		))
	}

	return list.front.value
}

func (list *LimitedSafeList[T]) IsEmpty() bool {
	list.frontMutex.RLock()
	defer list.frontMutex.RUnlock()

	return list.front == nil
}

func (list *LimitedSafeList[T]) Size() int {
	// Lock the front and back nodes to ensure that
	// edit operations are not performed while the size is being read.
	list.frontMutex.RLock()
	defer list.frontMutex.RUnlock()

	list.backMutex.RLock()
	defer list.backMutex.RUnlock()

	return list.size
}

func (list *LimitedSafeList[T]) ToSlice() []*T {
	// Lock everything to ensure that no other goroutine can modify the list while
	// it is being converted to a slice.
	list.frontMutex.RLock()
	defer list.frontMutex.RUnlock()

	list.backMutex.RLock()
	defer list.backMutex.RUnlock()

	slice := make([]*T, 0, list.size)

	// Add the values to the slice
	for node := list.front; node != nil; node = node.next {
		slice = append(slice, node.value)
	}

	return slice
}

func (list *LimitedSafeList[T]) Clear() {
	// Lock everything to ensure that no other goroutine can modify the list while
	// it is being cleared.
	list.frontMutex.RLock()
	defer list.frontMutex.RUnlock()

	if list.front == nil {
		return // nothing to clear
	}

	list.backMutex.Lock()
	defer list.backMutex.Unlock()

	// 1. First node
	node := list.front
	node.value = nil

	// 2. Subsequent nodes
	for node = node.next; node != nil; node = node.next {
		node.value = nil
		node.prev.next = nil
		node.prev = nil
	}

	// Reset the list's fields
	list.front = nil
	list.back = nil
	list.size = 0
}

func (list *LimitedSafeList[T]) IsFull() bool {
	list.frontMutex.RLock()
	defer list.frontMutex.RUnlock()

	list.backMutex.RLock()
	defer list.backMutex.RUnlock()

	return list.size >= list.capacity
}

func (list *LimitedSafeList[T]) String() string {
	list.frontMutex.RLock()
	defer list.frontMutex.RUnlock()
	list.backMutex.RLock()
	defer list.backMutex.RUnlock()

	var builder strings.Builder

	fmt.Fprintf(&builder, "LimitedSafeList[size=%d, capacity=%d, values=[", list.size, list.capacity)

	if list.front != nil {
		fmt.Fprintf(&builder, "%v", *list.front.value)

		for node := list.front.next; node != nil; node = node.next {
			fmt.Fprintf(&builder, ", %v", *node.value)
		}
	}

	fmt.Fprintf(&builder, "]]")

	return builder.String()
}

// The Prepend method adds a value of type T to the end of the list.
func (list *LimitedSafeList[T]) Prepend(value *T) {
	list.frontMutex.Lock()
	defer list.frontMutex.Unlock()

	list.backMutex.RLock()
	if list.size >= list.capacity {
		list.backMutex.RUnlock()

		panic(ers.NewErrOperationFailed(
			"prepend element", NewErrFullList(list),
		))
	} else {
		list.backMutex.RUnlock()
	}

	node := &linkedNode[T]{value: value}

	if list.front == nil {
		// The list is empty
		list.backMutex.Lock()
		list.back = node
		list.backMutex.Unlock()
	} else {
		node.next = list.front
		list.front.prev = node
	}

	list.front = node

	list.size++
}

// The DeleteLast method is a convenience method that deletelasts an element from the list
// and returns it.
// If the list is empty, it will panic.
func (list *LimitedSafeList[T]) DeleteLast() *T {
	list.backMutex.Lock()
	defer list.backMutex.Unlock()

	if list.back == nil {
		panic(ers.NewErrOperationFailed(
			"delete last element", NewErrEmptyList(list),
		))
	}

	value := list.back.value

	if list.back.prev == nil {
		// The list has only one element
		list.frontMutex.Lock()
		list.front = nil
		list.frontMutex.Unlock()
	} else {
		list.back.prev.next = nil
	}

	list.back = list.back.prev

	list.size--

	return value
}

// PeekLast is a method that returns the value at the front of the list without removing
// it.
// If the list is empty, it will panic.
func (list *LimitedSafeList[T]) PeekLast() *T {
	list.backMutex.RLock()
	defer list.backMutex.RUnlock()

	if list.back == nil {
		panic(ers.NewErrOperationFailed(
			"peek last element", NewErrEmptyList(list),
		))
	}

	return list.back.value
}
