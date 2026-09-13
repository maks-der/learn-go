# 13. Heaps and Priority Queues

## Description

A heap is a tree that is stored in an array. A priority queue is an ADT that removes the minimum or the maximum first. This topic explains the array layout, heapify, insert, extract, build-heap, heap sort, and two use-cases. It also names Fibonacci heaps and pairing heaps as later reading. Complete this topic after arrays and trees. Complete graphs before you implement Dijkstra with a heap.

Use one term for each concept. A **binary heap** is a complete binary tree with the heap property. A **min-heap** keeps each parent less than or equal to its children. A **priority queue** is the ADT. The heap is the usual implementation.

---

## Binary heap (array layout)

A **binary heap** is a complete binary tree. Each level is full except the last. The last level fills from the left. You store the tree in an array. You do not store child pointers.

For a 1-based array (index 1 is the root):

- left child of i is 2i
- right child of i is 2i + 1
- parent of i is floor(i / 2)

For a 0-based array (index 0 is the root):

- left child of i is 2i + 1
- right child of i is 2i + 2
- parent of i is floor((i − 1) / 2)

This handbook uses **0-based** indexes in examples unless a line says 1-based.

```text
min-heap values: 1, 3, 2, 7, 4, 5

        1
       / \
      3   2
     / \ /
    7  4 5

array: index  0  1  2  3  4  5
       value  1  3  2  7  4  5
```

The **heap property** is not the BST property. The heap does not keep a full sort in the array. The only guaranteed order is parent versus child. The left child can be larger or smaller than the right child.

A complete shape gives a compact array. There is no hole in the index range 0 .. n−1. Locality is better than a pointer tree.

A **max-heap** uses the opposite parent rule: each parent is greater than or equal to its children. Pick min or max and stay with that rule.

The root is the minimum in a min-heap. Find-min is Θ(1). Find of an arbitrary key is Θ(n) unless you add an extra map.

### Questions

#### Theoretical questions

1. What shape does a binary heap have?
2. What are the 0-based child and parent formulas?
3. How is the heap property different from the BST property?
4. Why does the complete shape fit an array without holes?
5. What is the cost of find-min in a min-heap?

#### Easy practical tasks

1. Draw the tree for array [2, 4, 3, 8, 5]. Write parent of index 4.
2. Make a table: "Index i", "Left", "Right", "Parent" for i = 0, 1, 2.
3. Write five sentences that define a min-heap in an array.
4. Label one pair of siblings that are not in sorted order. State that this is allowed.

#### Medium practical tasks

1. Write functions `left(i)`, `right(i)`, and `parent(i)` for 0-based indexes. Test them on n = 6.
2. Check the min-heap property on the 1–3–2 picture. Write pass or fail for each parent.
3. Explain in six sentences why a heap is not a search tree.

#### Advanced practical tasks

1. Store a heap in a 1-based slice (unused index 0). Rewrite the three formulas. Test the same tree.
2. Write a one-page note: cache behavior of an array heap versus a pointer binary tree of the same n.

---

## Heapify, insert, extract-min/max

**Sift-up** (also called bubble-up or swim) starts at a node and swaps with the parent while the heap property fails. Insert uses sift-up.

**Sift-down** (also called bubble-down, sink, or heapify at a node) starts at a node and swaps with the smaller child (min-heap) while the heap property fails. Extract-min uses sift-down. This handbook uses **sift-up** and **sift-down**. **Heapify** in this section means sift-down at one index.

**Insert** (min-heap):

1. Append the new value at index n.
2. Increase n.
3. Sift-up from the new last index.

Time is Θ(h) and h is Θ(log n).

**Extract-min** (min-heap):

1. Save the root as the result.
2. Move the last value to the root.
3. Decrease n.
4. Sift-down from the root.
5. Return the saved result.

Time is Θ(log n).

```text
siftDown(i)  // min-heap
  while i has a child:
      s = index of smaller child
      if a[i] <= a[s]: stop
      swap a[i] and a[s]
      i = s
```

```go
func parent(i int) int { return (i - 1) / 2 }
func left(i int) int   { return 2*i + 1 }
func right(i int) int  { return 2*i + 2 }
```

Extract-max on a max-heap is the same algorithm with the greater child.

Do not remove the root and shift the whole array. That cost is Θ(n). Use the last element and sift-down.

Decrease-key is a later need (Dijkstra). You sift-up after you decrease a value. You must know the index of that value. A plain array heap does not find the index in O(1). Add a map from key to index if you need decrease-key.

### Questions

#### Theoretical questions

1. When do you use sift-up?
2. When do you use sift-down?
3. What are the extract-min steps?
4. Why is extract-min Θ(log n) and not Θ(n)?
5. Why does decrease-key need an index map in a plain heap?

#### Easy practical tasks

1. Insert 0 into the 1–3–2 heap. Draw the array after sift-up.
2. Extract-min from that original heap. Draw the array after the last value moves to the root and after sift-down.
3. Make a table: "Operation", "Start index", "Direction". Rows: insert, extract-min.
4. Write five sentences that contrast sift-up and sift-down.

#### Medium practical tasks

1. Implement insert and extract-min on a min-heap of integers. Test empty, one element, and the 1–3–2 example.
2. Implement sift-down and show each swap on extract-min of [1, 3, 2, 7, 4, 5].
3. Write extract-max for a max-heap. Test with [9, 5, 8, 1].

#### Advanced practical tasks

1. Implement a priority queue ADT: Push, Pop, Peek, Len. Hide the array. Test Peek after each Push.
2. Add decrease-key with a map from value to index. Document the limit if values are not unique.

---

## Build-heap in linear time

You can build a heap of n elements in Θ(n) time. You do not insert n times. n inserts cost Θ(n log n).

**Build-heap** (Floyd):

1. Put the n values in the array in any order.
2. Find the last parent index: parent(n − 1).
3. For i from that index down to 0, call sift-down(i).

Leaves have no children. Sift-down on a leaf does no work. Most nodes are near the leaves. Their sift-down paths are short. The sum of those path lengths is Θ(n), not Θ(n log n).

```text
buildHeap(a):
  n = length of a
  for i = parent(n - 1) down to 0:
      siftDown(i)
```

```text
n = 6, last parent = parent(5) = 2
sift-down at 2, then 1, then 0
```

After build-heap, the array is a heap. It is not a fully sorted array.

Use build-heap when you already hold all elements. Use repeated insert when elements arrive one by one.

Do not prove the Θ(n) bound in this topic unless you want the sum. Remember the fact: build-heap is linear. Repeated insert is linearithmic.

### Questions

#### Theoretical questions

1. What is the cost of n inserts into an empty heap?
2. What is the cost of Floyd build-heap?
3. Why do you start sift-down at the last parent, not at the last leaf?
4. Why is the total work less than n times log n?
5. When do you choose build-heap instead of n inserts?

#### Easy practical tasks

1. For n = 10, write the last parent index in a 0-based array.
2. Make a table: "Method", "Time class". Rows: n inserts, build-heap.
3. Draw an unordered array of 5 values. Mark the last parent.
4. Write five sentences that define Floyd build-heap.

#### Medium practical tasks

1. Implement build-heap. Check the heap property after the loop.
2. Time n inserts versus build-heap for n = 100 000 if you can. Write the two times.
3. Explain in six sentences why leaves can be skipped.

#### Advanced practical tasks

1. Count sift-down swaps on one random array of n = 1 000. Compare with a bound of about 2n. Write the numbers.
2. Write the Θ(n) sum idea in eight STE sentences (nodes at height h, at most n / 2^{h+1} of them).

---

## Heap sort

**Heap sort** sorts an array in place with a max-heap (for ascending output).

Steps:

1. Build a max-heap on the full array.
2. For end = n − 1 down to 1:
   - Swap index 0 with index end.
   - Treat the heap size as end (the tail is already in final place).
   - Sift-down at index 0 in that smaller heap.

```text
max-heap root is the largest
swap root with last
shrink heap
sift-down
repeat
```

Time is Θ(n log n). Space is Θ(1) extra if you sort in the input array. The algorithm is not stable. Equal keys can change order.

Heap sort is a good exercise. Many libraries use a different sort (quick, merge, or a mix). You still implement heap sort once. You then know extract-max in a loop.

Do not use heap sort as your first lesson on heaps. First implement a priority queue. Then write heap sort as a client of sift-down.

A min-heap plus an extra result array also sorts. That form uses Θ(n) extra space. In-place heap sort uses a max-heap and the tail of the same array.

### Questions

#### Theoretical questions

1. Why does in-place heap sort use a max-heap for ascending order?
2. What does the swap of root and last achieve?
3. What is the time class of heap sort?
4. Is heap sort stable?
5. How does heap sort differ from n extract-min into a new array?

#### Easy practical tasks

1. Write the two-phase outline: build, then n − 1 extract steps.
2. Make a table: "Algorithm", "Extra space", "Stable?". Rows: heap sort, merge sort (from memory).
3. Draw one swap of root and last on a four-element max-heap.
4. Write five sentences that define heap sort.

#### Medium practical tasks

1. Implement in-place heap sort with a max-heap. Test empty, one element, two elements, and [3, 1, 2].
2. Show the array after build-max-heap and after the first swap+sift on [3, 1, 4, 2].
3. Compare your result with the language sort on the same input. They must match as sets. Record if equal keys change order.

#### Advanced practical tasks

1. Instrument swap counts for random n = 10 000. Write the count and n log₂ n.
2. Write a one-page note: when a library would not pick heap sort (locality, stability). Use STE.

---

## Dijkstra and scheduling use-cases

A **priority queue** removes the most urgent element. Urgency is a number: distance, time, or priority.

**Dijkstra** finds shortest paths in a graph with non-negative edge weights. The algorithm always extends the node with the smallest tentative distance. A min-heap stores nodes by that distance. Each extract-min gives the next node. Decrease-key (or a new insert of a better distance) updates a node when a shorter path appears. You will implement Dijkstra in the graph topic. Learn the ADT need here: extract-min and decrease-key (or extra inserts).

**Scheduling.** An operating system or a worker pool can store jobs in a heap. The next job is the one with the smallest time or the highest priority. Insert is a new job. Extract is "run this job now".

```text
priority queue ADT
  insert(element, priority)
  extract_extreme()     // min or max
  peek_extreme()
  optional: decrease_priority(element, new)
```

Other use-cases: Huffman codes (two extract-min, one insert), selection of the k smallest (a max-heap of size k), and event simulation (next event time).

A sorted array also gives extract-min in Θ(1) if you remove from the front. Insert is then Θ(n). A heap gives both insert and extract in Θ(log n). Select the heap when you mix the two operations.

Do not use a heap for find-by-key as in a map. Use a heap when you always need the extreme.

### Questions

#### Theoretical questions

1. What two heap operations does Dijkstra use most?
2. Why must Dijkstra wait for the graph topic for a full implementation?
3. How does a scheduler use extract-extreme?
4. Why is a heap better than a sorted array when inserts and extracts mix?
5. Why is a heap a poor dictionary for find-by-key?

#### Easy practical tasks

1. Write a priority-queue ADT card with four operations.
2. Make a table: "Use-case", "What the priority is". Rows: Dijkstra, job queue, next event.
3. List three jobs with times 5, 1, 3. Write the extract order in a min-heap of times.
4. Write five sentences that define a priority queue without the word "graph".

#### Medium practical tasks

1. Simulate a waiter: insert three orders with priorities. Extract twice. Write the queue after each step.
2. Explain in six sentences how k-smallest uses a max-heap of size k.
3. Write Dijkstra in eight STE sentences as a client of a min-heap. Do not write graph code.

#### Advanced practical tasks

1. Implement a job scheduler: insert jobs (name, time), extract the next job, print the run order.
2. Read a short Dijkstra outline. List each priority-queue call. Mark which calls need decrease-key versus extra insert.

---

## Fibonacci / pairing heaps (awareness only)

A **binary heap** is the structure that you implement. Other heaps exist. This section is awareness only. Do not implement them in this topic.

A **Fibonacci heap** supports decrease-key in amortized O(1) and extract-min in amortized O(log n). Some textbook Dijkstra bounds use that fact. Real systems often use a binary heap or a pairing heap. Fibonacci heaps have a large constant and a complex implementation.

A **pairing heap** is simpler than a Fibonacci heap. It is a tree (or a forest) with a simple meld. In practice it is often fast. The exact theoretical bounds are a research topic.

Other names that you may see: binomial heap, d-ary heap, leftist heap, skew heap. A **d-ary heap** is a useful small variant: more children per node, a shallower tree, a different sift cost. You can read d-ary after a binary heap.

```text
you implement:     binary heap in an array
you only name:     Fibonacci heap, pairing heap
you may read later: d-ary heap
```

When a text says "Dijkstra is O(m + n log n) with a Fibonacci heap", treat that as a bound, not as a homework spec. Your first Dijkstra can use a binary heap.

### Questions

#### Theoretical questions

1. Why is this section awareness only?
2. Which Fibonacci-heap operation is amortized O(1)?
3. Why do many systems not ship a Fibonacci heap?
4. What is a pairing heap in one sentence?
5. What Dijkstra bound should you expect with a binary heap in a first implementation?

#### Easy practical tasks

1. Make a table: "Heap", "You implement now?". Rows: binary, Fibonacci, pairing.
2. Write five sentences that warn a teammate not to start with a Fibonacci heap.
3. Copy the phrase "awareness only" into a note. Add one reason from this section.
4. Name two operations of a binary heap and their time class.

#### Medium practical tasks

1. Find one textbook page that states the Fibonacci Dijkstra bound. Write the bound in STE. Cite the page.
2. Explain in six sentences the difference between a bound and an implementation plan.
3. Read the name "d-ary heap". Write how child indexes change if d = 4 (formulas only).

#### Advanced practical tasks

1. Write a one-page compare sheet: binary versus pairing versus Fibonacci. Use only published facts that you cite. Do not implement.
2. After you finish Dijkstra in a later topic, add one paragraph that names which heap you used and why.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe the path from a priority-queue ADT to array indexes, then to Dijkstra as a client.
2. How do build-heap and heap sort share sift-down but produce different results?
3. A teammate says a heap "keeps the array sorted". Which facts do you use to correct that sentence?
4. What stays Θ(1), what is Θ(log n), and what is Θ(n) on a binary heap?
5. Why does this handbook separate "implement a binary heap" from "name a Fibonacci heap"?

#### Easy practical tasks

1. Write a one-page cheat sheet: layout formulas, heap property, insert, extract, build-heap, heap sort, two use-cases, awareness heaps.
2. Draw one min-heap of six values as a tree and as an array.
3. Make a table: "Need", "Use a heap?". Five rows that include find-by-key and extract-min.
4. Write the 0-based parent and child formulas from memory.

#### Medium practical tasks

1. Implement a min-heap: insert, extract-min, build-heap, peek. Tests: random n = 50 extract order is sorted.
2. Implement heap sort. Compare with your language sort on ten arrays.
3. Time build-heap versus n inserts for one large n. Write the two times and the two time classes.

#### Advanced practical tasks

1. Implement a max-heap and a min-heap with the same sift code and a compare function. Reuse it for heap sort and for a job queue.
2. Write a study plan of five steps from this heap to Dijkstra in the graph topic. Do not implement Dijkstra here.
