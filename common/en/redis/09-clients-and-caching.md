# 9. Clients and Caching

## Description

An application talks to Redis through a client library. A cache stores a result so that a later read avoids a slow path. This topic covers popular clients, connection pools, timeouts, retries, serialization, safe key iteration, and cache write patterns.

Complete topic 2 (expiry and eviction) and topic 8 (Cluster clients) first. This topic assumes Redis is the cache, not the only source of truth, unless a task says otherwise.

Use one term for each concept. A client is a library or tool that sends Redis commands. A connection is one TCP (or TLS) session to an instance. A pool is a set of reusable connections. A timeout is a limit on wait time. A retry is a second send of the same logical command. Serialization is the conversion of an application object to bytes that Redis stores. The source of truth is the durable store (often a SQL or document database). A hit is a read that finds a usable value in Redis. A miss is a read that does not. A stampede is many misses for the same key at the same time. Negative caching stores a "missing" marker. Jitter is a random extra TTL.

---

## Popular clients and connection pooling

Redis speaks a text protocol (RESP, then RESP3). You can send commands with `redis-cli`. Production applications use a library in the application language.

Common libraries:

- Go: `go-redis` (`github.com/redis/go-redis`)
- Python: `redis-py` (the `redis` package)
- Node.js: `ioredis` and `node-redis`
- Java: Jedis (synchronous) and Lettuce (asynchronous, reactive)

These names change. Confirm the current official or recommended client on [https://redis.io/docs/clients/](https://redis.io/docs/clients/).

A good client does more than `SET` and `GET`. It manages connections. It parses replies. It supports pipelines, Pub/Sub, and Cluster or Sentinel when you need those modes. Topic 8 explained why a Cluster client must follow `MOVED`.

Pick one client per language for a new project. Read that client documentation for:

- Standalone, Sentinel, and Cluster constructors
- Pool settings
- Timeout settings
- Pipeline and transaction APIs
- How the client maps Redis errors to language errors

Do not mix two clients in one process without a reason. Two pools double the connection count.

`redis-cli` remains the lab tool. When a library fails, reproduce the same command in `redis-cli`. That split shows whether the problem is the server or the client.

A client version and a server version must match for new commands (ACL, Functions, `GEOSEARCH`). Check both.

Each Redis command needs a connection. A new TCP handshake for every `GET` is slow. A pool keeps open connections and lends one connection per command (or per pipeline).

Typical pool settings:

- Maximum open connections
- Minimum idle connections
- Idle timeout
- Maximum lifetime of one connection

Set the maximum from the application concurrency and from `maxclients` on the server (`INFO clients`, `CONFIG GET maxclients`). If 50 application processes each open 100 connections, Redis must accept 5000 clients. That number is often too high.

One connection is sequential. Two goroutines or threads must not share one connection without a mutex. The pool gives each concurrent caller its own connection.

Cluster mode needs connections to many nodes. The pool is per node or a shared map of pools. Sentinel needs connections to Sentinel processes and to the current primary.

Pub/Sub often needs a dedicated connection. A connection that is in `SUBSCRIBE` mode does not run ordinary `GET` commands. Topic 10 covers Pub/Sub.

Close the pool when the process stops. Leaked connections stay in `connected_clients` until they time out.

Do not create a new client object inside a hot request handler. Create one client (one pool) at process start.

### Questions

#### Theoretical questions

1. What protocol does a Redis client speak?
2. Name one popular client for Go, one for Python, one for Java.
3. Why does a pool exist?
4. Why must two threads not share one connection without a lock?
5. When do you create the client object?

#### Easy practical tasks

1. Open [https://redis.io/docs/clients/](https://redis.io/docs/clients/). Write the listed client for your language.
2. Open your client docs. Write the pool option names (max connections, idle timeout).
3. Run `INFO clients`. Write `connected_clients` and `maxclients` (or `CONFIG GET maxclients`).
4. Draw a process, a pool of three connections, and Redis.

#### Medium practical tasks

1. Read the client page for Sentinel vs Cluster. Write the two constructor option names.
2. Set a tiny `maxclients` on a lab instance (for example `10`). Open more connections than that. Record the error.
3. Start an HTTP handler that creates a new client per request (lab only). Watch `connected_clients`. Then move the client to startup. Compare the two counts.

#### Advanced practical tasks

1. Load-test with 1, 10, and 100 concurrent callers. Record p50 latency and `connected_clients`. Choose a max pool size and justify it in five sentences.
2. Document a pool budget for three services that share one Redis: process count, max connections each, and `maxclients` headroom.

---

## Timeouts, retries, serialization

A timeout stops an infinite wait. A client usually has:

- Connect timeout — wait for the TCP (or TLS) handshake
- Read timeout — wait for a reply
- Write timeout — wait until the send buffer accepts the command
- Pool timeout — wait for a free connection

Set timeouts. A missing read timeout can hang a request thread forever if Redis or the network stalls.

A retry repeats a command after a timeout, a disconnect, or a `MOVED` (Cluster). Retries help after a short network blip. Retries also create risk.

Safe to retry without extra design: `GET`, `EXISTS`, `SET` of the same full value, `PING`.

Unsafe to retry blindly: `INCR`, `DECR`, `LPUSH`, `RPOP`, `XADD` without an explicit id. The first send may have succeeded. The retry then applies the command twice.

Idempotent design (topic 10) uses a request id key so that a second apply is a no-op.

Backoff: wait a short time, then retry. Cap the retry count (for example 2 or 3). Do not retry in a tight loop. That loop can amplify load during an outage.

Cluster clients retry on `MOVED` and `ASK`. That retry is a redirect, not a duplicate-write policy. You still need idempotency for writes after a TCP error.

Do not set a read timeout shorter than your slowest legitimate command (large `LRANGE`, a heavy Lua script). Use the slow log (topic 11) to learn real command times.

Redis stores bytes (bulk strings) and typed structures (hashes, lists, streams). Redis does not know your language objects. You encode. You decode.

Common encodings:

- Plain string or integer — counters, flags, ids
- JSON — readable, wide tool support, larger than binary
- MessagePack — compact binary maps and arrays
- Protocol Buffers (protobuf) — schema, smaller payloads, code generation

A Redis hash can store one field per attribute (`HSET user:42 name Ada`). That layout lets you update one field. JSON in one string (`SET user:42 '{"name":"Ada"}'`) updates as a whole unless you use RedisJSON (topic 12).

Choose with these questions:

- Must you read or update one field?
- Must humans debug values in `redis-cli`?
- Is payload size a memory problem?
- Do many languages read the same key?

JSON in `redis-cli` is easy to read. Protobuf in `redis-cli` looks like binary. You need a decode step.

Do not mix encodings on one key. A JSON `GET` decoded as protobuf fails. Put the encoding in the key catalog.

Large values cost RAM and network. Compress only when you measure a win. Compression makes `redis-cli` harder and can hide useful structure.

### Questions

#### Theoretical questions

1. Name four timeout kinds that a client may expose.
2. Why is a retry of `INCR` dangerous after a network error?
3. Which commands are safer to retry than `INCR`?
4. When is a hash better than a JSON string?
5. Why must one key have one encoding?

#### Easy practical tasks

1. Open your client docs. Write the default connect timeout and read timeout (or "no default").
2. Write a two-column table: "Command" and "Safe to retry?". Add `GET`, `SET`, `INCR`, `LPUSH`.
3. `SET lab:j` to a small JSON object. `GET` it. Save the reply.
4. `HSET lab:h name Ada age 1`. `HGET lab:h name`. Save the reply.

#### Medium practical tasks

1. Set a 1 ms read timeout in a lab client. Run `GET` and a blocking `BLPOP` with a long wait. Record which call fails.
2. Write a retry wrapper that retries `GET` three times with sleep, and that never retries `INCR`. Show the branch in comments.
3. Encode the same struct as JSON and as MessagePack (or protobuf). Write both byte lengths.

#### Advanced practical tasks

1. Read your client source or docs for automatic retries. Write which commands it retries and whether you can disable that.
2. Measure `MEMORY USAGE` for 1000 JSON user blobs vs 1000 hashes with the same fields. Write the two totals and one caveat.

---

## `SCAN` instead of `KEYS *`

`KEYS pattern` walks the keyspace on the main thread and returns all matches. On a large instance, `KEYS *` can stall Redis for a long time. Other clients wait. In production, treat `KEYS` as a dangerous command (topic 11).

`SCAN cursor [MATCH pattern] [COUNT hint] [TYPE type]` returns a new cursor and a batch of keys. Repeat until the cursor is `0`.

```text
SCAN 0 MATCH lab:* COUNT 20
```

`COUNT` is a hint, not an exact page size. Redis may return more or fewer keys.

`SCAN` is not a snapshot. If keys appear or disappear during the walk, you can see a key twice, or you can miss a key. That is acceptable for metrics and for cleanup jobs. That is not acceptable as a single consistent export unless you accept the limit.

Collection scans:

- `HSCAN` — large hashes
- `SSCAN` — large sets
- `ZSCAN` — large sorted sets

`HGETALL`, `SMEMBERS`, and a full `ZRANGE` have the same class of risk as `KEYS` when the collection is huge. Use the `*SCAN` family.

Cluster: `SCAN` is per node. You must scan every primary (topic 8).

`redis-cli --scan --pattern 'lab:*'` uses `SCAN` for you.

Do not run `KEYS *` on a shared or production instance. Use `SCAN` in the lab when you practice. Disable or ACL-deny `KEYS` in production (topic 11).

### Questions

#### Theoretical questions

1. Why can `KEYS *` stall Redis?
2. What two things does one `SCAN` reply contain?
3. What does `COUNT` mean?
4. Why is `SCAN` not a consistent snapshot?
5. Which commands scan a hash, a set, and a sorted set?

#### Easy practical tasks

1. Create five keys `lab:a` … `lab:e`. Run `SCAN 0 MATCH lab:*`. Save the cursor and the keys.
2. Repeat `SCAN` from the new cursor until the cursor is `0`. Write how many rounds you needed.
3. Run `redis-cli --scan --pattern 'lab:*'`. Save the output.
4. Write four sentences: `KEYS` vs `SCAN`.

#### Medium practical tasks

1. Load 1000 lab keys. Time `KEYS lab:*` and a full `SCAN` loop. Write both times. Delete the keys.
2. During a `SCAN` loop, add and delete keys. Record whether a key appeared twice or disappeared.
3. `HSET` 200 fields. Compare `HGETALL` time with an `HSCAN` loop. Write the times.

#### Advanced practical tasks

1. Write a small program that deletes keys by prefix with `SCAN` + `DEL` in batches. Test on lab keys only. Write the batch size and a safety check (prefix must start with `lab:`).
2. On a Cluster lab, scan all primaries. Write the node list and the total key count vs `DBSIZE` on each primary.

---

## Cache-aside, write-through, write-behind

Cache-aside (lazy load) puts the cache logic in the application:

1. The application `GET`s the cache key.
2. On a hit, the application uses the value.
3. On a miss, the application reads the source of truth.
4. The application `SET`s the cache key (often with a TTL).
5. The application returns the value.

```text
GET cache:product:1001
# miss
# ... read SQL ...
SET cache:product:1001 <bytes> EX 300
```

The application owns the policy. Redis does not call the database. Redis does not know SQL.

On write to the source of truth, the application updates or deletes the cache key. Two common choices:

- Delete (`DEL`) the cache key. The next read reloads. This is simple.
- Write the new value into the cache. The next read hits. This needs the same encoding as the read path.

A delete can race: process A reads SQL (old), process B writes SQL and deletes the cache, process A then `SET`s the old value. Short TTLs limit the stale window. Some teams use a version field or a compare step.

If Redis is down, the application must read the source of truth and skip the cache (topic 12 circuit breaking). Cache-aside makes that path natural: a miss and a Redis error can both fall back to SQL.

Do not treat a cache hit as proof that the source of truth still has that row. The cache can be stale until TTL or invalidation.

Read-through and write-through move the load and store steps into a cache module that the application calls as if it were the database.

Read-through: the application asks the module for a key. The module `GET`s Redis. On a miss, the module reads the source of truth, `SET`s Redis, and returns the value. The handler does not contain the miss branch.

Write-through: the application writes through the module. The module writes the source of truth and writes Redis in the same application request (two stores, one call). The user request waits for both.

In Redis you still implement this in the application or in a library. Redis does not load SQL by itself.

Read-through hides the miss path. That can be cleaner. It can also hide cost: a "simple get" may run a slow query.

Write-through keeps Redis and the source of truth closer at write time. The request is slower than a write that only touches SQL. If Redis fails, you must decide: fail the request, or write SQL and accept a cache miss later.

Do not confuse write-through with a Redis-only database. The source of truth remains the durable store unless you design otherwise.

Write-behind (write-back) writes Redis first (or only Redis in the request). A background worker writes the source of truth later.

The user request is fast. The durable store lags. If Redis loses the key before the worker runs (eviction, restart without persistence, a bug), the durable store never sees the write.

Use write-behind only when you accept that risk, or when Redis is persistent and you have a replay queue (streams, topic 10) that the worker consumes.

Do not use write-behind for money or legal records unless a durable log exists outside a volatile cache.

Cache-aside plus `DEL` is simpler when the source of truth must win every time.

### Questions

#### Theoretical questions

1. Who reads the source of truth in cache-aside?
2. What are the five steps of a miss?
3. What two stores does write-through update in one request?
4. Which store does the user request write first in write-behind?
5. What happens if Redis drops the key before the write-behind worker runs?

#### Easy practical tasks

1. Write the five cache-aside steps as a numbered list in your own words.
2. Draw application, Redis, and SQL. Label hit and miss arrows.
3. Write four sentences: write-through vs write-behind.
4. Name a cache key and a TTL for a product page.

#### Medium practical tasks

1. Implement cache-aside for one key in your language: miss loads a file or a mock DB. Show hit on the second read.
2. Implement write-through for one field: write SQL (mock) and `SET` Redis. Confirm a following `GET`.
3. Restart Redis without persistence after `SET` but before a write-behind worker. Record the lost update.

#### Advanced practical tasks

1. Demonstrate the stale-`SET` race with two processes (or two sleeps). Write the timeline and one mitigation (TTL, version, or lock).
2. Design write-behind with a stream, a consumer group, and SQL. Write failure cases: worker crash, Redis restart, SQL error.

---

## Stampede, negative caching, TTL jitter

A stampede (thundering herd) happens when many clients miss the same key at once. Each client loads the source of truth. The database receives a burst. The cache then receives a burst of `SET`s.

Typical causes:

- The key expires (`TTL` reaches 0) while traffic is high
- Eviction removes a hot key
- A deploy flushes the cache
- The process starts with an empty cache (cold start)

Effects: high SQL load, high Redis CPU, slower user requests.

Mitigations:

- Single-flight: only one loader per key; other callers wait for that loader
- A lock (`SET NX EX`) around the load (topic 10)
- Soft TTL plus early refresh
- TTL jitter so that many keys do not expire in the same second
- Warm the cache before a deploy (topic 12)
- Negative caching so that a missing row does not hit SQL on every request

A stampede is a coordination problem. Redis can hold a lock or a "loading" flag. The application must use it.

Do not flush all keys on deploy if the dataset is large and hot. Prefer versioned keys or gradual warming.

A miss can mean "the source of truth has no row." If you store nothing, every request for a bad id hits SQL.

Negative caching stores a marker for a short TTL:

```text
SET cache:product:9999:neg 1 EX 30
```

Or you store an empty JSON object with a type field `{ "missing": true }`. The application must distinguish "missing" from "product with empty name".

Use a short TTL. A negative entry that lives too long hides a row that you just inserted.

On a successful create in SQL, delete the negative key. Otherwise the next read can still show "missing".

Do not negative-cache errors such as "SQL is down." That marker would hide recovery. Cache "not found", not "cannot talk to SQL".

TTL jitter adds a random extra lifetime to each key. If 10,000 product keys all get `EX 300`, they can expire in the same second after a bulk load. The next second is a stampede.

```text
# conceptual: 300 seconds plus a random 0–60
SET cache:product:1001 <bytes> EX 347
```

Each key expires at a different time. Reloads spread across a minute.

Jitter is not a lock. Jitter does not help one hot key that all clients share. Jitter helps many keys that share the same TTL and the same load time.

Pick a jitter range from the base TTL. A 300 second TTL plus 0–30 seconds is common. A jitter larger than the base TTL makes expiry hard to reason about.

Eviction can still remove many keys at once if memory is full. Jitter does not replace `maxmemory` planning (topic 2).

Probabilistic early expiration refreshes a hot key before the TTL hits zero, with a probability that rises as remaining time shrinks. Redis does not compute the probability. The application does. Combine with a lock or single-flight so that two early refreshes do not both hit SQL.

### Questions

#### Theoretical questions

1. What is a stampede?
2. Name three causes of a sudden miss storm?
3. What problem does negative caching solve?
4. Why must you not cache "SQL is down"?
5. Which stampede shape does TTL jitter reduce, and which shape does it not solve?

#### Easy practical tasks

1. Write four sentences: miss vs stampede.
2. Write a negative key name and a 30 second `SET`.
3. Write a formula `ttl = base + random(0, jitter)` with example numbers.
4. Write a lock key name for `cache:product:1001`.

#### Medium practical tasks

1. Start five parallel miss loaders without a lock (sleep in the loader). Count how many times the mock DB runs. Repeat with `SET loader:product:1001 NX EX 10`.
2. Implement: missing id writes `neg`. Second read does not call the mock DB. Then insert the row and `DEL` the `neg` key.
3. Load 200 keys with the same `EX 15`. Sample how many remain each second. Repeat with jitter 0–10. Compare the expiry curve.

#### Advanced practical tasks

1. Implement single-flight in your language (mutex map or `singleflight` style). Prove one load for ten waiters.
2. Combine jitter, negative caching, and single-flight in one small cache layer. Write a test plan with five cases.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do pool size, `maxclients`, and process count fit together in one budget?
2. When does a timeout plus a retry make a counter wrong?
3. Why is `SCAN` the production iteration tool even if `KEYS` is easier in a tutorial?
4. Which cache pattern keeps Redis optional on every read path, and why does that matter when Redis is down?
5. How do stampede, jitter, early refresh, and negative caching work as one design, not as four isolated tricks?

#### Easy practical tasks

1. Write a cheat sheet: client name, pool max, four timeouts, `SCAN` starter command, seven cache pattern names.
2. Draw the path: HTTP request, pool, connection, Redis, reply, decode JSON.
3. Draw one product-read flowchart that includes hit, miss, negative hit, and Redis error.
4. Mark three commands in your app as "retry yes" or "retry no".

#### Medium practical tasks

1. Build a 20-line program: pool at startup, `PING`, `SET`/`GET` of a JSON blob, `SCAN` of a prefix. Record the outputs.
2. Build cache-aside with jitter and negative markers for one mock entity. Record hits, misses, and SQL calls for 50 reads.
3. Document invalidation for update, delete, and create (including negative keys).

#### Advanced practical tasks

1. Write an application Redis runbook: versions, constructors, pool math, retry policy, key catalog, and a ban on `KEYS`.
2. Load-test one hot key with and without single-flight. Write SQL call counts and p99 latency.
