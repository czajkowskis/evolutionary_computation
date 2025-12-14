package algorithms

import (
	"math"
	"math/rand"
)

// Calculate objective function value
func objective(D [][]int, costs []int, path []int) int {
	if len(path) == 0 {
		return math.MaxInt32 / 4
	}
	sum := 0
	n := len(path)
	for i := 0; i < n; i++ {
		a := path[i]
		b := path[(i+1)%n]
		sum += D[a][b]
	}
	for _, v := range path {
		sum += costs[v]
	}
	return sum
}

// randomConstruction builds a random solution by selecting random nodes
func randomConstruction(D [][]int, costs []int, n, targetSize int, rng *rand.Rand) Solution {
	// Create list of all nodes
	allNodes := make([]int, n)
	for i := 0; i < n; i++ {
		allNodes[i] = i
	}

	// Shuffle and take first targetSize nodes
	rng.Shuffle(n, func(i, j int) {
		allNodes[i], allNodes[j] = allNodes[j], allNodes[i]
	})

	path := make([]int, targetSize)
	copy(path, allNodes[:targetSize])

	return Solution{Path: path, Objective: objective(D, costs, path)}
}

// repair rebuilds the solution using nearest neighbor any position heuristic
func repair(partialPath []int, D [][]int, costs []int, targetSize int, rng *rand.Rand) Solution {
	n := len(costs)

	// Create set of nodes already in solution
	inSolution := make(map[int]bool)
	for _, node := range partialPath {
		inSolution[node] = true
	}

	// Create list of unvisited nodes
	unvisited := make(map[int]bool)
	for i := 0; i < n; i++ {
		if !inSolution[i] {
			unvisited[i] = true
		}
	}

	path := make([]int, len(partialPath))
	copy(path, partialPath)

	// Add nodes using greedy nearest neighbor any position
	for len(path) < targetSize && len(unvisited) > 0 {
		minIncrease := math.MaxInt32
		bestNode := -1
		bestPosition := -1

		for node := range unvisited {
			localMinIncrease := math.MaxInt32
			localBestPos := -1

			if len(path) == 0 {
				localMinIncrease = costs[node]
				localBestPos = 0
			} else if len(path) == 1 {
				// Insert before or after
				inc := D[node][path[0]] + costs[node]
				if inc < localMinIncrease {
					localMinIncrease = inc
					localBestPos = 0
				}
				inc = D[path[0]][node] + costs[node]
				if inc < localMinIncrease {
					localMinIncrease = inc
					localBestPos = 1
				}
			} else {
				// Try all positions
				for pos := 0; pos <= len(path); pos++ {
					var inc int
					if pos == 0 {
						// Insert at beginning
						inc = D[node][path[0]] + costs[node]
					} else if pos == len(path) {
						// Insert at end
						inc = D[path[len(path)-1]][node] + costs[node]
					} else {
						// Insert in middle
						a := path[pos-1]
						b := path[pos]
						deltaDist := D[a][node] + D[node][b] - D[a][b]
						inc = deltaDist + costs[node]
					}

					if inc < localMinIncrease {
						localMinIncrease = inc
						localBestPos = pos
					}
				}
			}

			if localMinIncrease < minIncrease {
				minIncrease = localMinIncrease
				bestNode = node
				bestPosition = localBestPos
			}
		}

		if bestNode != -1 {
			// Insert node at best position
			if bestPosition == len(path) {
				path = append(path, bestNode)
			} else {
				path = append(path[:bestPosition], append([]int{bestNode}, path[bestPosition:]...)...)
			}
			delete(unvisited, bestNode)
		} else {
			break
		}
	}

	return Solution{Path: path, Objective: objective(D, costs, path)}
}

// isDuplicate checks if solution already exists in population
func isDuplicate(sol Solution, population []Solution) bool {
	for _, existing := range population {
		if sol.Objective == existing.Objective {
			// Could also compare paths for exact duplicate
			return true
		}
	}
	return false
}

// findWorstIndex returns index of worst solution in population
func findWorstIndex(population []Solution) int {
	worstIdx := 0
	worstObj := population[0].Objective

	for i := 1; i < len(population); i++ {
		if population[i].Objective > worstObj {
			worstObj = population[i].Objective
			worstIdx = i
		}
	}

	return worstIdx
}

// reverseSlice reverses a slice in place
func reverseSlice(slice []int) {
	for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
		slice[i], slice[j] = slice[j], slice[i]
	}
}

// Helper functions for index manipulation
func prevIdx(i, n int) int {
	if i == 0 {
		return n - 1
	}
	return i - 1
}

func nextIdx(i, n int) int {
	return (i + 1) % n
}

// Calculate delta for 2-opt move (intra-route)
func deltaTwoOpt(D [][]int, path []int, i, j int) int {
	if i == j {
		return 0
	}
	n := len(path)
	// Adjacent edges -> degenerate 2-opt (no change)
	if nextIdx(i, n) == j || nextIdx(j, n) == i {
		return 0
	}
	a := path[i]
	b := path[nextIdx(i, n)]
	c := path[j]
	d := path[nextIdx(j, n)]
	before := D[a][b] + D[c][d]
	after := D[a][c] + D[b][d]
	return after - before
}

// Calculate delta for node exchange (inter-route)
func deltaExchangeSelected(D [][]int, costs []int, path []int, i int, u int) int {
	n := len(path)
	a := path[prevIdx(i, n)]
	v := path[i]
	b := path[nextIdx(i, n)]
	before := D[a][v] + D[v][b] + costs[v]
	after := D[a][u] + D[u][b] + costs[u]
	return after - before
}

// Apply 2-opt move
func applyTwoOpt(path []int, i, j int) {
	n := len(path)
	if i == j || nextIdx(i, n) == j || nextIdx(j, n) == i {
		return
	}
	if i > j {
		i, j = j, i
	}
	for l, r := i+1, j; l < r; l, r = l+1, r-1 {
		path[l], path[r] = path[r], path[l]
	}
}

// Apply node exchange
func applyExchangeSelected(path []int, i int, u int) {
	path[i] = u
}

// Steepest local search
func localSearchSteepest(D [][]int, costs []int, init Solution) Solution {
	path := append([]int(nil), init.Path...)
	n := len(path)

	for {
		bestDelta := 0
		var bestMove func()

		// Intra-route moves: 2-opt
		for i := 0; i < n; i++ {
			for j := i + 1; j < n; j++ {
				dl := deltaTwoOpt(D, path, i, j)
				if dl < bestDelta {
					ii, jj := i, j
					bestDelta = dl
					bestMove = func() { applyTwoOpt(path, ii, jj) }
				}
			}
		}

		// Inter-route moves: node exchange with non-selected nodes
		inSel := make([]bool, len(D))
		for _, v := range path {
			inSel[v] = true
		}
		nonSel := make([]int, 0, len(D)-n)
		for u := range D {
			if !inSel[u] {
				nonSel = append(nonSel, u)
			}
		}
		for i := 0; i < n; i++ {
			for _, u := range nonSel {
				dl := deltaExchangeSelected(D, costs, path, i, u)
				if dl < bestDelta {
					ii, uu := i, u
					bestDelta = dl
					bestMove = func() { applyExchangeSelected(path, ii, uu) }
				}
			}
		}

		if bestDelta < 0 {
			bestMove()
		} else {
			break
		}
	}
	return Solution{Path: path, Objective: objective(D, costs, path)}
}
