package main

import (
	"fmt"
	"log"
	"math/rand"
	"path/filepath"
	"time"

	"github.com/czajkowskis/evolutionary_computation/02_labs/greedy_regret_heuristics/pkg/algorithms"
	commonAlgorithms "github.com/czajkowskis/evolutionary_computation/pkg/common/algorithms"
	"github.com/czajkowskis/evolutionary_computation/pkg/common/config"
	"github.com/czajkowskis/evolutionary_computation/pkg/common/data"
	"github.com/czajkowskis/evolutionary_computation/pkg/common/utils"
	"github.com/czajkowskis/evolutionary_computation/pkg/common/visualisation"
)

func measureExecutionTime(algorithm func() []commonAlgorithms.Solution) ([]commonAlgorithms.Solution, time.Duration) {
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
	solutionSets := make(map[string][]commonAlgorithms.Solution)
	executionTimes := make(map[string]time.Duration)

	var solutions []commonAlgorithms.Solution
	var elapsed time.Duration

	solutions, elapsed = measureExecutionTime(func() []commonAlgorithms.Solution {
		return algorithms.NearestNeighborWeightedTwoRegret(distanceMatrix, nodeCosts, startNodeIndices, 1, 0)
	})
	solutionSets["Nearest_Neighbor_Two_Regret"] = solutions
	executionTimes["Nearest_Neighbor_Two_Regret"] = elapsed

	solutions, elapsed = measureExecutionTime(func() []commonAlgorithms.Solution {
		return algorithms.GreedyCycleWeightedTwoRegret(distanceMatrix, nodeCosts, startNodeIndices, 1, 0)
	})
	solutionSets["Greedy_Cycle_Two_Regret"] = solutions
	executionTimes["Greedy_Cycle_Two_Regret"] = elapsed

	solutions, elapsed = measureExecutionTime(func() []commonAlgorithms.Solution {
		return algorithms.NearestNeighborWeightedTwoRegret(distanceMatrix, nodeCosts, startNodeIndices, 0.5, 0.5)
	})
	solutionSets["Nearest_Neighbor_Weighted_Sum"] = solutions
	executionTimes["Nearest_Neighbor_Weighted_Sum"] = elapsed

	solutions, elapsed = measureExecutionTime(func() []commonAlgorithms.Solution {
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

			bestSolution := commonAlgorithms.FindBestSolution(solutions)
			fmt.Printf("Best path: %v\n", bestSolution.Path)
			name_to_title := map[string]string{
				"Nearest_Neighbor_Two_Regret":   "Nearest Neighbor (2-Regret)",
				"Greedy_Cycle_Two_Regret":       "Greedy Cycle (2-Regret)",
				"Nearest_Neighbor_Weighted_Sum": "Nearest Neighbor (Weighted Sum)",
				"Greedy_Cycle_Weighted_Sum":     "Greedy Cycle (Weighted Sum)",
			}

			plotTitle := fmt.Sprintf("Best %s Solution for Instance %s", name_to_title[name], instanceName)
			plotFileName := utils.SanitizeFileName(fmt.Sprintf("Best_%s_Solution_%s", name, instanceName))
			plotBounds := config.DefaultPlotBounds
			outputDir := filepath.Join("output", "02_labs", "greedy_regret_heuristics", "plots")

			if err := visualisation.PlotSolution(nodes, bestSolution.Path, plotTitle, plotFileName,
				plotBounds.XMin, plotBounds.XMax, plotBounds.YMin, plotBounds.YMax, outputDir); err != nil {
				log.Printf("Error plotting best solution for %s on instance %s: %v", name, instanceName, err)
			}
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	// Read nodes from CSV files
	instancePaths := config.DefaultInstancePaths()
	nodesA, err := data.ReadNodes(instancePaths.TSPA)
	if err != nil {
		log.Fatalf("Error reading TSPA.csv: %v", err)
	}

	nodesB, err := data.ReadNodes(instancePaths.TSPB)
	if err != nil {
		log.Fatalf("Error reading TSPB.csv: %v", err)
	}

	processInstance("A", nodesA)
	fmt.Println()
	processInstance("B", nodesB)
}
