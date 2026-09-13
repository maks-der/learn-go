# 9. Sets and Maps as ADTs

## Description

A set and a map are abstract data types. This topic explains unordered (hash) forms and ordered (tree) forms. You also learn a multiset, iterator invalidation, and how you choose hash versus tree. Complete this topic after hash tables. Complete the tree topics after this preview.

Use one term for each concept. A **set** stores unique elements. A **map** stores unique keys. Each key has one value. **Unordered** means the iteration order is not a sort order. **Ordered** means iteration follows a key order. A **multiset** (bag) allows duplicate elements.

---

## Unordered set / map (hash)

An **unordered set** supports insert, delete, and contains. The set does not keep a sort order. Two equal elements do not appear twice.

An **unordered map** supports insert(key, value), delete(key), and find(key). Keys are unique. A second insert of the same key updates the value.

The usual implementation is a hash table. Average contains and find are Θ(1). Worst case is Θ(n).

```text
set:   {"cat", "dog"}
map:   {"cat": 3, "dog": 1}
```

```go
seen := map[string]struct{}{"cat": {}, "dog": {}}
count := map[string]int{"cat": 3, "dog": 1}
```

```python
seen = {"cat", "dog"}
count = {"cat": 3, "dog": 1}
```

Iteration visits each element or each key one time. The order can change after a resize. Do not write a program that needs sorted iteration on a hash set.

A set is not a list. Contains on a list is Θ(n). Contains on a hash set is average Θ(1). If you only ask "is x present?", prefer a set.

A map is not a pair of parallel arrays unless n is tiny. Find by key on a hash map is average Θ(1). Find by key on an unsorted array is Θ(n).

Use `struct{}` values in Go when you need a set. The value uses no extra information. Use a `set` object in Python.

### Questions

#### Theoretical questions

1. What operations does an unordered set provide?
2. What happens when you insert an existing key in a map?
3. Why is iteration order not a sort order?
4. Why is a set better than a list for contains?
5. How do you make a set in Go with a map?

#### Easy practical tasks

1. Build a set of five words in your language. Test contains for one present word and one absent word.
2. Build a map from name to integer. Update one key. Print the value.
3. Make a table: "ADT", "Unique by", "Typical impl". Rows: unordered set, unordered map.
4. Write five sentences that contrast a set and a list.

#### Medium practical tasks

1. Count word frequencies in a short text with a map.
2. Deduplicate a list of integers with a set. Then compare length with the original list.
3. Time n contains checks on a list and on a set for n = 20 000.

#### Advanced practical tasks

1. Implement a set ADT on your chaining hash table from topic 8. Support insert, delete, contains, and iterate.
2. Explain in one page why a hash map is the default map in Go and in Python.

---

## Ordered set / map (tree)

An **ordered set** keeps elements in sorted order. Iteration visits the smallest element first (or a defined order).

An **ordered map** keeps keys in sorted order. Find still uses the key. You can also ask for the next key or for a range of keys.

The usual implementation is a balanced binary search tree, or a skip list in some libraries. Find, insert, and delete are Θ(log n) in the worst case if the tree is balanced. You will implement a basic BST in topic 11. Balance is a later topic.

```text
ordered set:   1, 3, 8, 20
range 4..10 →  8
```

Go has no ordered map in the language. You sort a slice of keys when you need order. The standard library has `slices.Sort`. Some extra libraries provide a tree map.

Python 3.7+ `dict` keeps **insertion** order. Insertion order is not key-sorted order. `dict` is not an ordered map in the tree sense. Use `sorted(d)` for a sorted view. Use a third-party tree if you need Θ(log n) next-key.

Java `TreeMap` and C++ `std::map` are ordered maps. Learn the idea even if your language has no built-in tree map.

Use an ordered map when you need:

- sorted iteration
- nearest key
- range queries

Use a hash map when you only need find by exact key.

### Questions

#### Theoretical questions

1. What extra property does an ordered set have?
2. What is a range query in one sentence?
3. What typical time class does a balanced tree map give?
4. Why is Python `dict` order not the same as key order?
5. When do you use an ordered map instead of a hash map?

#### Easy practical tasks

1. Put five integers in a set. Print them sorted. Write whether the set type did the sort or you did the sort.
2. Make a table: "Need", "Hash map", "Tree map". Add three needs.
3. Write five sentences that define an ordered map.
4. Name one language type that is a tree map, even if you do not use that language daily.

#### Medium practical tasks

1. Build a map. Print keys in sorted order with an extra sort of the key list.
2. Write a function that returns all keys in [lo, hi] from a sorted key slice. Do not use a tree yet.
3. Explain in six sentences why Go omits a language-level tree map.

#### Advanced practical tasks

1. After you finish topic 11, implement a tiny ordered map with a BST. Support find and in-order iterate. State that balance is missing.
2. Read Java `TreeMap` or C++ `std::map` documentation. Write ten STE sentences about operations and cost. Cite the page.

---

## Multiset / bag

A **multiset** (bag) stores elements and allows duplicates. The bag counts how many times an element appears.

Operations: add(x), remove(x) (remove one copy), count(x), and size (total copies).

```text
add 2, add 2, add 5
count(2) = 2
remove(2)
count(2) = 1
```

A common implementation is a map from element to count. Add increments the count. Remove decrements the count. If the count becomes 0, you delete the key.

```go
bag := map[int]int{}
bag[2]++
bag[2]++
bag[5]++
```

```python
from collections import Counter
bag = Counter([2, 2, 5])
```

A list also stores duplicates. Count on a list is Θ(n). Count on a map bag is average Θ(1).

A multiset is not a set. A set loses the second copy. If the problem says "two tokens of type A", you need a bag or a list.

Remove of one copy is not delete of all copies. Define the operation in the ADT. Provide `remove_all` if you need it.

### Questions

#### Theoretical questions

1. How is a multiset different from a set?
2. What does count(x) return?
3. How does a map implement a bag?
4. What do you do when a count reaches 0?
5. Why is a list a weak default for count(x)?

#### Easy practical tasks

1. Draw a bag after add 1, add 1, add 2, remove 1. Show counts.
2. Write the ADT contract for a bag in four lines.
3. Make a table: "ADT", "Duplicates?". Rows: set, bag, list.
4. Build a `Counter` or a map of counts for the letters in `banana`.

#### Medium practical tasks

1. Implement a bag with a map of counts. Include remove that fails when count is 0.
2. Compare `size` as number of keys versus number of copies. Write both functions.
3. Time count of one value in a list of n duplicates versus in a bag. Use a large n.

#### Advanced practical tasks

1. Implement a bag that also supports iterate of unique elements and iterate of all copies.
2. Write a short note: when you use a bag, when you use a list, and when you use a set. Give one program example for each.

---

## Iterator invalidation (language-specific, conceptual)

An **iterator** is an object or a loop state that visits elements. **Invalidation** means the iterator is no longer safe. A later use can skip elements, visit twice, crash, or give a stale view.

The rules depend on the language and the collection.

**Hash map resize.** A grow moves keys. An iterator that holds a bucket index can be wrong. Some languages forbid insert during iteration.

**Go maps.** You must not add a new key during `range` over the same map. An update of an existing value is allowed. Delete of the current key has defined rules. Read the current specification when you write production code. A simple student rule: do not insert new keys in a `range` loop.

**Python dict.** You must not change the set of keys during iteration. An insert or delete of a key can raise `RuntimeError`.

**Lists.** An iterator over a dynamic array can break when you insert or delete in the middle. Indexes shift.

```text
unsafe idea
  for each key in map:
      map.insert(new_key)   // can invalidate
```

Safe patterns:

- collect keys to add in a second list, then insert after the loop
- iterate a copy of the keys
- use a for loop with an index only when the language defines the rule

This handbook does not list every language rule. You must read the document for your collection. Learn the idea: a structural change can invalidate a walk.

### Questions

#### Theoretical questions

1. What is an iterator in this section?
2. What does invalidation mean?
3. Why can a hash-map resize break a walk?
4. What is a safe pattern for "add keys while you walk"?
5. Why are the rules language-specific?

#### Easy practical tasks

1. Write five sentences that explain invalidation without code.
2. Make a table: "Change", "Risk". Rows: insert new key in map loop, update value, delete current.
3. In Python, change a dict’s keys during `for k in d`. Record the error. In Go, skip if you do not want a panic; write the documented rule instead.
4. List two safe patterns from this section.

#### Medium practical tasks

1. Write a loop that collects keys to delete, then deletes after the loop. Use a map.
2. Iterate a list by index and delete the current index. Write why this is easy to get wrong. Prefer a second list.
3. Read the Go spec or Python docs on dict/map iteration. Write six STE sentences. Cite the page.

#### Advanced practical tasks

1. Cause one invalidation bug on purpose. Record the symptom. Then apply a safe pattern. Write the pair.
2. Compare iterator rules for a Go slice, a Go map, a Python list, and a Python dict. Make a four-row table.

---

## Choosing hash vs tree

Choose the structure from the operations.

**Choose a hash set or hash map** when:

- you need exact find, insert, and delete
- you do not need sorted keys
- you accept average Θ(1) and worst Θ(n)
- keys are hashable and stable

**Choose a tree set or tree map** when:

- you need sorted iteration
- you need next, previous, or range
- you need a worst-case Θ(log n) bound
- keys are comparable (an order exists)

**Choose a sorted array** when:

- you build once and find many times
- n is small or updates are rare
- binary search is enough
- you can pay Θ(n) to insert in the middle

```text
need only exact key?          → hash map
need order or range?          → tree map (or sort when you print)
need worst-case log n?        → balanced tree
n tiny and simple code?       → array + scan
```

Do not choose a tree because it looks advanced. A hash map is the right default for a dictionary.

Do not choose a hash map because it is popular if you need the next larger key. Sort or use a tree.

Memory: a hash table uses extra buckets. A tree uses extra pointers per node. Measure if memory is tight.

### Questions

#### Theoretical questions

1. When is a hash map the default choice?
2. When do you need a tree map?
3. When is a sorted array enough?
4. Why is "tree looks advanced" a bad reason?
5. What extra property must tree keys have that hash keys might not use?

#### Easy practical tasks

1. Make a table: "Problem", "Structure". Add four problems from this section.
2. Write five sentences that state the default and the exceptions.
3. For a phone book by exact name, pick a structure. Give one reason.
4. For "all events between 10:00 and 11:00", pick a structure. Give one reason.

#### Medium practical tasks

1. Implement the same word-count with a hash map and with a sort of pairs at the end. Write which part needs order.
2. List three APIs in a standard library: hash map, sort, and (if present) a tree. Write one line each.
3. Explain in six sentences a case where worst-case Θ(log n) matters more than average Θ(1).

#### Advanced practical tasks

1. Write a decision flowchart in text: six questions, each yes/no, that ends in hash, tree, or sorted array.
2. After topic 11, time find on a hash map and on your BST for n = 50 000 random keys. Write the times and a caution about balance.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe set, map, and bag as three ADTs. Name uniqueness and values.
2. How do unordered and ordered maps differ in iteration and in extra operations?
3. Why is iterator invalidation part of ADT use, not only part of implementation?
4. A teammate uses a list as a set of 100 000 ids. Which facts do you use to change the structure?
5. How do hashable keys and comparable keys appear in the hash-versus-tree choice?

#### Easy practical tasks

1. Write a one-page cheat sheet: set, map, bag, unordered, ordered, hash, tree, iterator, invalidation.
2. Draw a map of three keys and a bag of three copies of one key plus one other key.
3. Write one safe loop that updates values in a map and does not add keys.
4. Pick a structure for: unique visitors, word count, sorted leaderboard. Write one line each.

#### Medium practical tasks

1. Implement set, map, and bag facades on one hash table type. Keep the public names different.
2. Build a tiny API: `Add`, `Has`, `Get`, `RangePrint` where RangePrint sorts keys. Document that sort is extra.
3. Write tests that prove a set drops duplicates and a bag keeps counts.

#### Advanced practical tasks

1. Design an ordered-map ADT on paper: find, min, max, successor, range. Mark which operations a hash map cannot do well.
2. Read two language documents (Go map + Python dict, or one of them plus a tree map). Write a STE comparison of order, invalidation, and cost. Cite both sources.
