// Package data defines input data structures and helpers for reading and
// preprocessing TSP instances used in this lab.
package data

// Node represents a single point in the plane with an associated cost.
type Node struct {
	X, Y, Cost int
}
