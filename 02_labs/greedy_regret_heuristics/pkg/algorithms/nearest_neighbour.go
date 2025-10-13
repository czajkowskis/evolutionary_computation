package algorithms

import (
	"math"
)

func NearestNeighborWeightedTwoRegret(distanceMatrix [][]int, nodeCosts []int, startNodeIndices []int, regretWeight float64, objectiveWeight float64) []Solution {
	n := len(nodeCosts)

	if n == 0 {
		return nil
	}
	k := (n + 1) / 2
	var solutions []Solution

	for _, startNodeIndex := range startNodeIndices {
		path := []int{startNodeIndex}
		unvisited := make(map[int]bool)
		for j := 0; j < n; j++ {
			if j != startNodeIndex {
				unvisited[j] = true
			}
		}

		for len(path) < k {
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

				for pos := 0; pos <= len(path); pos++ {
					var insertionCost int
					if len(path) == 1 {
						insertionCost = distanceMatrix[path[0]][nodeIndex] + nodeCosts[nodeIndex]
					} else if pos == 0 {
						insertionCost = distanceMatrix[nodeIndex][path[0]] +
							distanceMatrix[nodeIndex][path[len(path)-1]] +
							nodeCosts[nodeIndex]
					} else if pos == len(path) {
						insertionCost = distanceMatrix[path[len(path)-1]][nodeIndex] +
							distanceMatrix[nodeIndex][path[0]] +
							nodeCosts[nodeIndex]
					} else {
						prev := path[pos-1]
						next := path[pos]
						insertionCost = distanceMatrix[prev][nodeIndex] +
							distanceMatrix[nodeIndex][next] -
							distanceMatrix[prev][next] +
							nodeCosts[nodeIndex]
					}

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

				// Calculate insertion costs for all positions
				for pos := 0; pos <= len(path); pos++ {
					var insertionCost int
					if len(path) == 1 {
						insertionCost = distanceMatrix[path[0]][nodeIndex] + nodeCosts[nodeIndex]
					} else if pos == 0 {
						insertionCost = distanceMatrix[nodeIndex][path[0]] +
							distanceMatrix[nodeIndex][path[len(path)-1]] +
							nodeCosts[nodeIndex]
					} else if pos == len(path) {
						insertionCost = distanceMatrix[path[len(path)-1]][nodeIndex] +
							distanceMatrix[nodeIndex][path[0]] +
							nodeCosts[nodeIndex]
					} else {
						prev := path[pos-1]
						next := path[pos]
						insertionCost = distanceMatrix[prev][nodeIndex] +
							distanceMatrix[nodeIndex][next] -
							distanceMatrix[prev][next] +
							nodeCosts[nodeIndex]
					}

					if insertionCost < bestLocalCost {
						secondBestLocalCost = bestLocalCost
						bestLocalCost = insertionCost
						bestPos = pos
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
				// Insert the node with highest score at its best position
				if bestPosition == len(path) {
					path = append(path, bestNodeIndex)
				} else {
					path = append(path[:bestPosition], append([]int{bestNodeIndex}, path[bestPosition:]...)...)
				}
				delete(unvisited, bestNodeIndex)
			} else {
				break
			}
		}

		totalDistance := 0
		if len(path) > 1 {
			for j := 0; j < len(path); j++ {
				next := (j + 1) % len(path)
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
