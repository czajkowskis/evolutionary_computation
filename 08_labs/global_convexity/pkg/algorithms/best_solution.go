package algorithms

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
