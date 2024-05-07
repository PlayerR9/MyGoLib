package Tree

import (
	"slices"
)

// CommonAncestor returns the first common ancestor of the two nodes.
//
// Parameters:
//   - n1: The first node.
//   - n2: The second node.
//
// Returns:
//   - *Node[T]: A pointer to the common ancestor. Nil if no such node is found.
func FindCommonAncestor[T any](n1, n2 *Node[T]) *Node[T] {
	if n1 == nil {
		return n2
	} else if n2 == nil {
		return n1
	} else if n1 == n2 {
		return n1
	}

	ancestors1 := n1.GetAncestors()
	ancestors2 := n2.GetAncestors()

	if len(ancestors1) > len(ancestors2) {
		ancestors1, ancestors2 = ancestors2, ancestors1
	}

	for _, node := range ancestors1 {
		if slices.Contains(ancestors2, node) {
			return node
		}
	}

	return nil
}

// ExtractData returns the values of the nodes in the slice.
//
// Parameters:
//   - nodes: The nodes to extract the values from.
//
// Returns:
//   - []T: A slice of the values of the nodes.
func ExtractData[T any](nodes []*Node[T]) []T {
	data := make([]T, 0, len(nodes))

	for _, node := range nodes {
		data = append(data, node.Data)
	}

	return data
}
