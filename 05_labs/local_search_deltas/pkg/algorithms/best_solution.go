package algorithms

// FindBestSolution returns the solution with the smallest objective value from
// the provided slice. For an empty slice it returns the zero-value Solution.
func FindBestSolution(solutions []Solution) Solution {
	if len(solutions) == 0 {
		return Solution{}
	}
	bestSolution := solutions[0]
	for _, sol := range solutions {
		if sol.Objective < bestSolution.Objective {
			bestSolution = sol
		}
	}
	return bestSolution
}
