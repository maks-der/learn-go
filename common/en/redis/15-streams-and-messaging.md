# 15. Streams and Messaging

## Description

A Redis stream is an append-only log of entries. This topic covers entry ids, consumer groups, pending entries, when to choose Kafka instead, and why Pub/Sub can lose messages.

Complete topic 7 (stream intro) and topic 14 (Pub/Sub vs streams) first. Topic 7 showed `XADD`, `XREAD`, and a first group. This topic is the operations model.

Use one term for each concept. An entry is one item in the stream. An id is the entry address (`milliseconds-sequence`). A consumer group is a named cursor plus a pending list. A consumer is a name inside a group. The PEL is the pending entries list. `XACK` removes an entry from the PEL. At-most-once means a delivery can be lost. At-least-once means a delivery can happen twice if you do not make work idempotent.

---

## Stream entries and IDs

`XADD key id field value [field value ...]` appends an entry. The usual id is `*`. Redis creates an id from the time and a sequence number.

```text
XADD lab:orders * item book qty 1
```

The reply is an id such as `1726248000000-0`. The part before `-` is milliseconds. The part after `-` is the sequence in that millisecond.

You can pass an explicit id. The new id must be greater than the last id in the stream. Explicit ids are for imports and tests. Prefer `*` in new producers.

Special ids when you read:

- `0-0` or `-` — from the start (command-dependent)
- `+` — the end (in `XRANGE`)
- `$` — only new entries after the current end (`XREAD`)
- `>` — new entries for a group (`XREADGROUP`)

```text
XRANGE lab:orders - +
XLEN lab:orders
XREAD COUNT 10 STREAMS lab:orders 0-0
```

`XREAD BLOCK milliseconds STREAMS lab:orders $` waits for a new entry. The connection is blocked. `INFO clients` can show blocked clients (topic 19).

Each field is a string. Put JSON in a field if you need a document. Keep a `type` field for routing.

An unbounded stream fills RAM. Use `MAXLEN` on `XADD` or `XTRIM` to cap the length. Approximate `MAXLEN` (`~`) is faster. You can lose older entries. That is a policy, not a surprise.

### Questions

#### Theoretical questions

1. What does `XADD` with `*` set for the id?
2. What are the two parts of a stream id?
3. What does `$` mean in `XREAD`?
4. Why must an explicit id increase?
5. Why do you cap a stream with `MAXLEN` or `XTRIM`?

#### Easy practical tasks

1. `XADD lab:s * msg hello`. Save the id.
2. `XADD` a second entry. `XRANGE lab:s - +`. Save both ids.
3. `XLEN lab:s`. Save the reply.
4. Write four sentences: stream id vs a list index.

#### Medium practical tasks

1. `XREAD COUNT 1 STREAMS lab:s 0-0`. Then `XREAD` from the id you received. Confirm you see the next entry.
2. `XADD lab:s MAXLEN 5` in a loop of 12. Write the final `XLEN`.
3. Try `XADD` with an id lower than the last id. Record the error.

#### Advanced practical tasks

1. `XADD` with `MINID` or a trim by id if your version supports it. Write how you drop entries older than a time.
2. Design fields for an order event. Write which field is JSON and which fields are plain. Give example `XADD` lines.

---

## Consumer groups, `XREADGROUP`, `XACK`

A consumer group lets many workers share a stream. Each entry goes to one consumer in the group (competing consumers). Redis stores the last delivered id for the group and the PEL.

Create a group:

```text
XGROUP CREATE lab:orders grpA 0 MKSTREAM
```

`0` means "start from the beginning". `$` means "start from new entries only". `MKSTREAM` creates the stream if it is missing.

Read as a consumer:

```text
XREADGROUP GROUP grpA worker1 COUNT 1 STREAMS lab:orders >
```

`>` means "give me entries that the group has not delivered yet." The consumer name is `worker1`. Two workers use two names.

After success, acknowledge:

```text
XACK lab:orders grpA 1726248000000-0
```

`XACK` removes the entry from the PEL. The entry can remain in the stream until trim. `XACK` is not `XDEL`.

If the process dies before `XACK`, the entry stays pending. Another worker can claim it (next section).

`XINFO GROUPS lab:orders` and `XINFO CONSUMERS lab:orders grpA` show lag-style fields and pending counts.

`XGROUP DESTROY` and `XGROUP DELCONSUMER` are operational commands. Use them in a lab or in a runbook, not in a hot path.

One stream can have many groups. Each group has its own cursor. The same entry can be delivered to every group. That is fan-out across groups, and competing consumers inside a group.

### Questions

#### Theoretical questions

1. What does a consumer group add that `XREAD` does not store?
2. What does `>` mean in `XREADGROUP`?
3. What does `XACK` remove, and what does it not delete?
4. What does `XGROUP CREATE ... $` skip?
5. How do two groups on one stream share an entry?

#### Easy practical tasks

1. `XGROUP CREATE lab:g g1 0 MKSTREAM`. `XADD` two entries. `XREADGROUP` as `c1`. Save the ids.
2. `XACK` those ids. `XPENDING lab:g g1`. Save the reply.
3. Write four sentences: `XREAD` vs `XREADGROUP`.
4. Run `XINFO GROUPS lab:g`. Write the group name and pending count.

#### Medium practical tasks

1. Start two consumers `c1` and `c2` in one group. `XADD` four entries. Write which consumer received which id.
2. Read with `>` and do not `XACK`. `XREADGROUP` again with `>`. Write whether you see the same entry.
3. Create a second group on the same stream. Show that the second group can read the same entries from its own cursor.

#### Advanced practical tasks

1. Implement a worker loop: `XREADGROUP BLOCK`, process, `XACK`. Crash before `XACK` (kill). Document the PEL.
2. Write a runbook: create group, scale out a consumer name, remove a dead consumer (`XGROUP DELCONSUMER`).

---

## Pending entries (`XPENDING`, `XCLAIM`)

The PEL holds entries that a group delivered and that no consumer has acknowledged.

```text
XPENDING lab:orders grpA
XPENDING lab:orders grpA - + 10
XPENDING lab:orders grpA - + 10 worker1
```

The short form returns a count and the idle range. The long form lists ids, the consumer, idle time, and delivery count.

A dead consumer leaves pending entries. Another consumer claims them:

```text
XCLAIM lab:orders grpA worker2 60000 1726248000000-0
```

The number `60000` is a minimum idle time in milliseconds. Redis does not claim an entry that is younger than that idle time (unless you force). This avoids a steal from a slow but live worker.

`XAUTOCLAIM` (Redis 6.2+) walks the PEL and claims a batch. Prefer it in new workers.

After `XCLAIM`, the new consumer owns the PEL entries. Process them. Then `XACK`.

Delivery count increases when Redis delivers again. A poison message (always fails) can loop. After N deliveries, write a dead-letter key or log, then `XACK` so that the PEL does not grow forever.

`XPENDING` is the health check for a group. A growing PEL means consumers are slow, crashed, or failing to `XACK`.

Do not `XACK` before the side effect is durable, if you need at-least-once processing. Do not skip `XACK` if you already completed the side effect, or the PEL grows and another worker will do the work again. Make the side effect idempotent (topic 14).

### Questions

#### Theoretical questions

1. What does the PEL store?
2. What does the idle time in `XCLAIM` protect?
3. What command claims a batch without listing each id?
4. What is a poison message in this context?
5. Why must processing be idempotent if you claim and run again?

#### Easy practical tasks

1. `XREADGROUP` one entry. Do not `XACK`. `XPENDING` the group. Save the output.
2. Write four sentences: `XACK` vs `XCLAIM`.
3. Open the `XCLAIM` command page. Write the unit of the idle argument.
4. List three reasons a PEL can grow.

#### Medium practical tasks

1. From a second consumer, `XCLAIM` the pending id with a small idle time (or wait). Process and `XACK`. Confirm `XPENDING` is empty.
2. Use `XAUTOCLAIM` if available. Save the reply shape.
3. Increment a delivery count by claiming or reading the pending id again (as the docs allow). Write the count.

#### Advanced practical tasks

1. Build a dead-letter path: after 5 deliveries, `XADD` to `lab:orders:dlq` and `XACK` the original. Test with a worker that always fails.
2. Measure PEL size and oldest idle time. Write alert thresholds for a lab stream.

---

## When to use Redis Streams vs Kafka

Redis Streams and Apache Kafka are both logs. They are not the same product.

Redis Streams fit when:

- You already run Redis
- The log is small enough for RAM (or you trim hard)
- Retention is short
- Consumer groups are simple
- Operational team is a Redis team
- You need low latency on a single cluster you already know

Kafka fits when:

- Retention is large (hours to days of high volume on disk)
- Many independent consumer groups read a long backlog
- You need a large broker cluster, partitions, and an ecosystem (Connect, schema registry)
- Replay of a huge history is a normal operation
- Isolation from the cache Redis is mandatory (topic 21)

Redis keeps stream data in memory (plus persistence if you enable it). A multi-day firehose of events will fill Redis. Kafka writes a disk log and scales partitions across brokers.

Kafka consumer offsets live in Kafka. Redis group state lives in Redis. Failure modes follow the rest of that system (Redis failover vs Kafka replica ISR).

You can use both: Redis for cache and short work queues, Kafka for the company event bus.

Do not pick Streams only because the command names are shorter. Measure volume, retention, and who will operate the system.

### Questions

#### Theoretical questions

1. Where does a Redis stream primarily live?
2. Why can long retention favor Kafka?
3. When is a Redis stream a good fit?
4. Why might you isolate a stream Redis from a cache Redis?
5. What does "replay a huge history" mean for Kafka vs a trimmed Redis stream?

#### Easy practical tasks

1. Write a six-row table: property, Redis Streams, Kafka. Fill retention, storage, groups, ops team, typical volume, latency.
2. Write four sentences: when you would not put the company event bus on Redis.
3. Open Redis Streams docs and a Kafka intro page. Write one official sentence from each (short quote, then your words).
4. Estimate RAM for 10 million entries at 200 bytes each. Write the product.

#### Medium practical tasks

1. List your last project events. Mark each "Redis stream", "Kafka", or "Pub/Sub". Give one reason each.
2. Read Kafka partition vs Redis stream (one key) scaling. Write how you scale Redis (shard keys) vs Kafka (partitions).
3. Compare persistence: Redis AOF/RDB vs Kafka disk log. Write a loss-window sentence for each.

#### Advanced practical tasks

1. Write a decision page for your team: event types on Redis vs Kafka, max stream `MAXLEN`, and who owns pages.
2. Prototype the same three events on a Redis stream and (if you have it) a Kafka topic. Write operational steps to replay last 1000.

---

## Pub/Sub at-most-once nature

Topic 14 showed `PUBLISH` and `SUBSCRIBE`. This section fixes the delivery word.

At-most-once: Redis tries to send the message to current subscribers. If a subscriber is disconnected, slow, or not yet subscribed, that subscriber does not get a later replay from Redis. Redis does not store the message.

Slow subscribers: Redis can disconnect a client that does not read fast enough (`client-output-buffer-limit` for pubsub). That client loses messages.

`SUBSCRIBE` after `PUBLISH` does not receive the old message.

There is no `ACK`. The publisher does not know that a given subscriber processed the payload.

Use Pub/Sub for hints. Example: "product 1001 changed." The subscriber `DEL`s a cache key or reloads flags. If the hint is lost, a TTL still repairs the cache.

Do not use Pub/Sub for "charge this card" or "send this email once." Use a stream (or Kafka) and idempotent workers.

Cluster Pub/Sub still does not add a backlog. The cluster bus fans out to nodes. Delivery to a subscriber remains at-most-once.

### Questions

#### Theoretical questions

1. What does at-most-once mean for Pub/Sub?
2. What happens if you `SUBSCRIBE` after `PUBLISH`?
3. What can happen to a slow Pub/Sub client?
4. Why is there no `XACK` equivalent in Pub/Sub?
5. What backup repairs a lost cache-invalidation hint?

#### Easy practical tasks

1. `PUBLISH` then `SUBSCRIBE` in that order. Write what the subscriber receives.
2. Write four sentences: at-most-once Pub/Sub vs at-least-once streams.
3. List three event types that must not use Pub/Sub.
4. Open `client-output-buffer-limit` docs. Write the pubsub limit in your own words.

#### Medium practical tasks

1. Subscribe, then publish 1000 large messages as fast as you can. Watch whether the subscriber stays connected. Record `INFO clients`.
2. Implement invalidation: on hint loss, show that TTL still expires the cache key.
3. Compare a list `BRPOP` queue with Pub/Sub: who holds the message if no worker is up?

#### Advanced practical tasks

1. Write a delivery table: Pub/Sub, list queue, streams, Kafka. Rows: stored if no consumer, ACK, competing consumers, typical loss.
2. Design a hybrid: Pub/Sub hint plus a stream of changes. Write what each path is for and how you avoid double work.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do entry ids, `>`, `$`, and the PEL fit on one diagram of a group?
2. When is `XACK` before the side effect worse than `XACK` after, and the reverse?
3. Why does a growing PEL change your scaling decision (more consumers vs a dead-letter path)?
4. Which volume and retention numbers push you from Redis Streams to Kafka?
5. How do you explain at-most-once Pub/Sub to a teammate who wants a "simple queue"?

#### Easy practical tasks

1. Write a cheat sheet: `XADD`, `XRANGE`, `XGROUP CREATE`, `XREADGROUP`, `XACK`, `XPENDING`, `XCLAIM`, `PUBLISH`.
2. Draw stream, group, two consumers, PEL, and a `PUBLISH` channel on the side.
3. Write five lab rules: trim streams, unique consumer names, `XACK` policy, no Pub/Sub for jobs, idempotent claims.
4. Run `XINFO STREAM` on a lab stream. Write two fields that you understand.

#### Medium practical tasks

1. Build a two-consumer group on one stream. Kill one consumer. Claim and finish its PEL. Document commands.
2. Write a one-page worker spec: block time, count, claim idle, max deliveries, trim policy.
3. Compare `MAXLEN` loss with Kafka retention. Write when each loss is acceptable.

#### Advanced practical tasks

1. Run a stream plus consumer group with two processes in your language. Restart one process. Prove no lost `XACK`ed work and list duplicate risk.
2. Write an architecture note: cache Redis, stream Redis, Kafka bus. Place six example events.
