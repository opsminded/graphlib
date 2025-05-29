package graphlib

import (
	"log/slog"
)

func (g *Graph) ClearHealthyStatus() {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.logger.Debug("core.Graph.ClearHealthyStatus")
	for k := range g.healthy {
		g.healthy[k] = true
	}
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

	g.propagateUnhealthy(v)

	return nil
}

func (g *Graph) propagateUnhealthy(v int) {
	g.healthy[v] = false
	g.lastCheck[v] = g.nowFn()

	for d := range g.dependents[v] {
		g.propagateUnhealthy(d)
	}
}
