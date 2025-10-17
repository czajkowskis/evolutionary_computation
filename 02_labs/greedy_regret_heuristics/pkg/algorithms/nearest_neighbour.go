package algorithms

import (
	"math"
	"sort"
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
			maxPossibleRegret := 0.0
			maxPossibleObjective := 0.0

			type insertionInfo struct {
				nodeIndex      int
				bestCost       int
				secondBestCost int
				bestPosition   int
			}
			var insertionInfos []insertionInfo

			// Single pass to calculate costs and find normalization values
			for nodeIndex := range unvisited {
				bestLocalCost := math.MaxInt32
				secondBestLocalCost := math.MaxInt32
				bestPos := -1

				for pos := 0; pos <= len(path); pos++ {
					var insertionCost int
					if pos == 0 {
						// Insert at the beginning of the path
						insertionCost = distanceMatrix[nodeIndex][path[0]] + nodeCosts[nodeIndex]
					} else if pos == len(path) {
						// Insert at the end of the path
						insertionCost = distanceMatrix[path[len(path)-1]][nodeIndex] + nodeCosts[nodeIndex]
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
				insertionInfos = append(insertionInfos, insertionInfo{nodeIndex, bestLocalCost, secondBestLocalCost, bestPos})

				if secondBestLocalCost == math.MaxInt32 {
					secondBestLocalCost = bestLocalCost
				}
				regret := secondBestLocalCost - bestLocalCost
				if float64(regret) > maxPossibleRegret {
					maxPossibleRegret = float64(regret)
				}
				if float64(bestLocalCost) > maxPossibleObjective {
					maxPossibleObjective = float64(bestLocalCost)
				}
			}

			// Sort insertionInfos to ensure deterministic iteration order for tie-breaking
			sort.Slice(insertionInfos, func(i, j int) bool {
				return insertionInfos[i].nodeIndex < insertionInfos[j].nodeIndex
			})

			// Avoid division by zero
			if maxPossibleRegret == 0 {
				maxPossibleRegret = 1
			}
			if maxPossibleObjective == 0 {
				maxPossibleObjective = 1
			}

			bestScore := math.Inf(-1)
			bestNodeIndex := -1
			bestPosition := -1

			// Find the best node to insert using the stored information
			for _, info := range insertionInfos {
				if info.secondBestCost == math.MaxInt32 {
					info.secondBestCost = info.bestCost
				}
				regret := info.secondBestCost - info.bestCost
				normalizedRegret := float64(regret) / maxPossibleRegret
				normalizedObjective := float64(info.bestCost) / maxPossibleObjective

				score := regretWeight*normalizedRegret - objectiveWeight*normalizedObjective
				if score > bestScore {
					bestScore = score
					bestNodeIndex = info.nodeIndex
					bestPosition = info.bestPosition
				}
			}

			if bestNodeIndex != -1 {
				path = append(path[:bestPosition], append([]int{bestNodeIndex}, path[bestPosition:]...)...)
				delete(unvisited, bestNodeIndex)
			} else {
				break // No more nodes can be inserted
			}
		}

		objective := 0
		if len(path) > 1 {
			for i := 0; i < len(path)-1; i++ {
				objective += distanceMatrix[path[i]][path[i+1]]
			}
			objective += distanceMatrix[path[len(path)-1]][path[0]]
		}
		for _, nodeIndex := range path {
			objective += nodeCosts[nodeIndex]
		}

		solutions = append(solutions, Solution{Path: path, Objective: objective})
	}

	return solutions
}
