# 10. Heaps

## Description

A heap is a tree that you store in an array. A priority queue is an ADT that removes the minimum or the maximum first. This topic explains the array layout, insert, extract, heapify, build-heap in linear time, heap sort, and priority-queue use cases.

Complete this topic after arrays and trees. Complete graphs before you implement Dijkstra with a heap (Topic 12).

Use one term for each concept. A **binary heap** is a complete binary tree with the heap property. A **min-heap** keeps each parent less than or equal to its children. A **max-heap** uses the opposite parent rule. A **priority queue** is the ADT. The heap is the usual implementation. **Heapify** restores the heap property at one node. **Build-heap** builds a heap from an unordered array.

---

## Binary heap in an array

A **binary heap** is a complete binary tree. Each level is full except the last. The last level fills from the left. You store the tree in an array. You do not store child pointers.

For a 0-based array (index 0 is the root):

- left child of i is 2i + 1
- right child of i is 2i + 2
- parent of i is floor((i − 1) / 2)

For a 1-based array (index 1 is the root):

- left child of i is 2i
- right child of i is 2i + 1
- parent of i is floor(i / 2)

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

The root is the minimum in a min-heap. Find-min is Θ(1). Find of an arbitrary key is Θ(n) unless you add an extra map.

Pick min or max and stay with that rule for the whole structure.

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

1. Write functions left(i), right(i), and parent(i) for 0-based indexes. Test them on n = 6.
2. Check the min-heap property on the 1–3–2 picture. Write pass or fail for each parent.
3. Explain in six sentences why a heap is not a search tree.

#### Advanced practical tasks

1. Convert a 1-based formula set to 0-based for the same tree. Write both index tables.
2. Estimate cache benefit: n heap nodes in an array versus n nodes with child pointers.

---

## Insert, extract, heapify

**Insert** appends the new value at index n, then **sifts up** (bubbles up). While the value is smaller than its parent (min-heap), swap with the parent. Time is Θ(log n) because the height is about log2 n.

```text
insert 0 into [1, 3, 2, 7, 4, 5]
append: [1, 3, 2, 7, 4, 5, 0]
sift up: 0 swaps toward the root
```

**Extract-min** removes the root. Move the last element to index 0. Decrease n. Then **sift down** (heapify the root). Compare the node with its children. Swap with the smaller child if that child is smaller than the node. Repeat. Time is Θ(log n).

```text
extract-min from [1, 3, 2, 7, 4, 5]
move 5 to root: [5, 3, 2, 7, 4]
sift down: 5 swaps with 2, then stops or continues
```

**Heapify** (sift down) at index i assumes the subtrees of i are already heaps. It restores the heap property at i. Time is Θ(h_i) where h_i is the height of node i.

Do not extract from an empty heap. Define an error.

Decrease-key and increase-key exist on some heaps. A binary heap needs the index of the key. A map from key to index can help. This topic only requires insert, extract, and heapify.

Push is insert. Pop is extract. Peek is read of index 0.

### Questions

#### Theoretical questions

1. Why does insert append first?
2. What does sift up do?
3. What does extract-min put at the root before sift down?
4. What does heapify assume about the children?
5. Why are insert and extract Θ(log n)?

#### Easy practical tasks

1. Insert 0 into the six-element picture on paper. Write the array after each swap.
2. Extract-min from [1, 3, 2]. Write the array after the move and after sift down.
3. Make a table: "Operation", "First step", "Repair". Rows: insert, extract-min.
4. Write five sentences about insert, extract, and heapify.

#### Medium practical tasks

1. Write the loop conditions for sift up and sift down on a min-heap.
2. Explain in six sentences why you swap with the smaller child, not an arbitrary child.
3. Show a case where sift down swaps twice.

#### Advanced practical tasks

1. Implement insert and extract-min. Test empty, one element, and a sequence of 20 inserts.
2. Implement decrease-key if you store a key-to-index map. Write why the map must update on each swap.

---

## Build-heap in linear time

**Build-heap** turns an unordered array of n elements into a heap in place.

A naive method inserts n elements into an empty heap. Each insert is O(log n). The total is O(n log n).

A faster method:

1. Treat the array as a complete tree (it already is, as an array).
2. Call heapify on each non-leaf, from the last parent down to the root.

The last parent has index floor(n/2) − 1 in a 0-based array.

```text
array: [4, 1, 3, 2, 16, 9, 10, 14, 8, 7]
heapify from the last parent toward index 0
result: a min-heap or a max-heap, depending on the compare
```

Why is the total Θ(n)? Most nodes have small height. About n/2 nodes are leaves (heapify does nothing). About n/4 nodes have height 1. The sum of heights over a complete tree is less than n. The standard proof bounds that sum by a constant times n.

Use build-heap when you already have all elements. Use repeated insert when elements arrive one by one.

Build-heap does not sort the array. It only installs the heap property.

Do not heapify from the root to the leaves. The child-heap assumption would fail.

### Questions

#### Theoretical questions

1. What is the naive build cost with n inserts?
2. Where does the linear build start (which index)?
3. Why do leaves need no heapify?
4. Why is the total time Θ(n), not Θ(n log n)?
5. Does build-heap sort the array?

#### Easy practical tasks

1. For n = 10, write the last parent index in 0-based form.
2. Mark leaves and non-leaves in a complete tree of 7 nodes.
3. Make a table: "Method", "Time class". Rows: n inserts, heapify-from-bottom.
4. Write five sentences that define build-heap.

#### Medium practical tasks

1. Run one bottom-up heapify pass on [4, 1, 3, 2] as a max-heap or min-heap. Write the array after each heapify.
2. Explain in six sentences why you must start at the last parent.
3. Write when you pick n inserts instead of build-heap.

#### Advanced practical tasks

1. Implement build-heap. Count the number of compares. Compare with n inserts for n = 1 000.
2. Read a proof that the sum of heights is O(n). Rewrite the idea in eight STE sentences.

---

## Heap sort

**Heap sort** sorts an array in place with a heap.

For **ascending** order, use a **max-heap**:

1. Build a max-heap on the array.
2. For i from n−1 down to 1: swap a[0] with a[i]. The maximum of the remaining heap sits at i. Reduce the heap size by 1. Heapify at index 0.

```text
max-heap: largest at index 0
swap a[0] with a[n-1]
heap size n-1
heapify root
repeat
array fills from the right with decreasing maxima
```

Time is Θ(n) for build plus (n−1) times Θ(log n) for extract. Total is Θ(n log n). Extra space is Θ(1) besides the array.

Heap sort is not stable in the usual implementation. Equal keys can change order.

Heap sort is a comparison sort. The Θ(n log n) class matches merge sort and typical quick sort. Constants and locality differ. Quick sort is often faster in practice. Heap sort has a strict worst-case Θ(n log n).

Do not forget: after each swap, the heap is only the prefix of length i. The suffix is already sorted.

A min-heap plus an extra array can also sort. That method uses extra space.

### Questions

#### Theoretical questions

1. Why does ascending heap sort use a max-heap?
2. What does the swap with a[i] do?
3. What is the time class of heap sort?
4. What extra space does in-place heap sort use?
5. Is the usual heap sort stable?

#### Easy practical tasks

1. Write the two main phases of heap sort.
2. Draw a max-heap of four values. Show the first swap and the new heap size.
3. Make a table: "Sort", "Worst class", "Extra space". Rows: heap sort, a naive n inserts into a list.
4. Write five sentences that define heap sort.

#### Medium practical tasks

1. Run heap sort by hand on [3, 1, 4, 2]. Write the array after build and after each extract swap.
2. Explain in six sentences why the suffix must not be heapified.
3. Write why heap sort is not stable with one example of two equal keys.

#### Advanced practical tasks

1. Implement in-place max-heap sort. Test empty, one element, duplicates, and reverse-sorted input.
2. Compare times with another Θ(n log n) sort in a language that you know. Write a careful conclusion.

---

## Priority-queue use cases

A **priority queue** is an ADT. The operations are insert (with a priority) and extract of the best priority. Peek of the best is optional. A heap is the usual implementation. Insert and extract are Θ(log n). Peek is Θ(1).

Use cases:

- **Scheduling.** The next job is the one with the highest priority or the earliest deadline.
- **Event simulation.** The next event is the one with the smallest time.
- **Graph algorithms.** Dijkstra and Prim extract the closest unset vertex (Topic 12).
- **Top-k.** Keep a heap of size k. For a min-heap of the current k largest, the root is the smallest of those k. A new value that is larger than the root can replace the root.
- **Merging sorted streams.** A min-heap of the current head of each stream gives the next global minimum.

```text
top-3 largest in a stream
max values live in a min-heap of size 3
root is the smallest of the top 3
a new x > root: extract root, insert x
```

Do not use a heap when you need search by key. Use a map. Do not use a heap when you need sorted walk of all keys. Use a balanced tree or sort once.

A sorted array can implement a priority queue. Insert is Θ(n). Extract-min from the front is Θ(n) if you shift, or extract-max from the end is Θ(1) if you sort descending. The heap is the usual compromise: both insert and extract are logarithmic.

A FIFO queue is not a priority queue. Arrival order is not priority unless you use arrival time as the priority.

### Questions

#### Theoretical questions

1. What operations does a priority queue provide?
2. Why is a heap a good implementation?
3. How do you keep top-k largest with a min-heap of size k?
4. Why is a heap the wrong structure for find-by-key?
5. How is a priority queue different from a FIFO queue?

#### Easy practical tasks

1. List five use cases from this section in one column.
2. Draw a min-heap of size 3 that stores the top-3 largest of the stream 1, 9, 3, 8. Write the heap after each number.
3. Make a table: "ADT", "Who comes out first". Rows: stack, queue, priority queue.
4. Write five sentences about priority-queue use cases.

#### Medium practical tasks

1. Write the steps of a discrete-event loop: extract next event, run it, insert new events.
2. Explain in six sentences the top-k min-heap trick.
3. Write why Dijkstra needs extract-min, not only a FIFO queue (preview of Topic 12).

#### Advanced practical tasks

1. Implement a priority queue and a top-k filter on a list of n numbers. Check the k results against a full sort.
2. Read a short note on pairing heaps or Fibonacci heaps. Write eight STE sentences about why decrease-key cost can matter. Do not implement them.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do complete shape, array indexes, and the heap property work together?
2. When do you use build-heap, and when do you use repeated insert?
3. How does heap sort reuse extract without a second array?
4. Why can siblings be out of order when the parent rule still holds?
5. Which later graph algorithms need a priority queue, and which only need a FIFO queue?

#### Easy practical tasks

1. Write a one-page cheat sheet: 0-based formulas, insert, extract, heapify, build-heap, heap sort, priority queue.
2. Draw one min-heap array and the tree picture side by side.
3. Make a table: "Operation", "Time". Add find-min, insert, extract, build-heap, heap sort.
4. List four defects: extract empty, heapify from the root first, treat heap as BST, use heap for key search.

#### Medium practical tasks

1. Implement a min-heap priority queue with insert, peek, and extract. Write tests.
2. Write a short report: heap sort versus build-heap. State what is sorted and what is not.
3. Design a scheduler ADT on top of a heap. Name the priority field.

#### Advanced practical tasks

1. Implement heap sort and a top-k heap. Share the heapify function.
2. Add a key-to-index map and decrease-key. Write a test that a graph-style relax would need.
