# 19. Spatial and Multi-Key Structures

## Description

These structures store points, rectangles, or approximate sets. This topic explains quadtrees, octrees, k-d trees, R-trees, Bloom filters, skip lists, and count-min sketches. Complete this topic after trees, hash tables, and lists.

Use one term for each concept. A **point** has coordinates. A **bounding box** is a rectangle (or a box) that contains a set of objects. A **false positive** is a "yes" that can be wrong. A **false negative** is a "no" that is wrong. A Bloom filter can have false positives. It must not have false negatives if you only insert.

---

## Quadtree / octree

A **quadtree** splits a 2-d rectangle into four equal quadrants. Each internal node has four children: NW, NE, SW, SE. A **leaf** stores a small number of points (or one point). When a leaf is too full, you split it into four children and you move the points down.

An **octree** is the same idea in 3-d. Each split makes eight octants.

```text
square
  NW | NE
  ---+---
  SW | SE

too many points in one leaf → four child squares
```

**Range query.** To find points in a query rectangle, you skip a node if its square does not intersect the query. You test points only in leaves that can intersect.

**Nearest neighbor.** You search the quadrant that should contain the target. You then prune squares that cannot beat the best distance so far.

A quadtree is a poor fit when all points sit in one corner. You get a deep chain of splits. A k-d tree or a well-built R-tree can handle some of those sets better.

Do not store a full 2-d array of cells for a sparse point set. A quadtree grows where points exist.

Time of query depends on the spread of points. Worst case can be linear. Average case on a uniform set is better.

A quadtree can also split **regions** (colored maps), not only points. The split rule is the same.

### Questions

#### Theoretical questions

1. How many children does a quadtree internal node have?
2. When do you split a leaf?
3. How does a range query skip work?
4. How is an octree different from a quadtree?
5. Why can clustered points make a deep tree?

#### Easy practical tasks

1. Draw one split of a square that holds five points if the leaf limit is 2.
2. Make a table: "Tree", "Dimensions", "Child count". Rows: quadtree, octree.
3. Write five sentences that define a quadtree.
4. Shade one quadrant that a query rectangle does not hit.

#### Medium practical tasks

1. Implement a point quadtree with a leaf capacity. Insert 20 random points in the unit square.
2. Implement a range query. Check against a slow scan of all points.
3. Explain in six sentences a nearest-neighbor prune.

#### Advanced practical tasks

1. Insert all points on a diagonal. Print the depth. Write why the depth is large.
2. Write a one-page note: quadtree versus a uniform grid of buckets.

---

## k-d tree

A **k-d tree** is a binary tree of points in k dimensions. Each level splits on one axis. A common rule cycles the axis: level 0 splits on x, level 1 on y, then x again (for k = 2). The split value is often the median of the points in that node.

```text
2-d points
root splits on x = median
left:  points with x < median
right: points with x >= median
children split on y
```

**Build.** If you pick the median each time, height is Θ(log n) at build. Build time is Θ(n log n) with linear-time median or Θ(n log² n) with sort at each node.

**Nearest neighbor.** Recurse into the side that contains the query. Then decide if the other side can hold a closer point. The split plane is the test.

**Range search.** Skip a subtree if its region does not intersect the query box.

A k-d tree stores **points**. It is less natural for overlapping rectangles. Use an R-tree for many boxes.

Updates after a static build can unbalance the tree. Rebuild or use a more complex update rule. For learning, build once from a fixed set.

A k-d tree is not a quadtree. A k-d tree is binary. A quadtree is 4-way and splits space in the middle of the square, not always at the data median.

Do not use a k-d tree as a first map from a single integer key. Use a BST or a hash table.

### Questions

#### Theoretical questions

1. What does one level of a k-d tree use as a split axis?
2. Why does a median split keep height small at build?
3. How do you prune the other side in nearest neighbor?
4. Why are overlapping rectangles a poor first use?
5. How is a k-d tree different from a quadtree?

#### Easy practical tasks

1. Draw a 2-d k-d tree of four points. Label the first x split.
2. Make a table: "Structure", "Branching", "Split choice". Rows: k-d, quadtree.
3. Write five sentences that define a k-d tree.
4. Mark the two half-planes of the root split.

#### Medium practical tasks

1. Build a 2-d k-d tree from a list of points. Print in-order of one axis to show the split idea.
2. Implement nearest neighbor for a small set. Check against a scan of all distances.
3. Explain in six sentences why a later insert can hurt balance.

#### Advanced practical tasks

1. Implement a box range query. Compare hits with a slow scan.
2. Time nearest neighbor versus a full scan for n = 5 000 in 2-d. Write both times.

---

## R-tree (GIS)

An **R-tree** stores spatial **objects** (often rectangles). Each node has a **minimum bounding rectangle (MBR)** that covers all children. A leaf stores object ids and their MBRs. An internal node stores child MBRs.

**Search.** To find objects that intersect a query box, you visit a child only if the child MBR intersects the query. Many children can intersect. You may visit more than one path.

**Insert.** You choose a child to hold the new box (least enlargement of MBR is a common rule). If a node overflows, you **split** the entries into two nodes and you repair MBRs up the path.

```text
query box
visit child if MBR ∩ query is not empty
skip child if MBR is disjoint
```

GIS and map databases use R-trees (and variants such as R* and packed R-trees) for "what is in this view" and "what intersects this road".

An R-tree is not a B+ tree of one number. It is a multi-dimensional sibling of the B-tree idea: high fanout, page-sized nodes, balance by split. The key is a box, not a single integer.

Do not implement a full R* split in the first week. Draw MBRs. Implement a tiny in-memory R-tree with a simple split, or study a diagram only.

Worst-case search can visit many nodes if MBRs overlap a lot. Good insert and split reduce overlap.

A point set can use an R-tree if you treat a point as a degenerate box.

### Questions

#### Theoretical questions

1. What does an MBR cover?
2. When do you skip a child during search?
3. Why can search follow more than one child?
4. What happens when a node overflows?
5. How is an R-tree like a B-tree, and how is it different?

#### Easy practical tasks

1. Draw two child MBRs and a query box that hits only one.
2. Make a table: "Structure", "What a key is". Rows: B+ tree, R-tree, k-d tree.
3. Write five sentences that define an R-tree.
4. List two GIS questions that match an R-tree.

#### Medium practical tasks

1. Compute the MBR of a set of rectangles by hand (min x, min y, max x, max y).
2. Write a search walk in pseudocode: intersect test, recurse, collect ids.
3. Explain in six sentences why overlap of MBRs hurts.

#### Advanced practical tasks

1. Implement a tiny R-tree: insert rectangles, search intersect. Use a simple split (for example linear split). Test with a slow scan.
2. Read one R-tree or R* overview. Write ten STE sentences. Cite the source.

---

## Bloom filter (probabilistic set)

A **Bloom filter** is a bit array of m bits plus k hash functions. It answers **maybe in the set** or **definitely not in the set**.

**Insert(x).** For each hash hi(x), set bit hi(x) mod m to 1.

**Contains(x).** If any bit hi(x) mod m is 0, x is not in the set. If all k bits are 1, x may be in the set, or the bits were set by other keys (**false positive**).

```text
bits: 0 1 0 1 1 0 0 1
insert x sets two positions
query y sees a 0 → y is absent
query z sees all 1s → z maybe present
```

There is **no false negative** if you only insert and you never clear bits (except in variants). If the filter says no, the key was not inserted.

You cannot list the keys. You cannot delete a key in the basic filter (a counting Bloom filter is a variant).

Space is small: a few bits per key for a chosen false-positive rate. Use a Bloom filter to skip a slow lookup (disk, network) when the answer is often no.

Do not use a Bloom filter as the only store of records that you must retrieve. You still need the real set or map.

Pick m and k from the expected n and the error rate. Formulas live in standard notes. This handbook requires the idea: more bits and good hashes lower the false-positive rate.

Hash functions must be independent enough. In practice you derive k hashes from two hashes.

### Questions

#### Theoretical questions

1. When can Contains return a wrong yes?
2. When is Contains a certain no?
3. Why can you not list the keys?
4. Why do you still need a real store behind the filter?
5. What does insert do to the bit array?

#### Easy practical tasks

1. Draw m = 8 bits. Insert one key that sets bits 1 and 4. Query a key that hits bit 2 = 0.
2. Make a table: "Answer", "Meaning". Rows: no, maybe yes.
3. Write five sentences that define a Bloom filter.
4. Write one system use: skip a disk read when the filter says no.

#### Medium practical tasks

1. Implement a tiny Bloom filter with two hash functions. Insert 20 words. Query 20 other words. Count false positives.
2. Change m. Show that a larger m can reduce false positives on the same set.
3. Explain in six sentences why delete is unsafe if you clear bits.

#### Advanced practical tasks

1. Measure false-positive rate versus the textbook formula for one (n, m, k). Write both numbers.
2. Write a one-page note: Bloom filter in front of a hash map or a file index. Cite one real system if you can.

---

## Skip list

A **skip list** is a layered linked list. Layer 0 is a sorted list of all keys. Higher layers hold fewer keys and act as express lanes.

**Insert.** You insert into layer 0. You then promote the key to layer 1, 2, ... with a coin flip (probability 1/2 is common) until the coin fails. Search starts at the top layer and drops down when the next key is too large.

```text
layer 2:  1 ----------- 9
layer 1:  1 ---- 5 ---- 9
layer 0:  1  3   5  7   9
```

Expected search, insert, and delete are O(log n). Extra space is expected O(n). Worst case can be linear if every coin is unlucky. Probability of a very bad shape is small.

A skip list is an alternative to a balanced BST. Implementation is often shorter than red-black. Redis sorted sets use a skip list (plus a hash map). That fact is a reading task later.

Forward pointers at each node form an array of next links, one per layer that the node joined.

Do not treat a skip list as a hash table. Keys stay sorted. Range scan is a walk on layer 0 from the first key in range.

Sentinel heads (and sometimes tails) simplify border cases.

The coin flips must be random. A fixed promote pattern can create a bad shape.

### Questions

#### Theoretical questions

1. What does layer 0 store?
2. How does search use a higher layer?
3. How does insert decide the height of a key?
4. What is the expected search class?
5. How is a skip list different from a hash table?

#### Easy practical tasks

1. Draw three layers for keys 1, 3, 5, 7. Promote 5 to layer 1 only.
2. Make a table: "Structure", "Sorted?", "Expected find". Rows: skip list, hash map, AVL.
3. Write five sentences that define a skip list.
4. Mark the path of search for key 7 on your drawing.

#### Medium practical tasks

1. Implement search and insert with coin-flip height. Test a sorted iterate on layer 0.
2. Implement delete. Keep layer links correct. Test find after delete.
3. Explain in six sentences why Redis can use a skip list for a sorted set.

#### Advanced practical tasks

1. Time insert and find versus a language balanced map for n = 50 000. Write both times.
2. Read a Redis skiplist comment or a skip-list paper abstract. Write ten STE sentences. Cite the source.

---

## Count-min sketch (awareness)

A **count-min sketch** estimates frequencies of items in a stream. It uses d rows of counters, each of width w. Each row has its own hash. On increment of key x, you add 1 to one counter in each row. The estimate of the count of x is the **minimum** of those d counters.

```text
rows:
  h1(x) → +1
  h2(x) → +1
  h3(x) → +1
estimate(x) = min of those cells
```

The estimate is never too small (for this simple form). It can be too large when other keys collide into the same cells. That error is a **overestimate**.

This section is **awareness only**. Do not implement a production sketch in this topic.

Use a count-min sketch when n is huge, you cannot store every key, and an approximate count is enough (heavy hitters, rough popularity).

A count-min sketch is not a Bloom filter. A Bloom filter stores membership bits. A count-min sketch stores counters.

A count-min sketch is not a heap. A heap can track exact top-k if you store keys. The sketch tracks approximate counts in small memory.

```text
you implement now:  Bloom filter, skip list, one spatial tree
you only name:      count-min sketch
```

Pick w and d from the error you allow. Formulas live in standard notes.

### Questions

#### Theoretical questions

1. How do you update the sketch for one key?
2. Why do you take the minimum of the d counters?
3. Can the estimate be smaller than the true count in the basic form?
4. Why is this section awareness only?
5. How is the sketch different from a Bloom filter?

#### Easy practical tasks

1. Draw 2 rows of 5 counters. Show +1 in two cells for one key.
2. Make a table: "Structure", "Stores". Rows: Bloom, count-min, hash map of counts.
3. Write five sentences that define a count-min sketch.
4. Write one stream problem that can use an approximate count.

#### Medium practical tasks

1. Find one textbook figure of a count-min sketch. Write the meanings of w and d. Cite the page.
2. Explain in six sentences why collisions only raise the min (basic increment form).
3. List three student errors: treating the sketch as exact, mixing it with Bloom bits, implementing it before a Bloom filter.

#### Advanced practical tasks

1. Write a one-page compare: exact hash map of counts, heap of top-k, count-min sketch. Do not implement the sketch.
2. After you implement a Bloom filter, write one paragraph on when you would add a sketch instead of more bits.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe how you pick quadtree, k-d tree, or R-tree from points versus boxes and from 2-d split style.
2. How do Bloom filters and count-min sketches both use hashes, but answer different questions?
3. A teammate says a skip list "is not a real balanced structure". Which facts do you use in a careful reply?
4. What can be a false positive, and what must not be a false negative, in a basic Bloom filter?
5. Why does this handbook implement a skip list or a Bloom filter before a count-min sketch?

#### Easy practical tasks

1. Write a one-page cheat sheet: quad/oct, k-d, R-tree MBR, Bloom, skip list, count-min awareness.
2. Draw one quad split, one k-d split, and one skip-list tower.
3. Make a table: "Need", "Structure". Six rows.
4. Write MBR, false positive, and layer 0 in one sentence each.

#### Medium practical tasks

1. Implement one spatial tree (quad or k-d) and one skip list. Write one test for each.
2. Implement a Bloom filter in front of a hash set. Count how often you skip the set on negative queries.
3. Classify six real features (map view, autocomplete points, cache filter, sorted leaderboard, packet counts, 3-d voxels) to a structure.

#### Advanced practical tasks

1. Build a tiny "map view": store rectangles in a simple R-tree or a grid, query a viewport, verify with a scan.
2. Read the next topic (advanced graphs). Write five sentences on why a skip list does not replace a graph shortest-path structure.
