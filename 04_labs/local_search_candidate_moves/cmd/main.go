package main

import (
	"fmt"
	"log"
	"math/rand"
	"path/filepath"
	"time"

	"github.com/czajkowskis/evolutionary_computation/04_labs/local_search_candidate_moves/pkg/algorithms"
	commonAlgorithms "github.com/czajkowskis/evolutionary_computation/pkg/common/algorithms"
	"github.com/czajkowskis/evolutionary_computation/pkg/common/config"
	"github.com/czajkowskis/evolutionary_computation/pkg/common/data"
	"github.com/czajkowskis/evolutionary_computation/pkg/common/utils"
	"github.com/czajkowskis/evolutionary_computation/pkg/common/visualisation"
)

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
		},
		// CANDIDATE MOVES - K = 5
		{
			Name:    "Candidates_Steepest_2opt_Random_K5",
			UseCand: true,
			CandK:   5,
		},
		// CANDIDATE MOVES - K = 10 (default value from the task)
		{
			Name:    "Candidates_Steepest_2opt_Random_K10",
			UseCand: true,
			CandK:   10,
		},
		// CANDIDATE MOVES - K = 15
		{
			Name:    "Candidates_Steepest_2opt_Random_K15",
			UseCand: true,
			CandK:   15,
		},
	}

	var rows []utils.Row

	for _, m := range methods {
		log.Printf("Starting method: %s for instance %s", m.Name, instanceName)
		start := time.Now()

		solutions := algorithms.RunLocalSearchBatch(D, costs, m, numSolutions)
		batchTime := time.Since(start)

		if len(solutions) == 0 {
			log.Printf("No solutions found for method %s", m.Name)
			continue
		}

		minVal, maxVal, avgVal := utils.CalculateStatistics(solutions)
		avgTimeMs := float64(batchTime.Nanoseconds()) / float64(numSolutions) / 1e6

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

		log.Printf("Completed method %s: best value %d, avg time %.2f ms", m.Name, best.Objective, avgTimeMs)

		// Wykres najlepszej trasy
		title := fmt.Sprintf("Best %s Solution for Instance %s", m.Name, instanceName)
		fileName := utils.SanitizeFileName(fmt.Sprintf("Best_%s_Solution_%s", m.Name, instanceName))
		plotBounds := config.DefaultPlotBounds
		outputDir := filepath.Join("output", "04_labs", "local_search_candidate_moves", "plots")
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
	outputDir := filepath.Join("output", "04_labs", "local_search_candidate_moves", "results")
	if err := utils.WriteResultsCSV(instanceName, rows, outputDir); err != nil {
		log.Printf("CSV write error for instance %s: %v", instanceName, err)
	} else {
		log.Printf("CSV results saved for instance %s", instanceName)
	}
}

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
