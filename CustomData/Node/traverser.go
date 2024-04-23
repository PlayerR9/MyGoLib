package Node

import (
	"slices"

	Queue "github.com/PlayerR9/MyGoLib/ListLike/Queue"
	Stack "github.com/PlayerR9/MyGoLib/ListLike/Stack"
)

// ObserverFunc is a function that observes a node.
//
// Parameters:
//
//   - node: The node to observe.
//
// Returns:
//
//   - error: An error if the observation fails.
type ObserverFunc[T any] func(node T) error

// NextsFunc is a function that returns the children of a node.
//
// Parameters:
//
//   - node: The node to compute the children of.
//
// Returns:
//
//   - []T: The children of the node.
//   - error: An error if children cannot be computed.
type NextsFunc[T any] func(node T) ([]T, error)

// Traverser is a struct that traverses a tree.
type Traverser[T any] struct {
	// The function that observes a node.
	observer ObserverFunc[T]

	// The function that computes the children of a node.
	nexts NextsFunc[T]
}

// NewTraverser creates a new traverser with the given observer and nexts functions.
//
// Parameters:
//
//   - observer: The function that observes a node.
//   - nexts: The function that computes the children of a node.
func NewTraverser[T any](observer ObserverFunc[T], nexts NextsFunc[T]) Traverser[T] {
	return Traverser[T]{
		observer: observer,
		nexts:    nexts,
	}
}

// DFS traverses the tree in depth-first order.
//
// Parameters:
//
//   - node: The root node of the tree.
//
// Returns:
//
//   - error: An error if the traversal fails.
func (t *Traverser[T]) DFS(node T) error {
	S := Stack.NewLinkedStack(node)

	var err error
	var children []T

	for !S.IsEmpty() {
		node := S.Pop()

		err = t.observer(node)
		if err != nil {
			return err
		}

		children, err = t.nexts(node)
		if err != nil {
			return err
		}

		slices.Reverse(children)

		for _, child := range children {
			S.Push(child)
		}
	}

	return nil
}

// BFS traverses the tree in breadth-first order.
//
// Parameters:
//
//   - node: The root node of the tree.
//
// Returns:
//
//   - error: An error if the traversal fails.
func (t *Traverser[T]) BFS(node T) error {
	Q := Queue.NewLinkedQueue(node)

	var err error
	var children []T

	for !Q.IsEmpty() {
		node := Q.Dequeue()

		err = t.observer(node)
		if err != nil {
			return err
		}

		children, err = t.nexts(node)
		if err != nil {
			return err
		}

		for _, child := range children {
			Q.Enqueue(child)
		}
	}

	return nil
}
