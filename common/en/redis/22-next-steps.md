# 22. Next Steps

## Description

This topic turns the handbook into practice. You complete official command drills, you build a cache-aside layer, you implement a rate limiter, you run a replica, and you compare a lock with a published critique.

Complete topics 1 through 21 first. Use a local lab. Do not use a shared production instance.

Use one term for each concept. A drill is a short official exercise. A layer is a small library in your language. A replica is a second Redis that copies a primary (topic 10). A critique is a public analysis of lock algorithms (topic 14). Official docs remain [https://redis.io/docs/](https://redis.io/docs/).

---

## Complete the official getting-started commands

Return to [https://redis.io/docs/](https://redis.io/docs/) and the command index [https://redis.io/commands/](https://redis.io/commands/). Run a minimum set in `redis-cli` until the replies are familiar without this handbook.

Minimum set:

- `PING`, `INFO server`, `INFO memory`
- `SET`, `GET`, `MSET`, `MGET`, `INCR`
- `EXPIRE`, `TTL`
- `HSET`, `HGET`, `HGETALL`
- `LPUSH`, `LRANGE`, `LPOP`
- `SADD`, `SMEMBERS`, `SISMEMBER`
- `ZADD`, `ZRANGE`
- `XADD`, `XRANGE`
- `SCAN 0`
- `MULTI` / `EXEC`

Open each command page. Write complexity and return type. Check `redis_version` if a flag is missing.

Redis University ([https://university.redis.io/](https://university.redis.io/)) is optional. It does not replace this path.

The suggested practice list in `redis.topics.md` is a second checklist. Finish those ten items if you have not.

Do not skip `SCAN` and still use `KEYS *` on any instance that is not a throwaway lab.

Bookmark the data types index: [https://redis.io/docs/data-types/](https://redis.io/docs/data-types/).

### Questions

#### Theoretical questions

1. Where is the official command reference?
2. Why do you read complexity on a command page?
3. Why must you check `redis_version` against a flag?
4. What is Redis University in relation to this handbook?
5. Why is `SCAN` in the minimum set?

#### Easy practical tasks

1. Run the minimum set. Save one reply per command family (string, hash, list, set, zset, stream).
2. Open five command pages. Write complexity and return type for each.
3. Write a one-page cheat sheet from memory. Then diff it with the docs.
4. Tick the ten items in `redis.topics.md` Suggested Practice Order. Mark done or date.

#### Medium practical tasks

1. Follow one official getting-started tutorial. Write three facts that this handbook already covered and one that it did not.
2. Export `INFO server` and a 20-line command log from your drill.
3. Use `HELP` or `COMMAND DOCS` if your version supports it. Compare with the web page for `SET`.

#### Advanced practical tasks

1. Build a personal index: topic number, official URL, one sentence. Cover topics 1–21.
2. Complete one Redis University course or an official learning path. Write what overlapped this handbook.

---

## Build a cache-aside layer in your language

Write a small module, not a framework.

Required behavior (topic 13, topic 12, topic 21):

- `Get(key)` — `GET` Redis; on miss, load from a source function; `SET` with TTL and jitter
- `Invalidate(key)` — `DEL`
- Redis errors — treat as miss if Redis is a cache; optional circuit breaker
- One client at process start (pool)
- No `KEYS`
- Documented key prefix (`cache:`)
- Tests: hit, miss, invalidation, Redis down

Keep the source of truth as a function or interface (`LoadProduct(id)`). Do not put SQL drivers in the first version if a mock is enough.

Add negative caching only after the simple path works.

Serialization: start with JSON strings. Note the encoding in a comment (topic 12).

Do not add write-behind.

Put the module in your language's usual layout (a package, a class). Keep it under a few hundred lines.

### Questions

#### Theoretical questions

1. What are the five cache-aside steps inside `Get`?
2. What does the layer do when Redis times out and Redis is a cache?
3. Why does the layer take a load function?
4. Why is `KEYS` forbidden in the layer?
5. Why do you delay write-behind?

#### Easy practical tasks

1. Write the function signatures for `Get` and `Invalidate`.
2. Write four test names (hit, miss, invalidate, redis down).
3. Write the key formula `cache:{type}:{id}`.
4. Write four sentences: this layer vs read-through (topic 13).

#### Medium practical tasks

1. Implement the layer with a mock loader. Show a second `Get` that does not call the loader.
2. Stop Redis. Show `Get` still returns the mock value.
3. Add TTL jitter. Print the `EX` value that you send.

#### Advanced practical tasks

1. Add single-flight for one key. Prove ten parallel `Get`s call the loader once.
2. Add a breaker and a catalog comment. Write a README of 20 lines: how to use, how to fail.

---

## Implement a rate limiter

Build one limiter that many processes can share (topic 14).

Pick one algorithm and document why:

- Fixed window — `INCR` + `EXPIRE` (Lua sets expiry on first increment)
- Sliding window — sorted set of timestamps
- Token bucket — Lua on one key

Required behavior:

- `Allow(subject) -> yes/no` (and optional remaining count)
- Cluster-safe key with a hash tag if you test Cluster (`{user:42}`)
- TTL or refill so that keys die
- Tests: under limit, at limit, after window reset
- No `KEYS`

Use the durable or a dedicated Redis, not the cache instance, if you cannot lose limits during a cache flush (topic 21). For a lab, one instance is fine.

Do not claim a perfect global limiter for a multi-region active-active design. State the instance as the scope.

### Questions

#### Theoretical questions

1. Why must the limiter live in Redis, not only in one process?
2. What burst does a fixed window allow at the boundary?
3. Why does first-increment expiry need Lua (or `SET` with a combined pattern)?
4. What is the scope of the limit (one Redis)?
5. Which instance role do you choose if a flush would reset limits?

#### Easy practical tasks

1. Write the key name for user `42` and a 60 second window.
2. Choose an algorithm. Write three reasons.
3. Write four test names.
4. Write four sentences: limiter vs idempotency key.

#### Medium practical tasks

1. Implement `Allow` for a fixed window of 5 per 10 seconds. Show yes, yes, no. Wait and show yes.
2. Port the increment+expire Lua. Test two increments, one TTL.
3. Add a second subject and show independent counters.

#### Advanced practical tasks

1. Implement token bucket Lua. Test burst and empty bucket. Save replies.
2. Add a hash tag and run against Cluster or write why a multi-key design would `CROSSSLOT`.

---

## Run a replica locally

Practice topic 10 on your machine.

Shape:

1. Start primary on port `6379`.
2. Start replica on port `6380`.
3. `REPLICAOF 127.0.0.1 6379` on the replica.
4. `SET` on the primary. `GET` on the replica.
5. `INFO replication` on both.

Docker example (adjust names):

```text
docker run --name redis-p -p 6379:6379 redis
docker run --name redis-r -p 6380:6379 redis
```

Then `REPLICAOF` using a host that the replica container can reach (`host.docker.internal` on Docker Desktop, or a user-defined network and the primary container name). Document the exact commands that worked.

Confirm:

- Replica is read-only by default
- A write on the replica fails
- After `REPLICAOF NO ONE`, the old replica is a standalone copy

Optional: enable AOF on the replica and restart. Topic 9.

Do not expose either port past localhost.

This drill is not Sentinel and not Cluster. Those are extra if you want them (topic 10, topic 11).

### Questions

#### Theoretical questions

1. What command makes an instance a replica?
2. Why does a write on the replica fail by default?
3. What does `REPLICAOF NO ONE` do?
4. What does `INFO replication` show on the primary?
5. Why is a replica not a backup?

#### Easy practical tasks

1. Start two instances. Write both ports.
2. `REPLICAOF` and `INFO replication` on both. Write `role` lines.
3. `SET` on the primary. `GET` on the replica. Save the value.
4. Write four sentences: replica vs Sentinel.

#### Medium practical tasks

1. Time the lag: `SET` then immediate `GET` on the replica in a loop. Write whether you ever saw a miss.
2. Stop the primary. Write what the replica serves and whether it accepts writes.
3. Copy an RDB from the replica. Write why that file can still be useful.

#### Advanced practical tasks

1. Document Docker networking that you used (`host`, bridge, `host.docker.internal`). Draw the path.
2. Add a second replica. Then try a simple failover drill (`REPLICAOF NO ONE` on one replica, point the other at it). Write risks vs Sentinel.

---

## Compare a lock implementation with a paper or Redlock critique

Topic 3 and topic 14 gave `SET NX EX`, Lua release, fencing, and the Redlock debate.

Now do a written comparison. You do not need to implement Redlock.

Read:

- The Redis documentation page on distributed locks
- A critique such as Martin Kleppmann's essay on Redlock (or another serious public critique)
- Optional: the Redis authors' reply

Write a table:

```text
Issue              Simple SET NX EX     Redlock (as described)    Your notes
Pause / GC
Clock / TTL
Failover loss
Fencing token
What it is for
```

Then implement or reuse a simple lock in your language:

- Unique holder id
- `SET NX EX`
- Lua release
- Optional `INCR` fence

List what your lock does not guarantee. State a use that is acceptable (single-instance cache rebuild) and a use that is not (money transfer without a durable fence or SQL constraint).

Do not ship Redlock in production only because a blog post praised it. Do not forbid all Redis locks because a critique exists. Match the tool to the risk.

Keep the write-up to one or two pages.

### Questions

#### Theoretical questions

1. What does `SET NX EX` guarantee on one healthy instance?
2. What extra claim did Redlock try to make?
3. What is a fencing token?
4. Why can a critique and the official page both be useful?
5. When is a Redis lock the wrong tool?

#### Easy practical tasks

1. Open the official lock page and one critique. Write the two URLs and the dates you read them.
2. Fill three rows of the issue table.
3. Write four sentences: lock vs Lua-only Redis update.
4. Write one acceptable use and one forbidden use for your team.

#### Medium practical tasks

1. Implement acquire + Lua release. Test a second acquirer. Test mismatch release.
2. Add a fence number. Write how a mock storage API rejects a stale fence.
3. Quote one claim from the critique and one from the official page in your own words. Mark each source.

#### Advanced practical tasks

1. Write the full comparison page (table plus recommendation). Include failover and pause.
2. Read a second source (Forrester, a vendor, or a conference talk notes). Mark it opinion or spec. Update your recommendation.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do the five next-step projects map to topics 10, 12, 13, and 14?
2. Which project proves Redis is optional as a cache, and which project proves a second process copies data?
3. Why do official command pages still matter after you finish this handbook?
4. What is the difference between finishing a drill and owning a production policy (topic 21)?
5. What do you still not know after these projects (Cluster at scale, Kafka, modules in prod)?

#### Easy practical tasks

1. Write a checklist of the five projects with a date column.
2. Draw a learning map: official docs, your cache layer, limiter, replica, lock note.
3. List three URLs you will keep (docs, commands, data types).
4. Write five lab safety rules (localhost, no `KEYS` in prod, no public `6379`, lab `FLUSHALL` only, two Redis roles).

#### Medium practical tasks

1. Put the cache layer and the limiter in one small HTTP demo. Record keys after 10 requests.
2. Run the replica drill and the lock tests on the same day. Write one page of results.
3. Compare your cheat sheet from the first section with topic 21 policies. Add the missing policy lines.

#### Advanced practical tasks

1. Publish an internal demo repo (or a local folder) with README, tests, and the lock comparison. No production secrets.
2. Plan the next quarter: Cluster lab, one module, a backup restore drill (topic 19), and a game-day breaker (topic 21). Write owners and dates.
