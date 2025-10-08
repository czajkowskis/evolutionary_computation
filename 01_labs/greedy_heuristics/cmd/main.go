package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/czajkowskis/evolutionary_computation/01_labs/greedy_heuristics/pkg/algorithms"
	"github.com/czajkowskis/evolutionary_computation/01_labs/greedy_heuristics/pkg/data"
	"github.com/czajkowskis/evolutionary_computation/01_labs/greedy_heuristics/pkg/utils"
	"github.com/czajkowskis/evolutionary_computation/01_labs/greedy_heuristics/pkg/visualisation"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	// Read nodes from CSV files
	nodesA, err := data.ReadNodes("../../TSPA.csv")
	if err != nil {
		log.Fatalf("Error reading TSPA.csv: %v", err)
	}

	nodesB, err := data.ReadNodes("../../TSPB.csv")
	if err != nil {
		log.Fatalf("Error reading TSPB.csv: %v", err)
	}

	// Calculate distance matrices
	distanceMatrixA := data.CalculateDistanceMatrix(nodesA)
	distanceMatrixB := data.CalculateDistanceMatrix(nodesB)

	// Generate solutions
	numSolutions := 200
	randomSolutionsA := algorithms.RandomSolution(nodesA, distanceMatrixA, numSolutions)
	nnEndSolutionsA := algorithms.NearestNeighborEnd(nodesA, distanceMatrixA, numSolutions)
	nnAnySolutionsA := algorithms.NearestNeighborAny(nodesA, distanceMatrixA, numSolutions)
	// greedyCycleSolutionsA := algorithms.GreedyCycle(nodesA, distanceMatrixA, numSolutions)

	randomSolutionsB := algorithms.RandomSolution(nodesB, distanceMatrixB, numSolutions)
	nnEndSolutionsB := algorithms.NearestNeighborEnd(nodesB, distanceMatrixB, numSolutions)
	nnAnySolutionsB := algorithms.NearestNeighborAny(nodesB, distanceMatrixB, numSolutions)
	// greedyCycleSolutionsB := algorithms.GreedyCycle(nodesB, distanceMatrixB, numSolutions)

	// Calculate statistics for instance A
	minRandomA, maxRandomA, avgRandomA := utils.CalculateStatistics(randomSolutionsA)
	minNnEndA, maxNnEndA, avgNnEndA := utils.CalculateStatistics(nnEndSolutionsA)
	minNnAnyA, maxNnAnyA, avgNnAnyA := utils.CalculateStatistics(nnAnySolutionsA)
	// minGreedyA, maxGreedyA, avgGreedyA := utils.CalculateStatistics(greedyCycleSolutionsA)

	// Calculate statistics for instance B
	minRandomB, maxRandomB, avgRandomB := utils.CalculateStatistics(randomSolutionsB)
	minNnEndB, maxNnEndB, avgNnEndB := utils.CalculateStatistics(nnEndSolutionsB)
	minNnAnyB, maxNnAnyB, avgNnAnyB := utils.CalculateStatistics(nnAnySolutionsB)
	// minGreedyB, maxGreedyB, avgGreedyB := utils.CalculateStatistics(greedyCycleSolutionsB)

	// Print statistics for instance A
	fmt.Println("Instance A Statistics:")
	fmt.Printf("Random Solution: min = %d, max = %d, average = %.2f\n", minRandomA, maxRandomA, avgRandomA)
	fmt.Printf("Nearest Neighbor (End Only): min = %d, max = %d, average = %.2f\n", minNnEndA, maxNnEndA, avgNnEndA)
	fmt.Printf("Nearest Neighbor (Any Position): min = %d, max = %d, average = %.2f\n", minNnAnyA, maxNnAnyA, avgNnAnyA)
	// fmt.Printf("Greedy Cycle: min = %d, max = %d, average = %.2f\n", minGreedyA, maxGreedyA, avgGreedyA)

	// Print statistics for instance B
	fmt.Println("\nInstance B Statistics:")
	fmt.Printf("Random Solution: min = %d, max = %d, average = %.2f\n", minRandomB, maxRandomB, avgRandomB)
	fmt.Printf("Nearest Neighbor (End Only): min = %d, max = %d, average = %.2f\n", minNnEndB, maxNnEndB, avgNnEndB)
	fmt.Printf("Nearest Neighbor (Any Position): min = %d, max = %d, average = %.2f\n", minNnAnyB, maxNnAnyB, avgNnAnyB)
	// fmt.Printf("Greedy Cycle: min = %d, max = %d, average = %.2f\n", minGreedyB, maxGreedyB, avgGreedyB)

	// Find and plot the best solutions for instance A
	bestRandomA := algorithms.FindBestSolution(randomSolutionsA)
	if err := visualisation.PlotSolution(nodesA, bestRandomA.Path, "Best_Random_Solution_A", 0, 4000, 0, 2000); err != nil {
		log.Printf("Error plotting best random solution for A: %v", err)
	}

	bestNnEndA := algorithms.FindBestSolution(nnEndSolutionsA)
	if err := visualisation.PlotSolution(nodesA, bestNnEndA.Path, "Best_Nearest_Neighbor_End_Solution_A", 0, 4000, 0, 2000); err != nil {
		log.Printf("Error plotting best nearest neighbor end solution for A: %v", err)
	}

	bestNnAnyA := algorithms.FindBestSolution(nnAnySolutionsA)
	if err := visualisation.PlotSolution(nodesA, bestNnAnyA.Path, "Best_Nearest_Neighbor_Any_Solution_A", 0, 4000, 0, 2000); err != nil {
		log.Printf("Error plotting best nearest neighbor any solution for A: %v", err)
	}

	// if greedyCycleSolutionsA != nil {
	// 	bestGreedyA := algorithms.FindBestSolution(greedyCycleSolutionsA)
	// 	if err := visualisation.PlotSolution(nodesA, bestGreedyA.Path, "Best_Greedy_Cycle_Solution_A"); err != nil {
	// 		log.Printf("Error plotting best greedy cycle solution for A: %v", err)
	// 	}
	// }

	// Find and plot the best solutions for instance B
	bestRandomB := algorithms.FindBestSolution(randomSolutionsB)
	if err := visualisation.PlotSolution(nodesB, bestRandomB.Path, "Best_Random_Solution_B", 0, 4000, 0, 2000); err != nil {
		log.Printf("Error plotting best random solution for B: %v", err)
	}

	bestNnEndB := algorithms.FindBestSolution(nnEndSolutionsB)
	if err := visualisation.PlotSolution(nodesB, bestNnEndB.Path, "Best_Nearest_Neighbor_End_Solution_B", 0, 4000, 0, 2000); err != nil {
		log.Printf("Error plotting best nearest neighbor end solution for B: %v", err)
	}

	bestNnAnyB := algorithms.FindBestSolution(nnAnySolutionsB)
	if err := visualisation.PlotSolution(nodesB, bestNnAnyB.Path, "Best_Nearest_Neighbor_Any_Solution_B", 0, 4000, 0, 2000); err != nil {
		log.Printf("Error plotting best nearest neighbor any solution for B: %v", err)
	}

	// if greedyCycleSolutionsB != nil {
	// 	bestGreedyB := algorithms.FindBestSolution(greedyCycleSolutionsB)
	// 	if err := visualisation.PlotSolution(nodesB, bestGreedyB.Path, "Best_Greedy_Cycle_Solution_B"); err != nil {
	// 		log.Printf("Error plotting best greedy cycle solution for B: %v", err)
	// 	}
	// }
}
