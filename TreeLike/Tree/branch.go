package Tree

import (
	uc "github.com/PlayerR9/MyGoLib/Units/common"
)

// BranchIterator is a generic data structure that represents an iterator
// for a branch in a tree.
type BranchIterator[T any] struct {
	// fromNode is the node from which the branch starts.
	fromNode *TreeNode[T]

	// toNode is the node to which the branch ends.
	toNode *TreeNode[T]

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

	if bi.current == bi.toNode {
		bi.current = nil
	} else {
		bi.current = bi.current.children[0]
	}

	return value, nil
}

// Restart implements the uc.Iterater interface.
func (bi *BranchIterator[T]) Restart() {
	bi.current = bi.fromNode
}

// Branch is a generic data structure that represents a branch in a tree.
type Branch[T any] struct {
	// fromNode is the node from which the branch starts.
	fromNode *TreeNode[T]

	// toNode is the node to which the branch ends.
	toNode *TreeNode[T]

	// size is the number of nodes in the branch.
	size int
}

// Copy implements the uc.Copier interface.
func (b *Branch[T]) Copy() uc.Copier {
	fromCopy := b.fromNode.Copy().(*TreeNode[T])
	toCopy := b.toNode.Copy().(*TreeNode[T])

	branchCopy := &Branch[T]{
		fromNode: fromCopy,
		toNode:   toCopy,
		size:     b.size,
	}

	return branchCopy
}

// Iterator implements the uc.Iterable interface.
func (b *Branch[T]) Iterator() uc.Iterater[T] {
	iter := &BranchIterator[T]{
		fromNode: b.fromNode,
		current:  b.fromNode,
		size:     b.size,
	}

	return iter
}

// Slice implements the uc.Slicer interface.
func (b *Branch[T]) Slice() []T {
	slice := make([]T, 0, b.size)

	for n := b.fromNode; n != b.toNode; n = n.children[0] {
		slice = append(slice, n.Data)
	}

	slice = append(slice, b.toNode.Data)

	return slice
}

// Size returns the number of nodes in the branch.
func (b *Branch[T]) Size() int {
	return b.size
}
