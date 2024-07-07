package Tree

import (
	"errors"

	uc "github.com/PlayerR9/MyGoLib/Units/common"
	us "github.com/PlayerR9/MyGoLib/Units/slice"
)

// Cleanup removes every node in the tree in a DFS traversal and sets the root to nil.
func (t *Tree[T]) Cleanup() {
	root := t.root
	if root == nil {
		return
	}

	root.Cleanup()

	t.root = nil
}

// recSnakeTraversal is an helper function that returns all the paths
// from n to the leaves of the tree rooted at n.
//
// Returns:
//   - result: A slice of slices of elements.
//
// Behaviors:
//   - The paths are returned in the order of a BFS traversal.
//   - It is a recursive function.
func recSnakeTraversal[T any](n *TreeNode[T]) (result [][]T) {
	uc.AssertParam("n", n != nil, errors.New("recSnakeTraversal: n is nil"))

	if n.FirstChild == nil {
		return [][]T{
			{n.Data},
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		subResults := recSnakeTraversal(c)

		for _, tmp := range subResults {
			result = append(result, append([]T{n.Data}, tmp...))
		}
	}

	return
}

// SnakeTraversal returns all the paths from the root to the leaves of the tree.
//
// Returns:
//   - result: A slice of slices of elements.
//
// Behaviors:
//   - The paths are returned in the order of a BFS traversal.
func (t *Tree[T]) SnakeTraversal() (result [][]T) {
	root := t.root
	if root == nil {
		return
	}

	sol := recSnakeTraversal(root)

	return sol
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
func recPruneFunc[T any](filter us.PredicateFilter[T], highest *TreeNode[T], n *TreeNode[T]) (*TreeNode[T], bool) {
	ok := filter(n.Data)

	if ok {
		// Delete all children
		n.Cleanup()

		ancestors := FindCommonAncestor(highest, n)

		return ancestors, true
	}

	var prev *TreeNode[T]

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		high, ok := recPruneFunc(filter, highest, c)
		if !ok {
			continue
		}

		prev_sibling := c.PrevSibling
		next_sibling := c.NextSibling

		if prev_sibling != nil {
			prev_sibling.NextSibling = next_sibling
		}

		if next_sibling != nil {
			next_sibling.PrevSibling = prev_sibling
		}

		c.PrevSibling = nil

		if prev != nil {
			prev.NextSibling = nil
		}

		highest = FindCommonAncestor(highest, high)

		prev = c
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
func (t *Tree[T]) PruneFunc(filter us.PredicateFilter[T]) bool {
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
