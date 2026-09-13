# 10. Replication and High Availability

## Description

Redis can copy data from a primary instance to one or more replicas. This topic covers `REPLICAOF`, asynchronous lag, Redis Sentinel, promotion, and stale reads from replicas.

Complete this topic before you study Cluster. Cluster is a sharding system. Sentinel is a failover system for one master (primary) and its replicas.

Use one term for each concept. A primary accepts writes. A replica copies the primary and can serve reads. Lag is the delay between a write on the primary and the same write on the replica. Sentinel is a set of processes that watch Redis and can promote a replica. Promotion is the change of a replica into a primary.

---

## Replica-of / `REPLICAOF`

`REPLICAOF host port` makes the current instance a replica of that primary. The replica connects, performs a sync, and then applies a stream of writes.

```text
REPLICAOF 127.0.0.1 6379
```

`REPLICAOF NO ONE` stops replication. The instance becomes a standalone primary with the data it already has.

Older docs and commands used `SLAVEOF`. Current Redis accepts `REPLICAOF`. Use `REPLICAOF` in new work. `INFO replication` still shows fields such as `role:master` and `role:slave` on many versions. Read the field names that your version prints.

A first sync often sends a snapshot (RDB) plus a backlog of commands. The replica loads the snapshot and then catches up. A large dataset makes the first sync heavy (CPU, disk, network).

By default a replica is read-only (`replica-read-only yes`). Writes on the replica fail. That is what you want. Do not write to a replica.

A replica can persist to disk. Persistence on the replica is independent of the primary. A replica is not a backup by itself until you copy files off the host.

Topology for this topic: one primary, one or more replicas. Do not chain many replicas unless you know the extra lag. Redis also supports replica-of-replica in some setups. Prefer a simple tree for learning.

`INFO replication` on the primary lists connected replicas. On the replica it shows the primary host, link status, and offset fields.

Authentication: the replica must use the password or user that the primary requires (`masteruser` / `masterauth` in config). TLS has extra settings.

### Questions

#### Theoretical questions

1. What does `REPLICAOF host port` do?
2. What does `REPLICAOF NO ONE` do?
3. Why should you not write to a replica?
4. What extra work happens on the first sync of a large dataset?
5. Is a replica a backup if you never copy files away?

#### Easy practical tasks

1. Run `INFO replication` on your single instance. Write `role`.
2. Open the replication docs. Write the modern command name and the old name.
3. Write four sentences: primary vs replica.
4. Find `replica-read-only` with `CONFIG GET`. Save the value.

#### Medium practical tasks

1. Start two Redis containers. Make B a replica of A. `SET` on A. `GET` on B. Record both.
2. Run `INFO replication` on A and B. Write `role` and one offset field from each.
3. Run `REPLICAOF NO ONE` on B. `SET` on B. Confirm A does not see that write.

#### Advanced practical tasks

1. Watch the first sync of a dataset with 100,000 keys. Time the sync. Write `master_sync_in_progress` observations.
2. Configure `masterauth` on a lab primary with a password. Connect a replica. Document the settings. Use only a private lab.

---

## Asynchronous replication and lag

Redis replication is asynchronous in the usual setup. The primary sends `OK` to the client when the write is in the primary memory (and after the primary persistence rules). The primary does not wait for the replica by default.

The replica receives the write later. The delay is lag. Lag grows when the replica is slow, the network is slow, or the replica is busy loading a snapshot.

`WAIT numreplicas timeout` on the primary can wait until at least `numreplicas` replicas acknowledge the write, or until the timeout. `WAIT` improves the chance that a replica has the write. `WAIT` is not a full synchronous commit across failures. Read the `WAIT` page for limits (and `WAITAOF` on Redis 7.2+).

`INFO replication` fields such as `master_repl_offset` and the replica's `slave_repl_offset` (names vary) help you see divergence. Sentinel and monitoring tools also track lag.

If the primary dies, a replica that you promote can miss the last writes. That is the cost of async replication. The application must accept that loss or use `WAIT` and still accept edge cases.

Do not assume that a read on a replica sees the write that the same user just sent to the primary. That is stale read (later section).

Disk persistence on the primary and replication to a replica are two paths. You can lose a write on crash of the primary even if a replica would have gotten it a moment later. You can also persist on the primary and still have a lagging replica.

Keep the replica in the same region for low lag if reads must be fresh enough. Measure lag under load.

### Questions

#### Theoretical questions

1. What does "asynchronous replication" mean for the client's `OK`?
2. What is lag?
3. What does `WAIT` ask the primary to do?
4. Can a promoted replica miss the last writes?
5. Why are persistence and replication not the same path?

#### Easy practical tasks

1. Write four sentences that explain async replication to a beginner.
2. Open the `WAIT` command page. Write the arguments and the return value.
3. Draw primary, replica, client, and two arrows (reply vs replicate).
4. List three causes of high lag.

#### Medium practical tasks

1. In a two-instance lab, `SET` then `WAIT 1 500` on the primary. Write the return value.
2. Slow the replica (CPU limit or a big `DEBUG SLEEP` only on a lab replica). Watch lag fields. Remove the slowdown.
3. Compare `INFO` offsets on primary and replica after 10,000 writes. Write the two numbers.

#### Advanced practical tasks

1. Read `WAIT` / `WAITAOF` limits in the docs. Write three failure cases where `WAIT` is not enough.
2. Script a lag monitor: print offset difference every second during a write burst. Save a short log.

---

## Redis Sentinel (monitoring + failover)

Redis Sentinel is a separate process (or several processes). Sentinel is not Cluster. Sentinel does not shard keys. Sentinel watches a named primary and its replicas.

A typical lab uses three Sentinel processes. They vote. A majority must agree that the primary is down. One Sentinel then starts failover: pick a replica, promote it, reconfigure other replicas, and update the primary name.

Clients ask Sentinel for the current primary address (`SENTINEL get-master-addr-by-name mymaster`). After failover, that address changes. Clients must support Sentinel or use a proxy that does.

Sentinel monitors:

- whether the primary answers
- replica state
- some configuration

Sentinel does not make replication synchronous. Failover can drop the last async writes.

You configure `sentinel monitor <name> <host> <port> <quorum>`. Quorum is the number of Sentinels that must agree. Three Sentinels and quorum 2 is a common learning setup.

Sentinel needs reachable network paths to Redis and to other Sentinels. Split brain is possible if networks partition in bad ways. Read the official Sentinel docs before production.

Do not run a single Sentinel and call that high availability. One Sentinel is a single watcher.

`INFO` on Redis is not enough to see Sentinel. Use `redis-cli -p 26379 SENTINEL masters` (default Sentinel port is `26379`).

### Questions

#### Theoretical questions

1. Does Sentinel shard data across slots?
2. Why do teams run at least three Sentinels?
3. What does quorum mean?
4. How does a client find the primary after failover?
5. Does Sentinel prevent loss of the last async writes?

#### Easy practical tasks

1. Write four sentences: Sentinel vs one primary with no watcher.
2. Open the Sentinel docs. Write the default port and one `SENTINEL` command.
3. Draw three Sentinels, one primary, two replicas.
4. Explain why one Sentinel is not enough (two sentences).

#### Medium practical tasks

1. Start a primary, a replica, and three Sentinels in Docker (follow an official or well-known compose example). Run `SENTINEL masters`. Save a summary.
2. Read `down-after-milliseconds` and `failover-timeout`. Write what each does.
3. Point `redis-cli` at Sentinel and resolve the primary host and port.

#### Advanced practical tasks

1. Perform a planned failover (`SENTINEL failover`). Record the new primary and the time until a `SET` works.
2. Read about Sentinel split-brain and `min-replicas-to-write`. Write a one-page risk note.

---

## Promotion of a replica

Promotion turns a replica into a primary. After promotion, that instance accepts writes. Other replicas should copy from the new primary.

Manual promotion (lab, no Sentinel):

1. Confirm the replica is as fresh as you can check (`INFO replication`)
2. `REPLICAOF NO ONE` on the chosen replica
3. Point applications to the new primary
4. `REPLICAOF new-primary` on the other replicas
5. Keep the old primary off or rebuild it as a replica to avoid two writers

If the old primary comes back as a primary, you have two writers (split brain). Data diverges. Sentinel tries to avoid this by reconfiguring the old primary as a replica. Manual operations must be careful.

Sentinel promotion picks a replica with rules (priority, offset, run id). `replica-priority 0` means "never promote this replica". Use that for a replica that is only a backup copy or is in a far region.

Promotion does not replay lost async writes from the dead primary if those writes never reached the replica. Clients that wrote and received `OK` can still lose those keys.

After promotion, replicas of the old primary need the new topology. Clients that cached the old IP must refresh (Sentinel, DNS, or a proxy).

Test promotion in a lab. Measure downtime (time when writes fail). Measure data loss with a write counter.

Do not promote a replica that is far behind because "it is the only one that is up" without a decision. You can lose a large suffix of writes.

### Questions

#### Theoretical questions

1. Which command makes a replica a standalone primary?
2. What is split brain after a failed primary returns?
3. What does `replica-priority 0` mean?
4. Can promotion restore writes that never reached the replica?
5. Why must clients refresh the primary address after promotion?

#### Easy practical tasks

1. Write the five manual promotion steps in your own words.
2. Make a table: action, command or tool (manual vs Sentinel).
3. Write four sentences on why two primaries are dangerous.
4. Find `replica-priority` in the docs. Write the default meaning.

#### Medium practical tasks

1. In a two-node lab, promote the replica by hand. `SET` a new key. Rebuild the old node as a replica. `GET` the key on both.
2. Simulate the old primary coming back without reconfiguration (start it as primary). Show two different values for the same key. Tear down the lab.
3. Set `replica-priority 0` on one replica. Read how Sentinel would treat it.

#### Advanced practical tasks

1. Automate a failover drill: write a counter, kill the primary, wait for Sentinel, read the counter. Report loss and downtime.
2. Write a runbook: who promotes, how to fence the old primary, how to tell clients, how to verify `INFO replication`.

---

## Read from replicas (stale data)

You can send read commands to a replica. The replica answers from its copy. That copy can be behind the primary. The application can see old values.

This is a stale read. Stale reads are acceptable for some dashboards, leaderboards, and cache fills. They are not acceptable for "read your own write" (a user saves a profile and the next page still shows the old city).

Patterns:

- All reads and writes to the primary — simple, no stale replica reads
- Writes to the primary, reads to replicas — scale read traffic, accept staleness
- Read your write from the primary, other reads from replicas — mixed

The client must know which connection is the primary. After failover, the old primary can become a replica. A write to the old address should fail if the node is now read-only. A read to a lagging replica can still be stale.

`WAIT` after a write can reduce the chance that a replica read is stale, if you then read that replica. It does not fix every race.

Replicas can have different lag. Do not load-balance reads and assume all replicas have the same data at the same time.

Monitor replica lag. Remove a replica from the read pool when lag exceeds a budget (for example 1 second).

`READONLY` is a Cluster command for replicas. On Sentinel topologies, `replica-read-only` is the usual control. Do not confuse the two.

### Questions

#### Theoretical questions

1. Why can a replica `GET` return an old value?
2. When is a stale read acceptable?
3. What is "read your own write"?
4. Why must a read pool drop a high-lag replica?
5. How does failover make a stale or failed write more likely if the client caches IPs?

#### Easy practical tasks

1. Write four sentences: replica read vs primary read.
2. Draw a user request that writes the primary and reads a replica. Mark the stale path.
3. List three features in an app that can use replica reads and two that must not.
4. Open docs on replica reads. Write one warning that the page gives.

#### Medium practical tasks

1. In a lab, `SET` on the primary and immediately `GET` on the replica in a tight loop. Record if you ever see a miss or old value.
2. Add `WAIT 1 100` after `SET`. Repeat the tight `GET` on the replica. Compare.
3. Configure a client (or a sketch) with a lag threshold. Write the `INFO` fields you would parse.

#### Advanced practical tasks

1. Build a small app: write primary, read replica, show a "may be stale" banner when offsets differ.
2. Read about `min-replicas-to-write` and `min-replicas-max-lag`. Write how they protect writes, not reads.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do `REPLICAOF`, async lag, and promotion combine when the primary dies?
2. What does Sentinel add that a manual `REPLICAOF NO ONE` does not add?
3. How do you explain data loss to a product owner after a failover with `everysec` AOF and async replicas?
4. When do you send reads to replicas, and how do you bound staleness?
5. Why is Sentinel the wrong tool if you need to split keys across many machines?

#### Easy practical tasks

1. Write a cheat sheet: `REPLICAOF`, `REPLICAOF NO ONE`, `WAIT`, `INFO replication`, `SENTINEL get-master-addr-by-name`.
2. Draw the happy path (one primary, two replicas) and the path after promotion.
3. Export `INFO replication` from a lab pair. Mark role, host, and offset lines.
4. Write five lab rules: no writes on replicas, three Sentinels, test failover, measure lag, fence old primary.

#### Medium practical tasks

1. Document a Docker topology (compose service names, ports 6379 and 26379) so that another beginner can start it.
2. Run a write load on the primary and a read load on the replica. Record lag and one stale-read example if you see one.
3. Write a client connection policy: Sentinel vs static replica list, timeouts, and retry after failover.

#### Advanced practical tasks

1. Kill the primary, wait for Sentinel, write, then start the old primary. Prove there is one primary and that the old node is a replica.
2. Compare Sentinel HA with "two primaries and application dual-write". Write six sentences on why dual-write is harder.
