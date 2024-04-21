package ListLike

import (
	"fmt"
	"strings"
	"sync"

	itff "github.com/PlayerR9/MyGoLib/Common/Interfaces"
	itf "github.com/PlayerR9/MyGoLib/CustomData/Iterators"
	gen "github.com/PlayerR9/MyGoLib/Utility/General"
)

// LimitedSafeList is a generic type that represents a thread-safe list data
// structure with or without a maximum capacity, implemented using a linked list.
type LimitedSafeList[T any] struct {
	// front and back are pointers to the first and last nodes in the safe list,
	// respectively.
	front, back *listLinkedNode[T]

	// frontMutex and backMutex are sync.RWMutexes, which are used to ensure that
	// concurrent reads and writes to the front and back nodes are thread-safe.
	frontMutex, backMutex sync.RWMutex

	// size is the current number of elements in the list.
	size int

	// capacity is the maximum number of elements that the list can hold.
	capacity int
}

// NewLimitedSafeList is a function that creates and returns a new instance of a
// LimitedSafeList.
//
// Parameters:
//
//   - values: A variadic parameter of type T, which represents the initial values to
//     be stored in the list.
//
// Returns:
//
//   - *LimitedSafeList[T]: A pointer to the newly created LimitedSafeList.
func NewLimitedSafeList[T any](values ...T) *LimitedSafeList[T] {
	list := new(LimitedSafeList[T])

	if len(values) == 0 {
		return list
	}

	list.size = len(values)

	// First node
	list_node := &listLinkedNode[T]{
		value: values[0],
	}

	list.front = list_node
	list.back = list_node

	// Subsequent nodes
	for _, element := range values {
		list_node := &listLinkedNode[T]{
			value: element,
			prev:  list.back,
		}

		list.back.next = list_node
		list.back = list_node
	}

	return list
}

// Append is a method of the LimitedSafeList type. It is used to add an element to the
// end of the list.
//
// Panics with an error of type *ErrCallFailed if the list is fu
//
// Parameters:
//
//   - value: The value of type T to be added to the list.
func (list *LimitedSafeList[T]) Append(value T) error {
	list.backMutex.Lock()
	defer list.backMutex.Unlock()

	if list.size >= list.capacity {
		return NewErrFullList(list)
	}

	node := &listLinkedNode[T]{value: value}

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

	return nil
}

// DeleteFirst is a method of the LimitedSafeList type. It is used to remove and return
// the first element from the list.
//
// Panics with an error of type *ErrCallFailed if the list is empty.
//
// Returns:
//
//   - T: The first element in the list.
func (list *LimitedSafeList[T]) DeleteFirst() (T, error) {
	list.frontMutex.Lock()
	defer list.frontMutex.Unlock()

	if list.front == nil {
		return *new(T), NewErrEmptyList(list)
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

	return toRemove.value, nil
}

// PeekFirst is a method of the LimitedSafeList type. It is used to return the first
// element from the list without removing it.
//
// Panics with an error of type *ErrCallFailed if the list is empty.
//
// Returns:
//
//   - T: The first element in the list.
func (list *LimitedSafeList[T]) PeekFirst() (T, error) {
	list.frontMutex.RLock()
	defer list.frontMutex.RUnlock()

	if list.front == nil {
		return *new(T), NewErrEmptyList(list)
	}

	return list.front.value, nil
}

// IsEmpty is a method of the LimitedSafeList type. It checks if the list is empty.
//
// Returns:
//
//   - bool: A boolean value that is true if the list is empty, and false otherwise.
func (list *LimitedSafeList[T]) IsEmpty() bool {
	list.frontMutex.RLock()
	defer list.frontMutex.RUnlock()

	return list.front == nil
}

// Size is a method of the LimitedSafeList type. It returns the number of elements in the
// list.
//
// Returns:
//
//   - int: An integer that represents the number of elements in the list.
func (list *LimitedSafeList[T]) Size() int {
	list.frontMutex.RLock()
	defer list.frontMutex.RUnlock()

	list.backMutex.RLock()
	defer list.backMutex.RUnlock()

	return list.size
}

// Capacity is a method of the LimitedSafeList type. It returns the maximum number of
// elements that the list can hold.
//
// Returns:
//
//   - optional.Int: An optional integer that represents the maximum number of
//     elements the list can hold.
func (list *LimitedSafeList[T]) Capacity() (int, bool) {
	return list.capacity, true
}

// Iterator is a method of the LimitedSafeList type. It is used to return an iterator
// for the list.
// However, the iterator does not share the list's thread safety.
//
// Returns:
//
//   - itf.Iterater[T]: An iterator for the list.
func (list *LimitedSafeList[T]) Iterator() itf.Iterater[T] {
	list.frontMutex.RLock()
	defer list.frontMutex.RUnlock()

	list.backMutex.RLock()
	defer list.backMutex.RUnlock()

	var builder itf.Builder[T]

	for node := list.front; node != nil; node = node.next {
		builder.Append(node.value)
	}

	return builder.Build()
}

// Clear is a method of the LimitedSafeList type. It is used to remove all elements from
// the list.
func (list *LimitedSafeList[T]) Clear() {
	list.frontMutex.Lock()
	defer list.frontMutex.Unlock()

	list.backMutex.Lock()
	defer list.backMutex.Unlock()

	if list.front == nil {
		return // List is already empty
	}

	// 1. First node
	list.front.prev = nil
	prev := list.front

	// 2. Subsequent nodes
	for node := list.front.next; node != nil; node = node.next {
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

// IsFull is a method of the LimitedSafeList type. It checks if the list is fu
//
// Returns:
//
//   - isFull: A boolean value that is true if the list is full, and false otherwise.
func (list *LimitedSafeList[T]) IsFull() (isFull bool) {
	return list.capacity <= list.size
}

// String is a method of the LimitedSafeList type. It returns a string representation of
// the list including information about its size, capacity, and elements.
//
// Returns:
//
//   - string: A string representation of the list.
func (list *LimitedSafeList[T]) String() string {
	list.frontMutex.RLock()
	defer list.frontMutex.RUnlock()

	list.backMutex.RLock()
	defer list.backMutex.RUnlock()

	var builder strings.Builder

	fmt.Fprintf(&builder, "LimitedSafeList[capacity=%d, ", list.capacity)

	if list.front == nil {
		fmt.Fprintf(&builder, "size=0, values=[]]")

		return builder.String()
	}

	fmt.Fprintf(&builder, "size=%d, values=[%v", list.size, list.front.value)

	fmt.Fprintf(&builder, "%v", list.front.value)

	for node := list.front.next; node != nil; node = node.next {
		fmt.Fprintf(&builder, ", %v", node.value)
	}

	fmt.Fprintf(&builder, "]]")

	return builder.String()
}

// Prepend is a method of the LimitedSafeList type. It is used to add an element to the
// front of the list.
//
// Panics with an error of type *ErrCallFailed if the list is fu
//
// Parameters:
//
//   - value: The value of type T to be added to the list.
func (list *LimitedSafeList[T]) Prepend(value T) error {
	list.frontMutex.Lock()
	defer list.frontMutex.Unlock()

	if list.size >= list.capacity {
		return NewErrFullList(list)
	}

	node := &listLinkedNode[T]{value: value}

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

	return nil
}

// DeleteLast is a method of the LimitedSafeList type. It is used to remove and return the
// last element from the list.
//
// Panics with an error of type *ErrCallFailed if the list is empty.
//
// Returns:
//
//   - T: The last element in the list.
func (list *LimitedSafeList[T]) DeleteLast() (T, error) {
	list.backMutex.Lock()
	defer list.backMutex.Unlock()

	if list.back == nil {
		return *new(T), NewErrEmptyList(list)
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

	return toRemove.value, nil
}

// PeekLast is a method of the LimitedSafeList type. It is used to return the last element
// from the list without removing it.
//
// Panics with an error of type *ErrCallFailed if the list is empty.
//
// Returns:
//
//   - T: The last element in the list.
func (list *LimitedSafeList[T]) PeekLast() (T, error) {
	list.backMutex.RLock()
	defer list.backMutex.RUnlock()

	if list.back == nil {
		return *new(T), NewErrEmptyList(list)
	}

	return list.back.value, nil
}

// CutNilValues is a method of the LimitedSafeList type. It is used to remove all nil
// values from the list.
func (list *LimitedSafeList[T]) CutNilValues() {
	list.frontMutex.Lock()
	defer list.frontMutex.Unlock()

	list.backMutex.Lock()
	defer list.backMutex.Unlock()

	if list.front == nil {
		return // List is empty
	}

	if gen.IsNil(list.front.value) && list.front == list.back {
		// Single node
		list.front = nil
		list.back = nil
		list.size = 0

		return
	}

	var toDelete *listLinkedNode[T] = nil

	// 1. First node
	if gen.IsNil(list.front.value) {
		toDelete = list.front

		list.front = list.front.next
		list.front.prev = nil

		toDelete.next = nil
		list.size--
	}

	prev := list.front

	// 2. Subsequent nodes (except last)
	for node := list.front.next; node.next != nil; node = node.next {
		if !gen.IsNil(node.value) {
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
	if gen.IsNil(list.back.value) {
		list.back = prev
		list.back.next = nil
		list.size--
	}
}

// Slice is a method of the LimitedSafeList type. It is used to return a slice of the
// elements in the list.
//
// Returns:
//
//   - []T: A slice of type T containing the elements of the list.
func (list *LimitedSafeList[T]) Slice() []T {
	list.frontMutex.RLock()
	defer list.frontMutex.RUnlock()

	list.backMutex.RLock()
	defer list.backMutex.RUnlock()

	slice := make([]T, 0, list.size)

	for node := list.front; node != nil; node = node.next {
		slice = append(slice, node.value)
	}

	return slice
}

// Copy is a method of the LimitedSafeList type. It is used to create a shallow copy of
// the list.
//
// Returns:
//
//   - itff.Copier: A copy of the list.
func (list *LimitedSafeList[T]) Copy() itff.Copier {
	list.frontMutex.RLock()
	defer list.frontMutex.RUnlock()

	list.backMutex.RLock()
	defer list.backMutex.RUnlock()

	listCopy := &LimitedSafeList[T]{
		size:     list.size,
		capacity: list.capacity,
	}

	if list.front == nil {
		return listCopy
	}

	// First node
	listCopy.front = &listLinkedNode[T]{
		value: list.front.value,
	}

	prev := listCopy.front

	// Subsequent nodes
	for node := list.front.next; node != nil; node = node.next {
		nodeCopy := &listLinkedNode[T]{
			value: node.value,
			prev:  prev,
		}

		prev.next = nodeCopy
		prev = nodeCopy
	}

	if listCopy.front.next != nil {
		listCopy.front.next.prev = listCopy.front
	}

	listCopy.back = prev

	return listCopy
}
