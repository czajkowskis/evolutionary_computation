package main

import (
	"fmt"
	"log"
	"math/rand"
	"path/filepath"
	"time"

	"github.com/czajkowskis/evolutionary_computation/07_labs/large_neighborhood_search/pkg/algorithms"
	commonAlgorithms "github.com/czajkowskis/evolutionary_computation/pkg/common/algorithms"
	"github.com/czajkowskis/evolutionary_computation/pkg/common/config"
	"github.com/czajkowskis/evolutionary_computation/pkg/common/data"
	"github.com/czajkowskis/evolutionary_computation/pkg/common/utils"
	"github.com/czajkowskis/evolutionary_computation/pkg/common/visualisation"
)

// Configuration constants
const (
	numLNSRuns = 20      // Number of LNS runs per instance
	timeLimitA = 3276.57 // Average running time of MSLS from the previous assignment for instance A
	timeLimitB = 2342.11 // Average running time of MSLS from the previous assignment for instance B
)

// InstanceResults stores all results for a single instance
type InstanceResults struct {
	Instance   string
	LNSResults []algorithms.LNSResult
}

// processInstance runs the full experimental pipeline for a single instance
func processInstance(instanceName string, nodes []data.Node) {
	log.Printf("Processing instance %s with %d nodes", instanceName, len(nodes))
	fmt.Printf("Instance %s Statistics:\n", instanceName)

	var timeLimit time.Duration
	if instanceName == "A" {
		timeLimit = time.Duration(timeLimitA * float64(time.Millisecond))
	} else {
		timeLimit = time.Duration(timeLimitB * float64(time.Millisecond))
	}

	D := data.CalculateDistanceMatrix(nodes)
	costs := make([]int, len(nodes))
	for i, node := range nodes {
		costs[i] = node.Cost
	}

	var rows []utils.Row

	// Define destroy methods to test
	destroyMethods := []string{"random_subpath"}

	// === PHASE 1: Run LNS using local search after destroy-repair operators ===
	for _, method := range destroyMethods {
		log.Printf("Starting LNS with LS using %s destroy method for instance %s", method, instanceName)
		start := time.Now()

		var lnsResults []algorithms.LNSResult
		totalLNSIterations := 0
		for run := 0; run < numLNSRuns; run++ {
			lnsResult := algorithms.LargeNeighborhoodSearch(D, costs, algorithms.LNSConfig{
				DestroyFraction: 0.3,
				UseLocalSearch:  true,
				TimeLimit:       timeLimit,
				DestroyMethod:   method,
			})
			lnsResults = append(lnsResults, lnsResult)
			totalLNSIterations += lnsResult.Iterations
		}

		totalLNSTime := time.Since(start)
		avgLNSTime := totalLNSTime / time.Duration(numLNSRuns)

		// Collect LNS solutions for statistics
		lnsSolutions := make([]commonAlgorithms.Solution, len(lnsResults))
		for i, r := range lnsResults {
			lnsSolutions[i] = r.BestSolution
		}

		lnsMin, lnsMax, lnsAvg := utils.CalculateStatistics(lnsSolutions)
		avgLNSTimeMs := float64(avgLNSTime.Nanoseconds()) / 1e6
		avgLNSIterations := float64(totalLNSIterations) / float64(numLNSRuns)
		bestLNS := commonAlgorithms.FindBestSolution(lnsSolutions)

		rows = append(rows, utils.Row{
			Name:        fmt.Sprintf("LNS+LS (%s)", method),
			AvgV:        lnsAvg,
			MinV:        lnsMin,
			MaxV:        lnsMax,
			AvgTms:      avgLNSTimeMs,
			AvgLNSIters: avgLNSIterations,
			BestPath:    bestLNS.Path,
			BestValue:   bestLNS.Objective,
		})

		log.Printf("Completed LNS+LS (%s): best value %d, avg time %.2f ms, avg iterations %.1f", method, bestLNS.Objective, avgLNSTimeMs, avgLNSIterations)
	}

	// === PHASE 2: Run LNS not using local search after destroy-repair operators ===
	for _, method := range destroyMethods {
		log.Printf("Starting LNS without LS using %s destroy method for instance %s", method, instanceName)
		start := time.Now()

		var lnsResults []algorithms.LNSResult
		totalLNSIterations := 0
		for run := 0; run < numLNSRuns; run++ {
			lnsResult := algorithms.LargeNeighborhoodSearch(D, costs, algorithms.LNSConfig{
				DestroyFraction: 0.3,
				UseLocalSearch:  false,
				TimeLimit:       timeLimit,
				DestroyMethod:   method,
			})
			lnsResults = append(lnsResults, lnsResult)
			totalLNSIterations += lnsResult.Iterations
		}

		totalLNSTime := time.Since(start)
		avgLNSTime := totalLNSTime / time.Duration(numLNSRuns)

		// Collect LNS solutions for statistics
		lnsSolutions := make([]commonAlgorithms.Solution, len(lnsResults))
		for i, r := range lnsResults {
			lnsSolutions[i] = r.BestSolution
		}

		lnsMin, lnsMax, lnsAvg := utils.CalculateStatistics(lnsSolutions)
		avgLNSTimeMs := float64(avgLNSTime.Nanoseconds()) / 1e6
		avgLNSIterations := float64(totalLNSIterations) / float64(numLNSRuns)
		bestLNS := commonAlgorithms.FindBestSolution(lnsSolutions)

		rows = append(rows, utils.Row{
			Name:        fmt.Sprintf("LNS-LS (%s)", method),
			AvgV:        lnsAvg,
			MinV:        lnsMin,
			MaxV:        lnsMax,
			AvgTms:      avgLNSTimeMs,
			AvgLNSIters: avgLNSIterations,
			BestPath:    bestLNS.Path,
			BestValue:   bestLNS.Objective,
		})

		log.Printf("Completed LNS-LS (%s): best value %d, avg time %.2f ms, avg iterations %.1f", method, bestLNS.Objective, avgLNSTimeMs, avgLNSIterations)
	}

	// Print console output
	fmt.Println("\nObjective value: av (min, max)")
	for _, r := range rows {
		fmt.Printf("%-34s  %.2f (%d, %d)\n", r.Name, r.AvgV, r.MinV, r.MaxV)
	}
	fmt.Println()

	fmt.Println("Average time per run [ms]:")
	for _, r := range rows {
		fmt.Printf("%-34s  %.4f\n", r.Name, r.AvgTms)
	}
	fmt.Println()

	fmt.Println("Average iterations per run:")
	for _, r := range rows {
		fmt.Printf("%-34s  %.1f\n", r.Name, r.AvgLNSIters)
	}

	// Save CSV
	outputDir := filepath.Join("output", "07_labs", "large_neighborhood_search", "results")
	if err := utils.WriteResultsCSV(instanceName, rows, outputDir); err != nil {
		log.Printf("CSV write error for instance %s: %v", instanceName, err)
	} else {
		log.Printf("CSV results saved for instance %s", instanceName)
	}

	// Plot best solutions for each method
	plotBounds := config.DefaultPlotBounds
	plotDir := filepath.Join("output", "07_labs", "large_neighborhood_search", "plots")
	for i, r := range rows {
		title := fmt.Sprintf("%s for Instance %s (Value: %d)", r.Name, instanceName, r.BestValue)
		fileName := utils.SanitizeFileName(fmt.Sprintf("%s_Instance_%s_%d", r.Name, instanceName, i))
		if err := visualisation.PlotSolution(nodes, r.BestPath, title, fileName,
			plotBounds.XMin, plotBounds.XMax, plotBounds.YMin, plotBounds.YMax, plotDir); err != nil {
			log.Printf("plot error for %s/%s: %v", instanceName, r.Name, err)
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	log.Println("Starting LNS local search experiments")

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

	log.Println("Program execution completed")
}
