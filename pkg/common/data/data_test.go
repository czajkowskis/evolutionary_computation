package data

import (
	"os"
	"testing"
)

func TestReadNodes(t *testing.T) {
	// Create a temporary CSV file for testing
	testCSV := `10;20;5
15;25;3
5;10;7
`
	tmpFile, err := os.CreateTemp("", "test_nodes_*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(testCSV); err != nil {
		t.Fatalf("Failed to write test data: %v", err)
	}
	tmpFile.Close()

	// Test reading nodes
	nodes, err := ReadNodes(tmpFile.Name())
	if err != nil {
		t.Fatalf("ReadNodes failed: %v", err)
	}

	if len(nodes) != 3 {
		t.Errorf("Expected 3 nodes, got %d", len(nodes))
	}

	// Check first node
	if nodes[0].X != 10 || nodes[0].Y != 20 || nodes[0].Cost != 5 {
		t.Errorf("First node incorrect: got (%d, %d, %d), expected (10, 20, 5)",
			nodes[0].X, nodes[0].Y, nodes[0].Cost)
	}
}

func TestCalculateDistanceMatrix(t *testing.T) {
	nodes := []Node{
		{X: 0, Y: 0, Cost: 1},
		{X: 3, Y: 4, Cost: 2}, // Distance should be 5 (3-4-5 triangle)
		{X: 0, Y: 0, Cost: 3}, // Same as first, distance 0
	}

	matrix := CalculateDistanceMatrix(nodes)

	if len(matrix) != 3 {
		t.Errorf("Expected 3x3 matrix, got %dx%d", len(matrix), len(matrix[0]))
	}

	// Check distance from node 0 to node 1 (should be 5)
	if matrix[0][1] != 5 {
		t.Errorf("Distance from node 0 to node 1: got %d, expected 5", matrix[0][1])
	}

	// Check distance from node 0 to node 2 (should be 0, same coordinates)
	if matrix[0][2] != 0 {
		t.Errorf("Distance from node 0 to node 2: got %d, expected 0", matrix[0][2])
	}

	// Check diagonal (should be 0)
	if matrix[0][0] != 0 || matrix[1][1] != 0 || matrix[2][2] != 0 {
		t.Error("Diagonal elements should be 0")
	}
}

