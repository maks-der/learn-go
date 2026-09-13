# 2. Complexity Basics

## Description

Complexity describes how the cost of an operation grows when the input grows. This topic shows how you read Big-O, Big-Theta, and Big-Omega. You learn best case, average case, and worst case. You also learn common growth classes and the idea of amortized cost. Complete this topic before you compare arrays and lists.

Use one term for each concept. **n** is the size of the input. A **growth class** is a family of functions such as constant or linear. **Locality** is the use of nearby memory cells in a short time.

---

## Big-O, Big-Theta, Big-Omega (practical reading)

Big-O, Big-Theta, and Big-Omega are notations for growth. You use them to compare operations. You do not use them as exact step counts.

**Big-O** is an upper bound. If an operation is O(n), the time does not grow faster than a constant times n, for large n. Big-O can be loose. O(n) is also O(n²). A tight upper bound is more useful.

**Big-Omega** is a lower bound. If an operation is Ω(n), the time does not grow slower than a constant times n, for large n.

**Big-Theta** is a tight bound. If an operation is Θ(n), the time grows like n. The operation is both O(n) and Ω(n).

In daily work you often see Big-O. Many texts write O(n) when they mean Θ(n). Read the text. If the author gives a typical cost, the author often means a tight bound.

```text
find in an unsorted array of n elements
  worst case:  Θ(n)   — you can visit every element
  you may write: O(n) — an upper bound
```

Ignore constant factors in the notation. 2n and 100n are both Θ(n). Do not ignore constants in a real program when n is small. Notation is for growth. Measurement is for a real machine.

Drop lower terms. n² + 3n + 10 is Θ(n²). The n² term dominates for large n.

### Questions

#### Theoretical questions

1. What does Big-O say about growth?
2. What does Big-Omega say about growth?
3. When do you write Big-Theta?
4. Why can O(n) also be O(n²)?
5. Why do you drop lower terms in a growth class?

#### Easy practical tasks

1. Write five sentences that define Big-O, Big-Theta, and Big-Omega. Use only facts from this section.
2. Make a table: "Notation", "Bound type", "One example". Add three rows.
3. Classify 5n + 20 as Θ of one simple function. Write the function.
4. Classify n² + n as Θ of one simple function. Write why you drop the n term.

#### Medium practical tasks

1. A loop visits each of n elements two times. Write Big-O, Big-Omega, and Big-Theta for the visit count.
2. Find three library functions. Write the Big-O that the documentation gives. If the documentation does not give Big-O, write "not stated".
3. Explain in six sentences why 1000n can be slower than n² for small n, and faster for large n.

#### Advanced practical tasks

1. Read a formal definition of Big-O (c, n0). Rewrite the definition in eight short STE sentences.
2. Prove or show with a table that 3n + 7 is O(n) and that n is O(3n + 7). Then write Θ.

---

## Best / average / worst case

One operation can have more than one time cost. The cost depends on the input.

**Best case** is the cheapest input of size n. Example: find a value that sits in the first cell of an array.

**Worst case** is the most expensive input of size n. Example: find a value that is not in the array. You visit every cell.

**Average case** is the expected cost for a typical input. You need a model of typical input. If each position is equally likely, the average find in an unsorted array visits about n/2 cells.

Write the case when you write a bound. "Find is O(n)" is not complete. "Find in an unsorted array is O(n) in the worst case" is complete.

Many structures give the worst-case bound in the documentation. Some structures give the average-case bound. A hash table is often O(1) on average and O(n) in the worst case.

```text
binary search in a sorted array of n elements
  best case:    Θ(1)   — the middle cell is the target
  worst case:   Θ(log n)
  average case: Θ(log n)  for a usual model
```

Do not mix cases. Do not compare the best case of one structure with the worst case of a different structure.

### Questions

#### Theoretical questions

1. What is the best case of an operation?
2. What is the worst case of an operation?
3. What extra fact do you need for an average case?
4. Why is "find is O(n)" not a complete statement?
5. Why must you not mix best case and worst case when you compare two structures?

#### Easy practical tasks

1. For linear search in an array, write one input for best case and one input for worst case.
2. Make a table: "Case", "Meaning". Add three rows.
3. Write the worst-case bound for a scan of n cells when the target is absent.
4. Give one sentence that states a complete bound for insert at the end of a linked list that has a tail pointer.

#### Medium practical tasks

1. Write a linear search in Go or in Python. Mark the best-case return and the worst-case return in comments.
2. Assume each index is equally likely. Write the average number of visits for a successful find in an array of n cells.
3. Find one hash-table document. Copy the average bound and the worst bound. Write them in your words.

#### Advanced practical tasks

1. Build a small timing test for linear search: target at index 0, target at index n-1, target absent. Use n = 100 000. Write the three times.
2. Write a one-page note that explains when a product team uses worst case and when the team uses average case.

---

## Constant, logarithmic, linear, linearithmic, quadratic

These names are common growth classes. Learn the names. Learn one picture for each class.

**Constant** — Θ(1). The cost does not grow with n. Example: read index i in an array.

**Logarithmic** — Θ(log n). The cost grows slowly. Each step throws away a fraction of the input. Example: binary search. In this handbook **log** means log base 2 unless a text says a different base. A change of base is a constant factor.

**Linear** — Θ(n). The cost grows with n. Example: visit each element one time.

**Linearithmic** — Θ(n log n). The cost is n times a log factor. Example: efficient comparison sort.

**Quadratic** — Θ(n²). The cost grows with n times n. Example: a nested loop that compares each pair.

```text
n        1     16      256       4096
1        1      1        1          1
log n    0      4        8         12
n        1     16      256       4096
n log n  0     64     2048      49152
n²       1    256    65536   16777216
```

A cubic class Θ(n³) and an exponential class Θ(2ⁿ) exist. You see them less often in basic structures. Avoid exponential growth in the inner loop of a large n.

When you nest a linear scan inside a linear scan, you often get quadratic cost. When you split the input in half each time, you often get logarithmic cost.

### Questions

#### Theoretical questions

1. What does constant time mean?
2. Why does binary search have logarithmic cost?
3. What loop shape often gives linear cost?
4. What does linearithmic mean?
5. What loop shape often gives quadratic cost?

#### Easy practical tasks

1. Match each class name to one Θ expression: 1, log n, n, n log n, n².
2. Draw a rough sketch of n and n² for n from 1 to 10. Label the axes.
3. Write one real operation for constant time and one real operation for linear time.
4. For n = 1024, write log₂ n as an integer.

#### Medium practical tasks

1. Write a nested pair of loops on an n by n grid. Count the inner-body runs as a function of n.
2. Write binary search on a sorted array. Count the worst-case compares for n = 1, 2, 4, 8, 16.
3. Sort n numbers with a simple pair-swap method and with a library sort. Write the two growth classes.

#### Advanced practical tasks

1. Fill a table for n = 10, 100, 1 000, 10 000 with values of n, n log n, and n². Use integer estimates.
2. Find one algorithm that is Θ(n log n) and one that is Θ(n²) for the same problem. Write when you still use the quadratic method.

---

## Amortized cost (plain idea)

Amortized cost is the average cost of one operation in a **sequence** of operations. The sequence can contain cheap operations and rare expensive operations.

A dynamic array is the usual example. Most append operations write one cell. Sometimes the array is full. Then the program copies all n elements to a larger block. That copy is expensive. If the block doubles, the expensive copies are rare. The average cost of one append in a long sequence is constant.

Amortized cost is not the same as average case. Average case uses a model of random input. Amortized cost uses a sequence of operations on one structure. The input can be chosen by an adversary. The bound still holds for the sequence.

```text
append to a dynamic array that doubles
  cheap append:  Θ(1)
  resize copy:   Θ(n) when size is n
  amortized append in a long sequence: Θ(1)
```

Write "amortized" when you use this idea. "Append is O(1)" is not complete. "Append is amortized O(1)" is complete.

You will use this idea again in the dynamic-array topic. Learn the word now. Learn the picture: many cheap steps pay for one expensive step.

### Questions

#### Theoretical questions

1. What is amortized cost?
2. How is amortized cost different from average case?
3. Why is a rare expensive resize still compatible with a cheap amortized bound?
4. Why must you write the word "amortized" in the bound?
5. What sequence is the usual example for amortized cost?

#### Easy practical tasks

1. Write five sentences that explain amortized append. Use only facts from this section.
2. Make a table: "Event", "Cost". Rows: cheap append, resize at size n.
3. A structure does 7 cheap steps and 1 step of cost 8. Write the amortized cost of one step in that sequence of 8.
4. Label this sentence as complete or not complete: "Push is O(1)". Fix the sentence if you can.

#### Medium practical tasks

1. Simulate a dynamic array that starts at capacity 1 and doubles. List the sizes where a copy occurs, from 1 to 32.
2. Sum the copy costs when you grow from 1 to 16 by doubling. Compare that sum with 16 times a constant.
3. Find the word "amortized" in one language document (Go slice, Java ArrayList, or Python list). Copy one sentence. Write it in your words.

#### Advanced practical tasks

1. Write a tiny dynamic array that doubles. Count assignments during 1 024 appends. Divide by 1 024. Write the amortized assignment count.
2. Explain in one page why a +1 growth strategy does not give amortized constant append. Use a cost sum.

---

## Memory locality (cache-friendly arrays vs pointer chasing)

The processor reads memory in blocks. A block in the cache is a **cache line**. Nearby addresses often sit in the same cache line.

**Locality** is good when the program reads nearby cells in a short time. An array stores elements in contiguous memory. A scan of an array has good locality. The cache helps.

A linked list stores nodes at addresses that the allocator selects. The nodes are often not nearby. A walk of a list jumps from pointer to pointer. That walk is **pointer chasing**. Pointer chasing has poor locality. The cache helps less.

```text
array scan:   cell 0, cell 1, cell 2, ...     — nearby addresses
list walk:    node A → node B → node C        — addresses can be far
```

Big-O can be the same for two walks. Both walks are Θ(n). The array walk is often faster on a real machine because of locality.

Do not ignore locality when n is large. Do not ignore Big-O when n is large. Use both ideas. Big-O is the growth class. Locality is a constant-factor and cache effect.

A structure that uses an array as the inner store is often cache-friendly. A structure that uses many small nodes is often not cache-friendly.

### Questions

#### Theoretical questions

1. What is a cache line in practical words?
2. What does good locality mean?
3. What is pointer chasing?
4. Why can two Θ(n) walks have different real times?
5. Why is an array scan often cache-friendly?

#### Easy practical tasks

1. Draw eight array cells in one row. Draw four list nodes with arrows that jump. Label which picture has better locality.
2. Write four sentences that contrast an array scan and a list walk.
3. Make a table: "Layout", "Locality". Rows: contiguous array, linked nodes.
4. List two structures that usually sit in contiguous memory. List two that usually use pointers.

#### Medium practical tasks

1. Write an array sum and a linked-list sum for the same n integers in Go or in Python. Do not time yet. Mark the memory access pattern in comments.
2. Read a short note on CPU cache (one page). Write five STE sentences about cache lines and arrays.
3. Explain why a two-dimensional array that you scan by rows can be faster than a scan by columns on a row-major layout.

#### Advanced practical tasks

1. Time an array sum and a linked-list sum for n = 1 000 000 integers if your machine allows it. Write the two times and one cause.
2. Write a one-page report: same Big-O, different locality. Use your timing or a published benchmark. Cite the source.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do you write a complete complexity statement for one operation? Name notation, case, and n.
2. When do you use amortized cost instead of a single-operation worst case?
3. Why does a change of log base not change the growth class?
4. How do growth class and locality work together when you select an array or a list?
5. A teammate writes only "O(n)" on a slide. Which facts do you add?

#### Easy practical tasks

1. Write a one-page cheat sheet: Big-O, Big-Theta, Big-Omega, best, average, worst, amortized, locality, and the five growth classes.
2. Classify these as a growth class: read a[i], binary search, full scan, merge sort, pair loop.
3. Draw a number line of growth from slow growth to fast growth. Place 1, log n, n, n log n, n².
4. For a list of n tasks, write one sentence each for worst-case scan and for amortized append.

#### Medium practical tasks

1. Write a short script that counts inner-loop runs for a linear scan and for a nested pair loop. Print the counts for n = 10, 20, 40.
2. Take one structure from a later topic (stack, queue, or hash table). Write best, worst, and average for one operation. Mark unknown if you do not know yet.
3. Explain in eight sentences why doubling a block gives a better amortized story than growth by one cell.

#### Advanced practical tasks

1. Measure three algorithms on the same machine for increasing n. Fit each run to a growth class. Write the method that you used.
2. Read the complexity notes of one standard collection. Rewrite every bound in the complete form that this topic teaches.
