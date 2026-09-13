# 21. Memory, Persistence, and Concurrency

## Description

This topic explains how structures share memory, how they keep old versions, how more than one thread uses them, and how they behave when data does not fit in RAM. You learn structural sharing, copy-on-write, concurrent queues and maps, cache-oblivious ideas, and external-memory sort. Complete this topic after arrays, trees, and the persistent segment tree preview.

Use one term for each concept. A **persistent** structure keeps old versions after an update. **Structural sharing** reuses unchanged nodes. **Copy-on-write** copies a page or a node only when a writer changes it. A **lock** lets one thread enter a critical section. **Lock-free** progress uses atomic operations without a mutex (the exact definition is a later course).

---

## Persistent / immutable data structures (structural sharing)

An **immutable** structure does not change in place. An "update" returns a **new root**. The old root still shows the old value.

**Structural sharing** means the new version reuses nodes that did not change. A persistent list push copies the new head and points at the old list. A persistent tree update copies the path from the root to the changed leaf. Unchanged subtrees stay shared.

```text
old list:  a → b → c
push d:    d → a → b → c
           a, b, c are the same nodes
```

```text
old tree root0
update one leaf
new path of O(height) nodes
side children still point into the old tree
```

Space of one update is O(height) for a tree, not O(n), if you share. Time is the same class as a mutable update plus allocation of the path.

Uses: undo, multiple versions, safer sharing between threads if versions are immutable (readers keep an old root).

A persistent structure is not automatically stored on disk. Persistence here means **versions in memory**. Disk durability is a different word in databases. This handbook says **versioned** or **structurally persistent** when the confusion is high. The topic title uses the CS structure sense.

Do not mutate a shared node. Mutation changes every version that points at that node.

Functional languages often use persistent maps and vectors (bit-partitioned tries). Read one picture after you understand path copy.

The persistent segment tree in topic 18 is one example. A persistent stack is a simpler first implementation.

### Questions

#### Theoretical questions

1. What does an immutable update return?
2. What does structural sharing reuse?
3. What is the extra space of one tree update that copies a path?
4. Why must you not mutate a shared node?
5. How is this "persistence" different from a database write to disk?

#### Easy practical tasks

1. Draw a list before and after push with sharing. Mark the new node.
2. Make a table: "Version", "Root", "Shared nodes". Two rows.
3. Write five sentences that define structural sharing.
4. Draw a three-node tree and the new path after you change one leaf.

#### Medium practical tasks

1. Implement a persistent stack: Push returns a new stack, Pop returns a new stack, the old stack still iterates the old values.
2. Implement a persistent binary tree insert that copies the path. Keep the old root. Query both versions.
3. Explain in six sentences why two threads can read two roots without a lock if no node is mutated.

#### Advanced practical tasks

1. Count allocated nodes for 1 000 sequential persistent inserts into a BST. Write the count versus n log n.
2. Read one note on a persistent vector (bit trie). Write ten STE sentences. Cite the source. Do not implement the full vector.

---

## Copy-on-write

**Copy-on-write (COW)** delays a copy until a writer changes a shared object. Readers share the same physical page or node. The first write copies the page, then changes the copy.

Operating systems use COW for `fork`: parent and child share pages until one process writes. Filesystems and some databases use COW for snapshots of pages.

```text
two users share page P
user A writes
system copies P to P'
A writes into P'
user B still reads P
```

COW is a cousin of structural sharing. Sharing is the default. Copy happens at write. The grain can be a memory page (OS) or a tree node (persistent tree).

**COW array** (simple): two structures hold a pointer to the same backing array and a refcount. A write that must not affect the other user copies the array first. This is expensive if you copy the whole array. A persistent vector copies only a path of blocks.

Do not confuse COW with a deep copy on every assignment. A deep copy always copies. COW copies only on write, and only the part that the design names.

Refcounts or garbage collection must free a page or node when no version uses it.

In a single-threaded editor, undo can be COW snapshots of buffers. Full-buffer copy is simple and slow. A rope or a persistent tree is a better share.

When you call `fork` on Unix, the child address space is COW. That fact is an OS example, not a data-structure homework.

### Questions

#### Theoretical questions

1. When does COW copy?
2. What do readers share before a write?
3. How does `fork` use COW in one sentence?
4. How is COW different from a deep copy on every assignment?
5. What must free a page that no version uses?

#### Easy practical tasks

1. Draw page P, then P' after one writer. Label readers.
2. Make a table: "Action", "Copy?". Rows: read, first write, second write by the same owner.
3. Write five sentences that define COW.
4. Name two grains: memory page and tree node.

#### Medium practical tasks

1. Implement a tiny COW array: two handles, write through one handle copies once, the other handle still sees old data.
2. Measure time of a full-array copy versus a persistent path copy on n = 100 000 (use topic 18 ideas). Write both times.
3. Explain in six sentences a filesystem snapshot as COW of blocks.

#### Advanced practical tasks

1. Add a refcount to your COW array. Free the backing store when the last handle is done. Write a test that checks that you do not free too early.
2. Read one OS or filesystem COW page. Write ten STE sentences. Cite the source.

---

## Concurrent queues and maps (conceptual: lock vs lock-free)

Two threads must not write the same mutable node without a rule. A **race** is a bug when the result depends on an unsafe interleaving.

A **mutex** (lock) around the whole structure is the simple rule. One thread holds the lock, does the operation, then releases the lock. The ADT stays correct. Threads wait. A large critical section can be slow.

A **finer lock** (one lock per bucket of a hash map) lets threads touch different buckets at the same time. You must still lock in a fixed order if you need two locks.

**Lock-free** structures use atomic compare-and-swap (CAS) on pointers or indexes. A failed CAS retries. The goal is that some thread completes an operation even if other threads halt (precise progress terms: lock-free, wait-free). Those terms belong to a concurrency course. This handbook needs the idea: no mutex, atomics, retry loops, hard to prove.

```text
locked queue
  lock
  enqueue
  unlock

lock-free queue (idea)
  CAS tail next from null to node
  CAS tail to node
  retry if CAS fails
```

A **concurrent queue** is a common need (work stealing, pipelines). A **concurrent map** is a common need (caches). Many languages give a well-tested type. Prefer the language type in production.

Do not invent a lock-free map as a first project. Implement a mutex-protected map. Then read one standard lock-free queue (Michael–Scott) as a study, not as a first write.

Immutable persistent roots can let many readers run without a lock. A single writer can publish a new root with an atomic pointer store. That is a useful pattern. It is not a full lock-free map.

Memory order (acquire, release) is a later topic. Do not ignore it if you write atomics. Use a language queue first.

### Questions

#### Theoretical questions

1. What is a race in one sentence?
2. What does a mutex around the ADT guarantee?
3. What does a per-bucket lock allow?
4. What does CAS do in a lock-free update?
5. Why does this handbook tell you not to invent a lock-free map first?

#### Easy practical tasks

1. Write five sentences that contrast lock and lock-free.
2. Make a table: "Tool", "Typical use". Rows: one mutex, bucket locks, language concurrent map.
3. Draw two threads and one queue. Mark the critical section.
4. List two language types (for example a Go channel or a Java concurrent map). Write one sentence each.

#### Medium practical tasks

1. Implement a mutex-protected queue. Start two threads that enqueue. Check that the final length matches the total enqueues.
2. Implement a map with one mutex. Do not write atomics.
3. Explain in six sentences how an atomic root pointer publishes a new persistent version.

#### Advanced practical tasks

1. Read a Michael–Scott queue outline. Write ten STE sentences. Do not write a production lock-free queue.
2. Write a one-page note: when a channel or a blocking queue is clearer than a shared lock-free structure.

---

## Cache-oblivious idea (awareness)

A **cache-aware** algorithm uses the cache line size or the block size B as a parameter. You tune the code for B.

A **cache-oblivious** algorithm is written without B. The same code is efficient for every level of the cache if the model holds. The analysis uses an ideal cache that the program does not name.

```text
cache-aware:   tile size = B, you write B in the code
cache-oblivious: recurse on halves, no B in the code
```

Classic examples (awareness): a van Emde Boas layout of a complete tree, cache-oblivious sorting, a cache-oblivious B-tree analogue. You do not implement them in this topic.

The idea that you already use: **scan an array** has good locality. **Pointer chasing** on a linked list has poor locality. A heap in an array is more cache-friendly than a pointer heap. A B+ tree is cache-aware (or disk-aware): the node size is a page.

A cache-oblivious layout tries to get multi-level cache gains without a page constant in the source.

This section is **awareness only**. Do not rewrite your libraries in van Emde Boas order now.

```text
you implement now:  array-friendly layouts that you already know
you only name:      cache-oblivious, van Emde Boas layout
```

When a text says "cache-oblivious", read "no B in the code, analysis still uses B".

### Questions

#### Theoretical questions

1. What does a cache-aware algorithm take as a parameter?
2. What does cache-oblivious mean for the source code?
3. Why is an array scan friendlier than a list walk?
4. Why is a B+ tree cache-aware or disk-aware?
5. Why is this section awareness only?

#### Easy practical tasks

1. Write five sentences that define cache-oblivious versus cache-aware.
2. Make a table: "Layout", "Locality". Rows: array heap, pointer tree, B+ page.
3. Copy "awareness only" into a note. Add one reason.
4. Name two cache levels on a machine (L1, RAM) in one sentence.

#### Medium practical tasks

1. Time a sum of an array versus a sum of a linked list of the same n. Write both times.
2. Find one van Emde Boas layout picture. Write what "recursive halves of the height" means in four STE sentences. Cite the page.
3. Explain in six sentences why a program can be fast on L1 and L2 without a B constant.

#### Advanced practical tasks

1. Write a one-page note: cache-aware tiles in matrix multiply versus cache-oblivious divide and conquer. Cite one source. Do not implement both.
2. After you implement external sort, write one paragraph on why that algorithm is cache-aware or disk-aware (it names a block size).

---

## External-memory algorithms (sort that does not fit RAM)

**External memory** means the data lives on disk or SSD. RAM holds only a part. The costly step is a **block transfer** between disk and RAM, not a CPU compare.

A **block** has size B. RAM holds M words. The input has N records. The number of transfers is the main cost.

**External merge sort** (idea):

1. Read chunks of about M / size(record) records. Sort each chunk in RAM. Write each sorted **run** to disk.
2. Merge k runs at a time. k is limited by how many blocks of size B fit in RAM (a k-way merge).
3. Repeat merge passes until one sorted file remains.

```text
pass 0:  many sorted runs of length ~ M
pass 1:  merge groups of k runs
...
last:    one run of length N
```

This is the same merge idea as in-memory merge sort. The difference is explicit runs and block I/O.

Do not load all N records into RAM if N does not fit. The program will page at random and become very slow, or it will fail.

A B+ tree is an external-memory **search** structure. External sort is an external-memory **sort**. Databases use both.

Replacement selection can make runs longer than M. That is an optional improvement.

Count I/Os in a model: scan of N records costs about N / B transfers. A good external sort costs about (N / B) log_{M/B}(N / B) transfers. Remember the shape of the formula. You do not need a proof here.

Implement a tiny simulation: files as run files, a small M, integer records. Do not claim a production sort.

### Questions

#### Theoretical questions

1. What is the costly operation in the external-memory model?
2. What is a run in external merge sort?
3. Why does k-way merge depend on M and B?
4. Why must you not load all N records if they do not fit?
5. How does a B+ tree differ from external sort in role?

#### Easy practical tasks

1. Write the three-step outline of external merge sort.
2. Make a table: "Symbol", "Meaning". Rows: N, M, B.
3. Write five sentences that define a run.
4. For N = 100, M = 20, write a possible number of first-pass runs.

#### Medium practical tasks

1. Implement a two-file merge of two sorted integer files. Write a sorted output file.
2. Implement pass 0: split a large file into sorted run files of a fixed RAM limit.
3. Explain in six sentences a k-way merge with a min-heap of k heads.

#### Advanced practical tasks

1. Chain run creation and k-way merges until one file remains. Test on integers that do not fit your chosen M.
2. Write a one-page note: I/O formula shape, why SSD still cares about sequential runs, how a database sort uses this idea. Cite one source.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe the path from a mutable array to a versioned tree, then to disk runs that do not fit RAM.
2. How do structural sharing and COW both avoid a full copy, but at a different grain?
3. A teammate says "lock-free is always faster". Which facts do you use in a careful reply?
4. What does cache-oblivious refuse to put in the source, and what does external sort name on purpose?
5. Why can many readers share an immutable root while a writer publishes a new root?

#### Easy practical tasks

1. Write a one-page cheat sheet: sharing, COW, mutex versus CAS, cache-oblivious awareness, external sort runs.
2. Draw one shared list, one COW page, and one disk merge of two runs.
3. Make a table: "Problem", "Tool". Five rows.
4. Write N, M, B and "new root" from memory in one sentence each.

#### Medium practical tasks

1. Implement a persistent stack and a mutex queue. Write tests for versions and for two threads.
2. Implement a small external merge of run files. Check the output against an in-memory sort of the same integers.
3. Time an array scan versus a pointer walk. Write both times and the locality reason.

#### Advanced practical tasks

1. Write a design of a tiny versioned key-value store in RAM (persistent tree) plus a note on when you would flush runs to disk. Do not build a database.
2. Read the last topic (practice). Write five structures from this path that you will implement again from memory.
