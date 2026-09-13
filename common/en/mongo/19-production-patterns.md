# 19. Production Patterns

## Description

This topic shows habits that keep a MongoDB application safe in production. You change schemas with **expand-contract**. You run **backfills**. You make writes **idempotent**. You choose a **multi-tenant** layout. You watch connections, replication lag, page faults, and slow operations. Complete modeling, drivers, indexes, replication, and performance first.

Use one term for each concept. **Expand-contract** is a schema change that adds the new shape, moves data, then removes the old shape. A **backfill** is a job that writes missing or new fields on existing documents. An **idempotent** write has the same result when you run it again. **Observability** is the practice of measuring the system so that you can see a problem.

These patterns are application and operations work together. A good model still fails if a deploy cannot roll back.

---

## Expand-contract schema changes

MongoDB does not require a table migration for a new field. The application still needs a plan. Old processes and new processes can run at the same time during a deploy.

**Expand-contract** (also called expand-migrate-contract):

1. **Expand.** Add the new field or the new collection. Deploy code that writes both the old shape and the new shape, or that writes the new field and still reads the old field.
2. **Migrate.** Backfill old documents (next section). Confirm readers can use the new field.
3. **Contract.** Deploy code that reads only the new field. Then stop writing the old field. Remove the old field later if you need the space.

Example: rename `n` to `quantity`.

- Version A writes `n` only
- Version B writes `n` and `quantity`
- Backfill copies `n` to `quantity` where missing
- Version C reads `quantity` and writes `quantity` only
- A later job `$unset`s `n`

Do not deploy a reader that requires `quantity` before the backfill finishes, unless the reader has a fallback.

Do not drop a field in the same release that first writes the new field.

A **schema version field** (`schemaVersion: 2`) helps the application pick a decode path. You still need expand-contract for rolling deploys.

Validation (JSON Schema) must expand first too. A strict schema that requires the new field will reject old documents. Raise `validationLevel` with care. Use `moderate` or a weaker rule during the move if you must.

### Questions

#### Theoretical questions

1. What are the three phases of expand-contract?
2. Why can a new reader fail if it requires a field that old documents lack?
3. Why do you write both fields in the expand phase of a rename?
4. How does `schemaVersion` help a rolling deploy?
5. Why can a strict JSON Schema block a migration?

#### Easy practical tasks

1. Write expand-contract steps for `email` → `emailAddress`.
2. Open a MongoDB schema-evolution or versioning page if you find one. Write the URL or write that you used this handbook.
3. Draw versions A, B, C for the `n` / `quantity` example.
4. Make a table: phase, what writers do, what readers do.

#### Medium practical tasks

1. In a lab, insert documents with `n`. Deploy a script that writes both fields. Backfill. Read only `quantity`.
2. Add `schemaVersion`. Write a reader that handles version 1 and version 2.
3. Change a JSON Schema from optional new field to required. Write when you flip the switch.

#### Advanced practical tasks

1. Plan a type change (string id to ObjectId) with expand-contract. Include indexes and queries.
2. Write a team checklist: expand, dual-write, backfill bar, contract, unset, validation.

---

## Backfills

A **backfill** is a batch job that updates existing documents. You use it after a schema expand, after a bug, or after you add a computed field.

Safe backfill habits:

- Work in **batches** (`find` with a range on `_id`, then `bulkWrite`)
- Use a filter that selects only documents that still need the change
- Make the job **idempotent** (next section)
- Limit rate so that you do not fill the WiredTiger cache or the oplog with a spike
- Log progress (last `_id`, count, errors)
- Prefer a secondary for the **read** of the scan if you accept lag, but **write** to the primary
- Run in a replica set so that you can abort and retry

Example idea:

```javascript
db.orders.updateMany(
  { quantity: { $exists: false }, n: { $exists: true } },
  [ { $set: { quantity: "$n" } } ]
)
```

A single `updateMany` on tens of millions of documents can be too large. Prefer looped batches with a pause.

Do not backfill in the web request path.

Do not start a backfill on Friday without a stop switch.

Do not hold a multi-document transaction across a huge backfill.

If you added an index that the backfill needs, create the index first. If the backfill makes an index useless, drop the old index in the contract phase.

Atlas has online tools and you can still run a script from a worker. Measure replication lag during the job.

### Questions

#### Theoretical questions

1. What does a backfill change?
2. Why do you batch by `_id` range?
3. Why must the filter skip documents that already have the new field?
4. Why is a backfill in an HTTP handler a poor design?
5. Why can one huge `updateMany` hurt the cluster?

#### Easy practical tasks

1. Write a filter that finds documents with `n` and without `quantity`.
2. Open `bulkWrite` or `updateMany` docs. Write the URL.
3. Write five log fields for a backfill worker.
4. Make a table: safe habit, risk if you skip it. Add four rows.

#### Medium practical tasks

1. Backfill 1000 lab documents in batches of 100. Write the loop and the duration.
2. Pause 50 ms between batches. Compare lag or time with a full `updateMany`. Write the observation.
3. Write a stop file or an env flag that the worker checks each batch.

#### Advanced practical tasks

1. Design a resume: store `lastId` in a `jobs` collection. Kill the process. Start again. Prove no double-bad writes.
2. Estimate oplog impact for 50 million updates. Propose a rate limit and a window.

---

## Idempotent writes

An **idempotent** write leaves the same stored state when you send it once or many times with the same intent.

Retries happen. Drivers retry some writes. Users double-click. Queues redeliver. Change-stream consumers replay. Without idempotency you get duplicate orders or double `$inc`.

Tools:

- A **natural key** and a **unique index**. `insertOne` of `{ _id: orderId }` or `{ idempotencyKey }` fails the second time. Treat duplicate key as success if the document is the same.
- **`upsert`** with a filter that is the identity of the entity.
- **`$set`** of an absolute value instead of `$inc` when you can send the final number.
- **Compare-and-set** (`updateOne` with a version field) when two writers race.
- For `$inc`, include a **dedupe set** (bounded) of request ids, or store one ledger document per request.

Example:

```javascript
db.orders.updateOne(
  { _id: "ord_123" },
  { $setOnInsert: { total: 10, createdAt: new Date() } },
  { upsert: true }
)
```

The second call does not create a second order.

Do not generate a new ObjectId in the client for each retry of the same user action. Send the id that the client created once.

Do not `$inc` a balance on every retry of the same payment.

Idempotency is not only for money. Email send flags and webhook handlers need it too.

Combine this with write concern majority for payments so that a retry after a timeout does not hide a committed write. Read the driver retry and `retryWrites` notes.

### Questions

#### Theoretical questions

1. What does idempotent mean for a write?
2. Why do retries create duplicates without a unique key?
3. How does `$setOnInsert` help an upsert?
4. Why is `$inc` dangerous on a retried request?
5. Why must the client reuse the same `_id` for the same action?

#### Easy practical tasks

1. Write two clicks of "Place order" with the same `orderId`. Write the unique index.
2. Open the upsert page. Write the URL.
3. Make a table: write type, idempotent method. Add insert, increment, email flag.
4. Write four sentences: retry, unique index, duplicate key, success.

#### Medium practical tasks

1. Implement upsert by `idempotencyKey`. Send the same body twice. Prove one document.
2. Show a double `$inc` bug. Then fix it with a request-id set or a ledger document. Write both.
3. Read `retryWrites` in your driver. Write the default and one error that retries.

#### Advanced practical tasks

1. Design webhook ingest: unique `eventId`, upsert, and a side effect that must not run twice (use an outbox flag).
2. Write a team standard: every POST that creates money or email must name its idempotency key.

---

## Multi-tenant: database vs collection vs `tenantId`

**Multi-tenant** means one application serves many customers (tenants). You must isolate data.

Three common layouts:

**Database per tenant.** `tenant_acme`, `tenant_globex`. Strong isolation. Custom indexes and backups per tenant. Many databases can stress metadata and connections. Best for few large tenants or strict isolation.

**Collection per tenant.** `orders_acme`. Isolation is weaker than a database. Collection count can explode. Rarely the best default.

**`tenantId` field** on shared collections. `{ tenantId: "acme", ... }`. One set of indexes. Every query **must** include `tenantId`. A unique index is `{ tenantId: 1, orderId: 1 }`, not `{ orderId: 1 }` alone if ids can collide. This is the usual default for many small tenants.

Rules for `tenantId`:

- The application sets `tenantId` from the session, not from the client body alone
- Indexes start with `tenantId` or include it
- If you shard, `tenantId` is often the prefix of the shard key (watch hot tenants)
- Encryption and RBAC can add a layer (user per tenant is hard at large scale; often the app enforces the filter)

Do not mix layouts without a written plan.

Do not trust `tenantId` from the JSON body without a server-side check.

A huge tenant can need its own database or shard later (outlier). Expand-contract applies to a move of one tenant.

Atlas and `maxIncomingConnections` still apply. One thousand databases can hurt. Measure.

### Questions

#### Theoretical questions

1. What does multi-tenant mean?
2. When is a database per tenant a good fit?
3. Why must every query include `tenantId` on a shared collection?
4. Why is a unique index on `orderId` alone unsafe if two tenants can use the same id?
5. Why must the server set `tenantId` from the session?

#### Easy practical tasks

1. Make a table: database per tenant, collection per tenant, `tenantId`. Add isolation, scale, ops cost.
2. Write a unique index for orders in a shared collection.
3. Draw a leaked query that omits `tenantId`. Mark the bug.
4. Open a MongoDB multi-tenant planning page if you find one. Write the URL or your notes source.

#### Medium practical tasks

1. Insert orders for two tenants. Write a find that cannot leak. Write a unit test that fails if `tenantId` is missing.
2. Plan indexes for `orders` and `users` with `tenantId` first.
3. Estimate 5000 tenants × 3 collections vs one shared set. Write a metadata concern.

#### Advanced practical tasks

1. Design a move of one huge tenant from shared collections to its own database. Use expand-contract and dual-write.
2. Combine RBAC and `tenantId`: what the database user can do vs what the app must enforce. Write a one-page policy.

---

## Observability: connections, replication lag, page faults, slow ops

**Observability** means you can see health with metrics, logs, and traces. For MongoDB, watch at least these signals.

**Connections.** Each client uses a pool. Too many application instances × large pools exhaust `maxIncomingConnections`. Metrics: current connections, available, percent used. Fix: smaller pools, fewer idle clients, do not create a `MongoClient` per request.

**Replication lag.** Secondaries behind the primary. Metrics: optime gap, Atlas replication lag. Causes: slow disk, network, huge writes, a backfill. Effects: stale reads, delayed backups, risk to the oplog window. Alert when lag stays high.

**Page faults / cache pressure.** The working set does not fit. On some OS setups you see page faults. On all setups you see latency, disk IOPS, and WiredTiger cache eviction. Fix: RAM, smaller documents, fewer indexes, archive, or shard.

**Slow operations.** Profiler, slow log, Atlas Query Insights. Look for `COLLSCAN`, huge sorts, and long writes. Fix with `explain` and indexes (topic 14).

Also watch:

- Disk free space
- CPU
- Oplog window
- Index build progress
- Ticket queues (WiredTiger read/write tickets) on busy servers

Put dashboards in Atlas or in your metrics stack (`serverStatus` scrape). Page a human on lag and on connection saturation, not only on CPU.

Do not alert on a single slow query in a lab. Do alert on a sustained change from the baseline.

Logs must not contain passwords. Redact.

### Questions

#### Theoretical questions

1. Why can many `MongoClient` objects exhaust connections?
2. What user-visible bug does replication lag cause on secondary reads?
3. What does cache pressure suggest about the working set?
4. Where do you find slow operations on Atlas?
5. Why is a baseline more useful than a single spike?

#### Easy practical tasks

1. Write four signals from this section and one metric name each.
2. Open Atlas Metrics or `serverStatus` help. Write the URL.
3. On your cluster, write current connection count if you can see it.
4. Make a table: signal, danger threshold idea, first check.

#### Medium practical tasks

1. Create a dashboard sketch: connections, lag, disk, slow query count. Mark alert lines.
2. Cause a slow `COLLSCAN` in a lab. Find it in the profiler or logs. Write the millis.
3. Read `maxIncomingConnections` and your driver pool size. Compute instances × pool.

#### Advanced practical tasks

1. Write an on-call runbook: high lag, connections 90%, disk 85%, p99 latency up. Four short procedures.
2. Instrument one API with a request id and log the MongoDB operation time. Sample 100 requests. Write p50 and p99.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do expand-contract, backfills, and idempotent writes work together in one rename-and-retry deploy?
2. When does a `tenantId` model fail observability or sharding in a way that a database-per-tenant model would not?
3. Why can a backfill without rate limits create lag, cache pressure, and slow ops at the same time?
4. How does a unique idempotency key protect you after a driver retry and after a user double-click?
5. A teammate ships a reader that requires a new field, runs `updateMany` on all docs at noon, and opens a new `MongoClient` per request. Which facts do you use in the review?

#### Easy practical tasks

1. Write a cheat sheet: expand-contract phases, batch backfill, unique key, three tenant layouts, four metrics.
2. Draw a deploy timeline: dual-write, backfill, contract, unset.
3. List five production writes in an app and mark idempotent or not.
4. Write your tenant choice for a fictional SaaS and one reason.

#### Medium practical tasks

1. Run a small expand-contract plus batched backfill on a lab collection. Record versions and counts.
2. Add a unique index for an idempotency key. Prove a second insert fails safely.
3. Write alert text for lag and for connections that an on-call engineer can follow at 03:00.

#### Advanced practical tasks

1. Produce a production handbook page: schema change, backfill SLA, idempotency, tenant isolation tests, dashboards, who approves.
2. Load-test a backfill and an API together. Write lag, connection count, and the rate limit that kept p99 stable.
