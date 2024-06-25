package Tree

import (
	uc "github.com/PlayerR9/MyGoLib/Units/common"
)

// StatusInfo is a generic data structure that represents a status node in a tree.
type StatusInfo[S uc.Enumer, T any] struct {
	// data is the data of the status node.
	data T

	// status is the status of the status node.
	status S
}

// NewStatusInfo creates a new status node.
//
// Parameters:
//   - data: The data of the status node.
//   - status: The status of the status node.
//
// Returns:
//   - *StatusInfo[S, T]: The new status node.
func NewStatusInfo[S uc.Enumer, T any](data T, status S) *StatusInfo[S, T] {
	si := &StatusInfo[S, T]{
		data:   data,
		status: status,
	}

	return si
}

// SetData sets the data of the tree node.
//
// Parameters:
//   - data: The data to set.
func (si *StatusInfo[S, T]) SetData(data T) {
	si.data = data
}

// ChangeStatus sets the status of the tree node.
//
// Parameters:
//   - status: The status to set.
func (si *StatusInfo[S, T]) ChangeStatus(status S) {
	si.status = status
}

// GetStatus returns the status of the tree node.
//
// Returns:
//   - S: The status of the tree node.
func (si *StatusInfo[S, T]) GetStatus() S {
	return si.status
}

// GetData returns the data of the tree node.
//
// Returns:
//   - T: The data of the tree node.
func (si *StatusInfo[S, T]) GetData() T {
	return si.data
}
