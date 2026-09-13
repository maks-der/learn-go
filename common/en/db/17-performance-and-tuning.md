# 17. Performance and Tuning

## Description

This topic shows how you find a slow database and how you change the right part. You learn measurement, slow-query logs, metrics, cardinality, statistics, hot rows, lock contention, pagination, caching, and connection storms.

Use one term for each concept. A metric is a number that you collect over time. Cardinality is the number of distinct values. Tuning is a change that you can measure. Complete this topic after you can read a query plan and use indexes. Plans and indexes are the usual first tools.

Do not change a setting before you have a baseline. Measure, change one thing, measure again.

---

## Measure first: slow query log, metrics

A slow query log records statements that exceed a time threshold. You set the threshold. The log gives the SQL text, the duration, and sometimes a plan. Start here when users say "it is slow" and you do not know which statement.

Metrics are time series: connections, CPU, disk I/O, cache hit ratio, lock waits, replication lag, transaction rate. Metrics show the shape of the problem (CPU bound, I/O bound, lock bound, connection bound). They do not always name the SQL.

```text
User report --> metrics (which resource) --> slow log or traces (which SQL) --> plan --> change
```

A baseline is a measurement before the change. Record: query text, duration, plan, row counts, and the time of day. Without a baseline you cannot know if the change helped.

Typical first questions:

1. Is the host saturated (CPU, disk, memory)?
2. Is one statement slow, or are many statements waiting?
3. Is the wait on I/O, on a lock, or on a sequential scan?

Do not enable a very low slow-log threshold on a busy production system without a plan. The log can fill the disk. Use a sample, a higher threshold, or a product-specific statement store.

Do not tune by folklore ("increase shared buffers and it will be fast"). Correlate a metric with a statement.

Application metrics matter. A 200 ms API time with 5 ms in the DBMS is not a database problem. Measure both sides.

### Questions

#### Theoretical questions

1. What does a slow query log record?
2. What do metrics show that a single log line may not show?
3. What is a baseline?
4. Why can a very low slow-log threshold harm a busy system?
5. Why must you measure the application and the DBMS?

#### Easy practical tasks

1. Find the slow-query log setting and the threshold name in your DBMS.
2. Write the four-step path from this section in your notes.
3. List five metrics that you would graph for a small server.
4. Write three baseline fields that you store before an index change.

#### Medium practical tasks

1. Turn on a slow log in a learning instance. Run a sequential scan on a large scratch table. Find the line in the log.
2. Record CPU and disk busy (OS or DBMS) during that scan. Write which resource rose.
3. Compare one API timer with the DBMS duration for the same request (log both). Write the two numbers.

#### Advanced practical tasks

1. Write a one-page measure playbook: logs, metrics, traces, and when you stop at the application.
2. Build a weekly review: top five statements by total time. Keep the list. Do not change production in this task unless you own it.

---

## Cardinality and statistics

Cardinality is the number of distinct values in a column (or in a group of columns). High cardinality means many distinct values. Example: a unique email column has cardinality equal to the number of rows. A boolean column has cardinality 2.

The optimizer uses statistics to estimate how many rows a filter returns. Statistics are samples or summaries: distinct counts, histograms, most-common values, null fraction. The optimizer picks a plan from those estimates.

```sql
-- High selectivity (few rows) if customer_id is unique in this filter
WHERE customer_id = 42

-- Low selectivity if status has two values and most rows are 'active'
WHERE status = 'active'
```

If statistics are stale, the estimate is wrong. A wrong estimate can pick a nested loop that should be a hash join, or a scan that should be a seek. After a large load, refresh statistics (product command: `ANALYZE` or equivalent).

Cardinality of a join result is not the product of the table sizes if keys match. The optimizer uses join statistics or assumptions. Bad join estimates are a common cause of a bad plan.

Do not create an index on a boolean that is true for 99 percent of rows and expect a seek to help those queries. Low selectivity favors a scan.

Do not confuse cardinality with table size. A table of 10 million rows can have a column of cardinality 10.

Read the estimate and the actual row count in a plan that shows both. A large gap means stale or weak statistics, or a correlation that the stats model misses.

### Questions

#### Theoretical questions

1. What is cardinality in this section?
2. What do optimizer statistics estimate?
3. What can stale statistics do to a plan?
4. Why is an index on a low-selectivity boolean often useless for that filter?
5. What does a large gap between estimated and actual rows suggest?

#### Easy practical tasks

1. Estimate cardinality for: primary key, country code, boolean `is_deleted` on a live-only table.
2. Find the command that refreshes statistics in your DBMS.
3. Write one high-selectivity and one low-selectivity `WHERE` on `orders`.
4. Make a glossary: cardinality, selectivity, histogram, stale stats.

#### Medium practical tasks

1. Load thousands of rows. Explain a filter. Refresh statistics. Explain again. Write what changed.
2. Compare estimate and actual rows on a join if your `EXPLAIN` shows both.
3. Compute: 1 million rows, a status value that matches 90 percent. Write why a scan can win.

#### Advanced practical tasks

1. Write a one-page note: when you run `ANALYZE` (or equivalent) after loads, and how you detect estimate errors.
2. Find a correlated pair of columns (city and postal code). Write why single-column stats can misestimate a pair filter.

---

## Hot rows and lock contention

A hot row is a row that many transactions try to change at the same time. Example: a single counter row for "tickets remaining," or one parent order that every line update locks.

Lock contention is wait time on locks. Sessions queue. Throughput drops. CPU can look idle while sessions wait.

```text
Session A: UPDATE counters SET n = n - 1 WHERE id = 1;  -- holds row lock
Session B: same row  -->  wait
Session C: same row  -->  wait
```

Isolation and row locks are correct. The design is wrong if one row is a global bottleneck.

Mitigations (choose after you measure):

1. Split the counter into many rows and sum them.
2. Use an atomic increment that the product documents, if it reduces transaction time.
3. Move the hot fact out of the transactional table (queue, cache with a later reconcile) when the business allows a short delay.
4. Shorten the transaction. Do not hold the hot row lock while you call a network service.

Do not raise isolation to serializable as a first fix for a hot row. Stronger isolation can increase waits.

Do not add more application servers that all update the same row. You add more waiters.

Read lock wait metrics and the blocking session. A slow query that holds a lock is often the cause, not the waiter that you saw first.

Deadlocks can appear when two hot rows are updated in opposite order. Use a stable lock order. Retry on deadlock (earlier topic).

### Questions

#### Theoretical questions

1. What is a hot row?
2. What do waiting sessions do under lock contention?
3. Why does a stronger isolation level often make a hot row worse?
4. Why do more application servers not fix one-row contention?
5. Why must you not hold a row lock during an external HTTP call?

#### Easy practical tasks

1. Name two hot-row examples in a shop or a ticket system.
2. Draw three sessions and one locked row.
3. Write four mitigations from this section in one line each.
4. Find a lock-wait view or metric in your DBMS docs.

#### Medium practical tasks

1. Open two sessions. Update the same row in an uncommitted transaction in session A. Update in session B. Record the wait. Commit A.
2. Time 100 sequential updates to one counter versus 100 updates spread across 20 counters. Write the two times.
3. Rewrite a long transaction so that the hot update is the last statement. Write the before and after.

#### Advanced practical tasks

1. Write a one-page design: ticket inventory without one global row. Include a race that you still prevent with a transaction.
2. Capture a lock graph (product tool). Identify the blocker. Write the fix (index, shorter txn, or split row).

---

## Pagination (`OFFSET` vs keyset)

Pagination returns a page of rows. Two common methods exist.

**OFFSET pagination.** You request `LIMIT n OFFSET k`. The DBMS still walks or sorts k + n rows, then discards k. Large k is expensive. Example: page 1000 with page size 20 uses `OFFSET 19980`.

```sql
SELECT order_id, ordered_at
FROM orders
ORDER BY ordered_at DESC, order_id DESC
LIMIT 20 OFFSET 40;
```

**Keyset pagination** (seek method). You ask for rows after the last key of the previous page. The query uses a range on the sort key. Cost stays near one page.

```sql
SELECT order_id, ordered_at
FROM orders
WHERE (ordered_at, order_id) < (:last_at, :last_id)
ORDER BY ordered_at DESC, order_id DESC
LIMIT 20;
```

Keyset needs a unique sort key (add the primary key to break ties). Keyset does not jump to page 500 without extra work. OFFSET can jump (at a cost).

OFFSET also drifts: a new insert at the top can shift rows between pages. Keyset is stable for "next page" if you use the last seen key.

Do not use `OFFSET` for deep pages on large tables. Do not use `SELECT *` on a wide table for a list page.

An index that matches `ORDER BY` (and the keyset predicate) keeps the page cheap. Read the plan.

Total count (`SELECT COUNT(*)`) is a separate cost. Do not run a full count on every page if the UI can show "more" without a total.

### Questions

#### Theoretical questions

1. Why is a large `OFFSET` expensive?
2. What does a keyset page use instead of `OFFSET`?
3. Why must the keyset sort key be unique?
4. How can `OFFSET` pages drift when inserts occur?
5. Why is `COUNT(*)` on every page a separate problem?

#### Easy practical tasks

1. Write an `OFFSET` query for page 3 of size 10.
2. Write a keyset query that continues after a given `(ordered_at, order_id)`.
3. Make a two-column table: OFFSET versus keyset. Add four rows.
4. Write one UI that needs jump-to-page (OFFSET) and one that needs only Next (keyset).

#### Medium practical tasks

1. Fill a table with thousands of rows. Time `OFFSET 0`, `OFFSET 1000`, and `OFFSET 10000` with `LIMIT 20`. Write the times.
2. Time the keyset query for the same late page. Write the time and the plan access method.
3. Add an index that matches the `ORDER BY`. Explain both queries again.

#### Advanced practical tasks

1. Implement Next/Previous with keyset in a small program. Prove that a new insert does not duplicate a row on Next.
2. Write a one-page pagination standard for the team: default method, when OFFSET is allowed, and how you avoid `COUNT(*)`.

---

## Caching in front of the database

A cache stores query results or objects in a fast store (memory, key-value product). The application reads the cache first. On a miss, it reads the DBMS and fills the cache.

Caching helps when:

- the same read is frequent
- the result can be a few seconds or minutes old
- the DBMS is the bottleneck

```text
Read:  cache hit --> return
       cache miss --> DBMS --> fill cache --> return
Write: DBMS commit --> invalidate or update cache keys
```

Invalidation is the hard part. If you update a row and leave the old cache key, users see stale data. If you invalidate too many keys, the cache is useless and the DBMS sees a storm of misses (thundering herd).

Cache keys must include every input that changes the result (user, page, locale). A key that omits the user can leak data across users. That leak is a security bug.

Do not cache a write path that must stay consistent with constraints unless you still commit in the DBMS. The cache is not the system of record.

Do not add a cache to hide an N+1 query or a missing index. Fix the query first. Measure. Then cache if the read is still hot.

A short TTL (time to live) limits staleness without perfect invalidation. TTL is a trade-off, not a full design.

The DBMS buffer cache is not this cache. Buffer cache holds pages. Application cache holds results.

### Questions

#### Theoretical questions

1. When does an application cache help?
2. What must a write do to cached data?
3. What is a thundering herd after a large invalidation?
4. Why must a cache key include the user when the result is per user?
5. Why must you not use a cache as the system of record?

#### Easy practical tasks

1. Draw hit, miss, and invalidate for a product-price cache.
2. Write a cache key for "orders for customer 42, page after id 100."
3. List three reads that may use a 30-second TTL and two that must not.
4. Write one sentence that distinguishes application cache from buffer cache.

#### Medium practical tasks

1. Implement a tiny in-memory cache for `GET product` in a learning app. Invalidate on update. Show stale data if you skip invalidate.
2. Measure DBMS query count with and without the cache for 50 identical reads.
3. Write a stampede rule: single fill on miss (lock or request coalesce). Describe it in six sentences.

#### Advanced practical tasks

1. Write a one-page cache policy: keys, TTL, invalidation, personal data, and what you never cache.
2. Design cache keys for a join-heavy report. List every invalidation event. Mark keys that you cannot maintain (then do not cache).

---

## Connection storms

A connection storm is a sudden large number of new DBMS sessions. Causes:

- application restart of many instances at once
- pool size too large times instance count
- retry loops without backoff after a timeout
- health checks that open a full connection each time
- a cache miss storm that opens work on every app thread

Each connection uses memory and a process or thread on the DBMS (product model differs). Too many connections slow every session. New connections fail. The outage looks like "database down." The trigger was the clients.

```text
100 app instances * 50 pool = 5000 sessions  -->  server stalls
```

Mitigations:

1. Cap pool size from a budget, not from "number of threads."
2. Use a pooler in front of the DBMS if the product ecosystem has one.
3. Restart instances in batches.
4. Retry with jitter and a limit.
5. Fail fast in the application when the pool is empty.

Do not raise `max_connections` as the first fix. A higher cap without a pool budget repeats the storm at a higher memory cost.

Do not hide a slow query behind more connections. You add waiters.

Measure: connection count, connection create rate, wait for a pool slot, and DBMS memory.

A storm can follow a failover. All clients reconnect together. Test reconnect with a cap.

### Questions

#### Theoretical questions

1. What is a connection storm?
2. How do pool size and instance count multiply?
3. Why can retries without backoff create a storm?
4. Why is a higher `max_connections` a weak first fix?
5. Why can failover start a storm?

#### Easy practical tasks

1. Compute worst-case sessions: instances times max pool. Write the product.
2. List five causes from this section.
3. Write a retry rule: three tries, increasing delay, random jitter (describe, do not invent a library).
4. Find `max_connections` (or equivalent) in your DBMS.

#### Medium practical tasks

1. Set a tiny `max_connections` on a disposable instance. Open more clients than the cap. Record the error.
2. Write a pool budget: DBMS memory, bytes per connection (from docs or a guess marked as a guess), max instances.
3. Design a rolling restart of four app instances so that connection create rate stays bounded.

#### Advanced practical tasks

1. Write a one-page connection policy: pool formula, pooler, restart, retry, and alerts on connect rate.
2. Simulate a reconnect wave (script of many short connections) on a disposable instance. Record when the server degrades. Then add a cap in the script.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do a slow log, statistics, and a plan work together before you add an index or a cache?
2. When is a hot row the cause, and when is a deep `OFFSET` the cause, of a slow page?
3. How can a cache and a connection storm make each other worse?
4. Which measurements tell you to fix SQL, and which tell you to fix the client pool?
5. A teammate raises buffers, adds a cache, and doubles `max_connections` without a baseline. Which facts do you use in the reply?

#### Easy practical tasks

1. Write a cheat sheet: measure, cardinality, hot rows, keyset, cache, storms.
2. Enable or locate a slow log. Copy one line into your notes (sanitize secrets).
3. Write one keyset query and one OFFSET query for the same list.
4. Compute your app's worst-case connection count.

#### Medium practical tasks

1. Build a scratch workload: large table, stale stats, then `ANALYZE`, then a keyset page. Record plans and times.
2. Produce a lock wait with two sessions. Produce a slow OFFSET. Write how the metrics differ.
3. Write a tuning ticket template: baseline, hypothesis, one change, result.

#### Advanced practical tasks

1. Tune one real slow statement (or a made-slow one): measure, stats, index or rewrite, pagination fix. Write before/after.
2. Write an operations review: dashboards you need (lag, locks, connections, slow statements) and the first action for each alert.
