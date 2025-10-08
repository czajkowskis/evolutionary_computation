package algorithms

import (
	"math"
	"math/rand"

	"github.com/czajkowskis/evolutionary_computation/01_labs/greedy_heuristics/pkg/data"
)

func NearestNeighborEnd(nodes []data.Node, distanceMatrix [][]int, numSolutions int) []Solution {
	n := len(nodes)
	k := (n + 1) / 2
	var solutions []Solution

	for i := 0; i < numSolutions; i++ {
		startNode := rand.Intn(n)
		path := []int{startNode}
		unvisited := make(map[int]bool)
		for j := 0; j < n; j++ {
			if j != startNode {
				unvisited[j] = true
			}
		}

		for len(path) < k {
			lastNode := path[len(path)-1]
			nearestNode := -1
			minDistance := math.MaxInt32

			for node := range unvisited {
				if distanceMatrix[lastNode][node] < minDistance {
					minDistance = distanceMatrix[lastNode][node]
					nearestNode = node
				}
			}

			if nearestNode != -1 {
				path = append(path, nearestNode)
				delete(unvisited, nearestNode)
			}
		}

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

func NearestNeighborAny(nodes []data.Node, distanceMatrix [][]int, numSolutions int) []Solution {
	n := len(nodes)
	k := (n + 1) / 2
	var solutions []Solution

	for i := 0; i < numSolutions; i++ {
		startNode := rand.Intn(n)
		path := []int{startNode}
		unvisited := make(map[int]bool)
		for j := 0; j < n; j++ {
			if j != startNode {
				unvisited[j] = true
			}
		}

		for len(path) < k {
			bestIncrease := math.MaxInt32
			bestNode := -1
			bestPosition := -1

			for node := range unvisited {
				for pos := 0; pos <= len(path); pos++ {
					var increase int
					if pos == 0 {
						increase = distanceMatrix[path[0]][node] + distanceMatrix[node][path[len(path)-1]] - distanceMatrix[path[0]][path[len(path)-1]]
					} else if pos == len(path) {
						increase = distanceMatrix[path[len(path)-1]][node] + distanceMatrix[node][path[0]] - distanceMatrix[path[len(path)-1]][path[0]]
					} else {
						increase = distanceMatrix[path[pos-1]][node] + distanceMatrix[node][path[pos]] - distanceMatrix[path[pos-1]][path[pos]]
					}

					if increase < bestIncrease {
						bestIncrease = increase
						bestNode = node
						bestPosition = pos
					}
				}
			}

			if bestNode != -1 {
				path = append(path[:bestPosition+1], path[bestPosition:]...)
				path[bestPosition] = bestNode
				delete(unvisited, bestNode)
			}
		}

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
