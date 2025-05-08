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

// Iterator define uma interface opcional para percorrer vértices do grafo.
type Iterator interface {
	HasNext() bool
	Next() *vertex
	Iterate(func(v *vertex) error) error
	Reset()
}

type vertex struct {
	id
	class
	label     string
	health    bool
	neighbors []*vertex
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

// NewVertex cria e registra um novo vértice com o label e classe fornecidos.
// Retorna erro se o label ou classe forem inválidos.
func (g *Graph) NewVertex(label, cla string) (*vertex, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if !validLabel.MatchString(label) {
		return nil, ErrInvalidLabel
	}

	if !validClassName.MatchString(cla) {
		return nil, ErrInvalidClassName
	}

	return g.newVertex(label, cla), nil
}

// NewEdge cria e registra uma nova aresta entre dois vértices existentes.
// Garante que não haja múltiplas arestas entre o mesmo par de vértices.
func (g *Graph) NewEdge(label string, cla string, source, destination *vertex) (*edge, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if source == nil || destination == nil {
		return nil, ErrNilVertices
	}

	if !validLabel.MatchString(label) {
		return nil, ErrInvalidLabel
	}

	if !validClassName.MatchString(cla) {
		return nil, ErrInvalidClassName
	}

	return g.newEdge(label, cla, source, destination), nil
}

func (g *Graph) GetVertexByLabel(label string) *vertex {
	for _, v := range g.vertices {
		if v.label == label {
			return v
		}
	}
	return nil
}

func (g *Graph) newVertex(label, cla string) *vertex {
	class := g.ensureVertexClass(cla)

	for _, v := range g.vertices {
		if v.label == label {
			return v
		}
	}

	nid := id(g.newID())

	v := &vertex{
		id:        nid,
		class:     class,
		label:     label,
		health:    false,
		neighbors: []*vertex{},
	}

	g.vertices[v.id] = v
	return v
}

func (g *Graph) newEdge(label string, cla string, source, destination *vertex) *edge {

	eclass := g.ensureEdgeClass(cla)

	// prevent edge-multiplicity
	if destMap, ok := g.edges[source.id]; ok {
		if e, ok := destMap[destination.id]; ok {
			return e
		}
	}

	var e *edge

	{
		eid := id(g.newID())

		source.neighbors = append(source.neighbors, destination)

		e = &edge{
			id:    eid,
			class: eclass,
			label: label,

			source:      source,
			destination: destination,
		}

		if _, ok := g.edges[source.id]; !ok {
			g.edges[source.id] = map[id]*edge{destination.id: e}
			return e
		}

		g.edges[source.id][destination.id] = e
	}

	return e
}

func (g *Graph) newID() uint32 {
	g.lastID++
	return g.lastID
}

func (g *Graph) ensureVertexClass(name string) class {
	for classID, cla := range g.vertexClasses {
		if cla == name {
			return classID
		}
	}

	classID := class(g.newID())
	g.vertexClasses[classID] = name
	return classID
}

func (g *Graph) ensureEdgeClass(name string) class {
	for classID, cla := range g.edgeClasses {
		if cla == name {
			return classID
		}
	}

	classID := class(g.newID())
	g.edgeClasses[classID] = name
	return classID
}
