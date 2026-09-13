# 1. Foundations

## Description

A data structure is a method to store and organize values in a program. This topic explains that idea. You learn the difference between an abstract data type and an implementation. You also learn Big-O in practice, common growth classes, amortized cost, and memory locality.

Complete this topic first. Complete this topic before you study arrays and lists.

Use one term for each concept. An **abstract data type (ADT)** is a list of operations. An **implementation** is the layout and the code that do those operations. **n** is the number of elements. A **growth class** is a family of functions such as constant or linear. **Locality** is the use of nearby memory cells in a short time.

---

## What a data structure is

A data structure is an organization of values in memory. The structure decides how you store values and how you get values. The structure also decides how you add values and how you remove values.

Examples of data structures: an array, a linked list, a hash table, and a tree. Each structure stores values. Each structure gives a different cost for each operation.

You do not select a structure only by name. You select a structure by the operations that you need. If you need fast access by index, you use an array. If you need fast insert at the front, you use a linked list.

A data structure is not an algorithm. An algorithm is a sequence of steps. A data structure is the layout of the values. An algorithm uses a data structure.

A value in a structure is an **element**. Some texts say "item". This handbook uses **element** for a stored value. This handbook uses **node** only for a linked cell that holds a value and a link.

```text
structure: array of 4 integers
indexes:   0    1    2    3
elements:  10   20   30   40
```

The same four numbers can live in a list of nodes. The values do not change. The layout in memory does change. That change of layout is the data structure.

A program uses time and space. Time is the work that the processor does. Space is the memory that the program holds. You compare structures by those costs. You measure time in steps, not in seconds. A step is one simple action: a compare, an assignment, or an index. Seconds change with the machine. Steps do not change.

### Questions

#### Theoretical questions

1. What is a data structure?
2. How is a data structure different from an algorithm?
3. What does this handbook mean by "element"?
4. When do you use the word "node"?
5. Why do you select a structure by operations, not only by name?

#### Easy practical tasks

1. Write five sentences that describe a data structure. Use only facts from this section.
2. Make a two-column table: "Structure" and "One typical operation". Add four rows.
3. List three program types that need a data structure (for example, a contact list). For each type, name one structure that can fit.
4. Draw the array of four integers from this section. Label indexes and elements.

#### Medium practical tasks

1. Compare an array and a linked list in six short sentences. Cover store, get, add, and remove.
2. Take a grocery list on paper. Write how an array stores that list. Write how a list of nodes stores that list.
3. Count the assignments in this sequence: set a to 1, set b to 2, set c to a + b. Write why this count is a step count, not a time in seconds.

#### Advanced practical tasks

1. Read the first chapter of one standard algorithms book. Write ten sentences that map the book terms to the terms in this section.
2. Select a small real program that you know. List every collection that the program uses. For each collection, write the data structure that it is.

---

## Abstract data type vs implementation

An abstract data type (ADT) is a contract. The contract names the operations. The contract does not name the layout in memory.

A stack is an ADT. The stack operations are push, pop, and peek. The contract does not say "array" or "list".

An implementation is the code that does the contract. One stack can use an array. A different stack can use a linked list. Both implementations are stacks if they do the same operations.

The ADT hides the layout. A user of the ADT calls the operations. The user does not read the inner fields. This split lets you change the implementation later.

```text
ADT: Stack
  push(x)  — put x on the top
  pop()    — remove the top and return it
  peek()   — return the top and do not remove it
  empty()  — true when the stack holds no element

Implementation A: array + top index
Implementation B: linked list + head node
```

A type in a programming language is not always an ADT. A language type is a tool. An ADT is a design idea. You can write an ADT as a type. You can also write an ADT as a set of functions.

Use one term for each idea. Say **ADT** when you talk about the contract. Say **implementation** when you talk about the code and the memory layout.

The ADT name does not tell you the cost. The implementation tells you the cost. Always ask which implementation you have.

One ADT can have many implementations. The operations stay the same. The costs change. You select an implementation for the operations that you do often.

### Questions

#### Theoretical questions

1. What is an abstract data type?
2. What is an implementation?
3. Why does a stack ADT not name "array" or "list"?
4. What does it mean to hide the layout?
5. Why does the ADT name not tell you the cost?

#### Easy practical tasks

1. Write the ADT contract for a queue: enqueue, dequeue, front, empty. Use four short lines.
2. Name two implementations that can do that queue contract.
3. Make a table with columns "ADT name", "Operations", and "Possible implementation". Add three rows.
4. In four sentences, explain why a user of a stack must not read the inner array.

#### Medium practical tasks

1. Write a stack ADT as a list of function names and inputs. Do not write the function bodies.
2. Draw two boxes: "ADT" and "implementation". Put the stack operations in the first box. Put "array" and "list" in the second box.
3. A teammate says "use a list". Write four questions that you ask before you select the implementation.

#### Advanced practical tasks

1. Write two small implementations of the same stack ADT in a language that you know: one array, one list. Keep the public function names the same.
2. Change the array implementation so that the capacity grows. Do not change the function names. Write why the ADT stays the same.

---

## Big-O in practice: best, average, worst

Complexity describes how the cost of an operation grows when the input grows. **Big-O** is the usual notation in daily work. Big-O is an upper bound. If an operation is O(n), the time does not grow faster than a constant times n, for large n.

Big-O can be loose. O(n) is also O(n²). A tight upper bound is more useful. **Big-Theta** is a tight bound. If an operation is Θ(n), the time grows like n. Many texts write O(n) when they mean Θ(n). Read the text.

Ignore constant factors in the notation. 2n and 100n are both Θ(n). Do not ignore constants in a real program when n is small. Notation is for growth. Measurement is for a real machine.

Drop lower terms. n² + 3n + 10 is Θ(n²). The n² term dominates for large n.

One operation can have more than one time cost. The cost depends on the input.

**Best case** is the cheapest input of size n. Example: find a value that sits in the first cell of an array.

**Worst case** is the most expensive input of size n. Example: find a value that is not in the array. You visit every cell.

**Average case** is the expected cost for a typical input. You need a model of typical input. If each position is equally likely, the average find in an unsorted array visits about n/2 cells.

Write the case when you write a bound. "Find is O(n)" is not complete. "Find in an unsorted array is O(n) in the worst case" is complete.

```text
binary search in a sorted array of n elements
  best case:    Θ(1)   — the middle cell is the target
  worst case:   Θ(log n)
  average case: Θ(log n)  for a usual model
```

Do not mix cases. Do not compare the best case of one structure with the worst case of a different structure.

### Questions

#### Theoretical questions

1. What does Big-O say about growth?
2. What is the best case of an operation?
3. What is the worst case of an operation?
4. What extra fact do you need for an average case?
5. Why is "find is O(n)" not a complete statement?

#### Easy practical tasks

1. Write five sentences that define Big-O, best case, and worst case. Use only facts from this section.
2. For linear search in an array, write one input for best case and one input for worst case.
3. Classify 5n + 20 as Θ of one simple function. Write the function.
4. Make a table: "Case", "Meaning". Add three rows.

#### Medium practical tasks

1. A loop visits each of n elements two times. Write Big-O and Big-Theta for the visit count.
2. Explain in six sentences why 1000n can be slower than n² for small n, and faster for large n.
3. Write the worst-case bound for a scan of n cells when the target is absent. Then write the best-case bound when the target is in cell 0.

#### Advanced practical tasks

1. Read a formal definition of Big-O (c, n0). Rewrite the definition in eight short STE sentences.
2. Show with a table that 3n + 7 is O(n) and that n is O(3n + 7). Then write Θ.

---

## Constant, logarithmic, linear, linearithmic, quadratic

A growth class is a name for a family of functions. You use these names when you compare operations.

**Constant** — Θ(1). The cost does not grow with n. Example: read index i in an array.

**Logarithmic** — Θ(log n). The cost grows like the number of times you can divide n by 2. Example: binary search in a sorted array. One step throws away about half of the remaining cells.

**Linear** — Θ(n). The cost grows like n. Example: visit each element one time.

**Linearithmic** — Θ(n log n). The cost is a constant times n times log n. Example: an efficient comparison sort of n elements.

**Quadratic** — Θ(n²). The cost grows like n times n. Example: compare each pair of elements.

```text
n        10     100     1 000    1 000 000
1        1      1       1        1
log n    ~3     ~7      ~10      ~20
n        10     100     1 000    1 000 000
n log n  ~33    ~664    ~10 000  ~20 000 000
n²       100    10 000  1 000 000  10^12
```

The table uses log base 2 and rounded values. The exact base does not change the class.

A constant-time operation is not free. It still uses steps. The word "constant" means the step count does not grow with n.

Logarithmic growth stays small for large n. Quadratic growth becomes large. For n = 1 000 000, n² is about one million million steps. That cost is often too large.

Do not say "fast" without a class. "Fast" is not a growth class.

Other classes exist: cubic, exponential, factorial. This topic uses the five names above for daily work.

### Questions

#### Theoretical questions

1. What does constant growth mean?
2. Why is binary search logarithmic?
3. What is a linear scan?
4. What does linearithmic mean?
5. Why is quadratic growth a problem for large n?

#### Easy practical tasks

1. Make a table: "Class", "Θ notation", "One example". Add five rows.
2. For n = 16, write approximate values of 1, log2 n, n, n log2 n, and n².
3. Write five sentences that name the five growth classes. Use only facts from this section.
4. Mark this loop as a class: for each of n elements, do one assignment.

#### Medium practical tasks

1. A nested pair of loops each run n times. Write the growth class of the inner body count.
2. Draw a sketch of n versus n log n versus n² for n = 2, 4, 8, 16. Do not use a computer if you can count by hand.
3. Write which class you expect for: array index, full visit, binary search, pair compare, efficient comparison sort.

#### Advanced practical tasks

1. Time a linear loop and a nested quadratic loop in a language that you know for n = 1 000 and n = 4 000. Write the two ratios. Write why the ratios are not exact powers.
2. Read a short note on log base. Show that log2 n and log10 n differ by a constant factor. Write why they are the same class.

---

## Amortized cost

**Amortized cost** is the average cost of one operation in a sequence. Some operations in the sequence are expensive. Many operations are cheap. The average can still be small.

The usual example is append on a dynamic array. Most appends write one cell. That cost is Θ(1). When the buffer is full, append copies all n elements to a larger buffer. That one append is Θ(n).

If the buffer doubles, the copies occur at sizes 1, 2, 4, ..., n. The sum of those copies is less than 2n. For n appends from empty, the total copy cost is Θ(n). The amortized cost of one append is Θ(1).

```text
n appends from empty, doubling
  cheap appends: write one cell
  grow at:       1, 2, 4, ..., n
  total copies:  < 2n
  amortized:     Θ(1) per append
```

Amortized cost is not the same as average case. Average case needs a model of random input. Amortized cost needs a sequence of operations. The sequence can be the worst sequence. Doubling still gives Θ(1) amortized append for any sequence of n appends from empty.

Amortized cost is not a guarantee for one operation. One append can still take Θ(n) time. If you need a strict worst-case bound for every single call, amortized Θ(1) is not enough.

A simple method to see amortized cost is the **aggregate method**: add the total cost of n operations, then divide by n.

Do not use growth by +1 for a general dynamic array. Each full append copies n elements. The total copy cost is about 1 + 2 + ... + n. That sum is Θ(n²). The amortized cost is then Θ(n), not Θ(1).

### Questions

#### Theoretical questions

1. What is amortized cost?
2. How is amortized cost different from average case?
3. Why can one append still take Θ(n) time when amortized append is Θ(1)?
4. What does the aggregate method do?
5. Why is growth by +1 a bad default for a dynamic array?

#### Easy practical tasks

1. Write five sentences that define amortized cost. Use only facts from this section.
2. Start at capacity 1. Double until the capacity is at least 16. List each capacity.
3. Make a table: "Append number", "Grow? (yes/no)" for eight appends that start at capacity 1 and double.
4. Draw a buffer of capacity 2 and length 2. Then draw the buffer after a doubling grow and one append.

#### Medium practical tasks

1. Add the copy counts at sizes 1, 2, 4, 8 for eight appends. Write the total. Write the amortized cost as total divided by 8.
2. Explain in six sentences why doubling keeps the total copy cost linear in n.
3. Write one sequence of operations where you care about the worst single call, not the amortized average.

#### Advanced practical tasks

1. Count element copies for 1 000 appends with growth by +1 and with growth by ×2. Write the two totals.
2. Write a one-page note that compares amortized Θ(1) with worst-case Θ(1). Give one structure that needs the worst-case bound.

---

## Memory locality: arrays vs pointer chasing

**Memory locality** is the use of nearby addresses in a short time. Processors load a **cache line**: a small block of nearby bytes. If the next element is in that block, the next access is cheap.

An **array** stores elements in contiguous cells. A walk from index 0 to n−1 uses nearby addresses. The walk has good locality. The cache can hold many elements in one line.

A **linked list** stores elements in nodes. Each node holds a value and a pointer to the next node. The next node can sit far away in memory. A walk follows pointers. That walk is **pointer chasing**. Pointer chasing has poor locality. Each node can miss the cache.

```text
array walk:     [10][20][30][40]   — cells sit in one block
list walk:      [10|•] -----> [20|•] -----> [30|null]
                nodes can sit far apart
```

Good locality does not change the growth class. An array walk and a list walk are both Θ(n) to visit n elements. The constant factors can differ a lot. On a real machine, the array walk is often faster.

Locality also matters for insert and delete. An array that shifts a block of cells still uses contiguous writes. A list that changes two pointers does less writing, but the later walk can stay slow.

Do not select a list only because insert in the middle is "a pointer change". Count the full cost: find the position, change the links, then walk later with poor locality.

This handbook uses **contiguous** for an array layout. This handbook uses **linked** for a node layout.

### Questions

#### Theoretical questions

1. What is memory locality?
2. What is a cache line at a high level?
3. Why does an array walk have good locality?
4. What is pointer chasing?
5. Why can two Θ(n) walks have different real times?

#### Easy practical tasks

1. Write five sentences that contrast an array walk and a list walk. Use only facts from this section.
2. Draw four array cells in one row. Draw four list nodes with arrows that jump.
3. Make a table: "Layout", "Locality (good/poor)", "Visit class". Add two rows.
4. Write why a cache line helps an array more than a list.

#### Medium practical tasks

1. Explain in six sentences why insert in a list can still lose to an array for a later full walk.
2. Write a short experiment plan: visit n integers in an array, then visit n integers in a list of nodes. Name the metric (time). Do not invent results.
3. Label this sentence as true or false and give a reason: "Poor locality changes the Big-O class of a full visit."

#### Advanced practical tasks

1. Implement an array visit and a linked-node visit in a language that you know. Time both for n = 100 000 if the machine allows it. Write the two times and one locality sentence.
2. Read a short note on cache lines (size is often 64 bytes). Estimate how many 8-byte integers fit in one line. Write how many lines an array of 1 000 integers needs.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe the path from a problem to a selected data structure. Name ADT, operations, implementation, and growth class.
2. How do time cost, space cost, and locality work together when you select a structure?
3. A teammate says that a queue "is a list". Which facts do you use to correct that sentence?
4. When do you write amortized cost, and when do you write worst-case cost?
5. Why do you name the case (best, average, worst) when you compare two structures?

#### Easy practical tasks

1. Write a one-page cheat sheet with these terms: data structure, element, node, ADT, implementation, Big-O, amortized, locality.
2. Draw three boxes in a row: Problem, ADT, Implementation. Put one example in each box.
3. Make a two-column table: "Term in this topic" and "One-sentence meaning". Add eight rows.
4. List five operations for a music playlist. Mark each operation as insert, delete, find, or iterate.

#### Medium practical tasks

1. Write two implementations of a tiny bag ADT that only supports insert and iterate. Use an array and a linked list. Keep the public names the same.
2. Design a table that maps each of the five growth classes to one structure or algorithm from later topics that you already know by name.
3. Write a short note: n = 10 versus n = 1 000 000. For each size, say whether you care more about class, constants, or locality.

#### Advanced practical tasks

1. Design a study plan for the next fourteen topics. For each topic, write one ADT and one implementation that you will build.
2. Select one collection from a program that you know. Write the ADT, the likely implementation, the growth class of find, and one locality comment.
