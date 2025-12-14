# Lab 09: Hybrid evolutionary algorithm

**Implementation of hybrid evolutionary algorithms, with elite population of 20, steady state algorithm, parents selected from the population with the uniform probability, local search applied to each offspring, and using two different recombination operators - for solving the Travelling Salesperson Problem with node costs and Euclidean distances.**

---

## Problem Overview

Given a set of nodes in a plane, each defined by:
- `(x, y)` coordinates,
- a node cost,

**the goal** is to select exactly 50% of the nodes (rounded up if odd) and form a Hamiltonian cycle through them, minimizing the sum of:
- the total Euclidean path length (rounded to integers),
- the total cost of the selected nodes.

---

## Validation
All best solutions are checked using the provided solution checker.
