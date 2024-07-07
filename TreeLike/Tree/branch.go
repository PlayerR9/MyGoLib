package Tree

import (
	uc "github.com/PlayerR9/MyGoLib/Units/common"
)

// BranchIterator is a generic data structure that represents an iterator
// for a branch in a tree.
type BranchIterator[T any] struct {
	// from_node is the node from which the branch starts.
	from_node *TreeNode[T]

	// to_node is the node to which the branch ends.
	to_node *TreeNode[T]

	// current is the current node of the iterator.
	current *TreeNode[T]

	// size is the number of nodes in the branch.
	size int
}

// Size implements the uc.Iterater interface.
func (bi *BranchIterator[T]) Size() (count int) {
	return bi.size
}

// Consume implements the uc.Iterater interface.
//
// This scans from the root node to the leaf node.
func (bi *BranchIterator[T]) Consume() (T, error) {
	if bi.current == nil {
		return *new(T), uc.NewErrExhaustedIter()
	}

	value := bi.current.Data

	if bi.current == bi.to_node {
		bi.current = nil
	} else {
		bi.current = bi.current.FirstChild
	}

	return value, nil
}

// Restart implements the uc.Iterater interface.
func (bi *BranchIterator[T]) Restart() {
	bi.current = bi.from_node
}

// Branch is a generic data structure that represents a branch in a tree.
type Branch[T any] struct {
	// from_node is the node from which the branch starts.
	from_node *TreeNode[T]

	// to_node is the node to which the branch ends.
	to_node *TreeNode[T]

	// size is the number of nodes in the branch.
	size int
}

// Copy implements the uc.Copier interface.
func (b *Branch[T]) Copy() uc.Copier {
	from_copy := b.from_node.Copy().(*TreeNode[T])
	to_copy := b.to_node.Copy().(*TreeNode[T])

	b_copy := &Branch[T]{
		from_node: from_copy,
		to_node:   to_copy,
		size:      b.size,
	}

	return b_copy
}

// Iterator implements the uc.Iterable interface.
func (b *Branch[T]) Iterator() uc.Iterater[T] {
	iter := &BranchIterator[T]{
		from_node: b.from_node,
		current:   b.from_node,
		size:      b.size,
	}

	return iter
}

// Slice implements the uc.Slicer interface.
func (b *Branch[T]) Slice() []T {
	slice := make([]T, 0, b.size)

	for n := b.from_node; n != b.to_node; n = n.FirstChild {
		slice = append(slice, n.Data)
	}

	slice = append(slice, b.to_node.Data)

	return slice
}

// Size returns the number of nodes in the branch.
func (b *Branch[T]) Size() int {
	return b.size
}
