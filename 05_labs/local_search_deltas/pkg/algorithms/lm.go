package algorithms

// MoveType distinguishes between 2-opt and exchange moves in the LM structure.
type MoveType int

const (
	MoveTwoOpt MoveType = iota
	MoveExchangeSelected
)

// edgeKey is a canonical representation of an undirected edge (x,y) with x < y.
type edgeKey struct {
	x, y int
}

type moveKey struct {
	e1 edgeKey
	e2 edgeKey
}

// MoveRecord stores a single improving move together with its precomputed delta.
type MoveRecord struct {
	kind  MoveType
	a, b  int // endpoints of first removed edge
	c, d  int // endpoints of second removed edge
	v, u  int // for exchange: v replaced by u (selected vertex v, new vertex u)
	delta int // precomputed delta value
	key   moveKey
}

// lmState stores a list-of-moves (LM) and an index to find moves by their edge keys.
type lmState struct {
	moves []MoveRecord
	index map[moveKey]int
}

// buildFullNeighborhoodLM builds the full improving neighborhood for the
// current solution and stores it in the LM structure. This is called once
// for the initial solution; subsequent iterations update LM incrementally.
func buildFullNeighborhoodLM(D [][]int, costs []int, path []int, nonSel []int, lm *lmState) {
	n := len(path)

	// intra: 2-opt
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			dl := deltaTwoOpt(D, path, i, j)
			if dl >= 0 {
				continue
			}
			a := path[i]
			b := path[nextIdx(i, n)]
			c := path[j]
			d := path[nextIdx(j, n)]
			rec := MoveRecord{
				kind:  MoveTwoOpt,
				a:     a,
				b:     b,
				c:     c,
				d:     d,
				v:     -1,
				u:     -1,
				delta: dl,
			}
			lm.addMove(rec)
		}
	}

	// inter: selected vertex with unselected one
	for i := 0; i < n; i++ {
		for _, u := range nonSel {
			dl := deltaExchangeSelected(D, costs, path, i, u)
			if dl >= 0 {
				continue
			}
			a := path[prevIdx(i, n)]
			v := path[i]
			b := path[nextIdx(i, n)]
			rec := MoveRecord{
				kind:  MoveExchangeSelected,
				a:     a,
				b:     v,
				c:     v,
				d:     b,
				v:     v,
				u:     u,
				delta: dl,
			}
			lm.addMove(rec)
		}
	}
}

// updateLMAfterMove updates the LM after applying a move by generating new
// improving moves in the vicinity of the modified edges/vertices instead of
// rebuilding the full neighborhood.
func updateLMAfterMove(D [][]int, costs []int, path []int, posOf []int, nonSel []int, bestMove MoveRecord, lm *lmState) {
	n := len(path)
	if n == 0 {
		return
	}

	edgeStarts := make(map[int]struct{}, 8)

	switch bestMove.kind {
	case MoveTwoOpt:
		i := bestMove.v
		j := bestMove.u
		if i < 0 || j < 0 || i >= n || j >= n {
			break
		}
		if i == j {
			break
		}
		if i > j {
			i, j = j, i
		}
		for k := i; k <= j; k++ {
			edgeStarts[k] = struct{}{}
			edgeStarts[prevIdx(k, n)] = struct{}{}
		}
	case MoveExchangeSelected:
		uNew := bestMove.u
		if uNew < 0 {
			break
		}
		pos := posOf[uNew]
		if pos < 0 || pos >= n {
			break
		}
		edgeStarts[pos] = struct{}{}
		edgeStarts[prevIdx(pos, n)] = struct{}{}
	}

	if len(edgeStarts) == 0 {
		return
	}

	// 2-opt moves touching affected edges.
	for i := range edgeStarts {
		for j := 0; j < n; j++ {
			if j == i {
				continue
			}
			if j == nextIdx(i, n) || j == prevIdx(i, n) {
				continue
			}
			dl := deltaTwoOpt(D, path, i, j)
			if dl >= 0 {
				continue
			}
			a := path[i]
			b := path[nextIdx(i, n)]
			c := path[j]
			d := path[nextIdx(j, n)]
			rec := MoveRecord{
				kind:  MoveTwoOpt,
				a:     a,
				b:     b,
				c:     c,
				d:     d,
				v:     -1,
				u:     -1,
				delta: dl,
			}
			lm.addMove(rec)
		}
	}

	// Exchange moves for positions adjacent to affected edges.
	for i := range edgeStarts {
		for _, u := range nonSel {
			dl := deltaExchangeSelected(D, costs, path, i, u)
			if dl >= 0 {
				continue
			}
			a := path[prevIdx(i, n)]
			v := path[i]
			b := path[nextIdx(i, n)]
			rec := MoveRecord{
				kind:  MoveExchangeSelected,
				a:     a,
				b:     v,
				c:     v,
				d:     b,
				v:     v,
				u:     u,
				delta: dl,
			}
			lm.addMove(rec)
		}
	}
}

func canonicalEdge(x, y int) edgeKey {
	if x > y {
		x, y = y, x
	}
	return edgeKey{x: x, y: y}
}

func canonicalMoveKey(e1, e2 edgeKey) moveKey {
	// order edges lexicographically to keep key canonical
	if e2.x < e1.x || (e2.x == e1.x && e2.y < e1.y) {
		e1, e2 = e2, e1
	}
	return moveKey{e1: e1, e2: e2}
}

func (lm *lmState) addMove(rec MoveRecord) {
	e1 := canonicalEdge(rec.a, rec.b)
	e2 := canonicalEdge(rec.c, rec.d)
	key := canonicalMoveKey(e1, e2)
	if _, exists := lm.index[key]; exists {
		return
	}
	rec.key = key
	lm.moves = append(lm.moves, rec)
	lm.index[key] = len(lm.moves) - 1
}

func (lm *lmState) remove(rec MoveRecord) {
	if len(lm.moves) == 0 {
		return
	}
	key := rec.key
	idx, ok := lm.index[key]
	if !ok {
		return
	}
	lastIdx := len(lm.moves) - 1
	if idx != lastIdx {
		// swap with last and update index
		lastRec := lm.moves[lastIdx]
		lm.moves[idx] = lastRec
		lm.index[lastRec.key] = idx
	}
	lm.moves = lm.moves[:lastIdx]
	delete(lm.index, key)
}

// findEdgeCut finds whether an undirected edge (x,y) appears in the current cycle defined by path/posOf.
// It returns:
//
//	ok       - true if the edge exists,
//	forward  - true if along the tour it goes x->y, false if y->x,
//	cutIndex - index i such that the removed edge can be represented as (path[i], path[next(i)]).
func findEdgeCut(path []int, posOf []int, x, y int) (ok bool, forward bool, cutIndex int) {
	n := len(path)
	px := posOf[x]
	py := posOf[y]

	// If neither endpoint is on the tour, the edge cannot appear.
	if px < 0 && py < 0 {
		return false, false, 0
	}

	if px >= 0 {
		if path[nextIdx(px, n)] == y {
			// edge x->y, cut at x
			return true, true, px
		}
		if path[prevIdx(px, n)] == y && py >= 0 {
			// along the tour the edge is y->x, so cutting at position of y
			return true, false, py
		}
	}

	if py >= 0 {
		if path[nextIdx(py, n)] == x {
			// along the tour the edge is y->x, cut at y
			return true, false, py
		}
		if path[prevIdx(py, n)] == x && px >= 0 {
			// along the tour the edge is x->y, cut at x
			return true, true, px
		}
	}

	return false, false, 0
}

// localSearchSteepestLM performs steepest local search with list-of-moves (LM)
// delta reuse, using the same neighborhood as the baseline variant.
func localSearchSteepestLM(D [][]int, costs []int, init Solution) Solution {
	path := append([]int(nil), init.Path...)
	n := len(path)
	if n == 0 {
		return init
	}

	dim := len(D)

	// quick lookup structures
	posOf := make([]int, dim)
	inSel := make([]bool, dim)
	for i := range posOf {
		posOf[i] = -1
	}
	for i, v := range path {
		posOf[v] = i
		inSel[v] = true
	}

	// heuristic preallocation: typical number of moves is O(n^2)
	prealloc := n * n
	lm := lmState{
		moves: make([]MoveRecord, 0, prealloc),
		index: make(map[moveKey]int, prealloc),
	}

	// maintain the list of non-selected vertices incrementally.
	nonSel := make([]int, 0, dim-n)
	for u := 0; u < dim; u++ {
		if !inSel[u] {
			nonSel = append(nonSel, u)
		}
	}

	// removeFromSlice removes value from xs using swap-and-pop (O(1)).
	// Order is not preserved, but this is acceptable for nonSel.
	removeFromSlice := func(xs []int, value int) []int {
		for i, v := range xs {
			if v == value {
				lastIdx := len(xs) - 1
				xs[i] = xs[lastIdx]
				return xs[:lastIdx]
			}
		}
		return xs
	}

	// Build full improving neighborhood once for the initial solution.
	buildFullNeighborhoodLM(D, costs, path, nonSel, &lm)

	for {
		bestDelta := 0
		var bestMove MoveRecord
		hasBest := false

		// 1) Browse LM, reusing stored deltas when applicable.
		//    While browsing, aggressively prune moves that are no longer
		//    applicable so we don't keep scanning dead entries.
		for idx := 0; idx < len(lm.moves); {
			rec := lm.moves[idx]
			removed := false
			switch rec.kind {
			case MoveTwoOpt:
				ok1, fwd1, cut1 := findEdgeCut(path, posOf, rec.a, rec.b)
				ok2, fwd2, cut2 := findEdgeCut(path, posOf, rec.c, rec.d)
				if !ok1 || !ok2 {
					lm.remove(rec)
					removed = true
				} else if fwd1 != fwd2 {
					// orientation changed relative to when the move was created, so
					// the stored delta is no longer valid; let it be regenerated
					// in the neighborhood phase and drop this stale entry.
					lm.remove(rec)
					removed = true
				} else {
					// same relative direction (both forward or both reversed) -> applicable now
					if rec.delta < bestDelta || !hasBest {
						hasBest = true
						bestDelta = rec.delta
						// normalise cut indices according to current orientation
						bestMove = rec
						// store current cut indices into unused fields v,u for convenience
						bestMove.v = cut1
						bestMove.u = cut2
					}
				}
			case MoveExchangeSelected:
				// Early checks before expensive findEdgeCut calls.
				// Verify v is still selected and u is not.
				if rec.v < 0 || rec.v >= dim || posOf[rec.v] == -1 {
					lm.remove(rec)
					removed = true
				} else if rec.u >= 0 && rec.u < dim && inSel[rec.u] {
					lm.remove(rec)
					removed = true
				} else {
					ok1, fwd1, _ := findEdgeCut(path, posOf, rec.a, rec.b)
					ok2, fwd2, _ := findEdgeCut(path, posOf, rec.c, rec.d)
					if !ok1 || !ok2 {
						lm.remove(rec)
						removed = true
					} else if fwd1 != fwd2 {
						lm.remove(rec)
						removed = true
					} else if rec.delta < bestDelta || !hasBest {
						hasBest = true
						bestDelta = rec.delta
						bestMove = rec
					}
				}
			}
			if !removed {
				idx++
			}
		}

		// 2) New moves are added incrementally in updateLMAfterMove after an
		// improving move is applied, so we do not rebuild the full
		// neighborhood here.
		if !hasBest || bestDelta >= 0 {
			break
		}

		// 3) Apply the best move and update structures; then update LM.
		switch bestMove.kind {
		case MoveTwoOpt:
			// indices were stored in v,u when best was selected
			i := bestMove.v
			j := bestMove.u
			applyTwoOptAndUpdatePos(path, posOf, i, j)
		case MoveExchangeSelected:
			posV := posOf[bestMove.v]
			if posV >= 0 {
				vOld := bestMove.v
				uNew := bestMove.u
				applyExchangeSelected(path, posV, uNew)
				inSel[vOld], inSel[uNew] = false, true
				posOf[vOld], posOf[uNew] = -1, posV

				// keep nonSel consistent with the incremental update policy.
				nonSel = removeFromSlice(nonSel, uNew)
				nonSel = append(nonSel, vOld)
			}
		}

		// remove the applied move from LM using its removed edges
		lm.remove(bestMove)

		// Incrementally add new moves affected by this modification instead of
		// rebuilding the neighborhood from scratch.
		updateLMAfterMove(D, costs, path, posOf, nonSel, bestMove, &lm)
	}

	return Solution{Path: path, Objective: objective(D, costs, path)}
}
