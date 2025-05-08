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
