# 1. Getting Started

## Description

A data structure is a method to store and organize values in a program. This topic explains that idea. You learn the difference between an abstract data type and an implementation. You also learn time and space as resources. Complete this topic before you study arrays and lists.

Use one term for each concept. An abstract data type (ADT) is a list of operations. An implementation is the code that does those operations. An operation is a named action such as insert, delete, find, or iterate.

---

## What a data structure is

A data structure is an organization of values in memory. The structure decides how you store values and how you get values. The structure also decides how you add values and how you remove values.

Examples of data structures: an array, a linked list, a hash table, and a tree. Each structure stores values. Each structure gives a different cost for each operation.

You do not select a structure only by name. You select a structure by the operations that you need. If you need fast access by index, you use an array. If you need fast insert at the front, you use a linked list.

A data structure is not an algorithm. An algorithm is a sequence of steps. A data structure is the layout of the values. An algorithm uses a data structure.

A value in a structure is an element. Some texts say "item". This handbook uses **element** for a stored value. This handbook uses **node** only for a linked cell that holds a value and a link.

```text
structure: array of 4 integers
indexes:   0    1    2    3
elements:  10   20   30   40
```

The same four numbers can live in a list of nodes. The values do not change. The layout in memory does change. That change of layout is the data structure.

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
2. Find one built-in collection in Go or in Python. Write the name and two operations that the collection gives.
3. Take a grocery list on paper. Write how an array stores that list. Write how a list of nodes stores that list.

#### Advanced practical tasks

1. Read the first chapter of one standard algorithms book (CLRS or Sedgewick). Write ten sentences that map the book terms to the terms in this section.
2. Select a small real program that you know. List every collection that the program uses. For each collection, write the data structure that it is.

---

## Abstract data type (ADT) vs implementation

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

### Questions

#### Theoretical questions

1. What is an abstract data type?
2. What is an implementation?
3. Why does a stack ADT not name "array" or "list"?
4. What does it mean to hide the layout?
5. How is an ADT different from a language type?

#### Easy practical tasks

1. Write the ADT contract for a queue: enqueue, dequeue, front, empty. Use four short lines.
2. Name two implementations that can do that queue contract.
3. Make a table with columns "ADT name", "Operations", and "Possible implementation". Add three rows.
4. In four sentences, explain why a user of a stack must not read the inner array.

#### Medium practical tasks

1. Write a stack ADT as a list of function signatures in Go or in Python. Do not write the function bodies.
2. Find the stack type in one standard library. Write which operations are public. Write which details are hidden.
3. Draw two boxes: "ADT" and "implementation". Put the stack operations in the first box. Put "array" and "list" in the second box.

#### Advanced practical tasks

1. Write two small implementations of the same stack ADT: one array, one list. Keep the public function names the same.
2. Change the array implementation so that the capacity grows. Do not change the function names. Write why the ADT stays the same.

---

## Time and space as resources

A program uses time and space. Time is the work that the processor does. Space is the memory that the program holds.

Each operation has a time cost. Each structure has a space cost. You compare structures by those costs.

Time cost grows when the input grows. Space cost also grows when the input grows. A structure that is fast can use more memory. A structure that uses little memory can be slow.

You measure time in steps, not in seconds. A step is one simple action: a compare, an assignment, or an index. Seconds change with the machine. Steps do not change.

You measure space in elements and in extra cells. Extra cells are pointers, indexes, and unused capacity.

```text
n = number of elements

array find by index:   few steps, does not grow with n
list find by value:    steps grow with n
array of n integers:   space grows with n
list of n integers:    space grows with n, plus one link per node
```

Do not optimize before you know the operation that is slow. First write the ADT. Then measure. Then select a faster implementation if you need it.

### Questions

#### Theoretical questions

1. What is time as a resource for a data structure?
2. What is space as a resource for a data structure?
3. Why do you count steps instead of seconds?
4. What is an extra cell?
5. Why can a fast structure use more memory?

#### Easy practical tasks

1. Write a two-column table: "Resource" and "How you measure it". Use the rows Time and Space.
2. Count the assignments in this code: `a = 1`, `b = 2`, `c = a + b`. Write the count.
3. An array holds 100 integers. A list holds 100 integers and 100 links. Which structure uses more extra cells?
4. List three operations on a contact list. For each operation, write "time grows with n" or "time does not grow with n". Give a reason.

#### Medium practical tasks

1. Write a loop that visits each element of a list of n names. Count the visits as a function of n.
2. Compare two structures for a list of 10 elements and for a list of 10 000 elements. Write which cost changes.
3. Find one function in a standard library that uses extra memory to go faster. Write the function name and the extra memory in one sentence.

#### Advanced practical tasks

1. Measure one small program with a timer for n = 100, 1 000, and 10 000. Write the three times. Write why the times are not the same as step counts.
2. Design a structure that is fast for find and slow for insert. Design a structure that is the opposite. Write the trade-off in eight sentences.

---

## Why the same ADT can have more than one implementation

One ADT can have many implementations. The operations stay the same. The costs change.

A list ADT can use a contiguous array. A list ADT can use linked nodes. The user still calls insert and find. The time for insert in the middle is not the same.

You select an implementation for the operations that you do often. If you do many finds by index, you use an array. If you do many inserts at the front, you use a linked list.

A second reason for more than one implementation is the language and the machine. One language gives a dynamic array. A different language gives a linked list in the standard library. The ADT is still a list.

A third reason is safety and simplicity. A simple implementation is easy to check. A complex implementation can be faster. You start with the simple implementation. You change it when a measurement shows a need.

```text
ADT: List
  get(i), set(i, x), insert(i, x), remove(i), length()

Implementation: fixed array     — get is fast; insert in the middle copies many elements
Implementation: linked nodes    — insert at the front is fast; get(i) walks i nodes
Implementation: dynamic array   — get is fast; append at the end is usually fast
```

The ADT name does not tell you the cost. The implementation tells you the cost. Always ask which implementation you have.

### Questions

#### Theoretical questions

1. Why can one ADT have more than one implementation?
2. What stays the same when you change the implementation?
3. What changes when you change the implementation?
4. How do you select an implementation?
5. Why does the ADT name not tell you the cost?

#### Easy practical tasks

1. Write three implementations of a list ADT. Use one sentence for each.
2. For a program that only reads by index, select one implementation. Give one reason.
3. For a program that only inserts at the front, select one implementation. Give one reason.
4. Make a table: "Implementation", "Fast operation", "Slow operation". Add three rows.

#### Medium practical tasks

1. Write the same three list operations as comments for an array and for a linked list. Do not write full code. Write the main steps.
2. Find two list types in one language (for example `list` and `array` in Python, or slice and container/list in Go). Write two cost differences.
3. A teammate says "use a list". Write four questions that you ask before you select the implementation.

#### Advanced practical tasks

1. Implement the same list ADT twice in one language: array and linked nodes. Time get(i) and insert(0) for n = 10 000.
2. Write a one-page note that explains when you keep the simple implementation. Use a numeric example from your timing.

---

## How to read "operations" (insert, delete, find, iterate)

An operation is a named action on a structure. Common operations are insert, delete, find, and iterate.

**Insert** adds an element. You must know the position. Insert at the front, insert at the end, or insert at index i are different operations. The costs are not the same.

**Delete** removes an element. Delete by position and delete by value are different operations. Delete by value must find the element first.

**Find** returns an element or reports that the element is not there. Find by index and find by value are different operations. Find by key is a third form. A map uses find by key.

**Iterate** visits each element one time. Iterate does not have to give a special order. Some structures iterate in index order. Some structures iterate in sorted order. Some structures iterate in an unspecified order.

When a text says "insert is fast", the text must name the position. When a text says "find is fast", the text must name the key. Read the full operation name.

```text
insert(x) at end     — append
insert(x) at index i — shift or relink
delete(i)            — remove the element at i
find(x)              — search for value x
find_key(k)          — search for key k
iterate              — visit each element
```

Write operations as a list when you design a structure. For each operation, write the input, the output, and the effect on the structure.

### Questions

#### Theoretical questions

1. What is an operation on a data structure?
2. Why are "insert at the front" and "insert at the end" different operations?
3. What extra work does delete by value do?
4. What is the difference between find by index and find by key?
5. Why must a text name the position when it says "insert is fast"?

#### Easy practical tasks

1. Write the four operations insert, delete, find, and iterate in one sentence each.
2. For an array of names, write one example of each of the four operations.
3. Make a table: "Operation", "Input", "Output". Add four rows.
4. Label this action: you walk a list of scores and print each score. Which operation is that?

#### Medium practical tasks

1. Write an ADT card for a phone book. List find by name, insert a contact, delete a contact, and iterate all names.
2. In Go or in Python, write a loop that iterates an array of five integers and finds the value 3. Mark which lines are iterate and which line is find.
3. Explain in six sentences why delete by index on an array copies elements.

#### Advanced practical tasks

1. Select one structure from a standard library. Write every public operation as a row: name, input, output, effect.
2. Design a small ADT for a waiting line at a shop. Name six operations. For each operation, write one sentence about time cost. Do not implement the ADT yet.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe the path from a problem to a selected data structure. Name ADT, operations, and implementation.
2. How do time cost and space cost work together when you select a structure?
3. A teammate says that a queue "is a list". Which facts do you use to correct that sentence?
4. What stays hidden when a user uses only the ADT operations?
5. Why do insert, delete, find, and iterate need a position, a value, or a key in the name?

#### Easy practical tasks

1. Write a one-page cheat sheet with these terms: data structure, element, node, ADT, implementation, operation, time, space.
2. Draw three boxes in a row: Problem, ADT, Implementation. Put one example in each box.
3. List five operations for a music playlist. Mark each operation as insert, delete, find, or iterate.
4. Create a two-column table: "Term in this topic" and "One-sentence meaning". Add six rows.

#### Medium practical tasks

1. Write two implementations of a tiny bag ADT that only supports insert and iterate. Use an array and a linked list. Keep the public names the same.
2. Interview a program that you wrote before. Replace each collection name with "ADT + implementation". Write the new list.
3. Write a short script that prints n integers from an array. Then write the same print with a chain of nodes. Count the steps for n = 5.

#### Advanced practical tasks

1. Read the documentation of one collection type in Go or in Python. Write the ADT in your words. Write the implementation notes that the documentation gives.
2. Design a study plan for the next ten topics. For each topic, write one ADT and one implementation that you will build.
