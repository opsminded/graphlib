package core_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/opsminded/graphlib/internal/core"
)

func TestVertexDependents(t *testing.T) {

	g := buildGraph()

	sg, err := g.VertexDependents("A", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	wantV := map[string]bool{"A": true, "B": true, "C": true}
	if got := setVerts(sg.Vertices); len(got) != len(wantV) {
		t.Fatalf("vertices mismatch got=%v want=%v", got, wantV)
	}
	wantE := map[e]bool{{"A", "B"}: true, {"A", "C"}: true}
	if got := setEdges(sg.Edges); len(got) != len(wantE) {
		t.Fatalf("edges mismatch got=%v want=%v", got, wantE)
	}
}

func TestVertexDependents_All(t *testing.T) {
	g := buildGraph()

	sg, err := g.VertexDependents("A", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	wantV := map[string]bool{"A": true, "B": true, "C": true, "D": true, "E": true}
	if got := setVerts(sg.Vertices); len(got) != len(wantV) {
		t.Fatalf("vertices mismatch got=%v want=%v", got, wantV)
	}
	wantE := map[e]bool{
		{"A", "B"}: true,
		{"A", "C"}: true,
		{"C", "D"}: true,
		{"D", "E"}: true,
	}
	if got := setEdges(sg.Edges); len(got) != len(wantE) {
		t.Fatalf("edges mismatch got=%v want=%v", got, wantE)
	}
}

func TestVertexDependents_NotFound(t *testing.T) {
	g := buildGraph()

	_, err := g.VertexDependents("X", false)
	var nf core.VertexNotFoundErr
	if !errors.As(err, &nf) || nf.Key != "X" {
		t.Fatalf("expected VertexNotFoundErr(X), got %v", err)
	}

}

func TestVertexDependencies(t *testing.T) {

	g := buildGraph()

	sg, err := g.VertexDependencies("D", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	wantV := map[string]bool{"C": true, "D": true}
	if got := setVerts(sg.Vertices); len(got) != len(wantV) {
		t.Fatalf("vertices mismatch got=%v want=%v", got, wantV)
	}
	wantE := map[e]bool{{"C", "D"}: true}
	if got := setEdges(sg.Edges); len(got) != len(wantE) {
		t.Fatalf("edges mismatch got=%v want=%v", got, wantE)
	}
}

func TestVertexDependencies_All(t *testing.T) {
	g := buildGraph()

	sg, err := g.VertexDependencies("E", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	wantV := map[string]bool{"E": true, "D": true, "C": true, "A": true}
	if got := setVerts(sg.Vertices); len(got) != len(wantV) {
		t.Fatalf("vertices mismatch got=%v want=%v", got, wantV)
	}
	wantE := map[e]bool{
		{"D", "E"}: true,
		{"C", "D"}: true,
		{"A", "C"}: true,
	}
	if got := setEdges(sg.Edges); len(got) != len(wantE) {
		t.Fatalf("edges mismatch got=%v want=%v", got, wantE)
	}
}

func TestVertexDependencies_NotFound(t *testing.T) {
	g := buildGraph()

	_, err := g.VertexDependencies("X", false)
	var nf core.VertexNotFoundErr
	if !errors.As(err, &nf) || nf.Key != "X" {
		t.Fatalf("expected VertexNotFoundErr(X), got %v", err)
	}

}

func TestPath(t *testing.T) {
	g := core.NewSoAGraph()
	for _, k := range []string{"A", "B", "C", "D"} {
		g.AddVertex(k, k, true)
	}
	g.AddEdge("A", "B")
	g.AddEdge("B", "C")
	g.AddEdge("A", "D")

	sg, err := g.Path("A", "C")
	if err != nil {
		t.Fatal(err)
	}

	wantV := map[string]bool{"A": true, "B": true, "C": true}
	if len(sg.Vertices) != 3 {
		t.Fatalf("want 3 vertices, got %d", len(sg.Vertices))
	}

	for _, v := range sg.Vertices {
		if !wantV[v.Key] {
			t.Fatalf("unexpected vertex %s", v.Key)
		}
		delete(wantV, v.Key)
	}

	if len(wantV) != 0 {
		t.Fatalf("missing vertices: %v", wantV)
	}

	wantE := map[string]bool{"A-B": true, "B-C": true}
	fmt.Println(sg.Edges)
	for _, e := range sg.Edges {
		if !wantE[e.Key] {
			t.Fatalf("unexpected edge %s", e.Key)
		}
		delete(wantE, e.Key)
	}

	if len(wantE) != 0 {
		t.Fatalf("missing edges: %v", wantE)
	}

}

/*** helpers ***************************************************************/

// transforma slice de vértices / arestas em conjunto para comparação
func setVerts(vs []core.Vertex) map[string]bool {
	m := make(map[string]bool, len(vs))
	for _, v := range vs {
		m[v.Key] = true
	}
	return m
}

type e struct{ src, dst string }

func setEdges(es []core.Edge) map[e]bool {
	m := make(map[e]bool, len(es))
	for _, ed := range es {
		m[e{ed.Source, ed.Target}] = true
	}
	return m
}

// grafo exemplo:   A → B
//
//	A → C → D
func buildGraph() *core.Graph {
	g := core.NewSoAGraph()
	for _, k := range []string{"A", "B", "C", "D", "E"} {
		g.AddVertex(k, "", true)
	}
	g.AddEdge("A", "B")
	g.AddEdge("A", "C")
	g.AddEdge("C", "D")
	g.AddEdge("D", "E")
	return g
}

// func TestNeighbors(t *testing.T) {
// 	// g := core.NewAoSGraph()

// 	// // vértice raiz
// 	// g.AddVertex("A", "A", true)
// 	// // dependentes
// 	// g.AddVertex("B", "B", true)
// 	// g.AddVertex("C", "C", true)
// 	// // dependências
// 	// g.AddVertex("D", "D", true)
// 	// g.AddVertex("E", "E", true)

// 	// // A → B , A → C
// 	// g.AddEdge("A", "B")
// 	// g.AddEdge("A", "C")
// 	// // D → A , E → A
// 	// g.AddEdge("D", "A")
// 	// g.AddEdge("E", "A")

// 	// sg, err := g.Neighbors("A")
// 	// if err != nil {
// 	// 	t.Fatalf("unexpected error: %v", err)
// 	// }

// 	// // ----- valida vértices ------------------------------------------------
// 	// wantV := map[string]bool{"A": true, "B": true, "C": true, "D": true, "E": true}
// 	// if len(sg.Vertices) != len(wantV) {
// 	// 	t.Fatalf("got %d vertices, want %d", len(sg.Vertices), len(wantV))
// 	// }
// 	// for _, v := range sg.Vertices {
// 	// 	if !wantV[v.Label] {
// 	// 		t.Fatalf("unexpected vertex %q in result", v.Label)
// 	// 	}
// 	// 	delete(wantV, v.Label)
// 	// }
// 	// if len(wantV) != 0 {
// 	// 	t.Fatalf("missing vertices: %v", wantV)
// 	// }

// 	// // ----- valida arestas -------------------------------------------------
// 	// type edge struct{ src, dst string }
// 	// wantE := map[edge]bool{
// 	// 	{"A", "B"}: true,
// 	// 	{"A", "C"}: true,
// 	// 	{"D", "A"}: true,
// 	// 	{"E", "A"}: true,
// 	// }
// 	// if len(sg.Edges) != len(wantE) {
// 	// 	t.Fatalf("got %d edges, want %d", len(sg.Edges), len(wantE))
// 	// }
// 	// for _, e := range sg.Edges {
// 	// 	key := edge{e.Source.Label, e.Target.Label}
// 	// 	if !wantE[key] {
// 	// 		t.Fatalf("unexpected edge %s→%s", key.src, key.dst)
// 	// 	}
// 	// 	delete(wantE, key)
// 	// }
// 	// if len(wantE) != 0 {
// 	// 	t.Fatalf("missing edges: %v", wantE)
// 	// }
// }

// func TestNeighborsVertexNotFound(t *testing.T) {
// 	g := core.NewAoSGraph()
// 	_, err := g.Neighbors("nonexistent")
// 	if err == nil {
// 		t.Fatalf("expected error for missing vertex, got nil")
// 	}
// }

// func TestVertexDependents(t *testing.T) {
// 	g := buildGraph()

// 	sg, err := g.VertexDependents("C", false)
// 	if err != nil {
// 		t.Fatalf("unexpected error: %v", err)
// 	}

// 	wantV := map[string]bool{"C": true, "B": true, "D": true}
// 	if gv := verts(sg.Vertices); len(gv) != len(wantV) {
// 		fmt.Println(sg.Vertices)
// 		t.Fatalf("got %d vertices, want %d", len(gv), len(wantV))
// 	}

// 	wantE := map[e]bool{
// 		{"B", "C"}: true,
// 		{"D", "C"}: true,
// 	}
// 	if ge := edges(sg.Edges); len(ge) != len(wantE) {
// 		t.Fatalf("got %d edges, want %d", len(ge), len(wantE))
// 	}
// }

// func TestVertexDependencies_NotFound(t *testing.T) {
// 	// g := buildGraph()
// 	// if _, err := g.VertexDependencies("X", false); err == nil {
// 	// 	t.Fatalf("expected error for missing vertex, got nil")
// 	// }

// 	// if _, err := g.VertexDependencies("X", true); err == nil {
// 	// 	t.Fatalf("expected error for missing vertex, got nil")
// 	// }
// }

// func buildGraph() *core.Graph {
// 	g := core.NewAoSGraph()
// 	g.AddVertex("A", "", true)
// 	g.AddVertex("B", "", true)
// 	g.AddVertex("C", "", true)
// 	g.AddVertex("D", "", true)

// 	g.AddEdge("A", "B")
// 	g.AddEdge("B", "C")
// 	g.AddEdge("D", "C")
// 	return g
// }

// func verts(vs []*core.Vertex) map[string]bool {
// 	m := make(map[string]bool, len(vs))
// 	for _, v := range vs {
// 		fmt.Println("v", v.Label)
// 		m[v.Label] = true
// 	}
// 	fmt.Println("m", m)
// 	return m
// }

// type e struct{ src, dst string }

// func edges(es []core.Edge) map[e]bool {
// 	m := make(map[e]bool, len(es))
// 	for _, ed := range es {
// 		m[e{ed.Source.Label, ed.Target.Label}] = true
// 	}
// 	return m
// }

// // utilitário para comparar conjuntos sem depender de ordem
// func vertsToSet(vs []core.Vertex) map[string]bool {
// 	m := make(map[string]bool, len(vs))
// 	for _, v := range vs {
// 		m[v.Label] = true
// 	}
// 	return m
// }

// type edge struct{ src, dst string }

// func edgesToSet(es []core.Edge) map[edge]bool {
// 	m := make(map[edge]bool, len(es))
// 	for _, e := range es {
// 		m[edge{e.Source.Label, e.Target.Label}] = true
// 	}
// 	return m
// }
