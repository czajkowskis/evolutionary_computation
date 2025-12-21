package main

import (
	"fmt"
	"log"
	"math/rand"
	"path/filepath"
	"time"

	"github.com/czajkowskis/evolutionary_computation/01_labs/greedy_heuristics/pkg/algorithms"
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
		return algorithms.RandomSolution(distanceMatrix, nodeCosts, startNodeIndices)
	})
	solutionSets["Random_Solution"] = solutions
	executionTimes["Random_Solution"] = elapsed

	solutions, elapsed = measureExecutionTime(func() []commonAlgorithms.Solution {
		return algorithms.NearestNeighborEnd(distanceMatrix, nodeCosts, startNodeIndices)
	})
	solutionSets["Nearest_Neighbor_End_Only"] = solutions
	executionTimes["Nearest_Neighbor_End_Only"] = elapsed

	solutions, elapsed = measureExecutionTime(func() []commonAlgorithms.Solution {
		return algorithms.NearestNeighborAny(distanceMatrix, nodeCosts, startNodeIndices)
	})
	solutionSets["Nearest_Neighbor_Any_Position"] = solutions
	executionTimes["Nearest_Neighbor_Any_Position"] = elapsed

	solutions, elapsed = measureExecutionTime(func() []commonAlgorithms.Solution {
		return algorithms.GreedyCycle(distanceMatrix, nodeCosts, startNodeIndices)
	})
	solutionSets["Greedy_Cycle"] = solutions
	executionTimes["Greedy_Cycle"] = elapsed

	for name, solutions := range solutionSets {
		if len(solutions) > 0 {
			min, max, avg := utils.CalculateStatistics(solutions)
			avgTime := float64(executionTimes[name].Nanoseconds()) / float64(numSolutions) / 1e6
			// fmt.Printf("%s: min = %d, max = %d, average = %.2f, avg_time = %.4f ms\n", name, min, max, avg, avgTime)
			fmt.Printf("%s: %.2f(%d,%d), avg_time = %.4f ms\n", name, avg, min, max, avgTime)

			bestSolution := commonAlgorithms.FindBestSolution(solutions)
			fmt.Printf("Best path: %v\n", bestSolution.Path)
			name_to_title := map[string]string{
				"Random_Solution":               "Random Solution",
				"Nearest_Neighbor_End_Only":     "Nearest Neighbor (End Only)",
				"Nearest_Neighbor_Any_Position": "Nearest Neighbor (Any Position)",
				"Greedy_Cycle":                  "Greedy Cycle",
			}

			plotTitle := fmt.Sprintf("Best %s Solution for Instance %s", name_to_title[name], instanceName)
			plotFileName := utils.SanitizeFileName(fmt.Sprintf("Best_%s_Solution_%s", name, instanceName))
			plotBounds := config.DefaultPlotBounds
			outputDir := filepath.Join("output", "01_labs", "greedy_heuristics", "plots")

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
