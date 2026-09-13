# 8. Hash Tables

## Description

A hash table maps a key to a value in average constant time. This topic explains hash functions, collisions, chaining, open addressing, load factor, key rules, and worst-case cost. Complete this topic after arrays and lists. You will use a hash table as a set and as a map in the next topic.

Use one term for each concept. A **key** is the lookup name. A **hash function** maps a key to an integer. A **bucket** is a slot in the table array. A **collision** occurs when two keys map to the same bucket. The **load factor** is n divided by the number of buckets.

---

## Hash function

A **hash function** h maps a key to an integer. The table uses `index = h(key) mod m` when the table has m buckets. Some tables use a bit mask when m is a power of 2.

The hash function must be **deterministic**. The same key must give the same hash while the key lives in the table.

The hash function must use the full key (or a defined part of the key). If you hash only the first letter of a string, many keys collide.

```text
h("cat")  →  81321
index     →  81321 mod m
```

```go
// Go: you do not call the map hash; the runtime hashes the key
m := map[string]int{"cat": 1}
```

```python
# Python: hash(obj) for a hashable object
h = hash("cat")
```

A language hash is often not a cryptographic hash. Do not use a table hash as a password hash. A table hash is for distribution in memory.

A good hash is fast. A slow hash can dominate the lookup cost.

Integer keys can use a mix function so that nearby integers do not all land in nearby buckets. String keys mix the characters. The exact mix is a later engineering topic. This topic requires the idea: map the key to an integer, then to a bucket.

### Questions

#### Theoretical questions

1. What does a hash function map?
2. How do you get a bucket index from a hash and m?
3. What does deterministic mean for a hash function?
4. Why must you not hash only the first letter of a string?
5. Why is a table hash not a password hash?

#### Easy practical tasks

1. Write a tiny hash for an integer key: `h(x) = x`. Then write `h(x) mod 8` for x = 1, 8, 16.
2. Make a table: "Key", "hash", "hash mod 5" for three strings. Use Python `hash` or a simple sum of letter codes.
3. Write five sentences that define a hash function for a table.
4. List two key types that your language can hash.

#### Medium practical tasks

1. Write `index(key, m)` with a simple string hash: sum of character codes. Show collisions for m = 4.
2. Compare `hash("a")` and `hash("b")` in Python, or print two Go map keys that you insert. Write that you do not control the mix.
3. Explain in six sentences why a fast hash matters for a large table.

#### Advanced practical tasks

1. Implement djb2 or a similar well-known string mix (look up the formula). Compare collision counts with sum-of-codes on 1 000 words.
2. Read how Go or Python hashes strings (high level). Write ten STE sentences. Cite the source.

---

## Uniformity and collisions

**Uniformity** means that keys spread across the buckets. Each bucket gets about n / m keys if the hash is uniform.

A **collision** occurs when two different keys give the same bucket index. Collisions are normal. You must resolve collisions. You cannot assume that h is perfect for all keys.

```text
m = 4
h(alice) mod 4 = 1
h(bob)   mod 4 = 1   ← collision
```

If the hash is not uniform, some buckets get long chains. Lookup becomes slow.

An **adversary** can pick keys that collide if the hash is predictable and the table is open to untrusted keys. Some languages randomize a seed per process. That seed makes collision attacks harder.

Two different keys can have the same hash integer before the mod. That event is also a collision. The table still compares keys for equality. The hash only selects the bucket.

Equality and hash must agree. If two keys are equal, they must have the same hash. If two keys are not equal, they should usually have different hashes. They can still collide.

### Questions

#### Theoretical questions

1. What is uniformity?
2. What is a collision?
3. Why are collisions normal?
4. What happens when a hash is not uniform?
5. Why must equal keys have the same hash?

#### Easy practical tasks

1. Draw four buckets. Place six keys with two collisions in one picture.
2. Write five sentences about why a perfect hash for all strings is not the usual plan.
3. Make a table: "Event", "What you do". Rows: uniform spread, collision, equal keys.
4. Give one example of two integers that land in the same bucket for m = 10 if h(x) = x.

#### Medium practical tasks

1. Hash 100 sequential integers with `h(x) = x` and m = 16. Count how many keys each bucket gets.
2. Hash the same integers with `h(x) = x * 31`. Compare the counts.
3. Explain in six sentences why the table still compares keys after it uses the hash.

#### Advanced practical tasks

1. Build a histogram of bucket sizes for a simple hash on a word list. Write the max chain length.
2. Read one note on hash-flooding attacks. Write eight STE sentences about seeds. Cite the source.

---

## Separate chaining

**Separate chaining** stores a list in each bucket. When two keys collide, both keys sit in the same list.

Lookup hashes the key, selects the bucket, then walks the list. The walk compares keys.

Insert hashes the key, finds the bucket, then adds a node if the key is new. If the key exists, you update the value.

Delete finds the key in the bucket list and removes the node.

```text
buckets
  0:  → ("ab", 1)
  1:  → ("cd", 2) → ("xy", 3)
  2:  empty
  3:  → ("zz", 4)
```

If the hash is uniform and m is about n, each list is short. Average lookup is Θ(1). If all keys land in one bucket, lookup is Θ(n).

A bucket can use a dynamic array instead of a linked list. An array bucket has better locality. Some languages switch a long bucket to a tree. That switch is an advanced design.

Chaining grows well. You can have n > m. The lists get longer. You still resize to keep the load factor small.

### Questions

#### Theoretical questions

1. What does a bucket store in separate chaining?
2. What steps does lookup take?
3. Why is average lookup Θ(1) when lists are short?
4. What is the worst-case lookup?
5. Why can a bucket use an array instead of a list?

#### Easy practical tasks

1. Draw three buckets. Put two keys in bucket 1 and one key in bucket 0.
2. Write the insert steps for a new key.
3. Make a table: "Operation", "Main steps". Rows: find, insert, delete.
4. Write what happens when you insert a key that already exists.

#### Medium practical tasks

1. Implement a chaining table for string keys and integer values. Use an array of lists.
2. Insert 20 keys with m = 5. Print each bucket length.
3. Time find on a uniform set and on a set that you force into one bucket (bad hash). Write the two times.

#### Advanced practical tasks

1. Implement resize of a chaining table: allocate a larger bucket array, re-insert all keys.
2. Compare list buckets and slice buckets for cache cost in a short written report. Measure if you can.

---

## Open addressing: linear / quadratic probing, Robin Hood (awareness)

**Open addressing** stores each key in the table array. There is no node list. When a collision occurs, the table probes another slot.

**Linear probing** tries index, index+1, index+2, ... and wraps around. Linear probing is simple. Keys form clusters. Clusters grow and slow the probe.

**Quadratic probing** tries offsets 1², 2², 3², ... (plus the base index). Quadratic probing reduces some clustering. You must pick m and the formula so that you visit empty slots.

**Robin Hood hashing** is an awareness topic. In Robin Hood, a key that travels far can displace a key that travels a short distance. The goal is more uniform probe lengths. You do not need to implement Robin Hood in this topic.

```text
linear probe, m = 8, start = 3
  try 3, 4, 5, 6, ... until empty or key found
```

Delete is harder than in chaining. A removed slot cannot be a simple empty if a probe chain must continue. Tables use a **tombstone** (deleted mark) or they rehash the cluster.

Open addressing needs more empty slots than chaining. A high load factor makes long probes. Resize earlier than you might in chaining.

You must not store more than m keys in a simple open table. Capacity is m. Load factor must stay below 1.

### Questions

#### Theoretical questions

1. Where does an open-addressing table store keys?
2. How does linear probing choose the next slot?
3. What is clustering in linear probing?
4. Why is delete harder than in chaining?
5. What is a tombstone?

#### Easy practical tasks

1. Draw eight slots. Insert three keys that collide at slot 2 with linear probing.
2. Write five sentences that contrast chaining and open addressing.
3. Make a table: "Method", "Next slot idea". Rows: linear, quadratic.
4. Write one sentence about Robin Hood hashing as awareness only.

#### Medium practical tasks

1. Implement insert and find with linear probing for integer keys. Stop when the table is full.
2. Show a delete with a tombstone. Then find a key that sits after the tombstone.
3. Count probe length for find after you insert n = m/2 keys.

#### Advanced practical tasks

1. Implement quadratic probing. Document the index formula and the full-table rule.
2. Read a short description of Robin Hood or Hopscotch hashing. Write eight STE sentences. Do not implement unless you want extra work.

---

## Load factor and resize

The **load factor** α is n / m. n is the number of keys. m is the number of buckets (or slots).

When α grows, collisions grow. Average lookup becomes slower. The table **resizes**: it allocates a larger m and **rehash** all keys.

```text
if n / m > max_alpha:
    new_m = next_size(m)   // often about 2 * m
    new_table = array of new_m buckets
    for each key in the old table:
        insert key into new_table
```

A typical max_alpha for chaining is about 0.75 to 1.0. A typical max_alpha for open addressing is lower, such as 0.5 to 0.7.

Resize is Θ(n). Resize is rare if m grows geometrically. Amortized insert can stay Θ(1) on average.

A language map resizes for you. You still pick an initial size if you know n. Go `make(map[K]V, n)` is a hint. Python `dict` grows on its own.

After resize, indexes change. You must not store a bucket index as a long-lived pointer into the table. Store the key.

### Questions

#### Theoretical questions

1. What is load factor?
2. Why does a large load factor slow lookup?
3. What does rehash do?
4. Why can amortized insert stay average Θ(1) if resize is Θ(n)?
5. Why must you not store a bucket index as a long-lived reference?

#### Easy practical tasks

1. Compute α for n = 30 and m = 64.
2. Draw a grow from m = 4 to m = 8 with three keys. Show new indexes if h(k) = k.
3. Make a table: "Strategy", "Typical max α". Rows: chaining, open addressing.
4. Write five sentences that describe resize.

#### Medium practical tasks

1. Add a load-factor check to your chaining table. Resize when α > 0.75.
2. Count how many times you resize during 1 000 inserts from empty. Start at m = 8, double m.
3. In Go, use `make(map[int]int, 1000)` versus an empty map. Time 1 000 inserts. In Python, time dict inserts only.

#### Advanced practical tasks

1. Implement shrink of a table when α is very small. Avoid thrashing with a hysteresis rule (grow and shrink thresholds differ).
2. Read the load-factor note for one language map. Write ten STE sentences. Cite the source.

---

## Keys must be hashable and stable

A key must be **hashable**. The language can compute h(key). The key must also support equality.

A key must be **stable** while it sits in the table. If you change a field that the hash uses, the key is in the wrong bucket. Lookup fails. The table is corrupt.

```text
bad idea
  put object in map
  change object.id
  look up by the new id  →  miss
```

Use immutable keys when you can. Strings and numbers are safe. A mutable list is a bad key. Python does not allow a list as a dict key. A tuple of numbers is allowed.

In Go, map keys must be comparable. Slices, maps, and functions are not comparable. Structs are comparable if all fields are comparable. Do not mutate a struct field that you used as a key.

If two objects compare as equal, they must have the same hash. If you write a custom equality, you must write a matching hash.

A float key can be a problem. NaN is not equal to NaN in usual IEEE rules. Avoid float keys.

### Questions

#### Theoretical questions

1. What does hashable mean?
2. What does stable mean for a key in a table?
3. Why is a mutable list a bad key?
4. What Go types are not map keys?
5. Why can a float be a bad key?

#### Easy practical tasks

1. List five good key types in your language.
2. List two bad key types and one reason for each.
3. Write five sentences about a mutated key.
4. Make a table: "Key", "Hashable in Python?", "Comparable map key in Go?".

#### Medium practical tasks

1. In Python, try to use a list as a dict key. Record the error. Use a tuple instead.
2. In Go, try to use a slice as a map key. Record the compiler error. Use a string or a struct of comparable fields.
3. Explain in six sentences the rule "equal keys, same hash".

#### Advanced practical tasks

1. Write a small key type (struct or class) with a custom hash and equality that agree. Put the key in a table. Prove a lookup works.
2. Write a broken pair: equality ignores a field that the hash uses. Show a failed lookup. Then fix the pair.

---

## Average O(1) vs worst-case degradation

A hash table gives **average Θ(1)** find, insert, and delete. That bound needs a uniform hash, a bounded load factor, and a simple key mix.

**Worst-case** find is Θ(n). All keys can land in one chain. Open addressing can also degrade to a long probe of many slots.

Degradation happens when:

- the hash is constant or very weak
- an adversary picks colliding keys
- you never resize
- almost every slot is full (open addressing)

```text
average (good hash, bounded α):  Θ(1)
worst case:                      Θ(n)
```

A balanced tree map is Θ(log n) in the worst case. If you need a worst-case bound, use a tree. If you need speed on typical data, use a hash table.

Some hash tables fall back to a tree in a long bucket. That design improves the worst case to Θ(log n) for that bucket. You can treat that design as awareness.

Do not write "O(1)" without "average" when you describe a hash table. Write "average Θ(1), worst Θ(n)" unless the implementation gives a stronger worst case.

### Questions

#### Theoretical questions

1. What conditions support average Θ(1)?
2. What is the usual worst-case find?
3. Name three causes of degradation.
4. When do you prefer a tree map?
5. Why must you write "average" with O(1) for a hash table?

#### Easy practical tasks

1. Write a complete complexity statement for find in a chaining table.
2. Make a table: "Case", "Find cost". Rows: average good, worst one chain.
3. Write five sentences that compare a hash table and a tree map at a high level.
4. List two program types that can accept average Θ(1) and one type that may need worst-case Θ(log n).

#### Medium practical tasks

1. Force worst-case chaining with a constant hash. Time n finds. Then use a normal hash. Time n finds.
2. Find the documentation bound for Go `map` or Python `dict`. Copy the claim. Write it in STE.
3. Explain in six sentences why resize policy is part of the average-time story.

#### Advanced practical tasks

1. Plot or table the find time as α grows from 0.2 to 0.95 in an open-addressing table.
2. Write a one-page decision note: hash table versus tree map versus sorted array. Use this topic and the next topic.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe the path of find from key to value in chaining and in linear probing.
2. How do hash, equality, and stability work together?
3. Why do load factor and growth factor both appear in an amortized insert story?
4. A teammate stores a pointer to a map bucket. Which facts do you use to stop that design?
5. How do you explain average O(1) to a person who only knows worst-case Big-O?

#### Easy practical tasks

1. Write a one-page cheat sheet: hash, bucket, collision, chaining, open addressing, probe, load factor, resize, hashable, stable.
2. Draw one collision resolved by a chain and the same collision resolved by linear probing.
3. Compute α before and after a double of m if n stays the same.
4. Write three rules for keys in one short list.

#### Medium practical tasks

1. Implement a mini map: chaining, resize at α > 0.75, string keys. Tests: insert, find, update, delete, missing key.
2. Implement linear probing with tombstones. Compare code complexity with chaining in a short note.
3. Build a word-count program with your map or with the language map. Time a medium text file.

#### Advanced practical tasks

1. Implement chaining and open addressing with the same tests. Compare average probe or chain length at α = 0.5 and α = 0.8.
2. Read one production hash-map notes (Go runtime map, or CPython dict). Write a STE summary of grow, hash, and collision method. Cite the source.
