package utils

// Row represents a single row of aggregated experiment results.
// Fields are optional - use only what's needed for each lab.
type Row struct {
	Name        string
	AvgV        float64
	MinV        int
	MaxV        int
	AvgTms      float64
	AvgLNSIters float64 // Optional: used by LNS labs
	BestPath    []int
	BestValue   int
}

