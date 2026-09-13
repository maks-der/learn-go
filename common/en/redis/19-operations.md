# 19. Operations

## Description

Operations keep Redis correct and available after the application works. This topic covers runtime config, disk fill, the `BGSAVE` memory spike, upgrades, live metrics, and file backups.

Complete topic 9 (persistence), topic 10 (replication), topic 16 (memory), and topic 17 (security) first.

Use one term for each concept. Runtime config is a setting that `CONFIG SET` changes in the live process. Disk fill is a full volume for RDB or AOF. Copy-on-write (COW) is the extra RAM during a fork. Observability is the set of metrics that tell you if Redis is healthy. A backup is a copy of RDB or AOF that lives off the Redis host.

---

## `CONFIG GET` / `SET` (careful)

`CONFIG GET pattern` reads live settings. `CONFIG SET name value` changes a setting until restart, unless you persist it.

```text
CONFIG GET maxmemory
CONFIG SET maxmemory 256mb
CONFIG REWRITE
```

`CONFIG REWRITE` writes the current settings into `redis.conf` (or the config file that Redis loaded). Without a rewrite or a file edit, a restart can restore old values.

Some settings do not apply through `CONFIG SET`. They need a restart. The command page and `CONFIG GET` comments tell you. Examples often include bind addresses, daemonize, and some TLS file paths. Check your version.

Risks:

- `CONFIG SET` is a dangerous command class (topic 17). Restrict it with ACLs.
- A typo in `dir` or `dbfilename` can write files to the wrong place.
- `maxmemory` too low starts eviction or write errors.
- `appendonly` and `save` changes change the loss window (topic 9).
- `protected-mode` and `requirepass` changes can lock you out or open the instance.

Prefer a change ticket, a lab test, and `CONFIG GET` before and after. Copy the old value.

`redis.conf` remains the source of truth for a new process. Document every `CONFIG SET` that you do not rewrite.

Do not experiment with `CONFIG SET` on a durable production primary.

### Questions

#### Theoretical questions

1. What is the difference between `CONFIG SET` and an edit of `redis.conf` without restart?
2. What does `CONFIG REWRITE` do?
3. Why do some settings need a restart?
4. Why is `CONFIG SET` restricted by ACL?
5. What happens if you lower `maxmemory` below `used_memory`?

#### Easy practical tasks

1. `CONFIG GET maxmemory`. Save the value.
2. `CONFIG GET save`. Save the value.
3. Write four sentences: `CONFIG SET` vs restart-only settings.
4. Open the `CONFIG SET` page. Write one setting that the page marks as not changeable at runtime (if listed).

#### Medium practical tasks

1. On a lab instance, change `slowlog-log-slower-than`. `CONFIG GET` it. Restart without `CONFIG REWRITE`. Write whether the value survived.
2. Repeat with `CONFIG REWRITE` if your lab file is writable. Confirm after restart.
3. List five settings you would never change on a live primary without a ticket.

#### Advanced practical tasks

1. Write an ops SOP: who may `CONFIG SET`, required `GET` before/after, rewrite yes/no, rollback.
2. Compare `CONFIG GET *` size on your instance. Write how you would scrape a subset for a dashboard (names only).

---

## Persistence disk fill

RDB and AOF write files under `dir`. If the volume is full, `BGSAVE` fails, AOF append fails, or rewrite fails. Redis can stop writes (`stop-writes-on-bgsave-error`, AOF error behavior). A cache may then reject `SET`. A durable Redis may refuse to lose data.

```text
CONFIG GET dir
INFO persistence
```

Watch disk used, inode count, and Redis log lines about write errors.

AOF rewrite and RDB write need extra space. A rewrite can need room for the new file beside the old file. A volume at 95 percent can fail even if the final file would fit.

`INFO persistence` fields such as `rdb_last_bgsave_status` and `aof_last_write_status` (names vary) tell you the last result.

If the disk is full:

- Free space (old dumps, logs) on a lab or with a known backup policy
- Do not delete the only live AOF without a plan
- Fix the volume size before you raise traffic

Replication does not replace a full disk on the primary. The primary still must write if persistence is on.

Do not put `dir` on a tiny container overlay disk without a mount.

### Questions

#### Theoretical questions

1. Where does Redis write RDB and AOF?
2. Why can a volume at 95 percent still fail a rewrite?
3. What can Redis do to writes when a save fails?
4. Which `INFO` section do you read for the last save status?
5. Why is deleting the live AOF dangerous?

#### Easy practical tasks

1. `CONFIG GET dir`. Save the path.
2. `INFO persistence`. Write `aof_enabled` and one `rdb_*` status field.
3. Write four sentences: full disk vs `maxmemory`.
4. List two files you expect under `dir` when AOF and RDB are on.

#### Medium practical tasks

1. On a disposable volume, fill the disk until `BGSAVE` fails (or set a tiny quota). Record the error and `rdb_last_bgsave_status`. Then free space.
2. Estimate space for AOF rewrite: current AOF size times two. Write the number for your lab file (`ls` or `dir`).
3. Draw primary disk, replica disk, and backup store. Label what must not fill.

#### Advanced practical tasks

1. Write disk alerts: percent full, inodes, `rdb_last_bgsave_status`, AOF write status. Include a first action.
2. Test `stop-writes-on-bgsave-error` yes vs no on a lab cache vs a lab durable instance. Write which setting you choose for each role.

---

## Fork / copy-on-write during `BGSAVE` (memory spike)

Topic 9 explained `BGSAVE`. The parent forks a child. The child writes the RDB. The operating system shares pages until the parent writes a page. Then COW copies that page. Many writes during a long save can raise RSS toward almost two copies of the dataset.

Effects:

- The host can swap or the Linux OOM killer can stop Redis
- Latency can rise during the fork on some kernels
- AOF rewrite (`BGREWRITEAOF`) also forks

Leave free RAM for the peak. A common planning rule is to keep headroom for a full extra copy if the write rate is high. Measure RSS during save. Do not guess only from `used_memory`.

```text
INFO memory
INFO persistence
```

Watch `used_memory_rss` while `rdb_bgsave_in_progress` is `1`.

Linux `overcommit_memory` and `vm.overcommit_memory` notes appear in Redis warnings. A `fork` error in the log means the kernel refused the child. Fix RAM, overcommit, or dataset size.

`SAVE` does not fork. `SAVE` blocks the main thread. Do not use `SAVE` in production.

Do not schedule a heavy write load and a manual `BGSAVE` at the same time on a tight host.

### Questions

#### Theoretical questions

1. Why does RSS grow during `BGSAVE` if the parent writes many keys?
2. What process writes the RDB file?
3. Why is `SAVE` worse for clients than `BGSAVE`?
4. What does a `fork` error in the log often mean?
5. Which other command family also forks?

#### Easy practical tasks

1. Run `BGSAVE`. Watch `rdb_bgsave_in_progress`. Write `LASTSAVE` after.
2. Write four sentences: fork vs `SAVE` on the main thread.
3. Record `used_memory_rss_human` idle.
4. Open Redis warnings about overcommit in the logs or docs. Write one sentence.

#### Medium practical tasks

1. During `BGSAVE`, rewrite a large set of keys in a lab. Record RSS peak vs idle.
2. Compare time of `PING` during `BGSAVE` vs idle. Write both.
3. Read `INFO persistence` child fields (`rdb_current_bgsave_time_sec` or similar). Write the names you see.

#### Advanced practical tasks

1. Size a host for 8 GB `used_memory` and a busy write rate. Write RAM and disk numbers with a COW caveat.
2. Write a playbook: `fork` failed, OOM killer, swap storm. First checks only (logs, `INFO`, free RAM).

---

## Upgrades

An upgrade changes the Redis binary or image. Plan it.

Typical replica-first path for a primary plus replicas:

1. Upgrade a replica. Check replication and `INFO`.
2. Upgrade the other replicas.
3. Fail over (Sentinel or Cluster) so that an upgraded replica becomes primary.
4. Upgrade the old primary (now a replica).
5. Test the application against new command behavior.

Read the release notes. Some versions change defaults (ACL, AOF, replica names). Some commands change.

Modules: upgrade Redis Stack as a bundle. Check module compatibility (topic 18).

Config: new versions add settings. Diff `redis.conf`. Do not copy an old file blindly if the new image expects new includes.

Clients: upgrade clients after you know the server version, or stay compatible. Test pipelines, Cluster, and AUTH.

Rollback: keep the previous image and a backup from before the upgrade. A data-format change can make rollback harder. Read the notes.

Do not upgrade the only primary in place during peak without a replica and a backup.

Kubernetes and cloud vendors have their own rolling procedures. Follow those plus the Redis notes.

### Questions

#### Theoretical questions

1. Why upgrade a replica before the primary?
2. What must you read before you jump a major version?
3. Why do modules change the upgrade?
4. What do you keep for rollback?
5. Why can a data-format change block a rollback?

#### Easy practical tasks

1. Write `redis_version` from `INFO server`.
2. Open the release notes for the next minor version. Write one change.
3. Write four sentences: rolling replica upgrade vs "replace the only container".
4. List three tests you run after an upgrade (`PING`, one `GET`, replication `INFO`).

#### Medium practical tasks

1. In Docker, start Redis version N as primary and N as replica (or two close versions if you can). Document `INFO replication`. Then replace the replica image. Confirm the link.
2. Diff two `redis.conf` files from two versions (or docs). Write three new or changed settings.
3. Write a change ticket template: version from/to, backup, failover, client tests, rollback.

#### Advanced practical tasks

1. Perform a lab failover after a replica upgrade (Sentinel or Cluster or manual `REPLICAOF`). Document the command order.
2. Write an upgrade matrix: server, Stack modules, go-redis (or your client), OS. Fill current and target.

---

## Observability: hit rate, evictions, connected clients, blocked clients

You cannot operate Redis from `PING` alone. Scrape `INFO` (or a vendor dashboard) on a interval.

Hit rate (cache):

```text
# conceptual
hit_rate = keyspace_hits / (keyspace_hits + keyspace_misses)
```

A low hit rate means the cache does not pay for itself, or TTLs are too short, or keys do not match. A high hit rate is not enough if the value is stale (application issue).

Evictions: `evicted_keys` rising means `maxmemory` and a policy are dropping keys (topic 4). For a cache, some eviction can be normal. For a durable Redis, eviction is an incident.

Connected clients: `connected_clients` vs `maxclients`. A leak in the application pool (topic 12) shows here. A sudden drop can mean a network cut or a restart.

Blocked clients: `blocked_clients` counts clients in blocking commands (`BRPOP`, `XREAD BLOCK`, `BZPOPMIN`). A number above zero can be normal. A number that equals all workers can mean they wait on an empty queue.

Other useful fields:

- `instantaneous_ops_per_sec`
- `rejected_connections`
- `expired_keys` vs `evicted_keys`
- `connected_slaves` / replica count
- `rdb_last_bgsave_status`
- `used_memory` vs `maxmemory`
- `latest_fork_usec`

Alert on: eviction on a durable instance, failed saves, `rejected_connections`, memory near `maxmemory`, replication lag (topic 10), certificate expiry (topic 17).

Do not alert on every `blocked_clients` > 0 if you use blocking reads.

### Questions

#### Theoretical questions

1. How do you compute a simple hit rate?
2. When is a rising `evicted_keys` an incident?
3. What does `connected_clients` near `maxclients` suggest?
4. When is `blocked_clients` normal?
5. Name three other `INFO` fields that belong on a dashboard.

#### Easy practical tasks

1. Run `INFO stats`. Write `keyspace_hits`, `keyspace_misses`, `evicted_keys`.
2. Run `INFO clients`. Write `connected_clients` and `blocked_clients`.
3. Write four sentences: expired vs evicted counters.
4. Compute hit rate if hits=90 and misses=10.

#### Medium practical tasks

1. Generate misses (`GET` missing keys) and hits. Show the two counters move.
2. Set a tiny `maxmemory` and a write load until `evicted_keys` grows (lab). Write the policy and the counter.
3. Run `BLPOP` with a timeout in one client. Watch `blocked_clients`. Write the number during the wait.

#### Advanced practical tasks

1. Design a dashboard with eight panels. Write the `INFO` field and the alert threshold idea for each.
2. Correlate a deploy with `connected_clients` and hit rate in a lab (restart app). Write what you expect to see.

---

## Backup of RDB/AOF

A replica is not a backup. A replica can delete data if the primary does. A backup is a file copy that you can restore after a bad `FLUSHALL`, a bad deploy, or a lost host.

Steps (shape):

1. Ensure persistence is on if you need history.
2. Trigger `BGSAVE` or wait for the scheduled save. Wait until `rdb_bgsave_in_progress` is `0` and status is OK.
3. Copy the RDB file off the host (object storage, another disk).
4. Checksum the copy.
5. If AOF is the source of truth, copy AOF after a rewrite or follow the official point-in-time notes for your version.
6. Test restore on a lab instance.

```text
CONFIG GET dir
CONFIG GET dbfilename
CONFIG GET appendfilename
```

Encrypt backups if they contain personal data. Restrict who can read the bucket.

Schedule: hourly or daily depends on the loss window. A cache may need no backup. A session store may need a short window. A durable Redis needs a tested restore.

Cluster: back up every primary (or use a vendor snapshot). Slots live on different nodes.

Do not only copy files while `BGSAVE` is in progress unless you know the file is consistent. Prefer wait-for-complete.

`DEBUG RELOAD` and similar are lab tools. Restore by starting Redis with the copied file in `dir` on an empty lab.

### Questions

#### Theoretical questions

1. Why is a replica not a backup?
2. When do you copy the RDB file?
3. Why do you checksum the copy?
4. Why encrypt backups that hold personal data?
5. What extra work does Cluster add?

#### Easy practical tasks

1. `CONFIG GET dir` and `dbfilename`. Write the full path.
2. `BGSAVE`. Wait. Copy the RDB to a lab folder (same machine is fine). Write the copy command.
3. Write four sentences: replica vs backup.
4. Write whether your learning instance needs backups (cache vs durable).

#### Medium practical tasks

1. Restore the copied RDB on a second Docker port. Confirm a known key. Write both commands.
2. Enable AOF. Write a key. Copy AOF after `BGREWRITEAOF`. Document file sizes.
3. Write a backup calendar: RPO (loss window) and RTO (restore time) for a cache and for a durable store.

#### Advanced practical tasks

1. Automate: `BGSAVE`, wait, copy, checksum, upload to a dummy folder, exit non-zero on fail. Test on lab data.
2. Drill a `FLUSHALL` recovery from backup on a lab (not production). Write the timeline and the gap since last backup.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do `CONFIG SET`, disk fill, and COW interact during one `BGSAVE` on a small host?
2. Which metrics tell you "cache is fine" vs "durable Redis is losing data"?
3. Why do upgrade order and backup exist in the same change window?
4. What must be true before you treat an RDB file as restorable?
5. Which operations tasks require ACL ops rights that the application must not have?

#### Easy practical tasks

1. Write a cheat sheet: `CONFIG GET/SET/REWRITE`, `INFO persistence`, `INFO memory`, `INFO stats`, `BGSAVE`, copy RDB.
2. Export `INFO` and highlight hit counters, evictions, clients, last save status, RSS.
3. Draw RAM (used + COW headroom) and disk (RDB + rewrite room).
4. Write five ops rules for your lab.

#### Medium practical tasks

1. Run a mini drill: change a harmless config, `BGSAVE`, copy RDB, restore on another port, read a key. Write the times.
2. Write an on-call one-pager: first 10 minutes for high memory, full disk, failed save, too many clients.
3. Map topic 10 lag and topic 16 slow log onto this topic's dashboard.

#### Advanced practical tasks

1. Write a full Redis operations runbook: config, disk, COW, upgrade, metrics, backup, security contacts.
2. Compare a vendor automated backup with your file-copy procedure. Write what you still test yourself (restore).
