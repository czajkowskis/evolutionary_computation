package algorithms

import (
	"math"
	"math/rand"
	"sort"
	"time"
)

type LSType int
type IntraType int
type StartType int

// MethodSpec defines the configuration for a local search run
type MethodSpec struct {
	Name    string
	UseCand bool // should use candidate moves?
	CandK   int  // how many nearest to include in candidate list
}

// objective computes the total cost: sum of edge distances in the cycle + sum of node costs
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
func nextIdx(i, n int) int { return (i + 1) % n }

func selectCount(n int) int { return (n + 1) / 2 }

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

// deltaTwoOpt computes the change in objective for a 2-opt move (intra-route).
// Removes edges (path[i], path[i+1]) and (path[j], path[j+1]),
// adds edges (path[i], path[j]) and (path[i+1], path[j+1]).
// Note: node costs do not change, only edge costs.
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

// deltaExchangeSelected computes the change in objective for swapping path[i] with node u (inter-route).
// Removes edges (prev, path[i]) and (path[i], next), removes cost[path[i]],
// adds edges (prev, u) and (u, next), adds cost[u].
func deltaExchangeSelected(D [][]int, costs []int, path []int, i int, u int) int {
	n := len(path)
	a := path[prevIdx(i, n)]
	v := path[i]
	b := path[nextIdx(i, n)]
	before := D[a][v] + D[v][b] + costs[v]
	after := D[a][u] + D[u][b] + costs[u]
	return after - before
}

func applyTwoOpt(path []int, i, j int) {
	n := len(path)
	if i == j || nextIdx(i, n) == j || nextIdx(j, n) == i {
		return
	}
	if i > j {
		i, j = j, i
	}
	for l, r := i+1, j; l < r; l, r = l+1, r-1 {
		path[l], path[r] = path[r], path[l]
	}
}

func applyTwoOptAndUpdatePos(path []int, posOf []int, i, j int) {
	n := len(path)
	if i == j || nextIdx(i, n) == j || nextIdx(j, n) == i {
		return
	}
	if i > j {
		i, j = j, i
	}
	// reverse segment [i+1..j]
	for l, r := i+1, j; l < r; l, r = l+1, r-1 {
		vl, vr := path[l], path[r]
		path[l], path[r] = vr, vl
		posOf[vl], posOf[vr] = r, l
	}
}

func applyExchangeSelected(path []int, i int, u int) { path[i] = u }

// Candidate data structure

type CandData struct {
	CandList [][]int
	isCand   map[uint64]struct{}
}

func packEdge(a, b int) uint64 {
	if a > b {
		a, b = b, a
	}
	return uint64(uint32(a))<<32 | uint64(uint32(b))
}

// Build K nearest per node by weight = D[u][v] + costs[v]
// For each node u, we find K nearest neighbors v ranked by edge weight + destination node cost
func buildCandidates(D [][]int, costs []int, K int) CandData {
	n := len(D)
	if K <= 0 {
		K = 10
	}
	cand := make([][]int, n)
	isCand := make(map[uint64]struct{}, n*K*2)

	for u := 0; u < n; u++ {
		type nb struct{ v, w int }
		nbs := make([]nb, 0, n-1)
		for v := 0; v < n; v++ {
			if v == u {
				continue
			}
			w := D[u][v] + costs[v]
			nbs = append(nbs, nb{v: v, w: w})
		}
		sort.Slice(nbs, func(i, j int) bool { return nbs[i].w < nbs[j].w })
		m := K
		if m > len(nbs) {
			m = len(nbs)
		}
		list := make([]int, m)
		for i := 0; i < m; i++ {
			v := nbs[i].v
			list[i] = v
			isCand[packEdge(u, v)] = struct{}{}
		}
		cand[u] = list
	}
	return CandData{CandList: cand, isCand: isCand}
}

func isCandidateEdge(cd CandData, a, b int) bool {
	_, ok := cd.isCand[packEdge(a, b)]
	return ok
}

// localSearchSteepestBaseline performs steepest-descent local search without candidate restrictions.
// Considers all 2-opt moves (intra-route) and all node exchanges (inter-route).
func localSearchSteepestBaseline(D [][]int, costs []int, init Solution) Solution {
	path := append([]int(nil), init.Path...)
	n := len(path)

	for {
		bestDelta := 0
		var bestMove func()

		// intra-route move - two-edges exchange: 2-opt
		for i := 0; i < n; i++ {
			for j := i + 1; j < n; j++ {
				dl := deltaTwoOpt(D, path, i, j)
				if dl < bestDelta {
					ii, jj := i, j
					bestDelta = dl
					bestMove = func() { applyTwoOpt(path, ii, jj) }
				}
			}
		}

		// inter-route moves - path[i] with u (u outside the current path)
		inSel := make([]bool, len(D))
		for _, v := range path {
			inSel[v] = true
		}
		nonSel := make([]int, 0, len(D)-n)
		for u := range D {
			if !inSel[u] {
				nonSel = append(nonSel, u)
			}
		}
		for i := 0; i < n; i++ {
			for _, u := range nonSel {
				dl := deltaExchangeSelected(D, costs, path, i, u)
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
	return Solution{Path: path, Objective: objective(D, costs, path)}
}

// localSearchSteepestCandidates performs steepest-descent local search with candidate move restrictions.
// Candidate moves prune the search space using K-nearest neighbor lists:
//   - Intra (2-opt): For each node n1 in cycle and each n2 in CandList[n1] also in cycle,
//     consider two 2-opt moves that introduce edge (n1, n2):
//     A) 2-opt(i, j) removes (n1, next(n1)) and (n2, next(n2))
//     B) 2-opt(prev(i), prev(j)) removes (prev(n1), n1) and (prev(n2), n2)
//   - Inter (exchange): For position i, only consider swapping with nodes u in CandList[prev(i)] or CandList[next(i)].
//
// This reduces evaluations while focusing on promising edges based on distance + node cost.
func localSearchSteepestCandidates(D [][]int, costs []int, init Solution, cd CandData) Solution {
	path := append([]int(nil), init.Path...)
	n := len(path)

	// quick lookup structures
	posOf := make([]int, len(D))
	inSel := make([]bool, len(D))
	for i := range posOf {
		posOf[i] = -1
	}
	for i, v := range path {
		posOf[v] = i
		inSel[v] = true
	}

	visitMark := make([]int, len(D))
	epoch := 0

	for {
		bestDelta := 0
		var bestMove func()

		// intra
		for i := 0; i < n; i++ {
			n1 := path[i]
			ip1 := nextIdx(i, n)
			im1 := prevIdx(i, n)

			for _, n2 := range cd.CandList[n1] {
				j := posOf[n2]
				if j == -1 {
					continue // n2 is not selected -> this is not intra, but inter
				}
				if j == i || j == ip1 || j == im1 {
					continue // neighbors/degenerate
				}

				// Prune symmetric duplicates: evaluate pair only when j > i for move A.
				if j > i {
					// MOVE A: 2-opt(i, j)  (cuts (i,i+1) & (j,j+1))
					if dlA := deltaTwoOpt(D, path, i, j); dlA < bestDelta {
						ii, jj := i, j
						bestDelta = dlA
						bestMove = func() {
							applyTwoOptAndUpdatePos(path, posOf, ii, jj)
						}
					}
				}

				// MOVE B: 2-opt(prev(i), prev(j)) (cuts (i-1,i) & (j-1,j))
				ii := prevIdx(i, n)
				jj := prevIdx(j, n)
				if dlB := deltaTwoOpt(D, path, ii, jj); dlB < bestDelta {
					iii, jjj := ii, jj
					bestDelta = dlB
					bestMove = func() {
						applyTwoOptAndUpdatePos(path, posOf, iii, jjj)
					}
				}
			}
		}

		// inter - consider exchanges only for candidate neighbors of the edges being modified
		// For position i, we check candidates of a=prev(i) and b=next(i)
		// Since u is in cd.CandList[a] or cd.CandList[b], at least one introduced edge is a candidate
		for i := 0; i < n; i++ {
			a := path[prevIdx(i, n)]
			b := path[nextIdx(i, n)]

			epoch++
			// Iterate cand[a] - all u here satisfy isCandidateEdge(cd, a, u)
			for _, u := range cd.CandList[a] {
				if visitMark[u] == epoch {
					continue
				}
				visitMark[u] = epoch
				if inSel[u] {
					continue
				}
				if dl := deltaExchangeSelected(D, costs, path, i, u); dl < bestDelta {
					ii, uu := i, u
					bestDelta = dl
					bestMove = func() {
						vOld := path[ii]
						applyExchangeSelected(path, ii, uu)

						inSel[vOld], inSel[uu] = false, true
						posOf[vOld], posOf[uu] = -1, ii
					}
				}
			}
			// Iterate cand[b] - all u here satisfy isCandidateEdge(cd, b, u)
			for _, u := range cd.CandList[b] {
				if visitMark[u] == epoch {
					continue
				}
				visitMark[u] = epoch
				if inSel[u] {
					continue
				}
				if dl := deltaExchangeSelected(D, costs, path, i, u); dl < bestDelta {
					ii, uu := i, u
					bestDelta = dl
					bestMove = func() {
						vOld := path[ii]
						applyExchangeSelected(path, ii, uu)
						inSel[vOld], inSel[uu] = false, true
						posOf[vOld], posOf[uu] = -1, ii
					}
				}
			}
		}

		if bestDelta < 0 {
			bestMove()
		} else {
			break
		}
	}

	return Solution{Path: path, Objective: objective(D, costs, path)}
}

// RunLocalSearchBatch runs local search multiple times with random initializations.
// Returns all solutions found (one per run).
func RunLocalSearchBatch(
	D [][]int,
	costs []int,
	m MethodSpec,
	numSolutions int,
) []Solution {
	if numSolutions <= 0 {
		return nil
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	results := make([]Solution, 0, numSolutions)

	var cd CandData
	if m.UseCand {
		K := m.CandK
		if K <= 0 {
			K = 10
		}
		cd = buildCandidates(D, costs, K)
	}

	for r := 0; r < numSolutions; r++ {
		init := startRandom(D, costs, rng)
		var sol Solution
		if m.UseCand {
			sol = localSearchSteepestCandidates(D, costs, init, cd)
		} else {
			sol = localSearchSteepestBaseline(D, costs, init)
		}
		results = append(results, sol)
	}
	return results
}
