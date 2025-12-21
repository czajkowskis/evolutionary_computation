# Evolutionary Computation - Travelling Salesperson Problem with Node Costs Optimization

Repository for algorithms solving the Traveling Salesperson problem with node costs.

## Problem Overview

Given a set of nodes in a plane, each defined by:
- `(x, y)` coordinates,
- a node cost,

**the goal** is to select exactly 50% of the nodes (rounded up if odd) and form a Hamiltonian cycle through them, minimizing the sum of:
- the total path length (rounded to integers),
- the total cost of the selected nodes.

The distance matrix is precomputed and used as the sole input for optimization methods.

## Project Structure

```
labs/
├── pkg/
│   └── common/              # Shared packages used across all labs
│       ├── algorithms/       # Common algorithm types (Solution, FindBestSolution)
│       ├── config/           # Configuration constants (plot bounds, instance paths)
│       ├── data/             # Data structures and I/O (Node, ReadNodes, CalculateDistanceMatrix)
│       ├── utils/            # Utility functions (stats, CSV writing, file sanitization)
│       └── visualisation/    # Plotting functionality
├── 01_labs/                  # Lab 01: Greedy Heuristics
├── 02_labs/                  # Lab 02: Greedy Regret Heuristics
├── 03_labs/                  # Lab 03: Local Search
├── 04_labs/                  # Lab 04: Local Search with Candidate Moves
├── 05_labs/                  # Lab 05: Local Search with Deltas
├── 06_labs/                  # Lab 06: Local Search Extensions (MSLS, ILS)
├── 07_labs/                  # Lab 07: Large Neighborhood Search
├── 08_labs/                  # Lab 08: Global Convexity
├── 09_labs/                  # Lab 09: Hybrid Evolutionary Algorithm
├── 10_labs/                  # Lab 10: Variable Neighborhood Search
├── instances/                 # Shared TSP instance files (TSPA.csv, TSPB.csv)
├── output/                   # Output directory for plots and results
├── go.mod                    # Single Go module for the entire project
└── Makefile                  # Build and run automation

```

Each lab directory contains:
- `cmd/main.go` - Main entry point for the lab
- `pkg/algorithms/` - Lab-specific algorithm implementations
- `README.md` - Lab-specific documentation

## Building and Running

### Prerequisites
- Go 1.25.0 or later
- Dependencies will be installed automatically via `go mod tidy`

### Quick Start

1. **Install dependencies:**
   ```bash
   make deps
   # or
   go mod tidy
   ```

2. **Run a specific lab:**
   ```bash
   make run-lab-01    # Run Lab 01: Greedy Heuristics
   make run-lab-02    # Run Lab 02: Greedy Regret Heuristics
   # ... and so on for labs 03-10
   ```

3. **Build all labs:**
   ```bash
   make build-all
   ```

4. **Run tests:**
   ```bash
   make test
   ```

5. **Clean build artifacts:**
   ```bash
   make clean
   ```

### Running Labs Directly

You can also run labs directly using Go:
```bash
cd 01_labs/greedy_heuristics/cmd && go run main.go
```

## Input Format

Input files should be in the following format (one row per node, semicolon-separated):

| x  | y  | cost |
|----|----|------|
| 10 | 20 | 5    |
| 15 | 25 | 3    |
| 5  | 10 | 7    |

Instance files are located in the `instances/` directory:
- `instances/TSPA.csv` - Instance A
- `instances/TSPB.csv` - Instance B

## Output

Each lab generates:
- **Plots**: Saved to `output/{lab_name}/plots/`
- **Results**: CSV files saved to `output/{lab_name}/results/`

## Adding a New Lab

1. Create a new directory: `{NN}_labs/{lab_name}/`
2. Create `cmd/main.go` with your lab's main function
3. Create `pkg/algorithms/` for lab-specific algorithms
4. Use common packages from `pkg/common/`:
   - `pkg/common/data` for Node, ReadNodes, CalculateDistanceMatrix
   - `pkg/common/algorithms` for Solution type and FindBestSolution
   - `pkg/common/utils` for statistics and CSV writing
   - `pkg/common/visualisation` for plotting
   - `pkg/common/config` for configuration constants

5. Update the Makefile to add build/run targets for your lab

## Testing

Run tests for common packages:
```bash
go test ./pkg/common/... -v
```

## Code Quality

- Format code: `make fmt` or `go fmt ./...`
- The project uses a single Go module for easier dependency management
- Common code is shared via `pkg/common/` to avoid duplication

---
