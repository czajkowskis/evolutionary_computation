package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/czajkowskis/evolutionary_computation/02_labs/greedy_regret_heuristics/pkg/algorithms"
	"github.com/czajkowskis/evolutionary_computation/02_labs/greedy_regret_heuristics/pkg/data"
	"github.com/czajkowskis/evolutionary_computation/02_labs/greedy_regret_heuristics/pkg/utils"
	"github.com/czajkowskis/evolutionary_computation/02_labs/greedy_regret_heuristics/pkg/visualisation"
)

func measureExecutionTime(algorithm func() []algorithms.Solution) ([]algorithms.Solution, time.Duration) {
	start := time.Now()
	solutions := algorithm()
	elapsed := time.Since(start)
	return solutions, elapsed
}

func processInstance(instanceName string, nodes []data.Node) {
	fmt.Printf("Instance %s Statistics:\n", instanceName)

	distanceMatrix := data.CalculateDistanceMatrix(nodes)
	startNodeIndices := utils.GenerateStartNodeIndices(len(nodes))
	numSolutions := len(startNodeIndices)

	nodeCosts := make([]int, len(nodes))
	for i, node := range nodes {
		nodeCosts[i] = node.Cost
	}

	// Apply algorithms
	solutionSets := make(map[string][]algorithms.Solution)
	executionTimes := make(map[string]time.Duration)

	var solutions []algorithms.Solution
	var elapsed time.Duration

	solutions, elapsed = measureExecutionTime(func() []algorithms.Solution {
		return algorithms.NearestNeighborWeightedTwoRegret(distanceMatrix, nodeCosts, startNodeIndices, 1, 0)
	})
	solutionSets["Nearest_Neighbor_Two_Regret"] = solutions
	executionTimes["Nearest_Neighbor_Two_Regret"] = elapsed

	solutions, elapsed = measureExecutionTime(func() []algorithms.Solution {
		return algorithms.GreedyCycleWeightedTwoRegret(distanceMatrix, nodeCosts, startNodeIndices, 1, 0)
	})
	solutionSets["Greedy_Cycle_Two_Regret"] = solutions
	executionTimes["Greedy_Cycle_Two_Regret"] = elapsed

	solutions, elapsed = measureExecutionTime(func() []algorithms.Solution {
		return algorithms.NearestNeighborWeightedTwoRegret(distanceMatrix, nodeCosts, startNodeIndices, 0.5, 0.5)
	})
	solutionSets["Nearest_Neighbor_Weighted_Sum"] = solutions
	executionTimes["Nearest_Neighbor_Weighted_Sum"] = elapsed

	solutions, elapsed = measureExecutionTime(func() []algorithms.Solution {
		return algorithms.GreedyCycleWeightedTwoRegret(distanceMatrix, nodeCosts, startNodeIndices, 0.5, 0.5)
	})
	solutionSets["Greedy_Cycle_Weighted_Sum"] = solutions
	executionTimes["Greedy_Cycle_Weighted_Sum"] = elapsed

	for name, solutions := range solutionSets {
		if len(solutions) > 0 {
			min, max, avg := utils.CalculateStatistics(solutions)
			avgTime := float64(executionTimes[name].Nanoseconds()) / float64(numSolutions) / 1e6
			// fmt.Printf("%s: min = %d, max = %d, average = %.2f, avg_time = %.4f ms\n", name, min, max, avg, avgTime)
			fmt.Printf("%s: %.2f(%d,%d), avg_time = %.4f ms\n", name, avg, min, max, avgTime)

			bestSolution := algorithms.FindBestSolution(solutions)
			fmt.Printf("Best path: %v\n", bestSolution.Path)
			name_to_title := map[string]string{
				"Nearest_Neighbor_Two_Regret":   "Nearest Neighbor (2-Regret)",
				"Greedy_Cycle_Two_Regret":       "Greedy Cycle (2-Regret)",
				"Nearest_Neighbor_Weighted_Sum": "Nearest Neighbor (Weighted Sum)",
				"Greedy_Cycle_Weighted_Sum":     "Greedy Cycle (Weighted Sum)",
			}

			plotTitle := fmt.Sprintf("Best %s Solution for Instance %s", name_to_title[name], instanceName)
			plotFileName := fmt.Sprintf("Best_%s_Solution_%s", name, instanceName)

			if err := visualisation.PlotSolution(nodes, bestSolution.Path, plotTitle, plotFileName, 0, 4000, 0, 2000); err != nil {
				log.Printf("Error plotting best solution for %s on instance %s: %v", name, instanceName, err)
			}
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	// Read nodes from CSV files
	nodesA, err := data.ReadNodes("./instances/TSPA.csv")
	if err != nil {
		log.Fatalf("Error reading TSPA.csv: %v", err)
	}

	nodesB, err := data.ReadNodes("./instances/TSPB.csv")
	if err != nil {
		log.Fatalf("Error reading TSPB.csv: %v", err)
	}

	processInstance("A", nodesA)
	fmt.Println()
	processInstance("B", nodesB)
}
