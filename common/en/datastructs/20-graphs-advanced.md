# 20. Graphs (Advanced)

## Description

This topic extends graph algorithms. You learn minimum spanning trees, A*, Floyd–Warshall, strongly connected components, max flow, and bipartite matching. Complete this topic after graph representation, core graph algorithms, heaps, and Union-Find.

Use one term for each concept. An **MST** is a minimum spanning tree. A **heuristic** in A* is a guess of remaining cost. An **SCC** is a strongly connected component. A **flow** is an assignment of load to edges under capacity.

---

## MST: Kruskal, Prim

A **minimum spanning tree** of a connected undirected weighted graph is a tree of n − 1 edges that joins all vertices and has minimum total weight.

**Kruskal** sorts edges by weight and adds an edge when it joins two components. Union-Find tests that condition. Topic 16 already described the steps. Implement Kruskal here if you did not finish it there.

**Prim** grows one tree from a start vertex.

1. Start from s. Dist[s] = 0. Dist[v] = ∞ for others. Dist here means the cheapest edge from the tree to v.
2. Repeat n times: extract the vertex u not in the tree with smallest dist. Add u to the tree. For each neighbor v outside the tree, set dist[v] = min(dist[v], weight(u, v)).
3. A min-heap implements extract-min and decrease-key (or extra inserts).

```text
Prim idea
  tree = empty
  add cheapest edge that leaves the tree
  repeat until n vertices are in the tree
```

```text
Kruskal idea
  forest of vertices
  add cheapest edge that does not make a cycle
```

Both algorithms give the same MST weight if all weights are distinct. If some weights are equal, more than one MST can exist. The weights still match.

Time of Kruskal is O(m log m) from the sort. Time of Prim with a binary heap is O(m log n).

Do not run these MST algorithms on a directed graph as your first MST. This handbook treats MST as undirected.

If the graph is disconnected, you get a **minimum spanning forest**. Prim from one start only fills one component. Kruskal fills all components.

A shortest-path tree from Dijkstra is not an MST. Different objective.

### Questions

#### Theoretical questions

1. What does an MST minimize?
2. How does Kruskal decide to add an edge?
3. How does Prim decide the next vertex?
4. Why can more than one MST exist?
5. Why is a Dijkstra tree not an MST?

#### Easy practical tasks

1. Run both algorithms by hand on a four-vertex graph. Write the edges.
2. Make a table: "Algorithm", "Main structure". Rows: Kruskal, Prim.
3. Write five sentences that define an MST.
4. Draw one skipped Kruskal edge that would form a cycle.

#### Medium practical tasks

1. Implement Kruskal. Print total weight and the edge list.
2. Implement Prim with a binary heap. Match the Kruskal weight on a connected graph.
3. Explain in six sentences what to do when the graph has two components.

#### Advanced practical tasks

1. Generate random connected graphs. Compare weights of Kruskal and Prim on 20 graphs.
2. Time both algorithms on a dense graph and on a sparse graph. Write n, m, and the times.

---

## A*

**A*** finds a shortest path from s to a target t when you have a **heuristic** h(v). h(v) estimates the remaining cost from v to t.

The priority of a vertex v is f(v) = g(v) + h(v). g(v) is the best known cost from s to v (the same role as dist in Dijkstra). The heap extracts the smallest f.

```text
f(v) = g(v) + h(v)
g = cost so far
h = estimate to t
```

**Admissible** heuristic: h(v) is never larger than the true remaining cost. Then A* finds a shortest path (with the usual non-negative weights and a correct implementation).

**Consistent** (monotone) heuristic: h(u) ≤ w(u, v) + h(v). Then A* behaves like Dijkstra with a reduced cost. You can mark v closed on first extract in the common form.

A grid with 4-way moves often uses Manhattan distance as h. A grid with 8-way moves often uses a matching distance. Euclidean distance is common in the plane.

If h(v) = 0 for every v, A* is Dijkstra.

A* is not BFS. BFS ignores weights and ignores h.

A bad heuristic (too large) can return a path that is not shortest. A heuristic that is always 0 is safe but not faster.

Do not use A* when you need distances from s to all vertices. Use Dijkstra. A* is for one target (or a small set) and a useful h.

Reconstruct the path with a parent array as in Dijkstra.

### Questions

#### Theoretical questions

1. What is f(v)?
2. What does admissible mean?
3. When is A* the same as Dijkstra?
4. Why can a too-large h be unsafe?
5. When do you prefer Dijkstra to A*?

#### Easy practical tasks

1. On a 3 × 3 grid, write Manhattan h to the opposite corner for each cell.
2. Make a table: "Algorithm", "Uses h?", "One target?". Rows: BFS, Dijkstra, A*.
3. Write five sentences that define A*.
4. For h = 0, write f in terms of g.

#### Medium practical tasks

1. Implement A* on a grid with obstacles. Use Manhattan h. Compare the path cost with BFS if every step costs 1.
2. Implement Dijkstra on the same grid. Count heap extracts for both. Write the two counts.
3. Explain in six sentences why Manhattan is admissible for 4-way unit cost.

#### Advanced practical tasks

1. Use a non-admissible h on purpose (2 × Manhattan). Show a case where the path cost is worse than optimal, or show that you cannot find one on your map and explain why.
2. Write a one-page note: consistent versus admissible. Cite one source.

---

## Floyd–Warshall

**Floyd–Warshall** computes shortest paths between **all pairs** of vertices. It uses an n × n matrix.

Idea: allow vertices 0 .. k as intermediate vertices, then k+1.

```text
for k in 0 .. n-1:
  for i in 0 .. n-1:
    for j in 0 .. n-1:
      if dist[i][k] + dist[k][j] < dist[i][j]:
          dist[i][j] = dist[i][k] + dist[k][j]
```

Init: dist[i][i] = 0. dist[i][j] = weight(i, j) if an edge exists. dist[i][j] = ∞ if i ≠ j and no edge.

Time is Θ(n³). Space is Θ(n²). Use Floyd–Warshall when n is small (hundreds, not millions) and you need all pairs.

Negative weights are allowed. A negative cycle exists if some dist[i][i] becomes negative.

You can store next[i][j] to reconstruct a path.

Do not use Floyd–Warshall on a huge sparse graph if you only need one source. Use Dijkstra n times (non-negative) or Johnson's algorithm (later reading).

Floyd–Warshall is the matrix method. It fits an adjacency matrix. It does not use a heap.

The three loops must use k as the outer loop. A wrong order is a common error.

### Questions

#### Theoretical questions

1. What does Floyd–Warshall output?
2. What does the k loop mean?
3. What is the time class?
4. How do you detect a negative cycle?
5. Why must k be the outer loop?

#### Easy practical tasks

1. Write the init matrix for a three-vertex path with weights 2 and 3.
2. Make a table: "Need", "Algorithm". Rows: one source non-negative, all pairs small n.
3. Write five sentences that define Floyd–Warshall.
4. After one k step, write which paths may change.

#### Medium practical tasks

1. Implement Floyd–Warshall. Test a small graph against n Dijkstra runs (non-negative).
2. Add a negative cycle. Show a negative diagonal.
3. Reconstruct one path with a next matrix.

#### Advanced practical tasks

1. Time Floyd–Warshall for n = 200 and n = 400 if you can. Write the two times and the n³ ratio.
2. Write a one-page note: Floyd–Warshall versus n times Dijkstra versus Johnson (awareness). Cite a source for Johnson.

---

## SCC (Kosaraju / Tarjan)

A **strongly connected component (SCC)** of a directed graph is a maximal set of vertices such that every vertex can reach every other vertex in the set along directed paths.

The graph of SCCs (each component contracted to one vertex) is a DAG.

**Kosaraju** (two DFS passes):

1. Run DFS on G. Record finish order (or push vertices on a stack at finish).
2. Build the transpose graph Gᵀ (reverse every edge).
3. Run DFS on Gᵀ. Start vertices in decreasing finish time of step 1. Each DFS tree is one SCC.

**Tarjan** (one DFS pass) uses a stack and low-link values. It finds SCCs during one walk. The code is shorter in a contest setting after you learn it. The idea is denser.

```text
G:   0 → 1 → 2
     ↑________↓     (if 2 → 0)
one SCC {0,1,2}

if you remove 2 → 0, three SCCs
```

Time of both methods is Θ(n + m).

Do not use undirected connected components as SCCs. Direction matters.

Condensation (the SCC DAG) is useful for "can A reach B?" after you check intra-component reachability, and for 2-SAT (later).

Implement Kosaraju first if you already have DFS finish times. Implement Tarjan when you want one pass.

Isolated vertices are SCCs of size 1. A single directed edge a → b without a return path gives two SCCs.

### Questions

#### Theoretical questions

1. What can each vertex in an SCC do to the others?
2. What shape is the condensation of SCCs?
3. What are the three Kosaraju steps?
4. Why do you reverse the edges?
5. Why is an undirected component not an SCC?

#### Easy practical tasks

1. Draw a directed cycle of 3 and one extra vertex that points into the cycle. Write the SCCs.
2. Make a table: "Method", "DFS passes". Rows: Kosaraju, Tarjan.
3. Write five sentences that define an SCC.
4. Reverse the edges of a three-edge drawing. Write Gᵀ.

#### Medium practical tasks

1. Implement Kosaraju. Print component ids. Test a cycle and a DAG (n SCCs).
2. Build the condensation DAG: an edge from component A to B if G has an edge from A to B.
3. Explain in six sentences the role of finish times.

#### Advanced practical tasks

1. Implement Tarjan on the same tests. Component partitions must match Kosaraju.
2. Read one 2-SAT reduction that uses SCCs (high level). Write ten STE sentences. Do not write a full solver unless you want a project.

---

## Max flow (Ford–Fulkerson / Dinic awareness)

A **flow network** is a directed graph with a **source** s, a **sink** t, and a **capacity** c(u, v) ≥ 0 on each edge. A **flow** assigns a number f(u, v) such that 0 ≤ f(u, v) ≤ c(u, v) and flow is conserved at every vertex except s and t. The **value** of the flow is the net flow out of s.

A **maximum flow** has maximum value.

**Residual capacity** of an edge is c − f. A **residual graph** also has a backward edge with capacity f. An **augmenting path** is an s–t path in the residual graph. You can push extra flow along that path.

**Ford–Fulkerson** repeats: find any augmenting path, push the bottleneck residual capacity, update residuals. If you find paths with BFS (Edmonds–Karp), each search is shortest in number of edges, and the algorithm is polynomial.

**Dinic** uses blocking flows and level graphs. It is faster on many graphs. This name is **awareness** for the first pass. Implement Edmonds–Karp (Ford–Fulkerson + BFS) before Dinic.

```text
you implement first:  Edmonds–Karp (BFS augmenting paths)
you only name:        Dinic, blocking flow
```

The **min cut** is a partition (S, T) with s in S and t in T. The cut capacity is the sum of capacities from S to T. The max-flow min-cut theorem says max flow value equals min cut capacity.

Do not ignore backward residual edges. They let you "undo" a bad push.

Capacities must be handled with care if they are real numbers. This handbook uses integers.

### Questions

#### Theoretical questions

1. What does conservation of flow mean at a vertex that is not s or t?
2. What is residual capacity?
3. What does one Ford–Fulkerson step do?
4. What is Edmonds–Karp?
5. What does the max-flow min-cut theorem state?

#### Easy practical tasks

1. Draw s → a → t with capacities 3 and 2. Write the max flow.
2. Make a table: "Term", "Meaning". Rows: capacity, flow, residual.
3. Write five sentences that define an augmenting path.
4. Draw a backward residual after you push 1 on an edge of capacity 3.

#### Medium practical tasks

1. Implement Edmonds–Karp. Test a small network. Print the max flow value.
2. After the run, mark vertices reachable from s in the residual graph. That set is S of a min cut.
3. Explain in six sentences why a backward edge exists.

#### Advanced practical tasks

1. Compare your max flow with a min cut capacity on two networks. They must match.
2. Read a Dinic outline. Write ten STE sentences. Do not implement Dinic until Edmonds–Karp is correct.

---

## Bipartite matching (awareness)

A **bipartite graph** has two sides L and R. Every edge joins a vertex in L to a vertex in R. There is no edge inside L or inside R.

A **matching** is a set of edges that do not share a vertex. A **maximum matching** has the maximum number of edges.

You can compute a maximum bipartite matching as a max flow: add s to all of L, add all of R to t, set every capacity to 1, compute max flow. The flow value is the matching size. Edges from L to R with flow 1 are the matching.

**Hopcroft–Karp** is a faster bipartite matching algorithm. Treat it as **awareness** until flow is clear.

**Hall's marriage condition** is a theorem about when a matching covers L. You do not need to implement a Hall checker.

```text
L: a b c
R: 1 2 3
edges: a-1, a-2, b-2, c-3
one maximum matching size 3 or less; find one by hand
```

Uses: assignment of jobs to workers, pairing students to slots, some scheduling models.

This section is **awareness plus a flow reduction**. Implement matching through Edmonds–Karp on a tiny bipartite graph. Do not start with Hopcroft–Karp.

A general-graph matching (non-bipartite) is harder (Blossom). That algorithm is out of scope.

### Questions

#### Theoretical questions

1. What is a bipartite graph?
2. What is a matching?
3. How do you reduce bipartite matching to max flow?
4. Why are capacities 1 in that reduction?
5. Why is general matching out of scope?

#### Easy practical tasks

1. Draw a bipartite graph of 3 + 3 vertices. Mark a matching of size 2.
2. Make a table: "Algorithm", "You implement now?". Rows: Edmonds–Karp matching, Hopcroft–Karp, Blossom.
3. Write five sentences that define bipartite matching.
4. Add s and t to your drawing. Write the new edges.

#### Medium practical tasks

1. Build the flow network for a small bipartite graph. Run your max-flow code. Read the matching from L–R flow.
2. Find one unmatched vertex. Write whether a larger matching exists by hand.
3. Explain in six sentences one assignment story (jobs and workers).

#### Advanced practical tasks

1. Write a one-page note: flow reduction, min cut meaning on the two sides, Hopcroft–Karp as later reading. Cite one source.
2. After Edmonds–Karp works, list three tests (empty graph, perfect matching, a bottleneck vertex).

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe how MST, shortest path, all-pairs, SCC, and flow ask different questions on a graph.
2. How do heaps, Union-Find, DFS finish times, and residual BFS each serve one algorithm in this topic?
3. A teammate says A* "is BFS with a score". Which facts do you use to correct that sentence?
4. What stays undirected (MST) versus directed (SCC, flow)?
5. Why does this handbook implement Edmonds–Karp before Dinic and Hopcroft–Karp?

#### Easy practical tasks

1. Write a one-page cheat sheet: Kruskal, Prim, A*, Floyd–Warshall, Kosaraju, flow, bipartite reduction.
2. Make a decision table: "Question", "Algorithm". Six rows.
3. Draw one MST, one SCC condensation, and one tiny flow network.
4. Write f = g + h and max-flow = min-cut from memory.

#### Medium practical tasks

1. Implement Kruskal, Prim, and Floyd–Warshall on one graph library. Write which graphs each function accepts.
2. Implement Kosaraju and Edmonds–Karp. Tests: a directed cycle, a small flow.
3. Implement A* on a grid and compare cost with Dijkstra.

#### Advanced practical tasks

1. Build a tool that reads a graph, a problem name, and prints MST weight, SCCs, or max flow.
2. Read topic 21. Write five sentences on why a residual graph is extra memory during flow, not a persistent version of the input graph.
