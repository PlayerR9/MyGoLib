package Graph

import (
	uc "github.com/PlayerR9/MyGoLib/Units/Common"
)

type VertexElementer interface {
	uc.Comparer
	uc.Equaler
	uc.Copier
}

type Vertex[T VertexElementer] struct {
	value T

	isInitial bool
	isFinal   bool
}

func (v *Vertex[T]) Equals(other uc.Equaler) bool {
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

func (v *Vertex[T]) Compare(other uc.Comparer) (int, bool) {
	if other == nil {
		return 0, false
	}

	otherV, ok := other.(*Vertex[T])
	if !ok {
		return 0, false
	}

	return v.value.Compare(otherV.value)
}

func (v *Vertex[T]) Copy() uc.Copier {
	return &Vertex[T]{
		value:     v.value.Copy().(T),
		isInitial: v.isInitial,
		isFinal:   v.isFinal,
	}
}

func NewVertex[T VertexElementer](value T) *Vertex[T] {
	return &Vertex[T]{
		value:     value,
		isInitial: false,
		isFinal:   false,
	}
}
