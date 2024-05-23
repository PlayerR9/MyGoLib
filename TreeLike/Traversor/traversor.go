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
func newTraversor[T any](data T, init uc.Objecter) *traversor[T] {
	t := &traversor[T]{
		elem: tr.NewTreeNode(data),
	}

	if init != nil {
		t.info = init.Copy()
	} else {
		t.info = nil
	}

	return t
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

	S := Stacker.NewLinkedStack(newTraversor(tree.Root().Data, init))

	for {
		top, err := S.Pop()
		if err != nil {
			break
		}

		ok, err := f(top.elem.Data, top.info)
		if err != nil {
			return err
		}

		if !ok {
			continue
		}

		children := top.elem.GetChildren()

		if len(children) == 0 {
			continue
		}

		for _, child := range children {
			newT := newTraversor(child.Data, top.info)

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

	Q := Queuer.NewLinkedQueue(newTraversor(tree.Root().Data, init))

	for {
		first, err := Q.Dequeue()
		if err != nil {
			break
		}

		ok, err := f(first.elem.Data, first.info)
		if err != nil {
			return err
		}

		if !ok {
			continue
		}

		for _, child := range first.elem.GetChildren() {
			newT := newTraversor(child.Data, first.info)

			err := Q.Enqueue(newT)
			if err != nil {
				panic(err)
			}
		}
	}

	return nil
}
