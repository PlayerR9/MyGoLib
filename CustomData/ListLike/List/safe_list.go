package ListLike

import (
	"fmt"
	"strings"
	"sync"

	itf "github.com/PlayerR9/MyGoLib/Interfaces"
	ers "github.com/PlayerR9/MyGoLib/Utility/Errors"
	gen "github.com/PlayerR9/MyGoLib/Utility/General"
	"github.com/markphelps/optional"
)

// SafeList is a generic type that represents a thread-safe list data
// structure with or without a maximum capacity, implemented using a linked list.
type SafeList[T any] struct {
	// front and back are pointers to the first and last nodes in the safe list,
	// respectively.
	front, back *linkedNode[T]

	// frontMutex and backMutex are sync.RWMutexes, which are used to ensure that
	// concurrent reads and writes to the front and back nodes are thread-safe.
	frontMutex, backMutex sync.RWMutex

	// size is the current number of elements in the list.
	size int

	// capacity is the maximum number of elements that the list can hold.
	capacity optional.Int

	// capacityMutex is a sync.RWMutex, which is used to ensure that concurrent
	// reads and writes to the capacity are thread-safe.
	capacityMutex sync.RWMutex
}

// NewSafeList is a function that creates and returns a new instance of a
// SafeList.
//
// Parameters:
//
//   - values: A variadic parameter of type T, which represents the initial values to
//     be stored in the list.
//
// Returns:
//
//   - *SafeList[T]: A pointer to the newly created SafeList.
func NewSafeList[T any](values ...T) *SafeList[T] {
	list := new(SafeList[T])

	if len(values) == 0 {
		return list
	}

	list.size = len(values)

	// First node
	list_node := &linkedNode[T]{
		value: &values[0],
	}

	list.front = list_node
	list.back = list_node

	// Subsequent nodes
	for _, element := range values {
		list_node := &linkedNode[T]{
			value: &element,
			prev:  list.back,
		}

		list.back.next = list_node
		list.back = list_node
	}

	return list
}

// WithCapacity is a method of the SafeList type. It is used to set the maximum
// number of elements the list can hold.
//
// Panics with an error of type *ErrCallFailed if the capacity is already set
// or with an error of type *ErrInvalidParameter if the provided capacity is negative
// or less than the current number of elements in the list.
//
// Parameters:
//
//   - capacity: An integer that represents the maximum number of elements the list
//     can hold.
//
// Returns:
//
//   - Lister[T]: A pointer to the list with the new capacity set.
func (list *SafeList[T]) WithCapacity(capacity int) Lister[T] {
	defer ers.PropagatePanic(ers.NewErrCallFailed("WithCapacity", list.WithCapacity))

	list.capacityMutex.Lock()
	defer list.capacityMutex.Unlock()

	list.capacity.If(func(cap int) {
		panic(ers.NewErrInvalidParameter("capacity").
			WithReason(fmt.Errorf("capacity is already set with a value of %d", cap)),
		)
	})

	// This prevents the list from being modified while the capacity is being set
	list.frontMutex.RLock()
	defer list.frontMutex.RUnlock()

	list.backMutex.RLock()
	defer list.backMutex.RUnlock()

	if capacity < 0 {
		panic(ers.NewErrInvalidParameter("capacity").
			WithReason(fmt.Errorf("negative capacity (%d) is not allowed", capacity)),
		)
	} else if list.size > capacity {
		panic(ers.NewErrInvalidParameter("capacity").WithReason(
			fmt.Errorf("capacity (%d) is less than the current number of elements (%d)",
				capacity, list.size)),
		)
	}

	list.capacity = optional.NewInt(capacity)

	return list
}

// Append is a method of the SafeList type. It is used to add an element to the
// end of the list.
//
// Panics with an error of type *ErrCallFailed if the list is full.
//
// Parameters:
//
//   - value: The value of type T to be added to the list.
func (list *SafeList[T]) Append(value T) {
	list.backMutex.Lock()
	defer list.backMutex.Unlock()

	list.capacityMutex.RLock()
	defer list.capacityMutex.RUnlock()

	list.capacity.If(func(cap int) {
		if list.size >= cap {
			panic(ers.NewErrCallFailed("Append", list.Append).
				WithReason(NewErrFullList(list)))
		}
	})

	node := &linkedNode[T]{value: &value}

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

// DeleteFirst is a method of the SafeList type. It is used to remove and return
// the first element from the list.
//
// Panics with an error of type *ErrCallFailed if the list is empty.
//
// Returns:
//
//   - T: The first element in the list.
func (list *SafeList[T]) DeleteFirst() T {
	list.frontMutex.Lock()
	defer list.frontMutex.Unlock()

	list.capacityMutex.RLock()
	defer list.capacityMutex.RUnlock()

	if list.front == nil {
		panic(ers.NewErrCallFailed("DeleteFirst", list.DeleteFirst).
			WithReason(NewErrEmptyList(list)),
		)
	}

	toRemove := list.front

	list.backMutex.Lock()

	list.front = list.front.next

	if list.front == nil {
		list.back = nil
	} else {
		list.front.prev = nil
	}

	list.backMutex.Unlock()

	list.size--

	toRemove.next = nil

	return *toRemove.value
}

// PeekFirst is a method of the SafeList type. It is used to return the first
// element from the list without removing it.
//
// Panics with an error of type *ErrCallFailed if the list is empty.
//
// Returns:
//
//   - T: The first element in the list.
func (list *SafeList[T]) PeekFirst() T {
	list.frontMutex.RLock()
	defer list.frontMutex.RUnlock()

	if list.front != nil {
		return *list.front.value
	}

	panic(ers.NewErrCallFailed("PeekFirst", list.PeekFirst).
		WithReason(NewErrEmptyList(list)),
	)
}

// IsEmpty is a method of the SafeList type. It checks if the list is empty.
//
// Returns:
//
//   - bool: A boolean value that is true if the list is empty, and false otherwise.
func (list *SafeList[T]) IsEmpty() bool {
	list.frontMutex.RLock()
	defer list.frontMutex.RUnlock()

	return list.front == nil
}

// Size is a method of the SafeList type. It returns the number of elements in the
// list.
//
// Returns:
//
//   - int: An integer that represents the number of elements in the list.
func (list *SafeList[T]) Size() int {
	list.frontMutex.RLock()
	defer list.frontMutex.RUnlock()

	list.backMutex.RLock()
	defer list.backMutex.RUnlock()

	return list.size
}

// Capacity is a method of the SafeList type. It returns the maximum number of
// elements that the list can hold.
//
// Returns:
//
//   - optional.Int: An optional integer that represents the maximum number of
//     elements the list can hold.
func (list *SafeList[T]) Capacity() optional.Int {
	list.capacityMutex.RLock()
	defer list.capacityMutex.RUnlock()

	return list.capacity
}

// Iterator is a method of the SafeList type. It is used to return an iterator
// for the list.
// However, the iterator does not share the list's thread safety.
//
// Returns:
//
//   - itf.Iterater[T]: An iterator for the list.
func (list *SafeList[T]) Iterator() itf.Iterater[T] {
	list.frontMutex.RLock()
	defer list.frontMutex.RUnlock()

	list.backMutex.RLock()
	defer list.backMutex.RUnlock()

	var builder itf.Builder[T]

	for node := list.front; node != nil; node = node.next {
		builder.Append(*node.value)
	}

	return builder.Build()
}

// Clear is a method of the SafeList type. It is used to remove all elements from
// the list.
func (list *SafeList[T]) Clear() {
	list.frontMutex.Lock()
	defer list.frontMutex.Unlock()

	list.backMutex.Lock()
	defer list.backMutex.Unlock()

	list.capacityMutex.RLock()
	defer list.capacityMutex.RUnlock()

	if list.front == nil {
		return // List is already empty
	}

	// 1. First node
	list.front.value = nil
	list.front.prev = nil
	prev := list.front

	// 2. Subsequent nodes
	for node := list.front.next; node != nil; node = node.next {
		node.value = nil
		node.prev = nil

		prev = node
		prev.next = nil
	}

	prev.next = nil

	// 3. Reset list fields
	list.front = nil
	list.back = nil
	list.size = 0
}

// IsFull is a method of the SafeList type. It checks if the list is full.
//
// Returns:
//
//   - isFull: A boolean value that is true if the list is full, and false otherwise.
func (list *SafeList[T]) IsFull() (isFull bool) {
	list.capacityMutex.RLock()
	defer list.capacityMutex.RUnlock()

	list.capacity.If(func(cap int) {
		isFull = cap <= list.size
	})

	return
}

// String is a method of the SafeList type. It returns a string representation of
// the list including information about its size, capacity, and elements.
//
// Returns:
//
//   - string: A string representation of the list.
func (list *SafeList[T]) String() string {
	list.frontMutex.RLock()
	defer list.frontMutex.RUnlock()

	list.backMutex.RLock()
	defer list.backMutex.RUnlock()

	var builder strings.Builder

	builder.WriteString("SafeList[")

	list.capacity.If(func(cap int) {
		fmt.Fprintf(&builder, "capacity=%d, ", cap)
	})

	if list.front == nil {
		fmt.Fprintf(&builder, "size=0, values=[]]")

		return builder.String()
	}

	fmt.Fprintf(&builder, "size=%d, values=[%v", list.size, *list.front.value)

	fmt.Fprintf(&builder, "%v", *list.front.value)

	for node := list.front.next; node != nil; node = node.next {
		fmt.Fprintf(&builder, ", %v", *node.value)
	}

	fmt.Fprintf(&builder, "]]")

	return builder.String()
}

// Prepend is a method of the SafeList type. It is used to add an element to the
// front of the list.
//
// Panics with an error of type *ErrCallFailed if the list is full.
//
// Parameters:
//
//   - value: The value of type T to be added to the list.
func (list *SafeList[T]) Prepend(value T) {
	list.frontMutex.Lock()
	defer list.frontMutex.Unlock()

	list.capacityMutex.RLock()
	defer list.capacityMutex.RUnlock()

	list.capacity.If(func(cap int) {
		if list.size >= cap {
			panic(ers.NewErrCallFailed("Prepend", list.Prepend).
				WithReason(NewErrFullList(list)),
			)
		}
	})

	node := &linkedNode[T]{value: &value}

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

// DeleteLast is a method of the SafeList type. It is used to remove and return the
// last element from the list.
//
// Panics with an error of type *ErrCallFailed if the list is empty.
//
// Returns:
//
//   - T: The last element in the list.
func (list *SafeList[T]) DeleteLast() T {
	list.backMutex.Lock()
	defer list.backMutex.Unlock()

	list.capacityMutex.RLock()
	defer list.capacityMutex.RUnlock()

	if list.back == nil {
		panic(ers.NewErrCallFailed("DeleteLast", list.DeleteLast).
			WithReason(NewErrEmptyList(list)),
		)
	}

	toRemove := list.back

	list.frontMutex.Lock()

	list.back = list.back.prev

	if list.back == nil {
		list.front = nil
	} else {
		list.back.next = nil
	}

	list.frontMutex.Unlock()

	list.size--

	toRemove.prev = nil

	return *toRemove.value
}

// PeekLast is a method of the SafeList type. It is used to return the last element
// from the list without removing it.
//
// Panics with an error of type *ErrCallFailed if the list is empty.
//
// Returns:
//
//   - T: The last element in the list.
func (list *SafeList[T]) PeekLast() T {
	list.backMutex.RLock()
	defer list.backMutex.RUnlock()

	if list.back != nil {
		return *list.back.value

	}

	panic(ers.NewErrCallFailed("PeekLast", list.PeekLast).
		WithReason(NewErrEmptyList(list)),
	)
}

// CutNilValues is a method of the SafeList type. It is used to remove all nil
// values from the list.
func (list *SafeList[T]) CutNilValues() {
	list.frontMutex.Lock()
	defer list.frontMutex.Unlock()

	list.backMutex.Lock()
	defer list.backMutex.Unlock()

	list.capacityMutex.RLock()
	defer list.capacityMutex.RUnlock()

	if list.front == nil {
		return // List is empty
	}

	if gen.IsNil(*list.front.value) && list.front == list.back {
		// Single node
		list.front = nil
		list.back = nil
		list.size = 0

		return
	}

	var toDelete *linkedNode[T] = nil

	// 1. First node
	if gen.IsNil(*list.front.value) {
		toDelete = list.front

		list.front = list.front.next
		list.front.prev = nil

		toDelete.next = nil
		list.size--
	}

	prev := list.front

	// 2. Subsequent nodes (except last)
	for node := list.front.next; node.next != nil; node = node.next {
		if !gen.IsNil(*node.value) {
			prev = node
		} else {
			prev.next = node.next
			node.next.prev = prev
			list.size--

			if toDelete != nil {
				toDelete.next = nil
			}

			toDelete = node
		}
	}

	if toDelete != nil {
		toDelete.next = nil
	}

	// 3. Last node
	if list.back.value == nil {
		list.back = prev
		list.back.next = nil
		list.size--
	}
}

// Slice is a method of the SafeList type. It is used to return a slice of the
// elements in the list.
//
// Returns:
//
//   - []T: A slice of type T containing the elements of the list.
func (list *SafeList[T]) Slice() []T {
	list.frontMutex.RLock()
	defer list.frontMutex.RUnlock()

	list.backMutex.RLock()
	defer list.backMutex.RUnlock()

	slice := make([]T, 0, list.size)

	for node := list.front; node != nil; node = node.next {
		slice = append(slice, *node.value)
	}

	return slice
}
