package data

import (
	"encoding/csv"
	"os"
	"strconv"
)

func ReadNodes(filename string) ([]Node, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var nodes []Node
	for _, record := range records {
		x, err := strconv.Atoi(record[0])
		if err != nil {
			return nil, err
		}
		y, err := strconv.Atoi(record[1])
		if err != nil {
			return nil, err
		}
		cost, err := strconv.Atoi(record[2])
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, Node{x, y, cost})
	}
	return nodes, nil
}
