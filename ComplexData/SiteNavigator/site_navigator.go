package SiteNavigator

import (
	"golang.org/x/net/html"

	tr "github.com/PlayerR9/MyGoLib/TreeLike/Tree"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
	slext "github.com/PlayerR9/MyGoLib/Units/slice"
)

var (
	// IsTextNodeSearch is a search criteria that matches text nodes.
	IsTextNodeSearch slext.PredicateFilter[*html.Node]
)

func init() {
	IsTextNodeSearch = NewSearchCriteria(html.TextNode).Build()
}

// GetDirectChildren returns a slice of the direct children of the provided node.
//
// Parameters:
//   - node: The HTML node to extract the children from.
//
// Returns:
//   - []*html.Node: A slice containing the direct children of the node.
func GetDirectChildren(node *html.Node) []*html.Node {
	if node == nil {
		return nil
	}

	children := make([]*html.Node, 0)
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		children = append(children, c)
	}
	return children
}

var (
	// GetChildrenFunc is a function that returns the children of an HTML node.
	GetChildrenFunc tr.NextsFunc = func(n tr.Noder, info uc.Copier) ([]tr.Noder, error) {
		if n == nil {
			return nil, uc.NewErrNilParameter("n")
		}

		elem, ok := n.(*tr.TreeNode[*html.Node])
		uc.Assert(ok, "GetChildrenFunc: n is not a *tr.TreeNode[*html.Node]")

		if elem.Data == nil {
			return nil, uc.NewErrNilParameter("n.Data")
		}

		var children []tr.Noder

		for c := elem.Data.FirstChild; c != nil; c = c.NextSibling {
			new_n := tr.NewTreeNode(c)

			children = append(children, new_n)
		}

		return children, nil
	}
)

// HtmlTree is a struct that represents an HTML tree.
type HtmlTree struct {
	// The tree constructed from the HTML node.
	tree *tr.Tree
}

// NewHtmlTree constructs a tree from an HTML node.
//
// Parameters:
//   - root: The root HTML node.
//
// Returns:
//   - *HtmlTree: The tree constructed from the HTML node.
//   - error: An error if the tree construction fails.
//
// Errors:
//   - *uc.ErrNilValue: If any html.Node is nil.
func NewHtmlTree(root *html.Node) (*HtmlTree, error) {
	var builder tr.Builder

	builder.SetNextFunc(GetChildrenFunc)

	n := tr.NewTreeNode(root)

	tree, err := builder.Build(n)
	if err != nil {
		return nil, err
	}

	ht := &HtmlTree{tree: tree}

	return ht, nil
}

// ExtractSpecificNode finds all nodes that match the given search criteria and
// that are direct children of the provided node.
//
// Parameters:
//   - matchFun: The search criteria to apply to each node.
//
// Returns:
//   - []*html.Node: A slice containing all nodes that match the search criteria.
//   - error: An error if the search fails.
//
// Behavior:
//   - If no criteria is provided, then any node will match.
func (t *HtmlTree) ExtractSpecificNode(matchFun slext.PredicateFilter[*html.Node]) ([]*html.Node, error) {
	if matchFun == nil {
		panic("Case not handled: matchFun is nil")
	}

	children, err := t.tree.GetDirectChildren()
	if err != nil {
		return nil, err
	}
	if len(children) == 0 {
		return nil, nil
	}

	S := make([]*html.Node, 0, len(children))

	for _, child := range children {
		n, ok := child.(*tr.TreeNode[*html.Node])
		uc.Assert(ok, "ExtractSpecificNode: child is not a *tr.TreeNode[*html.Node]")

		S = append(S, n.Data)
	}

	S = slext.SliceFilter(S, matchFun)
	return S, nil
}

// MatchNodes performs a breadth-first search on an HTML section returning a
// slice of nodes that match the provided search criteria.
//
// Parameters:
//   - matchFun: The search criteria to apply to each node.
//
// Returns:
//   - []*html.Node: A slice containing all nodes that match the search criteria.
//
// Behavior:
//   - It does not search the children of the nodes that match the criteria.
//   - If no criteria is provided, then the first node will match.
func (t *HtmlTree) MatchNodes(matchFun slext.PredicateFilter[*html.Node]) ([]*html.Node, error) {
	if matchFun == nil {
		panic("Case not handled: matchFun is nil")
	}

	var solution []*html.Node

	err := tr.BFS(t.tree, nil, func(node tr.Noder, info uc.Copier) (bool, error) {
		if node == nil {
			return false, uc.NewErrNilParameter("node")
		}

		n, ok := node.(*tr.TreeNode[*html.Node])
		uc.Assert(ok, "MatchNodes: node is not a *tr.TreeNode[*html.Node]")

		ok = matchFun(n.Data)
		if !ok {
			return true, nil
		}

		solution = append(solution, n.Data)
		return false, nil
	})
	if err != nil {
		return nil, err
	}

	return solution, nil
}

// ExtractContentFromDocument performs a depth-first search on an HTML document,
// finding the first node that matches the provided search criteria.
//
// Parameters:
//   - matchFun: The search criteria to apply to each node.
//
// Returns:
//   - *html.Node: The first node that matches the search criteria, nil if no
//     matching node is found.
func (t *HtmlTree) ExtractContentFromDocument(matchFun slext.PredicateFilter[*html.Node]) (*html.Node, error) {
	if matchFun == nil {
		panic("Case not handled: matchFun is nil")
	}

	var solution *html.Node

	err := tr.DFS(t.tree, nil, func(node tr.Noder, info uc.Copier) (bool, error) {
		if node == nil {
			return false, uc.NewErrNilParameter("node")
		}

		n, ok := node.(*tr.TreeNode[*html.Node])
		uc.Assert(ok, "ExtractContentFromDocument: node is not a *tr.TreeNode[*html.Node]")

		ok = matchFun(n.Data)
		if !ok {
			return true, nil
		}

		solution = n.Data
		return false, nil
	})
	if err != nil {
		return nil, err
	}

	return solution, nil
}

// ExtractNodes performs a breadth-first search on an HTML section returning a
// slice of nodes that match the provided search criteria.
//
// Parameters:
//   - criterias: A list of search criteria to apply to each node.
//
// Returns:
//   - []*html.Node: A slice containing all nodes that match the search criteria.
//
// Behavior:
//   - If no criteria is provided, then any node will match.
func (t *HtmlTree) ExtractNodes(criterias ...slext.PredicateFilter[*html.Node]) ([]*html.Node, error) {
	criterias = slext.FilterNilPredicates(criterias)
	if len(criterias) == 0 {
		return nil, nil
	}

	todo := []*HtmlTree{t}

	for i, criteria := range criterias {
		var new_todo []*html.Node

		for _, tree := range todo {
			result, err := tree.MatchNodes(criteria)
			if err != nil {
				return nil, uc.NewErrWhileAt("applying", i+1, "criteria", err)
			}

			if len(result) != 0 {
				new_todo = append(new_todo, result...)
			}
		}

		if len(new_todo) == 0 {
			return nil, nil
		}

		for i, node := range new_todo {
			new_tree, err := NewHtmlTree(node)
			if err != nil {
				return nil, uc.NewErrWhileAt("adding", i+1, "tree", err)
			}

			todo = append(todo, new_tree)
		}
	}

	var solution []*html.Node

	for _, t := range todo {
		root := t.tree.Root()

		val, ok := root.(*tr.TreeNode[*html.Node])
		uc.Assert(ok, "ExtractNodes: root is not a *tr.TreeNode[*html.Node]")

		solution = append(solution, val.Data)
	}

	return solution, nil
}
