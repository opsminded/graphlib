package graphlib

import "fmt"

type VertexNotFoundError struct {
	Label string
}

func (e VertexNotFoundError) Error() string {
	return fmt.Sprintf("vertex with label %s not found", e.Label)
}

type VertexNilError struct {
	position string
}

func (e VertexNilError) Error() string {
	return fmt.Sprintf("vertex in the %s is nil", e.position)
}

type EdgeNotFoundError struct {
	LabelA string
	LabelB string
}

func (e EdgeNotFoundError) Error() string {
	return fmt.Sprintf("edge between %s and %s not found", e.LabelA, e.LabelB)
}

type BidirectionalEdgeError struct {
	LabelA string
	LabelB string
}

func (e BidirectionalEdgeError) Error() string {
	return fmt.Sprintf("bidirectional edges are not allowed between %s and %s", e.LabelA, e.LabelB)
}

type CycleError struct {
	LabelA string
	LabelB string
}

func (e CycleError) Error() string {
	return fmt.Sprintf("adding this edge would create a cycle between %s and %s", e.LabelA, e.LabelB)
}
