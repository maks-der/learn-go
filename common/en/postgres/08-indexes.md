# 8. Indexes

## Description

This topic shows indexes in PostgreSQL 16 and PostgreSQL 17. You learn the default B-tree and when you use Hash, GiST, GIN, and BRIN. You also learn partial indexes, expression indexes, unique and composite indexes, `CREATE INDEX CONCURRENTLY`, and `REINDEX`.

Complete topic 7 first. You must know the operators that indexes support. Complete this topic before you read query plans.

Use one term for each concept. An index is a secondary structure that helps the server find rows. The heap is the table. A scan can use an index or read the heap in order (sequential scan). An index does not copy the full row unless you add extra columns (`INCLUDE`, topic 19).

---

## B-tree (default)

`CREATE INDEX` without `USING` builds a B-tree.

```sql
CREATE INDEX items_name_idx ON items (name);
CREATE INDEX items_qty_idx ON items (qty);
```

A B-tree supports comparison operators: `=`, `<`, `<=`, `>`, `>=`. It supports `ORDER BY` on the same column and direction. It supports `BETWEEN`. It supports `LIKE 'abc%'` (prefix). It does not support `LIKE '%abc'` by itself.

B-tree is the right index for:

- primary keys and unique keys (PostgreSQL already creates those)
- foreign key columns
- common equality and range filters
- sort keys that match the query

```sql
SELECT * FROM items WHERE id = 10;
SELECT * FROM items WHERE qty >= 100 ORDER BY qty;
SELECT * FROM items WHERE name LIKE 'na%';
```

The planner uses an index only when it estimates that the index is cheaper than a sequential scan. A filter that matches most rows often uses a sequential scan. That choice can be correct.

B-tree indexes are balanced trees. Lookups are logarithmic in the number of rows. Updates must maintain the index. Each extra index slows `INSERT`, `UPDATE`, and `DELETE`.

Name indexes with a clear pattern: `tablename_column_idx`. Unique indexes often use `tablename_column_key`.

`\d items` lists indexes. `pg_indexes` lists definitions.

Do not create a second B-tree on the same column as the primary key. Do not index a column when you have no query that needs that index.

### Questions

#### Theoretical questions

1. What index type does `CREATE INDEX` build by default?
2. Which operators can a B-tree support?
3. Can a B-tree support `LIKE '%abc'`?
4. Why does an extra index slow writes?
5. When can a sequential scan be the better plan?

#### Easy practical tasks

1. Create a B-tree on a `text` column. Run `\d` and find the index.
2. Create a B-tree on an `integer` column that you filter with `>=`.
3. Query `pg_indexes` for your table.
4. Write four sentences: B-tree, equality, range, prefix `LIKE`.

#### Medium practical tasks

1. Insert a few thousand rows. Run `EXPLAIN` on `WHERE id = 1` and on `WHERE qty > 0` if most rows match. Write the scan types.
2. Run `LIKE 'na%'` and `LIKE '%na%'`. Compare `EXPLAIN` (topic 9 if you need more detail).
3. Drop an unused practice index. Show `\d` after the drop.

#### Advanced practical tasks

1. Read "Multicolumn Indexes" preview: create `(customer_id, created_at)`. Write two queries: one that uses the first column, one that uses only the second. Compare `EXPLAIN`.
2. Measure `INSERT` time into a table with no extra index and with three B-trees. Write the two times.

---

## Hash, GiST, GIN, BRIN (when each class is for)

You select a class with `USING`:

```sql
CREATE INDEX ... ON table USING btree (col);  -- default
CREATE INDEX ... ON table USING hash (col);
CREATE INDEX ... ON table USING gist (col);
CREATE INDEX ... ON table USING gin (col);
CREATE INDEX ... ON table USING brin (col);
```

**Hash** supports equality (`=`) only. Hash indexes are WAL-logged in current PostgreSQL. A B-tree also supports `=` and also supports inequality and `ORDER BY`. Prefer B-tree unless you have a measured reason to use hash.

**GiST** (Generalized Search Tree) supports many operator classes: geometric types, range types, full-text (`tsvector` with GiST), `pg_trgm`, and exclusion constraints. Use GiST when the type docs say GiST. Exclusion constraints on ranges use GiST (topic 5).

**GIN** (Generalized Inverted Index) is for values that contain many elements: `jsonb`, arrays, `tsvector`, and `pg_trgm`. A GIN index maps elements to row pointers. Use GIN for `jsonb @>`, array `@>`, and full-text `@@`.

```sql
CREATE INDEX events_payload_gin ON events USING gin (payload);
CREATE INDEX packs_tags_gin ON packs USING gin (tags);
CREATE INDEX posts_tsv_gin ON posts USING gin (tsv);
```

**BRIN** (Block Range Index) stores summaries for page ranges. BRIN is small. It helps when the column correlates with physical order (time series appended to a table). Use BRIN on a large `timestamptz` that you insert in time order. Do not use BRIN on a random `uuid`.

```sql
CREATE INDEX readings_ts_brin ON readings USING brin (read_at);
```

Pick the class from the type and the operator, not from habit. If you only need `=` and `<` on an integer, use B-tree.

### Questions

#### Theoretical questions

1. What operator does a hash index support?
2. When do you use GiST?
3. When do you use GIN?
4. When do you use BRIN?
5. Why is B-tree still the usual choice for `integer` equality?

#### Easy practical tasks

1. Write one `CREATE INDEX` line for each class (btree, hash, gist, gin, brin) on a suitable column type. Create the ones your types allow.
2. Create a GIN index on `jsonb` or `text[]`.
3. Make a table: class, typical type, typical operator.
4. Find `USING` in `\d` or `pg_indexes`.

#### Medium practical tasks

1. On a `jsonb` column, compare `EXPLAIN` for `@>` with no index and with GIN.
2. Create a BRIN index on a `timestamptz` column. Read `pg_relation_size` for that index versus a B-tree on the same column.
3. Read the official index types page. Write one sentence for SP-GiST (know that it exists).

#### Advanced practical tasks

1. Create an exclusion constraint that uses GiST (topic 5 ranges). Show the index that PostgreSQL created.
2. Write a short report: three queries that are wrong for BRIN (random order, small table, leading wildcard on text).

---

## Partial indexes

A partial index uses a `WHERE` clause. The index stores only rows that match the predicate.

```sql
CREATE INDEX orders_open_idx
ON orders (created_at)
WHERE status = 'open';
```

The planner can use this index when the query `WHERE` implies the predicate:

```sql
SELECT * FROM orders
WHERE status = 'open' AND created_at > DATE '2026-01-01';
```

A query that omits `status = 'open'` does not use `orders_open_idx`.

Use a partial index when:

- a small subset is hot (`status = 'open'`, `deleted_at IS NULL`)
- you want a unique rule on a subset (`UNIQUE (email) WHERE deleted_at IS NULL`)

```sql
CREATE UNIQUE INDEX users_active_email_key
ON users (email)
WHERE deleted_at IS NULL;
```

This allows many rows with the same email if they are deleted, and one email among active rows.

Partial indexes are smaller than full indexes. They reduce write cost for rows that do not match the predicate.

Keep the predicate simple. Use immutable expressions. The query must use a compatible condition. `WHERE status IN ('open')` often matches `status = 'open'`. A function wrapper can hide the match.

### Questions

#### Theoretical questions

1. What extra clause does a partial index have?
2. When can the planner use a partial index?
3. Why is a partial index smaller?
4. How can a unique partial index model "unique when active"?
5. Why must the query `WHERE` imply the index predicate?

#### Easy practical tasks

1. Create a partial index on `status = 'open'` (add a `status` column if needed).
2. Run `EXPLAIN` on a query that includes `status = 'open'`.
3. Run `EXPLAIN` on a query that omits `status`. Write if the partial index appears.
4. Write four sentences: subset, predicate, query match, unique active email.

#### Medium practical tasks

1. Create `UNIQUE (...) WHERE deleted_at IS NULL`. Insert two rows with the same email: one deleted, one active. Then try two active rows.
2. Compare `pg_relation_size` of a full index and a partial index on the same column.
3. Change the query to `status <> 'closed'` and see if the `status = 'open'` index is used.

#### Advanced practical tasks

1. Read "Partial Indexes" in the docs. Write the official rule about predicate implication in your own words.
2. Design three partial indexes for a large `events` table (errors only, last day, unpaid). Write the `WHERE` clauses and the queries that hit them.

---

## Expression indexes

An expression index stores the result of an expression, not a raw column.

```sql
CREATE INDEX users_email_lower_idx ON users (lower(email));
```

The planner can use this index when the query uses the same expression:

```sql
SELECT * FROM users WHERE lower(email) = 'a@example.com';
```

`WHERE email = 'A@example.com'` does not use `lower(email)` unless you also write `lower(email)`.

Common expressions:

- `lower(email)` for case-insensitive equality
- `(payload->>'type')` for a JSON field
- `date_trunc('day', created_at)` for day buckets
- `(f(col))` when `f` is `IMMUTABLE`

The expression must be marked immutable for `CREATE INDEX`. `now()` is not immutable. You cannot index `now()`.

```sql
CREATE INDEX events_type_idx ON events ((payload->>'type'));
```

Parentheses are required around a non-column expression.

A unique expression index enforces uniqueness on the expression:

```sql
CREATE UNIQUE INDEX users_email_lower_key ON users (lower(email));
```

Expression indexes add write cost. The server computes the expression on each write.

Prefer a stored generated column plus a plain index when the expression is long and you also `SELECT` it (topic 13).

### Questions

#### Theoretical questions

1. What does an expression index store?
2. Why must the query repeat `lower(email)`?
3. Why can you not index `now()`?
4. Why do you need extra parentheses in `CREATE INDEX ON t ((expr))`?
5. What does a unique expression index enforce?

#### Easy practical tasks

1. Create `lower(email)` index. Query with `lower(email) = ...`.
2. Query with `email = ...` and write whether that index is used.
3. Create an index on `(payload->>'type')` if you have `jsonb`.
4. Run `\d` and copy the index definition.

#### Medium practical tasks

1. Insert two emails that differ only by case. With a unique `lower(email)` index, show the second insert fails.
2. Compare `ILIKE 'Foo'` with `lower(col) = 'foo'` and `EXPLAIN`.
3. Try `CREATE INDEX ON t (now())` or `DEFAULT now()` in an index expression. Record the error.

#### Advanced practical tasks

1. Index `date_trunc('day', created_at)`. Query a single day with the same `date_trunc`. Show `EXPLAIN`.
2. Read immutability rules. Write why `to_tsvector('english', body)` can be indexed but a `to_tsvector` call that uses a column for the config name might not.

---

## Unique and composite indexes

A **unique index** forbids duplicate keys. `PRIMARY KEY` and `UNIQUE` constraints create unique indexes. You can also write:

```sql
CREATE UNIQUE INDEX items_sku_key ON items (sku);
```

Prefer a `UNIQUE` constraint when the rule is a data rule. The constraint name appears in errors. The unique index is the implementation.

A **composite** (multicolumn) index stores more than one column:

```sql
CREATE INDEX orders_cust_created_idx
ON orders (customer_id, created_at);
```

Rule of thumb (B-tree):

- Equality on the leftmost columns, then range on the next column, matches well.
- `WHERE customer_id = 1 AND created_at > '2026-01-01'` can use `(customer_id, created_at)`.
- `WHERE created_at > '2026-01-01'` alone often cannot use that index well.

Column order matters. Put equality columns first. Put the range column after.

A composite unique index enforces uniqueness of the pair:

```sql
CREATE UNIQUE INDEX lines_order_sku_key
ON order_lines (order_id, sku);
```

`INCLUDE` adds non-key columns to a B-tree for index-only scans (topic 19). Know the word. Use it when you measure a heap fetch cost.

Do not duplicate a constraint index with a second handmade unique index on the same columns.

Two single-column indexes are not the same as one composite index. The planner can combine indexes with a bitmap scan (topic 9). A composite index is still better for a leading-equality plus range pattern.

### Questions

#### Theoretical questions

1. What does a unique index forbid?
2. Why prefer a `UNIQUE` constraint over a bare unique index for a business rule?
3. What is a composite index?
4. Why does column order matter in a B-tree?
5. Are two single-column indexes the same as one composite index?

#### Easy practical tasks

1. Create a unique index on a `sku` column. Insert a duplicate. Record the error.
2. Create `(customer_id, created_at)`. Run `\d`.
3. Write the leftmost-prefix rule in four sentences.
4. List indexes that your primary keys already created.

#### Medium practical tasks

1. Query with equality on the first column and a range on the second. Run `EXPLAIN`.
2. Query with a range only on the second column. Run `EXPLAIN`. Compare.
3. Create a composite unique pair. Insert two rows that share one column but not both.

#### Advanced practical tasks

1. Compare one composite index with two single-column indexes on the same filters. Use `EXPLAIN (ANALYZE)` on a few thousand rows.
2. Read `INCLUDE` in `CREATE INDEX` docs. Write one example and when topic 19 would use it.

---

## `CREATE INDEX CONCURRENTLY`

`CREATE INDEX` takes a lock that blocks writes on the table for the build. On a busy table that lock is a problem.

`CREATE INDEX CONCURRENTLY` builds the index without a long write lock. Reads and writes continue. The build scans the table twice. It cannot run inside a transaction block.

```sql
CREATE INDEX CONCURRENTLY items_name_idx ON items (name);
```

Rules:

- You cannot put `CONCURRENTLY` in `BEGIN ... COMMIT`.
- The build can fail. A failed concurrent index is invalid. `\d` shows `INVALID`.
- Drop an invalid index with `DROP INDEX CONCURRENTLY` and try again.
- Concurrent build uses more time and more resources than a normal build.

```sql
DROP INDEX CONCURRENTLY IF EXISTS items_name_idx;
```

`CREATE UNIQUE INDEX CONCURRENTLY` also exists. Unique checks run during the build. Concurrent writes must not violate uniqueness.

Use `CONCURRENTLY` in production on large, live tables. Use a normal `CREATE INDEX` in empty databases, migrations that already lock, and sessions that allow a write pause.

`REINDEX CONCURRENTLY` is the matching rebuild command (next section).

In PostgreSQL 16 and 17, concurrent index create is a standard operations tool. Test it on a copy. Watch load.

### Questions

#### Theoretical questions

1. What lock problem does `CONCURRENTLY` avoid?
2. Can you run `CREATE INDEX CONCURRENTLY` inside a transaction?
3. What is an invalid index?
4. How do you remove an invalid index?
5. When is a normal `CREATE INDEX` acceptable?

#### Easy practical tasks

1. Create a small table. Build a normal index. Then drop it and build with `CONCURRENTLY`.
2. Try `BEGIN; CREATE INDEX CONCURRENTLY ...`. Record the error. `ROLLBACK`.
3. Find `INVALID` in the docs or a failed test. Write what `\d` would show.
4. Write four sentences: lock, two scans, no transaction, invalid.

#### Medium practical tasks

1. Build a concurrent unique index. Insert a duplicate from a second session during the build if you can; otherwise write the expected error.
2. Compare wall time of `CREATE INDEX` and `CREATE INDEX CONCURRENTLY` on a table with 100000 rows.
3. Read the official caveats for concurrent index creation. List three.

#### Advanced practical tasks

1. Write a migration checklist: pre-check, `CONCURRENTLY`, validate `\d`, analyze, fallback `DROP INDEX CONCURRENTLY`.
2. Monitor `pg_stat_progress_create_index` during a long concurrent build. Write the columns that you saw.

---

## `REINDEX`

Indexes can become bloated or corrupt. `REINDEX` rebuilds an index from the table.

```sql
REINDEX INDEX items_name_idx;
REINDEX TABLE items;
REINDEX SCHEMA public;
REINDEX DATABASE current_database();
```

`REINDEX TABLE` rebuilds all indexes of the table, including the primary key index.

A normal `REINDEX` locks writes. `REINDEX INDEX CONCURRENTLY` and `REINDEX TABLE CONCURRENTLY` rebuild without a long write lock. Concurrent reindex cannot run in a transaction. It cannot process some system catalogs in all forms. Read the docs for the form that you use.

```sql
REINDEX INDEX CONCURRENTLY items_name_idx;
```

When to reindex:

- you measured index bloat (topic 11)
- you changed `fillfactor` or storage parameters
- the docs or a support process tell you to rebuild after corruption
- a concurrent create left an invalid index (drop and create again; reindex is a different path)

`REINDEX` is not a weekly habit for every table. Autovacuum and HOT updates (topic 19) keep many B-trees healthy. Measure first.

`VACUUM FULL` also rewrites a table and its indexes. That command is heavier. Topic 11 covers it.

PostgreSQL 17 improved some vacuum and reindex internals. The SQL you write stays the same.

### Questions

#### Theoretical questions

1. What does `REINDEX INDEX` rebuild?
2. What does `REINDEX TABLE` rebuild?
3. Why does `REINDEX CONCURRENTLY` exist?
4. When should you reindex?
5. How is `REINDEX` different from `VACUUM FULL`?

#### Easy practical tasks

1. `REINDEX INDEX` on a practice index. Show that `\d` still lists it.
2. `REINDEX TABLE` on a small table.
3. Try `BEGIN; REINDEX INDEX CONCURRENTLY ...`. Record the error.
4. Write the four `REINDEX` targets from this section.

#### Medium practical tasks

1. Compare `REINDEX INDEX` and `REINDEX INDEX CONCURRENTLY` times on a larger practice index.
2. Read `pg_index` for `indisvalid` and `indisready` after a normal reindex.
3. List indexes on a table from `pg_indexes` and reindex one by name.

#### Advanced practical tasks

1. Read the `REINDEX` reference for restrictions on system catalogs and concurrent mode. Write three restrictions.
2. Write an operations note: when you choose `REINDEX CONCURRENTLY` versus `VACUUM FULL` versus drop and `CREATE INDEX CONCURRENTLY`.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Pick an index class for `integer` range, `jsonb` containment, `tsvector` search, and a large append-only timestamp. Give one reason each.
2. How do partial indexes and expression indexes reduce unused index entries in different ways?
3. Why does leftmost prefix order matter for a composite B-tree but not for a GIN on `jsonb`?
4. What operational difference do `CONCURRENTLY` builds and reindexes share?
5. When is no extra index the correct design?

#### Easy practical tasks

1. On one schema, create a B-tree, a partial B-tree, an expression B-tree, and a GIN (if you have `jsonb` or an array).
2. Write a cheat sheet: five index classes, partial, expression, unique, composite, concurrently, reindex.
3. Show `\d` for a table and label each index as constraint-backed or extra.
4. Run `EXPLAIN` on two queries: one that uses an index and one that does not.

#### Medium practical tasks

1. Design indexes for `orders(customer_id, created_at, status)` given two real queries. Create them. Prove with `EXPLAIN`.
2. Break a concurrent create on purpose (transaction) and write the recovery command.
3. Document team rules: who may create indexes in production, and when `CONCURRENTLY` is required.

#### Advanced practical tasks

1. Load 100000 rows. Compare sizes of B-tree, BRIN (on a sequential column), and GIN (on tags or jsonb). Write `pg_relation_size` for each.
2. Read "Index Types" and "Building Indexes Concurrently". Write a one-page map from query pattern to index type.
