package Tree

import (
	"github.com/PlayerR9/MyGoLib/ListLike/Queuer"
	"github.com/PlayerR9/MyGoLib/ListLike/Stacker"
	ers "github.com/PlayerR9/MyGoLib/Units/Errors"
	itff "github.com/PlayerR9/MyGoLib/Units/Functions"
	slext "github.com/PlayerR9/MyGoLib/Utility/SliceExt"
)

// Tree is a generic data structure that represents a tree.
type Tree[T any] struct {
	// root is the root of the tree.
	root *Node[T]
}

// NewTree creates a new tree with the given root.
//
// Parameters:
//   - root: The root of the tree.
//
// Returns:
//   - *Tree[T]: A pointer to the newly created tree.
func NewTree[T any](root *Node[T]) *Tree[T] {
	return &Tree[T]{
		root: root,
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
	if t.root == nil {
		return nil
	}

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
				panic(ers.NewErrUnexpectedError(err))
			}
		}
	}

	return children
}

// Cleanup removes every node in the tree.
func (t *Tree[T]) Cleanup() {
	if t.root == nil {
		return
	}

	t.root.cleanup()

	t.root = nil
}

// GetLeaves returns all the leaves of the tree rooted at n in a BFS order.
//
// Returns:
//   - []T: A slice of the values of the nodes in the tree.
//
// Behaviors:
//   - If the tree does not have a root, it returns nil.
func (t *Tree[T]) GetLeaves() []T {
	if t.root == nil {
		return nil
	}

	leaves := make([]T, 0)

	Q := Queuer.NewLinkedQueue(t.root)

	for {
		node, err := Q.Dequeue()
		if err != nil {
			break
		}

		if len(node.children) == 0 {
			leaves = append(leaves, node.Data)
		} else {
			for _, child := range node.children {
				err := Q.Enqueue(child)
				if err != nil {
					panic(err)
				}
			}
		}
	}

	return leaves
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
	if observer == nil || t.root == nil {
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
	if observer == nil || t.root == nil {
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
	if t.root == nil {
		return make([][]T, 0)
	}

	return t.root.snakeTraversal()
}

// HasChild returns true if the tree has the given child in any of its nodes
// in a BFS order.
//
// Parameters:
//   - filter: The filter to apply.
//
// Returns:
//   - *Node[T]: The node that satisfies the filter, nil otherwise.
func (t *Tree[T]) HasChild(filter slext.PredicateFilter[T]) *Node[T] {
	if filter == nil || t.root == nil {
		return nil
	}

	Q := Queuer.NewLinkedQueue(t.root)

	for {
		node, err := Q.Dequeue()
		if err != nil {
			break
		}

		if filter(node.Data) {
			return node
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

// HasChildren returns all the children of the tree that satisfy the given filter
// in a BFS order.
//
// Parameters:
//   - filter: The filter to apply.
//
// Returns:
//   - []*Node[T]: A slice of pointers to the children of the node.
func (t *Tree[T]) HasChildren(filter slext.PredicateFilter[T]) []*Node[T] {
	if filter == nil || t.root == nil {
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

// PruneFunc removes all the children of the node that satisfy the given filter.
// The filter is a function that takes the value of a node and returns a boolean.
// If the filter returns true for a child, the child is removed along with its children.
//
// Parameters:
//   - filter: The filter to apply.
//
// Behaviors:
//   - If the root satisfies the filter, the tree is cleaned up.
//   - If the tree is empty, it does nothing.
//   - It is a recursive function.
func (t *Tree[T]) PruneFunc(filter slext.PredicateFilter[T]) {
	if filter == nil || t.root == nil {
		return
	}

	if t.root.pruneFunc(filter) {
		t.root = nil
	}
}

// SkipFunc removes all the children of the tree that satisfy the given filter
// without removing any of their children. Useful for removing unwanted nodes from the tree.
//
// Parameters:
//   - filter: The filter to apply.
//
// Behaviors:
//   - This function is recursive.
func (t *Tree[T]) SkipFunc(filter slext.PredicateFilter[T]) {
	if filter == nil || t.root == nil {
		return
	}

	newChildren, ok := t.root.skipFunc(filter)
	if ok {
		t.root = nil
		return
	}

	t.root.children = newChildren

	// Update the parent of the children
	for i := 0; i < len(t.root.children); i++ {
		t.root.children[i].parent = t.root
	}
}
