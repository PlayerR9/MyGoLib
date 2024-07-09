package Tree

import (
	"slices"

	ffs "github.com/PlayerR9/MyGoLib/Formatting/FString"

	lls "github.com/PlayerR9/MyGoLib/ListLike/Stacker"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
	us "github.com/PlayerR9/MyGoLib/Units/slice"
	uto "github.com/PlayerR9/MyGoLib/Utility/object"
)

// TreeNodeIterator is a iterator that iterates over the children of a tree node.
type TreeNodeIterator[T any] struct {
	// node is the node whose children are being iterated over.
	node *TreeNode[T]

	// current_node is the current node being iterated over.
	current_node *TreeNode[T]
}

// Size implements the common.Iterater interface.
func (t *TreeNodeIterator[T]) Consume() (Noder, error) {
	if t.current_node == nil {
		return nil, uc.NewErrExhaustedIter()
	}

	n := t.current_node

	t.current_node = t.current_node.NextSibling

	return n, nil
}

// Restart implements the common.Iterater interface.
func (t *TreeNodeIterator[T]) Restart() {
	t.current_node = t.node.FirstChild
}

// NewTreeNodeIterator creates a new iterator for the given node.
//
// Parameters:
//   - node: The node whose children are being iterated over.
//
// Returns:
//   - *TreeNodeIterator[T]: A pointer to the iterator. Nil if the node is nil.
func NewTreeNodeIterator[T any](node *TreeNode[T]) *TreeNodeIterator[T] {
	if node == nil {
		return nil
	}

	ti := &TreeNodeIterator[T]{
		node:         node,
		current_node: node.FirstChild,
	}

	return ti
}

// Noder is an interface that represents a node in a tree.
type Noder interface {
	// SetParent sets the parent of the node.
	//
	// Parameters:
	//   - parent: The parent node.
	//
	// Returns:
	//   - bool: True if the parent is set, false otherwise.
	SetParent(parent Noder) bool

	// GetParent returns the parent of the node.
	//
	// Returns:
	//   - Noder: The parent node.
	GetParent() Noder

	// LinkWithParent links the parent with the children. It also links the children
	// with each other.
	//
	// Parameters:
	//   - parent: The parent node.
	//   - children: The children nodes.
	//
	// Behaviors:
	//   - Only valid children are linked while the rest are ignored.
	LinkChildren(children []Noder)

	// IsLeaf returns true if the node is a leaf.
	//
	// Returns:
	//   - bool: True if the node is a leaf, false otherwise.
	IsLeaf() bool

	// IsSingleton returns true if the node is a singleton (i.e., has only one child).
	//
	// Returns:
	//   - bool: True if the node is a singleton, false otherwise.
	IsSingleton() bool

	// GetLeaves returns all the leaves of the tree rooted at the node.
	//
	// Should be a DFS traversal.
	//
	// Returns:
	//   - []Noder: A slice of pointers to the leaves of the tree.
	//
	// Behaviors:
	//   - The leaves are returned in the order of a DFS traversal.
	GetLeaves() []Noder

	// GetAncestors returns all the ancestors of the node.
	//
	// This excludes the node itself.
	//
	// Returns:
	//   - []Noder: A slice of pointers to the ancestors of the node.
	//
	// Behaviors:
	//   - The ancestors are returned in the opposite order of a DFS traversal.
	//     Therefore, the first element is the parent of the node.
	GetAncestors() []Noder

	// GetFirstChild returns the first child of the node.
	//
	// Returns:
	//   - Noder: The first child of the node. Nil if the node has no children.
	GetFirstChild() Noder

	// DeleteChild removes the given child from the children of the node.
	//
	// Parameters:
	//   - target: The child to remove.
	//
	// Returns:
	//   - []Noder: A slice of pointers to the children of the node. Nil if the node has no children.
	DeleteChild(target Noder) []Noder

	// Size returns the number of nodes in the tree rooted at n.
	//
	// Returns:
	//   - size: The number of nodes in the tree.
	Size() int

	// AddChild adds a new child to the node with the given data.
	//
	// Parameters:
	//   - child: The child to add.
	//
	// Behaviors:
	//   - If the child is not valid, it is ignored.
	AddChild(child Noder)

	// removeNode removes the node from the tree and shifts the children up
	// in the space occupied by the node.
	//
	// Returns:
	//   - []Noder: A slice of pointers to the children of the node if
	//     the node is the root. Nil otherwise.
	RemoveNode() []Noder

	Treeer
	uc.Iterable[Noder]
	uc.Copier
	ffs.FStringer
	uto.Cleaner
}

// TreeNode is a generic data structure that represents a node in a tree.
type TreeNode[T any] struct {
	// Parent, FirstChild, NextSibling, LastChild, PrevSibling are pointers to
	// the parent, first child, next sibling, last child, and previous sibling
	// of the node respectively.
	Parent, FirstChild, NextSibling, LastChild, PrevSibling *TreeNode[T]

	// Data is the value of the node.
	Data T
}

// Iterator implements Noder.
func (t *TreeNode[T]) Iterator() uc.Iterater[Noder] {
	iter := NewTreeNodeIterator(t)

	return iter
}

// FString implements the Noder interface.
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

// Copy implements the Noder interface.
func (n *TreeNode[T]) Copy() uc.Copier {
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
func (n *TreeNode[T]) SetParent(parent Noder) bool {
	if parent == nil {
		n.Parent = nil
		return true
	}

	p, ok := parent.(*TreeNode[T])
	if !ok {
		return false
	}

	n.Parent = p

	return true
}

// GetParent implements the Noder interface.
func (n *TreeNode[T]) GetParent() Noder {
	return n.Parent
}

// LinkWithParent implements the Noder interface.
//
// Children that are not *TreeNode[T] are ignored.
func (tn *TreeNode[T]) LinkChildren(children []Noder) {
	var valid_children []*TreeNode[T]

	for _, child := range children {
		c, ok := child.(*TreeNode[T])
		if !ok {
			continue
		}

		valid_children = append(valid_children, c)
	}
	if len(valid_children) == 0 {
		// Do nothing.
		return
	}

	valid_children = link_siblings(valid_children)

	for _, child := range valid_children {
		child.Parent = tn
	}

	tn.FirstChild = valid_children[0]
	tn.LastChild = valid_children[len(valid_children)-1]
}

// GetLeaves implements the Noder interface.
func (n *TreeNode[T]) GetLeaves() []Noder {
	S := lls.NewLinkedStack[Noder](n)

	var leaves []Noder

	for {
		top, ok := S.Pop()
		if !ok {
			break
		}

		val, ok := top.(*TreeNode[T])
		uc.Assert(ok, "GetLeaves: Invalid node type")

		if val.FirstChild == nil {
			leaves = append(leaves, top)

			continue
		}

		children := val.GetChildren()

		for _, child := range children {
			S.Push(child)
		}
	}

	return leaves
}

// ToTree implements the Noder interface.
func (n *TreeNode[T]) TreeOf() *Tree {
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

// GetAncestors implements the Noder interface.
func (n *TreeNode[T]) GetAncestors() []Noder {
	var ancestors []Noder

	for node := n; node.Parent != nil; node = node.Parent {
		ancestors = append(ancestors, node.Parent)
	}

	slices.Reverse(ancestors)

	return ancestors
}

// IsLeaf implements the Noder interface.
func (n *TreeNode[T]) IsLeaf() bool {
	return n.FirstChild == nil
}

// IsSingleton implements the Noder interface.
func (n *TreeNode[T]) IsSingleton() bool {
	if n.FirstChild == nil {
		return false
	}

	return n.FirstChild == n.LastChild
}

// GetFirstChild implements the Noder interface.
func (n *TreeNode[T]) GetFirstChild() Noder {
	return n.FirstChild
}

// DeleteChild implements the Noder interface.
func (n *TreeNode[T]) DeleteChild(target Noder) []Noder {
	if target == nil {
		return nil
	}

	tn, ok := target.(*TreeNode[T])
	if !ok {
		return nil
	}

	children := n.delete_child(tn)

	delink_with_parent(n, children)

	return children
}

// Size implements the Noder interface.
//
// This function is expensive since size is not stored.
func (n *TreeNode[T]) Size() int {
	S := lls.NewLinkedStack(n)

	var size int

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

	return size
}

// AddChild adds a new child to the node with the given data.
//
// Parameters:
//   - child: The child to add.
//
// Behaviors:
//   - If the child is nil, it does nothing.
func (n *TreeNode[T]) AddChild(child Noder) {
	if child == nil {
		return
	}

	c, ok := child.(*TreeNode[T])
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
func (n *TreeNode[T]) RemoveNode() []Noder {
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

	delink_with_parent(n, sub_roots)

	return sub_roots
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

// IsRoot returns true if the node does not have a parent.
//
// Returns:
//   - bool: True if the node is the root, false otherwise.
func (n *TreeNode[T]) IsRoot() bool {
	return n.Parent == nil
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
//   - []Noder: A slice of pointers to the children of the node.
func (n *TreeNode[T]) GetChildren() []Noder {
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
//   - []Noder: A slice of pointers to the children of the node.
//
// Behaviors:
//   - If the node has no children, it returns nil.
func (n *TreeNode[T]) delete_child(target *TreeNode[T]) []Noder {
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
func (n *TreeNode[T]) IsChildOf(target *TreeNode[T]) bool {
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
func (n *TreeNode[T]) GetData() T {
	return n.Data
}

// link_siblings links the siblings with each other. It also sets the prev and
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
func link_siblings[T any](nodes []*TreeNode[T]) []*TreeNode[T] {
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

// delink_with_parent delinks the parent with the children. It also delinks the children
// with each other.
//
// Parameters:
//   - parent: The parent node.
//   - children: The children nodes.
//
// Behaviors:
//   - If the parent has no children, it does nothing.
func delink_with_parent[T any](parent *TreeNode[T], children []Noder) {
	if len(children) == 0 {
		return
	}

	for i := 0; i < len(children); i++ {
		child := children[i]

		if child == nil {
			continue
		}

		c, ok := child.(*TreeNode[T])
		uc.Assert(ok, "delink_with_parent: Invalid child type")

		c.PrevSibling = nil
		c.NextSibling = nil
		c.Parent = nil
	}

	if parent != nil {
		parent.FirstChild = nil
		parent.LastChild = nil
	}
}
