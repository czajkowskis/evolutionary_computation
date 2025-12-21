package main

import (
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/czajkowskis/evolutionary_computation/09_labs/hybrid_evolutionary_algorithm/pkg/algorithms"
	commonAlgorithms "github.com/czajkowskis/evolutionary_computation/pkg/common/algorithms"
	"github.com/czajkowskis/evolutionary_computation/pkg/common/config"
	"github.com/czajkowskis/evolutionary_computation/pkg/common/data"
	"github.com/czajkowskis/evolutionary_computation/pkg/common/utils"
	"github.com/czajkowskis/evolutionary_computation/pkg/common/visualisation"
)

// Configuration constants
const (
	numRuns        = 20   // Number of runs per configuration
	timeLimitAMs   = 3277 // Average running time of MSLS for instance A (ms)
	timeLimitBMs   = 2342 // Average running time of MSLS for instance B (ms)
	populationSize = 20   // Elite population size as per requirements
)

// processInstance runs the full experimental pipeline for a single instance
func processInstance(instanceName string, nodes []data.Node) {
	log.Printf("Processing instance %s with %d nodes", instanceName, len(nodes))
	fmt.Printf("\n========================================\n")
	fmt.Printf("Instance %s Statistics:\n", instanceName)
	fmt.Printf("========================================\n")

	var timeLimitMs int
	if instanceName == "A" {
		timeLimitMs = timeLimitAMs
	} else {
		timeLimitMs = timeLimitBMs
	}
	timeLimit := time.Duration(timeLimitMs) * time.Millisecond

	D := data.CalculateDistanceMatrix(nodes)
	costs := make([]int, len(nodes))
	for i, node := range nodes {
		costs[i] = node.Cost
	}

	var rows []utils.Row

	// Define hybrid algorithm configurations to test
	configs := []struct {
		name           string
		operator       int
		useLocalSearch bool
	}{
		{"op_1_with_LS", 1, true},
		{"op_1_without_LS", 1, false},
		{"op_2_with_LS", 2, true},
		{"op_2_without_LS", 2, false},
	}

	// Run experiments for each configuration
	for _, cfg := range configs {
		log.Printf("Starting %s for instance %s", cfg.name, instanceName)
		start := time.Now()

		var solutions []commonAlgorithms.Solution
		totalIterations := 0

		for run := 0; run < numRuns; run++ {
			seed := time.Now().UnixNano() + int64(run)

			hybridConfig := algorithms.HybridConfig{
				PopulationSize: populationSize,
				TimeLimit:      timeLimit,
				UseLocalSearch: cfg.useLocalSearch,
				Operator:       cfg.operator,
				Seed:           seed,
			}

			result := algorithms.HybridEvolutionary(D, costs, hybridConfig)
			solutions = append(solutions, result.Solution)
			totalIterations += result.Iterations

			if run%5 == 0 {
				log.Printf("  Run %d/%d completed: objective = %d, iterations = %d",
					run+1, numRuns, result.Solution.Objective, result.Iterations)
			}
		}

		totalTime := time.Since(start)
		avgTime := totalTime / time.Duration(numRuns)

		minV, maxV, avgV := utils.CalculateStatistics(solutions)
		avgTimeMs := float64(avgTime.Nanoseconds()) / 1e6
		avgIterations := float64(totalIterations) / float64(numRuns)
		bestSolution := commonAlgorithms.FindBestSolution(solutions)

		rows = append(rows, utils.Row{
			Name:        cfg.name,
			AvgV:        avgV,
			MinV:        minV,
			MaxV:        maxV,
			AvgTms:      avgTimeMs,
			AvgLNSIters: avgIterations,
			BestPath:    bestSolution.Path,
			BestValue:   bestSolution.Objective,
		})

		log.Printf("Completed %s: best=%d, avg=%.2f, min=%d, max=%d, avg_time=%.2f ms, avg_iterations=%.1f",
			cfg.name, bestSolution.Objective, avgV, minV, maxV, avgTimeMs, avgIterations)
	}

	// Print detailed console output
	fmt.Println("\n--- Objective Value Statistics ---")
	fmt.Println("Configuration                           Avg (Min, Max)")
	fmt.Println("----------------------------------------------------------------")
	for _, r := range rows {
		fmt.Printf("%-38s  %.2f (%d, %d)\n", r.Name, r.AvgV, r.MinV, r.MaxV)
	}

	fmt.Println("\n--- Average Time per Run ---")
	fmt.Println("Configuration                           Time [ms]")
	fmt.Println("----------------------------------------------------------------")
	for _, r := range rows {
		fmt.Printf("%-38s  %.4f\n", r.Name, r.AvgTms)
	}

	fmt.Println("\n--- Average Iterations per Run ---")
	fmt.Println("Configuration                           Iterations")
	fmt.Println("----------------------------------------------------------------")
	for _, r := range rows {
		fmt.Printf("%-38s  %.1f\n", r.Name, r.AvgLNSIters)
	}

	fmt.Println("\n--- Best Solution Values ---")
	fmt.Println("Configuration                           Best Value")
	fmt.Println("----------------------------------------------------------------")
	for _, r := range rows {
		fmt.Printf("%-38s  %d\n", r.Name, r.BestValue)
	}
	fmt.Println()

	// Save CSV
	outputDir := filepath.Join("output", "09_labs", "hybrid_evolutionary_algorithm", "results")
	if err := utils.WriteResultsCSV(instanceName, rows, outputDir); err != nil {
		log.Printf("CSV write error for instance %s: %v", instanceName, err)
	} else {
		log.Printf("CSV results saved for instance %s", instanceName)
	}

	// Plot best solutions for each method
	plotBounds := config.DefaultPlotBounds
	plotDir := filepath.Join("output", "09_labs", "hybrid_evolutionary_algorithm", "plots")
	for i, r := range rows {
		title := fmt.Sprintf("%s - Instance %s (Value: %d)", r.Name, instanceName, r.BestValue)
		fileName := utils.SanitizeFileName(fmt.Sprintf("hybrid_%s_instance_%s_%d",
			r.Name, instanceName, i))

		if err := visualisation.PlotSolution(nodes, r.BestPath, title, fileName,
			plotBounds.XMin, plotBounds.XMax, plotBounds.YMin, plotBounds.YMax, plotDir); err != nil {
			log.Printf("Plot error for %s/%s: %v", instanceName, r.Name, err)
		} else {
			log.Printf("Plot saved: %s.png", fileName)
		}
	}
}

func main() {
	log.Println("========================================")
	log.Println("Starting Hybrid Evolutionary Algorithm Experiments")
	log.Println("========================================")
	log.Printf("Configuration: Population Size = %d, Runs per config = %d\n", populationSize, numRuns)

	instancePaths := config.DefaultInstancePaths()
	nodesA, err := data.ReadNodes(instancePaths.TSPA)
	if err != nil {
		log.Fatalf("Error reading TSPA.csv: %v", err)
	}
	nodesB, err := data.ReadNodes(instancePaths.TSPB)
	if err != nil {
		log.Fatalf("Error reading TSPB.csv: %v", err)
	}

	log.Printf("Loaded %d nodes from instance A", len(nodesA))
	log.Printf("Loaded %d nodes from instance B", len(nodesB))

	processInstance("A", nodesA)
	processInstance("B", nodesB)

	log.Println("\n========================================")
	log.Println("All experiments completed successfully")
	log.Println("========================================")
}
