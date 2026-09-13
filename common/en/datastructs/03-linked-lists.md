# 3. Linked Lists

## Description

A linked list stores elements in nodes. Each node holds a value and one or more pointers. This topic explains head, tail, singly and doubly linked lists, insert and delete at the ends, sentinel nodes, reverse, cycle detection, and the cases where arrays win.

Complete this topic after arrays. You will compare the two layouts.

Use one term for each concept. A **node** is a cell with a value and links. The **head** is the first node. The **tail** is the last node. A **singly linked list** has a next pointer. A **doubly linked list** has next and previous pointers. A **sentinel** is a dummy node that is not a real element. A **cycle** is a walk that returns to a node that you already visited.

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

The list object usually stores head, and often tail and length. Length makes the size operation Θ(1). If you do not store length, size walks the list.

Do not confuse the head pointer with the first value. The head pointer is a reference. The first value sits inside the first node.

A node is not an array cell. You cannot compute the address of the i-th node from the start. You walk i steps from the head. That walk is Θ(i).

Empty list, one-node list, and many-node list are different pictures. Write all three when you design operations.

### Questions

#### Theoretical questions

1. What fields does a singly linked node have?
2. What does a null head mean?
3. Why does a tail pointer make append Θ(1)?
4. Why does a stored length make size Θ(1)?
5. Why can you not compute the address of the i-th node?

#### Easy practical tasks

1. Draw a singly linked list of three integers. Label head and tail.
2. Draw an empty list. Show head and tail as null.
3. Make a table: "Pointer", "Names which node". Rows: head, tail, next of head.
4. Write five sentences that define node, head, and tail.

#### Medium practical tasks

1. Build a list of four nodes on paper. Walk from head and write the values.
2. Explain in six sentences when you store tail and when you omit tail.
3. Write the steps to get the value at index i. Count the nodes that you visit.

#### Advanced practical tasks

1. Implement a list object with head, tail, and length in a language that you know. Write a contract that keeps the three fields consistent.
2. Estimate bytes for n singly linked nodes versus n array cells. Use 8-byte values and 8-byte pointers.

---

## Singly and doubly linked lists

A **singly linked list** stores one next pointer per node. You walk only toward the tail. To delete a node, you need the previous node. A walk from the head finds that previous node.

A **doubly linked list** stores next and previous pointers. The head has a null previous pointer. The tail has a null next pointer. You can walk in two directions.

```text
singly:
head → [10|•] → [20|•] → [30|null]

doubly:
null ← [10] ⇄ [20] ⇄ [30] → null
        ↑               ↑
       head            tail
```

A doubly linked node uses more space. Each node holds two pointers. Insert and delete of a known node are easier. You can relink previous and next in Θ(1) if you already hold the node.

A singly linked list uses less space. Many algorithms only walk forward. A singly linked list is enough for those algorithms.

A **circular** list links the tail back to the head. This topic treats a circular list as a special form. A cycle in a list that should be linear is a defect. The reverse and cycle section covers detection.

Do not mix the two layouts in one picture. Draw next only, or draw next and previous.

### Questions

#### Theoretical questions

1. What extra field does a doubly linked node have?
2. Why do you need the previous node to delete in a singly linked list?
3. When is delete Θ(1) in a doubly linked list?
4. What space cost does the extra pointer add?
5. When is a singly linked list enough?

#### Easy practical tasks

1. Draw a doubly linked list of three integers. Label previous and next on the middle node.
2. Make a table: "List kind", "Pointers per node", "Walk directions". Add two rows.
3. Write five sentences that contrast singly and doubly linked lists.
4. From the tail of a doubly linked list, write the values if you walk previous links.

#### Medium practical tasks

1. Add previous fields to a three-node singly linked picture. Write each previous value (node or null).
2. Write the pointer assignments to insert a node between two existing doubly linked nodes.
3. Explain in six sentences why delete of a known node is harder in a singly linked list.

#### Advanced practical tasks

1. Implement both list kinds with the same public operations: append and walk. Write the extra fields that the doubly linked version stores.
2. Estimate the pointer bytes for n = 1 000 000 in each kind. Write when the extra memory is worth the easier delete.

---

## Insert and delete at ends

**Insert at the head** creates a new node and points it to the old head. Then head becomes the new node. The cost is Θ(1). You do not walk the list.

**Delete at the head** moves head to head.next. The old first node leaves the list. The cost is Θ(1). If the list becomes empty, set tail to null if you store tail.

**Insert at the tail** (append) needs the last node. With a tail pointer, you link the old tail to the new node and update tail. The cost is Θ(1). Without a tail pointer, you walk from the head. The cost is Θ(n).

**Delete at the tail** on a singly linked list needs the node before the tail. You walk from the head unless you store extra data. The cost is Θ(n) for a singly linked list. On a doubly linked list with a tail pointer, you use tail.previous. The cost is Θ(1).

```text
insert at head of 20 → 30
  new: 10 → 20 → 30
  head moves to 10

delete at head of 10 → 20 → 30
  head moves to 20
```

Keep length in sync if you store length. Empty and one-node lists are easy to get wrong. Write those cases first.

Insert and delete in the middle still need a walk to the position. The link change is Θ(1) after you hold the previous node. The walk is Θ(i).

### Questions

#### Theoretical questions

1. Why is insert at the head Θ(1)?
2. What extra pointer makes append Θ(1)?
3. Why is delete at the tail Θ(n) on a singly linked list?
4. How does a doubly linked tail make delete at the tail Θ(1)?
5. What must you set when delete at the head empties the list?

#### Easy practical tasks

1. Draw insert at the head of a two-node list. Show the new head.
2. Draw append with a tail pointer. Show the old tail and the new tail.
3. Make a table: "Operation", "Singly + tail", "Doubly + tail". Write Θ for insert head, insert tail, delete head, delete tail.
4. Write five sentences about end operations. Use only facts from this section.

#### Medium practical tasks

1. Write the pointer steps for append on an empty list and on a non-empty list. Include tail.
2. Write the steps for delete at the tail on a singly linked list of four nodes.
3. Explain in six sentences why middle insert is not Θ(1) from the start of the list.

#### Advanced practical tasks

1. Implement insert and delete at both ends for a singly linked list with tail. Test empty and one-node lists.
2. Implement delete at the tail for a doubly linked list. Compare the walk cost with the singly linked version.

---

## Sentinel nodes

A **sentinel** is a dummy node that is not a user element. The list always holds the sentinel. Head or a fixed pointer names the sentinel. Real elements start after the sentinel.

A sentinel removes many empty-list branches. Insert after the sentinel works for the first element and for later elements with the same code.

```text
without sentinel, empty:
  head = null

with head sentinel:
  sentinel → [10|•] → [20|null]
  the sentinel has no user value
```

A **header sentinel** sits before the first element. A **trailer sentinel** sits after the last element. A doubly linked list can use both. Then every real node has a previous node and a next node. The first real node has previous = header. The last real node has next = trailer.

Sentinels use extra space: one or two nodes. The extra space is a constant. The gain is simpler code and fewer null checks.

Do not count the sentinel as length. Length counts user elements only.

Do not print the sentinel value. The sentinel value is not data. Some implementations store no value in the sentinel.

A circular list can use one sentinel. The sentinel next points to the first element. The last element next points to the sentinel.

### Questions

#### Theoretical questions

1. What is a sentinel node?
2. Why does a sentinel reduce empty-list branches?
3. What is a header sentinel?
4. What is a trailer sentinel?
5. Why is the sentinel not part of length?

#### Easy practical tasks

1. Draw a list with a header sentinel and two user nodes.
2. Draw an empty list that still has a header sentinel.
3. Make a table: "Design", "Null checks for insert at head". Rows: no sentinel, header sentinel.
4. Write five sentences that define sentinels.

#### Medium practical tasks

1. Write the steps to insert the first user element after a header sentinel. Then write the steps to insert a second element at the head.
2. Explain in six sentences how header and trailer sentinels make every real node look the same.
3. Draw a circular list with one sentinel and three user nodes.

#### Advanced practical tasks

1. Implement a singly linked list with a header sentinel. Keep insert at head the same for empty and non-empty lists.
2. Implement a doubly linked list with header and trailer. Write why delete of the only user element does not need a special null tail.

---

## Reverse a list and cycle detection

**Reverse** a singly linked list changes every next pointer. The old head becomes the tail. The old tail becomes the head.

An iterative reverse uses three pointers: previous, current, and next. For each node, you set current.next to previous. Then you move the three pointers forward.

```text
before:  10 → 20 → 30 → null
after:   30 → 20 → 10 → null
```

Time is Θ(n). Extra space is Θ(1) for the iterative method. A recursive reverse uses the call stack. That extra space is Θ(n).

A **cycle** exists when a next pointer leads to an earlier node. A walk never reaches null. A cycle is usually a defect in a list that should be linear.

**Cycle detection** can use a slow pointer and a fast pointer. The slow pointer moves one node. The fast pointer moves two nodes. If they meet, a cycle exists. If the fast pointer reaches null, there is no cycle.

```text
cycle:
10 → 20 → 30
      ↑    |
      +----+
```

Time is Θ(n). Extra space is Θ(1). A set of visited node addresses also detects a cycle. That method uses Θ(n) extra space.

Do not reverse a list that can have a cycle until you detect the cycle. A reverse walk would not end.

After reverse, update tail if you store tail. The new tail is the old head.

### Questions

#### Theoretical questions

1. What does reverse do to head and tail?
2. Why is iterative reverse Θ(1) extra space?
3. What is a cycle in a list?
4. How do a slow pointer and a fast pointer detect a cycle?
5. Why can a set of visited addresses detect a cycle?

#### Easy practical tasks

1. Draw reverse of 1 → 2 → 3. Show the new head.
2. Draw a three-node cycle. Mark one back pointer.
3. Write five sentences about reverse and cycle detection.
4. Make a table: "Method", "Extra space". Rows: iterative reverse, recursive reverse, slow/fast detect, visited-set detect.

#### Medium practical tasks

1. Write the three-pointer steps of iterative reverse on a four-node list. Show previous, current, and next after each node.
2. Simulate slow and fast on a five-node list with a cycle at node 2. Write positions after each step until they meet or you stop.
3. Explain in six sentences why reverse must not run on an undetected cycle.

#### Advanced practical tasks

1. Implement iterative reverse and a slow/fast cycle test in a language that you know. Test a linear list and a list with a cycle.
2. Read a short note on finding the start of a cycle after a meet. Write the idea in eight STE sentences. Do not copy a full solution.

---

## When lists lose to arrays

A linked list wins when you insert or delete at the head often, and you already hold the position. A list also wins when you must not move other elements.

A list loses to an array in many common programs.

**Index access.** Array get(i) is Θ(1). List get(i) is Θ(i).

**Locality.** An array walk uses contiguous cells. A list walk chases pointers. The array walk is often much faster for the same Θ(n) class.

**Space.** Each list node stores one or two pointers. An array stores values in a tight block. For small elements, the pointers can dominate.

**Middle insert without a held node.** Both structures need a walk or a shift to reach index i. The array shift has better locality. The list walk has poor locality. For many machines, the array still wins for moderate n.

**Tail delete.** A dynamic array deletes the last element in Θ(1). A singly linked list needs a walk unless you use a more complex design.

```text
prefer array when
  you need get(i)
  you walk the full collection often
  elements are small
  you append and pop at the end

prefer list when
  you insert or delete at the head
  you hold a node and must relink without a shift
  you must keep stable node addresses
```

Do not select a list because it looks more "advanced". Select the layout from the operations.

Topic 1 defined locality. This section applies that fact to a real choice.

### Questions

#### Theoretical questions

1. When does a list win over an array?
2. Why does get(i) favor an array?
3. Why can pointer space dominate small elements?
4. Why can an array still win for middle insert at moderate n?
5. Why is "advanced" not a reason to select a list?

#### Easy practical tasks

1. Make a table: "Need", "Prefer array or list". Add five rows from this section.
2. Write five sentences about when lists lose. Use only facts from this section.
3. Draw an array of four small integers next to four list nodes with pointers. Mark extra pointer cells.
4. Write why pop of the last element favors a dynamic array.

#### Medium practical tasks

1. Write a decision note of eight sentences for a music playlist that jumps to track i and also inserts at the front.
2. Compare space for n = 100 integers of 4 bytes: array versus singly linked nodes with 8-byte pointers.
3. Explain in six sentences how locality can dominate a Θ(1) pointer change.

#### Advanced practical tasks

1. Time a full walk of n integers in an array and in a list for a large n. Write the two times and one conclusion.
2. Design a hybrid: a dynamic array of pointers to large records. Write when that hybrid is better than a list of records.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do head, tail, sentinels, and list kind change the cost of the four end operations?
2. When do you store length, tail, and previous pointers together, and when do you omit one of them?
3. How do reverse and cycle detection use walks without index formulas?
4. Why is a list node address stable when an array index is not after a middle insert?
5. Which facts from arrays (Topic 2) still apply when you compare the two layouts?

#### Easy practical tasks

1. Write a one-page cheat sheet: node, head, tail, singly, doubly, sentinel, reverse, cycle, when lists lose.
2. Draw one picture that includes a header sentinel, three user nodes, and a tail pointer.
3. List six operations on a waiting line. Mark each as cheap or expensive on a singly linked list without tail.
4. Make a table: "Defect", "Picture". Rows: empty list forgotten, tail not updated, cycle.

#### Medium practical tasks

1. Implement a doubly linked list with sentinels, insert at both ends, and iterative reverse. Write tests for empty and one-node lists.
2. Write a short report that maps each end operation to Θ for four designs: singly, singly+tail, doubly+tail, array.
3. Simulate a cycle detector and a reverse on the same three-node linear list. Write why reverse is safe.

#### Advanced practical tasks

1. Implement Floyd cycle detection and then reverse only if the list is linear. Document the order of the two operations.
2. Read one real list type in a standard library. Write which extra fields it stores and why lists still lose for random index access.
