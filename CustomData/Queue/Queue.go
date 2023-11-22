package Queue

import (
	"fmt"
	"strings"

	gen "github.com/PlayerR9/MyGoLib/Utility/General"
	"github.com/markphelps/optional"
)

const (
	// Queue Implementation Types

	// linked_queue with no maximum size Implementation
	LINKED int = iota

	// ArrayQueue with no maximum size Implementation
	ARRAY
)

type Queue[T any] struct {
	data           any
	implementation int

	size     int
	capacity optional.Int
}

func NewQueue[T any](implementation int, values ...T) Queue[T] {
	var queue Queue[T]

	switch implementation {
	case LINKED:
		queue.data = linked_queue[T]{
			front: nil,
			back:  nil,
		}
	case ARRAY:
		queue.data = make([]T, 0)
	default:
		panic("Invalid queue implementation type")
	}

	queue.implementation = implementation

	queue.size = 0
	queue.capacity = optional.Int{}

	for _, element := range values {
		queue.Enqueue(element)
	}

	return queue
}

func NewLimitedQueue[T any](implementation, capacity int, values ...T) Queue[T] {
	if capacity < 0 {
		panic("Cannot specify a negative capacity for a queue")
	}

	if len(values) > capacity {
		panic("Cannot specify more values than the capacity of the queue")
	}

	var queue Queue[T]

	switch implementation {
	case LINKED:
		queue.data = linked_queue[T]{
			front: nil,
			back:  nil,
		}
	case ARRAY:
		queue.data = make([]T, 0, capacity)
	default:
		panic("Invalid queue implementation type")
	}

	queue.implementation = implementation

	queue.size = 0
	queue.capacity = optional.NewInt(capacity)

	for _, element := range values {
		queue.Enqueue(element)
	}

	return queue
}

func (queue *Queue[T]) Enqueue(value T) {
	if queue.capacity.Present() && queue.size >= queue.capacity.MustGet() {
		panic(ErrFullQueue{})
	}

	switch queue.implementation {
	case LINKED:
		new_node := queue_node[T]{
			value: gen.DeepCopy(value).(T),
		}

		tmp := queue.data.(linked_queue[T])

		if tmp.back == nil {
			tmp.front = &new_node
		} else {
			tmp.back.next = &new_node
		}

		tmp.back = &new_node

		queue.data = tmp
	case ARRAY:
		tmp := queue.data.([]T)
		tmp = append(tmp, gen.DeepCopy(value).(T))
		queue.data = tmp
	}

	queue.size++
}

func (queue *Queue[T]) Dequeue() T {
	if queue.size <= 0 {
		panic(ErrEmptyQueue{})
	}

	var value T

	switch queue.implementation {
	case LINKED:
		tmp := queue.data.(linked_queue[T])

		value = tmp.front.value

		tmp.front = tmp.front.next

		if tmp.front == nil {
			tmp.back = nil
		}

		queue.data = tmp
	case ARRAY:
		tmp := queue.data.([]T)
		value = tmp[0]
		queue.data = tmp[1:]
	}

	queue.size--

	return value
}

func (queue Queue[T]) Peek() T {
	if queue.size == 0 {
		panic(ErrEmptyQueue{})
	}

	var value T

	switch queue.implementation {
	case LINKED:
		value = gen.DeepCopy(queue.data.(linked_queue[T]).front.value).(T)
	case ARRAY:
		value = gen.DeepCopy(queue.data.([]T)[0]).(T)
	}

	return value
}

func (queue Queue[T]) IsEmpty() bool {
	return queue.size == 0
}

func (queue Queue[T]) Size() int {
	return queue.size
}

func (queue *Queue[T]) ToSlice() []T {
	slice := make([]T, queue.size)

	switch queue.implementation {
	case LINKED:
		i := 0

		for queue_node := queue.data.(linked_queue[T]).front; queue_node != nil; queue_node = queue_node.next {
			slice[i] = gen.DeepCopy(queue_node.value).(T)
			i++
		}
	case ARRAY:
		for i, element := range queue.data.([]T) {
			slice[i] = gen.DeepCopy(element).(T)
		}
	}

	return slice
}

func (queue *Queue[T]) Clear() {
	switch queue.implementation {
	case LINKED:
		queue.data = linked_queue[T]{
			front: nil,
			back:  nil,
		}
	case ARRAY:
		if queue.capacity.Present() {
			queue.data = make([]T, 0, queue.capacity.MustGet())
		} else {
			queue.data = make([]T, 0)
		}
	}

	queue.size = 0
}

func (queue Queue[T]) IsFull() bool {
	return queue.capacity.Present() && queue.size == queue.capacity.MustGet()
}

func (queue *Queue[T]) String() string {
	if queue.size == 0 {
		return QueueHead
	}

	var builder strings.Builder

	builder.WriteString(QueueHead)

	switch queue.implementation {
	case LINKED:
		tmp := queue.data.(linked_queue[T])

		builder.WriteString(fmt.Sprintf("%v", tmp.front.value))

		for queue_node := tmp.front.next; queue_node != nil; queue_node = queue_node.next {
			builder.WriteString(QueueSep)
			builder.WriteString(fmt.Sprintf("%v", queue_node.value))
		}
	case ARRAY:
		tmp := queue.data.([]T)

		builder.WriteString(fmt.Sprintf("%v", tmp[0]))

		for _, element := range tmp[1:] {
			builder.WriteString(QueueSep)
			builder.WriteString(fmt.Sprintf("%v", element))
		}
	}

	return builder.String()
}

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

type queue_node[T any] struct {
	value T
	next  *queue_node[T]
}

type linked_queue[T any] struct {
	front, back *queue_node[T]
}
