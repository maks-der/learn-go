# 14. Distributed Coordination

## Description

Applications use Redis as a shared place for locks, limits, request ids, sessions, and flags. This topic covers those patterns and the difference between Pub/Sub and streams.

Complete topic 3 (string locks), topic 8 (Lua), and topic 12 (clients) first. Treat every pattern as application design plus Redis commands, not as a consensus product.

Use one term for each concept. A lock is a key that only one worker may hold. A rate limit is a cap on events in a time window. An idempotency key is a stored request id that makes a second apply a no-op. A session is per-user state with a TTL. A feature flag is a stored switch. Pub/Sub is fire-and-forget broadcast. A stream is a log with ids.

---

## Distributed lock (`SET NX EX`) and the Redlock debate

Topic 3 introduced a simple lock:

```text
SET lock:resource holder-id NX EX 30
```

`OK` means this worker holds the lock. A null reply means another worker holds it. The TTL frees the key if the worker dies.

Release must delete only your lock. A Lua script (topic 8) does `GET` then `DEL` if the value equals `holder-id`. A two-step `GET` plus `DEL` from the client can remove another worker's lock after expiry.

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

1. Run the lock-release Lua from topic 8. Test match and mismatch. Save replies.
2. Let the TTL expire while a "worker" sleeps. Acquire from a second client. Write the overlap risk in five sentences.
3. Add `INCR lock:resource:fence` on acquire. Write how a storage API would reject a lower token.

#### Advanced practical tasks

1. Read Kleppmann's Redlock essay (or an equivalent critique) and the Redis lock page. Write a six-row table: issue, lock side, critique side.
2. Implement acquire, fence, and Lua release in your language. List what the design still does not guarantee.

---

## Rate limiting (fixed window, sliding window, token bucket)

A rate limit rejects or delays events when a client exceeds a budget.

Fixed window:

```text
INCR rl:user:42:202609131857
EXPIRE rl:user:42:202609131857 60
```

Use the current minute (or second) in the key. If the count is above N, reject. `EXPIRE` on first increment (Lua can set expiry only when the value is `1`).

The limit resets at the window edge. A user can send N events at the end of one minute and N at the start of the next (2N in two seconds).

Sliding window: store event timestamps in a sorted set. `ZREMRANGEBYSCORE` drops old members. `ZCARD` is the count. `ZADD` the new event if under the limit. Topic 6 covered the commands. Memory grows with events in the window.

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

## Idempotency keys

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

### Questions

#### Theoretical questions

1. What problem does an idempotency key solve?
2. What does `SET NX` tell you on the first vs second request?
3. Why does the key need a long enough TTL?
4. What gap exists if the process dies after `SET NX` and before the side effect?
5. Why is Redis not a substitute for idempotency in the durable store?

#### Easy practical tasks

1. `SET idem:lab:r1 1 NX EX 60`. Run it twice. Save both replies.
2. Write four sentences: retry of `INCR` vs retry with an idempotency key.
3. Write a key name for a payment request id.
4. Draw first request vs retry: two clients, one key, one bank.

#### Medium practical tasks

1. Store JSON `{ "status":"done", "id":"x" }` after work. On `NX` fail, `GET` and return that JSON.
2. Simulate crash after `NX` success: leave `pending`. Write a reconcilers steps in six sentences.
3. Combine with topic 12 retries: a client retries `POST` with the same id. Show one side effect.

#### Advanced practical tasks

1. Design two keys or one hash: state, result, expiry. Write legal transitions (`new` → `pending` → `done` / `failed`).
2. Compare HTTP `Idempotency-Key` headers in a public API doc with this Redis pattern. Write three matches and one difference.

---

## Session store

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

Persistence: sessions can be a cache. Loss logs everyone out. If that is not acceptable, enable persistence or accept the risk (topic 9, topic 21).

### Questions

#### Theoretical questions

1. What does the cookie store, and what does Redis store?
2. Why do many processes share one session Redis?
3. What is a sliding session?
4. What does logout do to the key?
5. What happens to sessions if Redis restarts without persistence?

#### Easy practical tasks

1. `HSET` a session hash. `EXPIRE` 60. `HGETALL`. Save the reply.
2. Write four sentences: session hash vs session JSON string.
3. Draw browser cookie, application, Redis.
4. Write one field list that you must not put in a session.

#### Medium practical tasks

1. Implement create, read, slide `EXPIRE`, and `DEL` in your language. Show `TTL` after a slide.
2. Add an absolute `created` timestamp. Reject the session after 2 minutes even if you slide.
3. Compare `MEMORY USAGE` for 1000 small hashes vs 1000 JSON strings.

#### Advanced practical tasks

1. Write a session policy: TTL, slide vs absolute, cookie flags (`HttpOnly`, `Secure`), and Redis persistence choice.
2. Design session revocation for "log out all devices": a user version key that sessions must match.

---

## Feature flags

A feature flag is a stored switch. The application reads the flag and chooses a code path. Redis can store flags for fast, central change without a deploy.

```text
HSET flags checkout_v2 1 search_v3 0
HGET flags checkout_v2
```

Or one key per flag: `SET flag:checkout_v2 1`.

A hash keeps all flags in one `HGETALL` (small hash only). Many flags with large per-user overrides need a different layout (`flag:checkout_v2:user:42`).

Cache flags in the process for a few seconds if every request would otherwise `HGET`. Then accept a short delay after a change.

Flags are not a security boundary. A user who can write Redis can turn flags. Protect Redis (topic 17). Use ACLs so that the application user cannot `FLUSHALL`.

Do not leave dead flags forever. Name an owner and an expiry date in a catalog.

Complex targeting (percent rollout, user lists) can use sets, hashes, or a Lua script. Start with a global boolean.

### Questions

#### Theoretical questions

1. What is a feature flag?
2. Why can Redis change a path without a deploy?
3. When is one hash better than one key per flag?
4. Why is a flag not a security control?
5. Why do you cache flags in process for a few seconds?

#### Easy practical tasks

1. `HSET flags lab_feat 1`. `HGET flags lab_feat`. Save the reply.
2. Flip the flag to `0`. Read again.
3. Write four sentences: flag in Redis vs a config file on disk.
4. Write a catalog row: name, owner, default, Redis key.

#### Medium practical tasks

1. Read flags once per 5 seconds in a loop. Change the hash in `redis-cli`. Write when the process sees the change.
2. Add a set `flag:lab:allow` of user ids. Check membership with `SISMEMBER`. Write the rule.
3. Document an ACL that can `HGET` flags but cannot `HSET` (topic 17 preview). Write the intent in four sentences.

#### Advanced practical tasks

1. Design a 10 percent rollout: `INCR` a counter or hash a user id. Write how you keep the same user stable.
2. Write a flag cleanup process: list flags older than 90 days, owners, and delete steps.

---

## Pub/Sub (`PUBLISH`, `SUBSCRIBE`) vs Streams

Pub/Sub:

```text
SUBSCRIBE news
PUBLISH news hello
```

Redis delivers `hello` to current subscribers. If no subscriber is connected, the message is gone. There is no backlog, no `ACK`, and no consumer group.

`PSUBSCRIBE news*` matches patterns. `UNSUBSCRIBE` leaves a channel.

A connection in subscribe mode does not run ordinary commands. Use a dedicated connection (topic 12).

Streams (topic 7 intro, topic 15 full):

- Entries have ids
- A consumer can read later
- Groups and `XACK` track work
- `XTRIM` / `MAXLEN` bound memory

Use Pub/Sub for live signals when loss is acceptable (cache invalidation hints, "reload flags", presence). Use streams when a worker must process each event at least once and you need a pending list.

Pub/Sub in Cluster uses the cluster bus to fan out. Delivery is still at-most-once for each subscriber.

Do not use Pub/Sub as a durable job queue.

### Questions

#### Theoretical questions

1. What happens to a `PUBLISH` when no client is subscribed?
2. Why does Pub/Sub need a dedicated connection?
3. What do streams store that Pub/Sub does not store?
4. When is Pub/Sub a better fit than a stream?
5. Is Pub/Sub a durable job queue?

#### Easy practical tasks

1. In one `redis-cli`, `SUBSCRIBE lab:ch`. In another, `PUBLISH lab:ch hi`. Save the subscriber output.
2. `PUBLISH` before any `SUBSCRIBE`. Write what a late subscriber sees.
3. Write four sentences: Pub/Sub vs stream.
4. List two product signals that may use Pub/Sub and two that must use a stream.

#### Medium practical tasks

1. `PSUBSCRIBE lab:*`. Publish to `lab:a` and `lab:b`. Save what arrives.
2. Implement a cache-invalidation hint: `PUBLISH cache:inv product:1001`. A subscriber `DEL`s the cache key. Show one hit then a miss after publish.
3. Compare `INFO clients` blocked or pubsub fields (names vary by version) before and after `SUBSCRIBE`.

#### Advanced practical tasks

1. Draw a table: Pub/Sub, list queue, streams. Rows: retention, ACK, competing consumers, Cluster notes, loss.
2. Read Cluster Pub/Sub docs. Write how a publish on one node reaches a subscriber on another node (high level).

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Which patterns in this topic are correctness-critical, and what extra tool (fence, durable log, SQL constraint) do they need besides Redis?
2. How do a rate limit, an idempotency key, and a lock differ in the key lifetime and in the meaning of `NX`?
3. When do you choose Pub/Sub instead of a stream for the same event name?
4. How do sessions and feature flags differ in TTL policy and in who may write the key?
5. Why is "Redis is up" not enough to call a lock safe after failover?

#### Easy practical tasks

1. Write a cheat sheet: lock `SET`, rate-limit `INCR`, idempotency `SET NX`, session `HSET`, flag `HGET`, `PUBLISH`.
2. Draw one user request that checks a flag, a rate limit, and a session key.
3. Write five lab key names with prefixes `lock:`, `rl:`, `idem:`, `sess:`, `flag:`.
4. List commands that must use Lua for atomic read-then-write in this topic.

#### Medium practical tasks

1. Build a tiny API: session create, a fixed-window limit, and a flag gate. Record Redis keys after one success and one reject.
2. Write a decision page: lock vs Lua-only update vs SQL unique constraint.
3. Document retry + idempotency + rate limit together. Write the order of checks.

#### Advanced practical tasks

1. Implement lock + fencing + idempotency for a mock "create order" call. Write tests for double submit and lock expiry during work.
2. Write a production note: when this team allows Redlock, when it forbids it, and what it uses instead for payments.
