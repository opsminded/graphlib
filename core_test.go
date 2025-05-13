package graphlib

import (
	"errors"
	"testing"
)

func TestCore_BasicsOfEdge(t *testing.T) {
	g := new()

	a := g.newVertex(1, "A")
	b := g.newVertex(2, "B")
	c := g.newVertex(3, "C")

	e, err := g.newEdge(4, "AB", a, b)
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}

	if e.label != "AB" {
		t.Fatalf("Expected label to be 'AB', but got '%s'", e.label)
	}

	g.newEdge(5, "BC", b, c)

	{
		e1, _ := g.edges.find(a, b)
		e2, err := g.newEdge(4, "AB", a, b)
		if err != nil {
			t.Fatalf("Expected no error, but got %v", err)
		}

		if e1 != e2 {
			t.Fatalf("Expected edge to be the same, but got different edges")
		}
	}

	{
		_, err := g.newEdge(4, "AB", nil, b)
		if !errors.As(err, &VertexNilError{}) {
			t.Fatalf("Expected nil error but got %v", err)
		}
	}

	{
		_, err := g.newEdge(4, "AB", a, nil)
		if !errors.As(err, &VertexNilError{}) {
			t.Fatalf("Expected nil error but got %v", err)
		}
	}

	{
		_, err := g.newEdge(5, "BA", b, a)
		if !errors.As(err, &BidirectionalEdgeError{}) {
			t.Fatalf("Expected bidirectional edge error but got %v", err)
		}
	}

	{
		_, err := g.newEdge(6, "CA", c, a)
		if !errors.As(err, &CycleError{}) {
			t.Fatalf("Expected cycle error but got %v", err)
		}
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
	if a.healthy != true {
		t.Fatalf("Expected health to be true, but got false")
	}

	a.healthy = false
	if a.healthy != false {
		t.Fatalf("Expected health to be false, but got true")
	}

	b := g.newVertex(1, "A")
	if a != b {
		t.Fatalf("Expected vertex to be the same, but got different vertices")
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

func TestCore_AddNilVertexEdge(t *testing.T) {
	g := new()
	v := g.newVertex(1, "A")

	{
		_, err := g.newEdge(5, "AB", v, nil)
		if !errors.As(err, &VertexNilError{}) {
			t.Fatal("Expected nil error but got", err)
		}
	}

	{
		_, err := g.newEdge(5, "AB", nil, v)
		if !errors.As(err, &VertexNilError{}) {
			t.Fatal("Expected nil error but got", err)
		}
	}

	{
		_, err := g.newEdge(5, "AB", nil, nil)
		if !errors.As(err, &VertexNilError{}) {
			t.Fatal("Expected nil error but got", err)
		}
	}
}

func TestCore_EdgeCycleDetectionError(t *testing.T) {
	g := new()
	a := g.newVertex(1, "A")
	b := g.newVertex(2, "B")
	c := g.newVertex(3, "C")

	if _, err := g.newEdge(4, "AB", a, b); err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}

	if _, err := g.newEdge(5, "BC", b, c); err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}

	_, err := g.newEdge(6, "CA", c, a)
	if !errors.As(err, &CycleError{}) {
		t.Fatalf("Expected cycle error but got %v", err)
	}
}

func TestCore_BidirectionalEdge(t *testing.T) {
	g := new()
	a := g.newVertex(1, "A")
	b := g.newVertex(2, "B")

	if _, err := g.newEdge(4, "AB", a, b); err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}

	_, err := g.newEdge(5, "BA", b, a)
	if !errors.As(err, &BidirectionalEdgeError{}) {
		t.Fatalf("Expected bidirectional edge error but got %v", err)
	}
}
