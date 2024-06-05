package SiteNavigator

import (
	lls "github.com/PlayerR9/MyGoLib/ListLike/Stacker"
	us "github.com/PlayerR9/MyGoLib/Units/Slice"
	ue "github.com/PlayerR9/MyGoLib/Units/errors"
	"golang.org/x/net/html"
)

// FilterErrFunc is a function that returns an error if a condition is not met.
//
// Parameters:
//   - node: The node to check.
//
// Returns:
//   - error: An error if the condition is not met.
type FilterErrFunc func(node *html.Node) error

// FilterDataNode returns an FilterErrFunc that checks if the node has the specified data.
//
// Parameters:
//   - data: The data to check for.
//
// Returns:
//   - FilterErrFunc: The FilterErrFunc that checks if the node has the specified data.
func FilterDataNode(data string) FilterErrFunc {
	return func(node *html.Node) error {
		fn := NewSearchCriteria(html.ElementNode).SetData(data).Build()

		ok := fn(node)
		if !ok {
			return NewErrNoDataNodeFound(data)
		} else {
			return nil
		}
	}
}

// FilterTextNode returns an FilterErrFunc that checks if the node is a text node.
//
// Parameters:
//   - checkFirstChild: If true, the function checks if the first child is a text node.
//     Otherwise, it checks if the node itself is a text node.
//
// Returns:
//   - FilterErrFunc: The FilterErrFunc that checks if the node is a text node.
func FilterTextNode(checkFirstChild bool) FilterErrFunc {
	if checkFirstChild {
		return func(node *html.Node) error {
			ok := IsTextNodeSearch(node.FirstChild)
			if !ok {
				return NewErrNoTextNodeFound(true)
			} else {
				return nil
			}
		}
	} else {
		return func(node *html.Node) error {
			ok := IsTextNodeSearch(node)
			if !ok {
				return NewErrNoTextNodeFound(false)
			} else {
				return nil
			}
		}
	}
}

// filterValidNodes is a helper function that filters the valid nodes from a list.
//
// Parameters:
//   - list: The list of nodes to filter.
//   - filter: The function to check if a node is valid.
//
// Returns:
//   - []*html.Node: The list of valid nodes.
//   - error: An error if no valid nodes are found.
//
// Behaviors:
//   - If no valid nodes are found, the function returns the first error encountered.
//   - If list is empty, the function returns nil, nil.
//
// Assumptions:
//   - list has at least one element.
//   - filter is not nil.
func filterValidNodes(list []*html.Node, filter FilterErrFunc) ([]*html.Node, error) {
	var el ue.ErrOrSol[*html.Node]

	for i := 0; i < len(list); i++ {
		err := filter(list[i])
		if err == nil {
			el.AddSol(list[i], 0)
		} else {
			el.AddErr(err, 0)
		}
	}

	if el.HasError() {
		return nil, el.GetErrors()[0]
	} else {
		return el.GetSolutions(), nil
	}
}

// FilterValidNodes filters the valid nodes from a list.
//
// Parameters:
//   - list: The list of nodes to filter.
//   - filters: The functions to check if a node is valid.
//
// Returns:
//   - []*html.Node: The list of valid nodes.
//   - error: An error if no valid nodes are found.
//
// Behaviors:
//   - If no valid nodes are found, the function returns the first error encountered.
//   - If list is empty or filters is empty, the function returns list, nil.
func FilterValidNodes(list []*html.Node, filters []FilterErrFunc) ([]*html.Node, error) {
	if len(filters) == 0 {
		return list, nil
	}

	var el ue.ErrOrSol[*html.Node]

	for level, filter := range filters {
		if len(list) == 0 {
			break
		}

		for i := 0; i < len(list); i++ {
			err := filter(list[i])
			if err == nil {
				el.AddSol(list[i], level)
			} else {
				el.AddErr(err, level)
			}
		}

		if el.HasError() {
			return nil, el.GetErrors()[0]
		}

		list = el.GetSolutions()
	}

	return list, nil
}

// NodeListParser is a function that parses a list of nodes.
//
// Parameters:
//   - list: The list of nodes to parse.
//
// Returns:
//   - T: The parsed value.
//   - error: An error if the parsing fails.
type NodeListParser[T any] func(list []*html.Node) (T, error)

// CreateExtractor creates a NodeListParser from the given parameters.
//
// Parameters:
//   - parse: The function that parses the list of nodes.
//   - filters: The functions that filter the list of nodes.
//
// Returns:
//   - NodeListParser: The created NodeListParser.
//
// Behaviors:
//   - If parse is nil, the function returns a NodeListParser that returns
//     the error *errors.ErrInvalidParameter.
//   - Nil functions in filters are ignored.
func CreateExtractor[T any](parse NodeListParser[T], filters ...FilterErrFunc) NodeListParser[T] {
	if parse == nil {
		return func(list []*html.Node) (T, error) {
			return *new(T), ue.NewErrNilParameter("parse")
		}
	}

	filters = us.SliceFilter(filters, FilterNilFEFuncs)
	if len(filters) == 0 {
		return parse
	}

	return func(list []*html.Node) (T, error) {
		var err error

		for _, extract := range filters {
			if len(list) == 0 {
				break
			}

			list, err = filterValidNodes(list, extract)
			if err != nil {
				return *new(T), err
			}
		}

		res, err := parse(list)
		if err != nil {
			return *new(T), err
		}

		return res, nil
	}
}

// ActionType is an enumeration of the different actions that can be performed on a node.
type ActionType int8

const (
	// OnlyDirectChildren is an action that extracts only the direct children of a node.
	OnlyDirectChildren ActionType = iota

	// DFSOne is an action that extracts only one node using depth-first search.
	DFSOne

	// BFSMany is an action that extracts multiple nodes using breadth-first search.
	BFSMany
)

// GTEFunc is a function that extracts nodes from a tree.
//
// Parameters:
//   - tree: The tree to extract nodes from.
//
// Returns:
//   - []*html.Node: The list of nodes extracted from the tree.
type GTEFunc func(tree *HtmlTree) []*html.Node

// GenericTreeExtraction creates a GTEFunc from the given parameters.
//
// Parameters:
//   - search: The search criteria to use.
//   - action: The action to perform on the node.
//
// Returns:
//   - GTEFunc: The created GTEFunc.
//
// Behaviors:
//   - If search is nil, the function uses a nil filter.
//   - If action is not recognized, the function returns a GTEFunc that returns nil.
func GenericTreeExtraction(search *SearchCriteria, action ActionType) GTEFunc {
	var filter us.PredicateFilter[*html.Node]

	if search != nil {
		filter = search.Build()
	} else {
		filter = nil
	}

	switch action {
	case OnlyDirectChildren:
		return func(tree *HtmlTree) []*html.Node {
			return tree.ExtractSpecificNode(filter)
		}
	case BFSMany:
		return func(tree *HtmlTree) []*html.Node {
			return tree.ExtractNodes(filter)
		}
	case DFSOne:
		return func(tree *HtmlTree) []*html.Node {
			node := tree.ExtractContentFromDocument(filter)
			if node != nil {
				return []*html.Node{node}
			} else {
				return nil
			}
		}
	default:
		return func(tree *HtmlTree) []*html.Node {
			return nil
		}
	}
}

// CEWithSearch creates a NodeListParser from the given parameters.
//
// Parameters:
//   - search: The search criteria to use.
//   - action: The action to perform on the node.
//   - parse: The function that parses the list of nodes.
//   - filters: The functions that filter the list of nodes.
//
// Returns:
//   - NodeListParser: The created NodeListParser.
//
// Behaviors:
//   - If parse is nil, the function returns a NodeListParser that returns
//     the error *errors.ErrInvalidParameter.
//   - Nil functions in filters are ignored.
//   - Uses a Stack to traverse the tree and GTEFunc to extract nodes.
//   - It terminates as soon as a valid result is found.
func CEWithSearch[T any](search *SearchCriteria, action ActionType, parse NodeListParser[T], filters ...FilterErrFunc) NodeListParser[T] {
	if parse == nil {
		return func(list []*html.Node) (T, error) {
			return *new(T), ue.NewErrNilParameter("parse")
		}
	}

	filters = us.SliceFilter(filters, FilterNilFEFuncs)
	searchFunc := GenericTreeExtraction(search, action)

	return func(list []*html.Node) (T, error) {
		if len(list) == 0 {
			return *new(T), NewErrNoNodesFound()
		}

		S := lls.NewArrayStack(list...)

		var el ue.ErrOrSol[*html.Node]

		for {
			node, err := S.Pop()
			if err != nil {
				break
			}

			tree, err := NewHtmlTree(node)
			if err != nil {
				el.AddErr(err, 0)
				continue
			}

			newList := searchFunc(tree)

			if len(newList) != 0 {
				newList, err = FilterValidNodes(newList, filters)
				if err != nil {
					el.AddErr(err, 1)
					continue
				}
			}

			res, err := parse(newList)
			if err == nil {
				return res, nil
			}

			el.AddErr(err, 1+len(filters))
		}

		errList := el.GetErrors()

		return *new(T), ue.NewErrPossibleError(NewErrNoNodesFound(), errList[0])
	}
}
