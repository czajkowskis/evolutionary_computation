package algorithms

import (
	"math/rand"
)

// RandomSolution generates random solutions starting from a given set of start node indices.
func RandomSolution(distanceMatrix [][]int, nodeCosts []int, startNodeIndices []int) []Solution {
	n := len(nodeCosts)
	if n == 0 {
		return nil
	}
	k := (n + 1) / 2 // 50% rounded up
	var solutions []Solution

	for _, startNodeIndex := range startNodeIndices {
		path := []int{startNodeIndex}
		unvisited := make(map[int]bool)
		for i := 0; i < n; i++ {
			if i != startNodeIndex {
				unvisited[i] = true
			}
		}

		for len(path) < k && len(unvisited) > 0 {
			// Create a slice of unvisited node keys to enable random selection.
			unvisitedKeys := make([]int, 0, len(unvisited))
			for nodeIndex := range unvisited {
				unvisitedKeys = append(unvisitedKeys, nodeIndex)
			}

			// Explicitly select a random node from the list of unvisited nodes.
			randIdx := rand.Intn(len(unvisitedKeys))
			randomNodeIndex := unvisitedKeys[randIdx]

			// Add the randomly selected node to the path and remove it from the unvisited map.
			path = append(path, randomNodeIndex)
			delete(unvisited, randomNodeIndex)
		}
		// Shuffle the rest of the path to ensure randomness
		rand.Shuffle(len(path)-1, func(i, j int) {
			path[i+1], path[j+1] = path[j+1], path[i+1]
		})

		totalDistance := 0
		if k > 1 {
			for j := 0; j < k; j++ {
				next := (j + 1) % k
				totalDistance += distanceMatrix[path[j]][path[next]]
			}
		}

		totalCost := 0
		for _, idx := range path {
			totalCost += nodeCosts[idx]
		}

		objective := totalDistance + totalCost
		solutions = append(solutions, Solution{path, objective})
	}
	return solutions
}
