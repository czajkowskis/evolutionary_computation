package algorithms

import (
	"math/rand"
	"time"
)

// localSearchSteepestBaseline performs steepest local search on the full
// neighborhood (2-opt intra-route plus exchanges with unselected vertices).
func localSearchSteepestBaseline(D [][]int, costs []int, init Solution) Solution {
	path := append([]int(nil), init.Path...)
	n := len(path)

	for {
		bestDelta := 0
		var bestMove func()

		// intra-route move - two-edges exchange: 2-opt
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

		// inter-route moves - path[i] with u (u outside the current path)
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

// RunLocalSearchBatch runs a batch of independently initialised local searches
// for a given method specification and returns all final solutions together
// with per-run durations.
func RunLocalSearchBatch(
	D [][]int,
	costs []int,
	m MethodSpec,
	numSolutions int,
) ([]Solution, []time.Duration) {
	if numSolutions <= 0 {
		return nil, nil
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	results := make([]Solution, 0, numSolutions)
	durations := make([]time.Duration, 0, numSolutions)

	for r := 0; r < numSolutions; r++ {
		init := startRandom(D, costs, rng)

		start := time.Now()
		var sol Solution
		if m.UseCand {
			// candidate moves (2-opt + inter) with precomputed candidate lists
			K := m.CandK
			if K <= 0 {
				K = 10
			}
			cd := buildCandidates(D, costs, K)
			sol = localSearchSteepestCandidates(D, costs, init, cd)
		} else if m.UseLM {
			// steepest local search with list-of-moves (delta reuse)
			sol = localSearchSteepestLM(D, costs, init)
		} else {
			// baseline steepest local search (full neighborhood, no LM)
			sol = localSearchSteepestBaseline(D, costs, init)
		}
		results = append(results, sol)
		durations = append(durations, time.Since(start))
	}
	return results, durations
}
