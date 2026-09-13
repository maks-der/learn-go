# 17. Strings as Data

## Description

A string is a sequence of characters. This topic explains naive search, Knuth–Morris–Pratt, Rabin–Karp, tries, suffix arrays, and suffix trees. Complete this topic after arrays and hash tables. Complete trees before the trie section.

Use one term for each concept. The **text** T has length n. The **pattern** P has length m. A **match** is an index i where T[i .. i+m−1] equals P. A **prefix** is a leading segment. A **suffix** is a trailing segment.

---

## Naive search

**Naive search** tries the pattern at every start index i in 0 .. n−m.

At each i, compare P[0] with T[i], P[1] with T[i+1], ... until a mismatch or a full match.

```text
T = A B A B A C
P = A B A

i=0: ABA vs ABA  match
i=1: BAB vs ABA  mismatch at first letter
i=2: ABA vs ABA  match
i=3: BAC vs ABA  mismatch
```

```text
naive(T, P):
  n = length(T)
  m = length(P)
  for i = 0 to n - m:
      if T[i .. i+m-1] == P:
          report i
```

Worst-case time is Θ((n − m + 1) m). Example: T is n times A, P is (m−1) times A plus one B. Each start almost compares m characters.

Average time on random letters is closer to Θ(n). Many starts fail on the first letter.

Naive search is easy to write and easy to check. Implement it first. Use it as a test oracle for KMP and Rabin–Karp.

Do not use naive search as the only method on huge texts and long patterns if you can use a linear method.

Indexes are 0-based in this handbook. A match at the end starts at n − m.

An empty pattern is a special case. Define it: no match, or a match at every index. Write the rule in your tests.

### Questions

#### Theoretical questions

1. What start indexes does naive search try?
2. When do you report a match?
3. What is the worst-case time class?
4. Why is the all-A text a slow case?
5. Why do you implement naive search first?

#### Easy practical tasks

1. Run naive search by hand on T = ABABAC and P = ABA. Write the i values that match.
2. Make a table: "i", "Compare count", "Result" for that example.
3. Write five sentences that define naive search.
4. For n = 10, m = 3, write the number of start indexes.

#### Medium practical tasks

1. Implement naive search. Return all match indexes. Test the example and a no-match case.
2. Build the all-A worst case for small n and m. Count character compares.
3. Explain in six sentences how you treat an empty pattern.

#### Advanced practical tasks

1. Time naive search on a 1 000 000-character random string and m = 50. Write the time.
2. Use naive search as an oracle: later algorithms must return the same index list on 20 random tests.

---

## Knuth–Morris–Pratt (idea + implement)

**Knuth–Morris–Pratt (KMP)** searches in Θ(n + m) time. It does not restart the pattern from P[0] after every mismatch. It uses a **prefix table** (also called the failure function or lps: longest proper prefix that is also a suffix).

**Idea.** You compared T[i .. i+q−1] with P[0 .. q−1] and then T[i+q] != P[q]. Some tail of the matched piece can still be a prefix of P. The table says how far you can slide P so that that prefix stays aligned. You do not move i back on the text. You only decrease the matched length q.

**lps[k]** is the longest length j < k+1 such that P[0 .. j−1] equals P[k-j+1 .. k]. lps[0] = 0. (Exact index conventions vary. Pick one textbook and stay with it.)

```text
P = A B A B A C
lps idea: after a mismatch, slide to the next possible prefix
```

```text
kmp(T, P):
  lps = buildLps(P)
  q = 0
  for i in 0 .. n-1:
      while q > 0 and P[q] != T[i]:
          q = lps[q - 1]
      if P[q] == T[i]:
          q = q + 1
      if q == m:
          report i - m + 1
          q = lps[q - 1]
```

**buildLps** is KMP of the pattern against itself. Time is Θ(m).

Implement KMP after you can draw one slide on paper. Implement buildLps first. Print lps for ABA, AAAA, and ABC.

KMP does not hash. It only compares characters and uses integers in lps.

Worst-case time is Θ(n + m). Extra space is Θ(m) for lps.

Do not copy a table from a website without checking it against naive search.

### Questions

#### Theoretical questions

1. What does the prefix table store?
2. Why does KMP not move the text index backward?
3. What is the time class of KMP?
4. How do you build lps?
5. Why must you test KMP against naive search?

#### Easy practical tasks

1. Draw P = ABABAC. Write one mismatch case and the slide that a prefix allows.
2. Make a table: "Algorithm", "Worst-case time". Rows: naive, KMP.
3. Write five sentences that define the KMP idea without code.
4. For P = AAAA, explain why lps values grow.

#### Medium practical tasks

1. Implement buildLps. Print lps for ABA, ABABAC, and ABC.
2. Implement full KMP. Match naive search on the topic example and on AAAA / AAA.
3. Explain in six sentences the while loop that decreases q.

#### Advanced practical tasks

1. Implement KMP and time it versus naive on the all-A worst case. Write both times.
2. Read one standard lps definition (CLRS or a course). Map their indexes to yours in ten STE sentences. Cite the source.

---

## Rabin–Karp (rolling hash)

**Rabin–Karp** compares a **rolling hash** of the current window of length m with the hash of P. If the hashes differ, the window is not a match. If the hashes are equal, you compare the characters (or you accept a small collision risk in some variants).

A **rolling hash** updates the window hash in O(1) when the window moves by one character. You drop the left character and add the right character. You do not rescan m characters for the hash.

```text
hash of T[i .. i+m-1]
next = (hash - T[i] * base^{m-1}) * base + T[i+m]
modulo a prime
```

```text
rk(T, P):
  hp = hash(P)
  hw = hash(T[0 .. m-1])
  if hw == hp and T[0 .. m-1] == P: report 0
  for i = 1 to n - m:
      hw = roll(hw, T[i-1], T[i+m-1])
      if hw == hp and T[i .. i+m-1] == P: report i
```

Use a large prime modulus. Use a base such as 256 for bytes. Mind overflow. Use unsigned integers or explicit modular arithmetic.

Worst-case time is Θ(n m) if every window collides and you verify each one. Average time is Θ(n + m) with a good hash.

**Double hash** (two moduli) reduces collision risk.

Rabin–Karp is useful for multiple patterns of the same length and for the idea of rolling hash in other problems.

Do not skip the character verify if you need a correct match list. A hash collision is rare but real.

KMP is usually the first linear string match that you trust. Rabin–Karp is the hash lesson.

### Questions

#### Theoretical questions

1. What two hashes do you compare?
2. What does a roll step remove and add?
3. Why can equal hashes still need a character check?
4. What is the average time class with a good hash?
5. When is Rabin–Karp a better teaching tool than a faster matcher?

#### Easy practical tasks

1. Compute a tiny hash: base 10, modulus 97, pattern "12". Show one roll on text "123".
2. Make a table: "Event", "What you do". Rows: hash miss, hash hit.
3. Write five sentences that define a rolling hash.
4. Write the window start indexes for n = 8, m = 3.

#### Medium practical tasks

1. Implement rolling hash with verify. Match naive search on several strings.
2. Force a collision with a tiny modulus (for example 11). Show a false hash hit that verify rejects.
3. Explain in six sentences why you need base^{m-1} in the roll.

#### Advanced practical tasks

1. Implement two moduli. Compare collision counts with one tiny modulus on random texts.
2. Time Rabin–Karp versus KMP versus naive on one long text. Write the three times.

---

## Trie / prefix tree

A **trie** (prefix tree) stores a set of strings. Each edge is labeled with a character. A path from the root is a prefix. A node can mark "end of word".

```text
words: at, ate, be

      *
     / \
    a   b
    |   |
    t*  e*
    |
    e*
```

**Insert** walks or creates one child per character, then marks the last node as an end.

**Find** walks the characters. If a child is missing, the word is absent. If you finish the word and the end mark is set, the word is present.

**Prefix query.** Walk the prefix. If the walk succeeds, some stored word has that prefix. You can then iterate the subtree for autocomplete.

Time of insert or find is Θ(L) for a string of length L. It does not depend on the number of words n in the same way as a linear scan of all words. Space can be large if many nodes have few children.

An array of 26 children fits lowercase English. A map from character to child fits a larger alphabet.

A trie is not a suffix tree. A trie stores whole words as paths from the root. It does not store all suffixes of one text unless you insert them (that would waste space).

Do not use a trie when you only need exact find of a few words. A hash set is simpler. Use a trie when prefixes matter.

### Questions

#### Theoretical questions

1. What does one edge in a trie represent?
2. How do you know that a path is a complete word?
3. How do you answer "any word starts with this prefix?"
4. What is the time of find for a string of length L?
5. When do you prefer a hash set to a trie?

#### Easy practical tasks

1. Draw a trie for cat, car, and dog. Mark ends.
2. Make a table: "Operation", "Walk". Rows: insert, find, prefix.
3. Write five sentences that define a trie.
4. List the prefixes of "ate" that a node path can represent.

#### Medium practical tasks

1. Implement insert and find for lowercase words. Test the at/ate/be picture.
2. Implement a prefix function that returns true if any word has the prefix. Test "at" and "ba".
3. Explain in six sentences how autocomplete walks the subtree.

#### Advanced practical tasks

1. Implement autocomplete that returns up to k completions for a prefix. Test three prefixes.
2. Compare memory of a trie versus a hash set of the same word list. Write a rough node count.

---

## Suffix array (high-level)

A **suffix array** of text T is an array of the start indexes of all suffixes of T, sorted in lexicographic order.

```text
T = banana
suffixes:
  0 banana
  1 anana
  2 nana
  3 ana
  4 na
  5 a

sorted suffixes: a, ana, anana, banana, na, nana
suffix array:    5, 3, 1, 0, 4, 2
```

A suffix array uses Θ(n) integers. It uses less memory than a suffix tree. Many queries that a suffix tree answers can use a suffix array plus extra tables (LCP). Those extras are a later topic.

**Find a pattern.** Binary search the suffix array. Compare P with the suffix at the mid index. Time is O(m log n) compares if you compare from the start each time. Extra LCP data can reduce that cost.

You can build a suffix array in O(n log n) with sort of suffixes and a careful key, or in O(n) with a complex linear construction. This handbook requires the definition and a slow build for small n: sort the n suffixes with the language sort.

Do not build a suffix array of a huge text with n full string sorts if n is large. That can cost Θ(n² log n) character work. Use a known algorithm when n is large.

A suffix array is not a trie. It is a permutation of indexes 0 .. n−1.

### Questions

#### Theoretical questions

1. What does each entry of a suffix array store?
2. What order do the suffixes have?
3. How do you search a pattern with binary search?
4. Why can a naive sort of all suffixes be too slow?
5. Why does a suffix array use less memory than a suffix tree?

#### Easy practical tasks

1. Write all suffixes of "ana". Sort them. Write the suffix array.
2. Make a table: "Index", "Suffix" for "dog".
3. Write five sentences that define a suffix array.
4. For T = banana, check that index 5 is the suffix "a".

#### Medium practical tasks

1. Build a suffix array for a string of length ≤ 20 with a language sort of pairs (suffix, index). Print the array.
2. Binary-search a pattern on that array. Test a match and a miss.
3. Explain in six sentences why two suffixes that share a prefix sit near each other in the array.

#### Advanced practical tasks

1. Read one O(n log n) construction outline (prefix-doubling). Write ten STE sentences. Do not implement the full algorithm unless you want a project.
2. Write a one-page note: suffix array versus trie of all suffixes (space).

---

## Suffix tree (awareness)

A **suffix tree** is a compressed trie of all suffixes of T. A compact edge stores a substring of T (two indexes) instead of one character per edge. After linear-time construction, many string queries are fast.

This section is **awareness only**. Do not implement Ukkonen's algorithm in this topic.

Facts to remember:

- The tree has O(n) nodes if edges are compressed.
- You can find P in O(m) time after you have the tree (alphabet issues aside).
- A suffix tree can find the longest common substring of two texts (with a combined text and a separator).
- Construction in linear time is famous and easy to get wrong.
- A suffix array plus LCP is a common replacement in practice.

```text
you implement now:  naive search, KMP, rolling hash, trie, slow suffix array
you only name:      suffix tree, Ukkonen
```

If a paper says "build a suffix tree", ask whether a suffix array is enough for the query.

A suffix tree is not a B+ tree. The word "tree" is the only shared idea.

Do not confuse a suffix tree with a trie of a dictionary. A dictionary trie stores chosen words. A suffix tree stores every suffix of one text.

### Questions

#### Theoretical questions

1. What is a suffix tree at a high level?
2. Why do compressed edges keep O(n) nodes?
3. Why is this section awareness only?
4. What practical structure often replaces a suffix tree?
5. How is a suffix tree different from a dictionary trie?

#### Easy practical tasks

1. Write five sentences that define a suffix tree without construction steps.
2. Make a table: "Structure", "You implement now?". Rows: trie, suffix array, suffix tree.
3. Copy the phrase "awareness only" into a note. Add one reason.
4. Name two queries that a suffix tree can support (from this section).

#### Medium practical tasks

1. Find one diagram of a suffix tree for "banana". Write the number of leaves (should be n, or n+1 with a terminator). Cite the page.
2. Explain in six sentences why a terminator symbol (such as $) is useful.
3. List three student errors: implementing Ukkonen too early, mixing trie with suffix tree, mixing B+ with suffix tree.

#### Advanced practical tasks

1. Write a one-page compare sheet: suffix array, suffix tree, dictionary trie. Cite two sources. Do not implement a suffix tree.
2. After you finish a suffix array search, write one paragraph on which query would still want a suffix tree.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe the path from naive compare to a linear scan (KMP) and to a hash window (Rabin–Karp).
2. How do a trie and a suffix array both use prefixes, but store different objects?
3. A teammate says "put the text in a hash map to find a pattern". Which facts do you use to correct that sentence?
4. What must stay the same across naive, KMP, and Rabin–Karp on a correct test?
5. Why does this handbook delay suffix trees but ask you to implement KMP?

#### Easy practical tasks

1. Write a one-page cheat sheet: naive, lps, rolling hash, trie, suffix array, suffix tree awareness.
2. Draw a trie of three words and a suffix array of one four-letter word.
3. Make a table: "Need", "Structure". Rows: exact word set, prefix complete, one-pattern scan, all suffixes.
4. Write T, P, n, and m for one example from memory.

#### Medium practical tasks

1. Implement naive, KMP, and Rabin–Karp. They must agree on 15 tests.
2. Implement a trie for autocomplete of a 100-word list.
3. Build a slow suffix array and search three patterns. Check against naive search.

#### Advanced practical tasks

1. Time the three scanners on one long text and one worst-case text. Write a small report.
2. Read the next topic (segment trees and tries). Write five sentences on how a radix tree will extend this trie.
