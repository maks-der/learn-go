# 18. Advanced Trees and Tries

## Description

These structures answer range questions and compress string keys. This topic explains segment trees, Fenwick trees, sparse tables, Cartesian trees, persistent segment trees, and radix / Patricia tries. Complete this topic after arrays, binary trees, and the string trie.

Use one term for each concept. A **range query** returns a value from a contiguous index interval [L, R]. A **point update** changes one index. **Idempotent** merge (such as min) can use a sparse table. **Invertible** merge (such as sum) can use a Fenwick tree.

---

## Segment tree

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

Do not scan [L, R] in a loop if n is large and you have many queries. That is Θ(n) per query.

### Questions

#### Theoretical questions

1. What does one node of a segment tree store?
2. Why is query time Θ(log n)?
3. What is the identity value for sum and for min?
4. How is a segment tree different from a BST?
5. What extra idea does lazy propagation add?

#### Easy practical tasks

1. Draw the sum tree for [1, 3, 2, 4]. Write the cover of query [1, 2] (values 3 and 2).
2. Make a table: "Need", "Identity". Rows: sum, min, max.
3. Write five sentences that define a segment tree.
4. For n = 8, write a safe node-array size of 4n.

#### Medium practical tasks

1. Implement build and range sum with point update. Test the picture.
2. Implement range min on the same array. Test [0, 3] and [2, 2].
3. Explain in six sentences the three-way branch: outside, inside, partial.

#### Advanced practical tasks

1. Process 10 000 random updates and queries. Check each query against a slow loop on a copy of the array.
2. Read a short lazy-propagation note. Write ten STE sentences. Do not implement lazy until the oracle tests pass for point update.

---

## Fenwick tree (BIT)

A **Fenwick tree** (binary indexed tree, BIT) stores prefix sums in an array t[1 .. n]. It uses less code than a segment tree for **prefix sum** and **point add**.

Indexes are **1-based**. The next index after i uses the least set bit of i.

```text
i & -i     // least set bit (two's complement)
add i:     i += i & -i
parent i:  i -= i & -i
```

**Add(i, delta)** adds delta to a[i] (1-based). It updates t at i, i + lsb, and so on.

**PrefixSum(i)** returns a[1] + ... + a[i]. It sums t at i, i − lsb, and so on.

**Range sum [L, R]** is PrefixSum(R) − PrefixSum(L − 1).

```text
add(i, delta):
  while i <= n:
      t[i] += delta
      i += i & -i

sum(i):
  s = 0
  while i > 0:
      s += t[i]
      i -= i & -i
  return s
```

Time of add and prefix sum is Θ(log n). Space is Θ(n).

A Fenwick tree needs an **invertible** operation. Sum works. Min does not work in the basic form because you cannot subtract a min.

Build: add each a[i] once, or use a linear build if you learn it later. n adds are Θ(n log n).

Do not use 0-based i inside the bit loops. Convert 0-based user indexes to 1-based.

A Fenwick tree does not store an explicit binary tree of child pointers. The tree is the bit rule.

### Questions

#### Theoretical questions

1. What does a Fenwick tree store?
2. Why are indexes 1-based in the bit loops?
3. How do you get a range sum from prefix sums?
4. Why is min a poor fit for a basic Fenwick tree?
5. What is the time class of add and of prefix sum?

#### Easy practical tasks

1. For i = 6 (binary 110), write i & -i and the next two add indexes.
2. Make a table: "Structure", "Range min?", "Range sum + point add?". Rows: Fenwick, segment tree.
3. Write five sentences that define a BIT.
4. Write range [2, 5] as two prefix calls.

#### Medium practical tasks

1. Implement add and prefix sum. Test range sums against a slow array.
2. Implement point set as add(i, new − old). Test two sets.
3. Explain in six sentences why L − 1 must use a 1-based L.

#### Advanced practical tasks

1. Compare code size and times of Fenwick versus segment tree on the same sum tests.
2. Read one picture of Fenwick responsibility ranges. Write ten STE sentences that map t[i] to an interval. Cite the source.

---

## Sparse table

A **sparse table** answers **idempotent** range queries on a **static** array. Static means you do not update elements after you build the table.

The usual query is range min (RMQ) or range max. The merge must satisfy merge(x, x) = x so that overlap is safe.

```text
st[k][i] = merge of a[i .. i + 2^k - 1]
st[0][i] = a[i]
st[k][i] = merge(st[k-1][i], st[k-1][i + 2^{k-1}])
```

**Query [L, R]:** let len = R − L + 1. Let k = floor(log2(len)). Return merge(st[k][L], st[k][R − 2^k + 1]). The two intervals cover [L, R] and can overlap.

```text
[L, R] length 6, 2^k = 4
cover [L, L+3] and [R-3, R]
```

Build time is Θ(n log n). Query time is Θ(1) after you have k = log2(len). Space is Θ(n log n).

Precompute `log2[len]` in a table if you need many queries.

Do not use a sparse table when you must update. Rebuild is Θ(n log n). Use a segment tree for updates.

Do not use a sparse table for range sum with the overlap formula. Overlap would add the middle twice. Sum is not idempotent.

A sparse table is a good RMQ tool after you learn the log cover idea.

### Questions

#### Theoretical questions

1. What does st[k][i] store?
2. Why must the merge be idempotent for the overlap query?
3. Why is the array static?
4. What are build time and query time?
5. Why is range sum the wrong first query for a sparse table?

#### Easy practical tasks

1. For [L, R] = [1, 6], write len, k, and the two start indexes.
2. Make a table: "Structure", "Update?", "Min query class". Rows: sparse table, segment tree.
3. Write five sentences that define a sparse table.
4. Draw two overlapping blocks of length 4 on a line of 6 cells.

#### Medium practical tasks

1. Implement RMQ build and query. Test against a slow min loop.
2. Precompute floor(log2(i)) for i = 1 .. n. Use it in query.
3. Explain in six sentences why overlap is safe for min and unsafe for sum.

#### Advanced practical tasks

1. Time 1 000 000 RMQ queries versus a segment tree on a static array. Write both times.
2. Write a one-page note: when you pick sparse table, Fenwick, or segment tree. Use three decision questions.

---

## Cartesian tree (awareness)

A **Cartesian tree** of an array is a binary tree. In-order of the tree is the array order. The heap property on values holds (usually a min-heap: the smallest value is the root). The root is the minimum (or maximum). The left subtree is the Cartesian tree of the left subarray. The right subtree is the Cartesian tree of the right subarray.

```text
array values:  3  2  4  1  5
min root: 1
left of 1: Cartesian(3,2,4)    right of 1: Cartesian(5)
```

This section is **awareness only**. Do not implement a full Cartesian tree project in this topic unless the other structures already work.

Uses that you may read later:

- RMQ can reduce to LCA on the Cartesian tree.
- Some sort and quickselect talks use the same min-split idea.
- Treaps are a randomized BST that look like a Cartesian tree on (key, random priority).

A Cartesian tree is not a segment tree. It is not a BST on the values unless the array is sorted.

Linear construction uses a stack of the right spine. You can read that algorithm later.

```text
you implement now:  segment tree, Fenwick, sparse table
you only name:      Cartesian tree, treap link
```

### Questions

#### Theoretical questions

1. What two properties does a Cartesian tree mix?
2. Where does the minimum of the array sit?
3. Why is this section awareness only?
4. How does RMQ relate to LCA in one sentence?
5. How is a treap related in one sentence?

#### Easy practical tasks

1. Draw a min Cartesian tree for [3, 1, 2].
2. Make a table: "Tree", "Awareness only?". Rows: segment, Cartesian.
3. Write five sentences that define a Cartesian tree.
4. Write the in-order of your drawing. It must be 3, 1, 2.

#### Medium practical tasks

1. Find one diagram of a Cartesian tree. Check in-order and heap property. Cite the page.
2. Explain in six sentences why the tree is unique for distinct values.
3. List two student errors: mixing Cartesian with BST search, mixing Cartesian with segment indexes.

#### Advanced practical tasks

1. Write a one-page note on RMQ via Cartesian tree and LCA. Cite one source. Do not implement LCA here.
2. After you finish a sparse table RMQ, write one paragraph on why you did not need a Cartesian tree.

---

## Persistent segment tree (awareness)

A **persistent** segment tree keeps old versions after an update. An update creates O(log n) **new** nodes on the path that changed. The other nodes stay shared with the previous version. Each version is a root pointer.

```text
version 0: root0
update index i → version 1: root1
root1 path is new
side subtrees still point at version 0 nodes
```

This is **structural sharing**. Topic 21 expands persistence. Here you only need the picture.

Uses: you query a past version, or you build one version per prefix of an array and query the k-th order statistic on a subarray (a common contest pattern).

Space is Θ(n + u log n) for u updates if each update adds O(log n) nodes.

This section is **awareness only**. Do not implement persistence until a normal segment tree is correct.

A persistent tree is not a disk snapshot by itself. It is an in-memory sharing trick. Durability on disk is a different topic.

Do not mutate a shared child. Mutation would change an old version. You copy the node on the path (copy-on-write).

```text
you implement now:  mutable segment tree
you only name:      persistent roots, shared children
```

### Questions

#### Theoretical questions

1. What does one update allocate in a persistent segment tree?
2. What is structural sharing here?
3. Why must you not mutate a shared child?
4. What is the extra space class for u updates?
5. Why is this not the same as a database snapshot on disk?

#### Easy practical tasks

1. Draw root0 and root1 after one leaf update. Mark new nodes and shared nodes.
2. Make a table: "Version", "Root name". Two rows.
3. Write five sentences that define persistence for a segment tree.
4. Copy "awareness only" into a note. Add one implement-later reason.

#### Medium practical tasks

1. Read one picture of path copy. Count new nodes for n = 8. Cite the page.
2. Explain in six sentences how two versions can answer two different range sums.
3. List the order of study: mutable segment tree, then copy-on-write, then persistent roots.

#### Advanced practical tasks

1. Write a one-page design of persistent range sum (types for Node and Version). Do not write the full code.
2. Link this section to topic 21 in five sentences. Use the words structural sharing and copy-on-write.

---

## Radix tree / Patricia trie

A **radix tree** (a **Patricia trie** is a binary radix tree on bit strings) **compresses** a chain of single-child trie nodes into one edge. The edge stores a string (or a bit string), not one character.

```text
trie of "romane" and "romanus" wastes a long one-child chain
radix edge: "roman" then a split at e vs u
```

**Find** compares a slice of the key with the edge label. If the label matches, you continue. If the label differs at a position, you follow a split or you report absent.

A Patricia tree often uses bits. The skip count says how many bits are the same. IP routing tables use this idea.

Compare with a hash map: exact find of a full key is often faster in a hash map. A radix tree keeps **prefix order** and can share space among keys with a long common prefix.

Compare with a B+ tree of strings: a B+ tree is for pages and range scans of keys. A radix tree is a compressed trie. Some databases use both ideas in different places.

Implement a small radix tree after a simple trie. Compress only unary paths. Test insert and find of a few words.

Do not start with a production radix or a concurrent radix. Learn the compress rule first.

Time of find is O(L) in the key length, with fewer nodes than a simple trie.

### Questions

#### Theoretical questions

1. What does an edge store in a radix tree?
2. What chain does compression remove?
3. How does find use an edge label?
4. How is a Patricia trie related to bits?
5. When do you still pick a hash map?

#### Easy practical tasks

1. Draw a simple trie and a radix tree for "car" and "cart".
2. Make a table: "Structure", "Edge label". Rows: trie, radix tree.
3. Write five sentences that define a radix tree.
4. Mark the split node after you insert "cap" into a tree of "car".

#### Medium practical tasks

1. Implement a tiny radix insert/find for ASCII words. Compress unary paths. Test car, cart, cap.
2. Compare node counts with a simple trie on the same words.
3. Explain in six sentences one routing or prefix-table use (high level).

#### Advanced practical tasks

1. Implement delete of one key and keep compression correct. Test that other keys still find.
2. Read one Patricia or radix note. Write ten STE sentences. Cite the source.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe how you pick segment tree, Fenwick, or sparse table from update need and merge type.
2. How do persistence and radix compression both save space, but on different objects?
3. A teammate says a Fenwick tree "is a heap". Which facts do you use to correct that sentence?
4. What must you implement before Cartesian or persistent trees, and why?
5. Why is a radix tree a better follow-up to topic 17 than a suffix tree?

#### Easy practical tasks

1. Write a one-page cheat sheet: segment, Fenwick bits, sparse cover, Cartesian awareness, persistent awareness, radix.
2. Draw one sum segment tree and one radix split.
3. Make a table: "Query", "Best first structure". Five rows.
4. Write identity values and the Fenwick range-sum formula from memory.

#### Medium practical tasks

1. Implement segment tree and Fenwick on the same sum tests. They must match.
2. Implement a sparse table RMQ. Show that a Fenwick cannot replace it for min.
3. Implement a small radix tree and a simple trie. Compare node counts.

#### Advanced practical tasks

1. Build a tiny library: range sum (Fenwick), range min (segment or sparse), prefix search (radix). Document each pick.
2. Write a study plan of five steps from a mutable segment tree to a persistent one in topic 21. Do not implement persistence here.
