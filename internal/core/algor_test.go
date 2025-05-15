package core_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/opsminded/graphlib/internal/core"
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

func TestVertexNeighbors_NotFound(t *testing.T) {
	g := buildGraph()

	_, err := g.VertexNeighbors("X")
	var nf core.VertexNotFoundErr
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

func TestPath_VertexPathErr(t *testing.T) {
	g := core.NewSoAGraph()

	for _, k := range []string{"A", "B", "C"} {
		g.AddVertex(k, k, true)
	}
	g.AddEdge("A", "C")
	g.AddEdge("B", "C")

	_, err := g.Path("A", "B")
	var pe core.VertexPathErr
	if !errors.As(err, &pe) {
		t.Fatalf("expected VertexPathErr, got %v", err)
	}
}

func TestPath_NotFound(t *testing.T) {
	g := core.NewSoAGraph()

	for _, k := range []string{"A", "B", "C"} {
		g.AddVertex(k, k, true)
	}
	g.AddEdge("A", "C")
	g.AddEdge("B", "C")

	_, err := g.Path("A", "X")
	var nf core.VertexNotFoundErr
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
func buildGraph() *core.Graph {
	g := core.NewSoAGraph()
	for _, k := range []string{"A", "B", "C", "D", "E", "F"} {
		g.AddVertex(k, "", true)
	}
	g.AddEdge("A", "B")
	g.AddEdge("A", "C")
	g.AddEdge("C", "D")
	g.AddEdge("D", "E")
	g.AddEdge("F", "D")
	return g
}
