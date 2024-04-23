package SiteNavigator

import (
	"golang.org/x/net/html"

	Queue "github.com/PlayerR9/MyGoLib/ListLike/Queue"
	"github.com/PlayerR9/MyGoLib/ListLike/Stack"
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
		node := Q.Dequeue()

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
		node := S.Pop()

		if criteria == nil || criteria.Match(node) {
			return node
		}

		for c := node.FirstChild; c != nil; c = c.NextSibling {
			S.Push(c)
		}
	}

	return nil
}

// GetDirectChildren returns a slice of the direct children of the provided node.
//
// Parameters:
//
//   - node: The HTML node to extract the children from.
//
// Returns:
//
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
