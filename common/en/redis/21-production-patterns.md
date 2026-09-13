# 21. Production Patterns

## Description

Production Redis is a set of policies, not only a running process. This topic covers separate instances for cache and durable data, key budgets and TTLs, multi-tenant prefixes, circuit breaking when Redis is down, and how you warm or dual-write a cache.

Complete topics 9, 12, 13, and 17 first. Those topics give persistence, clients, cache patterns, and security. This topic combines them.

Use one term for each concept. A cache Redis is an instance that you may lose. A durable Redis is an instance that is a source of truth or a store you must restore. A key budget is a limit on count or bytes. A tenant prefix is a string at the start of a key. A circuit breaker stops calls to Redis after errors. Warming fills the cache before or during traffic. Dual-write updates two stores in one design.

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
- Persistence tested (topic 9)
- Replication and backup (topic 10, topic 19)
- No `FLUSHALL` as a deploy step

If both live on one process, a cache stampede or a `maxmemory` policy can evict durable keys. A `FLUSHALL` wipes both. A `KEYS` stall hurts both.

Pattern: two instances (or two vendor databases). Different hosts, different `maxmemory`, different ACLs, different backups. The application uses two clients.

A stream that must not drop (topic 15) belongs on the durable instance, or on Kafka. A product page cache belongs on the cache instance.

Cost: two pools, two dashboards, two upgrades. That cost is lower than one mixed incident.

Do not use logical `SELECT` to split cache and durable data (topic 1, topic 11 Cluster).

Name the instances in config: `REDIS_CACHE_URL` and `REDIS_DURABLE_URL`. Do not reuse one URL for both roles.

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

## Key budgets and TTLs as policy

A key budget is a written limit: maximum keys per prefix, maximum bytes per key, default TTL, jitter, and owner.

Without a policy, every feature adds keys forever. Redis grows until eviction or OOM.

Policy rows look like:

```text
prefix     role      ttl     jitter   max bytes   owner
cache:p:   cache     300s    0-60s    16KB        shop
sess:      session   1800s   0        4KB         iam
neg:       negative  30s     0-5s     16B         shop
```

Budgets use topic 13 (TTL, jitter, negative) and topic 16 (big keys). Ops alert when `INFO keyspace` or a `SCAN` estimate exceeds the budget.

`maxmemory` is the hard stop. The policy is the social stop. You hit the policy first.

Durable keys still need a budget. A leak of lock keys without TTL fills a `noeviction` instance until writes fail.

Every new prefix needs a review: TTL yes/no, who deletes, what happens on deploy.

Do not invent a prefix in a pull request without a catalog line.

Measure: `MEMORY USAGE` samples, `--bigkeys` on a replica, key count per prefix from a paced `SCAN` in a job, not `KEYS`.

### Questions

#### Theoretical questions

1. What fields belong in a key-policy row?
2. How does a budget differ from `maxmemory`?
3. Why do durable keys still need TTLs on locks?
4. When do you add a catalog line?
5. How do you count keys per prefix in production?

#### Easy practical tasks

1. Write a five-row policy table for a shop.
2. Write four sentences: policy vs `maxmemory`.
3. Add a missing owner to a prefix in your table.
4. Write one alert: `evicted_keys` on durable Redis.

#### Medium practical tasks

1. Estimate bytes: 1 million session hashes at 500 bytes each. Write RAM and a TTL that caps live sessions.
2. Find prefixes on a lab instance with `SCAN`. Mark which lack a TTL (`TTL` = `-1`).
3. Write a PR checklist: prefix, TTL, jitter, max size, owner, cache vs durable instance.

#### Advanced practical tasks

1. Build a weekly job spec: paced `SCAN`, estimate per prefix, compare to budget, ticket if over.
2. Write a full catalog for an app you know (or a fictional shop) with 12 prefixes and two instances.

---

## Multi-tenant key prefixes

Several tenants (customers, shops, teams) can share one Redis. You prefix keys:

```text
t:acme:cache:product:1
t:globex:cache:product:1
```

The prefix avoids collisions. Tenants still share:

- RAM and `maxmemory`
- CPU on the main thread
- `connected_clients`
- A `FLUSHALL` (if someone can run it)
- Network and disk

A prefix is not a security boundary (topic 2, topic 17). A client that can `GET t:acme:*` might still `GET t:globex:...` if ACLs do not restrict key patterns.

ACLs can limit a tenant user to `~t:acme:*`. That is a real control. The application user that serves all tenants must not have that isolation unless you use per-tenant Redis users (rare) or per-tenant instances (stronger).

Noisy neighbor: one tenant can fill memory or run a big `HGETALL`. Budgets per tenant prefix, rate limits (topic 14), and separate instances for large tenants fix this.

Cluster hash tags: `{t:acme}` can pin one tenant to a slot range if you design that way. It can also hot-spot one slot. Prefer tags for keys that must live together, not for all tenant keys unless you measured it.

Do not put tenant id only in the value and use the same key for all tenants.

### Questions

#### Theoretical questions

1. What does a tenant prefix prevent?
2. What does a tenant prefix not isolate?
3. How can ACLs make a prefix a real key restriction?
4. What is a noisy neighbor on Redis?
5. What is the risk of one hash tag for an entire tenant?

#### Easy practical tasks

1. Write keys for two tenants and the same product id.
2. Write four sentences: prefix vs separate instance.
3. Draw two tenants, one Redis, one shared `maxmemory` bar.
4. Write an ACL key pattern for tenant `acme` only.

#### Medium practical tasks

1. Implement a helper `key(tenant, parts...)` that rejects empty tenant ids. Show three keys.
2. Design a per-tenant budget: max keys, max ops/s. Write how you enforce (limits in Redis, topic 14).
3. Write when you move a tenant to its own instance (RAM, compliance, noisy neighbor).

#### Advanced practical tasks

1. Write a multi-tenant Redis standard: prefix format, ACL, budgets, Cluster tags, encryption/backup.
2. Compare shared Redis vs Redis-per-tenant for 500 small tenants. Write cost and isolation.

---

## Circuit breaking when Redis is down (app must still work if Redis is a cache)

If Redis is a cache, a Redis outage must not take down the product. The application reads the source of truth and skips the cache.

A circuit breaker:

- Counts errors or timeouts to Redis
- Opens after a threshold
- Fails fast (does not wait a long timeout on every request)
- Half-opens to try a probe (`PING` or one `GET`)
- Closes when Redis is healthy

Topic 12 timeouts matter. A 5 second timeout with no breaker can exhaust the application thread pool.

If Redis is durable (sessions you cannot rebuild, a job stream), the breaker does not invent data. The application returns an error or a degraded mode that you designed. That is a different product decision. Do not pretend a durable store is optional.

Cache-aside (topic 13) fits a breaker: treat Redis errors as misses.

Logs: when the breaker opens, log once per interval, not once per request.

Metrics: breaker state, Redis error rate, SQL load (SQL will rise when the cache is out).

Do not retry forever inside the breaker window (topic 12).

Test: stop Redis in a lab and confirm the HTTP API still serves cached-optional routes.

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

## Dual-write and warming

Warming fills cache keys before they are hot. Dual-write keeps two stores in sync during a migration or a cache fill.

Warming:

- A job reads the source of truth and `SET`s cache keys (with TTL and jitter)
- You warm after a flush, before a launch, or after a new region
- Pace the job so that Redis and SQL survive (pipelines with bounds, topic 16)

Dual-write:

- The application writes SQL and Redis on the same request (write-through, topic 13)
- Or writes two Redis clusters during a migration
- Failure handling: if Redis fails, decide whether the request fails (durable) or continues (cache)

A dual-write can diverge. You need a repair job (read SQL, `SET` Redis) and a TTL so that divergence dies.

Do not dual-write two durable systems without an idempotency story (topic 14) and a clear source of truth.

Cold start: an empty cache after restart causes a stampede (topic 13). Warm the top-N keys, or use single-flight, or accept SQL load.

Deploy: prefer versioned keys (`cache:p:1001:v2`) and warm `v2` before you switch. Then delete `v1` later. Do not `FLUSHALL` to "refresh the cache".

Read-repair: on a miss, fill Redis. That is cache-aside. Warming is the batch version.

### Questions

#### Theoretical questions

1. What does a warming job write?
2. What is dual-write in a cache migration?
3. Why can two writes diverge?
4. Why is `FLUSHALL` a poor deploy warm-up?
5. How do versioned keys help a deploy?

#### Easy practical tasks

1. Write four sentences: warming vs cache-aside miss fill.
2. Draw SQL, warm job, Redis, then users.
3. Write a top-N warm list idea (1000 product ids).
4. Write `v1` and `v2` key names for one product.

#### Medium practical tasks

1. Write a paced warmer: 100 keys, pipeline of 50, sleep 50 ms. Run on a lab mock.
2. Dual-write a mock SQL update and a Redis `SET`. Fail Redis once. Document the divergence and the repair `SET`.
3. Compare stampede after empty cache vs after a 1000-key warm. Use a mock SQL counter.

#### Advanced practical tasks

1. Design a region failover: warm the new cache from SQL, switch DNS or config, keep dual-write for 24 hours. Write checks.
2. Write a migration from instance A to B: dual-write, warm unread keys, switch reads, stop writes to A, backup A.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do two instances, a key catalog, and a circuit breaker form one production story?
2. When is a tenant prefix plus ACL enough, and when do you split instances?
3. Why do warming and jitter both exist in a launch plan?
4. What must stay true if Redis is optional versus if Redis is required?
5. How does a dual-write failure policy differ for cache Redis vs durable Redis?

#### Easy practical tasks

1. Write a cheat sheet: two URLs, catalog row, tenant prefix, breaker states, warm job, no `FLUSHALL`.
2. Draw users → app → cache Redis and SQL, plus durable Redis on the side.
3. Write five production rules for your team.
4. List the topic 13 patterns that this topic reused.

#### Medium practical tasks

1. Write a one-page production Redis policy that includes split instances, budgets, tenants, breaker, warming.
2. Game-day: stop cache Redis behind a breaker; confirm durable Redis still serves one key. Write the timeline.
3. Add tenant prefixes and budgets to the catalog from the key-budget section.

#### Advanced practical tasks

1. Review a real or sample architecture. Write gaps vs this topic (mixed instance, no breaker, no catalog).
2. Write an RFC: split a mixed Redis, migrate with dual-write and warm, add a breaker, set budgets. Include rollback.
