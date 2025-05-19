package graphlib

type Vertex struct {
	Key       string
	Label     string
	Class     string
	Healthy   bool
	LastCheck int64
}

type Edge struct {
	Key    string
	Source string
	Target string
}

type Subgraph struct {
	Vertices []Vertex
	Edges    []Edge
}

type Stats struct {
	TotalVertices          int
	TotalUnhealthyVertices int
	TotalEdges             int
	TotalHealthyVertices   int

	UnhealthyVertices []Vertex
}
