package Node

import (
	"fmt"
	"slices"

	Queue "github.com/PlayerR9/MyGoLib/ListLike/Queue"
	Stack "github.com/PlayerR9/MyGoLib/ListLike/Stack"
)

// Node is a generic data structure that represents a node in a tree.
type Node[T any] struct {
	// Data is the value of the node.
	Data T

	// parent is the parent of the node.
	parent *Node[T]

	// children are the children of the node.
	children []*Node[T]
}

// String is a method of fmt.Stringer interface.
// It should only be used for debugging and logging purposes.
//
// Returns:
//
//   - string: A string representation of the node.
func (n *Node[T]) String() string {
	if n == nil {
		return "Node[nil]"
	}

	return fmt.Sprintf("Node[%v]", n.Data)
}

// NewNode creates a new node with the given data.
//
// Parameters:
//
//   - data: The value of the node.
//
// Returns:
//
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
//
//   - data: The value of the new child.
func (n *Node[T]) AddChild(data T) {
	node := &Node[T]{
		Data:     data,
		parent:   n,
		children: make([]*Node[T], 0),
	}

	n.children = append(n.children, node)
}

// AddChildren adds zero or more children to the node with the given data.
//
// Parameters:
//
//   - children: The values of the new children.
func (n *Node[T]) AddChildren(children ...T) {
	if len(children) == 0 {
		return
	} else if n.children == nil {
		n.children = make([]*Node[T], 0)
	}

	var node *Node[T]

	for _, data := range children {
		node = NewNode(data)
		node.parent = n

		n.children = append(n.children, node)
	}
}

// GetChildren returns all the children of the node.
// If the node has no children, it returns nil.
//
// Returns:
//
//   - []*Node[T]: A slice of pointers to the children of the node.
func (n *Node[T]) GetChildren() []*Node[T] {
	return n.children
}

// Cleanup removes the node from the tree; including all its children.
func (n *Node[T]) Cleanup() {
	n.parent = nil

	for _, child := range n.children {
		if child != nil {
			child.Cleanup()
		}
	}

	for i := range n.children {
		n.children[i] = nil
	}

	n.children = nil
}

// GetLeaves returns all the leaves of the tree rooted at n.
// The leaves are returned in the order of a breadth-first traversal.
//
// Returns:
//
//   - []*Node[T]: A slice of pointers to the leaves of the tree.
func (n *Node[T]) GetLeaves() []*Node[T] {
	leaves := make([]*Node[T], 0)

	Q := Queue.NewArrayQueue(n)

	for !Q.IsEmpty() {
		node, _ := Q.Dequeue()

		if len(node.children) == 0 {
			leaves = append(leaves, node)
		} else {
			for _, child := range node.children {
				Q.Enqueue(child)
			}
		}
	}

	return leaves
}

// BFSTraversal traverses the tree rooted at n in a breadth-first manner.
// The traversal stops when the observer returns an error.
//
// Parameters:
//
//   - observer: A function that takes the value of a node and returns an error.
func (n *Node[T]) BFSTraversal(observer func(T) error) error {
	Q := Queue.NewArrayQueue(n)

	for !Q.IsEmpty() {
		node, _ := Q.Dequeue()

		if err := observer(node.Data); err != nil {
			return err
		}

		for _, child := range node.children {
			Q.Enqueue(child)
		}
	}

	return nil
}

// DFSTraversal traverses the tree rooted at n in a depth-first manner.
// The traversal stops when the observer returns an error.
//
// Parameters:
//
//   - observer: A function that takes the value of a node and returns an error.
func (n *Node[T]) DFSTraversal(observer func(T) error) error {
	S := Stack.NewArrayStack(n)

	for !S.IsEmpty() {
		node, _ := S.Pop()

		if err := observer(node.Data); err != nil {
			return err
		}

		for _, child := range node.children {
			S.Push(child)
		}
	}

	return nil
}

// SnakeTraversal returns all the paths from n to the leaves of the tree rooted at n.
// The paths are returned in the order of a depth-first traversal.
//
// Returns:
//
//   - [][]T: A slice of slices of the values of the nodes in the paths.
func (n *Node[T]) SnakeTraversal() [][]T {
	type StackNode struct {
		node *Node[T]
		path []T
	}

	S := Stack.NewLinkedStack(StackNode{
		node: n,
		path: []T{n.Data},
	})
	var result [][]T

	for !S.IsEmpty() {
		top, _ := S.Pop()

		if len(top.node.children) == 0 {
			result = append(result, top.path)
			continue
		}

		for _, child := range top.node.children {
			newPath := make([]T, len(top.path))
			copy(newPath, top.path)

			S.Push(StackNode{
				node: child,
				path: append(newPath, child.Data),
			})
		}
	}

	return result
}

// FindBranchingPoint returns the first node in the path from n to the root
// such that has more than one sibling.
// If no such node is found, it returns nil.
//
// Returns:
//
//   - *Node[T]: A pointer to the branching point.
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
//
//   - child: The child to check for.
//
// Returns:
//
//   - bool: True if the node has the child, false otherwise.
func (n *Node[T]) HasChild(child *Node[T]) bool {
	if child == nil {
		return false
	}

	for _, node := range n.children {
		if node == child {
			return true
		}
	}

	return false
}

// DeleteChild removes the given child from the children of the node.
// No op if the child is nil or not a child of the node.
//
// This does not remove the child from the tree; it only removes the
// parent-child relationship.
//
// Parameters:
//
//   - child: The child to remove.
//
// Returns:
//
//   - []*Node[T]: A slice of pointers to the children of the node.
func (n *Node[T]) DeleteChild(child *Node[T]) {
	if child == nil {
		return
	}

	index := slices.Index(n.children, child)
	if index == -1 {
		return
	}

	child.parent = nil

	n.children = slices.Delete(n.children, index, index+1)
}

// Parent is a getter for the parent of the node.
// If the node has no parent, it returns nil.
//
// Returns:
//
//   - *Node[T]: A pointer to the parent of the node.
func (n *Node[T]) Parent() *Node[T] {
	return n.parent
}

// PruneFunc removes all the children of the node that satisfy the given filter.
// The filter is a function that takes the value of a node and returns a boolean.
// If the filter returns true for a child, the child is removed.
//
// Parameters:
//
//   - filter: The filter to apply.
//
// Returns:
//
//   - bool: True if the node satisfies the filter, false otherwise.
func (n *Node[T]) PruneFunc(filter func(T) bool) bool {
	if filter(n.Data) {
		return true
	}

	var childrenToDelete []*Node[T]

	for _, child := range n.children {
		if child.PruneFunc(filter) {
			childrenToDelete = append(childrenToDelete, child)
		}
	}

	for _, child := range childrenToDelete {
		n.DeleteChild(child)
	}

	return false
}
