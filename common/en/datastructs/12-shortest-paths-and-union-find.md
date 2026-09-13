# 12. Shortest Paths and Union-Find

## Description

This topic explains weighted shortest paths and a disjoint-set structure. You learn Dijkstra, Bellman-Ford, DAG relaxation, the Union-Find parent array, union by rank, path compression, Kruskal as a client of Union-Find, and the idea of Kruskal and Prim for a minimum spanning tree.

Complete this topic after graphs, heaps, and arrays.

Use one term for each concept. **Relaxation** updates a tentative distance when a better path appears. A **non-negative** weight is greater than or equal to 0. A **negative cycle** is a directed cycle whose weights sum to a negative number. **Union-Find** stores a partition of n elements. **Find** returns the representative of an element. **Union** merges two sets. An **MST** is a minimum spanning tree.

---

## Dijkstra and Bellman-Ford

**Relaxation** of edge u → v with weight w:

```text
if dist[v] > dist[u] + w:
  dist[v] = dist[u] + w
  parent[v] = u
```

Both algorithms use relaxation. They differ in the order of edges and in the weight rules.

**Dijkstra** computes distances from s when **all weights are non-negative**.

Idea:

1. Set dist[s] = 0 and dist[v] = ∞ for v ≠ s.
2. Keep a set of settled vertices (the exact method uses a min priority queue of unsettled vertices).
3. Extract the unsettled vertex u with the smallest dist[u]. Settle u.
4. Relax every outgoing edge of u.

```text
s=A
A --1--> B --3--> C
A --4--> C

dist: A=0, B=1, C=min(4, 1+3)=4
```

With a binary heap, time is O((n + m) log n) on a list. Each vertex is extracted one time if you do not re-insert. Some implementations allow more than one queue entry per vertex.

Dijkstra is wrong when a negative edge can improve a path after you settle a vertex. Do not use Dijkstra on graphs that can have negative weights.

**Bellman-Ford** allows negative weights. It forbids a **negative cycle** that is reachable from s if you want finite distances.

Idea:

1. Set dist[s] = 0 and dist[v] = ∞ for v ≠ s.
2. Repeat n − 1 times: relax every edge.
3. One more pass: if an edge still relaxes, a negative cycle is reachable.

Time is Θ(n m).

Bellman-Ford is slower. It is the right default when weights can be negative.

Reconstruct a path with parent[], as in BFS. If dist[t] is still ∞, t is unreachable.

### Questions

#### Theoretical questions

1. What does relaxation do?
2. What weight rule does Dijkstra require?
3. Why can a negative edge break Dijkstra?
4. How does Bellman-Ford detect a negative cycle?
5. What are the usual time classes of the two algorithms?

#### Easy practical tasks

1. Run Dijkstra by hand on the A–B–C picture. Write dist after each extract.
2. Make a table: "Algorithm", "Negative weights?", "Time idea". Add two rows.
3. Write the relax step in four lines.
4. Write five sentences that contrast Dijkstra and Bellman-Ford.

#### Medium practical tasks

1. Give a two-edge example with one negative weight where Dijkstra (if you settle too early) would be wrong. Write the distances that Bellman-Ford should get.
2. Explain in six sentences the n − 1 passes of Bellman-Ford.
3. Implement Dijkstra with a binary heap on a small non-negative graph. Print dist and parent.

#### Advanced practical tasks

1. Implement Bellman-Ford and the extra pass. Test a negative cycle and a negative edge without a cycle.
2. Compare heap Dijkstra with an array-min Dijkstra (scan n vertices for the next u). Write when each wins.

---

## DAG relaxation

A DAG has no directed cycle (Topic 11). Shortest paths on a DAG can use a **single pass** in topological order.

Idea:

1. Compute a topological order.
2. Set dist[s] = 0 and dist[v] = ∞ for v ≠ s.
3. For each vertex u in topological order, relax every outgoing edge of u.

```text
A → B → D
A → C → D
weights on each edge

order: A, B, C, D
relax A's edges, then B's, then C's, then D's
```

Time is Θ(n + m) after you have the adjacency list. That class matches one DFS plus one edge scan.

Weights may be negative. There is no cycle, so there is no negative cycle.

If s is not a source that can reach everyone, some dist values stay ∞. Vertices before s in the order cannot help s. You can start the relax loop at s, or you can still scan all vertices; edges from unreachable vertices do not improve reachable ones if those vertices have dist ∞.

Longest paths on a DAG use the same order with a max relax, or with negated weights. This section only requires shortest paths.

Do not use DAG relaxation on a graph that might have a cycle. Run a cycle check first.

Dijkstra is unnecessary on a DAG. The topological pass is simpler and handles negative weights.

### Questions

#### Theoretical questions

1. Why can a DAG use one relax pass?
2. What order do you use?
3. May DAG shortest-path weights be negative?
4. What is the time class?
5. Why must you not use this method on a cyclic graph?

#### Easy practical tasks

1. Write a topological order for the A–B–C–D picture. Write the relax sequence.
2. Make a table: "Graph", "Shortest-path method". Rows: DAG, non-negative cycles allowed, negative weights no neg cycle.
3. Write five sentences that define DAG relaxation.
4. Draw a DAG with one negative edge. Write that the method still applies.

#### Medium practical tasks

1. Implement topo sort plus DAG relax. Test two paths to the same sink.
2. Explain in six sentences why a cycle would make "one pass" incomplete.
3. Write how you treat vertices that cannot reach from s.

#### Advanced practical tasks

1. Compute longest paths on a DAG by max-relax. Check against a full enumeration on a tiny graph.
2. Compare DAG relax with Bellman-Ford on the same DAG. Write the two times as functions of n and m.

---

## Parent array, union by rank, path compression

A **disjoint-set** structure (Union-Find) stores a partition of n elements. Elements are integers 0 .. n−1.

The **parent** array has length n. parent[i] is the parent of i. If parent[i] = i, then i is a **root**. Each root names one set.

At the start, each element is its own set:

```text
parent:  0  1  2  3  4
index:   0  1  2  3  4
```

**Find** without compression walks parent[i], parent[parent[i]], ... until it reaches a root. That root is the representative.

**Union** without rank attaches one root to the other:

```text
union(a, b):
  ra = find(a)
  rb = find(b)
  if ra == rb: return
  parent[rb] = ra
```

A bad attach order can make a line. Find then costs Θ(n).

**Union by rank** (or by size) attaches the shorter tree under the taller tree. Rank is an upper bound on height. The height stays O(log n) even without path compression.

**Path compression** during find sets every walked node to point at the root. Later finds are almost constant time.

```text
find(4) on 4 → 3 → 2 → 2 (root 2)
after compression: 4 → 2, 3 → 2
```

Together, union by rank and path compression give almost Θ(1) amortized find and union. The exact function is the inverse Ackermann function. For all practical n, treat the cost as very small.

The structure is a forest of trees. It is not a BST. There is no key order.

Space is Θ(n) for parent and rank.

Do not store the full member list in each set unless you need to iterate members.

### Questions

#### Theoretical questions

1. What does parent[i] = i mean?
2. How does find walk to the representative?
3. What does union by rank avoid?
4. What does path compression change?
5. Why is this forest not a search tree?

#### Easy practical tasks

1. Draw parent 0,0,1,3,4 after union(0,1) and union(1,2) for one attach order.
2. Make a table: "Idea", "What it fixes". Rows: union by rank, path compression.
3. Write five sentences that define parent, find, and union.
4. Show path compression of a three-node line.

#### Medium practical tasks

1. Implement find and union with rank only. Build a line on purpose without rank, then with rank. Write the two heights.
2. Explain in six sentences why compression does not break the partition.
3. Write same-set(a, b) as find(a) = find(b).

#### Advanced practical tasks

1. Implement rank and path compression. Run n unions and n finds. Write that you do not need a closed form of Ackermann.
2. Compare size-based union with rank-based union on the same sequence. Write the two max heights.

---

## Kruskal uses Union-Find

**Kruskal's algorithm** builds an MST of a connected undirected weighted graph.

Idea:

1. Sort all edges by increasing weight.
2. Start with an empty forest (n sets in Union-Find).
3. For each edge u–v in sorted order: if find(u) ≠ find(v), add the edge to the MST and union(u, v).
4. Stop when you added n − 1 edges, or after you scanned all edges.

```text
edges sorted: 1, 2, 2, 4, ...
first edge always added
later edge skipped if u and v are already in one set
that skip rejects a cycle
```

Union-Find answers "does this edge close a cycle?" in almost constant time. Without Union-Find, you would search a growing forest after each candidate edge.

If the graph is not connected, Kruskal produces a **minimum spanning forest**: an MST in each component.

Time is dominated by the sort: O(m log m), plus almost linear Union-Find work.

Kruskal does not compute shortest paths from one source. Distances on the MST are not shortest-path distances in the original graph.

Do not run Kruskal on a directed graph as a drop-in MST. MST is an undirected idea.

### Questions

#### Theoretical questions

1. In what order does Kruskal consider edges?
2. When does Kruskal add an edge?
3. How does Union-Find reject a cycle?
4. How many edges does an MST of a connected graph have?
5. What does Kruskal produce on a disconnected graph?

#### Easy practical tasks

1. Draw a triangle with weights 1, 2, 10. Write which two edges Kruskal keeps.
2. Make a table: "Edge", "Add or skip" for a four-vertex example that you draw.
3. Write five sentences that define Kruskal as a Union-Find client.
4. Write why the first (lightest) edge is always safe on a simple connected graph.

#### Medium practical tasks

1. Implement Kruskal with Union-Find. Print the MST weight.
2. Explain in six sentences why a skipped edge would have made a cycle.
3. Write a test: two components, no edge between them. Write the forest size.

#### Advanced practical tasks

1. Compare Kruskal MST weight with a hand MST on a 6-vertex graph.
2. Sort edges once and run Kruskal. Then change one weight and state what you must redo.

---

## MST: Kruskal and Prim (idea)

A **spanning tree** of a connected undirected graph is a subset of edges that connects all vertices and has no cycle. A **minimum spanning tree (MST)** has the minimum total weight among spanning trees.

**Kruskal** grows a forest by global lightest edges that do not form a cycle (previous section).

**Prim** grows one tree from a start vertex.

Idea of Prim:

1. Start from any vertex s. The tree set S = {s}.
2. Repeat until S holds all vertices: add the lightest edge that has one end in S and one end outside S. Add the new vertex to S.

```text
S starts {A}
lightest edge out of S: A—B weight 1
S = {A, B}
next lightest cut edge: B—C weight 2
...
```

A binary heap of candidate edges or of vertices with a key (like Dijkstra) implements Prim in O(m log n).

Prim and Kruskal give the same MST weight. The tree can differ when several MSTs exist (equal weights).

The **cut property** (idea): for a cut that splits V into S and V−S, a lightest edge across the cut is safe for some MST. Prim always picks such an edge. Kruskal's added edges are also safe by a similar argument.

Do not confuse MST with shortest-path tree. The shortest-path tree from s minimizes distances from s. The MST minimizes the sum of all tree edges. The two trees can differ.

```text
A--1--B--1--C
A------3------C

MST can be A—B, B—C (weight 2)
shortest path A to C can be the direct 3 or the path 2, depending on weights
here path A–B–C has length 2 < 3
```

Use MST for network design (connect all sites with minimum cable). Use shortest paths for travel from one source.

### Questions

#### Theoretical questions

1. What is a spanning tree?
2. What extra fact makes it minimum?
3. How does Prim choose the next edge?
4. Why can Prim and Kruskal produce different trees with the same weight?
5. How does an MST differ from a shortest-path tree?

#### Easy practical tasks

1. Draw a square with one diagonal. Assign weights. Write one MST.
2. Make a table: "Algorithm", "Grows", "Main helper". Rows: Kruskal, Prim.
3. Write five sentences that define MST, Kruskal, and Prim.
4. Write one sentence that states the cut-property idea.

#### Medium practical tasks

1. Run Prim by hand from two different starts on the same graph. Write whether the tree differs.
2. Explain in six sentences a case where the shortest-path tree is not an MST.
3. Write why Prim needs a connected graph (or why it only fills one component).

#### Advanced practical tasks

1. Implement Prim with a binary heap of vertices (key = dist to the tree). Compare total weight with Kruskal.
2. Read a short proof of the cut property. Rewrite it in ten STE sentences.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do you choose Dijkstra, Bellman-Ford, or DAG relaxation from weights and cycles?
2. Why is Union-Find the right structure for Kruskal's cycle test and the wrong structure for distances?
3. How do parent[] in shortest paths and parent[] in Union-Find differ in meaning?
4. When do you pick an MST, and when do you pick a shortest-path tree?
5. Which heap operations from Topic 10 appear in Dijkstra and Prim?

#### Easy practical tasks

1. Write a one-page cheat sheet: relax, Dijkstra, Bellman-Ford, DAG, find, union, rank, compression, Kruskal, Prim, MST.
2. Draw one weighted graph. Mark a shortest-path tree from A and one MST. Use two colors if you can.
3. Make a table: "Problem", "Algorithm". Add five rows from this topic.
4. List four defects: Dijkstra on negatives, DAG relax on a cycle, Kruskal without Union-Find scans, MST used as travel distances.

#### Medium practical tasks

1. Implement Dijkstra and Kruskal on the same weighted undirected graph. Write that the outputs answer different questions.
2. Write a test plan: negative cycle, unreachable vertex, disconnected MST forest, equal-weight MSTs.
3. Design an ADT for Union-Find: make(n), find, union, same. Write the parent and rank fields.

#### Advanced practical tasks

1. Implement Bellman-Ford, DAG relax, and heap Dijkstra. Route each test graph to the cheapest correct method.
2. Read a textbook chapter on MST or shortest paths. Map five lemmas to the words in this topic.
