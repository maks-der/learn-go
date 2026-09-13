# 9. Persistence

## Description

Redis can write data to disk so that a restart can reload keys. Persistence is optional. This topic covers RDB snapshots, the append-only file (AOF), `appendfsync` modes, hybrid RDB+AOF, data-loss windows, and `BGSAVE` versus `SAVE`.

Complete this topic before you study replication. Replicas do not replace a backup plan.

Use one term for each concept. A snapshot is a point-in-time RDB file. AOF is a log of write commands (or effects). `fsync` is the OS call that pushes file data to durable storage. A loss window is the time of writes that you can lose in a crash. Do not treat "persistence on" as "every `OK` is on disk".

---

## RDB snapshots

RDB is a compact binary snapshot of the dataset. Redis writes a file (often `dump.rdb`). On start, Redis loads that file into RAM. Keys that were not in the snapshot do not return.

You can trigger a snapshot with `BGSAVE` or `SAVE`. You can also set save rules in `redis.conf`:

```text
save 900 1
save 300 10
save 60 10000
```

Those rules mean: save after 900 seconds if at least 1 key changed, after 300 seconds if 10 keys changed, and after 60 seconds if 10,000 keys changed. `save ""` disables automatic RDB.

RDB is good for backups and for fast restarts of large datasets. The file is smaller than a verbose command log. You copy `dump.rdb` to backup storage.

RDB is not a per-command log. Writes after the last snapshot are missing if the process dies. The loss window is "time since last successful snapshot" plus any snapshot that did not finish.

A snapshot is a full copy of the dataset at a moment. Large datasets take time and extra RAM during `BGSAVE` (copy-on-write). Topic 19 in the full path covers the fork cost. Here, know that `BGSAVE` is not free.

`INFO persistence` shows `rdb_last_save_time`, `rdb_changes_since_last_save`, and last save status.

Do not edit an RDB file by hand. Use Redis or official tools to check the file.

### Questions

#### Theoretical questions

1. What does an RDB file contain?
2. What happens to writes that occur after the last snapshot if Redis crashes?
3. What do the three numbers in `save 300 10` mean?
4. Why is RDB useful as a backup format?
5. How do you disable automatic RDB saves?

#### Easy practical tasks

1. Run `INFO persistence`. Write `rdb_last_save_time` and `rdb_changes_since_last_save`.
2. Run `CONFIG GET save`. Save the value.
3. Open the persistence docs. Write the default RDB file name.
4. Write four sentences: RDB vs "every write on disk".

#### Medium practical tasks

1. `SET` a key. `BGSAVE`. Wait until `INFO` shows a new save time. Restart Redis. `GET` the key.
2. Change `save` on a lab instance to a short rule. Make N writes. Confirm a snapshot appears.
3. Copy `dump.rdb` out of a Docker volume (or data dir). Write the file size and the path.

#### Advanced practical tasks

1. Compare RDB file size with `used_memory` for 100,000 small keys. Write both numbers.
2. Read the RDB format overview (high level). Write three facts: magic, versions, why you must not mix versions blindly.

---

## AOF (append-only file)

AOF records writes so that Redis can replay them on start. The file grows with each write (until a rewrite). The default name is often `appendonly.aof` (Redis 7 may use a directory of AOF parts).

Enable AOF:

```text
appendonly yes
```

`CONFIG SET appendonly yes` can enable AOF on a running lab instance. Confirm on your version. Prefer `redis.conf` plus a restart for a real service.

On start, Redis replays the AOF (and may load an RDB preamble in hybrid mode). The dataset matches the log, except for a tail that was not `fsync`ed (next section).

AOF rewrite builds a shorter file that produces the same dataset. Redis can rewrite in the background (`BGREWRITEAOF`). Automatic rewrite uses `auto-aof-rewrite-percentage` and `auto-aof-rewrite-min-size`.

AOF is larger and slower to load than a compact RDB for the same data, unless you use hybrid mode. AOF can lose fewer writes than RDB if you `fsync` often.

`INFO persistence` shows `aof_enabled`, `aof_last_write_status`, and rewrite fields.

A damaged AOF can block start. Redis tools can fix the tail (`redis-check-aof`). Do not run repair on the only copy. Copy the file first.

AOF is not a query log for humans. Do not parse it as your audit system unless you build that on purpose.

### Questions

#### Theoretical questions

1. What does AOF store?
2. How does Redis use AOF at startup?
3. Why does AOF rewrite exist?
4. Which `INFO` field shows that AOF is on?
5. What should you do before you repair a broken AOF?

#### Easy practical tasks

1. Run `CONFIG GET appendonly`. Save the value.
2. Run `INFO persistence` and write `aof_enabled`.
3. Find `appendonly` in a sample `redis.conf`. Copy the line.
4. Write four sentences: AOF vs RDB.

#### Medium practical tasks

1. Enable AOF on a lab instance. `SET` a key. Restart. `GET` the key. Record the steps.
2. Run `BGREWRITEAOF`. Watch `INFO persistence` rewrite fields change.
3. Locate the AOF file or directory in the container. Write the path and size after 1,000 writes.

#### Advanced practical tasks

1. Read Redis 7 multi-part AOF notes. Write how a directory of files differs from one `appendonly.aof`.
2. Break the last bytes of a copy of an AOF in a lab (never production). Try start. Use `redis-check-aof` on the copy. Document the result.

---

## `appendfsync`: always / everysec / no

`appendfsync` controls when Redis asks the OS to flush AOF data to disk.

- `always` — `fsync` after every write. Loss window is about one command. This is the slowest mode. Latency of each write includes disk.
- `everysec` — `fsync` about once per second. You can lose about one second of writes in a crash. This is the common default. It is a balance of safety and speed.
- `no` — Redis does not `fsync` on a schedule. The OS flushes when it wants. The loss window can be longer. This is the fastest mode and the least safe.

```text
appendfsync everysec
```

`always` does not make Redis a full substitute for a disk-first database. Hardware and OS can still lose data in a power cut if the disk lies about flush. Use UPS and honest disks when the data matters.

`everysec` can stall writes if the `fsync` takes a long time (`appendfsync` and `no-appendfsync-on-rewrite` interact during rewrite). Read the docs if you see latency spikes during rewrite.

Check the running value with `CONFIG GET appendfsync`.

The application still receives `OK` after Redis applies the write in memory. `OK` is not a promise that `everysec` or `no` already flushed the disk. Only `always` waits for `fsync` on that write (with the limits above).

For a cache, AOF can stay off. For a store of record, pick `everysec` or `always` and test a kill of the process.

### Questions

#### Theoretical questions

1. What is the loss window of `appendfsync always` in the common description?
2. What is the loss window of `everysec`?
3. Who flushes the file when the mode is `no`?
4. Does `OK` mean the write is on disk in `everysec` mode?
5. Why is `always` slower?

#### Easy practical tasks

1. Run `CONFIG GET appendfsync`. Save the value.
2. Make a three-row table: mode, speed, typical loss window.
3. Write four sentences for a product owner about `everysec`.
4. Open the AOF docs. Write the default `appendfsync` that the page names.

#### Medium practical tasks

1. On a lab instance, set `everysec`, write keys, kill `-9` the process (Docker kill). Restart. Count how many of the last writes remain. Repeat a few times.
2. Compare latency of 1,000 `SET`s with `always` vs `everysec` vs `no` on local disk. Write three times.
3. Read `no-appendfsync-on-rewrite`. Write one sentence on latency during rewrite.

#### Advanced practical tasks

1. Graph or table fsync time vs `SET` p99 for `always` on a slow volume (or a Docker mount). Write a recommendation.
2. Read about `fsync` and disk caches. Write three sentences on "power loss vs Redis kill".

---

## Hybrid RDB+AOF

You can enable both RDB and AOF. Redis 7 uses an AOF that can start with an RDB preamble. The rewrite produces a binary snapshot plus a tail of commands. Startup loads the preamble, then replays the tail. Load time improves compared with a huge command-only AOF.

Older setups used RDB for backups and AOF for durability, with both files on disk. Understand your version:

- `appendonly yes` plus `save` rules
- Redis 7 `aof-use-rdb-preamble` (often yes by default)

Hybrid mode still has an `appendfsync` loss window on the AOF tail. The RDB preamble is a snapshot of an older moment. The tail is the rest.

Benefits:

- Faster restart than a long AOF replay
- Better durability than RDB only (with `everysec` or `always`)
- One restore story: enable AOF, keep RDB backups as extra copies

Costs:

- More disk
- More I/O
- More operational pieces to monitor (`INFO persistence`)

A replica can persist even when the primary persists. Persistence and replication are different. A replica with AOF is not a substitute for backups of the primary if you need point-in-time copies off the host.

When you turn AOF on for the first time, Redis can create a file from the current dataset (rewrite). Plan disk space before you enable AOF on a large instance.

### Questions

#### Theoretical questions

1. What is an RDB preamble in the AOF?
2. Why can hybrid mode start faster than a pure command AOF?
3. Does hybrid mode remove the `appendfsync` loss window?
4. What extra resource does hybrid mode use compared with RDB only?
5. Why enable AOF for the first time only after you check disk space?

#### Easy practical tasks

1. Run `CONFIG GET aof-use-rdb-preamble` if the key exists. Save the value or "missing".
2. Write four sentences: hybrid vs RDB only vs AOF only.
3. List three `INFO persistence` fields you would watch in hybrid mode.
4. Open Redis 7 persistence notes. Write one sentence about the AOF directory.

#### Medium practical tasks

1. Enable AOF and keep `save` rules. Trigger `BGSAVE` and `BGREWRITEAOF` on a lab instance. List the files in the data directory.
2. Restart and time the start (or watch logs). Write a rough duration for your dataset size.
3. Draw the load order: preamble, tail, ready for clients.

#### Advanced practical tasks

1. Compare start time of the same dataset with AOF off (RDB only), AOF without preamble if you can, and hybrid. Write a table.
2. Write a restore runbook: copy files, set `dir`, start Redis, `INFO keyspace`. Include a checksum step.

---

## Restart and data loss expectations

Write a clear contract for each environment.

Cache instance:

- Persistence can be off
- Restart loses all keys
- The application rebuilds from the source of truth
- Loss is expected

Durable instance:

- AOF on, `appendfsync everysec` or `always`
- RDB or hybrid for faster load and backups
- A process kill can lose up to about one second (`everysec`) or rare disk-level loss (`always`)
- A host disk failure loses data unless you have off-host backups and/or replicas

Test the contract. Do not guess.

Tests:

1. `SET` a marker key
2. `BGSAVE` or wait for AOF fsync
3. Restart Redis cleanly (`SHUTDOWN` or Docker stop)
4. `GET` the marker
5. Repeat with a hard kill after a burst of writes

`SHUTDOWN` can write a final RDB or AOF depending on configuration (`SHUTDOWN NOSAVE` does not). Read the `SHUTDOWN` page. A crash or `kill -9` skips a graceful save.

Replicas that were not promoted do not help if you destroy the only copy and you have no replica. Backups are files that you copied away.

Document the recovery time (load of RDB/AOF) for your size. A 20 GB dataset is not instant.

Clients must handle connection errors during restart. The application should retry with a cap.

### Questions

#### Theoretical questions

1. What is the expected data loss for a cache with persistence off?
2. What is the common loss window for AOF `everysec` after a crash?
3. How is a graceful `SHUTDOWN` different from `kill -9`?
4. Why is a replica not enough if you never copy files off the host?
5. Why must you measure restart time?

#### Easy practical tasks

1. Write a three-line contract for your lab instance (cache or durable).
2. Run `SHUTDOWN NOSAVE` only on a disposable lab container. Confirm data is gone after start. Recreate the container if needed.
3. Write four sentences that you would put in a runbook about loss.
4. Find `SHUTDOWN` options on the command page. List two.

#### Medium practical tasks

1. Perform the five-step test (marker, persist, clean restart, `GET`, hard kill burst). Record both outcomes.
2. Estimate load time from docs or a trial for your key count. Write the number.
3. List client retry settings you would use during a 10-second Redis restart.

#### Advanced practical tasks

1. Write an incident table: crash, disk full, bad AOF, accidental `SHUTDOWN NOSAVE`. For each, data loss and first command.
2. Automate a backup copy of RDB/AOF to another folder after `BGSAVE`. Restore into a second instance. Prove a key.

---

## `BGSAVE` vs blocking `SAVE`

`SAVE` writes an RDB file on the main thread. Redis does not run other commands until `SAVE` finishes. Clients wait. On a large dataset, `SAVE` can stall the instance for seconds or minutes. Do not use `SAVE` in production.

`BGSAVE` forks a child process. The child writes the RDB file. The parent keeps serving commands. The operating system uses copy-on-write (COW). If many keys change during the save, RSS can grow toward almost two copies of the dataset. Leave free RAM for that peak.

```text
BGSAVE
```

`LASTSAVE` returns the Unix time of the last successful snapshot.

`INFO persistence` shows `rdb_bgsave_in_progress`. Do not start another `BGSAVE` if one is already running unless you know the server allows it (it usually refuses).

`SHUTDOWN` can run a save. `BGREWRITEAOF` is the AOF analog of a background rewrite (also a fork).

Use `BGSAVE` for on-demand backups. Schedule automatic `save` rules or an external backup that copies the file after `BGSAVE` completes.

If `BGSAVE` fails (`fork` error, disk full), fix the cause. Redis can be configured to stop writes when the last `BGSAVE` failed (`stop-writes-on-bgsave-error`). That flag protects you from a silent "no snapshots" state. It can also block a cache. Know the setting.

### Questions

#### Theoretical questions

1. Why does `SAVE` block clients?
2. How does `BGSAVE` keep serving commands?
3. What RAM risk does copy-on-write create during `BGSAVE`?
4. What does `LASTSAVE` return?
5. What does `stop-writes-on-bgsave-error` protect you from?

#### Easy practical tasks

1. Run `LASTSAVE`. Save the number. Convert it to a date if you can.
2. Run `BGSAVE`. Check `rdb_bgsave_in_progress` then `LASTSAVE` again.
3. Write four sentences: `SAVE` vs `BGSAVE`.
4. Run `CONFIG GET stop-writes-on-bgsave-error`. Save the value.

#### Medium practical tasks

1. Time `SAVE` on a lab dataset of at least 10,000 keys. Then time `BGSAVE` and whether `PING` works during `BGSAVE`.
2. Fill disk in a disposable volume until `BGSAVE` fails (or simulate by setting a tiny quota). Record the error and `INFO`.
3. Read fork / COW notes in the docs. Draw parent, child, and shared pages.

#### Advanced practical tasks

1. During `BGSAVE`, rewrite many keys to force COW. Watch `used_memory_rss`. Write the peak vs baseline.
2. Write a backup script: `BGSAVE`, wait until `rdb_bgsave_in_progress` is 0, copy the RDB, checksum, exit non-zero on failure.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do RDB, AOF, and `appendfsync` combine into one durability story?
2. What can you lose in a crash with RDB only, with AOF `everysec`, and with persistence off?
3. When do you choose `BGSAVE` instead of `SAVE`, and when do you rewrite AOF?
4. What does hybrid persistence change about startup, and what does it not change about fsync?
5. How do you prove your restart contract with tests?

#### Easy practical tasks

1. Export `INFO persistence` to a file. Mark RDB fields, AOF fields, and last error fields.
2. Write a cheat sheet: `save` rules, `appendonly`, `appendfsync`, `BGSAVE`, `SAVE`, `BGREWRITEAOF`, `SHUTDOWN`.
3. Draw a timeline: write, `OK`, fsync (or snapshot), crash, restart, `GET`.
4. Write three lab rules: no `SAVE` in prod, copy files before repair, test kill -9.

#### Medium practical tasks

1. Configure a disposable instance with hybrid persistence and `everysec`. Run a 1,000-write test and a restart test. Write the results.
2. Document disk paths and backup copy commands for Docker Redis.
3. Write a one-page persistence policy for a cache Redis and a durable Redis.

#### Advanced practical tasks

1. Restore the same RDB on a second Redis version (or read compatibility notes). Write whether the file loaded and the risk of version skew.
2. Measure start time and peak RSS for `BGSAVE` on a 1 GB-class dataset if you can build one. Write hardware notes. If you cannot, write a sizing plan instead.
