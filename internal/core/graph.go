package core

import (
	"context"
	"log/slog"
	"os"
	"sync"
	"time"
)

type Graph struct {
	labels       []string
	healthy      []bool
	lastCheck    []int64
	keys         map[int]string
	lookup       map[string]int
	dependents   map[int]map[int]struct{}
	dependencies map[int]map[int]struct{}
	nowFn        func() int64
	logger       *slog.Logger
	mu           sync.RWMutex
}

func NewSoAGraph(logger *slog.Logger) *Graph {
	if logger == nil {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	}

	g := &Graph{
		labels:       make([]string, 0, 1000),
		healthy:      make([]bool, 0, 1000),
		lastCheck:    make([]int64, 0, 1000),
		keys:         make(map[int]string, 1000),
		lookup:       make(map[string]int, 1000),
		dependents:   make(map[int]map[int]struct{}, 1000),
		dependencies: make(map[int]map[int]struct{}, 1000),
		nowFn:        func() int64 { return time.Now().UnixNano() },
		logger:       logger,
		mu:           sync.RWMutex{},
	}

	return g
}

func (g *Graph) AddVertex(key string, label string, healthy bool) int {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.logger.Debug("core.Graph.AddVertex", slog.String("key", key), slog.String("label", label), slog.Bool("healthy", healthy))

	if k, ok := g.lookup[key]; ok {
		g.logger.Debug("core.Graph.AddVertex lookup found vertex. returning the integer id", slog.String("key", key), slog.Int("id", k))
		return k
	}

	idx := len(g.labels)
	g.logger.Debug("core.Graph.AddVertex could not find vertex. A new vertex will be created", slog.String("key", key), slog.Int("id", idx))

	g.keys[idx] = key
	g.lookup[key] = idx

	g.labels = append(g.labels, label)
	g.healthy = append(g.healthy, healthy)
	g.lastCheck = append(g.lastCheck, g.nowFn())

	return idx
}

func (g *Graph) AddEdge(src, tgt string) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.logger.Debug("core.Graph.AddEdge", slog.String("src", src), slog.String("tgt", tgt))

	ksrc, ok := g.lookup[src]
	if !ok {
		err := VertexNotFoundErr{Key: src}
		g.logger.Error("core.Graph.AddEdge src lookup error", slog.String("key", src), slog.String("err", err.Error()))
		return err
	}

	ktgt, ok := g.lookup[tgt]
	if !ok {
		err := VertexNotFoundErr{Key: tgt}
		g.logger.Error("core.Graph.AddEdge tgt lookup error", slog.String("key", tgt), slog.String("err", err.Error()))
		return err
	}

	g.logger.Debug("core.Graph.AddEdge src and tgt lookup success", slog.String("src", src), slog.String("tgt", tgt))

	// prevent edge multiplicity
	if g.exists(ksrc, ktgt) {
		g.logger.Info("core.Graph.AddEdge already exists", slog.String("src", src), slog.String("tgt", tgt))
		return nil
	}

	// prevent bidirectional edges
	if g.exists(ktgt, ksrc) {
		err := BidirectionalEdgeErr{Src: src, Tgt: tgt}
		g.logger.Error("core.Graph.AddEdge will cause a bidirectional relation", slog.String("src", src), slog.String("tgt", tgt), slog.String("err", err.Error()))
		return err
	}

	// prevent cycles
	if g.wouldCreateCycle(ksrc, ktgt) {
		err := CycleErr{Src: src, Tgt: tgt}
		g.logger.Error("core.Graph.AddEdge will cause a cycle", slog.String("src", src), slog.String("tgt", tgt), slog.String("err", err.Error()))
		return err
	}

	if g.dependencies[ksrc] == nil {
		g.dependencies[ksrc] = make(map[int]struct{}, 4)
	}

	if g.dependents[ktgt] == nil {
		g.dependents[ktgt] = make(map[int]struct{}, 4)
	}

	g.logger.Debug("core.Graph.AddEdge Edge will be created", slog.String("src", src), slog.String("tgt", tgt))

	g.dependencies[ksrc][ktgt] = struct{}{}
	g.dependents[ktgt][ksrc] = struct{}{}

	return nil
}

func (g *Graph) Find(key string) (Vertex, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	g.logger.Debug("core.Graph.Find", slog.String("key", key))

	v, ok := g.lookup[key]
	if !ok {
		err := VertexNotFoundErr{Key: key}
		g.logger.Error("core.Graph.Find lookup error", slog.String("key", key), slog.String("err", err.Error()))
		return Vertex{}, err
	}

	return Vertex{
		Key:       key,
		Label:     g.labels[v],
		Healthy:   g.healthy[v],
		LastCheck: g.lastCheck[v],
	}, nil
}

func (g *Graph) Stats() Stats {
	g.mu.RLock()
	defer g.mu.RUnlock()

	g.logger.Debug("core.Graph.GraphStats")

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
				LastCheck: g.lastCheck[i],
			})
		}
	}

	g.logger.Info("core.Graph.GraphStats",
		slog.Int("TotalVertices", stats.TotalVertices),
		slog.Int("TotalHealthyVertices", stats.TotalHealthyVertices),
		slog.Int("TotalUnhealthyVertices", stats.TotalUnhealthyVertices),
		slog.Int("TotalEdges", stats.TotalEdges))

	return stats
}

func (g *Graph) SetVertexHealth(key string, health bool) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.logger.Debug("core.Graph.SetVertexHealth", slog.String("key", key), slog.Bool("health", health))

	v, ok := g.lookup[key]
	if !ok {
		err := VertexNotFoundErr{Key: key}
		g.logger.Error("core.Graph.SetVertexHealth lookup error", slog.String("key", key), slog.String("err", err.Error()))
		return err
	}

	g.logger.Info("core.Graph.SetVertexHealth lookup success. The health status will be changed", slog.String("key", key), slog.Int("id", v), slog.Bool("health", health))

	g.healthy[v] = health
	g.lastCheck[v] = g.nowFn()

	return nil
}

func (g *Graph) ClearHealthyStatus() {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.logger.Debug("core.Graph.ClearHealthyStatus")
	for k := range g.healthy {
		g.healthy[k] = true
	}
}

func (g *Graph) StartHealthCheckLoop(ctx context.Context, checkInterval time.Duration) {
	g.logger.Debug("core.Graph.StartHealthCheckLoop", slog.Duration("checkInterval", checkInterval))

	go func() {
		ticker := time.NewTicker(checkInterval)
		defer ticker.Stop()
		g.logger.Info("core.Graph.StartHealthCheckLoop go routine started", slog.Duration("checkInterval", checkInterval))

		for {
			select {
			case <-ticker.C:
				g.updateHealthStatusAndPropagate(checkInterval)
			case <-ctx.Done():
				g.logger.Info("core.Graph.StartHealthCheckLoop context done")
				return
			}
		}
	}()
}

func (g *Graph) updateHealthStatusAndPropagate(checkInterval time.Duration) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.logger.Debug("core.Graph.updateHealthStatusAndPropagate", slog.Duration("checkInterval", checkInterval))

	now := g.nowFn()

	for i, ok := range g.healthy {
		lastCheck := g.lastCheck[i]
		duration := checkInterval.Nanoseconds()
		min := lastCheck + duration

		if ok && min < now {
			g.logger.Info("core.Graph.updateHealthStatusAndPropagate will change vertex health to false", slog.Int("id", i), slog.String("key", g.keys[i]))
			g.healthy[i] = false
			g.lastCheck[i] = g.nowFn()
		}
	}

	visited := make(map[int]struct{}, 10)
	for i, ok := range g.healthy {
		if !ok {
			g.logger.Info("core.Graph.updateHealthStatusAndPropagate will propagate now", slog.Int("id", i), slog.String("key", g.keys[i]))
			g.propagateUnhealthy(i, visited)
		}
	}
}

func (g *Graph) propagateUnhealthy(v int, visited map[int]struct{}) {
	if _, seen := visited[v]; seen {
		return
	}
	visited[v] = struct{}{}
	g.healthy[v] = false
	g.lastCheck[v] = g.nowFn()

	for d := range g.dependents[v] {
		g.propagateUnhealthy(d, visited)
	}
}

func (g *Graph) exists(src, tgt int) bool {
	g.logger.Debug("core.Graph.exists", slog.Int("src", src), slog.Int("tgt", tgt))
	_, ok := g.dependencies[src][tgt]
	return ok
}

func (g *Graph) wouldCreateCycle(src, tgt int) bool {
	g.logger.Debug("core.Graph.wouldCreateCycle", slog.Int("src", src), slog.Int("tgt", tgt))

	if src == tgt {
		g.logger.Info("core.Graph.wouldCreateCycle src and tgt is the same", slog.Int("src", src), slog.Int("tgt", tgt))
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
