# Lab 04: Local Search with Candidate Moves for Travelling Salesperson Problem with Node Costs 

**Implementation of the steepest local search algorithm with the use of candidate moves using the two-edges exchange intra-route moves neighborhood with random start for solving the Travelling Salesperson Problem with node costs and Euclidean distances.**

---

## Problem Overview

Given a set of nodes in a plane, each defined by:
- `(x, y)` coordinates,
- a node cost,

**the goal** is to select exactly 50% of the nodes (rounded up if odd) and form a Hamiltonian cycle through them, minimizing the sum of:
- the total Euclidean path length (rounded to integers),
- the total cost of the selected nodes.

---

## Implemented Method:

1. **Steepest Local Search with candidate moves using the two-edges exchange intra-route moves neighborhood with random start**

For this method, **200 solutions** are generated and compared with the same method without candidate moves.

---

## Validation
All best solutions are checked using the provided solution checker.
