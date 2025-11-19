package algorithms

import (
	"math/rand"
	"time"
)

// MSLSResult contains the results of MSLS algorithm
type MSLSResult struct {
	BestSolution    Solution
	NumLSIterations int
	Elapsed         time.Duration
	AllSolutions    []Solution
}

// MSLS performs Multiple Start Local Search
// It runs steepest local search multiple times (200 iterations) from random starting solutions
func MSLS(D [][]int, costs []int, iterations int) MSLSResult {
	startTime := time.Now()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Initialize with first random solution
	initialSolution := startRandom(D, costs, rng)
	initialSolution = localSearchSteepestBaseline(D, costs, initialSolution)

	bestSolution := initialSolution
	allSolutions := []Solution{initialSolution}

	// Run remaining iterations
	for i := 1; i < iterations; i++ {
		// Generate new random starting solution
		randomStart := startRandom(D, costs, rng)
		current := localSearchSteepestBaseline(D, costs, randomStart)

		allSolutions = append(allSolutions, current)

		// Update best solution if current is better
		if current.Objective < bestSolution.Objective {
			bestSolution = current
		}
	}

	elapsed := time.Since(startTime)

	return MSLSResult{
		BestSolution:    bestSolution,
		NumLSIterations: iterations,
		Elapsed:         elapsed,
		AllSolutions:    allSolutions,
	}
}
func RunMSLS(D [][]int, costs []int, numMSLSStarts int) MSLSResult {
	return MSLS(D, costs, numMSLSStarts)
}
