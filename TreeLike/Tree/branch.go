package Tree

import (
	uc "github.com/PlayerR9/MyGoLib/Units/common"
)

// BranchIterator is the pull-based iterator for the branch.
type BranchIterator struct {
	// from_node is the node from which the branch starts.
	from_node Noder

	// to_node is the node to which the branch ends.
	to_node Noder

	// current is the current node of the iterator.
	current Noder
}

// Consume implements the common.Iterater interface.
//
// This scans from the root node to the leaf node.
//
// Noder is never nil.
func (bi *BranchIterator) Consume() (Noder, error) {
	uc.Assert(bi.current != nil, "BranchIterator: current is nil")

	value := bi.current

	if bi.current == bi.to_node {
		return nil, uc.NewErrExhaustedIter()
	}

	bi.current = bi.current.GetFirstChild()
	return value, nil
}

// Restart implements the common.Iterater interface.
func (bi *BranchIterator) Restart() {
	bi.current = bi.from_node
}

// Branch represents a branch in a tree.
type Branch struct {
	// from_node is the node from which the branch starts.
	from_node Noder

	// to_node is the node to which the branch ends.
	to_node Noder
}

// Copy implements the uc.Copier interface.
func (b *Branch) Copy() uc.Copier {
	from_copy := b.from_node.Copy().(Noder)
	to_copy := b.to_node.Copy().(Noder)

	b_copy := &Branch{
		from_node: from_copy,
		to_node:   to_copy,
	}

	return b_copy
}

// Iterator implements the uc.Iterable interface.
func (b *Branch) Iterator() uc.Iterater[Noder] {
	iter := &BranchIterator{
		from_node: b.from_node,
		current:   b.from_node,
	}

	return iter
}

// Slice implements the uc.Slicer interface.
func (b *Branch) Slice() []Noder {
	var slice []Noder

	n := b.from_node
	for n != b.to_node {
		slice = append(slice, n)

		n = n.GetFirstChild()
	}

	slice = append(slice, b.to_node)

	return slice
}

// NewBranch works like GetAncestors but includes the node itself.
//
// The nodes are returned as a slice where [0] is the root node
// and [len(branch)-1] is the leaf node.
//
// Returns:
//   - *Branch: The branch from the node to the root.
func NewBranch(n Noder) *Branch {
	branch := &Branch{
		to_node: n,
	}

	node := n

	for {
		parent := node.GetParent()
		if parent == nil {
			break
		}

		node = parent
	}

	branch.from_node = node

	return branch
}
