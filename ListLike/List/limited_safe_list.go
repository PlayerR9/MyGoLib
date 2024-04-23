package List

import (
	"fmt"
	"strings"
	"sync"

	itf "github.com/PlayerR9/MyGoLib/CustomData/Iterators"
	"github.com/PlayerR9/MyGoLib/ListLike/Common"
	itff "github.com/PlayerR9/MyGoLib/Units/Interfaces"
	gen "github.com/PlayerR9/MyGoLib/Utility/General"
)

// LimitedSafeList is a generic type that represents a thread-safe list data
// structure with or without a maximum capacity, implemented using a linked list.
type LimitedSafeList[T any] struct {
	// front and back are pointers to the first and last nodes in the safe list,
	// respectively.
	front, back *Common.ListSafeNode[T]

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
	list_node := Common.NewListSafeNode(values[0])

	list.front = list_node
	list.back = list_node

	// Subsequent nodes
	for _, element := range values {
		list_node := Common.NewListSafeNode(element)

		list_node.SetPrev(list.back)

		list.back.SetNext(list_node)
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
		return Common.NewErrFullList(list)
	}

	node := Common.NewListSafeNode(value)

	if list.back != nil {
		list.back.SetNext(node)
		node.SetPrev(list.back)
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
		return *new(T), Common.NewErrEmptyList(list)
	}

	toRemove := list.front

	list.backMutex.Lock()

	list.front = list.front.Next()

	if list.front == nil {
		list.back = nil
	} else {
		list.front.SetPrev(nil)
	}

	list.backMutex.Unlock()

	list.size--

	toRemove.SetNext(nil)

	return toRemove.Value, nil
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
		return *new(T), Common.NewErrEmptyList(list)
	}

	return list.front.Value, nil
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
func (list *LimitedSafeList[T]) Capacity() int {
	return list.capacity
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

	for node := list.front; node != nil; node = node.Next() {
		builder.Append(node.Value)
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
	list.front.SetPrev(nil)
	prev := list.front

	// 2. Subsequent nodes
	for node := list.front.Next(); node != nil; node = node.Next() {
		node.SetPrev(nil)

		prev = node
		prev.SetNext(nil)
	}

	prev.SetNext(nil)

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

	fmt.Fprintf(&builder, "size=%d, values=[%v", list.size, list.front.Value)

	fmt.Fprintf(&builder, "%v", list.front.Value)

	for node := list.front.Next(); node != nil; node = node.Next() {
		fmt.Fprintf(&builder, ", %v", node.Value)
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
		return Common.NewErrFullList(list)
	}

	node := Common.NewListSafeNode(value)

	if list.front == nil {
		// The list is empty
		list.backMutex.Lock()
		list.back = node
		list.backMutex.Unlock()
	} else {
		node.SetNext(list.front)
		list.front.SetPrev(node)
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
		return *new(T), Common.NewErrEmptyList(list)
	}

	toRemove := list.back

	list.frontMutex.Lock()

	list.back = list.back.Prev()

	if list.back == nil {
		list.front = nil
	} else {
		list.back.SetNext(nil)
	}

	list.frontMutex.Unlock()

	list.size--

	toRemove.SetPrev(nil)

	return toRemove.Value, nil
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
		return *new(T), Common.NewErrEmptyList(list)
	}

	return list.back.Value, nil
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

	if gen.IsNil(list.front.Value) && list.front == list.back {
		// Single node
		list.front = nil
		list.back = nil
		list.size = 0

		return
	}

	var toDelete *Common.ListSafeNode[T] = nil

	// 1. First node
	if gen.IsNil(list.front.Value) {
		toDelete = list.front

		list.front = list.front.Next()
		list.front.SetPrev(nil)

		toDelete.SetNext(nil)
		list.size--
	}

	prev := list.front

	// 2. Subsequent nodes (except last)
	for node := list.front.Next(); node.Next() != nil; node = node.Next() {
		if !gen.IsNil(node.Value) {
			prev = node
		} else {
			prev.SetNext(node.Next())
			node.Next().SetPrev(prev)
			list.size--

			if toDelete != nil {
				toDelete.SetNext(nil)
			}

			toDelete = node
		}
	}

	if toDelete != nil {
		toDelete.SetNext(nil)
	}

	// 3. Last node
	if gen.IsNil(list.back.Value) {
		list.back = prev
		list.back.SetNext(nil)
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

	for node := list.front; node != nil; node = node.Next() {
		slice = append(slice, node.Value)
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
	listCopy.front = Common.NewListSafeNode(list.front.Value)

	prev := listCopy.front

	// Subsequent nodes
	for node := list.front.Next(); node != nil; node = node.Next() {
		nodeCopy := Common.NewListSafeNode(node.Value)
		nodeCopy.SetPrev(prev)

		prev.SetNext(nodeCopy)
		prev = nodeCopy
	}

	if listCopy.front.Next() != nil {
		listCopy.front.Next().SetPrev(listCopy.front)
	}

	listCopy.back = prev

	return listCopy
}
