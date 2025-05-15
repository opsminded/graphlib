package graphlib

import (
	"context"
	"sync"
	"time"

	"github.com/opsminded/graphlib/internal/core"
)

type Edge = core.Edge
type Vertex = core.Vertex
type Subgraph = core.Subgraph

type Graph struct {
	graph             *core.Graph
	mu                sync.RWMutex
	checkInterval     time.Duration
	unhealthyVertices []*Vertex
}

func NewGraph(ctx context.Context) *Graph {
	return &Graph{
		graph:             core.NewSoAGraph(),
		mu:                sync.RWMutex{},
		checkInterval:     5 * time.Second,
		unhealthyVertices: []*Vertex{},
	}
}

func (g *Graph) NewVertex(key, label string, healthy bool) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.graph.AddVertex(key, label, healthy)
}

func (g *Graph) NewEdge(src, tgt string) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.graph.AddEdge(src, tgt)
}

func (g *Graph) SetVertexHealth(key string, health bool) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.graph.SetVertexHealth(key, health)
}

func (g *Graph) ClearGraphHealthyStatus() {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.graph.ClearHealthyStatus()
}

func (g *Graph) GetVertex(key string) (Vertex, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	v, err := g.graph.Find(key)
	if err != nil {
		return Vertex{}, err
	}
	return v, nil
}

func (g *Graph) VertexDependents(key string, all bool) (Subgraph, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.graph.VertexDependents(key, all)
}

func (g *Graph) VertexDependencies(key string, all bool) (Subgraph, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.graph.VertexDependencies(key, all)
}

func (g *Graph) Path(src, tgt string) (Subgraph, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.graph.Path(src, tgt)
}
