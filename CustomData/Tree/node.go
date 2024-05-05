package Tree

import (
	"slices"

	"github.com/PlayerR9/MyGoLib/ListLike/Stacker"

	intf "github.com/PlayerR9/MyGoLib/Units/Interfaces"
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

// Copy returns a deep copy of the node.
//
// Returns:
//   - *Node[T]: A pointer to the deep copy of the node.
//
// Behaviors:
//   - This function is recursive.
//   - The parent is not copied.
//   - The data is shallow copied.
func (n *Node[T]) Copy() intf.Copier {
	node := &Node[T]{
		Data:     n.Data,
		parent:   nil,
		children: make([]*Node[T], 0, len(n.children)),
	}

	for _, child := range n.children {
		node.children = append(node.children, child.Copy().(*Node[T]))
	}

	return node
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

// IsLeaf returns true if the node is a leaf.
//
// Returns:
//   - bool: True if the node is a leaf, false otherwise.
func (n *Node[T]) IsLeaf() bool {
	return len(n.children) == 0
}

// IsRoot returns true if the node does not have a parent.
//
// Returns:
//   - bool: True if the node is the root, false otherwise.
func (n *Node[T]) IsRoot() bool {
	return n.parent == nil
}

// Leaves returns all the leaves of the tree rooted at n.
//
// Returns:
//   - []*Node[T]: A slice of pointers to the leaves of the tree.
//
// Behaviors:
//   - The leaves are returned in the order of a DFS traversal.
func (n *Node[T]) Leaves() []*Node[T] {
	leaves := make([]*Node[T], 0)

	S := Stacker.NewLinkedStack(n)

	for {
		top, err := S.Pop()
		if err != nil {
			break
		}

		if top.IsLeaf() {
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

	return leaves
}

// ToTree converts the node to a tree.
//
// Returns:
//   - *Tree[T]: A pointer to the tree.
func (n *Node[T]) ToTree() *Tree[T] {
	return &Tree[T]{
		root: n,
	}
}

// AddChildren adds a new child to the node with the given data.
//
// Parameters:
//   - children: The children to add.
func (n *Node[T]) AddChildren(children ...*Node[T]) {
	if len(children) == 0 {
		return
	}

	for _, child := range children {
		if child == nil {
			continue
		}

		child.parent = n

		n.children = append(n.children, child)
	}
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

// Size returns the number of nodes in the tree rooted at n.
//
// Returns:
//   - int: The number of nodes in the tree.
//
// Behaviors:
//   - This function is expensive since size is not stored.
func (n *Node[T]) Size() int {
	size := 0

	S := Stacker.NewLinkedStack(n)

	for {
		top, err := S.Pop()
		if err != nil {
			break
		}

		size++

		for _, child := range top.children {
			err := S.Push(child)
			if err != nil {
				panic(err)
			}
		}
	}

	return size
}

// GetAncestors returns all the ancestors of the node.
//
// Returns:
//   - []*Node[T]: A slice of pointers to the ancestors of the node.
//
// Behaviors:
//   - The ancestors are returned in the opposite order of a DFS traversal.
//     Therefore, the first element is the parent of the node.
func (n *Node[T]) GetAncestors() []*Node[T] {
	ancestors := make([]*Node[T], 0)

	for node := n; node.parent != nil; node = node.parent {
		ancestors = append(ancestors, node.parent)
	}

	slices.Reverse(ancestors)

	return ancestors
}

// IsChildOf returns true if the node is a child of the parent.
//
// Parameters:
//   - target: The target parent to check for.
//
// Returns:
//   - bool: True if the node is a child of the parent, false otherwise.
func (n *Node[T]) IsChildOf(target *Node[T]) bool {
	if target == nil {
		return false
	}

	parents := target.GetAncestors()

	for node := n; node.parent != nil; node = node.parent {
		if slices.Contains(parents, node.parent) {
			return true
		}
	}

	return false
}
