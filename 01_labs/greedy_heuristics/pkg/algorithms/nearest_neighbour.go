package algorithms

import (
	"math"
)

func NearestNeighborEnd(distanceMatrix [][]int, nodeCosts []int, startNodeIndices []int) []Solution {
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
			lastNodeIndex := path[len(path)-1]
			bestNodeIndex := -1
			minScore := math.MaxInt32

			for nodeIndex := range unvisited {
				score := distanceMatrix[lastNodeIndex][nodeIndex] + nodeCosts[nodeIndex]
				if score < minScore {
					minScore = score
					bestNodeIndex = nodeIndex
				}
			}

			if bestNodeIndex != -1 {
				path = append(path, bestNodeIndex)
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

func NearestNeighborAny(distanceMatrix [][]int, nodeCosts []int, startNodeIndices []int) []Solution {
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

		totalDistance := 0
		for len(path) < k {
			minIncrease := math.MaxInt32
			bestNodeIndex := -1
			bestPosition := -1

			for nodeIndex := range unvisited {
				localMinIncrease := math.MaxInt32
				localBestPos := -1

				if len(path) == 1 {
					// Path has only one node -> we can only insert before or after it
					inc := distanceMatrix[nodeIndex][path[0]] + nodeCosts[nodeIndex]
					if inc < localMinIncrease {
						localMinIncrease = inc
						localBestPos = 0
					}
					inc = distanceMatrix[path[0]][nodeIndex] + nodeCosts[nodeIndex]
					if inc < localMinIncrease {
						localMinIncrease = inc
						localBestPos = 1
					}
				} else {
					// Insert at the beginning
					incFront := distanceMatrix[nodeIndex][path[0]] + nodeCosts[nodeIndex]
					if incFront < localMinIncrease {
						localMinIncrease = incFront
						localBestPos = 0
					}
					// Insert at the end
					incBack := distanceMatrix[path[len(path)-1]][nodeIndex] + nodeCosts[nodeIndex]
					if incBack < localMinIncrease {
						localMinIncrease = incBack
						localBestPos = len(path)
					}
					// Insert in the middle
					for pos := 1; pos < len(path); pos++ {
						a := path[pos-1]
						b := path[pos]
						deltaDist := distanceMatrix[a][nodeIndex] + distanceMatrix[nodeIndex][b] - distanceMatrix[a][b]
						inc := deltaDist + nodeCosts[nodeIndex]
						if inc < localMinIncrease {
							localMinIncrease = inc
							localBestPos = pos
						}
					}
				}

				if localMinIncrease < minIncrease {
					minIncrease = localMinIncrease
					bestNodeIndex = nodeIndex
					bestPosition = localBestPos
				}
			}

			if bestNodeIndex != -1 {
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

		totalDistance = 0
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
