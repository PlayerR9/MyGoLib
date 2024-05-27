package Graph

import (
	uc "github.com/PlayerR9/MyGoLib/Units/Common"
	uos "github.com/PlayerR9/MyGoLib/Utility/Sorting"
)

// remapVertices shifts all vertices to the right starting from the given position.
//
// Parameters:
//   - pos: The position to start shifting vertices.
func (g *Graph[T]) remapVertices(pos int) {
	for i := len(g.edges) - 1; i >= pos; i-- {
		to := g.edges[i]

		// Shift all elements to the right
		for j := len(to) - 1; j >= pos; j-- {
			to[j+1] = to[j]
		}

		g.edges[i+1] = to
	}
}

type Graph[T uc.Comparer] struct {
	vertices *uos.Slice[*Vertex[T]]
	edges    map[int]map[int]bool
}

func NewGraph[T uc.Comparer]() *Graph[T] {
	return &Graph[T]{
		vertices: uos.NewSlice[*Vertex[T]](nil),
		edges:    make(map[int]map[int]bool),
	}
}

func (g *Graph[T]) AddVertex(v *Vertex[T]) {
	if v == nil {
		return
	}

	pos, ok := g.vertices.TryInsert(v)
	if !ok {
		g.vertices.Insert(v, false)
		g.remapVertices(pos)
	}

	g.edges[pos] = make(map[int]bool)
}

func (g *Graph[T]) AddEdge(from, to *Vertex[T]) {
	if to == nil {
		if from == nil {
			return
		}

		from.isFinal = true
	} else if from == nil {
		to.isInitial = true
	}

	g.AddVertex(from)
	g.AddVertex(to)

	if from == nil {
		return
	} else if to == nil {
		return
	}

	fromIndex := g.vertices.Find(from)
	toIndex := g.vertices.Find(to)

	g.edges[fromIndex][toIndex] = true
}
