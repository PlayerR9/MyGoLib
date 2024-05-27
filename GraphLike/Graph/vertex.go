package Graph

import (
	uc "github.com/PlayerR9/MyGoLib/Units/Common"
)

type Vertex[T uc.Comparer] struct {
	value T

	isInitial bool
	isFinal   bool
}

// String implements Common.Objecter.
func (v *Vertex[T]) String() string {
	panic("unimplemented")
}

func (v *Vertex[T]) Equals(other uc.Objecter) bool {
	if other == nil {
		return false
	}

	otherV, ok := other.(*Vertex[T])
	if !ok {
		return false
	}

	return v.value.Equals(otherV.value) &&
		v.isInitial == otherV.isInitial &&
		v.isFinal == otherV.isFinal
}

func (v *Vertex[T]) Compare(other uc.Comparer) int {
	if other == nil {
		return 1
	}

	otherV, ok := other.(*Vertex[T])
	if !ok {
		return 1
	}

	return v.value.Compare(otherV.value)
}

func (v *Vertex[T]) Copy() uc.Objecter {
	return &Vertex[T]{
		value:     v.value.Copy().(T),
		isInitial: v.isInitial,
		isFinal:   v.isFinal,
	}
}

func NewVertex[T uc.Comparer](value T) *Vertex[T] {
	return &Vertex[T]{
		value:     value,
		isInitial: false,
		isFinal:   false,
	}
}
