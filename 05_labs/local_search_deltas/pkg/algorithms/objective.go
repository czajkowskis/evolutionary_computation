package algorithms

import (
	"math"
	"math/rand"
)

// MethodSpec describes a single configured local search method used in experiments.
// It controls whether candidate moves or list-of-moves (LM) delta reuse is used.
type MethodSpec struct {
	Name    string
	UseCand bool // should use candidate moves?
	CandK   int  // how many nearest to include in candidate list
	UseLM   bool // should use list-of-moves (LM) delta reuse?
}

// objective computes the tour length plus node costs for a given path.
// An empty path is treated as a very large (effectively infinite) objective.
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

// prevIdx returns the previous index in a cyclic path of length n.
func prevIdx(i, n int) int {
	if i == 0 {
		return n - 1
	}
	return i - 1
}

// nextIdx returns the next index in a cyclic path of length n.
func nextIdx(i, n int) int { return (i + 1) % n }

// selectCount returns the number of nodes to select (50% rounded up).
func selectCount(n int) int { return (n + 1) / 2 }

// startRandom builds an initial solution by selecting selectCount(n) nodes at
// random and shuffling their order.
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
