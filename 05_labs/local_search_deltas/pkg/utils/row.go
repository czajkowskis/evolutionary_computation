// Package utils provides small helper types and functions for statistics,
// CSV output and filename handling.
package utils

// Row represents a single row of aggregated experiment results.
type Row struct {
	Name      string
	AvgV      float64
	MinV      int
	MaxV      int
	AvgTms    float64
	BestPath  []int
	BestValue int
}
