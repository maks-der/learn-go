# 16. Replication and High Availability

## Description

This topic shows how a DBMS copies changes to another server and how you keep service when one server fails. You learn primary and replica roles, synchronous and asynchronous replication, read scaling, replication lag, failover, split-brain, and multi-primary designs.

Use one term for each concept. A primary (leader) is the server that accepts writes. A replica (follower) applies those writes. High availability is a design that keeps service after a planned or unplanned stop of one node. Complete this topic after you understand WAL and backups. Replication is not a backup.

This path is vendor-neutral. Product terms differ (primary/standby, master/replica). This handbook uses primary and replica.

---

## Primary / replica (leader / follower)

A primary is the DBMS instance that accepts writes. Applications send `INSERT`, `UPDATE`, and `DELETE` to the primary. The primary writes WAL (or an equivalent stream).

A replica is an instance that receives that stream and applies the changes. A typical replica is read-only. You can send `SELECT` to the replica. A write on a read-only replica fails.

```text
Application writes  -->  primary  -->  log stream  -->  replica
Application reads   -->  primary or replica
```

Roles are not hardware labels. You can promote a replica to primary. After promote, the old primary must not keep the write role.

A replica starts from a base copy of the primary data. Then it follows the log. If the replica stops, it must catch up from the stored log position. If the gap is too large, you rebuild the replica from a new base backup.

Use cases for a replica:

- extra read capacity
- a warm standby for failover
- a reporting instance so that heavy `SELECT` does not hit the primary
- a source for backups (with product rules)

Do not write to a replica unless the product documents a special session. Do not assume that a replica has the same indexes or the same parameters unless you manage them.

A logical replica can apply row changes to a different instance and sometimes to a subset of tables. A physical replica follows pages or WAL closely. Physical replicas stay on the same major version in most products. Learn which kind you run.

Do not put the only replica on the same host as the primary. A host failure then stops both.

### Questions

#### Theoretical questions

1. Which instance accepts writes in a primary/replica pair?
2. What does a typical replica reject?
3. What two steps does a new replica need (base copy, then log)?
4. Name three use cases for a replica.
5. Why must the replica not share the only host with the primary?

#### Easy practical tasks

1. Draw primary, log stream, replica, and two arrows for write and read.
2. Find the product terms: primary, standby, publisher, or another pair. Write the mapping.
3. Write one `SELECT` you would send to a replica and one `INSERT` you would send only to the primary.
4. List two reasons a replica must rebuild from a new base backup.

#### Medium practical tasks

1. Read the official streaming-replica tutorial. Write the order of init, start, and attach.
2. If you can, start a primary and a replica with Docker Compose. Run `SELECT` on both. Run `INSERT` on the primary. Read the row on the replica.
3. Compare physical and logical replication in your manual. Write three differences.

#### Advanced practical tasks

1. Write a one-page role note: who writes, who reads, how you promote, and what you do with the old primary.
2. Break the replica (stop the process). Take more writes on the primary. Start the replica. Record how it catches up or why you rebuild.

---

## Synchronous vs asynchronous replication

Asynchronous replication means the primary commits locally, then sends the change to the replica. `COMMIT` on the primary does not wait for the replica. The replica can lag. If the primary fails before the replica received the last records, those commits can be lost on failover.

Synchronous replication means `COMMIT` waits until a replica (or a quorum) has the change at a defined level: received, flushed to disk, or applied. The exact level is a product setting.

| Mode | Commit wait | Typical lag | Failover data risk |
| --- | --- | --- | --- |
| Asynchronous | Primary only | Milliseconds to seconds or more | Last commits can be missing on the replica |
| Synchronous | Primary plus replica ack | Low if the replica is healthy | Lower loss; commit slower if the replica is slow |

```text
Async:  client <-- COMMIT -- primary      replica (later)
Sync:   client <-- COMMIT -- primary <-- ack -- replica
```

Synchronous replication reduces RPO toward zero for the protected commits. It increases commit latency. If the synchronous replica stops, writes can stall unless you have a fallback policy. That fallback can drop you back to asynchronous behavior. Document the policy.

Asynchronous replication is the usual default. It fits read replicas and distant sites. It does not give "zero data loss" on primary crash.

Do not mix the words. "We replicate" does not say whether commit waits. Name the mode and the ack level.

Do not put a synchronous replica on a slow or distant link if the application needs short write latency. Measure commit time.

A common pattern: one synchronous replica in the same site for failover, asynchronous replicas in other sites for reads and disaster recovery. That pattern is a design, not a default.

### Questions

#### Theoretical questions

1. When does an asynchronous primary return `COMMIT` to the client?
2. What can you lose if the primary fails in asynchronous mode?
3. What does a synchronous commit wait for?
4. What happens to writes if the only synchronous replica stops?
5. Why is "we replicate" an incomplete durability statement?

#### Easy practical tasks

1. Fill a two-row table: async versus sync. Add commit wait and failover risk.
2. Find the setting name for synchronous replication in your DBMS.
3. Write one use case for async and one use case for sync.
4. Draw the two commit paths from this section.

#### Medium practical tasks

1. Read the ack levels (received, flushed, applied). Write one sentence each.
2. Measure or estimate: a 50 ms network to a sync replica. Write the effect on commit time.
3. Write a fallback policy in five sentences for a dead synchronous replica.

#### Advanced practical tasks

1. On a disposable pair, compare commit time with async and with sync if the product allows both. Record the method and the numbers.
2. Write a one-page HA design: same-site sync replica, off-site async replica. State RPO for two failure types.

---

## Read scaling and replication lag

Read scaling means you send some `SELECT` traffic to replicas. The primary does fewer reads. Write traffic still goes to the primary. You scale reads, not writes, with a typical primary/replica pair.

Replication lag is the delay between a commit on the primary and the moment the replica can see that commit. Lag has two common measures:

- time (seconds behind)
- log position (bytes or sequence difference)

```text
Primary commits at t=0
Replica applies at t=lag
A read at t < lag can miss the new row
```

Lag grows when the replica is slow, when the network is slow, or when the replica runs heavy queries. Lag also grows after the replica restarts and must catch up.

Applications that write and then immediately read their own write must read the primary (or wait until the replica catches up). Example: a user saves a profile and reloads the page. A replica read can show the old profile.

Load balancers that send all reads to replicas without a "read your writes" rule create stale-read bugs. Those bugs look random.

Do not count replica CPU as write capacity. `INSERT` still serializes on the primary in this design.

Do not use a replica with a large lag for a dashboard that must match the primary. Use the primary, or wait, or show a "data as of" time.

Monitor lag. Alert when lag exceeds a budget. A replica that is 1 hour behind is a standby with a 1-hour RPO if you promote it.

### Questions

#### Theoretical questions

1. Which traffic can replicas reduce on the primary?
2. What is replication lag?
3. Why can a write-then-read on a replica miss the new row?
4. Why do replicas not increase write throughput in this design?
5. What RPO do you accept if you promote a replica that is 1 hour behind?

#### Easy practical tasks

1. Write two queries: one that must hit the primary, one that can hit a replica.
2. Find the lag metric name in your DBMS (example: replay lag).
3. Draw write, commit, apply, and a stale `SELECT`.
4. Write one sentence that you would show in a UI when you read a lagged replica.

#### Medium practical tasks

1. If you have a replica, insert a row on the primary. Select on the replica in a loop until it appears. Write the delay.
2. Run a heavy `SELECT` on the replica. Watch lag if a metric exists. Write what you see.
3. Design a routing rule: session after POST goes to primary for N seconds. Write the rule in six sentences.

#### Advanced practical tasks

1. Write a one-page read-routing policy: which endpoints use replicas, how you handle read-your-writes, and the lag alert.
2. Build a small demo that fails when it reads a replica immediately after write. Then fix the route. Record both behaviors.

---

## Failover and split-brain (high-level)

Failover is the change of role when the primary stops. A replica becomes the new primary. Clients must connect to the new primary.

Failover can be planned (switchover): you stop writes, wait for the replica to catch up, promote, then point clients. Failover can be unplanned: the primary crashes or the network fails. A manager process or an operator promotes a replica.

```text
Primary down --> choose replica --> promote --> redirect clients --> fence old primary
```

Fencing means the old primary must not accept writes after promote. If the old primary is alive but isolated, it can still accept writes. Two primaries that accept writes is split-brain.

Split-brain produces two histories. Merging them is hard. Unique keys collide. Counters diverge. The usual repair is to pick one history and discard or manually merge the other.

Prevention:

1. Use a cluster manager that gets a lock or a quorum.
2. Fence the old primary (STONITH, revoked disks, or forced read-only).
3. Do not let applications keep a cached write address without a retry policy that discovers the new primary.

Do not promote two replicas. Do not start the old primary as a primary again without a rewind or a rebuild.

DNS or a virtual IP can hide the new host. Clients still need a timeout and a reconnect. Long-lived connections to the dead primary fail. Pools must drop those connections.

Failover time is part of RTO. Test failover as you test restore. An untested promote is a hope.

### Questions

#### Theoretical questions

1. What role change does failover perform?
2. What is the difference between planned switchover and unplanned failover?
3. What is split-brain?
4. Why must you fence the old primary?
5. Why is an untested promote a hope?

#### Easy practical tasks

1. Number the failover steps from this section.
2. Write one sentence that defines fencing.
3. List three ways clients can find the new primary (VIP, DNS, connection string).
4. Write one data symptom of two primaries (example: two rows with the same natural key in different histories).

#### Medium practical tasks

1. Read the official failover or promote command. Write the command and the fence warning from the docs.
2. Write a client reconnect policy: timeout, retry, drop pool connections.
3. Describe a network partition: primary in site A, replica and manager in site B. Write who should win and why.

#### Advanced practical tasks

1. On a disposable pair, promote a replica. Try to write to the old primary. Record whether the product fences it. Write the gap if it does not.
2. Write a one-page incident runbook: detect down primary, promote, fence, check lag first, tell applications.

---

## Multi-primary (when it is hard)

A multi-primary (multi-leader) system accepts writes on more than one primary. Each primary replicates to the others. Clients can write in more than one site.

This design is hard because two primaries can change the same row at the same time. The system must resolve the conflict.

Common conflict cases:

- two updates to the same key
- insert of the same unique value in two sites
- increment of the same counter in two sites

Resolution methods exist: last-write-wins, application merge, conflict tables, or reserved key ranges (site A uses even ids, site B uses odd ids). Each method can lose a write or require custom code.

```text
Site A: UPDATE balance = 100
Site B: UPDATE balance = 80
-- both commit locally, then replicate
-- one value wins, or you merge by hand
```

Synchronous multi-primary across distant sites adds latency or availability loss (see CAP in a later topic). Asynchronous multi-primary adds conflict windows.

Most applications do not need multi-primary. A single primary plus replicas is simpler. Use multi-primary when a site must keep writes while the other site is down, and the business accepts conflict rules.

Do not treat "active-active" as a slogan. Ask: which keys can collide, what wins, and how you test a partition.

Sharding (a later topic) splits keys so that each key has one write owner. That design is not the same as two primaries for the same key.

### Questions

#### Theoretical questions

1. What extra ability does multi-primary add?
2. Why do two primaries create conflicts?
3. Name three conflict cases from this section.
4. What does last-write-wins risk?
5. When is a single primary plus replicas a better default?

#### Easy practical tasks

1. Write one story of a conflicting update on `customers.email`.
2. List three resolution methods. Give one cost each.
3. Draw two sites, two writes to the same row, and a merge cloud.
4. Write one sentence that distinguishes multi-primary from sharding with one writer per key.

#### Medium practical tasks

1. Read a product page on logical multi-primary or a conflict example. Write the resolution rule in your own words.
2. Design even/odd keys for two sites. Write how a sequence still collides if you ignore the rule.
3. Compare failover (one writer at a time) with multi-primary (two writers). Make a five-row table of operational costs.

#### Advanced practical tasks

1. Write a one-page decision record: reject multi-primary for a shop, or accept it for a specific table with a merge rule.
2. Table-top a 30-minute partition with writes on both sides. List the rows you must merge. Do not run this on production.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do primary role, sync mode, and lag determine what a promote will lose?
2. When do you add a replica for reads, and when must those reads still go to the primary?
3. How do fencing and a single-writer rule prevent the conflict problems of multi-primary?
4. Which problems does replication solve that a backup does not solve, and the reverse, in an HA design?
5. A teammate wants two writable primaries in two cities and "no lag." Which facts do you use in the reply?

#### Easy practical tasks

1. Write a cheat sheet: roles, sync versus async, lag, failover, split-brain, multi-primary.
2. Label your learning setup as single node or primary/replica. Write the write target.
3. Write RPO for async promote and for sync promote in one line each (assume the replica is current).
4. Draw a forbidden state: two primaries, one client pool to each.

#### Medium practical tasks

1. Write an HA diagram for a class shop: primary, one replica, backup off-host. Mark write path, read path, and backup path.
2. Write a promote checklist of eight items that includes lag check and fencing.
3. Choose sync or async for (a) same rack standby and (b) another continent report replica. Justify each.

#### Advanced practical tasks

1. If you can, run a planned switchover on a disposable pair. Measure downtime and whether the last write survived. Write the numbers.
2. Write a full HA note: modes, read routing, failover, split-brain, why you reject multi-primary. Review it against official docs.
