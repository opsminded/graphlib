package graphlib_test

import (
	"testing"

	"github.com/opsminded/graphlib"
)

func TestGraphBasics(t *testing.T) {
	g := graphlib.NewGraph()

	s1, err := g.NewVertex("SERVER01", "SERVER")
	if err != nil {
		t.Fatalf("NewVertex returned a unexpected error: %v", err)
	}

	s2, err := g.NewVertex("SERVER02", "SERVER")
	if err != nil {
		t.Fatalf("NewVertex returned a unexpected error: %v", err)
	}

	s3, err := g.NewVertex("SERVER03", "SERVER")
	if err != nil {
		t.Fatalf("NewVertex returned a unexpected error: %v", err)
	}

	e1, err := g.NewEdge("CONEXAO", "SERVER_CONN", s1, s2)

	if err != nil {
		t.Fatalf("failed to create edge: %v", err)
	}

	_, err = g.NewEdge("CONEXAO", "SERVER_CONN", s2, s3)
	if err != nil {
		t.Fatalf("failed to create edge: %v", err)
	}

	// add the same vertex
	s4, err := g.NewVertex("SERVER01", "SERVER")
	if err != nil {
		t.Fatalf("NewVertex returned a unexpected error: %v", err)
	}

	if s1 != s4 {
		t.Fatal("NewVertex called again with the same label should return the same vertex")
	}

	// add the same edge
	e3, err := g.NewEdge("CONEXAO", "SERVER_CONN", s1, s2)
	if err != nil {
		t.Fatalf("NewEdge returned a unexpected error: %v", err)
	}

	if e1 != e3 {
		t.Fatal("NewEdge called again with the same label should return the same edge")
	}

	// Add second edge to a existing vertex
	_, err = g.NewEdge("LBLTEST", "LBLTEST", s1, s3)
	if err != nil {
		t.Fatalf("NewEdge returned a unexpected error: %v", err)
	}
}

func TestValidations(t *testing.T) {
	g := graphlib.NewGraph()

	// Vertices
	if _, err := g.NewVertex("invalid", "SERVER"); err != graphlib.ErrInvalidLabel {
		t.Fatalf("invalid label expected. Nil returned")
	}

	if _, err := g.NewVertex("OK", "invalid"); err != graphlib.ErrInvalidClassName {
		t.Fatalf("invalid class expected. Nil returned")
	}

	// Edges
	s1, err := g.NewVertex("SERVER01", "SERVER")
	if err != nil {
		t.Fatalf("NewVertex returned a unexpected error: %v", err)
	}

	s2, err := g.NewVertex("SERVER02", "SERVER")
	if err != nil {
		t.Fatalf("NewVertex returned a unexpected error: %v", err)
	}

	if _, err := g.NewEdge("invalid", "SERVER", s1, s2); err != graphlib.ErrInvalidLabel {
		t.Fatalf("Error expected")
	}

	if _, err := g.NewEdge("OK", "invalid", s1, s2); err != graphlib.ErrInvalidClassName {
		t.Fatalf("Error expected")
	}

	if _, err := g.NewEdge("OK", "OK", nil, s2); err != graphlib.ErrNilVertices {
		t.Fatalf("Error expected")
	}

	if _, err := g.NewEdge("OK", "OK", s1, nil); err != graphlib.ErrNilVertices {
		t.Fatalf("Error expected")
	}
}
