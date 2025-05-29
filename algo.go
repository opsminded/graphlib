package graphlib

import (
	"fmt"
)

type edgeKey struct{ src, tgt int }

func (g *Graph) VertexNeighbors(key string) (Subgraph, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	// lookup
	rootID, ok := g.lookup[key]
	if !ok {
		return Subgraph{}, VertexNotFoundErr{Key: key}
	}

	// coletores
	verticesSet := map[int]struct{}{rootID: {}}
	edgesSet := make(map[edgeKey]struct{}, 8)
	addEdge := func(s, t int) { edgesSet[edgeKey{s, t}] = struct{}{} }

	// vizinhos diretos
	if outs, ok := g.dependencies[rootID]; ok {
		for tgt := range outs {
			verticesSet[tgt] = struct{}{}
			addEdge(rootID, tgt)
		}
	}

	if ins, ok := g.dependents[rootID]; ok {
		for src := range ins {
			verticesSet[src] = struct{}{}
			addEdge(src, rootID)
		}
	}

	// materializa DTO
	vertices := make([]Vertex, 0, len(verticesSet))
	for id := range verticesSet {
		vertices = append(vertices, Vertex{
			Key:       g.keys[id],
			Label:     g.labels[id],
			Class:     g.classLookup[g.classes[id]],
			Healthy:   g.healthy[id],
			LastCheck: g.lastCheck[id],
		})
	}

	edges := make([]Edge, 0, len(edgesSet))
	for k := range edgesSet {
		edges = append(edges, Edge{
			Key:    fmt.Sprintf("%s-%s", g.keys[k.src], g.keys[k.tgt]),
			Source: g.keys[k.src],
			Target: g.keys[k.tgt],
		})
	}

	return Subgraph{Vertices: vertices, Edges: edges}, nil
}

func (g *Graph) VertexDependencies(key string, all bool) (Subgraph, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	// lookup
	rootID, ok := g.lookup[key]
	if !ok {
		return Subgraph{}, VertexNotFoundErr{Key: key}
	}

	// coletores
	verticesSet := map[int]struct{}{rootID: {}}
	edgesSet := make(map[edgeKey]struct{}, 8)
	addEdge := func(s, t int) { edgesSet[edgeKey{s, t}] = struct{}{} }

	// dependentes diretos
	stack := make([]int, 0, 8)

	if outs, ok := g.dependencies[rootID]; ok {
		for tgt := range outs {
			verticesSet[tgt] = struct{}{}
			addEdge(rootID, tgt)
			if all {
				stack = append(stack, tgt)
			}
		}
	}

	// dependentes transitivos
	if all {
		seen := map[int]struct{}{rootID: {}}

		for len(stack) > 0 {
			n := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			if _, dup := seen[n]; dup {
				continue
			}
			seen[n] = struct{}{}

			if outs, ok := g.dependencies[n]; ok {
				for tgt := range outs {
					verticesSet[tgt] = struct{}{}
					addEdge(n, tgt)
					stack = append(stack, tgt)
				}
			}
		}
	}

	// materializar DTO
	vertices := make([]Vertex, 0, len(verticesSet))
	for id := range verticesSet {
		vertices = append(vertices, Vertex{
			Key:       g.keys[id],
			Label:     g.labels[id],
			Class:     g.classLookup[g.classes[id]],
			Healthy:   g.healthy[id],
			LastCheck: g.lastCheck[id],
		})
	}

	edges := make([]Edge, 0, len(edgesSet))
	for k := range edgesSet {
		edges = append(edges, Edge{
			Key:    fmt.Sprintf("%s-%s", g.keys[k.src], g.keys[k.tgt]),
			Source: g.keys[k.src],
			Target: g.keys[k.tgt],
		})
	}

	return Subgraph{Vertices: vertices, Edges: edges}, nil
}

func (g *Graph) VertexDependents(key string, all bool) (Subgraph, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	// lookup
	rootID, ok := g.lookup[key]
	if !ok {
		return Subgraph{}, VertexNotFoundErr{Key: key}
	}

	// coletores
	verticesSet := map[int]struct{}{rootID: {}}
	edgesSet := make(map[edgeKey]struct{}, 8)
	addEdge := func(s, t int) { edgesSet[edgeKey{s, t}] = struct{}{} }

	// dependências diretas
	stack := make([]int, 0, 8)

	if ins, ok := g.dependents[rootID]; ok {
		for src := range ins {
			verticesSet[src] = struct{}{}
			addEdge(src, rootID)
			if all {
				stack = append(stack, src)
			}
		}
	}

	// dependencias transitivas
	if all {
		seen := map[int]struct{}{rootID: {}}

		for len(stack) > 0 {
			n := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			if _, dup := seen[n]; dup {
				continue
			}
			seen[n] = struct{}{}

			if ins, ok := g.dependents[n]; ok {
				for src := range ins {
					verticesSet[src] = struct{}{}
					addEdge(src, n)
					stack = append(stack, src)
				}
			}
		}
	}

	// materializar DTO
	vertices := make([]Vertex, 0, len(verticesSet))
	for id := range verticesSet {
		vertices = append(vertices, Vertex{
			Key:       g.keys[id],
			Label:     g.labels[id],
			Class:     g.classLookup[g.classes[id]],
			Healthy:   g.healthy[id],
			LastCheck: g.lastCheck[id],
		})
	}

	edges := make([]Edge, 0, len(edgesSet))
	for k := range edgesSet {
		edges = append(edges, Edge{
			Key:    fmt.Sprintf("%s-%s", g.keys[k.src], g.keys[k.tgt]),
			Source: g.keys[k.src],
			Target: g.keys[k.tgt],
		})
	}

	return Subgraph{Vertices: vertices, Edges: edges}, nil
}

func (g *Graph) Path(srcKey, tgtKey string) (Subgraph, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	// lookup
	srcID, ok := g.lookup[srcKey]
	if !ok {
		return Subgraph{}, VertexNotFoundErr{Key: srcKey}
	}
	dstID, ok := g.lookup[tgtKey]
	if !ok {
		return Subgraph{}, VertexNotFoundErr{Key: tgtKey}
	}

	// coletores
	verts := map[int]struct{}{}
	edges := map[edgeKey]struct{}{}

	// DFS + memo
	memo := map[int]bool{} // id → existe caminho id→dst ?
	var dfs func(int) bool

	dfs = func(id int) bool {
		if res, ok := memo[id]; ok {
			return res
		}
		if id == dstID {
			verts[id] = struct{}{}
			memo[id] = true
			return true
		}
		found := false
		if outs, ok := g.dependencies[id]; ok {
			for tgt := range outs {
				if dfs(tgt) {
					found = true
					verts[id] = struct{}{}
					verts[tgt] = struct{}{}
					edges[edgeKey{id, tgt}] = struct{}{}
				}
			}
		}
		memo[id] = found
		return found
	}

	if !dfs(srcID) {
		return Subgraph{}, VertexPathErr{Src: srcKey, Dst: tgtKey}
	}

	// materializa DTO
	outV := make([]Vertex, 0, len(verts))
	for id := range verts {
		outV = append(outV, Vertex{
			Key:       g.keys[id],
			Label:     g.labels[id],
			Class:     g.classLookup[g.classes[id]],
			Healthy:   g.healthy[id],
			LastCheck: g.lastCheck[id],
		})
	}
	outE := make([]Edge, 0, len(edges))
	for k := range edges {
		outE = append(outE, Edge{
			Key:    fmt.Sprintf("%s-%s", g.keys[k.src], g.keys[k.tgt]),
			Source: g.keys[k.src],
			Target: g.keys[k.tgt],
		})
	}
	return Subgraph{Vertices: outV, Edges: outE}, nil
}
