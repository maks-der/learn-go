# 11. Indexes and Query Access

## Description

This topic shows how the DBMS finds rows. You learn why a full scan is slow, how a B-tree index works at a high level, unique indexes, composite indexes, covering indexes, selectivity, write cost, and how you read a query plan.

Use one term for each concept. An access path is the method that the optimizer chooses. A full table scan reads every row. An index seek uses an index to find a small set of rows. Complete this topic after you can write joins and `GROUP BY`. Indexes exist for those queries.

Do not add indexes before you have a query. Measure with a plan.

---

## Why a full table scan is slow

A full table scan reads every row in the table (more precisely: every page that holds the table). The DBMS then applies `WHERE`. If you need one customer by id from a million rows, a scan still reads the million rows.

Cost grows with table size. Twice the rows means about twice the work for a scan. Disk and memory bandwidth become the limit.

A scan is not always wrong. If you need most rows, a scan can be cheaper than many index lookups. Example: `SELECT * FROM products` with no filter. Example: a filter that matches 80 percent of the rows.

A scan is wrong when you need a few rows and an index exists on the filter column. Example: `WHERE customer_id = 42` on a large `orders` table.

```sql
SELECT *
FROM orders
WHERE customer_id = 42;
```

Without an index on `customer_id`, the DBMS reads all orders. With a useful index, it jumps to the matching entries.

`COUNT(*)` of a large table can also scan (or scan an index). Do not run heavy scans on production as a probe. Use a plan. Use a sample.

Do not guess. A table of 20 rows is small. A scan is fine. A table of 20 million rows is not small. The same SQL text has a different cost.

Joins can scan both sides if keys are not indexed. A nested loop join with a scan on the inner side is a common slow pattern. Index the join key (usually the foreign key).

### Questions

#### Theoretical questions

1. What does a full table scan read?
2. When is a scan a good access path?
3. When is a scan a poor access path?
4. Why does scan cost grow with table size?
5. Why do join keys need indexes on large tables?

#### Easy practical tasks

1. Write two queries: one that should scan, one that should use an index on a large table.
2. Explain in four sentences why `WHERE customer_id = 42` is expensive without an index.
3. Estimate: 1 µs per row versus 1 million rows. Write the scan time in words.
4. Label three of your practice queries as "scan is fine" because the table is small.

#### Medium practical tasks

1. On a table with thousands of rows (create them if needed), time `SELECT` with a selective `WHERE` before and after an index. Write the two times.
2. Run `SELECT *` with no filter. Write why an index does not help that query.
3. Draw a nested loop join: outer orders, inner customers by id, with and without an index on `customers.customer_id`.

#### Advanced practical tasks

1. Write a one-page note: when the optimizer prefers a scan even if an index exists (low selectivity).
2. Fill a scratch table with enough rows to see a plan change. Document row counts and plans.

---

## B-tree index idea

A B-tree index is a balanced tree of keys. Leaf nodes hold keys in order and pointers to rows (or include the row in a clustered index). Upper nodes guide the search.

To find `customer_id = 42`, the DBMS walks from the root to a leaf. The number of steps grows slowly with the number of keys (logarithmic). Then it follows pointers to the matching rows.

Because leaves are ordered, a B-tree also supports ranges:

```sql
WHERE customer_id = 42
WHERE ordered_at >= DATE '2026-01-01' AND ordered_at < DATE '2026-02-01'
WHERE name >= 'A' AND name < 'B'
```

A B-tree does not help every expression. `WHERE YEAR(ordered_at) = 2026` can hide the column from a simple index on `ordered_at`. Prefer a range on the column.

```text
Root
  ├─ ... 40
  └─ 41 ...
Leaf: 41, 42, 42, 43  → row pointers
```

Most relational products use a B-tree (or a B+ tree) as the default index. Other kinds exist (hash, bitmap, GIN). Learn B-tree first. This path stays vendor-neutral: think "ordered key lookup."

The primary key index is usually a B-tree. Secondary indexes are extra B-trees on other columns.

Do not picture the index as a copy of the full table. It stores the key and a pointer (or a subset of columns). Topic 8 said the index is a schema object. This section says how it finds rows.

The optimizer chooses whether to use the B-tree. You create the index. You do not pick the path in the `SELECT` text in standard SQL.

### Questions

#### Theoretical questions

1. What does a B-tree store in the leaves at a high level?
2. Why is a B-tree good for a range predicate?
3. Why can `YEAR(ordered_at) = 2026` miss an index on `ordered_at`?
4. What is the default index kind in most relational products?
5. Who chooses to use the index: you or the optimizer?

#### Easy practical tasks

1. Draw a tiny B-tree for keys `10, 20, 30, 40`.
2. Write three `WHERE` clauses that a B-tree on `ordered_at` can support.
3. Write one `WHERE` that is likely to hide that index.
4. Explain in four sentences the walk from root to leaf.

#### Medium practical tasks

1. Create a B-tree index (default `CREATE INDEX`) on a date column. Run a range query.
2. Compare equality and range on that column in the query plan if you can.
3. Read the product name for its default index type. Write it.

#### Advanced practical tasks

1. Write a one-page B-tree versus hash index note: equality only versus ranges.
2. Explain why a B-tree stays balanced after inserts (high-level: split nodes). No implementation code.

---

## Unique index vs primary key

A primary key is a logical rule: identify the row, not `NULL`, unique. The DBMS implements it with a unique index in almost every product.

A unique index is a physical object that forbids duplicate keys. You can create a unique index on a column that is not the primary key. That index supports a candidate key (`email`, `isbn`).

```sql
CREATE TABLE customers (
    customer_id INTEGER PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE
);
```

`PRIMARY KEY` creates a unique index on `customer_id`. `UNIQUE` creates a unique index on `email`.

You can also write:

```sql
CREATE UNIQUE INDEX idx_customers_email
    ON customers (email);
```

Do not create a second unique index on the same columns as the primary key. You waste writes and space.

Differences:

| Item | Primary key | Unique index / unique constraint |
| --- | --- | --- |
| `NULL` | Not allowed | Often allowed (product rules differ) |
| Count per table | One | Many |
| Foreign key target | Usual target | Allowed if unique |
| Meaning | The identifier | An extra candidate key |

A non-unique index allows duplicate keys. Use it for foreign keys (`orders.customer_id`) where many orders share one customer.

Lookups on a unique index return at most one row (except the `NULL` case on some products). The optimizer can use that fact.

Do not use a unique index as a substitute for a missing primary key. Give every entity table a primary key. Add unique indexes for other candidate keys.

### Questions

#### Theoretical questions

1. How does a primary key usually appear in storage?
2. Why do you add a unique index on `email` when `customer_id` is the primary key?
3. How does a unique index differ from a non-unique index?
4. Why must you not duplicate the primary-key index?
5. Can a foreign key point to a unique key that is not the primary key?

#### Easy practical tasks

1. Create a table with a primary key and a unique email.
2. Insert a duplicate email. Record the error.
3. Create a non-unique index on a foreign key. Insert two rows with the same foreign key.
4. Write four sentences that compare the two unique structures.

#### Medium practical tasks

1. List indexes after `CREATE TABLE` with `PRIMARY KEY` and `UNIQUE`. Name which objects appeared.
2. Try a unique index on a column that contains two `NULL` values. Write the product rule.
3. Point a foreign key at a unique `sku` instead of `product_id`. Show a valid insert.

#### Advanced practical tasks

1. Write a one-page key-and-index policy: one primary key, extra unique indexes for business keys, non-unique indexes for FKs.
2. Compare `UNIQUE` constraint versus `CREATE UNIQUE INDEX` in your DBMS docs. Write when each appears in errors.

---

## Composite indexes and leftmost prefix

A composite index (multicolumn index) stores more than one column in the key.

```sql
CREATE INDEX idx_orders_cust_date
    ON orders (customer_id, ordered_at);
```

The B-tree orders rows by `customer_id` first, then by `ordered_at`.

**Leftmost prefix.** A query can use the index when it constrains a prefix of the column list from the left.

Useful:

```sql
WHERE customer_id = 42
WHERE customer_id = 42 AND ordered_at >= DATE '2026-01-01'
```

Not useful as a simple seek (the index may still be scanned, but the first column is unused):

```sql
WHERE ordered_at >= DATE '2026-01-01'
```

An index on `(customer_id, ordered_at)` is not a substitute for an index on `(ordered_at)` if you often filter only by date.

Equality on the first column plus a range on the second column is a classic good pattern. A range on the first column can prevent a tight use of the second column.

Order of columns matters. Put the equality column first when one column is equality and the other is a range. Put the most selective equality column first when both are equalities, unless another rule (join, covering) wins.

Do not create both `(a, b)` and `(a)` without a reason. `(a, b)` already supports lookups on `a`. An extra index on `a` duplicates work. An index on `(b, a)` is different: it supports `b` as the leftmost prefix.

A composite unique index enforces uniqueness of the pair, not of each column.

```sql
CREATE UNIQUE INDEX idx_enroll_pair
    ON enrollments (student_id, course_id);
```

The same student can appear many times. The same pair cannot.

### Questions

#### Theoretical questions

1. What is a composite index?
2. What is the leftmost prefix rule?
3. Why does `WHERE ordered_at = ...` not use `(customer_id, ordered_at)` as a simple seek?
4. When do you put an equality column before a range column?
5. Why is `(a, b)` different from `(b, a)`?

#### Easy practical tasks

1. Create an index on `(customer_id, ordered_at)`.
2. Write two queries that use the leftmost prefix.
3. Write one query that filters only `ordered_at`.
4. Explain in four sentences why you might still need a second index on `ordered_at`.

#### Medium practical tasks

1. Compare plans (or times) for filter on `customer_id` versus filter on `ordered_at` with only the composite index.
2. Design a unique composite index for a join table. Show a duplicate pair fail.
3. Choose column order for `(status, created_at)` when queries are `status = ? AND created_at > ?`.

#### Advanced practical tasks

1. Write a one-page index-order guide with three query shapes and one index list that covers them without waste.
2. Explain skip-scan or extra features if your product can use a non-prefix column. Mark that feature as product-specific.

---

## Covering indexes (high-level)

A covering index contains every column that a query needs. The DBMS can answer the query from the index alone. It does not visit the table heap (or the base row).

```sql
CREATE INDEX idx_orders_cust_date
    ON orders (customer_id, ordered_at);

SELECT customer_id, ordered_at
FROM orders
WHERE customer_id = 42;
```

If the index is `(customer_id, ordered_at)`, this query can be covering. If the query also selects `status_code`, the index does not cover unless that column is in the index or in an include list.

Some products allow included (non-key) columns:

```sql
-- product-specific idea
CREATE INDEX idx_orders_cust
    ON orders (customer_id) INCLUDE (ordered_at, status_code);
```

`INCLUDE` is not standard SQL. The idea is portable: extra columns in the index leaf to avoid table lookups.

A covering index helps frequent, narrow queries. It costs more space and more write work. Do not cover `SELECT *`.

The primary-key index covers queries that only need the key columns. A secondary index that does not include the other selected columns needs a lookup to the row.

Read the plan. A plan that says "index only" or "covering" (words differ) means the table row was not required.

Do not add ten included columns. You rebuild a wide index on every write. Cover the columns of one or two hot queries.

### Questions

#### Theoretical questions

1. What does a covering index contain?
2. Why can a covering index skip the table row?
3. Why does `SELECT *` rarely get a covering index?
4. What is the idea of included columns?
5. What is the write cost of a wide covering index?

#### Easy practical tasks

1. Write a query that an index on `(customer_id, ordered_at)` can cover.
2. Add one extra selected column and write why the index may no longer cover.
3. Explain in four sentences the lookup from a secondary index to the row.
4. Mark "index only" as the plan phrase to look for (or the product equivalent).

#### Medium practical tasks

1. Create a narrow index. Select only those columns. Read the plan.
2. Select one more column. Read the plan again. Write the difference.
3. Read whether your DBMS supports `INCLUDE` or an equivalent.

#### Advanced practical tasks

1. Design one covering index for a hot customer-order list query. Justify each column.
2. Write a one-page trade-off: covering index versus an extra table (summary table) for a report.

---

## Selectivity

Selectivity is the fraction of rows that a predicate keeps. A highly selective predicate keeps few rows. A low-selectivity predicate keeps many rows.

```text
WHERE customer_id = 42          -- often highly selective
WHERE country = 'US'            -- can be low selectivity
WHERE active = TRUE             -- often low selectivity
WHERE 1 = 1                     -- no selectivity
```

The optimizer uses statistics to estimate selectivity. High selectivity favors an index seek. Low selectivity favors a scan.

An index on a boolean column is often a waste if half the rows are true. An index on a unique id is highly selective.

Selectivity depends on the value. `country = 'VA'` (Vatican) can be selective. `country = 'US'` may not be. The same column has different selectivity for different literals.

Composite indexes can be more selective than one column. `status = 'open' AND warehouse_id = 7` can be selective even if each column is not.

```sql
-- Low selectivity: many rows match
SELECT *
FROM orders
WHERE status_code = 'closed';

-- Higher selectivity: few rows match
SELECT *
FROM orders
WHERE order_id = 1001;
```

Do not index a column only because it appears in `WHERE`. Check how many rows typically match. Use a plan and row counts.

`DISTINCT` and `GROUP BY` also care about the number of distinct values (cardinality). A later performance topic covers statistics in more depth. Here, remember: few matches → index; most rows match → scan.

### Questions

#### Theoretical questions

1. What is selectivity?
2. Why is `order_id = 1001` more selective than `status_code = 'closed'` in a typical shop?
3. Why can the same column have different selectivity for two values?
4. When does low selectivity favor a scan?
5. Why is an index on a boolean often a poor default?

#### Easy practical tasks

1. Estimate selectivity: 100 matching rows in 1 million. Write the fraction.
2. Label five predicates as high or low selectivity in a shop.
3. Explain in four sentences why the optimizer needs statistics.
4. Write one composite predicate that is more selective than each part.

#### Medium practical tasks

1. Count rows for `status_code = 'open'` versus a single `order_id`. Compare fractions.
2. Read a plan for both filters if the table is large enough. Write which path you see.
3. Find the distinct count of a column (`COUNT(DISTINCT ...)`). Relate it to selectivity.

#### Advanced practical tasks

1. Write a one-page note: selectivity, cardinality, and when you create an index.
2. Show a case where an index exists but the optimizer scans because the estimate is high. Document the row counts.

---

## When an index hurts writes

Every index is extra work on `INSERT`, `UPDATE` of indexed columns, and `DELETE`. The DBMS must maintain each B-tree.

Costs:

- more disk space
- more log traffic
- more CPU to find the insert position and to split pages
- longer transactions and longer locks on index pages

A table with twelve indexes writes like twelve extra structures. A bulk load is much slower.

Updates that do not change indexed columns are cheaper than updates that change the key. A change of an indexed key is a delete plus an insert in the index.

Indexes also slow `COPY` / bulk insert. Some teams drop non-unique indexes, load, and recreate. That tactic is an operations choice. Recreate unique indexes only when you can still enforce uniqueness.

Do not add an index for a query that you run once a month on a table that you write a thousand times a second. The write tax is permanent. The query savings are rare.

Do not keep unused indexes. Unused indexes still cost writes. List indexes and compare them to real queries.

Foreign-key checks use indexes on the referenced key and, for deletes of the parent, on the child foreign-key column. That index has a write cost and a correctness benefit. Keep it.

```text
Writes per insert ≈ 1 table + N indexes
```

Measure write time with and without a candidate index on a scratch table. Then decide.

### Questions

#### Theoretical questions

1. Which write statements maintain indexes?
2. Why does a change of an indexed column cost more than a change of a non-indexed column?
3. Why do unused indexes still hurt?
4. Why do you still index a foreign key despite the write cost?
5. When is a monthly query a weak reason for an extra index on a hot table?

#### Easy practical tasks

1. List every index on one of your tables. Count them.
2. Write the "1 + N" cost model in four sentences.
3. Mark two indexes as required (PK, FK) and one as optional.
4. Explain why bulk load slows down with many indexes.

#### Medium practical tasks

1. Time 10000 inserts into a scratch table with no extra index and with three extra indexes. Write the two times.
2. Update an indexed column versus a non-indexed column. Compare times on a large scratch table.
3. Find an index that no query uses (practice database). Write whether you would drop it.

#### Advanced practical tasks

1. Write a one-page index budget for `orders`: required indexes, optional indexes, rejected indexes.
2. Document a bulk-load idea: drop optional indexes, load, recreate. List risks to uniqueness during the load.

---

## Explain / query plan as a reading skill

A query plan is the access path that the optimizer chose. `EXPLAIN` (the word differs: `EXPLAIN`, `EXPLAIN ANALYZE`, `SHOW PLAN`) prints that path. You read it. You do not guess.

Typical things to read:

- scan versus index seek (or index scan)
- join type (nested loop, hash, merge)
- estimated rows versus actual rows (when the product shows both)
- sort and aggregate steps
- whether the index was covering

```sql
EXPLAIN
SELECT o.order_id, c.name
FROM orders AS o
INNER JOIN customers AS c
    ON c.customer_id = o.customer_id
WHERE o.customer_id = 42;
```

The output is product-specific text or JSON. Learn the tokens of your DBMS. Look for the table name, the index name, and the filter.

`EXPLAIN` without execute shows estimates. `EXPLAIN ANALYZE` (if it exists) runs the query and shows actual times. Do not run `ANALYZE` on a heavy write statement on production without care. It performs the statement on some products. Read the manual.

Reading skill:

1. Find the leaf that reads `orders`.
2. See if it uses `idx_orders_customer_id` or a seq scan.
3. Find the join to `customers`.
4. See if `customers` is a seek on the primary key.
5. Compare estimated rows to what you expect.

A bad plan is a clue. Causes: missing index, hidden expression, stale statistics, wrong selectivity, too-wide `SELECT`.

Do not force an index hint as a first habit. Hints are vendor-specific and can go stale. Fix the SQL and the indexes. Update statistics if the product requires it.

Save plans in your notes for the queries that you tune. A plan is evidence.

### Questions

#### Theoretical questions

1. What does a query plan show?
2. What is the difference between `EXPLAIN` and a form that executes the query?
3. Which tokens do you look for first in a plan?
4. Why must you read the product manual before `EXPLAIN ANALYZE` on a write?
5. Why are index hints a weak first step?

#### Easy practical tasks

1. Run `EXPLAIN` on `SELECT * FROM customers`. Write the access method.
2. Run `EXPLAIN` on a filter on a primary key. Write the access method.
3. Find the official command name for explain in your DBMS.
4. List five tokens that you will search for in plans.

#### Medium practical tasks

1. Explain a join query. Name the join type and the index on each side.
2. Add an index. Explain the same query again. Write what changed.
3. Compare estimate and actual rows if your `EXPLAIN` can show both.

#### Advanced practical tasks

1. Write a one-page plan-reading guide for your DBMS with one annotated plan.
2. Take a slow query (or a made-slow query without an index). Record the plan, add an index, record the new plan.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe the path from a selective `WHERE` to a B-tree seek, and when the optimizer still scans.
2. How do leftmost prefix, covering, and selectivity change the value of one composite index?
3. When does a unique index implement a key, and when does a non-unique index only help access?
4. Why do writes and unused indexes belong in the same design review as `EXPLAIN`?
5. A teammate adds ten indexes after one slow page. Which facts do you use in the reply?

#### Easy practical tasks

1. Write a cheat sheet: scan, B-tree, unique versus PK, composite prefix, covering, selectivity, write cost, explain.
2. Create one non-unique index on a foreign key. Explain a query that uses it.
3. Write three `WHERE` shapes for an index on `(a, b)`.
4. Run `EXPLAIN` on two of your practice queries. Save the output.

#### Medium practical tasks

1. Build `orders` with enough rows. Index `(customer_id, ordered_at)`. Compare plans for prefix and non-prefix filters.
2. Time inserts with and without two extra indexes on a scratch table.
3. Design indexes for a shop: PK, unique sku, FK customer, date range. Justify each.

#### Advanced practical tasks

1. Tune one join report: record the plan, add or change one index, record the new plan and a time.
2. Write an index review for a three-table schema: keep, drop, or add, with a query and a plan for each decision.
