package algorithms

import (
	"math"
	"math/rand"
	"time"
)

// LNSConfig holds configuration for Large Neighborhood Search
type LNSConfig struct {
	DestroyFraction float64       // Fraction of nodes to destroy (default 0.3)
	UseLocalSearch  bool          // Whether to use local search after repair
	TimeLimit       time.Duration // Time limit for the algorithm
	DestroyMethod   string        // Method: "weighted", "worst_edges", "shaw", "random_subpath"
}

// LNSResult contains the result of LNS execution
type LNSResult struct {
	BestSolution Solution
	Iterations   int
	Duration     time.Duration
}

// LargeNeighborhoodSearch implements LNS algorithm
func LargeNeighborhoodSearch(D [][]int, costs []int, config LNSConfig) LNSResult {
	startTime := time.Now()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	n := len(costs)
	k := (n + 1) / 2

	if config.DestroyFraction == 0 {
		config.DestroyFraction = 0.3
	}
	if config.DestroyMethod == "" {
		config.DestroyMethod = "worst_edges" // Default to best performing method
	}

	// Generate initial random solution
	currentSolution := startRandom(D, costs, rng)

	// Apply local search to initial solution
	currentSolution = localSearchSteepest(D, costs, currentSolution)

	iterations := 0

	for time.Since(startTime) < config.TimeLimit {
		iterations++

		// Destroy: remove nodes from current solution using selected method
		var partialPath []int
		switch config.DestroyMethod {
		case "worst_edges":
			partialPath = destroyWorstEdges(currentSolution.Path, config.DestroyFraction, D, costs, rng)
		case "shaw":
			partialPath = destroyShaw(currentSolution.Path, config.DestroyFraction, D, costs, rng)
		case "random_subpath":
			partialPath = destroyRandomSubpath(currentSolution.Path, config.DestroyFraction, rng)
		case "weighted":
			partialPath = destroy(currentSolution.Path, config.DestroyFraction, D, costs, rng)
		default:
			partialPath = destroyWorstEdges(currentSolution.Path, config.DestroyFraction, D, costs, rng)
		}

		// Repair: rebuild solution using nearest neighbor any position
		repairedSolution := repair(partialPath, D, costs, k, rng)

		// Optional local search after repair
		if config.UseLocalSearch {
			repairedSolution = localSearchSteepest(D, costs, repairedSolution)
		}

		// Accept if improved
		if repairedSolution.Objective < currentSolution.Objective {
			currentSolution = repairedSolution
		}
	}
	elapsed := time.Since(startTime)
	return LNSResult{
		BestSolution: currentSolution,
		Iterations:   iterations,
		Duration:     elapsed,
	}
}

// destroyWorstEdges removes nodes incident to the longest/most expensive edges
// This is typically the most effective for TSP-like problems
func destroyWorstEdges(path []int, fraction float64, D [][]int, costs []int, rng *rand.Rand) []int {
	numToRemove := int(math.Ceil(float64(len(path)) * fraction))
	if numToRemove >= len(path) {
		numToRemove = len(path) - 1
	}
	if numToRemove < 1 {
		numToRemove = 1
	}

	// Calculate edge costs (distance + average node cost of endpoints)
	type edgeInfo struct {
		fromIdx int
		toIdx   int
		cost    float64
	}

	edges := make([]edgeInfo, len(path))
	for i := 0; i < len(path); i++ {
		next := (i + 1) % len(path)
		edgeDist := float64(D[path[i]][path[next]])
		avgNodeCost := float64(costs[path[i]]+costs[path[next]]) / 2.0
		edges[i] = edgeInfo{
			fromIdx: i,
			toIdx:   next,
			cost:    edgeDist + avgNodeCost,
		}
	}

	// Sort edges by cost (descending)
	for i := 0; i < len(edges)-1; i++ {
		for j := i + 1; j < len(edges); j++ {
			if edges[j].cost > edges[i].cost {
				edges[i], edges[j] = edges[j], edges[i]
			}
		}
	}

	// Select edges probabilistically, favoring worse edges
	removed := make(map[int]bool)
	edgeIdx := 0

	for len(removed) < numToRemove && edgeIdx < len(edges) {
		// Exponential decay probability: worse edges have higher probability
		prob := math.Exp(-float64(edgeIdx) / float64(len(edges)) * 3.0)

		if rng.Float64() < prob {
			// Remove one of the nodes from this edge (randomly choose which)
			if rng.Float64() < 0.5 {
				if !removed[edges[edgeIdx].fromIdx] {
					removed[edges[edgeIdx].fromIdx] = true
				}
			} else {
				if !removed[edges[edgeIdx].toIdx] {
					removed[edges[edgeIdx].toIdx] = true
				}
			}
		}
		edgeIdx++

		// Reset if we've gone through all edges
		if edgeIdx >= len(edges) {
			edgeIdx = 0
		}
	}

	// Build partial solution
	partial := make([]int, 0, len(path)-numToRemove)
	for i := 0; i < len(path); i++ {
		if !removed[i] {
			partial = append(partial, path[i])
		}
	}

	return partial
}

// destroyShaw removes related nodes based on Shaw removal heuristic
// Nodes that are similar (close in space and have similar costs) are removed together
func destroyShaw(path []int, fraction float64, D [][]int, costs []int, rng *rand.Rand) []int {
	numToRemove := int(math.Ceil(float64(len(path)) * fraction))
	if numToRemove >= len(path) {
		numToRemove = len(path) - 1
	}
	if numToRemove < 1 {
		numToRemove = 1
	}

	// Randomly select a seed node
	seedIdx := rng.Intn(len(path))
	seedNode := path[seedIdx]

	// Calculate relatedness of all other nodes to seed
	type nodeRelatedness struct {
		idx         int
		relatedness float64
	}

	relatedness := make([]nodeRelatedness, 0, len(path)-1)
	for i := 0; i < len(path); i++ {
		if i == seedIdx {
			continue
		}

		node := path[i]

		// Relatedness based on:
		// 1. Distance between nodes
		// 2. Cost similarity
		dist := float64(D[seedNode][node])
		costDiff := math.Abs(float64(costs[seedNode] - costs[node]))

		// Lower is more related
		rel := dist + costDiff
		relatedness = append(relatedness, nodeRelatedness{i, rel})
	}

	// Sort by relatedness (ascending - most related first)
	for i := 0; i < len(relatedness)-1; i++ {
		for j := i + 1; j < len(relatedness); j++ {
			if relatedness[j].relatedness < relatedness[i].relatedness {
				relatedness[i], relatedness[j] = relatedness[j], relatedness[i]
			}
		}
	}

	// Remove seed node and most related nodes
	removed := make(map[int]bool)
	removed[seedIdx] = true

	for i := 0; i < len(relatedness) && len(removed) < numToRemove; i++ {
		// Exponential probability: more related = higher probability
		prob := math.Exp(-float64(i) / float64(len(relatedness)) * 4.0)
		if rng.Float64() < prob {
			removed[relatedness[i].idx] = true
		}
	}

	// Ensure we remove exactly numToRemove nodes
	for i := 0; i < len(relatedness) && len(removed) < numToRemove; i++ {
		removed[relatedness[i].idx] = true
	}

	// Build partial solution
	partial := make([]int, 0, len(path)-len(removed))
	for i := 0; i < len(path); i++ {
		if !removed[i] {
			partial = append(partial, path[i])
		}
	}

	return partial
}

// destroyRandomSubpath removes a continuous segment of the path
// This maintains some structure while creating a large gap to fill
func destroyRandomSubpath(path []int, fraction float64, rng *rand.Rand) []int {
	numToRemove := int(math.Ceil(float64(len(path)) * fraction))
	if numToRemove >= len(path) {
		numToRemove = len(path) - 1
	}
	if numToRemove < 1 {
		numToRemove = 1
	}

	// Randomly select starting position
	startPos := rng.Intn(len(path))

	// Build partial solution by skipping the segment
	partial := make([]int, 0, len(path)-numToRemove)
	for i := 0; i < len(path); i++ {
		skipPos := (startPos + i) % len(path)
		if i >= numToRemove {
			partial = append(partial, path[skipPos])
		}
	}

	return partial
}

// destroy removes a fraction of nodes from the solution
// Nodes with higher costs and longer edges have higher probability of removal
func destroy(path []int, fraction float64, D [][]int, costs []int, rng *rand.Rand) []int {
	numToRemove := int(math.Ceil(float64(len(path)) * fraction))
	if numToRemove >= len(path) {
		numToRemove = len(path) - 1
	}
	if numToRemove < 1 {
		numToRemove = 1
	}

	// Calculate removal weights based on edge lengths and node costs
	weights := make([]float64, len(path))
	for i := 0; i < len(path); i++ {
		prev := (i - 1 + len(path)) % len(path)
		next := (i + 1) % len(path)

		// Weight based on adjacent edge lengths and node cost
		edgeLength1 := D[path[prev]][path[i]]
		edgeLength2 := D[path[i]][path[next]]
		avgEdgeLength := float64(edgeLength1+edgeLength2) / 2.0
		cost := float64(costs[path[i]])

		// Higher weight for longer edges and higher costs
		weights[i] = avgEdgeLength + cost
	}

	// Normalize weights to probabilities
	totalWeight := 0.0
	for _, w := range weights {
		totalWeight += w
	}
	if totalWeight > 0 {
		for i := range weights {
			weights[i] /= totalWeight
		}
	}

	// Select nodes to remove using weighted random selection
	removed := make(map[int]bool)
	for len(removed) < numToRemove {
		// Weighted random selection
		r := rng.Float64()
		for i := 0; i < len(path); i++ {
			if removed[i] {
				continue
			}
			if r <= weights[i] {
				removed[i] = true
				break
			}
		}
	}

	// Build partial solution
	partial := make([]int, 0, len(path)-numToRemove)
	for i := 0; i < len(path); i++ {
		if !removed[i] {
			partial = append(partial, path[i])
		}
	}

	return partial
}

// repair rebuilds the solution using nearest neighbor any position heuristic
func repair(partialPath []int, D [][]int, costs []int, targetSize int, rng *rand.Rand) Solution {
	n := len(costs)

	// Create set of nodes already in solution
	inSolution := make(map[int]bool)
	for _, node := range partialPath {
		inSolution[node] = true
	}

	// Create list of unvisited nodes
	unvisited := make(map[int]bool)
	for i := 0; i < n; i++ {
		if !inSolution[i] {
			unvisited[i] = true
		}
	}

	path := make([]int, len(partialPath))
	copy(path, partialPath)

	// Add nodes using greedy nearest neighbor any position
	for len(path) < targetSize && len(unvisited) > 0 {
		minIncrease := math.MaxInt32
		bestNode := -1
		bestPosition := -1

		for node := range unvisited {
			localMinIncrease := math.MaxInt32
			localBestPos := -1

			if len(path) == 0 {
				localMinIncrease = costs[node]
				localBestPos = 0
			} else if len(path) == 1 {
				// Insert before or after
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
				// Try all positions
				for pos := 0; pos <= len(path); pos++ {
					var inc int
					if pos == 0 {
						// Insert at beginning
						inc = D[node][path[0]] + costs[node]
					} else if pos == len(path) {
						// Insert at end
						inc = D[path[len(path)-1]][node] + costs[node]
					} else {
						// Insert in middle
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
			// Insert node at best position
			if bestPosition == len(path) {
				path = append(path, bestNode)
			} else {
				path = append(path[:bestPosition], append([]int{bestNode}, path[bestPosition:]...)...)
			}
			delete(unvisited, bestNode)
		} else {
			break
		}
	}

	return Solution{Path: path, Objective: objective(D, costs, path)}
}
