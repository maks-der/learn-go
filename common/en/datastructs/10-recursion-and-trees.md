# 10. Recursion and Trees Intro

## Description

Recursion and trees use the same idea: a problem has a small base and a smaller copy of the same problem. This topic explains base case, recursive case, the call tree, tree words, binary trees, and shape names. Complete this topic before binary search trees.

Use one term for each concept. A **base case** stops the recursion. A **recursive case** calls the same function on a smaller input. A **tree** is a connected structure of nodes without a cycle. The **root** is the top node. A **leaf** has no child.

---

## Base case and recursive case

A recursive function has two parts.

The **base case** handles the smallest input. The function returns without a further call. Example: the empty list. Example: n == 0 for factorial. Example: a null child pointer.

The **recursive case** solves a smaller input, then combines the result. The smaller input must move toward the base case. If the input does not get smaller, the function never stops.

```text
factorial(n)
  base:      n == 0  →  1
  recursive: n > 0   →  n * factorial(n - 1)
```

```go
func Fact(n int) int {
	if n == 0 {
		return 1
	}
	return n * Fact(n-1)
}
```

```python
def fact(n):
    if n == 0:
        return 1
    return n * fact(n - 1)
```

A list walk can be recursive:

```text
sum(node)
  base:      node is null  →  0
  recursive: node.value + sum(node.next)
```

You must write the base case first in your mind. Then you write the recursive case. Test the empty input.

Recursion uses the call stack. Each call has a frame. Deep n can overflow the call stack. A loop is often safer for a linear scan. Recursion is a natural fit for a tree: each child is a smaller tree.

Do not use recursion when a simple loop is enough, unless you study recursion.

### Questions

#### Theoretical questions

1. What does the base case do?
2. What does the recursive case do?
3. Why must the input get smaller?
4. What happens if you omit the base case?
5. Why can a linear scan use a loop instead of recursion?

#### Easy practical tasks

1. Write the base case and the recursive case for factorial in two sentences.
2. Write the two cases for sum of a linked list.
3. Make a table: "Input", "Which case?". Rows: n = 0, n = 3 for factorial.
4. Draw the idea of "smaller input" for `fact(3)` → `fact(2)` → `fact(1)` → `fact(0)`.

#### Medium practical tasks

1. Implement factorial with recursion. Test 0, 1, and 5.
2. Implement a recursive length of a linked list. Test empty and three nodes.
3. Write a recursive function with a missing base case on purpose. Record the error. Then add the base case.

#### Advanced practical tasks

1. Write recursive binary search on a sorted array. State the two base cases: empty range and found mid.
2. Convert one recursive list walk to a loop. Write why the loop uses less call-stack space.

---

## Call tree

A **call tree** is a picture of recursive calls. The root is the first call. Each child is a call that the parent makes. A leaf call is a base case.

```text
fib(4)
  fib(3)              fib(2)
    fib(2)  fib(1)      fib(1)  fib(0)
      fib(1) fib(0)
```

The naive Fibonacci function makes two recursive calls. The call tree has many repeated nodes. The number of calls grows like Θ(φⁿ). That growth is exponential. A loop or memoization avoids the repeat work.

A tree walk that calls left and right also has a call tree. If each node is visited a constant number of times, the total work is Θ(n) for n nodes.

Read a call tree from the top. The depth of the call tree is the depth of the call stack. A skinny deep tree can overflow the stack. A wide shallow tree uses less stack and more sibling work.

```text
depth of call tree ≈ maximum number of frames at one time
```

Draw a call tree when you debug a recursive function. Count the base-case leaves. Check that every path ends in a base case.

### Questions

#### Theoretical questions

1. What is a call tree?
2. What is a leaf in a call tree?
3. Why does naive Fibonacci repeat work?
4. How does call-tree depth relate to the call stack?
5. When is total work Θ(n) for a tree walk?

#### Easy practical tasks

1. Draw the call tree for `fact(3)`.
2. Draw the call tree for naive `fib(3)`.
3. Make a table: "Function", "Calls per parent". Rows: fact, naive fib, binary tree walk.
4. Write five sentences that define a call tree.

#### Medium practical tasks

1. Add a counter to naive Fibonacci. Print the number of calls for n = 6, 8, 10.
2. Draw the call tree of a recursive sum on a list of four nodes. Show that the tree is a line.
3. Explain in six sentences why a line-shaped call tree can overflow before a balanced tree walk of the same n.

#### Advanced practical tasks

1. Implement naive Fibonacci and a loop Fibonacci. Time both for a moderate n that still finishes.
2. Write a memoized Fibonacci (map from n to value). Compare the call count with the naive version.

---

## Tree terminology: root, leaf, parent, child, depth, height

A **tree** in this handbook is a set of nodes with parent and child links. There is one **root**. Every other node has exactly one parent. There is no cycle.

A **child** of node u is a node that has parent u. A **parent** of node v is the unique node above v. The root has no parent.

A **leaf** is a node with no child.

The **depth** of a node is the number of edges from the root to that node. The root has depth 0.

The **height** of a node is the number of edges on the longest path from that node down to a leaf. The height of a leaf is 0.

The **height of the tree** is the height of the root.

```text
        A          depth(A)=0, height(A)=2
       / \
      B   C        depth(B)=1, height(B)=1
     /
    D              depth(D)=2, height(D)=0  (leaf)
```

Some texts say the height of a leaf is 1, or they count nodes instead of edges. This handbook counts **edges**. Write the convention when you read a different book.

A **sibling** shares the same parent. An **ancestor** sits on the path from the node up to the root. A **descendant** sits in the subtree.

A **subtree** of node u is u plus all descendants of u.

An empty tree has no node. A null child is not a node. Many functions treat null as height −1 so that a leaf (two null children) has height 0.

### Questions

#### Theoretical questions

1. What is the root?
2. What is a leaf?
3. How does this handbook define depth?
4. How does this handbook define height of a tree?
5. What is a subtree?

#### Easy practical tasks

1. Draw a tree of five nodes. Label root and leaves.
2. Write depth and height for each node in the A–B–C–D picture.
3. Make a table: "Word", "Meaning". Add root, parent, child, leaf, depth, height.
4. Write one sentence that states the edge-count convention.

#### Medium practical tasks

1. Build a small tree in code with parent or child pointers. Print depth of each node with a walk.
2. Compute height of a tree with a recursive function. Use null height = −1.
3. Explain in six sentences why two books can disagree on height by 1.

#### Advanced practical tasks

1. Write functions: depth(node), height(node), is_leaf, is_root. Test on a known drawing.
2. Prove or argue in eight STE sentences that a tree with n nodes has n − 1 edges.

---

## Binary tree

A **binary tree** is a tree in which each node has at most two children. The children are **left** and **right**. Left and right are different. A node can have a left child and no right child.

```go
type BinNode struct {
	Value int
	Left  *BinNode
	Right *BinNode
}
```

```python
class BinNode:
    def __init__(self, value, left=None, right=None):
        self.value = value
        self.left = left
        self.right = right
```

A binary tree is not always a binary search tree. A binary search tree adds an order rule. That rule is topic 11.

A walk of a binary tree uses recursion in a natural way:

```text
walk(node)
  base: node is null → return
  walk(node.left)
  use node
  walk(node.right)
```

The order of the three lines gives in-order, pre-order, or post-order. Topic 11 defines those names in the BST setting. The same names apply to any binary tree.

A binary tree can represent an expression. The internal nodes are operators. The leaves are numbers. You can study that picture later.

The maximum number of nodes at depth d is 2ᵈ if every level is full. The maximum number of nodes in a binary tree of height h is 2^{h+1} − 1.

### Questions

#### Theoretical questions

1. How many children can a binary-tree node have?
2. Why are left and right different?
3. Is every binary tree a binary search tree?
4. What is the maximum number of nodes at depth d?
5. Why is recursion a natural walk of a binary tree?

#### Easy practical tasks

1. Draw a binary tree with a root and two children. Add one grandchild on the left.
2. Write a node type in Go or in Python.
3. Make a table: "Node", "Left", "Right" for a three-node tree.
4. Write the maximum node count for height 0, 1, and 2.

#### Medium practical tasks

1. Build a binary tree of four nodes in code. Walk and print values in any order that you document.
2. Count leaves with a recursive function.
3. Explain in six sentences the difference between a general tree (many children) and a binary tree.

#### Advanced practical tasks

1. Store a binary tree in an array as if it were complete (index i has children 2i+1 and 2i+2). Document missing nodes.
2. Write a recursive function that returns the node count and the height in one walk.

---

## Full / complete / perfect / balanced (definitions)

These words describe the **shape** of a binary tree. Learn the definitions. Different books disagree on "complete" and "full". This handbook uses the definitions below.

**Perfect.** Every internal node has two children. All leaves have the same depth. A perfect tree of height h has 2^{h+1} − 1 nodes.

**Full** (sometimes "proper"). Every node has 0 children or 2 children. No node has exactly one child.

**Complete.** Every level is full except possibly the last level. The last level is filled from the left. A binary heap uses a complete tree. You will see that layout in a later topic.

**Balanced** (practical). The height is small compared with n. A usual goal is height Θ(log n). One precise rule (AVL) is: for every node, the height of the left subtree and the height of the right subtree differ by at most 1. You do not implement AVL in this topic.

```text
perfect, height 1:     3 nodes (root and two leaves)
full, not perfect:     root with two children; one child has two leaves
                       (every node has 0 or 2 children; leaf depths differ)
complete, not perfect: last level filled from the left; one missing node on the right
```

A **degenerate** tree is a line. Each node has one child. Height is n − 1. A degenerate binary tree behaves like a linked list. Topic 11 uses this word again.

Do not say "balanced" when you have no height rule. Use a definition. For this topic, write height and n. If height is about log₂ n, the tree is shallow. If height is about n, the tree is degenerate.

### Questions

#### Theoretical questions

1. What is a perfect binary tree in this handbook?
2. What is a full binary tree in this handbook?
3. What is a complete binary tree in this handbook?
4. What practical meaning does balanced have here?
5. What is a degenerate tree?

#### Easy practical tasks

1. Draw one perfect tree of height 2.
2. Draw one complete tree that is not perfect.
3. Draw one degenerate tree of 4 nodes.
4. Make a table: "Shape", "Height versus n". Rows: perfect, degenerate.

#### Medium practical tasks

1. Count nodes in a perfect tree of height 3. Write 2^{3+1} − 1.
2. For each of four drawings (you make them), label perfect / full / complete / none. Follow this handbook.
3. Explain in six sentences why a heap later needs a complete shape.

#### Advanced practical tasks

1. Write predicates: is_perfect, is_full, is_complete. Test them on small trees. Document the definitions.
2. Find one textbook that uses a different meaning of "full" or "complete". Write the difference in four STE sentences. Cite the book.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do base case, recursive case, and call-tree depth work together?
2. Why is a tree walk a better use of recursion than a list sum?
3. Describe a binary-tree node and the words root, leaf, depth, and height in one short paragraph of STE sentences.
4. How do the shape names help you predict height?
5. A teammate says every recursive function is exponential. Which facts do you use to correct that claim?

#### Easy practical tasks

1. Write a one-page cheat sheet: base case, recursive case, call tree, root, leaf, parent, child, depth, height, binary tree, perfect, full, complete, balanced, degenerate.
2. Draw a call tree and a data tree side by side. Label which picture is which.
3. Write height of a perfect binary tree with 15 nodes.
4. Write one base case for a binary-tree function that takes a node pointer.

#### Medium practical tasks

1. Implement recursive height, node count, and leaf count. Test on a perfect tree of height 2 and on a line of 4 nodes.
2. Draw the call tree of height(root) on a small binary tree. Compare depth with tree height.
3. Write a loop that computes list length and a recursion that computes tree height. Comment why the tools differ.

#### Advanced practical tasks

1. Implement a safe recursive walk with an explicit depth limit. Return an error when depth exceeds the limit.
2. Write a one-page map from this topic to topic 11: which words you will reuse for BST search and for degenerate trees.
