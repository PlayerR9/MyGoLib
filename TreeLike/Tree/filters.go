package Tree

import (
	us "github.com/PlayerR9/MyGoLib/Units/slice"
)

var (
	// FilterNonNilTree is a filter that returns true if the tree is not nil.
	//
	// Parameters:
	//   - tree: The tree to filter.
	//
	// Returns:
	//   - bool: True if the tree is not nil, false otherwise.
	FilterNonNilTree us.PredicateFilter[*Tree]
)

func init() {
	FilterNonNilTree = func(tree *Tree) bool {
		if tree == nil {
			return false
		}

		return tree.root != nil
	}
}
