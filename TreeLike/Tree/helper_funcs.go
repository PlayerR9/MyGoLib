package Tree

import (
	"errors"

	uc "github.com/PlayerR9/MyGoLib/Units/common"
	us "github.com/PlayerR9/MyGoLib/Units/slice"
)

// recSnakeTraversal is an helper function that returns all the paths
// from n to the leaves of the tree rooted at n.
//
// Returns:
//   - [][]Noder: A slice of slices of elements.
//   - error: An error if the iterator fails.
//
// Behaviors:
//   - The paths are returned in the order of a BFS traversal.
//   - It is a recursive function.
func recSnakeTraversal(n Noder) ([][]Noder, error) {
	uc.AssertParam("n", n != nil, errors.New("recSnakeTraversal: n is nil"))

	ok := n.IsLeaf()
	if ok {
		return [][]Noder{
			{n},
		}, nil
	}

	iter := n.Iterator()
	if iter == nil {
		return nil, nil
	}

	var result [][]Noder

	for {
		value, err := iter.Consume()
		ok := uc.IsDone(err)
		if ok {
			break
		} else if err != nil {
			return nil, err
		}

		subResults, err := recSnakeTraversal(value)
		if err != nil {
			return nil, err
		}

		for _, tmp := range subResults {
			tmp = append([]Noder{n}, tmp...)
			result = append(result, tmp)
		}
	}

	return result, nil
}

// SnakeTraversal returns all the paths from the root to the leaves of the tree.
//
// Returns:
//   - [][]Noder: A slice of slices of elements.
//   - error: An error if the iterator fails.
//
// Behaviors:
//   - The paths are returned in the order of a BFS traversal.
func (t *Tree) SnakeTraversal() ([][]Noder, error) {
	root := t.root
	if root == nil {
		return nil, nil
	}

	sol, err := recSnakeTraversal(root)
	return sol, err
}

// recPruneFunc is an helper function that removes all the children of the
// node that satisfy the given filter including all of their children.
//
// Parameters:
//   - filter: The filter to apply.
//   - n: The node to prune.
//
// Returns:
//   - *Node[T]: A pointer to the highest ancestor of the pruned node.
//   - bool: True if the node satisfies the filter, false otherwise.
//
// Behaviors:
//   - This function is recursive.
func recPruneFunc(filter us.PredicateFilter[Noder], highest Noder, n Noder) (Noder, bool) {
	ok := filter(n)

	if ok {
		// Delete all children
		n.Cleanup()

		ancestors := FindCommonAncestor(highest, n)

		return ancestors, true
	}

	var prev Noder

	iter := n.Iterator()
	if iter == nil {
		return highest, false
	}

	for {
		value, err := iter.Consume()
		ok := uc.IsDone(err)
		if ok {
			break
		} else if err != nil {
			return highest, false
		}

		high, ok := recPruneFunc(filter, highest, c)
		if !ok {
			continue
		}

		prev_sibling := value.PrevSibling
		next_sibling := value.NextSibling

		if prev_sibling != nil {
			prev_sibling.NextSibling = next_sibling
		}

		if next_sibling != nil {
			next_sibling.PrevSibling = prev_sibling
		}

		value.PrevSibling = nil

		if prev != nil {
			prev.NextSibling = nil
		}

		highest = FindCommonAncestor(highest, high)

		prev = value
	}

	return highest, false
}

// PruneFunc removes all the children of the node that satisfy the given filter
// including all of their children.
//
// Parameters:
//   - filter: The filter to apply.
//
// Returns:
//   - bool: True if the node satisfies the filter, false otherwise.
//
// Behaviors:
//   - The root node is not pruned.
func (t *Tree) PruneFunc(filter us.PredicateFilter[Noder]) bool {
	if filter == nil {
		return false
	}

	root := t.root
	if root == nil {
		return false
	}

	highest, ok := recPruneFunc(filter, nil, root)
	if ok {
		return true
	}

	t.leaves = highest.Leaves()
	t.size = highest.Size()

	return false
}
