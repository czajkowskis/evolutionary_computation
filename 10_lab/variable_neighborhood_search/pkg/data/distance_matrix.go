package data

import (
	"log"
	"math"
)

func CalculateDistanceMatrix(nodes []Node) [][]int {
	n := len(nodes)
	distanceMatrix := make([][]int, n)
	for i := range distanceMatrix {
		distanceMatrix[i] = make([]int, n)
	}

	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if i != j {
				x1, y1 := nodes[i].X, nodes[i].Y
				x2, y2 := nodes[j].X, nodes[j].Y
				distance := math.Sqrt(math.Pow(float64(x2-x1), 2) + math.Pow(float64(y2-y1), 2))
				distanceMatrix[i][j] = int(math.Round(distance))
			}
		}
	}
	log.Printf("Calculated distance matrix for %d nodes", n)
	return distanceMatrix
}

