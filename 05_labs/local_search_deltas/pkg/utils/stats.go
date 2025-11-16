package utils

import (
	"math"

	"github.com/czajkowskis/evolutionary_computation/05_labs/local_search_deltas/pkg/algorithms"
)

// CalculateStatistics returns the minimum, maximum and average objective
// value across all provided solutions. For an empty slice it returns zeros.
func CalculateStatistics(solutions []algorithms.Solution) (int, int, float64) {
	if len(solutions) == 0 {
		return 0, 0, 0
	}
	minObj := math.MaxInt32
	maxObj := 0
	sum := 0
	for _, sol := range solutions {
		obj := sol.Objective
		if obj < minObj {
			minObj = obj
		}
		if obj > maxObj {
			maxObj = obj
		}
		sum += obj
	}
	avg := float64(sum) / float64(len(solutions))
	return minObj, maxObj, avg
}
