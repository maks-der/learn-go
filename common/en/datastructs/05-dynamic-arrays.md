# 5. Dynamic Arrays (Vectors / Slices / ArrayLists)

## Description

A dynamic array is an array that can grow. This topic explains growth, amortized append, the cost of geometric resize, optional shrink, iteration, and random access. Complete this topic before you compare lists and arrays in real code.

Use one term for each concept. A **dynamic array** is a buffer plus a length. The buffer has a **capacity**. **Append** adds one element at the end. **Resize** allocates a new buffer and copies the live elements. Languages use different names: Go **slice**, Python **list**, Java **ArrayList**, C++ **vector**.

---

## Growth strategy (typically ×2)

A dynamic array starts with a capacity of 0 or a small number. When you append and length equals capacity, the structure is full. The program allocates a larger buffer.

A common **growth strategy** is to multiply capacity by 2. Some libraries use 1.5. The exact factor can change. The idea is **geometric growth**: each new capacity is a constant factor times the old capacity.

```text
append when full
  new_capacity = max(1, old_capacity * 2)
  allocate new buffer
  copy length elements
  append the new element
```

```go
// Go: append may grow the slice
s = append(s, x)
```

```python
# Python: list.append grows the inner array
a.append(x)
```

Growth by +1 is a different strategy. Each append that is full copies all n elements. The total copy cost for n appends is about 1 + 2 + ... + n. That sum is Θ(n²). Do not use +1 growth for a general dynamic array.

Geometric growth keeps the total copy cost linear in n. The next section states the amortized cost of one append.

A library can pick a first capacity of 1, 4, or 8. A small first capacity saves space for empty lists. A larger first capacity can reduce the first few resizes.

### Questions

#### Theoretical questions

1. When does a dynamic array grow?
2. What is geometric growth?
3. Why is growth by +1 a bad default?
4. What does the program copy on resize?
5. Why can the first capacity be a small constant?

#### Easy practical tasks

1. Start at capacity 1. Double until the capacity is at least 20. List each capacity.
2. Make a table: "Language name", "Dynamic array". Rows: Go, Python, Java, C++.
3. Write five sentences that describe append when the buffer is full.
4. Draw a buffer of capacity 2 and length 2. Then draw the buffer after a doubling grow and one append.

#### Medium practical tasks

1. Write a tiny dynamic array that doubles. Print capacity after each grow during 16 appends.
2. Change the factor from 2 to 1.5 (use integers: new = old + old/2, minimum +1). List capacities from 2 to at least 20.
3. Find the growth note for Go slices or for CPython lists. Write the factor or the rule in your words.

#### Advanced practical tasks

1. Implement growth by +1 and growth by ×2. Count element copies for 1 000 appends in each strategy.
2. Write a one-page note that compares factor 2 and factor 1.5 for space waste and copy cost.

---

## Amortized append

Most append operations write one cell and add 1 to length. Those appends are Θ(1).

When the buffer is full, append resizes. Resize copies length elements. That append is Θ(n).

**Amortized append** is the average cost of one append in a sequence of n appends from empty. With doubling, the copies occur at sizes 1, 2, 4, ..., n. The sum of those copies is less than 2n. The amortized cost of one append is Θ(1).

```text
n appends from empty, doubling
  cheap writes:  n
  copies:        1 + 2 + 4 + ... + n/2  <  n
  total work:    Θ(n)
  per append:    Θ(1) amortized
```

Amortized constant time is not worst-case constant time. One append can still be Θ(n). A real-time system that cannot pause for a copy needs a different structure or a pre-sized buffer.

If you know n before you start, you **reserve** capacity. Reserve allocates one buffer of size n. Then each append is Θ(1) worst case until you exceed n. Go uses `make([]T, 0, n)`. Python uses tricks that are less direct. Some languages have `reserve`.

Write "amortized Θ(1)" for append when you use geometric growth and you do not reserve.

### Questions

#### Theoretical questions

1. What is the cheap path of append?
2. What is the expensive path of append?
3. Why is the sum of doubling copies less than 2n?
4. How is amortized Θ(1) different from worst-case Θ(1)?
5. What does reserve do?

#### Easy practical tasks

1. Write a sequence of 8 appends. Mark which appends resize if you start at capacity 1 and double.
2. Make a table: "Append number", "Resize?". Assume start capacity 1, factor 2, eight appends.
3. Write four sentences that define amortized append.
4. In Go, create a slice with `make([]int, 0, 100)` and append 100 times. In Python, build a list of 100 items. Write why reserve helps in Go.

#### Medium practical tasks

1. Count copies in a doubling array for n = 32 appends from capacity 1. Write the sum.
2. Time n appends with reserve and without reserve for n = 200 000 in Go, or compare `[]` growth in Python. Write the two times.
3. Explain in six sentences why a real-time loop can still fail if one append resizes.

#### Advanced practical tasks

1. Implement append and a running total of copied elements. Print total / n after n = 1, 2, 4, ..., 1024.
2. Read the Go `append` documentation or the Python list resize comment in the source. Write ten STE sentences about amortized append.

---

## Geometric resizing cost

**Geometric resizing cost** is the total work of all resizes when you grow from empty to n.

For factor 2, the last copy moves about n elements. The copy before that moves about n/2. The sum is less than 2n. The total resize cost is Θ(n). That total is why amortized append is Θ(1).

For factor 1.5, the sum is still Θ(n). The constant is larger than the constant for factor 2. You copy a larger fraction of the history.

For factor 1 (growth by +1), the sum is Θ(n²).

```text
factor 2, grow to n = 16
  copy sizes: 1, 2, 4, 8
  sum:        15
```

Space after a grow can waste capacity. If length is n and capacity is 2n, you can waste n cells. That waste is the price of geometric growth. A factor of 1.5 wastes a smaller fraction after a grow. A factor of 2 wastes more fraction and copies less often.

When you copy, you allocate a new block. The old block becomes free. A garbage collector or an allocator reclaims the old block. For a short time you hold two blocks. Peak space during grow is about length plus new capacity.

### Questions

#### Theoretical questions

1. What is geometric resizing cost for a full growth from empty to n?
2. Why is the total copy cost Θ(n) for doubling?
3. Why does factor 1 give Θ(n²) total copy cost?
4. What space waste can you see just after a doubling grow?
5. Why can peak space during grow be about two buffers?

#### Easy practical tasks

1. List copy sizes for doubling from 1 to 32.
2. Sum those copy sizes. Compare the sum with 32 and with 64.
3. Make a table: "Factor", "Total copy class". Rows: +1, ×1.5, ×2.
4. Draw old buffer and new buffer during a grow from capacity 4 to 8.

#### Medium practical tasks

1. Write a function that returns the total copies for n appends with factor 2, start capacity 1.
2. Compare wasted cells just after grow for factor 2 and for factor 1.5 at length 100. Use integer capacities.
3. Explain in six sentences the trade-off between waste and copy frequency.

#### Advanced practical tasks

1. Simulate peak live bytes during n appends of 8-byte elements. Include two buffers at each grow.
2. Write a report: pick a factor for a memory-tight program and a factor for a throughput-tight program. Give numbers.

---

## Shrinking (optional)

**Shrink** reduces capacity when length becomes much smaller than capacity. Shrink allocates a smaller buffer and copies length elements. Shrink is optional. Many libraries do not shrink on every delete.

If you never shrink, a list that grew to a large peak keeps a large buffer. Space stays high after you delete most elements.

If you shrink too early, a pattern of delete and append causes many copies. That pattern is **thrashing**.

A common rule is to shrink when length is less than capacity / 4. The new capacity can be capacity / 2. That rule leaves spare cells. A later append does not grow immediately.

```text
optional shrink rule
  if length * 4 <= capacity and capacity > min_capacity:
      new_capacity = max(min_capacity, capacity / 2)
      copy length elements
```

Go slices do not shrink when you reslice to a short length. The capacity can stay large. You copy into a new slice when you want a tight buffer: `append([]T(nil), s...)`.

Python lists can shrink in some resize paths. Do not rely on shrink without a test.

This handbook treats shrink as a policy. You document the policy. You test the policy with grow and shrink cycles.

### Questions

#### Theoretical questions

1. What does shrink do?
2. Why is shrink optional?
3. What is thrashing in this section?
4. Why can a 1/4 rule be safer than shrink at 1/2?
5. Why can a short Go subslice still hold a large capacity?

#### Easy practical tasks

1. Draw a buffer of capacity 16 and length 3. Mark unused cells.
2. Write five sentences about why a peak size can waste space later.
3. Make a table: "Action", "Capacity change". Rows: append when full, delete one element, shrink when 1/4 full.
4. Write one sentence that states a shrink policy you would use.

#### Medium practical tasks

1. Implement delete-from-end plus the 1/4 shrink rule in a tiny dynamic array.
2. Run a cycle: grow to 64, delete to 8, append to 64 again. Count resizes with shrink and without shrink.
3. In Go, show `cap` after you reslice to length 1. Then copy to a new slice and show `cap` again.

#### Advanced practical tasks

1. Design two shrink policies. Write a test that causes thrashing on the aggressive policy and not on the 1/4 policy.
2. Read how one standard list shrinks (or does not shrink). Write ten STE sentences. Cite the source.

---

## Iteration and random access

A dynamic array keeps the array operations.

**Random access** is `get(i)` and `set(i, x)` for 0 ≤ i < length. The cost is Θ(1). The address formula is the same as for a fixed array.

**Iteration** visits index 0, then 1, then 2, until length − 1. Iteration is Θ(n). Locality is good because the live elements are contiguous.

```go
for i := 0; i < len(s); i++ {
	use(s[i])
}
for _, v := range s {
	use(v)
}
```

```python
for i in range(len(a)):
    use(a[i])
for v in a:
    use(v)
```

Do not iterate past length. Do not use capacity as the loop bound. Unused cells are not live elements.

Insert and delete in the middle stay Θ(n). Append stays amortized Θ(1) with geometric growth. Those costs are why a dynamic array is the default list in many languages.

When you iterate and append to the same array, indexes can change. In some languages the iterator fails. In others the loop bound is stale. Prefer a simple index loop with a clear rule: do not append in a range loop unless you know the language rule.

### Questions

#### Theoretical questions

1. Why is `get(i)` Θ(1) on a dynamic array?
2. What is the usual iteration order?
3. Why must a loop use length, not capacity?
4. What is the cost of insert in the middle?
5. Why is a dynamic array the default list in many languages?

#### Easy practical tasks

1. Iterate a small dynamic array and print each index and value.
2. Write a `get` and a `set` that check bounds against length.
3. Make a table: "Operation", "Cost class". Rows: get(i), iterate all, append, insert at 0.
4. Draw length 3 in a capacity 8 buffer. Mark which cells a full iteration visits.

#### Medium practical tasks

1. Time random `get` on a dynamic array of 1 000 000 integers. Then time a full iteration. Write the two times.
2. Write insert at index 0 of a dynamic array. Count moves for n = 10 000.
3. Explain in six sentences what can go wrong if you append inside a range loop.

#### Advanced practical tasks

1. Implement a mini vector: append, get, set, len, iterate, insert, delete. Use doubling. Write a short complexity table.
2. Compare iteration of a dynamic array and of a linked list for n = 1 000 000 if you can. Write times and one cause.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe the full path of append from a full buffer to a completed write. Name allocate, copy, and length.
2. How do geometric growth and optional shrink work together as a policy pair?
3. Why does reserve change the worst case of the next n appends?
4. A teammate says "a slice is an array". Which facts do you use to make the sentence precise?
5. When do you still use a fixed-size array instead of a dynamic array?

#### Easy practical tasks

1. Write a one-page cheat sheet: dynamic array, length, capacity, append, resize, amortized, reserve, shrink, get(i).
2. Draw three snapshots: empty, after 3 appends with capacity 4, after a 4th append that doubles.
3. List five language names for this structure.
4. Write the complete complexity statement for append with doubling.

#### Medium practical tasks

1. Implement a dynamic array of integers without the language list grow (use a fixed buffer field that you replace). Support append and get.
2. Build with reserve for n and without reserve. Count allocations if you can log them, or count grows in your implementation.
3. Write a shrink-on-clear design: `clear` keeps capacity, `compact` copies to a tight buffer. Document both.

#### Advanced practical tasks

1. Implement factor-2 growth and 1/4 shrink. Test a random mix of append and pop for 100 000 operations. Record max capacity and copy count.
2. Read the implementation notes for Go `append` or CPython `list_resize`. Rewrite the growth rule in STE. Cite the file or page.
