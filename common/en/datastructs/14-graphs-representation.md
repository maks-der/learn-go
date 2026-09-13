# 14. Graphs (Representation)

## Description

A graph is a set of vertices and a set of edges. This topic explains vertices, edges, directed and undirected graphs, weights, adjacency lists, adjacency matrices, implicit graphs, and the words degree, path, cycle, and DAG. Complete this topic after stacks, queues, and trees. Complete this topic before graph algorithms.

Use one term for each concept. A **vertex** is a node of the graph. An **edge** joins two vertices. This handbook uses **vertex**, not "node", for graph points. It uses **node** for tree cells and list cells.

---

## Vertex and edge

A **graph** G is a pair (V, E). V is the set of vertices. E is the set of edges.

An **undirected edge** is a two-element set {u, v}. The pair has no direction. {u, v} is the same as {v, u}.

A **directed edge** (also called an **arc**) is an ordered pair (u, v). The edge goes from u to v. (u, v) is not the same as (v, u).

A **self-loop** is an edge from a vertex to itself. A **multigraph** allows more than one edge between the same pair. This handbook uses **simple graphs** unless a section says otherwise: no self-loops, no parallel edges.

```text
vertices:  A B C
edges:     A—B, B—C

A — B — C
```

You name vertices with integers 0 .. n−1 in code. You can also use strings and a map from name to index.

An edge can carry extra data later (a weight). The pair of endpoints is still the edge.

A tree is a special graph: connected, no cycle. You already used trees. A general graph can have cycles and more than one path between two vertices.

Do not call an edge a "line" in writing. Say **edge**. Do not call a vertex a "point" in writing. Say **vertex**.

### Questions

#### Theoretical questions

1. What is a graph as a pair of sets?
2. How is an undirected edge different from a directed edge?
3. What is a self-loop?
4. What does this handbook mean by a simple graph?
5. Why does this handbook use "vertex" instead of "node" in a graph?

#### Easy practical tasks

1. Draw three vertices and two edges. Write V and E as sets.
2. Make a table: "Term", "Meaning". Rows: vertex, edge, self-loop.
3. Write five sentences that define a graph.
4. Rename vertices A, B, C to 0, 1, 2. Rewrite E with integers.

#### Medium practical tasks

1. Write a Go or Python list of vertices and a list of pairs for a four-vertex path.
2. Draw one graph that is a tree and one graph that is not a tree. Mark a cycle if it exists.
3. Explain in six sentences when you use string names plus a map to indexes.

#### Advanced practical tasks

1. Design a Vertex type and an Edge type. Document whether Edge is directed. Do not write algorithms yet.
2. Read one standard definition (CLRS or a course note). Map their terms to vertex, edge, simple graph. Write ten STE sentences.

---

## Directed vs undirected; weighted vs unweighted

**Undirected.** An edge {u, v} means you can go from u to v and from v to u. A road map without one-way streets is undirected.

**Directed.** An edge (u, v) means you can go from u to v. You cannot use that edge to go back. A one-way street is directed. A web link is directed.

**Unweighted.** Each edge has the same cost, often 1. A shortest path then means the fewest edges. BFS solves that problem.

**Weighted.** Each edge has a number: length, time, or cost. Shortest path means smallest sum of weights. Dijkstra or Bellman-Ford solve that problem. Those algorithms are the next topic.

```text
undirected unweighted     A — B — C
directed unweighted       A → B → C
undirected weighted       A —3— B —1— C
directed weighted         A -3→ B -1→ C
```

A **non-negative weight** is greater than or equal to 0. Dijkstra needs that rule. A **negative weight** is allowed in some models (a credit, a reversal). Bellman-Ford can handle negative weights if there is no negative cycle.

You can store an undirected graph as a directed graph with two arcs for each undirected edge. You must keep the two arcs in sync when you edit.

Do not mix "weighted" with "directed". The two choices are independent. Four combinations exist.

When you write a problem, name all four words: directed or undirected, weighted or unweighted.

### Questions

#### Theoretical questions

1. What does an undirected edge allow that one directed edge does not?
2. What does a weight measure?
3. Why does unweighted shortest path use BFS?
4. Why does Dijkstra need non-negative weights?
5. How can you represent an undirected graph with directed arcs?

#### Easy practical tasks

1. Draw the same three vertices as undirected and as directed. Label the difference.
2. Make a table with four rows: the four combinations. Add one real example each.
3. Assign weights 2 and 5 to two edges. Write the path cost A–B–C.
4. Write five sentences that separate direction from weight.

#### Medium practical tasks

1. Convert an undirected edge list to a directed pair list (two arcs per edge). Write the list.
2. Mark which of these is weighted: "friend link", "flight time in minutes", "follows on a site".
3. Explain in six sentences why a negative flight price model needs a different algorithm than Dijkstra.

#### Advanced practical tasks

1. Write a small data file format: first line n m directed_flag weighted_flag, then m edge lines. Document each field.
2. Load that format into memory as lists of integer triples (u, v, w). Print the four flags that you stored.

---

## Adjacency list vs adjacency matrix

An **adjacency list** stores, for each vertex, a list of neighbors. For a directed graph, the list holds out-neighbors. For a weighted graph, each entry is (neighbor, weight).

An **adjacency matrix** is an n × n array. Slot [u][v] is 1 (or the weight) if an edge from u to v exists. Slot [u][v] is 0 (or ∞) if the edge is absent. For an undirected graph the matrix is symmetric.

```text
vertices 0, 1, 2
edges 0—1, 1—2

adj list:
  0: 1
  1: 0, 2
  2: 1

adj matrix (undirected):
    0 1 2
  0 0 1 0
  1 1 0 1
  2 0 1 0
```

```go
// adjacency list, unweighted undirected
g := [][]int{
	{1},
	{0, 2},
	{1},
}
```

Space of a list is Θ(n + m) where m is the number of edges. Space of a matrix is Θ(n²).

Edge query "is (u, v) present?" is Θ(1) in a matrix and Θ(degree(u)) in a list.

Scan of all neighbors of u is Θ(degree(u)) in a list and Θ(n) in a matrix.

Most sparse graphs (m much smaller than n²) use a list. Dense graphs and some closure algorithms use a matrix.

Do not store a huge matrix for a huge sparse graph. You will waste memory.

An edge list (a list of pairs) is a third form. It is good for input and for Kruskal. It is slow for neighbor scans. Build a list or a matrix from the edge list before BFS.

### Questions

#### Theoretical questions

1. What does an adjacency list store for each vertex?
2. What does slot [u][v] mean in an adjacency matrix?
3. What is the space class of each representation?
4. Which form answers "is (u, v) an edge?" in Θ(1)?
5. When do you keep only an edge list?

#### Easy practical tasks

1. Draw the list and the matrix for a triangle of 3 vertices.
2. Make a table: "Need", "Better form". Rows: sparse social graph, dense relation on 20 vertices, Kruskal input.
3. Write five sentences that compare the two forms.
4. For n = 4, m = 3, write n + m and n².

#### Medium practical tasks

1. Implement an undirected adjacency list from an edge list. Print neighbors of each vertex.
2. Implement an adjacency matrix from the same edges. Check symmetry.
3. Time a neighbor scan of one vertex on both forms for n = 2 000, m = 3 000 if you can.

#### Advanced practical tasks

1. Support weighted directed edges in both forms. Document the missing-edge value in the matrix (0 versus ∞).
2. Write a one-page note: pick a form for a grid map, a road map, and Floyd–Warshall. Give one reason each.

---

## Implicit graphs

An **implicit graph** has vertices and edges that you do not store as a full list. You compute neighbors when you need them.

Examples:

- A **grid**. A cell (r, c) is a vertex. Neighbors are up, down, left, and right (if inside the bounds). You do not store 4m edges.
- A **game state**. A position is a vertex. A legal move is an edge. The state space is too large to store.
- A **number graph**. From x you can go to x + 1 or 2x. You generate those two neighbors in code.

```text
grid 2 x 2

(0,0) — (0,1)
  |       |
(1,0) — (1,1)

neighbors((0,0)) = (0,1), (1,0)
```

BFS and DFS do not need a stored edge list. They need a function `neighbors(v)` and a way to mark visited vertices.

Visited storage must still exist. On a grid you use a 2-d boolean array. On a huge state space you use a hash set of seen states.

An implicit graph can be infinite if you do not bound it. Always define the vertex set. Example: integers from 0 to N, or cells of a finite grid.

Do not build a giant adjacency list for a grid only to run BFS. Write `neighbors` and a queue.

A stored graph is **explicit**. Use that word when you contrast it with implicit.

### Questions

#### Theoretical questions

1. What is an implicit graph?
2. How do you get neighbors if you do not store edges?
3. Why does a game-state graph often stay implicit?
4. What must you still store during BFS?
5. Why must the vertex set be bounded?

#### Easy practical tasks

1. Write the neighbors of cell (1, 1) on a 3 × 3 grid. Exclude out-of-bounds cells.
2. Make a table: "Domain", "Vertex", "Neighbor rule". Rows: grid, increment/double.
3. Draw a 2 × 2 grid as a picture and as four vertices.
4. Write five sentences that define implicit versus explicit.

#### Medium practical tasks

1. Write `neighbors(r, c, rows, cols)` that returns legal 4-way cells. Test a corner and a center.
2. Explain in six sentences how you mark visited cells on a grid.
3. Describe a knight on a chessboard as an implicit graph. List the move offsets. Do not search yet.

#### Advanced practical tasks

1. Run BFS on a small grid from (0, 0) with only `neighbors` and a queue. Print distances. Keep this as a preview of the next topic.
2. Write a one-page note: memory of an explicit list of all grid edges versus a `neighbors` function.

---

## Degree, path, cycle, DAG

The **degree** of a vertex in an undirected graph is the number of incident edges. In a directed graph, **out-degree** is the number of outgoing edges. **In-degree** is the number of incoming edges.

A **path** is a sequence of vertices v0, v1, ..., vk where each (vi, vi+1) is an edge (or {vi, vi+1} in the undirected case). The **length** of a path is the number of edges k, or the sum of weights if the graph is weighted. Name which length you use.

A **simple path** does not repeat a vertex. This handbook says **simple path** when that rule matters.

A **cycle** is a path with k ≥ 1 (directed: k ≥ 1; undirected often k ≥ 3) that starts and ends at the same vertex and does not repeat other vertices. An undirected graph that contains a cycle is **cyclic**. A graph with no cycle is **acyclic**.

A **DAG** is a directed acyclic graph. Tasks with prerequisites form a DAG. Topological sort exists if and only if the directed graph is a DAG. That algorithm is the next topic.

```text
path:   A-B-C
cycle:  A-B-C-A   (directed)
DAG:    A→B→C    (no back edge)
not DAG: A→B→C→A
```

A **connected** undirected graph has a path between every pair of vertices. A **weakly connected** directed graph is connected if you ignore direction. A **strongly connected** directed graph has a directed path from every vertex to every other vertex. Strong connectivity is a later topic.

Do not call a tree a DAG without care. A rooted tree as directed edges from parent to child is a DAG. An undirected tree is not a directed graph.

Handshaking lemma (undirected): the sum of degrees is 2m. Use it to check a drawing.

### Questions

#### Theoretical questions

1. What is degree in an undirected graph?
2. What is the difference between in-degree and out-degree?
3. What is a path, and what is a cycle?
4. What is a DAG?
5. When does a topological sort exist?

#### Easy practical tasks

1. On a square of 4 vertices, write the degree of each vertex.
2. Make a table: "Word", "Needs direction?". Rows: degree, path, DAG, cycle.
3. Draw one DAG of 4 vertices. Write one valid vertex order that respects edges.
4. Write five sentences that define path versus cycle.

#### Medium practical tasks

1. Compute in-degree and out-degree from a directed edge list. Print both arrays.
2. Mark a cycle on a drawing or write "no cycle". Give the vertex sequence if a cycle exists.
3. Explain in six sentences why a prerequisite graph must be a DAG.

#### Advanced practical tasks

1. Write `hasUndirectedCycle` later in the algorithm topic. Here, write only the definition and three test drawings: a tree, a unicycle, two components.
2. Write a one-page glossary: degree, path, simple path, cycle, DAG, connected, strongly connected. Use one sentence each.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe the path from a real map to a stored graph. Name vertex, edge, direction, weight, and representation.
2. How do adjacency lists, matrices, and implicit `neighbors` share the same ADT of "list of neighbors"?
3. A teammate says "the graph is a matrix". Which facts do you use to correct that sentence?
4. Why must you name directed or undirected before you name degree or cycle?
5. What from this topic does BFS need, and what can wait for the next topic?

#### Easy practical tasks

1. Write a one-page cheat sheet: V, E, four graph kinds, list, matrix, implicit, degree, path, cycle, DAG.
2. Draw one weighted directed graph of 4 vertices. Write it as a list and as a matrix.
3. Make a table: "Question about G", "Which representation answers it cheaply?". Four rows.
4. Convert a 3 × 3 grid into an explicit undirected edge list. Count m.

#### Medium practical tasks

1. Implement load of an edge list into an adjacency list and into a matrix. Print both.
2. Write `neighbors` for a grid and compare its output with an explicit list of grid edges.
3. Classify five real systems (web, roads, tasks, grid maze, friendships) with the four kind words plus list or implicit.

#### Advanced practical tasks

1. Build a tiny graph library: add vertex, add edge, list neighbors, edge exists. Support a directed flag. Tests on a path, a cycle, and a DAG.
2. Read the next topic list (BFS, DFS, Dijkstra). For each algorithm name, write which representation you will pass in. Do not implement the algorithms here.
