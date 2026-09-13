# 16. Backup and Replication

## Description

This topic shows backup and copies of a cluster in PostgreSQL 16 and PostgreSQL 17. You learn `pg_dump`, `pg_restore`, `pg_dumpall`, `pg_basebackup`, WAL archiving, point-in-time recovery, streaming replication, hot standby, logical replication, and failover tools at a high level.

Complete topic 15 first. You can detach a partition and dump a table. Complete this topic before you tune client pools.

Use one term for each concept. A logical dump is SQL or an archive of objects. A physical backup is a copy of the data directory plus WAL. WAL is the write-ahead log. PITR is recovery to a chosen time. A streaming replica is a physical stand-by that replays WAL. Logical replication copies a set of tables through publications and subscriptions. This handbook targets PostgreSQL 16 and PostgreSQL 17.

---

## `pg_dump` / `pg_restore` / `pg_dumpall`

`pg_dump` copies one database. It does not copy roles or tablespaces. `pg_dumpall` copies all databases plus cluster-wide objects (roles, tablespaces). `pg_restore` loads a non-plain archive from `pg_dump`.

Formats for `pg_dump`:

- plain (`-Fp`) — a SQL script; load with `psql`
- custom (`-Fc`) — compressed archive; load with `pg_restore`; you can pick objects
- directory (`-Fd`) — one folder; parallel dump and restore (`-j`)
- tar (`-Ft`) — tar archive; load with `pg_restore`

```text
pg_dump -h localhost -U postgres -d shop -Fc -f shop.dump
pg_restore -h localhost -U postgres -d shop_copy --clean --if-exists shop.dump
```

Schema only or data only:

```text
pg_dump -s -d shop -f shop_schema.sql
pg_dump -a -d shop -Fc -f shop_data.dump
```

One table (topic 15 archive):

```text
pg_dump -d shop -t shop.orders_2025_01 -Fc -f orders_2025_01.dump
```

`pg_dump` runs a consistent snapshot of that database. Other databases in the cluster are not in the file. Open transactions that hold old snapshots can delay the dump (topic 10 and topic 11).

```text
pg_dumpall -g -f roles.sql
```

`-g` dumps only globals. Restore globals before you restore a database that depends on those roles.

Match client and server major versions when you can. A newer `pg_dump` can often dump an older server. Restoring a 17 dump onto 16 can fail on new features.

Do not treat a filesystem copy of a running data directory as a backup. That copy is not consistent. Use `pg_dump` or `pg_basebackup`. Do not commit dump files that contain live data to git.

Official: [https://www.postgresql.org/docs/17/app-pgdump.html](https://www.postgresql.org/docs/17/app-pgdump.html), [https://www.postgresql.org/docs/17/app-pgrestore.html](https://www.postgresql.org/docs/17/app-pgrestore.html), [https://www.postgresql.org/docs/17/app-pg-dumpall.html](https://www.postgresql.org/docs/17/app-pg-dumpall.html).

### Questions

#### Theoretical questions

1. What does `pg_dump` omit that `pg_dumpall` includes?
2. Which formats do you load with `pg_restore`?
3. What is a consistent snapshot for `pg_dump`?
4. What does `pg_dumpall -g` write?
5. Why is a copy of a live data directory not a backup?

#### Easy practical tasks

1. `pg_dump -s` one lab database to a SQL file. Open the file and find `CREATE TABLE`.
2. `pg_dump -Fc` the same database. Run `pg_restore -l` on the dump.
3. Dump one table with `-t`. Restore it into a new database.
4. Write four sentences: `pg_dump`, `pg_restore`, `pg_dumpall`, format.

#### Medium practical tasks

1. Create database `shop_copy`. Restore a custom dump into it. Compare table counts.
2. Use directory format with `-j 2` if the machine has two cores. Write the command pair dump and restore.
3. Dump globals with `pg_dumpall -g`. Find `CREATE ROLE` in the file.

#### Advanced practical tasks

1. Restore only one table from a custom dump (`pg_restore -t`). Prove other tables did not load.
2. Read version compatibility in the `pg_dump` docs for 16 and 17. Write which direction (old server, new dump tool) the docs allow.

---

## `pg_basebackup`

`pg_basebackup` takes a physical backup of the whole cluster. It copies the data directory through the replication protocol. The backup is a starting point for a replica or for PITR.

```text
pg_basebackup -h primary.example -U replicator -D /var/lib/postgresql/17/backup -Fp -X stream -R
```

Useful flags:

- `-D` — destination directory
- `-Fp` or `-Ft` — plain files or tar
- `-X stream` — include WAL that arrives during the copy
- `-R` — write `standby.signal` and connection settings for a replica
- `-c fast` — fast checkpoint before the copy (more WAL, shorter start wait)

The role needs `REPLICATION` (or superuser). `pg_hba.conf` needs a `replication` line.

PostgreSQL 17 adds incremental physical backup. You take a full `pg_basebackup`, then later `pg_basebackup --incremental` against a prior manifest. `pg_combinebackup` builds a usable directory from the full backup plus incrementals. PostgreSQL 16 does not have this incremental tool set. On 16, use full `pg_basebackup` or a file-level incremental tool that understands PostgreSQL.

A physical backup is cluster-wide. You cannot restore one database from `pg_basebackup` without extra work. Use `pg_dump` when you need one database.

Do not run two primaries from copies of the same backup without a recovery plan. They will diverge and break replicas. Do not put the backup directory inside the live data directory.

Official: [https://www.postgresql.org/docs/17/app-pgbasebackup.html](https://www.postgresql.org/docs/17/app-pgbasebackup.html). PostgreSQL 17 also documents [https://www.postgresql.org/docs/17/app-pgcombinebackup.html](https://www.postgresql.org/docs/17/app-pgcombinebackup.html).

### Questions

#### Theoretical questions

1. What does `pg_basebackup` copy?
2. What privilege does the backup role need?
3. What does `-R` prepare?
4. What incremental feature exists in 17 that 16 does not ship?
5. Can you restore one database from a base backup as easily as from `pg_dump`?

#### Easy practical tasks

1. Read `\h` is not enough; open the `pg_basebackup` reference. List five options from this section.
2. Write a command that backs up `localhost` to a new directory (do not overwrite your live data directory).
3. Show `SELECT pg_is_in_recovery();` on your current server.
4. Write four sentences: physical backup, WAL stream, replica, incremental 17.

#### Medium practical tasks

1. Create a `REPLICATION` role. Add a `pg_hba` replication line in a lab. Run `pg_basebackup` to an empty folder.
2. Compare the size of a `pg_dump -Fc` of one database with a `pg_basebackup` of the cluster.
3. On a 17 lab only, read `--incremental` in the docs and write the extra files that a full backup must keep.

#### Advanced practical tasks

1. Start a second data directory from `-R` (next sections). Confirm `pg_is_in_recovery()` is true on the copy.
2. Compare 16 and 17 backup chapters. Write a team rule: which major version may use `pg_combinebackup`.

---

## WAL archiving and PITR

The server writes WAL before it writes data pages. If you keep every WAL segment from a base backup onward, you can replay to a target time. That process is point-in-time recovery (PITR).

Primary settings (high-level):

- `wal_level = replica` (or `logical` if you also need logical decoding)
- `archive_mode = on`
- `archive_command` — a shell command that copies `%p` to a safe store, or `archive_library` in current 16/17 installs

Example idea (adjust paths; do not copy a broken command into production):

```text
archive_command = 'test ! -f /archive/%f && cp %p /archive/%f'
```

`%p` is the path of the file. `%f` is the file name. The command must return success only when the copy is durable.

Recovery on PostgreSQL 16 and 17 uses `postgresql.conf` (or `postgresql.auto.conf`) plus a signal file:

- put `restore_command` in the config
- create empty `recovery.signal` for a recovery that becomes a primary
- create empty `standby.signal` to stay a stand-by

There is no `recovery.conf` in these versions. That file is an old pattern.

Target examples:

```text
recovery_target_time = '2026-09-13 16:00:00+00'
recovery_target_action = 'promote'
```

You restore a base backup, add WAL from the archive, set the target, and start the server.

Do not set `archive_command` to a command that can succeed without a real copy. Do not delete WAL from the archive before every base backup that still needs those files is expired. Do not test PITR only in your head. Restore into an empty directory.

Official: [https://www.postgresql.org/docs/17/continuous-archiving.html](https://www.postgresql.org/docs/17/continuous-archiving.html).

### Questions

#### Theoretical questions

1. What must you keep besides a base backup for PITR?
2. What do `%p` and `%f` mean in `archive_command`?
3. Which signal file do you use for recovery that promotes?
4. Where did `recovery.conf` go in PostgreSQL 16 and 17?
5. What does `recovery_target_time` select?

#### Easy practical tasks

1. `SHOW wal_level;` `SHOW archive_mode;` `SHOW archive_command;`.
2. Open the 17 continuous-archiving chapter. Write the names of the two signal files.
3. Write a restore checklist: restore backup, configure `restore_command`, signal file, target time, start.
4. Find `pg_wal` in the docs. Write what that directory holds.

#### Medium practical tasks

1. In a lab, set `archive_mode` and a local `archive_command` that copies to a folder. Switch WAL with `SELECT pg_switch_wal();`. List new files in the archive folder.
2. Plan a PITR drill with a written target time. Do not skip the empty-directory restore.
3. Compare `wal_level` values `minimal`, `replica`, and `logical` in the docs. Write which values allow archiving for PITR.

#### Advanced practical tasks

1. Restore a base backup and replay archive WAL to a target time on a lab cluster. Prove with a row that you inserted after the target (it must be absent) and a row before the target (it must be present).
2. Read `archive_library` versus `archive_command` in the 16 or 17 docs. Write when a team picks each.

---

## Streaming replication

Streaming replication sends WAL from a primary to a stand-by over a replication connection. The stand-by writes WAL and replays it. The copy is physical. The major version must match. The stand-by is read-only for normal SQL.

High-level steps:

1. Create a `REPLICATION` role on the primary.
2. Allow `replication` in `pg_hba.conf`.
3. Set `wal_level` to at least `replica`.
4. Take `pg_basebackup -R` into the stand-by data directory.
5. Start the stand-by.

The stand-by uses `primary_conninfo` (often written by `-R`) and `standby.signal`.

Watch lag:

```sql
-- on primary
SELECT client_addr, state, sent_lsn, replay_lsn
FROM pg_stat_replication;

-- on stand-by
SELECT pg_is_in_recovery(), pg_last_wal_replay_lsn();
```

Synchronous replication waits for a stand-by to flush WAL (`synchronous_commit` and `synchronous_standby_names`). That setting reduces data loss and increases write latency. Beginners start with asynchronous streaming.

Slots (`pg_create_physical_replication_slot`) keep WAL on the primary until the replica consumes it. A slot that nobody consumes fills the disk. Topic 18 covers disk-full incidents.

Do not write on the stand-by. Do not use streaming as a backup by itself; you still need a base backup and often an archive. Do not mix PostgreSQL 16 primary with a 17 stand-by.

Official: [https://www.postgresql.org/docs/17/warm-standby.html](https://www.postgresql.org/docs/17/warm-standby.html).

### Questions

#### Theoretical questions

1. What does streaming replication send?
2. Must the major version of primary and stand-by match?
3. What file marks a data directory as a stand-by?
4. Where do you read replica lag on the primary?
5. What risk does a forgotten replication slot create?

#### Easy practical tasks

1. Query `pg_stat_replication` on your lab (empty is acceptable).
2. Write the five high-level steps from this section.
3. `SHOW wal_level;` and write whether streaming is possible.
4. Open the 17 stand-by chapter. Write the section title that covers streaming.

#### Medium practical tasks

1. Build a primary and a replica with Docker or two data directories. Insert on the primary. Select on the replica.
2. Compare `sent_lsn` and `replay_lsn` after a burst of inserts.
3. Create a physical slot. Show it in `pg_replication_slots`. Drop it when you finish the lab.

#### Advanced practical tasks

1. Break the replica network for a minute. Watch lag. Restore the network. Watch catch-up.
2. Read synchronous replication in the docs. Write one sentence on data-loss versus latency. Do not enable it on a single-node laptop except as a short test.

---

## Hot standby and `hot_standby`

A hot stand-by is a replica that accepts read-only SQL while it recovers. `hot_standby` is `on` by default in PostgreSQL 16 and 17.

```sql
SHOW hot_standby;
SELECT pg_is_in_recovery();
```

On a hot stand-by you can `SELECT`. You cannot `INSERT`, `UPDATE`, `DELETE`, or DDL that writes.

Conflicts: the primary can vacuum or update a row that a long stand-by query still needs. The stand-by can cancel the query (`max_standby_streaming_delay`). A reporting query that runs for a long time can fail on a busy replica.

`hot_standby_feedback` tells the primary about stand-by snapshots. It reduces query cancels. It can delay vacuum on the primary (more bloat). Topic 11 explained dead tuples.

Use the replica for reporting that tolerates lag and cancel. Topic 19 covers that pattern.

Do not set `hot_standby = off` on a replica that applications must query. Do not point a write API at a stand-by.

### Questions

#### Theoretical questions

1. What does `hot_standby = on` allow?
2. What statements fail on a stand-by?
3. Why can a long `SELECT` on a replica fail?
4. What does `hot_standby_feedback` trade?
5. How do you test that a session is in recovery?

#### Easy practical tasks

1. `SHOW hot_standby;`
2. If you have a replica, run `SELECT 1;` and try `CREATE TABLE`. Record the error.
3. Write four sentences: hot stand-by, read only, conflict, feedback.
4. Find `max_standby_streaming_delay` with `SHOW`.

#### Medium practical tasks

1. On a replica, start a long `SELECT` (or `pg_sleep` in a transaction that already read a row). On the primary, update that row many times and vacuum. Write what happens to the replica query.
2. Compare `hot_standby_feedback` on and off in the docs. Write one ops rule.
3. Connect a GUI to the replica. Confirm that the session is read-only.

#### Advanced practical tasks

1. Document which application roles may connect to the replica. Include lag and cancel risk.
2. Read "Hot Standby" conflicts in the 17 docs. List two conflict types and one setting that relates to each.

---

## Logical replication and publications / subscriptions

Logical replication copies row changes for a set of tables. The publisher and subscriber can be different major versions in many supported pairs. The subscriber table is writable. Conflicts are your problem if both sides write the same key.

On the publisher:

```sql
CREATE PUBLICATION shop_pub FOR TABLE shop.items;
-- or FOR ALL TABLES
```

On the subscriber:

```sql
CREATE SUBSCRIPTION shop_sub
    CONNECTION 'host=primary port=5432 dbname=shop user=replicator'
    PUBLICATION shop_pub;
```

`wal_level` must be `logical` on the publisher. The publication can list tables, or all tables. PostgreSQL 15 and later (so 16 and 17) can add column lists and row filters on publications. Read the `CREATE PUBLICATION` page before you rely on a filter.

Initial sync copies existing rows, then applies changes. A replication slot on the publisher retains WAL for the subscriber. A dead subscriber fills the disk.

PostgreSQL 16 can perform logical decoding from a stand-by in supported setups. PostgreSQL 17 adds more failover helpers for logical slots (`pg_createsubscriber` and related tools). Treat those as operator features. Learn `PUBLICATION` and `SUBSCRIPTION` first.

Logical replication does not copy everything. Typical gaps: sequences (check your 16/17 version notes), large objects, DDL (you apply schema changes yourself), and some types. Always read the current restrictions page.

Do not use logical replication as the only backup. Do not subscribe a table to itself. Do not forget DDL: a new column on the publisher can break the subscriber until you add it.

Official: [https://www.postgresql.org/docs/17/logical-replication.html](https://www.postgresql.org/docs/17/logical-replication.html).

### Questions

#### Theoretical questions

1. What object do you create on the publisher?
2. What object do you create on the subscriber?
3. What `wal_level` does the publisher need?
4. Who applies `CREATE TABLE` on the subscriber?
5. What disk risk does a subscription slot create?

#### Easy practical tasks

1. Read `\h CREATE PUBLICATION` and `\h CREATE SUBSCRIPTION` in `psql`.
2. Write a publication for one table and a subscription stub (connection string can be fake in the file).
3. `SHOW wal_level;`. Write whether you can publish now.
4. List restrictions from the 17 logical-replication intro (three items).

#### Medium practical tasks

1. On two lab databases (same cluster is possible with care, two clusters is clearer), publish `shop.items` and subscribe. Insert on the publisher. Select on the subscriber.
2. Add a column on the publisher. Write what you must do on the subscriber.
3. Query `pg_publication`, `pg_subscription`, and `pg_replication_slots`.

#### Advanced practical tasks

1. Use a row filter or a column list on a PostgreSQL 16 or 17 publication. Prove that excluded rows or columns do not arrive.
2. Read 16 versus 17 release notes for logical replication. Write one feature that 17 adds for failover or setup.

---

## Failover tools (high-level: Patroni, Cloud vendor HA)

Streaming replication does not promote a replica by itself when the primary dies. A person or a tool must promote (`pg_ctl promote` or `SELECT pg_promote();`) and must move clients to the new primary.

**Patroni** is a common open-source agent. It uses etcd, Consul, or ZooKeeper to elect a leader. It manages `postgres` and a virtual IP or a proxy. You still need backups. Patroni is not a substitute for PITR.

**Cloud vendor HA** (Amazon RDS Multi-AZ, Google Cloud SQL HA, Azure, AlloyDB, Aurora PostgreSQL, Crunchy-managed) runs a vendor failover path. You learn the vendor runbook: who promotes, what the DNS name does, and what RPO/RTO the product states. Topic 22 lists those products. Topic 21 asks for a written runbook.

What a beginner must remember:

- failover is a process, not a dump file
- clients need one stable name (proxy, DNS, or vendor endpoint)
- a replica that you promote must not share the old primary without fencing
- test failover on a schedule
- keep `pg_dump` or base backups even when HA is on

Do not write a custom elector if your team can use Patroni or a vendor HA product. Do not promote two primaries. Do not skip backups because "we have a replica".

### Questions

#### Theoretical questions

1. Does streaming replication promote a replica automatically?
2. What SQL function promotes a stand-by?
3. What extra system does Patroni use for elections (name one)?
4. Why do you still need backups when HA is on?
5. What client problem must failover solve besides starting PostgreSQL?

#### Easy practical tasks

1. Open the Patroni project page or docs. Write one sentence on what it manages.
2. Open one cloud PostgreSQL HA page (RDS, Cloud SQL, or similar). Write the product name for failover.
3. Write a four-step human failover: detect, fence old primary, promote, move clients.
4. Find `pg_promote` in the 17 docs. Write the return type or the page title.

#### Medium practical tasks

1. In a lab with a replica, run `pg_promote()` on the replica (only if you accept that the old primary must not keep writes). Document how you fence the old primary.
2. Compare RPO of asynchronous streaming versus a vendor synchronous HA page. Write two sentences.
3. Draw: primary, replica, proxy or DNS, backup store. Label who fails over.

#### Advanced practical tasks

1. Write a one-page comparison: manual promote, Patroni, one cloud HA. Columns: elector, client move, backup still required.
2. Read split-brain in the Patroni docs or a vendor HA FAQ. Write how fencing prevents two primaries. Stay at a high level.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. When do you pick `pg_dump`, `pg_basebackup`, PITR, streaming, or logical replication for one problem?
2. How do WAL archives, slots, and forgotten replicas relate to disk-full risk?
3. Why is a hot stand-by not the same object as a logical subscriber?
4. How does PostgreSQL 17 incremental backup change a 16-era backup policy?
5. A teammate says a replica is enough and dumps are optional. Which facts from this topic do you use in the reply?

#### Easy practical tasks

1. Dump a lab database with `pg_dump -Fc` and restore it to a new name. Run `SELECT count(*)` on one table in both databases.
2. Write a cheat sheet: dump formats, base backup, archive, streaming, hot stand-by, publication, Patroni.
3. `SHOW wal_level;` `SHOW archive_mode;` `SELECT pg_is_in_recovery();`.
4. Bookmark docs for `pg_dump`, `pg_basebackup`, and logical replication for your major version.

#### Medium practical tasks

1. Draw a pipeline: nightly `pg_dump`, weekly `pg_basebackup`, continuous archive, one streaming replica. Write what each copy is for.
2. Create a publication and subscription on two lab databases for one table. Then take a `pg_dump` of the subscriber. Write why both exist.
3. Write a restore test plan with a target time and a success check (row present / row absent).

#### Advanced practical tasks

1. Build Docker Compose with a primary and a replica. Fail the primary in the lab. Promote. Point `psql` at the new primary. Keep a dump from before the test.
2. Map this topic to official chapters: backup, continuous archiving, high availability, logical replication. Add one 16-or-17 restriction that this topic did not include.
