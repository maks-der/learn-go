# 13. Caching Patterns

## Description

A cache stores a result so that a later read avoids a slow path. This topic covers how the application reads and writes that cache, how many clients can miss at once, and how TTL design reduces stampedes.

Complete topic 4 (expiry and eviction) and topic 12 (clients) first. This topic assumes Redis is the cache, not the only source of truth, unless a task says otherwise.

Use one term for each concept. The source of truth is the durable store (often a SQL or document database). A hit is a read that finds a usable value in Redis. A miss is a read that does not. A stampede is many misses for the same key at the same time. Negative caching stores a "missing" marker. Jitter is a random extra TTL.

---

## Cache-aside

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

If Redis is down, the application must read the source of truth and skip the cache (topic 21 circuit breaking). Cache-aside makes that path natural: a miss and a Redis error can both fall back to SQL.

Do not treat a cache hit as proof that the source of truth still has that row. The cache can be stale until TTL or invalidation.

### Questions

#### Theoretical questions

1. Who reads the source of truth in cache-aside?
2. What are the five steps of a miss?
3. Why does Redis not call SQL in this pattern?
4. What is one race after `DEL` plus a late `SET`?
5. What should the application do when Redis does not answer?

#### Easy practical tasks

1. Write the five steps as a numbered list in your own words.
2. Draw application, Redis, and SQL. Label hit and miss arrows.
3. Write four sentences: `DEL` after write vs `SET` after write.
4. Name a cache key and a TTL for a product page.

#### Medium practical tasks

1. Implement cache-aside for one key in your language: miss loads a file or a mock DB. Show hit on the second read.
2. Update the mock DB. `DEL` the cache key. Confirm the next read sees the new value.
3. Stop Redis. Confirm the application still returns the mock DB value.

#### Advanced practical tasks

1. Demonstrate the stale-`SET` race with two processes (or two sleeps). Write the timeline and one mitigation (TTL, version, or lock).
2. Write a one-page cache-aside guide: key names, TTL, invalidation, and the Redis-down path.

---

## Read-through / write-through (app-level)

Read-through and write-through move the load and store steps into a cache module that the application calls as if it were the database.

Read-through: the application asks the module for a key. The module `GET`s Redis. On a miss, the module reads the source of truth, `SET`s Redis, and returns the value. The handler does not contain the miss branch.

Write-through: the application writes through the module. The module writes the source of truth and writes Redis in the same application request (two stores, one call). The user request waits for both.

In Redis you still implement this in the application or in a library. Redis does not load SQL by itself. Some products add a proxy or a module. This handbook treats the pattern as application-level.

Read-through hides the miss path. That can be cleaner. It can also hide cost: a "simple get" may run a slow query.

Write-through keeps Redis and the source of truth closer at write time. The request is slower than a write that only touches SQL. If Redis fails, you must decide: fail the request, or write SQL and accept a cache miss later.

Do not confuse write-through with a Redis-only database. The source of truth remains the durable store unless you design otherwise.

Cache-aside and read-through can store the same bytes. The difference is who owns the miss and write functions.

### Questions

#### Theoretical questions

1. What does the handler call in read-through?
2. Who talks to SQL on a read-through miss?
3. What two stores does write-through update in one request?
4. Why can read-through hide latency?
5. How is write-through different from "Redis is the database"?

#### Easy practical tasks

1. Write four sentences: cache-aside vs read-through.
2. Draw write-through: client, module, SQL, Redis.
3. Make a two-column table: "Pattern" and "Who contains the miss branch?".
4. List two failure choices when Redis is down during write-through.

#### Medium practical tasks

1. Wrap cache-aside in a `GetProduct(id)` function so the handler never mentions Redis. That is read-through at app level. Show one caller.
2. Implement write-through for one field: write SQL (mock) and `SET` Redis. Confirm a following `GET`.
3. Fail the Redis `SET` after a successful mock SQL write. Document the state and the next read.

#### Advanced practical tasks

1. Compare request latency: SQL only, write-through (SQL+Redis), cache-aside write (`DEL` only). Write three numbers from a lab.
2. Write an interface with `Get` and `Put`. Provide a read-through / write-through implementation and a cache-aside implementation. List what the tests assert.

---

## Write-behind

Write-behind (write-back) writes Redis first (or only Redis in the request). A background worker writes the source of truth later.

The user request is fast. The durable store lags. If Redis loses the key before the worker runs (eviction, restart without persistence, a bug), the durable store never sees the write.

Use write-behind only when you accept that risk, or when Redis is persistent and you have a replay queue (streams, topic 15) that the worker consumes.

A typical shape:

1. The application writes the cache or a stream entry.
2. The application returns success to the user (if the product allows that).
3. A worker reads the queue and writes SQL.
4. The worker acknowledges the entry.

This is closer to a broker pattern than to a simple cache.

Do not use write-behind for money or legal records unless a durable log exists outside a volatile cache.

Cache-aside plus `DEL` is simpler when the source of truth must win every time.

### Questions

#### Theoretical questions

1. Which store does the user request write first in write-behind?
2. What happens if Redis drops the key before the worker runs?
3. What extra component does a safe write-behind design need?
4. When is write-behind a poor fit?
5. How does write-behind differ from write-through on request latency?

#### Easy practical tasks

1. Write four sentences: write-through vs write-behind.
2. Draw request path and worker path as two timelines.
3. List three data types that must not use volatile write-behind.
4. Name one Redis type that can hold the pending write (stream or list).

#### Medium practical tasks

1. `XADD` a change, then a second script reads and "writes SQL" (print or file). Document the delay.
2. Restart Redis without persistence after `SET` but before the worker. Record the lost update.
3. Write a table: pattern, user wait, loss risk, complexity. Add cache-aside, write-through, write-behind.

#### Advanced practical tasks

1. Design write-behind with a stream, a consumer group, and SQL. Write failure cases: worker crash, Redis restart, SQL error.
2. Read a public post-mortem about a lost write-behind cache. Write five sentences on the failure and the fix.

---

## Stampede / thundering herd

A stampede (thundering herd) happens when many clients miss the same key at once. Each client loads the source of truth. The database receives a burst. The cache then receives a burst of `SET`s.

Typical causes:

- The key expires (`TTL` reaches 0) while traffic is high
- Eviction removes a hot key
- A deploy flushes the cache
- The process starts with an empty cache (cold start)

Effects: high SQL load, high Redis CPU, slower user requests.

Mitigations (later sections add more):

- Single-flight: only one loader per key; other callers wait for that loader
- A lock (`SET NX EX`) around the load (topic 14)
- Soft TTL plus early refresh (next section)
- TTL jitter so that many keys do not expire in the same second
- Warm the cache before a deploy (topic 21)
- Negative caching so that a missing row does not hit SQL on every request

A stampede is a coordination problem. Redis can hold a lock or a "loading" flag. The application must use it.

Do not flush all keys on deploy if the dataset is large and hot. Prefer versioned keys or gradual warming.

### Questions

#### Theoretical questions

1. What is a stampede?
2. Name three causes of a sudden miss storm.
3. What does single-flight mean?
4. How can a lock reduce a stampede?
5. Why can a full flush at deploy time cause a stampede?

#### Easy practical tasks

1. Write four sentences: miss vs stampede.
2. Draw 20 clients, one expired key, and one SQL box with 20 arrows.
3. List five mitigations from this section.
4. Write a lock key name for `cache:product:1001`.

#### Medium practical tasks

1. Start five parallel miss loaders without a lock (sleep in the loader). Count how many times the mock DB runs.
2. Repeat with `SET loader:product:1001 NX EX 10`. Count mock DB runs again.
3. Expire 100 keys at the same second (same TTL). Watch a mock DB counter. Then add jitter (next section) and compare.

#### Advanced practical tasks

1. Implement single-flight in your language (mutex map or `singleflight` style). Prove one load for ten waiters.
2. Write a stampede runbook: metrics to watch (SQL QPS, Redis `instantaneous_ops_per_sec`, miss rate) and first actions.

---

## Probabilistic early expiration

A fixed TTL expires every copy of a hot key at one moment. Probabilistic early expiration refreshes the key before that moment, with a probability that rises as the remaining TTL shrinks.

Idea:

- Store the logical expiry time with the value (or use `PTTL`)
- On each hit, compute a probability from remaining time and the cost of a reload
- Sometimes reload from the source of truth and `SET` again, even though the key still exists

One client refreshes early. Other clients still hit. The key does not hit zero while all clients wait on SQL.

This is an application algorithm. Redis does not compute the probability. Papers and libraries describe variants (xfetch, perishable cache).

You still need a TTL. Early refresh is extra. If no client reads the key, it expires. That is fine for a cold key.

Cost: some extra SQL reads on a hot key. That cost is usually smaller than a stampede.

Do not early-refresh every hit. That removes the cache benefit.

Combine with a lock or single-flight so that two early refreshes do not both hit SQL.

### Questions

#### Theoretical questions

1. What problem does early expiration reduce?
2. Who computes the probability?
3. Why do you still set a TTL?
4. Why must you not refresh on every hit?
5. What extra coordination stops two early refreshes?

#### Easy practical tasks

1. Write four sentences: hard expiry vs probabilistic early refresh.
2. Draw a TTL timeline. Mark a window where refresh probability is high.
3. List what you store besides the payload (logical expiry, version).
4. Open a description of xfetch or a similar algorithm. Write one formula or rule in your own words.

#### Medium practical tasks

1. Implement a simple rule: if `PTTL` is below 10 percent of the original TTL, reload with probability 0.1. Log when you reload.
2. Hit a key in a tight loop. Count reloads vs hits. Adjust the threshold.
3. Compare a stampede at hard expiry vs early refresh under a synthetic load. Write both SQL counts.

#### Advanced practical tasks

1. Read one paper or article on probabilistic cache refresh. Write the parameters and one limit of the model.
2. Add single-flight to the early refresh path. Show that ten concurrent hits cause one reload.

---

## Negative caching

A miss can mean "the source of truth has no row." If you store nothing, every request for a bad id hits SQL.

Negative caching stores a marker for a short TTL:

```text
SET cache:product:9999:neg 1 EX 30
```

Or you store an empty JSON object with a type field `{ "missing": true }`. The application must distinguish "missing" from "product with empty name".

Use a short TTL. A negative entry that lives too long hides a row that you just inserted.

On a successful create in SQL, delete the negative key. Otherwise the next read can still show "missing".

Abuse: attackers can request many random ids. Negative caching protects SQL. It also writes many keys. Use a TTL, a prefix budget (topic 21), and rate limits (topic 14).

Do not negative-cache errors such as "SQL is down." That marker would hide recovery. Cache "not found", not "cannot talk to SQL".

### Questions

#### Theoretical questions

1. What problem does negative caching solve?
2. Why is the TTL usually short?
3. What must a create path do to a negative key?
4. Why must you not cache "SQL is down"?
5. How can negative caching increase Redis memory under abuse?

#### Easy practical tasks

1. Write a negative key name and a 30 second `SET`.
2. Write four sentences: miss of a missing row vs miss of an expired product.
3. Draw read path with a `neg` key hit (no SQL).
4. List two encodings for a negative marker.

#### Medium practical tasks

1. Implement: missing id writes `neg`. Second read does not call the mock DB. Then insert the row and `DEL` the `neg` key.
2. Leave a `neg` key in place. Insert the row. Read without `DEL`. Record the stale "missing" result.
3. Count Redis keys after 1000 random missing ids with TTL 30. Write the peak key count.

#### Advanced practical tasks

1. Design negative caching plus a rate limit on unknown ids. Write keys and TTLs.
2. Write test cases: not found, found, create after negative, SQL error (must not write `neg`).

---

## TTL jitter

TTL jitter adds a random extra lifetime to each key. If 10,000 product keys all get `EX 300`, they can expire in the same second after a bulk load. The next second is a stampede.

```text
# conceptual: 300 seconds plus a random 0–60
SET cache:product:1001 <bytes> EX 347
```

Each key expires at a different time. Reloads spread across a minute.

Jitter is not a lock. Jitter does not help one hot key that all clients share. Jitter helps many keys that share the same TTL and the same load time.

Pick a jitter range from the base TTL. A 300 second TTL plus 0–30 seconds is common. A jitter larger than the base TTL makes expiry hard to reason about.

Document the policy: `ttl = base + random(0, jitter)`. Use a clock that is good enough for seconds. Do not use a tiny jitter (0–1 second) on a huge bulk load.

Eviction can still remove many keys at once if memory is full. Jitter does not replace `maxmemory` planning (topic 4).

### Questions

#### Theoretical questions

1. What is TTL jitter?
2. Which stampede shape does jitter reduce?
3. Which stampede shape does jitter not solve?
4. Why is a jitter larger than the base TTL a problem?
5. Why does jitter not replace eviction planning?

#### Easy practical tasks

1. Write a formula `ttl = base + random(0, jitter)` with example numbers.
2. Write four sentences: jitter vs one lock on one hot key.
3. `SET` three keys with `EX` 20, 25, and 30. Watch `TTL` values.
4. List two events that still expire many keys together (flush, deploy).

#### Medium practical tasks

1. Load 200 keys with the same `EX 15`. Sample how many remain each second. Repeat with jitter 0–10. Compare the expiry curve.
2. Add jitter to your cache-aside `SET`. Print the TTL that you send.
3. Write a policy table: key prefix, base TTL, jitter, negative TTL.

#### Advanced practical tasks

1. Combine jitter, negative caching, and single-flight in one small cache layer. Write a test plan with five cases.
2. Model 50,000 keys loaded at once. Estimate peak SQL QPS with no jitter vs jitter of 60 seconds (back-of-envelope).

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Which pattern keeps Redis optional on every read path, and why does that matter when Redis is down?
2. How do stampede, jitter, early refresh, and negative caching work as one design, not as four isolated tricks?
3. When do you choose write-through instead of cache-aside `DEL`?
4. What data must not use write-behind into a volatile Redis?
5. How does a write race differ from a read stampede?

#### Easy practical tasks

1. Write a cheat sheet: seven pattern names from this topic, one sentence each.
2. Draw one product-read flowchart that includes hit, miss, negative hit, and Redis error.
3. Assign a pattern to session data, product pages, and payment capture. Write one reason each.
4. List the Redis commands this topic used (`GET`, `SET`, `DEL`, `PTTL`, `SET NX`).

#### Medium practical tasks

1. Build cache-aside with jitter and negative markers for one mock entity. Record hits, misses, and SQL calls for 50 reads.
2. Write a decision table: pattern vs consistency need vs latency vs loss risk. Fill seven rows.
3. Document invalidation for update, delete, and create (including negative keys).

#### Advanced practical tasks

1. Load-test one hot key with and without single-flight. Write SQL call counts and p99 latency.
2. Write a production cache policy page: patterns allowed, TTL catalog, stampede controls, and a ban on `FLUSHALL` as a deploy step.
