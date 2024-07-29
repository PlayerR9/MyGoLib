package Grammar

import (
	"fmt"

	fs "github.com/PlayerR9/MyGoLib/Formatting/Strings"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
)

type AstNoder interface {
	fs.Noder
}

func PrintAst(root AstNoder) string {
	if root == nil {
		return ""
	}

	str, err := fs.PrintTree(root)
	uc.AssertErr(err, "Strings.PrintTree(root)")

	return str
}

type LeftAstFunc[N AstNoder, T TokenTyper] func(children []*Token[T]) ([]N, error)

func LeftRecursive[N AstNoder, T TokenTyper](root *Token[T], lhs_type T, f LeftAstFunc[N, T]) ([]N, error) {
	uc.AssertNil(root, "root")

	var nodes []N

	for root != nil {
		if root.Type != lhs_type {
			return nodes, fmt.Errorf("expected %q, got %q instead", lhs_type.String(), root.Type.String())
		}

		children, ok := root.Data.([]*Token[T])
		if !ok {
			return nodes, fmt.Errorf("expected non-leaf node, got leaf node instead")
		} else if len(children) == 0 {
			return nodes, fmt.Errorf("expected at least 1 child, got 0 children instead")
		}

		last_child := children[len(children)-1]

		if last_child.Type == lhs_type {
			children = children[:len(children)-1]
			root = last_child
		} else {
			root = nil
		}

		sub_nodes, err := f(children)
		if len(sub_nodes) > 0 {
			nodes = append(nodes, sub_nodes...)
		}

		if err != nil {
			return nodes, fmt.Errorf("in %q: %w", root.Type.String(), err)
		}
	}

	return nodes, nil
}
