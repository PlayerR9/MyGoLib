package Tree

import (
	"slices"

	slext "github.com/PlayerR9/MyGoLib/Utility/SliceExt"
)

// Node is a generic data structure that represents a node in a tree.
type Node[T any] struct {
	// Data is the value of the node.
	Data T

	// parent is the parent of the node.
	parent *Node[T]

	// children is the children of the node.
	children []*Node[T]
}

// NewNode creates a new node with the given data.
//
// Parameters:
//   - data: The value of the node.
//
// Returns:
//   - *Node[T]: A pointer to the newly created node.
func NewNode[T any](data T) *Node[T] {
	return &Node[T]{
		Data:     data,
		children: make([]*Node[T], 0),
	}
}

// AddChild adds a new child to the node with the given data.
//
// Parameters:
//   - data: The value of the new child.
func (n *Node[T]) AddChild(child *Node[T]) {
	if child == nil {
		return
	}

	child.parent = n

	n.children = append(n.children, child)
}

// MakeChildren adds zero or more children to the node with the given data.
//
// Parameters:
//   - children: The values of the new children.
func (n *Node[T]) MakeChildren(children ...T) {
	for _, data := range children {
		child := &Node[T]{
			Data:   data,
			parent: n,
		}

		n.children = append(n.children, child)
	}
}

// GetChildren returns all the children of the node.
// If the node has no children, it returns nil.
//
// Returns:
//   - []*Node[T]: A slice of pointers to the children of the node.
func (n *Node[T]) GetChildren() []*Node[T] {
	if len(n.children) == 0 {
		return nil
	}

	children := make([]*Node[T], 0, len(n.children))
	copy(children, n.children)

	return children
}

// cleanup is a helper function that removes every node in the tree rooted at n.
//
// Behaviors:
//   - This function is recursive.
func (n *Node[T]) cleanup() {
	n.parent = nil

	for _, child := range n.children {
		child.cleanup()
	}

	for i := range n.children {
		n.children[i] = nil
	}

	n.children = nil
}

// snakeTraversal is an helper function that returns all the paths
// from n to the leaves of the tree rooted at n.
//
// Returns:
//   - [][]T: A slice of slices of the values of the nodes in the paths.
//
// Behaviors:
//   - The paths are returned in the order of a BFS traversal.
//   - It is a recursive function.
func (n *Node[T]) snakeTraversal() [][]T {
	if len(n.children) == 0 {
		return [][]T{
			{n.Data},
		}
	}

	result := make([][]T, 0)

	for _, child := range n.children {
		subResults := child.snakeTraversal()

		for _, tmp := range subResults {
			result = append(result, append([]T{n.Data}, tmp...))
		}
	}

	return result
}

// FindBranchingPoint returns the first node in the path from n to the root
// such that has more than one sibling.
//
// Returns:
//   - *Node[T]: A pointer to the branching point. Nil if no such node is found.
func (n *Node[T]) FindBranchingPoint() *Node[T] {
	for node := n; node.parent != nil; node = node.parent {
		if len(node.parent.children) > 1 {
			return node.parent
		}
	}

	return nil
}

// HasChild returns true if the node has the given child.
//
// Parameters:
//   - target: The child to check for.
//
// Returns:
//   - bool: True if the node has the child, false otherwise.
func (n *Node[T]) HasChild(target *Node[T]) bool {
	if target == nil {
		return false
	}

	for _, c := range n.children {
		if c == target {
			return true
		}
	}

	return false
}

// DeleteChild removes the given child from the children of the node.
//
// Parameters:
//   - target: The child to remove.
//
// Returns:
//   - []*Node[T]: A slice of pointers to the children of the node.
//
// Behaviors:
//   - If the node has no children, it returns nil.
func (n *Node[T]) DeleteChild(target *Node[T]) []*Node[T] {
	if target == nil {
		// No target to delete
		return nil
	}

	switch len(n.children) {
	case 0:
		// No children to delete
	case 1:
		if n.children[0] == target {
			n.children = nil
			target.parent = nil

			return target.children
		}
	default:
		index := slices.Index(n.children, target)
		if index == -1 {
			// target is not a child of n
			return nil
		}

		n.children = append(n.children[:index], n.children[index+1:]...)
		target.parent = nil

		return target.children
	}

	return nil
}

// Parent is a getter for the parent of the node.
//
// Returns:
//   - *Node[T]: A pointer to the parent of the node. Nil if the node has no parent.
func (n *Node[T]) Parent() *Node[T] {
	return n.parent
}

// pruneFunc is an helper function that removes all the children of the
// node that satisfy the given filter including all of their children.
//
// Parameters:
//   - filter: The filter to apply.
//
// Returns:
//   - bool: True if the node satisfies the filter, false otherwise.
//
// Behaviors:
//   - This function is recursive.
func (n *Node[T]) pruneFunc(filter slext.PredicateFilter[T]) bool {
	if filter(n.Data) {
		// Delete all children
		n.cleanup()

		return true
	}

	top := 0

	for i := 0; i < len(n.children); i++ {
		if n.children[i].pruneFunc(filter) {
			n.children[i] = nil
		} else {
			n.children[top] = n.children[i]
			top++
		}
	}

	n.children = n.children[:top]

	return false
}

// skipFunc is an helper function that removes all the children of the
// node that satisfy the given filter without removing their children.
//
// Parameters:
//   - filter: The filter to apply.
//
// Returns:
//   - []*Node[T]: A slice of pointers to the children of the node.
//   - bool: True if the node satisfies the filter, false otherwise.
func (n *Node[T]) skipFunc(filter slext.PredicateFilter[T]) ([]*Node[T], bool) {
	if filter(n.Data) {
		n.parent = nil

		return n.children, true
	}

	newChildren := make([]*Node[T], 0)

	for _, child := range n.children {
		subChildren, ok := child.skipFunc(filter)
		if !ok {
			newChildren = append(newChildren, child)
			continue
		}

		if len(subChildren) > 0 {
			newChildren = append(newChildren, subChildren...)
		}
	}

	n.children = newChildren

	// Update the parent of the children
	for i := 0; i < len(n.children); i++ {
		n.children[i].parent = n
	}

	return nil, false
}
