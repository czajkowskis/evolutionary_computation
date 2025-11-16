// Package algorithms contains local search methods and supporting structures
// used to solve the TSP with node costs in this lab.
package algorithms

// Solution represents a single TSP solution, including the path and its
// objective value (tour length plus node costs).
type Solution struct {
	Path      []int
	Objective int
}
