package Tree

import us "github.com/PlayerR9/MyGoLib/Units/slice"

// recCleanup is a helper function that removes every node in the tree rooted at n.
//
// Behaviors:
//   - This function is recursive.
func recCleanup[T any](n *TreeNode[T]) {
	n.parent = nil

	for _, child := range n.children {
		recCleanup(child)
	}

	for i := range n.children {
		n.children[i] = nil
	}

	n.children = nil
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
	if len(n.children) == 0 {
		return [][]T{
			{n.Data},
		}
	}

	for _, child := range n.children {
		subResults := recSnakeTraversal(child)

		for _, tmp := range subResults {
			result = append(result, append([]T{n.Data}, tmp...))
		}
	}

	return
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
	if filter(n.Data) {
		// Delete all children
		recCleanup(n)

		ancestors := FindCommonAncestor(highest, n)

		return ancestors, true
	}

	top := 0

	for i := 0; i < len(n.children); i++ {
		high, ok := recPruneFunc(filter, highest, n.children[i])
		if ok {
			n.children[i] = nil

			highest = FindCommonAncestor(highest, high)
		} else {
			n.children[top] = n.children[i]
			top++
		}
	}

	n.children = n.children[:top]

	return highest, false
}
