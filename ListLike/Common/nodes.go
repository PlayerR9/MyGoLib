package Common

// ListNode represents a node in a linked list. It holds a value of a generic type
// and a reference to the next node in the list.
type ListNode[T any] struct {
	// The Value stored in the node.
	Value T

	// A reference to the previous and next nodes in the list, respectively.
	prev, next *ListNode[T]
}

// NewListNode creates a new LinkedNode with the given value.
func NewListNode[T any](value T) *ListNode[T] {
	return &ListNode[T]{Value: value}
}

func (node *ListNode[T]) SetNext(next *ListNode[T]) {
	node.next = next
}

func (node *ListNode[T]) Next() *ListNode[T] {
	return node.next
}

func (node *ListNode[T]) SetPrev(prev *ListNode[T]) {
	node.prev = prev
}

func (node *ListNode[T]) Prev() *ListNode[T] {
	return node.prev
}

// QueueNode represents a node in a linked list.
type QueueNode[T any] struct {
	// Value is the value stored in the node.
	Value T

	// next is a pointer to the next linkedNode in the list.
	next *QueueNode[T]
}

// NewQueueNode creates a new LinkedNode with the given value.
func NewQueueNode[T any](value T) *QueueNode[T] {
	return &QueueNode[T]{
		Value: value,
	}
}

func (node *QueueNode[T]) SetNext(next *QueueNode[T]) {
	node.next = next
}

func (node *QueueNode[T]) Next() *QueueNode[T] {
	return node.next
}

// StackNode represents a node in a linked list.
type StackNode[T any] struct {
	// value is the value stored in the node.
	Value T

	// next is a pointer to the next linkedNode in the list.
	next *StackNode[T]
}

func NewStackNode[T any](value T) *StackNode[T] {
	return &StackNode[T]{
		Value: value,
	}
}

func (node *StackNode[T]) SetNext(next *StackNode[T]) {
	node.next = next
}

func (node *StackNode[T]) Next() *StackNode[T] {
	return node.next
}

// QueueSafeNode represents a node in a linked list.
type QueueSafeNode[T any] struct {
	// Value is the Value stored in the node.
	Value T

	// next is a pointer to the next queueLinkedNode in the list.
	next *QueueSafeNode[T]
}

// NewQueueSafeNode creates a new QueueSafeNode with the given value.
func NewQueueSafeNode[T any](value T) *QueueSafeNode[T] {
	return &QueueSafeNode[T]{Value: value}
}

// SetNext sets the next node in the list.
func (node *QueueSafeNode[T]) SetNext(next *QueueSafeNode[T]) {
	node.next = next
}

// Next returns the next node in the list.
func (node *QueueSafeNode[T]) Next() *QueueSafeNode[T] {
	return node.next
}

// ListSafeNode represents a node in a linked list. It holds a value of a
// generic type and a reference to the next and previous nodes in the list.
type ListSafeNode[T any] struct {
	// The Value stored in the node.
	Value T

	// A reference to the previous and next nodes in the list, respectively.
	prev, next *ListSafeNode[T]
}

// NewListSafeNode creates a new ListSafeNode with the given value.
func NewListSafeNode[T any](value T) *ListSafeNode[T] {
	return &ListSafeNode[T]{Value: value}
}

// SetNext sets the next node in the list.
func (node *ListSafeNode[T]) SetNext(next *ListSafeNode[T]) {
	node.next = next
}

// Next returns the next node in the list.
func (node *ListSafeNode[T]) Next() *ListSafeNode[T] {
	return node.next
}

// SetPrev sets the previous node in the list.
func (node *ListSafeNode[T]) SetPrev(prev *ListSafeNode[T]) {
	node.prev = prev
}

// Prev returns the previous node in the list.
func (node *ListSafeNode[T]) Prev() *ListSafeNode[T] {
	return node.prev
}
