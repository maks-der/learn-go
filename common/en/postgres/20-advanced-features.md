# 20. Advanced Features

## Description

This topic shows advanced SQL and extensions in PostgreSQL 16 and PostgreSQL 17. You learn window functions, `FILTER` on aggregates, range types and exclusion constraints, optional PostGIS, foreign data wrappers, parallel query, and logical decoding for change data capture.

Complete topic 19 first. You can read plans and indexes. Complete this topic before you harden a production cluster.

Use one term for each concept. A window function computes a value from a related set of rows without collapsing the result. `FILTER` restricts which rows enter one aggregate. A range type stores a lower and upper bound. An exclusion constraint prevents overlapping ranges. A foreign data wrapper (FDW) reads a remote source as a table. Logical decoding turns WAL into a change stream. This handbook targets PostgreSQL 16 and PostgreSQL 17.

---

## Window functions

A window function uses `OVER (...)`. The result set keeps one row per input row. `GROUP BY` collapses rows. Use a window when you need a rank, a running sum, or a neighbor row.

```sql
SELECT
    id,
    customer_id,
    total,
    sum(total) OVER (PARTITION BY customer_id) AS customer_sum,
    rank() OVER (PARTITION BY customer_id ORDER BY total DESC) AS rnk,
    lag(total) OVER (PARTITION BY customer_id ORDER BY id) AS prev_total
FROM shop.orders;
```

Clauses inside `OVER`:

- `PARTITION BY` — reset the window
- `ORDER BY` — order inside the partition (required for `rank`, `lag`, running sums)
- frame (`ROWS BETWEEN ...` or `RANGE BETWEEN ...`) — which neighbors count

Default frame for many aggregates with `ORDER BY` is `RANGE BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW`. That default is a running aggregate. If you wanted the sum of the whole partition, use `sum(total) OVER (PARTITION BY customer_id)` without `ORDER BY`, or set the frame to `ROWS BETWEEN UNBOUNDED PRECEDING AND UNBOUNDED FOLLOWING`.

Common functions: `row_number()`, `rank()`, `dense_rank()`, `lag()`, `lead()`, `first_value()`, `sum()`, `avg()`, `count()`.

Topic 6 covered `GROUP BY` and CTEs. You can put a window in a CTE and then filter `WHERE rnk = 1`. `DISTINCT ON` (topic 6) can replace a `row_number()` filter in simple "first row per group" cases. Prefer the form that you can test.

Do not use a window when a join to a grouped subquery is clearer and the planner is simpler. Do not forget `ORDER BY` inside `OVER` for `lag` and `rank`.

Official: [https://www.postgresql.org/docs/17/tutorial-window.html](https://www.postgresql.org/docs/17/tutorial-window.html), [https://www.postgresql.org/docs/17/functions-window.html](https://www.postgresql.org/docs/17/functions-window.html).

### Questions

#### Theoretical questions

1. How does a window differ from `GROUP BY` in the number of result rows?
2. What does `PARTITION BY` reset?
3. Why does `lag` need `ORDER BY` inside `OVER`?
4. What is the default frame for `sum(...) OVER (ORDER BY ...)`?
5. When can `DISTINCT ON` replace `row_number()`?

#### Easy practical tasks

1. Compute `row_number() OVER (ORDER BY id)` on a table.
2. Compute `sum(total) OVER (PARTITION BY customer_id)` and compare with a `GROUP BY` query.
3. Use `lag(id) OVER (ORDER BY id)`.
4. Write four sentences: window, partition, frame, `GROUP BY`.

#### Medium practical tasks

1. Rank orders per customer by `total`. Filter `rnk <= 3` in an outer query.
2. Compare a running `sum(total) OVER (ORDER BY id)` with a full-partition `sum` without `ORDER BY`. Write two numeric examples.
3. `EXPLAIN` a window query. Write the node name that computes the window.

#### Advanced practical tasks

1. Write a gaps-and-islands or sessionization query with `lag` and a cumulative sum. Use a small event table.
2. Read window-frame clauses in the 17 docs. Write one `ROWS` example and one `RANGE` example. State the difference.

---

## `FILTER` on aggregates

`FILTER` limits the rows that one aggregate sees. Other columns in the same `SELECT` can use a different filter.

```sql
SELECT
    customer_id,
    count(*) AS orders,
    count(*) FILTER (WHERE total >= 100) AS large_orders,
    sum(total) FILTER (WHERE created_on >= DATE '2026-01-01') AS ytd
FROM shop.orders
GROUP BY customer_id;
```

Without `FILTER` you write `sum(CASE WHEN total >= 100 THEN 1 ELSE 0 END)`. `FILTER` is shorter and states the intent.

`FILTER` works with `GROUP BY` and with window aggregates:

```sql
SELECT
    id,
    count(*) FILTER (WHERE total > 0)
        OVER (PARTITION BY customer_id) AS positive_count
FROM shop.orders;
```

`FILTER` is not a `WHERE` for the whole query. Rows that fail the filter still exist for other aggregates in the same row.

Do not replace a selective `WHERE` that removes most rows with a `FILTER` on a huge scan. `WHERE` still reduces the input set first. Do not mix `FILTER` with `DISTINCT` inside an aggregate until you read the docs for that combination.

Official: [https://www.postgresql.org/docs/17/sql-expressions.html#SYNTAX-AGGREGATES](https://www.postgresql.org/docs/17/sql-expressions.html#SYNTAX-AGGREGATES).

### Questions

#### Theoretical questions

1. What rows does `count(*) FILTER (WHERE total >= 100)` count?
2. How is `FILTER` different from `WHERE` on the query?
3. Can two aggregates in one `SELECT` use two filters?
4. Can a window aggregate use `FILTER`?
5. What older pattern does `FILTER` replace?

#### Easy practical tasks

1. Compute `count(*)` and `count(*) FILTER (WHERE qty = 0)` on `shop.items`.
2. Rewrite a `CASE` sum as `FILTER`.
3. Add `GROUP BY` and two filters.
4. Write four sentences: filter, group, window, `WHERE`.

#### Medium practical tasks

1. Build a report with four filtered sums in one pass (`GROUP BY` customer).
2. Compare `EXPLAIN` for `WHERE created_on >= $d` versus a `FILTER` on the same date with no `WHERE`. Write which input set is smaller.
3. Combine `FILTER` with `avg` and `NULL` rows. Write how `NULL` inputs behave (read the aggregate docs).

#### Advanced practical tasks

1. Write a single query that uses `FILTER` on a group aggregate and `FILTER` on a window. Explain each number.
2. Read the 17 aggregate-expression page. Write whether `FILTER` works with `jsonb_agg` and give an example.

---

## Range types and exclusion constraints

A range type stores a lower bound and an upper bound. Built-in types include `int4range`, `int8range`, `numrange`, `tsrange`, `tstzrange`, and `daterange`. Bounds can be inclusive or exclusive. A range can be empty. A range can be unbounded.

```sql
SELECT daterange('2026-01-01', '2026-02-01', '[)');
SELECT tstzrange('2026-01-01 00:00+00', '2026-01-01 12:00+00')
       && tstzrange('2026-01-01 11:00+00', '2026-01-01 18:00+00');
```

`&&` means overlaps. `@>` means contains. Topic 5 showed a high-level `EXCLUDE`. The usual booking constraint is:

```sql
CREATE TABLE shop.room_bookings (
    room_id bigint NOT NULL,
    during tstzrange NOT NULL,
    EXCLUDE USING gist (room_id WITH =, during WITH &&)
);
```

Two rows with the same `room_id` must not have overlapping `during`. The exclusion constraint needs a GiST index. PostgreSQL creates that index for the constraint.

Install `btree_gist` when you need `=` on a scalar next to `&&` on a range in one `EXCLUDE`:

```sql
CREATE EXTENSION IF NOT EXISTS btree_gist;
```

Topic 5 said to know `EXCLUDE`. This section is the working pattern. Prefer a range column plus `EXCLUDE` over two timestamp columns plus a tangle of `CHECK` clauses.

Do not store overlapping bookings and then filter them only in the application. Do not use `tsrange` if you need time zones; use `tstzrange`. Do not forget `btree_gist` when the docs require it.

Official: [https://www.postgresql.org/docs/17/rangetypes.html](https://www.postgresql.org/docs/17/rangetypes.html), [https://www.postgresql.org/docs/17/btree-gist.html](https://www.postgresql.org/docs/17/btree-gist.html).

### Questions

#### Theoretical questions

1. What does a `daterange` store?
2. What does `&&` mean for two ranges?
3. What does `EXCLUDE USING gist (room_id WITH =, during WITH &&)` forbid?
4. Why do you often need `btree_gist`?
5. When do you pick `tstzrange` over `tsrange`?

#### Easy practical tasks

1. Select a `daterange` for one month with `'[)'`.
2. Test `&&` on two `tstzrange` values that overlap and two that do not.
3. Create `btree_gist` if allowed. Create `room_bookings` with the exclusion constraint.
4. Insert two overlapping rows for the same room. Record the error.

#### Medium practical tasks

1. Insert a legal second booking that touches at the bound (`[)`). Write whether the exclusive upper bound allows the next start.
2. Query bookings that overlap a given window with `during && $r`.
3. Compare this design with two columns `starts_at` and `ends_at` plus a `CHECK (starts_at < ends_at)` only. Write what `CHECK` misses.

#### Advanced practical tasks

1. Add a `WHERE` predicate to a partial exclusion if the 16/17 docs allow it for your case, or document the limit. Test cancelled rows if you use a status.
2. Read range functions (`lower`, `upper`, `isempty`). Write a query that lists empty or unbounded ranges in a lab table.

---

## PostGIS (geospatial, optional)

PostGIS is an extension that adds geometry and geography types, spatial indexes (GiST), and spatial functions. It is optional. Many PostgreSQL 16 and 17 installs do not load it. Cloud vendors often offer it as a named extension.

```sql
CREATE EXTENSION IF NOT EXISTS postgis;
SELECT PostGIS_Version();
```

Typical objects:

- `geometry` — planar coordinates in a spatial reference system
- `geography` — measurements on a spheroid (meters)
- `ST_DWithin`, `ST_Contains`, `ST_Distance`
- `GIST` index on a geometry column

```sql
CREATE TABLE shop.stores (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name text NOT NULL,
    loc geography(Point, 4326) NOT NULL
);

CREATE INDEX stores_loc_gix ON shop.stores USING gist (loc);

SELECT name
FROM shop.stores
WHERE ST_DWithin(loc, ST_MakePoint(-0.12, 51.50)::geography, 1000);
```

If `CREATE EXTENSION postgis` fails, skip the practical tasks that need it. Write that the package is absent. You still must know that spatial search is not a `numeric` pair plus a B-tree.

Do not store latitude and longitude as `text`. Do not run `ST_Distance` on every row without a GiST index when the table is large. Do not treat PostGIS as required for the rest of this path.

Official: [https://postgis.net/documentation/](https://postgis.net/documentation/). Match the PostGIS version to your PostgreSQL 16 or 17.

### Questions

#### Theoretical questions

1. What types does PostGIS add?
2. Which index class is usual for spatial search?
3. What does `ST_DWithin` test?
4. Is PostGIS part of every PostgreSQL install?
5. Why is a B-tree on `lat` and `lon` a poor spatial index?

#### Easy practical tasks

1. Try `CREATE EXTENSION postgis`. Write success or the error.
2. If it loaded, run `SELECT PostGIS_Version();`.
3. Open the PostGIS documentation home. Write the current stable PostGIS version that the site lists.
4. Write four sentences: geometry, geography, GiST, optional.

#### Medium practical tasks

1. If PostGIS is present, create `shop.stores`, insert two points, query `ST_DWithin`.
2. `EXPLAIN` the `ST_DWithin` query with and without the GiST index (drop and recreate in a lab).
3. If PostGIS is absent, write a fallback schema (`lat numeric`, `lon numeric`) and list three queries you cannot do correctly.

#### Advanced practical tasks

1. Read `geometry` versus `geography` in PostGIS docs. Write when you pick each for a city-distance search.
2. List two extra PostGIS functions that this section did not name. Write one sentence each.

---

## Foreign data wrappers

A foreign data wrapper (FDW) lets you query a remote source as a foreign table. `postgres_fdw` is the common wrapper for another PostgreSQL server. Other wrappers exist for files and other databases. Quality varies.

```sql
CREATE EXTENSION postgres_fdw;

CREATE SERVER remote_shop
    FOREIGN DATA WRAPPER postgres_fdw
    OPTIONS (host 'remote.example', dbname 'shop', port '5432');

CREATE USER MAPPING FOR shop_app
    SERVER remote_shop
    OPTIONS (user 'remote_read', password 'secret');

CREATE FOREIGN TABLE shop.remote_items (
    id bigint,
    name text,
    qty int
)
SERVER remote_shop
OPTIONS (schema_name 'shop', table_name 'items');
```

`IMPORT FOREIGN SCHEMA` can create many foreign tables at once.

The planner can push some `WHERE` clauses to the remote server. A join of local and remote tables can pull many rows. `EXPLAIN` shows `Foreign Scan`.

FDW is not a substitute for replication when you need a local copy. It is a live query path. Network failure is a query failure. Transactions across local and remote have limits. Read `postgres_fdw` docs for the current 16/17 isolation behavior.

Do not put a password in `USER MAPPING` in a shared script. Do not query a huge remote table with `SELECT *` from a laptop. Do not use FDW as the only backup of the remote data.

Official: [https://www.postgresql.org/docs/17/postgres-fdw.html](https://www.postgresql.org/docs/17/postgres-fdw.html).

### Questions

#### Theoretical questions

1. What does an FDW add to a database?
2. What are `SERVER`, `USER MAPPING`, and `FOREIGN TABLE`?
3. What node name does `EXPLAIN` use for a foreign table?
4. Why is FDW not a backup?
5. What does `IMPORT FOREIGN SCHEMA` do?

#### Easy practical tasks

1. Read `\h CREATE SERVER` and `\h CREATE FOREIGN TABLE`.
2. Write the four SQL objects in the example in order.
3. Open the 17 `postgres_fdw` page. Write two `OPTIONS` keys for `CREATE SERVER`.
4. Write four sentences: wrapper, server, mapping, foreign table.

#### Medium practical tasks

1. If you have two databases, create `postgres_fdw` from one to the other. Select from the foreign table.
2. `EXPLAIN` a filtered query on the foreign table. Write whether the filter was pushed (remote SQL in the plan if shown).
3. Compare FDW with logical replication for a report. Write three differences.

#### Advanced practical tasks

1. Join a local table to a foreign table. `EXPLAIN`. Write a rule for when you copy the data instead.
2. Read `postgres_fdw` transaction and error notes in the 17 docs. Write two limits that an application must accept.

---

## Parallel query

PostgreSQL can use more than one worker for one query. You see `Gather`, `Parallel Seq Scan`, or `Parallel Hash` in `EXPLAIN`. Topic 9 said you may see `Parallel` in front of a node.

Settings (high-level):

- `max_parallel_workers_per_gather` — workers per gather node
- `max_parallel_workers` — cluster cap
- `max_worker_processes` — process slots
- `min_parallel_table_scan_size` / `min_parallel_index_scan_size` — size gates

```sql
SHOW max_parallel_workers_per_gather;
EXPLAIN (ANALYZE)
SELECT count(*) FROM shop.events;
```

Parallel query helps large scans and large hashes. It can hurt small OLTP lookups (start-up cost). The planner disables parallelism when the cost model says so, when the query is in a serial-only context, or when you `SET max_parallel_workers_per_gather = 0`.

Some operations do not parallelize. A function that is not `PARALLEL SAFE` can block workers. Mark functions correctly (`PARALLEL SAFE`, `RESTRICTED`, `UNSAFE`).

Do not raise `max_parallel_workers_per_gather` to 8 on a busy 4-core OLTP host without a measurement. Do not force parallel on a 10-row table. Do not forget that each worker uses `work_mem` (topic 18).

Official: [https://www.postgresql.org/docs/17/parallel-query.html](https://www.postgresql.org/docs/17/parallel-query.html).

### Questions

#### Theoretical questions

1. What `EXPLAIN` nodes hint that parallelism ran?
2. What does `max_parallel_workers_per_gather` cap?
3. Why can parallelism hurt a tiny query?
4. What function label can block parallel workers?
5. How does `work_mem` multiply with workers?

#### Easy practical tasks

1. `SHOW` the parallel settings in this section.
2. `EXPLAIN` `SELECT count(*)` on your largest lab table.
3. `SET max_parallel_workers_per_gather = 0;` Repeat `EXPLAIN`. Reset the setting.
4. Write four sentences: gather, worker, size gate, OLTP.

#### Medium practical tasks

1. Load 200000 rows. Compare `EXPLAIN ANALYZE` with gather on and off.
2. Read which operations cannot parallelize in the 17 docs. List three.
3. Create a `PARALLEL UNSAFE` SQL function and use it in a `WHERE`. Write whether the plan stays serial.

#### Advanced practical tasks

1. Write a policy: when reports may use 4 workers and when the API must stay at 0 or 1.
2. Read parallel safety of PL/pgSQL in the docs. Write the default for a PL/pgSQL function.

---

## Logical decoding / CDC

Logical decoding reads WAL and emits a change stream (insert, update, delete) for tables. Change data capture (CDC) tools use that stream. `wal_level` must be `logical`. A replication slot retains WAL for the consumer.

Built-in output plugin: `pgoutput` (used by logical replication). Contrib and external plugins exist (`test_decoding` for labs). Tools such as Debezium sit outside PostgreSQL and consume the slot.

```sql
SELECT * FROM pg_create_logical_replication_slot('cdc_lab', 'test_decoding');
-- make changes, then:
SELECT * FROM pg_logical_slot_get_changes('cdc_lab', NULL, NULL);
SELECT pg_drop_replication_slot('cdc_lab');
```

`test_decoding` may need the contrib package. If the plugin is absent, read the docs and skip the consume step.

CDC is not a backup. CDC is not a substitute for `LISTEN`. CDC consumers must handle schema changes, snapshots, and failures. A stuck slot fills the disk (topic 16 and topic 18).

PostgreSQL 16 can decode from a stand-by in supported setups. PostgreSQL 17 improves failover for logical slots. Learn slots and `wal_level` first.

Do not create a slot for a test and leave it. Do not parse WAL files by hand. Do not send the stream to an untrusted network without TLS and auth.

Official: [https://www.postgresql.org/docs/17/logicaldecoding.html](https://www.postgresql.org/docs/17/logicaldecoding.html).

### Questions

#### Theoretical questions

1. What does logical decoding read?
2. What `wal_level` do you need?
3. Why does a CDC slot retain WAL?
4. What built-in plugin does logical replication use?
5. How is CDC different from `LISTEN` / `NOTIFY`?

#### Easy practical tasks

1. `SHOW wal_level;`.
2. Query `pg_replication_slots`.
3. Open the 17 logical-decoding chapter. Write the purpose of an output plugin.
4. Write four sentences: WAL, slot, plugin, consumer.

#### Medium practical tasks

1. If `test_decoding` is available, create a slot, `INSERT` one row, get changes, drop the slot.
2. If the plugin is missing, write the `CREATE PUBLICATION` path as the supported CDC-adjacent feature (topic 16).
3. Document who owns slots in your team and how you alert on retained WAL size.

#### Advanced practical tasks

1. Read Debezium or another CDC tool's PostgreSQL page at a high level. Write how it uses a slot and a publication. Do not install malware; use official docs.
2. Compare 16 and 17 notes on logical decoding from a stand-by or slot failover. Write one sentence for operators.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. When do you pick a window, `FILTER`, or `GROUP BY` for the same report?
2. How do range exclusion, PostGIS GiST, and B-tree indexes differ in what they prevent or find?
3. When is FDW the wrong tool and logical decoding the right tool (or the reverse)?
4. How do parallel workers and `work_mem` change a report that you moved to a replica in topic 19?
5. A teammate wants PostGIS, three FDWs, and a permanent `test_decoding` slot on the OLTP primary. Which facts from this topic do you use in the reply?

#### Easy practical tasks

1. Run one window query, one `FILTER` aggregate, and one range overlap `SELECT` on lab data.
2. Write a cheat sheet: `OVER`, `FILTER`, `tstzrange`, `EXCLUDE`, PostGIS optional, FDW objects, parallel nodes, CDC slot.
3. `SHOW wal_level;` `SHOW max_parallel_workers_per_gather;`.
4. Try `CREATE EXTENSION postgis` and `CREATE EXTENSION postgres_fdw`. Record what exists.

#### Medium practical tasks

1. Write a booking report: exclusion constraint on `during`, `rank()` of bookings per room, `count(*) FILTER` for long bookings.
2. `EXPLAIN` a large `count(*)` with parallelism on, then a `Foreign Scan` plan from the docs if you have no remote server.
3. Create and drop a logical slot in a lab, or write why `wal_level` blocked it.

#### Advanced practical tasks

1. Design a "stores plus remote inventory" sketch: PostGIS or lat/lon, FDW or replica, and a CDC path to a search index. State one risk each.
2. Map this topic to official chapters: window, aggregates, range types, `postgres_fdw`, parallel query, logical decoding. Add one fact that this topic did not include.
