package Tree

import (
	"slices"

	ffs "github.com/PlayerR9/MyGoLib/Formatting/FString"
	lls "github.com/PlayerR9/MyGoLib/ListLike/Stacker"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
	us "github.com/PlayerR9/MyGoLib/Units/slice"
)

// StatusNodeIterator is a iterator that iterates over the children of a tree node.
type StatusNodeIterator[S uc.Enumer, T any] struct {
	// node is the node whose children are being iterated over.
	node *StatusNode[S, T]

	// current_node is the current node being iterated over.
	current_node *StatusNode[S, T]
}

// Size implements the common.Iterater interface.
func (t *StatusNodeIterator[S, T]) Consume() (Noder, error) {
	if t.current_node == nil {
		return nil, uc.NewErrExhaustedIter()
	}

	n := t.current_node

	t.current_node = t.current_node.NextSibling

	return n, nil
}

// Restart implements the common.Iterater interface.
func (t *StatusNodeIterator[S, T]) Restart() {
	t.current_node = t.node.FirstChild
}

// NewStatusNodeIterator creates a new iterator for the given node.
//
// Parameters:
//   - node: The node whose children are being iterated over.
//
// Returns:
//   - *StatusNodeIterator[S, T]: A pointer to the iterator. Nil if the node is nil.
func NewStatusNodeIterator[S uc.Enumer, T any](node *StatusNode[S, T]) *StatusNodeIterator[S, T] {
	if node == nil {
		return nil
	}

	ti := &StatusNodeIterator[S, T]{
		node:         node,
		current_node: node.FirstChild,
	}

	return ti
}

// StatusNode is a generic data structure that represents a node in a tree.
type StatusNode[S uc.Enumer, T any] struct {
	// Parent, FirstChild, NextSibling, LastChild, PrevSibling are pointers to
	// the parent, first child, next sibling, last child, and previous sibling
	// of the node respectively.
	Parent, FirstChild, NextSibling, LastChild, PrevSibling *StatusNode[S, T]

	// Data is the value of the node.
	Data T

	// Status is the status of the node.
	Status S
}

// Iterator implements Noder.
func (t *StatusNode[S, T]) Iterator() uc.Iterater[Noder] {
	iter := NewStatusNodeIterator(t)

	return iter
}

// FString implements the Noder interface.
func (t *StatusNode[S, T]) FString(trav *ffs.Traversor, opts ...ffs.Option) error {
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

// Copy implements the Noder interface.
func (n *StatusNode[S, T]) Copy() uc.Copier {
	var child_copy []Noder

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		c_copy := c.Copy().(Noder)
		child_copy = append(child_copy, c_copy)
	}

	d_copy := uc.CopyOf(n.Data).(T)

	nCopy := &TreeNode[T]{Data: d_copy}

	nCopy.LinkChildren(child_copy)

	return nCopy
}

// SetParent implements the Noder interface.
func (n *StatusNode[S, T]) SetParent(parent Noder) bool {
	if parent == nil {
		n.Parent = nil
		return true
	}

	p, ok := parent.(*StatusNode[S, T])
	if !ok {
		return false
	}

	n.Parent = p

	return true
}

// GetParent implements the Noder interface.
func (n *StatusNode[S, T]) GetParent() Noder {
	return n.Parent
}

// LinkWithParent implements the Noder interface.
//
// Children that are not *StatusNode[S, T] are ignored.
func (tn *StatusNode[S, T]) LinkChildren(children []Noder) {
	var valid_children []*StatusNode[S, T]

	for _, child := range children {
		c, ok := child.(*StatusNode[S, T])
		if !ok {
			continue
		}

		valid_children = append(valid_children, c)
	}
	if len(valid_children) == 0 {
		// Do nothing.
		return
	}

	valid_children = link_siblings_status(valid_children)

	for _, child := range valid_children {
		child.Parent = tn
	}

	tn.FirstChild = valid_children[0]
	tn.LastChild = valid_children[len(valid_children)-1]
}

// GetLeaves implements the Noder interface.
func (n *StatusNode[S, T]) GetLeaves() []Noder {
	stack := lls.NewLinkedStack[Noder](n)

	var leaves []Noder

	for {
		top, ok := stack.Pop()
		if !ok {
			break
		}

		val, ok := top.(*StatusNode[S, T])
		uc.Assert(ok, "GetLeaves: Invalid node type")

		if val.FirstChild == nil {
			leaves = append(leaves, top)

			continue
		}

		children := val.GetChildren()

		for _, child := range children {
			stack.Push(child)
		}
	}

	return leaves
}

// ToTree implements the Noder interface.
func (n *StatusNode[S, T]) TreeOf() *Tree {
	if n.FirstChild == nil {
		return &Tree{
			root:   n,
			leaves: []Noder{n},
			size:   1,
		}
	} else {
		return &Tree{
			root:   n,
			leaves: n.GetLeaves(),
			size:   n.Size(),
		}
	}
}

// Cleanup implements the Noder interface.
//
// Uses DFS traversal, does not use recursion, and does not remove the node itself.
func (tn *StatusNode[S, T]) Cleanup() {
	type Helper struct {
		Prev *StatusNode[S, T]
		Curr *StatusNode[S, T]
	}

	h := &Helper{
		Prev: nil,
		Curr: tn,
	}

	stack := lls.NewLinkedStack(h)

	for {
		h, ok := stack.Pop()
		if !ok {
			break
		}

		for c := h.Curr.FirstChild; c != nil; c = c.NextSibling {
			h := &Helper{
				Prev: c.PrevSibling,
				Curr: c,
			}

			stack.Push(h)
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

// GetAncestors implements the Noder interface.
func (n *StatusNode[S, T]) GetAncestors() []Noder {
	var ancestors []Noder

	for node := n; node.Parent != nil; node = node.Parent {
		ancestors = append(ancestors, node.Parent)
	}

	slices.Reverse(ancestors)

	return ancestors
}

// IsLeaf implements the Noder interface.
func (n *StatusNode[S, T]) IsLeaf() bool {
	return n.FirstChild == nil
}

// IsSingleton implements the Noder interface.
func (n *StatusNode[S, T]) IsSingleton() bool {
	if n.FirstChild == nil {
		return false
	}

	return n.FirstChild == n.LastChild
}

// GetFirstChild implements the Noder interface.
func (n *StatusNode[S, T]) GetFirstChild() Noder {
	return n.FirstChild
}

// DeleteChild implements the Noder interface.
func (n *StatusNode[S, T]) DeleteChild(target Noder) []Noder {
	if target == nil {
		return nil
	}

	tn, ok := target.(*StatusNode[S, T])
	if !ok {
		return nil
	}

	children := n.delete_child(tn)

	delink_with_parent_status(n, children)

	return children
}

// Size implements the Noder interface.
//
// This function is expensive since size is not stored.
func (n *StatusNode[S, T]) Size() int {
	stack := lls.NewLinkedStack(n)

	var size int

	for {
		top, ok := stack.Pop()
		if !ok {
			break
		}

		size++

		for c := top.FirstChild; c != nil; c = c.NextSibling {
			stack.Push(c)
		}
	}

	return size
}

// AddChild adds a new child to the node with the given data.
//
// Parameters:
//   - child: The child to add.
//
// Behaviors:
//   - If the child is nil, it does nothing.
func (n *StatusNode[S, T]) AddChild(child Noder) {
	if child == nil {
		return
	}

	c, ok := child.(*StatusNode[S, T])
	if !ok {
		return
	}

	// Make sure the child is not linked to any other node
	c.NextSibling = nil
	c.PrevSibling = nil

	last_child := n.LastChild

	if last_child == nil {
		n.FirstChild = c
	} else {
		last_child.NextSibling = c
		c.PrevSibling = last_child
	}

	c.Parent = n
	n.LastChild = c
}

// RemoveNode removes the node from the tree and shifts the children up
// in the space occupied by the node.
//
// Returns:
//   - []Noder: A slice of pointers to the children of the node if
//     the node is the root. Nil otherwise.
func (n *StatusNode[S, T]) RemoveNode() []Noder {
	prev_sibling := n.PrevSibling
	next_sibling := n.NextSibling

	var sub_roots []Noder

	if n.Parent == nil {
		sub_roots = n.GetChildren()
	} else {
		children := n.Parent.delete_child(n)

		for _, child := range children {
			child.SetParent(n.Parent)
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

	delink_with_parent_status(n, sub_roots)

	return sub_roots
}

// NewTreeNode creates a new node with the given data.
//
// Parameters:
//   - status: The status of the node.
//   - data: The value of the node.
//
// Returns:
//   - *StatusNode[S, T]: A pointer to the newly created node.
func NewStatusNode[S uc.Enumer, T any](status S, data T) *StatusNode[S, T] {
	tn := &StatusNode[S, T]{
		Data:   data,
		Status: status,
	}

	return tn
}

// GetLastSibling returns the last sibling of the node. If it has a parent,
// it returns the last child of the parent. Otherwise, it returns the last
// sibling of the node.
//
// Returns:
//   - *StatusNode[S, T]: A pointer to the last sibling. The node itself if it has no next sibling.
func (n *StatusNode[S, T]) GetLastSibling() *StatusNode[S, T] {
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
//   - *StatusNode[S, T]: A pointer to the first sibling. The node itself if it has no previous sibling.
func (n *StatusNode[S, T]) GetFirstSibling() *StatusNode[S, T] {
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

// IsRoot returns true if the node does not have a parent.
//
// Returns:
//   - bool: True if the node is the root, false otherwise.
func (n *StatusNode[S, T]) IsRoot() bool {
	return n.Parent == nil
}

// AddChildren adds zero or more children to the node.
//
// Parameters:
//   - children: The children to add.
//
// Behaviors:
//   - This is just a more efficient way to add multiple children.
func (n *StatusNode[S, T]) AddChildren(children ...*StatusNode[S, T]) {
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
//   - []Noder: A slice of pointers to the children of the node.
func (n *StatusNode[S, T]) GetChildren() []Noder {
	var children []Noder

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		children = append(children, c)
	}

	return children
}

// HasChild returns true if the node has the given child.
//
// Parameters:
//   - target: The child to check for.
//
// Returns:
//   - bool: True if the node has the child, false otherwise.
func (n *StatusNode[S, T]) HasChild(target *StatusNode[S, T]) bool {
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
//   - []Noder: A slice of pointers to the children of the node.
//
// Behaviors:
//   - If the node has no children, it returns nil.
func (n *StatusNode[S, T]) delete_child(target *StatusNode[S, T]) []Noder {
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

// IsChildOf returns true if the node is a child of the parent.
//
// Parameters:
//   - target: The target parent to check for.
//
// Returns:
//   - bool: True if the node is a child of the parent, false otherwise.
func (n *StatusNode[S, T]) IsChildOf(target *StatusNode[S, T]) bool {
	if target == nil {
		return false
	}

	parents := target.GetAncestors()

	for node := n; node.Parent != nil; node = node.Parent {
		parent := Noder(node.Parent)

		ok := slices.Contains(parents, parent)
		if ok {
			return true
		}
	}

	return false
}

// GetData is a getter for the data of the node.
//
// Returns:
//   - T: The data of the node.
func (n *StatusNode[S, T]) GetData() T {
	return n.Data
}

// link_siblings_status links the siblings with each other. It also sets the prev and
// next siblings of the first and last nodes to nil.
//
// Parameters:
//   - nodes: The nodes to link.
//
// Returns:
//   - []*StatusNode[S, T]: The linked nodes. (Without the nil values)
//
// Behaviors:
//   - If the nodes slice is empty, it does nothing. Nil values are filtered out.
func link_siblings_status[S uc.Enumer, T any](nodes []*StatusNode[S, T]) []*StatusNode[S, T] {
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

// delink_with_parent_status delinks the parent with the children. It also delinks the children
// with each other.
//
// Parameters:
//   - parent: The parent node.
//   - children: The children nodes.
//
// Behaviors:
//   - If the parent has no children, it does nothing.
func delink_with_parent_status[S uc.Enumer, T any](parent *StatusNode[S, T], children []Noder) {
	if len(children) == 0 {
		return
	}

	for i := 0; i < len(children); i++ {
		child := children[i]

		if child == nil {
			continue
		}

		c, ok := child.(*StatusNode[S, T])
		uc.Assert(ok, "delink_with_parent_status: Invalid child type")

		c.PrevSibling = nil
		c.NextSibling = nil
		c.Parent = nil
	}

	if parent != nil {
		parent.FirstChild = nil
		parent.LastChild = nil
	}
}
