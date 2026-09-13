# 11. Transactions and Consistency

## Description

This topic shows how MongoDB keeps writes correct. A single document write is always atomic. Multi-document **ACID** transactions exist on replica sets and on sharded clusters. You learn when you do not need a transaction. You set read concern, write concern, and read preference. You learn causal consistency at a high level. Complete CRUD, modeling, and drivers first.

Use one term for each concept. **Atomic** means the write to one document fully happens or does not happen. A **transaction** groups reads and writes across documents so that they commit or abort together. **Write concern** is how durable and how widely replicated a write must be before the server acknowledges it. **Read concern** is how up-to-date and how majority-committed a read must be. **Read preference** selects which replica you read from.

---

## Single-document atomicity (always)

A write to **one document** is atomic. The server does not leave half of the fields from that write on disk. `$set` of two fields in one `updateOne` applies both or neither.

This is true without a multi-document transaction. It is the first consistency tool in MongoDB.

Design for it. Put data that must change together in one document when the model allows it. An order with lines in one document can update `status` and `lines` in one write.

Atomicity is not isolation across documents. Two documents can change in an order that a second client sees as mixed.

Atomicity is not durability by itself. A write can be atomic in memory and still be lost if the node fails before the journal or the replica majority, depending on write concern.

`updateOne` of one match, `insertOne`, and `replaceOne` are single-document. `updateMany` is atomic **per document**, not as one all-or-nothing batch, unless you use a transaction.

Do not use a transaction only to change two fields in the same document. One update is enough.

### Questions

#### Theoretical questions

1. What does single-document atomicity guarantee?
2. Why can an order with embedded lines avoid a transaction?
3. Does atomicity mean the write is on a majority of replicas?
4. Is `updateMany` one atomic operation for all matches?
5. When is a transaction unnecessary for two field changes?

#### Easy practical tasks

1. `$set` two fields in one `updateOne`. Find the document. Confirm both fields changed.
2. Write one embed design that keeps a money transfer of two balances **out** of one document (they are two accounts). Write why atomicity of one document is not enough.
3. Write one embed design that keeps a status and a log array **in** one document. Write the one update.
4. Open the atomicity page in the manual. Write the URL.

#### Medium practical tasks

1. Run `updateMany` with `$inc` on three documents. Kill the client in the middle if you can, or discuss partial completion. Write what "per document" means.
2. Compare two `$set` calls vs one `$set` with two fields. Write an interleaving that a second reader can see in the two-call case.
3. Show `$inc` on one field as a single-document atomic counter. Write why the client must not read-add-write.

#### Advanced practical tasks

1. Read about document-level concurrency and `findAndModify`. Write one use (for example claim a job document).
2. Design a shopping cart as one document. Write which checkout steps still need a transaction (inventory in another collection).

---

## Multi-document ACID transactions (replica set / sharded)

A **multi-document transaction** lets you read and write more than one document as one commit.

Requirements:

- A **replica set** (even one member for learning) or a **sharded cluster**
- A driver API: `startSession`, `startTransaction`, operations, `commitTransaction` or `abortTransaction`
- Standalone `mongod` without a replica set does **not** support these transactions

ACID in this context:

- **Atomic** — all writes in the commit appear, or none appear
- **Consistency** — you still define application consistency; the server enforces constraints that you configured
- **Isolation** — snapshot isolation for the transaction (read the current manual for details)
- **Durability** — depends on write concern at commit

Example shape in `mongosh`:

```javascript
session = db.getMongo().startSession()
session.startTransaction()
const orders = session.getDatabase("shop").orders
const stock = session.getDatabase("shop").stock
orders.insertOne({ _id: 1, sku: "nail", qty: 1 })
stock.updateOne({ sku: "nail" }, { $inc: { qty: -1 } })
session.commitTransaction()
session.endSession()
```

If a write fails, abort. Retry the whole transaction on transient errors (the driver can help).

Transactions have limits: time, size, and number of operations. Do not hold a transaction while you wait for a user.

Sharded transactions are supported in current versions. They cost more. Keep them short.

### Questions

#### Theoretical questions

1. Why does a standalone `mongod` without a replica set reject multi-document transactions?
2. What does commit vs abort mean for the two writes in the example?
3. Why must you not wait for a user inside a transaction?
4. What is a session in this API?
5. Why can a sharded transaction cost more than a replica-set transaction?

#### Easy practical tasks

1. Confirm you have a replica set (`hello` or Atlas). Write if transactions are available.
2. Run a two-collection transaction that commits. Find both documents.
3. Run a transaction that aborts. Confirm neither write remains (or the abort left no change).
4. Open the transactions page. Write the URL.

#### Medium practical tasks

1. Cause a duplicate-key error inside a transaction. Abort. Write the collection state.
2. Retry a transaction after a transient error (or simulate abort and retry). Write the loop.
3. Compare the same two writes without a transaction. Show a state where only one write exists.

#### Advanced practical tasks

1. Read transaction limits (runtime, oplog). Write two limits for your version.
2. Run a transaction that touches two shards (if you have a sharded cluster) or write the extra care from the manual if you do not.

---

## When you do not need a multi-document transaction

You do not need a multi-document transaction when:

- One document write is enough (embed, computed field on the same document)
- The business accepts **eventual** consistency and a retry job (for example a log line that can arrive late)
- You can make the write **idempotent** and you can repair (outbox, retry)
- You only need a unique index to prevent a duplicate
- You read after write in a way that a session and majority concern already cover (see later sections)

Transactions have a cost. They hold locks or keep a snapshot. They can abort. They make the code harder.

Prefer model changes over transactions. If you always transaction-update a user and a counter, store the counter on the user, or accept a nightly repair.

Do not wrap every `insertOne` in a transaction.

Do not use a transaction to hide an unbounded array problem.

A single `bulkWrite` without a transaction is **not** all-or-nothing across documents.

When you do need a transaction: money between two accounts in two documents, inventory plus order if they are separate, two collections that must not show a mixed state to any reader.

### Questions

#### Theoretical questions

1. Why is a model change often better than a transaction?
2. When is eventual consistency acceptable?
3. Why is `bulkWrite` without a transaction not ACID across documents?
4. Why is "transaction on every insert" a bad default?
5. Give one domain that still needs a multi-document transaction.

#### Easy practical tasks

1. List five writes in a shop API. Mark transaction or single document.
2. Rewrite one two-document update as one document. Write the new shape.
3. Write an idempotent "create order with this `orderId`" without a transaction. Use a unique index.
4. Make a two-column table: "Need transaction" and "Do not". Add four rows.

#### Medium practical tasks

1. Design an outbox: write the order and an `outbox` event in one transaction, or show a single-document outbox embed. Write the trade-off.
2. Review a sample service that uses transactions for logging. Propose removal. Write the risk.
3. Show a unique-index conflict as a replacement for a "check then insert" transaction.

#### Advanced practical tasks

1. Read official guidance "when to use transactions". Write three sentences from that page in your own words.
2. Design account transfers. Compare one document per account plus transaction vs an event-sourced ledger collection. Write one page.

---

## Read and write concerns

**Write concern** is the acknowledgment rule for a write.

Common values:

- `w: 1` — the primary acknowledges. Fast. A primary crash can lose the write if it did not replicate.
- `w: "majority"` — a majority of the replica set voting members persist the write. Safer.
- `j: true` — the write is on the journal.

Atlas often defaults to majority. Set majority for money and for data that you cannot rebuild.

**Read concern** is the snapshot rule for a read.

Common values:

- `local` — data on the selected node, even if it is not majority-committed
- `majority` — data that a majority committed
- `snapshot` — a consistent snapshot, used in transactions

A read with `local` on the primary can still see a write that later rolls back in rare failover cases if the write was not majority-committed.

Do not mix `w: 1` writes with a belief that every other node already has the data.

You can set concerns in the URI, on the client, on the collection, or on one operation. Be explicit for important writes.

`w: 0` (unacknowledged) is for special cases. Do not use it in an API that must know the write worked.

### Questions

#### Theoretical questions

1. What does `w: "majority"` wait for?
2. What risk does `w: 1` have on failover?
3. What does read concern `majority` mean?
4. When do you use read concern `snapshot`?
5. Why is `w: 0` a poor default for an API?

#### Easy practical tasks

1. Insert with `{ writeConcern: { w: "majority" } }` on a replica set. Confirm success.
2. Find the default write concern for your Atlas cluster or your replica set. Write it.
3. Open the write-concern and read-concern pages. Write both URLs.
4. Make a table: concern name, what it waits for, one risk if you pick a weaker value.

#### Medium practical tasks

1. Set read concern `majority` on a find. Write the option in your driver.
2. Compare latency of `w: 1` vs `w: "majority"` on 100 inserts. Write the two times.
3. Explain a rollback of a `w: 1` write in words (primary dies before replication). Draw the timeline.

#### Advanced practical tasks

1. Read about `wtimeout`. Cause a timeout with an impossible `w` (for example `w: 5` on a three-node set). Record the error.
2. Read write concern majority and journaling. Write how `j` and `w` work together on WiredTiger.

---

## Read preference

**Read preference** selects which member of a replica set runs the read.

Common modes:

- `primary` — always the primary. Default. Strongest freshness on that node.
- `primaryPreferred` — primary if available, else a secondary
- `secondary` — a secondary
- `secondaryPreferred` — a secondary if available, else the primary
- `nearest` — a member with low latency that matches tags

Reads from a **secondary** can be **stale**. Replication is not instant. The application can miss a write that it just sent to the primary.

Use `primary` for reads that must see the latest write of this application (or use a session with causal consistency, next section).

Use secondaries for reports that accept lag, or to move read load. Measure lag. Do not guess.

`nearest` can still hit a secondary. Do not use `nearest` for "latest" reads.

Tag sets select members (for example a delayed member). Delayed members are for disaster practice, not for normal user reads.

Read preference is not read concern. Preference picks a node. Concern picks how committed the data must be on that node.

### Questions

#### Theoretical questions

1. What is the default read preference?
2. Why can a secondary miss a write that just succeeded?
3. When is `secondary` acceptable?
4. Why is `nearest` not a freshness guarantee?
5. How is read preference different from read concern?

#### Easy practical tasks

1. Run `hello` or `rs.status()`. Write which member is primary.
2. Set read preference `secondary` in `mongosh` or the driver (`db.getMongo().setReadPref("secondary")`). Run a find. Write which host you hit if the shell shows it.
3. Open the read-preference page. Write the five mode names.
4. Make a table: mode, can be stale, typical use.

#### Medium practical tasks

1. Insert on primary. Immediately read with `secondary`. Repeat until you see lag or write that you did not see lag. Write the observation.
2. Set `secondaryPreferred` and stop secondaries in a test set (if you can). Write where reads go.
3. Compare report-query load on primary vs secondary in a drawing. Mark the lag risk.

#### Advanced practical tasks

1. Read about max staleness (`maxStalenessSeconds`). Write how it limits how stale a secondary can be.
2. Design read preference per endpoint: checkout vs public catalog vs analytics. Write three modes.

---

## Causal consistency (high-level)

**Causal consistency** means that if operation B depends on operation A, every reader that uses the same rules sees A before B.

Example: a client writes a post, then reads the post list. Without care, a read on a stale secondary can omit the post.

MongoDB **sessions** can provide causal consistency. After a write, the session stores a cluster time. The next read in that session waits until the selected members are at least that far.

Default: many drivers enable causal consistency on a session. Read the driver page.

Causal consistency is not a full linearizability of the whole database for all clients. Another client without the session might still see a different order in edge cases, depending on concerns and preference.

You still pick write concern and read concern. Causal consistency plus majority concerns is a common pair for "read your writes" in an app session.

Do not confuse causal consistency with a multi-document transaction. A transaction is a different tool. You can use both.

High-level rule: use one session for a user's request chain when the next read must see the previous write, especially if you read from secondaries.

### Questions

#### Theoretical questions

1. What does causal consistency guarantee for operations A then B?
2. Why can "write then read" fail on a secondary without a session?
3. What does a session store for this purpose (high-level)?
4. Is causal consistency the same as a transaction?
5. Does every client in the cluster share one causal history?

#### Easy practical tasks

1. Open the causal-consistency page. Write one sentence definition and the URL.
2. In a driver, start a session, insert, find with the same session. Write the API names.
3. Write a "write then read on secondary" scenario in four steps. Mark where a session helps.
4. Make a table: tool, problem it solves: transaction, write concern, read preference, session causality.

#### Medium practical tasks

1. Try write on primary and immediate read on secondary without a session. Then with a session. Write the difference (or that you could not see lag).
2. Read whether your driver enables causal consistency by default on sessions. Write the default.
3. Explain to a teammate why two different app servers need a shared token or a primary read to see each other's writes.

#### Advanced practical tasks

1. Read about `afterClusterTime` and operation time. Write a short glossary of three terms.
2. Design a mobile API: which endpoints use a session for read-your-writes, and which endpoints accept stale reads.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do single-document atomicity, multi-document transactions, and write concern solve different problems?
2. When do read concern and read preference both matter for a "read after write" bug?
3. Why is a replica set a requirement for transactions but also the place where secondary stale reads appear?
4. How does causal consistency in a session differ from snapshot isolation in a transaction?
5. A teammate wraps every write in a transaction and reads from `nearest`. Which facts do you use in the review?

#### Easy practical tasks

1. Write a cheat sheet: single-doc atomic, transaction API, when not to use it, `w`/`readConcern`, read preference modes, session causality.
2. On a replica set, commit one transaction and one majority write. Find both documents.
3. Draw a timeline: client write, primary, secondary, a stale read, a majority read.
4. Mark three of your shop writes as: one document, transaction, or async repair.

#### Medium practical tasks

1. Implement a two-account transfer with a transaction and `w: "majority"`. Abort on insufficient funds. Write the tests.
2. Implement the same transfer as events in one ledger collection (single-document append). Compare code and failure modes.
3. Configure one API read as `primary` + `majority` and one report as `secondary` + `local`. Write why each pair fits.

#### Advanced practical tasks

1. Measure abort rate and duration of a contended transaction (two clients update the same two documents). Write the numbers and a model change that reduces contention.
2. Read the current consistency documentation for your major version. Write a one-page policy: defaults for write concern, read concern, read preference, and sessions.
