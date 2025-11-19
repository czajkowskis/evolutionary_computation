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

// Calculate number of nodes to select
func selectCount(n int) int {
	return (n + 1) / 2
}

// Generate random starting solution
func startRandom(D [][]int, costs []int, rng *rand.Rand) Solution {
	n := len(D)
	k := selectCount(n)
	idx := make([]int, n)
	for i := range idx {
		idx[i] = i
	}
	rng.Shuffle(n, func(i, j int) { idx[i], idx[j] = idx[j], idx[i] })
	path := append([]int(nil), idx[:k]...)
	rng.Shuffle(k, func(i, j int) { path[i], path[j] = path[j], path[i] })
	return Solution{Path: path, Objective: objective(D, costs, path)}
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

// Steepest local search baseline
func localSearchSteepestBaseline(D [][]int, costs []int, init Solution) Solution {
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
