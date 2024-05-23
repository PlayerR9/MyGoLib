package Traversor

import (
	uc "github.com/PlayerR9/MyGoLib/Units/Common"

	tr "github.com/PlayerR9/MyGoLib/TreeLike/Tree"
)

// stackElement is a stack element.
type stackElement[T any] struct {
	// prev is the previous node.
	prev *tr.TreeNode[T]

	// elem is the current node.
	elem *tr.TreeNode[T]

	// info is the info of the current node.
	info uc.Objecter

	// nextFunc is the function to get the next node.
	nextFunc NextsFunc[T]
}

// newStackElement creates a new stack element.
//
// Parameters:
//   - prev: The previous node.
//   - elem: The helper.
//
// Returns:
//   - *stackElement[T, E]: A pointer to the stack element.
func newStackElement[T any](prev *tr.TreeNode[T], data T, info uc.Objecter, nextFunc NextsFunc[T]) *stackElement[T] {
	se := &stackElement[T]{
		prev:     prev,
		elem:     tr.NewTreeNode(data),
		nextFunc: nextFunc,
	}

	if info == nil {
		se.info = nil
	} else {
		se.info = info.Copy()
	}

	return se
}

// apply applies the next function.
//
// Returns:
//   - []T: A slice of the next elements.
//   - error: An error if the function fails.
func (se *stackElement[T]) apply() ([]T, error) {
	if se.nextFunc == nil {
		return nil, nil
	}

	return se.nextFunc(se.elem.Data, se.info)
}
