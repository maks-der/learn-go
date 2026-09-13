# 20. Internals (Advanced)

## Description

This topic shows internal behavior at a high level. You learn the **WiredTiger cache**. You learn how documents live on disk now, and what **padding** meant in old engines. You learn **index structures** at a high level. You learn the **oplog window**. You learn that the **balancer** has an algorithm. Complete administration, indexes, replication, and sharding first.

Use one term for each concept. The **cache** is RAM that WiredTiger uses for data and indexes. **Padding** was extra unused space in old MMAPv1 records. The **oplog window** is how far back secondaries can catch up from the oplog. The **balancer algorithm** is the set of rules that pick which chunks to move.

This topic is awareness. It is not a license to change undocumented flags in production. Read the current manual. Internals change between major versions.

---

## WiredTiger cache

WiredTiger keeps recently used pages in a **cache**. The cache holds uncompressed-in-memory images of data and index pages (details depend on version). A read that hits the cache does not wait for disk.

Default size on self-managed `mongod` is a large part of RAM minus a reserved amount (a common formula is about 50% of RAM minus 1 GB, with a minimum). Atlas sets the cache from the cluster tier. You can set `storage.wiredTiger.engineConfig.cacheSizeGB`. Do not set it so large that the OS and connections have no RAM.

**Dirty** pages are changed pages that are not yet written to the data files. WiredTiger evicts and checkpoints. If application writes create dirty pages faster than eviction, latency grows. Tickets (internal queues) can stall writes.

Useful `serverStatus` ideas:

- Cache used vs max
- Dirty bytes
- Eviction counters
- Application threads that wait for cache

A **working set** that fits in the cache is the performance goal from topic 14. The cache is the mechanism.

Do not treat the cache as a user-controlled query cache of full result sets. It is a storage engine cache of pages.

Do not raise the cache to 95% of RAM because a blog said "more cache is better". The OS page cache and WiredTiger interact. Leave room.

Compression: on-disk pages are compressed. In-cache pages are larger. Disk size is not cache size.

### Questions

#### Theoretical questions

1. What does the WiredTiger cache hold?
2. Why must the cache leave RAM for the operating system?
3. What is a dirty page in this context?
4. How does the cache relate to the working set?
5. Why is on-disk size not equal to cache size?

#### Easy practical tasks

1. Open the WiredTiger cache configuration page. Write the URL and the default-size idea.
2. Read `serverStatus().wiredTiger.cache` if you can. Write two numbers.
3. Write four sentences: cache, dirty, eviction, RAM headroom.
4. Make a table: cache too small, cache too large. Add one symptom each.

#### Medium practical tasks

1. Compute a default cache size for 8 GB RAM and for 32 GB RAM with the formula in the manual.
2. Run a write-heavy loop in a lab. Watch dirty bytes. Write the change.
3. Compare Atlas tier RAM with the cache that Atlas documents (or a reasonable estimate). Write the pair.

#### Advanced practical tasks

1. Read eviction and checkpoint at a high level. Draw a cycle: write, dirty, checkpoint, evict.
2. Read WiredTiger tickets. Write what a thread waits for when tickets are 0.

---

## Document storage and padding (historical vs current)

**Historical (MMAPv1).** The old engine stored a document in a record. An update that grew the document might not fit. MMAPv1 used **padding** (extra unused bytes) so that small growth did not move the record. Moves caused extra I/O. Collection-level locks were common in that era. MMAPv1 is **not** in current MongoDB.

**Current (WiredTiger).** WiredTiger stores documents in a different on-disk format (B-trees / tables, with pages). An update does not use MMAPv1 padding. Growth of a document still has a cost: the engine writes new page data and frees old space. Very large documents and frequent growth still hurt. The 16 MB limit remains.

**You do not tune padding** on WiredTiger. Old advice (`usePowerOf2Sizes`, padding factor) does not apply.

**Free space.** After deletes, WiredTiger can reuse space. Files may not shrink until compact or a rewrite (topic 16). That is not padding. It is file reuse.

**BSON size** in memory is not the same as bytes on disk after compression.

When you read an old book or a 2014 talk, check the engine name. If the talk is about record moves and padding, it is historical.

Current care:

- Keep documents bounded
- Avoid patterns that grow one document forever
- Prefer `$set` of stable field sizes when you can

Do not add fake fields to "create padding" on WiredTiger. That is wasted space.

### Questions

#### Theoretical questions

1. What problem did MMAPv1 padding try to solve?
2. Does current MongoDB use MMAPv1?
3. Why is a padding-factor blog post a poor guide now?
4. How is leftover file space after delete different from padding?
5. Why can a document that grows on every update still be costly on WiredTiger?

#### Easy practical tasks

1. Write a two-column table: MMAPv1 (historical), WiredTiger (current). Add three rows.
2. Open a storage-engine or WiredTiger FAQ. Write the URL.
3. Find one old article that mentions padding. Write the date and "historical".
4. Write four sentences: 16 MB limit, compression, no padding tune, bounded growth.

#### Medium practical tasks

1. Update a lab document in a loop with a growing string. Watch `collStats` sizes. Write the trend.
2. Read `collStats` fields `size` vs `storageSize`. Write what each means on WiredTiger.
3. List three old MMAPv1 terms and the current term or "gone".

#### Advanced practical tasks

1. Read WiredTiger page and reconciliation at a very high level (official or WT docs). Write six sentences that a developer can use.
2. Write a "ignore this advice" list of five MMAPv1 tips that you must not apply.

---

## Index structures at a high level

A MongoDB index is a WiredTiger table (or similar structure) that maps **keys** to **record locators**. The common index is a **B-tree** (B+ tree style) on the index key.

High-level facts:

- Keys are ordered. Range scans walk adjacent keys.
- A **compound** index is one tree on a composite key, not two independent trees.
- A **unique** index rejects a second identical key.
- A **multikey** index has one key per array element (with limits).
- A **hashed** index stores the hash. It does not support a useful range on the source field.
- **TTL** uses an index on a date field and a background deleter.
- **Text** and **geospatial** indexes have their own structures and query planners.
- Atlas Search indexes are Lucene structures **outside** this WiredTiger B-tree story.

The query planner picks an index. `explain` shows `IXSCAN` of that tree.

Indexes live in the WiredTiger cache too. Extra indexes reduce cache for documents.

An index key has a size limit (historical 1024 bytes; read your version). Huge keys fail or are excluded. Do not index a giant string.

Hidden indexes stay in storage but the planner does not use them. They still cost writes and cache until you drop them.

This section does not replace the indexes topic. It explains why prefix order and range scans work.

Do not assume a hash table for a normal `{ a: 1 }` index. It is an ordered tree.

### Questions

#### Theoretical questions

1. What does an index map to a record locator?
2. Why does a B-tree help a range query?
3. Why is a compound index one structure?
4. Why do extra indexes compete for cache?
5. How is an Atlas Search index different from a WiredTiger index?

#### Easy practical tasks

1. Draw a tiny B-tree of three `userId` keys that point to documents.
2. Open the index-storage or index-properties page. Write the URL.
3. Make a table: index type, structure idea. Add single, compound, hashed, multikey.
4. Write why `{ a: 1, b: 1 }` is not `{ a: 1 }` plus `{ b: 1 }`.

#### Medium practical tasks

1. Create a compound index. Run a prefix query and a suffix-only query. Compare `explain`.
2. Read the index key size limit for your version. Write the number.
3. Hide an index. Confirm the planner ignores it and that writes still update it (`$indexStats` or docs).

#### Advanced practical tasks

1. Read WiredTiger `file_type` or index table notes at a high level. Write how an index page miss becomes disk I/O.
2. Estimate cache for two 5 GB indexes plus a 10 GB hot data set. Write if a 16 GB cache is enough.

---

## Oplog window

The **oplog window** is the time between the oldest and newest entries that still exist in `local.oplog.rs`.

Secondaries, change streams, and some delayed members need that history. If a consumer is slower than the window, it cannot resume from the oplog. Then you need **initial sync** (for a member) or a **full CDC rebuild** (for a listener).

The window shrinks when:

- The oplog is small
- The write rate is high
- Large writes create large entries

The window grows when you increase oplog size or when writes slow down.

Atlas shows oplog window and lets you grow the oplog on many tiers. Self-managed servers set oplog size in configuration or with the procedure in the manual.

A **delayed member** must have a delay that is **shorter** than the oplog window. If the delay is 24 hours and the window is 4 hours, the delayed member cannot apply.

Backups that rely on oplog tailing (some incremental tools) have the same limit.

Do not set a tiny oplog to save disk on a write-heavy cluster.

Do not ignore a shrinking window during a backfill.

`rs.printReplicationInfo()` estimates the window. Treat it as an estimate.

### Questions

#### Theoretical questions

1. What is the oplog window?
2. What happens to a secondary that is older than the window?
3. Why does a high write rate shrink the window for a fixed oplog size?
4. Why must a delayed member delay be less than the window?
5. How does a change-stream listener fail when the window is too short?

#### Easy practical tasks

1. Open the oplog sizing page. Write the URL.
2. Run `rs.printReplicationInfo()` or Atlas oplog metrics. Write the window or "unknown".
3. Write four consumers of the oplog (secondary, delayed, change stream, incremental backup).
4. Make a table: write rate up, oplog size up. Add effect on the window.

#### Medium practical tasks

1. Insert a large batch. Compare the estimated window before and after. Write the numbers.
2. Read how to resize the oplog on your version. Write the steps. Do not run them on a shared cluster.
3. Compute a rough window: oplog size / bytes per second of writes. Write the assumption.

#### Advanced practical tasks

1. Design oplog size for a backfill that writes 2 hours of extra traffic. Include change streams that must survive.
2. Write an alert: window below 2 hours. Include who grows the oplog and who pauses the backfill.

---

## Cluster balancer algorithm (awareness)

The **balancer** runs on the sharded cluster (a process associated with the config servers in current versions). It moves **chunks** from a shard that has too many to a shard that has too few.

Awareness, not source-code detail:

- The balancer uses the **chunk map** in config data
- It prefers to even **data size** (and, in current versions, can consider other imbalance signals that the docs describe)
- It moves a limited number of chunks at a time
- It can run in a **window** (hours of the day)
- It will not move a **jumbo** chunk
- A move copies the chunk to the target, then commits the new map, then deletes the source copy (high-level; the exact protocol has versions)
- `mongos` routes with a cached map and refreshes

You do not pick each chunk by hand in normal operations. You set policy: balancer on/off, window, maybe auto-split settings that your version documents.

If all writes hit one chunk (hot shard), the balancer **cannot** split a key that has no distinct values. The algorithm is not broken. The shard key is.

Do not stop the balancer forever to "reduce load" unless you are in a documented procedure. Imbalance then grows.

Do not expect the balancer to fix a monotonic ranged key. New data still lands on the latest chunk.

Read the balancer page for your major version. The scoring details change. Awareness means you know **what** it optimizes and **what** it cannot fix.

### Questions

#### Theoretical questions

1. What does the balancer move?
2. Where does the cluster store the chunk map?
3. Why does a jumbo chunk stop a move?
4. Why can a perfect balancer still leave a hot shard?
5. What is a balancer window?

#### Easy practical tasks

1. Open the balancer page. Write the URL and two settings that you can change.
2. Write four sentences: imbalance, move, jumbo, shard key.
3. Draw two shards, uneven chunks, one move, updated map.
4. Make a table: balancer can fix, balancer cannot fix. Add two rows each.

#### Medium practical tasks

1. On a sharded lab, run `sh.status()`. Write chunk counts. Note if the balancer is on.
2. Read the current move protocol at a high level (migrate + commit). Write five steps in your own words.
3. Propose a balancer window for a region that peaks 09:00–18:00.

#### Advanced practical tasks

1. Read changelog notes for balancer improvements in the last two major versions. Write two changes.
2. Write an incident guide: "balancer running but shard 0 is hot". Include key, jumbo, and window checks.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do the WiredTiger cache, index B-trees, and a working set that does not fit create one latency story?
2. Why must you ignore MMAPv1 padding advice when you size documents and disks on WiredTiger?
3. How do oplog window and balancer algorithm limit two different kinds of catch-up (a secondary vs a shard)?
4. When does growing the cache, growing the oplog, or starting the balancer each fail to fix the real cause?
5. A teammate sets cache to all RAM, adds fake padding fields, and stops the balancer because CPU moved. Which facts do you use in the review?

#### Easy practical tasks

1. Write a cheat sheet: cache formula, dirty pages, MMAPv1 vs WT, B-tree index, oplog window, jumbo, balancer window.
2. Draw one write: journal/checkpoint idea, cache dirty, oplog entry, secondary apply.
3. List five blog keywords that mean "historical" (MMAPv1, padding factor, record move, collection lock, power of two).
4. Record your cluster: cache size if known, oplog window if known, sharded or not.

#### Medium practical tasks

1. From `serverStatus` and `rs.printReplicationInfo()` (or Atlas), fill a one-page internals snapshot for your lab.
2. Explain to a teammate, in ten spoken sentences, why a hot shard is a key problem and not a balancer bug.
3. Estimate whether your indexes plus hot documents fit in cache. Write the arithmetic.

#### Advanced practical tasks

1. Write an internals primer for your team (two pages): cache, storage history, indexes, oplog window, balancer awareness, with links to the current manual.
2. Trace one production incident (real or invented) to cache, oplog, or balancer. Write the evidence you would collect first.
