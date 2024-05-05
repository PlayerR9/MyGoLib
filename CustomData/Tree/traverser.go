package Tree

import (
	"slices"

	"github.com/PlayerR9/MyGoLib/ListLike/Queuer"
	"github.com/PlayerR9/MyGoLib/ListLike/Stacker"
)

// ObserverFunc is a function that observes a node.
//
// Parameters:
//   - node: The node to observe.
//
// Returns:
//   - error: An error if the observation fails.
type ObserverFunc[T any] func(node T) error

// NextsFunc is a function that returns the children of a node.
//
// Parameters:
//   - node: The node to compute the children of.
//
// Returns:
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
//   - observer: The function that observes a node.
//   - nexts: The function that computes the children of a node.
//
// Returns:
//   - Traverser[T]: The newly created traverser.
func NewTraverser[T any](observer ObserverFunc[T], nexts NextsFunc[T]) Traverser[T] {
	return Traverser[T]{
		observer: observer,
		nexts:    nexts,
	}
}

// DFS traverses the tree in depth-first order.
//
// Parameters:
//   - node: The root node of the tree.
//
// Returns:
//   - error: An error if the traversal fails.
func (t *Traverser[T]) DFS(node T) error {
	S := Stacker.NewLinkedStack(node)

	var children []T

	for {
		node, err := S.Pop()
		if err != nil {
			break
		}

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
			err := S.Push(child)
			if err != nil {
				panic(err)
			}
		}
	}

	return nil
}

// BFS traverses the tree in breadth-first order.
//
// Parameters:
//   - node: The root node of the tree.
//
// Returns:
//   - error: An error if the traversal fails.
func (t *Traverser[T]) BFS(node T) error {
	Q := Queuer.NewLinkedQueue(node)

	var children []T

	for {
		node, err := Q.Dequeue()
		if err != nil {
			break
		}

		err = t.observer(node)
		if err != nil {
			return err
		}

		children, err = t.nexts(node)
		if err != nil {
			return err
		}

		for _, child := range children {
			err := Q.Enqueue(child)
			if err != nil {
				panic(err)
			}
		}
	}

	return nil
}

// MakeTree creates a tree from the given element.
//
// Parameters:
//   - elem: The element to start the tree from.
//
// Returns:
//   - *Tree[T]: The tree created from the element.
//   - error: An error if the tree cannot be created.
func (t *Traverser[T]) MakeTree(elem T) (*Tree[T], error) {
	// 1. Handle the first element
	err := t.observer(elem)
	if err != nil {
		return nil, err
	}

	tree := &Tree[T]{
		root: NewNode(elem),
	}

	nexts, err := t.nexts(elem)
	if err != nil {
		return nil, err
	}

	if len(nexts) == 0 {
		return tree, nil
	}

	// 2. Create a stack and push the first element
	type StackElement struct {
		Prev *Node[T]
		Elem T
	}

	S := Stacker.NewLinkedStack[StackElement]()

	for _, next := range nexts {
		err := S.Push(StackElement{Prev: tree.root, Elem: next})
		if err != nil {
			panic(err)
		}
	}

	for {
		se, err := S.Pop()
		if err != nil {
			break
		}

		err = t.observer(se.Elem)
		if err != nil {
			return nil, err
		}

		node := NewNode(se.Elem)
		se.Prev.AddChild(node)

		nexts, err := t.nexts(elem)
		if err != nil {
			return nil, err
		}

		for _, next := range nexts {
			err := S.Push(StackElement{Prev: node, Elem: next})
			if err != nil {
				panic(err)
			}
		}
	}

	return tree, nil
}
