package main

import (
	"fmt"
	"log"
	"math/rand"
	"path/filepath"
	"time"

	"github.com/czajkowskis/evolutionary_computation/06_labs/local_search_extensions/pkg/algorithms"
	commonAlgorithms "github.com/czajkowskis/evolutionary_computation/pkg/common/algorithms"
	"github.com/czajkowskis/evolutionary_computation/pkg/common/config"
	"github.com/czajkowskis/evolutionary_computation/pkg/common/data"
	"github.com/czajkowskis/evolutionary_computation/pkg/common/utils"
	"github.com/czajkowskis/evolutionary_computation/pkg/common/visualisation"
)

// Configuration constants
const (
	numMSLSRuns   = 20  // Number of MSLS runs per instance
	numMSLSStarts = 200 // Number of local search starts per MSLS run
	numILSRuns    = 20  // Number of ILS runs per instance
)

// InstanceResults stores all results for a single instance
type InstanceResults struct {
	Instance    string
	MSLSResults []algorithms.MSLSResult
	ILSResults  []algorithms.ILSResult
	AvgMSLSTime time.Duration
}

// processInstance runs the full experimental pipeline for a single instance
func processInstance(instanceName string, nodes []data.Node) {
	log.Printf("Processing instance %s with %d nodes", instanceName, len(nodes))
	fmt.Printf("Instance %s Statistics:\n", instanceName)

	D := data.CalculateDistanceMatrix(nodes)
	costs := make([]int, len(nodes))
	for i, node := range nodes {
		costs[i] = node.Cost
	}

	numMSLSRuns := 20
	numMSLSStarts := 200
	numILSRuns := 20

	var rows []utils.Row

	// === PHASE 1: Run MSLS ===
	log.Printf("Starting MSLS for instance %s", instanceName)
	start := time.Now()

	var mslsResults []algorithms.MSLSResult
	for run := 0; run < numMSLSRuns; run++ {
		mslsResult := algorithms.RunMSLS(D, costs, numMSLSStarts)
		mslsResults = append(mslsResults, mslsResult)
	}

	totalMSLSTime := time.Since(start)
	avgMSLSTime := totalMSLSTime / time.Duration(numMSLSRuns)

	// Collect MSLS solutions for statistics
	mslsSolutions := make([]commonAlgorithms.Solution, len(mslsResults))
	for i, r := range mslsResults {
		mslsSolutions[i] = r.BestSolution
	}

	mslsMin, mslsMax, mslsAvg := utils.CalculateStatistics(mslsSolutions)
	avgMSLSTimeMs := float64(avgMSLSTime.Nanoseconds()) / 1e6
	bestMSLS := commonAlgorithms.FindBestSolution(mslsSolutions)

	rows = append(rows, utils.Row{
		Name:      "MSLS",
		AvgV:      mslsAvg,
		MinV:      mslsMin,
		MaxV:      mslsMax,
		AvgTms:    avgMSLSTimeMs,
		BestPath:  bestMSLS.Path,
		BestValue: bestMSLS.Objective,
	})

	log.Printf("Completed MSLS: best value %d, avg time %.2f ms", bestMSLS.Objective, avgMSLSTimeMs)

	// === PHASE 2: Run ILS ===
	log.Printf("Starting ILS for instance %s with time limit %v", instanceName, avgMSLSTime)
	perturbType := algorithms.PerturbRandom4Opt

	var ilsResults []algorithms.ILSResult
	totalILSIterations := 0
	for run := 0; run < numILSRuns; run++ {
		ilsResult := algorithms.RunILS(D, costs, avgMSLSTime, perturbType)
		ilsResults = append(ilsResults, ilsResult)
		totalILSIterations += ilsResult.NumLSIterations
	}

	// Collect ILS solutions for statistics
	ilsSolutions := make([]commonAlgorithms.Solution, len(ilsResults))
	for i, r := range ilsResults {
		ilsSolutions[i] = r.BestSolution
	}

	ilsMin, ilsMax, ilsAvg := utils.CalculateStatistics(ilsSolutions)
	avgILSTimeMs := float64(avgMSLSTime.Nanoseconds()) / 1e6 // Same time limit as MSLS
	avgLSIterations := float64(totalILSIterations) / float64(numILSRuns)
	bestILS := commonAlgorithms.FindBestSolution(ilsSolutions)

	rows = append(rows, utils.Row{
		Name:      "ILS",
		AvgV:      ilsAvg,
		MinV:      ilsMin,
		MaxV:      ilsMax,
		AvgTms:    avgILSTimeMs,
		BestPath:  bestILS.Path,
		BestValue: bestILS.Objective,
	})

	log.Printf("Completed ILS: best value %d, avg LS iterations %.1f", bestILS.Objective, avgLSIterations)

	// Print console output
	fmt.Println("Objective value: av (min, max)")
	for _, r := range rows {
		fmt.Printf("%-34s  %.2f (%d, %d)\n", r.Name, r.AvgV, r.MinV, r.MaxV)
	}
	fmt.Println()

	fmt.Println("Average time per run [ms]:")
	for _, r := range rows {
		fmt.Printf("%-34s  %.4f\n", r.Name, r.AvgTms)
	}
	fmt.Printf("ILS - Average LS iterations per run: %.1f\n", avgLSIterations)

	// Save CSV
	outputDir := filepath.Join("output", "06_labs", "local_search_extensions", "results")
	if err := utils.WriteResultsCSV(instanceName, rows, outputDir); err != nil {
		log.Printf("CSV write error for instance %s: %v", instanceName, err)
	} else {
		log.Printf("CSV results saved for instance %s", instanceName)
	}

	// Plot best solutions
	plotBounds := config.DefaultPlotBounds
	plotDir := filepath.Join("output", "06_labs", "local_search_extensions", "plots")
	titleMSLS := fmt.Sprintf("Best MSLS Solution for Instance %s", instanceName)
	fileNameMSLS := utils.SanitizeFileName(fmt.Sprintf("Best_MSLS_Solution_%s", instanceName))
	if err := visualisation.PlotSolution(nodes, bestMSLS.Path, titleMSLS, fileNameMSLS,
		plotBounds.XMin, plotBounds.XMax, plotBounds.YMin, plotBounds.YMax, plotDir); err != nil {
		log.Printf("plot error for %s/MSLS: %v", instanceName, err)
	}

	titleILS := fmt.Sprintf("Best ILS Solution for Instance %s", instanceName)
	fileNameILS := utils.SanitizeFileName(fmt.Sprintf("Best_ILS_Solution_%s", instanceName))
	if err := visualisation.PlotSolution(nodes, bestILS.Path, titleILS, fileNameILS,
		plotBounds.XMin, plotBounds.XMax, plotBounds.YMin, plotBounds.YMax, plotDir); err != nil {
		log.Printf("plot error for %s/ILS: %v", instanceName, err)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	log.Println("Starting MSLS vs ILS local search experiments")

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
