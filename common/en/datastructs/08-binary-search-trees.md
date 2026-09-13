# 8. Binary Search Trees

## Description

A binary search tree (BST) is a binary tree with a key order. This topic explains the BST property, search, insert, delete, tree walks, degenerate trees, and why balance matters.

Complete this topic after recursion and trees. Complete this topic before balanced trees.

Use one term for each concept. The **BST property** is the key rule that makes search fast. A **tree walk** visits each node one time. A **degenerate** tree is a tree that has become a line. **Balance** keeps height logarithmic.

---

## BST property

A **binary search tree** is a binary tree of keys (and optional values) that satisfies the **BST property**:

- every key in the left subtree of node u is less than the key of u
- every key in the right subtree of node u is greater than the key of u

Some implementations allow equal keys. A common rule puts equals to the right, or stores a count. This handbook uses **distinct keys** unless a line says otherwise.

```text
        8
       / \
      3    10
     / \     \
    1   6     14
       / \
      4   7
```

Check one node: all keys left of 8 are less than 8. All keys right of 8 are greater than 8. The same rule holds at 3 and at every other node.

The property is **recursive**. It is not enough that the left child is smaller and the right child is larger. A deep left descendant must still be less than the ancestor.

```text
not a BST:
      5
     / \
    3   7
   /
  6      ← 6 is in the left subtree of 5, but 6 > 5
```

In-order walk of a BST visits keys in sorted order. That fact is a check and a feature.

A heap (Topic 10) is a binary tree with a different property. Do not mix the two.

A BST is an implementation of an ordered set or ordered map (Topic 6). The ADT does not name the tree. The tree gives Θ(h) find, where h is the height.

### Questions

#### Theoretical questions

1. What is the BST property?
2. Why is "left child smaller, right child larger" not enough?
3. What does in-order walk produce on a BST?
4. How does a BST implement an ordered map?
5. How is the BST property different from a heap property?

#### Easy practical tasks

1. Check the 8–3–10 picture. Write pass or fail for the property at node 3.
2. Draw the "not a BST" picture. Mark the key that breaks the property.
3. Write the in-order key list for the 8–3–10 picture.
4. Write five sentences that define the BST property.

#### Medium practical tasks

1. Insert the missing check: for each node, write the allowed key range (low, high). Use −∞ and +∞ at the root.
2. Explain in six sentences why a descendant can break the property when both children look correct.
3. Draw two different BST shapes that store the keys 1, 2, 3.

#### Advanced practical tasks

1. Write a recursive checker that carries a low and high bound. Test a valid tree and the 5–3–7–6 counterexample.
2. Read one textbook that allows equal keys. Write the equal-key rule in four STE sentences.

---

## Search, insert, delete

**Search** starts at the root. Compare the target with the current key. Go left if the target is smaller. Go right if the target is larger. Stop when you find the key or you reach null. Time is Θ(h).

**Insert** searches for the key. If the key is present, do not insert a duplicate (in this handbook). If you reach null, attach a new node there. Time is Θ(h).

```text
insert 5 into the 8–3–10 tree
search: 8 → 3 → 6 → 4, then 5 is right of 4
```

**Delete** has three cases when you found node u:

1. **No child.** Remove u. The parent pointer becomes null.
2. **One child.** Replace u with that child.
3. **Two children.** Find the **in-order successor** (the smallest key in the right subtree). Copy that key into u. Delete the successor node. The successor has no left child, so its delete is case 1 or case 2.

Some implementations use the in-order predecessor (the largest key in the left subtree). The idea is the same.

```text
delete 3 in the 8–3–10 tree
3 has two children
successor is 4
copy 4 into the node, then remove the old 4
```

Do not forget to update the parent link. A parent pointer, or a recursive return of a new child, both work.

Search, insert, and delete do not rebalance in a plain BST. Height can grow. The next section on degenerate trees explains the cost.

Find-min walks left until left is null. Find-max walks right. Both are Θ(h).

### Questions

#### Theoretical questions

1. How does search decide to go left or right?
2. Where does insert attach a new node?
3. What are the three delete cases?
4. What is the in-order successor?
5. Why is the successor easy to delete?

#### Easy practical tasks

1. Search for 7 and for 9 in the 8–3–10 picture. Write the nodes that you visit.
2. Draw the tree after insert of 5.
3. Make a table: "Delete case", "What you do". Add three rows.
4. Write five sentences about search, insert, and delete.

#### Medium practical tasks

1. Write the steps to find the in-order successor of 8 in the picture.
2. Explain in six sentences why you copy the successor key instead of moving many pointers.
3. Delete the root of a two-child tree on paper. Show the tree after the successor copy.

#### Advanced practical tasks

1. Implement search, insert, and the three delete cases in a language that you know. Test each delete case.
2. Write tests: empty tree, one node, insert duplicate (reject), delete missing key.

---

## Tree walks: in-order, pre-order, post-order, level-order

A **walk** visits each node one time. The four common walks are:

**In-order:** walk left, visit, walk right. On a BST, keys appear in sorted order.

**Pre-order:** visit, walk left, walk right. The root of each subtree appears first. Pre-order can copy the shape.

**Post-order:** walk left, walk right, visit. The root of each subtree appears last. Post-order can delete children before the parent.

**Level-order:** visit nodes by depth: the root, then depth 1, then depth 2. Level-order uses a **queue** (Topic 4). It is BFS on the tree.

```text
        8
       / \
      3    10
     / \
    1   6

in-order:    1, 3, 6, 8, 10
pre-order:   8, 3, 1, 6, 10
post-order:  1, 6, 3, 10, 8
level-order: 8, 3, 10, 1, 6
```

The first three walks use recursion or an explicit stack. The base case is null.

Time of each walk is Θ(n). Extra space is Θ(h) for recursion. Level-order extra space is Θ(w) where w is the maximum width of a level.

Do not mix the names. In-order is the sort walk on a BST. Pre-order is not sorted unless the tree is a line that happens to sort.

Level-order is not a DFS walk. It uses FIFO.

You can print values, collect keys, or evaluate an expression tree. The walk name tells the visit time, not the job.

### Questions

#### Theoretical questions

1. What is the visit order of in-order?
2. Why is in-order sorted on a BST?
3. When do you use post-order to delete a tree?
4. Which ADT does level-order use?
5. What is the time class of a full walk?

#### Easy practical tasks

1. Write all four walks for the picture in this section.
2. Make a table: "Walk", "Visit position", "Structure if iterative". Add four rows.
3. Draw a one-node tree. Write the four walks (they are the same).
4. Write five sentences that name the four walks.

#### Medium practical tasks

1. Write recursive in-order, pre-order, and post-order as three short procedures. Share the null base.
2. Explain in six sentences why level-order needs a queue.
3. Reconstruct the tree shape from the pre-order list and the in-order list of the picture. Write the idea, not a full proof.

#### Advanced practical tasks

1. Implement the four walks. Compare outputs with the picture.
2. Implement iterative in-order with an explicit stack. Write how you go left, then visit, then go right.

---

## Degenerate trees and why balance matters

Insert of sorted keys into a plain BST builds a **line**. Each new key is larger than the last. Each insert goes right. The result is a linked list that uses right pointers.

```text
insert 1, then 2, then 3, then 4

1
 \
  2
   \
    3
     \
      4
```

Height is n − 1. Search is Θ(n). The BST property still holds. The shape is **degenerate**.

A reverse-sorted sequence builds a line to the left. Many other sequences build a tall, thin tree. The expected height for random insert order is O(log n). Real keys are often sorted. Do not rely on luck.

**Balance** keeps height O(log n). Then search, insert, and delete stay logarithmic.

Topic 9 presents AVL trees, red-black trees, and B-trees. Those structures repair the shape after insert and delete.

A complete tree of n nodes has height about log2 n. That height is the goal.

Do not blame the BST property for a slow find. The property is correct. The height is the problem.

A periodic rebuild (sort the keys, build a balanced tree) can fix a degenerate tree. Rebuild is Θ(n). Balanced trees spread that work across operations.

### Questions

#### Theoretical questions

1. What insert order builds a right line?
2. What is the height of that line?
3. Why can the BST property hold when find is linear?
4. What does balance guarantee about height?
5. Why are real key streams often bad for a plain BST?

#### Easy practical tasks

1. Draw the tree after insert of 5, 4, 3. Write the height.
2. Draw a balanced tree of the same three keys. Write the height.
3. Make a table: "Shape", "Height", "Find class". Rows: line of n, perfect tree of n.
4. Write five sentences about degenerate trees.

#### Medium practical tasks

1. Insert 1 .. 8 in order on paper. Write how many compares search(8) needs.
2. Explain in six sentences why random order is not a reliable fix in a program.
3. Write the steps of "collect keys in-order, build a balanced tree" as a rebuild idea.

#### Advanced practical tasks

1. Time n sorted inserts into a plain BST versus n inserts in random order. Write the two times and the heights if you can measure them.
2. Write a one-page note: why Topic 9 exists. Name three balanced structures without their full rules.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do the BST property, height, and walk names work together in an ordered map?
2. Why does two-child delete use a successor that has at most one child?
3. When do you pick in-order versus level-order for a job?
4. How does a degenerate BST compare to a linked list (Topic 3)?
5. What stays true about a BST after a bad insert order, and what becomes false about cost?

#### Easy practical tasks

1. Write a one-page cheat sheet: BST property, search, insert, three deletes, four walks, degenerate, height.
2. Draw one BST. Write in-order keys and the path of search for one missing key.
3. Make a table: "Operation", "Class in terms of h". Add search, insert, delete, full walk.
4. List four defects: equal-key surprise, wrong successor, mix of heap and BST, sorted insert line.

#### Medium practical tasks

1. Implement a BST map: get, put, delete, in-order print. Write tests for the three delete cases.
2. Write a short report that maps each walk to one use: sort keys, copy shape, free nodes, print by layer.
3. Build the same keys as a line and as a more balanced tree. Count compares for find of the last key.

#### Advanced practical tasks

1. Implement successor and predecessor. Use them to write a range print [L, R] without a full scan of unrelated subtrees if you can.
2. Read one standard library ordered map. Write whether it uses a BST variant and what height bound the documentation suggests.
