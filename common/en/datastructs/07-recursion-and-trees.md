# 7. Recursion and Trees

## Description

Recursion is a function that calls itself. A tree is a linked structure with a root and no cycle. This topic explains base case and recursive case, tree words (root, leaf, parent, child, depth, height), binary trees, and the definitions of full, complete, perfect, and balanced.

Complete this topic after stacks. You will need the call-stack idea. Complete this topic before binary search trees.

Use one term for each concept. The **base case** is the input that does not recurse. The **recursive case** is the input that calls the same function on a smaller input. A **tree** is a set of nodes with parent–child links and exactly one root. A **binary tree** has at most two children per node. **Height** and **depth** are not the same number. This topic defines both.

---

## Base case and recursive case

A recursive function has two parts.

The **base case** stops the calls. It returns a value or it does nothing. Without a base case, the calls do not end. The call stack grows until the program fails.

The **recursive case** calls the same function on a smaller problem. The input must get closer to a base case. A usual measure is n, or the remaining list, or the remaining subtree.

```text
factorial(n):
  base:      if n = 0, return 1
  recursive: return n * factorial(n - 1)
```

```text
list length:
  base:      if head is null, return 0
  recursive: return 1 + length(head.next)
```

Each call has its own frame. The frame holds n or the current node. Topic 4 explained the call stack.

**Trust the recursive call.** You write the base case. You write how to combine the result of a smaller call. You do not replay the full chain in your head for large n.

A **wrong measure** fails to shrink. Example: factorial(n) that calls factorial(n) again. That call does not approach 0.

A **missing base** fails to stop. Example: a tree walk that always visits left and right and never checks for null.

Recursion is a method to write an algorithm. It is not a data structure. Trees fit recursion because each subtree is a smaller tree.

An iterative loop can replace some recursion. Use recursion when the structure is a tree or when the smaller-problem form is clear. Use an explicit stack when the call stack is too small.

### Questions

#### Theoretical questions

1. What is a base case?
2. What is a recursive case?
3. Why must the input shrink?
4. What happens if the base case is missing?
5. How does a tree fit recursion?

#### Easy practical tasks

1. Write the base case and the recursive case of factorial in two lines.
2. Write the base case and the recursive case of list length.
3. Make a table: "Function", "What shrinks". Rows: factorial, list length, tree walk.
4. Write five sentences that define base case and recursive case.

#### Medium practical tasks

1. Draw the call stack for factorial(3) before any return. Then write the return values from the bottom.
2. Explain in six sentences what "trust the recursive call" means on a left subtree plus a right subtree.
3. Write a recursive function that fails (does not shrink). Then write the one-line fix.

#### Advanced practical tasks

1. Write a recursive sum of an array using a low and high index. Then write an iterative sum. Compare extra stack space.
2. Convert a recursive list reverse idea to an iterative three-pointer reverse (Topic 3). Write which version uses the call stack.

---

## Root, leaf, parent, child, depth, height

A **tree** is a connected set of nodes with no cycle. One node is the **root**. Every other node has exactly one **parent**.

A **child** of node u is a node that has u as parent. A node can have zero or more children.

A **leaf** is a node with no children.

An **edge** in a tree is a parent–child link.

```text
        A          ← root
       / \
      B   C
     / \
    D   E          ← D and E and C are leaves
```

**Depth** of a node is the number of edges from the root to that node. The root has depth 0.

**Height** of a node is the number of edges on the longest path from that node down to a leaf. A leaf has height 0.

**Height of the tree** is the height of the root. In the picture, height(A) = 2. Depth(D) = 2.

Some texts count nodes on the path instead of edges. Then the root has depth 1. This handbook counts **edges**. Write the convention when you read another book.

A **path** in a tree is a sequence of distinct nodes that follow edges. There is exactly one path between any two nodes.

A **forest** is a set of trees. If you remove the root, the children become roots of a forest.

Do not mix depth and height. Depth is measured from the root. Height is measured toward the leaves.

### Questions

#### Theoretical questions

1. What is the root?
2. What is a leaf?
3. What is the depth of a node?
4. What is the height of a node?
5. What convention does this handbook use for depth of the root?

#### Easy practical tasks

1. Draw the A–B–C–D–E tree. Write depth of each node.
2. Write height of each node in the same tree.
3. Make a table: "Word", "Meaning". Add root, leaf, parent, child, depth, height.
4. Write five sentences that define the six words.

#### Medium practical tasks

1. Add a child F under C. Write the new height of the tree and the new depth of F.
2. Explain in six sentences why there is one path from the root to each node.
3. Write depth and height for a single-node tree.

#### Advanced practical tasks

1. Write recursive functions depth(node) from the root and height(node) from the node. State the base case for null.
2. Read one textbook definition of height. If it counts nodes, rewrite it in edge count in four sentences.

---

## Binary tree

A **binary tree** is a tree where each node has at most two children. The children have names: **left** and **right**. Left and right are different. A node with only a right child is not the same as a node with only a left child.

```text
    10
   /  \
  4    15
        \
         18
```

A node holds a value, a left pointer, and a right pointer. A null pointer means no child.

A **binary tree walk** uses recursion:

```text
walk(node):
  if node is null: return
  walk(node.left)
  visit node
  walk(node.right)
```

The position of "visit" changes the order. Topic 8 names in-order, pre-order, and post-order. This section only needs the null base case.

The maximum number of nodes at depth d is 2^d. Level 0 holds 1 node. Level 1 holds at most 2 nodes. Level 2 holds at most 4 nodes.

A binary tree is not a binary search tree unless it also keeps a key order. Topic 8 adds that property.

A node with three children is not a binary-tree node. A general tree can have many children. You can store children in a list.

Empty tree means the root pointer is null. A one-node tree has a root that is also a leaf.

### Questions

#### Theoretical questions

1. How many children can a binary-tree node have?
2. Why are left and right different?
3. What is the base case of a binary-tree walk?
4. What is the maximum number of nodes at depth d?
5. Why is a binary tree not automatically a search tree?

#### Easy practical tasks

1. Draw a binary tree of four nodes that is not complete. Label left and right.
2. Draw a node that has only a right child. Write that left is null.
3. Make a table: "Depth", "Max nodes". Rows: 0, 1, 2, 3.
4. Write five sentences that define a binary tree.

#### Medium practical tasks

1. Write a recursive count of nodes. Base: null returns 0. Recursive: 1 + left count + right count.
2. Explain in six sentences why a general tree uses a list of children.
3. Count leaves in the 10–4–15–18 picture.

#### Advanced practical tasks

1. Implement a binary-tree node and a recursive height function. Test the picture.
2. Write a function that returns whether a node has two children, one child, or no child. Use it to count each kind.

---

## Full, complete, perfect, and balanced (definitions)

These words describe shape. They do not describe key order.

A **full** binary tree: every node has 0 children or 2 children. No node has exactly one child.

A **complete** binary tree: every level is full except possibly the last. The last level fills from the left. There is no hole to the left of a node on the last level.

A **perfect** binary tree: every level is full. Every leaf has the same depth. A perfect tree is full and complete.

A **balanced** binary tree (this handbook): the height is O(log n). A common exact rule is the **AVL rule**: for every node, the height of the left subtree and the height of the right subtree differ by at most 1. Topic 9 uses that rule. Other books use other balance rules. Write the rule that you use.

```text
full, not complete:
      A
     / \
    B   C
       / \
      D   E

complete, not perfect:
      A
     / \
    B   C
   /
  D

perfect:
      A
     / \
    B   C
   / \ / \
  D  E F  G
```

A perfect tree of height h has 2^(h+1) − 1 nodes. Height counts edges.

Heaps (Topic 10) need a complete shape so that the tree fits in an array without holes.

Search trees need balance so that height stays logarithmic. A long line of nodes has height n − 1. Find then becomes linear.

Do not say "balanced" when you mean "perfect". A complete tree of n nodes is almost balanced, but check the definition that you use.

Do not say "full" when you mean "perfect". Full only bans single-child nodes.

### Questions

#### Theoretical questions

1. What is a full binary tree?
2. What is a complete binary tree?
3. What is a perfect binary tree?
4. What balance rule does this handbook name for AVL?
5. Why does a line of nodes have linear height?

#### Easy practical tasks

1. Label each of the three pictures in this section with the words that apply.
2. Make a table: "Word", "One-sentence meaning". Add four rows.
3. Draw a tree that is full and not complete. Use four or five nodes.
4. Write five sentences that separate the four words.

#### Medium practical tasks

1. For a perfect tree of height 2, write the number of nodes and the number of leaves.
2. Explain in six sentences why a heap wants a complete shape.
3. Draw a complete tree of 6 nodes. Mark the hole that you must not leave on the last level.

#### Advanced practical tasks

1. Write a checker for full, complete, and perfect on a small binary tree. Test all three pictures.
2. Read one AVL definition. Write the height-difference rule in four STE sentences. Do not write rotation code yet.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do base case, null child, and leaf work together in one recursive tree walk?
2. Why must you keep depth and height as two different numbers?
3. How does binary-tree shape (full, complete, perfect, balanced) differ from key order in Topic 8?
4. When do you replace recursion with an explicit stack (Topic 4)?
5. Why is a perfect tree always complete, and why is a complete tree not always perfect?

#### Easy practical tasks

1. Write a one-page cheat sheet: base case, recursive case, root, leaf, parent, child, depth, height, binary, full, complete, perfect, balanced.
2. Draw one tree. Label depth of each node, height of the root, and whether the tree is complete.
3. Make a table: "Idea", "Stops when". Rows: factorial, list length, tree walk.
4. List four defects: missing base, no shrink, mix of depth and height, one-child node in a "full" tree.

#### Medium practical tasks

1. Write recursive functions: node count, leaf count, and height. Share the same null base case.
2. Convert a recursive preorder visit to an explicit stack of nodes. Write the push order of left and right.
3. Build a complete tree of 5 nodes and a full tree of 5 nodes if you can. If you cannot, write why.

#### Advanced practical tasks

1. Measure call-stack depth for a recursive walk of a line of n nodes versus a perfect tree of n nodes. Write the two depths.
2. Write a one-page note that maps each shape word to one later topic (heap, AVL, B-tree) without explaining those topics in full.
