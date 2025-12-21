package utils

// GenerateStartNodeIndices creates a list of starting node indices for the algorithms,
// using each node index once.
func GenerateStartNodeIndices(n int) []int {
	startNodeIndices := make([]int, n)
	for i := 0; i < n; i++ {
		startNodeIndices[i] = i
	}
	return startNodeIndices
}

