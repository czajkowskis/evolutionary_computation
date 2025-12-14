package algorithms

import (
	"math/rand"
	"time"
)

// HybridConfig contains configuration for the hybrid algorithm
type HybridConfig struct {
	PopulationSize int
	TimeLimit      time.Duration
	UseLocalSearch bool
	Operator       int // 1 or 2
	Seed           int64
}

// HybridEvolutionary runs the hybrid evolutionary algorithm
func HybridEvolutionary(D [][]int, costs []int, config HybridConfig) Solution {
	rng := rand.New(rand.NewSource(config.Seed))
	n := len(costs)
	targetSize := (n + 1) / 2

	startTime := time.Now()

	// Initialize population
	population := initializePopulation(D, costs, targetSize, config.PopulationSize, rng)

	bestSolution := population[0]
	for _, sol := range population {
		if sol.Objective < bestSolution.Objective {
			bestSolution = sol
		}
	}

	iterations := 0
	for time.Since(startTime) < config.TimeLimit {
		// Select two parents uniformly at random
		parent1 := population[rng.Intn(len(population))]
		parent2 := population[rng.Intn(len(population))]

		// Apply recombination
		var offspring Solution
		if config.Operator == 1 {
			offspring = recombineOperator1(parent1, parent2, D, costs, targetSize, rng)
		} else {
			offspring = recombineOperator2(parent1, parent2, D, costs, targetSize, rng)
		}

		// Apply local search if enabled
		if config.UseLocalSearch {
			offspring = localSearchSteepest(D, costs, offspring)
		}

		// Update population if offspring is better and not duplicate
		if !isDuplicate(offspring, population) {
			worstIdx := findWorstIndex(population)
			if offspring.Objective < population[worstIdx].Objective {
				population[worstIdx] = offspring

				if offspring.Objective < bestSolution.Objective {
					bestSolution = offspring
				}
			}
		}

		iterations++
	}

	return bestSolution
}

// initializePopulation creates initial population using random start + local search
func initializePopulation(D [][]int, costs []int, targetSize, popSize int, rng *rand.Rand) []Solution {
	population := make([]Solution, 0, popSize)
	n := len(costs)

	for len(population) < popSize {
		// Create random initial solution
		sol := randomConstruction(D, costs, n, targetSize, rng)

		// Apply local search
		sol = localSearchSteepest(D, costs, sol)

		// Add if not duplicate
		if !isDuplicate(sol, population) {
			population = append(population, sol)
		}
	}

	return population
}

// recombineOperator1 implements the common edges/nodes operator
func recombineOperator1(parent1, parent2 Solution, D [][]int, costs []int, targetSize int, rng *rand.Rand) Solution {
	n := len(costs)

	// Find common nodes
	nodesP1 := make(map[int]bool)
	for _, node := range parent1.Path {
		nodesP1[node] = true
	}

	commonNodes := make([]int, 0)
	for _, node := range parent2.Path {
		if nodesP1[node] {
			commonNodes = append(commonNodes, node)
		}
	}

	// Find common edges and build subpaths
	subpaths := findCommonSubpaths(parent1, parent2, commonNodes)

	// Randomly reverse some subpaths
	for i := range subpaths {
		if rng.Float64() < 0.5 {
			reverseSlice(subpaths[i])
		}
	}

	// Get nodes not in common subpaths
	inSubpaths := make(map[int]bool)
	for _, subpath := range subpaths {
		for _, node := range subpath {
			inSubpaths[node] = true
		}
	}

	// Get all available nodes (not yet in subpaths)
	availableNodes := make([]int, 0)
	for i := 0; i < n; i++ {
		if !inSubpaths[i] {
			availableNodes = append(availableNodes, i)
		}
	}

	// Shuffle available nodes
	rng.Shuffle(len(availableNodes), func(i, j int) {
		availableNodes[i], availableNodes[j] = availableNodes[j], availableNodes[i]
	})

	// Create single-node subpaths from random nodes
	numNodesNeeded := targetSize - len(inSubpaths)
	for i := 0; i < numNodesNeeded && i < len(availableNodes); i++ {
		subpaths = append(subpaths, []int{availableNodes[i]})
	}

	// Shuffle all subpaths (common + random nodes) together
	rng.Shuffle(len(subpaths), func(i, j int) {
		subpaths[i], subpaths[j] = subpaths[j], subpaths[i]
	})

	// Concatenate all subpaths to form final path
	path := make([]int, 0, targetSize)
	for _, subpath := range subpaths {
		path = append(path, subpath...)
		if len(path) >= targetSize {
			break
		}
	}

	// Trim to exact targetSize if needed
	if len(path) > targetSize {
		path = path[:targetSize]
	}

	// If still short, add remaining random nodes
	if len(path) < targetSize {
		inSolution := make(map[int]bool)
		for _, node := range path {
			inSolution[node] = true
		}

		for i := 0; i < n && len(path) < targetSize; i++ {
			if !inSolution[i] {
				path = append(path, i)
				inSolution[i] = true
			}
		}
	}

	return Solution{Path: path, Objective: objective(D, costs, path)}
}

// findCommonSubpaths identifies common subpaths between two parents
func findCommonSubpaths(parent1, parent2 Solution, commonNodes []int) [][]int {
	if len(commonNodes) == 0 {
		return [][]int{}
	}

	// Build edge map for parent2
	edgesP2 := make(map[[2]int]bool)
	for i := 0; i < len(parent2.Path)-1; i++ {
		edgesP2[[2]int{parent2.Path[i], parent2.Path[i+1]}] = true
		edgesP2[[2]int{parent2.Path[i+1], parent2.Path[i]}] = true // undirected
	}

	// Find subpaths in parent1 that exist in parent2
	subpaths := make([][]int, 0)
	currentSubpath := make([]int, 0)

	commonSet := make(map[int]bool)
	for _, node := range commonNodes {
		commonSet[node] = true
	}

	for i := 0; i < len(parent1.Path); i++ {
		node := parent1.Path[i]
		if !commonSet[node] {
			if len(currentSubpath) > 0 {
				subpaths = append(subpaths, currentSubpath)
				currentSubpath = make([]int, 0)
			}
			continue
		}

		if len(currentSubpath) == 0 {
			currentSubpath = append(currentSubpath, node)
		} else {
			// Check if edge exists
			prev := currentSubpath[len(currentSubpath)-1]
			if edgesP2[[2]int{prev, node}] {
				currentSubpath = append(currentSubpath, node)
			} else {
				subpaths = append(subpaths, currentSubpath)
				currentSubpath = []int{node}
			}
		}
	}

	if len(currentSubpath) > 0 {
		subpaths = append(subpaths, currentSubpath)
	}

	return subpaths
}

// recombineOperator2 implements the parent-based repair operator
func recombineOperator2(parent1, parent2 Solution, D [][]int, costs []int, targetSize int, rng *rand.Rand) Solution {
	// Choose one parent as base
	var baseParent Solution
	if rng.Float64() < 0.5 {
		baseParent = parent1
	} else {
		baseParent = parent2
	}

	otherParent := parent1
	if baseParent.Path[0] == parent1.Path[0] && len(baseParent.Path) == len(parent1.Path) {
		otherParent = parent2
	}

	// Create set of nodes in other parent
	otherNodes := make(map[int]bool)
	for _, node := range otherParent.Path {
		otherNodes[node] = true
	}

	// Keep only common nodes in order
	partialPath := make([]int, 0)
	for _, node := range baseParent.Path {
		if otherNodes[node] {
			partialPath = append(partialPath, node)
		}
	}

	// Repair using greedy heuristic
	return repair(partialPath, D, costs, targetSize, rng)
}
