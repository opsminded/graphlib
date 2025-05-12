package graphlib_test

import (
	"testing"

	"github.com/opsminded/graphlib"
)

func TestGraph_Basics(t *testing.T) {
	g := graphlib.NewGraph()
	g.NewVertex("A")
	g.NewVertex("B")

	g.NewEdge("A-B", "A", "B")

	if g.GetVertexByLabel("A").Health != true {
		t.Errorf("Expected vertex A to be healthy, got %v", g.GetVertexByLabel("A").Health)
	}

	g.SetVertexHealth("A", false)
	if g.GetVertexByLabel("A").Health != false {
		t.Errorf("Expected vertex A to be unhealthy, got %v", g.GetVertexByLabel("A").Health)
	}

}

func TestGraph_NewEdge_Panic(t *testing.T) {
	g := graphlib.NewGraph()
	g.NewVertex("A")
	g.NewVertex("B")

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic, but did not get one")
		}
	}()

	g.NewEdge("A-B", "A", "C") // C does not exist
}

func TestGraph_Getters(t *testing.T) {
	g := graphlib.NewGraph()
	g.NewVertex("A")
	g.NewVertex("B")
	g.NewEdge("A-B", "A", "B")

	v := g.GetVertexByLabel("A")
	if v.Label != "A" {
		t.Errorf("Expected vertex label 'A', got '%s'", v.Label)
	}

	e := g.GetEdgeByLabel("A-B")
	if e.Label != "A-B" {
		t.Errorf("Expected edge label 'A-B', got '%s'", e.Label)
	}
}

func TestGraph_GetVertexByLabel_Panic(t *testing.T) {
	g := graphlib.NewGraph()
	g.NewVertex("A")
	g.NewVertex("B")
	g.NewEdge("A-B", "A", "B")

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic, but did not get one")
		}
	}()

	g.GetVertexByLabel("C") // C does not exist
}

func TestGraph_GetEdgeByLabel_Panic(t *testing.T) {
	g := graphlib.NewGraph()
	g.NewVertex("A")
	g.NewVertex("B")
	g.NewEdge("A-B", "A", "B")

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic, but did not get one")
		}
	}()

	g.GetEdgeByLabel("C") // C does not exist
}

func TestGraph_Lengths(t *testing.T) {
	g := graphlib.NewGraph()
	g.NewVertex("A")
	g.NewVertex("B")
	g.NewVertex("C")
	g.NewEdge("A-B", "A", "B")
	g.NewEdge("B-C", "B", "C")

	if g.VertexLen() != 3 {
		t.Errorf("Expected vertex A length 3, got %d", g.VertexLen())
	}

	if g.EdgeLen() != 2 {
		t.Errorf("Expected vertex B length 2, got %d", g.EdgeLen())
	}
}

func TestGraph_Neighbors(t *testing.T) {
	g := graphlib.NewGraph()

	g.NewVertex("A")
	g.NewVertex("B")
	g.NewVertex("C")
	g.NewVertex("D")

	g.NewEdge("A->B", "A", "B")
	g.NewEdge("B->C", "B", "C")
	g.NewEdge("C->D", "C", "D")

	// Executa a função
	subgraph := g.Neighbors("B")

	// Esperados
	expectedVertices := map[string]bool{
		"A": true,
		"B": true,
		"C": true,
	}

	expectedEdges := map[string]struct {
		Label       string
		Source      string
		Destination string
	}{
		"A->B": {"A->B", "A", "B"},
		"B->C": {"B->C", "B", "C"},
	}

	// Valida vértices
	for _, v := range subgraph.Vertices {
		if !expectedVertices[v.Label] {
			t.Errorf("unexpected vertex: %s", v.Label)
		}
		delete(expectedVertices, v.Label)
	}
	if len(expectedVertices) > 0 {
		t.Errorf("missing vertices: %v", keys(expectedVertices))
	}

	// Valida arestas
	for _, e := range subgraph.Edges {
		exp, ok := expectedEdges[e.Label]
		if !ok {
			t.Errorf("unexpected edge: %s", e.Label)
			continue
		}
		if e.Source.Label != exp.Source {
			t.Errorf("1. edge %s has wrong endpoints: got %s -> %s, want %s -> %s",
				e.Label, e.Source.Label, e.Destination.Label, exp.Source, exp.Destination)
		}
		if e.Destination.Label != exp.Destination {
			t.Errorf("2. edge %s has wrong endpoints: got %s -> %s, want %s -> %s",
				e.Label, e.Source.Label, e.Destination.Label, exp.Source, exp.Destination)
		}
		delete(expectedEdges, e.Label)
	}
	if len(expectedEdges) > 0 {
		t.Errorf("missing edges: %v", keysMap(expectedEdges))
	}
}

func TestGraph_Neighbors_VertexNotFound_Panic(t *testing.T) {
	g := graphlib.NewGraph()

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic, but did not get one")
		}
	}()

	g.Neighbors("B")
}

func TestGraph_GetVertexDependents(t *testing.T) {
	g := graphlib.NewGraph()

	// Montagem do grafo
	g.NewVertex("A")
	g.NewVertex("B")
	g.NewVertex("C")
	g.NewVertex("D")
	g.NewVertex("E")
	g.NewVertex("F")

	g.NewEdge("A->B", "A", "B")
	g.NewEdge("A->C", "A", "C")
	g.NewEdge("B->D", "B", "D")
	g.NewEdge("B->E", "B", "E")
	g.NewEdge("F->B", "F", "B")

	// Testar dependentes diretos de D
	sub := g.GetVertexDependents("D", false)

	expectedDirect := map[string]bool{
		"B": true,
	}

	for _, v := range sub.Vertices {
		if v.Label == "D" {
			continue
		}
		if !expectedDirect[v.Label] {
			t.Errorf("unexpected direct dependent: %s", v.Label)
		}
		delete(expectedDirect, v.Label)
	}
	if len(expectedDirect) > 0 {
		t.Errorf("missing direct dependents: %v", keys(expectedDirect))
	}

	// Testar dependentes todos de D (all = true)
	sub = g.GetVertexDependents("D", true)

	expectedAll := map[string]bool{
		"A": true,
		"B": true,
		"F": true,
	}

	for _, v := range sub.Vertices {
		if v.Label == "D" {
			continue
		}
		if !expectedAll[v.Label] {
			t.Errorf("unexpected dependent (all=true): %s", v.Label)
		}
		delete(expectedAll, v.Label)
	}
	if len(expectedAll) > 0 {
		t.Errorf("missing dependents (all=true): %v", keys(expectedAll))
	}
}

func TestGraph_GetVertexDependents_VertexNotFound_Panic(t *testing.T) {
	g := graphlib.NewGraph()

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic, but did not get one")
		}
	}()

	g.GetVertexDependents("X", false)
}

func TestGraph_GetVertexDependencies(t *testing.T) {
	g := graphlib.NewGraph()

	g.NewVertex("A")
	g.NewVertex("B")
	g.NewVertex("C")
	g.NewVertex("D")
	g.NewVertex("E")
	g.NewVertex("F")
	g.NewVertex("G")

	g.NewEdge("A->B", "A", "B")
	g.NewEdge("A->C", "A", "C")
	g.NewEdge("B->D", "B", "D")
	g.NewEdge("B->E", "B", "E")
	g.NewEdge("C->F", "C", "F")
	g.NewEdge("F->G", "F", "G")

	// Teste 1: dependências diretas de G (should return F)
	sub := g.GetVertexDependencies("C", false)

	expectedVertices := map[string]bool{
		"C": true,
		"F": true,
	}
	expectedEdges := map[string][2]string{
		"C->F": {"C", "F"},
	}

	checkSubgraph(t, sub, expectedVertices, expectedEdges)

	// Teste 2: dependências completas de G (should return F, C, A)
	sub = g.GetVertexDependencies("C", true)

	expectedVertices = map[string]bool{
		"C": true,
		"F": true,
		"G": true,
	}
	expectedEdges = map[string][2]string{
		"C->F": {"C", "F"},
		"F->G": {"F", "G"},
	}

	checkSubgraph(t, sub, expectedVertices, expectedEdges)
}

func TestGraph_GetVertexDependencies_VertexNotFound_Panic(t *testing.T) {
	g := graphlib.NewGraph()

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic, but did not get one")
		}
	}()

	g.GetVertexDependencies("X", false)
}

func TestGraph_Path(t *testing.T) {
	g := graphlib.NewGraph()

	// Montagem do grafo
	g.NewVertex("A")
	g.NewVertex("B")
	g.NewVertex("C")
	g.NewVertex("D")
	g.NewVertex("E")
	g.NewVertex("F")
	g.NewVertex("G")

	g.NewEdge("A->B", "A", "B")
	g.NewEdge("A->C", "A", "C")
	g.NewEdge("A->D", "A", "D")
	g.NewEdge("B->E", "B", "E")
	g.NewEdge("B->F", "B", "F")
	g.NewEdge("C->F", "C", "F")
	g.NewEdge("F->G", "F", "G")

	// Executa a função
	sub := g.Path("A", "G")

	// Esperados
	expectedVertices := map[string]bool{
		"A": true,
		"B": true,
		"C": true,
		"F": true,
		"G": true,
	}
	expectedEdges := map[string][2]string{
		"A->B": {"A", "B"},
		"A->C": {"A", "C"},
		"B->F": {"B", "F"},
		"C->F": {"C", "F"},
		"F->G": {"F", "G"},
	}

	checkSubgraph(t, sub, expectedVertices, expectedEdges)
}

func TestGraph_Path_VertexNotFound_Panic(t *testing.T) {
	g := graphlib.NewGraph()

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic, but did not get one")
		}
	}()

	g.Path("A", "X")
}

func keys(m map[string]bool) []string {
	k := make([]string, 0, len(m))
	for key := range m {
		k = append(k, key)
	}
	return k
}

func keysMap(m map[string]struct {
	Label       string
	Source      string
	Destination string
}) []string {
	k := make([]string, 0, len(m))
	for key := range m {
		k = append(k, key)
	}
	return k
}

func checkSubgraph(t *testing.T, sub graphlib.Subgraph, wantVertices map[string]bool, wantEdges map[string][2]string) {
	// Checa vértices
	seenVertices := map[string]bool{}
	for _, v := range sub.Vertices {
		seenVertices[v.Label] = true
	}

	for label := range wantVertices {
		if !seenVertices[label] {
			t.Errorf("expected vertex %s not found", label)
		}
	}
	for label := range seenVertices {
		if !wantVertices[label] {
			t.Errorf("unexpected vertex %s found", label)
		}
	}

	// Checa arestas
	seenEdges := map[string][2]string{}
	for _, e := range sub.Edges {
		seenEdges[e.Label] = [2]string{e.Source.Label, e.Destination.Label}
	}

	for label, pair := range wantEdges {
		if got, ok := seenEdges[label]; !ok {
			t.Errorf("expected edge %s not found", label)
		} else if got != pair {
			t.Errorf("edge %s has wrong direction: got %v -> %v, want %v -> %v",
				label, got[0], got[1], pair[0], pair[1])
		}
	}
	for label := range seenEdges {
		if _, ok := wantEdges[label]; !ok {
			t.Errorf("unexpected edge %s found", label)
		}
	}
}
