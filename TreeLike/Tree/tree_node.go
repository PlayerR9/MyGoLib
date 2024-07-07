package Tree

import (
	"slices"

	ffs "github.com/PlayerR9/MyGoLib/Formatting/FString"
	"github.com/PlayerR9/MyGoLib/ListLike/Stacker"
	lls "github.com/PlayerR9/MyGoLib/ListLike/Stacker"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
	us "github.com/PlayerR9/MyGoLib/Units/slice"
)

// TreeNode is a generic data structure that represents a node in a tree.
type TreeNode[T any] struct {
	// Parent, FirstChild, NextSibling, LastChild, PrevSibling are pointers to
	// the parent, first child, next sibling, last child, and previous sibling
	// of the node respectively.
	Parent, FirstChild, NextSibling, LastChild, PrevSibling *TreeNode[T]

	// Data is the value of the node.
	Data T
}

// FString implements the ffs.FStringer interface.
func (t *TreeNode[T]) FString(trav *ffs.Traversor, opts ...ffs.Option) error {
	err := trav.AddLine(uc.StringOf(t.Data))
	if err != nil {
		return err
	}

	children := t.GetChildren()

	err = ffs.ApplyFormMany(
		trav.GetConfig(
			ffs.WithModifiedIndent(1),
		),
		trav,
		children,
	)
	if err != nil {
		return err
	}

	return nil
}

// Copy implements the uc.Copier interface.
func (n *TreeNode[T]) Copy() uc.Copier {
	var childCopy []*TreeNode[T]

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		c_copy := c.Copy().(*TreeNode[T])
		childCopy = append(childCopy, c_copy)
	}

	d_copy := uc.CopyOf(n.Data).(T)

	nCopy := &TreeNode[T]{Data: d_copy}

	LinkWithParent(nCopy, childCopy)

	return nCopy
}

// NewTreeNode creates a new node with the given data.
//
// Parameters:
//   - data: The value of the node.
//
// Returns:
//   - *TreeNode[T]: A pointer to the newly created node.
func NewTreeNode[T any](data T) *TreeNode[T] {
	tn := &TreeNode[T]{
		Data: data,
	}

	return tn
}

// GetLastSibling returns the last sibling of the node. If it has a parent,
// it returns the last child of the parent. Otherwise, it returns the last
// sibling of the node.
//
// Returns:
//   - *TreeNode[T]: A pointer to the last sibling. The node itself if it has no next sibling.
func (n *TreeNode[T]) GetLastSibling() *TreeNode[T] {
	if n.Parent != nil {
		return n.Parent.LastChild
	} else if n.NextSibling == nil {
		return n
	}

	last_sibling := n

	for last_sibling.NextSibling != nil {
		last_sibling = last_sibling.NextSibling
	}

	return last_sibling
}

// GetFirstSibling returns the first sibling of the node. If it has a parent,
// it returns the first child of the parent. Otherwise, it returns the first
// sibling of the node.
//
// Returns:
//   - *TreeNode[T]: A pointer to the first sibling. The node itself if it has no previous sibling.
func (n *TreeNode[T]) GetFirstSibling() *TreeNode[T] {
	if n.Parent != nil {
		return n.Parent.FirstChild
	} else if n.PrevSibling == nil {
		return n
	}

	first_sibling := n

	for first_sibling.PrevSibling != nil {
		first_sibling = first_sibling.PrevSibling
	}

	return first_sibling
}

// IsLeaf returns true if the node is a leaf.
//
// Returns:
//   - bool: True if the node is a leaf, false otherwise.
func (n *TreeNode[T]) IsLeaf() bool {
	return n.FirstChild == nil
}

// IsRoot returns true if the node does not have a parent.
//
// Returns:
//   - bool: True if the node is the root, false otherwise.
func (n *TreeNode[T]) IsRoot() bool {
	return n.Parent == nil
}

// Leaves returns all the leaves of the tree rooted at n.
//
// Returns:
//   - leaves: A slice of pointers to the leaves of the tree.
//
// Behaviors:
//   - The leaves are returned in the order of a DFS traversal.
func (n *TreeNode[T]) Leaves() (leaves []*TreeNode[T]) {
	S := Stacker.NewLinkedStack(n)

	for {
		top, ok := S.Pop()
		if !ok {
			break
		}

		ok = top.IsLeaf()
		if ok {
			leaves = append(leaves, top)
		} else {
			children := top.GetChildren()

			for _, child := range children {
				S.Push(child)
			}
		}
	}

	return
}

// ToTree converts the node to a tree.
//
// Returns:
//   - *Tree[T]: A pointer to the tree.
func (n *TreeNode[T]) ToTree() *Tree[T] {
	if n.FirstChild == nil {
		return &Tree[T]{
			root:   n,
			leaves: []*TreeNode[T]{n},
			size:   1,
		}
	} else {
		return &Tree[T]{
			root:   n,
			leaves: n.Leaves(),
			size:   n.Size(),
		}
	}
}

// AddChild adds a new child to the node with the given data.
//
// Parameters:
//   - child: The child to add.
//
// Behaviors:
//   - If the child is nil, it does nothing.
func (n *TreeNode[T]) AddChild(child *TreeNode[T]) {
	if child == nil {
		return
	}

	// Make sure the child is not linked to any other node
	child.NextSibling = nil
	child.PrevSibling = nil

	last_child := n.LastChild

	if last_child == nil {
		n.FirstChild = child
	} else {
		last_child.NextSibling = child
		child.PrevSibling = last_child
	}

	child.Parent = n
	n.LastChild = child
}

// AddChildren adds zero or more children to the node.
//
// Parameters:
//   - children: The children to add.
//
// Behaviors:
//   - This is just a more efficient way to add multiple children.
func (n *TreeNode[T]) AddChildren(children ...*TreeNode[T]) {
	children = us.FilterNilValues(children)
	if len(children) == 0 {
		return
	}

	// Deal with the first child
	first_child := children[0]

	first_child.NextSibling = nil
	first_child.PrevSibling = nil

	last_child := n.LastChild

	if last_child == nil {
		n.FirstChild = first_child
	} else {
		last_child.NextSibling = first_child
		first_child.PrevSibling = last_child
	}

	first_child.Parent = n
	n.LastChild = first_child

	// Deal with the rest of the children
	for i := 1; i < len(children); i++ {
		child := children[i]

		child.NextSibling = nil
		child.PrevSibling = nil

		last_child := n.LastChild
		last_child.NextSibling = child
		child.PrevSibling = last_child

		child.Parent = n
		n.LastChild = child
	}
}

// GetChildren returns all the children of the node.
// If the node has no children, it returns nil.
//
// Returns:
//   - []*TreeNode[T]: A slice of pointers to the children of the node.
func (n *TreeNode[T]) GetChildren() []*TreeNode[T] {
	var children []*TreeNode[T]

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		children = append(children, c)
	}

	return children
}

// FindBranchingPoint returns the first node in the path from n to the root
// such that has more than one sibling.
//
// Returns:
//   - *treeNode[T]: A pointer to the one node before the branching point.
//   - *treeNode[T]: A pointer to the branching point.
//   - bool: True if the node has a branching point, false otherwise.
//
// Behaviors:
//   - If there is no branching point, it returns the root of the tree.
func (tn *TreeNode[T]) FindBranchingPoint() (*TreeNode[T], *TreeNode[T], bool) {
	if tn.Parent == nil {
		return nil, tn, false
	}

	node := tn
	parent := tn.Parent

	var hasBranchingPoint bool

	for parent.Parent != nil && !hasBranchingPoint {
		if parent.FirstChild != parent.LastChild {
			hasBranchingPoint = true
		} else {
			node = parent
			parent = parent.Parent
		}
	}

	return node, parent, hasBranchingPoint
}

// HasChild returns true if the node has the given child.
//
// Parameters:
//   - target: The child to check for.
//
// Returns:
//   - bool: True if the node has the child, false otherwise.
func (n *TreeNode[T]) HasChild(target *TreeNode[T]) bool {
	if target == nil || n.FirstChild == nil {
		return false
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
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
//   - []*treeNode[T]: A slice of pointers to the children of the node.
//
// Behaviors:
//   - If the node has no children, it returns nil.
func (n *TreeNode[T]) DeleteChild(target *TreeNode[T]) []*TreeNode[T] {
	children := n.delete_child(target)
	DelinkWithParent(n, children)

	return children
}

// DeleteChild removes the given child from the children of the node.
//
// Parameters:
//   - target: The child to remove.
//
// Returns:
//   - []*treeNode[T]: A slice of pointers to the children of the node.
//
// Behaviors:
//   - If the node has no children, it returns nil.
func (n *TreeNode[T]) delete_child(target *TreeNode[T]) []*TreeNode[T] {
	ok := n.HasChild(target)
	if !ok {
		// Nothing to delete
		return nil
	}

	prev_sibling := target.PrevSibling

	if prev_sibling != nil {
		target.PrevSibling.NextSibling = target.NextSibling
	}

	next_sibling := target.NextSibling

	if next_sibling != nil {
		target.NextSibling.PrevSibling = target.PrevSibling
	}

	if target == n.FirstChild {
		n.FirstChild = target.NextSibling

		if target.NextSibling == nil {
			n.LastChild = nil
		}
	} else if target == n.LastChild {
		n.LastChild = target.PrevSibling
	}

	target.Parent = nil
	target.PrevSibling = nil
	target.NextSibling = nil

	children := target.GetChildren()
	return children
}

// Size returns the number of nodes in the tree rooted at n.
//
// Returns:
//   - size: The number of nodes in the tree.
//
// Behaviors:
//   - This function is expensive since size is not stored.
func (n *TreeNode[T]) Size() (size int) {
	S := Stacker.NewLinkedStack(n)

	for {
		top, ok := S.Pop()
		if !ok {
			break
		}

		size++

		for c := top.FirstChild; c != nil; c = c.NextSibling {
			S.Push(c)
		}
	}

	return
}

// GetAncestors returns all the ancestors of the node.
//
// This excludes the node itself.
//
// Returns:
//   - ancestors: A slice of pointers to the ancestors of the node.
//
// Behaviors:
//   - The ancestors are returned in the opposite order of a DFS traversal.
//     Therefore, the first element is the parent of the node.
func (n *TreeNode[T]) GetAncestors() (ancestors []*TreeNode[T]) {
	for node := n; node.Parent != nil; node = node.Parent {
		ancestors = append(ancestors, node.Parent)
	}

	slices.Reverse(ancestors)

	return
}

// IsChildOf returns true if the node is a child of the parent.
//
// Parameters:
//   - target: The target parent to check for.
//
// Returns:
//   - bool: True if the node is a child of the parent, false otherwise.
func (n *TreeNode[T]) IsChildOf(target *TreeNode[T]) bool {
	if target == nil {
		return false
	}

	parents := target.GetAncestors()

	for node := n; node.Parent != nil; node = node.Parent {
		ok := slices.Contains(parents, node.Parent)
		if ok {
			return true
		}
	}

	return false
}

// removeNode removes the node from the tree and shifts the children up
// in the space occupied by the node.
//
// Returns:
//   - []*treeNode[T]: A slice of pointers to the children of the node if
//     the node is the root. Nil otherwise.
func (n *TreeNode[T]) removeNode() []*TreeNode[T] {
	prev_sibling := n.PrevSibling
	next_sibling := n.NextSibling

	var sub_roots []*TreeNode[T]

	if n.Parent == nil {
		sub_roots = n.GetChildren()
	} else {
		children := n.Parent.delete_child(n)

		for _, child := range children {
			child.Parent = n.Parent
		}
	}

	if prev_sibling != nil {
		prev_sibling.NextSibling = next_sibling
	} else {
		n.Parent.FirstChild = next_sibling
	}

	if next_sibling != nil {
		next_sibling.PrevSibling = prev_sibling
	} else {
		n.Parent.LastChild = prev_sibling
	}

	n.Parent = nil
	n.PrevSibling = nil
	n.NextSibling = nil

	DelinkWithParent(n, sub_roots)

	return sub_roots
}

// GetData is a getter for the data of the node.
//
// Returns:
//   - T: The data of the node.
func (n *TreeNode[T]) GetData() T {
	return n.Data
}

// GetBranch works like GetAncestors but includes the node itself.
//
// The nodes are returned as a slice where [0] is the root node
// and [len(branch)-1] is the leaf node.
//
// Returns:
//   - []*TreeNode[T]: A slice of pointers to the nodes in the branch.
func (n *TreeNode[T]) GetBranch() *Branch[T] {
	size := 1

	branch := &Branch[T]{
		to_node: n,
	}

	node := n

	for node.Parent != nil {
		size++
		node = node.Parent
	}

	branch.from_node = node
	branch.size = size

	return branch
}

// LinkWithParent links the parent with the children. It also links the children
// with each other.
//
// Parameters:
//   - parent: The parent node.
//   - children: The children nodes.
//
// Returns:
//   - []*TreeNode[T]: The linked children. (Without the nil values)
//
// Behaviors:
//   - If the parent has no children, it does nothing. Nil values are filtered out.
func LinkWithParent[T any](parent *TreeNode[T], children []*TreeNode[T]) []*TreeNode[T] {
	children = LinkSiblings(children)
	if len(children) == 0 {
		// Do nothing.
		return nil
	}

	for _, child := range children {
		child.Parent = parent
	}

	if parent == nil {
		// Do nothing.
		return children
	}

	parent.FirstChild = children[0]
	parent.LastChild = children[len(children)-1]

	return children
}

// LinkSiblings links the siblings with each other. It also sets the prev and
// next siblings of the first and last nodes to nil.
//
// Parameters:
//   - nodes: The nodes to link.
//
// Returns:
//   - []*TreeNode[T]: The linked nodes. (Without the nil values)
//
// Behaviors:
//   - If the nodes slice is empty, it does nothing. Nil values are filtered out.
func LinkSiblings[T any](nodes []*TreeNode[T]) []*TreeNode[T] {
	nodes = us.FilterNilValues(nodes)
	if len(nodes) == 0 {
		// Do nothing.
		return nil
	}

	nodes[0].PrevSibling = nil
	nodes[len(nodes)-1].NextSibling = nil

	if len(nodes) == 1 {
		// Do nothing.
		return nodes
	}

	for i := 0; i < len(nodes)-1; i++ {
		current := nodes[i]

		current.NextSibling = nodes[i+1]
	}

	for i := 1; i < len(nodes); i++ {
		current := nodes[i]

		current.PrevSibling = nodes[i-1]
	}

	return nodes
}

// DelinkWithParent delinks the parent with the children. It also delinks the children
// with each other.
//
// Parameters:
//   - parent: The parent node.
//   - children: The children nodes.
//
// Behaviors:
//   - If the parent has no children, it does nothing.
func DelinkWithParent[T any](parent *TreeNode[T], children []*TreeNode[T]) {
	if len(children) == 0 {
		return
	}

	for i := 0; i < len(children); i++ {
		child := children[i]

		if child == nil {
			continue
		}

		child.PrevSibling = nil
		child.NextSibling = nil
		child.Parent = nil
	}

	if parent != nil {
		parent.FirstChild = nil
		parent.LastChild = nil
	}
}

// Cleanup removes every child of the node in a DFS traversal.
//
// Behaviors:
//   - The node itself is not removed.
//   - It does not use recursion.
func (tn *TreeNode[T]) Cleanup() {
	type Helper struct {
		Prev *TreeNode[T]
		Curr *TreeNode[T]
	}

	h := &Helper{
		Prev: nil,
		Curr: tn,
	}

	S := lls.NewLinkedStack(h)

	for {
		h, ok := S.Pop()
		if !ok {
			break
		}

		for c := h.Curr.FirstChild; c != nil; c = c.NextSibling {
			h := &Helper{
				Prev: c.PrevSibling,
				Curr: c,
			}

			S.Push(h)
		}

		if h.Prev != nil {
			h.Prev.NextSibling = nil
			h.Prev.PrevSibling = nil
		}

		h.Curr.FirstChild = nil
		h.Curr.LastChild = nil
		h.Curr.Parent = nil
	}

	tn.NextSibling = nil
	tn.PrevSibling = nil
}
