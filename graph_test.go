package graphlib_test

// import (
// 	"errors"
// 	"testing"

// 	"github.com/opsminded/graphlib"
// )

// func TestGraph_Basics(t *testing.T) {
// 	g := graphlib.NewGraph()
// 	g.NewVertex("A")
// 	g.NewVertex("B")

// 	err := g.NewEdge("A", "B")
// 	if err != nil {
// 		t.Errorf("Expected no error, got %v", err)
// 	}

// 	v, err := g.GetVertexByLabel("A")
// 	if err != nil {
// 		t.Errorf("Expected no error, got %v", err)
// 	}
// 	if v.Healthy != true {
// 		t.Errorf("Expected vertex A to be healthy, got %v", v.Healthy)
// 	}

// 	g.SetVertexHealth("A", false)
// 	v, err = g.GetVertexByLabel("A")
// 	if err != nil {
// 		t.Errorf("Expected no error, got %v", err)
// 	}

// 	if v.Healthy != false {
// 		t.Errorf("Expected vertex A to be unhealthy, got %v", v.Healthy)
// 	}

// 	err = g.SetVertexHealth("X", false)
// 	if !errors.As(err, &graphlib.VertexNotFoundError{}) {
// 		t.Errorf("Expected vertex C to not be found, got %v", err)
// 	}
// }

// func TestGraph_ClearGraphHealthyStatus(t *testing.T) {
// 	g := graphlib.NewGraph()
// 	g.NewVertex("A")
// 	g.NewVertex("B")
// 	g.NewEdge("A", "B")

// 	err := g.SetVertexHealth("B", false)
// 	if err != nil {
// 		t.Errorf("Expected no error, got %v", err)
// 	}

// 	vB, _ := g.GetVertexByLabel("A")
// 	if vB.Healthy != false {
// 		t.Errorf("Expected vertex B to be unhealthy, got %v", vB.Healthy)
// 	}

// 	g.ClearGraphHealthyStatus()

// 	v, err := g.GetVertexByLabel("A")
// 	if err != nil {
// 		t.Errorf("Expected no error, got %v", err)
// 	}
// 	if v.Healthy != true {
// 		t.Errorf("Expected vertex A to be healthy, got %v", v.Healthy)
// 	}

// 	if err := g.SetVertexHealth("A", false); err != nil {
// 		t.Errorf("Expected no error, got %v", err)
// 	}

// 	if err := g.SetVertexHealth("A", true); err != nil {
// 		t.Errorf("Expected no error, got %v", err)
// 	}
// }

// func TestGraph_NewEdge_Error(t *testing.T) {
// 	g := graphlib.NewGraph()
// 	g.NewVertex("A")
// 	g.NewVertex("B")

// 	err := g.NewEdge("A", "C") // C does not exist
// 	if !errors.As(err, &graphlib.VertexNotFoundError{}) {
// 		t.Errorf("Expected vertex C to not be found, got %v", err)
// 	}

// 	err = g.NewEdge("C", "B") // C does not exist
// 	if !errors.As(err, &graphlib.VertexNotFoundError{}) {
// 		t.Errorf("Expected vertex C to not be found, got %v", err)
// 	}
// }

// func TestGraph_Getters(t *testing.T) {
// 	g := graphlib.NewGraph()
// 	g.NewVertex("A")
// 	g.NewVertex("B")
// 	g.NewEdge("A", "B")

// 	v, err := g.GetVertexByLabel("A")
// 	if err != nil {
// 		t.Errorf("Expected no error, got %v", err)
// 	}
// 	if v.Label != "A" {
// 		t.Errorf("Expected vertex label 'A', got '%s'", v.Label)
// 	}
// }

// func TestGraph_GetVertexByLabel_Error(t *testing.T) {
// 	g := graphlib.NewGraph()
// 	g.NewVertex("A")
// 	g.NewVertex("B")
// 	g.NewEdge("A", "B")
// 	_, err := g.GetVertexByLabel("C") // C does not exist
// 	if !errors.As(err, &graphlib.VertexNotFoundError{}) {
// 		t.Errorf("Expected vertex C to not be found, got %v", err)
// 	}
// }

// func TestGraph_Lengths(t *testing.T) {
// 	g := graphlib.NewGraph()
// 	g.NewVertex("A")
// 	g.NewVertex("B")
// 	g.NewVertex("C")
// 	g.NewEdge("A", "B")
// 	g.NewEdge("B", "C")

// 	if g.VertexLen() != 3 {
// 		t.Errorf("Expected vertex A length 3, got %d", g.VertexLen())
// 	}

// 	if g.EdgeLen() != 2 {
// 		t.Errorf("Expected vertex B length 2, got %d", g.EdgeLen())
// 	}
// }

// func TestGraph_Neighbors(t *testing.T) {
// 	g := graphlib.NewGraph()

// 	g.NewVertex("A")
// 	g.NewVertex("B")
// 	g.NewVertex("C")
// 	g.NewVertex("D")

// 	g.NewEdge("A", "B")
// 	g.NewEdge("B", "C")
// 	g.NewEdge("C", "D")

// 	// Executa a função
// 	subgraph, err := g.Neighbors("B")
// 	if err != nil {
// 		t.Errorf("Expected no error, got %v", err)
// 	}

// 	// Esperados
// 	expectedVertices := map[string]bool{
// 		"A": true,
// 		"B": true,
// 		"C": true,
// 	}

// 	expectedEdges := map[string]struct {
// 		Source string
// 		Target string
// 	}{
// 		"A-B": {"A", "B"},
// 		"B-C": {"B", "C"},
// 	}

// 	// Valida vértices
// 	for _, v := range subgraph.Vertices {
// 		if !expectedVertices[v.Label] {
// 			t.Errorf("unexpected vertex: %s", v.Label)
// 		}
// 		delete(expectedVertices, v.Label)
// 	}
// 	if len(expectedVertices) > 0 {
// 		t.Errorf("missing vertices: %v", expectedVertices)
// 	}

// 	// Valida arestas
// 	for _, e := range subgraph.Edges {
// 		exp, ok := expectedEdges[e.Source.Label+"-"+e.Target.Label]
// 		if !ok {
// 			t.Errorf("unexpected edge: %s", e.ID)
// 			continue
// 		}
// 		if e.Source.Label != exp.Source {
// 			t.Errorf("1. edge %s has wrong endpoints: got %s -> %s, want %s -> %s",
// 				e.ID, e.Source.Label, e.Target.Label, exp.Source, exp.Target)
// 		}
// 		if e.Target.Label != exp.Target {
// 			t.Errorf("2. edge %s has wrong endpoints: got %s -> %s, want %s -> %s",
// 				e.ID, e.Source.Label, e.Target.Label, exp.Source, exp.Target)
// 		}
// 		delete(expectedEdges, e.Source.Label+"-"+e.Target.Label)
// 	}
// 	if len(expectedEdges) > 0 {
// 		t.Errorf("missing edges: %v", expectedEdges)
// 	}
// }

// func TestGraph_Neighbors_VertexNotFound_Error(t *testing.T) {
// 	g := graphlib.NewGraph()
// 	_, err := g.Neighbors("B")
// 	if err == nil {
// 		t.Errorf("Expected error, got nil")
// 	}
// }

// func TestGraph_GetVertexDependents(t *testing.T) {
// 	g := graphlib.NewGraph()

// 	// Montagem do grafo
// 	g.NewVertex("A")
// 	g.NewVertex("B")
// 	g.NewVertex("C")
// 	g.NewVertex("D")
// 	g.NewVertex("E")
// 	g.NewVertex("F")

// 	g.NewEdge("A", "B")
// 	g.NewEdge("A", "C")
// 	g.NewEdge("B", "D")
// 	g.NewEdge("B", "E")
// 	g.NewEdge("F", "B")

// 	// Testar dependentes diretos de D
// 	sub, err := g.GetVertexDependents("D", false)
// 	if err != nil {
// 		t.Errorf("Expected no error, got %v", err)
// 	}

// 	expectedDirect := map[string]bool{
// 		"B": true,
// 	}

// 	for _, v := range sub.Vertices {
// 		if v.Label == "D" {
// 			continue
// 		}
// 		if !expectedDirect[v.Label] {
// 			t.Errorf("unexpected direct dependent: %s", v.Label)
// 		}
// 		delete(expectedDirect, v.Label)
// 	}
// 	if len(expectedDirect) > 0 {
// 		t.Errorf("missing direct dependents: %v", expectedDirect)
// 	}

// 	// Testar dependentes todos de D (all = true)
// 	sub, err = g.GetVertexDependents("D", true)
// 	if err != nil {
// 		t.Errorf("Expected no error, got %v", err)
// 	}

// 	expectedAll := map[string]bool{
// 		"A": true,
// 		"B": true,
// 		"F": true,
// 	}

// 	for _, v := range sub.Vertices {
// 		if v.Label == "D" {
// 			continue
// 		}
// 		if !expectedAll[v.Label] {
// 			t.Errorf("unexpected dependent (all=true): %s", v.Label)
// 		}
// 		delete(expectedAll, v.Label)
// 	}
// 	if len(expectedAll) > 0 {
// 		t.Errorf("missing dependents (all=true): %v", expectedAll)
// 	}
// }

// func TestGraph_GetVertexDependents_VertexNotFound_Error(t *testing.T) {
// 	g := graphlib.NewGraph()
// 	_, err := g.GetVertexDependents("X", false)
// 	if err == nil {
// 		t.Errorf("Expected error, got nil")
// 	}
// }

// func TestGraph_GetVertexDependencies(t *testing.T) {
// 	g := graphlib.NewGraph()

// 	g.NewVertex("A")
// 	g.NewVertex("B")
// 	g.NewVertex("C")
// 	g.NewVertex("D")
// 	g.NewVertex("E")
// 	g.NewVertex("F")
// 	g.NewVertex("G")

// 	g.NewEdge("A", "B")
// 	g.NewEdge("A", "C")
// 	g.NewEdge("B", "D")
// 	g.NewEdge("B", "E")
// 	g.NewEdge("C", "F")
// 	g.NewEdge("F", "G")

// 	// Teste 1: dependências diretas de G (should return F)
// 	sub, err := g.GetVertexDependencies("C", false)
// 	if err != nil {
// 		t.Errorf("Expected no error, got %v", err)
// 	}

// 	expectedVertices := map[string]bool{
// 		"C": true,
// 		"F": true,
// 	}
// 	expectedEdges := map[string][2]string{
// 		"C-F": {"C", "F"},
// 	}

// 	checkSubgraph(t, sub, expectedVertices, expectedEdges)

// 	// Teste 2: dependências completas de G (should return F, C, A)
// 	sub, err = g.GetVertexDependencies("C", true)
// 	if err != nil {
// 		t.Errorf("Expected no error, got %v", err)
// 	}

// 	expectedVertices = map[string]bool{
// 		"C": true,
// 		"F": true,
// 		"G": true,
// 	}
// 	expectedEdges = map[string][2]string{
// 		"C-F": {"C", "F"},
// 		"F-G": {"F", "G"},
// 	}

// 	checkSubgraph(t, sub, expectedVertices, expectedEdges)
// }

// func TestGraph_GetVertexDependencies_VertexNotFound_Error(t *testing.T) {
// 	g := graphlib.NewGraph()
// 	_, err := g.GetVertexDependencies("X", false)
// 	if err == nil {
// 		t.Errorf("Expected error, got nil")
// 	}
// }

// // func TestGraph_Path(t *testing.T) {
// // 	g := graphlib.NewGraph()

// // 	// Montagem do grafo
// // 	g.NewVertex("A")
// // 	g.NewVertex("B")
// // 	g.NewVertex("C")
// // 	g.NewVertex("D")
// // 	g.NewVertex("E")
// // 	g.NewVertex("F")
// // 	g.NewVertex("G")

// // 	g.NewEdge("A", "B")
// // 	g.NewEdge("A", "C")
// // 	g.NewEdge("A", "D")
// // 	g.NewEdge("B", "E")
// // 	g.NewEdge("B", "F")
// // 	g.NewEdge("C", "F")
// // 	g.NewEdge("F", "G")

// // 	// Executa a função
// // 	sub, err := g.Path("A", "G")
// // 	if err != nil {
// // 		t.Errorf("Expected no error, got %v", err)
// // 	}

// // 	// Esperados
// // 	expectedVertices := map[string]bool{
// // 		"A": true,
// // 		"B": true,
// // 		"C": true,
// // 		"F": true,
// // 		"G": true,
// // 	}
// // 	expectedEdges := map[string][2]string{
// // 		"A->B": {"A", "B"},
// // 		"A->C": {"A", "C"},
// // 		"B->F": {"B", "F"},
// // 		"C->F": {"C", "F"},
// // 		"F->G": {"F", "G"},
// // 	}

// // 	checkSubgraph(t, sub, expectedVertices, expectedEdges)
// // }

// // func TestGraph_Path_VertexNotFound_Error(t *testing.T) {
// // 	g := graphlib.NewGraph()
// // 	_, err := g.Path("A", "X")
// // 	if err == nil {
// // 		t.Errorf("Expected error, got nil")
// // 	}
// // }

// // func TestGraph_UnhealthyNodes(t *testing.T) {
// // 	g := graphlib.NewGraph()
// // 	g.NewVertex("A")
// // 	g.NewVertex("B")

// // 	if u := g.UnhealthyVertices(); len(u) != 0 {
// // 		t.Errorf("Expected no unhealthy nodes, got %d", len(u))
// // 	}

// // 	g.SetVertexHealth("A", false)
// // 	if u := g.UnhealthyVertices(); len(u) != 1 {
// // 		t.Errorf("Expected no unhealthy nodes, got %d", len(u))
// // 	}

// // 	g.SetVertexHealth("A", true)
// // 	if u := g.UnhealthyVertices(); len(u) != 0 {
// // 		t.Errorf("Expected no unhealthy nodes, got %d", len(u))
// // 	}

// // }

// func checkSubgraph(t *testing.T, sub graphlib.Subgraph, wantVertices map[string]bool, wantEdges map[string][2]string) {
// 	// Checa vértices
// 	seenVertices := map[string]bool{}
// 	for _, v := range sub.Vertices {
// 		seenVertices[v.Label] = true
// 	}

// 	for label := range wantVertices {
// 		if !seenVertices[label] {
// 			t.Errorf("expected vertex %s not found", label)
// 		}
// 	}
// 	for label := range seenVertices {
// 		if !wantVertices[label] {
// 			t.Errorf("unexpected vertex %s found", label)
// 		}
// 	}

// 	// Checa arestas
// 	seenEdges := map[string][2]string{}
// 	for _, e := range sub.Edges {
// 		seenEdges[e.Source.Label+"-"+e.Target.Label] = [2]string{e.Source.Label, e.Target.Label}
// 	}

// 	for label, pair := range wantEdges {
// 		if got, ok := seenEdges[label]; !ok {
// 			t.Errorf("expected edge %s not found", label)
// 		} else if got != pair {
// 			t.Errorf("edge %s has wrong direction: got %v -> %v, want %v -> %v",
// 				label, got[0], got[1], pair[0], pair[1])
// 		}
// 	}
// 	for label := range seenEdges {
// 		if _, ok := wantEdges[label]; !ok {
// 			t.Errorf("unexpected edge %s found", label)
// 		}
// 	}
// }
