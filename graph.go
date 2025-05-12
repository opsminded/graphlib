package graphlib

import (
	"fmt"
	"sync"
	"time"
)

type Vertex struct {
	Label string
}

type Edge struct {
	Label       string
	Source      Vertex
	Destination Vertex
}

type Subgraph struct {
	Vertices []Vertex
	Edges    []Edge
}

type Graph struct {
	graph         *graph
	mu            sync.RWMutex
	lastID        uint32
	checkInterval time.Duration
}

func NewGraph() *Graph {
	return &Graph{
		graph:         new(),
		mu:            sync.RWMutex{},
		lastID:        0,
		checkInterval: 5 * time.Second,
	}
}

func (g *Graph) NewVertex(label string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.lastID++
	g.graph.newVertex(g.lastID, label)
}

func (g *Graph) NewEdge(label, a, b string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.lastID++

	aV, aOK := g.graph.vertices.find(a)
	bV, bOK := g.graph.vertices.find(b)
	if !aOK || !bOK {
		panic("vertex not found")
	}

	g.graph.newEdge(g.lastID, label, aV, bV)
}

func (g *Graph) GetVertexByLabel(label string) Vertex {
	g.mu.RLock()
	defer g.mu.RUnlock()
	v, ok := g.graph.vertices.find(label)
	if !ok {
		panic("vertex not found")
	}

	return Vertex{
		Label: v.label,
	}
}

func (g *Graph) GetVertexHealth(label string) bool {
	g.mu.RLock()
	defer g.mu.RUnlock()
	v, ok := g.graph.vertices.find(label)
	if !ok {
		panic("vertex not found")
	}
	return v.health
}

func (g *Graph) SetVertexHealth(label string, health bool) {
	g.mu.Lock()
	defer g.mu.Unlock()

	v, ok := g.graph.vertices.find(label)
	if !ok {
		panic("vertex not found")
	}

	v.health = health
}

func (g *Graph) GetEdgeByLabel(label string) Edge {
	g.mu.RLock()
	defer g.mu.RUnlock()

	for _, e := range g.graph.edges {
		for _, edge := range e {
			if edge.label == label {
				return Edge{
					Label:       edge.label,
					Source:      Vertex{Label: edge.source.label},
					Destination: Vertex{Label: edge.destination.label},
				}
			}
		}
	}

	panic("edge not found")
}

func (g *Graph) VertexLen() int {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.graph.vertices.len()
}

func (g *Graph) EdgeLen() int {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.graph.edges.len()
}

func (g *Graph) Neighbors(label string) Subgraph {
	g.mu.RLock()
	defer g.mu.RUnlock()

	vertices := []Vertex{}
	edges := []Edge{}

	vSource, ok := g.graph.vertices.find(label)
	if !ok {
		panic("vertex not found")
	}

	source := Vertex{Label: vSource.label}

	for _, dep := range vSource.dependencies {
		vDestination, _ := g.graph.vertices.find(dep.label)
		vEdge, _ := g.graph.edges.find(vSource, vDestination)
		destination := Vertex{Label: vDestination.label}
		vertices = append(vertices, destination)
		new := Edge{
			Label:       vEdge.label,
			Source:      source,
			Destination: destination,
		}
		edges = append(edges, new)
	}

	for _, dep := range vSource.dependents {
		vDestination, _ := g.graph.vertices.find(dep.label)
		vEdge, _ := g.graph.edges.find(vDestination, vSource)
		destination := Vertex{Label: vDestination.label}
		vertices = append(vertices, destination)
		new := Edge{
			Label:       vEdge.label,
			Source:      destination,
			Destination: source,
		}
		edges = append(edges, new)
	}

	// Sempre incluir o próprio vértice
	vertices = append(vertices, Vertex{Label: vSource.label})

	return Subgraph{
		Vertices: vertices,
		Edges:    edges,
	}
}

func (g *Graph) GetVertexDependents(label string, all bool) Subgraph {
	g.mu.RLock()
	defer g.mu.RUnlock()

	vertices := map[string]Vertex{}
	edges := map[string]Edge{}

	v, ok := g.graph.vertices.find(label)
	if !ok {
		panic("vertex not found")
	}

	source := Vertex{Label: v.label}

	// Para dependentes diretos
	for _, dep := range v.dependents {
		destination := Vertex{Label: dep.label}
		vertices[destination.Label] = destination
		re, _ := g.graph.edges.find(dep, v)
		new := Edge{
			Label:       re.label,
			Source:      destination,
			Destination: source,
		}
		edges[new.Label] = new
	}

	// Se for para incluir todos os descendentes, recorremos recursivamente
	if all {
		var dfs func(v *vertex)
		dfs = func(v *vertex) {
			for _, dep := range v.dependents {
				destination := Vertex{Label: dep.label}
				vertices[destination.Label] = destination
				re, _ := g.graph.edges.find(dep, v)
				new := Edge{
					Label:       re.label,
					Source:      destination,
					Destination: Vertex{Label: v.label},
				}
				edges[new.Label] = new
				dfs(dep)
			}
		}
		dfs(v)
	}

	// Sempre incluir o próprio vértice
	vertices[v.label] = Vertex{Label: v.label}

	sub := Subgraph{
		Vertices: []Vertex{},
		Edges:    []Edge{},
	}

	for _, v := range vertices {
		sub.Vertices = append(sub.Vertices, v)
	}
	for _, e := range edges {
		sub.Edges = append(sub.Edges, e)
	}
	return sub
}

func (g *Graph) GetVertexDependencies(label string, all bool) Subgraph {
	g.mu.RLock()
	defer g.mu.RUnlock()

	vertices := map[string]Vertex{}
	edges := map[string]Edge{}

	v, ok := g.graph.vertices.find(label)
	if !ok {
		panic("vertex not found")
	}

	source := Vertex{Label: v.label}

	// Para dependências diretas
	for _, dep := range v.dependencies {
		destination := Vertex{Label: dep.label}
		vertices[destination.Label] = destination
		re, _ := g.graph.edges.find(v, dep)
		new := Edge{
			Label:       re.label,
			Source:      source,
			Destination: destination,
		}
		edges[new.Label] = new
	}

	// Se for para incluir todas as dependências, recorremos recursivamente
	if all {
		var dfs func(v *vertex)
		dfs = func(v *vertex) {
			for _, dep := range v.dependencies {
				destination := Vertex{Label: dep.label}
				vertices[destination.Label] = destination
				re, _ := g.graph.edges.find(v, dep)
				new := Edge{
					Label:       re.label,
					Source:      Vertex{Label: v.label},
					Destination: destination,
				}
				edges[new.Label] = new
				dfs(dep)
			}
		}
		dfs(v)
	}

	// Sempre incluir o próprio vértice
	vertices[v.label] = Vertex{Label: v.label}

	sub := Subgraph{
		Vertices: []Vertex{},
		Edges:    []Edge{},
	}

	for _, v := range vertices {
		sub.Vertices = append(sub.Vertices, v)
	}
	for _, e := range edges {
		sub.Edges = append(sub.Edges, e)
	}

	return sub
}

func (g *Graph) Path(from, to string) Subgraph {
	g.mu.RLock()
	defer g.mu.RUnlock()

	start, ok1 := g.graph.vertices.find(from)
	end, ok2 := g.graph.vertices.find(to)

	if !ok1 || !ok2 {
		panic("vertex not found")
	}

	vertexSet := make(map[uint32]*vertex)
	edgeSet := make(map[string]Edge)

	var dfs func(v *vertex, path []*vertex)

	dfs = func(v *vertex, path []*vertex) {
		path = append(path, v)

		if v.id == end.id {
			// Caminho completo encontrado, registra os vértices e arestas
			for i := 0; i < len(path)-1; i++ {
				from := path[i]
				to := path[i+1]
				vertexSet[from.id] = from
				vertexSet[to.id] = to
				if e, ok := g.graph.edges.find(from, to); ok {
					key := fmt.Sprintf("%d->%d", from.id, to.id)
					edgeSet[key] = Edge{
						Label:       e.label,
						Source:      Vertex{Label: from.label},
						Destination: Vertex{Label: to.label},
					}
				}
			}
			return
		}

		for _, dep := range v.dependencies {
			alreadyVisited := false
			for _, p := range path {
				if p.id == dep.id {
					alreadyVisited = true
					break
				}
			}
			if !alreadyVisited {
				dfs(dep, path)
			}
		}
	}

	dfs(start, []*vertex{})

	// Converter mapas para slices
	var vertices []Vertex
	for _, v := range vertexSet {
		vertices = append(vertices, Vertex{Label: v.label})
	}

	var edges []Edge
	for _, e := range edgeSet {
		edges = append(edges, e)
	}

	return Subgraph{
		Vertices: vertices,
		Edges:    edges,
	}
}
