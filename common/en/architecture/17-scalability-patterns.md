# 17. Scalability Patterns

## Description

Scalability is the ability of the system to keep required latency and throughput when load grows. This topic covers vertical versus horizontal scale, load balancing, sharding, caching and invalidation, content delivery networks (CDN), asynchronous offload, and queue-based load leveling.

Topic 2 defined scalability. Topic 4 covered stateful versus stateless services. Topic 9 covered caches and ownership. This topic names the patterns and their costs. Complete Topics 1 to 16 before this topic. Design for the load that you can measure plus a stated growth factor. Do not start with a global shard plan for 30 users.

Use one term for each concept. Vertical scale grows one machine. Horizontal scale grows the number of machines. A shard is a partition of data. A cache is not the source of truth. A CDN is an edge cache for files. Offload and load leveling are related. They are not the same pattern.

---

## Vertical vs horizontal scale

Vertical scale adds resources to one machine: more CPU, more memory, or a faster disk. The process stays one process (or a small set on that machine). The design stays simple. There is a hard ceiling. One machine remains one fault domain (Topic 15).

Horizontal scale adds more machines or more processes that share the work. You need a method to split work: a load balancer, a partition key, or competing consumers. Stateless workers make horizontal scale easier (Topic 4). Sticky sessions make it harder.

Pick vertical first when:

- The load fits one modern host
- The team is small
- The state lives in one database that you already operate

Pick horizontal when:

- One host hits a measured ceiling
- You need more than one fault domain
- A stateless tier can clone cheaply

Hybrid is common. You scale the web tier horizontally. You scale the database vertically for a long time. Then you add read replicas or shards (later sections).

Elasticity is automatic add and remove of capacity. Elasticity needs horizontal units and a metric that drives the policy. Elasticity is not free. Cold starts and rebalancing have a cost.

Write the growth axis. "More users" and "more data per user" need different patterns. A CDN helps the first for static files. Sharding helps the second for a huge table.

### Questions

#### Theoretical questions

1. What is vertical scale?
2. What is horizontal scale?
3. When do you pick vertical scale first?
4. Why do sticky sessions make horizontal scale harder?
5. How does "more users" differ from "more data per user" as a growth axis?

#### Easy practical tasks

1. Write five sentences that compare the two scale types.
2. Make a table: "Tier" and "Usual first move". Add web process and primary database.
3. List four limits of vertical scale.
4. Draw one host versus three web processes and one database. Label new failure points (Topic 2).

#### Medium practical tasks

1. Write a one-page scale plan for a campus catalog that grows 10 times in users. State what stays vertical.
2. Explain in eight sentences how a stateful in-memory cart blocks a second web process.
3. Write an ADR: scale the API horizontally, keep one primary database. Include the measure that reopens the database choice.

#### Advanced practical tasks

1. Write a capacity note: current load, 10× load, first bottleneck, and the pattern you apply.
2. Compare scheduled vertical resize with automatic horizontal elasticity. Write operations cost for a four-person team.

---

## Load balancing

A load balancer is a component that spreads incoming connections or requests across more than one healthy instance. Clients see one name. Instances come and go.

Common methods:

- Round robin: next instance in a list
- Least connections: pick the instance with fewer open connections
- Hash of a key: the same key prefers the same instance (careful with hot keys)

Health checks remove a dead instance (Topic 15). A bad check is an outage.

Layer 4 balancing forwards bytes without reading HTTP. Layer 7 balancing can route on path or host. Layer 7 can terminate TLS. Each layer adds policy and cost.

Load balancing does not fix a hot shared store. If all instances hit one saturated database, you only move the queue (Topic 2).

Session affinity (stickiness) sends one client to one instance. Use it as a temporary crutch. Prefer a shared session store or stateless tokens (Topics 4 and 14).

The balancer is a fault domain. Use more than one balancer or a platform service with a stated SLA. Timeouts at the balancer must match application timeouts (Topic 11).

Do not put domain rules in the balancer. Routing by path is enough. Price calculation is not a balancer job (Topic 12, gateway warning).

### Questions

#### Theoretical questions

1. What is a load balancer?
2. Name three spread methods.
3. How do Layer 4 and Layer 7 differ at a high level?
4. Why does balancing not repair a saturated shared store?
5. Why is session affinity a crutch?

#### Easy practical tasks

1. Write four sentences that define load balancing.
2. Make a table: "Method" and "Risk". Add round robin, least connections, and hash.
3. List five duties that belong to a balancer and two that do not.
4. Draw clients, one balancer, and three API instances.

#### Medium practical tasks

1. Write a health-check and timeout policy for a balancer in front of an API that talks to a database.
2. Explain in eight sentences how a hash on `user_id` can overload one instance.
3. Compare a cloud load-balancer service with a self-operated proxy. Write operations cost.

#### Advanced practical tasks

1. Write a one-page balancer standard: TLS, drains, health, and header forwarding of a request identifier (Topic 16).
2. Design a fail-open versus fail-closed behavior when all health checks fail. Write the user impact.

---

## Sharding

Sharding is a split of data across more than one store by a key. Each shard holds a subset of rows. Horizontal scale of the data tier often needs shards when one primary cannot hold the writes or the working set.

A shard key is the field that decides the shard. Examples: `tenant_id`, `user_id`. A poor key puts most writes on one shard. That shard is a hotspot.

Cross-shard queries are expensive. A query that must see all users needs scatter-gather. Transactions across shards are not a local ACID transaction (Topic 7). Sagas appear (Topic 12).

Resharding (move of keys to new nodes) is an operations project. Plan it before you need it. Topic 9 ownership still applies: one owner per write path.

Sharding is not the same as a read replica. A replica copies the same data for reads. A shard holds different data.

Do not shard a campus app with one database that still has CPU and disk left. Sharding is a late pattern. It complicates backups, joins, and support.

Directory or hash mapping must be available. If the map is wrong, you read the wrong shard and you can write the wrong shard. Treat the map as critical configuration.

Legal location rules can force shards by region (Topic 19). That split is a constraint, not a fashion.

### Questions

#### Theoretical questions

1. What is sharding?
2. What is a shard key?
3. Why are cross-shard queries expensive?
4. How does a shard differ from a read replica?
5. Why is sharding a late pattern?

#### Easy practical tasks

1. Write five sentences that define sharding.
2. Make a table: "Key" and "Hotspot risk". Add `tenant_id` for one huge tenant and `user_id` for even users.
3. List four operations tasks that get harder after a shard.
4. Draw two shards and one API. Label the map.

#### Medium practical tasks

1. Pick a shard key for a multi-campus grade book. Write two queries that become hard.
2. Explain in eight sentences why a global unique ISBN catalog may not shard well on `isbn` if lookups are random and small.
3. Write an ADR that rejects sharding and accepts a larger vertical database plus a replica.

#### Advanced practical tasks

1. Write a one-page reshard plan: new nodes, dual writes or expand-contract (Topic 18), and rollback.
2. Compare tenant shards with hash shards. Write isolation, hotspot, and support cost.

---

## Caching and cache invalidation

A cache is a store of copies that you can read faster than the source of truth. Topic 9 stated that a cache is not the source of truth. Pair details with `redis.topics.md`.

Caching improves latency and reduces load on the origin. It also creates staleness. Invalidation is the hard part: when and how the copy dies.

Common invalidation styles:

- Time to live (TTL): the copy expires after a duration
- Write-through or write-around with an explicit delete on change
- Version keys: a new version makes old keys unused

Invalidation bugs show wrong prices, old stock, or old permissions. For security-sensitive data, prefer short TTL or no cache (Topic 14).

Stampede risk: many clients miss at once and hit the origin together. Mitigations include a lock on refill and staggered TTL. Stay at the idea level.

Cache keys must include everything that changes the result: user role, language, and version. A key that omits role can leak a copy across users. That is a security defect.

Measure hit ratio and origin load. A cache that no one hits is cost without benefit.

Do not cache a write path as if it were a read path. After a write, the next read must meet the consistency rule that the stakeholder accepted (Topic 7).

### Questions

#### Theoretical questions

1. What is a cache in this handbook?
2. Why is invalidation hard?
3. Name three invalidation styles.
4. Why must a cache key include the role when results depend on the role?
5. What is a stampede at a high level?

#### Easy practical tasks

1. Write four sentences that define caching and invalidation.
2. Make a table: "Item" and "Cache? (yes/no)". Add public catalog page, live stock count, and an access token.
3. List five fields that a catalog cache key might include.
4. Write a TTL choice for a campus news page and one for a grade that just changed.

#### Medium practical tasks

1. Design cache plus delete-on-write for `GET /books/{id}`. Write the fail path if delete fails.
2. Explain in eight sentences how a missing role in the key becomes an information-disclosure risk (Topic 14). Stay defensive.
3. Write metrics: hit ratio, origin QPS, and stale-serve count.

#### Advanced practical tasks

1. Write a one-page cache policy: what you cache, TTL, key recipe, and forbidden data.
2. Compare TTL-only with explicit delete. Write which one you pick for prices that change twice a day.

---

## CDN

A CDN (content delivery network) is a set of edge caches that serve files close to users. Typical files: images, scripts, style sheets, and some public pages.

The origin holds the source. The CDN holds copies at many points of presence. Users hit the nearest edge. Latency and origin load drop for cacheable files.

CDN rules are HTTP cache rules: `Cache-Control`, `ETag`, and versioned file names. A file named `app.a1b2.js` can have a long TTL. A file named `app.js` that changes in place is hard to cache safely.

A CDN is a poor place for personalized HTML that includes private data. You can cache a public shell. You must not cache another user inbox (Topic 14).

TLS and headers still matter at the edge. The CDN is a trust boundary. You configure who can purge. You configure which headers the edge forwards.

Cost includes egress and request counts. A misconfigured TTL can either stampede the origin or serve a broken script for a day.

Do not buy a global CDN for an intranet app on one campus unless you have a measured latency problem or a large static set.

Purge is an operations action. Record who can purge production.

### Questions

#### Theoretical questions

1. What is a CDN?
2. Which files fit a CDN?
3. Why do versioned file names help?
4. Why must you not cache a private inbox at the edge?
5. Why is the CDN a trust boundary?

#### Easy practical tasks

1. Write five sentences that define a CDN.
2. Make a table: "Asset" and "CDN? (yes/no)". Add a logo, a grade JSON, and a hashed script.
3. List four HTTP cache headers or fields that you will read about.
4. Write a purge rule in four sentences.

#### Medium practical tasks

1. Design cache headers for `/static/*` and for `/api/me`. Write why they differ.
2. Explain in eight sentences how a one-day TTL on `app.js` can keep a defect visible after a deploy.
3. Write a cost note: origin egress versus CDN bill for a large image set.

#### Advanced practical tasks

1. Write a one-page static-asset standard: hash names, TTL, and rollback of a bad file.
2. Compare a single-region object store with a CDN for a campus video. Write when the extra product is worth it.

---

## Async offload

Async offload moves work out of the user request. The API accepts the command, records the intent, and returns. A worker does the heavy work later (Topics 4, 10, and 12).

Offload fits:

- Mail and export files
- Image processing
- Search index updates
- Reports

The user request stays inside the latency budget (Topic 2). The worker absorbs the long tail.

You must show status. "Accepted" is not "done". Give a job identifier or a later notification. Users who think the work finished will file defects.

Offload needs the same reliability tools: outbox, idempotency, poison path, and SLIs on lag (Topics 10, 15, and 16).

Do not offload a step that the user must finish in the same click (payment capture in many shops). Do offload the receipt mail.

Offload is not horizontal scale by itself. It changes when the work happens. You still need enough workers. Queue-based load leveling (next section) handles bursts.

Write the maximum lag that stakeholders accept. That lag is a quality attribute.

### Questions

#### Theoretical questions

1. What is async offload?
2. Which jobs fit offload?
3. Why must the API tell the user that work is not done?
4. What reliability tools does offload still need?
5. When must you not offload a step?

#### Easy practical tasks

1. Write four sentences that define async offload.
2. Make a table: "Job" and "In request or offload". Add charge card, send receipt, and build PDF.
3. List five fields of a job status record.
4. Draw API, store, queue, and worker for a thumbnail job.

#### Medium practical tasks

1. Design offload for "export my loans as CSV". Write accept response, lag SLO, and fail message.
2. Explain in eight sentences how missing outbox creates a lost export (Topic 10).
3. Write RED-style metrics for the worker (Topic 16) plus a lag SLI (Topic 15).

#### Advanced practical tasks

1. Write a one-page offload standard: which writes stay synchronous and which writes may wait.
2. Compare polling status with a push notification. Write client and operations cost (Topic 20).

---

## Queue-based load leveling

Queue-based load leveling uses a queue as a buffer between an uneven producer and a limited consumer. Short bursts enqueue. Workers drain at a safe rate. The origin or the database sees a smoother load.

This pattern is a form of backpressure plus delay (Topic 15). If the queue grows without a bound, you only moved the crash to the broker disk.

Leveling fits import jobs, webhooks that arrive in bursts, and nightly fan-out. It does not hide a permanent capacity gap. If the average arrive rate stays above the drain rate, the queue never recovers.

Write:

- Max depth
- Max age
- What you do at the bound (reject, shed, or spill)
- Worker concurrency
- Alert on depth and age (Topic 16)

Idempotent consumers are required. Bursts plus retries create duplicates (Topic 11).

Load leveling is not a CDN and not a cache. It does not make a read faster. It protects a write or a job path from a spike.

Do not put a queue in front of every HTTP GET. A GET that must answer now needs capacity or a cache, not a parking lot.

Pair with `kafka.topics.md` or a broker that the team can operate. The product is not the pattern.

### Questions

#### Theoretical questions

1. What is queue-based load leveling?
2. What happens if arrive rate stays above drain rate?
3. What bounds must you write?
4. Why must consumers be idempotent?
5. Why is a queue a poor front for a user-facing GET?

#### Easy practical tasks

1. Write five sentences that define load leveling.
2. Make a table: "Pattern" and "User waits?". Add cache, CDN, offload, and load leveling.
3. List four metrics for a leveling queue.
4. Draw a burst of webhooks, a queue, and two workers.

#### Medium practical tasks

1. Design leveling for a partner file that expands to 100 000 row jobs. Write depth, age, and reject policy.
2. Explain in eight sentences how this pattern relates to backpressure without being the same sentence as Topic 15.
3. Write an alert: queue age above the lag SLO for 15 minutes.

#### Advanced practical tasks

1. Write a one-page leveling playbook: scale workers, shed, and communicate to the partner.
2. Compare a broker queue with an event log for leveling. Write replay and operations cost (Topic 10).

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do you choose among vertical scale, more instances, cache, and shards for one measured bottleneck?
2. How do load balancing and sharding split different kinds of work?
3. How do cache, CDN, and the source of truth stay distinct ideas?
4. How do async offload and queue-based load leveling work together on one job path?
5. Why does this path place scalability patterns after observability and SLOs?

#### Easy practical tasks

1. Write a one-page cheat sheet of all patterns in this topic.
2. For a photo album with 100 users, pick two patterns that you refuse and one that you might add later. Write one sentence each.
3. Draw a system that uses a balancer, a cache, and a queue. Label the source of truth.
4. Write three ADR titles: no shards yet, cache catalog, offload exports.

#### Medium practical tasks

1. Write a two-page scale brief for a campus shop exam-week peak. Include balancer, cache, offload, and a rejected shard.
2. Take a teammate design that only says "we will be cloud scale". Rewrite it as a growth axis, a measure, and two patterns.
3. Map Topics 2, 9, and 15 onto one "catalog is slow" incident. Write which pattern you try first and why.

#### Advanced practical tasks

1. Design a year plan: vertical DB, read replica, cache, then a decision gate for shards. Include SLOs and cost.
2. Read a public CDN or cache vendor overview. Write ten sentences that extract patterns, not product names. Do not copy long passages.
