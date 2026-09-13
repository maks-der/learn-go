# 9. Query Planning

## Description

This topic shows how PostgreSQL 16 and PostgreSQL 17 plan and run SQL. You learn `EXPLAIN` and `EXPLAIN ANALYZE`. You learn sequential scans, index scans, and bitmap scans. You also learn join methods, statistics, `work_mem`, and why `SELECT *` hurts hot paths.

Complete topic 8 first. You need indexes to compare plans. Complete this topic before you study transactions in depth.

Use one term for each concept. The planner picks a plan. The executor runs the plan. A scan reads a table or an index. A join method combines two row sources. Statistics are estimates that `ANALYZE` stores. A hot path is a query that runs often or must stay fast.

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

EXPLAIN (ANALYZE, BUFFERS, FORMAT JSON)
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

1. Run `EXPLAIN` on `SELECT * FROM items`.
2. Run `EXPLAIN ANALYZE` on a `SELECT` with a `WHERE`.
3. Run `EXPLAIN (ANALYZE, BUFFERS)` on the same query.
4. Find `cost=` and `rows=` on the top node. Write those numbers.

#### Medium practical tasks

1. Run `EXPLAIN ANALYZE` on an `UPDATE` inside `BEGIN` ... `ROLLBACK`. Confirm the data did not stay changed.
2. Compare `FORMAT TEXT` and `FORMAT JSON` for the same query.
3. Find a query where estimated rows and actual rows differ by a large factor. Write both numbers.

#### Advanced practical tasks

1. Enable `SETTINGS` and `WAL` on an `EXPLAIN ANALYZE` of a small `INSERT`. Write one setting and one WAL line if they appear.
2. Read the official `EXPLAIN` page. List three options that this section named and one option that it did not name.

---

## Seq scan vs index scan vs bitmap scan

A **sequential scan** (Seq Scan) reads the heap from start to end. It applies the filter on each row. It is cheap when the table is small or the filter matches a large fraction of rows.

An **index scan** walks a B-tree (or other index) and fetches matching heap tuples. It is cheap when the filter is selective and the index matches the predicate.

An **index-only scan** reads the index and does not fetch the heap when the visibility map says the page is all-visible. `VACUUM` helps index-only scans (topic 11). `INCLUDE` columns help index-only scans (topic 19).

A **bitmap index scan** builds a bitmap of matching heap pages. A **bitmap heap scan** then reads those pages. PostgreSQL can combine bitmaps with `AND` and `OR` (BitmapAnd, BitmapOr). Bitmap scans work well when many rows match but not the whole table, and when two indexes combine.

```sql
EXPLAIN ANALYZE
SELECT * FROM items WHERE qty > 10 AND name = 'nail';
```

You might see:

- Seq Scan
- Index Scan using items_pkey
- Index Only Scan
- Bitmap Heap Scan + Bitmap Index Scan

The planner chooses. You do not write the scan type in SQL. You change SQL, indexes, and statistics.

`SELECT *` often blocks an index-only scan because the index does not hold every column.

A random `id = $1` on a large table should use an index scan on the primary key. If you see a seq scan, check that the table is not tiny, that `enable_seqscan` is on (keep it on), and that statistics exist.

Do not set `enable_seqscan = off` in production to force an index. That setting is a debug tool.

### Questions

#### Theoretical questions

1. What does a sequential scan read?
2. What does an index scan do after it finds a key?
3. What extra condition does an index-only scan need?
4. When does a bitmap scan help?
5. Why is `enable_seqscan = off` not a production fix?

#### Easy practical tasks

1. `EXPLAIN` a `SELECT` with no `WHERE` on a small table. Write the scan type.
2. `EXPLAIN` `WHERE id = 1` on a table with a primary key.
3. `EXPLAIN` a filter that matches almost all rows. Write the scan type.
4. Make a table: scan type, one sentence.

#### Medium practical tasks

1. Create two single-column indexes. Write a `WHERE` that uses both columns. Look for BitmapAnd or two conditions on one scan.
2. Compare `SELECT *` and `SELECT id` on a covering primary key. Look for Index Only Scan.
3. Run `VACUUM` on the table and retry the index-only query (topic 11). Write if the plan changed.

#### Advanced practical tasks

1. On 100000 rows, find a predicate that switches from index scan to seq scan as the selectivity changes. Write both `WHERE` clauses and both plans.
2. Read "Row Estimation Examples" or "Planner Method Configuration" in the docs. Write how `random_page_cost` affects index versus seq scan (high-level).

---

## Nested loop, hash join, merge join

When two row sources join, the planner picks a join method.

**Nested loop**: for each row on the outer side, find matches on the inner side. An index on the inner join key makes this fast. Nested loop is good when the outer side is small.

**Hash join**: build a hash table from the smaller side (the build side), then probe with the other side. Hash join is good for large equality joins without a useful inner index. Hash join needs memory (`work_mem`). If the hash table does not fit, PostgreSQL spills to disk.

**Merge join**: sort both sides on the join key (or use a sorted index) and merge. Merge join is good when both sides are already ordered, or when the join is a large equality join and sort is cheap enough.

```sql
EXPLAIN ANALYZE
SELECT c.email, o.id
FROM customers AS c
JOIN orders AS o ON o.customer_id = c.id;
```

Look for Nested Loop, Hash Join, or Merge Join in the plan.

Outer joins restrict some methods. The planner still has options. A cross join of huge tables is a nested loop of disaster. Filter first.

You do not pick the join method in SQL. You can:

- add an index on the join key (helps nested loop)
- rewrite the query so that one side stays small
- raise `work_mem` for a session that hashes or sorts a large set (next section)
- `ANALYZE` so that row estimates are right

Do not disable join types in production (`enable_hashjoin = off`) except for a short test.

PostgreSQL can also use parallel variants (Parallel Hash Join). Topic 20 covers parallel query. Know that you may see `Parallel` in front of a node.

### Questions

#### Theoretical questions

1. How does a nested loop join work?
2. How does a hash join work?
3. How does a merge join work?
4. Which join method uses a hash table in memory?
5. Why does an index on the inner join key help a nested loop?

#### Easy practical tasks

1. `EXPLAIN` an inner join of two small tables. Write the join method.
2. Add an index on the foreign key. `EXPLAIN` again. Write if the method changed.
3. Draw three boxes: nested loop, hash, merge. One sentence each.
4. Find the words Nested Loop, Hash Join, or Merge Join in a plan.

#### Medium practical tasks

1. Join a large table to a small table. Compare with a join of two large tables. Write the methods.
2. Force a test with `SET enable_hashjoin = off` in a session. `EXPLAIN` the same join. `RESET enable_hashjoin`.
3. Look for Sort nodes above a merge join. Write if a sort appeared.

#### Advanced practical tasks

1. Build two tables of 50000 rows and join on an unindexed column. Then add an index. Compare `EXPLAIN ANALYZE` times.
2. Read "Joins" in the planner docs. Write when the planner prefers hash join over merge join in your own words.

---

## Statistics: `ANALYZE`, `pg_stat_statements`

The planner uses statistics: row counts, most-common values, histograms, and null fractions. `ANALYZE` (and autovacuum analyze) collects them.

```sql
ANALYZE items;
ANALYZE;
```

`ANALYZE items` updates one table. `ANALYZE` without a name updates the database (all tables that the role can analyze).

See stored numbers:

```sql
SELECT relname, reltuples, relpages
FROM pg_class
WHERE relname = 'items';

SELECT * FROM pg_stats
WHERE tablename = 'items';
```

`reltuples` is an estimate. After a large load, run `ANALYZE`. Bad estimates cause bad plans (seq scan instead of index, or the wrong join).

`default_statistics_target` controls how much detail `ANALYZE` stores (default 100). You can raise it per column:

```sql
ALTER TABLE items ALTER COLUMN name SET STATISTICS 200;
ANALYZE items;
```

`pg_stat_statements` is an extension. It records how often each query text runs and how much time it uses.

```sql
CREATE EXTENSION IF NOT EXISTS pg_stat_statements;
```

The module must be in `shared_preload_libraries`. Many packaged servers already enable it. If `CREATE EXTENSION` fails, read the server docs.

```sql
SELECT query, calls, total_exec_time, mean_exec_time, rows
FROM pg_stat_statements
ORDER BY total_exec_time DESC
LIMIT 20;
```

Use this view to find hot paths. Then `EXPLAIN ANALYZE` those queries.

`pg_stat_statements` normalizes literals. Two `SELECT` statements that differ only by a constant can share one row.

Do not confuse `ANALYZE` the SQL command with `EXPLAIN ANALYZE`. The first updates statistics. The second runs a query for timing.

### Questions

#### Theoretical questions

1. What does `ANALYZE` collect?
2. Why do stale statistics cause bad plans?
3. What is `default_statistics_target`?
4. What does `pg_stat_statements` record?
5. How is `ANALYZE` different from `EXPLAIN ANALYZE`?

#### Easy practical tasks

1. Run `ANALYZE` on one table. Select `reltuples` from `pg_class`.
2. Select two rows from `pg_stats` for that table.
3. Try `CREATE EXTENSION pg_stat_statements`. Write success or the error.
4. Write four sentences: statistics, histogram, extension, hot query.

#### Medium practical tasks

1. Load many rows. Compare a plan before `ANALYZE` and after `ANALYZE`.
2. Raise `STATISTICS` on one column. `ANALYZE`. Compare `pg_stats` histogram bounds count.
3. If the extension works, run a query ten times. Find it in `pg_stat_statements`.

#### Advanced practical tasks

1. Read how to enable `shared_preload_libraries` for `pg_stat_statements`. Write the config line and whether you need a restart.
2. Find the top three statements by `total_exec_time`. Write an `EXPLAIN ANALYZE` for each (use `ROLLBACK` on writes).

---

## Work_mem and sorts / hashes (high-level)

`work_mem` is the amount of memory that one query operation may use for a sort or a hash table. The unit is memory size. Default is often `4MB`. You can set it in `postgresql.conf` or in a session:

```sql
SHOW work_mem;
SET work_mem = '32MB';
```

A query can use several sorts and hashes. Each operation can use up to `work_mem`. A session with many operations can use more than one `work_mem`. `max_connections` times many operations can exhaust RAM. Do not set a huge global `work_mem`.

When a sort does not fit, PostgreSQL spills to disk (external merge sort). When a hash does not fit, it spills in batches. Disk is slower than memory.

Plans show this:

- `Sort Method: quicksort Memory: ...`
- `Sort Method: external merge Disk: ...`
- `Hash Batches` greater than 1 means spill

Use `EXPLAIN (ANALYZE, BUFFERS)` to see temp reads and writes.

`maintenance_work_mem` is a separate setting for `VACUUM`, `CREATE INDEX`, and `ALTER TABLE`. Do not confuse it with `work_mem`. Topic 18 covers memory settings in more depth.

Raise `work_mem` for one heavy report session when you see disk sorts. Reset it after. Keep the global value modest.

`hash_mem_multiplier` (PostgreSQL 13+) can let hash operations use more than one `work_mem`. Know the name. Read the docs before you change it.

Do not raise `work_mem` to hide a missing index. Fix the query first.

### Questions

#### Theoretical questions

1. What operations use `work_mem`?
2. Why is a very large global `work_mem` a risk?
3. What happens when a sort does not fit in `work_mem`?
4. What is `maintenance_work_mem` for?
5. How do you raise `work_mem` for one session only?

#### Easy practical tasks

1. Run `SHOW work_mem;` and `SHOW maintenance_work_mem;`.
2. `SET work_mem = '16MB';` then `SHOW work_mem;`. `RESET work_mem`.
3. Find `Sort Method` in an `EXPLAIN ANALYZE` of `ORDER BY` on a larger table.
4. Write four sentences: per operation, spill, session set, global risk.

#### Medium practical tasks

1. Force a disk sort: set `work_mem` very low (example `'64kB'`) and `ORDER BY` a big result. Confirm `external merge`. Restore `work_mem`.
2. Compare the same query with a higher `work_mem`. Write both sort methods.
3. Read `hash_mem_multiplier` in the docs. Write the default.

#### Advanced practical tasks

1. Estimate a worst case: `max_connections` × number of hash nodes × `work_mem`. Use your server values. Write the product.
2. Find temp file activity in `EXPLAIN (ANALYZE, BUFFERS)` or `pg_stat_database`. Write what you used.

---

## Avoiding `SELECT *` in hot paths

`SELECT *` returns every column. It is fine in `psql` while you learn. It is a problem in hot paths.

Why it hurts:

- The executor reads columns that the application does not use.
- Index-only scans become less likely. The index does not contain every column.
- More data goes over the network.
- `CREATE TABLE` plus a new column changes the result shape. Client code that maps by position can break.
- `json` / `text` / `bytea` columns can be large. `SELECT *` always pulls them.

Write the column list:

```sql
SELECT id, email, created_at
FROM customers
WHERE id = 1;
```

Exceptions:

- Admin `psql` sessions
- `RETURNING *` on a narrow table when you need the full row
- Quick exploration

For `jsonb` payloads, select `id` and the keys that you need (`payload->>'type'`), not the whole document, when the document is large.

`SELECT t.*` in a join still pulls every column of `t`. Name columns.

Wide tables (many columns or large text) need extra care. Split cold columns to another table if you measure a problem (later design work).

Do not use `SELECT *` in a view that applications call on a hot path. Views that hide `*` still expand to all columns.

### Questions

#### Theoretical questions

1. Why does `SELECT *` block some index-only scans?
2. How does `SELECT *` affect the network?
3. Why can a new column break clients that use `SELECT *`?
4. When is `SELECT *` acceptable?
5. What is a hot path in this topic?

#### Easy practical tasks

1. Rewrite a `SELECT *` as a three-column list.
2. `EXPLAIN` `SELECT *` versus `SELECT id` on a table with a primary key. Write both scan types if they differ.
3. Add a `text` column. Show that `SELECT *` now includes it.
4. Write four sentences: columns, index-only, network, stability.

#### Medium practical tasks

1. Compare result size: `psql` with `\a` and a wide `SELECT *` versus two columns. Use a large `text` value.
2. Create a view as `SELECT * FROM items`. Add a column to `items`. Select from the view. Write what happened (`*` in views is expanded at create time in PostgreSQL — verify and write the fact).
3. Find three application queries (or write three) and mark any `*`.

#### Advanced practical tasks

1. Measure bytes and time for `SELECT *` versus a narrow list on 10000 rows with a 1 KB text column. Write both measurements.
2. Read how PostgreSQL stores views (`pg_get_viewdef`). Confirm whether `*` expanded at create time. Write the result.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe the path from a SQL string to a plan: parse, rewrite, plan, execute. One sentence each.
2. How do statistics, indexes, and `work_mem` change the same join?
3. When do you trust `EXPLAIN` without `ANALYZE`, and when must you run `ANALYZE`?
4. Why do bitmap scans and hash joins both appear on large filters or large joins?
5. How does `pg_stat_statements` tell you which query to explain first?

#### Easy practical tasks

1. Pick one join query. Run `EXPLAIN (ANALYZE, BUFFERS)`. Label scan type and join type.
2. Write a cheat sheet: `EXPLAIN` options, three scan types, three join types, `ANALYZE`, `work_mem`, `SELECT *`.
3. Run `ANALYZE` on your schema. Then explain one hot `SELECT`.
4. Rewrite two hot queries to drop `SELECT *`.

#### Medium practical tasks

1. Create a report that lists the five slowest statements from `pg_stat_statements` (or five hand-timed queries if the extension is off).
2. Show a disk sort, then fix it with an index or a session `work_mem`. Write which fix you chose and why.
3. Document a team rule: every new API query must have a saved `EXPLAIN (ANALYZE, BUFFERS)` on production-sized data.

#### Advanced practical tasks

1. Change `random_page_cost` in a session and re-explain one query. Write if the scan type changed. `RESET` the setting.
2. Read the "Using EXPLAIN" chapter. Write a one-page walkthrough of one of your plans with official terms only.
