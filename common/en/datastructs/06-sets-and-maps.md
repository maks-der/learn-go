# 6. Sets and Maps

## Description

A set stores unique keys. A map stores a value for each key. This topic explains unordered hash implementations, ordered tree implementations, multisets, iterator invalidation, and how you choose hash versus tree.

Complete this topic after hash tables. Complete binary search trees in Topic 8 after you know why an ordered map exists.

Use one term for each concept. A **set** answers "is this key present?". A **map** (also called a dictionary) answers "what value sits at this key?". **Unordered** means the walk order is not a sort of the keys. **Ordered** means the walk visits keys in sorted order. A **multiset** stores the same key more than one time. An **iterator** is a cursor that walks elements. **Invalidation** means that cursor is no longer safe to use.

---

## Unordered set/map (hash) vs ordered set/map (tree)

An **unordered set** is a hash table of keys. An **unordered map** is a hash table of key–value pairs. Find, insert, and delete are average Θ(1) with a good hash and a bounded load factor. Walk order is unspecified. Two walks can give two orders after a resize.

An **ordered set** is a balanced search tree of keys. An **ordered map** stores a value in each node. Find, insert, and delete are Θ(log n) in the worst case if the tree stays balanced. Walk in-order visits keys from small to large.

```text
unordered map (hash)
  "cat" → 1
  "eel" → 3
  "dog" → 2
  walk can be: eel, cat, dog

ordered map (tree)
  same three pairs
  walk is: cat, dog, eel
```

Both ADTs support get, put, delete, and "contains" for a key. The ordered ADT also supports:

- first and last key
- next key after k
- range walk from L to R

The hash ADT does not support those operations in a cheap way. You would scan all keys.

Space: a hash table stores a bucket array plus entries. A tree stores pointers in each node. Constants differ. Classes are the usual reason to choose.

Do not use a hash map when you need sorted keys. Do not use a tree map when you only need get and put and you want average constant time.

Topic 5 stated worst-case Θ(n) for a degraded hash table. A balanced tree gives a strict logarithmic bound. That bound is a reason to pick a tree even if you do not need order.

### Questions

#### Theoretical questions

1. What question does a set answer?
2. What question does a map answer?
3. What walk order does an unordered map give?
4. What extra operations does an ordered map give?
5. Why can a tree win even when you do not print keys in order?

#### Easy practical tasks

1. Make a table: "ADT", "Implementation", "Find class". Add four rows (set/map × hash/tree).
2. Write the in-order walk of keys 7, 2, 5.
3. List three operations that a hash map does not do cheaply.
4. Write five sentences that contrast unordered and ordered maps.

#### Medium practical tasks

1. Write an ADT card for an ordered map. Include get, put, delete, first, next, and range.
2. Explain in six sentences when resize of a hash map changes walk order.
3. For a phone book that must print names in alphabet order, select hash or tree. Give two reasons.

#### Advanced practical tasks

1. Implement a tiny unordered map with chaining and a tiny ordered map with an unbalanced BST. Compare walk orders.
2. Write a one-page note: worst-case find for hash versus balanced tree. Use Topic 5 facts.

---

## Multiset

A **multiset** (bag) stores keys and allows duplicates. The key 7 can appear three times. A set would store 7 one time.

A **count multiset** stores each distinct key one time and stores an integer count. Insert of an existing key adds 1 to the count. Delete subtracts 1. When the count reaches 0, you remove the key.

```text
multiset of integers:  2, 2, 5, 5, 5, 9
count form:            2 → 2,  5 → 3,  9 → 1
size (elements):       6
size (distinct):       3
```

Write which size you mean. **Total size** counts elements. **Distinct size** counts keys.

A hash map from key to count implements an unordered multiset. A tree map from key to count implements an ordered multiset. The walk of an ordered multiset can emit each key `count` times, or emit each distinct key one time.

A list of values is also a multiset if you do not remove duplicates. Find of one value is then Θ(n). The map-of-counts form finds a key in the class of the inner map.

Do not use a set when the problem counts copies (for example, votes or inventory). Do not use a multiset when uniqueness is the rule.

A map of key to list of records is a different design. That design stores extra data per copy. A count is enough when the copies are identical.

### Questions

#### Theoretical questions

1. How does a multiset differ from a set?
2. What does a count form store?
3. What is the difference between total size and distinct size?
4. How do you implement an unordered multiset with a map?
5. When is a list a poor multiset?

#### Easy practical tasks

1. Write the count form of 1, 1, 1, 4, 4.
2. Make a table: "Structure", "Allows duplicates", "Find class". Rows: set, count multiset, list.
3. Draw a small map 3 → 2, 8 → 1. Write the elements if you expand counts.
4. Write five sentences that define a multiset.

#### Medium practical tasks

1. Write the steps of insert and delete on a count map. Include the count-0 remove.
2. Explain in six sentences when you store a list per key instead of a count.
3. For the bag 2, 2, 5, write total size and distinct size after one delete of 2.

#### Advanced practical tasks

1. Implement a count multiset with a hash map. Support insert, delete-one, and count(key).
2. Implement an ordered multiset that walks each key `count` times in sorted order.

---

## Iterator invalidation (conceptual)

An **iterator** is a cursor. It names a position in a collection. Next moves the cursor. Some iterators also allow delete of the current element.

**Invalidation** means the cursor is no longer safe. A later next or read is a defect. The collection changed in a way that the cursor does not understand.

Typical causes:

- resize of a hash table (rehash moves entries)
- insert that grows a dynamic array (the buffer address can change)
- delete of the node that the cursor holds
- delete of an earlier array cell (later indexes move)

```text
walk a hash map
  mid-walk: insert many keys, table resizes
  old cursor points at a stale bucket
  next is invalid
```

Rules differ by implementation. This topic teaches the idea, not one language. A safe pattern is:

1. Do not insert or delete during a walk, except a defined "delete current" operation.
2. If you must change the collection, collect keys first, then change after the walk.
3. After resize or rebuild, create a new iterator.

An index integer into an array can also go stale. After insert at the front, old index i names a different element.

A pointer to a list node can stay valid when you insert elsewhere, because node addresses do not move. Delete of that node still invalidates that pointer.

Do not assume that a walk is safe under concurrent writes. Topic 15 covers concurrent maps at a high level.

### Questions

#### Theoretical questions

1. What is an iterator?
2. What does invalidation mean?
3. Why can hash-table resize invalidate a cursor?
4. Why can an array index go stale after insert at the front?
5. Why can a list-node pointer stay valid across an insert elsewhere?

#### Easy practical tasks

1. Write five sentences that define iterator invalidation.
2. Make a table: "Change", "Can invalidate a hash cursor?". Rows: find only, insert many, delete current if defined.
3. Draw an array [10, 20, 30] with a cursor at index 2. Show indexes after insert at 0.
4. List the three safe patterns from this section.

#### Medium practical tasks

1. Write a walk that would break if you insert during the walk. Then write the collect-then-change version.
2. Explain in six sentences why "delete current" can be safe when "insert anywhere" is not.
3. Compare a cursor that holds an index with a cursor that holds a node pointer. Write two invalidation examples.

#### Advanced practical tasks

1. Read documentation for one collection iterator in a language that you know. Write which operations invalidate it. Cite the source.
2. Implement a map walk that copies keys to an array, then deletes some keys after the copy. Write why the walk itself stays valid.

---

## How to choose hash vs tree

Choose a **hash set or hash map** when:

- you need get, put, delete, and contains
- you do not need sorted order or nearest key
- you accept average Θ(1) and a possible Θ(n) worst case
- keys are hashable and stable

Choose a **tree set or tree map** when:

- you need sorted walk, first, last, next, or a range
- you need a Θ(log n) bound on every call
- keys have a total order (you can compare them)
- you can pay extra pointers and a larger constant

```text
need "all names from K to M"     → ordered tree
need "value for this id"          → hash (usual)
need "never pause for Θ(n) chain" → balanced tree (or a hash with extra care)
need "most frequent key"          → neither alone; add a count or a heap
```

If keys cannot hash, you cannot use a hash table. If keys cannot compare, you cannot use a search tree.

A sorted array plus binary search is a third choice. Find is Θ(log n). Insert in the middle is Θ(n). Use it when the collection is almost static.

A list of pairs is a fourth choice. Find is Θ(n). Use it only for very small n.

Measure if n is moderate and the constants matter. Notation selects the family. Measurement selects the winner on a machine.

Do not choose a tree because it looks more formal. Do not choose a hash because it looks faster in a slogan.

Topic 9 explains balanced trees. Topic 8 explains the BST property that the ordered map uses.

### Questions

#### Theoretical questions

1. When do you choose a hash map?
2. When do you choose a tree map?
3. What must keys support for a hash table?
4. What must keys support for a search tree?
5. When is a sorted array a better map than a tree?

#### Easy practical tasks

1. Make a table: "Need", "Choose hash, tree, or sorted array". Add five rows.
2. Write five sentences that state the choice rules.
3. List two key types that compare but that you might not want to mutate (stability still matters for hash).
4. Write why a list of pairs is only for small n.

#### Medium practical tasks

1. For an autocomplete prefix walk, write why a hash map of full words is not enough. Name a later structure (trie) only as a label.
2. Explain in six sentences the difference between "keys can hash" and "keys can compare".
3. Write a decision for a leaderboard that must show ranks in order and also update one score.

#### Advanced practical tasks

1. Implement the same map ADT twice: hash chaining and a BST. Time n random finds. Write classes and times.
2. Write a one-page choice guide for a teammate who only knows the words "set" and "map".

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do uniqueness, order, and duplicate counts decide set versus map versus multiset?
2. Why is walk order part of the ADT contract for an ordered map and not for a hash map?
3. How does iterator invalidation change the way you delete while you walk?
4. Which facts from Topic 5 still apply when you wrap a hash table as a set?
5. When do compare, hash, and stability each become the blocking key rule?

#### Easy practical tasks

1. Write a one-page cheat sheet: set, map, unordered, ordered, multiset, count, iterator, invalidation, hash vs tree.
2. Draw one hash map and one tree map of the same three pairs. Show two different walks.
3. Make a table: "Collection", "Duplicate keys", "Sorted walk". Add set, map, multiset.
4. List four defects: mutated hash key, insert during hash walk, use of a set for inventory counts, hash map used for range.

#### Medium practical tasks

1. Implement a set as a map to a dummy value. Write why the value is unused.
2. Write a test plan for a multiset: insert twice, delete once, count, distinct size, total size.
3. Design an API note: "do not modify the map during iteration" plus one allowed exception.

#### Advanced practical tasks

1. Build a small ordered map that supports range print [L, R]. Then build a hash map that does the same by a full scan. Compare times for large n.
2. Read one standard library set and one map (hash or tree). Write which operations the documentation lists and which order it promises.
