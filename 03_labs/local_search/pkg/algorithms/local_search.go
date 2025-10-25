package algorithms

import (
	"log"
	"math"
	"math/rand"
	"time"
)

type LSType int
type IntraType int
type StartType int

const (
	LS_Steepest LSType = iota
	LS_Greedy
)

const (
	IntraSwap IntraType = iota
	Intra2Opt
)

const (
	StartRandom StartType = iota
	StartGreedy
)

type MethodSpec struct {
	LS    LSType
	Intra IntraType
	Start StartType
	Name  string
}

func objective(D [][]int, costs []int, path []int) int {
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

func prevIdx(i, n int) int {
	if i == 0 {
		return n - 1
	}
	return i - 1
}
func nextIdx(i, n int) int {
	return (i + 1) % n
}

// START TYPES

// choose 50% of the nodes - ceil(n/2)
func selectCount(n int) int { return (n + 1) / 2 }

// random selection of K nodes and random permutation (cycle)
func startRandom(D [][]int, costs []int, rng *rand.Rand) Solution {
	n := len(D)
	k := selectCount(n)

	idx := make([]int, n)
	for i := range idx {
		idx[i] = i
	}
	rng.Shuffle(n, func(i, j int) { idx[i], idx[j] = idx[j], idx[i] })
	path := append([]int(nil), idx[:k]...)
	rng.Shuffle(k, func(i, j int) { path[i], path[j] = path[j], path[i] })

	return Solution{Path: path, Objective: objective(D, costs, path)}
}

// make a greedy path (using NN Any Position algortihm) - from a chosen node - as a start for LS
func startGreedy(distanceMatrix [][]int, nodeCosts []int, startNodeIndex, k int) Solution {
	n := len(nodeCosts)

	path := []int{startNodeIndex}
	unvisited := make(map[int]bool)
	for j := 0; j < n; j++ {
		if j != startNodeIndex {
			unvisited[j] = true
		}
	}

	for len(path) < k {
		minIncrease := math.MaxInt32
		bestNodeIndex := -1
		bestPosition := -1

		for nodeIndex := range unvisited {
			localMinIncrease := math.MaxInt32
			localBestPos := -1

			if len(path) == 1 {
				// Path has only one node -> we can only insert before or after it
				inc := distanceMatrix[nodeIndex][path[0]] + nodeCosts[nodeIndex]
				if inc < localMinIncrease {
					localMinIncrease = inc
					localBestPos = 0
				}
				inc = distanceMatrix[path[0]][nodeIndex] + nodeCosts[nodeIndex]
				if inc < localMinIncrease {
					localMinIncrease = inc
					localBestPos = 1
				}
			} else {
				// Insert at the beginning
				incFront := distanceMatrix[nodeIndex][path[0]] + nodeCosts[nodeIndex]
				if incFront < localMinIncrease {
					localMinIncrease = incFront
					localBestPos = 0
				}
				// Insert at the end
				incBack := distanceMatrix[path[len(path)-1]][nodeIndex] + nodeCosts[nodeIndex]
				if incBack < localMinIncrease {
					localMinIncrease = incBack
					localBestPos = len(path)
				}
				// Insert in the middle
				for pos := 1; pos < len(path); pos++ {
					a := path[pos-1]
					b := path[pos]
					deltaDist := distanceMatrix[a][nodeIndex] + distanceMatrix[nodeIndex][b] - distanceMatrix[a][b]
					inc := deltaDist + nodeCosts[nodeIndex]
					if inc < localMinIncrease {
						localMinIncrease = inc
						localBestPos = pos
					}
				}
			}

			if localMinIncrease < minIncrease {
				minIncrease = localMinIncrease
				bestNodeIndex = nodeIndex
				bestPosition = localBestPos
			}
		}

		if bestNodeIndex != -1 {
			if bestPosition == len(path) {
				path = append(path, bestNodeIndex)
			} else {
				path = append(path[:bestPosition], append([]int{bestNodeIndex}, path[bestPosition:]...)...)
			}
			delete(unvisited, bestNodeIndex)
		} else {
			break
		}

	}
	return Solution{Path: path, Objective: objective(distanceMatrix, nodeCosts, path)}
}

// DELTAS FOR DIFFERENT MOVES

// intra-route move - two-nodes exchange: change path[i] with path[j]
func deltaSwap(D [][]int, path []int, i, j int) int {
	if i == j {
		return 0
	}
	n := len(path)
	if i > j {
		i, j = j, i
	}

	a, b := path[i], path[j]
	im1, ip1 := path[prevIdx(i, n)], path[nextIdx(i, n)]
	jm1, jp1 := path[prevIdx(j, n)], path[nextIdx(j, n)]

	// Adjacent along the cycle
	if nextIdx(i, n) == j { // ... im1 -> a -> b -> jp1 ...
		before := D[im1][a] + D[a][b] + D[b][jp1]
		after := D[im1][b] + D[b][a] + D[a][jp1]
		return after - before
	}
	if nextIdx(j, n) == i { // ... jm1 -> b -> a -> ip1 ...
		before := D[jm1][b] + D[b][a] + D[a][ip1]
		after := D[jm1][a] + D[a][b] + D[b][ip1]
		return after - before
	}

	// Non-adjacent case: four edges change
	before := D[im1][a] + D[a][ip1] + D[jm1][b] + D[b][jp1]
	after := D[im1][b] + D[b][ip1] + D[jm1][a] + D[a][jp1]
	return after - before
}

// intra-route move - two edges exchange: 2-opt between path[i] and path[j]
func deltaTwoOpt(D [][]int, path []int, i, j int) int {
	if i == j {
		return 0
	}
	n := len(path)
	// Adjacent edges -> degenerate 2-opt (no change)
	if nextIdx(i, n) == j || nextIdx(j, n) == i {
		return 0
	}
	a := path[i]
	b := path[nextIdx(i, n)]
	c := path[j]
	d := path[nextIdx(j, n)]
	before := D[a][b] + D[c][d]
	after := D[a][c] + D[b][d]
	return after - before
}

// inter-route move - two-nodes exchange - path[i] with u (u outside the current path)
func deltaExchangeSelected(distanceMatrix [][]int, nodeCosts []int, path []int, i int, u int) int {
	n := len(path)
	a := path[prevIdx(i, n)]
	v := path[i]
	b := path[nextIdx(i, n)]

	before := distanceMatrix[a][v] + distanceMatrix[v][b] + nodeCosts[v]
	after := distanceMatrix[a][u] + distanceMatrix[u][b] + nodeCosts[u]
	return after - before
}

func applySwap(path []int, i, j int) { path[i], path[j] = path[j], path[i] }

func applyTwoOpt(path []int, i, j int) {
	if i > j {
		i, j = j, i
	}
	for l, r := i+1, j; l < r; l, r = l+1, r-1 {
		path[l], path[r] = path[r], path[l]
	}
}

func applyExchangeSelected(path []int, i int, u int) { path[i] = u }

// LS TYPES - Steepest / Greedy

func localSearchSteepest(distanceMatrix [][]int, nodeCosts []int, init Solution, intra IntraType) Solution {
	path := append([]int(nil), init.Path...)
	n := len(path)

	for {
		bestDelta := 0
		bestMove := func() {}

		// INTRA
		switch intra {
		case IntraSwap:
			for i := 0; i < n; i++ {
				for j := i + 1; j < n; j++ {
					dl := deltaSwap(distanceMatrix, path, i, j)
					if dl < bestDelta {
						ii, jj := i, j
						bestDelta = dl
						bestMove = func() { applySwap(path, ii, jj) }
					}
				}
			}
		case Intra2Opt:
			for i := 0; i < n; i++ {
				for j := i + 1; j < n; j++ {
					dl := deltaTwoOpt(distanceMatrix, path, i, j)
					if dl < bestDelta {
						ii, jj := i, j
						bestDelta = dl
						bestMove = func() { applyTwoOpt(path, ii, jj) }
					}
				}
			}
		}

		// INTER
		inSel := make([]bool, len(distanceMatrix))
		for _, v := range path {
			inSel[v] = true
		}
		nonSel := make([]int, 0, len(distanceMatrix)-n)
		for u := range distanceMatrix {
			if !inSel[u] {
				nonSel = append(nonSel, u)
			}
		}
		for i := 0; i < n; i++ {
			for _, u := range nonSel {
				dl := deltaExchangeSelected(distanceMatrix, nodeCosts, path, i, u)
				if dl < bestDelta {
					ii, uu := i, u
					bestDelta = dl
					bestMove = func() { applyExchangeSelected(path, ii, uu) }
				}
			}
		}
		if bestDelta < 0 {
			bestMove()
		} else {
			break
		}
	}
	return Solution{Path: path, Objective: objective(distanceMatrix, nodeCosts, path)}
}

func localSearchGreedy(distanceMatrix [][]int, nodeCosts []int, init Solution, intra IntraType, rng *rand.Rand) Solution {
	path := append([]int(nil), init.Path...)
	n := len(path)

	for {
		improved := false

		// Random order of neighborhood types (0=intra,1=inter)
		order := []int{0, 1}
		rng.Shuffle(2, func(i, j int) { order[i], order[j] = order[j], order[i] })

		tryIntra := func() bool {
			switch intra {
			case IntraSwap:
				pi := randPerm(rng, n)
				for _, i := range pi {
					pj := randPermFrom(rng, i+1, n)
					for _, j := range pj {
						if deltaSwap(distanceMatrix, path, i, j) < 0 {
							applySwap(path, i, j)
							return true
						}
					}
				}
			case Intra2Opt:
				pi := randPerm(rng, n)
				for _, i := range pi {
					pj := randPermFrom(rng, i+1, n)
					for _, j := range pj {
						if deltaTwoOpt(distanceMatrix, path, i, j) < 0 {
							applyTwoOpt(path, i, j)
							return true
						}
					}
				}
			}
			return false
		}

		tryInter := func() bool {
			inSel := make([]bool, len(distanceMatrix))
			for _, v := range path {
				inSel[v] = true
			}
			nonSel := make([]int, 0, len(distanceMatrix)-n)
			for u := range distanceMatrix {
				if !inSel[u] {
					nonSel = append(nonSel, u)
				}
			}
			rng.Shuffle(len(nonSel), func(i, j int) { nonSel[i], nonSel[j] = nonSel[j], nonSel[i] })
			pi := randPerm(rng, n)
			for _, i := range pi {
				for _, u := range nonSel {
					if deltaExchangeSelected(distanceMatrix, nodeCosts, path, i, u) < 0 {
						applyExchangeSelected(path, i, u)
						return true
					}
				}
			}
			return false
		}

		for _, which := range order {
			if which == 0 {
				if tryIntra() {
					improved = true
					break
				}
			} else {
				if tryInter() {
					improved = true
					break
				}
			}
		}
		if !improved {
			break
		}
	}
	return Solution{Path: path, Objective: objective(distanceMatrix, nodeCosts, path)}
}

func randPerm(r *rand.Rand, n int) []int {
	p := make([]int, n)
	for i := range p {
		p[i] = i
	}
	r.Shuffle(n, func(i, j int) { p[i], p[j] = p[j], p[i] })
	return p
}
func randPermFrom(r *rand.Rand, start, n int) []int {
	if start >= n {
		return nil
	}
	m := n - start
	p := make([]int, m)
	for i := 0; i < m; i++ {
		p[i] = start + i
	}
	r.Shuffle(m, func(i, j int) { p[i], p[j] = p[j], p[i] })
	return p
}

// BATCH RUNNING OF LS

// startNodeIndices:
//   - for StartGreedy we will use consecutive indices as starting nodes for greedy construction
//   - for StartRandom we ignore this list (generating numSolutions randomly)
func RunLocalSearchBatch(
	distanceMatrix [][]int,
	nodeCosts []int,
	startNodeIndices []int,
	m MethodSpec,
	numSolutions int,
) []Solution {
	if numSolutions <= 0 {
		return nil
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	n := len(distanceMatrix)
	k := selectCount(n)

	results := make([]Solution, 0, numSolutions)

	switch m.Start {
	case StartRandom:
		for r := 0; r < numSolutions; r++ {
			init := startRandom(distanceMatrix, nodeCosts, rng)
			var sol Solution
			switch m.LS {
			case LS_Steepest:
				sol = localSearchSteepest(distanceMatrix, nodeCosts, init, m.Intra)
			case LS_Greedy:
				sol = localSearchGreedy(distanceMatrix, nodeCosts, init, m.Intra, rng)
			}
			results = append(results, sol)
		}
	case StartGreedy:
		for r := 0; r < numSolutions; r++ {
			start := 0
			if len(startNodeIndices) > 0 {
				start = startNodeIndices[r%len(startNodeIndices)] % n
			} else {
				start = r % n
			}
			init := startGreedy(distanceMatrix, nodeCosts, start, k)
			var sol Solution
			switch m.LS {
			case LS_Steepest:
				sol = localSearchSteepest(distanceMatrix, nodeCosts, init, m.Intra)
			case LS_Greedy:
				sol = localSearchGreedy(distanceMatrix, nodeCosts, init, m.Intra, rng)
			}
			results = append(results, sol)

			if sol.Objective > init.Objective {
				log.Fatalf("LS worsened solution: init=%d, after=%d", init.Objective, sol.Objective)
			}

		}
	default:
		for r := 0; r < numSolutions; r++ {
			init := startRandom(distanceMatrix, nodeCosts, rng)
			var sol Solution
			switch m.LS {
			case LS_Steepest:
				sol = localSearchSteepest(distanceMatrix, nodeCosts, init, m.Intra)
			case LS_Greedy:
				sol = localSearchGreedy(distanceMatrix, nodeCosts, init, m.Intra, rng)
			}
			results = append(results, sol)
		}
	}
	return results
}
