# 18. Configuration and Operations

## Description

This topic shows daily operations for PostgreSQL 16 and PostgreSQL 17. You learn memory settings, checkpoints, statement logging, activity and lock views, major-version upgrades, and disk-full incidents.

Complete topic 17 first. You can name sessions with `application_name`. Complete this topic before you apply performance patterns.

Use one term for each concept. A setting is a name that `SHOW` and `postgresql.conf` use. `shared_buffers` is the server page cache inside PostgreSQL. `work_mem` is memory for one sort or hash step. A checkpoint is a flush of dirty pages and a WAL position. `pg_stat_activity` is the live session view. A major upgrade changes the first version number. This handbook targets PostgreSQL 16 and PostgreSQL 17.

---

## Memory: `shared_buffers`, `work_mem`, `maintenance_work_mem`, `effective_cache_size`

PostgreSQL uses several memory settings. They are not one "RAM percent" knob.

**`shared_buffers`** is the shared cache of table and index pages. A common start is about 25 percent of machine RAM on a dedicated host, with an upper bound that you measure. The operating system cache still holds pages. Do not set `shared_buffers` to almost all RAM.

**`work_mem`** is the limit for one sort or hash operation in one query step. A query can use several steps. Many sessions can sort at the same time. A host with `max_connections = 200` and `work_mem = 100MB` can request more RAM than the machine has. Start low (for example `8MB` or `16MB`). Raise it for a known heavy query in that session: `SET work_mem = '64MB';`.

**`maintenance_work_mem`** is for `VACUUM`, `CREATE INDEX`, and `ALTER TABLE`. One maintenance command uses this budget. You can set it higher than `work_mem` (for example `256MB` or more on a large host) because few maintenance jobs run at once.

**`effective_cache_size`** is a planner hint. It is not an allocation. Set it to about the RAM that you expect for OS cache plus `shared_buffers` (often 50 to 75 percent of RAM on a dedicated host). A wrong value does not allocate RAM. It can change whether the planner picks an index.

```sql
SHOW shared_buffers;
SHOW work_mem;
SHOW maintenance_work_mem;
SHOW effective_cache_size;
```

Topic 9 introduced `work_mem` for sorts and hashes. This section is the operations view.

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

Do not set `max_wal_size` to a tiny value to "stay safe". You increase write amplification. Do not ignore checkpoint spikes after a bulk `COPY`. Topic 17 bulk loads create dirty pages.

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

## Logging: `log_min_duration_statement`

PostgreSQL can write a log line for every statement that runs at least N milliseconds.

```sql
SHOW log_min_duration_statement;
-- session test
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

Read logs to find slow SQL. Then use `EXPLAIN ANALYZE` (topic 9) and `pg_stat_statements` (below). Do not guess.

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
2. Enable `log_lock_waits` in a session or lab. Open a lock (topic 10). Wait from a second session. Find the log if the wait is long enough.
3. Write a production starting value for `log_min_duration_statement` and a reason.

#### Advanced practical tasks

1. Design a log rotation rule (size or time) from `log_filename` and `log_rotation_*` docs.
2. Compare logging plus `pg_stat_statements`. Write what each source is for.

---

## `pg_stat_activity`, locks, `pg_stat_statements`

**`pg_stat_activity`** shows one row per backend. Useful columns: `pid`, `usename`, `application_name`, `client_addr`, `state`, `wait_event_type`, `wait_event`, `query`, `backend_xid`, `xact_start`, `query_start`.

```sql
SELECT pid, usename, application_name, state, wait_event, left(query, 80)
FROM pg_stat_activity
WHERE datname = current_database()
  AND pid <> pg_backend_pid();
```

`state` values include `active` and `idle in transaction`. A long `idle in transaction` blocks vacuum (topic 11).

**Locks:** `pg_locks` joined to `pg_stat_activity` shows who waits. `log_lock_waits` writes a log line. `pg_blocking_pids(pid)` lists blockers.

```sql
SELECT * FROM pg_blocking_pids(12345);
```

**`pg_stat_statements`** (topic 9 and topic 13) aggregates time per normalized query. You need `shared_preload_libraries` and `CREATE EXTENSION`.

```sql
SELECT calls, round(total_exec_time::numeric, 1) AS total_ms, left(query, 80)
FROM pg_stat_statements
ORDER BY total_exec_time DESC
LIMIT 10;
```

`pg_stat_statements_reset()` clears the totals. Use it after a test, not at random on production.

`pg_cancel_backend(pid)` asks a query to stop. `pg_terminate_backend(pid)` ends the session. Prefer cancel first. Do not terminate a needed autovacuum worker without a reason (topic 11 wraparound).

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

## Upgrading major versions (`pg_upgrade`, dump/restore, logical)

A major upgrade (16 to 17) changes catalogs. You cannot start 17 binaries on a 16 data directory.

Three supported paths:

1. **`pg_upgrade`** — converts the data directory in place or with links. Fast. Needs both binaries installed. Run `--check` first. Read the 17 `pg_upgrade` page for `--link`, `--clone`, and `--copy`.
2. **dump and restore** — `pg_dumpall` or `pg_dump` of each database, new `initdb`, restore. Slow. Simple. Good when you change machines or encoding.
3. **logical replication** — new 17 cluster as subscriber, then switch writes. More moving parts. Useful for low downtime. Topic 16 restrictions still apply.

Minor upgrades (16.x to 16.y) do not use `pg_upgrade`. You install new binaries and restart. Topic 1 covered that.

After any major upgrade:

- `ANALYZE` the new cluster (or use `vacuumdb --analyze-in-stages`)
- run `REINDEX` if the notes say so (collation or glibc changes)
- test extensions
- update `shared_preload_libraries` and `postgresql.conf` for new names
- test the application on 17 before you drop 16

PostgreSQL 16 and 17 each have a "Migration" section in the release notes. Read it.

Do not skip `--check`. Do not upgrade production before a restore test. Do not develop on 17 and deploy to 16 without a feature list.

Official: [https://www.postgresql.org/docs/17/pgupgrade.html](https://www.postgresql.org/docs/17/pgupgrade.html), [https://www.postgresql.org/docs/17/upgrading.html](https://www.postgresql.org/docs/17/upgrading.html).

### Questions

#### Theoretical questions

1. Why can a 17 server not open a 16 data directory?
2. What does `pg_upgrade --check` do?
3. When is dump and restore the simpler path?
4. When is logical replication the path for less downtime?
5. What must you run after a major upgrade so that plans stay sane?

#### Easy practical tasks

1. Open the 17 migration notes. Write three compatibility items.
2. Write a table: method, needs extra cluster, typical downtime (low/medium/high).
3. `SHOW server_version;`
4. Find `vacuumdb --analyze-in-stages` in the docs. Write its purpose.

#### Medium practical tasks

1. Plan a 16 to 17 upgrade in eight steps for `pg_upgrade`. Include `--check` and analyze.
2. Compare `--link` and `--copy` in the `pg_upgrade` docs. Write the risk of `--link`.
3. List extensions in `pg_extension`. Mark which ones you must install on the new cluster first.

#### Advanced practical tasks

1. Upgrade a lab data directory from 16 to 17 with `pg_upgrade` or dump/restore. Save `SHOW server_version` before and after.
2. Write a rollback plan for each of the three methods. State when rollback is impossible without a backup.

---

## Disk fill-up incidents

When the disk that holds `PGDATA` or `pg_wal` is full, PostgreSQL can stop writes, panic, or fail to start. This is a production emergency. Topic 21 asks for a runbook. This section is the technical picture.

Common fillers:

- WAL growth: archive_command failing, replication slot holding WAL, `max_wal_size` plus a stuck checkpoint
- bloat: topic 11
- logs: `log_min_duration_statement = 0`
- `COPY` or `CREATE INDEX` on the same volume
- a forgotten `pg_basebackup` destination on the same disk

First looks (if the server still accepts a superuser connection):

```sql
SELECT pg_size_pretty(pg_database_size(current_database()));
SELECT slot_name, slot_type, active, pg_wal_lsn_diff(pg_current_wal_lsn(), restart_lsn) AS retained
FROM pg_replication_slots;
SELECT * FROM pg_stat_archiver;
```

On the host, find which directory is full (`du` on `pg_wal`, `log`, base). Free space in a way that does not destroy the only copy of WAL that you need for a replica.

Safe ideas (high-level):

- move or compress old log files if `logging_collector` owns them
- drop a disposable database or table only if you accept the data loss
- remove a stale replication slot that you proved is unused
- fix `archive_command` so that WAL can recycle

Unsafe ideas:

- delete random files in `pg_wal`
- `rm` the data directory
- `VACUUM FULL` when you have no extra space (it needs a copy)

If the server is down, free space at the OS level, then start PostgreSQL and follow the logs. Restore from backup if files are gone.

Do not wait for 100 percent disk in production without an alert. Do not delete WAL because a wiki said "old files". Topic 16 explained why WAL matters.

### Questions

#### Theoretical questions

1. Name three causes of a full data disk.
2. Why can a stale replication slot fill `pg_wal`?
3. Why is `VACUUM FULL` a poor first step on a full disk?
4. What does `pg_stat_archiver` help you prove?
5. Why must you not delete random files in `pg_wal`?

#### Easy practical tasks

1. Run `pg_database_size` and `pg_size_pretty` on your lab database.
2. Query `pg_replication_slots`.
3. Query `pg_stat_archiver`.
4. Write a five-line checklist: size, slots, archiver, logs, host `du`.

#### Medium practical tasks

1. On a lab, fill a dummy tablespace or a dummy directory (careful). Practice the queries above. Do not destroy a real cluster.
2. Write which person may drop a slot and what they must check first (replica gone, not lagging).
3. Compare log size settings with WAL size settings. Write two alerts: log disk, WAL disk.

#### Advanced practical tasks

1. Write a full incident card: detect, freeze writes if needed, free space, verify, postmortem. Topic 21 will reuse this.
2. Read "Disk Full" discussions in the official docs or wiki with the docs as the winner. Write two recovery paths: server still up, server already down.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do `work_mem`, checkpoints, and a bulk `COPY` interact on a busy host?
2. How do logs, `pg_stat_activity`, and `pg_stat_statements` form one debug path?
3. Why is a major upgrade also a memory-and-extension review?
4. How can a failed archive and a high `max_wal_size` end as a disk incident?
5. A teammate sets `shared_buffers` to 90 percent of RAM and `log_min_duration_statement` to `0`. Which facts from this topic do you use in the reply?

#### Easy practical tasks

1. `SHOW` eight settings from this topic (memory, WAL, logging). Save the list.
2. Write a cheat sheet: four memory names, checkpoint pair, slow-log, three views, three upgrade paths, disk fillers.
3. Select `pg_stat_activity` and `pg_database_size`.
4. Open the 16 and 17 runtime-config pages. Bookmark resource, WAL, and logging.

#### Medium practical tasks

1. Write a weekly ops checklist: memory `SHOW`, checkpoint stats, slow log sample, top `pg_stat_statements`, slot sizes, disk free.
2. Simulate a lock wait and a slow statement. Find both in views or logs.
3. Draft a 16-to-17 upgrade plan that includes `shared_preload_libraries` and `ANALYZE`.

#### Advanced practical tasks

1. Build a lab script that prints: version, the four memory settings, `max_wal_size`, replication slots, top five statements, database size.
2. Map this topic to official chapters: resource config, WAL, logging, monitoring, upgrading. Add one `pg_stat_io` fact for 16 or 17 that this topic did not include.
