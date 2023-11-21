package Queue

import (
	"fmt"
	"strings"

	gen "github.com/PlayerR9/MyGoLib/Utility/General"
)

type ErrEmptyQueue struct{}

func (e ErrEmptyQueue) Error() string {
	return "Empty queue"
}

type ErrFullQueue struct{}

func (e ErrFullQueue) Error() string {
	return "Full queue"
}

var (
	QueueHead string = "← | "
	QueueSep  string = " | "
)

type node[T any] struct {
	value T
	next  *node[T]
}

type Queue[T any] interface {
	Enqueue(value T) Queue[T]
	Dequeue() (T, Queue[T])
	Peek() T
	IsEmpty() bool
	Size() int
	ToSlice() []T
	DeepCopy() Queue[T]
	Clear() Queue[T]
}

type LinkedQueue[T any] struct {
	front *node[T]
	back  *node[T]
	size  int
}

func NewLinkedQueue[T any](elements ...T) (queue LinkedQueue[T]) {
	for _, element := range elements {
		queue = queue.Enqueue(element).(LinkedQueue[T])
	}

	return
}

func (queue LinkedQueue[T]) Enqueue(value T) Queue[T] {
	var node node[T]

	node.value = gen.DeepCopy(value).(T)

	if queue.back == nil {
		queue.front = &node
	} else {
		queue.back.next = &node
	}

	queue.back = &node

	queue.size++

	return queue
}

func (queue LinkedQueue[T]) Dequeue() (T, Queue[T]) {
	if queue.front == nil {
		panic(ErrEmptyQueue{})
	}

	value := queue.front.value
	queue.front = queue.front.next

	if queue.front == nil {
		queue.back = nil
	}

	queue.size--

	return value, queue
}

func (queue LinkedQueue[T]) Peek() T {
	if queue.front == nil {
		panic(ErrEmptyQueue{})
	}

	return gen.DeepCopy(queue.front.value).(T)
}

func (queue LinkedQueue[T]) IsEmpty() bool {
	return queue.front == nil
}

func (queue LinkedQueue[T]) Size() int {
	return queue.size
}

func (queue LinkedQueue[T]) ToSlice() (slice []T) {
	for node := queue.front; node != nil; node = node.next {
		slice = append(slice, gen.DeepCopy(node.value).(T))
	}

	return
}

func (queue LinkedQueue[T]) String() string {
	if queue.front == nil {
		return QueueHead
	}

	var builder strings.Builder

	builder.WriteString(QueueHead)

	builder.WriteString(fmt.Sprintf("%v", queue.front.value))

	for node := queue.front.next; node != nil; node = node.next {
		builder.WriteString(fmt.Sprintf("%s%v", QueueSep, node.value))
	}

	return builder.String()
}

func (queue LinkedQueue[T]) DeepCopy() (other Queue[T]) {
	for node := queue.front; node != nil; node = node.next {
		other = other.Enqueue(gen.DeepCopy(node.value).(T))
	}

	return
}

func (queue LinkedQueue[T]) Clear() Queue[T] {
	return LinkedQueue[T]{
		front: nil,
		back:  nil,
		size:  0,
	}
}

type ArrayQueue[T any] struct {
	elements   []T
	hasMaxSize bool
}

func NewArrayQueue[T any](elements ...T) ArrayQueue[T] {
	other := make([]T, len(elements))

	for i, element := range elements {
		other[i] = gen.DeepCopy(element).(T)
	}

	return ArrayQueue[T]{other, false}
}

func NewArrayQueueWithMaxSize[T any](maxSize int, elements ...T) ArrayQueue[T] {
	if maxSize < 0 {
		panic(fmt.Errorf("maxSize cannot be negative"))
	}

	if len(elements) > maxSize {
		panic(fmt.Errorf("maxSize is smaller than the number of elements"))
	}

	others := make([]T, len(elements), maxSize)

	for i, element := range elements {
		others[i] = gen.DeepCopy(element).(T)
	}

	return ArrayQueue[T]{others, true}
}

func (queue ArrayQueue[T]) Enqueue(value T) Queue[T] {
	if queue.hasMaxSize && len(queue.elements) == cap(queue.elements) {
		panic(ErrFullQueue{})
	}

	queue.elements = append(queue.elements, gen.DeepCopy(value).(T))

	return queue
}

func (queue ArrayQueue[T]) Dequeue() (T, Queue[T]) {
	if len(queue.elements) == 0 {
		panic(ErrEmptyQueue{})
	}

	value := queue.elements[0]
	queue.elements = queue.elements[1:]

	return value, queue
}

func (queue ArrayQueue[T]) Peek() T {
	if len(queue.elements) == 0 {
		panic(ErrEmptyQueue{})
	}

	return gen.DeepCopy(queue.elements[0]).(T)
}

func (queue ArrayQueue[T]) IsEmpty() bool {
	return len(queue.elements) == 0
}

func (queue ArrayQueue[T]) Size() int {
	return len(queue.elements)
}

func (queue ArrayQueue[T]) ToSlice() []T {
	slice := make([]T, len(queue.elements))

	for i, element := range queue.elements {
		slice[i] = gen.DeepCopy(element).(T)
	}

	return slice
}

func (queue ArrayQueue[T]) String() string {
	if len(queue.elements) == 0 {
		return QueueHead
	}

	var builder strings.Builder

	builder.WriteString(QueueHead)

	builder.WriteString(fmt.Sprintf("%v", queue.elements[0]))

	for _, element := range queue.elements[1:] {
		builder.WriteString(QueueSep)
		builder.WriteString(fmt.Sprintf("%v", element))
	}

	return builder.String()
}

func (queue ArrayQueue[T]) DeepCopy() Queue[T] {
	obj := ArrayQueue[T]{
		hasMaxSize: queue.hasMaxSize,
	}

	if queue.hasMaxSize {
		obj.elements = make([]T, 0, cap(queue.elements))
	} else {
		obj.elements = make([]T, 0)
	}

	for _, element := range queue.elements {
		obj.elements = append(obj.elements, gen.DeepCopy(element).(T))
	}

	return obj
}

func (queue ArrayQueue[T]) Clear() Queue[T] {
	other := ArrayQueue[T]{
		hasMaxSize: queue.hasMaxSize,
	}

	if queue.hasMaxSize {
		other.elements = make([]T, 0, cap(queue.elements))
	} else {
		other.elements = make([]T, 0)
	}

	return other
}
