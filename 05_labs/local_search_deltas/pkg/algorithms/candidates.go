package algorithms

import "sort"

// CandData stores precomputed candidate neighbors for each node and a fast
// lookup structure to check whether an undirected edge is a candidate edge.
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

// buildCandidates builds K nearest neighbors for each node, using the weight
// D[u][v] + costs[v]. It also builds a fast lookup map for candidate edges.
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

// localSearchSteepestCandidates performs steepest-descent local search using
// candidate moves (2-opt intra-route and exchanges with unselected vertices).
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

		// inter - allow only if at least one of the introduced edges is a candidate edge (a,u) or (u,b), where a = prev(path[i]), b = next(path[i])
		for i := 0; i < n; i++ {
			a := path[prevIdx(i, n)]
			b := path[nextIdx(i, n)]

			epoch++
			// Iterate cand[a]
			for _, u := range cd.CandList[a] {
				if visitMark[u] == epoch {
					continue
				}
				visitMark[u] = epoch
				if inSel[u] {
					continue
				}
				// must introduce at least one candidate edge
				if !(isCandidateEdge(cd, a, u) || isCandidateEdge(cd, b, u)) {
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
			// Iterate cand[b]
			for _, u := range cd.CandList[b] {
				if visitMark[u] == epoch {
					continue
				}
				visitMark[u] = epoch
				if inSel[u] {
					continue
				}
				if !(isCandidateEdge(cd, a, u) || isCandidateEdge(cd, b, u)) {
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
