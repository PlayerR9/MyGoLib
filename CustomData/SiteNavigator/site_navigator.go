package SiteNavigator

import (
	"golang.org/x/net/html"

	tr "github.com/PlayerR9/MyGoLib/CustomData/Tree"
	ers "github.com/PlayerR9/MyGoLib/Units/Errors"

	slext "github.com/PlayerR9/MyGoLib/Utility/SliceExt"
)

// IsTextNodeSearch is a search criteria that matches text nodes.
var IsTextNodeSearch slext.PredicateFilter[*html.Node] = NewSearchCriteria(html.TextNode).Build()

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

// HtmlTree is a struct that represents an HTML tree.
type HtmlTree struct {
	// The tree constructed from the HTML node.
	tree *tr.Tree[*html.Node]
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
//   - *ers.ErrNilValue: If any html.Node is nil.
func NewHtmlTree(root *html.Node) (*HtmlTree, error) {
	tree, err := tr.NoInfoMakeTree(
		root,
		func(elem *html.Node) ([]*html.Node, error) {
			if elem == nil {
				return nil, ers.NewErrNilValue()
			}

			children := make([]*html.Node, 0)

			for c := elem.FirstChild; c != nil; c = c.NextSibling {
				children = append(children, c)
			}

			return children, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return &HtmlTree{tree: tree}, nil
}

// ExtractSpecificNode finds all nodes that match the given search criteria and
// that are direct children of the provided node.
//
// Parameters:
//   - criteria: The search criteria to apply to each node.
//
// Returns:
//   - nodes: A slice containing all nodes that match the search criteria.
//
// Behavior:
//   - If no criteria is provided, then any node will match.
//   - If the node is nil, then a nil slice is returned.
func (t *HtmlTree) ExtractSpecificNode(matchFun slext.PredicateFilter[*html.Node]) []*html.Node {
	children := t.tree.GetDirectChildren()
	if len(children) == 0 {
		return nil
	}

	S := make([]*html.Node, 0, len(children))

	for _, child := range children {
		S = append(S, child.GetData())
	}

	S = slext.SliceFilter(S, matchFun)
	if len(S) == 0 {
		return nil
	}

	return S
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
func (t *HtmlTree) MatchNodes(matchFun slext.PredicateFilter[*html.Node]) []*html.Node {
	solution := make([]*html.Node, 0)

	err := tr.NoInfoTraverse(
		t.tree,
		func(node *html.Node) (bool, error) {
			if !matchFun(node) {
				return true, nil
			}

			solution = append(solution, node)
			return false, nil
		},
	).BFS()
	if err != nil {
		panic(err)
	}

	return solution
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
func (t *HtmlTree) ExtractContentFromDocument(matchFun slext.PredicateFilter[*html.Node]) *html.Node {
	if matchFun == nil {
		return nil
	}

	var solution *html.Node = nil

	err := tr.NoInfoTraverse(
		t.tree,
		func(node *html.Node) (bool, error) {
			if !matchFun(node) {
				return true, nil
			}

			solution = node
			return false, ers.NewErrNoError(nil)
		},
	).DFS()
	if err != nil {
		panic(err)
	}

	return solution
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
func (t *HtmlTree) ExtractNodes(criterias []slext.PredicateFilter[*html.Node]) []*html.Node {
	criterias = slext.FilterNilPredicates(criterias)
	if len(criterias) == 0 {
		return nil
	}

	todo := []*HtmlTree{t}

	for _, criteria := range criterias {
		newTodo := make([]*html.Node, 0)

		for _, tree := range todo {
			result := tree.MatchNodes(criteria)
			if len(result) != 0 {
				newTodo = append(newTodo, result...)
			}
		}

		if len(newTodo) == 0 {
			return nil
		}

		for _, node := range newTodo {
			newTree, err := NewHtmlTree(node)
			if err != nil {
				panic(err)
			}

			todo = append(todo, newTree)
		}
	}

	solution := make([]*html.Node, 0)

	for _, t := range todo {
		solution = append(solution, t.tree.Root().GetData())
	}

	return solution
}
