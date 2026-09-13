# 16. Administration

## Description

This topic shows how you run a MongoDB server. You configure `mongod`. You learn the **WiredTiger** storage engine. You learn **journaling**. You take **backups** and you practice **restore**. You learn when **compact** helps. You plan **major version upgrades**. Complete replication and security basics in parallel. Do not run administration commands on a shared cluster without a role and a change window.

Use one term for each concept. **`mongod`** is the database process. **WiredTiger** is the default storage engine. The **journal** is the write-ahead log that WiredTiger uses for crash recovery. A **restore drill** is a practice restore that you time and that you verify.

Administration is operations work. The application developer still must know backup and upgrade risk.

---

## `mongod` configuration

`mongod` reads a YAML configuration file. The default path depends on the operating system. You can also pass flags. Prefer a file for anything that must survive a restart.

Common settings (names are examples; read your version):

- `storage.dbPath` — data files
- `systemLog.path` — log file
- `net.port` — listen port (`27017` is the default)
- `net.bindIp` — interfaces
- `replication.replSetName` — replica-set name
- `security.authorization` — user access on or off
- `processManagement.fork` — daemon mode on Unix (when you use it)

Example shape:

```yaml
storage:
  dbPath: /var/lib/mongo
systemLog:
  destination: file
  path: /var/log/mongodb/mongod.log
net:
  port: 27017
  bindIp: 127.0.0.1
```

Start:

```text
mongod --config /etc/mongod.conf
```

`mongosh` is the client. `mongod` is the server. Do not confuse the two processes.

On Windows, MongoDB often runs as a service. On Docker, flags or an env file set the same ideas. Atlas hides `mongod.conf`. You set equivalent ideas in the Atlas UI (cluster tier, IP list, backup).

Validate a change on a test host. A wrong `dbPath` or `bindIp` can start an empty instance or expose the port.

Do not enable a public `bindIp` without authentication and TLS.

Store the config in version control without secrets. Put keys and passwords in a secret store.

### Questions

#### Theoretical questions

1. What process is `mongod`?
2. What format is the usual configuration file?
3. What does `storage.dbPath` select?
4. Why can a wrong `dbPath` look like an empty database?
5. Who manages `mongod.conf` on Atlas?

#### Easy practical tasks

1. Find the config file or the Docker command that starts your `mongod`. Write the path or the command.
2. Open the configuration-file page. Write the URL.
3. Write five settings from this section and one sentence each.
4. Make a table: Community Server vs Atlas. Add who sets port, disk, and replica-set name.

#### Medium practical tasks

1. Change the log path on a lab instance. Restart. Confirm new log lines.
2. Compare `mongod --help` with the YAML reference. Write two flags that map to YAML keys.
3. Write a minimal config for a one-member replica set on localhost. Do not use it on a public host.

#### Advanced practical tasks

1. Read configuration validation and include files if your version has them. Write how you test a file before restart.
2. Document a production config checklist: bind IP, auth, TLS, replica set, paths, log rotation.

---

## Storage engine: WiredTiger

**WiredTiger** is the default storage engine. It stores documents and indexes on disk. It uses a cache in RAM. It supports document-level concurrency. Two updates to two documents can run at the same time.

WiredTiger compresses data. Compression reduces disk. Compression uses CPU. The default is fine for most learning systems.

Each collection and each index can have WiredTiger options. You rarely change them at the start.

WiredTiger is not MMAPv1. MMAPv1 is gone from current MongoDB. Old blog posts about record padding and collection-level locks do not apply to WiredTiger.

`db.serverStatus().wiredTiger` shows cache and other counters. Atlas shows storage charts.

Checkpoints write a consistent snapshot to the data files. Between checkpoints, the journal holds recent writes. The next section covers the journal.

Do not delete files under `dbPath` to "free space" while `mongod` runs.

Do not mix data files from two engines or two major versions without the official upgrade procedure.

If you see "storage engine" in old training, confirm the page is for WiredTiger and for your version.

### Questions

#### Theoretical questions

1. What is the default storage engine?
2. What does document-level concurrency mean?
3. Why do old MMAPv1 lock articles mislead you now?
4. What is a checkpoint at a high level?
5. Why must you not delete files in `dbPath` while the process runs?

#### Easy practical tasks

1. Run `db.serverStatus().storageEngine` or the equivalent. Write the name.
2. Open the WiredTiger page. Write the URL.
3. Write four sentences: cache, compression, document locks, checkpoints.
4. Make a table: MMAPv1 (historical), WiredTiger. Add two rows.

#### Medium practical tasks

1. Read `wiredTiger.cache` in `serverStatus`. Write current and max bytes.
2. Create two collections. Compare `collStats` storage sizes. Write the numbers.
3. Find the compression default in the manual. Write the name.

#### Advanced practical tasks

1. Read collection-level WiredTiger options (`block_compressor`). Write when a team might change it.
2. Explain how WiredTiger cache and the OS file cache both hold data. Write six sentences.

---

## Journaling

**Journaling** writes intended changes to a journal before (or as part of) making them durable. After a crash, WiredTiger replays the journal and returns to a consistent state.

Journaling is on by default for WiredTiger. Do not turn it off on a system that you care about.

Write concern `j: true` waits until the write is in the journal (and follows the rules for that concern). `w: "majority"` interacts with journal and replica acknowledgment. Read the write-concern page again with journaling in mind.

The journal lives under `dbPath`. It is not a user backup. You cannot ship only the journal and ignore data files.

A power loss with journaling on does not corrupt WiredTiger metadata if the disk honors flush. Cheap USB disks and some virtual disks lie about flush. Use storage that is safe for databases.

Atlas manages journaling. You still choose write concern.

Do not disable the journal to "go faster" on production.

Do not copy data files for backup while writes run unless you use a supported snapshot method (next sections).

### Questions

#### Theoretical questions

1. What problem does the journal solve after a crash?
2. Is journaling on by default for WiredTiger?
3. How does `j: true` relate to the journal?
4. Why is the journal not a backup?
5. Why does a disk that ignores flush break crash recovery?

#### Easy practical tasks

1. Open the journaling page. Write the URL and whether you can disable the journal.
2. Write four sentences: journal, checkpoint, crash, write concern `j`.
3. Find `dbPath` on your lab. Write the data directory name. Do not delete files.
4. Make a table: durable write, what you wait for. Add `w:1`, `j:true`, `w:majority`.

#### Medium practical tasks

1. Insert with `{ writeConcern: { w: 1, j: true } }` and with `{ w: "majority" }`. Write if both succeed on your set.
2. Read how WiredTiger uses the journal with checkpoints. Draw a three-step timeline.
3. Compare latency of 1000 inserts with `j: true` and without it on a lab. Write the two times.

#### Advanced practical tasks

1. Read recovery after an unclean shutdown. Write the operator steps at a high level.
2. Write a storage policy: disk type, flush, journal always on, who may change `j`.

---

## Backup: `mongodump`, snapshots, Atlas backup

A **backup** is a copy that you can restore. A secondary is not a backup. A delayed member is not a full backup.

Common methods:

**`mongodump` / `mongorestore`.** Logical export of BSON. Good for small data and for selected collections. Slow and heavy on large clusters. You can dump from a hidden secondary.

**Filesystem or volume snapshots.** Snapshot the disk (or the cloud volume) with a consistent method. For self-managed replica sets, the manual describes snapshot plus the journal, or snapshot on a secondary with flush procedures. Cloud disks (EBS, Azure disks, GCP disks) have snapshot products.

**Atlas backup.** Atlas offers cloud snapshots and, on some tiers, more frequent backups and point-in-time restore. Use Atlas backup for Atlas clusters unless you have a written exception.

Pick a method that matches size and RPO (how much data you can lose). A nightly `mongodump` of 2 TB is often the wrong tool.

Encrypt backups. Control who can download them. A backup has all the data.

Test the backup. A backup that you never restore is a hope, not a plan. The next section is the drill.

Do not run `mongodump` against the primary at peak without a plan. Prefer a hidden secondary.

Do not copy `dbPath` with `cp` or zip while `mongod` writes, unless the official procedure allows that exact method.

### Questions

#### Theoretical questions

1. Why is a secondary not a backup?
2. What does `mongodump` write?
3. When is a disk snapshot better than `mongodump`?
4. What does Atlas backup replace for Atlas users?
5. Why must you encrypt backup files?

#### Easy practical tasks

1. Open the backup methods page and the Atlas backup page. Write both URLs.
2. Run `mongodump --help`. Write two options (`--db`, `--out`).
3. Make a table: method, good for, poor for. Add dump, snapshot, Atlas.
4. Write your RPO in one sentence for a learning project (for example "one day").

#### Medium practical tasks

1. Dump one small database with `mongodump`. Restore it to a new database name with `mongorestore`. Confirm a document.
2. If you use Atlas, open Backup. Write the snapshot schedule that you see or "none on this tier".
3. Write who may run dump in your team and from which member (hidden secondary or Atlas).

#### Advanced practical tasks

1. Read consistent snapshot steps for a replica set on your OS or cloud. Write the steps in order. Do not skip journal notes.
2. Estimate dump time for 100 GB from a blog or from a lab scaling test. Propose snapshot instead if the time is too long.

---

## Restore drills

A **restore drill** is a practice restore. You restore to a safe target. You measure time. You check data. You write what failed.

A drill answers:

- How long to restore (RTO)
- Who has the passwords and the roles
- Whether the backup is complete
- Whether applications can point to the restored data

Steps in spirit:

1. Pick a backup (dump, snapshot, or Atlas snapshot).
2. Restore to a new cluster or a new database name. Do not overwrite production on the first drill.
3. Run checks: document counts, a sample find, an application smoke test.
4. Record the clock and the problems.
5. Repeat on a schedule (for example every 90 days).

Atlas: use restore to a new cluster when you learn. Production point-in-time restore is a named procedure. Practice it in a non-production project.

`mongorestore` can drop collections if you pass drop flags. Read the flags. A wrong restore can delete data.

After a restore, replica-set identity and user credentials can surprise you. Follow the manual for restore into a replica set.

Do not wait for the first real incident to learn the UI.

Do not store the only restore instructions in one person's head.

### Questions

#### Theoretical questions

1. What is a restore drill?
2. What is RTO in this context?
3. Why do you restore to a new name first?
4. Why can `mongorestore` be dangerous?
5. Why is a 90-day drill useful?

#### Easy practical tasks

1. Write a drill checklist with eight steps.
2. Open the `mongorestore` page. Write the URL and the drop-related flag name.
3. Write who you would call if production data is gone (role names, not personal data if you share notes).
4. Make a table: check after restore, command or action. Add four rows.

#### Medium practical tasks

1. Restore yesterday's lab dump to `learn_restore`. Compare `countDocuments` with the source.
2. Time the restore. Write RTO for that size.
3. Write a one-page drill report template: date, backup id, duration, errors, sign-off.

#### Advanced practical tasks

1. On Atlas (or on a local replica set), practice restore to a new cluster or new port. Point `mongosh` at it. Write the connection change.
2. Design a yearly drill that includes a dropped collection and a point-in-time target. Write the success criteria.

---

## Compact / compacting awareness

WiredTiger can leave free space inside files after deletes. The disk file does not always shrink. **Compact** is an operation that rewrites data to release space to the file system (with limits and with version-specific behavior).

```javascript
db.runCommand({ compact: "orders" })
```

Compact can take a long time. Compact uses extra disk while it runs. Compact can affect latency. Current versions improved compact, but it is still an administration event, not a weekly habit.

Awareness:

- Deletes do not always return disk to the OS
- Drop a collection if you can replace it (that returns space)
- TTL deletes still need compact or a drop-and-recreate in some cases if files stay large
- Atlas disk charts can show used vs allocated

Do not run compact on the primary at peak. Prefer a maintenance window. On a replica set, some teams compact secondaries first. Read the current procedure.

Do not compact to fix a bad model (unbounded arrays, huge documents). Fix the model.

`compact` is not `repairDatabase`. Old repair advice can be harmful. Use official pages for your version.

Atlas may limit or hide compact. You still need the idea when disk does not drop after a large delete.

### Questions

#### Theoretical questions

1. Why can disk stay large after a big delete?
2. What does `compact` try to do?
3. Why is compact not a weekly default?
4. Why is drop of a collection sometimes better?
5. Does compact fix an unbounded-array model?

#### Easy practical tasks

1. Open the `compact` command page. Write the URL and one warning.
2. Write four sentences: delete, file size, compact, drop.
3. Run `collStats` on a collection. Write `size` and `storageSize`.
4. Make a table: action, likely disk effect. Add delete many, drop, compact, insert.

#### Medium practical tasks

1. Insert many documents. Delete them. Compare `storageSize` before and after. Write if the file shrank.
2. Read whether compact is blocking on your version. Write the answer.
3. Write a rule: when your team may run compact (window, member, who approves).

#### Advanced practical tasks

1. Read compact on secondaries and rollback risk. Write a safe order for a three-member set.
2. Compare compact with "create new collection, copy, rename" for a huge delete. Write trade-offs.

---

## Major version upgrades

A **major version** is the first number (for example 6 to 7, or 7 to 8). MongoDB supports an upgrade path. You do not jump over unsupported gaps. Read **Upgrade** for your from-version and to-version.

Typical ideas:

- Upgrade binaries on secondaries first, then the primary (rolling)
- Set **feature compatibility version** (`setFeatureCompatibilityVersion`) only when the docs say so
- Read compatibility of drivers and of Atlas
- Test on a copy of the data
- Read deprecated commands and removed options

Atlas: you pick a version in the UI. Atlas still needs a change window and an application test. Read the Atlas upgrade notes.

Community Server: download the new packages. Follow the replica-set rolling upgrade page. Do not upgrade only the primary first.

Drivers: an old driver can fail on a new server. Read the compatibility matrix.

Feature compatibility version (FCV) controls which new on-disk features are on. A downgrade can be hard or impossible after you enable new FCV features. Do not set FCV in production on the same day as the binary upgrade unless the procedure says that you must.

Keep notes: previous version, new version, FCV, date, who, link to the procedure.

Do not upgrade production on a Friday without a team and a rollback plan that the manual still allows.

### Questions

#### Theoretical questions

1. What is a major version in MongoDB?
2. Why do you upgrade secondaries before the primary?
3. What is feature compatibility version?
4. Why can a driver fail after a server upgrade?
5. Why can downgrade fail after FCV change?

#### Easy practical tasks

1. Write the major version of your server (`db.version()`).
2. Open the upgrade page for the next major version. Write the URL.
3. Open the driver compatibility matrix. Write if your driver major is listed.
4. Make a table: step, why. Add backup, upgrade secondary, upgrade primary, FCV.

#### Medium practical tasks

1. Read removed features between your version and the next. Write two items that could break an app.
2. Write a test plan: CRUD, index build, transaction, change stream after upgrade.
3. If you use Atlas, write the UI path to schedule a version change. Do not apply it on a shared project.

#### Advanced practical tasks

1. Perform a major upgrade on a local replica set in Docker. Record commands, FCV, and one failed write if any.
2. Write a team upgrade standard: freeze window, canary app, rollback limits, who sets FCV.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do `mongod` config, WiredTiger, and the journal work together on one acknowledged write?
2. When do you choose `mongodump`, a snapshot, or Atlas backup for the same data size?
3. Why do restore drills and compact both need a maintenance window but solve different problems?
4. How does a major upgrade change risk if backups and FCV are not in the plan?
5. A teammate disables the journal, copies `dbPath` with zip during writes, and runs compact on the primary at noon. Which facts do you use in the review?

#### Easy practical tasks

1. Write a cheat sheet: config keys, WiredTiger, journal, dump vs snapshot vs Atlas, drill, compact, rolling upgrade, FCV.
2. Record your environment: Community, Docker, or Atlas, version, backup method, last drill date or "never".
3. Draw crash recovery: write, journal, checkpoint, restart, replay.
4. List five administration actions that you must not run on production without a ticket.

#### Medium practical tasks

1. Run one dump/restore drill on a lab database. Complete the report template from the drill section.
2. Write an operations runbook page: restart `mongod`, where logs live, how to see replica health, who has Atlas admin.
3. Compare disk `storageSize` after a large delete and after a drop of a test collection. Write the lesson for compact.

#### Advanced practical tasks

1. Produce an administration standard: config baseline, backup RPO/RTO, drill calendar, compact rules, upgrade path, Atlas vs self-managed split.
2. Practice a rolling restart of a three-member lab set. Confirm no data loss and write the write-concern you used during the restart.
