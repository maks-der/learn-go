# 11. Graphs and Traversal

## Description

A graph is a set of vertices and a set of edges. This topic explains directed and undirected graphs, weights, adjacency lists and matrices, degree, path, cycle, DAG, BFS and unweighted shortest path, DFS, cycle detection, and topological sort.

Complete this topic after stacks, queues, and trees. Complete this topic before shortest paths and Union-Find.

Use one term for each concept. A **vertex** is a node of the graph. An **edge** is a link between two vertices. A **path** is a sequence of vertices that follows edges. A **cycle** is a path that returns to its start and has at least one edge. A **DAG** is a directed acyclic graph. **BFS** is breadth-first search. **DFS** is depth-first search.

---

## Vertex, edge, directed vs undirected, weighted

A **graph** G = (V, E). V is the set of vertices. E is the set of edges.

An **undirected** edge {u, v} has no direction. You can walk from u to v and from v to u. {u, v} is the same edge as {v, u}.

A **directed** edge (u, v) goes from u to v. You cannot walk back unless (v, u) also exists.

A **weighted** edge has a number: a length, a cost, or a capacity. An **unweighted** graph treats each edge as length 1.

```text
undirected unweighted:
  A — B
  |   |
  C — D

directed weighted:
  A --2--> B
  A --5--> C
```

A **self-loop** is an edge from a vertex to itself. A **multiedge** is more than one edge between the same pair. This handbook uses simple graphs unless a line says otherwise: no self-loops, at most one undirected edge per pair, at most one directed edge per ordered pair.

n = |V| is the number of vertices. m = |E| is the number of edges.

A **dense** graph has m near n². A **sparse** graph has m much smaller than n², often O(n).

Do not mix directed and undirected in one picture. Draw arrows or draw plain lines.

A tree is a special undirected graph: connected and acyclic. A directed tree is a different picture. Use the words from Topic 7 when you mean a rooted tree.

### Questions

#### Theoretical questions

1. What is a vertex?
2. What is the difference between {u, v} and (u, v)?
3. What does a weight represent at a high level?
4. What is a simple graph in this handbook?
5. What is the difference between dense and sparse?

#### Easy practical tasks

1. Draw an undirected graph of four vertices and four edges.
2. Draw a directed graph of three vertices with two arrows.
3. Make a table: "Kind", "How you walk an edge". Rows: undirected, directed.
4. Write five sentences that define vertex, edge, directed, undirected, and weighted.

#### Medium practical tasks

1. Count n and m on your two drawings.
2. Explain in six sentences when you model a street as directed.
3. Write why a self-loop is not in a simple graph.

#### Advanced practical tasks

1. Encode a small weighted digraph as a list of triples (u, v, w). Write n and m.
2. Read one definition of simple graph in a textbook. Map it to this section in four sentences.

---

## Adjacency list vs matrix

An **adjacency list** stores, for each vertex, a list of neighbors (and weights if needed).

Space is Θ(n + m). Walk of all neighbors of u is Θ(degree(u)). Walk of all edges is Θ(n + m).

An **adjacency matrix** is an n × n table. Cell [u][v] is 1 (or the weight) if the edge exists, else 0 (or ∞).

Space is Θ(n²). Test "is (u, v) an edge?" is Θ(1). Walk of all neighbors of u is Θ(n).

```text
vertices 0, 1, 2
edges 0—1, 1—2

list:
  0: 1
  1: 0, 2
  2: 1

matrix (undirected):
    0 1 2
  0 0 1 0
  1 1 0 1
  2 0 1 0
```

Use a list for a sparse graph. Use a matrix when n is small or when you need constant-time edge tests and you can pay n² space.

For an undirected graph, the list stores each edge in two lists. The matrix is symmetric.

For a directed graph, the list of u stores only out-neighbors. The matrix cell [u][v] is the arc u → v.

Do not scan the full matrix if you only need existing edges on a large sparse graph. That scan is Θ(n²).

An **edge list** is a third form: a list of pairs or triples. It is simple to store. Neighbor walk needs a scan of all m edges unless you build a list.

### Questions

#### Theoretical questions

1. What does an adjacency list store for one vertex?
2. What is the space of an adjacency matrix?
3. When is an edge test Θ(1)?
4. Why does an undirected list store each edge twice?
5. When do you prefer a list for a sparse graph?

#### Easy practical tasks

1. Write the list and the matrix for the 0–1–2 picture.
2. Make a table: "Need", "List or matrix". Rows: sparse walk, edge test, tiny n.
3. Draw a directed list for 0→1 and 1→2.
4. Write five sentences that contrast list and matrix.

#### Medium practical tasks

1. For n = 10 000 and m = 30 000, write the two space classes.
2. Explain in six sentences why neighbor walk on a matrix is Θ(n).
3. Write how you store weights in a list and in a matrix.

#### Advanced practical tasks

1. Implement both representations for a small graph. Time "list all edges" on a sparse instance.
2. Build a list from an edge list in Θ(n + m). Write the steps.

---

## Degree, path, cycle, DAG

The **degree** of a vertex in an undirected graph is the number of incident edges. In a directed graph, **out-degree** is the number of outgoing edges. **In-degree** is the number of incoming edges.

A **path** is a sequence of distinct vertices v0, v1, ..., vk where each (vi, vi+1) is an edge (or {vi, vi+1} in the undirected case). The **length** of the path can count edges (k) or the sum of weights.

A **cycle** is a closed walk with at least one edge that does not repeat vertices except the start and the end. (Exact textbook wording varies. This handbook uses a simple cycle: no repeated vertex except the close.)

An undirected graph is **connected** if every pair of vertices has a path. A directed graph is **strongly connected** if every pair has a directed path in both directions. **Weakly connected** ignores direction.

A **DAG** is a directed graph with no directed cycle.

```text
DAG:
  A → B → D
  A → C → D

not a DAG:
  A → B → A
```

DAGs model tasks with prerequisites. Topological sort (later in this topic) needs a DAG.

The **handshaking** fact: the sum of degrees in an undirected graph is 2m. Each edge adds 1 to two degrees.

Do not call a directed cycle an undirected cycle without care. The same drawing can have a directed cycle and still be connected as an undirected picture.

### Questions

#### Theoretical questions

1. What is degree in an undirected graph?
2. What is a path?
3. What is a simple cycle in this handbook?
4. What is a DAG?
5. Why is the sum of degrees equal to 2m?

#### Easy practical tasks

1. Write degrees for a square of four vertices and four edges.
2. Draw a DAG of four vertices. Draw a directed cycle of three vertices.
3. Make a table: "Word", "Needs direction?". Rows: degree, out-degree, DAG, connected.
4. Write five sentences that define degree, path, cycle, and DAG.

#### Medium practical tasks

1. Write one path and one cycle on a cycle graph of five vertices.
2. Explain in six sentences why a task graph with a cycle cannot have a valid finish order.
3. Compute in-degree and out-degree for A, B, C in A→B, A→C, B→C.

#### Advanced practical tasks

1. Write a checker that a directed edge list is a DAG by a method you will specify after you read the DFS section (plan the tests now).
2. Prove the handshaking fact in six short sentences.

---

## BFS and unweighted shortest path

**Breadth-first search (BFS)** visits vertices in order of unweighted distance from a start s. BFS uses a **queue**. Each edge has length 1.

Idea:

1. Set dist[s] = 0. Set dist[v] = ∞ for every other vertex.
2. Enqueue s. Mark s as visited.
3. While the queue is not empty, dequeue u. For each unvisited neighbor v of u, set dist[v] = dist[u] + 1, mark v visited, enqueue v.

```text
s=0
0 — 1 — 2
    |
    3

queue order: 0, 1, 2, 3   (if neighbor order is 1 then 2 from 0, and 3 from 1)
dist: 0:0  1:1  2:1  3:2
```

Visited-on-enqueue is the usual unweighted form. You mark v when you first find v. That first find is a shortest path in number of edges.

Store **parent[v] = u** when you set dist[v]. Reconstruct the path from s to t by walking parent from t back to s, then reverse the list.

BFS on an implicit grid is the same algorithm. Vertices are cells. Neighbors come from a neighbor function.

Time is Θ(n + m) on an adjacency list. Each vertex enters the queue at most one time. Each edge is scanned a constant number of times.

BFS does not compute weighted shortest paths. A small weight on a long path can win. Use Dijkstra or Bellman-Ford (Topic 12).

Do not use a stack for this order. A stack gives DFS, not shortest unweighted paths.

### Questions

#### Theoretical questions

1. Which ADT does BFS use?
2. Why is the first time you reach v a shortest unweighted path?
3. How do you reconstruct the path with a parent array?
4. What is the time class on an adjacency list?
5. Why does BFS fail as a weighted shortest-path algorithm?

#### Easy practical tasks

1. Run BFS on the 0–1–2–3 picture. Write the queue after each dequeue.
2. Make a table: "Vertex", "dist", "parent".
3. Write five sentences that define unweighted shortest path.
4. Draw a 2 × 2 grid. Write BFS distances from (0, 0).

#### Medium practical tasks

1. Implement BFS with a queue and a parent array. Print the path from s to each reachable vertex.
2. Run BFS on a grid with one blocked cell. Print dist of the opposite corner or "unreachable".
3. Explain in six sentences why you mark visited when you enqueue, not when you dequeue, for unweighted shortest path.

#### Advanced practical tasks

1. Implement BFS from every vertex on a small undirected graph and compare dist[u][v] with dist[v][u].
2. Write tests: empty neighbor list, disconnected target, a cycle. Path length must match a hand count.

---

## DFS, cycle detect, topological sort

**Depth-first search (DFS)** goes as deep as possible along one path, then backtracks. DFS uses the **call stack** (recursion) or an **explicit stack**.

A useful DFS records two times per vertex:

- **discover** time: the step when you first enter the vertex
- **finish** time: the step when you leave the vertex after all descendants

Color or state helps:

- **white:** not seen
- **gray:** discovered, not finished (on the recursion stack)
- **black:** finished

**Cycle detection in a directed graph.** A **back edge** goes to a gray vertex. That edge closes a directed cycle.

**Cycle detection in an undirected graph.** An edge to an already discovered vertex that is not the parent is a cycle edge. The parent check is required because the tree edge appears in both directions in the list.

**Topological sort** of a DAG is an order of vertices such that every directed edge goes from an earlier vertex to a later vertex.

One method: run DFS. After the graph is finished, list vertices in **decreasing finish time**. That list is a topological order if the graph is a DAG.

```text
A → B → D
A → C → D

one topological order: A, B, C, D
another: A, C, B, D
```

A second method: Kahn's algorithm. Repeatedly dequeue a vertex with in-degree 0. Decrease in-degrees of its neighbors. If you cannot dequeue n vertices, a cycle exists.

Time of DFS on a list is Θ(n + m). Time of topological sort is the same.

Do not run topological sort on a graph that can have a cycle without a cycle check.

DFS does not compute unweighted shortest paths. Visit order is not by distance.

### Questions

#### Theoretical questions

1. What do discover time and finish time record?
2. What is a back edge in a directed DFS?
3. Why must undirected cycle detection ignore the parent?
4. What is a topological order?
5. How does decreasing finish time give a topological order on a DAG?

#### Easy practical tasks

1. Run a recursive DFS on the A–B–C–D DAG. Write one possible discover/finish table.
2. Draw a directed cycle of three vertices. Mark a back edge.
3. Make a table: "Problem", "BFS or DFS". Rows: unweighted dist, cycle in a digraph, topological order.
4. Write five sentences about DFS, cycle detect, and topological sort.

#### Medium practical tasks

1. Implement directed cycle detection with gray vertices. Test a DAG and a cycle.
2. Implement topological sort by finish times. Test two valid orders on the A–B–C–D picture.
3. Explain in six sentences Kahn's algorithm and the leftover vertices when a cycle exists.

#### Advanced practical tasks

1. Implement both topological methods. Check that they reject the same cyclic graphs.
2. Write a compiler-style task list (compile A before B). Model it as a DAG and print one legal order.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do n, m, and density decide list versus matrix for BFS and DFS?
2. When is a path a shortest unweighted path, and when is it only a walk from DFS?
3. How do gray vertices and in-degree 0 both detect that a directed graph is not a DAG?
4. Why does parent[] from BFS reconstruct a path while finish times reconstruct an order?
5. Which facts from Topic 4 (queue and stack) are required to implement this topic?

#### Easy practical tasks

1. Write a one-page cheat sheet: vertex, edge, directed, weight, list, matrix, degree, path, cycle, DAG, BFS, DFS, topo.
2. Draw one graph. Show BFS layers from a start and one DFS path from the same start.
3. Make a table: "Algorithm", "ADT", "Result". Add BFS, DFS, Kahn.
4. List four defects: stack used for unweighted shortest path, matrix scan on huge sparse graphs, topo on a cycle, undirected cycle check without parent.

#### Medium practical tasks

1. Implement BFS shortest path and DFS cycle detect on the same adjacency list type.
2. Write a report: connected components via BFS or DFS on an undirected graph (idea: restart from an unvisited vertex).
3. Design tests for a grid: walls, no path, multiple shortest paths (parent records one).

#### Advanced practical tasks

1. Implement implicit BFS on a large grid without storing all cells in a matrix of edges. Write the neighbor function.
2. Read a textbook page on edge classification (tree, back, forward, cross). Map each class to cycle detection in eight sentences.
