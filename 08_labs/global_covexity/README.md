# Lab 03: Local Search for Travelling Salesperson Problem with Node Costs 

**Implementation of local search algorithms for solving the Travelling Salesperson Problem with node costs and Euclidean distances.**

---

## Problem Overview

Given a set of nodes in a plane, each defined by:
- `(x, y)` coordinates,
- a node cost,

**the goal** is to select exactly 50% of the nodes (rounded up if odd) and form a Hamiltonian cycle through them, minimizing the sum of:
- the total Euclidean path length (rounded to integers),
- the total cost of the selected nodes.

---

## Implemented Local Search combinations (LS type, intra-moves, start type)

1. **Steepest, Two-nodes, Random**
2. **Steepest, Two-nodes, Greedy**
3. **Steepest, Two-edges, Random**
4. **Steepest, Two-edges, Greedy**
5. **Greedy, Two-nodes, Random**
6. **Greedy, Two-nodes, Greedy**
7. **Greedy, Two-edges, Random**
8. **Greedy, Two-edges, Greedy**

For each method, **200 solutions** are generated.

---

## Validation
All best solutions are checked using the provided solution checker.
