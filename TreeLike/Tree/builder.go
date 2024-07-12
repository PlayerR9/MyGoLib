package Tree

import (
	lls "github.com/PlayerR9/MyGoLib/ListLike/Stacker"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
)

// Infoer is an interface that provides the info of the element.
type Infoer interface {
	uc.Copier
}

// NextsFunc is a function that returns the next elements.
//
// Parameters:
//   - elem: The element to get the next elements from.
//   - info: The info of the element.
//
// Returns:
//   - []Noder: The next elements.
//   - error: An error if the function fails.
type NextsFunc func(elem Noder, info Infoer) ([]Noder, error)

// Builder is a struct that builds a tree.
type Builder struct {
	// info is the info of the builder.
	info Infoer

	// f is the next function.
	f NextsFunc
}

// SetInfo sets the info of the builder.
//
// Parameters:
//   - info: The info to set.
func (b *Builder) SetInfo(info Infoer) {
	b.info = info
}

// SetNextFunc sets the next function of the builder.
//
// Parameters:
//   - f: The function to set.
func (b *Builder) SetNextFunc(f NextsFunc) {
	b.f = f
}

// MakeTree creates a tree from the given element.
//
// Parameters:
//   - elem: The element to start the tree from.
//   - info: The info of the element.
//   - f: The function that, given an element and info, returns the next elements.
//     (i.e., the children of the element).
//
// Returns:
//   - *Tree: The tree created from the element.
//   - error: An error if the next function fails.
//
// Behaviors:
//   - The 'info' parameter is copied for each node and it specifies the initial info
//     before traversing the tree.
func (b *Builder) Build(elem Noder) (*Tree, error) {
	if b.f == nil {
		return nil, nil
	}

	// 1. Handle the root node
	nexts, err := b.f(elem, b.info)
	if err != nil {
		return nil, err
	}

	tree := NewTree(elem)

	if len(nexts) == 0 {
		return tree, nil
	}

	S := lls.NewLinkedStack[*stackElement]()

	for _, next := range nexts {
		root := tree.Root()

		se := new_stack_element(root, next, b.info)

		S.Push(se)
	}

	for {
		top, ok := S.Pop()
		if !ok {
			break
		}

		data, ok := top.get_data()
		uc.Assert(ok, "Missing data")

		top_inf := top.get_info()

		nexts, err := b.f(data, top_inf)
		if err != nil {
			return nil, err
		}

		ok = top.link_to_prev()
		uc.Assert(ok, "Cannot link to previous node")

		if len(nexts) == 0 {
			continue
		}

		top_elem := top.get_elem()

		for _, next := range nexts {
			se := new_stack_element(top_elem, next, top_inf)

			S.Push(se)
		}
	}

	b.Reset()

	tree.RegenerateLeaves()

	return tree, nil
}

// Reset resets the builder.
func (b *Builder) Reset() {
	b.info = nil
	b.f = nil
}