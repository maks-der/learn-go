# 21. Operations and Architecture

## Description

This topic shows how you place a database in a larger system and how you run it day to day. You learn OLTP versus OLAP, warehouses versus operational databases, ETL and ELT, change data capture (CDC), observability, capacity, and schema change in production.

Use one term for each concept. An operational database serves live transactions. A warehouse serves analysis. Observability is the set of signals that tell you if the system is healthy. Complete this topic after you understand backups, replication, and performance measurement. Architecture choices use those operations.

This path is vendor-neutral. Tool names differ. The jobs do not.

---

## OLTP vs OLAP

OLTP means online transaction processing. The workload is many small reads and writes. Each transaction touches a few rows. The user waits. Examples: place an order, check in a guest, post a payment.

OLAP means online analytical processing. The workload is fewer queries that read many rows. The query aggregates. The user is an analyst or a report. Examples: sales by month, conversion funnel, inventory aging.

```text
OLTP:  UPDATE one order, COMMIT, 10 ms
OLAP:  SUM sales for two years, 30 s, scan many pages
```

Schema shape differs. OLTP favors normalized tables and current state. OLAP favors wide fact tables, dimensions, and history that does not change (or changes in a controlled way).

If you run heavy OLAP on the OLTP primary, you pollute the buffer cache, you hold locks or you consume I/O, and you slow checkout. Use a replica, a warehouse, or a night batch.

Some products try to serve both (HTAP). Treat that as a product claim. Measure. Isolate the workloads if checkout latency rises.

Do not design an OLTP table for a 50-column report first. Do not design a star schema for a 5 ms checkout.

Indexes differ. OLTP indexes support point lookups and small ranges. OLAP indexes and column stores support scans and aggregates.

### Questions

#### Theoretical questions

1. What is a typical OLTP transaction?
2. What is a typical OLAP query?
3. Why does a large aggregate harm an OLTP primary?
4. How do schemas differ at a high level?
5. What must you measure if a product claims to serve both?

#### Easy practical tasks

1. Label six statements as OLTP or OLAP: checkout, monthly tax report, password login, cohort analysis, add-to-cart, year-over-year sales.
2. Draw two boxes: OLTP primary and OLAP store. Write one query on each.
3. Write three OLTP indexes and one OLAP-style access (scan plus aggregate).
4. Write one sentence you would use to refuse a 2-minute report on the checkout database.

#### Medium practical tasks

1. Run a large `GROUP BY` on a practice table on the same instance as small updates. Write a qualitative effect (or a time if you can measure).
2. Design a normalized checkout model and a fact table for "sales by day and sku." List columns of the fact table.
3. Decide: replica versus warehouse for a daily 20-minute report. Write four reasons.

#### Advanced practical tasks

1. Write a one-page workload split for a shop: which queries stay on OLTP, which move, and the freshness need.
2. Compare row store versus column store for the fact table from docs. Write five sentences.

---

## Data warehouse vs operational database

An operational database (system of record) holds current business state. It accepts writes from applications. Constraints and transactions protect that state. Backup and RPO apply here first.

A data warehouse holds a copy of data for analysis. It is not the place where a customer completes checkout. Loads arrive in batches or as a stream. Users run reports and models.

```text
Apps --> operational DBMS --> (ETL/ELT/CDC) --> warehouse --> reports
```

Differences:

| Topic | Operational | Warehouse |
| --- | --- | --- |
| Writes | Application transactions | Bulk load or pipeline |
| Schema | Normalized, current | Historical, dimensional or wide |
| Users | Applications, a few admins | Analysts, BI tools |
| Freshness | Now | Minutes to a day (by design) |
| Failure | Outage stops the business | Outage stops reports |

A replica of the operational database is not a warehouse. A replica has the same schema and the same row layout. A warehouse transforms, history-tracks, and serves different tools.

Do not grant analysts `SELECT` on production operational tables as the only architecture. You leak personal data, you load the primary, and you cannot rebuild history.

Do not write operational updates in the warehouse and expect the application to read them back. One direction of truth: operational is the record. The warehouse is derived.

Some teams add a lake (files) in front of or beside the warehouse. The idea stays: derived analytical storage, not checkout.

Access control still applies. A warehouse can hold more history and more joined personal data. Treat it as sensitive.

### Questions

#### Theoretical questions

1. What writes does an operational database accept?
2. What job does a warehouse do?
3. Why is a replica not a warehouse?
4. Why must analysts not use production OLTP as the only report store?
5. Which system is the system of record for an order?

#### Easy practical tasks

1. Fill the comparison table in your notes with one extra row.
2. Draw apps, operational DBMS, pipeline, warehouse, dashboard.
3. Write two users of each system.
4. Write a freshness number for a warehouse that you would accept for a shop (example: 1 hour). Give one reason.

#### Medium practical tasks

1. List five operational tables. Write which become facts and which become dimensions.
2. Write a data-access rule: who can read production, who can read the warehouse, how you mask personal data.
3. Describe a warehouse outage versus an operational outage in four sentences each.

#### Advanced practical tasks

1. Write a one-page platform note: record store, replica, warehouse, who owns each, RPO for each.
2. Design a slowly changing customer dimension (name change). Write how history stays in the warehouse and not in checkout.

---

## ETL / ELT

ETL means extract, transform, load. You extract from the source, you transform in a pipeline, you load into the target (often a warehouse).

ELT means extract, load, transform. You extract, you load raw data into the target, you transform inside the target with SQL or a warehouse engine.

```text
ETL:  source --> transform (job) --> warehouse
ELT:  source --> raw load --> transform (SQL in warehouse)
```

ETL fits when:

- you must reduce data before it leaves the source
- the target cannot express the transform
- you must hide columns before they land

ELT fits when:

- the warehouse can process large SQL
- you want to keep a raw copy and change the transform later
- you accept more storage of raw data

Both need:

- a schedule or a trigger
- idempotent loads or a watermark (load from timestamp T)
- a failure retry
- a data contract (column types and meaning)

Do not transform in ad-hoc notebooks as the only production path. The job must be versioned like application code.

Do not run a full extract every hour if an incremental watermark exists. Full extracts load the operational system.

Quality checks belong in the pipeline: row counts, null rates, unique keys. A silent wrong load is worse than a failed job.

Time zones and money types must survive the pipeline. A float conversion in the middle reintroduces money errors.

### Questions

#### Theoretical questions

1. What is the order of steps in ETL?
2. What is the order of steps in ELT?
3. When do you hide columns in ETL before the load?
4. What is a watermark?
5. Why must a production transform not live only in a notebook?

#### Easy practical tasks

1. Draw ETL and ELT side by side.
2. Write a watermark column for `orders` (`updated_at` or a log sequence).
3. List three quality checks for a nightly load.
4. Write one transform that must happen before the data leaves the source (example: drop a secret column).

#### Medium practical tasks

1. Write a pseudo-job: extract new orders since last watermark, load, update watermark, on failure do not move the watermark.
2. Compare full extract versus incremental for a 10 million row table. Write I/O on the source.
3. Write a data contract of six fields for `orders` in the warehouse (name, type, null, grain).

#### Advanced practical tasks

1. Write a one-page pipeline standard: ETL versus ELT choice, idempotency, checks, owners, and on-call.
2. Design a late-arriving fact (an order that updates after load). Write how the warehouse updates the fact.

---

## CDC (change data capture)

CDC reads changes from the operational database (usually the WAL or a change stream) and publishes those changes to consumers: warehouse, search, cache, another service.

```text
Primary WAL --> CDC reader --> events --> consumers
```

CDC benefits:

- low extra load versus a full-table poll
- near-real-time derived stores
- a full change story (insert, update, delete) if the product provides it

CDC needs:

- a stable primary key in the event
- a schema-change policy (a new column must appear in the stream)
- idempotent consumers (at-least-once is common)
- access to the log (privilege and retention)

CDC is not a backup. CDC is not a replica with a complete transaction story unless you design it that way. A consumer that applies events can lag. A consumer that drops an event must repair from a snapshot plus the stream.

Snapshot plus stream is the usual start: copy the table, then apply later events. The snapshot and the first event position must line up.

Do not parse application logs as CDC if you can read the database log. Application logs miss direct SQL and miss other writers.

Do not enable CDC on every table "just in case." Each table is storage, privilege, and consumer cost.

Deletes: a hard delete must produce a delete event or you keep ghosts in the warehouse. Soft deletes appear as updates. Agree on the meaning.

### Questions

#### Theoretical questions

1. Where does CDC usually read changes?
2. Name three consumers of CDC events.
3. Why is CDC lighter than a full-table poll?
4. Why must consumers be idempotent?
5. Why is CDC not a backup?

#### Easy practical tasks

1. Draw WAL, CDC reader, two consumers.
2. Write an event for `UPDATE products SET price = 10 WHERE id = 1`.
3. Find a CDC product or built-in logical decoding name for your DBMS (docs).
4. Write snapshot-then-stream in four sentences.

#### Medium practical tasks

1. Write how you line up a snapshot and the first log sequence.
2. Design a consumer for search: handle insert, update, delete. Write the idempotency key.
3. Write a policy for `ALTER TABLE ADD COLUMN` while CDC runs.

#### Advanced practical tasks

1. Write a one-page CDC runbook: privileges, retention, lag alert, snapshot repair, PII in the stream.
2. Compare poll-every-minute ETL with CDC for the same `orders` table. Write load, delay, and operational risk.

---

## Observability: connections, locks, replication lag, disk

Observability for a DBMS is a small set of signals that you watch and that you can act on.

**Connections.** Count of sessions, connect rate, waiting for a pool slot, rejected connections. Action: find a storm, a leak, or a low `max_connections`.

**Locks.** Wait time, blocked sessions, long transactions. Action: find the blocker, kill only with a rule, fix the query or the hot row.

**Replication lag.** Time or log-sequence behind. Action: stop replica reads, delay promote, find a heavy query on the replica.

**Disk.** Free bytes, IOPS, WAL volume, bloat. Action: grow the volume before it hits zero, find a missing vacuum, find a runaway log.

```text
Alert --> signal --> owner --> first action (from a runbook)
```

Add also: error rate, CPU, memory, backup age, last successful restore test, checkpoint time, cache hit (as a hint, not a goal).

A dashboard without a runbook is decoration. Each red signal needs a first step.

Logs and traces complete metrics. A lock metric without the blocking query text is incomplete.

Do not alert on every metric. Alert on symptoms that users feel or on risks that destroy data (disk full, backup fail, lag beyond RPO).

Do not collect only on the application. The DBMS can be the wait.

Capacity and observability work together. A disk-full alert is late if you never graph growth.

### Questions

#### Theoretical questions

1. What four signal groups does this section name first?
2. What is the first action when replication lag exceeds a budget?
3. Why is a dashboard without a runbook weak?
4. Which alerts protect data, not only latency?
5. Why is cache-hit ratio a hint and not a goal?

#### Easy practical tasks

1. Write one metric and one first action for each of the four groups.
2. Find a connection-count and a disk-free view in your DBMS docs.
3. Write three alerts that you would page a human for.
4. Write three graphs that you would only review weekly.

#### Medium practical tasks

1. Build a one-page dashboard list (names only) for a class production shop.
2. Simulate a disk-low condition in words: growth rate, days to full, who you call.
3. Pair a lock-wait metric with the query text that you would capture.

#### Advanced practical tasks

1. Write a runbook of four pages (or four short sections): connections, locks, lag, disk. Include "do not" steps.
2. Review a public DBMS exporter metric list. Pick 12 metrics. Map each to a user symptom or a data risk.

---

## Capacity: disk, IOPS, memory

Capacity planning answers "when do we run out?" for disk, I/O, and memory.

**Disk.** Data pages, WAL, dumps, snapshots, temp sorts. Leave headroom. A full disk stops writes. Measure growth per day. Project weeks to full.

**IOPS and throughput.** Disk operations per second and bytes per second. A cheap volume can meet size and fail latency. OLTP needs low-latency random I/O. Large scans need sequential throughput. WAL needs reliable writes.

**Memory.** Buffer cache, connections, sorts, the OS, the application if colocated. Do not colocate a hungry app and a DBMS on a small host without a budget.

```text
Used disk + daily growth * 30  <  0.7 * volume size   -- example headroom rule
```

Indexes use disk and they use cache. Extra indexes raise write IOPS.

Temp files appear when sorts do not fit in work memory. A sudden temp spike is a capacity event.

Plan for backup space on a different volume. A dump that fills the data disk is an outage.

Cloud volumes can grow. Growth is not instant in every product. Test a resize. Memory resize needs a restart in some products.

Do not buy CPU only because a query is slow. Measure. A missing index is cheaper than a larger host.

Write a capacity review on a schedule. Incidents that start as "slow disk" often started as "we never graphed IOPS."

### Questions

#### Theoretical questions

1. What happens when the data volume is full?
2. Why can a large disk still fail an OLTP workload?
3. Which memory consumers compete with the buffer cache?
4. Why must backups not fill the data volume?
5. Why is a missing index a capacity decision as well as a query decision?

#### Easy practical tasks

1. Write a headroom rule for disk in one formula (you may use the example).
2. List five things that consume disk besides table rows.
3. Find how to read free disk and IOPS for your environment (OS or cloud).
4. Write a memory budget in percents for a 16 GiB DBMS host (rough).

#### Medium practical tasks

1. Measure size of a practice database and WAL. Estimate days to fill a 10 GiB volume at that growth (invent a growth if needed, mark it).
2. Explain a query that spills to temp. Write how you see temp use in your product.
3. Compare two volume types (cloud docs): IOPS versus size. Pick one for OLTP and one for dumps.

#### Advanced practical tasks

1. Write a one-page capacity plan: 12-month row growth, indexes, backups, replica, IOPS, memory. State assumptions.
2. Design an alert: 30 days to disk full at current growth. Write the query or metric idea.

---

## Schema evolution in production

Schema evolution is the change of tables while the application stays up (or with a planned short pause). You use migrations. You pair them with application versions.

Expand-contract (additive change):

1. Add a nullable column or a new table. Deploy the migration.
2. Deploy application code that writes both old and new (or reads both).
3. Backfill data.
4. Deploy code that reads only the new shape.
5. Drop the old column in a later migration.

```text
v1 app + v1 schema
v1 app + v2 schema (compatible)
v2 app + v2 schema
v3 app + v3 schema (old objects gone)
```

Dangerous changes:

- `DROP COLUMN` while old code still reads it
- `ALTER` that rewrites a huge table and locks writes
- a new `NOT NULL` without a default and without a backfill
- a new unique constraint on dirty data
- a long transaction that holds a lock during a migration

Read the lock and rewrite behavior in your product. Some `ALTER` types are instant. Some are not. Test on a copy with production-like size.

Do not run an untested migration on production on Friday without a rollback plan. Rollback of a drop is a restore or a new migration, not a down script that recreates data.

Feature flags can hide new columns from users while you backfill.

Multiple application versions run during a rolling deploy. The schema must be compatible with both versions in that window.

Coordinate CDC, ORMs, and replicas. A new column must not break consumers.

### Questions

#### Theoretical questions

1. What is the expand step in expand-contract?
2. Why must two application versions work with one schema during a rolling deploy?
3. Why is `NOT NULL` without a backfill dangerous?
4. Why is a down migration a weak rollback after `DROP COLUMN`?
5. What must you test on a production-sized copy?

#### Easy practical tasks

1. Write a five-step expand-contract for rename of `phone` to `phone_e164`.
2. List five dangerous changes from this section.
3. Write one migration that you would refuse on Friday night and why.
4. Find one `ALTER` that the manual calls instant and one that it calls rewriting.

#### Medium practical tasks

1. Practice add-column, dual-write (or dual-read), backfill, drop-old on a disposable database with a small app or SQL only.
2. Add a unique constraint in two steps: clean duplicates, then `ADD CONSTRAINT`. Write the failed one-step error.
3. Write a migration review checklist of ten items (locks, size, rollback, CDC, ORM).

#### Advanced practical tasks

1. Write a one-page production-change standard: windows, copies, expand-contract, who approves `DROP`.
2. Plan a type change from `INTEGER` cents to `NUMERIC` money with no downtime. Write versions and backfill batches.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do OLTP, a warehouse, and CDC form one path from checkout to a report?
2. When do you choose ETL, ELT, or CDC for the same `orders` table?
3. Which observability signals tell you that capacity will fail before users feel it?
4. How does expand-contract protect a rolling deploy while a pipeline still reads the old column?
5. A teammate runs a year-long `GROUP BY` on the primary, grants analysts production `SELECT`, and drops a column in the same hour as a deploy. Which facts do you use in the reply?

#### Easy practical tasks

1. Write a cheat sheet: OLTP/OLAP, warehouse, ETL/ELT, CDC, four signals, capacity, expand-contract.
2. Draw the production data path for a shop with names on each arrow.
3. Write RPO for the operational database and freshness for the warehouse.
4. List disk, IOPS, and memory in one budget sentence each.

#### Medium practical tasks

1. Design a minimal platform: OLTP, replica, nightly ELT, three alerts, one expand-contract example.
2. Write a failed-load story (watermark not advanced) and a CDC-lag story. Write the user impact of each.
3. Size a 6-month disk plan from a made-up daily growth. Show the math.

#### Advanced practical tasks

1. Write an architecture-and-ops review of a class system against this topic. File gaps by risk.
2. Produce a 90-day operations plan: dashboards, restore test, capacity review, one schema change, one pipeline change.
