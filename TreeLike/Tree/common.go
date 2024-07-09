package Tree

import (
	"slices"
)

// Treeer is an interface for types that can be converted to a tree.
type Treeer interface {
	// TreeOf converts the type to a tree.
	//
	// Returns:
	//   - *Tree: A pointer to the tree.
	TreeOf() *Tree
}

// TreeOf converts the element to a tree.
//
// Parameters:
//   - elem: The element to convert.
//
// Returns:
//   - *Tree: A pointer to the tree. Nil if the element is nil.
//
// Behaviors:
//   - If the element implements the Treeer interface, the function calls the TreeOf method.
//   - Otherwise, the function creates a new tree with the element as the root.
func TreeOf(elem Noder) *Tree {
	if elem == nil {
		return nil
	}

	var tree *Tree

	switch elem := elem.(type) {
	case Treeer:
		tree = elem.TreeOf()
	default:
		tree = NewTree(elem)
	}

	return tree
}

// CommonAncestor returns the first common ancestor of the two nodes.
//
// Parameters:
//   - n1: The first node.
//   - n2: The second node.
//
// Returns:
//   - *Node[T]: A pointer to the common ancestor. Nil if no such node is found.
func FindCommonAncestor(n1, n2 Noder) Noder {
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
		ok := slices.Contains(ancestors2, node)
		if ok {
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
func ExtractData[T Noder](nodes []*TreeNode[T]) []T {
	data := make([]T, 0, len(nodes))

	for _, node := range nodes {
		data = append(data, node.Data)
	}

	return data
}
