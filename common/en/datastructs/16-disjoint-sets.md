# 16. Disjoint Sets (Union-Find)

## Description

A disjoint-set structure (Union-Find) stores a partition of n elements. This topic explains the parent array, union by rank, path compression, Kruskal as a client, and connectivity queries. Complete this topic after arrays and the graph component preview.

Use one term for each concept. A **set** here is one block of the partition. **Find** returns the representative of an element. **Union** merges two sets. The **representative** (root) is the parent that names the set.

---

## Parent array

Store n elements as integers 0 .. n−1. The **parent** array has length n. parent[i] is the parent of i. If parent[i] = i, then i is a **root**. Each root names one set.

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
  if ra == rb: return   // already one set
  parent[rb] = ra
```

```text
union(0, 1) then union(1, 2)

0 ← 1 ← 2     or a star, depending on attach order
parent example: 0, 0, 1, 3, 4   // if 1→0 and 2→1
```

The structure is a forest of trees. It is not a BST. There is no key order. The only links are parent links.

A bad attach order can make a line. Find then costs Θ(n). The next two sections stop that decay.

Space is Θ(n) for the parent array. You can add a rank or size array.

Do not store the full member list in each set unless you need to iterate members. Find and union do not need that list.

### Questions

#### Theoretical questions

1. What does parent[i] = i mean?
2. How does find walk to the representative?
3. What does union do to two roots?
4. Why can a line of parents make find slow?
5. Why is this forest not a search tree?

#### Easy practical tasks

1. Draw parent [0, 0, 1, 3, 4] as trees. Circle the roots.
2. Make a table: "Index", "Parent", "Is root?". Five rows.
3. Write five sentences that define the parent array.
4. Start from five singletons. Write parent after union(2, 3) if you attach 3 to 2.

#### Medium practical tasks

1. Implement init, find (no compression), and union (no rank). Print parent after each union.
2. Force a line: union(0,1), union(1,2), union(2,3) with attach of the second root to the first. Count find(3) steps.
3. Explain in six sentences when union does nothing.

#### Advanced practical tasks

1. Write a function that lists all members of the set of x by a scan of all i with find(i) == find(x). Time it for n = 10 000. Write why this scan is Θ(n).
2. Document an API: MakeSet, Find, Union. Hide the parent array.

---

## Union by rank

**Union by rank** attaches the tree with the smaller rank under the tree with the larger rank. Rank is an integer on each root. Rank is not always the exact height after you add path compression. Without compression, rank equals height.

Rules:

1. Init rank[i] = 0 for every i.
2. Find the two roots ra and rb.
3. If ra == rb, stop.
4. If rank[ra] < rank[rb], set parent[ra] = rb.
5. If rank[ra] > rank[rb], set parent[rb] = ra.
6. If the ranks are equal, attach one to the other and increase the winner rank by 1.

```text
unionByRank(a, b):
  ra = find(a)
  rb = find(b)
  if ra == rb: return
  if rank[ra] < rank[rb]:
      parent[ra] = rb
  else if rank[ra] > rank[rb]:
      parent[rb] = ra
  else:
      parent[rb] = ra
      rank[ra] = rank[ra] + 1
```

**Union by size** is a close variant. You store the number of elements in the tree. You attach the smaller tree under the larger root. You add the sizes.

Both methods keep the tree shallow. Height is O(log n) if you union by rank or size and you do not compress paths.

```text
equal ranks 0 and 0
attach 1 under 0, rank[0] becomes 1
```

Do not increase rank when the ranks are different. Only the equal-rank case grows rank.

Rank of a non-root is unused. You can leave the old number. Only roots need a correct rank.

### Questions

#### Theoretical questions

1. What does rank store on a root?
2. Which tree becomes the child when ranks differ?
3. When do you increase rank by 1?
4. How is union by size different from union by rank?
5. What height class do you get without path compression if you union by rank?

#### Easy practical tasks

1. Start five singletons. Union 0 with 1 (equal ranks). Write parent and rank.
2. Make a table: "Case", "What you attach". Three rows: less, greater, equal.
3. Write five sentences that define union by rank.
4. Draw two trees of rank 1 and rank 0 before and after union.

#### Medium practical tasks

1. Implement union by rank with find that does not compress. Print rank of roots after a sequence of unions.
2. Implement union by size on the same tests. Compare the two parent arrays. Both must keep one partition.
3. Explain in six sentences why you must not add 1 to rank when ranks differ.

#### Advanced practical tasks

1. Build a line of n unions without rank (always attach b to a). Then build the same n unions with rank. Compare find depth of the last element.
2. Prove in eight STE sentences that equal-rank union increases height by at most 1.

---

## Path compression

**Path compression** changes parent links during find. After find(x), every node on the path from x to the old root points to the root. The next find on that path is short.

```text
find(x):
  if parent[x] != x:
      parent[x] = find(parent[x])   // recursive compression
  return parent[x]
```

Iterative form: walk to the root, then walk again and set each parent to the root.

```text
before find(3):  3 → 2 → 1 → 0
after:           3 → 0, 2 → 0, 1 → 0
```

Path compression does not merge two sets. It only flattens one tree.

**Together with union by rank (or size),** find and union are almost O(1) amortized. The exact function is the inverse Ackermann function α(n). α(n) is at most 4 for any size that you will use. This handbook says **almost O(1)** after you name α(n) one time.

You can compress without rank. You can rank without compress. The pair is the standard.

Do not compress on a structure that must keep a meaningful tree shape for another algorithm. Union-Find trees have no extra meaning.

Recursive find can overflow the call stack on a very deep chain before the first compression. For a first implementation, n is small. For a large n with a bad forest, use an iterative find.

### Questions

#### Theoretical questions

1. What links does path compression change?
2. Does compression merge two sets?
3. What is the amortized cost class with rank plus compression?
4. What is α(n) in one sentence for practice?
5. Why can a recursive find fail on a deep chain?

#### Easy practical tasks

1. Draw the 3→2→1→0 chain before and after find(3).
2. Make a table: "Technique", "What it improves". Rows: union by rank, path compression.
3. Write five sentences that define path compression.
4. Write the recursive find in four lines of pseudocode.

#### Medium practical tasks

1. Implement find with compression and union by rank. After find of a deep node, print parent of each node on the old path.
2. Implement iterative two-pass compression. Match the roots of the recursive find.
3. Time 100 000 mixed find/union operations. Write the time. Do not claim a proof from one number.

#### Advanced practical tasks

1. Disable compression and keep rank. Then enable both. Compare total parent walks on the same operation log.
2. Read one statement of the inverse Ackermann bound. Write ten STE sentences. Cite the source. Do not copy the proof.

---

## Kruskal MST uses this

A **minimum spanning tree (MST)** of a connected undirected weighted graph is a set of n − 1 edges. The set joins all vertices. The sum of weights is minimum.

**Kruskal:**

1. Sort all edges by increasing weight.
2. Init Union-Find on n vertices.
3. For each edge {u, v} in that order: if find(u) != find(v), add the edge to the MST and union(u, v).
4. Stop when you have n − 1 edges, or report that the graph is not connected.

```text
edges sorted:  (0,1,1), (1,2,2), (0,2,3)
add (0,1), add (1,2), skip (0,2)  // same set
```

Union-Find answers "does this edge close a cycle?" in almost O(1). If u and v are already in one set, the edge is a cycle in the forest and you skip it.

You will see Prim in the advanced graph topic. Prim grows one tree with a heap. Kruskal grows a forest with Union-Find. Both compute an MST if the graph is connected.

Time of Kruskal is dominated by the sort: O(m log m), plus almost O(m) Union-Find work.

Do not use Kruskal on a directed graph as a first MST. MST is an undirected idea in this handbook.

If the graph has several components, Kruskal produces a **minimum spanning forest**. You get one tree per component.

### Questions

#### Theoretical questions

1. What does an MST contain?
2. Why does Kruskal sort edges?
3. When do you skip an edge?
4. How does Union-Find detect a cycle here?
5. What is a minimum spanning forest?

#### Easy practical tasks

1. Run Kruskal by hand on a triangle with weights 1, 2, 3. Write the two MST edges.
2. Make a table: "Step", "find(u)==find(v)?", "Action". Three edges.
3. Write five sentences that define Kruskal as a client of Union-Find.
4. Draw a skipped edge that would make a cycle.

#### Medium practical tasks

1. Implement Kruskal with your Union-Find and a sorted edge list. Test a connected graph of 5 vertices.
2. Remove edges until the graph is disconnected. Show that you get fewer than n − 1 MST edges.
3. Explain in six sentences why the heaviest edge of a cycle is not needed.

#### Advanced practical tasks

1. Compare your MST weight with a second method (Prim or a brute set of spanning trees on n ≤ 6). Weights must match.
2. Write Kruskal as a function that returns the edge list of the forest. Document the disconnected case.

---

## Connectivity queries

A **connectivity query** asks whether u and v are in the same component.

Offline form:

1. Init Union-Find.
2. Union every edge of the graph.
3. Answer find(u) == find(v) for each query.

```text
edges: 0-1, 1-2
query (0,2) → true
query (0,3) → false
```

Online form: edges and queries mix. You union when an edge appears. You find when a query appears. Union-Find supports this if edges only add (no delete).

Delete of an edge is not a standard Union-Find operation. If edges disappear, you need a different structure or a rebuild.

**Dynamic connectivity** with deletes is an advanced topic. This handbook does not require it.

You can store an extra **component size** on each root. Then you answer "how many vertices are with u?" as size[find(u)].

You can assign a **component id** as the root. Roots can change after a union. Do not store an old root as a forever id unless you add a separate id array that you update with care. Prefer find at query time.

BFS also answers connectivity. Union-Find is better when the input is an edge list and you do not need the vertex visit order.

### Questions

#### Theoretical questions

1. How do you answer "are u and v connected?" after you union all edges?
2. What is the difference between offline and online queries here?
3. Why is edge delete not a basic Union-Find operation?
4. How do you answer component size if you store size on roots?
5. Why can you not cache a root as a permanent id without care?

#### Easy practical tasks

1. After unions 0-1 and 2-3, write true or false for queries (0,1), (0,2), (2,3).
2. Make a table: "Question", "Operation". Rows: same component, merge edge, component size.
3. Write five sentences that define a connectivity query.
4. List two systems that ask "same group?" (accounts, network hosts). Use your own words.

#### Medium practical tasks

1. Implement a query loop: commands UNION u v and QUERY u v on stdin or in a test list.
2. Add size on roots. Print size after a chain of unions.
3. Explain in six sentences when you pick BFS instead of Union-Find for connectivity.

#### Advanced practical tasks

1. Process a mixed log of 10 000 unions and queries. Write the answers for a known small suffix by hand and match them.
2. Write a one-page note: what breaks if a later topic asks you to delete an edge. Name rebuild as one valid student answer.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe the path from a partition of n people to Kruskal. Name parent, rank, compression, and skip-on-same-find.
2. How do a line of parents and a star of parents change find cost before and after the two optimizations?
3. A teammate says Union-Find "is a heap". Which facts do you use to correct that sentence?
4. What stays correct if you omit compression but keep rank? What becomes slower?
5. Why do connectivity queries and MST skip use the same find comparison?

#### Easy practical tasks

1. Write a one-page cheat sheet: parent, find, union, rank, size, compression, Kruskal skip, queries.
2. Draw one forest before and after a compressing find.
3. Make a table: "Client", "What it asks find". Rows: Kruskal, connectivity, component size.
4. Write init parent and rank for n = 6 from memory.

#### Medium practical tasks

1. Implement full Union-Find: rank plus compression. Tests: n singletons, a merge of two large sets, 1 000 random unions.
2. Implement Kruskal and a query checker on the same structure. MST endpoints of one tree must all share one find.
3. Compare operation counts with and without compression on one recorded sequence.

#### Advanced practical tasks

1. Implement union by size plus compression. Use it for Kruskal and for a social-group query tool.
2. Read the next topic (strings). Write five sentences on why Union-Find is the wrong tool for substring search.
