# 10. Locks, Limits, and Messaging

## Description

Applications use Redis as a shared place for locks, limits, request ids, and sessions. Redis also delivers messages with Pub/Sub and streams. This topic covers those patterns, consumer groups, and when to choose Kafka.

Complete topic 2 (string locks), topic 5 (stream intro), topic 6 (Lua), and topic 9 (clients) first. Treat every pattern as application design plus Redis commands, not as a consensus product.

Use one term for each concept. A lock is a key that only one worker may hold. A rate limit is a cap on events in a time window. An idempotency key is a stored request id that makes a second apply a no-op. A session is per-user state with a TTL. Pub/Sub is fire-and-forget broadcast. A stream is a log with ids. A consumer group is a named cursor plus a pending list. The PEL is the pending entries list. `XACK` removes an entry from the PEL.

---

## Distributed lock (`SET NX EX`) and the Redlock debate

Topic 2 introduced a simple lock:

```text
SET lock:resource holder-id NX EX 30
```

`OK` means this worker holds the lock. A null reply means another worker holds it. The TTL frees the key if the worker dies.

Release must delete only your lock. A Lua script (topic 6) does `GET` then `DEL` if the value equals `holder-id`. A two-step `GET` plus `DEL` from the client can remove another worker's lock after expiry.

Limits of one-instance locks:

- Redis pause, eviction, or a long GC in the worker can expire the lock while the worker still works.
- A second worker then acquires the lock. Two workers run the critical section.
- One Redis instance is not a multi-node consensus system. A failover can lose a lock that was not persisted.

Redlock is an algorithm that uses several independent Redis instances. The Redis documentation described it. Martin Kleppmann and others published critiques. The debate covers clocks, pauses, fencing, and whether Redlock is safe for correctness-critical work.

This handbook does not settle the debate. For learning:

- Use one instance, a unique holder id, a TTL, and atomic release.
- For correctness-critical work (money, unique side effects), add a fencing token (a monotonic number from `INCR`) that the resource checks, or use a system that is designed for consensus.
- Do not treat `SET NX EX` as a complete distributed database transaction.

Do not lock when one Redis command already makes the update safe (`INCR`, a Lua update of Redis-only data).

### Questions

#### Theoretical questions

1. Which `SET` options acquire a simple lock?
2. Why does release need Lua (or an equivalent atomic check)?
3. What can cause two workers to hold the logical lock over time?
4. What is Redlock in one sentence?
5. What is a fencing token for?

#### Easy practical tasks

1. Acquire `lock:lab` with `SET NX EX 20` and a holder id. Save the reply.
2. Try a second acquire with a different holder. Save the reply.
3. Write four sentences: lock vs `INCR` for a counter.
4. Open a Redlock page and a critique page. Write one claim from each. Mark "doc" or "critique".

#### Medium practical tasks

1. Run the lock-release Lua from topic 6. Test match and mismatch. Save replies.
2. Let the TTL expire while a "worker" sleeps. Acquire from a second client. Write the overlap risk in five sentences.
3. Add `INCR lock:resource:fence` on acquire. Write how a storage API would reject a lower token.

#### Advanced practical tasks

1. Read Kleppmann's Redlock essay (or an equivalent critique) and the Redis lock page. Write a six-row table: issue, lock side, critique side.
2. Implement acquire, fence, and Lua release in your language. List what the design still does not guarantee.

---

## Rate limiting

A rate limit rejects or delays events when a client exceeds a budget.

Fixed window:

```text
INCR rl:user:42:202609131857
EXPIRE rl:user:42:202609131857 60
```

Use the current minute (or second) in the key. If the count is above N, reject. `EXPIRE` on first increment (Lua can set expiry only when the value is `1`).

The limit resets at the window edge. A user can send N events at the end of one minute and N at the start of the next (2N in two seconds).

Sliding window: store event timestamps in a sorted set. `ZREMRANGEBYSCORE` drops old members. `ZCARD` is the count. `ZADD` the new event if under the limit. Topic 4 covered the commands. Memory grows with events in the window.

Token bucket: a key stores tokens (or last refill time plus tokens). Each request takes one token. Tokens refill over time up to a cap. A Lua script updates the key atomically. Bursts are allowed up to the bucket size.

Choose:

- Fixed window — simple, cheap, burst at boundaries
- Sliding window — smoother, more CPU and memory
- Token bucket — burst control with a steady refill

Cluster: keep all keys for one limit on one slot (`{user:42}`).

Do not rate-limit only in the application process. Many processes share one Redis budget.

### Questions

#### Theoretical questions

1. What two commands implement a simple fixed window?
2. What extra burst can a fixed window allow at the boundary?
3. Which Redis type often implements a sliding window of events?
4. What does a token bucket refill?
5. Why must many application processes share one Redis limit?

#### Easy practical tasks

1. Build `rl:lab` with `INCR` and `EXPIRE 10`. Hit it five times. Save the count and `TTL`.
2. Write four sentences: fixed window vs sliding window.
3. Draw a token bucket: capacity 10, refill 1 per second, one request.
4. Write key names for "100 requests per user per minute".

#### Medium practical tasks

1. Write Lua that `INCR`s and `EXPIRE`s only when the value is `1`. Test two increments. Confirm one TTL.
2. Implement a sliding window with a sorted set and a 10 second window. Add 8 members. Confirm `ZCARD`.
3. Compare memory: 10,000 users with a fixed-window integer vs a sorted set of 100 timestamps. Estimate orders of magnitude.

#### Advanced practical tasks

1. Implement one token-bucket Lua script. Test empty bucket and refill after a sleep. Save replies.
2. Design Cluster-safe keys for user, IP, and global limits. Write hash tags and a denial order (which limit you check first).

---

## Idempotency keys and session store

A client can send the same request twice (retry, double click, timeout). An idempotency key makes the second apply return the first result (or a reserved error) without a second side effect.

Pattern:

```text
SET idem:pay:req-9f3a 1 NX EX 86400
```

If the reply is `OK`, this is the first request. Do the work. Then `SET` the stored result (or a hash) on the same key or on a sibling key.

If the reply is null, the key exists. Return the stored result. Do not charge again.

Use a TTL long enough for retries (hours for payments is common). Document the TTL.

The work itself must still be safe if the process dies after `SET NX` and before the side effect completes. Some designs store `pending`, then `done` plus the result. A reconciler finishes `pending` keys.

Redis is a coordinator. The bank or SQL write still needs its own idempotency when you can use it.

Do not use a short TTL for a payment id that the client can retry after 10 minutes.

Cluster: `{pay:req-9f3a}` on both the lock key and the result key if you use two keys.

A session is data that belongs to a browser or device after login. Redis is a common session store: low latency, TTL, and a shared store for many application processes.

Typical layout:

```text
HSET sess:a1b2 user 42 role staff
EXPIRE sess:a1b2 1800
```

Or one JSON string: `SET sess:a1b2 '{"user":42}' EX 1800`.

The application sets a cookie with the session id (a random value). The cookie is not the user id. Redis holds the mapping.

Sliding session: on each request, `EXPIRE` again (or `SET` with a new TTL). Absolute session: keep a `created` field and refuse after a max age even if `EXPIRE` slides.

Logout: `DEL` the session key.

Do not store passwords in the session. Store a user id and a small set of claims. Sign or encrypt the cookie if the id must be unguessable; the Redis key still needs a long random id.

Many processes share one Redis. That is the point. A process-local session map fails when the load balancer changes the process.

Persistence: sessions can be a cache. Loss logs everyone out. If that is not acceptable, enable persistence or accept the risk (topic 7, topic 12).

### Questions

#### Theoretical questions

1. What problem does an idempotency key solve?
2. What does `SET NX` tell you on the first vs second request?
3. What does the session cookie store, and what does Redis store?
4. What is a sliding session?
5. What happens to sessions if Redis restarts without persistence?

#### Easy practical tasks

1. `SET idem:lab:r1 1 NX EX 60`. Run it twice. Save both replies.
2. `HSET` a session hash. `EXPIRE` 60. `HGETALL`. Save the reply.
3. Write four sentences: retry of `INCR` vs retry with an idempotency key.
4. Draw browser cookie, application, Redis.

#### Medium practical tasks

1. Store JSON `{ "status":"done", "id":"x" }` after work. On `NX` fail, `GET` and return that JSON.
2. Implement create, read, slide `EXPIRE`, and `DEL` in your language. Show `TTL` after a slide.
3. Add an absolute `created` timestamp. Reject the session after 2 minutes even if you slide.

#### Advanced practical tasks

1. Design two keys or one hash: state, result, expiry. Write legal transitions (`new` → `pending` → `done` / `failed`).
2. Write a session policy: TTL, slide vs absolute, cookie flags (`HttpOnly`, `Secure`), and Redis persistence choice.

---

## Pub/Sub vs Streams

Pub/Sub:

```text
SUBSCRIBE news
PUBLISH news hello
```

Redis delivers `hello` to current subscribers. If no subscriber is connected, the message is gone. There is no backlog, no `ACK`, and no consumer group.

`PSUBSCRIBE news*` matches patterns. `UNSUBSCRIBE` leaves a channel.

A connection in subscribe mode does not run ordinary commands. Use a dedicated connection (topic 9).

At-most-once: Redis tries to send the message to current subscribers. If a subscriber is disconnected, slow, or not yet subscribed, that subscriber does not get a later replay from Redis. Redis does not store the message.

Slow subscribers: Redis can disconnect a client that does not read fast enough (`client-output-buffer-limit` for pubsub). That client loses messages.

`SUBSCRIBE` after `PUBLISH` does not receive the old message.

There is no `ACK`. The publisher does not know that a given subscriber processed the payload.

Streams (topic 5 intro):

- Entries have ids
- A consumer can read later
- Groups and `XACK` track work
- `XTRIM` / `MAXLEN` bound memory

Use Pub/Sub for live signals when loss is acceptable (cache invalidation hints, "reload flags", presence). Use streams when a worker must process each event at least once and you need a pending list.

Pub/Sub in Cluster uses the cluster bus to fan out. Delivery is still at-most-once for each subscriber.

Do not use Pub/Sub as a durable job queue. Do not use Pub/Sub for "charge this card" or "send this email once."

Use Pub/Sub for hints. Example: "product 1001 changed." The subscriber `DEL`s a cache key. If the hint is lost, a TTL still repairs the cache.

### Questions

#### Theoretical questions

1. What happens to a `PUBLISH` when no client is subscribed?
2. Why does Pub/Sub need a dedicated connection?
3. What do streams store that Pub/Sub does not store?
4. What does at-most-once mean for Pub/Sub?
5. Is Pub/Sub a durable job queue?

#### Easy practical tasks

1. In one `redis-cli`, `SUBSCRIBE lab:ch`. In another, `PUBLISH lab:ch hi`. Save the subscriber output.
2. `PUBLISH` before any `SUBSCRIBE`. Write what a late subscriber sees.
3. Write four sentences: Pub/Sub vs stream.
4. List two product signals that may use Pub/Sub and two that must use a stream.

#### Medium practical tasks

1. `PSUBSCRIBE lab:*`. Publish to `lab:a` and `lab:b`. Save what arrives.
2. Implement a cache-invalidation hint: `PUBLISH cache:inv product:1001`. A subscriber `DEL`s the cache key. Show one hit then a miss after publish.
3. `PUBLISH` then `SUBSCRIBE` in that order. Write what the subscriber receives.

#### Advanced practical tasks

1. Draw a table: Pub/Sub, list queue, streams. Rows: retention, ACK, competing consumers, Cluster notes, loss.
2. Read Cluster Pub/Sub docs. Write how a publish on one node reaches a subscriber on another node (high level).

---

## Consumer groups, `XACK`, pending entries

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

If the process dies before `XACK`, the entry stays pending. Another worker can claim it.

`XINFO GROUPS lab:orders` and `XINFO CONSUMERS lab:orders grpA` show lag-style fields and pending counts.

One stream can have many groups. Each group has its own cursor. The same entry can be delivered to every group. That is fan-out across groups, and competing consumers inside a group.

The PEL holds entries that a group delivered and that no consumer has acknowledged.

```text
XPENDING lab:orders grpA
XPENDING lab:orders grpA - + 10
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

Do not `XACK` before the side effect is durable, if you need at-least-once processing. Do not skip `XACK` if you already completed the side effect, or the PEL grows and another worker will do the work again. Make the side effect idempotent.

`XADD key * field value` appends an entry. An unbounded stream fills RAM. Use `MAXLEN` on `XADD` or `XTRIM` to cap the length.

### Questions

#### Theoretical questions

1. What does a consumer group add that `XREAD` does not store?
2. What does `>` mean in `XREADGROUP`?
3. What does `XACK` remove, and what does it not delete?
4. What does the PEL store?
5. What does the idle time in `XCLAIM` protect?

#### Easy practical tasks

1. `XGROUP CREATE lab:g g1 0 MKSTREAM`. `XADD` two entries. `XREADGROUP` as `c1`. Save the ids.
2. `XACK` those ids. `XPENDING lab:g g1`. Save the reply.
3. `XREADGROUP` one entry. Do not `XACK`. `XPENDING` the group. Save the output.
4. Write four sentences: `XACK` vs `XCLAIM`.

#### Medium practical tasks

1. Start two consumers `c1` and `c2` in one group. `XADD` four entries. Write which consumer received which id.
2. From a second consumer, `XCLAIM` the pending id with a small idle time (or wait). Process and `XACK`. Confirm `XPENDING` is empty.
3. Create a second group on the same stream. Show that the second group can read the same entries from its own cursor.

#### Advanced practical tasks

1. Implement a worker loop: `XREADGROUP BLOCK`, process, `XACK`. Crash before `XACK` (kill). Document the PEL.
2. Build a dead-letter path: after 5 deliveries, `XADD` to `lab:orders:dlq` and `XACK` the original. Test with a worker that always fails.

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
- Isolation from the cache Redis is mandatory (topic 12)

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

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Which patterns in this topic are correctness-critical, and what extra tool (fence, durable log, SQL constraint) do they need besides Redis?
2. How do a rate limit, an idempotency key, and a lock differ in the key lifetime and in the meaning of `NX`?
3. When do you choose Pub/Sub instead of a stream for the same event name?
4. How do entry ids, `>`, `$`, and the PEL fit on one diagram of a group?
5. Which volume and retention numbers push you from Redis Streams to Kafka?

#### Easy practical tasks

1. Write a cheat sheet: lock `SET`, rate-limit `INCR`, idempotency `SET NX`, session `HSET`, `PUBLISH`, `XREADGROUP`, `XACK`, `XPENDING`.
2. Draw stream, group, two consumers, PEL, and a `PUBLISH` channel on the side.
3. Write five lab key names with prefixes `lock:`, `rl:`, `idem:`, `sess:`, plus one stream key.
4. Write five lab rules: trim streams, unique consumer names, `XACK` policy, no Pub/Sub for jobs, idempotent claims.

#### Medium practical tasks

1. Build a tiny API: session create, a fixed-window limit, and a lock around one mock side effect. Record Redis keys after one success and one reject.
2. Build a two-consumer group on one stream. Kill one consumer. Claim and finish its PEL. Document commands.
3. Document retry + idempotency + rate limit together. Write the order of checks.

#### Advanced practical tasks

1. Implement lock + fencing + idempotency for a mock "create order" call. Write tests for double submit and lock expiry during work.
2. Write an architecture note: cache Redis, stream Redis, Kafka bus. Place six example events.
