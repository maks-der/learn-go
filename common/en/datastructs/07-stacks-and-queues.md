# 7. Stacks and Queues

## Description

A stack and a queue are abstract data types for ordered access. This topic explains LIFO and FIFO. You learn array and list implementations, a ring buffer, a deque, the call stack, and a preview of BFS and DFS. Complete this topic after arrays and lists.

Use one term for each concept. **LIFO** means last in, first out. **FIFO** means first in, first out. A **stack** is a LIFO ADT. A **queue** is a FIFO ADT. A **deque** is a double-ended queue. A **ring buffer** is a fixed buffer with wrap-around indexes.

---

## LIFO stack: array vs list

A **stack** supports push, pop, and peek (top). Push puts an element on the top. Pop removes the top and returns it. Peek returns the top and does not remove it. Empty reports whether the stack has no element.

**LIFO** means the last push is the first pop.

```text
push 1, push 2, push 3
stack top → 3, 2, 1
pop → 3
```

**Array implementation.** Use a dynamic array. Push is append. Pop is delete at the end. Peek is get(length − 1). With geometric growth, push is amortized Θ(1). Pop is Θ(1). Locality is good.

**List implementation.** Use insert at head as push. Use delete at head as pop. Each operation is Θ(1) worst case. Locality is poor. Extra space includes pointers.

```go
// array-backed idea
func (s *Stack) Push(x int) { s.data = append(s.data, x) }
func (s *Stack) Pop() int {
	n := len(s.data) - 1
	x := s.data[n]
	s.data = s.data[:n]
	return x
}
```

Do not pop an empty stack. Return an error or a boolean. Define the empty behavior in the ADT.

A stack does not give get(i) in the ADT. If you need index access, you need a different ADT.

Use the array stack as the default. Use the list stack when you need guaranteed Θ(1) push without a rare resize, or when you already use nodes.

### Questions

#### Theoretical questions

1. What does LIFO mean?
2. What three operations does a stack ADT provide besides empty?
3. How does a dynamic array implement push and pop?
4. How does a singly linked list implement push and pop?
5. Why is the array stack the default?

#### Easy practical tasks

1. Draw a stack after push 10, push 20, pop, push 30. Show the top.
2. Write the ADT contract for stack in four lines.
3. Make a table: "Implementation", "Push cost", "Locality". Two rows: array, list.
4. Write what pop must do when the stack is empty.

#### Medium practical tasks

1. Implement a stack on a dynamic array. Include empty checks.
2. Implement a stack on a linked list. Keep the same public function names.
3. Time n push and n pop on both implementations. Write the two times.

#### Advanced practical tasks

1. Implement a stack that grows a fixed buffer by doubling. Do not use the language list type for storage.
2. Use a stack to check balanced parentheses in a string. Document the algorithm in STE. Do not look up a full solution if you can design it.

---

## FIFO queue: ring buffer vs list

A **queue** supports enqueue, dequeue, and front. Enqueue puts an element at the back. Dequeue removes the element at the front. Front returns the front and does not remove it. FIFO means the first enqueue is the first dequeue.

```text
enqueue 1, enqueue 2, enqueue 3
front → 1, then 2, then 3
dequeue → 1
```

**List implementation.** Enqueue at tail. Dequeue at head. You need head and tail. Each operation is Θ(1). Do not enqueue at head and dequeue at tail on a singly linked list. Dequeue at tail is Θ(n).

**Ring buffer implementation.** Use a fixed array plus two indexes: head and tail, or head and length. When an index moves past the last cell, the index wraps to 0.

```text
capacity 4
cells:  _  10  20  30
head=1  tail=0  (next write at 0)
enqueue 40 → 40  10  20  30
dequeue → 10, head=2
```

A ring buffer enqueue is Θ(1) when the buffer is not full. When the buffer is full, you reject the enqueue, or you allocate a larger buffer and copy in order.

A dynamic array that dequeues from index 0 is a bad queue. Each dequeue shifts n elements. That cost is Θ(n). Do not implement a queue that way.

Go has no ring buffer in the language. You write one, or you use a list. Python `collections.deque` is a block chain that is efficient at both ends.

### Questions

#### Theoretical questions

1. What does FIFO mean?
2. Why does a singly linked queue need a tail pointer?
3. What does wrap-around mean in a ring buffer?
4. Why is dequeue at index 0 of a compact dynamic array a bad queue?
5. What can you do when a ring buffer is full?

#### Easy practical tasks

1. Draw a queue after enqueue 1, enqueue 2, dequeue, enqueue 3. Show front and back.
2. Write the ADT contract for queue in four lines.
3. Draw a ring of four cells with two live elements and wrap-around.
4. Make a table: "Implementation", "Enqueue", "Dequeue". Rows: list + tail, ring buffer, compact array shift.

#### Medium practical tasks

1. Implement a queue with a singly linked list and tail.
2. Implement a ring buffer of fixed capacity. Reject enqueue when full.
3. Time n enqueue and n dequeue on the list queue and on the ring buffer.

#### Advanced practical tasks

1. Implement a growing ring buffer. When full, allocate 2 × capacity and copy from front in FIFO order.
2. Compare Python `deque` or a Go list queue with your ring buffer. Write a short report.

---

## Deque

A **deque** (double-ended queue) supports insert and delete at both the front and the back. The usual operations are push_front, push_back, pop_front, and pop_back. Peek at both ends is optional.

A deque can implement a stack (use one end) and a queue (use opposite ends).

```text
push_back 1, push_back 2, push_front 0
front → 0, 1, 2 ← back
pop_front → 0
pop_back → 2
```

**List implementation.** A doubly linked list with head and tail gives Θ(1) at both ends.

**Block or ring implementation.** A ring buffer can serve as a deque if you wrap both indexes. Python `deque` uses linked blocks of arrays. That design has better locality than one node per element.

A compact dynamic array is cheap at the back and expensive at the front. Do not use it as a general deque.

Use a deque when you need both ends. If you need only one end, use a stack. If you need only FIFO, use a queue. A smaller ADT is easier to check.

### Questions

#### Theoretical questions

1. What operations does a deque add compared with a queue?
2. How do you use a deque as a stack?
3. How do you use a deque as a queue?
4. Why does a doubly linked list fit a deque?
5. Why is a compact dynamic array a poor general deque?

#### Easy practical tasks

1. Draw a deque after push_back 1, push_front 2, pop_back.
2. Write the ADT contract for deque in five lines.
3. Make a table: "ADT", "Ends in use". Rows: stack, queue, deque.
4. Name one language type that is a deque.

#### Medium practical tasks

1. Implement a deque on a doubly linked list. Support all four update operations.
2. Implement a fixed-capacity ring deque. Wrap both ends.
3. Explain in six sentences why Python `deque` uses blocks.

#### Advanced practical tasks

1. Implement a deque that grows. Document the growth of the ring or of the block list.
2. Solve a sliding-window maximum sketch with a deque of indexes (awareness). Write the invariant in STE. Implement a small n example.

---

## Call stack vs explicit stack

The **call stack** is the memory region for function frames. Each call pushes a frame. Each return pops a frame. The call stack is LIFO. Recursion uses the call stack.

An **explicit stack** is a stack ADT that you create. You push your own values. You pop your own values. The values can be indexes, nodes, or frames that you design.

```text
function a calls b, b calls c
call stack (top): c frame, b frame, a frame

explicit stack for the same idea:
  push a, push b, push c
  pop c, pop b, pop a
```

A recursive algorithm uses the call stack. An iterative algorithm can use an explicit stack and do the same walk. The iterative form avoids a deep call-stack overflow. The iterative form also lets you see the stack in the debugger as your data.

The call stack has a limited size. A deep recursion can crash. An explicit stack lives on the heap. The heap limit is larger. You still must not grow without a bound.

Do not mix the two terms. Say **call stack** for frames. Say **stack** or **explicit stack** for the ADT.

A stack ADT does not replace the call stack for return addresses. The language still uses the call stack for calls. Your ADT is extra data.

### Questions

#### Theoretical questions

1. What does the call stack store?
2. What does an explicit stack store?
3. Why can deep recursion fail?
4. Why can an iterative stack walk survive a deeper input?
5. Why must you not mix the two terms?

#### Easy practical tasks

1. Draw three frames on a call stack for `main` → `f` → `g`.
2. Write five sentences that contrast call stack and explicit stack.
3. Make a table: "Kind", "Who pushes", "Limit". Two rows.
4. Write one recursive function and name the call-stack growth in n.

#### Medium practical tasks

1. Write a recursive factorial or a recursive list walk. Then write an iterative version with an explicit stack if the walk needs one, or explain why a loop is enough.
2. Cause a deep recursion on purpose with a large n. Record the error. Then write a loop.
3. Explain in six sentences when you convert recursion to an explicit stack.

#### Advanced practical tasks

1. Implement depth-first walk of a binary tree (or a tiny fake tree) with an explicit stack of nodes. Do not use recursion.
2. Write a one-page note: stack overflow versus heap growth of an explicit stack. Use one measurement if you can.

---

## BFS uses a queue; DFS uses a stack (preview)

This section is a preview. You will study graphs and trees in later topics. Learn the pairing now.

**Breadth-first search (BFS)** visits nodes by distance from a start. BFS uses a **queue**. You enqueue neighbors. You dequeue the next node. FIFO order visits level 0, then level 1, then level 2.

**Depth-first search (DFS)** goes deep along one path. DFS uses a **stack**. Recursion uses the call stack. An iterative DFS uses an explicit stack. LIFO order continues the newest node first.

```text
BFS idea
  enqueue start
  while queue not empty:
      v = dequeue
      enqueue unvisited neighbors of v

DFS idea
  push start
  while stack not empty:
      v = pop
      push unvisited neighbors of v
```

The same graph can use BFS or DFS. The visit order is different. The ADT is different.

BFS finds a shortest path in an unweighted graph. That fact is a later topic. Remember that the queue is the tool.

A maze with "keep one hand on the wall" is closer to DFS. A flood fill by rings is closer to BFS.

Do not implement a full graph library in this topic. Implement a tiny example: a grid of four cells, or a tree of three nodes.

### Questions

#### Theoretical questions

1. Which ADT does BFS use?
2. Which ADT does DFS use?
3. Why does FIFO visit nodes by level?
4. Why does LIFO go deep first?
5. How can DFS use the call stack instead of an explicit stack?

#### Easy practical tasks

1. Draw a small tree: root A, children B and C. Write a BFS order and a DFS order. State your neighbor order.
2. Make a table: "Search", "ADT", "Order idea". Two rows.
3. Write five sentences that explain the preview without graph theory words that you do not know.
4. Label a queue picture as BFS and a stack picture as DFS.

#### Medium practical tasks

1. Implement BFS on a 2 by 2 grid with four neighbors. Print the visit order.
2. Implement DFS on the same grid with an explicit stack. Print the visit order.
3. Compare the two orders in six sentences.

#### Advanced practical tasks

1. Implement both searches on an adjacency-list of a tiny graph with 6 nodes. Write the two orders.
2. Write a short preview of why BFS gives shortest unweighted paths. Use one drawing. Do not write a full proof.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe how LIFO and FIFO decide the implementation at the ends of a sequence.
2. Why can one deque implementation serve as both stack and queue?
3. How do a ring buffer and a linked list both give Θ(1) queue operations?
4. When do you replace a recursive DFS with an explicit stack?
5. A teammate implements a queue with `insert(0, x)` on a dynamic array. Which facts do you use to reject that design?

#### Easy practical tasks

1. Write a one-page cheat sheet: LIFO, FIFO, stack, queue, deque, ring buffer, call stack, BFS, DFS.
2. Draw one array stack and one list queue. Mark the live end or ends.
3. Write empty-stack and empty-queue behavior for your ADT.
4. List four operations of a deque and which ADT they replace.

#### Medium practical tasks

1. Implement stack and queue with the same public test style: empty, one element, many elements, error on empty pop.
2. Convert a small recursive tree print to an explicit-stack print.
3. Time a bad array queue (dequeue at 0 with shift) against a ring buffer for n = 20 000.

#### Advanced practical tasks

1. Implement stack, queue, and deque without language collections (only arrays or your nodes). Write a complexity table.
2. Use a queue for BFS distance on a small unweighted grid. Use a stack for a DFS path dump. Document both.
