# 9. Balanced Trees and B-trees

## Description

A balanced search tree keeps height logarithmic after insert and delete. This topic explains the idea of AVL rotations, red-black rules (study the rules before you implement AVL and red-black), the idea of 2-3 and 2-3-4 trees, B-trees and B+ trees, and why databases use B+ trees.

Complete this topic after binary search trees. Do not implement AVL and red-black in the same week until you can state both rule sets.

Use one term for each concept. A **rotation** is a local pointer change that preserves the BST property and changes height. A **red-black tree** is a BST with color rules that bound the height. A **B-tree** is a search tree with many keys per node. A **B+ tree** stores all records in leaves and links the leaves.

---

## AVL rotations (idea)

An **AVL tree** is a BST that keeps the **balance factor** of every node in {−1, 0, +1}. The balance factor is height(left) − height(right), or the opposite sign in some books. This handbook uses **left height minus right height**.

After insert or delete, one or more nodes can leave that set. An AVL tree repairs the shape with **rotations**.

A **right rotation** at y lifts the left child x. The right subtree of x becomes the left subtree of y.

```text
before:     y              after:    x
           /                    / \
          x                    A   y
         / \                      /
        A   B                    B
```

A **left rotation** is the mirror.

The BST property stays. All keys in A stay less than x. All keys in B stay between x and y.

Four insert cases exist:

- **Left-left:** one right rotation
- **Right-right:** one left rotation
- **Left-right:** left rotation on the child, then right rotation on the node
- **Right-left:** the mirror pair

You do not need the full code in this section. You need the idea: a constant number of pointer changes restores the AVL rule on the path to the root. Height stays O(log n). Find, insert, and delete are Θ(log n).

AVL is stricter than red-black. AVL can do more rotations. Lookups can be a little faster. Inserts can be a little slower. Measure if you care.

Do not rotate without a height or balance-factor update. The next operation would see a wrong factor.

### Questions

#### Theoretical questions

1. What is the AVL balance-factor rule?
2. What does a rotation preserve?
3. What does a rotation change?
4. Why are there four insert cases?
5. What height class does AVL keep?

#### Easy practical tasks

1. Draw the right-rotation picture. Label A, x, B, y.
2. Make a table: "Case", "Rotations". Add the four insert cases.
3. Write five sentences that define AVL at the idea level.
4. Write why the BST property still holds after a right rotation.

#### Medium practical tasks

1. Insert 1, 2, 3 into an AVL tree on paper. Show the left rotation that removes the line.
2. Explain in six sentences the left-right case as two rotations.
3. Compare a plain BST line of 3 keys with the AVL tree of the same keys.

#### Advanced practical tasks

1. Implement insert with the four rotation cases. Test sorted keys 1 .. 20. Write the final height.
2. Read one AVL delete note. Write eight STE sentences about why delete can rotate on the way up. Do not paste code.

---

## Red-black rules (study before you implement both)

A **red-black tree** is a BST where each node has a color: **red** or **black**. Null children count as black.

A common rule set is:

1. Each node is red or black.
2. The root is black.
3. Every leaf (null) is black.
4. A red node has two black children. (No two reds in a parent–child pair.)
5. Every path from a node to a descendant null has the same number of black nodes.

```text
example shape (B = black, R = red):
          7B
         /  \
       3R    18B
            /  \
          10R   22R
```

These rules force height to stay O(log n). A long path cannot hold too many blacks. It also cannot hold too many reds because reds cannot stack.

Repair after insert uses recolor and rotations. The details are longer than AVL. Study the rules first. Draw small trees. Then implement.

**Study order.** Learn AVL rotations and the AVL factor. Learn the five red-black rules and one insert repair picture. Do not implement both structures in parallel until you can state both contracts without notes.

A red-black tree is a common implementation of an ordered map in libraries. It does fewer rotations than AVL on insert in many models. Height can be larger than AVL (still logarithmic).

Do not treat color as data. Color is a balance mark.

A 2-3-4 tree (next section) maps to a red-black tree. That map is a study aid. You can learn one view and then the other.

### Questions

#### Theoretical questions

1. List the five red-black rules in your words.
2. Why do null children count as black?
3. Why can two reds not sit as parent and child?
4. Why do you study rules before you implement AVL and red-black together?
5. How does a red-black height bound compare to a degenerate BST?

#### Easy practical tasks

1. Color a three-node balanced BST so that the rules hold. Write the colors.
2. Draw a parent–child red pair. Write which rule fails.
3. Make a table: "Structure", "What you store extra". Rows: AVL, red-black.
4. Write five sentences that state why you must not implement both in the same first week.

#### Medium practical tasks

1. Count black nodes on two root-to-null paths in a small legal tree. Write that the counts match.
2. Explain in six sentences why the rules ban a long red chain and a long black imbalance.
3. Write a study plan of five days: rules, pictures, AVL insert, red-black insert, tests.

#### Advanced practical tasks

1. Implement red-black insert only after you write the five rules from memory. Test sorted inserts.
2. Compare heights of AVL and red-black on the same key set. Write the two heights and one sentence on rotations.

---

## 2-3 / 2-3-4 trees (idea)

A **2-3 tree** is a search tree whose nodes are **2-nodes** or **3-nodes**.

- A 2-node holds 1 key and 2 children (or it is a leaf).
- A 3-node holds 2 keys and 3 children (or it is a leaf).

All leaves sit at the same depth. The tree is perfectly balanced in height.

```text
3-node [10 | 20]
children:  <10    10..20    >20
```

Insert adds a key into a leaf. If a node would hold too many keys, it **splits**. The middle key moves to the parent. Split can propagate to the root. A new root increases height by 1. Every leaf still has the same depth.

A **2-3-4 tree** also allows a **4-node**: 3 keys and 4 children. Insert can **split a 4-node** on the way down. That method keeps overflow local.

```text
4-node [5 | 10 | 15] splits
  parent receives 10
  left 2-node [5], right 2-node [15]
```

2-3-4 trees and red-black trees represent the same idea. A 2-node is a black node. A 3-node is a black node with one red child. A 4-node is a black node with two red children. You can study splits in 2-3-4 form, then map them to colors and rotations.

You do not need a full implementation in this section. You need the idea: many keys in a node, split, shared leaf depth.

B-trees generalize 2-3-4 trees to a large number of keys per node.

### Questions

#### Theoretical questions

1. What does a 2-node hold?
2. What does a 3-node hold?
3. Why do all leaves have the same depth in a 2-3 tree?
4. What happens when a node splits?
5. How does a 4-node map to red-black colors at a high level?

#### Easy practical tasks

1. Draw a single 3-node with two keys and three empty children.
2. Make a table: "Node kind", "Keys", "Children". Rows: 2-node, 3-node, 4-node.
3. Write five sentences that define 2-3 and 2-3-4 trees.
4. Draw a 4-node split. Show the key that moves up.

#### Medium practical tasks

1. Insert 1, 2, 3 into an empty 2-3 tree on paper. Show a split.
2. Explain in six sentences why split keeps all leaves at one depth.
3. Write the mapping: 2-node, 3-node, 4-node to black and red children.

#### Advanced practical tasks

1. Simulate 2-3-4 insert of ten keys on paper. Count splits.
2. Read a short note that maps 2-3-4 insert to red-black insert. Write ten STE sentences.

---

## B-trees and B+ trees

A **B-tree** of **order** m (definitions vary) stores many keys in each node. Each internal node has a number of children in a fixed range, except the root. All leaves sit at the same depth.

A node holds a sorted list of keys and a list of child pointers. Search uses binary search (or a scan) inside the node, then follows one child.

```text
internal node:
  keys:     10      20      30
  children: c0  c1      c2      c3
  keys in c1 are between 10 and 20
```

Insert can split a full node. Delete can merge or borrow keys. Height grows only when the root splits. Height stays small because the branching factor is large.

A **B+ tree** is a B-tree variant:

- **Internal nodes** store keys (separators) and child pointers. They do not store full records.
- **Leaves** store all keys and records (or pointers to records).
- Leaves form a **linked list** (often doubly linked) so that a range walk is a linear scan of leaves.

```text
internal:     [  20  ]
             /        \
leaves:  [1,4,9] ⇄ [20,21,25]
```

Search in a B+ tree always goes to a leaf. A range [L, R] finds the leaf for L, then walks sibling leaves until keys pass R.

B-tree height is O(log_b n) where b is the branching factor. For large b, height is a small integer: 2, 3, or 4 even for large n.

Do not mix B-tree and binary tree. The B is a name, not "binary".

Order and minimum occupancy have more than one textbook definition. When you read a page, copy that page's rule.

### Questions

#### Theoretical questions

1. Why does a B-tree node hold many keys?
2. What happens when a full node splits?
3. Where does a B+ tree store records?
4. Why are B+ leaves linked?
5. Why is height small when the branching factor is large?

#### Easy practical tasks

1. Draw an internal node with two keys and three children. Write the key ranges of the children.
2. Draw a tiny B+ tree with one internal node and two leaves. Link the leaves.
3. Make a table: "Structure", "Where records sit", "Range walk". Rows: B-tree, B+ tree.
4. Write five sentences that define B-trees and B+ trees.

#### Medium practical tasks

1. Search for key 21 in the picture in this section. Write the nodes that you visit.
2. Explain in six sentences how a range query uses the leaf list.
3. Write why an internal key in a B+ tree can be a copy of a leaf key.

#### Advanced practical tasks

1. Implement a tiny B+ tree with a small order (for example, max 4 keys per node). Support insert and range print.
2. Compare height of a binary BST and a B-tree of order 8 for n = 1000 as a formula. Write the two log expressions.

---

## Why databases use B+ trees

A database index often lives on **disk** or on flash. A read of one node is a **page** I/O. I/O is much slower than a compare in memory.

A B+ tree fits many keys in one page. One I/O moves a large node. Height is 3 or 4 for a large table. A binary tree of the same n keys would need many more I/Os along the path.

Leaves store rows in key order (or pointers to rows). A range query (`WHERE id BETWEEN 100 AND 200`) walks leaf pages in order. The sibling links avoid a new root-to-leaf walk for each next key.

Internal pages stay hot in cache. The root and the first levels fit in memory. Most I/O hits leaves.

```text
binary tree on disk:  many small nodes, many I/Os per find
B+ tree on disk:      few large pages, few I/Os per find
                      range = walk leaf pages
```

A hash index can find one key fast. It does not give a sorted range walk. Databases still use B+ trees for primary keys and for many secondary indexes.

Do not implement a production database index from this section. Use the idea: large nodes, all data in leaves, linked leaves, few I/Os.

Memory-only maps often use red-black trees or hash tables. They do not pay page I/O. The B+ design is for external memory. Topic 15 returns to data that does not fit in RAM.

### Questions

#### Theoretical questions

1. Why is one page I/O expensive compared to a memory compare?
2. Why do many keys per page reduce find I/O?
3. How does a B+ leaf list help a range query?
4. Why can a hash index be a poor range index?
5. Why do memory maps not need B+ trees for the same reason?

#### Easy practical tasks

1. Write five sentences that explain B+ trees in a database.
2. Make a table: "Need", "Hash index or B+ index". Rows: one id, a range of ids, sorted scan.
3. Draw three levels: root page, internal page, leaf page. Mark one find path.
4. Write why the root stays in cache.

#### Medium practical tasks

1. Estimate height if each page holds 100 child pointers and n = 1 000 000 keys. Write a log100 estimate.
2. Explain in six sentences why a binary tree of 1 000 000 nodes is a poor disk index.
3. Write a short note: clustered leaf order versus a pointer to a heap row (idea only).

#### Advanced practical tasks

1. Read a database documentation page on B-tree indexes. Write ten STE sentences. Cite the source.
2. Compare a memory red-black map and a B+ tree for n keys that fit in RAM. Write when each design wins.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do rotations, colors, and splits each keep height logarithmic?
2. Why does this topic tell you to study red-black rules before you implement AVL and red-black together?
3. How is a 2-3-4 node a study view of a red-black node?
4. What does a B+ tree add to a B-tree that a database range query uses?
5. When do you use a memory balanced BST, and when do you use a B+ tree?

#### Easy practical tasks

1. Write a one-page cheat sheet: AVL factor, rotation, five red-black rules, 2-3-4 split, B-tree, B+ leaf list, disk page.
2. Draw one rotation, one 4-node split, and one B+ leaf pair.
3. Make a table: "Tree", "Extra data", "Height idea". Add AVL, red-black, B-tree.
4. List four study mistakes: implement both color trees at once, mix B and binary, store records in B+ internals, ignore I/O.

#### Medium practical tasks

1. Write a comparison table: find class, insert repair idea, typical use (memory map vs disk index).
2. Simulate sorted insert of 1 .. 7 in AVL and in a 2-3 tree on paper. Write the two final shapes.
3. Write an ADT card for an ordered map. Then write which balanced tree you would pick in RAM versus on disk.

#### Advanced practical tasks

1. Implement one structure only: AVL insert or a tiny B+ tree. Write why you did not implement both in the same pass.
2. Read CLRS or a similar chapter on red-black or B-trees. Map each textbook term to this topic in a ten-row table.
