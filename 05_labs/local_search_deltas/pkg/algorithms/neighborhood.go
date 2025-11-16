package algorithms

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
func deltaExchangeSelected(D [][]int, costs []int, path []int, i int, u int) int {
	n := len(path)
	a := path[prevIdx(i, n)]
	v := path[i]
	b := path[nextIdx(i, n)]
	before := D[a][v] + D[v][b] + costs[v]
	after := D[a][u] + D[u][b] + costs[u]
	return after - before
}

// applyTwoOpt performs a 2-opt move on the path between indices i and j, in-place.
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

// applyTwoOptAndUpdatePos performs a 2-opt move and keeps the position index
// array `posOf` in sync with the modified path.
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

// applyExchangeSelected replaces the selected vertex at position i with a new
// vertex u (which must be outside the current path).
func applyExchangeSelected(path []int, i int, u int) { path[i] = u }
