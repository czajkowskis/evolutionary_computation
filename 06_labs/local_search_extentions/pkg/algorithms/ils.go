package algorithms

import (
	"math/rand"
	"time"
)

// ILSResult contains results from ILS run
type ILSResult struct {
	BestSolution    Solution
	NumLSIterations int
	Elapsed         time.Duration
	AllSolutions    []Solution
}

// PerturbationType defines the type of perturbation
type PerturbationType int

const (
	PerturbDoubleExchange PerturbationType = iota
	PerturbRandom4Opt
	PerturbPathDestroy
)

// applyPerturbation applies a perturbation to escape local optimum
func applyPerturbation(D [][]int, costs []int, sol Solution, perturbType PerturbationType, rng *rand.Rand) Solution {
	switch perturbType {
	case PerturbDoubleExchange:
		return perturbDoubleExchange(D, costs, sol, rng)
	case PerturbRandom4Opt:
		return perturbRandom4Opt(D, costs, sol, rng)
	case PerturbPathDestroy:
		return perturbPathDestroy(D, costs, sol, rng)
	default:
		return perturbDoubleExchange(D, costs, sol, rng)
	}
}

// perturbDoubleExchange - Exchange 2 pairs of selected/non-selected nodes
func perturbDoubleExchange(D [][]int, costs []int, sol Solution, rng *rand.Rand) Solution {
	path := append([]int(nil), sol.Path...)
	n := len(path)

	if n < 2 {
		return Solution{Path: path, Objective: objective(D, costs, path)}
	}

	inSel := make([]bool, len(D))
	for _, v := range path {
		inSel[v] = true
	}
	nonSel := make([]int, 0, len(D)-n)
	for u := 0; u < len(D); u++ {
		if !inSel[u] {
			nonSel = append(nonSel, u)
		}
	}

	if len(nonSel) < 2 {
		return Solution{Path: path, Objective: objective(D, costs, path)}
	}

	numExchanges := 2
	if n < 2 {
		numExchanges = 1
	}

	for ex := 0; ex < numExchanges; ex++ {
		i := rng.Intn(n)
		u := nonSel[rng.Intn(len(nonSel))]

		oldNode := path[i]
		path[i] = u

		inSel[oldNode] = false
		inSel[u] = true

		newNonSel := make([]int, 0, len(nonSel))
		for _, v := range nonSel {
			if v != u {
				newNonSel = append(newNonSel, v)
			}
		}
		newNonSel = append(newNonSel, oldNode)
		nonSel = newNonSel
	}

	return Solution{Path: path, Objective: objective(D, costs, path)}
}

// perturbRandom4Opt - Apply multiple random 2-opt moves
func perturbRandom4Opt(D [][]int, costs []int, sol Solution, rng *rand.Rand) Solution {
	path := append([]int(nil), sol.Path...)
	n := len(path)

	if n < 4 {
		return perturbDoubleExchange(D, costs, sol, rng)
	}

	numMoves := 2 + rng.Intn(2)

	for m := 0; m < numMoves; m++ {
		i := rng.Intn(n)
		j := rng.Intn(n)
		if i > j {
			i, j = j, i
		}
		if j-i > 1 && j-i < n-1 {
			applyTwoOpt(path, i, j)
		}
	}

	return Solution{Path: path, Objective: objective(D, costs, path)}
}

// perturbPathDestroy - Destroy 25-30% of path and reconstruct randomly
func perturbPathDestroy(D [][]int, costs []int, sol Solution, rng *rand.Rand) Solution {
	path := append([]int(nil), sol.Path...)
	n := len(path)

	if n < 4 {
		return perturbDoubleExchange(D, costs, sol, rng)
	}

	numRemove := n / 4
	if numRemove < 1 {
		numRemove = 1
	}
	if numRemove > n/3 {
		numRemove = n / 3
	}

	inSel := make([]bool, len(D))
	for _, v := range path {
		inSel[v] = true
	}

	for i := 0; i < numRemove && len(path) > 2; i++ {
		idx := rng.Intn(len(path))
		inSel[path[idx]] = false
		path = append(path[:idx], path[idx+1:]...)
	}

	nonSel := make([]int, 0, len(D)-len(path))
	for u := 0; u < len(D); u++ {
		if !inSel[u] {
			nonSel = append(nonSel, u)
		}
	}

	targetSize := n
	for len(path) < targetSize && len(nonSel) > 0 {
		idx := rng.Intn(len(nonSel))
		u := nonSel[idx]
		path = append(path, u)
		nonSel = append(nonSel[:idx], nonSel[idx+1:]...)
	}

	k := len(path) / 3
	if k > 0 {
		start := rng.Intn(len(path) - k)
		for i := 0; i < k; i++ {
			j := start + rng.Intn(k)
			path[start+i], path[j] = path[j], path[start+i]
		}
	}

	return Solution{Path: path, Objective: objective(D, costs, path)}
}

// RunILS - Iterated Local Search
func RunILS(D [][]int, costs []int, timeLimit time.Duration, perturbType PerturbationType) ILSResult {
	startTime := time.Now()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	current := startRandom(D, costs, rng)
	current = localSearchSteepestBaseline(D, costs, current)

	bestSolution := current
	numLSIterations := 1
	allSolutions := []Solution{current}

	for time.Since(startTime) < timeLimit {
		perturbed := applyPerturbation(D, costs, current, perturbType, rng)
		localOpt := localSearchSteepestBaseline(D, costs, perturbed)
		numLSIterations++
		allSolutions = append(allSolutions, localOpt)

		if localOpt.Objective < current.Objective {
			current = localOpt
		}

		if localOpt.Objective < bestSolution.Objective {
			bestSolution = localOpt
		}

		if time.Since(startTime) >= timeLimit {
			break
		}
	}

	elapsed := time.Since(startTime)
	return ILSResult{
		BestSolution:    bestSolution,
		NumLSIterations: numLSIterations,
		Elapsed:         elapsed,
		AllSolutions:    allSolutions,
	}
}
