package SiteNavigator

import (
	"slices"

	"github.com/markphelps/optional"
	"golang.org/x/net/html"

	Queue "github.com/PlayerR9/MyGoLib/CustomData/ListLike/Queue"
	Stack "github.com/PlayerR9/MyGoLib/CustomData/ListLike/Stack"
)

var IsTextNodeSearch *SearchCriteria = NewSearchCriteria(html.TextNode)

// ExtractSpecificNode finds all nodes that match the given search criteria and
// that are direct children of the provided node.
//
// Panics with an error of type *ers.InvalidParameterError if the node is nil.
//
// Parameters:
//
//   - node: The HTML node to search within.
//   - criteria: The search criteria to apply to each node.
//
// Returns:
//
//   - nodes: A slice containing all nodes that match the search criteria.
//
// If no criteria is provided, then any node will match.
func ExtractSpecificNode(node *html.Node, criteria *SearchCriteria) []*html.Node {
	if node == nil {
		return nil
	}

	nodes := make([]*html.Node, 0)

	// If no criteria is provided, then any node will match
	if criteria == nil {
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			nodes = append(nodes, c)
		}
	} else {
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			if criteria.Match(c) {
				nodes = append(nodes, c)
			}
		}
	}

	return nodes
}

// MatchNodes performs a breadth-first search on an HTML section returning a
// slice of nodes that match the provided search criteria.
// It does not search the children of the nodes that match the criteria.
// If no criteria is provided, then the first node will match.
//
// Parameters:
//
//   - section: The HTML section to search within.
//   - criteria: The search criteria to apply to each node.
//
// Returns:
//
//   - []*html.Node: A slice containing all nodes that match the search criteria.
func MatchNodes(section *html.Node, criteria *SearchCriteria) []*html.Node {
	if section == nil {
		return nil // No nodes to extract
	} else if criteria == nil {
		return []*html.Node{section}
	}

	solution := make([]*html.Node, 0)
	Q := Queue.NewLinkedQueue(section)

	for !Q.IsEmpty() {
		node, _ := Q.Dequeue()

		if criteria.Match(node) {
			solution = append(solution, node)
		} else {
			// Search the children of the node
			for c := node.FirstChild; c != nil; c = c.NextSibling {
				Q.Enqueue(c)
			}
		}
	}

	return solution
}

// ExtractNodes performs a breadth-first search on an HTML section returning a
// slice of nodes that match the provided search criteria.
//
// Parameters:
//
//   - section: The HTML section to search within.
//   - criterias: A list of search criteria to apply to each node.
//
// Returns:
//
//   - []*html.Node: A slice containing all nodes that match the search criteria.
//
// If no criteria is provided, then any node will match.
func ExtractNodes(section *html.Node, criterias []*SearchCriteria) []*html.Node {
	if section == nil {
		return nil // No nodes to extract
	}

	solution := []*html.Node{section}

	for _, criteria := range criterias {
		partialSol := make([]*html.Node, 0)

		for _, node := range solution {
			partialSol = append(partialSol, MatchNodes(node, criteria)...)
		}

		if len(partialSol) == 0 {
			return nil
		}

		solution = partialSol
	}

	return solution
}

// ExtractContentFromDocument performs a depth-first search on an HTML document,
// finding the first node that matches the provided search criteria.
//
// Parameters:
//
//   - doc: The HTML document to search within.
//   - criteria: The search criteria to apply to each node.
//
// Returns:
//
//   - *html.Node: The first node that matches the search criteria, nil if no
//     matching node is found.
func ExtractContentFromDocument(doc *html.Node, criteria *SearchCriteria) *html.Node {
	if doc == nil {
		return nil
	}

	S := Stack.NewLinkedStack(doc)

	for !S.IsEmpty() {
		node, _ := S.Pop()

		if criteria == nil || criteria.Match(node) {
			return node
		}

		for c := node.FirstChild; c != nil; c = c.NextSibling {
			S.Push(c)
		}
	}

	return nil
}

// SearchCriteria is a struct that encapsulates the parameters for searching
// within an HTML node.
type SearchCriteria struct {
	// NodeType specifies the type of the HTML node to search for.
	NodeType html.NodeType

	// Data represents the data contained within the node.
	Data optional.String

	// AttrKey and AttrVal define the attribute key-value pair to match in the
	// node.
	AttrKey []string
	AttrVal []func(string) bool
}

// SearchCriteriaOption is a functional option type for the SearchCriteria struct.
//
// Parameters:
//
//   - *SearchCriteria: The SearchCriteria instance to apply the option to.
type SearchCriteriaOption func(*SearchCriteria)

// WithData is a functional option that sets the data field of the SearchCriteria
// instance to the provided string.
//
// Parameters:
//
//   - data: The data to set in the SearchCriteria instance.
//
// Returns:
//
//   - SearchCriteriaOption: A functional option that sets the data field of the
func WithData(data string) SearchCriteriaOption {
	return func(sc *SearchCriteria) {
		sc.Data = optional.NewString(data)
	}
}

// WithAttr is a functional option that sets the attribute key-value pair to match
// in the SearchCriteria instance.
//
// Parameters:
//
//   - key: The attribute key to match.
//   - val: The attribute value to match.
//
// Returns:
//
//   - SearchCriteriaOption: A functional option that sets the attribute key-value
//     pair to match in the SearchCriteria instance.
func WithAttr(key string, val func(string) bool) SearchCriteriaOption {
	return func(sc *SearchCriteria) {
		sc.AttrKey = append(sc.AttrKey, key)
		sc.AttrVal = append(sc.AttrVal, val)
	}
}

// NewSearchCriteria constructs a new SearchCriteria instance using the provided
// parameters.
//
// Parameters:
//
//   - node_type: The type of the HTML node to search for.
//   - options: A variadic list of functional options to apply to the SearchCriteria
//     instance.
//
// Returns:
//
//   - *SearchCriteria: The newly created SearchCriteria instance.
func NewSearchCriteria(node_type html.NodeType, options ...SearchCriteriaOption) *SearchCriteria {
	sc := &SearchCriteria{
		NodeType: node_type,
	}

	for _, option := range options {
		option(sc)
	}

	return sc
}

// Match is a method of the SearchCriteria type that evaluates whether a given HTML
// node matches the search criteria encapsulated by the SearchCriteria instance.
//
// Parameters:
//
//   - node: The HTML node to match against the search criteria.
//
// Returns:
//
//   - bool: True if the node matches the search criteria, otherwise false.
func (sc *SearchCriteria) Match(node *html.Node) bool {
	if node == nil || node.Type != sc.NodeType {
		return false
	}

	matchesData := true
	sc.Data.If(func(data string) {
		matchesData = node.Data == data
	})
	if !matchesData {
		return false
	}

	// Check if it matches the attribute
	for i, key := range sc.AttrKey {
		if !slices.ContainsFunc(node.Attr, func(a html.Attribute) bool {
			return a.Key == key && sc.AttrVal[i](a.Val)
		}) {
			return false
		}
	}

	return true
}

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
