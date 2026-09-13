# 15. Graph Algorithms (Core)

## Description

A graph algorithm walks vertices and edges to answer a question. This topic explains BFS, unweighted shortest path, DFS, timestamps, cycle detection, topological sort, Dijkstra, Bellman-Ford, DAG relaxation, and a preview of connected components. Complete this topic after graph representation, queues, stacks, and heaps.

Use one term for each concept. A **visit** marks a vertex as found. **Relaxation** updates a tentative distance when a better path appears. A **timestamp** is a step counter during DFS.

---

## BFS and shortest unweighted path

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

BFS on an implicit grid is the same algorithm. Vertices are cells. Neighbors come from `neighbors(cell)`.

Time is Θ(n + m) on an adjacency list. Each vertex enters the queue at most one time. Each edge is scanned a constant number of times.

BFS does not compute weighted shortest paths. A small weight on a long path can win. Use Dijkstra or Bellman-Ford for weights.

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

1. Implement BFS with a queue and a parent array. Print the path from s to each vertex.
2. Run BFS on a grid with one blocked cell. Print dist of the opposite corner or "unreachable".
3. Explain in six sentences why you mark visited when you enqueue, not when you dequeue, for unweighted shortest path.

#### Advanced practical tasks

1. Implement BFS from every vertex on a small graph and compare dist[u][v] with dist[v][u] on an undirected graph.
2. Write tests: empty neighbor list, disconnected target, a cycle. Path length must match a hand count.

---

## DFS, timestamps, cycle detect

**Depth-first search (DFS)** goes as deep as possible along one path, then backtracks. DFS uses the **call stack** (recursion) or an **explicit stack**.

A useful DFS records two times per vertex:

- **discover** time: the step when you first enter the vertex
- **finish** time: the step when you leave the vertex after all descendants

```text
dfs(u):
  time = time + 1
  disc[u] = time
  color[u] = gray
  for each neighbor v:
      if color[v] is white: dfs(v)
      else if color[v] is gray:  // back edge, directed cycle
          report cycle
  color[u] = black
  time = time + 1
  fin[u] = time
```

Colors:

- **white**: not discovered
- **gray**: discovered, not finished
- **black**: finished

In a **directed** graph, a gray neighbor is a **back edge**. A back edge means a directed cycle.

In an **undirected** graph, a neighbor that is the parent is not a cycle by itself. A gray or visited neighbor that is not the parent means a cycle.

```text
directed cycle:  0 → 1 → 2 → 0
dfs(0) discovers 1, then 2, then edge to gray 0
```

Timestamps order vertices. Finish times are the key for topological sort and for some strongly connected component methods.

Time of DFS is Θ(n + m) on an adjacency list if you visit each vertex once (a forest of DFS trees).

Do not use DFS when you need unweighted shortest paths. Use BFS.

### Questions

#### Theoretical questions

1. What do discover time and finish time record?
2. What does gray mean?
3. How do you detect a directed cycle with colors?
4. Why is the parent not a cycle edge in an undirected graph?
5. What is the time class of DFS on an adjacency list?

#### Easy practical tasks

1. Draw a directed path 0→1→2. Write disc and fin if you start at 0.
2. Make a table: "Color", "Meaning". Three rows.
3. Mark a back edge on a three-vertex directed cycle.
4. Write five sentences that contrast DFS order with BFS order.

#### Medium practical tasks

1. Implement recursive DFS with disc, fin, and color. Print the two times.
2. Implement directed cycle detection. Test a DAG and a cycle.
3. Implement undirected cycle detection with a parent argument. Test a tree and a square.

#### Advanced practical tasks

1. Implement iterative DFS with an explicit stack. Match the visit set of recursive DFS (order can differ). Document the difference.
2. Write a one-page note: four edge types in directed DFS (tree, back, forward, cross) from a textbook. Cite the source. Do not invent names.

---

## Topological sort

A **topological sort** of a DAG is a list of all vertices. For every edge u → v, u appears before v in the list. The list is a valid work order for prerequisites.

A topological sort exists if and only if the directed graph has no cycle.

**DFS method.** Run DFS. After the search, list vertices in decreasing finish time. That list is a topological order if the graph is a DAG. If DFS finds a back edge, stop. There is no topological order.

**Kahn method.**

1. Compute in-degree of each vertex.
2. Enqueue every vertex with in-degree 0.
3. While the queue is not empty, dequeue u, append u to the order, and decrease in-degree of each neighbor. If a neighbor reaches 0, enqueue it.
4. If the order length is less than n, a cycle exists.

```text
edges: 0→2, 1→2, 2→3
one order: 0, 1, 2, 3
also valid: 1, 0, 2, 3
```

More than one order can exist. Any valid order is a topological sort.

Time is Θ(n + m) for both methods.

Do not run topological sort on an undirected graph. Direction is required.

Use the order in DAG relaxation: process vertices in topological order, then relax outgoing edges.

### Questions

#### Theoretical questions

1. What does a topological sort guarantee for each edge?
2. When does a topological sort not exist?
3. How do DFS finish times give an order?
4. What is the role of in-degree 0 in Kahn's method?
5. Why must the graph be directed?

#### Easy practical tasks

1. Draw a four-vertex DAG. Write two different valid orders if they exist.
2. Make a table: "Method", "Main tool". Rows: DFS, Kahn.
3. Write five sentences that define topological sort.
4. For edges A→B, A→C, write the positions that A, B, and C may take.

#### Medium practical tasks

1. Implement Kahn. Print the order or "cycle".
2. Implement DFS topological sort. Compare the two orders on one DAG. Both must be valid.
3. Add one back edge. Show that both methods report failure.

#### Advanced practical tasks

1. Generate a random DAG (edges only from lower index to higher). Check that Kahn returns n vertices.
2. Write a course-planner: tasks with prerequisites. Print one legal study order.

---

## Dijkstra

**Dijkstra** computes shortest paths from s in a directed or undirected graph with **non-negative** edge weights.

Idea:

1. Set dist[s] = 0 and dist[v] = ∞ for other vertices.
2. Put vertices in a min-heap by dist (or insert s and insert others as you find them).
3. While the heap is not empty, extract the vertex u with smallest dist. For each edge u → v of weight w, **relax**: if dist[u] + w < dist[v], set dist[v] = dist[u] + w and update the heap (decrease-key or insert a new pair).

```text
relax(u, v, w):
  if dist[v] > dist[u] + w:
      dist[v] = dist[u] + w
      parent[v] = u
```

When you extract u, dist[u] is final if all weights are ≥ 0. You do not need to process u again in the classic form. If you use extra inserts instead of decrease-key, you may extract stale pairs. Skip a pair if its heap distance is worse than the current dist[u].

```text
s=0
0 -1→ 1 -4→ 3
0 -2→ 2 -1→ 3

dist[3] = 3  (path 0-2-3)
```

Time with a binary heap is O((n + m) log n) in the extra-insert form, or O(m log n) with a careful decrease-key form. Learn the idea first. Then pick one heap API.

Do not run Dijkstra with a negative weight. A later extract can be wrong.

Unweighted graphs can use BFS. Dijkstra still works if every weight is 1, but BFS is simpler.

Store parent as in BFS to reconstruct paths.

### Questions

#### Theoretical questions

1. What weight rule does Dijkstra need?
2. What does relax do?
3. Why is dist[u] final when you extract u?
4. How do you handle stale heap pairs if you insert duplicates?
5. Why is BFS enough when every weight is 1?

#### Easy practical tasks

1. Run Dijkstra by hand on the 0–1–2–3 picture. Write dist after each extract.
2. Make a table: "Algorithm", "Weights". Rows: BFS, Dijkstra.
3. Write five sentences that define relaxation.
4. Draw one negative edge. Write why this handbook forbids Dijkstra there.

#### Medium practical tasks

1. Implement Dijkstra with a binary heap of (dist, vertex) pairs. Skip stale pairs. Test the picture.
2. Reconstruct the path to each vertex with parent.
3. Compare Dijkstra with weight 1 on all edges to BFS distances. They must match.

#### Advanced practical tasks

1. Implement decrease-key with a map from vertex to heap index, or document why you use extra inserts.
2. Time Dijkstra on a larger sparse graph. Write n, m, and the time.

---

## Bellman-Ford

**Bellman-Ford** computes shortest paths from s when edge weights can be **negative**. It also detects a **negative cycle** that is reachable from s.

Idea:

1. Set dist[s] = 0 and dist[v] = ∞ for other vertices.
2. Repeat n − 1 times: relax every edge (u, v, w) in the graph.
3. Do one more pass over every edge. If any relax still succeeds, a negative cycle is reachable from s (or from a vertex that this pass can improve). Report "negative cycle".

```text
for i = 1 to n - 1:
    for each edge (u, v, w):
        relax(u, v, w)
for each edge (u, v, w):
    if dist[v] > dist[u] + w:  // still improvable
        report negative cycle
```

After n − 1 successful rounds, a shortest simple path (at most n − 1 edges) is found if no negative cycle affects s.

Time is Θ(n m). That is slower than Dijkstra. Use Bellman-Ford when a weight can be negative or when you must detect a negative cycle.

A negative cycle that you cannot reach from s does not always matter. Define whether you care about the whole graph or only about vertices reachable from s.

Do not use Bellman-Ford as your first shortest-path code. Implement BFS, then Dijkstra, then Bellman-Ford.

You can stop early if one full edge pass does no update.

### Questions

#### Theoretical questions

1. What extra problem can Bellman-Ford solve that Dijkstra cannot?
2. Why do you relax all edges n − 1 times?
3. What does a successful relax in pass n mean?
4. What is the time class?
5. When do you still prefer Dijkstra?

#### Easy practical tasks

1. Write the two-phase outline: n − 1 passes, then one detect pass.
2. Make a table: "Algorithm", "Negative weights?", "Negative cycle detect?".
3. Draw a three-vertex negative cycle. Write the three weights so the sum is negative.
4. Write five sentences that define Bellman-Ford.

#### Medium practical tasks

1. Implement Bellman-Ford from an edge list. Test a graph with one negative edge and no negative cycle.
2. Add a negative cycle reachable from s. Show that the detect pass fires.
3. Explain in six sentences why n − 1 is enough for a simple path.

#### Advanced practical tasks

1. Implement early stop when a pass does no relax. Compare pass counts on a DAG versus a dense graph.
2. Mark vertices that a negative cycle can reach (extra BFS or extra passes). Write which dist values you must treat as −∞.

---

## DAG relaxation

A **DAG** has a topological order. Shortest paths (and longest paths) on a DAG use that order.

**Shortest paths on a DAG:**

1. Compute a topological order. If the graph is not a DAG, stop.
2. Set dist[s] = 0 and dist[v] = ∞ for other vertices.
3. For each vertex u in the topological order, relax every outgoing edge u → v.

```text
order: 0, 1, 2, 3
for u in order:
    for each edge u → v:
        relax(u, v, w)
```

Each edge relaxes one time. Time is Θ(n + m) after you have the order.

Weights can be negative. A DAG has no cycle, so a negative cycle cannot exist.

**Longest paths on a DAG** use the same walk. You maximize instead of minimize, or you negate weights and run shortest path. Longest simple path on a general graph is a different, hard problem. The DAG case is linear.

Use DAG relaxation for task duration: longest path is the critical path in a project DAG.

Do not replace Dijkstra with DAG relaxation on a graph that can have a cycle.

If s is not the first vertex in the order, vertices before s stay at ∞ unless another source exists. That is correct if you only care about paths from s.

### Questions

#### Theoretical questions

1. Why can you relax each edge only one time on a DAG?
2. What must you compute before the relax loop?
3. Why are negative weights safe on a DAG?
4. How do you get a longest path on a DAG?
5. Why is longest path on a general graph not this algorithm?

#### Easy practical tasks

1. Draw a four-vertex DAG with weights. Write a topological order and the relax sequence.
2. Make a table: "Graph class", "Shortest-path method". Rows: DAG, non-negative, negative allowed.
3. Write five sentences that define DAG relaxation.
4. Mark the critical path idea on a three-task DAG (tasks as vertices).

#### Medium practical tasks

1. Implement topological sort plus one relax pass. Test shortest paths from s.
2. Implement longest path on the same DAG by a max update. Print the two dist arrays.
3. Feed a cyclic graph to the function. Show that you reject it.

#### Advanced practical tasks

1. Build a project chart: tasks, durations on vertices or edges (pick one and document it). Print the critical path length.
2. Write a one-page compare: DAG relaxation versus Dijkstra versus Bellman-Ford. Use one sentence for when you pick each.

---

## Connected components / Union-Find preview

An undirected graph splits into **connected components**. Two vertices are in the same component if a path joins them.

**DFS or BFS method.** Repeat: pick an unvisited vertex s, run DFS or BFS, and assign all reached vertices the same component id. Time is Θ(n + m).

```text
0—1    3—4
   \    |
    2   5

components: {0,1,2} and {3,4,5}
```

**Union-Find preview.** A disjoint-set structure stores a partition of vertices. `find(u)` names the component of u. `union(u, v)` merges the two components when you process edge {u, v}. After you union every edge, `find(u) == find(v)` means u and v are connected. The next topic implements Union-Find. This section only states the need.

Use DFS/BFS when you already have adjacency lists and you need to visit vertices. Use Union-Find when edges arrive as a list and you ask many connectivity queries, or when you run Kruskal.

Directed graphs use **strongly connected components** (SCC). That algorithm is an advanced graph topic. Do not treat a directed DFS tree as an undirected component.

Count of components is the number of BFS/DFS starts (plus isolated vertices, which each start one search).

### Questions

#### Theoretical questions

1. When are two vertices in the same connected component?
2. How do BFS or DFS find all components?
3. What two operations does the Union-Find preview name?
4. When is Union-Find a better fit than a full BFS?
5. Why is SCC not this section?

#### Easy practical tasks

1. Draw two components. Write a component id for each vertex.
2. Make a table: "Method", "Input form". Rows: BFS flood, Union-Find.
3. Write five sentences that define a connected component.
4. Count components in a graph of 5 isolated vertices.

#### Medium practical tasks

1. Implement component ids with BFS. Print the id array and the count.
2. Write a loop `union` on an edge list in words (no path compression yet). Show the partition after each edge on a tiny graph.
3. Explain in six sentences the difference between "visited in one BFS" and `find(u) == find(v)`.

#### Advanced practical tasks

1. Compare component count from BFS with a later Union-Find implementation on the same edge list. They must match. You can finish this after topic 16.
2. Write a study plan of five steps from this preview to Kruskal MST. Do not implement MST here.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe how you pick BFS, DFS, Dijkstra, Bellman-Ford, or DAG relaxation from the weight and cycle facts.
2. How do parent arrays, colors, and dist arrays work together across these algorithms?
3. A teammate says "DFS finds shortest paths". Which facts do you use to correct that sentence?
4. What do topological sort and DAG relaxation share, and what does each add?
5. Why does this topic preview Union-Find instead of finishing Kruskal?

#### Easy practical tasks

1. Write a one-page cheat sheet: BFS, DFS times, cycle tests, topo, Dijkstra, Bellman-Ford, DAG relax, components.
2. Make a decision table: "Weights", "Cycles?", "Algorithm". Five rows.
3. Draw one graph and name which algorithm you would run for four different questions.
4. Write one valid topological order and one BFS dist array on two different drawings.

#### Medium practical tasks

1. Implement BFS path, DFS cycle detect, and Kahn on the same tiny library from topic 14.
2. Implement Dijkstra and Bellman-Ford. On a non-negative graph they must match. On a negative-edge DAG they must match DAG relaxation.
3. Time BFS versus Dijkstra with weights 1 on one sparse graph. Write both times.

#### Advanced practical tasks

1. Build a command-line tool: read a graph file, run a named algorithm, print dist or order or "cycle".
2. Read the next two topics (Union-Find, strings). Write which graph problems you will revisit with Union-Find. Do not implement MST here.
