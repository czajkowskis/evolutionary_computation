# Lab 10: Variable Neighborhood Search (VNS) for Travelling Salesperson Problem with Node Costs

**Implementation of Variable Neighborhood Search algorithm for solving the Travelling Salesperson Problem with node costs and Euclidean distances.**

---

## Problem Overview

Given a set of nodes in a plane, each defined by:
- `(x, y)` coordinates,
- a node cost,

**the goal** is to select exactly 50% of the nodes (rounded up if odd) and form a Hamiltonian cycle through them, minimizing the sum of:
- the total Euclidean path length (rounded to integers),
- the total cost of the selected nodes.

---

## Variable Neighborhood Search (VNS) Method

VNS is a metaheuristic that systematically explores multiple neighborhood structures to escape local optima. The algorithm alternates between:

1. **Shaking**: Apply a neighborhood operator to escape local optimum
2. **Local Search**: Intensify the search using steepest descent
3. **Neighborhood Change**: Move to next neighborhood if no improvement

### Neighborhood Structures (Shaking Operators)

1. **N1: Node Exchange** - Exchange selected nodes with non-selected nodes
2. **N2: Random 2-opt** - Apply random 2-opt moves
3. **N3: Destroy-Repair** - Remove 20-30% of nodes and rebuild with greedy nearest neighbor
4. **N4: Double-Bridge** - Apply 4-opt variant (double-bridge move)

### Neighborhood Change Strategies

- **Sequential**: Systematically try neighborhoods 1→2→3→4→5→1...
- **Random**: Randomly select neighborhood
- **Adaptive**: Weight neighborhoods based on recent performance

### Configuration Options

**Stopping Conditions** (at least one should be set):
- **Time Limit**: Maximum execution time per run (0 = no time limit)
- **Max Iterations**: Maximum number of VNS iterations (0 = no limit)
- **Max Iterations No Improve**: Maximum iterations without improvement (0 = no limit)

**Other Options**:
- **Shaking Intensity**: Number of moves in shaking phase (default: 3)
- **Local Search**: Enable/disable local search after shaking
- **Neighborhood Change**: Strategy for selecting neighborhoods (sequential, random, adaptive)

---

## Implemented Configurations

1. **VNS_Sequential_LS** - Sequential neighborhood change with local search
2. **VNS_Adaptive_LS** - Adaptive neighborhood selection with local search
3. **VNS_Random_LS** - Random neighborhood selection with local search

For each configuration, **20 runs** are executed per instance.

---

## Advantages Over Previous Methods

1. **Systematic Exploration**: VNS systematically tries different neighborhoods instead of random restarts
2. **Escapes Local Optima**: Shaking operators provide controlled diversification
3. **Combines Strengths**: Uses best local search from previous labs (steepest descent with 2-opt)
4. **Adaptive Learning**: Can learn which neighborhoods work best for the problem
5. **Flexible Framework**: Easy to add new neighborhood structures
6. **Flexible Stopping**: Multiple stopping conditions available (time, iterations, or iterations without improvement)

---

## Running the Program

```bash
cd variable_neighborhood_search
go run cmd/main.go
```

The program will:
- Process both instances (A and B)
- Run all VNS configurations
- Generate CSV results in `output/results/`
- Generate visualization plots in `output/plots/`
- Print statistics to console

---

## Output Format

Results are saved in CSV format with the following columns:
- `instance`: Instance name (A or B)
- `method`: VNS configuration name
- `avg_objective`: Average objective value
- `av(min,max)`: Average with min and max
- `min_objective`: Minimum objective value
- `max_objective`: Maximum objective value
- `avg_time_ms`: Average execution time in milliseconds
- `avg_iterations`: Average number of VNS iterations
- `best_objective`: Best objective value found
- `best_path`: Best solution path

---

## Validation

All best solutions are checked using the provided solution checker.

---

## References

- Mladenović, N., & Hansen, P. (1997). Variable neighborhood search. *Computers & operations research*, 24(11), 1097-1100.
- Hansen, P., & Mladenović, N. (2001). Variable neighborhood search: Principles and applications. *European journal of operational research*, 130(3), 449-467.

