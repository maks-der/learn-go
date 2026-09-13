# 12. Modules and Production Use

## Description

A module adds commands and types that core Redis does not ship. Production Redis is a set of policies, not only a running process. This topic surveys RedisJSON, RediSearch, TimeSeries, and Bloom, then covers separate cache and durable instances, key budgets, tenant prefixes, circuit breaking, and forks such as Valkey.

Complete topic 5 (module survey), topic 7 (persistence), topic 9 (clients and cache patterns), and topic 11 (security) first.

Use one term for each concept. A module is a server plugin that registers commands. Redis Stack is Redis plus a bundle of modules. A cache Redis is an instance that you may lose. A durable Redis is an instance that is a source of truth or a store you must restore. A key budget is a limit on count or bytes. A tenant prefix is a string at the start of a key. A circuit breaker stops calls to Redis after errors. An alternative implementation is a different server that speaks Redis protocols or commands.

---

## RedisJSON, RediSearch, TimeSeries, Bloom

RedisJSON stores JSON documents as a first-class value. You address paths. You do not rewrite the whole document for a small change.

```text
JSON.SET lab:user:1 $ '{"name":"Ada","n":1}'
JSON.GET lab:user:1 $.name
JSON.NUMINCRBY lab:user:1 $.n 1
```

`$` is the document root in JSONPath. Path syntax is on the RedisJSON command pages. Check your module version.

Why not `SET` a JSON string? A string `SET` replaces the whole blob. Concurrent updates of two fields need a lock or a Lua read-modify-write. RedisJSON can update one path.

`TYPE` on a JSON key is not `string`. Commands such as `GET` do not apply. Clients need RedisJSON support or raw command calls.

Use RedisJSON when documents are nested, updates are partial, and you already run Redis Stack or a hosted JSON feature. Use a hash when the model is a flat map of small fields.

RediSearch (Search) builds an index on hashes or JSON. You query the index with `FT.SEARCH`. This is not SQL. You must create a schema.

```text
FT.CREATE idx:lab ON HASH PREFIX 1 lab:doc: SCHEMA title TEXT body TEXT score NUMERIC
FT.SEARCH idx:lab "hello"
```

`FT.CREATE` defines fields and types (`TEXT`, `TAG`, `NUMERIC`, `GEO`, and others by version). Documents are ordinary Redis keys that match the prefix. The module maintains inverted indexes.

Search uses RAM for the index. Large text corpora cost memory beyond the documents. Aggregations (`FT.AGGREGATE`) exist. They are not a warehouse.

Use Search when you need full-text or filter queries that hashes and `SCAN` cannot do. Do not use Search as a substitute for PostgreSQL when you need joins, constraints, and ad-hoc SQL.

RedisTimeSeries stores timestamped samples. Each key is one series.

```text
TS.CREATE lab:cpu
TS.ADD lab:cpu * 0.42
TS.RANGE lab:cpu - +
```

`TS.CREATERULE` downsamples into another series (avg, min, max, and others). Retention (`RETENTION`) deletes old samples.

Why not a sorted set? A sorted set can store `score=timestamp` and `member=value`. At large volume, TimeSeries is more compact and has aggregation rules. Topic 4 windows are fine for small cases.

Labels let you query many series (`TS.MRANGE`). High-cardinality labels (one series per user id at huge scale) explode memory.

Use TimeSeries for infrastructure or product metrics that you already want in Redis. Use a dedicated metrics system when retention and query volume exceed a memory store.

RedisBloom adds probabilistic structures: Bloom filters, Cuckoo filters, Count-Min sketches, Top-K, and related types.

A Bloom filter answers "definitely not in the set" or "possibly in the set." It does not store the items. False positives exist. False negatives do not (in the standard Bloom filter).

```text
BF.RESERVE lab:bf 0.01 10000
BF.ADD lab:bf user:42
BF.EXISTS lab:bf user:42
```

`0.01` is a target error rate. `10000` is a capacity hint. After too many inserts, the error rate rises unless you reserved enough capacity.

HyperLogLog (topic 5) estimates unique counts. A Bloom filter tests membership. Do not mix the two.

Do not use a Bloom filter as an access-control list. A false positive can mean "possibly allowed." That is the wrong tool for security.

Build it yourself (core types + application) when a hash, stream, or sorted set already matches the access pattern, or when you must run vanilla Redis. Use a module when the module implements a known algorithm well and your platform already ships it.

A module command still runs on the main thread unless the module docs say otherwise. A huge `FT.SEARCH` can stall.

Licensing and vendor lock-in are part of the choice. Read the current license before you ship. Do not load an unmaintained module on a production primary.

A module changes the operations story. Redis version, module version, Redis Stack image tag, and client library must match. A replica without the module cannot load the data. Test `BGSAVE`, restart, and replica sync with module keys. Pin an image digest in production. Do not run `latest` without a pin.

`MODULE LOAD` is a dangerous command (topic 11). Load only signed or vendor-supported modules.

Use Redis Stack in Docker for learning:

```text
docker run --name stack-learn -p 6379:6379 redis/redis-stack
```

### Questions

#### Theoretical questions

1. What problem does RedisJSON solve that a JSON string does not solve well?
2. What must you create before `FT.SEARCH` works?
3. What does one TimeSeries key represent?
4. What two answers can a Bloom filter give?
5. Why must a replica load the same module?

#### Easy practical tasks

1. If Stack is available, run the three RedisJSON commands in this section. Save the replies. If not, write the Docker image name `redis/redis-stack` and a `docker run` line.
2. Write a two-column table: module and one command prefix (`JSON.`, `FT.`, `TS.`, `BF.`).
3. Run `INFO modules` or `MODULE LIST`. Write the module names (empty is a valid result).
4. Write four sentences: RedisJSON vs `SET` of JSON text.

#### Medium practical tasks

1. If Stack is available, create an index on three hashes. `FT.SEARCH` a word that appears in one title. Save the reply. Drop the index in the lab.
2. Compare `MEMORY USAGE` of 1000 `TS.ADD` samples vs 1000 `ZADD` members (same numbers). Write both.
3. Insert 1000 items in a Bloom filter if present. Test 100 missing items. Count false positives.

#### Advanced practical tasks

1. Write a decision page: core type vs module for documents, search, metrics, and unique counts (HLL vs Bloom vs set).
2. Write an upgrade plan: Stack image A to B, replica first, module compatibility check, rollback. Include a restore drill.

---

## Separate cache Redis from durable Redis

One instance that is both a huge cache and the only copy of jobs or sessions creates a conflict.

A cache wants:

- Eviction (`allkeys-lru` or similar)
- Short TTLs
- Optional persistence or none
- Tolerance for flush and restart

A durable store wants:

- `noeviction` or careful `volatile-*` only
- Persistence tested (topic 7)
- Replication and backup (topic 8, topic 11)
- No `FLUSHALL` as a deploy step

If both live on one process, a cache stampede or a `maxmemory` policy can evict durable keys. A `FLUSHALL` wipes both. A `KEYS` stall hurts both.

Pattern: two instances (or two vendor databases). Different hosts, different `maxmemory`, different ACLs, different backups. The application uses two clients.

A stream that must not drop (topic 10) belongs on the durable instance, or on Kafka. A product page cache belongs on the cache instance.

Cost: two pools, two dashboards, two upgrades. That cost is lower than one mixed incident.

Do not use logical `SELECT` to split cache and durable data (topic 1, topic 8 Cluster).

Name the instances in config: `REDIS_CACHE_URL` and `REDIS_DURABLE_URL`. Do not reuse one URL for both roles.

Warming fills cache keys before they are hot. A job reads the source of truth and `SET`s cache keys with TTL and jitter. Pace the job. Prefer versioned keys (`cache:p:1001:v2`) and warm `v2` before you switch. Do not `FLUSHALL` to "refresh the cache".

### Questions

#### Theoretical questions

1. What does a cache instance accept that a durable instance must not accept?
2. What can eviction do if cache and durable keys share `maxmemory`?
3. Why is `SELECT` a poor split?
4. Where do you put a stream that must not drop?
5. Why do you use two client objects?

#### Easy practical tasks

1. Write a two-column table: cache Redis vs durable Redis. Add eviction, persistence, backup, flush.
2. Write four sentences: two instances vs one mixed instance.
3. Write two environment variable names for the two URLs.
4. Assign six key prefixes to cache or durable (`cache:`, `sess:`, `lock:`, `jobs:`, `flags:`, `rl:`). Give one reason each.

#### Medium practical tasks

1. Run two Docker Redis ports (`6379`, `6380`). Point a tiny app at both. Kill the cache container. Confirm the durable `GET` still works.
2. Set `maxmemory` small on a mixed lab instance with one durable key and many cache keys. Write whether the durable key survived. Then explain why production must not mix.
3. Write ACL users: cache app vs durable app. Different key patterns.

#### Advanced practical tasks

1. Write a production topology: cache primary+replica (or none), durable primary+replica+backup. Include memory sizes.
2. Review an existing app that uses one Redis for everything. Write a split plan in ten steps.

---

## Key budgets and multi-tenant prefixes

A key budget is a written limit: maximum keys per prefix, maximum bytes per key, default TTL, jitter, and owner.

Without a policy, every feature adds keys forever. Redis grows until eviction or OOM.

Policy rows look like:

```text
prefix     role      ttl     jitter   max bytes   owner
cache:p:   cache     300s    0-60s    16KB        shop
sess:      session   1800s   0        4KB         iam
neg:       negative  30s     0-5s     16B         shop
```

Budgets use topic 9 (TTL, jitter, negative) and topic 11 (big keys). Ops alert when `INFO keyspace` or a `SCAN` estimate exceeds the budget.

`maxmemory` is the hard stop. The policy is the social stop. You hit the policy first.

Durable keys still need a budget. A leak of lock keys without TTL fills a `noeviction` instance until writes fail.

Every new prefix needs a review: TTL yes/no, who deletes, what happens on deploy. Do not invent a prefix in a pull request without a catalog line.

Measure: `MEMORY USAGE` samples, `--bigkeys` on a replica, key count per prefix from a paced `SCAN` in a job, not `KEYS`.

Several tenants (customers, shops, teams) can share one Redis. You prefix keys:

```text
t:acme:cache:product:1
t:globex:cache:product:1
```

The prefix avoids collisions. Tenants still share RAM, CPU on the main thread, `connected_clients`, a `FLUSHALL` if someone can run it, network, and disk.

A prefix is not a security boundary (topic 1, topic 11). A client that can `GET t:acme:*` might still `GET t:globex:...` if ACLs do not restrict key patterns.

ACLs can limit a tenant user to `~t:acme:*`. That is a real control. The application user that serves all tenants must not have that isolation unless you use per-tenant Redis users (rare) or per-tenant instances (stronger).

Noisy neighbor: one tenant can fill memory or run a big `HGETALL`. Budgets per tenant prefix, rate limits (topic 10), and separate instances for large tenants fix this.

Cluster hash tags: `{t:acme}` can pin one tenant to a slot range if you design that way. It can also hot-spot one slot. Prefer tags for keys that must live together, not for all tenant keys unless you measured it.

Do not put tenant id only in the value and use the same key for all tenants.

### Questions

#### Theoretical questions

1. What fields belong in a key-policy row?
2. How does a budget differ from `maxmemory`?
3. What does a tenant prefix prevent?
4. What does a tenant prefix not isolate?
5. What is a noisy neighbor on Redis?

#### Easy practical tasks

1. Write a five-row policy table for a shop.
2. Write keys for two tenants and the same product id.
3. Write four sentences: prefix vs separate instance.
4. Write an ACL key pattern for tenant `acme` only.

#### Medium practical tasks

1. Estimate bytes: 1 million session hashes at 500 bytes each. Write RAM and a TTL that caps live sessions.
2. Find prefixes on a lab instance with `SCAN`. Mark which lack a TTL (`TTL` = `-1`).
3. Design a per-tenant budget: max keys, max ops/s. Write how you enforce (limits in Redis, topic 10).

#### Advanced practical tasks

1. Write a full catalog for an app you know (or a fictional shop) with 12 prefixes and two instances.
2. Write a multi-tenant Redis standard: prefix format, ACL, budgets, Cluster tags, encryption/backup.

---

## Circuit breaking when the cache is down

If Redis is a cache, a Redis outage must not take down the product. The application reads the source of truth and skips the cache.

A circuit breaker:

- Counts errors or timeouts to Redis
- Opens after a threshold
- Fails fast (does not wait a long timeout on every request)
- Half-opens to try a probe (`PING` or one `GET`)
- Closes when Redis is healthy

Topic 9 timeouts matter. A 5 second timeout with no breaker can exhaust the application thread pool.

If Redis is durable (sessions you cannot rebuild, a job stream), the breaker does not invent data. The application returns an error or a degraded mode that you designed. That is a different product decision. Do not pretend a durable store is optional.

Cache-aside (topic 9) fits a breaker: treat Redis errors as misses.

Logs: when the breaker opens, log once per interval, not once per request.

Metrics: breaker state, Redis error rate, SQL load (SQL will rise when the cache is out).

Do not retry forever inside the breaker window (topic 9).

Test: stop Redis in a lab and confirm the HTTP API still serves cached-optional routes.

Dual-write keeps two stores in sync during a migration or a cache fill. The application writes SQL and Redis on the same request, or writes two Redis clusters. If Redis fails, decide whether the request fails (durable) or continues (cache). A dual-write can diverge. You need a repair job and a TTL so that divergence dies.

Cold start: an empty cache after restart causes a stampede (topic 9). Warm the top-N keys, or use single-flight, or accept SQL load.

### Questions

#### Theoretical questions

1. What must a cache-aside app do when Redis is down?
2. What does an open circuit do to the next request?
3. Why do timeouts without a breaker exhaust threads?
4. Why can a durable Redis outage still fail the user request?
5. Why do you log breaker opens with a rate limit?

#### Easy practical tasks

1. Write four sentences: miss vs Redis error vs open circuit.
2. Draw closed, open, half-open.
3. List three routes that must survive cache Redis down, and one that cannot (if you have a durable Redis).
4. Write a probe command (`PING`).

#### Medium practical tasks

1. Implement a small breaker around `GET`: after 3 errors, skip Redis for 10 seconds, then probe. Test with Redis stopped.
2. Compare p99 of 100 requests with timeout 2s and no breaker vs with a breaker (Redis down). Write both.
3. Document SQL load expectation when the cache breaker opens.

#### Advanced practical tasks

1. Integrate a known breaker library in your language. Write settings: error percent, sleep window, fallback.
2. Write a game-day script: stop cache Redis, watch breaker and error rate, start Redis, watch close. No production.

---

## Valkey and other forks (awareness)

The Redis protocol and command set inspired other servers.

Valkey is a Linux Foundation project that forked from Redis after a license change in the Redis project. Many commands stay familiar. It is an alternative you may see in vendors and distros. Read Valkey docs for Cluster, modules, and license. Do not assume every Redis module loads.

KeyDB is a fork that emphasizes multi-threaded command execution. The concurrency model is not the Redis main-thread model. Race assumptions that are true on Redis (`INCR` atomic on one thread) need a careful read of KeyDB docs. Modules and Cluster behavior can differ.

Other proxies and compatible caches exist (vendor caches, `redis-stack` images, cloud engines). Compatibility is a spectrum: common `GET`/`SET` work; `MODULE`, `FUNCTION`, and exact `INFO` fields may not.

Awareness rules:

- Speak the server name in runbooks ("Valkey 8" not "Redis" if it is Valkey)
- Test the client against that server
- Read the license for your use
- Do not mix replicas of different engines unless a vendor supports it

This handbook remains a Redis learning path. Commands here target Redis. When you use an alternative, verify each command on that engine.

Do not treat a fork as Redis for compliance or for support contracts without a check.

Official docs remain [https://redis.io/docs/](https://redis.io/docs/). Redis University is [https://university.redis.io/](https://university.redis.io/). The suggested practice list in `redis.topics.md` is a second checklist.

### Questions

#### Theoretical questions

1. What is Valkey in one sentence?
2. What design emphasis is KeyDB known for?
3. Why might a Redis module fail on a fork?
4. Why do you name the engine in a runbook?
5. Why is this handbook still written for Redis?

#### Easy practical tasks

1. Open the Valkey site. Write the project one-line description in your own words.
2. Open a KeyDB page. Write one sentence about threads.
3. Write four sentences: protocol compatible vs fully compatible.
4. Write the license names you see today for Redis and Valkey (date your note).

#### Medium practical tasks

1. If you can run a Valkey container, `PING` and `INFO server`. Write the server name field. If not, write a `docker` image name from the docs.
2. Compare one feature (Functions, ACL, Cluster) in Redis docs vs Valkey docs. Write match or differ.
3. Write a client test list: 10 commands you would run before you switch an app.

#### Advanced practical tasks

1. Write a migration note: Redis to Valkey (or the reverse) for a cache instance. Include license, image, modules, and rollback.
2. Read a reliable comparison of KeyDB threading vs Redis I/O threads. Write six sentences on what stays atomic and what you must re-read.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do RedisJSON, Search, TimeSeries, and Bloom map to "document", "query", "metric", and "maybe-seen"?
2. How do two instances, a key catalog, and a circuit breaker form one production story?
3. When is a tenant prefix plus ACL enough, and when do you split instances?
4. What must stay true if Redis is optional versus if Redis is required?
5. What do you tell a teammate who says "KeyDB is Redis with more cores, so our Lua races go away"?

#### Easy practical tasks

1. Write a cheat sheet: `JSON.`, `FT.`, `TS.`, `BF.`, two URLs, catalog row, tenant prefix, breaker states, Valkey, KeyDB.
2. Draw users → app → cache Redis and SQL, plus durable Redis on the side.
3. Check your instance for modules. Write present or absent for each of the four.
4. Write five production rules: pin images, two instances, no `FLUSHALL`, breaker on cache, name the engine.

#### Medium practical tasks

1. If Stack is available, run one command from each module. Save replies. If not, write a compose file that starts Stack.
2. Write a one-page production Redis policy that includes split instances, budgets, tenants, breaker, and allowed modules.
3. Game-day: stop cache Redis behind a breaker; confirm durable Redis still serves one key. Write the timeline.

#### Advanced practical tasks

1. Write a production module and fork policy: allowed modules, owners, license review, upgrade cadence, restore drill, when Valkey is acceptable.
2. Write an RFC: split a mixed Redis, migrate with dual-write and warm, add a breaker, set budgets. Include rollback.
