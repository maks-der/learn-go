# 4. Arrays

## Description

An array is a contiguous block of elements of one type. This topic explains a fixed-size array. You learn constant-time indexing, insert and delete in the middle, two-dimensional layout, bounds checking, and the words buffer and capacity. Complete this topic before you study dynamic arrays.

Use one term for each concept. The **length** is the number of elements in use. The **index** is an integer position. The first index is 0 in Go and in Python. A **buffer** is the block of memory. **Capacity** is the number of elements that the buffer can hold.

---

## Fixed-size array

A **fixed-size array** has a length that you set when you create the array. The length does not grow. The length does not shrink.

The elements sit in contiguous memory. Each element has the same size. The program can compute the address of index i from the start address.

```text
array a of length 4
index:  0    1    2    3
value:  10   20   30   40
```

In Go, the type `[4]int` is a fixed-size array. The size is part of the type. `[4]int` and `[5]int` are different types.

In Python, a `tuple` has a fixed length. A Python `list` is not a fixed-size array. A Python `array.array` or a library such as NumPy can hold a compact array.

A fixed-size array is simple. You use it when you know n. You use it for a small table. You use it as the inner buffer of a dynamic array.

You cannot append past the length. You must create a new array if you need a different length. Copy the elements that you keep.

### Questions

#### Theoretical questions

1. What is a fixed-size array?
2. Why do all elements have the same size?
3. Why are `[4]int` and `[5]int` different types in Go?
4. When do you use a fixed-size array?
5. What must you do when you need a different length?

#### Easy practical tasks

1. Write a fixed array of five integers in Go, or a tuple of five integers in Python. Print length.
2. Make a table: "Language construct", "Fixed length?". Rows: Go `[n]T`, Go slice, Python list, Python tuple.
3. Draw an array of length 3. Label indexes and values.
4. Write one sentence that states what happens if you try to grow a fixed-size array in place.

#### Medium practical tasks

1. Copy a Go array of 4 integers into a new `[6]int`. Put zeros in the extra cells. In Python, build a new tuple that is longer.
2. Write a function that takes a fixed array and returns the sum. Do not use a dynamic collection inside the function.
3. Explain in six sentences why a fixed size can be part of the type.

#### Advanced practical tasks

1. In Go, compare assignment of `[1000]int` with assignment of a slice header. Write what is copied. In Python, compare tuple assignment with list assignment.
2. Design a small table of month names as a fixed array of 12 strings. Write why a fixed size fits this table.

---

## Indexing in O(1)

**Indexing** is the operation that reads or writes the element at index i.

The address of element i is:

```text
address(i) = start + i * element_size
```

The formula uses a multiply and an add. The formula does not walk the first i elements. The time is Θ(1). The time does not grow with n. The time does not grow with i.

This constant-time index is the main reason to use an array. A linked list does not have this formula. A list walk is Θ(i) to reach index i.

```go
// Go: a is []int or [n]int
x := a[i]
a[i] = x
```

```python
# Python list uses an array of references; index is still O(1)
x = a[i]
a[i] = x
```

The index must be an integer in the valid range. The next section on bounds checking covers an invalid index.

Random access means you can read any index in constant time. Sequential access means you read index 0, then 1, then 2. An array supports both. Sequential access also has good locality.

### Questions

#### Theoretical questions

1. What is the address formula for index i?
2. Why is indexing Θ(1)?
3. Why does a linked list not use that formula?
4. What is random access?
5. Why does sequential access have good locality on an array?

#### Easy practical tasks

1. An array starts at address 2000. Each element is 8 bytes. Write the address of index 3.
2. Write five sentences that explain constant-time indexing.
3. Make a table: "i", "address" for i = 0 to 4, start = 100, size = 4.
4. Read and write index 2 of a small array in Go or in Python.

#### Medium practical tasks

1. Time 1 000 000 reads of a random valid index in a large array. Then time 1 000 000 reads of index 0. Write the two times.
2. Write a function `get(a, i)` that returns `a[i]`. Document that the cost is Θ(1).
3. Explain in six sentences why a Python list index is O(1) even if each element is a large object.

#### Advanced practical tasks

1. Implement a tiny array as a byte buffer plus an element size. Write `get` and `set` with the address formula. Use Go or Python.
2. Compare index cost of an array and of a linked list for i = n/2. Use n = 100 000 if you can. Write the cause.

---

## Insertion and deletion in the middle

**Insert** at index i puts a new element at i. The elements from i to the end move one step to the right. If the array is full, you cannot insert unless you create a larger array.

**Delete** at index i removes the element at i. The elements from i+1 to the end move one step to the left. The length decreases by 1 if the array tracks length. A fixed-size type may keep a hole. You usually compact the hole.

```text
insert 99 at index 1
before:  10  20  30  40
after:   10  99  20  30  40

delete index 2
before:  10  99  20  30  40
after:   10  99  30  40
```

The number of moves is n − i in the usual compact array. The worst case is insert or delete at index 0. You move n elements. The cost is Θ(n).

Insert at the end of a fixed array that is not full is Θ(1) if you track length. That operation is append. Append is not the same as insert in the middle.

Delete by value is not the same as delete by index. Delete by value finds the index first. Then it deletes by index.

### Questions

#### Theoretical questions

1. What must you move when you insert at index i?
2. What is the worst-case index for insert in a compact array?
3. Why is insert in the middle Θ(n)?
4. How is append different from insert in the middle?
5. What extra work does delete by value do?

#### Easy practical tasks

1. Draw insert of 5 at index 0 in the array `[1, 2, 3]`.
2. Draw delete of index 1 in the array `[1, 2, 3, 4]`.
3. Count the moves for insert at i = 2 in an array of length 10.
4. Write four sentences that explain why index 0 is the expensive end.

#### Medium practical tasks

1. Write `insert(a, i, x)` and `delete(a, i)` on a Go slice or a Python list. Do not call the library insert if you study the moves. Use a loop.
2. Time insert at index 0 and append at the end for n = 50 000. Write the two times.
3. Explain in six sentences how a hole in a fixed buffer is different from a compact delete.

#### Advanced practical tasks

1. Implement a compact array with a length field. Support insert and delete. Reject insert when length equals capacity.
2. Write a report that compares delete-by-index and delete-by-swap-with-last when order does not matter. Give costs.

---

## Two-dimensional arrays / row-major vs column-major

A **two-dimensional array** is a table of rows and columns. You write `a[r][c]` or `a[r, c]`.

The table still sits in a one-dimensional memory. The layout maps (row, column) to one index.

**Row-major** layout stores a full row, then the next row. C, C++, and Go nested arrays use row-major order. Many Python nested lists are lists of rows, but each row is a separate object.

**Column-major** layout stores a full column, then the next column. Fortran and some numeric libraries use column-major order.

```text
table 2 by 3
  1 2 3
  4 5 6

row-major one-dimensional:    1 2 3 4 5 6
column-major one-dimensional: 1 4 2 5 3 6
```

Index formulas for a table with R rows and C columns:

```text
row-major:    index = r * C + c
column-major: index = c * R + r
```

A scan that follows the layout has good locality. A scan that goes against the layout jumps and has poor locality.

A true rectangular block is not the same as an array of pointers to rows. An array of pointers to rows can have rows of different lengths. That form is a jagged table. Locality is weaker.

### Questions

#### Theoretical questions

1. What is a two-dimensional array?
2. What does row-major layout store first?
3. What is the row-major index formula?
4. Why does a scan against the layout have poor locality?
5. What is a jagged table?

#### Easy practical tasks

1. Write a 2 by 3 table on paper. Write the row-major sequence.
2. Write the column-major sequence for the same table.
3. Compute the row-major index of row 1, column 2 in a table with 4 columns.
4. Make a table: "Layout", "Formula". Two rows.

#### Medium practical tasks

1. Store a 3 by 3 grid in a one-dimensional array with the row-major formula. Print `get(r, c)`.
2. Sum a large grid by rows, then by columns. Time both loops if the grid is contiguous.
3. Explain in six sentences why Go `[][]int` is not always one contiguous block.

#### Advanced practical tasks

1. Implement a rectangular grid as one slice plus rows and columns. Provide `get` and `set`. Document the layout.
2. Read how one numeric library stores matrices. Write whether the library is row-major or column-major. Cite the page.

---

## Bounds checking

A **valid index** for an array of length n is an integer i with 0 ≤ i < n. Some languages also allow a negative index. Python maps −1 to n − 1. Go does not allow a negative index.

**Bounds checking** is the test that i is valid before the program reads `a[i]`.

If the check fails, the language stops the access. Go **panics**. Python raises `IndexError`. The program does not read a random cell.

```go
a := [3]int{1, 2, 3}
_ = a[3] // panic: index out of range
```

```python
a = [1, 2, 3]
# a[3] raises IndexError
```

Bounds checking costs a compare. Compilers can remove a check when they prove that i is valid. You still write code that stays in range.

An **off-by-one error** uses n or uses −1 by mistake. Loops must use `i < n` when indexes start at 0.

Do not turn off bounds checking to go faster unless you are in a special low-level tool and you can prove safety. For this handbook, keep the checks.

### Questions

#### Theoretical questions

1. What is a valid index for length n when indexes start at 0?
2. What is bounds checking?
3. What does Go do on an invalid index?
4. What does Python do on an invalid index?
5. What is an off-by-one error?

#### Easy practical tasks

1. Write a loop that prints every valid index of an array of length 4.
2. Cause one out-of-range access on purpose. Record the error name. Then fix the index.
3. Make a table: "Index", "Valid for n = 3?". Rows: −1, 0, 2, 3.
4. Write four sentences about why a language checks bounds.

#### Medium practical tasks

1. Write a safe `get(a, i)` that returns an error or a boolean instead of a panic or exception.
2. Find one off-by-one bug in a loop that you write on purpose. Fix the loop. Write the wrong bound and the right bound.
3. Explain in six sentences how Python negative indexes differ from Go indexes.

#### Advanced practical tasks

1. Write tests for indexes: 0, n−1, n, −1. State the expected result in your language.
2. Read how one compiler can remove bounds checks. Write eight STE sentences. Cite the source.

---

## Buffer and capacity

A **buffer** is the memory block that holds elements. The program reads and writes the buffer.

**Capacity** is the number of elements that the buffer can hold. **Length** is the number of elements in use. Length is less than or equal to capacity.

```text
buffer cells:  a  b  c  _  _  _
length:        3
capacity:      6
```

A fixed-size array often has length equal to capacity. A dynamic array (next topic) has length that can be less than capacity.

A buffer can hold raw bytes. A file read fills a byte buffer. You then interpret the bytes as integers or as text. The capacity is the size of the byte block. The length is the number of bytes in use.

Do not read past length. The cells between length and capacity can be zeros or old values. Those cells are not part of the logical array.

When you reuse a buffer, you can set length to 0 and keep the capacity. That reuse avoids a new allocation. You still must not read stale values as live data.

### Questions

#### Theoretical questions

1. What is a buffer?
2. What is capacity?
3. What is the relation between length and capacity?
4. Why must you not read past length?
5. Why can you reuse a buffer after you set length to 0?

#### Easy practical tasks

1. Draw a buffer of capacity 5 and length 2. Mark live cells and unused cells.
2. Write five sentences that contrast length and capacity.
3. Make a table: "Word", "Meaning". Rows: buffer, length, capacity, index.
4. In Go, print `len` and `cap` of a slice that you make with `make([]int, 2, 5)`. In Python, write that a list has length and an over-allocated inner array.

#### Medium practical tasks

1. Write a small byte buffer that stores n bytes and a capacity. Provide `appendByte` that fails when the buffer is full.
2. Fill a buffer, set length to 0, then append again. Show that capacity can stay the same.
3. Explain in six sentences why unused capacity still uses memory.

#### Advanced practical tasks

1. Implement a fixed buffer of integers with length and capacity. Operations: append, get, clear (length = 0). Reject overflow.
2. Read the Go slice header (`ptr`, `len`, `cap`) or the CPython list over-allocation note. Write ten STE sentences. Cite the source.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe the path from a start address to `a[r][c]` in a row-major rectangular block.
2. Why do insert in the middle and constant-time index exist in the same structure?
3. How do bounds checking and the length field work together?
4. A teammate says a Python list is a fixed-size array. Which facts do you use in a reply?
5. When do you store a table as one buffer instead of many row objects?

#### Easy practical tasks

1. Write a one-page cheat sheet: array, index, length, capacity, buffer, row-major, column-major, bounds check.
2. Draw a 2 by 2 table in row-major form inside one buffer of four cells. Label each index.
3. Count moves to delete index 0 in an array of length 8.
4. Write a valid-index rule for your language in one sentence. Include negative indexes if they exist.

#### Medium practical tasks

1. Build a small matrix ADT on one slice. Support `get`, `set`, and a bounds error. Use row-major order.
2. Time a row-wise sum and a column-wise sum on a large contiguous grid. Write the two times.
3. Write insert and delete on a compact buffer with a length field. Add tests for i = 0, i = length−1, and i = length (insert as append).

#### Advanced practical tasks

1. Implement a jagged table as an array of row arrays. Then implement a rectangular view on one buffer. Compare get cost and locality in writing.
2. Read one standard-library array or slice document. Rewrite every operation with a complete complexity statement.
