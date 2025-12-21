package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"time"

	"github.com/czajkowskis/evolutionary_computation/03_labs/local_search/pkg/algorithms"
	commonAlgorithms "github.com/czajkowskis/evolutionary_computation/pkg/common/algorithms"
	"github.com/czajkowskis/evolutionary_computation/pkg/common/config"
	"github.com/czajkowskis/evolutionary_computation/pkg/common/data"
	"github.com/czajkowskis/evolutionary_computation/pkg/common/utils"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	// Load instances
	instancePaths := config.DefaultInstancePaths()
	instances := []struct {
		Name string
		Path string
	}{
		{"A", instancePaths.TSPA},
		{"B", instancePaths.TSPB},
	}

	for _, inst := range instances {
		log.Printf("Processing instance %s...", inst.Name)
		nodes, err := data.ReadNodes(inst.Path)
		if err != nil {
			log.Fatalf("Error reading %s: %v", inst.Path, err)
		}

		D := data.CalculateDistanceMatrix(nodes)
		costs := make([]int, len(nodes))
		for i, node := range nodes {
			costs[i] = node.Cost
		}

		// 1. Generate 1000 random local optima
		log.Println("Generating 1000 random local optima...")
		numSolutions := 1000

		// Method for generating random local optima: Random Start + Greedy LS + Swap (as per assignment "random local optima obtained from random solutions using greedy local search")
		// "greedy local search" usually implies Steepest or Greedy. The assignment says "greedy local search".
		// In the provided code `LS_Greedy` exists. I will use `LS_Greedy`.
		// Intra operator is not specified, I'll use IntraSwap as it's standard.
		genMethod := algorithms.MethodSpec{
			LS:    algorithms.LS_Greedy,
			Intra: algorithms.IntraSwap,
			Start: algorithms.StartRandom,
			Name:  "Generator",
		}

		// We don't need startNodeIndices for StartRandom
		solutions := algorithms.RunLocalSearchBatch(D, costs, nil, genMethod, numSolutions)

		// 2. Identify "Best of 1000"
		bestOf1000 := solutions[0]
		for _, s := range solutions {
			if s.Objective < bestOf1000.Objective {
				bestOf1000 = s
			}
		}
		log.Printf("Best of 1000 Objective: %d", bestOf1000.Objective)

		// 3. Generate "Best Known" (Very good solution)
		log.Println("Using provided Best Known solution...")
		var bestKnown algorithms.Solution

		if inst.Name == "A" {
			path := []int{127, 123, 162, 133, 151, 51, 118, 59, 65, 116, 43, 42, 184, 35, 84, 112, 4, 190, 10, 177, 54, 48, 160, 34, 181, 146, 22, 18, 108, 69, 159, 193, 41, 139, 115, 46, 68, 140, 93, 117, 0, 143, 183, 89, 186, 23, 137, 176, 80, 79, 63, 94, 124, 148, 9, 62, 102, 144, 14, 49, 178, 106, 52, 55, 185, 40, 119, 165, 90, 81, 196, 179, 57, 129, 92, 145, 78, 31, 56, 113, 175, 171, 16, 25, 44, 120, 2, 152, 97, 1, 101, 75, 86, 26, 100, 53, 180, 154, 135, 70}
			bestKnown = algorithms.Solution{
				Path:      path,
				Objective: calculateObjective(D, costs, path),
			}
		} else if inst.Name == "B" {
			path := []int{29, 0, 109, 35, 143, 106, 124, 62, 18, 55, 34, 170, 152, 183, 140, 4, 149, 28, 20, 60, 148, 47, 94, 66, 179, 22, 99, 130, 95, 185, 86, 166, 194, 176, 113, 114, 137, 127, 89, 103, 163, 187, 153, 81, 77, 141, 91, 61, 36, 177, 5, 45, 142, 78, 175, 80, 190, 136, 73, 54, 31, 193, 117, 198, 156, 1, 16, 27, 38, 135, 63, 40, 107, 133, 122, 131, 121, 51, 90, 147, 6, 188, 169, 132, 70, 3, 15, 145, 13, 195, 168, 139, 11, 138, 33, 160, 144, 104, 8, 111}
			bestKnown = algorithms.Solution{
				Path:      path,
				Objective: calculateObjective(D, costs, path),
			}
		} else {
			log.Println("Generating Best Known solution with Strong LS...")
			strongMethod := algorithms.MethodSpec{
				LS:    algorithms.LS_Steepest,
				Intra: algorithms.Intra2Opt,
				Start: algorithms.StartGreedy,
				Name:  "Strong",
			}
			// Run a small batch to get a good one
			strongSolutions := algorithms.RunLocalSearchBatch(D, costs, utils.GenerateStartNodeIndices(len(nodes)), strongMethod, 200)
			bestKnown = commonAlgorithms.FindBestSolution(strongSolutions)
		}
		log.Printf("Best Known Objective: %d", bestKnown.Objective)

		// 4. Calculate Similarities and Correlations
		// We need 6 combinations per instance:
		// Target: Average, BestOf1000, BestKnown
		// Measure: Edges, Nodes

		var results []ResultRow

		// Pre-calculate edges/nodes sets for all solutions to speed up "Average" calculation?
		// O(N^2) comparisons might be slow for 1000x1000. 1000x1000 = 1M comparisons.
		// Each comparison takes O(K). K ~ 100-200. 1M * 100 operations is fine (100M ops ~ 1-2 seconds).

		log.Println("Calculating similarities...")

		// Helper to get edges map
		getEdges := func(path []int) map[[2]int]bool {
			edges := make(map[[2]int]bool)
			n := len(path)
			for i := 0; i < n; i++ {
				u, v := path[i], path[(i+1)%n]
				if u > v {
					u, v = v, u
				}
				edges[[2]int{u, v}] = true
			}
			return edges
		}

		// Helper to get nodes map
		getNodes := func(path []int) map[int]bool {
			nodes := make(map[int]bool)
			for _, u := range path {
				nodes[u] = true
			}
			return nodes
		}

		// Precompute maps for all solutions
		allEdges := make([]map[[2]int]bool, len(solutions))
		allNodes := make([]map[int]bool, len(solutions))
		for i, s := range solutions {
			allEdges[i] = getEdges(s.Path)
			allNodes[i] = getNodes(s.Path)
		}

		best1000Edges := getEdges(bestOf1000.Path)
		best1000Nodes := getNodes(bestOf1000.Path)
		bestKnownEdges := getEdges(bestKnown.Path)
		bestKnownNodes := getNodes(bestKnown.Path)

		for i, s := range solutions {
			// Skip the bestOf1000 itself when calculating similarity to it?
			// Assignment: "In the results with similarity to a single good solution do not include this solution itself"
			// We will filter later or handle it here.
			// Actually, if s == bestOf1000, we should probably exclude it from that specific chart.
			// But for the "Average" chart, it's just one of 1000.
			// Let's collect all data and filter during correlation/plotting.

			// Sim to Best1000
			sBest1000E := commonEdges(allEdges[i], best1000Edges)
			sBest1000N := commonNodes(allNodes[i], best1000Nodes)

			// Sim to BestKnown
			sBestKnownE := commonEdges(allEdges[i], bestKnownEdges)
			sBestKnownN := commonNodes(allNodes[i], bestKnownNodes)

			// Avg Sim
			sumEdges := 0
			sumNodes := 0
			count := 0
			for j := range solutions {
				if i == j {
					continue
				}
				sumEdges += commonEdges(allEdges[i], allEdges[j])
				sumNodes += commonNodes(allNodes[i], allNodes[j])
				count++
			}
			avgEdges := float64(sumEdges) / float64(count)
			avgNodes := float64(sumNodes) / float64(count)

			results = append(results, ResultRow{
				Objective:          s.Objective,
				SimAvgEdges:        avgEdges,
				SimAvgNodes:        avgNodes,
				SimBest1000Edges:   sBest1000E,
				SimBest1000Nodes:   sBest1000N,
				SimBestKnownEdges:  sBestKnownE,
				SimBestKnownNodes:  sBestKnownN,
				BestKnownObjective: bestKnown.Objective,
			})
		}

		// Save to CSV
		csvFile, err := os.Create(fmt.Sprintf("convexity_results_%s.csv", inst.Name))
		if err != nil {
			log.Fatal(err)
		}
		writer := csv.NewWriter(csvFile)
		writer.Write([]string{
			"Objective",
			"SimAvgEdges", "SimAvgNodes",
			"SimBest1000Edges", "SimBest1000Nodes",
			"SimBestKnownEdges", "SimBestKnownNodes",
			"IsBestOf1000", // Flag to help filtering
			"BestKnownObjective",
		})

		for i, r := range results {
			isBest := (solutions[i].Objective == bestOf1000.Objective) // Approximate check
			// Better check: compare pointers or index if we tracked it, but objective equality is likely unique enough or fine if multiple
			// Actually, let's just mark the one we used.
			// But wait, bestOf1000 is a copy.
			// Let's just use objective.

			writer.Write([]string{
				fmt.Sprintf("%d", r.Objective),
				fmt.Sprintf("%.4f", r.SimAvgEdges),
				fmt.Sprintf("%.4f", r.SimAvgNodes),
				fmt.Sprintf("%d", r.SimBest1000Edges),
				fmt.Sprintf("%d", r.SimBest1000Nodes),
				fmt.Sprintf("%d", r.SimBestKnownEdges),
				fmt.Sprintf("%d", r.SimBestKnownNodes),
				fmt.Sprintf("%t", isBest),
				fmt.Sprintf("%d", r.BestKnownObjective),
			})
		}
		writer.Flush()
		csvFile.Close()

		// Calculate Correlations
		// 1. Avg Edges
		printCorrelation(results, "Avg Edges", func(r ResultRow) float64 { return r.SimAvgEdges }, nil)
		// 2. Avg Nodes
		printCorrelation(results, "Avg Nodes", func(r ResultRow) float64 { return r.SimAvgNodes }, nil)

		// Filter out BestOf1000 for its own correlations
		filterBest := func(r ResultRow) bool { return r.Objective == bestOf1000.Objective } // Simple filter

		// 3. Best1000 Edges
		printCorrelation(results, "Best1000 Edges", func(r ResultRow) float64 { return float64(r.SimBest1000Edges) }, filterBest)
		// 4. Best1000 Nodes
		printCorrelation(results, "Best1000 Nodes", func(r ResultRow) float64 { return float64(r.SimBest1000Nodes) }, filterBest)

		// 5. BestKnown Edges
		printCorrelation(results, "BestKnown Edges", func(r ResultRow) float64 { return float64(r.SimBestKnownEdges) }, nil) // Don't filter best known (it's likely not in the set, or if it is, it's rare)
		// 6. BestKnown Nodes
		printCorrelation(results, "BestKnown Nodes", func(r ResultRow) float64 { return float64(r.SimBestKnownNodes) }, nil)

		fmt.Println("------------------------------------------------")
	}

	// Automatically regenerate plots
	log.Println("Regenerating plots...")
	cmd := exec.Command("python3", "plot_convexity.py")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Printf("Error running plot script: %v", err)
	} else {
		log.Println("Plots regenerated successfully.")
	}
}

func commonEdges(e1, e2 map[[2]int]bool) int {
	count := 0
	for e := range e1 {
		if e2[e] {
			count++
		}
	}
	return count
}

func commonNodes(n1, n2 map[int]bool) int {
	count := 0
	for n := range n1 {
		if n2[n] {
			count++
		}
	}
	return count
}

type ResultRow struct {
	Objective          int
	SimAvgEdges        float64
	SimAvgNodes        float64
	SimBest1000Edges   int
	SimBest1000Nodes   int
	SimBestKnownEdges  int
	SimBestKnownNodes  int
	BestKnownObjective int
}

func printCorrelation(results []ResultRow, name string, valueExtractor func(r ResultRow) float64, excludeFilter func(r ResultRow) bool) {

	var xs, ys []float64
	for _, r := range results {
		if excludeFilter != nil && excludeFilter(r) {
			continue
		}
		xs = append(xs, float64(r.Objective))
		ys = append(ys, valueExtractor(r))
	}

	corr := correlation(xs, ys)
	fmt.Printf("Correlation (%s): %.4f\n", name, corr)
}

func correlation(xs, ys []float64) float64 {
	n := float64(len(xs))
	if n < 2 {
		return 0
	}

	sumX, sumY := 0.0, 0.0
	for _, x := range xs {
		sumX += x
	}
	for _, y := range ys {
		sumY += y
	}
	meanX, meanY := sumX/n, sumY/n

	num := 0.0
	denX := 0.0
	denY := 0.0

	for i := 0; i < len(xs); i++ {
		dx := xs[i] - meanX
		dy := ys[i] - meanY
		num += dx * dy
		denX += dx * dx
		denY += dy * dy
	}

	if denX == 0 || denY == 0 {
		return 0
	}

	return num / math.Sqrt(denX*denY)
}

func calculateObjective(D [][]int, costs []int, path []int) int {
	if len(path) == 0 {
		return math.MaxInt32 / 4
	}
	sum := 0
	n := len(path)
	for i := 0; i < n; i++ {
		a := path[i]
		b := path[(i+1)%n]
		sum += D[a][b]
	}
	for _, v := range path {
		sum += costs[v]
	}
	return sum
}
