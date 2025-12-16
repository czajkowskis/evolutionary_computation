package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/czajkowskis/evolutionary_computation/10_lab/variable_neighborhood_search/pkg/algorithms"
	"github.com/czajkowskis/evolutionary_computation/10_lab/variable_neighborhood_search/pkg/data"
	"github.com/czajkowskis/evolutionary_computation/10_lab/variable_neighborhood_search/pkg/utils"
	"github.com/czajkowskis/evolutionary_computation/10_lab/variable_neighborhood_search/pkg/visualisation"
)

// Configuration constants
const (
	numVNSRuns = 20      // Number of VNS runs per instance
	timeLimitA = 3276.57 // Average running time of MSLS from the previous assignment for instance A
	timeLimitB = 2342.11 // Average running time of MSLS from the previous assignment for instance B
)

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

	// Define VNS configurations to test (only with local search)
	configs := []struct {
		name               string
		neighborhoodChange string
		shakingIntensity   int
		useLocalSearch     bool
	}{
		{"VNS_Sequential_LS", "sequential", 3, true},
		{"VNS_Adaptive_LS", "adaptive", 3, true},
		{"VNS_Random_LS", "random", 3, true},
	}

	// Run VNS with different configurations
	for _, cfg := range configs {
		log.Printf("Starting VNS with config: %s for instance %s", cfg.name, instanceName)
		start := time.Now()

		var vnsResults []algorithms.VNSResult
		totalVNSIterations := 0
		for run := 0; run < numVNSRuns; run++ {
			vnsResult := algorithms.VariableNeighborhoodSearch(D, costs, algorithms.VNSConfig{
				TimeLimit:          timeLimit,
				MaxNeighborhoods:   4,
				ShakingIntensity:   cfg.shakingIntensity,
				NeighborhoodChange: cfg.neighborhoodChange,
				UseLocalSearch:     cfg.useLocalSearch,
			})
			vnsResults = append(vnsResults, vnsResult)
			totalVNSIterations += vnsResult.Iterations
		}

		totalVNSTime := time.Since(start)
		avgVNSTime := totalVNSTime / time.Duration(numVNSRuns)

		// Collect VNS solutions for statistics
		vnsSolutions := make([]algorithms.Solution, len(vnsResults))
		for i, r := range vnsResults {
			vnsSolutions[i] = r.BestSolution
		}

		vnsMin, vnsMax, vnsAvg := utils.CalculateStatistics(vnsSolutions)
		avgVNSTimeMs := float64(avgVNSTime.Nanoseconds()) / 1e6
		avgVNSIterations := float64(totalVNSIterations) / float64(numVNSRuns)
		bestVNS := algorithms.FindBestSolution(vnsSolutions)

		rows = append(rows, utils.Row{
			Name:        cfg.name,
			AvgV:        vnsAvg,
			MinV:        vnsMin,
			MaxV:        vnsMax,
			AvgTms:      avgVNSTimeMs,
			AvgVNSIters: avgVNSIterations,
			BestPath:    bestVNS.Path,
			BestValue:   bestVNS.Objective,
		})

		log.Printf("Completed %s: best value %d, avg time %.2f ms, avg iterations %.1f", cfg.name, bestVNS.Objective, avgVNSTimeMs, avgVNSIterations)
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
		fmt.Printf("%-34s  %.1f\n", r.Name, r.AvgVNSIters)
	}

	// Save CSV
	if err := utils.WriteResultsCSV(instanceName, rows); err != nil {
		log.Printf("CSV write error for instance %s: %v", instanceName, err)
	} else {
		log.Printf("CSV results saved for instance %s", instanceName)
	}

	// Plot best solutions for each method
	for i, r := range rows {
		title := fmt.Sprintf("%s for Instance %s (Value: %d)", r.Name, instanceName, r.BestValue)
		fileName := utils.SanitizeFileName(fmt.Sprintf("%s_Instance_%s_%d", r.Name, instanceName, i))
		if err := visualisation.PlotSolution(nodes, r.BestPath, title, fileName, 0, 4000, 0, 2000); err != nil {
			log.Printf("plot error for %s/%s: %v", instanceName, r.Name, err)
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	log.Println("Starting Variable Neighborhood Search experiments")

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

	log.Println("Program execution completed")
}
