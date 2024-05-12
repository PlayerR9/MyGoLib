package SiteNavigator

import (
	"slices"

	"golang.org/x/net/html"

	cdp "github.com/PlayerR9/MyGoLib/CustomData/Pair"
	slext "github.com/PlayerR9/MyGoLib/Utility/SliceExt"
)

// SearchCriteria is a struct that encapsulates the parameters for searching
// within an HTML node.
type SearchCriteria struct {
	// NodeType specifies the type of the HTML node to search for.
	NodeType html.NodeType

	// Data represents the data contained within the node.
	Data *string

	// Attrs is a slice of attribute key-value pairs to match.
	Attrs []*cdp.Pair[string, slext.PredicateFilter[string]]
}

// NewSearchCriteria constructs a new SearchCriteria instance using the provided
// parameters.
//
// Parameters:
//   - node_type: The type of the HTML node to search for.
//
// Returns:
//   - *SearchCriteria: A new SearchCriteria instance.
func NewSearchCriteria(node_type html.NodeType) *SearchCriteria {
	return &SearchCriteria{
		NodeType: node_type,
	}
}

// SetData sets the data field of the SearchCriteria instance.
//
// Parameters:
//   - data: The data to set in the SearchCriteria instance.
//
// Returns:
//   - *SearchCriteria: The SearchCriteria instance with the data field set.
func (sc *SearchCriteria) SetData(data string) *SearchCriteria {
	sc.Data = &data

	return sc
}

// AppendAttr is a method of the SearchCriteria type that appends an attribute key-value
// pair to the SearchCriteria instance.
//
// Parameters:
//   - key: The attribute key to match.
//   - val: The attribute value to match.
//
// Returns:
//   - *SearchCriteria: The SearchCriteria instance with the attribute key-value pair appended.
func (sc *SearchCriteria) AppendAttr(key string, val slext.PredicateFilter[string]) *SearchCriteria {
	sc.Attrs = append(sc.Attrs, cdp.NewPair(key, val))

	return sc
}

// Build is a method of the SearchCriteria type that constructs a slext.PredicateFilter
// function using the search criteria.
//
// Returns:
//   - slext.PredicateFilter: A function that matches the search criteria.
func (sc *SearchCriteria) Build() slext.PredicateFilter[*html.Node] {
	attrsFunc := sc.Attrs
	nt := sc.NodeType

	if sc.Data != nil {
		data := *sc.Data

		return func(node *html.Node) bool {
			if node == nil || node.Type != nt || node.Data != data {
				return false
			}

			// Check if it matches the attribute
			for _, key := range attrsFunc {
				ok := slices.ContainsFunc(node.Attr, func(a html.Attribute) bool {
					return a.Key == key.First && key.Second(a.Val)
				})

				if !ok {
					return false
				}
			}

			return true
		}
	} else {
		return func(node *html.Node) bool {
			if node == nil || node.Type != nt {
				return false
			}

			for _, key := range attrsFunc {
				ok := slices.ContainsFunc(node.Attr, func(a html.Attribute) bool {
					return a.Key == key.First && key.Second(a.Val)
				})

				if !ok {
					return false
				}
			}

			return true
		}
	}
}
