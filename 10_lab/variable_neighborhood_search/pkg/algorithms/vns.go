package algorithms

import (
	"math"
	"math/rand"
	"time"
)

// VNSConfig holds configuration for Variable Neighborhood Search
type VNSConfig struct {
	TimeLimit               time.Duration // Time limit for the algorithm (0 = no time limit)
	MaxIterations           int           // Maximum number of iterations (0 = no limit)
	MaxIterationsNoImprove  int           // Maximum iterations without improvement (0 = no limit)
	MaxNeighborhoods        int           // Maximum number of neighborhoods to try
	ShakingIntensity        int           // Number of moves in shaking (default: 3)
	NeighborhoodChange      string        // Strategy: "sequential", "random", "adaptive"
	UseLocalSearch          bool          // Whether to use local search after shaking
	AdaptiveIntensity       bool          // Enable adaptive shaking intensity based on success rate
	InitialSolutionStrategy string        // "random" or "greedy" - strategy for initial solution
	UseMemory               bool          // Enable solution memory to avoid revisiting recent solutions
	BestImprovement         bool          // Use best improvement within cycle instead of first improvement
}

// VNSResult contains the result of VNS execution
type VNSResult struct {
	BestSolution               Solution
	Iterations                 int
	Duration                   time.Duration
	NeighborhoodUsage          []int   // Count of each neighborhood used
	ImprovementsByNeighborhood []int   // Improvements per neighborhood
	AvgShakingIntensity        float64 // Average intensity used
}

// NeighborhoodType represents different shaking operators
type NeighborhoodType int

const (
	NeighborhoodNodeExchange NeighborhoodType = iota
	NeighborhoodTwoOpt
	NeighborhoodDestroyRepair
	NeighborhoodDoubleBridge
)

// VariableNeighborhoodSearch implements VNS algorithm
func VariableNeighborhoodSearch(D [][]int, costs []int, config VNSConfig) VNSResult {
	startTime := time.Now()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Set defaults
	if config.MaxNeighborhoods == 0 {
		config.MaxNeighborhoods = 4
	}
	if config.ShakingIntensity == 0 {
		config.ShakingIntensity = 3
	}
	if config.NeighborhoodChange == "" {
		config.NeighborhoodChange = "sequential"
	}
	if config.InitialSolutionStrategy == "" {
		config.InitialSolutionStrategy = "random"
	}

	// Generate initial solution
	var currentSolution Solution
	if config.InitialSolutionStrategy == "greedy" {
		currentSolution = startGreedy(D, costs, rng)
		currentSolution = localSearchSteepest(D, costs, currentSolution)
	} else {
		currentSolution = startRandom(D, costs, rng)
		currentSolution = localSearchSteepest(D, costs, currentSolution)
	}

	bestSolution := currentSolution
	iterations := 0
	iterationsNoImprove := 0

	// Statistics tracking
	neighborhoodUsage := make([]int, config.MaxNeighborhoods)
	improvementsByNeighborhood := make([]int, config.MaxNeighborhoods)
	totalIntensity := 0.0
	intensityCount := 0

	// Stopping condition check
	shouldContinue := func() bool {
		// Check time limit
		if config.TimeLimit > 0 && time.Since(startTime) >= config.TimeLimit {
			return false
		}
		// Check max iterations
		if config.MaxIterations > 0 && iterations >= config.MaxIterations {
			return false
		}
		// Check iterations without improvement
		if config.MaxIterationsNoImprove > 0 && iterationsNoImprove >= config.MaxIterationsNoImprove {
			return false
		}
		return true
	}

	for shouldContinue() {
		iterations++
		k := 1 // Start with first neighborhood

		for k <= config.MaxNeighborhoods && shouldContinue() {
			// Select neighborhood based on strategy
			neighborhoodIdx := selectNeighborhood(k, config.NeighborhoodChange, rng)
			neighborhoodUsage[neighborhoodIdx]++

			// Shaking: apply neighborhood operator
			shakingIntensity := config.ShakingIntensity
			shakenSolution := shake(D, costs, currentSolution, NeighborhoodType(neighborhoodIdx), shakingIntensity, rng)

			// Local search intensification
			var localOpt Solution
			if config.UseLocalSearch {
				localOpt = localSearchSteepest(D, costs, shakenSolution)
			} else {
				localOpt = shakenSolution
			}

			// Track intensity
			totalIntensity += float64(shakingIntensity)
			intensityCount++

			// Acceptance: first improvement
			if localOpt.Objective < currentSolution.Objective {
				currentSolution = localOpt
				improvementsByNeighborhood[neighborhoodIdx]++

				// Update best solution
				if localOpt.Objective < bestSolution.Objective {
					bestSolution = localOpt
					iterationsNoImprove = 0
				} else {
					iterationsNoImprove++
				}

				// Return to first neighborhood (first improvement strategy)
				k = 1
				break
			} else {
				iterationsNoImprove++
				// Move to next neighborhood
				k++
			}
		}
	}

	elapsed := time.Since(startTime)
	avgIntensity := 0.0
	if intensityCount > 0 {
		avgIntensity = totalIntensity / float64(intensityCount)
	}
	return VNSResult{
		BestSolution:               bestSolution,
		Iterations:                 iterations,
		Duration:                   elapsed,
		NeighborhoodUsage:          neighborhoodUsage,
		ImprovementsByNeighborhood: improvementsByNeighborhood,
		AvgShakingIntensity:        avgIntensity,
	}
}

// selectNeighborhood chooses which neighborhood to use based on strategy
func selectNeighborhood(k int, strategy string, rng *rand.Rand) int {
	switch strategy {
	case "random":
		return rng.Intn(4)
	case "adaptive":
		// For default version, fall back to sequential
		return (k - 1) % 4
	default: // "sequential"
		return (k - 1) % 4
	}
}

// shake applies a shaking operator to escape local optimum
func shake(D [][]int, costs []int, sol Solution, nType NeighborhoodType, intensity int, rng *rand.Rand) Solution {
	switch nType {
	case NeighborhoodNodeExchange:
		return shakeNodeExchange(D, costs, sol, intensity, rng)
	case NeighborhoodTwoOpt:
		return shakeTwoOpt(D, costs, sol, intensity, rng)
	case NeighborhoodDestroyRepair:
		return shakeDestroyRepair(D, costs, sol, intensity, rng)
	case NeighborhoodDoubleBridge:
		return shakeDoubleBridge(D, costs, sol, intensity, rng)
	default:
		return shakeNodeExchange(D, costs, sol, intensity, rng)
	}
}

// N2: Random 2-opt moves
func shakeTwoOpt(D [][]int, costs []int, sol Solution, numMoves int, rng *rand.Rand) Solution {
	path := append([]int(nil), sol.Path...)
	n := len(path)

	if n < 4 {
		return Solution{Path: path, Objective: objective(D, costs, path)}
	}

	numMoves = min(numMoves, n/2)
	for i := 0; i < numMoves; i++ {
		idx1 := rng.Intn(n)
		idx2 := rng.Intn(n)
		if idx1 > idx2 {
			idx1, idx2 = idx2, idx1
		}
		if idx2-idx1 > 1 && idx2-idx1 < n-1 {
			applyTwoOpt(path, idx1, idx2)
		}
	}

	return Solution{Path: path, Objective: objective(D, costs, path)}
}

// N1: Node exchange (exchange selected nodes with non-selected nodes)
func shakeNodeExchange(D [][]int, costs []int, sol Solution, numSwaps int, rng *rand.Rand) Solution {
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

	if len(nonSel) == 0 {
		return Solution{Path: path, Objective: objective(D, costs, path)}
	}

	numSwaps = min(numSwaps, n, len(nonSel))
	for i := 0; i < numSwaps; i++ {
		pathIdx := rng.Intn(len(path))
		nonSelIdx := rng.Intn(len(nonSel))

		oldNode := path[pathIdx]
		path[pathIdx] = nonSel[nonSelIdx]

		nonSel[nonSelIdx] = oldNode
	}

	return Solution{Path: path, Objective: objective(D, costs, path)}
}

// N3: Destroy-repair (remove nodes and rebuild with greedy)
func shakeDestroyRepair(D [][]int, costs []int, sol Solution, intensity int, rng *rand.Rand) Solution {
	path := append([]int(nil), sol.Path...)
	n := len(path)
	k := selectCount(len(D))

	if n < 4 {
		return Solution{Path: path, Objective: objective(D, costs, path)}
	}

	// Destroy: remove 20-30% of nodes
	destroyFraction := 0.2 + rng.Float64()*0.1
	numToRemove := int(math.Ceil(float64(n) * destroyFraction))
	if numToRemove >= n {
		numToRemove = n - 1
	}
	if numToRemove < 1 {
		numToRemove = 1
	}

	inSel := make([]bool, len(D))
	for _, v := range path {
		inSel[v] = true
	}

	// Remove random nodes
	removed := make(map[int]bool)
	for len(removed) < numToRemove && len(path) > 2 {
		idx := rng.Intn(len(path))
		node := path[idx]
		if !removed[node] {
			inSel[node] = false
			removed[node] = true
			path = append(path[:idx], path[idx+1:]...)
		}
	}

	// Repair: add nodes using greedy nearest neighbor any position
	nonSel := make([]int, 0, len(D)-len(path))
	for u := 0; u < len(D); u++ {
		if !inSel[u] {
			nonSel = append(nonSel, u)
		}
	}

	for len(path) < k && len(nonSel) > 0 {
		minIncrease := math.MaxInt32
		bestNode := -1
		bestPosition := -1

		for _, node := range nonSel {
			localMinIncrease := math.MaxInt32
			localBestPos := -1

			if len(path) == 0 {
				localMinIncrease = costs[node]
				localBestPos = 0
			} else if len(path) == 1 {
				inc := D[node][path[0]] + costs[node]
				if inc < localMinIncrease {
					localMinIncrease = inc
					localBestPos = 0
				}
				inc = D[path[0]][node] + costs[node]
				if inc < localMinIncrease {
					localMinIncrease = inc
					localBestPos = 1
				}
			} else {
				for pos := 0; pos <= len(path); pos++ {
					var inc int
					if pos == 0 {
						inc = D[node][path[0]] + costs[node]
					} else if pos == len(path) {
						inc = D[path[len(path)-1]][node] + costs[node]
					} else {
						a := path[pos-1]
						b := path[pos]
						deltaDist := D[a][node] + D[node][b] - D[a][b]
						inc = deltaDist + costs[node]
					}

					if inc < localMinIncrease {
						localMinIncrease = inc
						localBestPos = pos
					}
				}
			}

			if localMinIncrease < minIncrease {
				minIncrease = localMinIncrease
				bestNode = node
				bestPosition = localBestPos
			}
		}

		if bestNode != -1 {
			if bestPosition == len(path) {
				path = append(path, bestNode)
			} else {
				path = append(path[:bestPosition], append([]int{bestNode}, path[bestPosition:]...)...)
			}
			// Remove from nonSel
			newNonSel := make([]int, 0, len(nonSel)-1)
			for _, u := range nonSel {
				if u != bestNode {
					newNonSel = append(newNonSel, u)
				}
			}
			nonSel = newNonSel
		} else {
			break
		}
	}

	return Solution{Path: path, Objective: objective(D, costs, path)}
}

// N4: Double-bridge move (4-opt variant)
func shakeDoubleBridge(D [][]int, costs []int, sol Solution, intensity int, rng *rand.Rand) Solution {
	path := append([]int(nil), sol.Path...)
	n := len(path)

	if n < 8 {
		return shakeTwoOpt(D, costs, sol, intensity, rng)
	}

	minGap := n / 8
	maxAttempts := 10

	// Try multiple position sets before falling back
	for attempt := 0; attempt < maxAttempts; attempt++ {
		// Select 4 positions
		positions := make([]int, 4)
		for i := 0; i < 4; i++ {
			positions[i] = rng.Intn(n)
		}

		// Sort positions
		for i := 0; i < 3; i++ {
			for j := i + 1; j < 4; j++ {
				if positions[j] < positions[i] {
					positions[i], positions[j] = positions[j], positions[i]
				}
			}
		}

		// Check if gaps are sufficient
		gap1 := positions[1] - positions[0]
		gap2 := positions[2] - positions[1]
		gap3 := positions[3] - positions[2]
		gap4 := n - positions[3] + positions[0]

		if gap1 >= minGap && gap2 >= minGap && gap3 >= minGap && gap4 >= minGap {
			// Valid positions found - apply double-bridge
			newPath := make([]int, 0, n)
			newPath = append(newPath, path[:positions[0]+1]...)
			newPath = append(newPath, path[positions[2]+1:positions[3]+1]...)
			newPath = append(newPath, path[positions[1]+1:positions[2]+1]...)
			newPath = append(newPath, path[positions[0]+1:positions[1]+1]...)
			if positions[3] < n-1 {
				newPath = append(newPath, path[positions[3]+1:]...)
			}

			return Solution{Path: newPath, Objective: objective(D, costs, newPath)}
		}
	}

	// Fallback to 2-opt if no valid positions found
	return shakeTwoOpt(D, costs, sol, intensity, rng)
}

func min(a, b int, rest ...int) int {
	m := a
	if b < m {
		m = b
	}
	for _, v := range rest {
		if v < m {
			m = v
		}
	}
	return m
}
