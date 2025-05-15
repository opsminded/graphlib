package core

import (
	"context"
	"fmt"
	"log"
	"time"
)

type Vertex struct {
	Key       string
	Label     string
	Healthy   bool
	LastCheck int64
}

type Edge struct {
	Key    string
	Source string
	Target string
}

type Subgraph struct {
	Vertices []Vertex
	Edges    []Edge
}

type Stats struct {
	TotalVertices          int
	TotalUnhealthyVertices int
	TotalEdges             int
	TotalHealthyVertices   int

	UnhealthyVertices []Vertex
}

type Graph struct {
	labels       []string
	healthy      []bool
	LastCheck    []int64
	keys         map[int]string
	lookup       map[string]int
	dependents   map[int]map[int]struct{}
	dependencies map[int]map[int]struct{}
	nowFn        func() int64
}

func NewSoAGraph() *Graph {
	g := &Graph{
		labels:       make([]string, 0, 1000),
		healthy:      make([]bool, 0, 1000),
		LastCheck:    make([]int64, 0, 1000),
		keys:         make(map[int]string, 1000),
		lookup:       make(map[string]int, 1000),
		dependents:   make(map[int]map[int]struct{}, 1000),
		dependencies: make(map[int]map[int]struct{}, 1000),
		nowFn:        func() int64 { return time.Now().UnixNano() },
	}

	return g
}

func (g *Graph) AddVertex(key string, label string, healthy bool) int {
	if k, ok := g.lookup[key]; ok {
		return k
	}

	idx := len(g.labels)
	g.keys[idx] = key
	g.lookup[key] = idx

	g.labels = append(g.labels, label)
	g.healthy = append(g.healthy, healthy)
	g.LastCheck = append(g.LastCheck, g.nowFn())

	return idx
}

func (g *Graph) AddEdge(src, tgt string) error {
	ksrc, ok := g.lookup[src]
	if !ok {
		return VertexNotFoundErr{Key: src}
	}

	ktgt, ok := g.lookup[tgt]
	if !ok {
		return VertexNotFoundErr{Key: tgt}
	}

	// prevent edge multiplicity
	if g.exists(ksrc, ktgt) {
		return nil
	}

	// prevent bidirectional edges
	if g.exists(ktgt, ksrc) {
		return BidirectionalEdgeErr{Src: src, Tgt: tgt}
	}

	// prevent cycles
	if g.wouldCreateCycle(ksrc, ktgt) {
		return CycleErr{Src: src, Tgt: tgt}
	}

	if g.dependencies[ksrc] == nil {
		g.dependencies[ksrc] = make(map[int]struct{}, 4)
	}

	if g.dependents[ktgt] == nil {
		g.dependents[ktgt] = make(map[int]struct{}, 4)
	}

	g.dependencies[ksrc][ktgt] = struct{}{}
	g.dependents[ktgt][ksrc] = struct{}{}

	return nil
}

func (g *Graph) Find(key string) (Vertex, error) {
	v, ok := g.lookup[key]
	if !ok {
		return Vertex{}, VertexNotFoundErr{Key: key}
	}

	return Vertex{
		Key:       key,
		Label:     g.labels[v],
		Healthy:   g.healthy[v],
		LastCheck: g.LastCheck[v],
	}, nil
}

func (g *Graph) GraphStats() Stats {
	stats := Stats{
		TotalVertices:          len(g.keys),
		TotalHealthyVertices:   0,
		TotalUnhealthyVertices: 0,
		TotalEdges:             0,
	}

	for _, healthy := range g.healthy {
		if healthy {
			stats.TotalHealthyVertices++
		} else {
			stats.TotalUnhealthyVertices++
		}
	}

	for _, deps := range g.dependencies {
		stats.TotalEdges += len(deps)
	}

	stats.UnhealthyVertices = make([]Vertex, 0, stats.TotalUnhealthyVertices)
	for i, healthy := range g.healthy {
		if !healthy {
			stats.UnhealthyVertices = append(stats.UnhealthyVertices, Vertex{
				Key:       g.keys[i],
				Label:     g.labels[i],
				Healthy:   g.healthy[i],
				LastCheck: g.LastCheck[i],
			})
		}
	}

	return stats
}

func (g *Graph) SetVertexHealth(key string, health bool) error {
	v, ok := g.lookup[key]
	if !ok {
		return VertexNotFoundErr{Key: key}
	}

	g.healthy[v] = health
	g.LastCheck[v] = g.nowFn()

	return nil
}

func (g *Graph) ClearHealthyStatus() {
	for k := range g.healthy {
		g.healthy[k] = true
	}
}

func (g *Graph) StartHealthCheckLoop(ctx context.Context, checkInterval time.Duration) {
	go func() {
		ticker := time.NewTicker(checkInterval)
		defer ticker.Stop()
		log.Println("[core] health check loop started")

		for {
			select {
			case <-ticker.C:
				g.updateHealthStatusAndPropagate(checkInterval)
			case <-ctx.Done():
				log.Println("[core] health check loop stopped")
				return
			}
		}
	}()
}

func (g *Graph) updateHealthStatusAndPropagate(checkInterval time.Duration) {
	now := g.nowFn()

	log.Println("g now", now)
	log.Println("checkInterval.Nanoseconds", checkInterval.Nanoseconds())
	log.Println("go time", time.Now().UnixNano())

	for i, ok := range g.healthy {

		lastCheck := g.LastCheck[i]
		duration := checkInterval.Nanoseconds()
		min := lastCheck + duration
		log.Println("lastCheck", lastCheck)
		log.Println("duration", duration)
		log.Println("min", min)

		if ok && min < now {
			log.Println("healthy", g.keys[i], "unhealthy")
			g.healthy[i] = false
			g.LastCheck[i] = g.nowFn()
		}
	}

	visited := make(map[int]struct{}, 10)
	for i, ok := range g.healthy {
		if !ok {
			g.propagateUnhealthy(i, visited)
		}
	}
}

func (g *Graph) propagateUnhealthy(v int, visited map[int]struct{}) {
	log.Println("propagateUnhealthy", g.keys[v])

	if _, seen := visited[v]; seen {
		return
	}
	visited[v] = struct{}{}
	g.healthy[v] = false
	g.LastCheck[v] = g.nowFn()

	for d := range g.dependents[v] {
		g.propagateUnhealthy(d, visited)
	}
}

func (g *Graph) exists(src, tgt int) bool {
	_, ok := g.dependencies[src][tgt]
	return ok
}

func (g *Graph) wouldCreateCycle(src, tgt int) bool {
	if src == tgt {
		return true
	}

	// Check if the target vertex is reachable from the source vertex
	visited := make(map[int]struct{}, 10)
	stack := []int{tgt}

	for len(stack) > 0 {
		n := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if _, seen := visited[n]; seen {
			continue
		}
		visited[n] = struct{}{}

		if n == src {
			return true
		}

		if deps, ok := g.dependencies[n]; ok {
			for v := range deps {
				stack = append(stack, v)
			}
		}
	}

	return false
}

type VertexNotFoundErr struct {
	Key string
}

func (e VertexNotFoundErr) Error() string {
	return fmt.Sprintf("vertex %q not found", e.Key)
}

type BidirectionalEdgeErr struct {
	Src, Tgt string
}

func (e BidirectionalEdgeErr) Error() string {
	return fmt.Sprintf("bidirectional edge %s ↔ %s not allowed", e.Src, e.Tgt)
}

type CycleErr struct {
	Src string
	Tgt string
}

func (e CycleErr) Error() string {
	return fmt.Sprintf("edge %s → %s would create a cycle", e.Src, e.Tgt)
}

type VertexPathErr struct {
	Src string
	Dst string
}

func (e VertexPathErr) Error() string {
	return fmt.Sprintf("no path from %s to %s", e.Src, e.Dst)
}
