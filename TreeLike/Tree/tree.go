package Tree

import (
	fsp "github.com/PlayerR9/MyGoLib/Formatting/FString"
	"github.com/PlayerR9/MyGoLib/ListLike/Queuer"
	"github.com/PlayerR9/MyGoLib/ListLike/Stacker"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
	us "github.com/PlayerR9/MyGoLib/Units/slice"
)

// Tree is a generic data structure that represents a tree.
type Tree[T any] struct {
	// root is the root of the tree.
	root *TreeNode[T]

	// leaves is the leaves of the tree.
	leaves []*TreeNode[T]

	// size is the number of nodes in the tree.
	size int
}

// FString implements the FString.FStringer interface for the Tree type.
func (t *Tree[T]) FString(trav *fsp.Traversor, opts ...fsp.Option) error {
	root := t.root
	if root == nil {
		return nil
	}

	form := fsp.NewFormatter(
		fsp.NewIndentConfig("| ", 0),
	)

	err := fsp.ApplyForm(form, trav, root)
	if err != nil {
		return err
	}

	return nil
}

// Copy creates a deep copy of the tree.
//
// Returns:
//   - uc.Copier: A deep copy of the tree.
func (t *Tree[T]) Copy() uc.Copier {
	root := t.root.Copy().(*TreeNode[T])

	tree := &Tree[T]{
		root:   root,
		leaves: t.leaves,
		size:   t.size,
	}

	return tree
}

// NewTree creates a new tree with the given root.
//
// Parameters:
//   - data: The value of the root.
//
// Returns:
//   - *Tree[T]: A pointer to the newly created tree.
func NewTree[T any](data T) *Tree[T] {
	root := NewTreeNode(data)

	tree := &Tree[T]{
		root:   root,
		leaves: []*TreeNode[T]{root},
		size:   1,
	}

	return tree
}

// SetChildren sets the children of the root of the tree.
//
// Parameters:
//   - children: The children to set.
//
// Returns:
//   - error: An error of type *ErrMissingRoot if the tree does not have a root.
func (t *Tree[T]) SetChildren(children []*Tree[T]) error {
	children = us.SliceFilter(children, FilterNilTree)
	if len(children) == 0 {
		return nil
	}

	root := t.root
	if root == nil {
		return NewErrMissingRoot()
	}

	var leaves, subchildren []*TreeNode[T]

	t.size = 1

	for _, child := range children {
		leaves = append(leaves, child.leaves...)
		t.size += child.Size()
		child.root.parent = root
		subchildren = append(subchildren, child.root)
	}

	root.children = subchildren
	t.leaves = leaves

	return nil
}

// IsSingleton returns true if the tree has only one node.
//
// Returns:
//   - bool: True if the tree has only one node, false otherwise.
func (t *Tree[T]) IsSingleton() bool {
	return t.size == 1
}

// Size returns the number of nodes in the tree.
//
// Returns:
//   - int: The number of nodes in the tree.
func (t *Tree[T]) Size() int {
	return t.size
}

// Root returns the root of the tree.
//
// Returns:
//   - *Node[T]: A pointer to the root of the tree.
func (t *Tree[T]) Root() *TreeNode[T] {
	return t.root
}

// GetChildren returns all the children of the tree in a DFS order.
//
// Returns:
//   - children: A slice of the values of the children of the tree.
//
// Behaviors:
//   - The root is the first element in the slice.
//   - If the tree does not have a root, it returns nil.
func (t *Tree[T]) GetChildren() (children []T) {
	root := t.root
	if root == nil {
		return nil
	}

	S := Stacker.NewLinkedStack(root)

	for {
		node, ok := S.Pop()
		if !ok {
			break
		}

		children = append(children, node.Data)

		for _, child := range node.children {
			S.Push(child)
		}
	}

	return children
}

// Cleanup removes every node in the tree.
//
// Behaviors:
//   - This function is recursive and so, it is expensive.
func (t *Tree[T]) Cleanup() {
	root := t.root

	recCleanup(root)

	t.root = nil
}

// GetLeaves returns all the leaves of the tree.
//
// Returns:
//   - []*Node[T]: A slice of pointers to the leaves of the tree.
//
// Behaviors:
//   - Always returns at least one leaf.
//   - It returns the leaves that are stored in the tree. Make sure to call
//     any update function before calling this function if the tree has been modified
//     unexpectedly.
func (t *Tree[T]) GetLeaves() []*TreeNode[T] {
	return t.leaves
}

// RegenerateLeaves regenerates the leaves of the tree and returns them.
//
// Returns:
//   - []*Node[T]: A slice of pointers to the leaves of the tree.
//
// Behaviors:
//   - The leaves are updated in a DFS order.
//   - Expensive operation; use it only when necessary (i.e., leaves changed unexpectedly.)
//   - This also updates the size of the tree.
func (t *Tree[T]) RegenerateLeaves() []*TreeNode[T] {
	root := t.root

	if root == nil {
		t.leaves = make([]*TreeNode[T], 0)
		t.size = 0

		return t.leaves
	}

	var leaves []*TreeNode[T]

	S := Stacker.NewLinkedStack(root)

	t.size = 0

	for {
		top, ok := S.Pop()
		if !ok {
			break
		}

		t.size++

		if len(top.children) == 0 {
			leaves = append(leaves, top)
		} else {
			for _, child := range top.children {
				S.Push(child)
			}
		}
	}

	t.leaves = leaves

	return leaves
}

// UpdateLeaves updates the leaves of the tree and returns them.
//
// Returns:
//   - []*Node[T]: A slice of pointers to the leaves of the tree.
//
// Behaviors:
//   - The leaves are updated in a DFS order.
//   - Less expensive than RegenerateLeaves. However, if nodes has been deleted
//     from the tree, this may give unexpected results.
//   - This also updates the size of the tree.
func (t *Tree[T]) UpdateLeaves() []*TreeNode[T] {
	if len(t.leaves) == 0 {
		t.size = 0

		return t.leaves
	}

	var leaves []*TreeNode[T]

	S := Stacker.NewLinkedStack(t.leaves...)

	t.size -= len(t.leaves)

	for {
		top, ok := S.Pop()
		if !ok {
			break
		}

		t.size++

		if len(top.children) == 0 {
			leaves = append(leaves, top)
		} else {
			for _, child := range top.children {
				S.Push(child)
			}
		}
	}

	t.leaves = leaves

	return leaves
}

// SnakeTraversal returns all the paths from the root to the leaves of the tree.
//
// Returns:
//   - [][]T: A slice of slices of the values of the nodes in the paths.
//
// Behaviors:
//   - The paths are returned in the order of a DFS traversal.
//   - If the tree is empty, it returns an empty slice.
func (t *Tree[T]) SnakeTraversal() [][]T {
	root := t.root
	if root == nil {
		return nil
	}

	trav := recSnakeTraversal(root)
	return trav
}

// HasChild returns true if the tree has the given child in any of its nodes
// in a BFS order.
//
// Parameters:
//   - filter: The filter to apply.
//
// Returns:
//   - bool: True if the tree has the child, false otherwise.
func (t *Tree[T]) HasChild(filter us.PredicateFilter[T]) bool {
	if filter == nil {
		return false
	}

	root := t.root
	if root == nil {
		return false
	}

	Q := Queuer.NewLinkedQueue(root)

	for {
		node, ok := Q.Dequeue()
		if !ok {
			break
		}

		ok = filter(node.Data)
		if ok {
			return true
		}

		for _, child := range node.children {
			Q.Enqueue(child)
		}
	}

	return false
}

// FilterChildren returns all the children of the tree that satisfy the given filter
// in a BFS order.
//
// Parameters:
//   - filter: The filter to apply.
//
// Returns:
//   - sols: A slice of pointers to the nodes that satisfy the filter.
func (t *Tree[T]) FilterChildren(filter us.PredicateFilter[T]) (sols []*TreeNode[T]) {
	if filter == nil {
		return nil
	}

	root := t.root
	if root == nil {
		return
	}

	Q := Queuer.NewLinkedQueue(root)

	for {
		node, ok := Q.Dequeue()
		if !ok {
			break
		}

		ok = filter(node.Data)
		if ok {
			sols = append(sols, node)
		}

		for _, child := range node.children {
			Q.Enqueue(child)
		}
	}

	return
}

// PruneBranches removes all the children of the node that satisfy the given filter.
// The filter is a function that takes the value of a node and returns a boolean.
// If the filter returns true for a child, the child is removed along with its children.
//
// Parameters:
//   - filter: The filter to apply.
//
// Returns:
//   - bool: True if the whole tree can be deleted, false otherwise.
//
// Behaviors:
//   - If the root satisfies the filter, the tree is cleaned up.
//   - It is a recursive function.
func (t *Tree[T]) PruneBranches(filter us.PredicateFilter[T]) bool {
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

// SkipFunc removes all the children of the tree that satisfy the given filter
// without removing any of their children. Useful for removing unwanted nodes from the tree.
//
// Parameters:
//   - filter: The filter to apply.
//
// Returns:
//   - []*Tree[T]: A slice of pointers to the trees obtained after removing the nodes.
//
// Behaviors:
//   - If this function returns only one tree, this is the updated tree. But, if
//     it returns more than one tree, then we have deleted the root of the tree and
//     obtained a forest.
func (t *Tree[T]) SkipFilter(filter us.PredicateFilter[T]) (forest []*Tree[T]) {
	frontier := make([]*TreeNode[T], len(t.leaves))
	copy(frontier, t.leaves)

	seen := make(map[*TreeNode[T]]bool)
	var leaves []*TreeNode[T]

	f := func(n *TreeNode[T]) bool {
		return !seen[n]
	}

	for len(frontier) > 0 {
		leaf := frontier[0]
		seen[leaf] = true

		// Remove any node that has been seen from the frontier.
		frontier = us.SliceFilter(frontier, f)

		ok := filter(leaf.Data)

		if !ok {
			if leaf.parent == nil {
				// We reached the root
				frontier = frontier[1:]
			} else {
				if len(leaf.children) == 0 {
					leaves = append(leaves, leaf)
				}

				if !seen[leaf.parent] {
					frontier[0] = leaf.parent
				} else {
					frontier = frontier[1:]
				}
			}
		} else {
			children := leaf.removeNode()

			if len(children) != 0 {
				// We obtained a forest as we reached the root

				for _, child := range children {
					tree := child.ToTree()

					forest = append(forest, tree)
				}

				// We reached the root
				frontier = frontier[1:]
			} else {
				if !seen[leaf.parent] {
					frontier[0] = leaf.parent
				} else {
					frontier = frontier[1:]
				}

				t.size--
			}
		}
	}

	if len(forest) == 0 {
		t.leaves = leaves

		forest = []*Tree[T]{t}
	}

	return
}

// replaceLeafWithTree is a helper function that replaces a leaf with a tree.
//
// Parameters:
//   - at: The index of the leaf to replace.
//   - children: The children of the leaf.
//
// Behaviors:
//   - The leaf is replaced with the children.
//   - The size of the tree is updated.
func (t *Tree[T]) replaceLeafWithTree(at int, children []T) {
	leaf := t.leaves[at]

	// Make the subtree
	leaf.children = make([]*TreeNode[T], 0, len(children))
	for _, child := range children {
		node := NewTreeNode(child)
		node.parent = leaf

		leaf.children = append(leaf.children, node)
	}

	// Update the size of the tree
	t.size += len(leaf.children)

	// Replace the current leaf with the leaf's children
	if at == len(t.leaves)-1 {
		t.leaves = append(t.leaves[:at], leaf.children...)
	} else if at == 0 {
		t.leaves = append(leaf.children, t.leaves[at+1:]...)
	} else {
		t.leaves = append(t.leaves[:at], append(leaf.children, t.leaves[at+1:]...)...)
	}
}

// ProcessLeaves applies the given function to the leaves of the tree and replaces
// the leaves with the children returned by the function.
//
// Parameters:
//   - f: The function to apply to the leaves.
//
// Returns:
//   - error: An error returned by the function.
//
// Behaviors:
//   - The function is applied to the leaves in order.
//   - The function must return a slice of values of type T.
//   - If the function returns an error, the process stops and the error is returned.
//   - The leaves are replaced with the children returned by the function.
func (t *Tree[T]) ProcessLeaves(f uc.EvalManyFunc[T, T]) error {
	for i, leaf := range t.leaves {
		children, err := f(leaf.Data)
		if err != nil {
			return err
		}

		if len(children) != 0 {
			t.replaceLeafWithTree(i, children)
		}
	}

	return nil
}

// SearchNodes searches for the first node that satisfies the given filter in a BFS order.
//
// Parameters:
//   - f: The filter to apply.
//
// Returns:
//   - *treeNode[T]: A pointer to the node that satisfies the filter.
func (t *Tree[T]) SearchNodes(f us.PredicateFilter[T]) *TreeNode[T] {
	root := t.root

	Q := Queuer.NewLinkedQueue(root)

	for {
		first, ok := Q.Dequeue()
		if !ok {
			break
		}

		ok = f(first.Data)
		if ok {
			return first
		}

		for _, child := range first.children {
			Q.Enqueue(child)
		}
	}

	return nil
}

// DeleteBranchContaining deletes the branch containing the given node.
//
// Parameters:
//   - tn: The node to delete.
//
// Returns:
//   - error: An error if the node is not a part of the tree.
func (t *Tree[T]) DeleteBranchContaining(tn *TreeNode[T]) error {
	root := t.root

	child, parent, hasBranching := tn.FindBranchingPoint()
	if !hasBranching {
		if parent != root {
			return NewErrNodeNotPartOfTree()
		}

		t.Cleanup()
	}

	children := parent.DeleteChild(child)

	for _, child := range children {
		recCleanup(child)
	}

	t.leaves = t.RegenerateLeaves()

	return nil
}

// GetDirectChildren returns the direct children of the root of the tree.
//
// Returns:
//   - []*Node[T]: A slice of pointers to the direct children of the root.
func (t *Tree[T]) GetDirectChildren() []*TreeNode[T] {
	root := t.root
	if root == nil {
		return nil
	}

	return root.children
}

// PruneTree prunes the tree using the given filter.
//
// Parameters:
//   - filter: The filter to use to prune the tree.
//
// Returns:
//   - bool: True if no nodes were pruned, false otherwise.
func (t *Tree[T]) Prune(filter us.PredicateFilter[T]) bool {
	for t.Size() != 0 {
		target := t.SearchNodes(filter)
		if target == nil {
			return true
		}

		t.DeleteBranchContaining(target)
	}

	return false
}
