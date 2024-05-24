package WeightedGraph

import (
	tlt "github.com/PlayerR9/MyGoLib/TreeLike/Traversor"
	tr "github.com/PlayerR9/MyGoLib/TreeLike/Tree"
	uc "github.com/PlayerR9/MyGoLib/Units/Common"
)

// WeightFunc is a function that calculates the weight of an edge.
//
// Parameters:
//   - from: the source vertex.
//   - to: the destination vertex.
//
// Returns:
//   - float64: the weight of the edge.
//   - bool: true if the edge exists, otherwise false.
type WeightFunc[T uc.Objecter] func(from, to T) (float64, bool)

// WeightedGraph represents a graph.
type WeightedGraph[T uc.Objecter] struct {
	// Vertices in the graph.
	Vertices []T

	// Edges in the graph.
	Edges [][]*float64
}

// NewWeightedGraph creates a new graph with the given vertices.
//
// Parameters:
//   - vertices: vertices in the graph.
//
// Returns:
//   - *WeightedGraph: the new graph.
func NewWeightedGraph[T uc.Objecter](vertices []T, f WeightFunc[T]) *WeightedGraph[T] {
	if len(vertices) == 0 {
		return &WeightedGraph[T]{
			Vertices: make([]T, 0),
			Edges:    make([][]*float64, 0),
		}
	}

	g := &WeightedGraph[T]{
		Vertices: vertices,
		Edges:    make([][]*float64, 0, len(vertices)),
	}

	for _, from := range vertices {
		edge := make([]*float64, 0, len(vertices))

		for _, to := range vertices {
			w, ok := f(from, to)
			if !ok {
				edge = append(edge, nil)
			} else {
				edge = append(edge, &w)
			}
		}

		g.Edges = append(g.Edges, edge)
	}

	return g
}

// IndexOf returns the index of the given element in the graph.
//
// Parameters:
//   - elem: the element to find.
//
// Returns:
//   - int: the index of the element, or -1 if not found.
func (g *WeightedGraph[T]) IndexOf(elem T) int {
	for i, x := range g.Vertices {
		if x.Equals(elem) {
			return i
		}
	}

	return -1
}

// AdjacentOf returns the adjacent vertices of the given vertex.
//
// Parameters:
//   - from: the source vertex.
//
// Returns:
//   - []T: the adjacent vertices.
func (g *WeightedGraph[T]) AdjacentOf(from T) []T {
	index := g.IndexOf(from)
	if index == -1 {
		return nil
	}

	adj := make([]T, 0)

	for j, distance := range g.Edges[index] {
		if distance != nil {
			adj = append(adj, g.Vertices[j])
		}
	}

	return adj
}

// MakeTree creates a tree of the graph with the given root.
//
// Parameters:
//   - root: the root of the tree.
//   - f: the nexts function.
//
// Returns:
//   - *WeightedGraphTree: the tree of the graph.
//   - error: an error if the tree creation fails.
func (g *WeightedGraph[T]) MakeTree(root T, info uc.Objecter, f tlt.NextsFunc[T]) (*tr.Tree[T], error) {
	var builder tlt.Builder[T]

	builder.SetInfo(info)
	builder.SetNextFunc(f)

	return builder.Build(root)
}
