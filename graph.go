// Package graphlib fornece uma estrutura de grafo direcionado com suporte
// para vértices e arestas classificados, identificação única e verificação
// de validade baseada em rótulos.
package graphlib

import (
	"errors"
	"regexp"
	"sync"
)

type (
	id    uint32
	class uint32
)

var (
	validLabel     = regexp.MustCompile(`^[A-Z_0-9]+$`)
	validClassName = regexp.MustCompile(`^[A-Z_0-9]+$`)
)

var (
	ErrNilVertices      = errors.New("vertices are nil")
	ErrInvalidLabel     = errors.New("invalid label")
	ErrInvalidClassName = errors.New("invalid class name")
)

type vertex struct {
	id
	class
	label     string
	health    bool
	neighbors []*vertex
	lastCheck int64 // unix timestamp of last health check
}

func (v vertex) Label() string {
	return v.label
}

type edge struct {
	id
	class
	label       string
	source      *vertex
	destination *vertex
}

type resultSet struct {
	Principal *vertex
	All       bool

	Edges    []*edge
	Vertices []*vertex
}

func (e edge) Label() string {
	return e.label
}

// Graph representa a estrutura principal contendo vértices e arestas.
type Graph struct {
	vertices map[id]*vertex
	edges    map[id]map[id]*edge

	lastID        uint32
	edgeClasses   map[class]string
	vertexClasses map[class]string

	mu sync.Mutex
}

// NewGraph cria uma nova instância de grafo vazio.
func NewGraph() *Graph {
	return &Graph{
		vertices: make(map[id]*vertex),
		edges:    make(map[id]map[id]*edge),

		lastID:        0,
		edgeClasses:   make(map[class]string),
		vertexClasses: make(map[class]string),
	}
}

func (g *Graph) Path(label, destination string) resultSet {
	g.mu.Lock()
	defer g.mu.Unlock()

	var start, end *vertex
	for _, v := range g.vertices {
		if v.label == label {
			start = v
		}
		if v.label == destination {
			end = v
		}
	}

	if start == nil || end == nil {
		return resultSet{}
	}

	var allVertices []*vertex
	var allEdges []*edge
	visited := make(map[id]bool)
	pathStack := []*vertex{}

	var dfs func(current *vertex)
	dfs = func(current *vertex) {
		visited[current.id] = true
		pathStack = append(pathStack, current)

		if current == end {
			// Registrar caminho atual
			for _, v := range pathStack {
				if !containsVertex(allVertices, v) {
					allVertices = append(allVertices, v)
				}
			}
			for i := 0; i < len(pathStack)-1; i++ {
				edgesFrom := g.edges[pathStack[i].id]
				if edge, ok := edgesFrom[pathStack[i+1].id]; ok && !containsEdge(allEdges, edge) {
					allEdges = append(allEdges, edge)
				}
			}
		} else {
			for _, neighbor := range current.neighbors {
				if !visited[neighbor.id] {
					dfs(neighbor)
				}
			}
		}

		// Backtrack
		visited[current.id] = false
		pathStack = pathStack[:len(pathStack)-1]
	}

	dfs(start)

	return resultSet{
		Principal: start,
		Vertices:  allVertices,
		Edges:     allEdges,
	}
}

func containsVertex(list []*vertex, v *vertex) bool {
	for _, item := range list {
		if item.id == v.id {
			return true
		}
	}
	return false
}

func containsEdge(list []*edge, e *edge) bool {
	for _, item := range list {
		if item.id == e.id {
			return true
		}
	}
	return false
}
