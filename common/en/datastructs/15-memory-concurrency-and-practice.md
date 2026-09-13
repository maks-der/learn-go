# 15. Memory, Concurrency, and Practice

## Description

This topic connects data structures to memory, shared-memory programs, and data that does not fit in RAM. You learn persistent and immutable structures, concurrent queues and maps at a conceptual level (lock versus lock-free), external-memory sort, a required implement list, and how to read one real implementation.

Complete this topic last. Use Topics 1 to 14 when you implement and when you read production code.

Use one term for each concept. A **persistent** structure keeps old versions after an update. An **immutable** value does not change in place. A **lock** is a mutual-exclusion mechanism. **Lock-free** is a progress property of some concurrent algorithms. **External memory** is storage that is not the main RAM (disk or similar). This topic is conceptual for concurrency. Do not treat it as a guide to write lock-free code for production.

---

## Persistent / immutable structures

An **in-place** update changes a node or an array cell. The old value is gone.

An **immutable** structure never changes a published object. An "update" builds a new object. The old object stays valid.

A **persistent** structure lets you use more than one version after updates. **Partial persistence** lets you read all old versions and update only the newest. **Full persistence** lets you update old versions (branching history). This handbook focuses on the common case: immutable trees that share structure.

```text
version 0:      4
               / \
              2   6

insert 5 → version 1 copies the path 4–6
              4'
             / \
            2   6'
               /
              5
node 2 is shared
```

**Structural sharing** copies O(height) nodes on a tree update. The rest of the tree is reused. Time and extra space of one update are O(log n) for a balanced tree, not O(n).

A full array copy is O(n). A persistent array can use a tree of blocks so that a point update copies a path of blocks.

Why use persistence:

- undo and time travel
- multiple threads can read a version with no lock if that version never changes
- functional programs pass new maps and keep old maps

Cost: more allocations, more pointer chasing, more garbage. An in-place array still wins for a tight numeric loop.

Do not mutate a node that other versions still share. That mutation would change history.

A language can enforce immutability. The idea is independent of the language: new version, shared unchanged parts.

### Questions

#### Theoretical questions

1. What does an in-place update do to the old value?
2. What does structural sharing copy on a tree insert?
3. What is partial persistence?
4. Why can many readers use one version without a lock?
5. Why can a persistent update cost more allocations than an in-place update?

#### Easy practical tasks

1. Draw two versions of a three-node tree after one insert. Mark the shared node.
2. Make a table: "Style", "Old version", "Typical extra work". Rows: in-place array, persistent tree.
3. Write five sentences that define persistent and immutable structures.
4. Write why you must not mutate a shared node.

#### Medium practical tasks

1. Explain in six sentences how undo is a pointer to an old root.
2. Compare a full array copy with a path copy in a binary tree of n nodes.
3. Write one program type that needs history (for example, a document). Name the versions.

#### Advanced practical tasks

1. Implement a persistent BST insert that returns a new root and reuses unchanged nodes. Test that the old root still finds the old keys.
2. Read a short note on persistent segment trees or ropes. Write ten STE sentences.

---

## Concurrent queues and maps (lock vs lock-free, conceptual)

More than one thread can call operations on one structure. Without a rule, races destroy the structure: lost updates, broken links, reads of half-written fields.

A **lock** (mutex) lets one thread run a critical section. Other threads wait.

```text
lock
  enqueue or map put
unlock
```

A **coarse lock** covers the whole structure. It is simple. It limits parallelism.

**Fine locks** cover buckets or nodes. More than one thread can work on different parts. Deadlock and lock order become your problem. Hash-map bucket locks are a common picture.

**Lock-free** algorithms use atomic read-modify-write (for example compare-and-swap) so that at least one thread makes progress even if others pause. They are hard to write and hard to prove. They still need a memory model. This topic does not give a lock-free recipe.

```text
conceptual contrast
  lock:           "I own the structure for this operation"
  lock-free:      "I retry an atomic step; no thread holds a mutex"
```

A **concurrent queue** is a common building block (work queue, pipeline). A locked ring or list is a valid design. A lock-free queue exists in the literature. Prefer a library implementation.

A **concurrent map** often uses a lock per shard (a group of buckets) or a lock-free table. Resize is the hard moment: many threads and a growing bucket array.

Rules you can use now:

- document which thread owns a structure
- if more than one thread writes, use a lock or a known concurrent type
- do not iterate a map while another thread resizes it (Topic 6 invalidation plus races)
- do not write your own lock-free map for a first project

Progress words (wait-free, obstruction-free) exist. You do not need them to lock a queue correctly.

This section is defensive. Use it to avoid data races. Do not use it to attack a system.

### Questions

#### Theoretical questions

1. What can a race do to a linked list?
2. What is a coarse lock?
3. What is a shard lock on a map?
4. What does lock-free mean at a high level?
5. Why is resize a hard moment for a concurrent map?

#### Easy practical tasks

1. Write five sentences that contrast lock and lock-free at the idea level.
2. Make a table: "Design", "Who waits". Rows: coarse lock, lock-free (conceptual).
3. Draw two threads and one queue. Mark the lock around enqueue.
4. List three rules from this section.

#### Medium practical tasks

1. Explain in six sentences why iteration plus concurrent insert is unsafe without a contract.
2. Write why a library concurrent queue is a better first choice than a hand-written CAS loop.
3. Design a shard plan: 16 locks, hash to a shard. Write one limit of that plan (a hot key).

#### Advanced practical tasks

1. Use a mutex around a simple queue in a language that you know. Write a test with two producers and one consumer. Do not write a lock-free queue.
2. Read documentation of one concurrent map. Write ten STE sentences about its lock or shard model. Cite the source.

---

## External-memory sort (data larger than RAM)

**RAM** holds a limited array. **External memory** (disk, flash) holds the rest. A disk **page** or block is the unit of I/O. One I/O is expensive compared with a compare in RAM (Topic 9).

**External-memory sort** (external merge sort) sorts n records when n does not fit in RAM.

Idea:

1. Read a chunk that fits in RAM. Sort the chunk. Write a **sorted run** to disk. Repeat until all input is in runs.
2. **Merge** k runs at a time into a larger sorted run. Use a k-way merge (a heap of the next record from each run).
3. Repeat merge passes until one sorted file remains.

```text
RAM holds 3 pages
input: 12 pages
pass 0: write 4 sorted runs of 3 pages
pass 1: merge to fewer, longer runs
...
output: 1 sorted file
```

The number of passes depends on n, RAM size, and the merge width k. k is limited by how many run heads you can keep in RAM and by file handles.

A **heap** (Topic 10) selects the next minimum among k run heads. Each popped record reads the next record from that run (buffered by pages).

Do not assume that you can mmap the whole file as a random-access array. Random I/O on disk is slow. Sequential reads and writes of runs are the usual pattern.

A B+ tree (Topic 9) is an external index. External sort is a one-time (or batch) algorithm. You can build a B+ tree from sorted runs.

Count **I/Os**, not only compares, when you judge an external algorithm.

If the data fits in RAM, use an in-memory sort. External sort is for the overflow case.

### Questions

#### Theoretical questions

1. What is a sorted run?
2. Why does a chunk size equal about the RAM you can use?
3. How does a k-way merge use a heap?
4. Why do you prefer sequential I/O?
5. How do you judge cost: I/Os or only compares?

#### Easy practical tasks

1. Write the two phases: make runs, merge runs.
2. Make a table: "Memory", "Method". Rows: fits in RAM, larger than RAM.
3. Draw three runs and one merge step that outputs the next smallest key.
4. Write five sentences that define external-memory sort.

#### Medium practical tasks

1. For 12 pages of data and 3 pages of RAM, write how many runs pass 0 creates.
2. Explain in six sentences why k cannot be arbitrarily large.
3. Write how a B+ tree build can consume a sorted file.

#### Advanced practical tasks

1. Simulate external sort on paper with 8 records and RAM of 3 records. Write the runs and the merges.
2. Read a short note on replacement selection (longer first runs). Write eight STE sentences.

---

## Implement: dynamic array, hash map, BST, heap, BFS/DFS

Implement these five pieces in a language that you know. Do not copy a full solution from this handbook. The handbook has no solutions. Use the earlier topics as contracts.

**1. Dynamic array** (Topic 2). Fields: buffer, length, capacity. Operations: get, set, append, insert at i, delete at i. Grow by a constant factor. Reject out-of-bounds indexes. Iterate only to length.

**2. Hash map** (Topics 5–6). Chaining is enough. Operations: get, put, delete, size. Resize when the load factor passes a threshold. Keys must be hashable and stable. Write tests for collision, missing key, and overwrite.

**3. BST** (Topic 8). Operations: search, insert, delete (three cases), in-order walk. Test a degenerate insert order and a mixed order. You may skip AVL. State that height can be linear.

**4. Heap** (Topic 10). Array binary min-heap. Operations: insert, peek-min, extract-min, build-heap. Use 0-based index formulas. Test empty extract as an error.

**5. BFS and DFS** (Topics 4 and 11). Adjacency list. BFS: queue, dist, parent, unweighted shortest path. DFS: recursion or explicit stack, directed cycle detect or topological sort on a DAG. Test a grid or a small graph with a cycle and without a cycle.

Write tests before you call a structure done. Tests are not solutions. They are checks.

Suggested order matches the path: array, then hash map, then BST, then heap, then BFS/DFS.

You may add Union-Find, a trie, or a segment tree after the five. Those are extra.

Do not spend the first week on a lock-free map or a production B+ tree. Finish the five.

### Questions

#### Theoretical questions

1. Which topic is the contract for the dynamic array?
2. Which collision method is enough for the required hash map?
3. What BST delete cases must you test?
4. Why is find-min on the heap Θ(1)?
5. Which ADT does BFS require?

#### Easy practical tasks

1. Write a checklist of operations for each of the five pieces.
2. Make a table: "Piece", "First test". Add five rows.
3. Write five sentences that state the implement list.
4. List two extras that you will not do until the five pass.

#### Medium practical tasks

1. Write a test plan only (no implementation in this file): bounds on the array, collision on the map, two-child delete on the BST.
2. Explain in six sentences the suggested order.
3. Draw one graph that you will use for both BFS path and DFS cycle tests.

#### Advanced practical tasks

1. Implement the five pieces. Keep them in separate modules. Write a one-page status: pass or fail per test family.
2. Add one extra: Union-Find or a trie. Write why it was extra.

---

## Read one real implementation (Go slice/map, Java HashMap, or Redis skiplist)

Pick **one** real implementation and read it. Use the documentation and, if you can, the source. Write notes in STE. Do not paste large copyrighted source into your notes.

Choose one:

- **Go slice** (dynamic array) and optionally **Go map** (hash table)
- **Java HashMap** (hash table with a published resize and bucket story)
- **Redis skiplist** (skip list used inside sorted sets)

What to extract:

1. The ADT that users see
2. The layout (array, buckets, nodes, layers)
3. Growth or rehash or random height
4. The stated costs (amortized or expected)
5. One detail that the textbook omitted (a sentinel, a treeify of long chains, a max layer)

```text
note card
  name: ...
  ADT: ...
  layout: ...
  grow / height: ...
  cost words: ...
  extra detail: ...
```

Go slice: length and capacity, append can allocate, more than one slice can share a backing array until a grow copies. That share rule is a real-world locality and alias fact.

Java HashMap: load factor, power-of-two capacity, tree bins in some versions for long chains. That last point is worst-case care (Topic 5).

Redis skiplist: ordered by score and member, extra backward pointers for reverse range. That is a skip list in a real system (Topic 14).

Read only as a learner. Follow the license of the source. Do not copy the code into this course as if it were yours.

If you cannot open those three, read one standard-library list or map in another language. The method stays: ADT, layout, grow, cost, one extra.

### Questions

#### Theoretical questions

1. Why do you read one real implementation after the textbook structures?
2. What five facts belong on the note card?
3. Why must you not paste large source into your notes?
4. What extra fact can a Go slice teach about shared backing arrays?
5. What extra fact can a Java HashMap teach about long chains?

#### Easy practical tasks

1. Pick one of the three targets. Write why you picked it.
2. Make an empty note card with the five headings.
3. Write five sentences that define this reading task.
4. List one license or attribution action that you will follow.

#### Medium practical tasks

1. Fill the note card from documentation only (no source). Write what the docs omit.
2. Explain in six sentences one difference between your hash map and the real map that you read.
3. Map Redis skiplist operations to Topic 14 words: layer, expected log, ordered walk.

#### Advanced practical tasks

1. Read a portion of the source (if the license allows). Add three STE bullets that the docs did not state. Do not quote long code.
2. Compare your dynamic array grow factor with the real slice or list grow rule. Write the two rules.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do immutability and a coarse lock both reduce races, and how do their costs differ?
2. Why do B+ trees (Topic 9) and external sort both care about pages?
3. In which order do you implement the five pieces, and which topic each piece uses?
4. What does a real HashMap or skiplist add that a classroom structure omits?
5. When do you refuse to write lock-free code and use a lock or a library instead?

#### Easy practical tasks

1. Write a one-page cheat sheet: persistent, share path, lock, lock-free, run, k-way merge, five implements, one real read.
2. Draw a persistent tree pair, a locked queue, and two disk runs.
3. Make a table: "Task", "Done (yes/no)". Add the five implements and the one reading.
4. List four defects: mutate a shared persistent node, iterate a map during resize on two threads, mmap-sort a huge file at random, skip tests.

#### Medium practical tasks

1. Write a personal lab plan of two weeks: days for each implement, one day for the reading note.
2. Design a tiny pipeline: producer threads, a locked queue, a consumer that inserts into your hash map. Write the lock scope.
3. Write an I/O budget for an external sort of N pages with M pages of RAM (symbols only).

#### Advanced practical tasks

1. Finish the five implements and the note card. Store them in your own repository. This handbook still contains no solutions.
2. Read CLRS, Sedgewick, VisuAlgo, or CP-Algorithms for one structure that you implemented. Write a ten-row term map.
