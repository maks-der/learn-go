# 21. Anti-Patterns

## Description

This topic shows designs that fail as data grows. You learn why **unbounded arrays** hurt. You learn why documents near **16 MB** hurt. You learn why a **low-cardinality shard key** hurts. You learn why a **transaction on every write** hurts. You learn why MongoDB is **not a drop-in SQL replacement**. Complete modeling, sharding, transactions, and performance first.

Use one term for each concept. An **anti-pattern** is a design that looks convenient and then fails. **Unbounded** means the array or the document can grow without a planned limit. **Low cardinality** means few distinct values. A **drop-in replacement** is the idea that you can keep SQL habits and only change the driver.

If you already use an anti-pattern, you change the model. You do not add more hardware first.

---

## Unbounded arrays that grow forever

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

Do not rely on "we will archive later" without a job and a date.

`explain` and `collStats` do not always make this obvious until one document is huge. Query the max array length in aggregation (`$size`) on a sample.

### Questions

#### Theoretical questions

1. What is an unbounded array?
2. Why does a large array make each `$push` more expensive?
3. Why does a list API suffer if comments live in one huge array?
4. What does `$slice` on `$push` guarantee?
5. When is a second collection the better home for child items?

#### Easy practical tasks

1. Write three unbounded-array examples and three bounded-array examples.
2. Open the unbounded-array or data-modeling page that warns about growing documents. Write the URL.
3. Draw one post document that failed at 16 MB. Draw the split into `posts` and `comments`.
4. Make a table: pattern, cap, when to use. Add last-N, reference, bucket.

#### Medium practical tasks

1. `$push` in a loop until the document is large (stay safe in a lab). Time the last 10 pushes vs the first 10. Write the times.
2. Rewrite the same data as a comments collection. Write the two finds (post, last 20 comments).
3. Write an aggregation that finds the document with the largest array length in a lab collection.

#### Advanced practical tasks

1. Design a migration from an unbounded `events[]` to a time series or bucket collection. Include expand-contract.
2. Review an application model (yours or a sample). List every array. Mark bounded or unbounded. Fix one.

---

## Huge documents near 16 MB

The BSON **document size limit is 16 MB**. A document that is "only" 10 MB is already a problem. It uses cache, bandwidth, and lock time on that document.

Causes:

- Unbounded arrays
- Embedded files (use GridFS or object storage)
- Huge strings (HTML dumps, base64 images)
- Nested copies of the same related data

Effects:

- Insert or update fails with a size error
- Every find of that document moves megabytes
- Working set shrinks (few huge docs fill RAM)
- Replication and backups take longer

Prevention:

- Store binaries outside the document
- Reference large related graphs
- Use subset pattern (store a small preview, keep the rest elsewhere)
- Validate max length in the application

If you are near the limit, you are late. Split before you hit the error in production.

Do not store a PDF as a 14 MB BinData field because "it is under 16 MB".

Do not `find` the full document when the API needs three fields. Projection helps, but the server may still fetch the document.

`Object.bsonsize` in `mongosh` (or a driver size helper) measures a document. Watch `p99` size, not only the average.

### Questions

#### Theoretical questions

1. What is the BSON document size limit?
2. Why is a 10 MB document already a risk?
3. How do base64 images inflate a document?
4. Why does a huge document shrink the working set?
5. Why is "under 16 MB" a poor success test?

#### Easy practical tasks

1. Insert a document. Run `Object.bsonsize` on it in `mongosh`. Write the number.
2. Open the document size-limit page. Write the URL.
3. Make a table: content type, store in document?, store where instead.
4. Write four sentences: 16 MB, cache, projection, GridFS or S3.

#### Medium practical tasks

1. Build a 2 MB string field in a lab. Time find with and without projection. Write the times.
2. Design a product document with a 20-page description. Write what you embed and what you reference.
3. Write a validation or application check that rejects documents above 200 KB for a given collection.

#### Advanced practical tasks

1. Find the largest documents in a lab (`$bsonSize` in aggregation if your version has it). Write the `_id` and size.
2. Plan a split of a 12 MB document into a summary plus parts. Include a client read path.

---

## Low-cardinality shard keys

A **low-cardinality** shard key has few distinct values. Examples: `country` with three values, `status` with four values, `isActive` with two values.

Effects:

- Chunks cannot split past the distinct values
- Those chunks become **jumbo**
- All documents with `status: "open"` live on one shard
- That shard is **hot**
- The balancer cannot save you

Hashing a boolean does not create cardinality. You still have two hashes.

A compound key `{ status: 1, orderId: 1 }` can add cardinality **if** `orderId` is distinct. Refine can add a suffix. The prefix `status` still groups opens together. Writes of new "open" orders can still concentrate.

High cardinality is necessary and not sufficient. You also need even frequency and a key that queries use.

Do not shard on a field because "every query has it" if the field has five values.

Do not confuse unique `_id` (high cardinality) with a good **ranged** key if `_id` is monotonic. That is a different anti-pattern (hot last chunk). Low cardinality is the few-values problem.

If you already sharded on a poor key, use refine or reshard (topic 13). That is a project.

### Questions

#### Theoretical questions

1. What does low cardinality mean for a shard key?
2. Why can those chunks become jumbo?
3. Why does hashing `isActive` fail?
4. How can a compound suffix help, and what problem remains?
5. Why can the balancer not fix this key?

#### Easy practical tasks

1. Write five poor shard keys and one reason each.
2. Open the shard-key cardinality page. Write the URL.
3. Draw three shards and all `status: "open"` on shard 1.
4. Make a table: field, distinct count guess, shard key? Add five fields.

#### Medium practical tasks

1. On paper, estimate cardinality of `tenantId`, `country`, `orderId`, `createdAt` for an app that you know.
2. Propose a better key for a collection that someone sharded on `{ type: 1 }`.
3. Write when refine is enough vs when you must reshard.

#### Advanced practical tasks

1. Invent metrics that prove a low-cardinality key (chunk count stuck, jumbo flag, ops on one shard). Write the dashboard.
2. Plan a reshard from `{ country: 1 }` to `{ tenantId: 1, _id: 1 }`. Include queries that must change.

---

## Transactions for every write

A **multi-document transaction** has a cost. It holds a snapshot. It can abort. It has time and size limits. It is necessary for some money moves across documents.

The anti-pattern is to wrap **every** `insertOne` or `updateOne` in a transaction "for safety".

Problems:

- Extra latency
- More aborts under load
- Developers hide a bad model (two documents that belong in one document)
- People hold the transaction while they call HTTP or wait for a user
- `bulkWrite` inside a transaction is not a reason to transaction every small write

Single-document writes are already atomic. Use that.

Use a transaction when two documents must change together and you cannot model them as one document or as an idempotent repair.

Do not start a session and a transaction in a generic repository `save()` that every feature calls.

Do not copy a relational habit of "always BEGIN". MongoDB is not that runtime.

If you need a transaction often, first ask if the documents can be one document.

Write concern and unique indexes prevent many bugs without a transaction.

### Questions

#### Theoretical questions

1. What cost does a multi-document transaction add?
2. Why is a transaction unnecessary for two `$set` fields on one document?
3. Why is a generic `save()` that always opens a transaction a problem?
4. Why must you not wait for a user inside a transaction?
5. Which two tools often replace a "just in case" transaction?

#### Easy practical tasks

1. List five writes. Mark transaction or single-document.
2. Open the transactions limits page. Write the URL and one limit.
3. Make a table: need transaction, do not. Add four rows.
4. Write four sentences: atomic document, transaction, abort, model first.

#### Medium practical tasks

1. Time 1000 inserts with and without a transaction per insert. Write the two times.
2. Refactor a two-document update into one document. Write the new shape.
3. Review a sample service that uses transactions for logging. Write the removal plan.

#### Advanced practical tasks

1. Measure abort rate when two clients transaction-update the same two documents. Propose a model that reduces conflict.
2. Write a team rule: who may add a transaction, with a checklist of model alternatives first.

---

## Treating MongoDB as a drop-in SQL replacement

A **drop-in SQL replacement** is the idea that you keep tables, normalize every entity, join in every read, and only swap the database.

MongoDB does not run SQL as the primary language. `$lookup` is not a full SQL join planner for every ad-hoc report. There are no declarative foreign keys as the main integrity tool. Transactions exist but they are not the default write.

If you map each table to a collection one-for-one, you often get:

- Many round trips or many `$lookup` stages
- Slow reports that SQL did better
- A schema that never uses embed
- Shard keys that do not match access
- Frustration that "MongoDB is slow"

MongoDB fits when you design **documents for access**. Embed what you read together. Reference what you share and update independently. Use a warehouse or a relational database for heavy ad-hoc analytics if that is the workload.

You can migrate from SQL. You must **redesign**. You do not only migrate types.

Do not implement a SQL emulator on MongoDB for the whole application.

Do not forbid embed because "normalization is always right". Do not embed everything because "joins are forbidden". Both extremes are anti-patterns.

ORMs that hide documents behind tables can recreate this anti-pattern. Read what SQL they generate or what queries they emit.

### Questions

#### Theoretical questions

1. What is a drop-in SQL replacement in this section?
2. Why is a 1:1 table-to-collection map often slow?
3. Why is `$lookup` not a complete SQL report engine?
4. When does a relational database or a warehouse stay in the system?
5. Why must a SQL migration include a redesign?

#### Easy practical tasks

1. Write five SQL habits and the MongoDB habit that replaces each one.
2. Open the data-modeling introduction. Write the URL and the access-together rule in your own words.
3. Make a table: embed, reference, SQL join. Add one example each.
4. Draw a 4-table shop schema and a 2-collection document design.

#### Medium practical tasks

1. Take a 5-table relational schema. Write a document design. List joins that disappear.
2. Compare one report in SQL vs an aggregation. Write which system you would use in production.
3. Review an ODM model that looks like tables. Write three queries that will hurt.

#### Advanced practical tasks

1. Write a one-page "we are not a SQL database" guide for teammates who know only PostgreSQL.
2. Plan a coexistence: MongoDB for the app, warehouse for BI. Write the CDC path (topic 15) in five steps.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do unbounded arrays and huge documents become the same incident at 16 MB?
2. Why do a low-cardinality shard key and a transaction-on-every-write both look like "the database is slow"?
3. How does a 1:1 SQL copy create unbounded arrays or huge documents if someone embeds "to avoid joins"?
4. When do bounded arrays, references, and a warehouse each replace a different anti-pattern?
5. A teammate shards on `status`, wraps every insert in a transaction, and `$push`es every click onto the user. Which facts do you use in the review?

#### Easy practical tasks

1. Write a cheat sheet: array cap, 16 MB, cardinality, transactions last, access-together design.
2. Audit one of your collections for the five anti-patterns. Write a yes/no row for each.
3. Draw a "stop" sign list of five commands or patterns (unbounded `$push`, `find(req.body)` is topic 17, `transaction` in `save()`, shard on boolean, 14 MB BinData).
4. Open the five official pages that you used in this topic. Write the five URLs.

#### Medium practical tasks

1. Pick one anti-pattern in a lab. Implement the bad version and the better version. Write sizes or times.
2. Write a code-review checklist of ten items from this topic.
3. Redesign one SQL-shaped schema that you wrote earlier in this path. List which anti-patterns you removed.

#### Advanced practical tasks

1. Produce a model review for a sample app: findings, severity, rewrite plan, and a test that fails if an array has no cap.
2. Write a team standard that forbids the five anti-patterns with a named exception process (who approves a transaction-heavy path or a shard key).
