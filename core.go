package graphlib

import (
	"time"
)

type edgesMap map[uint32]map[uint32]*edge

func (e edgesMap) add(i uint32, l string, a, b *vertex) *edge {

	// prevent nil vertices
	if a == nil || b == nil {
		panic("vertices are nil")
	}

	// prevent bidirectional edges
	if _, ok := e.find(b, a); ok {
		panic("bidirectional edges are not allowed")
	}

	// prevent cycles
	if e.wouldCreateCycle(a, b) {
		panic("adding this edge would create a cycle")
	}

	// prevent edge-multiplicity
	if x, ok := e.find(a, b); ok {
		return x
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

	return new
}

func (e edgesMap) find(a, b *vertex) (*edge, bool) {
	if destMap, ok := e[a.id]; ok {
		if e, ok := destMap[b.id]; ok {
			return e, true
		}
	}
	return nil, false
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

func (e *edge) Label() string {
	return e.label
}

type vertexList []*vertex

func (l *vertexList) add(i uint32, lbl string) *vertex {
	if v, ok := l.find(lbl); ok {
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

func (l *vertexList) find(lbl string) (*vertex, bool) {
	for _, v := range *l {
		if v.label == lbl {
			return v, true
		}
	}
	return nil, false
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

func (v *vertex) Label() string {
	return v.label
}

func (v *vertex) Health() bool {
	return v.health
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

func (g *graph) newEdge(i uint32, l string, a, b *vertex) *edge {
	return g.edges.add(i, l, a, b)
}
