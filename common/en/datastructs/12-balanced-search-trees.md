# 12. Balanced Search Trees

## Description

A balanced search tree keeps height near log n after insert and delete. This topic explains AVL rotations, red-black rules, 2-3 and 2-3-4 trees, B-trees, B+ trees, and why a database uses a B+ tree. Complete this topic after binary search trees.

Use one term for each concept. A **rotation** is a local change of parent and child that keeps the BST property. The **balance factor** of a node is the height of the left subtree minus the height of the right subtree. The **order** of a B-tree is the maximum number of children of one node.

---

## AVL rotations (idea + one implementation)

An **AVL tree** is a BST. For every node, the heights of the two child subtrees differ by at most 1. The balance factor is in {−1, 0, 1}. After insert or delete, a node can leave that set. You restore the rule with one or two rotations.

A **right rotation** on node y moves the left child x up. Node y becomes the right child of x. The old right subtree of x becomes the left subtree of y. The BST property stays true because all keys in that middle subtree sit between x and y.

A **left rotation** is the mirror. You use it when the right side is too high.

Four insert cases exist:

1. **LL.** The extra key is in the left child of the left child. Do one right rotation.
2. **RR.** The extra key is in the right child of the right child. Do one left rotation.
3. **LR.** The extra key is in the right child of the left child. Do a left rotation on the left child, then a right rotation on the node.
4. **RL.** The extra key is in the left child of the right child. Do a right rotation on the right child, then a left rotation on the node.

```text
right rotation on y

      y                x
     / \              / \
    x   C    →       A   y
   / \                  / \
  A   B                B   C
```

```text
rotateRight(y):
  x = y.left
  y.left = x.right
  x.right = y
  update height of y
  update height of x
  return x
```

Store the height in each node. After each local change, write height = 1 + max(height of left, height of right). A null child has height −1.

Implement **insert plus rotations** first. Walk the insert path from the new leaf back to the root. Repair the first node that is not AVL. One insert needs at most two rotations.

Do not implement delete in the same week if insert is not clear. Delete uses the same rotations. The walk up the path is longer.

Search in an AVL tree is the same as search in a BST. Time is Θ(log n) because height is Θ(log n).

### Questions

#### Theoretical questions

1. What is the AVL height rule for every node?
2. What is the balance factor?
3. What does a right rotation change, and what does it keep?
4. Which two rotations does the LR case use, and in which order?
5. Why is search in an AVL tree Θ(log n)?

#### Easy practical tasks

1. Draw the right-rotation picture. Label A, B, C, x, and y.
2. Write the four case names LL, RR, LR, and RL. Write one sentence for each case.
3. Make a table: "Case", "First rotation", "Second rotation". Four rows.
4. For keys 1, 2, 3 inserted in that order, draw the tree before and after the RR repair.

#### Medium practical tasks

1. Implement a node with key, left, right, and height. Write `heightOf` and `balanceFactor`.
2. Implement `rotateRight` and `rotateLeft`. Test them on a three-node line and check the BST property.
3. Implement AVL insert for the four cases. Insert 1, 2, 3 and 3, 2, 1 and 3, 1, 2. Draw each result.

#### Advanced practical tasks

1. Implement AVL insert for n random keys. After each insert, check the AVL rule and the BST property on every node.
2. Compare height of a raw BST and an AVL tree on the same sorted insert of n = 1 000. Write both heights.

---

## Red-black tree (rules; implement or study, not both at first)

A **red-black tree** is a BST. Each node has a color: red or black. Null children count as black. The rules keep the longest root-to-null path at most twice the shortest root-to-null path. Height is then Θ(log n).

The usual rules are:

1. Each node is red or black.
2. The root is black.
3. Every null leaf is black.
4. A red node has two black children. Two red nodes do not sit as parent and child.
5. Every path from a node to a null in that subtree has the same number of black nodes. That number is the **black height**.

```text
example shape (B = black, R = red)

        7B
       /  \
     3R    10B
    /  \     \
  1B    5B    12R
```

A red-black tree is the binary form of a 2-3-4 tree. A black node with a red child is one 3-node. A black node with two red children is one 4-node. You can study that map before you write code.

Many libraries use a red-black tree for an ordered map. Examples: Java `TreeMap`, C++ `std::map`. The library code is large. Do not start there.

**First week rule.** Either implement a small red-black insert from a textbook, or study the five rules and the 2-3-4 map. Do not do both at the same time. AVL insert is a better first implementation if you want to write code.

Repair uses recolor and rotation. Recolor changes red to black or black to red. Rotation is the same geometric change as in AVL.

Search does not use color. Search is the BST walk. Color matters only for insert and delete.

### Questions

#### Theoretical questions

1. What are the five red-black rules in this handbook?
2. What is black height?
3. Why do the rules bound the height by a constant times log n?
4. How does a black parent with one red child relate to a 3-node?
5. Why must you not implement and study every proof in the same week?

#### Easy practical tasks

1. Color a three-node BST so that it obeys the five rules. Write the colors.
2. Make a table: "Rule number", "Rule in one sentence". Five rows.
3. Mark one illegal pair of red parent and red child on a drawing. Write which rule fails.
4. Write five sentences that contrast AVL (height numbers) and red-black (colors).

#### Medium practical tasks

1. Take the 7–3–10 picture. Check each rule. Write pass or fail for each rule.
2. Draw the 2-3-4 node that matches a black node with two red children. Label the three keys.
3. Read one textbook page on red-black insert. Write the first three repair steps in STE. Do not copy the page.

#### Advanced practical tasks

1. Pick one path: implement red-black insert from a single reference, or write a two-page rule sheet with six legal and six illegal drawings. Do not do both.
2. Find the ordered-map type in one language. Write whether the documentation names red-black, AVL, or only "balanced tree".

---

## 2-3 / 2-3-4 trees (idea)

A **2-3 tree** and a **2-3-4 tree** store keys in nodes that can hold more than one key. All leaves sit at the same depth. The tree stays balanced because a split or a merge keeps that depth rule.

Node types:

- A **2-node** holds one key and has two children (or is a leaf).
- A **3-node** holds two keys in sorted order and has three children.
- A **4-node** holds three keys in sorted order and has four children. A 2-3 tree does not use 4-nodes. A 2-3-4 tree does.

Search in a node compares the key with the one, two, or three keys, then follows the correct child.

Insert adds the key into a leaf. If the leaf becomes too full, you **split**. The middle key moves to the parent. The left keys stay in a new left node. The right keys stay in a new right node. If the parent becomes too full, you split again. A split at the root grows the height by 1. Every leaf stays at the same new depth.

```text
split a 4-node with keys a < b < c

   [a b c]     →      [b]
                     /   \
                  [a]     [c]
```

A 2-3-4 tree maps to a red-black tree. You do not need both implementations. Learn the split idea here. Use it again when you read B-trees.

Do not write a full 2-3-4 library as your first balanced tree. Draw splits on paper. Then implement AVL or a B-tree of small order.

### Questions

#### Theoretical questions

1. What is a 2-node, a 3-node, and a 4-node?
2. Why are all leaves at the same depth?
3. What happens when you split a full node?
4. When does the height of the tree grow by 1?
5. How does a 2-3-4 tree relate to a red-black tree?

#### Easy practical tasks

1. Draw a 2-node, a 3-node, and a 4-node. Write the key counts.
2. Make a table: "Tree name", "Largest node type". Two rows: 2-3, 2-3-4.
3. Draw the split of keys 2, 4, 6 as in the picture.
4. Write five sentences that define search in a 3-node.

#### Medium practical tasks

1. Insert 1, 2, 3, 4 into an empty 2-3-4 tree on paper. Show each split.
2. Explain in six sentences why a split keeps all leaves at one depth.
3. Map one 4-node to a black node with two red children. Draw both.

#### Advanced practical tasks

1. Write a one-page note that compares split (2-3-4) with rotation (AVL). Use one insert sequence of four keys for both ideas.
2. Implement a 2-node and 3-node search only (no insert). Test a fixed drawing.

---

## B-trees and B+ trees (databases and filesystems)

A **B-tree** is a balanced search tree with a large **fanout**. Fanout is the number of children of a node. A node is one disk page or one filesystem block. One node holds many keys. Height is small: often 3 or 4 for a large file.

A B-tree of order m has at most m children and at most m − 1 keys. The exact min and max counts follow the textbook definition that you pick. The idea is the same: a wide node, a short path.

Search starts at the root. You binary-search the keys in the node (or you scan them). You follow one child pointer. You repeat until you find the key or you reach a leaf.

Insert can split a full node. The middle key goes to the parent. This is the 2-3-4 split at a larger size.

A **B+ tree** is a B-tree variant. All **records** (the full values) live in the **leaves**. Internal nodes store **separator keys** and child pointers. They do not store the full record. Leaves link to the next leaf. A range scan walks the leaf list.

```text
B-tree idea: keys and values in internal nodes and in leaves

B+ tree idea:
  internal:  separators + child pointers
  leaves:    all records, linked left to right
```

A filesystem directory and a database index often use a B+ tree or a close form. The page size matches the disk or the SSD read.

Do not confuse a B-tree with a binary tree. The "B" is not "binary". A B-tree node is multi-way.

### Questions

#### Theoretical questions

1. What is fanout in a B-tree?
2. Why is the height of a B-tree small for a large file?
3. Where do records live in a B+ tree?
4. What extra link do B+ leaves have?
5. Why is a B-tree not a binary tree?

#### Easy practical tasks

1. Draw a B-tree node with three keys and four children. Label keys and pointers.
2. Make a table: "Tree", "Where the record lives". Rows: B-tree, B+ tree.
3. Write five sentences that define a range scan on B+ leaves.
4. For page size 4 KiB and 8-byte keys plus 8-byte pointers, estimate a rough fanout. Show the division.

#### Medium practical tasks

1. Search for a key on a drawing of a small B+ tree. Write the nodes that you visit.
2. Split a full leaf of four keys (use a tiny order). Show the separator that moves up.
3. Explain in six sentences why one disk page should hold one node.

#### Advanced practical tasks

1. Read one B+ diagram in a database or filesystem text. Write ten STE sentences that name root, internal node, leaf, and sibling link.
2. Implement an in-memory B+ tree of a tiny order (for example max 4 keys). Support insert and range iterate. Do not add disk I/O yet.

---

## Why databases use B+ trees

A database index must find one key and must scan a range of keys. A B+ tree supports both.

**Point find.** Height is small. Each level is one page read if the page is not in cache. A few reads find one row.

**Range scan.** You find the first leaf. Then you follow leaf links. You do not walk back to the root for each next key. Sequential leaf pages are friendly to prefetch.

**Packing.** Internal pages store only separators and pointers. They store more keys per page than a B-tree that also stores full records in internal nodes. Higher fanout means fewer levels.

**Balance.** Splits keep all leaves at the same depth. Cost does not become linear in n after many inserts.

A hash index is fast for exact key find. A hash index does not give a sorted range. A database that needs `WHERE id BETWEEN a AND b` or `ORDER BY` on the key prefers a B+ tree.

A raw BST or an AVL tree of single rows is a poor disk index. Each node can sit in a different page. You would read one page per node. Fanout 2 makes height large.

```text
need
  exact key     →  B+ tree or hash
  sorted range  →  B+ tree
  few page I/Os →  high fanout, short height
```

Filesystems use the same idea for large directories and for some file-block maps. The page is a disk block.

### Questions

#### Theoretical questions

1. Why does a small height help a database point find?
2. How does a leaf link help a range scan?
3. Why can a B+ internal page hold more keys than a B-tree page that stores full records?
4. When is a hash index not enough?
5. Why is a binary AVL tree a poor disk index?

#### Easy practical tasks

1. Write five sentences that answer "why B+ trees" with find, range, fanout, and balance.
2. Make a table: "Need", "Structure". Rows: exact key only, sorted range, disk pages.
3. Draw three levels: root, one internal node, two leaves with a link. Mark a range scan.
4. List two SQL-style needs that match a B+ tree (use your own words, not a full SQL course).

#### Medium practical tasks

1. Compare in eight sentences: hash table in RAM, AVL in RAM, B+ tree on disk.
2. Count page reads for a point find of height 3 if the root is in cache. Write the count.
3. Explain in six sentences why leaf order matters for `BETWEEN`.

#### Advanced practical tasks

1. Write a one-page design note: an index on a 1 000 000-row table. Estimate height if fanout is 100. Show the arithmetic.
2. Read why one real database documents B+ (or a close name). Write ten STE sentences. Cite the source.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe the path from a raw BST problem (sorted inserts) to a structure that a database can use. Name height, rotation or split, and pages.
2. How do AVL height numbers, red-black colors, and 2-3-4 node widths solve the same problem?
3. A teammate says "use a balanced tree" for a disk index of 10⁹ keys. Which facts do you use to pick B+ over AVL?
4. What stays the same in AVL, red-black, and B+ search? What changes in repair?
5. Why does this handbook tell you to implement one rotation tree or study red-black rules, not both at first?

#### Easy practical tasks

1. Write a one-page cheat sheet: AVL cases, red-black rules, 2-3-4 split, B versus B+, database reasons.
2. Draw one insert that needs an AVL rotation and the same keys in a 2-3-4 split. Do not force the drawings to match in shape.
3. Make a table: "Structure", "Balance tool", "Typical home". Five rows.
4. Write one complete cost sentence for find in a balanced RAM tree and one for find in a B+ tree on disk.

#### Medium practical tasks

1. Implement AVL insert. Keep a raw BST insert. Time sorted insert of n = 10 000 on both. Write heights and times.
2. Walk a teammate through a B+ range scan with only a whiteboard. Record the six steps that you said.
3. Map one red-black drawing to 2-3-4 nodes. Write which red edges you grouped.

#### Advanced practical tasks

1. Implement one RAM balanced tree (AVL) and one tiny in-memory B+ tree. Use the same 200 keys. Compare iterate-in-order code paths.
2. Read the next topic list (heaps). Write five sentences on what a heap does not replace in a search tree.
