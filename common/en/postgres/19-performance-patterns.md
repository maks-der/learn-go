# 19. Performance Patterns

## Description

This topic shows common performance patterns in PostgreSQL 16 and PostgreSQL 17. You learn keyset pagination, partial indexes, covering indexes with `INCLUDE`, heap-only tuple (HOT) updates, why a high `max_connections` thrashes, and how you use read replicas for reporting.

Complete topic 18 first. You can read `EXPLAIN` and memory settings. Complete this topic before you study window functions and other advanced SQL.

Use one term for each concept. Keyset pagination is a seek on the last seen sort key. A partial index is an index with a `WHERE` clause. A covering index includes extra columns so that an index-only scan can avoid the heap. HOT is a heap-only tuple update. Thrashing here is too many backends for the RAM and CPU. A read replica is a hot stand-by that you use for reads. This handbook targets PostgreSQL 16 and PostgreSQL 17.

---

## Keyset pagination

`LIMIT` plus `OFFSET` skips rows each time. Page 100 with `OFFSET 1000` still reads and discards 1000 rows. That cost grows with the page number.

Keyset pagination (also called seek method) filters on the last row that the client already saw:

```sql
SELECT id, created_at, title
FROM shop.posts
WHERE (created_at, id) < ('2026-09-01 12:00:00+00', 8400)
ORDER BY created_at DESC, id DESC
LIMIT 21;
```

The client uses the last `created_at` and `id` of the previous page. A unique tie-breaker (`id`) is required. Without it, two rows with the same timestamp can be skipped or repeated.

A matching B-tree is:

```sql
CREATE INDEX posts_created_id_idx ON shop.posts (created_at DESC, id DESC);
```

`EXPLAIN` must show an index scan that stops after `LIMIT`, not a full sort of the table.

`OFFSET` is acceptable for small admin pages. It is the wrong default for infinite scroll on a large table.

Do not mix `OFFSET` with a keyset in a way that the client cannot explain. Do not paginate on a non-unique column alone. Do not use `LIMIT` / `OFFSET` to hide a missing index.

Topic 3 introduced `LIMIT`, `OFFSET`, and `FETCH`. This section is the large-table pattern.

### Questions

#### Theoretical questions

1. Why does a large `OFFSET` get slower?
2. What extra column do you add when `created_at` is not unique?
3. What index matches a keyset `ORDER BY created_at DESC, id DESC`?
4. When is `OFFSET` still acceptable?
5. What must `EXPLAIN` show for a healthy keyset query?

#### Easy practical tasks

1. Select the first page with `ORDER BY created_at DESC, id DESC LIMIT 5`.
2. Write the second-page query with a keyset from the last row of page one.
3. Write the same second page with `OFFSET 5`. Compare `EXPLAIN`.
4. Create the composite index. `\d shop.posts`.

#### Medium practical tasks

1. Load 20000 rows. Compare `EXPLAIN ANALYZE` for `OFFSET 15000` and a keyset at the same position.
2. Omit `id` from the `WHERE` tuple. Write a case with two equal timestamps that can skip a row.
3. Rewrite the keyset with `FETCH FIRST 21 ROWS ONLY`. Confirm the plan still uses the index.

#### Advanced practical tasks

1. Implement next and previous pages (two directions). Write both predicates.
2. Read "Use the Index, Luke" on keyset pagination (PostgreSQL). Write one extra rule that this section did not include.

---

## Partial indexes for hot filters

A partial index stores only rows that match a predicate. It is smaller. It stays useful when the predicate matches the query.

```sql
CREATE INDEX orders_open_idx
    ON shop.orders (customer_id)
    WHERE status = 'open';
```

A query that uses the index:

```sql
SELECT id, total
FROM shop.orders
WHERE status = 'open' AND customer_id = 42;
```

A query with `status = 'closed'` cannot use `orders_open_idx`. A query that omits `status` cannot use it either, unless the planner can prove the predicate.

Partial indexes fit:

- a small hot subset (`status = 'open'`, `deleted_at IS NULL`)
- a rare flag that you search (`is_error = true`)
- a tenant or region that you query often

Topic 8 introduced partial indexes. This section is the "hot filter" rule: index the slice that production actually filters.

`EXPLAIN` must mention the partial index. If the planner uses a sequential scan, the predicate does not match or the table is tiny.

Do not create a partial index for every status value without a measurement. Do not use a partial unique index as a substitute for a real business rule unless that rule is "unique among live rows":

```sql
CREATE UNIQUE INDEX users_email_live_uidx
    ON shop.users (email)
    WHERE deleted_at IS NULL;
```

### Questions

#### Theoretical questions

1. What rows does a partial index store?
2. When can a query use `WHERE status = 'open'` index?
3. Why is a partial index often smaller than a full index?
4. What unique pattern does `WHERE deleted_at IS NULL` allow?
5. Why might a tiny table ignore a partial index?

#### Easy practical tasks

1. Create a partial index on a status or a boolean. Run `\d`.
2. `EXPLAIN` a query that matches the predicate.
3. `EXPLAIN` a query that uses a different status. Write whether the partial index appears.
4. Compare `pg_relation_size` of a full index and a partial index on the same column if you create both.

#### Medium practical tasks

1. Load mixed statuses. Prove with `EXPLAIN ANALYZE` that the open-only query is cheaper with the partial index.
2. Create a unique partial index on `email` where `deleted_at IS NULL`. Insert two live rows with the same email. Record the error. Insert a live and a deleted pair if your rule allows it.
3. Write three production filters from an app you know. Mark which ones deserve a partial index.

#### Advanced practical tasks

1. Combine a partial index with an expression (`WHERE` plus `(lower(email))`). Prove with `EXPLAIN`.
2. Read partial-index examples in the 17 docs. Add one example that this section did not show.

---

## Covering indexes (`INCLUDE`)

A covering index is a B-tree whose leaf can satisfy the query without a heap fetch. PostgreSQL calls the extra columns **included columns**. They are not sort keys.

```sql
CREATE INDEX items_name_cover_idx
    ON shop.items (name)
    INCLUDE (qty, price);
```

A query that can become an index-only scan:

```sql
SELECT name, qty, price
FROM shop.items
WHERE name = 'nail';
```

`name` is the search key. `qty` and `price` are payload. `WHERE` and `ORDER BY` still use only the key columns (and their order).

Topic 8 and topic 9 previewed `INCLUDE` and index-only scans. Index-only scans also need a visibility map that says the page is all-visible. `VACUUM` helps (topic 11). If `EXPLAIN` shows "Heap Fetches" above zero, the covering index is not enough or vacuum is behind.

Included columns make the index larger. They slow writes. They can break HOT if you include a column that updates often (next section).

Do not `INCLUDE` every column to copy the table into an index. Do not expect `INCLUDE` to help `SELECT *` on a wide row. Do not add `INCLUDE` before you see heap fetches in `EXPLAIN (ANALYZE)`.

Official: [https://www.postgresql.org/docs/17/indexes-index-only-scans.html](https://www.postgresql.org/docs/17/indexes-index-only-scans.html).

### Questions

#### Theoretical questions

1. What is an included column?
2. Can you `ORDER BY` an included column using that index as a sort key?
3. What extra condition does an index-only scan need besides the index?
4. How do you see leftover heap fetches?
5. Why can `INCLUDE` of a hot column hurt updates?

#### Easy practical tasks

1. Create a B-tree with `INCLUDE`. `\d shop.items`.
2. `EXPLAIN` a query that selects only the key and included columns.
3. `EXPLAIN` `SELECT *` on the same filter. Write the difference.
4. Write four sentences: key, include, index-only, vacuum.

#### Medium practical tasks

1. `EXPLAIN (ANALYZE)` before and after `INCLUDE`. Compare heap fetches after `VACUUM`.
2. Compare index size with and without `INCLUDE`.
3. Try `ORDER BY` on an included column. Write whether the plan uses an extra sort.

#### Advanced practical tasks

1. Build a covering index for a keyset page query that returns three display columns. Prove index-only or write why heap fetches remain.
2. Read index-only scans in the 17 docs. Write the role of the visibility map in two sentences.

---

## Heap-only tuples (HOT) updates (awareness)

A normal `UPDATE` writes a new heap tuple and new index entries for every index that includes a changed column. A **HOT** update writes the new tuple on the same heap page and skips index updates when:

1. no indexed column changes (key columns and `INCLUDE` columns count as indexed)
2. the new tuple fits on the same page (free space; `fillfactor` helps)

HOT reduces index bloat. Topic 11 introduced HOT and fillfactor. This section is the application rule: do not index columns that you update on every request if you can avoid it.

```sql
-- last_seen changes often; a B-tree on last_seen blocks HOT for that update
UPDATE shop.users SET last_seen = now() WHERE id = 1;
```

If `last_seen` has an index, those updates are not HOT. If you only filter `last_seen` in a rare report, drop that index or use a partial index that excludes the hot path.

`SELECT n_tup_hot_upd, n_tup_upd FROM pg_stat_user_tables WHERE relname = 'users';` shows HOT versus all updates (names can vary slightly; read `pg_stat_all_tables` columns on your version).

`EXPLAIN` does not say "HOT". You infer HOT from stats and from index design.

Do not add covering columns that you update as often as you read them. Do not set `fillfactor = 50` on every table. Measure a hot table first.

Official page: [https://www.postgresql.org/docs/17/storage-hot.html](https://www.postgresql.org/docs/17/storage-hot.html) (same path under `/docs/16/`).

### Questions

#### Theoretical questions

1. What two conditions does a HOT update need?
2. Does an `INCLUDE` column count as an indexed column for HOT?
3. How does `fillfactor` help HOT?
4. Where do you read `n_tup_hot_upd`?
5. Why can an index on `last_seen` hurt a busy `UPDATE`?

#### Easy practical tasks

1. Query `pg_stat_user_tables` for `n_tup_upd` and the HOT column on one table.
2. List indexes on a table that you update often. Mark columns that those updates change.
3. Write four sentences: new tuple, same page, skip index, fillfactor.
4. Open the HOT storage page or the 17 storage chapter. Write one official sentence in your own words.

#### Medium practical tasks

1. Update a non-indexed column many times. Record HOT updates. Add an index on that column. Repeat. Compare.
2. `INCLUDE` a hot column in a covering index. Repeat the update test. Write the effect.
3. Set `fillfactor = 80` on a copy table. Repeat a tight update loop. Compare HOT counts with fillfactor 100.

#### Advanced practical tasks

1. Write an index policy: which columns may be indexed on a table with frequent updates of `last_seen` and rare searches of `email`.
2. Read HOT and summary indexes (BRIN) in the docs. Write why BRIN is not a HOT substitute.

---

## Avoiding thrashing with too-high `max_connections`

`max_connections` is the cap on server backends. Each backend uses memory (`work_mem` per step, plus per-backend overhead). The operating system spends more time on context switches when hundreds of backends are active.

Thrashing here means: too many backends, too little RAM, heavy swap or cache eviction, and every query slower than a smaller pool.

Topic 17 said to use PgBouncer. Topic 18 said `work_mem` times connections can exceed RAM. The pattern:

1. set `max_connections` to a number the host can run (often 100 to 200 on a mid host, not 2000)
2. put PgBouncer or a driver pool in front
3. size the pool to a few dozen server connections for OLTP
4. use a separate small pool for admin and migrations

```sql
SHOW max_connections;
SELECT count(*) FROM pg_stat_activity;
```

Cloud instances often default `max_connections` by RAM class. Do not raise the default to "remove too many connections" errors without a pool. That error is a signal to pool, not to multiply backends.

Do not set `max_connections` to 1000 on a 4 GB host. Do not open a session per HTTP request (topic 17). Do not forget that replicas have their own `max_connections`.

### Questions

#### Theoretical questions

1. What does `max_connections` cap?
2. How does `work_mem` change the risk of a high cap?
3. What error is a signal to add a pool rather than raise the cap?
4. Why is 2000 backends a poor default on a small host?
5. Do replicas share the primary `max_connections` value?

#### Easy practical tasks

1. `SHOW max_connections;` Count `pg_stat_activity`.
2. Write a table: host RAM, suggested starting cap, pool in front yes/no.
3. Read your cloud or Docker default for `max_connections`.
4. Write four sentences: backend, pool, RAM, context switch.

#### Medium practical tasks

1. Compute worst-case RAM: `max_connections * work_mem` plus `shared_buffers`. Write whether the host can hold it.
2. Plan pool sizes: API 40, reports 10, admin 5, `max_connections` 80. Explain the headroom.
3. From `pg_stat_activity`, group by `application_name` and `state`. Write who holds idle sessions.

#### Advanced practical tasks

1. In a lab, lower `max_connections` to a small number (restart). Hit the limit with many `psql` sessions. Then add PgBouncer or stop extra sessions. Document the error text.
2. Read `superuser_reserved_connections` in the 17 docs. Write why you keep a reserve.

---

## Read replicas for reporting

A reporting query that sorts a large table can compete with OLTP on the primary. A hot stand-by (topic 16) can run read-only SQL. You point the report at the replica.

Rules:

- the report must accept lag (seconds or more)
- the report must accept query cancel on recovery conflicts
- the report must not write
- you still need indexes on the primary (they replicate physically)
- `hot_standby_feedback` can protect the report and delay vacuum on the primary

```sql
-- on replica
SELECT pg_is_in_recovery();
SELECT now() - pg_last_xact_replay_timestamp() AS replica_lag;
```

`pg_last_xact_replay_timestamp()` can be null when there is no replay yet. Treat lag as a display, not a transaction guarantee.

Logical subscribers (topic 16) can also serve reports if the publication includes the tables. They can have different indexes. They are a different product path.

Do not run the nightly `VACUUM FULL` on the replica as a substitute for primary maintenance. Do not use the replica for the next key of an order number. Do not ignore lag in a dashboard that looks "live".

Topic 21 monitors replication lag. This section is the performance reason to have a replica.

### Questions

#### Theoretical questions

1. What problem does a reporting replica solve on the primary?
2. What two replica limits must a report accept?
3. Do physical replica indexes differ from the primary?
4. What trade does `hot_standby_feedback` make?
5. Why must an order-number generator stay on the primary?

#### Easy practical tasks

1. If you have a replica, run `pg_is_in_recovery()` and a lag query. If not, write the two SQL statements for later.
2. Write three report types that fit a replica and two that do not.
3. Draw: OLTP app → primary, report job → replica.
4. `SHOW hot_standby_feedback;` on primary or replica if available.

#### Medium practical tasks

1. Run a heavy `GROUP BY` on a lab primary. Then describe how you would move it to a replica (connection string only).
2. Compare physical replica versus logical subscriber for a report that needs an extra index. Write three sentences.
3. Document a maximum lag that your fake business accepts. Write what the job does when lag is larger.

#### Advanced practical tasks

1. Build a report connection policy: URI, `application_name`, read-only role, lag check before the report starts.
2. Read standby conflict settings in the 17 docs. Write which setting you change first when reports cancel.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do keyset pagination and a covering index work together on a feed query?
2. How can a partial index and HOT conflict or cooperate on a `status` plus `updated_at` table?
3. Why is a high `max_connections` a performance bug even when the SQL is indexed?
4. When do you add a replica instead of another index on the primary?
5. A teammate pages with `OFFSET 50000`, indexes every column, and sets `max_connections = 2000`. Which facts from this topic do you use in the reply?

#### Easy practical tasks

1. On one table, create a composite sort index and a partial index. Run one keyset query and one hot-filter query. Save `EXPLAIN`.
2. Write a cheat sheet: keyset, partial, `INCLUDE`, HOT, `max_connections`, replica lag.
3. `SHOW max_connections;` Query HOT stats for one table.
4. Write a keyset `WHERE` for `(created_at, id)` descending.

#### Medium practical tasks

1. Load 50000 rows. Compare `OFFSET` and keyset at a deep page. Add `INCLUDE` for the selected columns. Compare heap fetches after `VACUUM`.
2. Measure HOT before and after you drop an index on a hot column.
3. Write a connection map: API pool, report replica, admin direct.

#### Advanced practical tasks

1. Design indexes for `posts(id, created_at, author_id, status, title)` with a public feed, an author draft filter, and frequent `status` updates. Justify HOT and partial choices.
2. Map this topic to official chapters: `LIMIT`/`OFFSET`, partial indexes, index-only scans, HOT, connections, hot stand-by. Add one planner fact that this topic did not include.
