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


## Input Format

Input files should be in the following format (one row per node):

<custom-element data-json="%7B%22type%22%3A%22table-metadata%22%2C%22attributes%22%3A%7B%22title%22%3A%22Input%20Format%22%7D%7D" />

| x  | y  | cost |
|----|----|------|
| 10 | 20 | 5    |
| 15 | 25 | 3    |
| 5  | 10 | 7    |

---
