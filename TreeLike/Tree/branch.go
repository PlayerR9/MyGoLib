package Tree

import (
	uc "github.com/PlayerR9/MyGoLib/Units/common"
)

// BranchIterator is a generic data structure that represents an iterator
// for a branch in a tree.
type BranchIterator struct {
	// from_node is the node from which the branch starts.
	from_node Noder

	// to_node is the node to which the branch ends.
	to_node Noder

	// current is the current node of the iterator.
	current Noder

	// size is the number of nodes in the branch.
	size int
}

// Consume implements the uc.Iterater interface.
//
// This scans from the root node to the leaf node.
func (bi *BranchIterator) Consume() (Noder, error) {
	value := bi.current

	if bi.current == bi.to_node {
		return nil, uc.NewErrExhaustedIter()
	}

	bi.current = bi.current.FirstChild
	return value, nil
}

// Restart implements the uc.Iterater interface.
func (bi *BranchIterator) Restart() {
	bi.current = bi.from_node
}

// Branch is a generic data structure that represents a branch in a tree.
type Branch struct {
	// from_node is the node from which the branch starts.
	from_node Noder

	// to_node is the node to which the branch ends.
	to_node Noder

	// size is the number of nodes in the branch.
	size int
}

// Copy implements the uc.Copier interface.
func (b *Branch) Copy() uc.Copier {
	from_copy := b.from_node.Copy().(Noder)
	to_copy := b.to_node.Copy().(Noder)

	b_copy := &Branch{
		from_node: from_copy,
		to_node:   to_copy,
		size:      b.size,
	}

	return b_copy
}

// Iterator implements the uc.Iterable interface.
func (b *Branch) Iterator() uc.Iterater[Noder] {
	iter := &BranchIterator{
		from_node: b.from_node,
		current:   b.from_node,
		size:      b.size,
	}

	return iter
}

// Slice implements the uc.Slicer interface.
func (b *Branch) Slice() []Noder {
	slice := make([]Noder, 0, b.size)

	for n := b.from_node; n != b.to_node; n = n.FirstChild {
		slice = append(slice, n)
	}

	slice = append(slice, b.to_node)

	return slice
}

// Size returns the number of nodes in the branch.
//
// Returns:
//   - int: The number of nodes in the branch.
func (b *Branch) Size() int {
	return b.size
}
