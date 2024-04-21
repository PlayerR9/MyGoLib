package Node

import (
	Queue "github.com/PlayerR9/MyGoLib/ListLike/Queue"
	Stack "github.com/PlayerR9/MyGoLib/ListLike/Stack"
	slext "github.com/PlayerR9/MyGoLib/Utility/SliceExt"

	itff "github.com/PlayerR9/MyGoLibUnits/Functions"
)

// Node is a generic data structure that represents a node in a tree.
type Node[T any] struct {
	// Data is the value of the node.
	Data T

	// parent is the parent of the node.
	parent *Node[T]

	// firstChild is the first child of the node.
	firstChild *Node[T]

	// nextSibling is the next sibling of the node.
	nextSibling *Node[T]
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
		Data: data,
	}
}

// AddChild adds a new child to the node with the given data.
//
// Parameters:
//
//   - data: The value of the new child.
func (n *Node[T]) AddChild(data T) {
	node := &Node[T]{
		Data:   data,
		parent: n,
	}

	if n.firstChild == nil {
		n.firstChild = node
	} else {
		current := n.firstChild

		for current.nextSibling != nil {
			current = current.nextSibling
		}

		current.nextSibling = node
	}
}

// AddChildren adds zero or more children to the node with the given data.
//
// Parameters:
//
//   - children: The values of the new children.
func (n *Node[T]) AddChildren(children ...T) {
	if len(children) == 0 {
		return
	}

	node := &Node[T]{
		Data:   children[0],
		parent: n,
	}

	if n.firstChild == nil {
		n.firstChild = node
	} else {
		child := n.firstChild

		for child.nextSibling != nil {
			child = child.nextSibling
		}

		child.nextSibling = node
	}

	current := node

	for _, data := range children[1:] {
		node := &Node[T]{
			Data:   data,
			parent: n,
		}

		current.nextSibling = node
		current = node
	}
}

// GetChildren returns all the children of the node.
// If the node has no children, it returns nil.
//
// Returns:
//
//   - []*Node[T]: A slice of pointers to the children of the node.
func (n *Node[T]) GetChildren() []*Node[T] {
	children := make([]*Node[T], 0)

	current := n.firstChild

	for current != nil {
		children = append(children, current)
		current = current.nextSibling
	}

	return children
}

// Cleanup removes the node from the tree; including all its children.
func (n *Node[T]) Cleanup() {
	n.parent = nil

	if n.firstChild != nil {
		n.firstChild.Cleanup()
		n.firstChild = nil
	}

	for s := n.nextSibling; s != nil; s = s.nextSibling {
		s.Cleanup()
	}

	n.nextSibling = nil
}

// GetLeaves returns all the leaves of the tree rooted at n.
// The leaves are returned in the order of a breadth-first traversal.
//
// Returns:
//
//   - []*Node[T]: A slice of pointers to the leaves of the tree.
func (n *Node[T]) GetLeaves() []*Node[T] {
	leaves := make([]*Node[T], 0)

	Q := Queue.NewLinkedQueue(n)

	for !Q.IsEmpty() {
		node := Q.MustDequeue()

		if node.firstChild == nil {
			leaves = append(leaves, node)
		} else {
			current := node.firstChild

			for current != nil {
				Q.Enqueue(current)

				current = current.nextSibling
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
//
// Returns:
//
//   - error: An error returned by the observer.
func (n *Node[T]) BFSTraversal(observer itff.ObserverFunc[T]) error {
	Q := Queue.NewLinkedQueue(n)

	for !Q.IsEmpty() {
		node := Q.MustDequeue()

		if err := observer(node.Data); err != nil {
			return err
		}

		for child := node.firstChild; child != nil; child = child.nextSibling {
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
//
// Returns:
//
//   - error: An error returned by the observer.
func (n *Node[T]) DFSTraversal(observer itff.ObserverFunc[T]) error {
	S := Stack.NewLinkedStack(n)

	for !S.IsEmpty() {
		node := S.MustPop()

		if err := observer(node.Data); err != nil {
			return err
		}

		for child := node.firstChild; child != nil; child = child.nextSibling {
			S.Push(child)
		}
	}

	return nil
}

// SnakeTraversal returns all the paths from n to the leaves of the tree rooted at n.
// The paths are returned in the order of a breadth-first traversal.
//
// Returns:
//
//   - [][]T: A slice of slices of the values of the nodes in the paths.
func (n *Node[T]) SnakeTraversal() [][]T {
	if n.firstChild == nil {
		return [][]T{
			{n.Data},
		}
	}

	result := make([][]T, 0)

	for child := n.firstChild; child != nil; child = child.nextSibling {
		for _, tmp := range child.SnakeTraversal() {
			result = append(result, append([]T{n.Data}, tmp...))
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
		if node.parent.firstChild.nextSibling != nil {
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

	for c := n.firstChild; c != nil; c = c.nextSibling {
		if c == child {
			return true
		}
	}

	return false
}

// DeleteChild removes the given child from the children of the node.
// No op if the child is nil or not a child of the node.
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

	for c := n.firstChild; c != nil; c = c.nextSibling {
		if c != child {
			continue
		}

		if c.nextSibling != nil {
			c.nextSibling.parent = n
		}

		break
	}

	child.parent = nil
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
func (n *Node[T]) PruneFunc(filter slext.PredicateFilter[T]) bool {
	if filter(n.Data) {
		return true
	}

	for child := n.firstChild; child != nil; child = child.nextSibling {
		if child.PruneFunc(filter) {
			n.DeleteChild(child)
		}
	}

	return false
}
