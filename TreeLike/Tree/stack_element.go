package Tree

import (
	uc "github.com/PlayerR9/MyGoLib/Units/common"
)

// stackElement is a stack element.
type stackElement struct {
	// prev is the previous node.
	prev Noder

	// elem is the current node.
	elem Noder

	// info is the info of the current node.
	info uc.Copier
}

// new_stack_element creates a new stack element.
//
// Parameters:
//   - prev: The previous node.
//   - elem: The helper.
//
// Returns:
//   - *stackElement: A pointer to the stack element.
func new_stack_element(prev, data Noder, info uc.Copier) *stackElement {
	tree := NewTreeNode(data)

	se := &stackElement{
		prev: prev,
		elem: tree,
		info: info,
	}

	return se
}

// get_data returns the data of the stack element.
//
// Returns:
//   - Tree.Noder: The data of the stack element.
//   - bool: True if the data is valid, otherwise false.
func (se *stackElement) get_data() (Noder, bool) {
	if se.elem == nil {
		return nil, false
	}

	return se.elem, true
}

// get_info returns the info of the stack element.
//
// Returns:
//   - common.Copier: The info of the stack element.
func (se *stackElement) get_info() uc.Copier {
	return se.info
}

// link_to_prev links the current node to the previous node.
//
// Returns:
//   - bool: True if the link is successful, otherwise false.
func (se *stackElement) link_to_prev() bool {
	if se.prev == nil {
		return false
	}

	se.prev.AddChild(se.elem)

	return true
}

// get_elem returns the current node.
//
// Returns:
//   - Tree.Noder: The current node.
func (se *stackElement) get_elem() Noder {
	return se.elem
}
