# 5. Indexes and Query Plans

## Description

This topic shows indexes and query plans in PostgreSQL 16 and PostgreSQL 17. You learn B-tree, Hash, GiST, GIN, and BRIN. You learn partial, expression, unique, and composite indexes. You also learn `CREATE INDEX CONCURRENTLY`, `REINDEX`, `EXPLAIN`, join methods, `ANALYZE`, and `pg_stat_statements`.

Complete topic 4 first. You know the operators that indexes support. Complete this topic before you study transactions and vacuum.

Use one term for each concept. An index is a secondary structure that helps the server find rows. The heap is the table. The planner picks a plan. The executor runs the plan. A scan reads a table or an index. A join method combines two row sources. Statistics are estimates that `ANALYZE` stores.

---

## B-tree, Hash, GiST, GIN, BRIN

`CREATE INDEX` without `USING` builds a B-tree.

```sql
CREATE INDEX items_name_idx ON items (name);
CREATE INDEX items_qty_idx ON items (qty);
```

A B-tree supports comparison operators: `=`, `<`, `<=`, `>`, `>=`. It supports `ORDER BY` on the same column and direction. It supports `BETWEEN`. It supports `LIKE 'abc%'` (prefix). It does not support `LIKE '%abc'` by itself.

B-tree is the right index for primary keys, unique keys, foreign-key columns, equality, ranges, and sort keys that match the query.

You select another class with `USING`:

| Class | Typical use |
| --- | --- |
| **B-tree** | default; equality and range |
| **Hash** | equality only (`=`); rare; no ordering |
| **GiST** | ranges, geometry, some full-text and trigram |
| **GIN** | arrays, `jsonb`, full-text `tsvector` |
| **BRIN** | very large tables with physical correlation (time-append) |

```sql
CREATE INDEX items_id_hash ON items USING HASH (id);
CREATE INDEX events_payload_gin ON events USING GIN (payload);
CREATE INDEX readings_created_brin ON readings USING BRIN (created_at);
```

Hash indexes are WAL-logged in current versions. They still cannot support `ORDER BY` or ranges. Prefer B-tree unless you measure a win.

GiST is a balanced tree for types that are not a simple scalar order. PostGIS (topic 14) uses GiST. `tstzrange` and similar range types often use GiST.

GIN is an inverted index. It is good when one row contains many keys (array elements, JSON keys, lexemes). GIN is large. Writes are slower.

BRIN stores min/max summaries per page range. It is small. It helps when rows that share a key sit together on disk (insert by time). It fails when values are random on disk.

The planner uses an index only when it estimates that the index is cheaper than a sequential scan. A filter that matches most rows often uses a sequential scan. That choice can be correct.

Each extra index slows `INSERT`, `UPDATE`, and `DELETE`. Do not create a second B-tree on the same column as the primary key. Name indexes with a clear pattern: `tablename_column_idx`.

### Questions

#### Theoretical questions

1. What index type does `CREATE INDEX` build by default?
2. Which operators can a B-tree support?
3. When do you pick GIN instead of B-tree?
4. When does BRIN work well?
5. Why does an extra index slow writes?

#### Easy practical tasks

1. Create a B-tree on a `text` column. Run `\d` and find the index.
2. Create a GIN index on a `text[]` or `jsonb` column.
3. Query `pg_indexes` for your table.
4. Write five sentences: B-tree, Hash, GiST, GIN, BRIN.

#### Medium practical tasks

1. Insert a few thousand rows. `EXPLAIN` `WHERE id = 1` and `WHERE qty > 0` if most rows match. Write the scan types.
2. Compare `LIKE 'na%'` and `LIKE '%na%'`. Write whether a B-tree can help each.
3. Create a BRIN index on an insert-time column. `EXPLAIN` a range filter.

#### Advanced practical tasks

1. Read "Index Types" in the 17 docs. Write one operator class example for GiST or GIN that this section omitted.
2. Measure `INSERT` time into a table with no extra index and with three B-trees. Write the two times.

---

## Partial, expression, unique, and composite indexes

A **partial** index covers a subset of rows:

```sql
CREATE INDEX orders_open_idx ON orders (id)
    WHERE status = 'open';
```

The planner can use it only when the query `WHERE` implies the predicate. Partial indexes stay small.

An **expression** index stores a computed value:

```sql
CREATE INDEX users_email_lower_idx ON users (lower(email));
```

The query must use the same expression: `WHERE lower(email) = 'a@b.com'`. Topic 4 used `(payload ->> 'type')` as an expression index.

A **unique** index forbids duplicates. `CREATE UNIQUE INDEX` and a `UNIQUE` constraint both build a unique index. A unique index can be partial:

```sql
CREATE UNIQUE INDEX users_one_active_email
    ON users (email)
    WHERE deleted_at IS NULL;
```

A **composite** (multicolumn) index has a column list:

```sql
CREATE INDEX orders_cust_created_idx
    ON orders (customer_id, created_at);
```

B-tree composite indexes are left-prefix. A query that filters `customer_id` can use this index. A query that filters only `created_at` usually cannot use it well.

Column order matters. Put equality columns first. Put range or sort columns after.

You can add non-key columns with `INCLUDE` (covering index). Topic 13 covers `INCLUDE` with pagination. Know that `INCLUDE` columns are not part of the search key.

```sql
CREATE INDEX orders_cust_inc
    ON orders (customer_id) INCLUDE (created_at, total);
```

Do not build five overlapping composite indexes "just in case". Do not unique-index a column that must allow duplicates. Do not forget that functions in queries must match expression indexes.

### Questions

#### Theoretical questions

1. When can the planner use a partial index?
2. What must a query write to use `lower(email)` as an index?
3. Can a unique index be partial?
4. What is the left-prefix rule for a B-tree on `(a, b)`?
5. What does `INCLUDE` add?

#### Easy practical tasks

1. Create a partial index `WHERE qty > 0`. Run `\d`.
2. Create `(lower(name))` and query with `lower(name) = ...`.
3. Create `(customer_id, created_at)`.
4. Write four sentences: partial, expression, unique, composite.

#### Medium practical tasks

1. `EXPLAIN` a query that matches the partial predicate and one that does not. Write which plan uses the index.
2. Query only the second column of a composite B-tree. Write the scan type.
3. Create a partial unique index that allows many `NULL` or deleted emails.

#### Advanced practical tasks

1. Design indexes for `WHERE customer_id = $1 AND created_at >= $2 ORDER BY created_at`. Write the column order and why.
2. Read "Multicolumn Indexes" and "Partial Indexes" in the 17 docs. Write two planner rules in your own words.

---

## `CREATE INDEX CONCURRENTLY` and `REINDEX`

`CREATE INDEX` on a live table takes a strong lock that blocks writes. `CREATE INDEX CONCURRENTLY` builds the index without a long write lock. It waits for old transactions. It scans the table twice.

```sql
CREATE INDEX CONCURRENTLY items_name_idx ON items (name);
```

Rules:

- You cannot run `CONCURRENTLY` inside a transaction block.
- The build can fail and leave an **invalid** index. `\d` or `pg_index.indisvalid` shows it.
- Drop an invalid index and try again: `DROP INDEX CONCURRENTLY items_name_idx;`
- `CONCURRENTLY` is slower than a normal build.

`CREATE UNIQUE INDEX CONCURRENTLY` is valid. Unique violations fail the build.

`REINDEX` rebuilds an index. Use it after bloat or corruption. `REINDEX INDEX name` rebuilds one index. `REINDEX TABLE name` rebuilds all indexes on that table.

```sql
REINDEX INDEX CONCURRENTLY items_name_idx;
REINDEX TABLE CONCURRENTLY items;
```

`REINDEX CONCURRENTLY` is available in current versions. It has similar transaction rules. It cannot run in a transaction block.

PostgreSQL 16 and 17 also have `REINDEX` options such as `TABLESPACE`. System catalogs have extra limits. Do not `REINDEX` system catalogs casually on a shared host.

A failed concurrent reindex can leave a temporary index. Read the docs and `pg_index` before you drop objects.

Do not use `CONCURRENTLY` on an empty lab table just to type the word. Use it when writes must continue. Do not ignore an invalid index. The planner will not use it.

### Questions

#### Theoretical questions

1. What lock problem does `CREATE INDEX CONCURRENTLY` avoid?
2. Why can you not run `CONCURRENTLY` in a `BEGIN` block?
3. What is an invalid index?
4. What does `REINDEX TABLE` rebuild?
5. Is `REINDEX CONCURRENTLY` valid in PostgreSQL 16 and 17?

#### Easy practical tasks

1. Build a normal index. Then drop it. Build it `CONCURRENTLY`.
2. Query `indisvalid` in `pg_index` for that index.
3. Run `REINDEX INDEX` on a practice index.
4. Write four sentences: lock, concurrent, invalid, reindex.

#### Medium practical tasks

1. Try `BEGIN; CREATE INDEX CONCURRENTLY ...`. Record the error.
2. Open a second session with a long `BEGIN`. Start `CREATE INDEX CONCURRENTLY` and watch `pg_stat_activity` for wait.
3. `REINDEX TABLE CONCURRENTLY` a small table. Confirm indexes stay valid.

#### Advanced practical tasks

1. Read `CREATE INDEX` and `REINDEX` in the 17 docs. List three limits of `CONCURRENTLY`.
2. Write a short runbook: how you detect an invalid index and how you replace it. Do not run it on production.

---

## `EXPLAIN` and `EXPLAIN ANALYZE`

`EXPLAIN` prints the plan. It does not run the query (except some forms that still prepare work). `EXPLAIN ANALYZE` runs the query and prints actual times and row counts.

```sql
EXPLAIN
SELECT * FROM items WHERE name = 'nail';

EXPLAIN ANALYZE
SELECT * FROM items WHERE name = 'nail';
```

Useful options:

```sql
EXPLAIN (ANALYZE, BUFFERS, FORMAT TEXT)
SELECT * FROM items WHERE qty > 10;
```

Common options in PostgreSQL 16 and 17:

- `ANALYZE` — run the query
- `BUFFERS` — show shared and read buffers
- `WAL` — show WAL records (useful on writes)
- `SETTINGS` — show planner settings that differ from default
- `VERBOSE` — extra detail
- `COSTS` — show costs (default on)
- `TIMING` — timer per node (default on with `ANALYZE`)
- `SUMMARY` — total time summary

`EXPLAIN ANALYZE` on `INSERT`, `UPDATE`, or `DELETE` **runs** the change. Use a transaction and `ROLLBACK` when you only want the plan:

```sql
BEGIN;
EXPLAIN ANALYZE UPDATE items SET qty = qty + 1;
ROLLBACK;
```

Read a text plan from the top. The first node is the root. Indentation shows children. PostgreSQL runs children first.

Two numbers on each node: estimated cost (`cost=startup..total`) and estimated rows (`rows=`). After `ANALYZE`, you also see `actual time` and `actual rows`. A large gap between estimated rows and actual rows means bad statistics or a bad estimate for the predicate.

Do not tune a query from cost numbers alone. Use actual time from `ANALYZE`. Do not run `EXPLAIN ANALYZE` on a heavy write in production without a plan.

`EXPLAIN` does not include network time to your client. It does not include application loops.

### Questions

#### Theoretical questions

1. What is the difference between `EXPLAIN` and `EXPLAIN ANALYZE`?
2. Why must you wrap a write `EXPLAIN ANALYZE` in `ROLLBACK` when you test?
3. What does `BUFFERS` add?
4. What do estimated `rows` versus actual rows tell you?
5. Does `EXPLAIN` include client network time?

#### Easy practical tasks

1. `EXPLAIN` a `SELECT` with a `WHERE` on a primary key.
2. `EXPLAIN ANALYZE` the same query. Write estimated and actual rows.
3. Run `EXPLAIN (BUFFERS)` and find `shared hit` or `read`.
4. Write five sentences: planner, executor, cost, actual time, `ROLLBACK`.

#### Medium practical tasks

1. Compare `FORMAT TEXT` and `FORMAT JSON` for the same query.
2. `EXPLAIN ANALYZE` an `UPDATE` inside `BEGIN`/`ROLLBACK`.
3. Change a parameter with `SET enable_seqscan = off` in one session. `EXPLAIN` again. Reset the setting.

#### Advanced practical tasks

1. Find a query where estimated rows and actual rows differ by 10× or more. Run `ANALYZE` on the table. Explain again.
2. Read `EXPLAIN` in the 17 docs. List three options that this section named and one that it did not.

---

## Seq scan, index scan, nested loop, hash join, merge join

A **sequential scan** reads the heap from start to end. It is cheap per row when you need most rows. It is expensive when you need one row from a large table.

An **index scan** walks the index and fetches heap tuples. An **index-only scan** can use the visibility map and skip the heap when all needed columns are in the index (and visibility allows it). A **bitmap index scan** plus **bitmap heap scan** is efficient when many rows match.

```sql
EXPLAIN SELECT * FROM items WHERE id = 1;
EXPLAIN SELECT * FROM items WHERE qty > 0;
```

Join methods:

| Method | Idea |
| --- | --- |
| **Nested loop** | For each outer row, find inner rows. Good for small outer sets and indexed inner lookups. |
| **Hash join** | Build a hash table on one side. Probe with the other. Good for larger equality joins. |
| **Merge join** | Sort both sides (or scan ordered indexes) and merge. Good for sorted equality or some ranges. |

The planner picks a method from costs and statistics. You can force methods for study:

```sql
SET enable_hashjoin = off;
SET enable_mergejoin = off;
```

Reset after the test. Do not leave these settings in production.

`work_mem` limits memory for a sort or hash step (topic 13). A small `work_mem` can spill to disk. `EXPLAIN ANALYZE` shows `Sort Method` and hash batches.

A join that looks like a nested loop in `EXPLAIN` is not a programming loop that you write. It is the executor algorithm.

Do not add an index for every sequential scan that you see. Measure selectivity. Do not disable sequential scan globally.

### Questions

#### Theoretical questions

1. When is a sequential scan a good plan?
2. What extra condition does an index-only scan need?
3. When is a nested loop a typical choice?
4. What does a hash join build?
5. What does a merge join need from its inputs?

#### Easy practical tasks

1. `EXPLAIN` a primary-key lookup. Write the scan node name.
2. `EXPLAIN` `SELECT *` from a small table. Write whether you see `Seq Scan`.
3. Join two tables. Write the join method.
4. Make a three-column table: join method, one-sentence idea, one typical case.

#### Medium practical tasks

1. Disable hash join in a session. `EXPLAIN` the same join. Write the new method.
2. Compare `SELECT id FROM items WHERE id = 1` with `SELECT *` on a covering-friendly index. Look for `Index Only Scan`.
3. Find `Bitmap Heap Scan` on a query that matches many rows via an index.

#### Advanced practical tasks

1. Build two tables of 100000 rows. Compare nested loop, hash, and merge by toggling `enable_*`. Write times from `EXPLAIN ANALYZE`.
2. Read "How the Planner Uses Statistics" preview. Write how row estimates change the join method.

---

## `ANALYZE` and `pg_stat_statements`

`ANALYZE` collects statistics: row counts, most-common values, histograms, and null fractions. The planner reads `pg_statistic` (through views such as `pg_stats`).

```sql
ANALYZE items;
ANALYZE;
```

`ANALYZE` without a table name analyzes all tables in the database that you may read. Autovacuum also runs analyze (topic 7). After a bulk load, run `ANALYZE` yourself. Do not wait if you will `EXPLAIN` at once.

```sql
SELECT attname, n_distinct, most_common_vals
FROM pg_stats
WHERE tablename = 'items';
```

Bad statistics cause bad row estimates. Bad estimates cause wrong join methods and wrong scan types.

`pg_stat_statements` records a fingerprint of SQL and cumulative time. It is an extension (topic 9). Enable it in `postgresql.conf`:

```text
shared_preload_libraries = 'pg_stat_statements'
```

Then in the database:

```sql
CREATE EXTENSION IF NOT EXISTS pg_stat_statements;

SELECT query, calls, total_exec_time, mean_exec_time, rows
FROM pg_stat_statements
ORDER BY total_exec_time DESC
LIMIT 20;
```

PostgreSQL 16 and 17 use `total_exec_time` and related exec columns. Older text said `total_time`.

`pg_stat_statements` normalizes literals. Many similar queries become one row. It does not replace `EXPLAIN` on a single call. It tells you which statements burn time in production.

`pg_stat_user_indexes` shows index scans. `idx_scan = 0` after a week can mean an unused index.

Do not run `ANALYZE` in a tight loop. Do not enable `pg_stat_statements` without `shared_preload_libraries`. Do not treat the top query as "broken" without `EXPLAIN ANALYZE`.

### Questions

#### Theoretical questions

1. What does `ANALYZE` collect?
2. Why do you run `ANALYZE` after a bulk load?
3. What does `pg_stat_statements` store?
4. Why do you need `shared_preload_libraries` for `pg_stat_statements`?
5. What can `idx_scan = 0` mean?

#### Easy practical tasks

1. Run `ANALYZE` on one table. Select `pg_stats` for that table.
2. `EXPLAIN` a query before and after `ANALYZE` on a table that you just filled.
3. If the extension is available, `CREATE EXTENSION pg_stat_statements` and select one row.
4. Write four sentences: statistics, planner, extension, unused index.

#### Medium practical tasks

1. Load 10000 rows. Compare plans before and after `ANALYZE`.
2. Query `pg_stat_user_indexes` for `idx_scan` and `idx_tup_read`.
3. Read `pg_stat_statements` columns in `\d pg_stat_statements`. Write three columns that you will watch.

#### Advanced practical tasks

1. Enable `pg_stat_statements` on a lab cluster (restart if you must). Run a workload. List the top five statements by `total_exec_time`.
2. Write a weekly check: tables with stale `last_analyze`, indexes with zero scans, and the top statements. Use catalog views.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do you choose between a B-tree on a typed column and a GIN on `jsonb`?
2. Why can `CREATE INDEX CONCURRENTLY` still fail after a long wait?
3. What does a large gap between `EXPLAIN` rows and `EXPLAIN ANALYZE` rows tell you to do first?
4. When is a hash join a better guess than a nested loop?
5. How do `ANALYZE` and `pg_stat_statements` answer different questions?

#### Easy practical tasks

1. Create a table with 1000 rows, one B-tree, and one query. Save `EXPLAIN ANALYZE`.
2. Add a partial index. Show `\d` and one matching `EXPLAIN`.
3. Run `ANALYZE` and select `n_distinct` from `pg_stats`.
4. Write a cheat sheet: five index classes and three join methods.

#### Medium practical tasks

1. Build `(customer_id, created_at)` and a GIN on `tags`. Write two queries and their scan types.
2. Rebuild one index `CONCURRENTLY`. Confirm `indisvalid`.
3. Compare sequential scan and index scan on the same column by changing the filter so that selectivity changes.

#### Advanced practical tasks

1. Take one slow lab query. Change only statistics or only one index. Prove the plan change with `EXPLAIN ANALYZE`.
2. Read "Using EXPLAIN" in the 17 docs. Annotate one real plan line by line (cost, rows, width, actual time).
