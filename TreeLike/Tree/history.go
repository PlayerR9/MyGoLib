package Tree

import (
	"errors"

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
//   - *Tree[T]: A pointer to the newly created tree.
func NewTreeWithHistory[T any](data T) *ud.History[*Tree[T]] {
	root := NewTreeNode(data)

	tree := &Tree[T]{
		root:   root,
		leaves: []*TreeNode[T]{root},
		size:   1,
	}

	h := ud.NewHistory(tree)

	return h
}

// SetChildrenCmd is a command that sets the children of a node.
type SetChildrenCmd[T any] struct {
	// children is a slice of pointers to the children to set.
	children []*Tree[T]

	// size is the size of the tree before setting the children.
	size int

	// prevChildren is a slice of pointers to the previous children of the node.
	prevChildren []*TreeNode[T]

	// prevLeaves is a slice of pointers to the previous leaves of the tree.
	prevLeaves []*TreeNode[T]
}

// Execute implements the Debugging.Commander interface.
func (c *SetChildrenCmd[T]) Execute(data *Tree[T]) error {
	c.size = data.size
	c.prevChildren = data.root.GetChildren()
	c.prevLeaves = data.leaves

	err := data.SetChildren(c.children)
	if err != nil {
		return err
	}

	return nil
}

// Undo implements the Debugging.Commander interface.
func (c *SetChildrenCmd[T]) Undo(data *Tree[T]) error {
	root := data.root
	if root == nil {
		return NewErrMissingRoot()
	}

	LinkWithParent(root, c.prevChildren)
	data.leaves = c.prevLeaves
	data.size = c.size

	return nil
}

// Copy implements the Debugging.Commander interface.
func (c *SetChildrenCmd[T]) Copy() uc.Copier {
	cCopy := &SetChildrenCmd[T]{
		children:     make([]*Tree[T], len(c.children)),
		size:         c.size,
		prevChildren: make([]*TreeNode[T], len(c.prevChildren)),
		prevLeaves:   make([]*TreeNode[T], len(c.prevLeaves)),
	}

	for i, child := range c.children {
		childCopy := child.Copy().(*Tree[T])

		cCopy.children[i] = childCopy
	}

	for i, child := range c.prevChildren {
		childCopy := child.Copy().(*TreeNode[T])
		childCopy.Parent = child.Parent

		cCopy.prevChildren[i] = childCopy
	}

	for i, leaf := range c.prevLeaves {
		leafCopy := leaf.Copy().(*TreeNode[T])
		leafCopy.Parent = leaf.Parent

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
func NewSetChildrenCmd[T any](children []*Tree[T]) *SetChildrenCmd[T] {
	children = us.SliceFilter(children, FilterNilTree)
	if len(children) == 0 {
		return nil
	}

	return &SetChildrenCmd[T]{
		children: children,
	}
}

// CleanupCmd is a command that cleans up the tree.
type CleanupCmd[T any] struct {
	// root is a pointer to the root of the tree.
	root *TreeNode[T]
}

// Execute implements the Debugging.Commander interface.
//
// Never errors.
func (c *CleanupCmd[T]) Execute(data *Tree[T]) error {
	c.root = data.root.Copy().(*TreeNode[T])

	data.Cleanup()

	return nil
}

// Undo implements the Debugging.Commander interface.
//
// Never errors.
func (c *CleanupCmd[T]) Undo(data *Tree[T]) error {
	data.root = c.root

	return nil
}

// Copy implements the Debugging.Commander interface.
func (c *CleanupCmd[T]) Copy() uc.Copier {
	return &CleanupCmd[T]{
		root: c.root.Copy().(*TreeNode[T]),
	}
}

// NewCleanupCmd creates a new CleanupCmd.
//
// Returns:
//   - *CleanupCmd: A pointer to the new CleanupCmd.
func NewCleanupCmd[T any]() *CleanupCmd[T] {
	return &CleanupCmd[T]{}
}

// RegenerateLeavesCmd is a command that regenerates the leaves of the tree.
type RegenerateLeavesCmd[T any] struct {
	// leaves is a slice of pointers to the leaves of the tree.
	leaves []*TreeNode[T]

	// size is the size of the tree before regenerating the leaves.
	size int
}

// Execute implements the Debugging.Commander interface.
func (c *RegenerateLeavesCmd[T]) Execute(data *Tree[T]) error {
	c.leaves = data.leaves
	c.size = data.size

	data.RegenerateLeaves()

	return nil
}

// Undo implements the Debugging.Commander interface.
func (c *RegenerateLeavesCmd[T]) Undo(data *Tree[T]) error {
	data.leaves = c.leaves
	data.size = c.size

	return nil
}

// Copy implements the Debugging.Commander interface.
func (c *RegenerateLeavesCmd[T]) Copy() uc.Copier {
	leaves := make([]*TreeNode[T], len(c.leaves))
	copy(leaves, c.leaves)

	cCopy := &RegenerateLeavesCmd[T]{
		leaves: leaves,
		size:   c.size,
	}

	return cCopy
}

// NewRegenerateLeavesCmd creates a new RegenerateLeavesCmd.
//
// Returns:
//   - *RegenerateLeavesCmd: A pointer to the new RegenerateLeavesCmd.
func NewRegenerateLeavesCmd[T any]() *RegenerateLeavesCmd[T] {
	return &RegenerateLeavesCmd[T]{}
}

// UpdateLeavesCmd is a command that updates the leaves of the tree.
type UpdateLeavesCmd[T any] struct {
	// leaves is a slice of pointers to the leaves of the tree.
	leaves []*TreeNode[T]

	// size is the size of the tree before updating the leaves.
	size int
}

// Execute implements the Debugging.Commander interface.
func (c *UpdateLeavesCmd[T]) Execute(data *Tree[T]) error {
	c.leaves = data.leaves
	c.size = data.size

	data.UpdateLeaves()

	return nil
}

// Undo implements the Debugging.Commander interface.
func (c *UpdateLeavesCmd[T]) Undo(data *Tree[T]) error {
	data.leaves = c.leaves
	data.size = c.size

	return nil
}

// Copy implements the Debugging.Commander interface.
func (c *UpdateLeavesCmd[T]) Copy() uc.Copier {
	leaves := make([]*TreeNode[T], len(c.leaves))
	copy(leaves, c.leaves)

	cCopy := &UpdateLeavesCmd[T]{
		leaves: leaves,
		size:   c.size,
	}

	return cCopy
}

// NewUpdateLeavesCmd creates a new UpdateLeavesCmd.
//
// Returns:
//   - *UpdateLeavesCmd: A pointer to the new UpdateLeavesCmd.
func NewUpdateLeavesCmd[T any]() *UpdateLeavesCmd[T] {
	return &UpdateLeavesCmd[T]{}
}

// PruneBranchesCmd is a command that prunes the branches of the tree.
type PruneBranchesCmd[T any] struct {
	// tree is a pointer to the tree before pruning the branches.
	tree *Tree[T]

	// filter is the filter to apply to prune the branches.
	filter us.PredicateFilter[T]

	// ok is true if the whole tree can be deleted, false otherwise.
	ok bool
}

// Execute implements the Debugging.Commander interface.
//
// Never errors.
func (c *PruneBranchesCmd[T]) Execute(data *Tree[T]) error {
	c.tree = data.Copy().(*Tree[T])

	c.ok = data.PruneBranches(c.filter)

	return nil
}

// Undo implements the Debugging.Commander interface.
//
// Never errors.
func (c *PruneBranchesCmd[T]) Undo(data *Tree[T]) error {
	data.root = c.tree.root
	data.leaves = c.tree.leaves
	data.size = c.tree.size

	return nil
}

// Copy implements the Debugging.Commander interface.
func (c *PruneBranchesCmd[T]) Copy() uc.Copier {
	tree := c.tree.Copy().(*Tree[T])

	cmdCopy := &PruneBranchesCmd[T]{
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
func NewPruneBranchesCmd[T any](filter us.PredicateFilter[T]) *PruneBranchesCmd[T] {
	if filter == nil {
		return nil
	}

	cmd := &PruneBranchesCmd[T]{
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
func (c *PruneBranchesCmd[T]) GetOk() bool {
	return c.ok
}

// SkipFuncCmd is a command that skips the nodes of the tree that
// satisfy the given filter.
type SkipFuncCmd[T any] struct {
	// tree is a pointer to the tree before skipping the nodes.
	tree *Tree[T]

	// filter is the filter to apply to skip the nodes.
	filter us.PredicateFilter[T]

	// trees is a slice of pointers to the trees obtained after
	// skipping the nodes.
	trees []*Tree[T]
}

// Execute implements the Debugging.Commander interface.
//
// Never errors.
func (c *SkipFuncCmd[T]) Execute(data *Tree[T]) error {
	c.tree = data.Copy().(*Tree[T])

	c.trees = data.SkipFilter(c.filter)

	return nil
}

// Undo implements the Debugging.Commander interface.
//
// Never errors.
func (c *SkipFuncCmd[T]) Undo(data *Tree[T]) error {
	data.root = c.tree.root
	data.leaves = c.tree.leaves
	data.size = c.tree.size

	return nil
}

// Copy implements the Debugging.Commander interface.
func (c *SkipFuncCmd[T]) Copy() uc.Copier {
	tree := c.tree.Copy().(*Tree[T])

	trees := make([]*Tree[T], len(c.trees))
	for i, t := range c.trees {
		tCopy := t.Copy().(*Tree[T])

		trees[i] = tCopy
	}

	cmdCopy := &SkipFuncCmd[T]{
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
func NewSkipFuncCmd[T any](filter us.PredicateFilter[T]) *SkipFuncCmd[T] {
	if filter == nil {
		return nil
	}

	cmd := &SkipFuncCmd[T]{
		filter: filter,
	}

	return cmd
}

// GetTrees returns the value of the trees field.
//
// Call this function after executing the command.
//
// Returns:
//   - []*Tree[T]: A slice of pointers to the trees obtained after
//     skipping the nodes.
func (c *SkipFuncCmd[T]) GetTrees() []*Tree[T] {
	return c.trees
}

// ProcessLeavesCmd is a command that processes the leaves of the tree.
type ProcessLeavesCmd[T any] struct {
	// leaves is a slice of pointers to the leaves of the tree.
	leaves []*TreeNode[T]

	// f is the function to apply to the leaves.
	f uc.EvalManyFunc[T, T]
}

// Execute implements the Debugging.Commander interface.
func (c *ProcessLeavesCmd[T]) Execute(data *Tree[T]) error {
	c.leaves = data.leaves

	err := data.ProcessLeaves(c.f)
	if err != nil {
		return err
	}

	return nil
}

// Undo implements the Debugging.Commander interface.
func (c *ProcessLeavesCmd[T]) Undo(data *Tree[T]) error {
	data.leaves = c.leaves

	return nil
}

// Copy implements the Debugging.Commander interface.
func (c *ProcessLeavesCmd[T]) Copy() uc.Copier {
	leaves := make([]*TreeNode[T], len(c.leaves))
	copy(leaves, c.leaves)

	cCopy := &ProcessLeavesCmd[T]{
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
func NewProcessLeavesCmd[T any](f uc.EvalManyFunc[T, T]) *ProcessLeavesCmd[T] {
	if f == nil {
		return nil
	}

	cmd := &ProcessLeavesCmd[T]{
		f: f,
	}

	return cmd
}

// DeleteBranchContainingCmd is a command that deletes the branch containing
// the given node.
type DeleteBranchContainingCmd[T any] struct {
	// tree is a pointer to the tree before deleting the branch.
	tree *Tree[T]

	// tn is a pointer to the node to delete.
	tn *TreeNode[T]
}

// Execute implements the Debugging.Commander interface.
func (c *DeleteBranchContainingCmd[T]) Execute(data *Tree[T]) error {
	c.tree = data.Copy().(*Tree[T])

	err := data.DeleteBranchContaining(c.tn)
	if err != nil {
		return err
	}

	return nil
}

// Undo implements the Debugging.Commander interface.
func (c *DeleteBranchContainingCmd[T]) Undo(data *Tree[T]) error {
	data.root = c.tree.root
	data.leaves = c.tree.leaves
	data.size = c.tree.size

	return nil
}

// Copy implements the Debugging.Commander interface.
func (c *DeleteBranchContainingCmd[T]) Copy() uc.Copier {
	tree := c.tree.Copy().(*Tree[T])
	tn := c.tn.Copy().(*TreeNode[T])

	cmdCopy := &DeleteBranchContainingCmd[T]{
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
func NewDeleteBranchContainingCmd[T any](tn *TreeNode[T]) *DeleteBranchContainingCmd[T] {
	if tn == nil {
		return nil
	}

	cmd := &DeleteBranchContainingCmd[T]{
		tn: tn,
	}

	return cmd
}

// PruneTreeCmd is a command that prunes the tree using the given filter.
type PruneTreeCmd[T any] struct {
	// tree is a pointer to the tree before pruning.
	tree *Tree[T]

	// filter is the filter to use to prune the tree.
	filter us.PredicateFilter[T]

	// ok is true if no nodes were pruned, false otherwise.
	ok bool
}

// Execute implements the Debugging.Commander interface.
func (c *PruneTreeCmd[T]) Execute(data *Tree[T]) error {
	c.tree = data.Copy().(*Tree[T])

	c.ok = data.Prune(c.filter)

	return nil
}

// Undo implements the Debugging.Commander interface.
func (c *PruneTreeCmd[T]) Undo(data *Tree[T]) error {
	data.root = c.tree.root
	data.leaves = c.tree.leaves
	data.size = c.tree.size

	return nil
}

// Copy implements the Debugging.Commander interface.
func (c *PruneTreeCmd[T]) Copy() uc.Copier {
	tree := c.tree.Copy().(*Tree[T])

	cmdCopy := &PruneTreeCmd[T]{
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
func NewPruneTreeCmd[T any](filter us.PredicateFilter[T]) *PruneTreeCmd[T] {
	if filter == nil {
		return nil
	}

	cmd := &PruneTreeCmd[T]{
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
func (c *PruneTreeCmd[T]) GetOk() bool {
	return c.ok
}

// ExtractBranchCmd is a command that extracts the branch containing the
// given node.
type ExtractBranchCmd[T any] struct {
	// leaf is a pointer to the leaf to extract the branch from.
	leaf *TreeNode[T]

	// branch is a pointer to the branch extracted.
	branch *Branch[T]
}

// Execute implements the Debugging.Commander interface.
func (c *ExtractBranchCmd[T]) Execute(data *Tree[T]) error {
	c.branch = data.ExtractBranch(c.leaf, true)

	return nil
}

// Undo implements the Debugging.Commander interface.
func (c *ExtractBranchCmd[T]) Undo(data *Tree[T]) error {
	return nil
}

// Copy implements the Debugging.Commander interface.
func (c *ExtractBranchCmd[T]) Copy() uc.Copier {
	leafCopy := c.leaf.Copy().(*TreeNode[T])
	branchCopy := c.branch.Copy().(*Branch[T])

	cmdCopy := &ExtractBranchCmd[T]{
		leaf:   leafCopy,
		branch: branchCopy,
	}

	return cmdCopy
}

// NewExtractBranchCmd creates a new ExtractBranchCmd.
//
// Parameters:
//   - leaf: The leaf to extract the branch from.
//
// Returns:
//   - *ExtractBranchCmd: A pointer to the new ExtractBranchCmd.
func NewExtractBranchCmd[T any](leaf *TreeNode[T]) *ExtractBranchCmd[T] {
	cmd := &ExtractBranchCmd[T]{
		leaf: leaf,
	}
	return cmd
}

// GetBranch returns the value of the branch field.
//
// Call this function after executing the command.
//
// Returns:
//   - *Branch[T]: A pointer to the branch extracted.
func (c *ExtractBranchCmd[T]) GetBranch() *Branch[T] {
	branch := c.branch
	return branch
}

// InsertBranchCmd is a command that inserts a branch into the tree.
type InsertBranchCmd[T any] struct {
	// branch is a pointer to the branch to insert.
	branch *Branch[T]

	// hasError is true if an error occurred during execution, false otherwise.
	hasError bool
}

// Execute implements the Debugging.Commander interface.
func (c *InsertBranchCmd[T]) Execute(data *Tree[T]) error {
	ok := data.InsertBranch(c.branch)
	if !ok {
		c.hasError = true
	}

	return nil
}

// Undo implements the Debugging.Commander interface.
func (c *InsertBranchCmd[T]) Undo(data *Tree[T]) error {
	err := data.DeleteBranchContaining(c.branch.from_node)
	if err != nil {
		if c.hasError {
			return nil
		}

		return err
	} else if c.hasError {
		return errors.New("error occurred during execution")
	}

	return nil
}

// Copy implements the Debugging.Commander interface.
func (c *InsertBranchCmd[T]) Copy() uc.Copier {
	branchCopy := c.branch.Copy().(*Branch[T])

	cmdCopy := &InsertBranchCmd[T]{
		branch:   branchCopy,
		hasError: c.hasError,
	}

	return cmdCopy
}

// NewInsertBranchCmd creates a new InsertBranchCmd.
//
// Parameters:
//   - branch: The branch to insert.
//
// Returns:
//   - *InsertBranchCmd: A pointer to the new InsertBranchCmd.
func NewInsertBranchCmd[T any](branch *Branch[T]) *InsertBranchCmd[T] {
	if branch == nil {
		return nil
	}

	cmd := &InsertBranchCmd[T]{
		branch: branch,
	}

	return cmd
}
