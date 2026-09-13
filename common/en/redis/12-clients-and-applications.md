# 12. Clients and Applications

## Description

An application talks to Redis through a client library. This topic covers popular clients, connection pools, timeouts, retries, key names, serialization, and safe key iteration.

Complete this topic before you study cache patterns and distributed coordination. Those topics assume a client that you control in code.

Use one term for each concept. A client is a library or tool that sends Redis commands. A connection is one TCP (or TLS) session to an instance. A pool is a set of reusable connections. A timeout is a limit on wait time. A retry is a second send of the same logical command. Serialization is the conversion of an application object to bytes that Redis stores.

---

## Official and popular clients (go-redis, redis-py, ioredis, Jedis/Lettuce)

Redis speaks a text protocol (RESP, then RESP3). You can send commands with `redis-cli`. Production applications use a library in the application language.

Common libraries:

- Go: `go-redis` (`github.com/redis/go-redis`)
- Python: `redis-py` (the `redis` package)
- Node.js: `ioredis` and `node-redis`
- Java: Jedis (synchronous) and Lettuce (asynchronous, reactive)

These names change. Confirm the current official or recommended client on [https://redis.io/docs/clients/](https://redis.io/docs/clients/).

A good client does more than `SET` and `GET`. It manages connections. It parses replies. It supports pipelines, Pub/Sub, and Cluster or Sentinel when you need those modes. Topic 11 explained why a Cluster client must follow `MOVED`. Topic 10 explained Sentinel failover.

Pick one client per language for a new project. Read that client documentation for:

- Standalone, Sentinel, and Cluster constructors
- Pool settings
- Timeout settings
- Pipeline and transaction APIs
- How the client maps Redis errors to language errors

Do not mix two clients in one process without a reason. Two pools double the connection count.

`redis-cli` remains the lab tool. When a library fails, reproduce the same command in `redis-cli`. That split shows whether the problem is the server or the client.

A client version and a server version must match for new commands (ACL, Functions, `GEOSEARCH`). Check both.

### Questions

#### Theoretical questions

1. What protocol does a Redis client speak?
2. Name one popular client for Go, one for Python, one for Java.
3. What extra work does a Cluster client do that a standalone client does not do?
4. Why do you confirm a command in `redis-cli` when a library fails?
5. Why must you check the client version against the server version?

#### Easy practical tasks

1. Open [https://redis.io/docs/clients/](https://redis.io/docs/clients/). Write the listed client for your language.
2. Open the README of that client. Write the constructor or `NewClient` example in one line.
3. Write a two-column table: "Client" and "Language". Add four rows from this section.
4. Run `PING` from `redis-cli` and from a five-line program in your language. Save both replies.

#### Medium practical tasks

1. Read the client page for Sentinel vs Cluster. Write the two constructor option names.
2. List five command helpers in the client (`Get`, `Set`, `HGetAll`, pipeline, `Eval`). Write the method names as the library spells them.
3. Print the client library version and `INFO server` `redis_version`. Write both numbers.

#### Advanced practical tasks

1. Compare two clients in the same language (for example Jedis vs Lettuce, or `node-redis` vs `ioredis`) in a six-row table: API style, Cluster, pool, pipeline, maintenance, license.
2. Write a one-page "client choice" note for your team: library name, minimum version, standalone only or Cluster, and who owns upgrades.

---

## Connection pooling

Each Redis command needs a connection. A new TCP handshake for every `GET` is slow. A pool keeps open connections and lends one connection per command (or per pipeline).

Typical pool settings:

- Maximum open connections
- Minimum idle connections
- Idle timeout
- Maximum lifetime of one connection

Set the maximum from the application concurrency and from `maxclients` on the server (`INFO clients`, `CONFIG GET maxclients`). If 50 application processes each open 100 connections, Redis must accept 5000 clients. That number is often too high.

One connection is sequential. Two goroutines or threads must not share one connection without a mutex. The pool gives each concurrent caller its own connection.

Cluster mode needs connections to many nodes. The pool is per node or a shared map of pools. Sentinel needs connections to Sentinel processes and to the current primary.

Pub/Sub often needs a dedicated connection. A connection that is in `SUBSCRIBE` mode does not run ordinary `GET` commands. Topic 14 covers Pub/Sub.

Close the pool when the process stops. Leaked connections stay in `connected_clients` until they time out.

Do not create a new client object inside a hot request handler. Create one client (one pool) at process start.

### Questions

#### Theoretical questions

1. Why does a pool exist?
2. What happens if many processes each open a large pool?
3. Why must two threads not share one connection without a lock?
4. Why does Pub/Sub often need a dedicated connection?
5. When do you create the client object?

#### Easy practical tasks

1. Open your client docs. Write the pool option names (max connections, idle timeout).
2. Run `INFO clients`. Write `connected_clients` and `maxclients` (or `CONFIG GET maxclients`).
3. Write four sentences: one connection vs a pool.
4. Draw a process, a pool of three connections, and Redis.

#### Medium practical tasks

1. Set a tiny `maxclients` on a lab instance (for example `10`). Open more connections than that. Record the error.
2. Start an HTTP handler that creates a new client per request (lab only). Watch `connected_clients`. Then move the client to startup. Compare the two counts.
3. Read Cluster pool behavior in your client. Write whether the library opens one pool per node.

#### Advanced practical tasks

1. Load-test with 1, 10, and 100 concurrent callers. Record p50 latency and `connected_clients`. Choose a max pool size and justify it in five sentences.
2. Document a pool budget for three services that share one Redis: process count, max connections each, and `maxclients` headroom.

---

## Timeouts and retries

A timeout stops an infinite wait. A client usually has:

- Connect timeout — wait for the TCP (or TLS) handshake
- Read timeout — wait for a reply
- Write timeout — wait until the send buffer accepts the command
- Pool timeout — wait for a free connection

Set timeouts. A missing read timeout can hang a request thread forever if Redis or the network stalls.

A retry repeats a command after a timeout, a disconnect, or a `MOVED` (Cluster). Retries help after a short network blip. Retries also create risk.

Safe to retry without extra design: `GET`, `EXISTS`, `SET` of the same full value, `PING`.

Unsafe to retry blindly: `INCR`, `DECR`, `LPUSH`, `RPOP`, `XADD` without an explicit id. The first send may have succeeded. The retry then applies the command twice.

Idempotent design (topic 14) uses a request id key so that a second apply is a no-op.

Backoff: wait a short time, then retry. Cap the retry count (for example 2 or 3). Do not retry in a tight loop. That loop can amplify load during an outage.

Cluster clients retry on `MOVED` and `ASK`. That retry is a redirect, not a duplicate-write policy. You still need idempotency for writes after a TCP error.

Do not set a read timeout shorter than your slowest legitimate command (large `LRANGE`, a heavy Lua script). Use the slow log (topic 16) to learn real command times.

### Questions

#### Theoretical questions

1. Name four timeout kinds that a client may expose.
2. What is the risk of no read timeout?
3. Why is a retry of `INCR` dangerous after a network error?
4. Which commands are safer to retry than `INCR`?
5. Why must a retry loop have a cap and a backoff?

#### Easy practical tasks

1. Open your client docs. Write the default connect timeout and read timeout (or "no default").
2. Write a two-column table: "Command" and "Safe to retry?". Add `GET`, `SET`, `INCR`, `LPUSH`.
3. Write four sentences: timeout vs retry.
4. List three errors that your client can raise on a closed connection.

#### Medium practical tasks

1. Set a 1 ms read timeout in a lab client. Run `GET` and a blocking `BLPOP` with a long wait. Record which call fails.
2. Write a retry wrapper that retries `GET` three times with sleep, and that never retries `INCR`. Show the branch in comments.
3. Simulate a dropped connection during `SET` (stop Redis, then start it). Record whether your client retried and what the key value is.

#### Advanced practical tasks

1. Read your client source or docs for automatic retries. Write which commands it retries and whether you can disable that.
2. Design an idempotent increment: a request-id key plus a Lua script or a compare step. Write the keys and the failure cases. Do not claim it is perfect.

---

## Key naming design

Redis has one flat keyspace per logical DB. You invent the names. A name is a contract between all writers and readers.

Rules that scale:

- Use one separator (`:` is common). Keep it.
- Put a domain first: `cache:`, `sess:`, `lock:`, `rl:`.
- Put the object type and the id: `user:42:profile`.
- Keep names short enough to read, long enough to avoid collisions.
- Do not put secrets or personal data in the key string. Names appear in `MONITOR`, in `SCAN`, and in support dumps.
- Document the scheme on one page.

Cluster adds hash tags. Keys that must join in one multi-key command share `{tag}` (topic 11). Example: `{user:42}:profile` and `{user:42}:cart`.

A shared instance needs a product or tenant prefix: `shopA:user:42`. That prefix is a convention. It is not a security boundary. Topic 21 covers multi-tenant policy.

Temporary keys need a purpose, a unique suffix, and a TTL:

```text
tmp:export:2026-09-13:job9f3a
```

Avoid names such as `data`, `tmp`, or `key1` in shared environments.

Version a schema in the name when the value layout changes and old keys may still exist: `user:42:profile:v2`. Plan a migration. Do not silently change the meaning of a live key.

### Questions

#### Theoretical questions

1. Why is a key name a contract?
2. What belongs in a typical production key (domain, type, id)?
3. Why do you avoid secrets in key names?
4. When do you add a `{hash tag}`?
5. Why is a tenant prefix not a security boundary?

#### Easy practical tasks

1. Write ten key names for a shop: user, cart, session, cache, lock. Use one separator.
2. Rewrite two bad names (`data`, `tmp`) into documented names with TTLs.
3. Write four sentences: prefix vs logical `SELECT` (topic 1).
4. Open topic 2 in this handbook. Copy the separator rule into your own words (three sentences).

#### Medium practical tasks

1. Draw a key catalog table: prefix, type, TTL policy, example key. Add six rows.
2. Design Cluster keys for one user profile and cart that must share a slot. Write the names.
3. Find five keys on your lab instance with `SCAN`. Mark which names violate your catalog.

#### Advanced practical tasks

1. Write a one-page key standard for a team: separator, prefixes, hash tags, forbidden names, and review steps.
2. Plan a `v1` to `v2` profile key migration. Write dual-read order and when you delete `v1`.

---

## Serialization: strings, JSON, protobuf, MessagePack

Redis stores bytes (bulk strings) and typed structures (hashes, lists, streams). Redis does not know your language objects. You encode. You decode.

Common encodings:

- Plain string or integer — counters, flags, ids
- JSON — readable, wide tool support, larger than binary
- MessagePack — compact binary maps and arrays
- Protocol Buffers (protobuf) — schema, smaller payloads, code generation

A Redis hash can store one field per attribute (`HSET user:42 name Ada`). That layout lets you update one field. JSON in one string (`SET user:42 '{"name":"Ada"}'`) updates as a whole unless you use RedisJSON (topic 18).

Choose with these questions:

- Must you read or update one field?
- Must humans debug values in `redis-cli`?
- Is payload size a memory problem?
- Do many languages read the same key?

JSON in `redis-cli` is easy to read. Protobuf in `redis-cli` looks like binary. You need a decode step.

Do not mix encodings on one key. A JSON `GET` decoded as protobuf fails. Put the encoding in the key catalog (`user:42:profile` is JSON, `user:42:blob` is protobuf).

Large values cost RAM and network. Compress only when you measure a win. Compression makes `redis-cli` harder and can hide useful structure.

Hashes, lists, and streams still need a per-field encoding. A stream field is a string. You can put JSON in a field named `payload`.

### Questions

#### Theoretical questions

1. What does Redis store if you "set an object"?
2. When is a hash better than a JSON string?
3. What advantage does JSON have over protobuf in a lab?
4. Why must one key have one encoding?
5. What extra module updates a JSON path without a full rewrite?

#### Easy practical tasks

1. `SET lab:j` to a small JSON object. `GET` it. Save the reply.
2. `HSET lab:h name Ada age 1`. `HGET lab:h name`. Save the reply.
3. Write a three-column table: "Encoding", "Readable in redis-cli", "Schema".
4. Write four sentences: hash fields vs one JSON string.

#### Medium practical tasks

1. Encode the same struct as JSON and as MessagePack (or protobuf). Write both byte lengths.
2. Document decode errors: put JSON in a key, decode as the other format. Record the error.
3. Design stream fields: `type` (plain) and `payload` (JSON). `XADD` one entry. `XRANGE` and parse.

#### Advanced practical tasks

1. Write a team rule: which prefixes use hashes, which use JSON, which use protobuf. Give one reason each.
2. Measure `MEMORY USAGE` for 1000 JSON user blobs vs 1000 hashes with the same fields. Write the two totals and one caveat (encoding internals, topic 16).

---

## Avoiding huge `KEYS *` in production (`SCAN`)

`KEYS pattern` walks the keyspace on the main thread and returns all matches. On a large instance, `KEYS *` can stall Redis for a long time. Other clients wait. Topic 2 and topic 1 already warn about this. In production, treat `KEYS` as a dangerous command (topic 17).

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

Cluster: `SCAN` is per node. You must scan every primary (topic 11).

`redis-cli --scan --pattern 'lab:*'` uses `SCAN` for you.

Do not run `KEYS *` on a shared or production instance. Use `SCAN` in the lab when you practice. Disable or ACL-deny `KEYS` in production (topic 17).

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

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do pool size, `maxclients`, and process count fit together in one budget?
2. How do key naming, serialization, and Cluster hash tags appear on one catalog page?
3. When does a timeout plus a retry make a counter wrong?
4. Why is `SCAN` the production iteration tool even if `KEYS` is easier in a tutorial?
5. What do you check first when an application error might be the client or the server?

#### Easy practical tasks

1. Write a cheat sheet: client name, pool max, four timeouts, key prefix list, `SCAN` starter command.
2. Draw the path: HTTP request, pool, connection, Redis, reply, decode JSON.
3. Mark three commands in your app as "retry yes" or "retry no".
4. Export `INFO clients` and `INFO server`. Highlight `connected_clients`, `redis_version`, and `tcp_port`.

#### Medium practical tasks

1. Build a 20-line program: pool at startup, `PING`, `SET`/`GET` of a JSON blob, `SCAN` of a prefix. Record the outputs.
2. Document failure modes: Redis down, timeout, `WRONGTYPE`, `NOSCRIPT`. Write the user-visible behavior for each.
3. Review an existing codebase (or a sample) for `KEYS` and per-request client creation. Write findings.

#### Advanced practical tasks

1. Write an application Redis runbook: versions, constructors, pool math, retry policy, key catalog, and a ban on `KEYS`.
2. Compare standalone client config and Cluster client config for the same app in a table: addresses, pool, redirects, multi-key, `SCAN`.
