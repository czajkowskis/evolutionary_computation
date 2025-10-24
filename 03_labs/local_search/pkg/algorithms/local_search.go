package algorithms

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
)

// ---------------------- interfejs publiczny ----------------------

// Uwaga: typy Solution / FindBestSolution pochodzą z Twoich istniejących plików
// (solution.go, best_solution.go) w tym samym pakiecie `algorithms`.

// MethodSpec — wybór wariantu + opcje
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
	StartGreedy           // jeśli startNodeIndices puste/nil → losujemy start node
)

// Pola opcjonalne: Seed/InterFirstProb/StrictValidate nie zmieniają starego zachowania,
// jeśli zostaną pozostawione z wartościami zero.
type MethodSpec struct {
	LS    LSType
	Intra IntraType
	Start StartType
	Name  string

	// Opcjonalnie: ziarno losowości (domyślnie deterministyczne na bazie kombinacji parametrów).
	Seed int64
	// Prawdopodobieństwo, że w LS-Greedy jako pierwszy będzie rozpatrywany „inter” (0..1). Domyślnie 0.5.
	InterFirstProb float64
	// Jeśli true, walidacja wymaga idealnej symetrii D[i][j]==D[j][i]; w przeciwnym razie toleruje różnice ≤1.
	StrictValidate bool
}

func (m MethodSpec) String() string {
	name := m.Name
	if name == "" {
		ls := map[LSType]string{LS_Steepest: "Steepest", LS_Greedy: "Greedy"}[m.LS]
		intra := map[IntraType]string{IntraSwap: "Swap", Intra2Opt: "2Opt"}[m.Intra]
		start := map[StartType]string{StartRandom: "StartRandom", StartGreedy: "StartGreedy"}[m.Start]
		name = fmt.Sprintf("%s+%s+%s", ls, intra, start)
	}
	return name
}

// ValidateInstance — sprawdza spójność danych wejściowych.
// - D musi być kwadratowe n×n, D[i][i]==0, D[i][j]≥0
// - koszty muszą mieć długość n, C[v]≥0
// - (opcjonalnie) |D[i][j]-D[j][i]| ≤ 1 lub (strict) dokładnie 0
// - k=(n+1)/2 zgodny z definicją (sprawdzenie pochodne — informacyjne)
func ValidateInstance(D [][]int, costs []int, strict bool) error {
	n := len(D)
	if n == 0 {
		// pusty przypadek jest dopuszczalny — LS nic nie zrobi
		return nil
	}
	for i := 0; i < n; i++ {
		if len(D[i]) != n {
			return fmt.Errorf("D must be square: row %d has len=%d, expected %d", i, len(D[i]), n)
		}
	}
	if len(costs) != n {
		return fmt.Errorf("costs length mismatch: got %d, want %d", len(costs), n)
	}
	for i := 0; i < n; i++ {
		if D[i][i] != 0 {
			return fmt.Errorf("D[%d][%d] must be 0 (got %d)", i, i, D[i][i])
		}
		if costs[i] < 0 {
			return fmt.Errorf("costs[%d] must be non-negative (got %d)", i, costs[i])
		}
		for j := 0; j < n; j++ {
			if D[i][j] < 0 {
				return fmt.Errorf("D[%d][%d] must be non-negative (got %d)", i, j, D[i][j])
			}
			if i < j {
				if strict {
					if D[i][j] != D[j][i] {
						return fmt.Errorf("strict symmetry violated at (%d,%d): %d vs %d", i, j, D[i][j], D[j][i])
					}
				} else {
					diff := D[i][j] - D[j][i]
					if diff < 0 {
						diff = -diff
					}
					if diff > 1 {
						return fmt.Errorf("symmetry tolerance exceeded at (%d,%d): %d vs %d", i, j, D[i][j], D[j][i])
					}
				}
			}
		}
	}
	// informacyjnie (nie błąd): k z definicji
	_ = (n + 1) / 2
	return nil
}

// RunLocalSearchBatch uruchamia dany wariant 'runs' razy i zwraca rozwiązania.
// Dodatkowo waliduje instancję i jest odporne na n=0/1/2.
// Jeśli StartGreedy i lista startów pusta/nil — losuje węzeł startowy.
func RunLocalSearchBatch(D [][]int, costs []int, startNodeIndices []int, m MethodSpec, runs int) []Solution {
	// Walidacja (nie panikujemy; zwracamy puste wyniki przy błędzie)
	if err := ValidateInstance(D, costs, m.StrictValidate); err != nil {
		// Możesz tu zamiast „pustych” rzucić panic lub zwrócić błąd — zostaję przy
		// „brak wyników” dla kompatybilności sygnatury.
		// fmt.Println("instance validation error:", err)
		return nil
	}

	n := len(D)
	if n == 0 || runs <= 0 {
		return nil
	}
	k := (n + 1) / 2
	if k <= 0 {
		return nil
	}
	ins := &instance{n: n, k: k, D: D, C: costs}

	// Ustalenie bazowego seeda deterministycznie, ale modyfikowalnie z zewnątrz
	baseSeed := m.Seed
	if baseSeed == 0 {
		// deterministyczna baza po polach metody
		baseSeed = int64(1337 + int(m.LS)*1009 + int(m.Intra)*917 + int(m.Start)*701)
	}

	interFirstProb := m.InterFirstProb
	if interFirstProb <= 0 || interFirstProb >= 1 {
		interFirstProb = 0.5
	}

	solutions := make([]Solution, 0, runs)
	for r := 0; r < runs; r++ {
		fmt.Printf("Node %d:\n", r)
		var s *solution
		switch m.Start {
		case StartRandom:
			rng := rand.New(rand.NewSource(baseSeed + int64(r*7919)))
			s = buildRandomStart(ins, rng)
		case StartGreedy:
			var startNode int
			if len(startNodeIndices) == 0 {
				// Losujemy start node deterministycznie z seeda batcha
				rng := rand.New(rand.NewSource(baseSeed + int64(r*104729)))
				startNode = rng.Intn(n)
			} else {
				startNode = startNodeIndices[r%len(startNodeIndices)]
				if startNode < 0 || startNode >= n {
					// nieprawidłowy indeks — awaryjnie losujemy
					rng := rand.New(rand.NewSource(baseSeed + int64(r*104729)))
					startNode = rng.Intn(n)
				}
			}
			s = buildGreedyStart(ins, startNode)
		default:
			// fallback
			rng := rand.New(rand.NewSource(baseSeed + int64(r*7919)))
			s = buildRandomStart(ins, rng)
		}

		if m.LS == LS_Steepest {
			lsSteepest(s, m.Intra) // steepest nie potrzebuje RNG
		} else {
			rng := rand.New(rand.NewSource(baseSeed + 2025 + int64(r*977)))
			lsGreedy(s, m.Intra, rng, interFirstProb)
		}

		// eksport do publicznego typu Solution (z solution.go)
		solutions = append(solutions, Solution{Path: append([]int(nil), s.T...), Objective: s.obj})
	}
	return solutions
}

// ---------------------- reprezentacja wewnętrzna ----------------------

type instance struct {
	n int
	k int
	D [][]int
	C []int
}

type solution struct {
	T   []int  // cykl (k węzłów)
	sel []bool // czy węzeł jest wybrany
	pos []int  // pozycja węzła w T lub -1
	obj int    // wartość celu
	ins *instance
}

func newSolution(ins *instance) *solution {
	s := &solution{
		T:   make([]int, 0, ins.k),
		sel: make([]bool, ins.n),
		pos: make([]int, ins.n),
		obj: 0,
		ins: ins,
	}
	for i := range s.pos {
		s.pos[i] = -1
	}
	return s
}

func (s *solution) k() int { return len(s.T) }

// pełne przeliczenie celu — do ewentualnych testów spójności
func (s *solution) recompute() int {
	k := s.k()
	if k == 0 {
		return 0
	}
	sum := 0
	for i := 0; i < k; i++ {
		a := s.T[i]
		b := s.T[(i+1)%k]
		sum += s.ins.D[a][b]
	}
	for v := 0; v < s.ins.n; v++ {
		if s.sel[v] {
			sum += s.ins.C[v]
		}
	}
	return sum
}

func prevIdx(i, k int) int {
	if i == 0 {
		return k - 1
	}
	return i - 1
}
func nextIdx(i, k int) int {
	if i == k-1 {
		return 0
	}
	return i + 1
}

// ---------------------- delty ruchów ----------------------

// intra: zamiana dwóch węzłów (pozycji i<j)
func (s *solution) deltaSwap(i, j int) int {
	k := s.k()
	if k < 2 || i == j {
		return 0
	}
	if i > j {
		i, j = j, i
	}
	u := s.T[i]
	v := s.T[j]
	D := s.ins.D
	pi := prevIdx(i, k)
	ni := nextIdx(i, k)
	pj := prevIdx(j, k)
	nj := nextIdx(j, k)

	if ni == j { // sąsiadujące
		a := s.T[pi]
		b := s.T[nj]
		old := D[a][u] + D[u][v] + D[v][b]
		new := D[a][v] + D[v][u] + D[u][b]
		return new - old
	}
	a := s.T[pi]
	b := s.T[ni]
	c := s.T[pj]
	d := s.T[nj]
	old := D[a][u] + D[u][b] + D[c][v] + D[v][d]
	new := D[a][v] + D[v][b] + D[c][u] + D[u][d]
	return new - old
}

func (s *solution) applySwap(i, j int, delta int) {
	if i > j {
		i, j = j, i
	}
	u := s.T[i]
	v := s.T[j]
	s.T[i], s.T[j] = s.T[j], s.T[i]
	s.pos[u], s.pos[v] = j, i
	s.obj += delta
}

// intra: 2-opt (odwrócenie segmentu (i+1..j)); dla poprawności wymagamy k>=4.
func (s *solution) delta2opt(i, j int) int {
	k := s.k()
	if k < 4 || i == j {
		return 0
	}
	if i > j {
		i, j = j, i
	}
	if (i+1)%k == j { // sąsiednie krawędzie – bez sensu
		return 0
	}
	a := s.T[i]
	u := s.T[(i+1)%k]
	v := s.T[j]
	d := s.T[(j+1)%k]
	D := s.ins.D
	old := D[a][u] + D[v][d]
	new := D[a][v] + D[u][d]
	return new - old
}

func (s *solution) apply2opt(i, j int, delta int) {

	if i > j {
		i, j = j, i
	}
	// odwrócenie segmentu (i+1..j)
	for l, r := i+1, j; l < r; l, r = l+1, r-1 {
		s.T[l], s.T[r] = s.T[r], s.T[l]
	}
	// szybsza aktualizacja pos tylko dla dotkniętego zakresu
	for idx := i + 1; idx <= j; idx++ {
		s.pos[s.T[idx]] = idx
	}
	s.obj += delta
}

// inter: wymiana selNode (w turze) z newNode (poza turą), wstawienie w miejsce selNode
// (zostawiamy wersję 1:1 pozycji — zgodną z literalnym "exchange")
func (s *solution) deltaExchange(selNode, newNode int) (int, int) {
	posS := s.pos[selNode]
	k := s.k()
	p := s.T[prevIdx(posS, k)]
	q := s.T[nextIdx(posS, k)]
	D := s.ins.D
	dlen := D[p][newNode] + D[newNode][q] - D[p][selNode] - D[selNode][q]
	dcost := s.ins.C[newNode] - s.ins.C[selNode]
	return dlen + dcost, posS
}

func (s *solution) applyExchange(selNode, newNode int, posS int, delta int) {
	s.T[posS] = newNode
	s.sel[selNode] = false
	s.sel[newNode] = true
	s.pos[newNode] = posS
	s.pos[selNode] = -1
	s.obj += delta
}

// ---------------------- starty ----------------------

// losowy: wybór k węzłów + losowa permutacja
func buildRandomStart(ins *instance, rng *rand.Rand) *solution {
	// przypadki n=1/2: działa bez zmian
	idx := rng.Perm(ins.n)
	T := append([]int(nil), idx[:ins.k]...)
	rng.Shuffle(ins.k, func(i, j int) { T[i], T[j] = T[j], T[i] })
	s := newSolution(ins)
	s.T = T
	for i, v := range T {
		s.sel[v] = true
		s.pos[v] = i
		s.obj += ins.C[v]
	}
	for i := 0; i < ins.k; i++ {
		a := T[i]
		b := T[(i+1)%ins.k]
		s.obj += ins.D[a][b]
	}
	return s
}

// greedy + regret-2: start od zadanego węzła (jeśli n<3 — bezpieczne skróty)
func buildGreedyStart(ins *instance, startNode int) *solution {
	n, k := ins.n, ins.k
	selected := make([]bool, n)
	T := make([]int, 0, k)
	T = append(T, startNode)
	selected[startNode] = true

	if k == 1 {
		// jednowęzłowy cykl
		s := newSolution(ins)
		s.T = T
		s.sel[startNode] = true
		s.pos[startNode] = 0
		s.obj = ins.C[startNode] // brak krawędzi D[a][b], bo k=1 (a==b, ale D[i][i]==0)
		return s
	}

	// drugi: min D + cost — jeśli n==1 to już wyszliśmy; jeśli n==2: wybieramy jedyny pozostały
	bestU := -1
	bestScore := math.MaxInt
	if n >= 2 {
		for u := 0; u < n; u++ {
			if selected[u] {
				continue
			}
			score := ins.D[startNode][u] + ins.C[u]
			if score < bestScore {
				bestScore = score
				bestU = u
			}
		}
	}
	if bestU == -1 { // bezpieczeństwo dla degeneratów
		bestU = startNode
	}
	selected[bestU] = true
	if bestU != startNode {
		T = append(T, bestU)
	}

	// trzeci: najlepsze wstawienie do cyklu 2-węzłowego — tylko gdy mamy sensownie n>=3 i k>=3
	if len(T) == 2 && k >= 3 && n >= 3 {
		uBest := -1
		best := math.MaxInt
		for u := 0; u < n; u++ {
			if selected[u] {
				continue
			}
			a, b := T[0], T[1]
			ins1 := ins.D[a][u] + ins.D[u][b] - ins.D[a][b]
			ins2 := ins.D[b][u] + ins.D[u][a] - ins.D[b][a]
			v := ins1
			if ins2 < v {
				v = ins2
			}
			score := v + ins.C[u]
			if score < best {
				best = score
				uBest = u
			}
		}
		if uBest != -1 {
			selected[uBest] = true
			a, b := T[0], T[1]
			ins1 := ins.D[a][uBest] + ins.D[uBest][b] - ins.D[a][b]
			ins2 := ins.D[b][uBest] + ins.D[uBest][a] - ins.D[b][a]
			if ins1 <= ins2 {
				T = []int{a, uBest, b}
			} else {
				T = []int{b, uBest, a}
			}
		}
	}

	// dobijamy do k: regret-2 (max różnica 2nd-best - best; remis: min(best + cost))
	for len(T) < k {
		type cand struct {
			u, bestPos int
			bestIns    int
			secondIns  int
			regret     int
			tiebreak   int
		}
		C := make([]cand, 0, n-len(T))
		curK := len(T)
		// jeśli curK==1 — rozpatrujemy jedyne „wstawienie” między (T[0],T[0])
		for u := 0; u < n; u++ {
			if selected[u] {
				continue
			}
			best := math.MaxInt
			second := math.MaxInt
			bestPos := -1
			if curK == 1 {
				a := T[0]
				insc := ins.D[a][u] + ins.D[u][a] - ins.D[a][a] // D[a][a]==0
				best, bestPos = insc, 1
				second = insc
			} else {
				for i := 0; i < curK; i++ {
					a := T[i]
					b := T[(i+1)%curK]
					insc := ins.D[a][u] + ins.D[u][b] - ins.D[a][b]
					if insc < best {
						second = best
						best = insc
						bestPos = i + 1
					} else if insc < second {
						second = insc
					}
				}
			}
			C = append(C, cand{
				u: u, bestPos: bestPos, bestIns: best, secondIns: second,
				regret: second - best, tiebreak: best + ins.C[u],
			})
		}
		if len(C) == 0 {
			break
		}
		// wybór: max regret; remis: min tiebreak
		bestIdx := 0
		for i := 1; i < len(C); i++ {
			if C[i].regret > C[bestIdx].regret ||
				(C[i].regret == C[bestIdx].regret && C[i].tiebreak < C[bestIdx].tiebreak) {
				bestIdx = i
			}
		}
		ch := C[bestIdx]
		pos := ch.bestPos % (len(T) + 1)
		T = append(T, 0)
		copy(T[pos+1:], T[pos:])
		T[pos] = ch.u
		selected[ch.u] = true
	}

	s := newSolution(ins)
	s.T = T
	for i, v := range T {
		s.sel[v] = true
		s.pos[v] = i
		s.obj += ins.C[v]
	}
	for i := 0; i < len(T); i++ {
		a := T[i]
		b := T[(i+1)%len(T)]
		s.obj += ins.D[a][b]
	}
	return s
}

// ---------------------- silniki LS ----------------------

// Steepest: pełne sąsiedztwo (intra wybranego typu + inter), wybieramy najlepszą poprawę
func lsSteepest(s *solution, intra IntraType) {
	k := s.k()
	n := s.ins.n
	if k <= 1 {
		return
	}
	for {
		bestDelta := 0
		moveKind := 0 // 0 none, 1 swap, 2 2opt, 3 exchange
		a, b, posS := 0, 0, -1

		// INTRA
		if intra == IntraSwap {
			for i := 0; i < k; i++ {
				for j := i + 1; j < k; j++ {
					d := s.deltaSwap(i, j)
					if d < bestDelta {
						bestDelta = d
						moveKind = 1
						a, b = i, j
					}
				}
			}
		} else if k >= 4 {
			for i := 0; i < k; i++ {
				for j := i + 1; j < k; j++ {
					if (i+1)%k == j {
						continue
					}
					d := s.delta2opt(i, j)
					if d < bestDelta {
						bestDelta = d
						moveKind = 2
						a, b = i, j
					}
				}
			}
		}

		// INTER
		for _, selNode := range s.T {
			for u := 0; u < n; u++ {
				if s.sel[u] {
					continue
				}
				d, p := s.deltaExchange(selNode, u)
				if d < bestDelta {
					bestDelta = d
					moveKind = 3
					a, b, posS = selNode, u, p
				}
			}
		}

		if moveKind == 0 || bestDelta >= 0 {
			return
		}
		switch moveKind {
		case 1:
			s.applySwap(a, b, bestDelta)
		case 2:
			s.apply2opt(a, b, bestDelta)
		case 3:
			s.applyExchange(a, b, posS, bestDelta)
		}
		// kolejne iteracje…
	}
}

// Greedy (first-improvement): typ ruchu i kolejność kandydatów losowane co iterację.
// Optymalizacje: bufory wielokrotnego użytku, brak zbędnych alokacji.
func lsGreedy(s *solution, intra IntraType, rng *rand.Rand, interFirstProb float64) {
	k := s.k()
	n := s.ins.n
	if k <= 1 {
		return
	}

	// Bufory wielokrotnego użycia
	S := make([]int, 0, k)
	U := make([]int, 0, n-k)
	J := make([]int, 0, k) // indeksy pomocnicze dla intra

	nextIsInter := func() bool {
		// jeśli interFirstProb==0.5, rzut monetą
		return rng.Float64() < interFirstProb
	}

	for {
		improved := false

		// Losujemy kolejność typów ruchów (z wagą interFirstProb)
		order := [2]int{0, 1} // 0=intra,1=inter
		if nextIsInter() {
			order[0], order[1] = 1, 0
		}

		for _, typ := range order {
			if typ == 0 {
				// INTRA
				if intra == IntraSwap {
					// Permutacja i — bez alokacji: wypełnij J = 0..k-1, tasuj
					J = J[:k]
					for i := 0; i < k; i++ {
						J[i] = i
					}
					rng.Shuffle(k, func(a, b int) { J[a], J[b] = J[b], J[a] })
					done := false
					for _, i := range J {
						// J2: wartości j>i — zamiast budować slice, iterujemy i tasujemy zakres [i+1..k-1]
						if i+1 >= k {
							continue
						}
						// zbuduj pomocniczą listę [i+1..k-1] w J[:m]
						m := k - (i + 1)
						J = J[:m]
						for t := 0; t < m; t++ {
							J[t] = i + 1 + t
						}
						rng.Shuffle(m, func(a, b int) { J[a], J[b] = J[b], J[a] })
						for idx := 0; idx < m; idx++ {
							j := J[idx]
							d := s.deltaSwap(i, j)
							if d < 0 {
								s.applySwap(i, j, d)
								improved = true
								done = true
								break
							}
						}
						if done {
							break
						}
					}
				} else { // Intra2Opt
					if k >= 4 {
						// Permutacja i
						J = J[:k]
						for i := 0; i < k; i++ {
							J[i] = i
						}
						rng.Shuffle(k, func(a, b int) { J[a], J[b] = J[b], J[a] })
						done := false
						for _, i := range J {
							if i+1 >= k {
								continue
							}
							// zbuduj [i+1..k-1] bez krawędzi sąsiedniej (i+1)%k==j
							m := 0
							J = J[:0]
							for j := i + 1; j < k; j++ {
								if (i+1)%k == j {
									continue
								}
								J = append(J, j)
								m++
							}
							if m == 0 {
								continue
							}
							rng.Shuffle(m, func(a, b int) { J[a], J[b] = J[b], J[a] })
							for _, j := range J {
								d := s.delta2opt(i, j)
								if d < 0 {
									s.apply2opt(i, j, d)
									improved = true
									done = true
									break
								}
							}
							if done {
								break
							}
						}
					}
				}
			} else {
				// INTER
				S = S[:0]
				for _, v := range s.T {
					S = append(S, v)
				}
				U = U[:0]
				for u := 0; u < n; u++ {
					if !s.sel[u] {
						U = append(U, u)
					}
				}
				rng.Shuffle(len(S), func(i, j int) { S[i], S[j] = S[j], S[i] })
				rng.Shuffle(len(U), func(i, j int) { U[i], U[j] = U[j], U[i] })
				done := false
				for _, selNode := range S {
					for _, u := range U {
						d, posS := s.deltaExchange(selNode, u)
						if d < 0 {
							s.applyExchange(selNode, u, posS, d)
							improved = true
							done = true
							break
						}
					}
					if done {
						break
					}
				}
			}
			if improved {
				break
			}
		}

		if !improved {
			return
		}
	}
}

// ---------------------- drobne pomocnicze ----------------------

// ensureStartNodes — przydatne, gdy chcesz jawnie wymusić listę startów (opcjonalnie)
func ensureStartNodes(n, runs int, startNodeIndices []int) ([]int, error) {
	if len(startNodeIndices) == 0 {
		return nil, errors.New("startNodeIndices is empty; provide indices or use StartRandom/StartGreedy with RNG fallback")
	}
	for _, v := range startNodeIndices {
		if v < 0 || v >= n {
			return nil, fmt.Errorf("invalid start node index: %d (n=%d)", v, n)
		}
	}
	if len(startNodeIndices) < runs {
		// dopuszczamy powtórki modulo — nie jest to błąd
	}
	return startNodeIndices, nil
}

// LocalSearch represents the main solver structure
type LocalSearch struct {
	instance *instance
	method   MethodSpec
}

// NewLocalSearch creates a new local search solver
func NewLocalSearch(D [][]int, costs []int, method MethodSpec) (*LocalSearch, error) {
	if err := ValidateMethodSpec(method); err != nil {
		return nil, err
	}
	if err := ValidateInstance(D, costs, method.StrictValidate); err != nil {
		return nil, err
	}

	n := len(D)
	k := (n + 1) / 2
	ins := &instance{n: n, k: k, D: D, C: costs}

	return &LocalSearch{
		instance: ins,
		method:   method,
	}, nil
}

// ValidateMethodSpec validates method configuration
func ValidateMethodSpec(m MethodSpec) error {
	if m.InterFirstProb < 0 || m.InterFirstProb > 1 {
		return &LSError{Msg: "InterFirstProb must be between 0 and 1"}
	}

	switch m.LS {
	case LS_Steepest, LS_Greedy:
	default:
		return &LSError{Msg: "Invalid LS type"}
	}

	switch m.Intra {
	case IntraSwap, Intra2Opt:
	default:
		return &LSError{Msg: "Invalid Intra type"}
	}

	switch m.Start {
	case StartRandom, StartGreedy:
	default:
		return &LSError{Msg: "Invalid Start type"}
	}

	return nil
}

// Solve runs the local search algorithm with given parameters
func (ls *LocalSearch) Solve(startNodes []int, runs int) ([]Solution, error) {
	if runs <= 0 {
		return nil, &LSError{Msg: "runs must be positive"}
	}

	solutions := make([]Solution, 0, runs)
	for r := 0; r < runs; r++ {
		sol, err := ls.solveSingle(startNodes, r)
		if err != nil {
			return nil, err
		}
		solutions = append(solutions, sol)
	}

	return solutions, nil
}

func (ls *LocalSearch) solveSingle(startNodes []int, run int) (Solution, error) {
	// Initialize solution
	var s *solution
	switch ls.method.Start {
	case StartRandom:
		rng := ls.getRNG(run)
		s = buildRandomStart(ls.instance, rng)
	case StartGreedy:
		startNode := ls.getStartNode(startNodes, run)
		s = buildGreedyStart(ls.instance, startNode)
	}

	// Apply local search
	if ls.method.LS == LS_Steepest {
		lsSteepest(s, ls.method.Intra)
	} else {
		rng := ls.getRNG(run + 1000) // Different seed for LS
		lsGreedy(s, ls.method.Intra, rng, ls.method.InterFirstProb)
	}

	return Solution{
		Path:      append([]int(nil), s.T...),
		Objective: s.obj,
	}, nil
}

func (ls *LocalSearch) getRNG(run int) *rand.Rand {
	seed := ls.method.Seed
	if seed == 0 {
		seed = int64(1337 + int(ls.method.LS)*1009 +
			int(ls.method.Intra)*917 +
			int(ls.method.Start)*701)
	}
	return rand.New(rand.NewSource(seed + int64(run)))
}

func (ls *LocalSearch) getStartNode(nodes []int, run int) int {
	if len(nodes) == 0 {
		rng := ls.getRNG(run)
		return rng.Intn(ls.instance.n)
	}
	node := nodes[run%len(nodes)]
	if node < 0 || node >= ls.instance.n {
		rng := ls.getRNG(run)
		return rng.Intn(ls.instance.n)
	}
	return node
}

// LSError represents local search specific errors
type LSError struct {
	Msg string
}

func (e LSError) Error() string {
	return e.Msg
}
