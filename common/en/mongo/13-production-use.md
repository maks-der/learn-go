# 13. Production Use

## Description

This topic shows habits and tools that keep a MongoDB application safe in production. You open change streams and you store resume tokens. You learn Atlas Search, time series collections, and GridFS. You change schemas with expand-contract, you run backfills, and you make writes idempotent. You choose a multi-tenant layout. You avoid common anti-patterns.

Complete modeling, drivers, indexes, replication, and performance first.

Use one term for each concept. A **change stream** is an official cursor of change events. A **resume token** lets you continue after a disconnect. **Expand-contract** is a schema change that adds the new shape, moves data, then removes the old shape. An **idempotent** write has the same result when you run it again. An **anti-pattern** is a design that looks convenient and then fails.

These patterns are application and operations work together. A good model still fails if a deploy cannot roll back.

---

## Change streams and resume tokens

A change stream is a long-lived cursor. The server sends a new event when a matching write commits. Change streams need a replica set or a sharded cluster. A standalone `mongod` is not enough.

Open a stream on a collection in `mongosh`:

```javascript
const cs = db.orders.watch()
cs.tryNext()
```

Drivers use `watch` or an equivalent helper. You can watch:

- One collection
- One database
- The whole deployment (with the right API and privileges)

A typical event includes:

- `_id` — the resume token
- `operationType` — `insert`, `update`, `replace`, `delete`, and others
- `fullDocument` — the document after the change when you request it
- `documentKey` — `_id` and shard key fields
- `updateDescription` — fields that changed, for some updates
- `clusterTime` — when the change occurred in the cluster

By default, an update event may omit the full document. Open the stream with `fullDocument: "updateLookup"` when you need the document after the update. That lookup has a cost and can race. Read the current options.

The stream follows the oplog. If the listener is down longer than the oplog window, the stream cannot resume. You must do a full resync of the downstream system.

Filter in `watch` with an aggregation pipeline (`$match` on `operationType` or on fields). Do not pull every event if you need only inserts.

Do not treat a change stream as a queue with infinite retention. It is not Kafka. Retention is the oplog window unless you add another store.

Do not run a change-stream worker that does slow work inside the receive loop without back-pressure. You can fall off the oplog.

A **resume token** is the `_id` of a change event. You pass it to `watch` so that the next cursor continues after that event.

```javascript
const token = event._id
const cs2 = db.orders.watch([], { resumeAfter: token })
```

Related options (names depend on version):

- `resumeAfter` — continue after that token
- `startAfter` — start after an invalidate in some cases
- `startAtOperationTime` — start at a cluster timestamp

Store the token only after you finish the downstream work for that event. If you store the token first and then crash, you can skip an event. If you never store the token, you replay events. Prefer **at-least-once** delivery and make the consumer **idempotent**.

Tokens are opaque. Do not build them by hand. Do not store a random ObjectId as a token.

If the token is outside the oplog, `watch` fails. You must snapshot the source collection and then start a new stream from a known time.

After `invalidate` (drop, rename, drop database), many cursors die. You may need `startAfter` or a new full sync.

Teams use change streams for **CDC** (change data capture) into search or a warehouse. A typical flow:

1. Take an initial snapshot of the collection.
2. Record a cluster time or a token.
3. Apply the snapshot to the sink.
4. Start the change stream from that time or token.
5. Apply events in order per document. Make each apply idempotent.

Do not tail `local.oplog.rs` in the application. Use the official change-stream API.

### Questions

#### Theoretical questions

1. What is a change stream?
2. Why does a standalone `mongod` reject `watch`?
3. What field is the resume token?
4. Why must you persist the token after downstream success?
5. What happens if the token is older than the oplog?

#### Easy practical tasks

1. Confirm you have a replica set. Open `watch` on a test collection. Insert one document. Print the event and the `operationType`.
2. Print `_id` from one change event. Write that it is the resume token.
3. Open the change-streams page and the resume page. Write both URLs.
4. Draw snapshot → token → stream → search upsert/delete.

#### Medium practical tasks

1. Open a stream with `fullDocument: "updateLookup"`. Update a field. Write whether the event includes the full document.
2. Receive one event. Close the cursor. Reopen with `resumeAfter`. Insert again. Confirm you do not see the first event twice (or write what you saw).
3. Store tokens in a small `cdc_state` collection. Write the document shape `{ listenerId, token, updatedAt }`.

#### Advanced practical tasks

1. Read change-stream event types for your version (including `invalidate`). Write when the cursor closes.
2. Design at-least-once CDC with a unique key on the sink. Write the retry and the skip rules.

---

## Atlas Search, time series, and GridFS

**Atlas Search** is a search engine that Atlas builds on **Apache Lucene**. You create a search index on a collection. You query with the aggregation stage `$search` (and related stages such as `$searchMeta`).

Atlas Search is not the same as a MongoDB text index (`$text`). The text index is a server feature with limited scoring and language support. Atlas Search is stronger for full-text relevance, fuzzy match, autocomplete, facets, and analyzers.

You define mappings. Dynamic mappings index many fields. Static mappings name the fields and the analyzer. Static mappings are easier to control.

A typical pipeline:

```javascript
db.articles.aggregate([
  {
    $search: {
      index: "default",
      text: { query: "replica set", path: "body" }
    }
  },
  { $limit: 20 },
  { $project: { title: 1, score: { $meta: "searchScore" } } }
])
```

Search is eventually consistent with the collection. A write is not always in the index in the same millisecond. Community Server does not include Atlas Search.

Do not create a search index on every field "just in case". Do not use `$search` as the first idea for `{ _id: id }`. Use `find`.

A **time series collection** stores documents that have a **time field** and, usually, a **meta field**. MongoDB groups documents into **buckets** on disk. The bucket layout saves space and can speed some range scans on time.

```javascript
db.createCollection("readings", {
  timeseries: {
    timeField: "ts",
    metaField: "sensorId",
    granularity: "seconds"
  }
})
```

Each measurement is still a document in the API. You insert one reading per event. You do not manage buckets by hand.

Good fit: sensors, metrics, click or event streams with a clear timestamp.

Poor fit: documents that you update often in place, highly relational operational data, data without a real time field.

You cannot treat a time series collection as a normal collection for every command. Some update and delete patterns are restricted. Read the limitations page before you migrate a hot collection.

Do not store a huge unbounded array of readings inside one sensor document.

The document size limit is 16 MB. **GridFS** stores a larger file as many **chunks** plus one **file** metadata document.

Default collections:

- `fs.files` — file name, length, content type, upload date
- `fs.chunks` — `{ files_id, n, data }` where `data` is a binary chunk

The default chunk size is 255 KB on many versions. Drivers implement upload and download. You do not assemble chunks by hand in ordinary application code.

Use GridFS when you must keep the file in MongoDB and the file is larger than 16 MB, or you want a stream API.

Do not use GridFS when object storage (S3 or similar) is available and the file is large or numerous. Do not use GridFS for a few kilobytes. Store a BinData field or a URL.

GridFS needs the indexes that the driver creates (`files_id` + `n` unique on chunks). Do not drop them.

Backup size grows with every file. A dump of GridFS is heavy.

Pick the tool that matches the problem. Do not use GridFS for 1 KB JSON. Do not use Atlas Search for a single equality on `_id`.

### Questions

#### Theoretical questions

1. How is Atlas Search different from a `$text` index?
2. What aggregation stage runs an Atlas Search query?
3. What two fields do you set when you create a time series collection?
4. What size limit makes GridFS necessary for one file?
5. When is object storage a better place for the file?

#### Easy practical tasks

1. Open the Atlas Search, time series, and GridFS pages. Write the three URLs.
2. Write three query types that fit Lucene search and three that fit `find`.
3. Write a sample reading document with `ts` and `sensorId`.
4. Draw `fs.files` (one row) and three `fs.chunks` for one file.

#### Medium practical tasks

1. Create a search index on a test collection (or write the JSON mapping). Run one `$search`. Write the first `_id`.
2. Create a time series collection. Insert 100 readings. Run a time-range find. Write `explain` stage names if you can.
3. Compare the bucket pattern from topic 6 with a time series collection. Write two similarities and two differences.

#### Advanced practical tasks

1. Design mappings for a product catalog: title, description, sku, facets for brand. Write static mappings in outline form.
2. Design retention: 7 days hot, TTL, and a monthly archive collection. Write the jobs. Then compare GridFS vs S3 for the same files.

---

## Expand-contract changes, backfills, and idempotent writes

MongoDB does not require a table migration for a new field. The application still needs a plan. Old processes and new processes can run at the same time during a deploy.

**Expand-contract** (also called expand-migrate-contract):

1. **Expand.** Add the new field or the new collection. Deploy code that writes both the old shape and the new shape, or that writes the new field and still reads the old field.
2. **Migrate.** Backfill old documents. Confirm readers can use the new field.
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

Validation (JSON Schema) must expand first too. A strict schema that requires the new field will reject old documents. Use `moderate` or a weaker rule during the move if you must.

A **backfill** is a batch job that updates existing documents. You use it after a schema expand, after a bug, or after you add a computed field.

Safe backfill habits:

- Work in **batches** (`find` with a range on `_id`, then `bulkWrite`)
- Use a filter that selects only documents that still need the change
- Make the job **idempotent**
- Limit rate so that you do not fill the WiredTiger cache or the oplog with a spike
- Log progress (last `_id`, count, errors)
- Prefer a secondary for the **read** of the scan if you accept lag, but **write** to the primary

Example idea:

```javascript
db.orders.updateMany(
  { quantity: { $exists: false }, n: { $exists: true } },
  [ { $set: { quantity: "$n" } } ]
)
```

A single `updateMany` on tens of millions of documents can be too large. Prefer looped batches with a pause.

Do not backfill in the web request path. Do not hold a multi-document transaction across a huge backfill.

An **idempotent** write leaves the same stored state when you send it once or many times with the same intent.

Retries happen. Drivers retry some writes. Users double-click. Queues redeliver. Change-stream consumers replay. Without idempotency you get duplicate orders or double `$inc`.

Tools:

- A **natural key** and a **unique index**. `insertOne` of `{ _id: orderId }` fails the second time. Treat duplicate key as success if the document is the same.
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

### Questions

#### Theoretical questions

1. What are the three phases of expand-contract?
2. Why can a new reader fail if it requires a field that old documents lack?
3. Why do you batch a backfill by `_id` range?
4. What does idempotent mean for a write?
5. Why is `$inc` dangerous on a retried request?

#### Easy practical tasks

1. Write expand-contract steps for `email` → `emailAddress`.
2. Write a filter that finds documents with `n` and without `quantity`.
3. Write two clicks of "Place order" with the same `orderId`. Write the unique index.
4. Make a table: phase, what writers do, what readers do.

#### Medium practical tasks

1. In a lab, insert documents with `n`. Deploy a script that writes both fields. Backfill in batches of 100. Read only `quantity`.
2. Implement upsert by `idempotencyKey`. Send the same body twice. Prove one document.
3. Show a double `$inc` bug. Then fix it with a request-id set or a ledger document. Write both.

#### Advanced practical tasks

1. Plan a type change (string id to ObjectId) with expand-contract. Include indexes, queries, and a resume for the backfill (`lastId`).
2. Design webhook ingest: unique `eventId`, upsert, and a side effect that must not run twice (use an outbox flag).

---

## Multi-tenant layout

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

## Anti-patterns: unbounded arrays, huge documents, low-cardinality shard keys, SQL-shaped usage

An **anti-pattern** is a design that looks convenient and then fails. If you already use an anti-pattern, you change the model. You do not add more hardware first.

An **unbounded array** is a field that receives `push` on every event with no cap. Examples: all comments on a viral post, all page views, all log lines, all payments on one customer document.

Problems:

- Each update rewrites a larger document
- The document approaches 16 MB and then writes fail
- Indexes on the array become large (multikey)
- Every reader loads the whole array even when they need the last five items
- Concurrent `$push` and other updates contend on one document

The modeling rule remains: data that you access together can live together **if the array stays bounded**.

Bounded patterns:

- Keep the last N items (`$slice` on `$push`)
- Store children in a second collection with a parent id
- Use the bucket pattern or a time series collection for events
- Use the outlier pattern for one popular parent

A "small" array that grows 1 item per day becomes large in a few years. Write the cap in the design.

Do not `$push` a view event into the user document for every click.

The BSON **document size limit is 16 MB**. A document that is "only" 10 MB is already a problem. It uses cache, bandwidth, and lock time on that document.

Causes of huge documents: unbounded arrays, embedded files, huge strings (HTML dumps, base64 images), nested copies of the same related data.

Prevention: store binaries outside the document, reference large related graphs, use the subset pattern, validate max length in the application.

Do not store a PDF as a 14 MB BinData field because "it is under 16 MB". If you are near the limit, you are late. Split before you hit the error in production.

A **low-cardinality** shard key has few distinct values. Examples: `country` with three values, `status` with four values, `isActive` with two values.

Effects:

- Chunks cannot split past the distinct values
- Those chunks become **jumbo**
- All documents with `status: "open"` live on one shard
- That shard is **hot**
- The balancer cannot save you

Hashing a boolean does not create cardinality. You still have two hashes.

A compound key `{ status: 1, orderId: 1 }` can add cardinality **if** `orderId` is distinct. The prefix `status` still groups opens together.

Do not shard on a field because "every query has it" if the field has five values.

If you already sharded on a poor key, use refine or reshard. That is a project.

A **drop-in SQL replacement** is the idea that you keep tables, normalize every entity, join in every read, and only swap the database.

MongoDB does not run SQL as the primary language. `$lookup` is not a full SQL join planner for every ad-hoc report. There are no declarative foreign keys as the main integrity tool. Transactions exist but they are not the default write.

If you map each table to a collection one-for-one, you often get many round trips or many `$lookup` stages, slow reports that SQL did better, a schema that never uses embed, and shard keys that do not match access.

MongoDB fits when you design **documents for access**. Embed what you read together. Reference what you share and update independently. Use a warehouse or a relational database for heavy ad-hoc analytics if that is the workload.

You can migrate from SQL. You must **redesign**. You do not only migrate types.

Do not implement a SQL emulator on MongoDB for the whole application.

Do not forbid embed because "normalization is always right". Do not embed everything because "joins are forbidden". Both extremes are anti-patterns.

A related anti-pattern is to wrap **every** write in a multi-document transaction. Single-document writes are already atomic. Use a transaction when two documents must change together and you cannot model them as one document.

### Questions

#### Theoretical questions

1. What is an unbounded array?
2. Why is a 10 MB document already a risk?
3. What does low cardinality mean for a shard key?
4. Why is a 1:1 table-to-collection map often slow?
5. Why must a SQL migration include a redesign?

#### Easy practical tasks

1. Write three unbounded-array examples and three bounded-array examples.
2. Insert a document. Run `Object.bsonsize` on it in `mongosh`. Write the number.
3. Write five poor shard keys and one reason each.
4. Write five SQL habits and the MongoDB habit that replaces each one.

#### Medium practical tasks

1. `$push` in a loop until the document is large (stay safe in a lab). Time the last 10 pushes vs the first 10. Then rewrite the same data as a comments collection.
2. On paper, estimate cardinality of `tenantId`, `country`, `orderId`, `createdAt` for an app that you know. Propose a better key than `{ type: 1 }`.
3. Take a 5-table relational schema. Write a document design. List joins that disappear.

#### Advanced practical tasks

1. Design a migration from an unbounded `events[]` to a time series or bucket collection. Include expand-contract.
2. Write a one-page "we are not a SQL database" guide for teammates who know only PostgreSQL. Include when a warehouse stays in the system.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do the oplog, change streams, and resume tokens form one pipeline?
2. How do expand-contract, backfills, and idempotent writes work together in one rename-and-retry deploy?
3. When does a `tenantId` model fail sharding in a way that a database-per-tenant model would not?
4. How do unbounded arrays and huge documents become the same incident at 16 MB?
5. A teammate tails `local.oplog.rs`, shards on `status`, wraps every insert in a transaction, and `$push`es every click onto the user. Which facts do you use in the review?

#### Easy practical tasks

1. Write a cheat sheet: `watch` targets, token persist order, `$search` vs `find`, time series fields, GridFS collections, expand-contract phases, three tenant layouts, five anti-patterns.
2. Open one stream. Capture one insert event. Save the JSON (redact secrets).
3. Draw a deploy timeline: dual-write, backfill, contract, unset.
4. Audit one of your collections for the anti-patterns in this topic. Write a yes/no row for each.

#### Medium practical tasks

1. Build a small listener that prints `operationType` and `_id`. Restart it with a saved token. Prove continuity.
2. Run a small expand-contract plus batched backfill on a lab collection. Add a unique index for an idempotency key. Prove a second insert fails safely.
3. Pick one anti-pattern in a lab. Implement the bad version and the better version. Write sizes or times.

#### Advanced practical tasks

1. Implement idempotent CDC into a second MongoDB collection that mirrors `orders`. Handle insert, update, replace, delete. Resume after kill.
2. Produce a production handbook page: schema change, backfill SLA, idempotency, tenant isolation tests, change-stream token storage, and a named exception process for transactions or shard keys.
