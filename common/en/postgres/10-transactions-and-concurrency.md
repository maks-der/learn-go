# 10. Transactions and Concurrency

## Description

This topic shows transactions and concurrency in PostgreSQL 16 and PostgreSQL 17. You learn MVCC, isolation levels, row locks, advisory locks, deadlocks, and snapshots. Readers do not block writers under MVCC for normal `SELECT`.

Complete topic 9 first. You can read plans. Complete this topic before you study vacuum.

Use one term for each concept. A transaction is a group of statements that commit or roll back together. Isolation is how much one transaction sees of another. A snapshot is the set of transactions that are visible. A lock waits or fails when another session holds a conflicting lock. MVCC means Multi-Version Concurrency Control.

---

## MVCC in PostgreSQL (readers do not block writers)

PostgreSQL keeps more than one version of a row. An `UPDATE` writes a new tuple and marks the old tuple as dead. A `DELETE` marks a tuple as dead. A `SELECT` uses a snapshot. The `SELECT` reads the version that is visible to that snapshot. It does not wait for a writer that changes a different version.

Rules that beginners must remember:

- A normal `SELECT` does not take a row lock.
- A normal `SELECT` does not wait for `UPDATE` or `DELETE` on the same row.
- An `UPDATE` waits for another `UPDATE` (or `SELECT FOR UPDATE`) on the same row.
- Dead tuples stay until `VACUUM` (topic 11).

Each tuple has hidden system columns. You can select some of them:

```sql
SELECT xmin, xmax, ctid, * FROM items WHERE id = 1;
```

`xmin` is the inserting transaction id (32-bit xid display). `ctid` is the physical location `(page, offset)`. After an `UPDATE`, `ctid` changes. The logical `id` stays the same.

Because readers do not block writers, PostgreSQL can serve many reads during writes. Writers still conflict with writers on the same row.

`INSERT` of a new key does not wait for a `SELECT`. Unique indexes still force two inserters of the same key to serialize.

MVCC is not the same as "no locks". DDL takes stronger locks. `SELECT FOR UPDATE` takes row locks. Foreign keys take locks. Topic later sections cover those.

Do not use `SELECT` as a lock. If you must stop another writer, use `SELECT ... FOR UPDATE` or an `UPDATE`.

### Questions

#### Theoretical questions

1. What does MVCC store for an updated row?
2. Does a normal `SELECT` wait for an `UPDATE` of the same row?
3. When does one `UPDATE` wait for another `UPDATE`?
4. What is `ctid`?
5. Why is `SELECT` not a lock?

#### Easy practical tasks

1. Select `xmin`, `xmax`, and `ctid` for one row.
2. `UPDATE` that row. Select `ctid` again. Write if `ctid` changed.
3. Write five sentences: version, snapshot, dead tuple, reader, writer.
4. Draw two sessions: one `SELECT`, one `UPDATE`, same row. Mark wait or no wait.

#### Medium practical tasks

1. Open two `psql` sessions. In session A, `BEGIN; UPDATE` a row; do not commit. In session B, `SELECT` that row. Write what B sees. Then `COMMIT` A and `SELECT` again in B.
2. In B, try `UPDATE` of the same row while A still holds the open `UPDATE`. Write whether B waits.
3. Read "Concurrency Control" introduction in the docs. Write the official sentence about readers and writers in your own words.

#### Advanced practical tasks

1. Show that a unique `INSERT` waits when another transaction holds an uncommitted insert of the same key. Use `ROLLBACK` on the first session.
2. Query `pg_stat_activity` while a session waits. Write `wait_event_type` and `wait_event`.

---

## Isolation levels; default is read committed

SQL defines isolation levels. PostgreSQL implements:

| Level | PostgreSQL behavior |
| --- | --- |
| `READ UNCOMMITTED` | Same as `READ COMMITTED`. Dirty reads do not occur. |
| `READ COMMITTED` | Default. Each statement sees rows committed before that statement starts. |
| `REPEATABLE READ` | Snapshot from the first statement. Later statements in the same transaction see the same snapshot. |
| `SERIALIZABLE` | Repeatable read plus checks that prevent more anomalies. A transaction can fail with a serialization error. |

Set the level:

```sql
BEGIN TRANSACTION ISOLATION LEVEL READ COMMITTED;
-- statements
COMMIT;

BEGIN;
SET TRANSACTION ISOLATION LEVEL REPEATABLE READ;
```

`default_transaction_isolation` is `read committed` unless you change it.

**Read committed** facts:

- Statement 1 can see row version V1. After another session commits an update, statement 2 in the same transaction can see V2.
- Non-repeatable reads are possible. Phantom rows are possible.
- `UPDATE` finds the latest committed version and waits if needed. Then it re-evaluates `WHERE` on that version.

Most transactional applications use read committed. It is the right default.

You cannot read uncommitted data from another transaction. There is no dirty read.

Show the current level:

```sql
SHOW transaction_isolation;
SELECT current_setting('transaction_isolation');
```

Do not raise the global default to `serializable` without a plan for retries. Do not assume `READ UNCOMMITTED` gives dirty reads. It does not.

### Questions

#### Theoretical questions

1. What is the default isolation level?
2. What does `READ UNCOMMITTED` do in PostgreSQL?
3. What can change between two statements in read committed?
4. Are dirty reads possible?
5. How do you set the isolation level for one transaction?

#### Easy practical tasks

1. Run `SHOW transaction_isolation;`.
2. `BEGIN; SHOW transaction_isolation; COMMIT;`
3. Start a transaction at `REPEATABLE READ`. Show the setting. Commit.
4. Make a table of four SQL levels and the PostgreSQL behavior.

#### Medium practical tasks

1. In read committed, show a non-repeatable read: session A `SELECT`, session B `UPDATE` and `COMMIT`, session A `SELECT` again in the same transaction.
2. Repeat the demo at `REPEATABLE READ`. Write if the second `SELECT` changed.
3. Read the official isolation table. Copy the anomaly names (dirty, nonrepeatable, phantom) in your own words.

#### Advanced practical tasks

1. Demonstrate an `UPDATE` in read committed that waits, then updates the new version. Use two sessions. Write the final value.
2. Find `default_transaction_isolation` in `postgresql.conf` or `SHOW ALL`. Write why you would not change it globally without retries.

---

## `REPEATABLE READ` and `SERIALIZABLE`

**Repeatable read** takes a snapshot at the first query (or first statement that needs data). All later `SELECT` statements see that snapshot. If you `UPDATE` a row that another transaction committed after your snapshot, PostgreSQL can raise:

```text
ERROR:  could not serialize access due to concurrent update
```

That error is a serialization failure. You roll back and retry the transaction.

Repeatable read in PostgreSQL uses snapshot isolation. Some write-skew anomalies can still occur. Two transactions can write different rows that together break a rule that neither transaction saw.

**Serializable** uses Serializable Snapshot Isolation (SSI). The server tracks read/write dependencies. If a dangerous pattern appears, one transaction fails with a serialization error. When all transactions that touch the same data use `SERIALIZABLE` and retry on failure, the result matches one serial order.

```sql
BEGIN TRANSACTION ISOLATION LEVEL SERIALIZABLE;
-- reads and writes
COMMIT;
```

Rules:

- The application must retry on serialization failure (`SQLSTATE 40001`).
- All sessions that need the guarantee must use `SERIALIZABLE`.
- Serializable has more overhead than read committed.
- Read-only serializable transactions can use a safer snapshot mode.

Use repeatable read for consistent reports that run many statements. Use serializable when the business rule cannot accept write skew and you will write a retry loop. Use read committed for most transactional work.

Do not ignore serialization errors. Do not retry forever without a backoff and a limit.

### Questions

#### Theoretical questions

1. When does repeatable read take its snapshot?
2. What is a serialization failure?
3. What extra anomalies can occur in repeatable read but not in serializable?
4. What is SSI?
5. Why must the application retry serializable transactions?

#### Easy practical tasks

1. Start `REPEATABLE READ`. `SELECT` a row. In another session, update and commit that row. `SELECT` again. Write the values.
2. Try `UPDATE` of that row in the first session. Write success or the serialize error.
3. Start `SERIALIZABLE`. `SHOW transaction_isolation;`. Commit.
4. Write the `SQLSTATE` that you must retry.

#### Medium practical tasks

1. Write a two-session script that produces `could not serialize access due to concurrent update` under `REPEATABLE READ`.
2. Read "Serializable Isolation Level" in the docs. Write one write-skew story in four sentences.
3. Compare duration of a read-only report in read committed versus repeatable read on your data (rough timer).

#### Advanced practical tasks

1. Implement a retry loop in a language that you know: on `40001`, retry up to three times. Use a practice `UPDATE`.
2. Read about `default_transaction_read_only` and serializable read-only. Write when a read-only serializable transaction is useful.

---

## Row locks: `SELECT ... FOR UPDATE`

`SELECT ... FOR UPDATE` locks the selected rows. Another session that wants to `UPDATE`, `DELETE`, or `SELECT FOR UPDATE` those rows waits (or fails with `NOWAIT`).

```sql
BEGIN;
SELECT * FROM items WHERE id = 1 FOR UPDATE;
-- other statements
COMMIT;
```

Lock strengths (strongest last among these four):

- `FOR KEY SHARE` — allows updates that do not change the key; used by foreign keys
- `FOR SHARE` — share lock; blocks `FOR UPDATE`
- `FOR NO KEY UPDATE` — like update but does not block `KEY SHARE`
- `FOR UPDATE` — blocks other exclusive row locks and updates

Options:

- `NOWAIT` — error instead of wait
- `SKIP LOCKED` — skip rows that are already locked (work queues)
- `OF table` — lock rows of one table in a join

```sql
SELECT id FROM jobs
WHERE status = 'new'
ORDER BY id
FOR UPDATE SKIP LOCKED
LIMIT 1;
```

The lock lasts until the end of the transaction. `COMMIT` or `ROLLBACK` releases it.

`FOR UPDATE` does not lock the whole table. Other rows stay available.

Use row locks when you read a value and then write a decision (balance check, job claim). Do not hold a transaction open while you wait for a user in a browser.

See locks:

```sql
SELECT * FROM pg_locks;
SELECT * FROM pg_stat_activity WHERE wait_event_type = 'Lock';
```

### Questions

#### Theoretical questions

1. What does `SELECT ... FOR UPDATE` lock?
2. When is the lock released?
3. What does `SKIP LOCKED` do?
4. What does `NOWAIT` do?
5. Why is a long `FOR UPDATE` transaction a problem?

#### Easy practical tasks

1. In a transaction, `SELECT ... FOR UPDATE` one row. Commit.
2. From a second session, try `UPDATE` of that row before commit. Write that it waits. Then commit the first session.
3. Repeat with `NOWAIT`. Record the error.
4. Write the four lock strengths in order.

#### Medium practical tasks

1. Build a two-row job table. Use `FOR UPDATE SKIP LOCKED LIMIT 1` from two sessions. Confirm they claim different rows.
2. Join two tables and use `FOR UPDATE OF items`. Write which rows lock.
3. Query `pg_locks` while a lock is held. Write `locktype` and `mode`.

#### Advanced practical tasks

1. Compare `FOR SHARE` and `FOR UPDATE` from two sessions (two shares succeed; share plus update waits).
2. Read the official row-level lock table. Write which lock a foreign key check takes when you insert a child row.

---

## Advisory locks

Advisory locks are locks that the application defines. PostgreSQL does not bind them to a table row. You pass integer keys.

Session-level locks stay until you unlock or the session ends:

```sql
SELECT pg_advisory_lock(42);
-- work
SELECT pg_advisory_unlock(42);
```

Transaction-level locks stay until `COMMIT` or `ROLLBACK`:

```sql
BEGIN;
SELECT pg_advisory_xact_lock(42);
COMMIT;  -- lock released
```

Two-argument forms use a pair of `int4` keys: `pg_advisory_lock(1, 2)`.

`pg_try_advisory_lock` returns `false` instead of waiting.

Uses:

- single-job cron: lock key 1001 at start; skip if the lock is held
- migrate a resource that has no row yet
- serialize work that spans more than one table

Rules:

- Pick a key map and write it down. Avoid collisions.
- Prefer transaction locks when the work is one transaction.
- Session locks must have a matching unlock on every path. A forgotten unlock blocks others until disconnect.
- Advisory locks do not survive a crash as data. They are not constraints.

See held advisory locks in `pg_locks` where `locktype = 'advisory'`.

Do not use advisory locks instead of a unique constraint. Do not use them instead of `FOR UPDATE` when a row exists.

### Questions

#### Theoretical questions

1. What is an advisory lock?
2. What is the difference between `pg_advisory_lock` and `pg_advisory_xact_lock`?
3. What does `pg_try_advisory_lock` return when the lock is held?
4. When is an advisory lock a better fit than `FOR UPDATE`?
5. Why must session unlock run on every path?

#### Easy practical tasks

1. Take `pg_advisory_lock(1)`. Query `pg_locks` for `advisory`. Unlock.
2. In a transaction, take `pg_advisory_xact_lock(2)`. Commit. Confirm the lock is gone.
3. Use `pg_try_advisory_lock` from two sessions on the same key.
4. Write four sentences: key, session, transaction, try.

#### Medium practical tasks

1. Simulate a cron lock: session A holds 1001; session B tries and skips.
2. Forget unlock on a session lock (test only). Connect a second session and show the wait. Disconnect A. Confirm B proceeds.
3. Use a two-integer key pair. Document your key map.

#### Advanced practical tasks

1. Write a small program or `psql` script that takes a lock, runs a `SELECT`, and always unlocks (use a transaction lock to make this easy).
2. Read the advisory lock docs. List blocking and try functions for both session and transaction scopes.

---

## Deadlocks in logs

A deadlock happens when session A waits for session B, and session B waits for session A. PostgreSQL detects the cycle. It aborts one transaction:

```text
ERROR:  deadlock detected
```

The other transaction continues. The aborted session must retry.

Typical cause: two transactions update the same two rows in opposite order.

```text
Session A: UPDATE t SET ... WHERE id = 1;
Session B: UPDATE t SET ... WHERE id = 2;
Session A: UPDATE t SET ... WHERE id = 2;  -- waits
Session B: UPDATE t SET ... WHERE id = 1;  -- deadlock
```

Prevention:

- Always update rows in the same order (example: order by `id`)
- Keep transactions short
- Do not hold locks while you wait for the network

PostgreSQL writes a detail line that shows the two processes. `log_lock_waits` and `deadlock_timeout` control wait logging. `deadlock_timeout` is how long a wait lasts before the server checks for a deadlock (default `1s`).

```sql
SHOW deadlock_timeout;
SHOW log_lock_waits;
```

The server log (or `stderr` in Docker) contains the deadlock report. `psql` shows the error to the client that lost.

A deadlock is not a crash. It is a normal error under concurrency. The application retries.

Do not raise `deadlock_timeout` to hide a lock-order bug. Fix the order.

### Questions

#### Theoretical questions

1. What is a deadlock?
2. Which transaction fails when PostgreSQL detects a deadlock?
3. What lock order rule prevents the two-row example?
4. What is `deadlock_timeout`?
5. Is a deadlock a server crash?

#### Easy practical tasks

1. `SHOW deadlock_timeout;` and `SHOW log_lock_waits;`.
2. Write the two-session update order that deadlocks. Do not run it yet.
3. Write the same updates in `id` order. Explain why that avoids a cycle.
4. Find "deadlock detected" in the official docs. Write the `SQLSTATE` (`40P01`).

#### Medium practical tasks

1. Produce a deadlock on a practice table with two sessions. Save both client messages. `ROLLBACK` as needed.
2. Find the deadlock line in the server log. Copy the process ids.
3. Set `log_lock_waits = on` in a session or in a test config. Cause a wait longer than `deadlock_timeout` without a deadlock. Write if a log line appeared.

#### Advanced practical tasks

1. Write application rules: lock order, retry on `40P01`, max retries.
2. Read `deadlock_timeout` and lock logging docs. Write how a wait log differs from a deadlock log.

---

## `txid` / snapshots (awareness)

A snapshot says which transaction ids are visible. You do not set a snapshot by hand in normal SQL. Isolation levels set it.

Older functions use the `txid_` prefix. Current PostgreSQL also has `xid8` functions:

```sql
SELECT pg_current_xact_id();           -- xid8, assigns an xid if needed
SELECT pg_current_xact_id_if_assigned();
SELECT pg_current_snapshot();
SELECT pg_snapshot_xmin(pg_current_snapshot());
SELECT pg_snapshot_xmax(pg_current_snapshot());
```

Older aliases still exist (`txid_current()`, `txid_current_snapshot()`). Prefer the `pg_current_*` names in new work on PostgreSQL 16 and 17.

`pg_current_xact_id()` assigns a transaction id if the transaction does not have one yet. A read-only transaction can avoid an xid until it writes. Use `pg_current_xact_id_if_assigned()` when you must not assign.

A snapshot has `xmin`, `xmax`, and a list of in-progress xids. Visibility rules use those values. You do not need the full rules to write application SQL. You need to know:

- Vacuum and wraparound use xids (topic 11)
- `xmin` on a tuple is the inserting xid
- Repeatable read and serializable freeze visibility to one snapshot

```sql
SELECT xmin, xmax FROM items;
```

Do not use xids as application keys. Xids wrap (topic 11). Do not compare xids as integers in app code.

This section is awareness. You will read xids again when you study freeze and wraparound.

### Questions

#### Theoretical questions

1. What does a snapshot decide?
2. What is the modern function that returns the current xid?
3. Why can `pg_current_xact_id()` change a read-only transaction?
4. Why must you not use an xid as an application key?
5. Which isolation levels keep one snapshot for many statements?

#### Easy practical tasks

1. Run `SELECT pg_current_snapshot();`.
2. Run `SELECT pg_current_xact_id_if_assigned();` in a fresh `BEGIN` before any write. Then `INSERT` and run it again.
3. Select `xmin` from one row.
4. Write four sentences: snapshot, xid8, assign, wrap.

#### Medium practical tasks

1. In repeatable read, print `pg_current_snapshot()` twice with a write from another session in between. Write if the snapshot text changed.
2. Compare `txid_current()` and `pg_current_xact_id()` in one transaction. Write both results.
3. Read `pg_snapshot_xmin` docs. Write what `xmin` in a snapshot means.

#### Advanced practical tasks

1. Read "Transaction ID Wraparound" preview in the official docs (topic 11). Write how snapshots and freeze relate in three sentences.
2. Query `age(datfrozenxid)` from `pg_database`. Write the value for your database. Topic 11 explains the risk.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do MVCC, read committed, and `FOR UPDATE` each change what a second session does?
2. When do you pick advisory locks, when row locks, and when serializable isolation?
3. Why must a deadlock and a serialization failure both retry, and how do their causes differ?
4. What does a snapshot have to do with vacuum and dead tuples?
5. Why does PostgreSQL implement `READ UNCOMMITTED` as read committed?

#### Easy practical tasks

1. Run one demo: open transaction, `UPDATE`, other session `SELECT` (no wait), other session `UPDATE` (wait). Commit.
2. Write a cheat sheet: MVCC, four isolation levels, `FOR UPDATE` options, advisory lock functions, deadlock, snapshot functions.
3. `SHOW` isolation, `deadlock_timeout`, and `transaction_isolation`.
4. Take and release one transaction advisory lock.

#### Medium practical tasks

1. Write a two-session script file (comments allowed) that demos non-repeatable read, then the same steps under repeatable read.
2. Produce either a deadlock or a serialization error. Save the client text and the `SQLSTATE`.
3. Document team defaults: isolation level, retry policy, max transaction time.

#### Advanced practical tasks

1. Build a tiny job queue with `FOR UPDATE SKIP LOCKED` and prove two workers never claim the same row.
2. Read the full "Concurrency Control" chapter table of contents. Map each section of this topic to a chapter heading.
