# 11. Binary Search Trees

## Description

A binary search tree (BST) stores keys in a binary tree with an order rule. This topic explains the BST property, search, insert, delete, tree walks, degenerate shape, and why balance matters. Complete this topic after recursion and trees. Complete hash maps before you compare exact-key find.

Use one term for each concept. The **BST property** is the order rule on left and right subtrees. A **successor** is the next key in sorted order. A **degenerate** BST is a line. **Balance** keeps height near log n.

---

## BST property

A **binary search tree** is a binary tree of keys (and optional values). For every node u:

- every key in the left subtree of u is **less than** the key of u
- every key in the right subtree of u is **greater than** the key of u

This handbook uses **strict** inequality. Duplicate keys are not stored. If you need duplicates, store a count in the node or store equal keys on one side and write that rule.

```text
        8
       / \
      3   10
     / \    \
    1   6    14
```

The BST property is stronger than "left child < u < right child". The rule applies to the whole subtree, not only to the child.

The in-order walk of a BST visits keys in sorted order. That fact is the main reason to use a BST.

Compare a BST with a hash table. A hash table does not keep keys sorted. A BST does. Find in a balanced BST is Θ(log n). Find in a hash table is average Θ(1).

Compare a BST with a sorted array. Find in both is Θ(log n) with binary search or with tree search. Insert in a sorted array is Θ(n). Insert in a BST is Θ(height).

Keys must be **comparable**. You need a total order: less, equal, greater.

### Questions

#### Theoretical questions

1. What is the BST property for node u?
2. Why is the rule about subtrees, not only about children?
3. How does this handbook treat duplicate keys?
4. What order does an in-order walk produce?
5. Why must keys be comparable?

#### Easy practical tasks

1. Draw a BST with keys 2, 1, 3. Label left and right.
2. Check the BST property on the 8–3–10 picture. Write yes or no for each node.
3. Make a table: "Structure", "Sorted keys?", "Find class (typical)". Rows: hash map, BST, sorted array.
4. Write five sentences that define a BST.

#### Medium practical tasks

1. Draw two binary trees with the same keys. One is a BST. One is not. Mark the violation.
2. Write a recursive function `isBST(node, min, max)` that checks the property. Test both drawings.
3. Explain in six sentences why "left child < parent < right child" is not enough.

#### Advanced practical tasks

1. Prove in eight STE sentences that in-order traversal of a BST is sorted. Use the subtree rule.
2. Implement a BST node with a key and a value. Document that the property uses the key only.

---

## Search, insert, delete (successor)

**Search** starts at the root. If the key equals the node, you found it. If the key is less, you go left. If the key is greater, you go right. If you reach null, the key is absent. Time is Θ(h) where h is the height.

**Insert** searches for the key. If the key exists, you reject or update the value. If you reach null, you attach a new node at that place. Time is Θ(h).

**Delete** has three cases:

1. The node is a leaf. Remove it. The parent link becomes null.
2. The node has one child. Replace the node with that child.
3. The node has two children. Find the **successor**. The successor is the minimum key in the right subtree (the next key in sorted order). Copy the successor key (and value) into the node. Then delete the successor node. The successor has no left child, so case 1 or case 2 applies.

```text
successor of 8 in the picture: minimum of the right subtree = 10
successor of 3: minimum of right subtree of 3 = 6
```

The **predecessor** is the maximum of the left subtree. Some delete implementations use the predecessor. Both methods are correct if you stay consistent.

```go
func Search(n *BinNode, k int) *BinNode {
	for n != nil {
		if k == n.Value {
			return n
		}
		if k < n.Value {
			n = n.Left
		} else {
			n = n.Right
		}
	}
	return nil
}
```

Do not delete a node with two children by a simple unlink. You would lose a subtree. Use the successor (or predecessor) method.

### Questions

#### Theoretical questions

1. How does search choose left or right?
2. Where do you attach a new node on insert?
3. What are the three delete cases?
4. What is the successor of a node with two children?
5. Why does the successor have no left child?

#### Easy practical tasks

1. Search for 6 and for 7 in the 8–3–10 picture. Write the path of keys.
2. Insert 4 into that picture. Draw the new tree.
3. Draw delete of a leaf. Draw delete of a node with one child.
4. Make a table: "Delete case", "What you do". Three rows.

#### Medium practical tasks

1. Implement search and insert. Test a sequence of inserts. Check `isBST`.
2. Implement minimum(node) as a walk to the left. Use it as successor when you delete a two-child node.
3. Delete the root of a small tree in all three cases. Draw before and after.

#### Advanced practical tasks

1. Implement delete with successor. Write tests for leaf, one child, two children, and missing key.
2. Implement delete with predecessor. Show that both trees stay a BST. Compare the two results on one example.

---

## In-order / pre-order / post-order / level-order

A **traversal** visits each node one time.

**In-order.** Walk left, visit the node, walk right. On a BST, the visit order is sorted keys.

**Pre-order.** Visit the node, walk left, walk right. Pre-order is useful to copy a tree or to print a prefix form.

**Post-order.** Walk left, walk right, visit the node. Post-order is useful to delete a tree or to evaluate an expression tree. You visit children before the parent.

**Level-order.** Visit the root, then depth 1, then depth 2. Level-order uses a **queue** (BFS). Level-order is not recursive in the simple form.

```text
        8
       / \
      3   10
     / \
    1   6

in-order:     1, 3, 6, 8, 10
pre-order:    8, 3, 1, 6, 10
post-order:   1, 6, 3, 10, 8
level-order:  8, 3, 10, 1, 6
```

```text
in-order(node):
  if node is null: return
  in-order(node.left)
  visit node
  in-order(node.right)
```

Time of each walk is Θ(n). Extra space is Θ(h) for the recursive walks (call stack). Extra space is Θ(width) for level-order (queue).

Use in-order when you need sorted keys. Use level-order when you need levels. Use post-order when you free nodes from the leaves up.

### Questions

#### Theoretical questions

1. What is the visit order of in-order on a BST?
2. When do you visit the node in pre-order?
3. Why does post-order help when you free a tree?
4. Which ADT does level-order use?
5. What extra space does a recursive walk use?

#### Easy practical tasks

1. Write the four orders for a tiny tree: root 2, left 1, right 3.
2. Make a table: "Traversal", "Left / visit / right order".
3. Draw a queue during the first three steps of level-order on the 8–3–10 tree.
4. Write five sentences that contrast in-order and level-order.

#### Medium practical tasks

1. Implement the three recursive walks. Print keys.
2. Implement level-order with a queue. Print keys.
3. Rebuild a BST from a sequence of inserts. Show that in-order is sorted and that pre-order is not always sorted.

#### Advanced practical tasks

1. Reconstruct a binary tree from pre-order and in-order of a known example (no duplicates). Document the method.
2. Implement iterative in-order with an explicit stack. Compare the output with recursive in-order.

---

## Degenerate tree (linked list)

A **degenerate** BST is a line. Each node has only a left child, or each node has only a right child. Height is n − 1.

You get a degenerate tree when you insert keys in sorted order into an empty BST:

```text
insert 1, 2, 3, 4
1
 \
  2
   \
    3
     \
      4
```

Search, insert, and delete become Θ(n). The BST is then as slow as a linked list for find. You still have extra pointers and poorer locality than a simple list or array.

A reverse-sorted insert sequence also makes a line (all left children).

Random insert order usually gives a height of about 2 ln n for a simple BST. That fact is average-case theory. Do not depend on a random insert order in production. An adversary can insert sorted keys on purpose.

A sorted array plus binary search does not become a line. The array height of the search is always about log n. The cost is insert, not find.

Detect a bad shape: compare height with log₂ n. If n is 1 000 and height is 999, the tree is degenerate.

### Questions

#### Theoretical questions

1. What is a degenerate BST?
2. Which insert order creates a right line?
3. What is the find cost in a degenerate BST?
4. Why can an adversary force this shape?
5. Why does a sorted array avoid this find cost?

#### Easy practical tasks

1. Draw the BST after insert 5, 4, 3.
2. Draw the BST after insert 3, 4, 5.
3. Make a table: "Insert order", "Shape". Two rows: sorted, mixed.
4. Write five sentences about find cost when height is n − 1.

#### Medium practical tasks

1. Insert n = 1, 2, ..., 20 into a BST. Print height. Compare with log₂ 20.
2. Insert a shuffled sequence of the same keys. Print height.
3. Time search of the last key in a sorted-insert tree of n = 5 000 if you can. Write the time.

#### Advanced practical tasks

1. Write a function that reports whether height > 2 * log2(n+1) (a simple alarm). Test on a line and on a mixed tree.
2. Write a one-page note: why production maps use a hash table or a balanced tree, not a raw BST.

---

## Why balance matters

**Balance** keeps height small. If height is Θ(log n), search, insert, and delete are Θ(log n). If height is Θ(n), those operations are Θ(n).

n = 1 000 000. log₂ n is about 20. A balanced tree does about 20 compares. A degenerate tree can do 1 000 000 compares.

Balance does not change the BST property. Balance changes the shape. AVL trees and red-black trees restore balance after insert and delete. Those methods are the next topic after this one. This topic only states the need.

```text
same keys, two shapes
  balanced: height ≈ log n    → fast
  line:     height = n - 1    → slow
```

A complete tree is balanced in a practical sense. A random BST is often shallow enough for a learning program. A real system that must not degrade uses a balanced tree or a hash table.

Rotation is the usual local change that shortens a long path. You will learn rotation in the balanced-tree topic. Do not invent ad-hoc swaps that break the BST property.

When you implement a BST for learning, always print height and in-order. Height tells you the shape. In-order tells you the property.

### Questions

#### Theoretical questions

1. What does balance keep small?
2. Why does height control find time in a BST?
3. Does balance change the BST property?
4. Why is a random insert order not enough for a real system?
5. What will a later topic add (AVL or red-black) in one sentence?

#### Easy practical tasks

1. For n = 16, write log₂ 16 and n − 1. Write both as possible heights.
2. Make a table: "Shape", "Height class", "Find class". Rows: balanced, degenerate.
3. Write five sentences that explain why 20 compares versus 1 000 000 compares matters.
4. List two production structures that avoid a raw BST.

#### Medium practical tasks

1. Build two trees of the same 15 keys: one shuffled, one sorted. Print both heights.
2. Time n finds on both trees. Write the two times.
3. Explain in six sentences what a rotation must preserve (the BST property) and what it changes (height).

#### Advanced practical tasks

1. Read a short AVL or red-black overview (rules only). Write ten STE sentences about why rotations exist. Do not implement yet.
2. Write a report: for a dictionary of exact-key finds, when you pick a hash table, when you pick a balanced tree, and when a learning BST is enough.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe the path of insert from the root to a new leaf. Name compares and null.
2. How do successor delete and in-order order work together?
3. Why do four traversals exist if in-order already sorts a BST?
4. How do a degenerate shape and the lack of balance produce the same cost class?
5. A teammate says a BST is always O(log n). Which facts do you use to fix that sentence?

#### Easy practical tasks

1. Write a one-page cheat sheet: BST property, search, insert, three delete cases, successor, four traversals, height, degenerate, balance.
2. Draw insert of 2, 1, 3, 4 and then delete 3 (two children if you insert a different set so that 3 has two children).
3. Write the four traversal names and one use each.
4. Write a complete complexity statement for search: include height, and include n when the tree is a line.

#### Medium practical tasks

1. Implement a BST: insert, search, in-order print, height, delete. Tests: empty, one node, a small mixed tree, a sorted insert line.
2. Compare in-order of the BST with `sort` of the same keys. They must match.
3. Time insert of 10 000 sorted keys versus 10 000 shuffled keys. Write heights and times.

#### Advanced practical tasks

1. Implement level-order and iterative in-order. Use them to dump a tree for debugging after each delete.
2. Read the next topic list (balanced search trees). Write a study plan of five steps from this raw BST to AVL or red-black. Do not implement balance in this file.
