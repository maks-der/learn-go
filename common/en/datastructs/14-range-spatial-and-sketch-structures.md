# 14. Range, Spatial, and Sketch Structures

## Description

This topic surveys structures for range queries, probabilistic skip lists, spatial trees, and Bloom filters. You learn segment trees and Fenwick trees, the skip-list idea, a short survey of quadtrees, k-d trees, and R-trees, and Bloom filters.

Complete this topic after arrays, trees, and hash tables. Treat spatial trees as a survey: definitions and uses, not full production code.

Use one term for each concept. A **range query** returns a value from a contiguous index interval [L, R]. A **point update** changes one index. A **skip list** is a layered linked list with random extra pointers. A **spatial tree** splits space or groups rectangles. A **Bloom filter** is a bit array that can say "possibly in the set" or "definitely not in the set".

---

## Segment tree and Fenwick tree

A **segment tree** is a binary tree over an array a[0 .. n−1]. Each node stores a merge of a contiguous segment. The root stores the merge of the full array. A leaf stores one element.

The usual merge is sum, min, or max. The merge must be associative. You combine two child values into the parent.

```text
array:  1  3  2  4
sums:
            10
          /    \
        4        6
       / \      / \
      1   3    2   4
```

**Query [L, R]** walks O(log n) nodes that together cover [L, R] and no extra index. **Update** at index i walks the path from that leaf to the root and repairs O(log n) nodes.

Time of build is Θ(n). Time of query or point update is Θ(log n). Space is about 4n for a heap-style array of nodes (a safe bound for a power-of-two layout).

```text
query(node, segL, segR, L, R):
  if the node segment is outside [L,R]: return identity
  if the node segment is inside [L,R]: return node.value
  return merge(query(left), query(right))
```

The **identity** is 0 for sum, +∞ for min, −∞ for max.

A segment tree can add **lazy propagation** for range updates. That is a second skill. Implement point update and range sum first.

A segment tree is not a BST of keys. Indexes are the keys. The split is always the mid of the index range.

A **Fenwick tree** (binary indexed tree, BIT) stores prefix sums in an array t[1 .. n]. Indexes are **1-based**. The next index after i uses the least set bit of i.

**Point add** at i updates i, i + lsb(i), i + lsb(i+lsb(i)), ... while the index is ≤ n.

**Prefix sum** of 1 .. i reads i, i − lsb(i), ... until 0.

```text
prefix(R) − prefix(L−1)  =  sum of a[L .. R]
```

Fenwick code is short. It fits invertible merges such as sum. Min without extra structure is not a natural Fenwick operation. A segment tree fits min and max.

Time of Fenwick point add and prefix is Θ(log n). Space is Θ(n).

Do not scan [L, R] in a loop if n is large and you have many queries. That is Θ(n) per query.

### Questions

#### Theoretical questions

1. What does one node of a segment tree store?
2. Why is segment-tree query time Θ(log n)?
3. What is the identity value for sum and for min?
4. What operations does a Fenwick tree support well?
5. How do you get a range sum from two prefixes?

#### Easy practical tasks

1. Draw the sum tree for [1, 3, 2, 4]. Write the cover of query [1, 2] (values 3 and 2).
2. Make a table: "Need", "Segment tree or Fenwick". Rows: range min, range sum, point add.
3. Write five sentences that define both trees.
4. For n = 8, write a safe segment-tree array size of 4n.

#### Medium practical tasks

1. Implement build and range sum with point update on a segment tree. Test the picture.
2. Implement Fenwick prefix and point add. Check range sum against a slow loop.
3. Explain in six sentences the three-way branch: outside, inside, partial.

#### Advanced practical tasks

1. Process many random updates and queries. Check each query against a slow array copy.
2. Read a short lazy-propagation note. Write ten STE sentences. Do not implement lazy until point update tests pass.

---

## Skip list

A **skip list** is a linked list with extra layers. Each node has a **height** (a number of forward pointers). Height is chosen at random when you insert, often so that about half of the nodes that appear at layer k also appear at layer k+1.

Layer 0 is a full sorted list. Higher layers skip nodes. Search starts at the top-left and walks right while the next key is too small, then drops one layer.

```text
layer 2:  1 ----------------> 9
layer 1:  1 ----> 5 --------> 9
layer 0:  1 → 3 → 5 → 7 → 9
```

Expected search, insert, and delete are O(log n) if the random heights are independent and the p-parameter is a constant such as 1/2.

You do not need to rotate or recolor. Insert splices pointers at each layer of the new node. Delete unsplices them.

A skip list can implement an ordered set or ordered map (Topic 6). Walk of layer 0 is sorted order.

Worst case can be linear if every height is 1, or if random choices are unlucky. The probability of a very tall worst case is small for large n. This is a randomized guarantee, not an AVL-style worst-case bound.

Some databases use skip lists in memory (for example, an in-memory index). This topic only needs the idea.

Do not treat a skip list as a hash table. Keys stay ordered. There is no hash.

A **deterministic** skip list exists in theory. The usual teaching structure is randomized.

### Questions

#### Theoretical questions

1. What does layer 0 store?
2. How does search use higher layers?
3. When do you choose a node's height?
4. What is the expected class of search?
5. How is a skip list different from a hash table?

#### Easy practical tasks

1. Draw a three-layer skip list of five keys. Show one search path.
2. Make a table: "Structure", "Balance method". Rows: AVL, skip list.
3. Write five sentences that define a skip list.
4. Write why layer 0 must contain every key.

#### Medium practical tasks

1. Write the search steps: go right, drop down, repeat.
2. Explain in six sentences how insert splices more than one layer.
3. Write why the worst case can still be linear.

#### Advanced practical tasks

1. Implement a skip list with random height. Test sorted walk and find.
2. Read a Redis skiplist note later in Topic 15. For this section, write eight STE sentences about expected height.

---

## Quadtree, k-d tree, R-tree (survey)

These structures store **points** or **rectangles** in the plane or in more dimensions. This section is a survey. Learn what each name means. Do not implement all three.

A **quadtree** splits a square into four smaller squares (NW, NE, SW, SE). A node holds a capacity of points. When the square is full, it splits. Search in a region visits only squares that meet that region.

```text
world square
  NW | NE
  ---+---
  SW | SE
each cell can split again
```

Use a quadtree for points in 2D when the space is a box that you can split evenly. Uneven point clouds can make a deep tree in a dense pocket.

A **k-d tree** (k-dimensional tree) is a binary tree of points. Level 0 splits on x (or axis 0). Level 1 splits on y (axis 1). Axes cycle. Each node holds one point (in a common form). A split is a hyperplane.

Nearest-neighbor search prunes a subtree when the split plane is farther than the best distance so far.

k-d trees fit nearest neighbor and range search in moderate dimension. In high dimension, prune fails more often (distance concentration).

An **R-tree** stores **minimum bounding rectangles** (MBRs). Leaves hold rectangles (or pointers to objects). Internal nodes hold MBRs that cover their children. Rectangles can overlap. Search visits every child whose MBR meets the query box.

```text
query box hits two overlapping MBRs
  you must search both children
```

Databases use R-trees for geographic boxes and for "objects that intersect this region".

```text
need points in a city grid     → quadtree (often)
need nearest neighbor in 2D/3D → k-d tree (often)
need many overlapping boxes    → R-tree (often)
```

Do not use these names for a 1D index range on an array. Use a segment tree or a Fenwick tree.

A grid (fixed cells) is a simpler spatial index. Use a grid when cells have even load.

### Questions

#### Theoretical questions

1. How does a quadtree split space?
2. How does a k-d tree choose the split axis?
3. What does an R-tree node store?
4. Why can R-tree children overlap?
5. When do you use a segment tree instead of a spatial tree?

#### Easy practical tasks

1. Draw one quadtree split of four points in a square.
2. Make a table: "Tree", "What you store", "Typical query". Add three rows.
3. Write five sentences that survey the three trees.
4. Write why a dense pocket can deepen a quadtree.

#### Medium practical tasks

1. Explain in six sentences how k-d nearest neighbor can skip a subtree.
2. Draw two overlapping MBRs and a query that hits both.
3. Write a choice for: taxi positions on a map, 3D molecule points, land-parcel polygons.

#### Advanced practical tasks

1. Implement a tiny quadtree insert and a region count for a small point set. Do not implement R-tree.
2. Read one survey page on R-trees or k-d trees. Write ten STE sentences. Cite the source.

---

## Bloom filter

A **Bloom filter** is a compact, probabilistic set. It stores **no keys**. It stores a bit array of size m and k hash functions.

**Insert** of key x: for each hash hi(x), set bit (hi(x) mod m) to 1.

**Query** of key x: if any of the k bits is 0, x is **definitely not** in the set. If all k bits are 1, x is **possibly** in the set.

```text
bits:  0 1 0 1 1 0 1 0
insert "cat" sets three positions
query "dog": a 0 bit → absent
query "cow": all 1s → maybe (could be a false positive)
```

A **false positive** is a "maybe" for a key that you never inserted. Other keys set the same bits.

There is **no false negative** in a standard Bloom filter: if you inserted x and you did not delete, query(x) is always "maybe" (the bits stay 1).

You cannot list the keys. You cannot delete in the basic form (unless you use a counting Bloom filter).

The false-positive rate depends on n, m, and k. More bits and a good k reduce the rate. A full table of formulas is later reading. The idea: you trade space and a small error for fast membership tests.

Use a Bloom filter to skip a slow lookup (disk, network) when most queries miss. If the filter says no, do not go to disk. If the filter says maybe, do the real lookup.

Do not use a Bloom filter as the only store of passwords or of records that you must retrieve.

Hash functions here are table hashes, not password hashes (Topic 5).

### Questions

#### Theoretical questions

1. What does a Bloom filter store?
2. When can you say a key is definitely absent?
3. What is a false positive?
4. Why is there no false negative in the basic filter?
5. Why can you not list the keys?

#### Easy practical tasks

1. Draw 8 bits. Insert one key that sets bits 1, 4, 6. Query a second key that hits a 0.
2. Make a table: "Answer", "Meaning". Rows: any bit 0, all bits 1.
3. Write five sentences that define a Bloom filter.
4. Write one use: skip a disk read on a miss.

#### Medium practical tasks

1. Explain in six sentences why two keys can share bits and cause a false positive.
2. Write why delete is unsafe if you clear bits that another key still needs.
3. Compare a Bloom filter with a hash set: space, error, list keys.

#### Advanced practical tasks

1. Implement a Bloom filter with two or three integer hashes. Measure false positives on a held-out key set.
2. Read a short note on choosing m and k. Write eight STE sentences. Cite the source.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do segment trees and Fenwick trees answer 1D index ranges that spatial trees do not target?
2. When is a skip list a reasonable ordered map compared with AVL (Topic 9)?
3. How do quadtree, k-d tree, and R-tree split or group space differently?
4. Why is a Bloom filter a sketch and not a set ADT in the Topic 6 sense?
5. Which structures in this topic give exact answers, and which can err?

#### Easy practical tasks

1. Write a one-page cheat sheet: segment tree, Fenwick, skip list, quadtree, k-d, R-tree, Bloom, false positive.
2. Draw one sum tree, one skip-list search, and one 8-bit Bloom array.
3. Make a table: "Need", "Structure". Add six rows from this topic.
4. List four defects: Fenwick for min without extra design, Bloom as only store, R-tree for 1D prefix sums, skip list as a hash map.

#### Medium practical tasks

1. Implement one exact range structure (segment or Fenwick) and one Bloom filter. Write the two contracts.
2. Write a choice guide of one page for a teammate: range sum, nearest point, bounding boxes, membership sketch.
3. Design tests: segment identity, Fenwick 1-based indexes, Bloom never says no after insert.

#### Advanced practical tasks

1. Combine a Bloom filter in front of a hash set: measure skipped lookups on a miss-heavy workload.
2. Read one database page on R-trees or inverted indexes. Map five terms to this survey.
