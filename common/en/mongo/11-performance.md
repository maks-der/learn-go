# 11. Performance

## Description

This topic shows how you measure and improve MongoDB speed. You read `explain`. You compare `COLLSCAN` and `IXSCAN`. You keep the working set in RAM. You cut payload with projection. You use covered queries. You use the profiler and Atlas Performance Advisor.

Complete indexes, aggregation, and querying first.

Use one term for each concept. A **collection scan** (`COLLSCAN`) reads every document. An **index scan** (`IXSCAN`) reads an index. The **working set** is the data and indexes that the application uses often. A **covered query** is a query that the server answers from the index alone.

Do not tune from a guess. Measure with `explain` and with real load.

---

## `explain("executionStats")`

`explain` shows how the server plans to run a find or an aggregation. The `executionStats` mode runs the plan and reports counts and time.

```javascript
db.orders.find({ userId: "u1" }).explain("executionStats")
```

Useful fields:

- `queryPlanner.winningPlan` — the chosen plan
- `executionStats.nReturned` — documents returned
- `executionStats.totalDocsExamined` — documents read
- `executionStats.totalKeysExamined` — index keys read
- `executionStats.executionTimeMillis` — time for that run

Modes:

- `queryPlanner` — plan only, no execution counts
- `executionStats` — plan plus counts for that execution
- `allPlansExecution` — also shows rejected plans

Use `executionStats` when you learn. The run does real work. Do not run a heavy `explain("executionStats")` on a huge collection in production without care.

A good sign: `nReturned` is close to `totalDocsExamined` for a selective filter. A bad sign: you examine 1 000 000 documents to return 10.

Aggregation uses `db.orders.explain("executionStats").aggregate([ ... ])` or the helper that your version documents.

`explain` is not a load test. One call on a warm cache is not your peak. Use it to see the plan shape.

Read the explain page for your major version. Field names can change.

### Questions

#### Theoretical questions

1. What does `executionStats` add that `queryPlanner` does not add?
2. What does `totalDocsExamined` mean?
3. Why is a large gap between examined and returned a warning?
4. Why is one `explain` not a load test?
5. Why must the manual version match your server?

#### Easy practical tasks

1. Run `explain("executionStats")` on a find with no index on the filter. Write `nReturned` and `totalDocsExamined`.
2. Open the explain results page. Write the URL.
3. Write the three explain modes in a table with one sentence each.
4. Run `explain("queryPlanner")` on the same find. Write one field that execution stats had and this mode does not.

#### Medium practical tasks

1. Add an index. Run `executionStats` again. Write the two `totalDocsExamined` values.
2. Explain an aggregation that `$match`es then `$group`s. Write the first stage name in the plan.
3. Compare `allPlansExecution` with `executionStats` on a query that has two possible indexes. Write one rejected plan field.

#### Advanced practical tasks

1. Read `executionTimeMillis` vs `executionTimeMillisEstimate`. Write when each appears.
2. Capture `explain` JSON for a slow query from production (or a copy). Annotate five fields that decide your next change.

---

## COLLSCAN vs IXSCAN

`COLLSCAN` means the server reads the collection. It does not use an index for that filter (or it chooses not to).

`IXSCAN` means the server reads an index. A **`FETCH`** stage then loads full documents when the index does not hold all needed fields.

Typical winning plan for a good find:

```text
IXSCAN → FETCH → PROJECTION (optional)
```

Typical poor plan:

```text
COLLSCAN
```

`COLLSCAN` is acceptable on a tiny collection. `COLLSCAN` on millions of documents is a problem for an API path.

An index that exists is not enough. The filter must match the index (including the compound prefix rule from the indexes topic). A function on the field, a poor `$or`, or a leading `$not` can block the index.

Sort can use an index or can add a blocking `SORT`. A large `SORT` that does not fit in memory fails unless you allow disk use (more common in aggregation).

`explain` shows `stage: "COLLSCAN"` or `stage: "IXSCAN"` in the winning plan. Start there.

Do not add ten indexes to remove one `COLLSCAN` without a write-cost check. Fix the query shape first.

A **write-heavy** workload inserts or updates many documents. Extra indexes slow writes. A **read-heavy** workload runs many finds. Indexes, projection, and a working set in RAM matter more. Measure, then change one thing.

### Questions

#### Theoretical questions

1. What does `COLLSCAN` mean?
2. What does `IXSCAN` mean?
3. What does `FETCH` do after `IXSCAN`?
4. When is `COLLSCAN` acceptable?
5. Why can a filter ignore an index that exists?

#### Easy practical tasks

1. Force a `COLLSCAN` (find on an unindexed field). Write the stage name from `explain`.
2. Create an index. Show `IXSCAN` on the same filter.
3. Open the query-plans page. Write the URL.
4. Make a table: stage, what it reads. Add `COLLSCAN`, `IXSCAN`, `FETCH`.

#### Medium practical tasks

1. Write a filter that cannot use your compound index (wrong prefix). Prove `COLLSCAN` or a different plan with `explain`.
2. Add a sort that the index cannot cover. Write whether the plan shows `SORT`.
3. Time the same find with `COLLSCAN` and with `IXSCAN` on at least 10 000 documents. Write the two times.

#### Advanced practical tasks

1. Read about `COUNT_SCAN` or `IDHACK` if your `explain` shows them. Write one sentence each.
2. Find a `$or` or regex query that stays on `COLLSCAN`. Rewrite it or add an index. Prove the new stage.

---

## Working set, projection, and covered queries

The **working set** is the documents and index pages that the application reads and writes often. MongoDB (WiredTiger) keeps as much as it can in the **WiredTiger cache**. The operating system also caches file pages.

If the working set fits in RAM, most reads avoid slow disk. If the working set is larger than RAM, the server pages data in and out. Latency grows. **Page faults** and cache eviction increase.

Fit the working set:

- Keep documents small
- Index only what you query
- Archive cold data
- Add RAM or add shards when the hot data is truly large

The cache size has a default on self-managed `mongod` (a large fraction of RAM, with a reserved amount). Atlas sets cache from the cluster tier. Do not set cache to 100% of RAM. The operating system and connections need memory.

Working set is not the same as total data size. A 2 TB collection can have a 20 GB hot set.

Do not buy a shard only because total disk is large. Measure the hot set first.

**Projection** selects fields in a find or in `$project`. A smaller document on the wire is faster. The application also parses less data.

```javascript
db.orders.find({ userId: "u1" }, { total: 1, status: 1 })
```

Include the fields that the caller needs. Exclude large unused fields. Projection does not always avoid a `FETCH`. The server may still load the full document from storage, then drop fields. Projection still helps the network and the client.

Do not return the full document in a list API "for later". Add a detail endpoint.

`find({}, {})` with no limit is still unbounded. Projection does not replace `limit`.

Put `$match` first in aggregation. Project early when you drop large fields before a `$group` or a `$lookup`.

A **covered query** uses an index that contains every field in the filter, the sort, and the projection. The server does not `FETCH` the full document.

`explain` shows `IXSCAN` and then a projection or a count, without `FETCH`. `totalDocsExamined` is `0` in many covered plans.

Requirements (typical):

- All filter fields are in the index
- All projected fields are in the index
- The projection excludes fields that are not in the index
- You do not need fields outside the index

Example:

```javascript
db.orders.createIndex({ userId: 1, status: 1 })
db.orders.find({ userId: "u1" }, { _id: 0, userId: 1, status: 1 })
```

`_id` is in the document and is included by default. If `_id` is not in the index, exclude `_id` in the projection or add `_id` to the index.

Covered queries help hot read paths with a small field set. They do not help when you need the whole document.

Multikey indexes (arrays) often cannot cover in the way you expect. Read the current limits.

Do not add a huge compound index only to cover a rare query.

### Questions

#### Theoretical questions

1. What is the working set?
2. How is working set different from total collection size?
3. Does projection always avoid reading the full document from disk?
4. What makes a query covered?
5. Why does default `_id` break coverage if `_id` is not in the index?

#### Easy practical tasks

1. Write four sentences: working set, cache, cold data, total size.
2. Find three documents with and without projection. Compare printed size by eye. Write which fields you dropped.
3. Create the covered-query example index. Run the example find with `explain`. Write if you see `FETCH`.
4. Repeat the covered find with `_id` included. Write the plan change.

#### Medium practical tasks

1. In Atlas or `serverStatus`, find cache bytes. Write the used and max values.
2. Store a 100 KB string in a field. Time find with and without that field in the projection on 1000 documents. Write the two times.
3. Build a covered count or exists check for `{ email: 1 }` with projection `{ _id: 0, email: 1 }`. Prove it with `explain`.

#### Advanced practical tasks

1. Read the default cache-size formula for self-managed `mongod`. Compute it for a 16 GB host.
2. Read coverage rules for embedded fields and arrays in your version. Write three limits.

---

## Profiler and Atlas Performance Advisor

The **database profiler** records slow operations into `system.profile`. You set a level and a slow-ms threshold.

```javascript
db.setProfilingLevel(1, { slowms: 100 })
```

Level `0` is off. Level `1` records slow operations. Level `2` records all operations (heavy). Do not leave level `2` on a busy production database.

Read entries:

```javascript
db.system.profile.find().sort({ ts: -1 }).limit(5)
```

Each entry shows the command, millis, and plan summary. Use it to find `COLLSCAN` and long locks.

**Atlas Performance Advisor** reads slow logs and suggests indexes. It is a helper. You still verify with `explain`. A suggested index can be a duplicate or can hurt writes.

Atlas also shows:

- Slow query logs
- Real-time performance
- Query Insights on some tiers

`serverStatus` and `currentOp` help on self-managed servers. `db.currentOp()` shows operations that run now.

Log slow operations even when the profiler is off (`operationProfiling.slowOpThresholdMs` or Atlas default). The log is enough for many teams.

Do not add every Advisor index on the same day. Add one. Measure. Keep or drop.

### Questions

#### Theoretical questions

1. What does profiler level `1` record?
2. Why is level `2` dangerous on a busy database?
3. What does Performance Advisor suggest?
4. Why must you still run `explain` after an Advisor index?
5. What does `currentOp` show that the profiler does not show?

#### Easy practical tasks

1. Set profiling level `1` with `slowms` 50 on a learning database. Run a slow find. Read `system.profile`.
2. Turn the profiler back to `0`. Confirm with `db.getProfilingStatus()`.
3. Open the profiler page and the Atlas Performance Advisor page. Write both URLs.
4. Write three fields that you will read first in a profile document.

#### Medium practical tasks

1. If you have Atlas, open Performance Advisor. Write one suggested index or write that there is none. Do not apply it yet.
2. Compare a profiler entry with `explain` for the same filter. Write matching stage names.
3. Run `db.currentOp()` during a long aggregation (or a sleep in a lab). Write one field that identifies the operation.

#### Advanced practical tasks

1. Read profiler overhead and `system.profile` capped size. Write how you would enable it for one hour in production.
2. Build a weekly review: top five slow ops, one index added, one index unused (`$indexStats`). Write the template.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do `explain`, the profiler, and Performance Advisor answer different questions?
2. When do `COLLSCAN`, a large working set, and a fat payload each look "slow" with different metrics?
3. How do covered queries and projection both cut work, and how do they differ?
4. Why does a write-heavy collection reject the same index plan that a read-heavy collection needs?
5. A teammate adds five Advisor indexes and sets profiler level `2` on production. Which facts do you use in the review?

#### Easy practical tasks

1. Write a cheat sheet: `executionStats` fields, `COLLSCAN` vs `IXSCAN`, working set, projection, covered, profiler levels, Advisor.
2. On one collection, capture one `explain("executionStats")` and write the winning stage.
3. Draw a flow: slow API → profiler → `explain` → index or projection → measure again.
4. Mark three queries as write-path or read-path. Write one tune for each.

#### Medium practical tasks

1. Take a find that does `COLLSCAN`. Add one index and a projection. Prove `IXSCAN` and a smaller payload. Write the numbers.
2. Estimate whether your learning data set fits in cache. Write RAM, cache, and data size.
3. Write a performance checklist of ten items that you run before you add a shard.

#### Advanced practical tasks

1. Load 100 000 documents. Create a slow report and a fast API find. Document `explain`, one index, and one projection change with before/after millis.
2. Write a team policy: when to use Advisor, when to profile, max indexes per collection, and who approves a covering index.
