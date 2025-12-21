package utils

import (
	"testing"

	"github.com/czajkowskis/evolutionary_computation/pkg/common/algorithms"
)

func TestCalculateStatistics(t *testing.T) {
	tests := []struct {
		name     string
		solutions []algorithms.Solution
		wantMin  int
		wantMax  int
		wantAvg  float64
	}{
		{
			name: "empty solutions",
			solutions: []algorithms.Solution{},
			wantMin: 0,
			wantMax: 0,
			wantAvg: 0,
		},
		{
			name: "single solution",
			solutions: []algorithms.Solution{
				{Path: []int{0, 1, 2}, Objective: 100},
			},
			wantMin: 100,
			wantMax: 100,
			wantAvg: 100.0,
		},
		{
			name: "multiple solutions",
			solutions: []algorithms.Solution{
				{Path: []int{0, 1, 2}, Objective: 100},
				{Path: []int{0, 2, 1}, Objective: 150},
				{Path: []int{1, 0, 2}, Objective: 120},
			},
			wantMin: 100,
			wantMax: 150,
			wantAvg: 123.33333333333333,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMin, gotMax, gotAvg := CalculateStatistics(tt.solutions)
			if gotMin != tt.wantMin {
				t.Errorf("CalculateStatistics() min = %v, want %v", gotMin, tt.wantMin)
			}
			if gotMax != tt.wantMax {
				t.Errorf("CalculateStatistics() max = %v, want %v", gotMax, tt.wantMax)
			}
			if gotAvg != tt.wantAvg {
				t.Errorf("CalculateStatistics() avg = %v, want %v", gotAvg, tt.wantAvg)
			}
		})
	}
}

func TestGenerateStartNodeIndices(t *testing.T) {
	tests := []struct {
		name string
		n    int
		want []int
	}{
		{
			name: "zero nodes",
			n:    0,
			want: []int{},
		},
		{
			name: "one node",
			n:    1,
			want: []int{0},
		},
		{
			name: "five nodes",
			n:    5,
			want: []int{0, 1, 2, 3, 4},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateStartNodeIndices(tt.n)
			if len(got) != len(tt.want) {
				t.Errorf("GenerateStartNodeIndices() length = %v, want %v", len(got), len(tt.want))
			}
			for i, v := range got {
				if v != tt.want[i] {
					t.Errorf("GenerateStartNodeIndices()[%d] = %v, want %v", i, v, tt.want[i])
				}
			}
		})
	}
}

func TestSanitizeFileName(t *testing.T) {
	tests := []struct {
		name string
		input string
		want  string
	}{
		{
			name:  "simple name",
			input: "test_file",
			want:  "test_file",
		},
		{
			name:  "with spaces",
			input: "test file name",
			want:  "test_file_name",
		},
		{
			name:  "with parentheses",
			input: "test(file)",
			want:  "testfile",
		},
		{
			name:  "with commas",
			input: "test,file,name",
			want:  "testfilename",
		},
		{
			name:  "complex name",
			input: "Best (2-opt) Solution, Instance A",
			want:  "Best_2-opt_Solution_Instance_A",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SanitizeFileName(tt.input)
			if got != tt.want {
				t.Errorf("SanitizeFileName() = %v, want %v", got, tt.want)
			}
		})
	}
}

