package graphlib

import (
	"errors"
	"fmt"
	"testing"
)

func TestGraphBasics(t *testing.T) {
	g := NewSoAGraph(nil)

	g.AddVertex("A", "A", "server", true)
	g.AddVertex("B", "B", "server", true)

	err := g.AddEdge("A", "B")
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}

	v, err := g.GetVertex("A")
	if err != nil {
		t.Fatalf("Expected to find vertex A, but got error %v", err)
	}
	if v.Key != "A" {
		t.Fatalf("Expected vertex A, but got %v", v.Key)
	}
	if v.Label != "A" {
		t.Fatalf("Expected label A, but got %v", v.Label)
	}
	if v.Healthy != true {
		t.Fatalf("Expected health true, but got %v", v.Healthy)
	}
}

func TestVertexClasses(t *testing.T) {
	g := NewSoAGraph(nil)

	g.AddVertex("A", "A", "server", true)
	g.AddVertex("B", "B", "server", true)
	g.AddVertex("C", "C", "client", true)

	a, _ := g.GetVertex("A")
	b, _ := g.GetVertex("B")
	c, _ := g.GetVertex("C")

	fmt.Println(a, b, c)

	if a.Class != "server" {
		t.Fatalf("Expected class 'server' for vertex A, but got '%s'", a.Class)
	}
	if b.Class != "server" {
		t.Fatalf("Expected class 'server' for vertex B, but got '%s'", b.Class)
	}
	if c.Class != "client" {
		t.Fatalf("Expected class 'client' for vertex C, but got '%s'", c.Class)
	}
}

func TestGraphAddVertex(t *testing.T) {
	g := NewSoAGraph(nil)

	g.AddVertex("A", "A", "server", true)
	g.AddVertex("A", "A", "server", true)

	a1, _ := g.GetVertex("A")
	a2, _ := g.GetVertex("A")

	if a1 != a2 {
		t.Fatalf("Expected to find vertex A, but got different vertices")
	}

	a3, err := g.GetVertex("B")
	if err == nil {
		t.Fatalf("Expected error for non-existent vertex B, but got %v", a3)
	}
}

func TestGraphAddEdge(t *testing.T) {
	g := NewSoAGraph(nil)

	g.AddVertex("A", "A", "server", true)
	g.AddVertex("B", "B", "server", true)

	err := g.AddEdge("A", "B")
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}

	err = g.AddEdge("A", "C")
	if err == nil {
		t.Fatal("Expected error for non-existent target vertex, but got none")
	}

	err = g.AddEdge("C", "B")
	if err == nil {
		t.Fatal("Expected error for non-existent source vertex, but got none")
	}

	// Test adding the same edge again
	err = g.AddEdge("A", "B")
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}

	err = g.AddEdge("B", "A")
	if err == nil {
		t.Fatalf("Expected bidirectional edge error, but got none")
	}
}

func TestGraphCycleDetection(t *testing.T) {
	g := NewSoAGraph(nil)
	g.AddVertex("A", "A", "server", true)
	g.AddVertex("B", "B", "server", true)
	g.AddVertex("C", "C", "server", true)

	err := g.AddEdge("A", "B")
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}

	err = g.AddEdge("B", "C")
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}

	err = g.AddEdge("C", "A")
	if err == nil {
		t.Fatal("Expected cycle detection error, but got none")
	}

	err = g.AddEdge("A", "A")
	if err == nil {
		t.Fatal("Expected self-loop error, but got none")
	}
}

func TestGraphLosangleCycle(t *testing.T) {
	g := NewSoAGraph(nil)

	g.AddVertex("A", "A", "server", true)
	g.AddVertex("B", "B", "server", true)
	g.AddVertex("C", "C", "server", true)
	g.AddVertex("D", "D", "server", true)
	g.AddVertex("E", "E", "server", true)
	g.AddVertex("F", "F", "server", true)
	g.AddVertex("G", "G", "server", true)

	g.AddEdge("B", "C")

	g.AddEdge("C", "D")
	g.AddEdge("C", "E")

	g.AddEdge("D", "F")
	g.AddEdge("E", "F")
	g.AddEdge("F", "G")

	g.AddEdge("A", "B")

	err := g.AddEdge("G", "A")
	if err == nil {
		t.Fatal("Expected cycle detection error, but got none")
	}
}

func TestGraphFind(t *testing.T) {
	g := NewSoAGraph(nil)

	g.AddVertex("A", "A", "server", true)

	_, err := g.GetVertex("A")
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}

	_, err = g.GetVertex("B")
	if err == nil {
		t.Fatal("Expected error for non-existent vertex B, but got none")
	}

	v1, err := g.GetVertex("A")
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}
	if v1.Key != "A" {
		t.Fatalf("Expected vertex A, but got %v", v1.Key)
	}
	if v1.Label != "A" {
		t.Fatalf("Expected label A, but got %v", v1.Label)
	}
	if v1.Healthy != true {
		t.Fatalf("Expected health true, but got %v", v1.Healthy)
	}

	_, err = g.GetVertex("B")
	if err == nil {
		t.Fatal("Expected error for non-existent vertex B, but got none")
	}
}

func TestGraphStats(t *testing.T) {
	g := NewSoAGraph(nil)

	g.AddVertex("A", "A", "server", true)
	g.AddVertex("B", "B", "server", true)
	g.AddVertex("C", "C", "server", false)

	g.AddEdge("A", "B")
	g.AddEdge("B", "C")

	stats := g.Stats()

	if stats.TotalVertices != 3 {
		t.Fatalf("Expected 3 vertices, but got %d", stats.TotalVertices)
	}
	if stats.TotalHealthyVertices != 2 {
		t.Fatalf("Expected 2 healthy vertices, but got %d", stats.TotalHealthyVertices)
	}
	if stats.TotalUnhealthyVertices != 1 {
		t.Fatalf("Expected 1 unhealthy vertex, but got %d", stats.TotalUnhealthyVertices)
	}
	if stats.TotalEdges != 2 {
		t.Fatalf("Expected 2 edges, but got %d", stats.TotalEdges)
	}
}

func TestSetVertexHealth(t *testing.T) {
	g := NewSoAGraph(nil)

	g.AddVertex("A", "A", "server", true)
	g.AddVertex("B", "B", "server", true)

	v1, _ := g.GetVertex("A")
	v2, _ := g.GetVertex("B")

	if (!v1.Healthy) || (!v2.Healthy) {
		t.Fatal("Expected vertices to be healthy, but got unhealthy")
	}

	g.AddEdge("A", "B")

	g.SetVertexHealth("B", false)

	v1, _ = g.GetVertex("A")
	v2, _ = g.GetVertex("B")

	if v1.Healthy || v2.Healthy {
		t.Fatal("Expected vertices to be unhealthy, but got healthy")
	}
}

func TestSetVertexHealth_NotFound(t *testing.T) {
	g := NewSoAGraph(nil)

	err := g.SetVertexHealth("A", false)
	if err == nil {
		t.Fatal("Expected error for non-existent vertex A, but got none")
	}

	var vErr VertexNotFoundErr
	if !errors.As(err, &vErr) {
		t.Fatalf("expected VertexNotFoundErr, got %v", err)
	}

	want := fmt.Sprintf("vertex %q not found", "A")
	if err.Error() != want {
		t.Fatalf("expected vertex not found error message, got %v", err.Error())
	}
}

func TestGraphClearHealthyStatus(t *testing.T) {
	g := NewSoAGraph(nil)

	g.AddVertex("A", "A", "server", false)
	g.AddVertex("B", "B", "server", false)

	v1, _ := g.GetVertex("A")
	v2, _ := g.GetVertex("B")

	if v1.Healthy || v2.Healthy {
		t.Fatal("Expected vertices to be unhealthy, but got healthy")
	}

	g.ClearHealthyStatus()

	v1, _ = g.GetVertex("A")
	v2, _ = g.GetVertex("B")

	if (!v1.Healthy) || (!v2.Healthy) {
		t.Fatal("Expected vertices to be healthy, but got healthy")
	}
}

func TestVertexNotFoundErr(t *testing.T) {
	g := NewSoAGraph(nil)

	_, err := g.GetVertex("A")
	if err == nil {
		t.Fatal("Expected error for non-existent vertex A, but got none")
	}

	var vErr VertexNotFoundErr
	if !errors.As(err, &vErr) {
		t.Fatalf("expected VertexNotFoundErr, got %v", err)
	}

	want := fmt.Sprintf("vertex %q not found", "A")
	if err.Error() != want {
		t.Fatalf("expected vertex not found error message, got %v", err.Error())
	}

}

func TestBidirectionalEdgeErr(t *testing.T) {
	g := NewSoAGraph(nil)

	g.AddVertex("A", "A", "server", true)
	g.AddVertex("B", "B", "server", true)

	err := g.AddEdge("A", "B")
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}

	err = g.AddEdge("B", "A")

	var bidirectionalErr BidirectionalEdgeErr
	if !errors.As(err, &bidirectionalErr) {
		t.Fatalf("expected BidirectionalEdgeErr, got %v", err)
	}

	want := fmt.Sprintf("bidirectional edge %s ↔ %s not allowed", "B", "A")
	if err.Error() != want {
		t.Fatalf("expected bidirectional edge error message, got %v", err.Error())
	}
}
func TestVertexCycleErr(t *testing.T) {
	g := NewSoAGraph(nil)

	g.AddVertex("A", "A", "server", true)
	g.AddVertex("B", "B", "server", true)
	g.AddVertex("C", "C", "server", true)
	g.AddEdge("A", "B")
	g.AddEdge("B", "C")

	err := g.AddEdge("C", "A")

	var cycleErr CycleErr
	if !errors.As(err, &cycleErr) {
		t.Fatalf("expected CycleError, got %v", err)
	}

	want := fmt.Sprintf("edge %s → %s would create a cycle", "C", "A")
	if err.Error() != want {
		t.Fatalf("expected cycle error message, got %v", err.Error())
	}
}
