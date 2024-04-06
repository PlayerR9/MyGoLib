package SiteNavigator

import (
	"slices"

	"github.com/markphelps/optional"
	"golang.org/x/net/html"
)

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
