
# Lab 01: Greedy Heuristics for Travelling Salesperson Problem with Node Costs 

**Implementation of greedy heuristics for solving the Travelling Salesperson Problem with node costs and Euclidean distances.**

---

## Problem Overview

Given a set of nodes in a plane, each defined by:
- `(x, y)` coordinates,
- a node cost,

**the goal** is to select exactly 50% of the nodes (rounded up if odd) and form a Hamiltonian cycle through them, minimizing the sum of:
- the total Euclidean path length (rounded to integers),
- the total cost of the selected nodes.

---

## Implemented Methods

1. **Random Solution**
2. **Nearest Neighbor (End only)**
3. **Nearest Neighbor (Any position)**
4. **Greedy Cycle**

For each method, **200 solutions** are generated starting from each node. Additionally, **200 random solutions** are generated.

---

## Validation
All best solutions are checked using the provided solution checker.