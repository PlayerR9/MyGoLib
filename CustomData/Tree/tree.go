package Tree

import (
	"github.com/PlayerR9/MyGoLib/ListLike/Queuer"
	"github.com/PlayerR9/MyGoLib/ListLike/Stacker"
	itff "github.com/PlayerR9/MyGoLib/Units/Common"
	ers "github.com/PlayerR9/MyGoLib/Units/Errors"
	slext "github.com/PlayerR9/MyGoLib/Utility/SliceExt"
)

// Tree is a generic data structure that represents a tree.
type Tree[T any] struct {
	// root is the root of the tree.
	root *Node[T]

	// leaves is the leaves of the tree.
	leaves []*Node[T]

	// size is the number of nodes in the tree.
	size int
}

// NewTree creates a new tree with the given root.
//
// Parameters:
//   - root: The root of the tree.
//
// Returns:
//   - *Tree[T]: A pointer to the newly created tree.
//   - error: An error of type *ers.ErrInvalidParameter if the root is nil.
//
// Behaviors:
//   - This function is expensive as it copies the root and its children.
//     Use Node[T].ToTree() to create a tree without copying the root.
func NewTree[T any](root *Node[T]) (*Tree[T], error) {
	if root == nil {
		return nil, ers.NewErrNilParameter("root")
	}

	tree := &Tree[T]{
		root:   root.Copy().(*Node[T]),
		leaves: root.Leaves(),
		size:   root.Size(),
	}

	return tree, nil
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
//
// Behaviors:
//   - Never returns nil.
func (t *Tree[T]) Root() *Node[T] {
	return t.root
}

// MergeTree merges the given tree into the tree.
//
// Parameters:
//   - tree: The tree to merge.
//
// Behaviors:
//   - The root of the other tree is added as a child of the root of the first tree.
func (t *Tree[T]) MergeTree(tree *Tree[T]) {
	if tree == nil {
		return
	}

	if t.leaves[0] == t.root {
		t.leaves = t.leaves[1:]
	}

	t.root.children = append(t.root.children, tree.root)
	t.root.parent = tree.root

	t.size += tree.size

	t.leaves = append(t.leaves, tree.leaves...)
}

// MergeTrees merges the given trees into the tree.
//
// Parameters:
//   - trees: The trees to merge.
//
// Behaviors:
//   - This function is a shortcut for calling MergeTree multiple times.
func (t *Tree[T]) MergeTrees(trees ...*Tree[T]) {
	trees = slext.SliceFilter(trees, FilterNilTree)
	if len(trees) == 0 {
		return
	}

	if t.leaves[0] == t.root {
		t.leaves = t.leaves[1:]
	}

	for _, tree := range trees {
		t.root.children = append(t.root.children, tree.root)
		t.root.parent = tree.root

		t.size += tree.size

		t.leaves = append(t.leaves, tree.leaves...)
	}
}

// AddChild adds a new child to the root of the tree.
//
// Parameters:
//   - child: The child to add.
func (t *Tree[T]) AddChild(child *Node[T]) {
	if child == nil {
		return
	}

	if t.leaves[0] == t.root {
		t.leaves = t.leaves[1:]
	}

	child.parent = t.root
	t.root.children = append(t.root.children, child)
	t.leaves = append(t.leaves, child.Leaves()...)

	t.size += child.Size()
}

// AddChildren adds new children to the root of the tree.
//
// Parameters:
//   - children: The children to add.
//
// Behaviors:
//   - This function is a shortcut for calling AddChild multiple times.
func (t *Tree[T]) AddChildren(children ...*Node[T]) {
	children = slext.SliceFilter(children, FilterNilNode)
	if len(children) == 0 {
		return
	}

	if t.leaves[0] == t.root {
		t.leaves = t.leaves[1:]
	}

	for _, child := range children {
		child.parent = t.root
		t.root.children = append(t.root.children, child)
		t.leaves = append(t.leaves, child.Leaves()...)

		t.size += child.Size()
	}
}

// GetChildren returns all the children of the tree in a DFS order.
//
// Returns:
//   - []T: A slice of the values of the nodes in the tree.
//
// Behaviors:
//   - The root is the first element in the slice.
//   - If the tree does not have a root, it returns nil.
func (t *Tree[T]) GetChildren() []T {
	children := make([]T, 0)

	S := Stacker.NewLinkedStack(t.root)

	for {
		node, err := S.Pop()
		if err != nil {
			break
		}

		children = append(children, node.Data)

		for _, child := range node.children {
			err := S.Push(child)
			if err != nil {
				panic(err)
			}
		}
	}

	return children
}

// Cleanup removes every node in the tree.
//
// Behaviors:
//   - This function is recursive and so, it is expensive.
func (t *Tree[T]) Cleanup() {
	recCleanup(t.root)
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
func (t *Tree[T]) GetLeaves() []*Node[T] {
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
func (t *Tree[T]) RegenerateLeaves() []*Node[T] {
	leaves := make([]*Node[T], 0)

	S := Stacker.NewLinkedStack(t.root)

	t.size = 0

	for {
		top, err := S.Pop()
		if err != nil {
			break
		}

		t.size++

		if len(top.children) == 0 {
			leaves = append(leaves, top)
		} else {
			for _, child := range top.children {
				err := S.Push(child)
				if err != nil {
					panic(err)
				}
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
func (t *Tree[T]) UpdateLeaves() []*Node[T] {
	newLeaves := make([]*Node[T], 0)

	S := Stacker.NewLinkedStack(t.leaves...)

	t.size -= len(t.leaves)

	for {
		top, err := S.Pop()
		if err != nil {
			break
		}

		t.size++

		if len(top.children) == 0 {
			newLeaves = append(newLeaves, top)
		} else {
			for _, child := range top.children {
				err := S.Push(child)
				if err != nil {
					panic(err)
				}
			}
		}
	}

	t.leaves = newLeaves

	return newLeaves
}

// BFSTraversal traverses the tree in BFS order.
//
// Parameters:
//   - observer: The observer function to apply to the nodes of the tree.
//
// Returns:
//   - error: An error returned by the observer.
//
// Behaviors:
//   - The traversal stops as soon as the observer returns an error.
func (t *Tree[T]) BFSTraversal(observer itff.ObserverFunc[T]) error {
	if observer == nil {
		return nil
	}

	Q := Queuer.NewLinkedQueue(t.root)

	for {
		node, err := Q.Dequeue()
		if err != nil {
			break
		}

		if err := observer(node.Data); err != nil {
			return err
		}

		for _, child := range node.children {
			err := Q.Enqueue(child)
			if err != nil {
				panic(err)
			}
		}
	}

	return nil
}

// DFSTraversal traverses the tree rooted at n in a DFS order.
//
// Parameters:
//   - observer: The observer function to apply to the nodes of the tree.
//
// Returns:
//   - error: An error returned by the observer.
//
// Behaviors:
//   - The traversal stops as soon as the observer returns an error.
func (t *Tree[T]) DFSTraversal(observer itff.ObserverFunc[T]) error {
	if observer == nil {
		return nil
	}

	S := Stacker.NewLinkedStack(t.root)

	for {
		node, err := S.Pop()
		if err != nil {
			break
		}

		if err := observer(node.Data); err != nil {
			return err
		}

		for _, child := range node.children {
			err := S.Push(child)
			if err != nil {
				panic(err)
			}
		}
	}

	return nil
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
	return recSnakeTraversal(t.root)
}

// HasChild returns true if the tree has the given child in any of its nodes
// in a BFS order.
//
// Parameters:
//   - filter: The filter to apply.
//
// Returns:
//   - bool: True if the tree has the child, false otherwise.
func (t *Tree[T]) HasChild(filter slext.PredicateFilter[T]) bool {
	if filter == nil {
		return false
	}

	Q := Queuer.NewLinkedQueue(t.root)

	for {
		node, err := Q.Dequeue()
		if err != nil {
			break
		}

		if filter(node.Data) {
			return true
		}

		for _, child := range node.children {
			err := Q.Enqueue(child)
			if err != nil {
				panic(err)
			}
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
//   - []*Node[T]: A slice of pointers to the children of the node.
func (t *Tree[T]) FilterChildren(filter slext.PredicateFilter[T]) []*Node[T] {
	if filter == nil {
		return nil
	}

	Q := Queuer.NewLinkedQueue(t.root)

	solutions := make([]*Node[T], 0)

	for {
		node, err := Q.Dequeue()
		if err != nil {
			break
		}

		if filter(node.Data) {
			solutions = append(solutions, node)
		}

		for _, child := range node.children {
			err := Q.Enqueue(child)
			if err != nil {
				panic(err)
			}
		}
	}

	return solutions
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
func (t *Tree[T]) PruneBranches(filter slext.PredicateFilter[T]) bool {
	if filter == nil {
		return false
	}

	highest, ok := recPruneFunc(filter, nil, t.root)
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
// Behaviors:
//   - This function is recursive.
func (t *Tree[T]) SkipFilter(filter slext.PredicateFilter[T]) {
	if filter == nil || t.root == nil {
		return
	}

	newChildren, amount, ok := recSkipFunc(filter, t.root)
	if ok {
		t.root = nil
		return
	}

	t.root.children = newChildren
	t.size -= amount

	// Update the parent of the children
	for i := 0; i < len(t.root.children); i++ {
		t.root.children[i].parent = t.root
	}

	// FIXME: Expensive operation. Find a better way to update the leaves.
	t.RegenerateLeaves()
}
