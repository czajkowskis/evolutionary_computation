package algorithms

// FindBestSolution finds the solution with the minimum objective value
func FindBestSolution(solutions []Solution) Solution {
	if len(solutions) == 0 {
		return Solution{}
	}
	best := solutions[0]
	for _, s := range solutions {
		if s.Objective < best.Objective {
			best = s
		}
	}
	return best
}

