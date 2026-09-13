# 5. Hash Tables

## Description

A hash table maps a key to a value in average constant time. This topic explains hash functions, collisions, separate chaining, open addressing, load factor, resize, key rules, and worst-case cost.

Complete this topic after arrays and lists. You will use a hash table as a set and as a map in the next topic.

Use one term for each concept. A **key** is the lookup name. A **hash function** maps a key to an integer. A **bucket** is a slot in the table array. A **collision** occurs when two keys map to the same bucket. The **load factor** is n divided by the number of buckets. A key is **hashable** when you can compute a hash from it. A key is **stable** when its hash and its equality do not change while the key lives in the table.

---

## Hash function and collisions

A **hash function** h maps a key to an integer. The table uses `index = h(key) mod m` when the table has m buckets. Some tables use a bit mask when m is a power of 2.

The hash function must be **deterministic**. The same key must give the same hash while the key lives in the table.

The hash function must use the full key (or a defined part of the key). If you hash only the first letter of a string, many keys collide.

```text
h("cat")  →  81321
index     →  81321 mod m
```

A language hash is often not a cryptographic hash. Do not use a table hash as a password hash. A table hash is for distribution in memory.

A good hash is fast. A slow hash can dominate the lookup cost.

A good hash is **uniform**. Each bucket should receive about the same number of keys for typical input. Uniform does not mean random in a security sense. It means the indexes spread out.

A **collision** occurs when two different keys get the same bucket index. Collisions are normal. You cannot prevent all collisions when the number of keys is large. You must store more than one key per bucket, or you must probe another bucket.

Integer keys can use a mix function so that nearby integers do not all land in nearby buckets. String keys mix the characters. This topic requires the idea: map the key to an integer, then to a bucket.

### Questions

#### Theoretical questions

1. What does a hash function map?
2. How do you get a bucket index from a hash and m?
3. What does deterministic mean for a hash function?
4. What is a collision?
5. Why is a table hash not a password hash?

#### Easy practical tasks

1. Write a tiny hash for an integer key: h(x) = x. Then write h(x) mod 8 for x = 1, 8, 16.
2. Make a table: "Key", "simple hash", "hash mod 5" for three short words. Use the sum of letter positions (A=1).
3. Write five sentences that define a hash function and a collision.
4. List two key types that you can hash (for example integer and string).

#### Medium practical tasks

1. Write index(key, m) with a simple string hash: sum of character codes. Show collisions for m = 4.
2. Explain in six sentences why a fast hash matters for a large table.
3. Write why hashing only the first character is a poor function.

#### Advanced practical tasks

1. Implement a well-known string mix (look up a published formula). Compare collision counts with sum-of-codes on 1 000 words.
2. Read a short note on uniform hashing. Write ten STE sentences. Cite the source.

---

## Separate chaining and open addressing

**Separate chaining** stores a small collection in each bucket. The usual collection is a linked list of entries. Each entry holds a key, a value, and a next pointer.

To find a key, you hash to a bucket, then you walk the chain until you find an equal key or you reach the end.

```text
m = 4
bucket 0:  (cat → 1) → (dog → 2)
bucket 1:  empty
bucket 2:  (eel → 3)
bucket 3:  empty
```

Insert adds an entry to the chain. Delete removes the entry from the chain.

If the hash is uniform, the average chain length is the load factor. Find is then average Θ(1 + α) where α is n/m.

**Open addressing** stores entries in the table array itself. There is no chain. When a bucket is full, you **probe** another bucket.

**Linear probing** tries index, index+1, index+2, ... and wraps with modulus m.

```text
h(key) mod 5 = 2
cells:  [ ] [ ] [A] [B] [ ]
insert C that also hashes to 2
C goes to the first empty cell after 2
```

Find follows the same probe sequence until it finds the key or an empty cell.

Delete in open addressing is not a simple clear. A cleared cell can stop a later probe. Many tables store a **tombstone** (a deleted mark). Find skips tombstones. Insert can reuse a tombstone.

Open addressing needs a load factor below 1. The table must have empty cells. Chaining can grow past m entries because chains sit outside the array.

**Clustering** is a run of full cells in linear probing. Long clusters slow find and insert. Other probe methods exist (quadratic probing, double hashing). This topic needs the idea: a second formula picks the next cell.

### Questions

#### Theoretical questions

1. What does one chain store in separate chaining?
2. How does find walk in chaining?
3. What is a probe in open addressing?
4. Why does delete need a tombstone in open addressing?
5. Why must open addressing keep empty cells?

#### Easy practical tasks

1. Draw four buckets. Put two keys in bucket 0 as a chain. Put one key in bucket 2.
2. Draw linear probing of capacity 5 after three keys that all hash to 1.
3. Make a table: "Method", "Where the entry sits", "Load factor limit". Add two rows.
4. Write five sentences that contrast chaining and open addressing.

#### Medium practical tasks

1. Write the find steps for chaining and for linear probing. Include the not-found case.
2. Explain in six sentences why a tombstone is not the same as an empty cell.
3. Show a cluster of four full cells in a table of size 8. Write how a new key that hashes into the cluster behaves.

#### Advanced practical tasks

1. Implement chaining and linear probing for integer keys. Count probes for the same insert sequence.
2. Read a short note on quadratic probing or double hashing. Write eight STE sentences about clustering.

---

## Load factor and resize

The **load factor** α is n / m. n is the number of keys. m is the number of buckets (or cells in open addressing).

When α is small, many buckets are empty. Find is fast. Space is wasted.

When α is large, chains are long, or probes are long. Find is slow.

A table **resizes** when α passes a **threshold**. A common threshold is about 0.75 for open addressing. Chaining can use a similar threshold or a larger one.

Resize allocates a new array with a larger m. A common new size is about 2m. Then you **rehash**: you insert each old key into the new table. You do not copy chains as they are. The new modulus changes the indexes.

```text
old m = 4, n = 3, α = 0.75
resize to m = 8
rehash each key with mod 8
```

Rehash of n keys is Θ(n). That cost is like a dynamic-array grow. If you double m, the amortized cost of insert can stay Θ(1) for a good hash.

Do not resize on every insert. Do not pick a new m that is too small. α would stay large.

Some tables shrink when n becomes small. Shrink also rehashes. Shrink is optional.

m is often a power of 2 so that the index can use a mask. Some tables pick a prime m. The exact rule is an implementation detail. The idea stays: m must grow when n grows.

### Questions

#### Theoretical questions

1. What is the load factor?
2. What happens when the load factor is large?
3. When does a table resize?
4. Why must you rehash after resize?
5. Why is rehash Θ(n)?

#### Easy practical tasks

1. For n = 6 and m = 8, write α.
2. Draw a table of m = 4 with three keys. Then draw m = 8 after a doubling resize. Do not invent hashes. Mark that indexes can change.
3. Make a table: "α", "Space", "Find speed (qualitative)". Rows: 0.1, 0.7, 2.0 (chaining).
4. Write five sentences about load factor and resize.

#### Medium practical tasks

1. Write the steps of insert when insert would pass the threshold.
2. Explain in six sentences how doubling m keeps amortized insert cheap.
3. Write why you cannot move a chain to the same bucket index after m changes.

#### Advanced practical tasks

1. Implement chaining with resize at α > 0.75. Print m and α after each grow during 100 inserts.
2. Compare threshold 0.5 and 0.9 on the same key set. Write probe or chain-length totals.

---

## Keys must be hashable and stable

A key must be **hashable**. You must compute an integer hash from the key. You must also test **equality**. Find uses equality, not hash alone. Two keys can collide. Only equal keys are the same entry.

A key must be **stable** while it lives in the table. Stability means:

- the hash of the key does not change
- equality with other keys does not change

If you put a mutable object in the table and then change a field that the hash uses, the key sits in the old bucket. Find uses the new hash. Find looks in the wrong bucket. The entry is lost.

```text
bad:
  insert object {name: "Ann"}
  change name to "Bob"
  find "Bob" → wrong bucket
  find "Ann" → equality can also fail
```

Safe keys: integers, immutable strings, tuples of immutable values. Unsafe keys: a list or a record that you change after insert.

Two keys that are equal must have the same hash. That rule is the **hash-equality contract**. If equal keys have different hashes, find can miss.

Two keys that are not equal may have the same hash. That event is a collision. It is allowed.

Do not use a changing counter or a current time as part of a stored key.

Some languages forbid unhashable types as keys. Some languages do not. The rule is still yours: hashable and stable.

### Questions

#### Theoretical questions

1. What does hashable mean?
2. Why is equality required in addition to hash?
3. What does stable mean for a key?
4. What is the hash-equality contract?
5. Why can a mutated key become lost?

#### Easy practical tasks

1. List four safe key examples and two unsafe key examples.
2. Write five sentences about hashable and stable keys.
3. Make a table: "Pair of keys", "Same hash required?", "Collision allowed?". Rows: equal keys, unequal keys.
4. Draw the "lost key" picture: old bucket versus new hash after a field change.

#### Medium practical tasks

1. Write a small record with a name field. Describe insert, then a name change, then find. Write which step breaks.
2. Explain in six sentences why collisions are allowed but hash mismatch for equal keys is not.
3. Write a rule list of five lines for keys in a table that you design.

#### Advanced practical tasks

1. Implement a table that rejects a key type that you mark as mutable. Write the check that you use.
2. Read a language rule for hashable keys (documentation). Write ten STE sentences. Cite the source.

---

## Average O(1) vs worst-case degradation

A good hash table gives **average Θ(1)** find, insert, and delete. The average assumes a uniform hash and a bounded load factor.

The **worst case** can be Θ(n). All keys can land in one bucket. Chaining then walks n entries. Open addressing can also probe many cells if the table is full of collisions or tombstones.

```text
average:  hash spreads keys, α is bounded, few probes
worst:    all keys in one chain, find scans n entries
```

Degradation has more than one cause:

- a poor hash function
- a load factor that is too large
- many tombstones
- an attacker who picks keys that collide (in some old tables)

Resize and a better hash reduce the first two causes. Periodic rebuild can clear tombstones.

Average Θ(1) is not a guarantee for one operation. One find can still scan a long chain.

If you need a strict worst-case bound, a balanced tree map can give Θ(log n) for every call. Topic 6 and Topic 9 cover that choice.

Do not say "hash table is O(1)" without the words average and the assumptions.

Some tables switch a long chain to a tree. That hybrid improves the worst chain. It is an extra design. The basic model remains: average constant, worst linear unless you add more structure.

### Questions

#### Theoretical questions

1. What assumptions sit behind average Θ(1)?
2. What is the worst-case find cost in a chain of n keys?
3. List three causes of degradation.
4. Why is average Θ(1) not a guarantee for one find?
5. When do you pick a tree map for a worst-case bound?

#### Easy practical tasks

1. Write five sentences that contrast average Θ(1) and worst-case Θ(n).
2. Make a table: "Cause", "What you see". Add three rows.
3. Draw one bucket that holds all five keys. Write the find cost for the last key.
4. Write a complete sentence that replaces "the table is O(1)".

#### Medium practical tasks

1. Insert n keys that all use h(x) = 0 into a chaining table. Write the class of find.
2. Explain in six sentences how resize and a uniform hash prevent that picture.
3. Write when tombstones make open-address find slower than the live load factor suggests.

#### Advanced practical tasks

1. Build a chaining table with a bad hash and a good hash on the same keys. Plot or table the maximum chain length.
2. Write a one-page note: realtime code that must not pause for Θ(n) rehash. Name one alternative structure.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do hash, collision method, load factor, and key rules work together in one successful find?
2. When do you prefer chaining, and when do you prefer open addressing?
3. Why is rehash not an optional extra after you change m?
4. How does a mutated key break the same contract that collisions do not break?
5. What do you write instead of "hash tables are always O(1)"?

#### Easy practical tasks

1. Write a one-page cheat sheet: hash, bucket, collision, chaining, probe, tombstone, load factor, resize, stable key, average vs worst.
2. Draw one chaining picture and one linear-probing picture for the same three keys.
3. Make a table: "Operation", "Average", "Worst". Rows: find, insert, delete.
4. List five key-design rules in one column.

#### Medium practical tasks

1. Implement a chaining map with resize. Write tests: collide, missing key, overwrite same key, grow.
2. Write a short report that maps α to a qualitative find cost for chaining and for linear probing.
3. Design an ADT card for a map: get, put, delete, size. Write the key rules in the card.

#### Advanced practical tasks

1. Implement linear probing with tombstones and a rebuild that drops tombstones when their count is high.
2. Compare a hash map and a sorted array of pairs on n = 10 000 random finds. Write classes and one measured time pair.
