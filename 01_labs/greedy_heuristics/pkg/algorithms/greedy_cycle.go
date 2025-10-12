package algorithms

import (
	"math"
)

// GreedyCycle generates solutions using a greedy cycle algorithm.
func GreedyCycle(distanceMatrix [][]int, nodeCosts []int, startNodeIndices []int) []Solution {
	n := len(nodeCosts)
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
			bestNodeIndex := -1
			minScore := math.MaxInt32
			for nodeIndex := range unvisited {
				score := 2*distanceMatrix[startNodeIndex][nodeIndex] + nodeCosts[nodeIndex]
				if score < minScore {
					minScore = score
					bestNodeIndex = nodeIndex
				}
			}
			if bestNodeIndex != -1 {
				path = append(path, bestNodeIndex)
				delete(unvisited, bestNodeIndex)
			}
		}

		// Build the rest of the path using greedy insertion
		for len(path) < k && len(unvisited) > 0 {
			bestNodeIndex := -1
			bestPosition := -1
			minIncreaseScore := math.MaxInt32

			for nodeIndex := range unvisited {
				for i := 0; i < len(path); i++ {
					p1 := path[i]
					p2 := path[(i+1)%len(path)]
					deltaDist := distanceMatrix[p1][nodeIndex] + distanceMatrix[nodeIndex][p2] - distanceMatrix[p1][p2]
					increaseScore := deltaDist + nodeCosts[nodeIndex]
					if increaseScore < minIncreaseScore {
						minIncreaseScore = increaseScore
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
			totalCost += nodeCosts[idx]
		}

		objective := totalDistance + totalCost
		solutions = append(solutions, Solution{Path: path, Objective: objective})
	}

	return solutions
}
