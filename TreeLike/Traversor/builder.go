package Traversor

import (
	"github.com/PlayerR9/MyGoLib/ListLike/Stacker"
	uc "github.com/PlayerR9/MyGoLib/Units/Common"

	tr "github.com/PlayerR9/MyGoLib/TreeLike/Tree"
)

// Builder is a struct that builds a tree.
type Builder[T any] struct {
	// info is the info of the builder.
	info uc.Objecter

	// f is the next function.
	f NextsFunc[T]
}

// SetInfo sets the info of the builder.
//
// Parameters:
//   - info: The info to set.
func (b *Builder[T]) SetInfo(info uc.Objecter) {
	b.info = info
}

// NextsFunc is a function that returns the next elements.
//
// Parameters:
//   - elem: The element to get the next elements from.
//   - info: The info of the element.
//
// Returns:
//   - []T: A slice of the next elements.
//   - error: An error if the function fails.
type NextsFunc[T any] func(elem T, info uc.Objecter) ([]T, error)

// SetNextFunc sets the next function of the builder.
//
// Parameters:
//   - f: The function to set.
func (b *Builder[T]) SetNextFunc(f NextsFunc[T]) {
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
//   - *Tree[T]: The tree created from the element.
//   - error: An error if the tree cannot be created.
//
// Behaviors:
//   - The 'info' parameter is copied for each node and it specifies the initial info
//     before traversing the tree.
func (b *Builder[T]) Build(elem T) (*tr.Tree[T], error) {
	if b.f == nil {
		return nil, nil
	}

	// 1. Initialize the tree
	tree := tr.NewTree(elem)

	S := Stacker.NewLinkedStack(
		newStackElement[T](nil, elem, b.info, b.f),
	)

	for {
		top, err := S.Pop()
		if err != nil {
			break
		}

		nexts, err := top.apply()
		if err != nil {
			return nil, err
		}

		if len(nexts) == 0 {
			continue
		}

		if top.prev != nil {
			top.prev.AddChildren(top.elem)

			// Update the leaves
			tree.UpdateLeaves()
		}

		for _, next := range nexts {
			var se *stackElement[T]

			if b.info == nil {
				se = newStackElement(top.elem, next, nil, b.f)
			} else {
				se = newStackElement(top.elem, next, top.info, b.f)
			}

			err := S.Push(se)
			if err != nil {
				panic(err)
			}
		}
	}

	b.Reset()

	return tree, nil
}

// Reset resets the builder.
func (b *Builder[T]) Reset() {
	b.info = nil
	b.f = nil
}
