package graphlib

import (
	"errors"
	"time"
)

type edgesMap map[uint32]map[uint32]*edge

func (e edgesMap) add(i uint32, l string, a, b *vertex) (*edge, error) {
	// prevent edge multiplicity
	if ex, err := e.find(a, b); err == nil {
		return ex, nil
	} else if !errors.As(err, &EdgeNotFoundError{}) {
		return nil, err
	}

	// prevent bidirectional edges
	if e.wouldCreateBidirectionalEdge(a, b) {
		return nil, BidirectionalEdgeError{
			LabelA: a.label,
			LabelB: b.label,
		}
	}

	// prevent cycles
	if e.wouldCreateCycle(a, b) {
		return nil, CycleError{
			LabelA: a.label,
			LabelB: b.label,
		}
	}

	a.dependencies = append(a.dependencies, b)
	b.dependents = append(b.dependents, a)

	new := &edge{
		id:          i,
		label:       l,
		source:      a,
		destination: b,
	}

	if _, ok := e[a.id]; !ok {
		e[a.id] = map[uint32]*edge{}
	}
	e[a.id][b.id] = new

	return new, nil
}

func (e edgesMap) find(a, b *vertex) (*edge, error) {
	// prevent nil vertices
	if a == nil {
		return nil, VertexNilError{position: "first"}
	}
	if b == nil {
		return nil, VertexNilError{position: "second"}
	}

	if destMap, ok := e[a.id]; ok {
		if e, ok := destMap[b.id]; ok {
			return e, nil
		}
	}
	return nil, EdgeNotFoundError{
		LabelA: a.label,
		LabelB: b.label,
	}
}

func (e edgesMap) exists(a, b *vertex) bool {
	if _, err := e.find(a, b); err == nil {
		return true
	}
	return false
}

func (e edgesMap) wouldCreateBidirectionalEdge(from, to *vertex) bool {
	// Swap the arguments to verify whether an edge already exists in the reverse direction.
	return e.exists(to, from)
}

func (e edgesMap) wouldCreateCycle(from, to *vertex) bool {

	var dfs func(v *vertex) bool

	visited := make(map[uint32]bool)

	dfs = func(v *vertex) bool {
		if v.id == from.id {
			return true
		}
		visited[v.id] = true
		for _, dep := range v.dependencies {
			if dfs(dep) {
				return true
			}
		}
		return false
	}

	return dfs(to)
}

func (e edgesMap) len() int {
	return len(e)
}

type edge struct {
	id          uint32
	label       string
	source      *vertex
	destination *vertex
}

type vertexList []*vertex

func (l *vertexList) add(i uint32, lbl string) *vertex {
	if v, err := l.find(lbl); err == nil {
		return v
	}

	new := &vertex{
		id:           i,
		label:        lbl,
		health:       true,
		dependents:   vertexList{},
		dependencies: vertexList{},
		lastCheck:    time.Now(),
	}

	*l = append(*l, new)
	return new
}

func (l *vertexList) find(lbl string) (*vertex, error) {
	for _, v := range *l {
		if v.label == lbl {
			return v, nil
		}
	}
	return nil, VertexNotFoundError{Label: lbl}
}

func (l *vertexList) len() int {
	return len(*l)
}

type vertex struct {
	id           uint32
	label        string
	health       bool
	dependents   vertexList
	dependencies vertexList
	lastCheck    time.Time
}

type graph struct {
	edges    edgesMap
	vertices vertexList
}

func new() *graph {
	return &graph{
		edges:    edgesMap{},
		vertices: vertexList{},
	}
}

func (g *graph) newVertex(i uint32, l string) *vertex {
	return g.vertices.add(i, l)
}

func (g *graph) newEdge(i uint32, l string, a, b *vertex) (*edge, error) {
	return g.edges.add(i, l, a, b)
}
