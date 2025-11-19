# Lab 06: Multiple start local search (MSLS) and iterated local search (ILS)

**Implementation of two simple extensions of local search for solving the Travelling Salesperson Problem with node costs and Euclidean distances.**

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

1. **Multiple start local search (MSLS) – we will use steepest local search starting from random solutions.**
2. **Iterated local search (ILS).**

Each of the methods (MSLS and ILS) is run 20 times for each instance. In MSLS we perform 200 iterations of basic local search. For ILS as the stopping condition we use the average running time of MSLS. For ILS as the starting solution (one for each run of ILS) we use random solution. The results of a single run of MSLS is the best solution among a given number of runs of local search.

---

## Validation
All best solutions are checked using the provided solution checker.
