# 13. Strings and Tries

## Description

A string is a sequence of characters. This topic explains naive search, the idea of Knuth–Morris–Pratt (KMP), Rabin–Karp, tries (prefix trees), and suffix arrays at a high level.

Complete this topic after arrays, hash tables, and trees.

Use one term for each concept. The **text** T has length n. The **pattern** P has length m. A **match** is an index i where T[i .. i+m−1] equals P. A **trie** is a tree of prefixes. A **suffix array** is a sorted list of suffix start indexes.

---

## Naive search, KMP (idea), Rabin–Karp

**Naive search** tries P at every start i from 0 to n − m. At each i, it compares characters until a mismatch or a full match.

```text
T = A B A B C A B A B
P = A B A B
starts: 0 (match), 1 (fail), 2 (match), ...
```

Worst-case time is Θ((n − m + 1) m). A repeated letter can force almost m compares at almost every i.

**KMP (idea).** KMP never moves the text pointer backward. After a mismatch, it uses a **prefix table** (failure function) π. π[j] is the length of the longest proper prefix of P[0 .. j] that is also a suffix of that substring.

When T[i] ≠ P[j], KMP sets j to π[j − 1] (with a boundary at 0) and tries again. It does not restart at i = i+1 and j = 0 if a prefix of P already matches a suffix of the text that you read.

Time is Θ(n + m) after you build π in Θ(m).

You do not need a full KMP implementation in the first pass. You need the idea: preprocess P, reuse the matched prefix, scan T once.

**Rabin–Karp** hashes the pattern and hashes each window of length m in T. A **rolling hash** updates the window hash in Θ(1) when the window moves by one character.

```text
window hash at i+1 from hash at i
  drop T[i]
  add T[i+m]
  use a modulus and a base
```

If the hashes differ, the window is not a match. If the hashes are equal, compare characters to confirm. A collision can make the hashes equal when the strings differ.

Average time is O(n + m) with a good hash. Worst case is still Θ(n m) if every window collides and you verify.

Rabin–Karp extends to multiple patterns with a set of pattern hashes. KMP is for one pattern (Aho–Corasick generalizes the automaton idea).

Do not skip the verify step after a hash hit unless the hash is defined to be unique (for example, a large fingerprint with an accepted risk).

### Questions

#### Theoretical questions

1. What is the worst-case class of naive search?
2. What does the KMP prefix table store?
3. Why does KMP not move backward in T?
4. What is a rolling hash?
5. Why must Rabin–Karp verify a hash hit?

#### Easy practical tasks

1. Run naive search of P = AB on T = AABAB. Write each start i and match or fail.
2. Make a table: "Method", "Preprocess", "Scan idea". Add naive, KMP, Rabin–Karp.
3. Write five sentences that define the three methods.
4. Write one collision example: two different windows, same small hash modulus.

#### Medium practical tasks

1. Build π for P = ABAB by hand if you can, or write the meaning of each π[j] in words.
2. Explain in six sentences how a rolling hash drops one character and adds the next.
3. Compare naive and KMP on a text of many A's and a pattern of many A's plus one B.

#### Advanced practical tasks

1. Implement naive search and a rolling-hash Rabin–Karp with verify. Test a collision modulus if you pick a tiny modulus.
2. Read a KMP π construction note. Write ten STE sentences. Then implement π and the scan.

---

## Trie / prefix tree

A **trie** (prefix tree) stores a set of strings. Each edge is labeled with a character. A path from the root is a prefix. A node can mark **end of word**.

```text
words: can, cat, cot

        *
       /
      c
     / \
    a   o
   / \   \
  n   t   t
  $   $   $
```

**Insert** walks from the root. It creates a missing child for each next character. It marks the last node as an end.

**Find** walks the same path. If a child is missing, the word is absent. If the walk ends without an end mark, the string is only a prefix of another word (or a prefix that you never inserted as a word).

**Prefix search** walks the prefix, then walks the subtree to list all words with that prefix. Autocomplete uses this walk.

Time of insert or find is Θ(L) where L is the length of the string. It does not depend on the number of other words, except through memory and branching.

Space can be large: one node per character of all unique prefixes. An alphabet of size σ gives up to σ child pointers per node. A map from character to child saves space when the alphabet is large.

A trie is not a BST of full strings. Compare is per character on an edge, not a full-string compare at a node (unless you store a compressed edge).

A **compressed trie** (radix / Patricia) merges a unary path into one edge that holds a string. This section only requires the uncompressed idea.

Do not store the same word twice as two ends unless you want a multiset of words.

Delete must not remove nodes that other words still need. Delete is harder than insert. You can mark end = false and leave nodes, or you can prune unused tails.

### Questions

#### Theoretical questions

1. What does one path from the root represent?
2. Why does a node need an end mark?
3. What is the time of find in terms of L?
4. How does prefix search use a subtree?
5. Why can a large alphabet waste space?

#### Easy practical tasks

1. Draw a trie for at, ate, tea.
2. Mark end nodes. Write whether "te" is a word in that trie.
3. Make a table: "Operation", "Walk". Rows: insert, find, prefix list.
4. Write five sentences that define a trie.

#### Medium practical tasks

1. Insert "cat" and "car" on paper. Write which nodes they share.
2. Explain in six sentences why find("ca") can fail when "cat" is present.
3. Write how you list all words after a prefix node (a walk).

#### Advanced practical tasks

1. Implement a trie with insert, find, and startswith. Test shared prefixes.
2. Replace an array of 26 children with a small map. Write the space change for a sparse alphabet.

---

## Suffix array (high-level)

A **suffix** of T is T[i .. n−1] for a start index i.

A **suffix array** SA is an array of the n start indexes, sorted by the suffix string at that index.

```text
T = banana
indexes: 0:banana  1:anana  2:nana  3:ana  4:na  5:a

sorted suffixes:
  a        → 5
  ana      → 3
  anana    → 1
  banana   → 0
  na       → 4
  nana     → 2

SA = [5, 3, 1, 0, 4, 2]
```

You do not store n full suffix strings. You store n integers. A compare of two suffixes reads T from those starts.

**Find a pattern.** Binary search on SA. Each mid compare is a string compare of P with the suffix T[SA[mid] ..]. Time is O(m log n) compares if each compare is O(m). Extra tables (LCP) can speed this. This topic only needs the idea: sorted suffixes, binary search.

**Longest common prefix** of two suffixes is useful for repeated substrings. An **LCP array** stores LCP of SA[k] and SA[k−1]. Construction algorithms exist. Treat them as later reading.

A **suffix tree** is a compressed trie of all suffixes. It answers the same questions with more memory and a more complex build. The suffix array is a compact cousin.

Naive build: create n indexes, sort with a full string compare. Time is O(n² log n) in a simple model. Linear and n log n construction algorithms exist. You do not need to implement a linear builder in this topic.

Do not confuse a suffix array with a prefix trie of a dictionary. The suffix array indexes one text. The trie stores many words.

### Questions

#### Theoretical questions

1. What is a suffix of T?
2. What does the suffix array store?
3. Why do you not store n full strings?
4. How do you search a pattern in SA at a high level?
5. How does a suffix array differ from a dictionary trie?

#### Easy practical tasks

1. Write all suffixes of T = aba. Then write a sorted SA.
2. Make a table: "Structure", "Indexes", "Typical use". Rows: trie, suffix array.
3. Write five sentences that define a suffix array at a high level.
4. From the banana example, write the suffix at SA[0].

#### Medium practical tasks

1. Binary-search the pattern "ana" on the banana SA on paper. Write the compares.
2. Explain in six sentences why a naive sort of suffixes can be quadratic per compare.
3. Write one use of LCP in one sentence (repeated substring idea).

#### Advanced practical tasks

1. Build a suffix array with naive sort for a short T. Search two patterns, one present and one absent.
2. Read a short note on suffix trees versus suffix arrays. Write ten STE sentences. Do not copy construction code.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. When do you pick naive search, KMP, or Rabin–Karp for one pattern in one text?
2. How do a trie and a suffix array both use order of characters, but on different inputs?
3. Why does an end mark matter in a trie but not in a suffix-array binary search the same way?
4. How does hashing in Rabin–Karp relate to Topic 5 (collisions and verify)?
5. What does "high-level" mean for the suffix array in this topic?

#### Easy practical tasks

1. Write a one-page cheat sheet: naive, KMP π, rolling hash, trie, end mark, prefix, SA, suffix.
2. Draw a trie of three words and the SA of one short text on the same page.
3. Make a table: "Need", "Structure or algorithm". Add find substring, autocomplete, all suffixes of one T.
4. List four defects: skip Rabin–Karp verify, no end mark, naive O(n m) on huge T, mix trie with SA.

#### Medium practical tasks

1. Implement a trie autocomplete and a naive substring search. Write when each is the right tool.
2. Write a test plan: overlapping matches (AAAA in AAAAA), empty pattern policy, missing prefix.
3. Design an ADT: insert word, find word, find by prefix. State that SA is not this ADT.

#### Advanced practical tasks

1. Implement Rabin–Karp and naive search. Compare times on a long random text and on a repetitive text.
2. Read KMP or a suffix-array construction survey. Map five terms to this topic in a table.
