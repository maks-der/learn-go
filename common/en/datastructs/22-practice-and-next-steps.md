# 22. Practice and Next Steps

## Description

This topic turns the path into a practice plan. You implement a core set of structures, you solve problems that force a choice, you read one real implementation, and you plan deeper algorithm study. Complete this topic after topics 1 to 21. Use it as a checklist, not as a new theory chapter.

Use one term for each concept. A **core implement** is a structure that you write from memory and test. A **forced choice** is a problem that does not name the structure. A **real implementation** is production source that you read. A **next course** is algorithms beyond this data-structure path.

---

## Implement: dynamic array, stack, queue, hash map, BST, heap, graph BFS/DFS

Write each item in a language that you know. Hide the layout behind an ADT. Test empty, one element, and a larger case. Print costs where the topic asked for height, load factor, or distance.

**Dynamic array.** Manual resize (often ×2). Append, get, set, length. Know amortized append.

**Stack.** Push, pop, peek. Array or list. Evaluate one postfix expression.

**Queue.** Enqueue, dequeue, front. Ring buffer or list. Do not use a stack as a queue.

**Hash map.** Separate chaining is enough. Hash, insert, find, delete, resize at a load factor. Keys must be stable.

**BST.** Insert, search, in-order, height, delete with successor. Also build a degenerate case on purpose. Then keep one balanced tree (AVL) if you completed topic 12.

**Heap.** Array min-heap: insert, extract-min, build-heap. Use it as a priority queue.

**Graph BFS and DFS.** Adjacency list. BFS distances and parent path. DFS colors or timestamps. Detect a directed cycle.

```text
minimum kit
  array+grow   stack   queue
  hash map     BST     heap
  graph: list + BFS + DFS
```

Suggested order (same as the topic list):

1. Dynamic array with manual resize
2. Stack and a postfix evaluator
3. Hash map with chaining
4. BST with in-order print
5. Binary heap priority queue
6. BFS shortest path on a grid
7. Union-Find on connectivity queries
8. Trie for autocomplete
9. Read a B+ tree diagram. Write why a database uses one
10. Segment tree for range sums

Items 7 to 10 extend the minimum kit. Do them after the seven names in the heading.

Do not copy a full library and call it done. Write the operations. Keep the files. You will reread them.

A language slice or map is not a substitute for the implement. Use the language type after you can write the student version.

### Questions

#### Theoretical questions

1. What is a core implement in this handbook?
2. Why do you test the empty case for every ADT?
3. Why is a language map not a substitute for a hash map that you write?
4. What extra BST case must you build on purpose?
5. Which two graph walks are in the minimum kit?

#### Easy practical tasks

1. Make a table: "Structure", "File name", "Done?". Seven rows for the heading list.
2. Write five sentences that define the minimum kit.
3. List the ten suggested practice items in order.
4. Mark which of the ten items are not in the seven-name heading.

#### Medium practical tasks

1. Implement or reopen the seven structures. For each, write three tests that you ran.
2. Add Union-Find and a trie if they are missing. Write one test each.
3. Explain in six sentences why postfix evaluation proves that your stack is a stack.

#### Advanced practical tasks

1. Time get on your dynamic array, find on your hash map, and search on your BST for n = 50 000. Write the three times and the expected classes.
2. Put all seven ADTs in one module. Write a one-page README that names operations and file paths.

---

## Solve problems that force you to pick a structure

A **forced choice** problem states operations and constraints. It does not say "use a heap". You select the ADT and the implementation.

How you choose:

1. Write the operations (insert, delete, find, extract-extreme, range, prefix, connect).
2. Write the constraints (sorted?, weighted?, memory?, many queries?).
3. Pick the ADT. Then pick the implementation. Name the cost class.
4. If two structures fit, pick the simpler one that meets the bound.

```text
examples of a force
  "return the minimum each time"           →  heap or sorted set
  "is x in the set, average fast"          →  hash set
  "iterate keys in order"                  →  balanced tree or skip list
  "shortest unweighted path"               →  BFS
  "range sum with updates"                 →  Fenwick or segment tree
  "are u and v connected after unions"     →  Union-Find
  "words with this prefix"                 →  trie
```

Contest problems and job tasks often hide the structure in a story. Translate the story to operations first.

Do not start from a structure that you want to practice if it fails the operations. Practice that structure on a different problem.

Keep a log: problem name, operations, structure, one cost sentence. The log is the learning product.

Wrong picks are useful. If a list is too slow, write why. Then replace it.

Mix easy problems (one structure) with medium problems (two structures, for example hash map plus heap).

This handbook does not include solutions to problems. It includes a method. Use an online judge or a book only as a source of statements.

### Questions

#### Theoretical questions

1. What is a forced choice problem?
2. What four steps does this section use to pick a structure?
3. Why do you write operations before you name a heap or a map?
4. When do you pick the simpler structure?
5. What does the problem log store?

#### Easy practical tasks

1. Translate these three stories to operations: a printer queue, a phone book, a maze.
2. Make a table: "Story", "ADT", "Implementation". Three rows.
3. Write five sentences that define a forced choice.
4. Add two lines to a practice log template.

#### Medium practical tasks

1. Take five problems from a judge or a book. Fill the log. Do not paste solutions.
2. Find one problem that you first solved with a list. Write why a set or a heap is better.
3. Explain in six sentences a problem that needs two structures.

#### Advanced practical tasks

1. Solve or design eight forced-choice tasks that cover heap, hash, BFS, Union-Find, trie, and a range tree. Write only the pick and the cost, not a solution write-up of code.
2. Review a teammate pick (or your old code). Write four questions that this section would ask before you accept the pick.

---

## Read one real implementation (Go slice/map, Java HashMap, Redis SDS/skiplist)

Read **one** production implementation in depth. Skim is not enough. Trace one operation from the public function to the memory layout.

Pick one track:

- **Go slice.** Growth, capacity, copy on grow, the slice header (pointer, length, capacity).
- **Go map.** Buckets, overflow, hash, grow. The runtime map is not a small file. Read the comments first.
- **Java HashMap.** Table, treeify of long chains, resize. Read the class comment.
- **Redis SDS.** Simple Dynamic String: length, alloc, flags. Why Redis does not use only C strings.
- **Redis skiplist.** Sorted set backing (with a dict). Level, span, backward pointer.

```text
how to read
  1. public operation
  2. struct fields
  3. one happy path (insert or append)
  4. one grow or rebalance path
  5. comments that state invariants
```

Write a one-page note in STE: fields, invariants, one cost fact, one difference from your student version.

Do not read all five in the same week. One deep read beats five skims.

Cite the version or commit if you can. Source moves.

The real code can use bit tricks and unsafe conversions. You do not copy those into a first student file. You record why they exist (speed, compactness).

A second read later (a B-tree in a database, or a language red-black tree) is optional after this one.

### Questions

#### Theoretical questions

1. Why does this section say one implementation, not five?
2. What five read steps does this handbook list?
3. What belongs in the one-page STE note?
4. Why do you cite a version or commit?
5. Why must you not copy unsafe tricks into a first student structure?

#### Easy practical tasks

1. Pick one track. Write the path to the source file on your machine or the URL.
2. Make a table: "Field", "Role" with four rows from that source.
3. Write five sentences that contrast that implementation with your student version.
4. Copy one invariant from a comment. Rewrite it in STE. Do not copy a long block.

#### Medium practical tasks

1. Trace append (slice), put (map/HashMap), or insert (skiplist/SDS grow). Write the step list.
2. Trace a grow or resize. Write what happens to old memory.
3. Explain in six sentences one design choice that your student code omitted.

#### Advanced practical tasks

1. Write a two-page compare: your hash map or skip list versus the real one. Cite the source.
2. Read a second implementation only after the first note exists. Add a half-page "what transferred".

---

## Then learn algorithms more deeply (sorting proofs, DP, NP preview)

This path taught structures and the algorithms that they serve (BFS, Dijkstra, MST, KMP). A next course is **algorithms as proofs and classes**.

**Sorting proofs.** Why heap sort is Θ(n log n). Why comparison sort needs Ω(n log n) compares in the worst case. Why counting sort is not a comparison sort. You already implemented heap sort. Now read a proof.

**Dynamic programming (DP).** A DP table is often an array or a matrix. The new skill is the recurrence and the argument that subproblems are enough. Structures stay simple. The design is the algorithm.

**NP preview.** Some problems have no known polynomial algorithm. Decision problems in NP have short certificates. NP-complete problems are as hard as each other under reductions. You do not prove P ≠ NP here. You learn to recognize "this looks like an exact cover / SAT / TSP" so that you do not spend a month on a brute force that cannot scale.

```text
keep
  data structures from this path
add
  proof of a bound
  DP recurrence
  a sense of NP-complete names
```

Official starting points (from the topic list):

- CLRS — Introduction to Algorithms
- Sedgewick & Wayne — Algorithms
- VisuAlgo
- CP-Algorithms
- Language docs for built-in collections

Do not drop structure practice while you read proofs. Reimplement one structure each month from memory.

A job that uses graphs still needs BFS and Union-Find. A job that uses DP still needs arrays. The next course does not replace this path.

If you want systems next, go to OS memory, databases and B+ trees, or concurrent maps. Those tracks reuse topic 12, 16, 19, and 21.

### Questions

#### Theoretical questions

1. What does a sorting proof add that heap sort code does not?
2. What is new in DP if the table is only an array?
3. What does an NP preview help you avoid?
4. Why do you still reimplement a structure after you start CLRS?
5. Which later systems topics reuse B+ trees or concurrent maps?

#### Easy practical tasks

1. Write five sentences that define the next course versus this path.
2. Make a table: "Direction", "First book or site". Rows: proofs, DP practice, visualizations.
3. Open CLRS or Sedgewick. Write the chapter names for sorting and for graphs.
4. List three NP-complete names from a preview page. Do not write reductions.

#### Medium practical tasks

1. Read one comparison-sort lower-bound outline. Write ten STE sentences. Cite the source.
2. Solve two DP exercises that use only arrays (for example coin change or a path on a grid). Write the recurrence, not a copied solution essay.
3. Explain in six sentences why TSP brute force fails for n = 30.

#### Advanced practical tasks

1. Write a three-month plan: month 1 sorting proofs and one structure rewrite, month 2 DP, month 3 graph algorithms from a book with proofs. Name the resources.
2. Pick OS, databases, or concurrency as a systems follow-up. Write five structures from this path that you will need there.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe the path from the minimum kit to a forced choice, then to a real source file, then to a proof course.
2. How do the ten suggested practice items cover more than the seven-name kit, and why is that order useful?
3. A teammate says "I know data structures because I import them". Which facts from this topic do you use to reply?
4. What stays the same when you move from this handbook to CLRS: ADT, cost, and test habit?
5. Why does this path end with practice instead of one more structure name?

#### Easy practical tasks

1. Write a one-page cheat sheet: seven implements, ten practice items, choice steps, one real-source track, next-course trio.
2. Fill a status table: implement, forced problems (count), source note (yes/no), next-course plan (yes/no).
3. Draw four boxes in a row: Kit, Choice, Real code, Proofs. Put one example in each box.
4. List the official resource names from this topic.

#### Medium practical tasks

1. Complete any missing kit item this week. Add two forced-choice log lines. Start the real-source note.
2. Interview your own older program. Replace each collection with "ADT + implementation + why". Write the new list.
3. Write a short talk of twelve sentences that teaches a beginner how to pick a structure. Use only terms from this path.

#### Advanced practical tasks

1. From memory, reimplement three of: dynamic array, hash map, heap, BFS. Compare with your old files. Write what you omitted.
2. Produce a personal syllabus for the next six months that cites this handbook's file names and one algorithms book. Do not add new structure chapters that this path already covered.
