package algorithms

import (
	"math/rand"

	"github.com/czajkowskis/evolutionary_computation/01_labs/greedy_heuristics/pkg/data"
)

func RandomSolution(nodes []data.Node, distanceMatrix [][]int, numSolutions int) []Solution {
	n := len(nodes)
	k := (n + 1) / 2 // 50% rounded up
	var solutions []Solution

	for i := 0; i < numSolutions; i++ {
		path := rand.Perm(n)[:k]

		totalDistance := 0
		for j := 0; j < k; j++ {
			next := (j + 1) % k
			totalDistance += distanceMatrix[path[j]][path[next]]
		}

		totalCost := 0
		for _, idx := range path {
			totalCost += nodes[idx].Cost
		}

		objective := totalDistance + totalCost
		solutions = append(solutions, Solution{path, objective})
	}
	return solutions
}
