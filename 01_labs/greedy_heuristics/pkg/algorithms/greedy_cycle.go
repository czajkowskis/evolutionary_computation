package algorithms

import (
	"math"

	"github.com/czajkowskis/evolutionary_computation/01_labs/greedy_heuristics/pkg/data"
)

// GreedyCycle generates solutions using a greedy cycle algorithm.
func GreedyCycle(nodes []data.Node, distanceMatrix [][]int, startNodeIndices []int) []Solution {
	n := len(nodes)
	if n == 0 {
		return nil
	}
	k := (n + 1) / 2 // 50% of nodes rounded up
	var solutions []Solution

	for _, startNodeIndex := range startNodeIndices {
		path := []int{startNodeIndex}
		unvisited := make(map[int]bool)
		for i := 0; i < n; i++ {
			if i != startNodeIndex {
				unvisited[i] = true
			}
		}

		// Second node: choose the nearest neighbor to the start node
		if len(unvisited) > 0 {
			nearestNodeIndex := -1
			minDistance := math.MaxInt32
			for nodeIndex := range unvisited {
				if distanceMatrix[startNodeIndex][nodeIndex] < minDistance {
					minDistance = distanceMatrix[startNodeIndex][nodeIndex]
					nearestNodeIndex = nodeIndex
				}
			}
			if nearestNodeIndex != -1 {
				path = append(path, nearestNodeIndex)
				delete(unvisited, nearestNodeIndex)
			}
		}

		// Build the rest of the path using greedy insertion
		for len(path) < k && len(unvisited) > 0 {
			bestNodeIndex := -1
			bestPosition := -1
			minIncrease := math.MaxInt32

			for nodeIndex := range unvisited {
				for i := 0; i < len(path); i++ {
					p1 := path[i]
					p2 := path[(i+1)%len(path)]
					increase := distanceMatrix[p1][nodeIndex] + distanceMatrix[nodeIndex][p2] - distanceMatrix[p1][p2]
					if increase < minIncrease {
						minIncrease = increase
						bestNodeIndex = nodeIndex
						bestPosition = i + 1
					}
				}
			}

			if bestNodeIndex != -1 {
				path = append(path[:bestPosition], append([]int{bestNodeIndex}, path[bestPosition:]...)...)
				delete(unvisited, bestNodeIndex)
			} else {
				break
			}
		}

		totalDistance := 0
		if len(path) > 1 {
			for i := 0; i < len(path); i++ {
				totalDistance += distanceMatrix[path[i]][path[(i+1)%len(path)]]
			}
		}

		totalCost := 0
		for _, idx := range path {
			totalCost += nodes[idx].Cost
		}

		objective := totalDistance + totalCost
		solutions = append(solutions, Solution{Path: path, Objective: objective})
	}

	return solutions
}
