package graphlib

import (
	"context"
	"log/slog"
	"time"
)

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
