package Traversor

import (
	"github.com/PlayerR9/MyGoLib/ListLike/Queuer"
	"github.com/PlayerR9/MyGoLib/ListLike/Stacker"
	uc "github.com/PlayerR9/MyGoLib/Units/Common"

	tr "github.com/PlayerR9/MyGoLib/TreeLike/Tree"
)

// ObserverFunc is a function that observes a node.
//
// Parameters:
//   - data: The data of the node.
//   - info: The info of the node.
//
// Returns:
//   - bool: True if the traversal should continue, otherwise false.
//   - error: An error if the observation fails.
type ObserverFunc[T any] func(data T, info uc.Objecter) (bool, error)

// traversor is a struct that traverses a tree.
type traversor[T any] struct {
	// elem is the current node.
	elem *tr.TreeNode[T]

	// info is the info of the current node.
	info uc.Objecter
}

// newTraversor creates a new traversor for the tree.
//
// Parameters:
//   - tree: The tree to traverse.
//   - init: The initial info.
//
// Returns:
//   - Traversor[T, I]: The traversor.
func newTraversor[T any](node *tr.TreeNode[T], init uc.Objecter) *traversor[T] {
	t := &traversor[T]{
		elem: node,
	}

	if init != nil {
		t.info = init.Copy()
	} else {
		t.info = nil
	}

	return t
}

// getData returns the data of the traversor.
//
// Returns:
//   - T: The data of the traversor.
//   - bool: True if the data is valid, otherwise false.
func (t *traversor[T]) getData() (T, bool) {
	if t.elem == nil {
		return *new(T), false
	}

	return t.elem.Data, true
}

// getInfo returns the info of the traversor.
//
// Returns:
//   - uc.Objecter: The info of the traversor.
func (t *traversor[T]) getInfo() uc.Objecter {
	return t.info
}

// getChildren returns the children of the traversor.
//
// Returns:
//   - []*tr.TreeNode[T]: The children of the traversor.
func (t *traversor[T]) getChildren() []*tr.TreeNode[T] {
	if t.elem == nil {
		return nil
	}

	return t.elem.GetChildren()
}

// DFS traverses the tree in depth-first order.
//
// Parameters:
//   - t: The traversor.
//
// Returns:
//   - error: An error if the traversal fails.
func DFS[T any](tree *tr.Tree[T], init uc.Objecter, f ObserverFunc[T]) error {
	if f == nil || tree == nil {
		return nil
	}

	S := Stacker.NewLinkedStack(newTraversor(tree.Root(), init))

	for {
		top, err := S.Pop()
		if err != nil {
			break
		}

		topData, ok := top.getData()
		if !ok {
			panic("Missing data")
		}
		topInfo := top.getInfo()

		ok, err = f(topData, topInfo)
		if err != nil {
			return err
		} else if !ok {
			continue
		}

		children := top.getChildren()
		if len(children) == 0 {
			continue
		}

		for _, child := range children {
			newT := newTraversor(child, topInfo)

			err := S.Push(newT)
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
//   - t: The traversor.
//
// Returns:
//   - error: An error if the traversal fails.
func BFS[T any](tree *tr.Tree[T], init uc.Objecter, f ObserverFunc[T]) error {
	if f == nil || tree == nil {
		return nil
	}

	Q := Queuer.NewLinkedQueue(newTraversor(tree.Root(), init))

	for {
		first, err := Q.Dequeue()
		if err != nil {
			break
		}

		firstData, ok := first.getData()
		if !ok {
			panic("Missing data")
		}
		firstInfo := first.getInfo()

		ok, err = f(firstData, firstInfo)
		if err != nil {
			return err
		} else if !ok {
			continue
		}

		children := first.getChildren()
		if len(children) == 0 {
			continue
		}

		for _, child := range children {
			newT := newTraversor(child, firstInfo)

			err := Q.Enqueue(newT)
			if err != nil {
				panic(err)
			}
		}
	}

	return nil
}
