# 15. Change Streams and Integration

## Description

This topic shows how MongoDB pushes data changes to listeners. A **change stream** is an official cursor of change events. A **resume token** lets you continue after a disconnect. Teams use change streams for **CDC** (change data capture) into search or a warehouse. **Atlas Triggers** run code when data changes. Complete replication, drivers, and CRUD first.

Use one term for each concept. A **change event** is one document that describes an insert, update, replace, delete, or some cluster events. The **resume token** is the `_id` of an event. **CDC** is the practice of copying those changes into another system.

Change streams need a replica set or a sharded cluster. A standalone `mongod` is not enough.

---

## Change streams

A change stream is a long-lived cursor. The server sends a new event when a matching write commits.

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

By default, an update event may omit the full document. Open the stream with `fullDocument: "updateLookup"` when you need the document after the update. That lookup has a cost and can race. Read the current options (`whenAvailable`, `required` on newer versions).

The stream follows the oplog. If the listener is down longer than the oplog window, the stream cannot resume. You must do a full resync of the downstream system.

Filter in `watch` with an aggregation pipeline (`$match` on `operationType` or on fields). Do not pull every event if you need only inserts.

Do not treat a change stream as a queue with infinite retention. It is not Kafka. Retention is the oplog window unless you add another store.

Do not run a change-stream worker that does slow work inside the receive loop without back-pressure. You can fall off the oplog.

### Questions

#### Theoretical questions

1. What is a change stream?
2. Why does a standalone `mongod` reject `watch`?
3. What field is the resume token?
4. Why can an update event omit the full document?
5. How does the oplog window limit a listener that is down?

#### Easy practical tasks

1. Confirm you have a replica set. Open `watch` on a test collection. Insert one document. Print the event.
2. Write the `operationType` that you see.
3. Open the change-streams page. Write the URL.
4. Make a table: watch target, example use. Add collection, database, deployment.

#### Medium practical tasks

1. Open a stream with `fullDocument: "updateLookup"`. Update a field. Write whether the event includes the full document.
2. Add a pipeline `$match` for `operationType: "insert"` only. Prove that an update does not appear.
3. In one official driver, write the `watch` call and how you iterate events.

#### Advanced practical tasks

1. Read change-stream event types for your version (including `invalidate`). Write when the cursor closes.
2. Design a worker that writes events to a file. List failure modes: lag, process crash, `invalidate`.

---

## Resume tokens

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

On a sharded cluster, tokens still work, but you must use the driver APIs. Do not open a raw stream on one shard only.

After `invalidate` (drop, rename, drop database), many cursors die. You may need `startAfter` or a new full sync. Read the page for your version.

Atlas and drivers can help with resume. You still persist the token in your application store.

### Questions

#### Theoretical questions

1. What do you store to resume a stream?
2. Why must you persist the token after downstream success?
3. What happens if the token is older than the oplog?
4. Why must you not invent a token?
5. Why does `invalidate` need extra care?

#### Easy practical tasks

1. Print `_id` from one change event. Write that it is the resume token.
2. Open the resume page. Write the URL and the option names that you see.
3. Write four sentences: token, persist order, oplog, idempotent consumer.
4. Make a table: `resumeAfter`, `startAfter`, `startAtOperationTime`. Add one sentence each (from the manual).

#### Medium practical tasks

1. Receive one event. Close the cursor. Reopen with `resumeAfter`. Insert again. Confirm you do not see the first event twice (or write what you saw).
2. Store tokens in a small `cdc_state` collection. Write the document shape `{ listenerId, token, updatedAt }`.
3. Simulate a crash after downstream write but before token save. Write if the next start replays. Write how you make the write safe.

#### Advanced practical tasks

1. Read token format and version compatibility notes. Write one upgrade warning.
2. Design at-least-once CDC with a unique key on the sink. Write the retry and the skip rules.

---

## CDC into search or a warehouse

**CDC** copies database changes into another system. Common targets:

- A search engine or Atlas Search index that you maintain yourself
- A data warehouse or lake
- A cache or a read model
- Another microservice

Change streams are one CDC source. Atlas also offers other pipeline tools on some tiers. Kafka Connect and similar tools exist. Pick one pipeline. Do not run three independent listeners that write the same sink.

A typical flow:

1. Take an initial snapshot of the collection (or of a query).
2. Record a cluster time or a token.
3. Apply the snapshot to the sink.
4. Start the change stream from that time or token.
5. Apply events in order per document. Make each apply idempotent.

Search CDC: map each insert/update/replace to an upsert in the index. Map delete to a delete in the index. If you miss events, the index is wrong until you rebuild.

Warehouse CDC: many teams write events to files or to a queue, then load in batches. Do not issue a tiny warehouse statement per event if the warehouse is slow.

Ordering: events for one document are ordered. Events for two documents can interleave. The sink must tolerate that.

Do not use CDC to replace a transaction across MongoDB and the sink. The sink is asynchronous. The application that needs a search hit "now" must accept delay or must write to search in the request path (with a different consistency story).

Secure the listener. A change stream can read all fields. Use a role with only the rights that you need. Do not log full documents if they contain secrets.

### Questions

#### Theoretical questions

1. What does CDC copy?
2. Why do you take a snapshot before you start the stream?
3. How do you apply a delete event to a search index?
4. Why is the sink not in the same transaction as MongoDB?
5. Why can two documents arrive in an order that surprises a report?

#### Easy practical tasks

1. Draw snapshot → token → stream → search upsert/delete.
2. Open a MongoDB CDC or Kafka connector page (official or Atlas). Write the URL.
3. Write three sink types and one risk for each.
4. List fields that you must not send to a public search index.

#### Medium practical tasks

1. Write a fake sink (a JSON file). Apply insert, update, and delete from a change stream. Show the file after each event.
2. Design a rebuild: drop the index, snapshot, resume. Write the user-visible downtime.
3. Compare Atlas Search automatic sync (next topic) with a custom change-stream indexer. Write two differences.

#### Advanced practical tasks

1. Read about `updateLookup` races vs using `updateDescription` to patch a sink. Write which you choose for search.
2. Design warehouse load: events to object storage every 5 minutes. Write partition keys and exactly-once vs at-least-once.

---

## Triggers (Atlas)

**Atlas Triggers** run a function when an event occurs. A database trigger can fire on insert, update, replace, or delete in a linked Atlas cluster. Scheduled triggers fire on a timer. Authentication triggers fire on user events in Atlas App Services.

A database trigger receives a change event. The function can call an API, write to another collection, or call an Atlas service.

Triggers are useful for:

- Small fan-out (send a notification)
- Light denormalization
- A first CDC step into an Atlas service

Triggers are not useful for:

- Long CPU work
- Unbounded loops
- The only audit log if you need a hard guarantee without monitoring

Limits exist: execution time, compute, and event volume. Read the Atlas Triggers page for your project. A busy collection can exceed the trigger budget. Then you need a self-managed worker or a queue.

You configure the trigger in the Atlas UI or with infrastructure-as-code that Atlas supports. You select the collection, the operation types, and whether you need the full document.

Failures retry by Atlas rules. Make the function idempotent. A retry can run the side effect twice.

Do not put secrets in trigger source that you commit to a public repository. Use Atlas secret storage.

If you do not use Atlas, you do not have Atlas Triggers. Use a change-stream worker on your own host.

### Questions

#### Theoretical questions

1. What does a database trigger run?
2. Which operation types can fire a database trigger?
3. Why must a trigger function be idempotent?
4. When do you replace a trigger with your own worker?
5. What do you use if you are not on Atlas?

#### Easy practical tasks

1. Open the Atlas Triggers documentation. Write the URL and the trigger types.
2. Write three good trigger uses and three poor uses.
3. Draw: insert → trigger function → HTTP call.
4. Write where you store an API key for a trigger.

#### Medium practical tasks

1. If you have Atlas App Services, create a trigger on a test collection that writes `{ ok: 1 }` to a log collection. Insert a document. Confirm the log.
2. Read time and compute limits. Write two numbers from the page.
3. Compare a trigger with a change-stream process that you run in Docker. Write who restarts the process.

#### Advanced practical tasks

1. Design a trigger that calls an external API. Include retry, idempotency key, and a dead-letter collection.
2. Read whether triggers use change streams internally. Write how oplog retention still matters.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do the oplog, change streams, and resume tokens form one pipeline?
2. When do you choose a change-stream worker, Atlas Triggers, or a warehouse connector?
3. Why is "store the token first" the opposite of a safe at-least-once consumer?
4. How do `invalidate`, oplog window loss, and a missed delete each corrupt a search index in a different way?
5. A teammate tails `local.oplog.rs` in the application. Which facts do you use in the review?

#### Easy practical tasks

1. Write a cheat sheet: `watch` targets, event fields, `fullDocument`, token options, snapshot+CDC, trigger types.
2. Open one stream. Capture one insert event. Save the JSON (redact secrets).
3. Draw two architectures: trigger notification vs CDC to a warehouse.
4. Write a one-page glossary: change stream, resume token, CDC, trigger, invalidate.

#### Medium practical tasks

1. Build a small listener that prints `operationType` and `_id`. Restart it with a saved token. Prove continuity.
2. Write a CDC runbook: initial snapshot, token store, rebuild, who alerts on lag.
3. Map three integration targets (search, warehouse, email) to trigger, worker, or batch job. Give one reason each.

#### Advanced practical tasks

1. Implement idempotent CDC into a second MongoDB collection that mirrors `orders`. Handle insert, update, replace, delete. Resume after kill.
2. Write a production standard: max handler time, token storage, oplog size, secrets, and when Triggers are forbidden.
