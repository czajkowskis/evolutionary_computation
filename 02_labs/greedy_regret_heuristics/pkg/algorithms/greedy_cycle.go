package algorithms

import (
	"math"
)

// GreedyCycle generates solutions using a greedy cycle algorithm.
func GreedyCycleWeightedTwoRegret(distanceMatrix [][]int, nodeCosts []int, startNodeIndices []int, regretWeight float64, objectiveWeight float64) []Solution {
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

		// Build the rest of the path using greedy insertion with regret
		for len(path) < k && len(unvisited) > 0 {
			bestScore := math.Inf(-1)
			bestNodeIndex := -1
			bestPosition := -1

			// Find maximum possible regret and objective for normalization
			maxPossibleRegret := 0.0
			maxPossibleObjective := 0.0

			// First pass to find normalization values
			for nodeIndex := range unvisited {
				bestLocalCost := math.MaxInt32
				secondBestLocalCost := math.MaxInt32

				for i := 0; i < len(path); i++ {
					p1 := path[i]
					p2 := path[(i+1)%len(path)]
					deltaDist := distanceMatrix[p1][nodeIndex] + distanceMatrix[nodeIndex][p2] - distanceMatrix[p1][p2]
					insertionCost := deltaDist + nodeCosts[nodeIndex]

					if insertionCost < bestLocalCost {
						secondBestLocalCost = bestLocalCost
						bestLocalCost = insertionCost
					} else if insertionCost < secondBestLocalCost {
						secondBestLocalCost = insertionCost
					}
				}

				regret := secondBestLocalCost - bestLocalCost
				if float64(regret) > maxPossibleRegret {
					maxPossibleRegret = float64(regret)
				}
				if float64(bestLocalCost) > maxPossibleObjective {
					maxPossibleObjective = float64(bestLocalCost)
				}
			}

			// Avoid division by zero
			if maxPossibleRegret == 0 {
				maxPossibleRegret = 1
			}
			if maxPossibleObjective == 0 {
				maxPossibleObjective = 1
			}

			// Second pass to find best node using normalized values
			for nodeIndex := range unvisited {
				bestLocalCost := math.MaxInt32
				secondBestLocalCost := math.MaxInt32
				bestPos := -1

				for i := 0; i < len(path); i++ {
					p1 := path[i]
					p2 := path[(i+1)%len(path)]
					deltaDist := distanceMatrix[p1][nodeIndex] + distanceMatrix[nodeIndex][p2] - distanceMatrix[p1][p2]
					insertionCost := deltaDist + nodeCosts[nodeIndex]

					if insertionCost < bestLocalCost {
						secondBestLocalCost = bestLocalCost
						bestLocalCost = insertionCost
						bestPos = i + 1
					} else if insertionCost < secondBestLocalCost {
						secondBestLocalCost = insertionCost
					}
				}

				// Normalize values
				normalizedRegret := float64(secondBestLocalCost-bestLocalCost) / maxPossibleRegret
				normalizedObjective := float64(bestLocalCost) / maxPossibleObjective

				// Calculate weighted score
				score := (regretWeight * normalizedRegret) - (objectiveWeight * normalizedObjective)

				if score > bestScore {
					bestScore = score
					bestNodeIndex = nodeIndex
					bestPosition = bestPos
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
