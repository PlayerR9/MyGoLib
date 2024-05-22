package Graph

import (
	slext "github.com/PlayerR9/MyGoLib/Units/Slices"
)

type Builder struct {
	vertices []string
	edges    map[string][]string
}

func (b *Builder) AddVertex(v string) {
	if v == "" {
		return
	}

	b.vertices = append(b.vertices, v)
}

func (b *Builder) AddEdge(from, to string) {
	if from == "" || to == "" {
		return
	}

	if b.edges == nil {
		b.edges = make(map[string][]string)
	}

	b.edges[from] = append(b.edges[from], to)
}

func (b *Builder) Build() *Graph {
	g := &Graph{
		vertices: make([]string, 0),
		edges:    make(map[string][]string),
	}

	b.vertices = slext.RemoveDuplicates(b.vertices)

	verticesToRemove := make([]string, 0, len(b.edges))

	for from := range b.edges {
		verticesToRemove = append(verticesToRemove, from)
	}

	return &Graph{
		vertices: b.vertices,
		edges:    b.edges,
	}
}

func (b *Builder) Reset() {
	b.vertices = nil
	b.edges = nil
}
