package algorithms

import (
	"math"

	"github.com/czajkowskis/evolutionary_computation/01_labs/greedy_heuristics/pkg/data"
)

func NearestNeighborEnd(nodes []data.Node, distanceMatrix [][]int, startNodeIndices []int) []Solution {
	n := len(nodes)
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
			nearestNodeIndex := -1
			minDistance := math.MaxInt32

			for nodeIndex := range unvisited {
				if distanceMatrix[lastNodeIndex][nodeIndex] < minDistance {
					minDistance = distanceMatrix[lastNodeIndex][nodeIndex]
					nearestNodeIndex = nodeIndex
				}
			}

			if nearestNodeIndex != -1 {
				path = append(path, nearestNodeIndex)
				delete(unvisited, nearestNodeIndex)
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
			totalCost += nodes[idx].Cost
		}

		objective := totalDistance + totalCost
		solutions = append(solutions, Solution{path, objective})
	}
	return solutions
}

func NearestNeighborAny(nodes []data.Node, distanceMatrix [][]int, startNodeIndices []int) []Solution {
	n := len(nodes)
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
				for pos := 0; pos <= len(path); pos++ {
					var increase int
					if len(path) == 1 {
						increase = distanceMatrix[path[0]][nodeIndex] * 2
					} else {
						prevIdx := (pos - 1 + len(path)) % len(path)
						currIdx := pos % len(path)
						increase = distanceMatrix[path[prevIdx]][nodeIndex] + distanceMatrix[nodeIndex][path[currIdx]] - distanceMatrix[path[prevIdx]][path[currIdx]]
					}

					if increase < minIncrease {
						minIncrease = increase
						bestNodeIndex = nodeIndex
						bestPosition = pos
					}
				}
			}

			if bestNodeIndex != -1 {
				if bestPosition == len(path) {
					path = append(path, bestNodeIndex)
				} else {
					path = append(path[:bestPosition], append([]int{bestNodeIndex}, path[bestPosition:]...)...)
				}
				totalDistance += minIncrease
				delete(unvisited, bestNodeIndex)
			} else {
				break
			}
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
