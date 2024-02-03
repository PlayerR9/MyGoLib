package ListLike

import (
	"fmt"
	"strings"
	"sync"

	ers "github.com/PlayerR9/MyGoLib/Utility/Errors"
)

// SafeList is a generic type in Go that represents a thread-safe list data
// structure implemented using a linked list.
type SafeList[T any] struct {
	// front and back are pointers to the first and last nodes in the safe list,
	// respectively.
	front, back *linkedNode[T]

	// frontMutex and backMutex are sync.RWMutexes, which are used to ensure that
	// concurrent reads and writes to the front and back nodes are thread-safe.
	frontMutex, backMutex sync.RWMutex

	// size is the current number of elements in the list.
	size int
}

// NewSafeList is a function that creates and returns a new instance of a SafeList.
// It takes a variadic parameter of type T, which represents the initial values to be
// stored in the list.
//
// If no initial values are provided, the function simply returns a new SafeList with
// all its fields set to their zero values.
//
// If initial values are provided, the function creates a new SafeList and initializes
// its size. It then creates a linked list of linkedNodes from the initial values, with
// each node holding one value, and sets the front and back pointers of the list.
// The new SafeList is then returned.
func NewSafeList[T any](values ...*T) *SafeList[T] {
	if len(values) == 0 {
		return new(SafeList[T])
	}

	list := new(SafeList[T])
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

func (list *SafeList[T]) Append(value *T) {
	node := &linkedNode[T]{value: value}

	list.backMutex.Lock()
	defer list.backMutex.Unlock()

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

func (list *SafeList[T]) DeleteFirst() *T {
	if list.IsEmpty() {
		panic(ers.NewErrOperationFailed(
			"delete first element", NewErrEmptyList(list),
		))
	}

	list.frontMutex.Lock()
	defer list.frontMutex.Unlock()

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

func (list *SafeList[T]) PeekFirst() *T {
	list.frontMutex.RLock()
	defer list.frontMutex.RUnlock()

	if list.front == nil {
		panic(ers.NewErrOperationFailed(
			"peek first element", NewErrEmptyList(list),
		))
	}

	return list.front.value
}

func (list *SafeList[T]) IsEmpty() bool {
	list.frontMutex.RLock()
	defer list.frontMutex.RUnlock()

	return list.front == nil
}

func (list *SafeList[T]) Size() int {
	// Lock the front and back nodes to ensure that
	// edit operations are not performed while the size is being read.
	list.frontMutex.RLock()
	defer list.frontMutex.RUnlock()

	list.backMutex.RLock()
	defer list.backMutex.RUnlock()

	return list.size
}

func (list *SafeList[T]) ToSlice() []*T {
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

func (list *SafeList[T]) Clear() {
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

func (list *SafeList[T]) IsFull() bool {
	return false
}

func (list *SafeList[T]) String() string {
	list.frontMutex.RLock()
	defer list.frontMutex.RUnlock()
	list.backMutex.RLock()
	defer list.backMutex.RUnlock()

	var builder strings.Builder

	fmt.Fprintf(&builder, "SafeList[size=%d, values=[", list.size)

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
func (list *SafeList[T]) Prepend(value *T) {
	list.frontMutex.Lock()
	defer list.frontMutex.Unlock()

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
func (list *SafeList[T]) DeleteLast() *T {
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
func (list *SafeList[T]) PeekLast() *T {
	list.backMutex.RLock()
	defer list.backMutex.RUnlock()

	if list.back == nil {
		panic(ers.NewErrOperationFailed(
			"peek last element", NewErrEmptyList(list),
		))
	}

	return list.back.value
}
