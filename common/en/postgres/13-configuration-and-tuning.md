# 13. Configuration and Tuning

## Description

This topic shows configuration and common tuning patterns in PostgreSQL 16 and PostgreSQL 17. You learn memory settings, checkpoints, slow-statement logging, activity and lock views, keyset pagination, covering indexes, `max_connections`, and read replicas.

Complete topic 12 first. You can name sessions with `application_name`. Complete this topic before you harden production.

Use one term for each concept. A setting is a name that `SHOW` and `postgresql.conf` use. `shared_buffers` is the server page cache inside PostgreSQL. `work_mem` is memory for one sort or hash step. A checkpoint is a flush of dirty pages and a WAL position. Keyset pagination is a seek on the last seen sort key. A covering index includes extra columns so that an index-only scan can avoid the heap. This handbook targets PostgreSQL 16 and PostgreSQL 17.

---

## `shared_buffers`, `work_mem`, `effective_cache_size`

PostgreSQL uses several memory settings. They are not one "RAM percent" knob.

**`shared_buffers`** is the shared cache of table and index pages. A common start is about 25 percent of machine RAM on a dedicated host, with an upper bound that you measure. The operating system cache still holds pages. Do not set `shared_buffers` to almost all RAM.

**`work_mem`** is the limit for one sort or hash operation in one query step. A query can use several steps. Many sessions can sort at the same time. A host with `max_connections = 200` and `work_mem = 100MB` can request more RAM than the machine has. Start low (for example `8MB` or `16MB`). Raise it for a known heavy query in that session: `SET work_mem = '64MB';`.

**`maintenance_work_mem`** is for `VACUUM`, `CREATE INDEX`, and `ALTER TABLE`. One maintenance command uses this budget. You can set it higher than `work_mem` because few maintenance jobs run at once.

**`effective_cache_size`** is a planner hint. It is not an allocation. Set it to about the RAM that you expect for OS cache plus `shared_buffers` (often 50 to 75 percent of RAM on a dedicated host). A wrong value does not allocate RAM. It can change whether the planner picks an index.

```sql
SHOW shared_buffers;
SHOW work_mem;
SHOW maintenance_work_mem;
SHOW effective_cache_size;
```

Do not copy a blog that sets `shared_buffers = 80%` and `work_mem = 1GB` on a shared laptop. Do not raise `work_mem` globally to fix one report. Do not treat `effective_cache_size` as a memory reserve.

Official: [https://www.postgresql.org/docs/17/runtime-config-resource.html](https://www.postgresql.org/docs/17/runtime-config-resource.html).

### Questions

#### Theoretical questions

1. What does `shared_buffers` hold?
2. Why can a high `work_mem` plus a high `max_connections` exhaust RAM?
3. What commands use `maintenance_work_mem`?
4. Does `effective_cache_size` allocate memory?
5. How do you raise `work_mem` for one session only?

#### Easy practical tasks

1. `SHOW` the four settings in this section. Write the values.
2. Write a two-column table: setting, allocates RAM yes/no.
3. In one session `SET work_mem = '32MB'; SHOW work_mem;` Then reconnect and `SHOW` again.
4. Open the 17 resource-consumption page. Write the unit that each setting uses.

#### Medium practical tasks

1. For a host with 16 GB RAM and `max_connections = 100`, propose starting values for the four settings. Show the arithmetic for worst-case `work_mem`.
2. Run `EXPLAIN (ANALYZE)` on a sort query. Change `work_mem` so that the plan stops spilling to disk (or write that it already fits).
3. Compare `shared_buffers` on your Docker instance with a dedicated-host rule. Write why they differ.

#### Advanced practical tasks

1. Write a memory policy for development, test, and production. Include who may `ALTER SYSTEM`.
2. Read `huge_pages` and `shared_buffers` on Linux in the docs. Write one sentence on when a team enables huge pages.

---

## Checkpoints and `max_wal_size`

A checkpoint flushes dirty buffers and records a position. Recovery starts from the latest checkpoint and then replays WAL. Checkpoints that run too often cause extra writes. Checkpoints that wait too long make WAL grow and make crash recovery longer.

**`max_wal_size`** is a target for WAL growth between checkpoints. When WAL grows toward this size, the server starts a checkpoint. **`min_wal_size`** is a floor for recycled WAL files.

**`checkpoint_timeout`** is the time limit between checkpoints (default is often 5 minutes). The server checkpoints at the sooner of the time limit and the WAL-size trigger.

```sql
SHOW max_wal_size;
SHOW min_wal_size;
SHOW checkpoint_timeout;
SELECT * FROM pg_stat_bgwriter;
```

`pg_stat_bgwriter` and, in PostgreSQL 17, related I/O views help you see checkpoint activity. PostgreSQL 16 also has `pg_stat_io` for I/O stats. Use the views that your major version ships.

`CHECKPOINT` runs a checkpoint now. Use it in labs. Do not run it in a loop on production.

A replica and `archive_command` must keep up. A large `max_wal_size` can be fine if disks are fast and recovery time is acceptable.

Do not set `max_wal_size` to a tiny value to "stay safe". You increase write amplification. Do not ignore checkpoint spikes after a bulk `COPY`.

Official: [https://www.postgresql.org/docs/17/wal-configuration.html](https://www.postgresql.org/docs/17/wal-configuration.html).

### Questions

#### Theoretical questions

1. What does a checkpoint flush?
2. What does `max_wal_size` trigger?
3. What does `checkpoint_timeout` trigger?
4. Why is a very small `max_wal_size` a problem?
5. What does `CHECKPOINT` do?

#### Easy practical tasks

1. `SHOW max_wal_size;` `SHOW checkpoint_timeout;`.
2. Run `CHECKPOINT;` in a lab. Record the time it takes.
3. Select a few columns from `pg_stat_bgwriter`.
4. Write four sentences: WAL, checkpoint, timeout, `max_wal_size`.

#### Medium practical tasks

1. `COPY` or insert many rows. Watch `pg_stat_bgwriter` before and after. Write what changed.
2. Read `pg_stat_io` on 16 or 17. Write one row that relates to WAL or checkpoints.
3. Compare recovery-time versus write-amp in two sentences for a larger `max_wal_size`.

#### Advanced practical tasks

1. Propose `max_wal_size` and `checkpoint_timeout` for a reporting host versus an OLTP host. Give reasons.
2. Read WAL configuration in the 17 docs. Add `checkpoint_completion_target`. Write what it spreads.

---

## Logging slow statements

PostgreSQL can write a log line for every statement that runs at least N milliseconds.

```sql
SHOW log_min_duration_statement;
SET log_min_duration_statement = 200;
```

`0` logs every statement. `-1` turns the feature off (common default). A global `0` on a busy server fills the disk and costs I/O.

Related settings (high-level):

- `logging_collector` — capture logs into files
- `log_directory` and `log_filename` — where files go
- `log_line_prefix` — add time, user, `application_name`, pid
- `log_checkpoints` — log checkpoint start and stop
- `log_lock_waits` — log waits longer than `deadlock_timeout`
- `log_connections` / `log_disconnections` — useful in labs; noisy in production

`SET` in a session is enough to learn. `ALTER SYSTEM` plus `SELECT pg_reload_conf();` changes the running server when the setting is reloadable.

Read logs to find slow SQL. Then use `EXPLAIN ANALYZE` (topic 5) and `pg_stat_statements`. Do not guess.

Do not leave `log_min_duration_statement = 0` on production. Do not log passwords. Client `SET` statements can appear in logs; keep secrets out of SQL text.

Official: [https://www.postgresql.org/docs/17/runtime-config-logging.html](https://www.postgresql.org/docs/17/runtime-config-logging.html).

### Questions

#### Theoretical questions

1. What does `log_min_duration_statement = 200` write?
2. What does `0` mean for that setting?
3. What does `-1` mean?
4. Why can full statement logging fill the disk?
5. What does `log_lock_waits` add?

#### Easy practical tasks

1. `SHOW log_min_duration_statement;` `SHOW log_line_prefix;`.
2. `SET log_min_duration_statement = 0;` Run `SELECT pg_sleep(0.05);`. Find the log line if you have file access.
3. Write five settings from this section and one purpose each.
4. Open the 17 logging page. Write which settings need a restart versus a reload (pick two).

#### Medium practical tasks

1. Set a prefix that includes `application_name`. Connect with that name. Run a slow `pg_sleep`. Match the log line to the app name.
2. Enable `log_lock_waits` in a session or lab. Open a lock (topic 6). Wait from a second session. Find the log if the wait is long enough.
3. Write a production starting value for `log_min_duration_statement` and a reason.

#### Advanced practical tasks

1. Design a log rotation rule (size or time) from `log_filename` and `log_rotation_*` docs.
2. Compare logging plus `pg_stat_statements`. Write what each source is for.

---

## `pg_stat_activity` and locks

**`pg_stat_activity`** shows one row per backend. Useful columns: `pid`, `usename`, `application_name`, `client_addr`, `state`, `wait_event_type`, `wait_event`, `query`, `xact_start`, `query_start`.

```sql
SELECT pid, usename, application_name, state, wait_event, left(query, 80)
FROM pg_stat_activity
WHERE datname = current_database()
  AND pid <> pg_backend_pid();
```

`state` values include `active` and `idle in transaction`. A long `idle in transaction` blocks vacuum (topic 7).

**Locks:** `pg_locks` joined to `pg_stat_activity` shows who waits. `log_lock_waits` writes a log line. `pg_blocking_pids(pid)` lists blockers.

```sql
SELECT * FROM pg_blocking_pids(12345);
```

**`pg_stat_statements`** aggregates time per normalized query. You need `shared_preload_libraries` and `CREATE EXTENSION` (topics 5 and 9).

```sql
SELECT calls, round(total_exec_time::numeric, 1) AS total_ms, left(query, 80)
FROM pg_stat_statements
ORDER BY total_exec_time DESC
LIMIT 10;
```

`pg_cancel_backend(pid)` asks a query to stop. `pg_terminate_backend(pid)` ends the session. Prefer cancel first.

Do not run `SELECT * FROM pg_stat_activity` in a tight loop in the application. Do not grant these views to the world if they expose SQL text with secrets.

Official: [https://www.postgresql.org/docs/17/monitoring-stats.html](https://www.postgresql.org/docs/17/monitoring-stats.html).

### Questions

#### Theoretical questions

1. What does `pg_stat_activity.state` tell you?
2. Why is `idle in transaction` an operations problem?
3. How do you find the pid that blocks pid 12345?
4. What extra install step does `pg_stat_statements` need?
5. What is the difference between cancel and terminate?

#### Easy practical tasks

1. Select your own row from `pg_stat_activity`.
2. Open a second session. Find it from the first session.
3. If the extension exists, select the top five `pg_stat_statements` rows by `total_exec_time`.
4. Write four sentences: activity, lock, statements, cancel.

#### Medium practical tasks

1. In session A `BEGIN; UPDATE` a row. In session B update the same row. Query `pg_locks` and `pg_blocking_pids` from a third session. Rollback A.
2. Filter `pg_stat_activity` for `wait_event IS NOT NULL`.
3. Reset statements in a lab (`pg_stat_statements_reset`). Run two queries. Show the new totals.

#### Advanced practical tasks

1. Write a monitoring query that lists sessions idle in transaction for more than five minutes.
2. Read `pg_stat_io` (16+) or `pg_stat_bgwriter`. Write one alert idea that is not "CPU high".

---

## Keyset pagination and covering indexes (`INCLUDE`)

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

A **covering** index adds non-key columns with `INCLUDE`. The planner can use an index-only scan when the query needs only the key columns plus the included columns, and the visibility map allows it.

```sql
CREATE INDEX posts_created_id_inc
    ON shop.posts (created_at DESC, id DESC)
    INCLUDE (title);
```

`INCLUDE` columns are not part of the search order. They still count as indexed columns for HOT updates (topic 7). Do not `INCLUDE` columns that you update on every request.

`OFFSET` is acceptable for small admin pages. It is the wrong default for infinite scroll on a large table.

Do not paginate on a non-unique column alone. Do not add five covering indexes that overlap. Do not use `LIMIT` / `OFFSET` to hide a missing index.

### Questions

#### Theoretical questions

1. Why does a large `OFFSET` get slower?
2. What extra column do you add when `created_at` is not unique?
3. What index matches a keyset `ORDER BY created_at DESC, id DESC`?
4. What does `INCLUDE` add to a B-tree?
5. Why can `INCLUDE` of a hot column block HOT updates?

#### Easy practical tasks

1. Select the first page with `ORDER BY created_at DESC, id DESC LIMIT 5`.
2. Write the second-page query with a keyset from the last row of page one.
3. Create the composite index. Run `\d`.
4. Add `INCLUDE (title)` on a copy of the index. `EXPLAIN` a query that selects those columns.

#### Medium practical tasks

1. Load 20000 rows. Compare `EXPLAIN ANALYZE` for `OFFSET 15000` and a keyset at the same position.
2. Compare `Index Scan` and `Index Only Scan` with and without `INCLUDE` of the selected columns.
3. Omit `id` from the `WHERE` tuple. Write a case with two equal timestamps that can skip a row.

#### Advanced practical tasks

1. Implement next and previous pages (two directions). Write both predicates and the matching index.
2. Read "Use the Index, Luke" on keyset pagination (PostgreSQL). Write one extra rule that this section did not include.

---

## `max_connections` and read replicas

`max_connections` is the cap on server backends. Each backend uses memory (`work_mem` per step, plus per-backend overhead). The operating system spends more time on context switches when hundreds of backends are active.

Thrashing here means: too many backends, too little RAM, heavy swap or cache eviction, and every query slower than a smaller pool.

The pattern:

1. set `max_connections` to a number the host can run (often 100 to 200 on a mid host, not 2000)
2. put PgBouncer or a driver pool in front (topic 12)
3. size the pool to a few dozen server connections for OLTP
4. use a separate small pool for admin and migrations

```sql
SHOW max_connections;
SELECT count(*) FROM pg_stat_activity;
```

Cloud instances often default `max_connections` by RAM class. Do not raise the default to "remove too many connections" errors without a pool. That error is a signal to pool, not to multiply backends.

A **read replica** is a hot stand-by (topic 11). You send reporting `SELECT`s there. Writes stay on the primary.

Rules for replica reads:

- tolerate lag
- tolerate query cancel (`hot_standby` conflicts)
- do not send writes
- give the replica its own `max_connections` and pool
- use `application_name` so that you can see report traffic

Do not set `max_connections` to 1000 on a 4 GB host. Do not open a session per HTTP request. Do not treat a replica as a second primary.

### Questions

#### Theoretical questions

1. What does `max_connections` cap?
2. How does `work_mem` change the risk of a high cap?
3. What error is a signal to add a pool rather than raise the cap?
4. What must a report query tolerate on a replica?
5. Do replicas share the primary `max_connections` value?

#### Easy practical tasks

1. `SHOW max_connections;` Count `pg_stat_activity`.
2. Write a table: host RAM, suggested starting cap, pool in front yes/no.
3. Write four sentences: backend, pool, replica lag, cancel.
4. If you have a replica, run `SELECT pg_is_in_recovery();` on both nodes.

#### Medium practical tasks

1. Compute worst-case RAM: `max_connections * work_mem` plus `shared_buffers`. Write whether the host can hold it.
2. Plan pool sizes: API 40, reports 10, admin 5, `max_connections` 80. Explain the headroom.
3. From `pg_stat_activity`, group by `application_name` and `state`. Write who holds idle sessions.

#### Advanced practical tasks

1. Write a routing rule: which queries go to the primary and which go to the replica. Include lag SLA.
2. Read `max_connections` and reserved connections (`superuser_reserved_connections`) in the 17 docs. Write why you leave headroom for admin.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do `shared_buffers`, `work_mem`, and `max_connections` form one RAM budget?
2. When do you change `max_wal_size` instead of adding an index?
3. How do slow-query logs and `pg_stat_activity` answer different questions?
4. Why do keyset pagination and `INCLUDE` belong in the same tuning topic?
5. A teammate raises `max_connections` to fix timeouts and points reports at a replica with no lag check. Which facts do you use in the reply?

#### Easy practical tasks

1. `SHOW` `shared_buffers`, `work_mem`, `max_wal_size`, `log_min_duration_statement`, and `max_connections`. Save the five values.
2. Write a cheat sheet: memory knobs, checkpoint, slow log, activity, keyset, `INCLUDE`, replica.
3. Run one `EXPLAIN ANALYZE` on a paged query. Write the scan type.
4. Select `count(*)` from `pg_stat_activity`.

#### Medium practical tasks

1. Tune one lab query: `ANALYZE`, one index (maybe `INCLUDE`), keyset instead of `OFFSET`. Prove with `EXPLAIN ANALYZE`.
2. Cause a lock wait. Find it in `pg_stat_activity` and `pg_blocking_pids`. Log it if `log_lock_waits` is on.
3. Write a pool-and-replica diagram with `application_name` on each path.

#### Advanced practical tasks

1. Write a one-page tuning order: measure (`pg_stat_statements`, logs), plan (`EXPLAIN`), index, memory, checkpoints, pool, replica. Give one example at each step.
2. Map this topic to official resource, WAL, logging, and monitoring chapters. Add one 16-or-17 view that this topic named only in passing.
