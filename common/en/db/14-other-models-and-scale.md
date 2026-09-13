# 14. Other Models and Scale

## Description

This topic is a survey of data models that are not the relational table model, and of how you split work across stores and nodes. You learn key-value, document, wide-column, graph, and search stores, polyglot persistence, CAP in practical words, sharding, eventual consistency, sagas, outbox, idempotent writes, and OLTP versus OLAP.

Use one term for each concept. A data model is the way a product organizes values and queries. A system of record is the store that wins when copies disagree. A shard is a subset of the data on one node. Complete this topic after you can model tables and keys. You need the relational model as a comparison.

This path is vendor-neutral. Product names are examples of a category. Do not pick a product only by popularity.

---

## Key-value, document, wide-column, graph, and search stores

A **key-value** store maps a key to a value. The primary operation is get and set by key. The store does not parse the value as a first-class query model.

```text
SET session:42 --> { user_id: 7, expiry: ... }
GET session:42
DEL session:42
```

Use a key-value store when you always know the key, you need low latency, and the value is a cache, a session, a lock, or a flag. Do not use it as the only store for orders with many query shapes. You would scan keys or you would build secondary indexes by hand. Keys need a naming scheme. A collision overwrites a value. Durability differs. Some configurations keep data only in memory. Read the persistence setting before you store something that you cannot rebuild.

A **document** store keeps a document, often JSON, as the unit of data. One document can nest arrays and objects. Documents in one collection can differ in fields (flexible schema).

```text
{
  "order_id": 1001,
  "customer": { "name": "Ada", "city": "Kyiv" },
  "lines": [
    { "sku": "A-10", "qty": 2 }
  ]
}
```

Use a document store when the unit of read and write is a whole document, the shape changes often, and you load one aggregate by id. Integrity is often in the application. Duplication is common. You copy a customer name into the order document. An update of the customer name does not change old orders unless you write that job. Do not store unbounded growing arrays in one document. Documents have size limits.

A **wide-column** store organizes data by a row key and a flexible set of columns, often grouped in column families. Queries that match one row key are efficient. Queries that scan many keys can be expensive.

```text
Row key: user#42
  family profile: name=Ada, city=Kyiv
  family events:  2026-09-13T10:00=login
```

Use a wide-column store when you write a large volume of events or metrics and you read by a known row key or a key range. The primary-key design is the schema. If you need a different access path later, you add another table and you write twice. Do not pick this model for a small shop with many join-heavy reports.

A **graph** database stores nodes and edges. A query walks paths: neighbors, shortest path, or a pattern of types.

```text
(Ada:Person)-[:WORKS_AT]->(Acme:Company)
(Ada:Person)-[:KNOWS]->(Grace:Person)
```

Use a graph database when the question is the relationship: friends of friends, who can reach this resource, recommendation along a path. A relational database can store edges as a join table. Deep paths become recursive SQL. Do not model a simple one-to-many invoice as a graph only to follow a trend. Store money in a system of record. A supernode (a node with millions of edges) is a hot spot.

A **search** engine indexes tokens and fields so that you can rank text. Users type words. The engine returns documents by relevance, filters, and facets.

Use a search engine when users search text, when you need fuzzy match, or when you need relevance. Do not use a search index as the only copy of an order. Search indexes are often eventually consistent with the system of record. Writes go to the system of record first. Then you update the index. Access control must apply to hits. A search that ignores row grants can leak documents.

```text
Relational   : tables, rows, SQL
Document     : JSON documents
Key-value    : get / set by key
Wide-column  : row key plus columns
Graph        : nodes and edges
Search       : inverted index, ranked hits
```

### Questions

#### Theoretical questions

1. What is the primary access path in a key-value store, and when is that store a good fit?
2. What is the unit of data in a document store, and why does a copied customer name go stale?
3. Which queries are efficient in a wide-column store, and what does query-first key design mean?
4. What are the two basic units in a graph database, and why is a graph a poor system of record for money?
5. Why is a search index a poor only copy of an order, and in which order do you write the record store and the index?

#### Easy practical tasks

1. Design keys for: session, feature flag, page cache. Write a JSON document for a blog post with tags and an author object.
2. Design a row key for "metrics for host H on day D." Draw five nodes and six edges for a classroom.
3. Write three user queries that need search, not `LIKE '%x%'`. Draw write path: app, DBMS, index.
4. Name one product in each category from public knowledge. Mark them as examples only.

#### Medium practical tasks

1. Design an order document versus 3NF tables. Write which queries are easier in each model. Store a JSON blob under a key. Show that you cannot query a field without reading the blob.
2. Design two tables (two key layouts) for inbox-by-user and message-by-id. Write the double-write rule. Express "friends of friends who are not already friends" as a graph pattern and as SQL.
3. Compare `LIKE` and a token index for a 1 million row article table. Write four differences. Design a tenant filter on every search query. Write the field and the risk if you omit it.

#### Advanced practical tasks

1. Write a one-page design: sessions in a key-value store, orders in a relational database, catalog search in an index. Include expiry, rebuild, and retry.
2. Write a one-page model for a product catalog in documents plus a graph for recommendations. Include the sync path and a supernode risk.

---

## Polyglot persistence

Polyglot persistence means one application uses more than one data model on purpose. Example: relational orders, key-value sessions, search for the catalog, a graph for recommendations.

Each store has a job. One store is the system of record for each fact. Other stores are derived.

```text
Order placed --> relational DBMS (record)
             --> search index (catalog availability text)
             --> cache invalidate (session cart)
```

Benefits: each query uses a model that fits. Costs: more operations, more failure modes, more consistency work. You must define the write order and the repair job.

Rules:

1. Name the system of record for each fact.
2. Do not write a price in three stores as three equals without a sync rule.
3. Start with one store. Add a second store when a measured need exists (search latency, session volume).
4. Keep transactions in the record store. Use outbox or a retry job for derived stores (later section).

Do not add a store because a blog post praised it. Do not run five products for a class project.

A single relational DBMS can cover a large part of a small system. Polyglot is an architecture choice under load or under a special query, not a badge.

### Questions

#### Theoretical questions

1. What does polyglot persistence mean?
2. What is a system of record in this section?
3. Name two costs of extra stores.
4. When do you add a second store?
5. Where do transactions live in a polyglot design?

#### Easy practical tasks

1. Assign stores to: invoice, session, full-text help pages.
2. Draw one user action that touches two stores.
3. Write one fact that must have a single record store.
4. List three operational tasks that double when you add a DBMS product.

#### Medium practical tasks

1. Design a small shop with two stores only. Justify why a third store is not needed yet.
2. Write a sync failure: record committed, search not updated. Write the user-visible symptom and a repair.
3. Compare "JSON in a relational DBMS" versus a document DBMS for one collection. Write when one product is enough.

#### Advanced practical tasks

1. Write a one-page polyglot map: facts, record store, derived stores, write order, repair.
2. Review a public architecture that uses three stores. Mark each as record or derived. Note any unclear owner.

---

## CAP in practical words

CAP is a way to talk about a partition (a network split) in a distributed store.

- **C (consistency)** here means all clients see the same latest write (linearizability, in informal words: one up-to-date copy).
- **A (availability)** means every request to a non-failing node receives a response (not an error that says "I cannot decide").
- **P (partition tolerance)** means the system continues while the network drops messages between nodes.

On a single node, you do not choose CAP. When the network between nodes fails, a product must choose: refuse some writes or reads (keep one answer), or keep serving and accept that copies diverge.

```text
Partition: site A cannot talk to site B
  Option 1: A stops writes so B cannot miss them  (prefer one copy)
  Option 2: A and B both accept writes            (prefer availability, diverge)
```

Practical words:

- A relational primary that refuses writes when it cannot reach a sync replica prefers one copy.
- A multi-primary that accepts writes on both sides of a split prefers availability. You merge later.
- "We are AP" or "we are CP" as a slogan hides the real setting (quorum size, retry, client failover).

CAP does not say that you drop durability or that SQL is obsolete. It does not apply to a single SQLite file on one disk in the usual way.

When you read a vendor page, ask: what happens when two sites cannot talk? Who accepts writes? How do clients find a node? What do you repair after the split?

Do not use CAP to end a design talk. Use it to start the partition story.

### Questions

#### Theoretical questions

1. What does partition mean in this section?
2. What do you give up if both sides accept writes during a split?
3. What do you give up if a primary refuses writes during a split?
4. Why is "we are AP" a weak requirement?
5. Why does CAP not describe a single-node embedded file in the usual way?

#### Easy practical tasks

1. Write C, A, and P in one short sentence each (practical words).
2. Label option 1 and option 2 from this section on a two-site drawing.
3. Write one product setting that is closer to "one copy" and one that is closer to "keep serving."
4. Write one sentence that you will not use as a slogan.

#### Medium practical tasks

1. Take a primary plus async replica. Write what clients see if the primary is in the minority partition.
2. Take a quorum of three nodes. Write how a majority still accepts a write.
3. Read a CAP paragraph in a vendor doc. Rewrite it as a partition story in six sentences.

#### Advanced practical tasks

1. Write a one-page partition table-top for a shop: two sites, who writes, what users see, how you merge.
2. Compare two stores (relational HA versus a distributed key-value) on one partition. Use manuals. No slogans.

---

## Sharding and eventual consistency

Sharding (horizontal partitioning) splits the rows of a dataset by a shard key. Each shard holds a disjoint subset of keys. A router or the application sends a request to the shard that owns the key.

```text
shard key = customer_id
customers 1..999    --> shard A
customers 1000..1999 --> shard B
```

Use sharding when one primary cannot hold the data or the write rate. Sharding is an operational and design cost. Do not shard a small database.

A good shard key appears in almost every query, spreads writes (no single hot key), and keeps related rows that you join on the same shard when you can.

A poor shard key creates a hot shard (all new orders go to "today" if you shard only on date) or forces cross-shard joins. Cross-shard queries must fan out. They are slower and harder to isolate. Transactions that touch two shards are distributed transactions or they are application sagas.

Do not shard "for high availability." Replication gives copies. Sharding gives capacity. You often need both.

Range shards are simple and can unbalance. Hash shards spread better and make range scans harder. Resharding (move keys to new nodes) is a planned project.

Eventual consistency means copies can differ for a time. If writes stop, the copies converge to the same values (under the product's rules). Readers can see stale data during the delay.

```text
Write on node A at t=0
Read on node B at t=50ms  -->  old value
Read on node B at t=2s    -->  new value (if replication caught up)
```

This model appears in asynchronous replicas, search indexes after a write, caches with TTL, and saga steps that have not finished.

Use eventual consistency when the business can accept a short stale read. Example: a view count, a search result, a recommendation. Do not use it as the only story for a bank balance unless you have a precise conflict rule and a ledger.

Give users a truth they can understand: "Your change is saved. Search can take a minute." A silent stale read looks like a lost write.

**Read your writes** is a stronger need. After a save, read the primary or wait for a version.

Do not say "eventual consistency" to mean "we do not know." Write the delay budget and the repair if convergence fails (a stuck replica).

### Questions

#### Theoretical questions

1. What does a shard key decide, and when do you shard?
2. What makes a shard key poor, and why are cross-shard joins expensive?
3. Why is sharding not a substitute for replication?
4. What does eventual consistency promise if writes stop?
5. Why can a stale read look like a lost write, and what does read-your-writes require after a save?

#### Easy practical tasks

1. Pick a shard key for `orders` in a shop. Write one reason. Draw two shards and a router.
2. Write one query that stays on one shard and one that fans out. Write one hot-key example (a celebrity user id).
3. Write two features that can be eventually consistent and two that cannot for a shop. Write a user-facing sentence for a delayed search index.
4. Draw a stale read on a replica. Define a delay budget for cart count versus for payment status.

#### Medium practical tasks

1. Design range shards versus hash shards for `customer_id`. Write a range report that hash shards make hard. Estimate: 10 million customers, 4 shards. Show a skew if ids are sequential and you use range shards badly.
2. Demo a write on a primary and a read on an async replica (or simulate with a delay). Record the stale window. Design a "saved" page that reads the primary while the list page reads a replica.
3. Write a monitor: if replica lag exceeds the budget, stop sending reads to the replica. Write a reshard step list at a high level: add node, move keys, update router, verify counts.

#### Advanced practical tasks

1. Write a one-page shard plan for a multi-tenant app (tenant_id as key). Include a large tenant that needs its own shard, and a consistency menu: strong local txn, read-your-writes, eventual search.
2. Table-top a stuck replica that never converges. Write detection, user impact, and repair. Table-top a cross-shard order: customer on A, inventory on B. Write why a single local transaction is not enough.

---

## Sagas, outbox, and idempotent writes

A saga is a sequence of local transactions. Each step has a compensating action that undoes the business effect if a later step fails.

```text
1. Reserve inventory   (compensate: release)
2. Charge payment      (compensate: refund)
3. Create shipment     (compensate: cancel shipment)
```

The saga is not ACID isolation across steps. Other readers can see a reserved item before the charge finishes. You design status fields (`reserved`, `paid`, `failed`).

An **orchestration** saga has a coordinator process that calls each step. A **choreography** saga lets each service react to events. Orchestration is easier to see. Choreography avoids one boss process and is harder to debug.

The **outbox** pattern writes the business row and an "event to send" row in the same local transaction. A publisher reads the outbox and sends the event to a bus or to another store. Then it marks the outbox row as sent.

```text
BEGIN
  UPDATE orders SET status = 'paid'
  INSERT INTO outbox (event) VALUES ('OrderPaid')
COMMIT
-- publisher sends OrderPaid; at-least-once
```

Without an outbox, you can commit the order and crash before you send the event. The other service never hears. Dual-write without a transaction is the common bug.

An idempotent write can run more than once and leave the same final state as one successful run. Networks retry. Sagas retry. Outbox consumers retry. Without idempotency, a retry charges a card twice.

```text
POST /payments { idempotency_key: "ord-1001-pay" }
-- first call: charges
-- retry: returns the same payment id, no second charge
```

Methods:

1. A client idempotency key stored with the result.
2. A natural unique key (`order_id` unique on `payments`).
3. A state machine: `pending` to `paid` only once.
4. Upsert with a version number.

At-least-once delivery (outbox, message bus) requires idempotent consumers. Exactly-once is a combination of at-least-once plus idempotency.

`UPDATE balances SET n = n + 10` is not idempotent. Two retries add 20. Compensations must be safe to retry. A refund that runs twice is a money bug unless the refund is idempotent.

Do not use a saga to hide a missing local constraint. Do not rely on "the user will not double-click." The timeout retry is enough to double-apply.

A distributed transaction (two-phase commit) across products is expensive: extra round trips, locks held until the coordinator finishes, a blocking period if the coordinator crashes. Many teams avoid 2PC across products. They use a single local transaction in one system of record, then they propagate with saga and outbox.

### Questions

#### Theoretical questions

1. What does a compensating action do, and why is a saga not one ACID transaction?
2. What two writes does an outbox put in one local transaction, and what dual-write bug does the outbox prevent?
3. What does an idempotent write guarantee on retry, and why do outbox consumers need idempotency?
4. Why is `n = n + 10` not idempotent, and what is exactly-once in practical terms in this section?
5. What is the difference between orchestration and choreography?

#### Easy practical tasks

1. Write a three-step saga for a bookstore order. Name each compensate. Draw order table plus outbox table plus publisher.
2. Add an `idempotency_key` column idea to a `payments` table. Write the unique constraint.
3. Write a retry story for "request timed out, user does not know if the charge ran." Label three statements as idempotent or not: `INSERT` without key, `DELETE FROM t WHERE id = 1`, `n = n + 1`.
4. Write one status value that other services can see in the middle of a saga. Label a flow as orchestration or choreography.

#### Medium practical tasks

1. Implement an outbox table in a practice schema. Insert an order and an event in one transaction. Write a publisher loop in pseudo-code.
2. Implement a unique payment per `order_id` in SQL. Run the insert twice. Record the second result. Write a consumer that stores `event_id` as processed. Skip a duplicate.
3. Compare "send HTTP in the same request after COMMIT" with outbox. Write two failure cases. Rewrite an increment as a versioned set. Show a double retry that does not double-add.

#### Advanced practical tasks

1. Write a one-page saga for checkout with timeouts and a manual repair queue. Design outbox cleanup: sent rows, retention, and poison events. Include idempotent consumers.
2. Write a one-page idempotency standard for HTTP writes and for queue consumers. Include key lifetime. Design a refund that is idempotent and a compensate that uses it in a saga.

---

## OLTP vs OLAP, warehouses, ETL/ELT, and CDC

OLTP means online transaction processing. The workload is many small reads and writes. Each transaction touches a few rows. The user waits. Examples: place an order, check in a guest, post a payment.

OLAP means online analytical processing. The workload is fewer queries that read many rows. The query aggregates. The user is an analyst or a report. Examples: sales by month, conversion funnel, inventory aging.

```text
OLTP:  UPDATE one order, COMMIT, 10 ms
OLAP:  SUM sales for two years, 30 s, scan many pages
```

Schema shape differs. OLTP favors normalized tables and current state. OLAP favors wide fact tables, dimensions, and history that does not change (or changes in a controlled way).

If you run heavy OLAP on the OLTP primary, you pollute the buffer cache, you hold locks or you consume I/O, and you slow checkout. Use a replica, a warehouse, or a night batch.

An operational database (system of record) holds current business state. It accepts writes from applications. Constraints and transactions protect that state.

A data warehouse holds a copy of data for analysis. It is not the place where a customer completes checkout. Loads arrive in batches or as a stream. Users run reports and models.

```text
Apps --> operational DBMS --> (ETL/ELT/CDC) --> warehouse --> reports
```

A replica of the operational database is not a warehouse. A replica has the same schema and the same row layout. A warehouse transforms, history-tracks, and serves different tools.

Do not grant analysts `SELECT` on production operational tables as the only architecture. You leak personal data, you load the primary, and you cannot rebuild history. One direction of truth: operational is the record. The warehouse is derived.

ETL means extract, transform, load. You extract from the source, you transform in a pipeline, you load into the target.

ELT means extract, load, transform. You extract, you load raw data into the target, you transform inside the target with SQL or a warehouse engine.

```text
ETL:  source --> transform (job) --> warehouse
ELT:  source --> raw load --> transform (SQL in warehouse)
```

ETL fits when you must reduce data before it leaves the source or hide columns before they land. ELT fits when the warehouse can process large SQL and you want to keep a raw copy.

Both need a schedule or a trigger, idempotent loads or a watermark (load from timestamp T), a failure retry, and a data contract. Do not transform in ad-hoc notebooks as the only production path. Do not run a full extract every hour if an incremental watermark exists. Quality checks belong in the pipeline: row counts, null rates, unique keys. Time zones and money types must survive the pipeline.

CDC (change data capture) reads changes from the operational database (usually the WAL or a change stream) and publishes those changes to consumers: warehouse, search, cache, another service.

```text
Primary WAL --> CDC reader --> events --> consumers
```

CDC benefits: near-real-time derived stores, less load than full extracts, a natural event log. CDC costs: ordering, schema change, and the need for idempotent consumers. CDC is not a backup. CDC is not a substitute for a tested dump.

Do not design an OLTP table for a 50-column report first. Do not design a star schema for a 5 ms checkout.

### Questions

#### Theoretical questions

1. What is a typical OLTP transaction, and what is a typical OLAP query?
2. Why does a large aggregate harm an OLTP primary, and why is a replica not a warehouse?
3. What is the order of steps in ETL versus ELT, and when do you hide columns in ETL before the load?
4. What is a watermark, and why must a production transform not live only in a notebook?
5. Which system is the system of record for an order, and why must analysts not use production OLTP as the only report store?

#### Easy practical tasks

1. Label six statements as OLTP or OLAP: checkout, monthly tax report, password login, cohort analysis, add-to-cart, year-over-year sales.
2. Draw apps, operational DBMS, pipeline, warehouse, dashboard. Draw ETL and ELT side by side.
3. Write a watermark column for `orders` (`updated_at` or a log sequence). List three quality checks for a nightly load.
4. Write a freshness number for a warehouse that you would accept for a shop (example: 1 hour). Give one reason. Write one sentence you would use to refuse a 2-minute report on the checkout database.

#### Medium practical tasks

1. Design a normalized checkout model and a fact table for "sales by day and sku." List columns of the fact table. List five operational tables. Write which become facts and which become dimensions.
2. Write a pseudo-job: extract new orders since last watermark, load, update watermark, on failure do not move the watermark. Compare full extract versus incremental for a 10 million row table.
3. Decide: replica versus warehouse versus CDC for a daily 20-minute report and for a search index. Write four reasons. Write a data contract of six fields for `orders` in the warehouse.

#### Advanced practical tasks

1. Write a one-page workload split for a shop: which queries stay on OLTP, which move, the freshness need, ETL versus ELT choice, and CDC consumers.
2. Design a slowly changing customer dimension (name change) plus a late-arriving fact. Write how history stays in the warehouse and not in checkout. Write how CDC plus idempotent consumers keep search in sync.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do you choose key-value, document, wide-column, graph, or search from the access path, not from the product name?
2. When does polyglot persistence need a system of record, and how does CAP describe a split between two of those stores?
3. How do a shard key, eventual consistency, and a local transaction work together to avoid a distributed commit?
4. When do you choose a saga and an outbox instead of a two-phase commit, and how do idempotent consumers protect retries?
5. A teammate wants one document store for invoices, search, sessions, and friend graphs, shards on `created_at`, and a 2-minute report on the checkout primary. Which facts do you use in the reply?

#### Easy practical tasks

1. Write a cheat sheet: five models, polyglot, CAP partition story, shard, eventual, saga, outbox, idempotent, OLTP/OLAP, warehouse, ETL/ELT, CDC.
2. Fill a table: model, unit of data, typical query, poor query. Assign one example product per model. Mark all as examples.
3. Pick a shard key and an idempotency key for `PlaceOrder`. Pick a system of record for a library loan. Name one derived store if you need search.
4. Draw outbox plus a retrying consumer with a processed-event table. Write one user sentence for a stale derived store.

#### Medium practical tasks

1. Split a social shop across two models. Write write order, one repair, and a partition story. State who accepts checkout writes.
2. Design a two-shard shop that still charges once (idempotency) and updates search later (outbox). Write the tables. Rebuild one aggregate (order) as tables and as a document.
3. Write a pipeline: operational `orders` to a warehouse fact table with a watermark, plus CDC to search. Write freshness and a quality check.

#### Advanced practical tasks

1. Write an end-to-end design: shard map, one-shard checkout transaction, outbox, saga compensate, idempotent pay, warehouse load, no slogans. Review it against official docs.
2. Write an architecture note: three stores, record versus derived, operations cost, a CAP table-top, and which queries stay on OLTP versus the warehouse.
