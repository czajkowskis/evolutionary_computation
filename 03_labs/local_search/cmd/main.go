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
	fmt.Printf("Instance %s Statistics:\n", instanceName)

	// Odległości (po wczytaniu pracujemy wyłącznie na macierzy D i wektorze kosztów)
	D := data.CalculateDistanceMatrix(nodes)
	startNodeIndices := utils.GenerateStartNodeIndices(len(nodes))
	numSolutions := 200 // liczba uruchomień na metodę

	costs := make([]int, len(nodes))
	for i, node := range nodes {
		costs[i] = node.Cost
	}

	methods := []algorithms.MethodSpec{
		{LS: algorithms.LS_Steepest, Intra: algorithms.IntraSwap, Start: algorithms.StartRandom, Name: "Steepest + Swap + Random"},
		// {LS: algorithms.LS_Steepest, Intra: algorithms.IntraSwap, Start: algorithms.StartGreedy, Name: "Steepest + Swap + GreedyStart"},
		// {LS: algorithms.LS_Steepest, Intra: algorithms.Intra2Opt, Start: algorithms.StartRandom, Name: "Steepest + 2-opt + Random"},
		// {LS: algorithms.LS_Steepest, Intra: algorithms.Intra2Opt, Start: algorithms.StartGreedy, Name: "Steepest + 2-opt + GreedyStart"},
		// {LS: algorithms.LS_Greedy, Intra: algorithms.IntraSwap, Start: algorithms.StartRandom, Name: "Greedy + Swap + Random"},
		// {LS: algorithms.LS_Greedy, Intra: algorithms.IntraSwap, Start: algorithms.StartGreedy, Name: "Greedy + Swap + GreedyStart"},
		// {LS: algorithms.LS_Greedy, Intra: algorithms.Intra2Opt, Start: algorithms.StartRandom, Name: "Greedy + 2-opt + Random"},
		// {LS: algorithms.LS_Greedy, Intra: algorithms.Intra2Opt, Start: algorithms.StartGreedy, Name: "Greedy + 2-opt + GreedyStart"},
	}

	// Zbierz wyniki do tabel
	type Row struct {
		Name      string
		AvgV      float64
		MinV      int
		MaxV      int
		AvgTms    float64
		BestPath  []int
		BestValue int
	}
	var rows []Row

	for _, m := range methods {
		start := time.Now()
		solutions := algorithms.RunLocalSearchBatch(D, costs, startNodeIndices, m, numSolutions)
		batchTime := time.Since(start)

		if len(solutions) == 0 {
			continue
		}
		minVal, maxVal, avgVal := utils.CalculateStatistics(solutions)
		avgTimeMs := float64(batchTime.Nanoseconds()) / float64(numSolutions) / 1e6

		best := algorithms.FindBestSolution(solutions)

		rows = append(rows, Row{
			Name:      m.Name,
			AvgV:      avgVal,
			MinV:      minVal,
			MaxV:      maxVal,
			AvgTms:    avgTimeMs,
			BestPath:  best.Path,
			BestValue: best.Objective,
		})

		// Wykres najlepszej ścieżki (opcjonalny; granice dopasuj do swoich instancji)
		title := fmt.Sprintf("Best %s Solution for Instance %s", m.Name, instanceName)
		fileName := utils.SanitizeFileName(fmt.Sprintf("Best_%s_Solution_%s", m.Name, instanceName))
		if err := visualisation.PlotSolution(nodes, best.Path, title, fileName, 0, 4000, 0, 2000); err != nil {
			log.Printf("plot error for %s/%s: %v", instanceName, m.Name, err)
		}
	}

	// Tabela wartości celu
	fmt.Println("Objective value: av (min, max)")
	for _, r := range rows {
		fmt.Printf("%-34s  %.2f (%d, %d)\n", r.Name, r.AvgV, r.MinV, r.MaxV)
	}
	fmt.Println()

	// Tabela czasów (średni czas jednego uruchomienia)
	fmt.Println("Average time per run [ms]:")
	for _, r := range rows {
		fmt.Printf("%-34s  %.4f\n", r.Name, r.AvgTms)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// Wczytanie dwóch instancji (ścieżki jak w Twoim projekcie)
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
