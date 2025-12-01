package utils

import (
	"math"

	"github.com/czajkowskis/evolutionary_computation/07_labs/large_neighborhood_search/pkg/algorithms"
)

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
