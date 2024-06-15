package TreeExplorer

import (
	tr "github.com/PlayerR9/MyGoLib/TreeLike/Tree"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
	ers "github.com/PlayerR9/MyGoLib/Units/errors"
	us "github.com/PlayerR9/MyGoLib/Units/slice"
)

// TreeEvaluator is a tree evaluator that uses a grammar to tokenize a string.
type TreeEvaluator[R MatchResulter[O], M Matcher[R, O], O any] struct {
	// root is the root node of the tree evaluator.
	root *tr.Tree[*CurrentEval[O]]

	// matcher is the matcher used by the tree evaluator.
	matcher M

	// filterBranches is a list of functions that filter branches.
	filters []FilterBranchesFunc[O]
}

// NewTreeEvaluator creates a new tree evaluator.
//
// Parameters:
//   - matcher: The matcher that the tree evaluator will use.
//
// Returns:
//   - *TreeEvaluator: A pointer to the new tree evaluator.
func NewTreeEvaluator[R MatchResulter[O], M Matcher[R, O], O any](filters ...FilterBranchesFunc[O]) *TreeEvaluator[R, M, O] {
	te := &TreeEvaluator[R, M, O]{
		filters: filters,
	}

	return te
}

// addMatchLeaves adds the matches to a root tree as leaves.
//
// Parameters:
//   - root: The root of the tree to add the leaves to.
//   - matches: The matches to add to the tree evaluator.
func (te *TreeEvaluator[R, M, O]) addMatchLeaves(root *tr.Tree[*CurrentEval[O]], matches []R) {
	// Get the longest match.
	matches = te.matcher.SelectBestMatches(matches)

	children := make([]*tr.Tree[*CurrentEval[O]], 0, len(matches))

	for _, match := range matches {
		currMatch := match.GetMatch()
		ht := NewCurrentEval(currMatch)

		tree := tr.NewTree(ht)

		children = append(children, tree)
	}

	root.SetChildren(children)
}

// processLeaves processes the leaves in the tree evaluator.
//
// Returns:
//   - bool: True if all leaves are complete, false otherwise.
//   - error: An error of type *ErrAllMatchesFailed if all matches failed.
func (te *TreeEvaluator[R, M, O]) processLeaves() uc.EvalManyFunc[*CurrentEval[O], *CurrentEval[O]] {
	filterFunc := func(data *CurrentEval[O]) ([]*CurrentEval[O], error) {
		nextAt := te.matcher.GetNext(data.Elem)

		ok := te.matcher.IsDone(nextAt)
		if ok {
			data.SetStatus(EvalComplete)

			return nil, nil
		}

		matches, err := te.matcher.Match(nextAt)
		if err != nil {
			data.SetStatus(EvalError)

			return nil, nil
		}

		// Get the longest match.
		matches = te.matcher.SelectBestMatches(matches)

		children := make([]*CurrentEval[O], 0, len(matches))

		for _, match := range matches {
			curr := match.GetMatch()
			ht := NewCurrentEval(curr)

			children = append(children, ht)
		}

		data.SetStatus(EvalComplete)

		return children, nil
	}

	return filterFunc
}

// canContinue returns true if the tree evaluator can continue.
//
// Returns:
//   - bool: True if the tree evaluator can continue, false otherwise.
func (te *TreeEvaluator[R, M, O]) canContinue() bool {
	leaves := te.root.GetLeaves()

	for _, leaf := range leaves {
		if leaf.Data.Status == EvalIncomplete {
			return true
		}
	}

	return false
}

// pruneTree prunes the tree evaluator.
//
// Parameters:
//   - filter: The filter to use to prune the tree.
//
// Returns:
//   - bool: True if no nodes were pruned, false otherwise.
func (te *TreeEvaluator[R, M, O]) pruneTree(filter us.PredicateFilter[*CurrentEval[O]]) bool {
	for te.root.Size() != 0 {
		target := te.root.SearchNodes(filter)
		if target == nil {
			return true
		}

		te.root.DeleteBranchContaining(target)
	}

	return false
}

// Evaluate is the main function of the tree evaluator.
//
// Parameters:
//   - source: The source to evaluate.
//   - root: The root of the tree evaluator.
//
// Returns:
//   - error: An error if lexing fails.
//
// Errors:
//   - *ErrEmptyInput: The source is empty.
//   - *ers.ErrAt: An error occurred at a specific index.
//   - *ErrAllMatchesFailed: All matches failed.
func (te *TreeEvaluator[R, M, O]) Evaluate(matcher M, root O) error {
	te.matcher = matcher

	ce := NewCurrentEval(root)
	te.root = tr.NewTree(ce)

	matches, err := te.matcher.Match(0)
	if err != nil {
		return ers.NewErrAt(0, "position", err)
	}

	te.addMatchLeaves(te.root, matches)

	te.root.Root().Data.SetStatus(EvalComplete)

	shouldContinue := true

	for shouldContinue {
		err := te.root.ProcessLeaves(te.processLeaves())
		if err != nil {
			return err
		}

		ok := te.pruneTree(FilterErrorLeaves)
		if !ok {
			return NewErrAllMatchesFailed()
		}

		shouldContinue = te.canContinue()
	}

	ok := te.pruneTree(FilterIncompleteLeaves)
	if !ok {
		return NewErrAllMatchesFailed()
	}

	return nil
}

// GetBranches returns the tokens that have been lexed.
//
// Remember to use Lexer.RemoveToSkipTokens() to remove tokens that
// are not needed for the parser (i.e., marked as to skip in the grammar).
//
// Returns:
//   - result: The tokens that have been lexed.
//   - reason: An error if the tree evaluator has not been run yet.
func (te *TreeEvaluator[R, M, O]) GetBranches() ([][]*CurrentEval[O], error) {
	if te.root == nil {
		return nil, ers.NewErrInvalidUsage(
			ers.NewErrNilValue(),
			"must call TreeEvaluator.Evaluate() first",
		)
	}

	tokenBranches := te.root.SnakeTraversal()

	branches, invalidTokIndex := filterInvalidBranches(tokenBranches)
	if invalidTokIndex != -1 {
		return branches, ers.NewErrAt(invalidTokIndex, "token", NewErrInvalidElement())
	}

	var err error

	for _, filter := range te.filters {
		branches, err = filter(branches)
		if err != nil {
			return branches, err
		}
	}

	te.root = nil

	return branches, nil
}
