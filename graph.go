package graphlib

import (
	"context"
	"time"

	"github.com/opsminded/graphlib/v2/internal/core"
)

type (
	Edge     = core.Edge
	Vertex   = core.Vertex
	Subgraph = core.Subgraph
	Stats    = core.Stats
)

type Graph struct {
	graph *core.Graph
}

func NewGraph() *Graph {
	g := &Graph{
		graph: core.NewSoAGraph(nil),
	}
	return g
}

func (g *Graph) AddVertex(key, label string, healthy bool) {
	g.graph.AddVertex(key, label, healthy)
}

func (g *Graph) AddEdge(src, tgt string) error {
	return g.graph.AddEdge(src, tgt)
}

func (g *Graph) GetVertex(key string) (Vertex, error) {
	return g.graph.Find(key)
}

func (g *Graph) SetVertexHealth(key string, health bool) error {
	return g.graph.SetVertexHealth(key, health)
}

func (g *Graph) ClearGraphHealthyStatus() {
	g.graph.ClearHealthyStatus()
}

func (g *Graph) StartHealthCheckLoop(ctx context.Context, check time.Duration) {
	g.graph.StartHealthCheckLoop(ctx, check)
}

func (g *Graph) GraphStats() Stats {
	return g.graph.Stats()
}

func (g *Graph) VertexDependents(key string, all bool) (Subgraph, error) {
	return g.graph.VertexDependents(key, all)
}

func (g *Graph) VertexDependencies(key string, all bool) (Subgraph, error) {
	return g.graph.VertexDependencies(key, all)
}

func (g *Graph) Path(src, tgt string) (Subgraph, error) {
	return g.graph.Path(src, tgt)
}

func (g *Graph) VertexNeighbors(key string) (Subgraph, error) {
	return g.graph.VertexNeighbors(key)
}
