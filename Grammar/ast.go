package Grammar

import (
	"fmt"
	"strconv"
	"strings"

	fs "github.com/PlayerR9/MyGoLib/Formatting/Strings"
	uc "github.com/PlayerR9/lib_units/common"
	"github.com/PlayerR9/lib_units/slices"
)

// AstNoder is an interface that defines the behavior of an AST node.
type AstNoder interface {
	fs.Noder
}

// NodeTyper is an interface that defines the behavior of a node type.
type NodeTyper interface {
	~int

	fmt.Stringer
}

// Node is a node in the AST.
type Node[N NodeTyper] struct {
	// Parent is the parent of the node.
	Parent *Node[N]

	// Children is the children of the node.
	Children []*Node[N]

	// Type is the type of the node.
	Type N

	// Data is the data of the node.
	Data string
}

// IsLeaf implements the AstNoder interface.
func (n *Node[N]) IsLeaf() bool {
	return len(n.Children) == 0
}

// Iterator implements the AstNoder interface.
func (n *Node[N]) Iterator() uc.Iterater[fs.Noder] {
	if len(n.Children) == 0 {
		return nil
	}

	nodes := make([]fs.Noder, 0, len(n.Children))
	for _, child := range n.Children {

		nodes = append(nodes, child)
	}

	return uc.NewSimpleIterator(nodes)
}

// String implements the AstNoder interface.
func (n *Node[N]) String() string {
	var builder strings.Builder

	builder.WriteString("Node[")
	builder.WriteString(n.Type.String())

	if n.Data != "" {
		builder.WriteString(" (")
		builder.WriteString(strconv.Quote(n.Data))
		builder.WriteRune(')')
	}

	builder.WriteRune(']')

	return builder.String()
}

// NewNode creates a new node.
//
// Parameters:
//   - t: The type of the node.
//   - data: The data of the node.
//
// Returns:
//   - *Node[N]: The new node. Never returns nil.
func NewNode[N NodeTyper](t N, data string) *Node[N] {
	return &Node[N]{
		Type: t,
		Data: data,
	}
}

// SetChildren sets the children of the node. Nil children are ignored.
func (n *Node[N]) SetChildren(children []*Node[N]) {
	children = slices.FilterNilValues(children)
	if len(children) == 0 {
		return
	}

	for _, child := range children {
		child.Parent = n
	}
	n.Children = children
}

// PrintAst stringifies the AST.
//
// Parameters:
//   - root: The root of the AST.
//
// Returns:
//   - string: The AST as a string.
func PrintAst(root AstNoder) string {
	if root == nil {
		return ""
	}

	str, err := fs.PrintTree(root)
	uc.AssertErr(err, "Strings.PrintTree(root)")

	return str
}

// LeftAstFunc is a function that parses the left-recursive AST.
//
// Parameters:
//   - children: The children of the current node.
//
// Returns:
//   - []N: The left-recursive AST.
//   - error: An error if the left-recursive AST could not be parsed.
type LeftAstFunc[N AstNoder, T TokenTyper] func(children []*Token[T]) ([]N, error)

// LeftRecursive parses the left-recursive AST.
//
// Parameters:
//   - root: The root of the left-recursive AST.
//   - lhs_type: The type of the left-hand side.
//   - f: The function that parses the left-recursive AST.
//
// Returns:
//   - []N: The left-recursive AST.
//   - error: An error if the left-recursive AST could not be parsed.
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

// ToAstFunc is a function that parses the AST.
//
// Parameters:
//   - root: The root of the AST.
//
// Returns:
//   - []N: The AST.
//   - error: An error if the AST could not be parsed.
type ToAstFunc[N AstNoder, T TokenTyper] func(root *Token[T]) ([]N, error)

// ToAst parses the AST.
//
// Parameters:
//   - root: The root of the AST.
//   - to_ast: The function that parses the AST.
//
// Returns:
//   - []N: The AST.
//   - error: An error if the AST could not be parsed.
//
// Errors:
//   - *common.ErrInvalidParameter: If the root is nil or the to_ast is nil.
//   - error: Any error returned by the to_ast function.
func ToAst[N AstNoder, T TokenTyper](root *Token[T], to_ast ToAstFunc[N, T]) ([]N, error) {
	if root == nil {
		return nil, uc.NewErrNilParameter("root")
	} else if to_ast == nil {
		return nil, uc.NewErrNilParameter("to_ast")
	}

	nodes, err := to_ast(root)
	if err != nil {
		return nodes, err
	}

	return nodes, nil
}

// ExtractData extracts the data from a token.
//
// Parameters:
//   - node: The token to extract the data from.
//
// Returns:
//   - string: The data of the token.
//   - error: An error if the data is not of type string or if the token is nil.
func ExtractData[T TokenTyper](node *Token[T]) (string, error) {
	if node == nil {
		return "", uc.NewErrNilParameter("node")
	}

	data, ok := node.Data.(string)
	if !ok {
		return "", fmt.Errorf("expected string, got %T instead", node.Data)
	}

	return data, nil
}

// ExtractChildren extracts the children from a token.
//
// Parameters:
//   - node: The token to extract the children from.
//
// Returns:
//   - []*gr.Token[T]: The children of the token.
//   - error: An error if the children is not of type []*gr.Token[T] or if the token is nil.
func ExtractChildren[T TokenTyper](node *Token[T]) ([]*Token[T], error) {
	if node == nil {
		return nil, uc.NewErrNilParameter("node")
	}

	children, ok := node.Data.([]*Token[T])
	if !ok {
		return nil, fmt.Errorf("expected []*Token, got %T instead", node.Data)
	}

	return children, nil
}
