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

	// linked_queue with a maximum size Implementation
	LINKED_SIZE

	// ArrayQueue with no maximum size Implementation
	ARRAY

	// ArrayQueue with a maximum size Implementation
	ARRAY_SIZE
)

type Queue[T any] struct {
	data     any
	size     int
	capacity optional.Int
	methods  *queue_methods[T]
}

func (queue *Queue[T]) Initialize(values ...T) {
	queue.size = 0
	queue.data = queue.methods.new()

	for _, element := range values {
		queue.methods.enqueue(queue.data, element)
	}
}

func (queue *Queue[T]) Enqueue(value T) {
	if queue.capacity.Present() && queue.size >= queue.capacity.MustGet() {
		panic(ErrFullQueue{})
	}

	queue.methods.enqueue(queue.data, value)
	queue.size++
}

func (queue *Queue[T]) Dequeue() T {
	if queue.size <= 0 {
		panic(ErrEmptyQueue{})
	}

	queue.size--

	return queue.methods.dequeue(queue.data)
}

func (queue Queue[T]) Peek() T {
	if queue.size == 0 {
		panic(ErrEmptyQueue{})
	}

	return queue.methods.peek(queue.data)
}

func (queue Queue[T]) IsEmpty() bool {
	return queue.size == 0
}

func (queue Queue[T]) Size() int {
	return queue.size
}

func (queue *Queue[T]) ToSlice() []T {
	slice := make([]T, queue.size)

	queue.methods.to_slice(queue.data, slice)

	return slice
}

func (queue *Queue[T]) Clear() {
	queue.size = 0
	queue.data = queue.methods.new()
}

func (queue Queue[T]) IsFull() bool {
	return queue.capacity.Present() && queue.size == queue.capacity.MustGet()
}

func (queue *Queue[T]) String() string {
	var str strings.Builder

	str.WriteString(QueueHead)

	if queue.size != 0 {
		queue.methods.stringer(queue.data, &str)
	}

	return str.String()
}

type queue_methods[T any] struct {
	new      func() any
	enqueue  func(any, T)
	dequeue  func(any) T
	peek     func(any) T
	to_slice func(any, []T)
	stringer func(any, *strings.Builder)
}

func NewQueue[T any](implementation int, capacity optional.Int) Queue[T] {
	var queue Queue[T]

	switch implementation {
	case LINKED, LINKED_SIZE:
		queue.data = linked_queue[T]{
			front: nil,
			back:  nil,
		}

		queue.size = 0

		if implementation == LINKED_SIZE {
			if !capacity.Present() {
				panic("Must specify capacity for a linked queue with a maximum size")
			}

			if capacity.MustGet() < 0 {
				panic("Cannot specify a negative capacity for a linked queue")
			}

			queue.capacity = capacity
		}

		if capacity.Present() && implementation != LINKED_SIZE {
			panic("Cannot specify capacity for a linked queue with no maximum size")
		}

		queue.methods = &queue_methods[T]{
			new: func() any {
				return linked_queue[T]{
					front: nil,
					back:  nil,
				}
			},

			enqueue: func(data any, value T) {
				new_node := queue_node[T]{
					value: gen.DeepCopy(value).(T),
				}

				tmp := data.(linked_queue[T])

				if tmp.back == nil {
					tmp.front = &new_node
				} else {
					tmp.back.next = &new_node
				}

				tmp.back = &new_node

				data = tmp
			},

			dequeue: func(data any) T {
				tmp := data.(linked_queue[T])

				value := tmp.front.value

				tmp.front = tmp.front.next

				if tmp.front == nil {
					tmp.back = nil
				}

				data = tmp

				return value
			},

			peek: func(data any) T {
				return gen.DeepCopy(data.(linked_queue[T]).front.value).(T)
			},

			to_slice: func(data any, slice []T) {
				i := 0

				for queue_node := data.(linked_queue[T]).front; queue_node != nil; queue_node = queue_node.next {
					slice[i] = gen.DeepCopy(queue_node.value).(T)
					i++
				}
			},

			stringer: func(data any, str *strings.Builder) {
				tmp := data.(linked_queue[T])

				str.WriteString(fmt.Sprintf("%v", tmp.front.value))

				for queue_node := tmp.front.next; queue_node != nil; queue_node = queue_node.next {
					str.WriteString(QueueSep)
					str.WriteString(fmt.Sprintf("%v", queue_node.value))
				}
			},
		}
	case ARRAY, ARRAY_SIZE:
		queue.data = make([]T, 0)

		queue.size = 0

		if capacity.Present() {
			if implementation == ARRAY_SIZE {
				queue.capacity = capacity
			} else {
				panic("Cannot specify capacity for an array queue with no maximum size")
			}
		} else if implementation == ARRAY_SIZE {
			panic("Must specify capacity for an array queue with a maximum size")
		}

		queue.methods = &queue_methods[T]{
			new: func() any {
				return make([]T, 0)
			},

			enqueue: func(data any, value T) {
				tmp := data.([]T)
				tmp = append(tmp, gen.DeepCopy(value).(T))
				data = tmp
			},

			dequeue: func(data any) T {
				tmp := data.([]T)
				value := tmp[0]
				data = tmp[1:]

				return value
			},

			peek: func(data any) T {
				tmp := data.([]T)
				return gen.DeepCopy(tmp[0]).(T)
			},

			to_slice: func(data any, slice []T) {
				for i, element := range data.([]T) {
					slice[i] = gen.DeepCopy(element).(T)
				}
			},

			stringer: func(data any, str *strings.Builder) {
				tmp := data.([]T)

				str.WriteString(fmt.Sprintf("%v", tmp[0]))

				for _, element := range tmp[1:] {
					str.WriteString(QueueSep)
					str.WriteString(fmt.Sprintf("%v", element))
				}
			},
		}
	default:
		panic(fmt.Sprintf("Invalid queue implementation type: %d", implementation))
	}

	return queue
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
