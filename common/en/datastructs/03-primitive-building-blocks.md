# 3. Primitive Building Blocks

## Description

A data structure is built from bits, words, primitive types, and links. This topic explains those building blocks. You learn contiguous memory and linked nodes. You also learn the idea of stack allocation and heap allocation. Complete this topic before you implement arrays and lists.

Use one term for each concept. A **bit** is a 0 or a 1. A **byte** is 8 bits. A **word** is the natural unit that the processor moves. A **pointer** or **reference** is a value that names an address. The exact rules depend on the language.

---

## Bits, bytes, words

A computer stores information as bits. One bit has two values: 0 and 1.

A **byte** is a group of 8 bits. Memory addresses usually name bytes. One address is one byte unless the architecture says a different rule.

A **word** is a group of bytes that the processor prefers. Common word sizes today are 4 bytes (32-bit) and 8 bytes (64-bit). A word can hold an integer or an address on many machines.

```text
1 bit    = 0 or 1
1 byte   = 8 bits
1 word   = 4 bytes or 8 bytes on many machines
1 KiB    = 1024 bytes
```

You count memory in bytes when you talk about space. You count bits when you talk about flags and compact codes.

An integer of 32 bits can represent 2³² different patterns. An integer of 64 bits can represent 2⁶⁴ different patterns. A larger width uses more space and can hold a larger range.

Do not mix **bit** and **byte** in speech. "A 32-bit integer" is a width in bits. "A 4-byte integer" is the same width in bytes on a common machine.

Network data and files also use bytes. A data structure in memory uses the same byte model.

### Questions

#### Theoretical questions

1. What is a bit?
2. What is a byte?
3. What is a word on a common machine?
4. Why do memory addresses usually name bytes?
5. How many different patterns can 32 bits represent?

#### Easy practical tasks

1. Convert 16 bits to bytes. Write the integer.
2. Make a table: "Unit", "Size in bits". Rows: bit, byte, 32-bit word, 64-bit word.
3. Write five sentences that use bit, byte, and word correctly.
4. An array holds 10 integers of 4 bytes each. Write the data size in bytes. Ignore extra headers.

#### Medium practical tasks

1. In Go, print `bits.UintSize` or compare `32 << (^uint(0) >> 63)`. Write the word size that you see. In Python, print `sys.maxsize` and write what it tells you.
2. Count the bytes of a string of 5 ASCII letters in one language. Write whether a header exists.
3. Explain in six sentences why a 64-bit address can name more bytes than a 32-bit address.

#### Advanced practical tasks

1. Write a small program that packs 8 boolean flags into one byte with shifts and masks. Document each bit.
2. Read the integer sizes in the language specification. Make a table of type name, bits, and minimum range. Cite the source.

---

## Primitive types and their size

A **primitive type** is a built-in type that the language gives. Common primitive types are integers, floating-point numbers, booleans, and bytes.

Each primitive type has a size. The size is a number of bits or bytes. The size can be fixed. The size can depend on the machine.

```text
typical sizes (many languages, many machines)
  bool           1 byte in a struct, or 1 bit in a packed bitset
  8-bit integer  1 byte
  32-bit integer 4 bytes
  64-bit integer 8 bytes
  32-bit float   4 bytes
  64-bit float   8 bytes
```

Go has `int` and `uint`. The size of `int` is 32 bits or 64 bits. The size matches the machine word in usual builds. Go also has `int32` and `int64` when you need a fixed size.

Python has `int` that grows. A small Python integer is an object with a header. The header is larger than 4 bytes. Do not treat a Python `int` as a 32-bit cell.

A character encoding is not a data structure. A character still uses bytes. UTF-8 uses 1 to 4 bytes per code point.

When you estimate space, multiply the element size by n. Then add headers, pointers, and unused capacity. The primitive size is only the first term.

### Questions

#### Theoretical questions

1. What is a primitive type?
2. Why does the size of `int` change in some languages?
3. Why is a Python `int` not a 4-byte cell?
4. What extra terms do you add after n times element size?
5. Why do you use a fixed-width integer in a binary file?

#### Easy practical tasks

1. List five primitive types in Go or in Python. Write a typical size for each.
2. Make a table: "Type", "Size in bytes", "One use". Add four rows.
3. Estimate the data bytes of 1 000 values of type 64-bit integer.
4. Write one sentence that contrasts `int` and `int64` in Go.

#### Medium practical tasks

1. In Go, use `unsafe.Sizeof` on `int`, `int32`, `int64`, and `float64`. Write the sizes. In Python, use `sys.getsizeof` on `0`, `1`, and a large integer. Write the sizes.
2. Estimate an array of 1 000 000 `float64` values in bytes and in mebibytes.
3. Explain why a boolean in a struct can use a full byte even if the value needs one bit.

#### Advanced practical tasks

1. Measure object size for a list of 10 000 small integers in Python, or a slice of 10 000 `int` in Go. Compare with n times 8 bytes.
2. Write a one-page note: when you pick a small integer type to save space, and when that choice is a mistake.

---

## Pointers / references

A **pointer** is a value that holds an address. The address names a cell or an object in memory.

A **reference** is a language name for a value that refers to an object. Java and Python use references. Go uses pointers as explicit types (`*T`). The idea is the same: the program stores a name of a location, not a full copy of a large object.

```text
variable p  →  address 0x100
memory 0x100 holds the integer 42
```

A pointer has a size. On a 64-bit machine the pointer is often 8 bytes. A list of n nodes stores n pointers if each node has one next link. Those pointers are extra space.

A **null** pointer (or `nil` in Go, `None` in Python) means "no object". A walk of a list stops at null.

When you copy a pointer, you do not copy the object. Two pointers can name the same object. A change through one pointer is visible through the other pointer. That sharing is useful. That sharing is also a source of errors.

An **array of pointers** is not the same as an array of values. The values can sit far from each other. Locality can be poor.

Do not dereference a null pointer. The program crashes or raises an error.

### Questions

#### Theoretical questions

1. What does a pointer store?
2. How is a reference similar to a pointer?
3. What extra space does one next-pointer add to a node?
4. What does a null pointer mean?
5. What happens when two pointers name the same object and you change the object?

#### Easy practical tasks

1. Draw a box for a node. Draw an arrow from a pointer variable to the box. Label address and value.
2. Write four sentences that explain copy of a pointer versus copy of an integer.
3. Make a table: "Language", "Null name". Rows: Go, Python, and one more language that you know.
4. A node holds one 8-byte integer and one 8-byte pointer. Write the data size of the node. Ignore alignment.

#### Medium practical tasks

1. In Go, make a pointer to an integer, change the integer through the pointer, and print the original variable. In Python, put a list in two variables and append to one. Write what you see.
2. Draw three nodes in a chain. Write the number of next pointers that are not null.
3. Explain in six sentences why an array of pointers to large objects can have poor locality.

#### Advanced practical tasks

1. Write a small program that creates a cycle of two nodes. Each node points to the other. Draw the picture. Write why a naive walk does not stop.
2. Read the pointer or reference model in the language documentation. Write ten STE sentences. Cite the page.

---

## Contiguous memory vs linked nodes

**Contiguous memory** is a block of cells with consecutive addresses. An array uses contiguous memory. Index i is at a fixed offset from the start.

**Linked nodes** are separate cells. Each node holds a value and one or more pointers to other nodes. A singly linked node holds a next pointer. A doubly linked node holds next and previous pointers.

```text
contiguous:  [ a | b | c | d ]     addresses in a row
linked:      (a)→(b)→(c)→(d)→null  addresses can be anywhere
```

Contiguous memory gives:

- index in constant time
- good locality
- a need to shift elements when you insert in the middle
- a fixed block, or a new block when you grow

Linked nodes give:

- insert and delete at a known node without a shift of the other values
- a walk from the head to index i
- extra space for pointers
- poor locality

Alignment can add unused bytes inside a node. The unused bytes make the node larger than the sum of the field sizes.

You will implement both layouts. Select the layout from the operations. Do not select the layout from habit.

### Questions

#### Theoretical questions

1. What is contiguous memory?
2. What does a linked node store besides the value?
3. Why is index i fast in a contiguous block?
4. Why is insert in the middle of a contiguous block expensive?
5. What is alignment in one sentence?

#### Easy practical tasks

1. Draw four integers in a contiguous block. Write the address of each cell if the start is 1000 and each integer is 4 bytes.
2. Draw the same four integers as linked nodes. Show next arrows.
3. Make a table: "Property", "Contiguous", "Linked". Add four properties.
4. Write which layout fits "get the 500th element often".

#### Medium practical tasks

1. Write a struct or class for a singly linked node in Go or in Python. Fields: value and next.
2. Estimate space for 1 000 integers of 8 bytes in an array and in nodes with one 8-byte next pointer. Ignore headers.
3. Explain in six sentences when you accept extra pointer space to avoid shifts.

#### Advanced practical tasks

1. Measure `unsafe.Sizeof` of a Go struct `{int64; *Node}` or `sys.getsizeof` of a small Python object with two fields. Write padding if you see it.
2. Write a one-page comparison of a file buffer (contiguous) and a chain of packet buffers (linked). Name operations for each.

---

## Stack vs heap allocation (language-dependent, conceptual)

**Allocation** is the act of reserving memory for a value.

The **call stack** is a region that holds frames. Each function call pushes a frame. The frame holds local variables that the language puts on the stack. When the function returns, the frame goes away. Stack allocation is fast. The size of a stack frame is usually known at compile time. A large array in a frame can overflow the stack.

The **heap** is a region for values that live after the call, or that have a size that you know only at run time. The allocator finds a free block. The program later frees the block, or a garbage collector frees the block. Heap allocation is slower than stack allocation. Heap allocation is more flexible.

```text
call f()
  stack frame of f: small locals
  heap:             a dynamic array that f creates and returns
return from f
  stack frame of f is gone
  heap array still lives if a caller holds a reference
```

The rules depend on the language:

- Go: many values escape to the heap. The compiler decides. You do not call `malloc`. A garbage collector frees unused heap objects.
- Python: almost every object lives on the heap. Names are references.
- C: you call `malloc` and `free`. You own the lifetime.

This handbook uses **stack** for short-lived frame storage and **heap** for the allocator region. Do not mix this stack with the stack ADT. The stack ADT is a later topic. The call stack is a memory region.

A data structure of n nodes usually lives on the heap. The variable that names the head can live in a frame.

### Questions

#### Theoretical questions

1. What is allocation?
2. What does a call-stack frame hold?
3. When does a stack frame go away?
4. Why does a dynamic structure usually live on the heap?
5. How is the call stack different from the stack ADT?

#### Easy practical tasks

1. Write five sentences that contrast stack allocation and heap allocation.
2. Make a table: "Region", "Lifetime", "Typical content". Two rows: stack, heap.
3. Draw a call of function `g` from `main`. Show two frames.
4. Write where a linked-list node of a long-lived list usually lives.

#### Medium practical tasks

1. In Go, write a function that returns a pointer to a local integer. Run it. Write that the value must live after the return. In Python, write a function that returns a new list. Write where the list lives.
2. Explain in six sentences why a very large local array can be a bad idea.
3. Find the words "escape" or "heap" in Go documentation, or "object" in Python documentation. Write four STE sentences.

#### Advanced practical tasks

1. In Go, build with `-gcflags=-m` on a tiny program. Copy one escape line. Write what it means. In Python, skip that flag and write how reference counting or GC frees a list.
2. Write a one-page note: lifetime of a node from allocate to free in your language. Include who frees the node.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe the path from a bit to a linked node. Name byte, word, primitive type, and pointer.
2. Why does a 64-bit machine make each pointer more expensive than a 32-bit machine?
3. How do contiguous layout and heap allocation work together in a dynamic array?
4. A teammate says that Python has no pointers. Which facts do you use in a reply?
5. Why must you add headers and padding when you estimate real space?

#### Easy practical tasks

1. Write a one-page cheat sheet: bit, byte, word, primitive type, pointer, reference, null, contiguous, node, stack, heap.
2. Draw one picture that shows a stack frame that holds a head pointer and two heap nodes.
3. Estimate bytes for 50 nodes: 4-byte value + 8-byte pointer. Then add 16 bytes of header per node as a guess.
4. Label each item as primitive or structure: `int64`, array of `int64`, pointer to `int64`, linked list.

#### Medium practical tasks

1. Write two tiny programs: one sums a contiguous array, one sums a chain of nodes. Use the same n values. Comment every pointer use.
2. Compare Go `int` size and Python `int` object size on your machine. Write a short table.
3. Explain in eight sentences why a structure of many small heap nodes can lose to one array of primitives.

#### Advanced practical tasks

1. Build a node type and print its size. Then pack the same data into an array of primitives. Compare bytes for n = 10 000.
2. Read one article on memory allocators for your language. Write ten STE sentences about heap cost. Cite the source.
