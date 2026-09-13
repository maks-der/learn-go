# 10. Sharding

## Description

This topic shows how MongoDB splits one data set across many machines. A **sharded cluster** has **shards**, a **mongos** router, and **config servers**. You learn when to shard, how you choose a shard key, hashed vs ranged keys, chunks, the balancer, hot shards, and resharding.

Complete replication, modeling, and indexes first.

Use one term for each concept. A **shard** is a replica set that holds a subset of the data. The **shard key** is the field list that decides which shard owns a document. A **chunk** is a range of shard-key values that the balancer can move. **`mongos`** is the process that routes queries.

Sharding increases capacity. Sharding also increases operational cost. Do not shard a collection that still fits on one replica set.

---

## When to shard

Shard when one replica set is not enough. Common signals:

- The working set does not fit in RAM on one replica set
- Write throughput saturates the primary
- Disk on one replica set is full or will be full soon
- A single replica set cannot meet the latency target under peak load

Do not shard because "large systems shard". Many applications stay on one replica set for a long time. Atlas and modern hardware hold large data sets on one replica set.

Sharding costs:

- You must choose a shard key
- Some queries become scatter-gather (they hit many shards)
- Transactions and joins (`$lookup`) have extra limits and cost
- You operate many replica sets plus config servers plus `mongos`

Shard a collection, not "the whole database" as one unit. Other collections can stay unsharded. Unsharded collections live on a primary shard.

Complete modeling first. A bad document shape plus a shard is still a bad shape.

For learning, you can start a sharded cluster in Docker or use an Atlas sharded cluster. On paper, you can design a shard key before you pay for extra shards.

### Questions

#### Theoretical questions

1. What problem does sharding solve that replication does not solve?
2. Why is a full disk on one replica set a reason to consider sharding?
3. Why is "we might grow" a weak reason to shard today?
4. Where does an unsharded collection live?
5. What extra processes does a sharded cluster add?

#### Easy practical tasks

1. Write five sentences: replica set vs sharded cluster.
2. Open the sharding introduction in the manual. Write the URL.
3. List three metrics that would push you to shard. List three that would not.
4. Draw: client, `mongos`, two shards (replica sets), config servers.

#### Medium practical tasks

1. Estimate RAM, disk, and write rate for a collection that you know. Write if one replica set is enough.
2. Read Atlas shard options for your tier. Write the minimum shard count and what you pay for.
3. Interview a teammate or read a case study. Record one reason they sharded or they did not.

#### Advanced practical tasks

1. Read limits of `$lookup` and transactions on a sharded cluster for your version. Write three limits.
2. Write a one-page "shard or not" decision for a 400 GB collection with a 40 GB working set.

---

## Shard key choice

The **shard key** is one or more fields that exist in every document. MongoDB hashes or ranges those values to place the document.

The shard key is the hard part. A poor key creates a **hot shard** (one shard gets most writes or reads) or **jumbo chunks** that do not move.

A useful shard key is:

- **High cardinality** — many distinct values
- **Even frequency** — values do not all pile on one value
- **Non-monotonic for writes** if you need even write spread (or use a hashed key)
- **Present in the query filter** for targeted reads

A poor shard key is:

- A boolean or a status with three values
- A monotonically increasing `_id` or timestamp for ranged sharding (all inserts go to the last chunk)
- A field that most queries do not include

The shard key values in a document have rules. You cannot unset the shard key. You can change a document's shard-key value in current versions, but the change is expensive. Plan the key as if it will stay.

Compound shard keys are common. Example: `{ tenantId: 1, orderId: 1 }`. The first field targets a tenant. The second field adds cardinality.

The application must send the shard key in queries when it can. A find without the shard key can scatter to all shards.

Do not pick a key only because it is unique. `_id` is unique and can still be a poor ranged key if it increases with time.

Read the official shard-key pages. Test with production-like data before you shard a live collection.

### Questions

#### Theoretical questions

1. What is a shard key?
2. Why is low cardinality a problem?
3. Why can a time-increasing `_id` be a poor ranged key?
4. Why must many queries include the shard key?
5. Why is a compound shard key often better than one low-cardinality field?

#### Easy practical tasks

1. For `orders`, propose two shard keys. Write one reason for each.
2. Open the shard-key page. Write the URL and two properties of a good key.
3. Make a table: field, cardinality, in queries?, monotonic? Score three fields.
4. Write why a boolean `isActive` is a poor shard key.

#### Medium practical tasks

1. On paper, shard `orders` by `{ customerId: 1, _id: 1 }`. Write which finds target one shard.
2. Take a collection that you modeled in topic 3. Propose a key. List two queries that scatter.
3. Read shard-key immutability and change rules for your version. Write what you can change.

#### Advanced practical tasks

1. Design a shard key for a multi-tenant SaaS. Include a tenant with 100 times the data of others (outlier). Write the plan.
2. Review an official "choosing a shard key" page. Rewrite their checklist in eight short sentences in your own words.

---

## Hashed vs ranged keys

A **ranged** shard key places documents by the order of the key values. Adjacent values can live in the same chunk. Range queries on the shard key can hit one shard or a few shards.

A **hashed** shard key computes a hash of the key value. MongoDB places documents by the hash. Inserts spread even if the source field increases with time.

Use a ranged key when:

- Many queries use a range on that key
- Values already spread writes
- You want nearby keys on the same shard

Use a hashed key when:

- The natural key is monotonic (ObjectId, timestamp)
- You need even write spread
- You accept that a range query on the source field will scatter

Hashed keys do not help a range query on the original field. Example: hashed `_id` does not make `find({ _id: { $gt: ... } })` a single-shard range.

You create a hashed index as the shard key:

```javascript
sh.shardCollection("shop.orders", { customerId: "hashed" })
```

Compound keys can mix a ranged prefix and other fields. Read the current rules. Hashed compound keys have version-specific limits.

Do not hash a low-cardinality field and expect magic. A boolean still has two hashes.

Atlas and `mongosh` helpers (`sh.status()`, `sh.shardCollection`) are the usual tools.

### Questions

#### Theoretical questions

1. How does a ranged shard key place documents?
2. How does a hashed shard key place documents?
3. When does a monotonic field need a hashed key?
4. Why can a range query fail to target one shard with a hashed key?
5. Does hashing a boolean fix low cardinality?

#### Easy practical tasks

1. Write four sentences: hashed vs ranged.
2. Open the hashed-sharding page. Write the URL.
3. For `createdAt` inserts, write which key type spreads writes. Write which type helps a date range.
4. Make a table: key type, write spread, range query on that field.

#### Medium practical tasks

1. If you have a sharded lab, shard one collection hashed and one ranged (or write the two `shardCollection` commands).
2. Draw 12 inserts with increasing `_id` on a ranged key and on a hashed key. Mark the target shard.
3. Choose hashed or ranged for `tenantId` plus `orderId`. Write one paragraph.

#### Advanced practical tasks

1. Read compound hashed shard-key rules for your major version. Write two allowed patterns and one forbidden pattern.
2. Estimate scatter-gather cost for a date-range report on a hashed `_id` collection. Propose a ranged field that reports can use.

---

## Chunks, balancers, hot shards

A **chunk** is a contiguous range of shard-key values (or hashed values) that belongs to one shard. The cluster splits chunks when they grow. The **balancer** moves chunks so that shards have a similar amount of data (and, in current versions, similar load by the policy that you configure).

Config servers store the chunk map. `mongos` caches that map and routes each query.

Typical operations:

```javascript
sh.status()
sh.isBalancerRunning()
sh.startBalancer()
sh.stopBalancer()
```

Do not stop the balancer for a long time in production without a reason. Chunks grow uneven. Then a catch-up move is large.

Do not run heavy schema changes, some backups, or some upgrades without the procedure in the manual. The procedure may ask you to stop the balancer for a short window.

Chunk size has a default (often 128 MB in many versions). Large documents make fewer documents per chunk. Small documents make more documents per chunk.

A query that includes the full shard key can target one chunk or one shard. A query with no shard key can hit all shards.

Moving a chunk uses resources on the source shard, the destination shard, and the network. Balance during a peak write period can add latency.

A **hot shard** is a shard that receives much more traffic or data than the others. Causes:

- A monotonic ranged key (all inserts go to the latest chunk)
- A popular shard-key value (one tenant, one country)
- A query pattern that always targets one prefix

A **jumbo chunk** is a chunk that is too large to move under the balancer rules. Causes:

- Too many documents with the same shard-key prefix (the chunk cannot split)
- Very large documents
- A key with almost no distinct values

Jumbo chunks stay on one shard. The balancer cannot fix that shard by moving that chunk. You must change the key, refine the key, split if the version allows it, or reshard.

Detect hot shards with metrics: opcounters, CPU, disk, and chunk counts that do not match traffic. Atlas has shard metrics. `sh.status()` shows jumbo flags on some versions.

Fix the model first when one tenant is huge. The outlier pattern from schema design can apply. Sharding does not remove a hot key.

Do not ignore jumbo warnings. They grow. Prevention is easier than repair. Pick cardinality and even frequency before you shard.

### Questions

#### Theoretical questions

1. What is a chunk?
2. What does the balancer move?
3. What is a hot shard?
4. What is a jumbo chunk?
5. Why can the balancer fail to cool a hot shard?

#### Easy practical tasks

1. Open the chunks, balancer, and jumbo-chunk pages. Write the URLs.
2. Draw two shards and four chunks. Show one chunk move. Then draw a ranged `_id` key with all new writes on shard B. Label "hot".
3. Write three causes of a hot shard. Write three causes of a jumbo chunk.
4. List three `sh.` helpers from this section and one purpose each.

#### Medium practical tasks

1. On a sharded lab or Atlas, run `sh.status()`. Write chunk counts per shard.
2. Read the default chunk size for your version. Write the number and what a 16 MB document means for documents per chunk.
3. In a lab, insert many documents with the same shard-key value. Inspect `sh.status()` for uneven chunks or jumbo notes.

#### Advanced practical tasks

1. Read balancer windows (active window). Propose a window for a region that peaks at 18:00. Write the config idea.
2. Plan a production response: jumbo chunk on the primary write shard. Include detect, mitigate, and long-term key change.

---

## Resharding (modern versions)

Older MongoDB versions made the shard key almost fixed after `shardCollection`. Current versions add two important tools.

**Refine a shard key** (`refineCollectionShardKey`) adds one or more suffix fields to an existing key. Example: `{ customerId: 1 }` becomes `{ customerId: 1, orderId: 1 }`. You must have a supporting index. Refine helps cardinality when the old prefix stays. It does not change the existing prefix field.

**Reshard a collection** (`reshardCollection`) changes the shard key to a new key. The cluster copies data to a new distribution. The operation needs extra disk and time. Read the current limits (size, disk headroom, blocking behavior).

These tools reduce the cost of a first-key mistake. They do not make the first key cheap. A reshard of a large collection is still a project.

Feature availability depends on the major version and on Atlas. Use the manual that matches your version.

Steps in spirit:

1. Create the required index for the new key.
2. Run refine or reshard as the manual shows.
3. Watch balancer and disk.
4. Change application queries to include the new key.

Do not refine instead of reshard when the prefix itself is the problem (for example the prefix is a boolean).

Do not run reshard without a disk and rollback plan.

### Questions

#### Theoretical questions

1. What does refine add to a shard key?
2. What does reshard change?
3. Why do you need an index before refine?
4. When is refine not enough?
5. Why is reshard still a project on a large collection?

#### Easy practical tasks

1. Open the refine and reshard pages for your version. Write both URLs.
2. Write one example refine: old key, new key, new index.
3. Make a table: refine vs reshard. Add three rows (what changes, disk, when to use).
4. Write the major version that first documents reshard in the manual that you open.

#### Medium practical tasks

1. On a lab cluster that supports it, shard a small collection, refine the key, and run `sh.status()`. Write the new key.
2. Read disk-space requirements for reshard. Write the headroom number or rule.
3. Write an application change list for a new shard key (filters, unique indexes, transactions).

#### Advanced practical tasks

1. Plan a reshard from `{ _id: 1 }` ranged to `{ tenantId: 1, _id: 1 }`. Include index, disk, query changes, and a test plan.
2. Compare "shard late with a good key" vs "shard early and reshard later" for a startup. Write a recommendation with two risks.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do shard key, chunks, and `mongos` decide where one insert lives?
2. When do hashed keys, ranged keys, and compound keys solve different hot-shard problems?
3. Why can jumbo chunks and a monotonic key both make the balancer look "broken"?
4. How do refine and reshard change the old rule "you cannot change a shard key"?
5. A teammate shards on `{ status: 1 }` because every query filters status. Which facts do you use in the review?

#### Easy practical tasks

1. Write a cheat sheet: when to shard, good key tests, hashed vs ranged, chunk, balancer, jumbo, refine, reshard.
2. Draw a full cluster: three shards, config replica set, two `mongos`, one targeted find, one scatter find.
3. For one real collection name, write a candidate key and one query that would scatter.
4. Open `sh.status()` docs. Write two fields that you will read first in an incident.

#### Medium practical tasks

1. Design shard keys for `users`, `orders`, and `auditLog`. Justify hashed or ranged for each.
2. Write a pre-shard checklist of ten items (metrics, key, indexes, balancer window, rollback).
3. On paper, show how a hot tenant plus jumbo chunks would look in `sh.status()` and in CPU graphs.

#### Advanced practical tasks

1. Run or script a small sharded cluster. Shard a collection, generate skew, and record metrics. Then propose a refine or reshard.
2. Write a production standard: when the team may shard, who approves the key, how you load-test, and how you watch jumbo chunks.
