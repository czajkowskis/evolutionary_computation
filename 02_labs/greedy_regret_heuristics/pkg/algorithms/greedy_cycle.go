package algorithms

import (
	"math"
	"sort"
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

				for i := 0; i < len(path); i++ {
					p1 := path[i]
					p2 := path[(i+1)%len(path)]
					deltaDist := distanceMatrix[p1][nodeIndex] + distanceMatrix[nodeIndex][p2] - distanceMatrix[p1][p2]
					insertionCost := deltaDist + nodeCosts[nodeIndex]

					if insertionCost < bestLocalCost {
						secondBestLocalCost = bestLocalCost
						bestLocalCost = insertionCost
						bestPos = i
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
				if score >= bestScore {
					bestScore = score
					bestNodeIndex = info.nodeIndex
					bestPosition = info.bestPosition
				}
			}

			if bestNodeIndex != -1 {
				// Insert the best node at the best position
				bestPositionIndex := bestPosition + 1
				path = append(path[:bestPositionIndex], append([]int{bestNodeIndex}, path[bestPositionIndex:]...)...)
				delete(unvisited, bestNodeIndex)
			} else {
				break // No more nodes can be inserted
			}
		}

		// Calculate final objective value
		objective := 0
		if len(path) > 0 {
			for i := 0; i < len(path); i++ {
				objective += distanceMatrix[path[i]][path[(i+1)%len(path)]]
			}
			for _, nodeIndex := range path {
				objective += nodeCosts[nodeIndex]
			}
		}

		solutions = append(solutions, Solution{Path: path, Objective: objective})
	}

	return solutions
}
