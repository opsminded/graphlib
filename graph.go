package graphlib

import (
	"fmt"
	"slices"
	"sync"
	"time"
)

type Vertex struct {
	Label   string
	Healthy bool
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
	graph             *graph
	mu                sync.RWMutex
	lastID            uint32
	checkInterval     time.Duration
	unhealthyVertices []*vertex
}

func NewGraph() *Graph {
	return &Graph{
		graph:             new(),
		mu:                sync.RWMutex{},
		lastID:            0,
		checkInterval:     5 * time.Second,
		unhealthyVertices: []*vertex{},
	}
}

func (g *Graph) NewVertex(label string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.lastID++
	g.graph.newVertex(g.lastID, label)
}

func (g *Graph) NewEdge(label, a, b string) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.lastID++

	aV, err := g.graph.vertices.find(a)
	if err != nil {
		return err
	}
	bV, err := g.graph.vertices.find(b)
	if err != nil {
		return err
	}

	_, err = g.graph.newEdge(g.lastID, label, aV, bV)
	return err
}

func (g *Graph) GetVertexByLabel(label string) (Vertex, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	v, err := g.graph.vertices.find(label)
	if err != nil {
		return Vertex{}, err
	}

	return Vertex{
		Label:   v.label,
		Healthy: v.healthy,
	}, nil
}

func (g *Graph) SetVertexHealth(label string, health bool) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	v, err := g.graph.vertices.find(label)
	if err != nil {
		return err
	}

	return g.setVertexHealth(v, health)
}

func (g *Graph) ClearGraphHealthyStatus() {
	g.mu.Lock()
	defer g.mu.Unlock()

	for _, v := range g.graph.vertices {
		v.healthy = true
	}
	g.unhealthyVertices = []*vertex{}
}

func (g *Graph) setVertexHealth(v *vertex, health bool) error {
	v.healthy = health
	if v.healthy {
		g.unhealthyVertices = slices.DeleteFunc(g.unhealthyVertices, func(x *vertex) bool {
			return x.label == v.label
		})
	} else {
		g.unhealthyVertices = append(g.unhealthyVertices, v)
		for _, dep := range v.dependents {
			if dep.healthy {
				g.setVertexHealth(dep, false)
			}
		}
	}
	return nil
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

func (g *Graph) Neighbors(label string) (Subgraph, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	vertices := []Vertex{}
	edges := []Edge{}

	vSource, err := g.graph.vertices.find(label)
	if err != nil {
		return Subgraph{}, err
	}

	source := Vertex{
		Label:   vSource.label,
		Healthy: vSource.healthy,
	}

	for _, dep := range vSource.dependencies {
		vDestination, _ := g.graph.vertices.find(dep.label)
		vEdge, _ := g.graph.edges.find(vSource, vDestination)
		destination := Vertex{
			Label:   vDestination.label,
			Healthy: vDestination.healthy,
		}
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
		destination := Vertex{
			Label:   vDestination.label,
			Healthy: vDestination.healthy,
		}
		vertices = append(vertices, destination)
		new := Edge{
			Label:       vEdge.label,
			Source:      destination,
			Destination: source,
		}
		edges = append(edges, new)
	}

	// Sempre incluir o próprio vértice
	vertices = append(vertices, Vertex{
		Label:   vSource.label,
		Healthy: vSource.healthy,
	})

	return Subgraph{
		Vertices: vertices,
		Edges:    edges,
	}, nil
}

func (g *Graph) GetVertexDependents(label string, all bool) (Subgraph, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	vertices := map[string]Vertex{}
	edges := map[string]Edge{}

	v, err := g.graph.vertices.find(label)
	if err != nil {
		return Subgraph{}, err
	}

	source := Vertex{
		Label:   v.label,
		Healthy: v.healthy,
	}

	// Para dependentes diretos
	for _, dep := range v.dependents {
		destination := Vertex{
			Label:   dep.label,
			Healthy: dep.healthy,
		}
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
				destination := Vertex{
					Label:   dep.label,
					Healthy: dep.healthy,
				}
				vertices[destination.Label] = destination
				re, _ := g.graph.edges.find(dep, v)
				new := Edge{
					Label:  re.label,
					Source: destination,
					Destination: Vertex{
						Label:   v.label,
						Healthy: v.healthy,
					},
				}
				edges[new.Label] = new
				dfs(dep)
			}
		}
		dfs(v)
	}

	// Sempre incluir o próprio vértice
	vertices[v.label] = Vertex{
		Label:   v.label,
		Healthy: v.healthy,
	}

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
	return sub, nil
}

func (g *Graph) GetVertexDependencies(label string, all bool) (Subgraph, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	vertices := map[string]Vertex{}
	edges := map[string]Edge{}

	v, err := g.graph.vertices.find(label)
	if err != nil {
		return Subgraph{}, err
	}

	source := Vertex{
		Label:   v.label,
		Healthy: v.healthy,
	}

	// Para dependências diretas
	for _, dep := range v.dependencies {
		destination := Vertex{
			Label:   dep.label,
			Healthy: dep.healthy,
		}
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
				destination := Vertex{
					Label:   dep.label,
					Healthy: dep.healthy,
				}
				vertices[destination.Label] = destination
				re, _ := g.graph.edges.find(v, dep)
				new := Edge{
					Label: re.label,
					Source: Vertex{
						Label:   v.label,
						Healthy: v.healthy,
					},
					Destination: destination,
				}
				edges[new.Label] = new
				dfs(dep)
			}
		}
		dfs(v)
	}

	// Sempre incluir o próprio vértice
	vertices[v.label] = Vertex{
		Label:   v.label,
		Healthy: v.healthy,
	}

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

	return sub, nil
}

func (g *Graph) Path(from, to string) (Subgraph, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	start, err := g.graph.vertices.find(from)
	if err != nil {
		return Subgraph{}, err
	}

	end, err := g.graph.vertices.find(to)
	if err != nil {
		return Subgraph{}, err
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
				if e, err := g.graph.edges.find(from, to); err == nil {
					key := fmt.Sprintf("%d->%d", from.id, to.id)
					edgeSet[key] = Edge{
						Label: e.label,
						Source: Vertex{
							Label:   from.label,
							Healthy: from.healthy,
						},
						Destination: Vertex{
							Label:   to.label,
							Healthy: to.healthy,
						},
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
		vertices = append(vertices, Vertex{
			Label:   v.label,
			Healthy: v.healthy,
		})
	}

	var edges []Edge
	for _, e := range edgeSet {
		edges = append(edges, e)
	}

	return Subgraph{
		Vertices: vertices,
		Edges:    edges,
	}, nil
}

func (g *Graph) UnhealthyVertices() []Vertex {
	g.mu.RLock()
	defer g.mu.RUnlock()

	unhealthy := []Vertex{}

	for _, v := range g.unhealthyVertices {
		if !v.healthy {
			if v.dependencies.len() == 0 || v.dependents.len() == 0 {
				unhealthy = append(unhealthy, Vertex{
					Label:   v.label,
					Healthy: v.healthy,
				})
			}
		}
	}
	return unhealthy
}

func (g *Graph) Lineage(from string) (Subgraph, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	start, err := g.graph.vertices.find(from)
	if err != nil {
		return Subgraph{}, err
	}

	// 1) cache de profundidades
	depthMemo := map[uint32]int{}
	var depth func(v *vertex) int
	depth = func(v *vertex) int {
		if d, ok := depthMemo[v.id]; ok {
			return d
		}
		if len(v.dependencies) == 0 {
			depthMemo[v.id] = 0
			return 0
		}
		maxD := 0
		for _, dep := range v.dependencies {
			if d := depth(dep) + 1; d > maxD {
				maxD = d
			}
		}
		depthMemo[v.id] = maxD
		return maxD
	}

	// 2) caminhada única
	verts := []Vertex{{Label: start.label}}
	edges := []Edge{}

	curr := start
	for {
		if len(curr.dependencies) == 0 {
			break
		}
		// escolhe dependência(s) mais profundas
		bestDepth := -1
		var best *vertex
		tie := false
		for _, dep := range curr.dependencies {
			if d := depth(dep); d > bestDepth {
				bestDepth, best, tie = d, dep, false
			} else if d == bestDepth {
				tie = true
			}
		}
		// empate?  pára aqui
		if tie {
			break
		}
		// inclui aresta e prossegue
		edgeCore, err := g.graph.edges.find(curr, best)
		if err != nil {
			return Subgraph{}, err
		}

		edges = append(edges, Edge{
			Label: edgeCore.label,
			Source: Vertex{
				Label:   curr.label,
				Healthy: curr.healthy,
			},
			Destination: Vertex{
				Label:   best.label,
				Healthy: best.healthy,
			},
		})
		verts = append(verts, Vertex{
			Label:   best.label,
			Healthy: best.healthy,
		})
		curr = best
	}

	return Subgraph{
		Vertices: verts,
		Edges:    edges,
	}, nil
}
