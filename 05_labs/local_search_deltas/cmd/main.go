// Command local_search_deltas runs TSP local search experiments for this lab.
package main

import (
	"fmt"
	"log"
	"math/rand"
	"path/filepath"
	"time"

	"github.com/czajkowskis/evolutionary_computation/05_labs/local_search_deltas/pkg/algorithms"
	commonAlgorithms "github.com/czajkowskis/evolutionary_computation/pkg/common/algorithms"
	"github.com/czajkowskis/evolutionary_computation/pkg/common/config"
	"github.com/czajkowskis/evolutionary_computation/pkg/common/data"
	"github.com/czajkowskis/evolutionary_computation/pkg/common/utils"
	"github.com/czajkowskis/evolutionary_computation/pkg/common/visualisation"
)

// processInstance runs the full experimental pipeline for a single instance:
// build distance matrix, run all configured methods, print stats, plot best
// solutions and persist CSV summaries.
func processInstance(instanceName string, nodes []data.Node) {
	log.Printf("Processing instance %s with %d nodes", instanceName, len(nodes))
	fmt.Printf("Instance %s Statistics:\n", instanceName)

	D := data.CalculateDistanceMatrix(nodes)

	costs := make([]int, len(nodes))
	for i, node := range nodes {
		costs[i] = node.Cost
	}

	numSolutions := 200

	methods := []algorithms.MethodSpec{
		// BASELINE (no candidate moves)
		{
			Name:    "Baseline_Steepest_2opt_Random",
			UseCand: false,
			CandK:   0,
			UseLM:   false,
		},
		// LM-based steepest local search (no candidate moves)
		{
			Name:    "LM_Steepest_2opt_Random",
			UseCand: false,
			CandK:   0,
			UseLM:   true,
		},
		// CANDIDATE MOVES - K = 5
		{
			Name:    "Candidates_Steepest_2opt_Random_K5",
			UseCand: true,
			CandK:   5,
			UseLM:   false,
		},
		// CANDIDATE MOVES - K = 10 (default value from the task)
		{
			Name:    "Candidates_Steepest_2opt_Random_K10",
			UseCand: true,
			CandK:   10,
			UseLM:   false,
		},
		// CANDIDATE MOVES - K = 15
		{
			Name:    "Candidates_Steepest_2opt_Random_K15",
			UseCand: true,
			CandK:   15,
			UseLM:   false,
		},
	}

	var rows []utils.Row

	for _, m := range methods {
		log.Printf("Starting method: %s for instance %s", m.Name, instanceName)
		start := time.Now()

		solutions, durations := algorithms.RunLocalSearchBatch(D, costs, m, numSolutions)
		batchTime := time.Since(start)

		if len(solutions) == 0 {
			log.Printf("No solutions found for method %s", m.Name)
			continue
		}

		minVal, maxVal, avgVal := utils.CalculateStatistics(solutions)
		// Per-run timing statistics (always compute average; for LM also print min/max).
		avgTimeMs := float64(batchTime.Nanoseconds()) / float64(numSolutions) / 1e6
		var minTimeMs, maxTimeMs float64
		if len(durations) > 0 {
			minDur := durations[0]
			maxDur := durations[0]
			var total time.Duration
			for _, d := range durations {
				if d < minDur {
					minDur = d
				}
				if d > maxDur {
					maxDur = d
				}
				total += d
			}
			avgTimeMs = float64(total.Nanoseconds()) / float64(len(durations)) / 1e6
			minTimeMs = float64(minDur.Nanoseconds()) / 1e6
			maxTimeMs = float64(maxDur.Nanoseconds()) / 1e6
		}

		best := commonAlgorithms.FindBestSolution(solutions)

		rows = append(rows, utils.Row{
			Name:      m.Name,
			AvgV:      avgVal,
			MinV:      minVal,
			MaxV:      maxVal,
			AvgTms:    avgTimeMs,
			BestPath:  best.Path,
			BestValue: best.Objective,
		})

		if m.UseLM && len(durations) > 0 {
			log.Printf("Completed method %s: best value %d, avg time %.2f ms (min: %.2f, max: %.2f)",
				m.Name, best.Objective, avgTimeMs, minTimeMs, maxTimeMs)
		} else {
			log.Printf("Completed method %s: best value %d, avg time %.2f ms",
				m.Name, best.Objective, avgTimeMs)
		}

		// Wykres najlepszej trasy
		title := fmt.Sprintf("Best %s Solution for Instance %s", m.Name, instanceName)
		fileName := utils.SanitizeFileName(fmt.Sprintf("Best_%s_Solution_%s", m.Name, instanceName))
		plotBounds := config.DefaultPlotBounds
		outputDir := filepath.Join("output", "05_labs", "local_search_deltas", "plots")
		if err := visualisation.PlotSolution(nodes, best.Path, title, fileName,
			plotBounds.XMin, plotBounds.XMax, plotBounds.YMin, plotBounds.YMax, outputDir); err != nil {
			log.Printf("plot error for %s/%s: %v", instanceName, m.Name, err)
		}
	}

	// 5) Wyniki — konsola
	fmt.Println("Objective value: av (min, max)")
	for _, r := range rows {
		fmt.Printf("%-34s  %.2f (%d, %d)\n", r.Name, r.AvgV, r.MinV, r.MaxV)
		fmt.Printf("Best path: %v\n", r.BestPath)
	}
	fmt.Println()

	fmt.Println("Average time per run [ms]:")
	for _, r := range rows {
		fmt.Printf("%-34s  %.4f\n", r.Name, r.AvgTms)
	}

	// 6) CSV
	outputDir := filepath.Join("output", "05_labs", "local_search_deltas", "results")
	if err := utils.WriteResultsCSV(instanceName, rows, outputDir); err != nil {
		log.Printf("CSV write error for instance %s: %v", instanceName, err)
	} else {
		log.Printf("CSV results saved for instance %s", instanceName)
	}
}

// main seeds RNG and runs the local search experiments for both provided
// instances A and B.
func main() {
	rand.Seed(time.Now().UnixNano())
	log.Println("Starting evolutionary computation local search program")

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
