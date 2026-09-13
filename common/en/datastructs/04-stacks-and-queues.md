# 4. Stacks and Queues

## Description

A stack and a queue are abstract data types. A stack removes the last element that you inserted. A queue removes the first element that you inserted. This topic explains array and list implementations, ring buffers, deques, the call stack, and the use of a queue in BFS and a stack in DFS.

Complete this topic after arrays and linked lists.

Use one term for each concept. **LIFO** means last in, first out. **FIFO** means first in, first out. A **deque** is a double-ended queue. You insert and delete at both ends. A **ring buffer** is a fixed array that wraps indexes with a modulus. The **call stack** is the implicit stack that a language uses for function calls. An **explicit stack** is a stack that you store in your data.

---

## LIFO stack: array vs list

A **stack** is an ADT. The main operations are **push** (insert at the top), **pop** (remove the top), and **peek** (read the top and do not remove it). Empty reports whether the stack holds no element.

LIFO is the order rule. The last push is the first pop.

```text
push 10, push 20, push 30
stack top → 30
            20
            10
pop → 30, now top is 20
```

**Array implementation.** Store elements in a dynamic array. The top is the last index. Push is append. Pop deletes the last element. Peek reads the last element. With geometric growth, push is amortized Θ(1). Pop and peek are Θ(1). Locality is good.

**List implementation.** Store elements in a singly linked list. The top is the head. Push inserts at the head. Pop deletes at the head. Peek reads the head. Each operation is Θ(1). Locality is poor. Each node stores a pointer.

Both implementations are stacks if they do the same operations. The costs and the locality differ.

Do not pop an empty stack. Define an error or a clear empty result in the ADT.

A stack does not give get(i) for an arbitrary i in the ADT. If you need index access, you selected the wrong ADT.

### Questions

#### Theoretical questions

1. What does LIFO mean?
2. What are push, pop, and peek?
3. How does an array store the top?
4. How does a list store the top?
5. Why is get(i) not a stack operation?

#### Easy practical tasks

1. Start empty. Write the stack after push 1, push 2, pop, push 3.
2. Draw an array stack of three integers. Mark the top index.
3. Draw a list stack of the same three integers. Mark the head as top.
4. Make a table: "Implementation", "Push cost", "Locality". Add two rows.

#### Medium practical tasks

1. Write the steps of push and pop for both implementations. Include the empty-stack case.
2. Explain in six sentences when you prefer the array stack.
3. Write four sentences about amortized push on a growing array stack.

#### Advanced practical tasks

1. Implement both stack implementations with the same public names. Test five push/pop sequences.
2. Time n pushes and n pops on both implementations for a large n. Write the two times and one locality comment.

---

## FIFO queue: ring buffer vs list

A **queue** is an ADT. The main operations are **enqueue** (insert at the back), **dequeue** (remove the front), and **front** (read the front and do not remove it). Empty reports whether the queue holds no element.

FIFO is the order rule. The first enqueue is the first dequeue.

```text
enqueue 10, enqueue 20, enqueue 30
front 10, then 20, then 30 at the back
dequeue → 10, now front is 20
```

**List implementation.** Use a singly linked list with head and tail. Enqueue inserts at the tail. Dequeue deletes at the head. Both are Θ(1). Do not enqueue at the head and dequeue at the tail. Delete at the tail is Θ(n) on a singly linked list.

**Ring buffer implementation.** Use a fixed array of capacity m. Store a head index and a length, or a head index and a tail index. Enqueue writes at the tail index, then moves tail with (tail + 1) mod m. Dequeue reads at the head index, then moves head with (head + 1) mod m.

```text
capacity 4
cells:  [10] [20] [  ] [  ]
head=0  tail=2  length=2

enqueue 30
cells:  [10] [20] [30] [  ]
head=0  tail=3

dequeue
cells:  [  ] [20] [30] [  ]
head=1  tail=3
```

When head or tail passes the last cell, the index wraps to 0. That wrap is the ring.

A full ring buffer cannot enqueue until you dequeue or you grow. A dynamic ring can allocate a larger array and copy in queue order. Copy must start at head, not at cell 0, if the live elements wrap.

Array queue without a ring is a poor default. Dequeue at index 0 shifts all elements. That dequeue is Θ(n).

Do not dequeue an empty queue. Define an error in the ADT.

### Questions

#### Theoretical questions

1. What does FIFO mean?
2. Why does a list queue use head and tail?
3. What is a ring buffer?
4. Why do indexes use modulus m?
5. Why is dequeue at index 0 a poor array queue?

#### Easy practical tasks

1. Start empty. Write the queue after enqueue 1, enqueue 2, dequeue, enqueue 3.
2. Draw a ring of capacity 4 with two live elements that wrap across cell 0.
3. Make a table: "Implementation", "Enqueue", "Dequeue". Write Θ.
4. Write five sentences that define a FIFO queue.

#### Medium practical tasks

1. Write head, tail, and cells after each of: enqueue A, enqueue B, dequeue, enqueue C on capacity 3.
2. Explain in six sentences how you copy a wrapped ring into a larger buffer.
3. Write why a list queue must not dequeue from the tail of a singly linked list.

#### Advanced practical tasks

1. Implement a ring buffer queue with full and empty checks. Test wrap-around.
2. Implement a list queue with head and tail. Compare unused space with a ring of fixed capacity.

---

## Deque

A **deque** (double-ended queue) is an ADT. You insert and delete at the front and at the back.

Typical operations: push_front, push_back, pop_front, pop_back, peek_front, peek_back.

A deque can do stack work and queue work. A stack uses one end only. A queue uses two ends, one for insert and one for delete.

```text
deque:  front [10, 20, 30] back
push_front 5  →  [5, 10, 20, 30]
pop_back      →  [5, 10, 20]
```

**Doubly linked list implementation.** Insert and delete at both ends are Θ(1) if you store head and tail, or if you use sentinels.

**Ring buffer implementation.** Head and tail move in both directions. Decrement uses (i − 1 + m) mod m. All four end operations are Θ(1) when the buffer is not full.

**Dynamic array at one end only is not enough.** push_front on a normal dynamic array shifts all elements. That cost is Θ(n). A ring or a block deque avoids that shift.

A common block deque stores a dynamic array of blocks. Each block is a small array. This topic only needs the idea: the ADT is four end operations. The implementation must make those four operations cheap.

Do not mix deque with dequeue. Dequeue is the queue remove operation. Deque is the ADT name.

### Questions

#### Theoretical questions

1. What operations does a deque add beyond a queue?
2. How can a deque act as a stack?
3. Why is a doubly linked list a natural deque?
4. Why does a normal dynamic array make push_front expensive?
5. How does a ring support both ends?

#### Easy practical tasks

1. Start empty. Apply push_back 1, push_front 2, pop_back. Write the remaining elements from front to back.
2. Draw a deque of three values. Mark front and back.
3. Make a table: "ADT", "Ends that you use". Rows: stack, queue, deque.
4. Write five sentences that define a deque.

#### Medium practical tasks

1. Write the four end operations as pointer steps on a doubly linked list with sentinels.
2. Explain in six sentences why a ring increment and decrement both need modulus.
3. Write a short note: when do you use a deque instead of a stack plus a queue?

#### Advanced practical tasks

1. Implement a deque with a ring buffer. Test wrap on both ends.
2. Read a short note on a block deque. Write eight STE sentences about why blocks reduce shift cost.

---

## Call stack vs explicit stack

The **call stack** is the stack that the language runtime uses. Each **call frame** holds the return address, local variables, and arguments for one function call. A call pushes a frame. A return pops a frame.

You do not name this stack in your program. The compiler and the runtime manage it. Recursion uses the call stack. A deep recursion can overflow the call stack.

An **explicit stack** is a stack that you create. You push values that you choose. You pop them in your code. Topic 7 uses an explicit stack to walk a tree without recursion. Topic 11 uses an explicit stack for DFS.

```text
function f calls g, g calls h
call stack (top is h):
  h frame
  g frame
  f frame

explicit stack in your data:
  you push node C
  you push node B
  you push node A
  pop → A
```

The two stacks are the same ADT idea (LIFO). They store different things. The call stack stores frames. The explicit stack stores your elements.

A tail call can reuse a frame in some languages. Do not assume that reuse. Count depth as one frame per active call unless you measured otherwise.

Convert recursion to a loop plus an explicit stack when you need control of memory, or when the language stack is small.

Do not confuse the call stack with a stack data structure in the heap. The call stack has a limited size that the operating system or the runtime sets.

### Questions

#### Theoretical questions

1. What does one call frame hold?
2. What is an explicit stack?
3. How does recursion use the call stack?
4. Why can deep recursion fail?
5. Why convert recursion to an explicit stack?

#### Easy practical tasks

1. Draw the call stack for f → g → h before any return.
2. Draw an explicit stack after push A, push B, pop.
3. Make a table: "Stack kind", "Who manages it", "What it stores". Add two rows.
4. Write five sentences that contrast the two stacks.

#### Medium practical tasks

1. Write a recursive countdown and an iterative countdown that uses an explicit stack of integers. Write which stack each version uses.
2. Explain in six sentences why a large local array in a frame can overflow the call stack.
3. List two algorithms in later topics that use an explicit stack.

#### Advanced practical tasks

1. Write a recursive list walk and an iterative list walk with an explicit stack. Compare maximum extra memory.
2. Find the call-stack size limit in a language that you know (documentation). Write the limit and one safe depth estimate.

---

## BFS uses a queue; DFS uses a stack

**Breadth-first search (BFS)** visits vertices in order of distance from a start, when each edge has length 1. BFS uses a **queue**. You enqueue neighbors. You dequeue the next vertex to visit. The FIFO order visits layer 0, then layer 1, then layer 2.

**Depth-first search (DFS)** goes as deep as possible along one path, then backtracks. DFS uses a **stack**. Recursion uses the call stack. An iterative DFS uses an explicit stack.

```text
graph:
  A — B — D
  |
  C

BFS from A, queue: A, then B and C, then D
order can be: A, B, C, D

DFS from A, stack (or recursion): A, then B, then D, then back, then C
order can be: A, B, D, C
```

The neighbor order can change the exact visit list. The ADT choice does not change: BFS needs FIFO. DFS needs LIFO.

Do not use a stack for unweighted shortest path. A stack does not visit by distance. Topic 11 states the full algorithms. This section only ties the ADTs to the two walks.

A deque can implement both if you use it as a queue or as a stack. The meaning comes from which ends you use.

Mark a vertex when you first put it in the structure, or when you first take it out. The two mark times are not the same. Topic 11 explains the usual rule for BFS: mark when you enqueue.

### Questions

#### Theoretical questions

1. Which ADT does BFS use?
2. Which ADT does DFS use?
3. Why does FIFO visit vertices by unweighted distance?
4. Why does LIFO go deep first?
5. Why is a stack the wrong structure for unweighted shortest path?

#### Easy practical tasks

1. Run BFS on the A–B–C–D picture. Write the queue after each dequeue. Use neighbor order B then C from A.
2. Run DFS on the same picture with an explicit stack. Write the stack after each pop. Use the same neighbor order.
3. Make a table: "Walk", "ADT", "Order idea". Add two rows.
4. Write five sentences that tie BFS to a queue and DFS to a stack.

#### Medium practical tasks

1. Draw a 2 × 2 grid of cells. Write a BFS visit order from the top-left cell. Write a DFS visit order.
2. Explain in six sentences how a deque can serve both walks.
3. Write why mark-on-enqueue and mark-on-dequeue can enqueue the same vertex two times if you mark late.

#### Advanced practical tasks

1. Implement BFS with a queue and DFS with an explicit stack on a small graph. Print visit orders.
2. Write a one-page note: same graph, same start, different ADTs, different orders. State one problem that needs BFS and one problem that needs DFS.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do LIFO and FIFO decide the implementation that you pick for array, ring, and list?
2. When is a deque the right ADT, and when is a stack or a queue enough?
3. How does the call stack relate to recursive DFS?
4. Why must a ring copy in queue order when it grows?
5. Which end operations stay Θ(1) for stack, queue, and deque in the designs in this topic?

#### Easy practical tasks

1. Write a one-page cheat sheet: LIFO, FIFO, push, pop, enqueue, dequeue, ring, deque, call stack, BFS, DFS.
2. Draw one array stack, one ring queue, and one list deque. Label the ends.
3. Make a table: "Need", "ADT". Add five rows (undo, waiting line, both ends, recursion, layer walk).
4. List four defects: pop empty, dequeue empty, ring full, wrong end on a list queue.

#### Medium practical tasks

1. Implement a stack and a queue with the same ring buffer code and different end rules. Write the two rule sets.
2. Convert a short recursive walk to an explicit stack. Write the mapping from frames to pushed values.
3. Write a test plan for wrap-around, empty, full, and one-element ring queues.

#### Advanced practical tasks

1. Implement BFS and DFS on a grid with walls. Use a queue and a stack. Do not print a solution path in this handbook. Write how you record parent cells.
2. Read one real queue or deque in a standard library. Write whether it uses a ring, blocks, or a list, from the documentation.
