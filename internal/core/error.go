package core

import "fmt"

type VertexNotFoundErr struct {
	Key string
}

func (e VertexNotFoundErr) Error() string {
	return fmt.Sprintf("vertex %q not found", e.Key)
}

type BidirectionalEdgeErr struct {
	Src, Tgt string
}

func (e BidirectionalEdgeErr) Error() string {
	return fmt.Sprintf("bidirectional edge %s ↔ %s not allowed", e.Src, e.Tgt)
}

type CycleErr struct {
	Src string
	Tgt string
}

func (e CycleErr) Error() string {
	return fmt.Sprintf("edge %s → %s would create a cycle", e.Src, e.Tgt)
}

type VertexPathErr struct {
	Src string
	Dst string
}

func (e VertexPathErr) Error() string {
	return fmt.Sprintf("no path from %s to %s", e.Src, e.Dst)
}
