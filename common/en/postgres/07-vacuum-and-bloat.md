# 7. Vacuum and Bloat

## Description

This topic shows why PostgreSQL 16 and PostgreSQL 17 need vacuum. You learn dead tuples, autovacuum, `VACUUM` versus `VACUUM FULL`, transaction id wraparound, and bloat.

Complete topic 6 first. You must know MVCC and snapshots. Complete this topic before you study roles and security.

Use one term for each concept. A dead tuple is a row version that no snapshot needs. Vacuum removes dead tuples and updates the visibility map. Bloat is extra space that remains after heavy updates or deletes. Wraparound is the risk when 32-bit transaction ids reuse old numbers. Autovacuum is the background process that runs vacuum and analyze.

---

## Dead tuples after `UPDATE` / `DELETE`

MVCC does not overwrite a row in place for a normal `UPDATE`. The server writes a new tuple. It marks the old tuple as dead for future snapshots. `DELETE` marks the tuple as dead. The space does not return to the operating system. It stays in the table file until vacuum marks it reusable.

```sql
UPDATE items SET qty = qty + 1 WHERE id = 1;
DELETE FROM items WHERE id = 2;
```

Each of these statements can create dead tuples. A long-running transaction that still sees the old snapshot delays vacuum. Vacuum cannot remove a tuple that some snapshot still needs.

See dead tuple counts:

```sql
SELECT relname, n_live_tup, n_dead_tup, last_autovacuum, last_autoanalyze
FROM pg_stat_user_tables
WHERE relname = 'items';
```

`n_dead_tup` is an estimate. Run `ANALYZE` or wait for activity so that counters update.

Heap-only tuple (HOT) updates can avoid a new index entry when the indexed columns do not change and the new tuple fits on the same page. Know that some updates cost less than others.

`INSERT` does not create a dead tuple on success. A rolled-back `INSERT` leaves a tuple that vacuum must clean.

Dead tuples also exist in indexes. Vacuum cleans index pointers to dead heap tuples.

Do not run `UPDATE` of every row every minute without a vacuum plan. Do not keep a `BEGIN` open for hours on a busy database.

### Questions

#### Theoretical questions

1. Why does `UPDATE` create a dead tuple?
2. Why does `DELETE` not shrink the file at once?
3. What delays removal of a dead tuple?
4. Where do you read `n_dead_tup`?
5. Does `INSERT` create a dead tuple on success?

#### Easy practical tasks

1. Update one row several times. Select `n_dead_tup` for that table.
2. Delete a few rows. Select `n_dead_tup` again.
3. Select `xmin` and `ctid` before and after an `UPDATE`.
4. Write four sentences: new tuple, dead tuple, snapshot, vacuum.

#### Medium practical tasks

1. Open a transaction in session A and `SELECT` a row. In session B, update that row many times. Watch `n_dead_tup`. Commit A. Vacuum and watch again.
2. Compare an update that changes an indexed column with an update that changes only a non-indexed column. Read about HOT. Write what you expect for index bloat.
3. Query `pg_stat_activity` for a transaction that is open longer than five minutes (`xact_start`).

#### Advanced practical tasks

1. Load 10000 rows. Update all of them. Record `n_dead_tup` and `pg_relation_size` before vacuum.
2. Read "Heap-Only Tuples" in the official docs (or wiki). Write the two conditions for a HOT update.

---

## Autovacuum

Autovacuum is on by default. A launcher process starts worker processes. Workers run `VACUUM` and `ANALYZE` on tables that pass a threshold.

Important settings (high-level):

- `autovacuum` — on/off (keep on)
- `autovacuum_naptime` — how often the launcher looks at tables
- `autovacuum_vacuum_threshold` and `autovacuum_vacuum_scale_factor` — when to vacuum
- `autovacuum_analyze_threshold` and `autovacuum_analyze_scale_factor` — when to analyze
- `autovacuum_vacuum_insert_threshold` and insert scale factor — vacuum after many inserts (visibility map)
- `autovacuum_freeze_max_age` — wraparound protection (later section)

The vacuum trigger is about:

```text
threshold + scale_factor * reltuples
```

A large table needs more dead tuples before autovacuum starts. A busy small table vacuums sooner.

You can set storage parameters on one table:

```sql
ALTER TABLE items SET (autovacuum_vacuum_scale_factor = 0.05);
```

See workers:

```sql
SELECT * FROM pg_stat_activity
WHERE backend_type = 'autovacuum worker';
```

`pg_stat_progress_vacuum` shows a running vacuum.

Do not turn autovacuum off. If you turn it off for a bulk load, turn it on again and run `VACUUM ANALYZE` yourself.

Autovacuum can be busy after a large `UPDATE`. That is normal. If autovacuum never keeps up, you have too many updates, too few workers, or a long snapshot.

PostgreSQL 16 and 17 improved vacuum memory and related internals. You still configure thresholds the same way.

### Questions

#### Theoretical questions

1. What two jobs does autovacuum run?
2. What is the role of `scale_factor`?
3. Why must you keep `autovacuum` on?
4. How do you change autovacuum for one table?
5. Where do you see a running autovacuum worker?

#### Easy practical tasks

1. `SHOW autovacuum;` and `SHOW autovacuum_naptime;`.
2. Select `last_autovacuum` and `n_dead_tup` from `pg_stat_user_tables`.
3. Read `autovacuum_vacuum_scale_factor` with `SHOW`.
4. Write the threshold formula in one line.

#### Medium practical tasks

1. Set a lower `autovacuum_vacuum_scale_factor` on a practice table. Update many rows. Wait or watch for autovacuum. Restore the setting if you want defaults.
2. Query `pg_stat_progress_vacuum` during a large `VACUUM` (start `VACUUM` in another session).
3. List autovacuum-related settings from `pg_settings` with `WHERE name LIKE 'autovacuum%'`.

#### Advanced practical tasks

1. Read the official autovacuum docs. Write how `autovacuum_max_workers` and `autovacuum_vacuum_cost_limit` interact (high-level).
2. Design per-table autovacuum settings for a tiny queue table and a large append-only log. Write the `ALTER TABLE` lines.

---

## `VACUUM` vs `VACUUM FULL`

`VACUUM` removes dead tuples and makes space reusable inside the file. It does not return disk space to the operating system in the common case. It can run in parallel with reads and writes. It takes a modest lock.

```sql
VACUUM items;
VACUUM (VERBOSE) items;
VACUUM (ANALYZE) items;
VACUUM;
```

`VACUUM` without a table name processes the database. Superusers and certain roles can do that. Table owners can `VACUUM` their tables.

`VACUUM (ANALYZE)` also updates planner statistics. That pair is common after a bulk load. `ANALYZE` alone does not remove dead tuples.

`VACUUM FULL` rewrites the table into a new file. It returns space to the operating system. It builds new indexes. It takes an `ACCESS EXCLUSIVE` lock. No session can read or write the table during the rewrite. It needs extra disk space for the new copy.

```sql
VACUUM FULL items;
```

Use `VACUUM` (plain) for normal maintenance. Use `VACUUM FULL` only when you measured a large empty space and you can lock the table (or you use a tool that rewrites with less lock, such as `pg_repack`, which is not in core).

Options in current versions:

- `VACUUM (INDEX_CLEANUP ON)` — default cleanup of indexes
- `VACUUM (TRUNCATE ON)` — try to truncate empty pages at the end of the file
- `VACUUM (PARALLEL n)` — parallel index vacuum for a single table (plain vacuum)
- `VACUUM (FREEZE, VERBOSE)` — freeze old xids (next section)

Do not schedule `VACUUM FULL` every night on large production tables. Do not run `VACUUM FULL` instead of fixing a long transaction that blocks autovacuum.

### Questions

#### Theoretical questions

1. What does plain `VACUUM` do with dead tuples?
2. Does plain `VACUUM` usually shrink the file on disk?
3. What lock does `VACUUM FULL` take?
4. Why does `VACUUM FULL` need extra disk space?
5. When is `VACUUM (ANALYZE)` useful?

#### Easy practical tasks

1. Run `VACUUM items;` on a practice table.
2. Run `VACUUM (VERBOSE, ANALYZE) items;` and read the messages.
3. Write the difference between `VACUUM` and `VACUUM FULL` in four sentences.
4. List three `VACUUM` options from this section.

#### Medium practical tasks

1. Create dead tuples with `UPDATE`. Run plain `VACUUM`. Compare `n_dead_tup` and `pg_relation_size`.
2. Run `VACUUM FULL` on a small copy table. Compare size before and after a large delete.
3. Try `VACUUM FULL` while another session reads the table in a long transaction. Write what you observe (wait or error).

#### Advanced practical tasks

1. Read the `VACUUM` reference for `PARALLEL` and `INDEX_CLEANUP`. Write when you would set them.
2. Compare `VACUUM FULL` with `REINDEX` (topic 5). Write which problem each command fixes.

---

## Transaction ID wraparound

PostgreSQL labels row versions with a 32-bit transaction id (xid). About two billion xids are usable in the circular space. If the cluster assigns xids and never freezes old rows, old xids can become "in the future". Visibility then breaks. PostgreSQL stops new transactions to protect data:

```text
ERROR:  database is not accepting commands to avoid wraparound data loss
```

**Freeze** replaces old xids on tuples with a frozen marker. Vacuum freeze does this work. Autovacuum always runs a freeze vacuum when a table age approaches `autovacuum_freeze_max_age` (default 200 million). This emergency vacuum can be aggressive. It is not optional.

Watch age:

```sql
SELECT datname, age(datfrozenxid), datfrozenxid
FROM pg_database
ORDER BY age(datfrozenxid) DESC;

SELECT c.relname, age(c.relfrozenxid)
FROM pg_class AS c
JOIN pg_namespace AS n ON n.oid = c.relnamespace
WHERE n.nspname NOT IN ('pg_catalog', 'information_schema')
  AND c.relkind = 'r'
ORDER BY age(c.relfrozenxid) DESC;
```

`age(datfrozenxid)` should stay far below two billion. Many teams alert at a few hundred million.

Causes of wraparound risk:

- autovacuum off
- a very long transaction or a forgotten prepared transaction
- a slot or replica that holds a snapshot (logical slots, stale replicas)
- so much update traffic that freeze never finishes

Fixes: turn autovacuum on, end long transactions, vacuum freeze the old tables, add capacity.

```sql
VACUUM (FREEZE, VERBOSE) items;
```

You do not freeze by hand every day. You monitor age. You keep autovacuum healthy.

`xid8` (64-bit) functions in topic 6 do not remove 32-bit wraparound from the heap. Freeze still matters in PostgreSQL 16 and 17.

This risk is real. Treat wraparound alerts as urgent.

### Questions

#### Theoretical questions

1. Why can 32-bit xids wrap?
2. What does freeze do?
3. What happens when PostgreSQL stops commands to avoid wraparound?
4. What default setting forces an anti-wraparound vacuum?
5. Name three causes of wraparound risk.

#### Easy practical tasks

1. Query `age(datfrozenxid)` for all databases.
2. Query `age(relfrozenxid)` for your user tables.
3. `SHOW autovacuum_freeze_max_age;`.
4. Write four sentences: xid, freeze, autovacuum, alert.

#### Medium practical tasks

1. Run `VACUUM (FREEZE, VERBOSE)` on a practice table. Read the freeze messages.
2. Find long transactions (`xact_start`) and prepared transactions (`pg_prepared_xacts`) if any.
3. Read the official wraparound section. Write the meaning of `vacuum_freeze_min_age` in one sentence.

#### Advanced practical tasks

1. Write a monitoring query that lists the five tables with the highest `age(relfrozenxid)` and the five databases with the highest `age(datfrozenxid)`.
2. Read about `idle_in_transaction_session_timeout`. Write how it reduces wraparound risk.

---

## Table and index bloat

Bloat is unused space inside a table or index file. Dead tuples, failed HOT chains, and index splits cause it. The file stays large. Queries read more pages. Cache hit rate can drop.

Plain `VACUUM` makes space reusable. It does not always shrink the file. New inserts can reuse the space. If the table shrinks in row count and stays huge on disk, you have bloat.

Check sizes:

```sql
SELECT
    relname,
    pg_size_pretty(pg_relation_size(oid)) AS heap,
    pg_size_pretty(pg_indexes_size(oid)) AS indexes
FROM pg_class
WHERE relname = 'items';
```

`pgstattuple` and `pgstatindex` (extension `pgstattuple`) estimate free space. Many sites use queries from the PostgreSQL wiki (bloat estimates from `pg_class` and `reltuples`). Estimates can be wrong. Measure with the extension when you must know.

Reduce table bloat:

- autovacuum that keeps up
- fewer needless `UPDATE`s (update only changed columns; avoid touch-all jobs)
- `VACUUM FULL` or `CLUSTER` when you can lock
- `pg_repack` or similar tools when you cannot lock (external)

Reduce index bloat:

- `REINDEX INDEX CONCURRENTLY` (topic 5)
- drop unused indexes
- prefer HOT-friendly updates (do not change indexed columns)

Fillfactor below 100 leaves space on each page for HOT updates. A lower fillfactor can reduce bloat on hot update columns. It uses more disk. Set it with a measurement.

```sql
ALTER TABLE items SET (fillfactor = 80);
```

Do not chase zero bloat. Some free space is useful. Chase bloat that you measured (size versus row count, cache, I/O).

PostgreSQL 17 reduced some vacuum overhead. Bloat still appears if you update more than vacuum can clean.

### Questions

#### Theoretical questions

1. What is bloat?
2. Why can a table stay large after `DELETE` and plain `VACUUM`?
3. How do you read heap and index size?
4. How do you rebuild a bloated index without a long write lock?
5. What does `fillfactor` change?

#### Easy practical tasks

1. Run `pg_relation_size` and `pg_indexes_size` on one table.
2. Delete half the rows of a copy table. Compare size before `VACUUM`, after `VACUUM`, and after `VACUUM FULL`.
3. List indexes and their `pg_relation_size`.
4. Write four sentences: reusable space, file size, reindex, fillfactor.

#### Medium practical tasks

1. Update every row in a loop. Watch size and `n_dead_tup`. Vacuum. Watch again.
2. `REINDEX INDEX CONCURRENTLY` on one index. Compare size before and after if the index was bloated.
3. Try `CREATE EXTENSION pgstattuple` and `SELECT * FROM pgstattuple('items');` if your role allows it.

#### Advanced practical tasks

1. Find a wiki or docs query that estimates bloat percent. Run it. Write that it is an estimate.
2. Write an operations rule: when you vacuum, when you reindex concurrently, when you accept `VACUUM FULL` downtime.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do dead tuples, autovacuum, and wraparound form one maintenance story?
2. When is `VACUUM FULL` the wrong tool even if the table looks large?
3. Why must `ANALYZE` follow a bulk load even if you do not vacuum yet?
4. How does a long `idle in transaction` session harm both bloat and wraparound?
5. What is the difference between reusable heap space and index bloat?

#### Easy practical tasks

1. After a batch of updates, run `VACUUM (ANALYZE, VERBOSE)` and save the output. Then read `n_dead_tup` and `age(relfrozenxid)`.
2. Write a cheat sheet: dead tuples, autovacuum, `VACUUM`, `VACUUM FULL`, wraparound, bloat.
3. `SHOW` five autovacuum settings. Write their values.
4. Compute `pg_size_pretty` for your largest practice table and its indexes.

#### Medium practical tasks

1. Write a weekly checklist: ages, `n_dead_tup`, sizes, long transactions, autovacuum workers.
2. Simulate a load: insert, update, vacuum, analyze, explain a query. Record sizes and plans at each step.
3. Document who may run `VACUUM FULL` in production and what they must check first (disk, lock time, replicas).

#### Advanced practical tasks

1. Create an alert sketch (SQL only) for wraparound age, dead tuple ratio, and a vacuum that has not run for N days.
2. Read "Routine Vacuuming" in the official docs end to end. Map each heading to a section of this topic. Add one fact that this topic did not include.
