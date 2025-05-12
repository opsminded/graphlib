package graphlib

import (
	"testing"
)

func TestCore_BasicsOfEdge(t *testing.T) {
	g := new()
	a := g.newVertex(1, "A")
	b := g.newVertex(2, "B")
	e := g.newEdge(4, "AB", a, b)

	if e.label != "AB" {
		t.Fatalf("Expected label to be 'AB', but got '%s'", e.label)
	}
}

func TestCore_EdgesMapLength(t *testing.T) {
	g := new()
	a := g.newVertex(1, "A")
	b := g.newVertex(2, "B")
	c := g.newVertex(3, "C")
	d := g.newVertex(4, "D")

	g.newEdge(4, "AB", a, b)
	g.newEdge(5, "BC", b, c)
	g.newEdge(6, "CD", c, d)

	if g.edges.len() != 3 {
		t.Fatalf("Expected edges map length to be 3, but got %d", len(g.edges))
	}
}

func TestCore_BasicsOfVertex(t *testing.T) {
	g := new()
	a := g.newVertex(1, "A")
	if a.label != "A" {
		t.Fatalf("Expected label to be 'A', but got '%s'", a.label)
	}
	if a.health != true {
		t.Fatalf("Expected health to be true, but got false")
	}

	a.health = false
	if a.health != false {
		t.Fatalf("Expected health to be false, but got true")
	}
}

func TestCore_VertexListLength(t *testing.T) {
	g := new()
	g.newVertex(1, "A")
	g.newVertex(2, "B")
	g.newVertex(3, "C")
	g.newVertex(4, "D")
	g.newVertex(4, "D")
	g.newVertex(4, "D")

	if g.vertices.len() != 4 {
		t.Fatalf("Expected vertices map length to be 4, but got %d", len(g.vertices))
	}
}

func TestCore_BasicAddVertex(t *testing.T) {
	g := new()
	a := g.newVertex(1, "A")
	b := g.newVertex(2, "A")
	if a != b {
		t.Fatalf("Expected vertex to be the same, but got different vertices")
	}
}

func TestCore_BasicAddEdge(t *testing.T) {
	g := new()
	a := g.newVertex(1, "A")
	b := g.newVertex(2, "B")
	c := g.newVertex(3, "C")
	d := g.newVertex(4, "D")

	e := g.newEdge(4, "AB", a, b)
	if e == nil {
		t.Fatalf("Expected edge to be created, but got nil")
	}

	f := g.newEdge(5, "AB", a, b)
	if f != e {
		t.Fatalf("Expected edge to be the same, but got different edges")
	}

	g.newEdge(6, "BC", b, c)
	g.newEdge(7, "CD", c, d)
}

func TestCore_AddNilVertexEdge(t *testing.T) {
	g := new()
	v := g.newVertex(1, "A")

	// test if newEdge panics
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("Expected panic, but did not panic")
		}
	}()
	g.newEdge(5, "AB", v, nil)
	g.newEdge(5, "AB", nil, v)
	g.newEdge(5, "AB", nil, nil)
}

func TestCore_EdgeCycleDetectionPanic(t *testing.T) {
	g := new()
	a := g.newVertex(1, "A")
	b := g.newVertex(2, "B")
	c := g.newVertex(3, "C")

	g.newEdge(4, "AB", a, b)
	g.newEdge(5, "BC", b, c)

	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("Expected panic, but did not panic")
		}
	}()
	g.newEdge(6, "CA", c, a)
}

func TestCore_BidirectionalEdge(t *testing.T) {
	g := new()
	a := g.newVertex(1, "A")
	b := g.newVertex(2, "B")

	g.newEdge(4, "AB", a, b)

	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("Expected panic, but did not panic")
		}
	}()
	g.newEdge(5, "BA", b, a)
}
