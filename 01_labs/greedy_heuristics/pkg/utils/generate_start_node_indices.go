package utils

import "math/rand"

// GenerateStartNodeIndices creates a list of random starting node indices for the algorithms.
func GenerateStartNodeIndices(n, numSolutions int) []int {
	startNodeIndices := make([]int, numSolutions)
	for i := 0; i < numSolutions; i++ {
		startNodeIndices[i] = rand.Intn(n)
	}
	return startNodeIndices
}
