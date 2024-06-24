package StatusTree

import (
	ud "github.com/PlayerR9/MyGoLib/Units/Debugging"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
	us "github.com/PlayerR9/MyGoLib/Units/slice"
)

// NewTree creates a new tree with the given root.
//
// Parameters:
//   - data: The value of the root.
//
// Returns:
//   - *Tree[S, T]: A pointer to the newly created tree.
func NewTreeWithHistory[S uc.Enumer, T any](status S, data T) *ud.History[*Tree[S, T]] {
	root := NewTreeNode(status, data)

	tree := &Tree[S, T]{
		root:   root,
		leaves: []*TreeNode[S, T]{root},
		size:   1,
	}

	h := ud.NewHistory(tree)

	return h
}

// SetChildrenCmd is a command that sets the children of a node.
type SetChildrenCmd[S uc.Enumer, T any] struct {
	// children is a slice of pointers to the children to set.
	children []*Tree[S, T]

	// size is the size of the tree before setting the children.
	size int

	// prevChildren is a slice of pointers to the previous children of the node.
	prevChildren []*TreeNode[S, T]

	// prevLeaves is a slice of pointers to the previous leaves of the tree.
	prevLeaves []*TreeNode[S, T]
}

// Execute implements the Debugging.Commander interface.
func (c *SetChildrenCmd[S, T]) Execute(data *Tree[S, T]) error {
	c.size = data.size
	c.prevChildren = data.root.children
	c.prevLeaves = data.leaves

	err := data.SetChildren(c.children)
	if err != nil {
		return err
	}

	return nil
}

// Undo implements the Debugging.Commander interface.
func (c *SetChildrenCmd[S, T]) Undo(data *Tree[S, T]) error {
	root := data.root
	if root == nil {
		return NewErrMissingRoot()
	}

	root.children = c.prevChildren
	data.leaves = c.prevLeaves
	data.size = c.size

	return nil
}

// Copy implements the Debugging.Commander interface.
func (c *SetChildrenCmd[S, T]) Copy() uc.Copier {
	cCopy := &SetChildrenCmd[S, T]{
		children:     make([]*Tree[S, T], len(c.children)),
		size:         c.size,
		prevChildren: make([]*TreeNode[S, T], len(c.prevChildren)),
		prevLeaves:   make([]*TreeNode[S, T], len(c.prevLeaves)),
	}

	for i, child := range c.children {
		childCopy := child.Copy().(*Tree[S, T])

		cCopy.children[i] = childCopy
	}

	for i, child := range c.prevChildren {
		childCopy := child.Copy().(*TreeNode[S, T])
		childCopy.setParent(child.parent)

		cCopy.prevChildren[i] = childCopy
	}

	for i, leaf := range c.prevLeaves {
		leafCopy := leaf.Copy().(*TreeNode[S, T])
		leafCopy.setParent(leaf.parent)

		cCopy.prevLeaves[i] = leafCopy
	}

	return cCopy
}

// NewSetChildrenCmd creates a new SetChildrenCmd.
//
// Parameters:
//   - children: The children to set.
//
// Returns:
//   - *SetChildrenCmd: A pointer to the new SetChildrenCmd.
func NewSetChildrenCmd[S uc.Enumer, T any](children []*Tree[S, T]) *SetChildrenCmd[S, T] {
	children = us.SliceFilter(children, FilterNilTree)
	if len(children) == 0 {
		return nil
	}

	return &SetChildrenCmd[S, T]{
		children: children,
	}
}

// CleanupCmd is a command that cleans up the tree.
type CleanupCmd[S uc.Enumer, T any] struct {
	// root is a pointer to the root of the tree.
	root *TreeNode[S, T]
}

// Execute implements the Debugging.Commander interface.
//
// Never errors.
func (c *CleanupCmd[S, T]) Execute(data *Tree[S, T]) error {
	c.root = data.root.Copy().(*TreeNode[S, T])

	data.Cleanup()

	return nil
}

// Undo implements the Debugging.Commander interface.
//
// Never errors.
func (c *CleanupCmd[S, T]) Undo(data *Tree[S, T]) error {
	data.root = c.root

	return nil
}

// Copy implements the Debugging.Commander interface.
func (c *CleanupCmd[S, T]) Copy() uc.Copier {
	return &CleanupCmd[S, T]{
		root: c.root.Copy().(*TreeNode[S, T]),
	}
}

// NewCleanupCmd creates a new CleanupCmd.
//
// Returns:
//   - *CleanupCmd: A pointer to the new CleanupCmd.
func NewCleanupCmd[S uc.Enumer, T any]() *CleanupCmd[S, T] {
	return &CleanupCmd[S, T]{}
}

// RegenerateLeavesCmd is a command that regenerates the leaves of the tree.
type RegenerateLeavesCmd[S uc.Enumer, T any] struct {
	// leaves is a slice of pointers to the leaves of the tree.
	leaves []*TreeNode[S, T]

	// size is the size of the tree before regenerating the leaves.
	size int
}

// Execute implements the Debugging.Commander interface.
func (c *RegenerateLeavesCmd[S, T]) Execute(data *Tree[S, T]) error {
	c.leaves = data.leaves
	c.size = data.size

	data.RegenerateLeaves()

	return nil
}

// Undo implements the Debugging.Commander interface.
func (c *RegenerateLeavesCmd[S, T]) Undo(data *Tree[S, T]) error {
	data.leaves = c.leaves
	data.size = c.size

	return nil
}

// Copy implements the Debugging.Commander interface.
func (c *RegenerateLeavesCmd[S, T]) Copy() uc.Copier {
	leaves := make([]*TreeNode[S, T], len(c.leaves))
	copy(leaves, c.leaves)

	cCopy := &RegenerateLeavesCmd[S, T]{
		leaves: leaves,
		size:   c.size,
	}

	return cCopy
}

// NewRegenerateLeavesCmd creates a new RegenerateLeavesCmd.
//
// Returns:
//   - *RegenerateLeavesCmd: A pointer to the new RegenerateLeavesCmd.
func NewRegenerateLeavesCmd[S uc.Enumer, T any]() *RegenerateLeavesCmd[S, T] {
	return &RegenerateLeavesCmd[S, T]{}
}

// UpdateLeavesCmd is a command that updates the leaves of the tree.
type UpdateLeavesCmd[S uc.Enumer, T any] struct {
	// leaves is a slice of pointers to the leaves of the tree.
	leaves []*TreeNode[S, T]

	// size is the size of the tree before updating the leaves.
	size int
}

// Execute implements the Debugging.Commander interface.
func (c *UpdateLeavesCmd[S, T]) Execute(data *Tree[S, T]) error {
	c.leaves = data.leaves
	c.size = data.size

	data.UpdateLeaves()

	return nil
}

// Undo implements the Debugging.Commander interface.
func (c *UpdateLeavesCmd[S, T]) Undo(data *Tree[S, T]) error {
	data.leaves = c.leaves
	data.size = c.size

	return nil
}

// Copy implements the Debugging.Commander interface.
func (c *UpdateLeavesCmd[S, T]) Copy() uc.Copier {
	leaves := make([]*TreeNode[S, T], len(c.leaves))
	copy(leaves, c.leaves)

	cCopy := &UpdateLeavesCmd[S, T]{
		leaves: leaves,
		size:   c.size,
	}

	return cCopy
}

// NewUpdateLeavesCmd creates a new UpdateLeavesCmd.
//
// Returns:
//   - *UpdateLeavesCmd: A pointer to the new UpdateLeavesCmd.
func NewUpdateLeavesCmd[S uc.Enumer, T any]() *UpdateLeavesCmd[S, T] {
	return &UpdateLeavesCmd[S, T]{}
}

// PruneBranchesCmd is a command that prunes the branches of the tree.
type PruneBranchesCmd[S uc.Enumer, T any] struct {
	// tree is a pointer to the tree before pruning the branches.
	tree *Tree[S, T]

	// filter is the filter to apply to prune the branches.
	filter us.PredicateFilter[uc.Pair[S, T]]

	// ok is true if the whole tree can be deleted, false otherwise.
	ok bool
}

// Execute implements the Debugging.Commander interface.
//
// Never errors.
func (c *PruneBranchesCmd[S, T]) Execute(data *Tree[S, T]) error {
	c.tree = data.Copy().(*Tree[S, T])

	c.ok = data.PruneBranches(c.filter)

	return nil
}

// Undo implements the Debugging.Commander interface.
//
// Never errors.
func (c *PruneBranchesCmd[S, T]) Undo(data *Tree[S, T]) error {
	data.root = c.tree.root
	data.leaves = c.tree.leaves
	data.size = c.tree.size

	return nil
}

// Copy implements the Debugging.Commander interface.
func (c *PruneBranchesCmd[S, T]) Copy() uc.Copier {
	tree := c.tree.Copy().(*Tree[S, T])

	cmdCopy := &PruneBranchesCmd[S, T]{
		tree:   tree,
		filter: c.filter,
		ok:     c.ok,
	}

	return cmdCopy
}

// NewPruneBranchesCmd creates a new PruneBranchesCmd.
//
// Parameters:
//   - filter: The filter to apply to prune the branches.
//
// Returns:
//   - *PruneBranchesCmd: A pointer to the new PruneBranchesCmd.
func NewPruneBranchesCmd[S uc.Enumer, T any](filter us.PredicateFilter[uc.Pair[S, T]]) *PruneBranchesCmd[S, T] {
	if filter == nil {
		return nil
	}

	cmd := &PruneBranchesCmd[S, T]{
		filter: filter,
	}

	return cmd
}

// GetOk returns the value of the ok field.
//
// Call this function after executing the command.
//
// Returns:
//   - bool: The value of the ok field.
func (c *PruneBranchesCmd[S, T]) GetOk() bool {
	return c.ok
}

// SkipFuncCmd is a command that skips the nodes of the tree that
// satisfy the given filter.
type SkipFuncCmd[S uc.Enumer, T any] struct {
	// tree is a pointer to the tree before skipping the nodes.
	tree *Tree[S, T]

	// filter is the filter to apply to skip the nodes.
	filter us.PredicateFilter[uc.Pair[S, T]]

	// trees is a slice of pointers to the trees obtained after
	// skipping the nodes.
	trees []*Tree[S, T]
}

// Execute implements the Debugging.Commander interface.
//
// Never errors.
func (c *SkipFuncCmd[S, T]) Execute(data *Tree[S, T]) error {
	c.tree = data.Copy().(*Tree[S, T])

	c.trees = data.SkipFilter(c.filter)

	return nil
}

// Undo implements the Debugging.Commander interface.
//
// Never errors.
func (c *SkipFuncCmd[S, T]) Undo(data *Tree[S, T]) error {
	data.root = c.tree.root
	data.leaves = c.tree.leaves
	data.size = c.tree.size

	return nil
}

// Copy implements the Debugging.Commander interface.
func (c *SkipFuncCmd[S, T]) Copy() uc.Copier {
	tree := c.tree.Copy().(*Tree[S, T])

	trees := make([]*Tree[S, T], len(c.trees))
	for i, t := range c.trees {
		tCopy := t.Copy().(*Tree[S, T])

		trees[i] = tCopy
	}

	cmdCopy := &SkipFuncCmd[S, T]{
		tree:   tree,
		filter: c.filter,
		trees:  trees,
	}

	return cmdCopy
}

// NewSkipFuncCmd creates a new SkipFuncCmd.
//
// Parameters:
//   - filter: The filter to apply to skip the nodes.
//
// Returns:
//   - *SkipFuncCmd: A pointer to the new SkipFuncCmd.
func NewSkipFuncCmd[S uc.Enumer, T any](filter us.PredicateFilter[uc.Pair[S, T]]) *SkipFuncCmd[S, T] {
	if filter == nil {
		return nil
	}

	cmd := &SkipFuncCmd[S, T]{
		filter: filter,
	}

	return cmd
}

// GetTrees returns the value of the trees field.
//
// Call this function after executing the command.
//
// Returns:
//   - []*Tree[S, T]: A slice of pointers to the trees obtained after
//     skipping the nodes.
func (c *SkipFuncCmd[S, T]) GetTrees() []*Tree[S, T] {
	return c.trees
}

// ProcessLeavesCmd is a command that processes the leaves of the tree.
type ProcessLeavesCmd[S uc.Enumer, T any] struct {
	// leaves is a slice of pointers to the leaves of the tree.
	leaves []*TreeNode[S, T]

	// f is the function to apply to the leaves.
	f uc.EvalManyFunc[*TreeNode[S, T], uc.Pair[S, T]]
}

// Execute implements the Debugging.Commander interface.
func (c *ProcessLeavesCmd[S, T]) Execute(data *Tree[S, T]) error {
	c.leaves = data.leaves

	err := data.ProcessLeaves(c.f)
	if err != nil {
		return err
	}

	return nil
}

// Undo implements the Debugging.Commander interface.
func (c *ProcessLeavesCmd[S, T]) Undo(data *Tree[S, T]) error {
	data.leaves = c.leaves

	return nil
}

// Copy implements the Debugging.Commander interface.
func (c *ProcessLeavesCmd[S, T]) Copy() uc.Copier {
	leaves := make([]*TreeNode[S, T], len(c.leaves))
	copy(leaves, c.leaves)

	cCopy := &ProcessLeavesCmd[S, T]{
		leaves: leaves,
		f:      c.f,
	}

	return cCopy
}

// NewProcessLeavesCmd creates a new ProcessLeavesCmd.
//
// Parameters:
//   - f: The function to apply to the leaves.
//
// Returns:
//   - *ProcessLeavesCmd: A pointer to the new ProcessLeavesCmd.
func NewProcessLeavesCmd[S uc.Enumer, T any](f uc.EvalManyFunc[*TreeNode[S, T], uc.Pair[S, T]]) *ProcessLeavesCmd[S, T] {
	if f == nil {
		return nil
	}

	cmd := &ProcessLeavesCmd[S, T]{
		f: f,
	}

	return cmd
}

// DeleteBranchContainingCmd is a command that deletes the branch containing
// the given node.
type DeleteBranchContainingCmd[S uc.Enumer, T any] struct {
	// tree is a pointer to the tree before deleting the branch.
	tree *Tree[S, T]

	// tn is a pointer to the node to delete.
	tn *TreeNode[S, T]
}

// Execute implements the Debugging.Commander interface.
func (c *DeleteBranchContainingCmd[S, T]) Execute(data *Tree[S, T]) error {
	c.tree = data.Copy().(*Tree[S, T])

	err := data.DeleteBranchContaining(c.tn)
	if err != nil {
		return err
	}

	return nil
}

// Undo implements the Debugging.Commander interface.
func (c *DeleteBranchContainingCmd[S, T]) Undo(data *Tree[S, T]) error {
	data.root = c.tree.root
	data.leaves = c.tree.leaves
	data.size = c.tree.size

	return nil
}

// Copy implements the Debugging.Commander interface.
func (c *DeleteBranchContainingCmd[S, T]) Copy() uc.Copier {
	tree := c.tree.Copy().(*Tree[S, T])
	tn := c.tn.Copy().(*TreeNode[S, T])

	cmdCopy := &DeleteBranchContainingCmd[S, T]{
		tree: tree,
		tn:   tn,
	}

	return cmdCopy
}

// NewDeleteBranchContainingCmd creates a new DeleteBranchContainingCmd.
//
// Parameters:
//   - tn: The node to delete.
//
// Returns:
//   - *DeleteBranchContainingCmd: A pointer to the new DeleteBranchContainingCmd.
func NewDeleteBranchContainingCmd[S uc.Enumer, T any](tn *TreeNode[S, T]) *DeleteBranchContainingCmd[S, T] {
	if tn == nil {
		return nil
	}

	cmd := &DeleteBranchContainingCmd[S, T]{
		tn: tn,
	}

	return cmd
}

// PruneTreeCmd is a command that prunes the tree using the given filter.
type PruneTreeCmd[S uc.Enumer, T any] struct {
	// tree is a pointer to the tree before pruning.
	tree *Tree[S, T]

	// filter is the filter to use to prune the tree.
	filter us.PredicateFilter[uc.Pair[S, T]]

	// ok is true if no nodes were pruned, false otherwise.
	ok bool
}

// Execute implements the Debugging.Commander interface.
func (c *PruneTreeCmd[S, T]) Execute(data *Tree[S, T]) error {
	c.tree = data.Copy().(*Tree[S, T])

	c.ok = data.Prune(c.filter)

	return nil
}

// Undo implements the Debugging.Commander interface.
func (c *PruneTreeCmd[S, T]) Undo(data *Tree[S, T]) error {
	data.root = c.tree.root
	data.leaves = c.tree.leaves
	data.size = c.tree.size

	return nil
}

// Copy implements the Debugging.Commander interface.
func (c *PruneTreeCmd[S, T]) Copy() uc.Copier {
	tree := c.tree.Copy().(*Tree[S, T])

	cmdCopy := &PruneTreeCmd[S, T]{
		tree:   tree,
		filter: c.filter,
		ok:     c.ok,
	}

	return cmdCopy
}

// NewPruneTreeCmd creates a new PruneTreeCmd.
//
// Parameters:
//   - filter: The filter to use to prune the tree.
//
// Returns:
//   - *PruneTreeCmd: A pointer to the new PruneTreeCmd.
func NewPruneTreeCmd[S uc.Enumer, T any](filter us.PredicateFilter[uc.Pair[S, T]]) *PruneTreeCmd[S, T] {
	if filter == nil {
		return nil
	}

	cmd := &PruneTreeCmd[S, T]{
		filter: filter,
	}

	return cmd
}

// GetOk returns the value of the ok field.
//
// Call this function after executing the command.
//
// Returns:
//   - bool: The value of the ok field.
func (c *PruneTreeCmd[S, T]) GetOk() bool {
	return c.ok
}
