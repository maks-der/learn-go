# 12. Backup and Replication

## Description

This topic shows how you copy a database, how you bring it back, and how a DBMS copies changes to another server. You learn logical dumps, physical backups, point-in-time recovery, RPO, RTO, why replication is not a backup, primary and replica roles, failover, and how you test a restore.

Use one term for each concept. A backup is a copy that you can restore without the original server. A primary is the server that accepts writes. A replica applies those writes. Complete this topic after you understand transactions. Replication is not a backup.

This path is vendor-neutral. Command names differ. The jobs do not. This handbook uses primary and replica.

---

## Logical dump vs physical backup

A logical dump copies database objects as SQL (or as another logical format). The tool reads tables and writes `CREATE` and `INSERT` (or `COPY`) statements. You restore a dump into a running DBMS. The restore replays those statements.

A physical backup copies the files that the DBMS uses: data files, and usually the WAL needed to make those files consistent. You restore files, then start the DBMS. The server performs recovery.

| Kind | What you copy | Restore method | Typical use |
| --- | --- | --- | --- |
| Logical dump | Objects and rows as statements | Run SQL against a server | Move data, small databases, version change |
| Physical backup | Data files plus needed log | Replace files, start, recover | Large databases, exact page image, PITR base |

Logical dump benefits: you can restore one table. You can often load into a newer major version. The dump is readable.

Logical dump costs: dump and restore take time proportional to row volume and index rebuild. A dump of a busy database needs a consistent snapshot (transaction or export mode). The dump is not a page-level clone.

Physical backup benefits: faster for large volumes. The copy matches the storage layout. Combined with WAL archive, you get point-in-time recovery.

Physical backup costs: the copy is tied to the product and often to the major version. You cannot restore one table with a simple file copy. You must follow the product procedure so that files and WAL match.

```text
Logical :  DBMS --read rows--> dump file --SQL--> new DBMS
Physical:  data files + WAL  -->  same product, then recover
```

Do not call a file-system snapshot a backup unless you follow the product rules for a consistent snapshot (freeze, snapshot, thaw, or an equivalent). A raw copy of open data files without WAL is not a reliable restore source.

Do not store the only backup on the same disk as the live data. A disk failure then removes both.

### Questions

#### Theoretical questions

1. What does a logical dump contain?
2. What does a physical backup copy?
3. When is a logical dump a better fit than a physical backup?
4. Why is a raw copy of open data files a weak backup?
5. Why must the only backup not live on the same disk as the data?

#### Easy practical tasks

1. Find the official dump command and the official physical backup method for your DBMS.
2. Make a two-column table: logical versus physical. Add four differences from this section.
3. Write one sentence that distinguishes a dump file from a data-directory copy.
4. Name two restore targets: empty new instance, and replace-in-place. Mark which backup kind fits each.

#### Medium practical tasks

1. Dump a practice database. Restore it under a new database name. Run a row count. Compare with the source.
2. Read the product page on consistent dump of a live system. Write the option that gives a consistent snapshot.
3. Estimate dump time versus file-copy time for a 1 GB database in words (order of magnitude is enough).

#### Advanced practical tasks

1. Write a one-page choice record: dump, physical backup, or both, for a class project and for a 500 GB production system.
2. Follow the official physical backup procedure on a disposable instance. Restore to a second data directory. Confirm that the server starts.

---

## Point-in-time recovery, RPO, and RTO

Point-in-time recovery (PITR) restores the database to a chosen instant. You start from a base physical backup. You replay archived WAL (or the equivalent log) up to that instant.

PITR needs three parts:

1. A base backup taken at time T0.
2. A continuous archive of WAL from T0 to the target time T1.
3. A recover command that stops at T1 (timestamp, transaction id, or named restore point).

```text
Base backup @ T0     WAL archive ------> T1 (target)
                 replay, then stop
```

A logical dump is usually a single point. You cannot replay WAL onto a dump in the usual way. For PITR you use a physical base plus logs.

If the WAL archive has a gap, you cannot replay through the gap. You can recover only to the last continuous point. Treat archive success as a production check.

PITR does not undo an application bug by itself. You choose T1 before the bad write. You lose committed work after T1. That loss is the RPO for this restore.

RPO is recovery point objective. RPO is the maximum amount of data loss that the business accepts, measured in time. An RPO of 1 hour means that a restore may lose up to 1 hour of committed work.

RTO is recovery time objective. RTO is the maximum time to restore service after a failure. An RTO of 4 hours means that users may wait up to 4 hours.

```text
Failure ---- restore work ---- service back
         |<------ RTO ------>|

Last good restore point ---- failure
|<----------- RPO ---------->|  (lost work)
```

Backup frequency and WAL archive set a lower bound on RPO. A nightly dump without log archive cannot beat "since last night" as RPO. Continuous WAL archive can make RPO seconds, if the archive is off-host and you test it.

Restore method, database size, and automation set a lower bound on RTO. A 2 TB physical restore plus replay can exceed a short RTO. A replica promote can meet a short RTO. That promote is not a backup.

Write RPO and RTO as numbers with units. "We need good backups" is not an objective. "RPO 15 minutes, RTO 2 hours" is an objective.

Different systems can have different objectives. A session cache can accept a large RPO. A payments ledger cannot.

Do not set RPO to zero unless you paid for the design: synchronous copies, tested failover, and a written exception for dual-site loss.

Do not set RTO shorter than the measured restore time. Measure a restore. Then write the objective, or change the method.

Do not keep base backups without WAL archive if you advertised PITR. The base alone is a restore to T0 only.

### Questions

#### Theoretical questions

1. What three parts does PITR need?
2. Why is a logical dump not the usual base for WAL replay?
3. What does RPO measure, and what does RTO measure?
4. Why can a nightly dump not meet a 15-minute RPO?
5. What committed work do you lose when you stop at T1, and what can you restore if you have a base backup and no WAL archive?

#### Easy practical tasks

1. Write the T0 / T1 timeline from this section in four sentences. Draw base backup, WAL files, and a stop time.
2. Write RPO and RTO in one sentence each with a numeric example. Assign rough RPO/RTO to: personal notes, shop orders, a public blog cache.
3. Find the product name for WAL archive and for restore-to-time. List three events that you might choose as T1 (bad deploy, bad `DELETE`, drop of a table).
4. List three factors that increase RTO for a physical restore.

#### Medium practical tasks

1. Read the official PITR tutorial for your DBMS. Write the stop-target types that the product allows. Write a checklist: how you verify that WAL archive is continuous for one day.
2. Time a dump-and-restore of your practice database. Write that time as a lower bound on RTO for that method.
3. Map backup methods to best-case RPO: nightly dump, hourly dump, WAL archive every minute. Explain in five sentences why PITR to "one minute before the bad delete" still loses later good writes.

#### Advanced practical tasks

1. On a disposable instance, take a base backup, archive WAL, run a write, restore to a time before that write. Record the steps and the proof query.
2. Write a one-page service sheet: two databases, RPO, RTO, backup method, and a gap you cannot meet yet. Design a backup schedule that meets RPO 30 minutes and RTO 3 hours for a 50 GB database. State assumptions.

---

## Replication is not a backup

Replication copies changes to another running DBMS (a replica). The replica is online. You use it for reads or for failover.

A backup is a copy that you can restore after you lose the original data. The backup can be offline. The backup can be a week old. The backup is not a live process.

| Need | Replication | Backup |
| --- | --- | --- |
| Fast failover | Yes, if designed | Slow: restore first |
| Extra read capacity | Yes, with lag | No |
| Recover from bad `DELETE` | No, if already replicated | Yes, if you have an older copy or PITR |
| Recover from ransomware on the server | No, if the replica was also changed | Yes, if the backup is isolated |
| Meet RPO after site loss | Only if the replica is in another site and durable | Yes, if the backup is off-site |

A replica applies the same writes. A bad transaction that deletes a table also deletes the table on an asynchronous replica a moment later. Replication does not give you yesterday.

A backup does not serve live queries. A backup does not replace a replica for high availability.

Many production systems need both: replicas for availability and reads, backups for history and isolation.

Do not point the only backup at the same replica disk and call the system safe. Isolate backups. Use retention (how many days you keep copies).

Do not stop backups because "we have a replica." Write the failure that each tool covers.

### Questions

#### Theoretical questions

1. What is the primary job of replication in this section?
2. What is the primary job of a backup?
3. Why does a replica not protect you from a committed bad `DELETE`?
4. Which tool helps after ransomware on the live servers, if the copy is isolated?
5. Why do many systems run both replicas and backups?

#### Easy practical tasks

1. Fill the comparison table from this section in your notes with one extra row of your own.
2. Label four events as "replica helps", "backup helps", or both: disk loss, bad deploy SQL, read scaling, off-site history.
3. Write one sentence that you would say to a teammate who wants to drop backups.
4. Define retention in one sentence. Pick a retention of 7 days or 30 days for a class project. Give one reason.

#### Medium practical tasks

1. Write a failure matrix with five failures and a yes/no for replica and for backup.
2. Read a vendor high-availability page. Mark each feature as replica, backup, or both.
3. Describe how a delayed replica (if the product has one) is still not a full backup. Write four sentences.

#### Advanced practical tasks

1. Write a one-page architecture note: primary, replica, WAL archive, off-site backup. Label RPO/RTO for two failure types.
2. Design an isolated backup path (different account, different host). List three isolation rules.

---

## Primary and replica; sync vs async

A primary is the DBMS instance that accepts writes. Applications send `INSERT`, `UPDATE`, and `DELETE` to the primary. The primary writes WAL (or an equivalent stream).

A replica is an instance that receives that stream and applies the changes. A typical replica is read-only. You can send `SELECT` to the replica. A write on a read-only replica fails.

```text
Application writes  -->  primary  -->  log stream  -->  replica
Application reads   -->  primary or replica
```

Roles are not hardware labels. You can promote a replica to primary. After promote, the old primary must not keep the write role.

A replica starts from a base copy of the primary data. Then it follows the log. If the replica stops, it must catch up from the stored log position. If the gap is too large, you rebuild the replica from a new base backup.

Use cases for a replica: extra read capacity, a warm standby for failover, a reporting instance so that heavy `SELECT` does not hit the primary, a source for backups (with product rules).

Do not write to a replica unless the product documents a special session. Do not put the only replica on the same host as the primary. A host failure then stops both.

Asynchronous replication means the primary commits locally, then sends the change to the replica. `COMMIT` on the primary does not wait for the replica. The replica can lag. If the primary fails before the replica received the last records, those commits can be lost on failover.

Synchronous replication means `COMMIT` waits until a replica (or a quorum) has the change at a defined level: received, flushed to disk, or applied. The exact level is a product setting.

| Mode | Commit wait | Typical lag | Failover data risk |
| --- | --- | --- | --- |
| Asynchronous | Primary only | Milliseconds to seconds or more | Last commits can be missing on the replica |
| Synchronous | Primary plus replica ack | Low if the replica is healthy | Lower loss; commit slower if the replica is slow |

```text
Async:  client <-- COMMIT -- primary      replica (later)
Sync:   client <-- COMMIT -- primary <-- ack -- replica
```

Synchronous replication reduces RPO toward zero for the protected commits. It increases commit latency. If the synchronous replica stops, writes can stall unless you have a fallback policy.

Asynchronous replication is the usual default. It fits read replicas and distant sites. It does not give "zero data loss" on primary crash.

Replication lag is the delay between a commit on the primary and the moment the replica can see that commit. Applications that write and then immediately read their own write must read the primary (or wait until the replica catches up).

Do not mix the words. "We replicate" does not say whether commit waits. Name the mode and the ack level.

A common pattern: one synchronous replica in the same site for failover, asynchronous replicas in other sites for reads and disaster recovery. That pattern is a design, not a default.

### Questions

#### Theoretical questions

1. Which instance accepts writes in a primary/replica pair, and what does a typical replica reject?
2. When does an asynchronous primary return `COMMIT` to the client, and what can you lose if the primary fails in asynchronous mode?
3. What does a synchronous commit wait for, and what happens to writes if the only synchronous replica stops?
4. What is replication lag, and why can a write-then-read on a replica miss the new row?
5. Why is "we replicate" an incomplete durability statement?

#### Easy practical tasks

1. Draw primary, log stream, replica, and two arrows for write and read.
2. Fill a two-row table: async versus sync. Add commit wait and failover risk.
3. Write one `SELECT` you would send to a replica and one `INSERT` you would send only to the primary. Write one use case for async and one use case for sync.
4. Find the setting name for synchronous replication in your DBMS. Find the product terms: primary, standby, publisher, or another pair. Write the mapping.

#### Medium practical tasks

1. If you can, start a primary and a replica. Run `SELECT` on both. Run `INSERT` on the primary. Read the row on the replica. Insert a row on the primary. Select on the replica in a loop until it appears. Write the delay.
2. Read the ack levels (received, flushed, applied). Write one sentence each. Write a fallback policy in five sentences for a dead synchronous replica.
3. Design a routing rule: session after POST goes to primary for N seconds. Write the rule in six sentences. Compare physical and logical replication in your manual. Write three differences.

#### Advanced practical tasks

1. Write a one-page role note: who writes, who reads, how you promote, and what you do with the old primary. Include a read-routing policy and the lag alert.
2. Write a one-page HA design: same-site sync replica, off-site async replica. State RPO for two failure types. On a disposable pair, compare commit time with async and with sync if the product allows both.

---

## Failover and split-brain (high-level)

Failover is the change of role when the primary stops. A replica becomes the new primary. Clients must connect to the new primary.

Failover can be planned (switchover): you stop writes, wait for the replica to catch up, promote, then point clients. Failover can be unplanned: the primary crashes or the network fails. A manager process or an operator promotes a replica.

```text
Primary down --> choose replica --> promote --> redirect clients --> fence old primary
```

Fencing means the old primary must not accept writes after promote. If the old primary is alive but isolated, it can still accept writes. Two primaries that accept writes is split-brain.

Split-brain produces two histories. Merging them is hard. Unique keys collide. Counters diverge. The usual repair is to pick one history and discard or manually merge the other.

Prevention:

1. Use a cluster manager that gets a lock or a quorum.
2. Fence the old primary (STONITH, revoked disks, or forced read-only).
3. Do not let applications keep a cached write address without a retry policy that discovers the new primary.

Do not promote two replicas. Do not start the old primary as a primary again without a rewind or a rebuild.

DNS or a virtual IP can hide the new host. Clients still need a timeout and a reconnect. Long-lived connections to the dead primary fail. Pools must drop those connections.

Failover time is part of RTO. Test failover as you test restore. An untested promote is a hope.

### Questions

#### Theoretical questions

1. What role change does failover perform?
2. What is the difference between planned switchover and unplanned failover?
3. What is split-brain?
4. Why must you fence the old primary?
5. Why is an untested promote a hope?

#### Easy practical tasks

1. Number the failover steps from this section.
2. Write one sentence that defines fencing.
3. List three ways clients can find the new primary (VIP, DNS, connection string).
4. Write one data symptom of two primaries (example: two rows with the same natural key in different histories).

#### Medium practical tasks

1. Read the official failover or promote command. Write the command and the fence warning from the docs.
2. Write a client reconnect policy: timeout, retry, drop pool connections.
3. Describe a network partition: primary in site A, replica and manager in site B. Write who should win and why.

#### Advanced practical tasks

1. On a disposable pair, promote a replica. Try to write to the old primary. Record whether the product fences it. Write the gap if it does not.
2. Write a one-page incident runbook: detect down primary, promote, fence, check lag first, tell applications.

---

## Test a restore

A backup that you never restore is a hope, not a procedure. Test the restore. Write the time. Write the proof query.

A restore test does this:

1. Take a backup with the official tool.
2. Restore onto a different instance or a different data directory.
3. Start the DBMS.
4. Run checks: the database exists, a row count, a checksum or a known value, application login.
5. Record duration. Compare the duration with RTO.
6. Record the newest transaction time. Compare it with RPO.

```text
Backup --> isolated host --> restore --> checks --> notes --> destroy the test copy
```

Test on a schedule. A monthly restore test finds broken credentials, missing WAL, and version mismatches. A test only on the first install is not enough.

Test the path that you would use in an incident. If production restore is physical plus PITR, do not test only a logical dump of a toy schema.

Do not restore over the live data directory to "test." Restore to a copy. Do not expose the test copy on a public network. The test copy contains the same personal data as production.

Automate the checks. A human "it started" is a weak check. Count rows in a critical table. Compare with a number that you stored at backup time.

If the test fails, the incident is now, not during the outage. Fix the backup. Repeat the test.

### Questions

#### Theoretical questions

1. Why is an untested backup a hope?
2. What six steps does this section name for a restore test?
3. Why must you restore to a different instance?
4. Why must a test copy follow the same data-protection rules as production?
5. What do you do when a scheduled restore test fails?

#### Easy practical tasks

1. Write a proof query for your practice schema (example: `COUNT(*)` of `orders`).
2. Number the six steps. Put them in a checklist file.
3. Write two checks that are stronger than "the process started."
4. Write one reason to run a restore test every month.

#### Medium practical tasks

1. Dump or back up your `learn` database. Restore to a new name or instance. Run your proof query. Write the duration.
2. Store the proof count at backup time. After restore, compare the counts.
3. Fail a restore on purpose (wrong file path). Record the error. Fix the path. Restore again.

#### Advanced practical tasks

1. Write a restore-test runbook with roles, isolated host, checks, and how you destroy the copy. Run it once on a disposable system.
2. Time logical restore and, if you can, physical restore. Write which method you would use for an RTO that is half of the slower time.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do logical dump, physical backup, and PITR map to different RPO and RTO numbers?
2. When is a replica the wrong tool and a tested backup the right tool, in one incident story?
3. How do primary role, sync mode, and lag determine what a promote will lose?
4. How do WAL archive gaps break both PITR and a promised RPO?
5. A teammate says the nightly dump and a local replica are enough for payments with RPO 5 minutes. Which facts do you use in the reply?

#### Easy practical tasks

1. Write a cheat sheet: dump versus files, PITR parts, RPO, RTO, replica versus backup, sync versus async, failover, restore-test steps.
2. Take one backup of `learn`. Restore it once. Write the command names only.
3. Write numeric RPO and RTO for your class project. Name the backup method.
4. Label your learning setup as single node or primary/replica. Write the write target.

#### Medium practical tasks

1. Build a mini runbook: backup command, restore command, proof query, measured time, RPO/RTO targets.
2. Write an HA diagram for a class shop: primary, one replica, backup off-host. Mark write path, read path, and backup path. Write a promote checklist of eight items that includes lag check and fencing.
3. Write a matrix of four failures versus dump, PITR, and replica. Choose sync or async for (a) same rack standby and (b) another continent report replica. Justify each.

#### Advanced practical tasks

1. Run a scheduled-style test: backup, wait, more writes, restore to a copy, show that extra writes are absent or present as you intend.
2. Write an operations review of a public backup and HA doc (your DBMS). Mark RPO, RTO, isolation, fencing, and test frequency. Fill any gap with a proposed control.
