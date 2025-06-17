package graphlib_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/opsminded/graphlib/v2"
)

func TestVertexNeighbors(t *testing.T) {

	g := buildGraph()

	sg, err := g.VertexNeighbors("C")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	wantV := map[string]bool{"A": true, "C": true, "D": true}
	if got := setVerts(sg.Vertices); len(got) != len(wantV) {
		t.Fatalf("vertices mismatch got=%v want=%v", got, wantV)
	}
	wantE := map[e]bool{{"A", "C"}: true, {"C", "D"}: true}
	if got := setEdges(sg.Edges); len(got) != len(wantE) {
		t.Fatalf("edges mismatch got=%v want=%v", got, wantE)
	}
}

func TestVertexNeighbors_IsolatedVertex(t *testing.T) {
	g := graphlib.NewSoAGraph(nil)
	g.AddVertex("A", "A", "server", true)
	g.AddVertex("B", "B", "server", true) // Another vertex to ensure graph isn't empty

	sg, err := g.VertexNeighbors("A")
	if err != nil {
		t.Fatalf("VertexNeighbors(\"A\") returned error: %v", err)
	}

	if len(sg.Vertices) != 1 {
		t.Fatalf("VertexNeighbors(\"A\") expected 1 vertex, got %d", len(sg.Vertices))
	}
	if sg.Vertices[0].Key != "A" {
		t.Fatalf("VertexNeighbors(\"A\") expected vertex A, got %s", sg.Vertices[0].Key)
	}
	if len(sg.Edges) != 0 {
		t.Fatalf("VertexNeighbors(\"A\") expected 0 edges, got %d", len(sg.Edges))
	}
}

func TestVertexNeighbors_OnEmptyGraph(t *testing.T) {
	g := graphlib.NewSoAGraph(nil)
	_, err := g.VertexNeighbors("A")
	if err == nil {
		t.Fatal("VertexNeighbors(\"A\") on empty graph expected error, got nil")
	}
	var nf graphlib.VertexNotFoundErr
	if !errors.As(err, &nf) {
		t.Fatalf("VertexNeighbors(\"A\") on empty graph expected VertexNotFoundErr, got %T", err)
	}
	if nf.Key != "A" {
		t.Fatalf("VertexNeighbors(\"A\") on empty graph expected VertexNotFoundErr for key A, got %s", nf.Key)
	}
}

func TestVertexDependents_IsolatedVertex(t *testing.T) {
	g := graphlib.NewSoAGraph(nil)
	g.AddVertex("A", "A", "server", true)
	g.AddVertex("B", "B", "server", true)

	// Test with all = true
	sg, err := g.VertexDependents("A", true)
	if err != nil {
		t.Fatalf("VertexDependents(\"A\", true) returned error: %v", err)
	}
	if len(sg.Vertices) != 1 || sg.Vertices[0].Key != "A" {
		t.Fatalf("VertexDependents(\"A\", true) expected subgraph with vertex A, got %v", sg.Vertices)
	}
	if len(sg.Edges) != 0 {
		t.Fatalf("VertexDependents(\"A\", true) expected 0 edges, got %d", len(sg.Edges))
	}

	// Test with all = false
	sg, err = g.VertexDependents("A", false)
	if err != nil {
		t.Fatalf("VertexDependents(\"A\", false) returned error: %v", err)
	}
	if len(sg.Vertices) != 1 || sg.Vertices[0].Key != "A" {
		t.Fatalf("VertexDependents(\"A\", false) expected subgraph with vertex A, got %v", sg.Vertices)
	}
	if len(sg.Edges) != 0 {
		t.Fatalf("VertexDependents(\"A\", false) expected 0 edges, got %d", len(sg.Edges))
	}
}

func TestVertexDependents_OnEmptyGraph(t *testing.T) {
	g := graphlib.NewSoAGraph(nil)
	_, err := g.VertexDependents("A", false)
	if err == nil {
		t.Fatal("VertexDependents(\"A\", false) on empty graph expected error, got nil")
	}
	var nf graphlib.VertexNotFoundErr
	if !errors.As(err, &nf) {
		t.Fatalf("VertexDependents(\"A\", false) on empty graph expected VertexNotFoundErr, got %T", err)
	}
	if nf.Key != "A" {
		t.Fatalf("VertexDependents(\"A\", false) on empty graph expected VertexNotFoundErr for key A, got %s", nf.Key)
	}
}

func TestPath_Self(t *testing.T) {
	g := graphlib.NewSoAGraph(nil)
	g.AddVertex("A", "A", "server", true)

	sg, err := g.Path("A", "A")
	if err != nil {
		t.Fatalf("g.Path(\"A\", \"A\") returned error: %v", err)
	}

	if len(sg.Vertices) != 1 {
		t.Fatalf("g.Path(\"A\", \"A\") expected 1 vertex, got %d", len(sg.Vertices))
	}

	if sg.Vertices[0].Key != "A" {
		t.Fatalf("g.Path(\"A\", \"A\") expected vertex with Key \"A\", got %s", sg.Vertices[0].Key)
	}

	if len(sg.Edges) != 0 {
		t.Fatalf("g.Path(\"A\", \"A\") expected 0 edges, got %d", len(sg.Edges))
	}
}

func TestVertexNeighbors_NotFound(t *testing.T) {
	g := buildGraph()

	_, err := g.VertexNeighbors("X")
	var nf graphlib.VertexNotFoundErr
	if !errors.As(err, &nf) || nf.Key != "X" {
		t.Fatalf("expected VertexNotFoundErr(X), got %v", err)
	}
}

func TestVertexDependents(t *testing.T) {

	g := buildGraph()

	sg, err := g.VertexDependents("D", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	wantV := map[string]bool{"C": true, "D": true, "F": true}
	if got := setVerts(sg.Vertices); len(got) != len(wantV) {
		t.Fatalf("vertices mismatch got=%v want=%v", got, wantV)
	}
	wantE := map[e]bool{{"C", "D"}: true, {"F", "D"}: true}
	if got := setEdges(sg.Edges); len(got) != len(wantE) {
		t.Fatalf("edges mismatch got=%v want=%v", got, wantE)
	}
}

func TestVertexDependents_All(t *testing.T) {
	g := buildGraph()

	sg, err := g.VertexDependents("D", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	wantV := map[string]bool{"A": true, "C": true, "D": true, "F": true}
	if got := setVerts(sg.Vertices); len(got) != len(wantV) {
		t.Fatalf("vertices mismatch got=%v want=%v", got, wantV)
	}
	wantE := map[e]bool{
		{"A", "C"}: true,
		{"C", "D"}: true,
		{"F", "D"}: true,
	}
	if got := setEdges(sg.Edges); len(got) != len(wantE) {
		t.Fatalf("edges mismatch got=%v want=%v", got, wantE)
	}
}

func TestVertexDependents_NotFound(t *testing.T) {
	g := buildGraph()

	_, err := g.VertexDependents("X", false)
	var nf graphlib.VertexNotFoundErr
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

	sg, err := g.VertexDependencies("A", true)
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

func TestVertexDependencies_NotFound(t *testing.T) {
	g := buildGraph()

	_, err := g.VertexDependencies("X", false)
	var nf graphlib.VertexNotFoundErr
	if !errors.As(err, &nf) || nf.Key != "X" {
		t.Fatalf("expected VertexNotFoundErr(X), got %v", err)
	}

}

func TestVertexDependencies_IsolatedVertex(t *testing.T) {
	g := graphlib.NewSoAGraph(nil)
	g.AddVertex("A", "A", "server", true)
	g.AddVertex("B", "B", "server", true)

	// Test with all = true
	sg, err := g.VertexDependencies("A", true)
	if err != nil {
		t.Fatalf("VertexDependencies(\"A\", true) returned error: %v", err)
	}
	if len(sg.Vertices) != 1 || sg.Vertices[0].Key != "A" {
		t.Fatalf("VertexDependencies(\"A\", true) expected subgraph with vertex A, got %v", sg.Vertices)
	}
	if len(sg.Edges) != 0 {
		t.Fatalf("VertexDependencies(\"A\", true) expected 0 edges, got %d", len(sg.Edges))
	}

	// Test with all = false
	sg, err = g.VertexDependencies("A", false)
	if err != nil {
		t.Fatalf("VertexDependencies(\"A\", false) returned error: %v", err)
	}
	if len(sg.Vertices) != 1 || sg.Vertices[0].Key != "A" {
		t.Fatalf("VertexDependencies(\"A\", false) expected subgraph with vertex A, got %v", sg.Vertices)
	}
	if len(sg.Edges) != 0 {
		t.Fatalf("VertexDependencies(\"A\", false) expected 0 edges, got %d", len(sg.Edges))
	}
}

func TestVertexDependencies_OnEmptyGraph(t *testing.T) {
	g := graphlib.NewSoAGraph(nil)
	_, err := g.VertexDependencies("A", false)
	if err == nil {
		t.Fatal("VertexDependencies(\"A\", false) on empty graph expected error, got nil")
	}
	var nf graphlib.VertexNotFoundErr
	if !errors.As(err, &nf) {
		t.Fatalf("VertexDependencies(\"A\", false) on empty graph expected VertexNotFoundErr, got %T", err)
	}
	if nf.Key != "A" {
		t.Fatalf("VertexDependencies(\"A\", false) on empty graph expected VertexNotFoundErr for key A, got %s", nf.Key)
	}
}

func TestPath(t *testing.T) {
	g := graphlib.NewSoAGraph(nil)

	for _, k := range []string{"A", "B", "C", "D"} {
		g.AddVertex(k, k, "server", true)
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

func TestPath_OnEmptyGraph(t *testing.T) {
	g := graphlib.NewSoAGraph(nil)
	_, err := g.Path("A", "B")
	if err == nil {
		t.Fatal("Path(\"A\", \"B\") on empty graph expected error, got nil")
	}
	var nf graphlib.VertexNotFoundErr
	if !errors.As(err, &nf) {
		t.Fatalf("Path(\"A\", \"B\") on empty graph expected VertexNotFoundErr, got %T", err)
	}
	if nf.Key != "A" { // Path checks source first
		t.Fatalf("Path(\"A\", \"B\") on empty graph expected VertexNotFoundErr for key A, got %s", nf.Key)
	}
}

func TestPath_Disconnected(t *testing.T) {
	g := graphlib.NewSoAGraph(nil)
	g.AddVertex("A", "A", "server", true)
	g.AddVertex("B", "B", "server", true)
	g.AddVertex("C", "C", "server", true)
	g.AddVertex("D", "D", "server", true)
	g.AddEdge("A", "B")

	_, err := g.Path("A", "C")
	if err == nil {
		t.Fatal("Path(\"A\", \"C\") for disconnected vertices expected error, got nil")
	}
	var pe graphlib.VertexPathErr
	if !errors.As(err, &pe) {
		t.Fatalf("Path(\"A\", \"C\") for disconnected vertices expected VertexPathErr, got %T", err)
	}
	if pe.Src != "A" || pe.Dst != "C" {
		t.Fatalf("Path(\"A\", \"C\") expected VertexPathErr for A to C, got Src=%s, Dst=%s", pe.Src, pe.Dst)
	}
}

func TestPath_VertexPathErr(t *testing.T) {
	g := graphlib.NewSoAGraph(nil)

	for _, k := range []string{"A", "B", "C"} {
		g.AddVertex(k, k, "server", true)
	}
	g.AddEdge("A", "C")
	g.AddEdge("B", "C")

	_, err := g.Path("A", "B")
	var pe graphlib.VertexPathErr
	if !errors.As(err, &pe) {
		t.Fatalf("expected VertexPathErr, got %v", err)
	}
}

func TestPath_NotFound(t *testing.T) {
	g := graphlib.NewSoAGraph(nil)

	for _, k := range []string{"A", "B", "C"} {
		g.AddVertex(k, k, "server", true)
	}
	g.AddEdge("A", "C")
	g.AddEdge("B", "C")

	_, err := g.Path("A", "X")
	var nf graphlib.VertexNotFoundErr
	if !errors.As(err, &nf) || nf.Key != "X" {
		t.Fatalf("expected VertexNotFoundErr(X), got %v", err)
	}

	_, err = g.Path("X", "A")
	if !errors.As(err, &nf) || nf.Key != "X" {
		t.Fatalf("expected VertexNotFoundErr(X), got %v", err)
	}
}

/*** helpers ***************************************************************/

// transforma slice de vértices / arestas em conjunto para comparação
func setVerts(vs []graphlib.Vertex) map[string]bool {
	m := make(map[string]bool, len(vs))
	for _, v := range vs {
		m[v.Key] = true
	}
	return m
}

type e struct{ src, dst string }

func setEdges(es []graphlib.Edge) map[e]bool {
	m := make(map[e]bool, len(es))
	for _, ed := range es {
		m[e{ed.Source, ed.Target}] = true
	}
	return m
}

/*
	digraph G {
	    A;
	    B;
	    C;
	    D;
	    E;
	    F;

		    A -> B;
		    A -> C;
		    C -> D;
		    D -> E;
		    F -> D;
		}
*/
func buildGraph() *graphlib.Graph {
	g := graphlib.NewSoAGraph(nil)
	for _, k := range []string{"A", "B", "C", "D", "E", "F"} {
		g.AddVertex(k, "", "server", true)
	}
	g.AddEdge("A", "B")
	g.AddEdge("A", "C")
	g.AddEdge("C", "D")
	g.AddEdge("D", "E")
	g.AddEdge("F", "D")
	return g
}
