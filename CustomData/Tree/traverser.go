package Tree

import (
	"github.com/PlayerR9/MyGoLib/ListLike/Queuer"
	"github.com/PlayerR9/MyGoLib/ListLike/Stacker"
	intf "github.com/PlayerR9/MyGoLib/Units/Common"
	up "github.com/PlayerR9/MyGoLib/Units/Pair"

	ers "github.com/PlayerR9/MyGoLib/Units/Errors"
)

// ObserverFunc is a function that observes a node.
//
// Parameters:
//   - data: The data of the node.
//   - info: The info of the node.
//
// Returns:
//   - bool: True if the traversal should continue, otherwise false.
//   - error: An error if the observation fails.
type ObserverFunc[T any, I intf.Copier] func(data T, info I) (bool, error)

// helper is a helper struct for traversing the tree.
type helper[T any, I intf.Copier] up.Pair[*TreeNode[T], I]

// newHelper creates a new helper.
//
// Parameters:
//   - node: The current node.
//   - info: The info of the current node.
//
// Returns:
//   - *helper[T, I]: A pointer to the helper.
func newHelper[T any, I intf.Copier](node *TreeNode[T], info I) *helper[T, I] {
	return &helper[T, I]{
		First:  node,
		Second: info.Copy().(I),
	}
}

// Traversor is a struct that traverses a tree.
type Traversor[T any, I intf.Copier] struct {
	// The helper struct.
	h *helper[T, I]

	// The observer function.
	observe ObserverFunc[T, I]
}

// DFS traverses the tree in depth-first order.
//
// Returns:
//   - error: An error if the traversal fails.
func (t *Traversor[T, I]) DFS() error {
	if t.h.First == nil || t.observe == nil {
		return nil
	}

	S := Stacker.NewLinkedStack(newHelper(t.h.First, t.h.Second))

	for {
		top, err := S.Pop()
		if err != nil {
			break
		}

		ok, err := t.observe(top.First.Data, top.Second)
		if err != nil {
			return err
		}

		if !ok {
			continue
		}

		for _, child := range top.First.children {
			err := S.Push(newHelper(child, top.Second))
			if err != nil {
				panic(err)
			}
		}
	}

	return nil
}

// BFS traverses the tree in breadth-first order.
//
// Returns:
//   - error: An error if the traversal fails.
func (t *Traversor[T, I]) BFS() error {
	if t.h.First == nil || t.observe == nil {
		return nil
	}

	Q := Queuer.NewLinkedQueue(newHelper(t.h.First, t.h.Second))

	for {
		first, err := Q.Dequeue()
		if err != nil {
			break
		}

		ok, err := t.observe(first.First.Data, first.Second)
		if err != nil {
			return err
		}

		if !ok {
			continue
		}

		for _, child := range first.First.children {
			err := Q.Enqueue(newHelper(child, first.Second))
			if err != nil {
				panic(err)
			}
		}
	}

	return nil
}

// Traverse creates a new traversor for the tree.
//
// Parameters:
//   - tree: The tree to traverse.
//   - init: The initial info.
//   - f: The observer function.
//
// Returns:
//   - Traversor[T, I]: The traversor.
func Traverse[T any, I intf.Copier](tree *Tree[T], init I, f ObserverFunc[T, I]) *Traversor[T, I] {
	var root *TreeNode[T]

	if tree == nil {
		root = nil
	} else {
		root = tree.root
	}

	return &Traversor[T, I]{
		h: &helper[T, I]{
			First:  root,
			Second: init,
		},
		observe: f,
	}
}

// NextsFunc is a function that returns the next elements.
//
// Parameters:
//   - elem: The current element.
//   - info: The info of the current element.
//
// Returns:
//   - []T: The next elements.
//   - error: An error if the next elements cannot be found.
type NextsFunc[T any, I intf.Copier] func(elem T, info I) ([]T, error)

// MakeTree creates a tree from the given element.
//
// Parameters:
//   - elem: The element to start the tree from.
//
// Returns:
//   - *Tree[T]: The tree created from the element.
//   - error: An error if the tree cannot be created.
func MakeTree[T any, I intf.Copier](elem T, info I, f NextsFunc[T, I]) (*Tree[T], error) {
	// 1. Handle the first element
	h := newHelper(newTreeNode(elem), info)

	nexts, err := f(h.First.Data, h.Second)
	if err != nil {
		return nil, err
	}

	tree := &Tree[T]{
		root: h.First,
	}

	if len(nexts) == 0 {
		return tree, nil
	}

	// 2. Create a stack and push the first element
	type StackElement struct {
		Prev *TreeNode[T]
		Elem *helper[T, I]
	}

	S := Stacker.NewLinkedStack[StackElement]()

	for _, next := range nexts {
		err := S.Push(StackElement{Prev: tree.root, Elem: newHelper(newTreeNode(next), h.Second)})
		if err != nil {
			panic(err)
		}
	}

	for {
		top, err := S.Pop()
		if err != nil {
			break
		}

		nexts, err := f(top.Elem.First.Data, top.Elem.Second)
		if err != nil {
			return nil, err
		}

		top.Prev.AddChildren(top.Elem.First)

		for _, next := range nexts {
			err := S.Push(StackElement{Prev: top.Elem.First, Elem: newHelper(newTreeNode(next), top.Elem.Second)})
			if err != nil {
				panic(err)
			}
		}
	}

	return tree, nil
}

// NoInfo is a struct that contains no information.
type NoInfo struct{}

// Copy creates a copy of the NoInfo.
//
// Returns:
//   - Copier: The copy of the NoInfo.
func (n *NoInfo) Copy() intf.Copier {
	return &NoInfo{}
}

// NewNoInfo creates a new NoInfo. This is useful for when traversing a tree
// but no additional information is needed.
//
// Returns:
//   - *NoInfo: The NoInfo.
func NewNoInfo() *NoInfo {
	return &NoInfo{}
}

// ObserverFunc is a function that observes a node.
//
// Parameters:
//   - data: The data of the node.
//   - info: The info of the node.
//
// Returns:
//   - bool: True if the traversal should continue, otherwise false.
//   - error: An error if the observation fails.
type NoInfoObserverFunc[T any] func(data T) (bool, error)

// NoInfoTraversor is a struct that traverses a tree.
type NoInfoTraversor[T any] struct {
	// The helper struct.
	node *TreeNode[T]

	// The observer function.
	observe NoInfoObserverFunc[T]
}

// DFS traverses the tree in depth-first order.
//
// Returns:
//   - error: An error if the traversal fails.
//
// Behavior:
//   - Use *ers.ErrNoError to stop the traversal without an error.
func (t *NoInfoTraversor[T]) DFS() error {
	if t.node == nil || t.observe == nil {
		return nil
	}

	S := Stacker.NewLinkedStack(t.node)

	for {
		top, err := S.Pop()
		if err != nil {
			break
		}

		ok, err := t.observe(top.Data)
		if err != nil {
			if ers.As[*ers.ErrNoError](err) {
				return nil
			} else {
				return err
			}
		}

		if !ok {
			continue
		}

		for _, child := range top.children {
			err := S.Push(child)
			if err != nil {
				panic(err)
			}
		}
	}

	return nil
}

// BFS traverses the tree in breadth-first order.
//
// Returns:
//   - error: An error if the traversal fails.
//
// Behavior:
//   - Use *ers.ErrNoError to stop the traversal without an error.
func (t *NoInfoTraversor[T]) BFS() error {
	if t.node == nil || t.observe == nil {
		return nil
	}

	Q := Queuer.NewLinkedQueue(t.node)

	for {
		first, err := Q.Dequeue()
		if err != nil {
			break
		}

		ok, err := t.observe(first.Data)
		if err != nil {
			if ers.As[*ers.ErrNoError](err) {
				return nil
			} else {
				return err
			}
		}

		if !ok {
			continue
		}

		for _, child := range first.children {
			err := Q.Enqueue(child)
			if err != nil {
				panic(err)
			}
		}
	}

	return nil
}

// Traverse creates a new traversor for the tree.
//
// Parameters:
//   - tree: The tree to traverse.
//   - init: The initial info.
//   - f: The observer function.
//
// Returns:
//   - Traversor[T, I]: The traversor.
func NoInfoTraverse[T any](tree *Tree[T], f NoInfoObserverFunc[T]) *NoInfoTraversor[T] {
	var root *TreeNode[T]

	if tree == nil {
		root = nil
	} else {
		root = tree.root
	}

	return &NoInfoTraversor[T]{
		node:    root,
		observe: f,
	}
}

// NextsFunc is a function that returns the next elements.
//
// Parameters:
//   - elem: The current element.
//   - info: The info of the current element.
//
// Returns:
//   - []T: The next elements.
//   - error: An error if the next elements cannot be found.
type NoInfoNextsFunc[T any] func(elem T) ([]T, error)

// MakeTree creates a tree from the given element.
//
// Parameters:
//   - elem: The element to start the tree from.
//
// Returns:
//   - *Tree[T]: The tree created from the element.
//   - error: An error if the tree cannot be created.
func NoInfoMakeTree[T any](elem T, f NoInfoNextsFunc[T]) (*Tree[T], error) {
	// 1. Handle the first element
	nexts, err := f(elem)
	if err != nil {
		return nil, err
	}

	tree := NewTree(elem)

	if len(nexts) == 0 {
		return tree, nil
	}

	// 2. Create a stack and push the first element
	type StackElement struct {
		Prev *TreeNode[T]
		Elem *TreeNode[T]
	}

	S := Stacker.NewLinkedStack[StackElement]()

	for _, next := range nexts {
		err := S.Push(StackElement{Prev: tree.root, Elem: newTreeNode(next)})
		if err != nil {
			panic(err)
		}
	}

	for {
		top, err := S.Pop()
		if err != nil {
			break
		}

		nexts, err := f(top.Elem.Data)
		if err != nil {
			return nil, err
		}

		top.Prev.AddChildren(top.Elem)

		for _, next := range nexts {
			err := S.Push(StackElement{Prev: top.Elem, Elem: newTreeNode(next)})
			if err != nil {
				panic(err)
			}
		}
	}

	return tree, nil
}
