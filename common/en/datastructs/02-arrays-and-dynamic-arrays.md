# 2. Arrays and Dynamic Arrays

## Description

An array is a contiguous block of elements of one type. A dynamic array is an array that can grow. This topic explains fixed-size arrays, constant-time indexing, insert and delete in the middle, bounds, row-major and column-major layout, growth, amortized append, length, and capacity.

Complete this topic after foundations. Complete this topic before you study linked lists.

Use one term for each concept. The **length** is the number of elements in use. The **index** is an integer position. This handbook uses 0 as the first index. A **buffer** is the block of memory. **Capacity** is the number of elements that the buffer can hold. **Append** adds one element at the end. **Resize** allocates a new buffer and copies the live elements.

---

## Fixed-size array and O(1) indexing

A **fixed-size array** has a length that you set when you create the array. The length does not grow. The length does not shrink.

The elements sit in contiguous memory. Each element has the same size. The program can compute the address of index i from the start address.

```text
array a of length 4
index:  0    1    2    3
value:  10   20   30   40
```

**Indexing** is the operation that reads or writes the element at index i.

The address of element i is:

```text
address(i) = start + i * element_size
```

The formula uses a multiply and an add. The formula does not walk the first i elements. The time is Θ(1). The time does not grow with n. The time does not grow with i.

This constant-time index is the main reason to use an array. A linked list does not have this formula. A list walk is Θ(i) to reach index i.

A fixed-size array is simple. You use it when you know n. You use it for a small table. You use it as the inner buffer of a dynamic array.

You cannot append past the length. You must create a new array if you need a different length. Copy the elements that you keep.

### Questions

#### Theoretical questions

1. What is a fixed-size array?
2. Why do all elements have the same size?
3. What is the address formula for index i?
4. Why is indexing Θ(1)?
5. What must you do when you need a different length?

#### Easy practical tasks

1. Draw an array of length 3. Label indexes and values.
2. Write one sentence that states what happens if you try to grow a fixed-size array in place.
3. Make a table: "Index i", "Address if start = 1000 and size = 4" for i = 0, 1, 2, 3.
4. Write five sentences that define a fixed-size array and indexing. Use only facts from this section.

#### Medium practical tasks

1. Write the steps of a function that returns the sum of a fixed array. Do not use a second collection.
2. Explain in six sentences why a linked list cannot use the address formula.
3. Copy an array of 4 integers into a new array of length 6 on paper. Put zeros in the extra cells. Write how many elements you copy.

#### Advanced practical tasks

1. Design a small table of month names as a fixed array of 12 strings. Write why a fixed size fits this table.
2. Compare assignment of a large fixed array with assignment of a small header that points to a buffer. Write what is copied in each case.

---

## Insert and delete in the middle

**Insert at index i** puts a new element at i. The elements from i to the old end must move one cell to the right. That shift is Θ(n − i). In the worst case you insert at index 0. The shift is Θ(n).

**Delete at index i** removes the element at i. The elements from i + 1 to the end must move one cell to the left. That shift is also Θ(n − i). Delete at index 0 is Θ(n).

```text
array:  10  20  30  40
insert 25 at index 2
shift:  30 and 40 move right
result: 10  20  25  30  40

delete index 1 (value 20)
shift:  25, 30, 40 move left
result: 10  25  30  40
```

Insert and delete at the end do not shift a large block. Append writes one cell. Delete of the last element decreases length. Those operations are Θ(1) on a fixed array that still has room, or on a dynamic array that does not resize.

Find by value is a scan. The scan is Θ(n) in the worst case. After you find the index, delete still shifts.

Do not say "array insert is slow" without a position. Insert at the end is cheap when there is room. Insert at the front is expensive.

A hole in the middle is not allowed if you want the address formula to stay simple. After delete, the live elements stay contiguous from 0 to length − 1.

### Questions

#### Theoretical questions

1. Why does insert at index i shift elements?
2. What is the worst-case cost of insert at the front?
3. Why is delete at the end cheaper than delete at the front?
4. What extra work does delete by value do?
5. Why must live elements stay contiguous after delete?

#### Easy practical tasks

1. Start with [10, 20, 30, 40]. Write the array after insert of 15 at index 1.
2. Start with [10, 20, 30, 40]. Write the array after delete of index 2.
3. Make a table: "Operation", "Cells that move" for insert at 0, insert at end, delete at 0, delete at end. Use n = 5.
4. Write five sentences about insert and delete in the middle. Use only facts from this section.

#### Medium practical tasks

1. Count the assignments for insert at index i in an array of length n. Write the count as a formula.
2. Write the steps to delete the first element that equals x. Include the scan and the shift.
3. Explain in six sentences why append is not the same operation as insert at index 0.

#### Advanced practical tasks

1. Implement insert and delete at index i on a fixed buffer in a language that you know. Test front, middle, and end.
2. Time insert at 0 versus append for n = 10 000. Write the two times and the growth classes that you expect.

---

## Bounds, row-major vs column-major

A valid index i satisfies 0 ≤ i ≤ length − 1. An index outside that range is an **out-of-bounds** access.

Out-of-bounds access is a defect. Some languages stop the program. Some languages do not check and write a wrong cell. Do not use an index that you did not validate.

```text
length = 4
valid indexes: 0, 1, 2, 3
invalid:       -1, 4, 99
```

A **two-dimensional array** stores a table of rows and columns. The machine memory is still one line of cells. You must choose an order.

**Row-major** order stores a full row, then the next row. Element (r, c) sits at:

```text
index = r * number_of_columns + c
```

**Column-major** order stores a full column, then the next column. Element (r, c) sits at:

```text
index = c * number_of_rows + r
```

```text
table 2 rows × 3 columns
values:
  A B C
  D E F

row-major block:    A B C D E F
column-major block: A D B E C F
```

Walk in the same order as the layout if you want good locality. A row walk is good in row-major layout. A column walk is good in column-major layout.

Bounds apply to each dimension. A valid pair (r, c) needs 0 ≤ r < rows and 0 ≤ c < columns.

Do not mix the two orders when you pass a block to another component. Write the order in the contract.

### Questions

#### Theoretical questions

1. What is a valid index for length n?
2. What is an out-of-bounds access?
3. What is row-major order?
4. What is column-major order?
5. Why does walk order matter for locality?

#### Easy practical tasks

1. For length 5, list all valid indexes. List three invalid indexes.
2. Draw the 2 × 3 table from this section. Write the row-major block.
3. Write the column-major block for the same table.
4. Make a table: "Order", "Index of (1, 2) for 2 × 3". Add two rows.

#### Medium practical tasks

1. Write a formula for element (r, c) in a 4 × 5 row-major array. Compute the index of (2, 3).
2. Explain in six sentences why a column walk on a row-major block can miss the cache more often.
3. Write two checks that you do before you read a[r][c].

#### Advanced practical tasks

1. Flatten a 3 × 3 table to a one-dimensional array in both orders. Write both arrays.
2. Write a short note: when do you store a matrix in row-major order, and when do you store it in column-major order? Give one example for each.

---

## Growth strategy and amortized append

A dynamic array starts with a capacity of 0 or a small number. When you append and length equals capacity, the structure is full. The program allocates a larger buffer.

A common **growth strategy** is to multiply capacity by 2. Some libraries use 1.5. The exact factor can change. The idea is **geometric growth**: each new capacity is a constant factor times the old capacity.

```text
append when full
  new_capacity = max(1, old_capacity * 2)
  allocate new buffer
  copy length elements
  append the new element
```

Growth by +1 is a different strategy. Each append that is full copies all n elements. The total copy cost for n appends is about 1 + 2 + ... + n. That sum is Θ(n²). Do not use +1 growth for a general dynamic array.

Geometric growth keeps the total copy cost linear in n. Most append operations write one cell and add 1 to length. Those appends are Θ(1). When the buffer is full, append resizes. Resize copies length elements. That append is Θ(n).

With doubling, the copies occur at sizes 1, 2, 4, ..., n. The sum of those copies is less than 2n. The amortized cost of one append is Θ(1).

A library can pick a first capacity of 1, 4, or 8. A small first capacity saves space for empty lists. A larger first capacity can reduce the first few resizes.

Some libraries shrink the buffer when length becomes much smaller than capacity. Shrink is optional. Shrink can copy again. Do not shrink on every delete. A common rule is to shrink only when length is a small fraction of capacity.

### Questions

#### Theoretical questions

1. When does a dynamic array grow?
2. What is geometric growth?
3. Why is growth by +1 a bad default?
4. What does the program copy on resize?
5. Why is amortized append Θ(1) with doubling?

#### Easy practical tasks

1. Start at capacity 1. Double until the capacity is at least 20. List each capacity.
2. Write five sentences that describe append when the buffer is full.
3. Draw a buffer of capacity 2 and length 2. Then draw the buffer after a doubling grow and one append.
4. Make a table: "Growth rule", "Total copy class for n appends". Rows: +1, ×2.

#### Medium practical tasks

1. Write a tiny dynamic array that doubles. Print capacity after each grow during 16 appends.
2. Change the factor from 2 to 1.5 (use integers: new = old + old/2, minimum +1). List capacities from 2 to at least 20.
3. Explain in six sentences why shrink on every delete can be expensive.

#### Advanced practical tasks

1. Implement growth by +1 and growth by ×2. Count element copies for 1 000 appends in each strategy.
2. Write a one-page note that compares factor 2 and factor 1.5 for space waste and copy cost.

---

## Length vs capacity

**Length** is the number of live elements. Live elements sit in indexes 0 .. length − 1.

**Capacity** is the number of cells in the buffer. Capacity is always greater than or equal to length.

```text
buffer:   10  20  30  ??  ??
length:   3
capacity: 5
live:     indexes 0, 1, 2
unused:   indexes 3, 4
```

The unused cells are extra space. They make the next appends cheap. They use memory.

Do not iterate past length. The unused cells can hold old values. Those values are not part of the collection.

After resize, length stays the same until you append. Capacity becomes the new buffer size. After append, length increases by 1.

A fixed-size array has length equal to capacity. There is no unused cell in the contract. A dynamic array separates the two numbers.

When you report size to a user, you report length. When you plan memory, you look at capacity.

Clear or reset can set length to 0 and keep the buffer. The next appends reuse the capacity. That reuse avoids a new allocate.

### Questions

#### Theoretical questions

1. What is length?
2. What is capacity?
3. Why can unused cells hold old values?
4. What happens to length and capacity on a doubling resize, before the new append?
5. Why do you report length, not capacity, as the size of the collection?

#### Easy practical tasks

1. Draw a buffer of capacity 6 and length 2. Mark live cells and unused cells.
2. Make a table: "Event", "Length", "Capacity". Start at 0, 0. Then append until one grow at capacity 2.
3. Write five sentences that contrast length and capacity.
4. Write why a loop must stop at length, not at capacity.

#### Medium practical tasks

1. Write the steps of append when length < capacity and when length = capacity.
2. Explain in six sentences how a reset that keeps capacity helps a later fill.
3. For a doubling array, write length and capacity after 5 appends from empty if the first grow creates capacity 1, then 2, then 4, then 8.

#### Advanced practical tasks

1. Implement a dynamic array with length and capacity fields in a language that you know. Keep unused cells out of iterate.
2. Measure memory of a collection that you grew to n = 1 000 000 and then cleared without shrink. Write why capacity can stay large.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do the address formula, middle insert, and geometric growth work together in one dynamic array?
2. When do you keep a fixed-size array, and when do you use a dynamic array?
3. Why must you name the index when you say that array insert is cheap or expensive?
4. How do bounds checks and layout order protect correctness and locality?
5. What extra space does capacity buy, and what cost does that space have?

#### Easy practical tasks

1. Write a one-page cheat sheet: fixed array, index, bounds, row-major, column-major, grow, append, length, capacity.
2. Draw one picture that shows a 2 × 2 row-major block inside a dynamic buffer with unused cells.
3. List five operations on a playlist stored as a dynamic array. Write Θ for each if you can.
4. Make a table: "Need", "Use fixed or dynamic". Add four rows.

#### Medium practical tasks

1. Implement a dynamic array of integers with get, set, append, insert at i, and delete at i. Do not write a solution here. Write a test list first.
2. Flatten a 2 × 4 table, then append two more values to the one-dimensional buffer. Write length, capacity, and the row-major meaning of the first eight cells.
3. Write a short report: one defect from an out-of-bounds index, and one defect from iterate to capacity.

#### Advanced practical tasks

1. Compare three growth factors (2, 1.5, +1) on the same n appends. Report copies and final unused cells.
2. Design an ADT card for a vector: operations, errors for bad indexes, and a note on amortized append. Then implement the card in a language that you know.
