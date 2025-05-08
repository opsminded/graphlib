package graphlib

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

func (g *Graph) GetVertexClass(v *vertex) string {
	return g.vertexClasses[v.class]
}

func (g *Graph) GetEdgeClass(e *edge) string {
	return g.edgeClasses[e.class]
}

func (g *Graph) VertexLen() int {
	return len(g.vertices)
}

func (g *Graph) EdgeLen() int {
	return len(g.edges)
}

func (g *Graph) UnhealthVertices() []*vertex {
	g.mu.Lock()
	defer g.mu.Unlock()

	var list []*vertex
	for _, v := range g.vertices {
		if !v.health {
			list = append(list, v)
		}
	}

	return list
}

func (g *Graph) Neighbors(v *vertex) resultSet {
	g.mu.Lock()
	defer g.mu.Unlock()

	rs := resultSet{
		Principal: v,
	}

	// Dependências: arestas saindo de v
	if outgoing, ok := g.edges[v.id]; ok {
		for _, e := range outgoing {
			rs.Edges = append(rs.Edges, e)
			rs.Vertices = append(rs.Vertices, e.destination)
		}
	}

	// Dependentes: arestas chegando em v
	for _, dests := range g.edges {
		if e, ok := dests[v.id]; ok {
			rs.Edges = append(rs.Edges, e)
			rs.Vertices = append(rs.Vertices, e.source)
		}
	}

	return rs
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

func (g *Graph) EdgeSourceLabel(e *edge) string {
	return g.vertices[e.source.id].label
}

func (g *Graph) EdgeDestinationLabel(e *edge) string {
	return g.vertices[e.destination.id].label
}

func (g *Graph) GetVertexDependencies(label string, all bool) resultSet {
	g.mu.Lock()
	defer g.mu.Unlock()

	var principal *vertex
	for _, v := range g.vertices {
		if v.label == label {
			principal = v
			break
		}
	}
	if principal == nil {
		return resultSet{} // Vértice não encontrado
	}

	visited := make(map[id]bool)
	var vertices []*vertex
	var edges []*edge

	var dfs func(v *vertex)
	dfs = func(v *vertex) {
		for _, neighbor := range v.neighbors {
			if !visited[neighbor.id] {
				visited[neighbor.id] = true
				vertices = append(vertices, neighbor)

				if edgeMap, ok := g.edges[v.id]; ok {
					if e, ok := edgeMap[neighbor.id]; ok {
						edges = append(edges, e)
					}
				}

				if all {
					dfs(neighbor)
				}
			}
		}
	}

	visited[principal.id] = true
	dfs(principal)

	return resultSet{
		Principal: principal,
		All:       all,
		Vertices:  vertices,
		Edges:     edges,
	}
}
