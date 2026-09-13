# 6. Singly and Doubly Linked Lists

## Description

A linked list stores elements in nodes. Each node holds a value and one or more pointers. This topic explains head, tail, insert, delete, sentinel nodes, reverse, cycle detection, and cache cost. Complete this topic after arrays. You will compare the two layouts.

Use one term for each concept. A **node** is a cell with a value and links. The **head** is the first node. The **tail** is the last node. A **singly linked list** has a next pointer. A **doubly linked list** has next and previous pointers. A **sentinel** is a dummy node that is not a real element.

---

## Node, head, tail

A **node** holds a value and a link to the next node. The last node holds a null next pointer.

The **head** pointer names the first node. If the list is empty, head is null.

The **tail** pointer names the last node. Tail is optional. If you store tail, append at the end is Θ(1). If you do not store tail, append walks from the head. That walk is Θ(n).

```text
head → [10|•] → [20|•] → [30|null]
                     ↑
                    tail
```

```go
type Node struct {
	Value int
	Next  *Node
}
```

```python
class Node:
    def __init__(self, value, next=None):
        self.value = value
        self.next = next
```

A **doubly linked node** also holds a previous pointer. The head has a null previous pointer. The tail has a null next pointer. You can walk in two directions.

```text
null ← [10] ⇄ [20] ⇄ [30] → null
        ↑               ↑
       head            tail
```

The list object usually stores head, and often tail and length. Length makes `len` a Θ(1) operation. If you do not store length, `len` walks the list.

Do not confuse the head pointer with the first value. The head pointer is a reference. The first value sits inside the first node.

### Questions

#### Theoretical questions

1. What fields does a singly linked node have?
2. What does a null head mean?
3. Why does a tail pointer make append Θ(1)?
4. What extra field does a doubly linked node have?
5. Why does a stored length make `len` Θ(1)?

#### Easy practical tasks

1. Draw a singly linked list of three integers. Label head and tail.
2. Write a node type in Go or in Python.
3. Make a table: "Pointer", "Names which node". Rows: head, tail, next of head.
4. Draw an empty list. Show head and tail as null.

#### Medium practical tasks

1. Build a list of four nodes by hand (no insert function yet). Walk from head and print values.
2. Add a previous field. Set previous links for three nodes. Walk from tail to head.
3. Explain in six sentences when you store tail and when you omit tail.

#### Advanced practical tasks

1. Implement a list struct with head, tail, and length. Write `len` in Θ(1). Keep the three fields consistent in a comment contract.
2. Estimate bytes for n singly linked nodes versus n doubly linked nodes. Use 8-byte values and 8-byte pointers.

---

## Insert/delete at head and tail

**Insert at head** creates a node and points it to the old head. Then head becomes the new node. The cost is Θ(1). If the list was empty, tail also becomes the new node.

**Delete at head** moves head to head.next. The old first node is unused. The cost is Θ(1). If the list becomes empty, tail becomes null.

**Insert at tail** with a tail pointer: the old tail.next becomes the new node. Then tail becomes the new node. The cost is Θ(1). Without a tail pointer, you walk to the last node first. That walk is Θ(n).

**Delete at tail** in a singly linked list is Θ(n). You must find the node before the tail. A tail pointer does not give you the previous node. In a doubly linked list, delete at tail is Θ(1). You use tail.previous.

```text
insert 5 at head
before:  10 → 20
after:    5 → 10 → 20

insert 30 at tail (tail stored)
before:   5 → 10 → 20
after:    5 → 10 → 20 → 30
```

Insert in the middle at a known node is Θ(1) for the link updates. Find of that node is still Θ(n) if you start from the head.

Delete in the middle of a singly linked list needs the previous node. You cannot relink if you only have the target node, unless you copy the next value (a special trick). Prefer to keep the previous node.

### Questions

#### Theoretical questions

1. What pointer updates does insert at head do?
2. Why is delete at tail slow in a singly linked list?
3. Why is delete at tail fast in a doubly linked list?
4. When is insert in the middle Θ(1)?
5. Why do you need the previous node to delete in a singly linked list?

#### Easy practical tasks

1. Draw insert at head of 7 into `1 → 2`.
2. Draw delete at head of `7 → 1 → 2`.
3. Write the steps of insert at tail when tail is stored.
4. Make a table: "Operation", "Singly + tail", "Doubly + head/tail". Add insert head, insert tail, delete head, delete tail.

#### Medium practical tasks

1. Implement insert_head, delete_head, and insert_tail on a singly linked list with tail.
2. Implement delete_tail on a singly linked list. Show the walk to the previous node.
3. Time n insert_head versus n insert_tail without a tail pointer. Use a moderate n.

#### Advanced practical tasks

1. Implement a doubly linked list with insert and delete at both ends in Θ(1).
2. Write tests for empty, one node, and many nodes. Include tail and length after each operation.

---

## Dummy / sentinel nodes

A **sentinel** (dummy node) is a node that is not an element. The sentinel simplifies empty-list cases.

A **head sentinel** sits before the first real element. Head of the list object always names the sentinel. The first real element is sentinel.next. An empty list has sentinel.next == null (or sentinel.next == tail sentinel).

A **tail sentinel** can sit after the last real element. Some lists use two sentinels. Some circular lists use one sentinel that links to itself.

```text
head sentinel → [10] → [20] → null

empty list:
head sentinel → null
```

Without a sentinel, insert at head and insert after a node are different functions. With a head sentinel, insert after a node covers insert at the front: you insert after the sentinel.

Delete also becomes uniform. You always delete the node after some predecessor. The predecessor of the first element is the sentinel.

Sentinels use extra space. The extra space is one or two nodes. The gain is fewer branches in the code. Fewer branches reduce empty-list bugs.

A circular list can use a sentinel so that no next pointer is null. A walk stops when you return to the sentinel.

Do not count the sentinel as a live element. Length does not include the sentinel.

### Questions

#### Theoretical questions

1. What is a sentinel node?
2. How do you find the first real element when a head sentinel exists?
3. Why can insert after a node cover insert at the front?
4. Does length include the sentinel?
5. What empty-list bug does a sentinel help you avoid?

#### Easy practical tasks

1. Draw a list with a head sentinel and two real nodes.
2. Draw the empty list with only a head sentinel.
3. Write five sentences that explain why sentinels reduce branches.
4. Make a table: "List state", "sentinel.next". Rows: empty, one element, two elements.

#### Medium practical tasks

1. Implement insert_after(pred, x) on a list with a head sentinel. Use it to insert at the front and in the middle.
2. Implement delete_after(pred). Delete the first real element by delete_after(sentinel).
3. Compare code size of insert_head with sentinel and without sentinel. Write the two versions.

#### Advanced practical tasks

1. Implement a circular doubly linked list with one sentinel. Insert and delete at both ends. Walk until you return to the sentinel.
2. Write a short report: when a sentinel is worth the extra node, and when a null head is simpler.

---

## Reverse a list

**Reverse** changes the direction of the next pointers. The old head becomes the tail. The old tail becomes the head.

An **iterative reverse** uses three pointers: previous, current, and next. You walk the list one time. At each node you set current.next = previous. Then you move forward. The cost is Θ(n). The extra space is Θ(1).

```text
before:  1 → 2 → 3 → null
after:   3 → 2 → 1 → null
```

```text
steps (idea)
  prev = null
  curr = head
  while curr != null:
      next = curr.next
      curr.next = prev
      prev = curr
      curr = next
  head = prev
```

A **recursive reverse** reverses the tail, then points the old tail to the old head. Recursion uses the call stack. The extra space is Θ(n). Prefer the iterative reverse unless you study recursion.

If you store tail, swap head and tail after the pointer reverse, or set tail to the old head.

A doubly linked list reverse can swap next and previous on each node, or swap head and tail and walk. Still Θ(n).

Reverse does not copy values into a new array unless you choose that method. An array method uses Θ(n) extra space. The in-place pointer method uses Θ(1) extra space.

### Questions

#### Theoretical questions

1. What does reverse do to head and tail?
2. Why is iterative reverse Θ(n) time and Θ(1) extra space?
3. Why does recursive reverse use Θ(n) extra space?
4. What three pointers does the iterative method use?
5. How is an array copy different from an in-place pointer reverse?

#### Easy practical tasks

1. Draw reverse of `1 → 2 → 3`.
2. Write the iterative loop in comments only. Name prev, curr, next.
3. Make a table: "Method", "Extra space". Rows: iterative, recursive, copy to array.
4. Write the list after reverse of a one-node list and of an empty list.

#### Medium practical tasks

1. Implement iterative reverse on a singly linked list. Update tail if you store tail.
2. Implement reverse by push of values onto a dynamic array and rebuild. Compare extra space.
3. Walk a list before and after reverse. Write a test that checks the sequence.

#### Advanced practical tasks

1. Reverse a sublist from position left to right (1-based or 0-based, document the rule). Keep the rest of the list.
2. Write recursive reverse. Show that a long list can overflow the call stack. Then use the iterative method.

---

## Cycle detection (Floyd)

A **cycle** exists when a walk of next pointers returns to a node that you already visited. A correct list of finite length has no cycle. A bug or a graph edge can create a cycle.

If a cycle exists, a naive walk does not reach null. The program loops forever.

**Floyd’s algorithm** (tortoise and hare) uses two pointers. The slow pointer moves one node per step. The fast pointer moves two nodes per step. If there is no cycle, the fast pointer reaches null. If there is a cycle, the fast pointer meets the slow pointer inside the cycle.

```text
slow = head
fast = head
while fast != null and fast.next != null:
    slow = slow.next
    fast = fast.next.next
    if slow == fast:
        return has cycle
return no cycle
```

The time is Θ(n). The extra space is Θ(1). You do not store a visited set.

A visited set also detects a cycle. You store each node identity. Extra space is Θ(n). Floyd’s algorithm is better when you want constant extra space.

Floyd’s algorithm can also find the start of the cycle. After a meet, move one pointer to the head. Move both pointers one step at a time. They meet at the cycle start. You can learn that second phase later. This topic requires the detection idea.

Do not use Floyd’s algorithm as a substitute for correct list code. Use it to detect a corrupted list or to solve a cycle problem.

### Questions

#### Theoretical questions

1. What is a cycle in a linked list?
2. Why does a naive walk fail when a cycle exists?
3. How do the slow pointer and the fast pointer move?
4. What extra space does Floyd’s algorithm use?
5. How does a visited set detect a cycle?

#### Easy practical tasks

1. Draw a list of four nodes with a cycle from the last node to the second node.
2. Write five sentences that explain why two speeds meet in a cycle.
3. Make a table: "Method", "Time", "Extra space". Rows: Floyd, visited set.
4. Draw a list with no cycle. Write where the fast pointer stops.

#### Medium practical tasks

1. Implement Floyd cycle detection. Test one list with a cycle and one list without a cycle.
2. Implement visited-set detection with a set of node identities. Compare the code to Floyd.
3. Explain in six sentences why the fast pointer moves two steps.

#### Advanced practical tasks

1. After a meet, implement the second phase that finds the cycle start. Test with a known start.
2. Write a proof sketch in eight STE sentences: why the two pointers meet if a cycle exists.

---

## When lists lose to arrays (cache)

A linked list loses to an array for many real programs. The Big-O class can still look good.

**Locality.** Nodes sit at heap addresses that are far apart. A walk is pointer chasing. An array walk reads contiguous cells. The cache helps the array.

**Extra space.** Each node stores one or two pointers and often a header. An array of integers stores the integers. The list uses more bytes per element.

**Random access.** `get(i)` is Θ(i) on a list and Θ(1) on an array.

**Append.** A dynamic array append is amortized Θ(1) and cache-friendly. A list append is Θ(1) with a tail pointer, but the next scan is still slow.

Lists win when you insert or delete at a known node and you must not move the other elements. Lists also win when you splice two lists in Θ(1) by pointer updates. Those cases are real. They are less common than "I need a sequence of values".

```text
prefer a dynamic array when:
  you iterate often
  you index by i
  you append at the end
  n is large and elements are small

prefer a list when:
  you insert or delete at the head often
  you splice lists
  you already hold the node
```

Measure when n is large. Do not select a list only because insert in the middle is Θ(1) after you have the node. Find of the node can still be Θ(n).

### Questions

#### Theoretical questions

1. Why does a list walk have poor locality?
2. What extra space does a node add?
3. Why is `get(i)` slower on a list than on an array?
4. When does a list win?
5. Why is "insert is O(1)" incomplete for a list?

#### Easy practical tasks

1. Write five sentences that contrast list walk and array walk.
2. Make a table: "Need", "Prefer". Add four rows from this section.
3. Estimate pointers for 1 000 singly linked nodes.
4. Draw one array of four integers and one list of four nodes. Mark contiguous versus scattered.

#### Medium practical tasks

1. Time a full scan of n integers in a dynamic array and in a linked list. Use n = 100 000 or more if you can.
2. Time `get(n/2)` on both structures. Write the two times.
3. Explain in six sentences why splice of two lists can be Θ(1).

#### Advanced practical tasks

1. Implement both structures. Run scan, get(i), insert_head, and append. Write a result table.
2. Write a one-page policy: default to a dynamic array. List the exceptions with a reason for each exception.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe a complete singly linked list object. Name node fields and list fields.
2. How do sentinels and head/tail pointers work together in insert and delete?
3. Why does reverse need a careful update of tail?
4. How do you use Floyd’s algorithm after a suspected pointer bug?
5. When do cache effects cancel an O(1) list operation in a real program?

#### Easy practical tasks

1. Write a one-page cheat sheet: node, head, tail, singly, doubly, sentinel, reverse, cycle, Floyd, locality.
2. Draw insert at head, insert at tail, and reverse for a two-node list.
3. Write the cost of delete at tail for singly versus doubly, with tail stored.
4. List three bugs: null head, lost tail, and a cycle. Write one symptom for each bug.

#### Medium practical tasks

1. Implement a singly linked list with sentinel, tail, and length. Support insert_head, insert_tail, delete_head, reverse, and has_cycle.
2. Build the same sequence in a list and in a dynamic array. Compare scan times.
3. Write delete of the first node with value x. Handle "not found" and empty list.

#### Advanced practical tasks

1. Implement a doubly linked list and a splice operation that joins two lists in Θ(1). Test length and ends.
2. Read one standard linked-list type (Go `container/list` or Java `LinkedList`). Map each public method to this topic. Cite the page.
