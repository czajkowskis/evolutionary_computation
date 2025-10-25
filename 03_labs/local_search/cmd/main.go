package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/czajkowskis/evolutionary_computation/03_labs/local_search/pkg/algorithms"
	"github.com/czajkowskis/evolutionary_computation/03_labs/local_search/pkg/data"
	"github.com/czajkowskis/evolutionary_computation/03_labs/local_search/pkg/utils"
	"github.com/czajkowskis/evolutionary_computation/03_labs/local_search/pkg/visualisation"
)

func processInstance(instanceName string, nodes []data.Node) {
	log.Printf("Processing instance %s with %d nodes", instanceName, len(nodes))

	fmt.Printf("Instance %s Statistics:\n", instanceName)

	D := data.CalculateDistanceMatrix(nodes)
	startNodeIndices := utils.GenerateStartNodeIndices(len(nodes))
	numSolutions := 200

	costs := make([]int, len(nodes))
	for i, node := range nodes {
		costs[i] = node.Cost
	}

	methods := []algorithms.MethodSpec{
		{LS: algorithms.LS_Steepest, Intra: algorithms.IntraSwap, Start: algorithms.StartRandom, Name: "Steepest_Swap_Random"},
		{LS: algorithms.LS_Steepest, Intra: algorithms.IntraSwap, Start: algorithms.StartGreedy, Name: "Steepest_Swap_GreedyStart"},
		{LS: algorithms.LS_Steepest, Intra: algorithms.Intra2Opt, Start: algorithms.StartRandom, Name: "Steepest_2-opt_Random"},
		{LS: algorithms.LS_Steepest, Intra: algorithms.Intra2Opt, Start: algorithms.StartGreedy, Name: "Steepest_2-opt_GreedyStart"},
		{LS: algorithms.LS_Greedy, Intra: algorithms.IntraSwap, Start: algorithms.StartRandom, Name: "Greedy_Swap_Random"},
		{LS: algorithms.LS_Greedy, Intra: algorithms.IntraSwap, Start: algorithms.StartGreedy, Name: "Greedy_Swap_GreedyStart"},
		{LS: algorithms.LS_Greedy, Intra: algorithms.Intra2Opt, Start: algorithms.StartRandom, Name: "Greedy_2-opt_Random"},
		{LS: algorithms.LS_Greedy, Intra: algorithms.Intra2Opt, Start: algorithms.StartGreedy, Name: "Greedy_2-opt_GreedyStart"},
	}

	var rows []utils.Row

	for _, m := range methods {
		log.Printf("Starting method: %s for instance %s", m.Name, instanceName)
		start := time.Now()

		solutions := algorithms.RunLocalSearchBatch(D, costs, startNodeIndices, m, numSolutions)
		batchTime := time.Since(start)

		if len(solutions) == 0 {
			log.Printf("No solutions found for method %s", m.Name)
			continue
		}
		minVal, maxVal, avgVal := utils.CalculateStatistics(solutions)
		avgTimeMs := float64(batchTime.Nanoseconds()) / float64(numSolutions) / 1e6

		best := algorithms.FindBestSolution(solutions)

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

		// Plots for the best solutions
		title := fmt.Sprintf("Best %s Solution for Instance %s", m.Name, instanceName)
		fileName := utils.SanitizeFileName(fmt.Sprintf("Best_%s_Solution_%s", m.Name, instanceName))
		if err := visualisation.PlotSolution(nodes, best.Path, title, fileName, 0, 4000, 0, 2000); err != nil {
			log.Printf("plot error for %s/%s: %v", instanceName, m.Name, err)
		}
	}

	// Print results to console
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

	// Save results to CSV
	if err := utils.WriteResultsCSV(instanceName, rows); err != nil {
		log.Printf("CSV write error for instance %s: %v", instanceName, err)
	} else {
		log.Printf("CSV results saved for instance %s", instanceName)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	log.Println("Starting evolutionary computation local search program")

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
