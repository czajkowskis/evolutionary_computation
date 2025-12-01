package utils

type Row struct {
	Name        string
	AvgV        float64
	MinV        int
	MaxV        int
	AvgTms      float64
	AvgLNSIters float64
	BestPath    []int
	BestValue   int
}
