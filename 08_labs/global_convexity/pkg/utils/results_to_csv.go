package utils

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const outputDir = "output/results"

func intsToDashString(nums []int) string {
	if len(nums) == 0 {
		return ""
	}
	sb := strings.Builder{}
	sb.WriteString("[")
	for i, n := range nums {
		if i > 0 {
			sb.WriteString(" ")
		}
		sb.WriteString(strconv.Itoa(n))
	}
	sb.WriteString("]")
	return sb.String()
}

func WriteResultsCSV(instanceName string, rows []Row) error {

	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return fmt.Errorf("make dir %s: %w", outputDir, err)
	}

	filename := filepath.Join(
		outputDir,
		fmt.Sprintf("results_instance_%s.csv", instanceName),
	)

	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("create csv: %w", err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	if err := w.Write([]string{
		"instance",
		"method",
		"avg_objective",
		"av(min,max)",
		"min_objective",
		"max_objective",
		"avg_time_ms",
		"best_objective",
		"best_path",
	}); err != nil {
		return fmt.Errorf("write header: %w", err)
	}

	for _, r := range rows {
		avg4 := fmt.Sprintf("%.4f", r.AvgV)
		avgSummary := fmt.Sprintf("%.4f (%d, %d)", r.AvgV, r.MinV, r.MaxV)

		rec := []string{
			instanceName,
			r.Name,
			avg4,       // avg_objective
			avgSummary, // av(min,max)
			strconv.Itoa(r.MinV),
			strconv.Itoa(r.MaxV),
			fmt.Sprintf("%.2f", r.AvgTms),
			strconv.Itoa(r.BestValue),
			intsToDashString(r.BestPath),
		}
		if err := w.Write(rec); err != nil {
			return fmt.Errorf("write row: %w", err)
		}
	}

	log.Printf("CSV saved: %s", filename)
	return nil
}
