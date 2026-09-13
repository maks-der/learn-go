# 8. Replication, Sentinel, and Cluster

## Description

Redis can copy data from a primary instance to replicas. Redis can also split keys across many primaries. This topic covers asynchronous lag, Redis Sentinel, failover, Cluster hash slots, hash tags, multi-key limits, and what a client must do in cluster mode.

Complete the earlier topics first. Cluster is not Sentinel. Cluster shards data. Sentinel fails over one primary.

Use one term for each concept. A primary accepts writes. A replica copies the primary and can serve reads. Lag is the delay between a write on the primary and the same write on the replica. Sentinel is a set of processes that watch Redis and can promote a replica. A slot is an integer from `0` to `16383`. A hash tag is the `{...}` part of a key that Redis uses to compute the slot. A redirect is a `MOVED` or `ASK` reply that tells the client another node.

---

## Asynchronous replicas and lag

`REPLICAOF host port` makes the current instance a replica of that primary. The replica connects, performs a sync, and then applies a stream of writes.

```text
REPLICAOF 127.0.0.1 6379
```

`REPLICAOF NO ONE` stops replication. The instance becomes a standalone primary with the data it already has.

Older docs and commands used `SLAVEOF`. Current Redis accepts `REPLICAOF`. Use `REPLICAOF` in new work. `INFO replication` still shows fields such as `role:master` and `role:slave` on many versions.

A first sync often sends a snapshot (RDB) plus a backlog of commands. The replica loads the snapshot and then catches up. A large dataset makes the first sync heavy (CPU, disk, network).

By default a replica is read-only (`replica-read-only yes`). Writes on the replica fail. That is what you want. Do not write to a replica.

A replica can persist to disk. Persistence on the replica is independent of the primary. A replica is not a backup by itself until you copy files off the host.

Redis replication is asynchronous in the usual setup. The primary sends `OK` to the client when the write is in the primary memory (and after the primary persistence rules). The primary does not wait for the replica by default.

The replica receives the write later. The delay is lag. Lag grows when the replica is slow, the network is slow, or the replica is busy loading a snapshot.

`WAIT numreplicas timeout` on the primary can wait until at least `numreplicas` replicas acknowledge the write, or until the timeout. `WAIT` improves the chance that a replica has the write. `WAIT` is not a full synchronous commit across failures. Read the `WAIT` page for limits (and `WAITAOF` on Redis 7.2+).

If the primary dies, a replica that you promote can miss the last writes. That is the cost of async replication. The application must accept that loss or use `WAIT` and still accept edge cases.

Do not assume that a read on a replica sees the write that the same user just sent to the primary. That is a stale read. Stale reads are acceptable for some dashboards. They are not acceptable for "read your own write".

`INFO replication` on the primary lists connected replicas. On the replica it shows the primary host, link status, and offset fields.

Disk persistence on the primary and replication to a replica are two paths. You can lose a write on crash of the primary even if a replica would have gotten it a moment later.

Monitor replica lag. Remove a replica from the read pool when lag exceeds a budget (for example 1 second).

### Questions

#### Theoretical questions

1. What does `REPLICAOF host port` do?
2. What does "asynchronous replication" mean for the client's `OK`?
3. What is lag?
4. What does `WAIT` ask the primary to do?
5. Can a promoted replica miss the last writes?

#### Easy practical tasks

1. Run `INFO replication` on your single instance. Write `role`.
2. Write four sentences that explain async replication to a beginner.
3. Draw primary, replica, client, and two arrows (reply vs replicate).
4. Find `replica-read-only` with `CONFIG GET`. Save the value.

#### Medium practical tasks

1. Start two Redis containers. Make B a replica of A. `SET` on A. `GET` on B. Record both.
2. In a two-instance lab, `SET` then `WAIT 1 500` on the primary. Write the return value.
3. Compare `INFO` offsets on primary and replica after 10,000 writes. Write the two numbers.

#### Advanced practical tasks

1. Read `WAIT` / `WAITAOF` limits in the docs. Write three failure cases where `WAIT` is not enough.
2. Watch the first sync of a dataset with 100,000 keys. Time the sync. Write `master_sync_in_progress` observations.

---

## Sentinel and failover

Redis Sentinel is a separate process (or several processes). Sentinel is not Cluster. Sentinel does not shard keys. Sentinel watches a named primary and its replicas.

A typical lab uses three Sentinel processes. They vote. A majority must agree that the primary is down. One Sentinel then starts failover: pick a replica, promote it, reconfigure other replicas, and update the primary name.

Clients ask Sentinel for the current primary address (`SENTINEL get-master-addr-by-name mymaster`). After failover, that address changes. Clients must support Sentinel or use a proxy that does.

Sentinel does not make replication synchronous. Failover can drop the last async writes.

You configure `sentinel monitor <name> <host> <port> <quorum>`. Quorum is the number of Sentinels that must agree. Three Sentinels and quorum 2 is a common learning setup.

Do not run a single Sentinel and call that high availability. One Sentinel is a single watcher.

`INFO` on Redis is not enough to see Sentinel. Use `redis-cli -p 26379 SENTINEL masters` (default Sentinel port is `26379`).

Promotion turns a replica into a primary. After promotion, that instance accepts writes. Other replicas should copy from the new primary.

Manual promotion (lab, no Sentinel):

1. Confirm the replica is as fresh as you can check (`INFO replication`)
2. `REPLICAOF NO ONE` on the chosen replica
3. Point applications to the new primary
4. `REPLICAOF new-primary` on the other replicas
5. Keep the old primary off or rebuild it as a replica to avoid two writers

If the old primary comes back as a primary, you have two writers (split brain). Data diverges. Sentinel tries to avoid this by reconfiguring the old primary as a replica.

`replica-priority 0` means "never promote this replica". Use that for a replica that is only a backup copy or is in a far region.

Promotion does not replay lost async writes from the dead primary if those writes never reached the replica. Clients that wrote and received `OK` can still lose those keys.

After promotion, clients that cached the old IP must refresh (Sentinel, DNS, or a proxy).

Test promotion in a lab. Measure downtime. Measure data loss with a write counter.

Split brain is possible if networks partition in bad ways. Read the official Sentinel docs before production.

### Questions

#### Theoretical questions

1. Does Sentinel shard data across slots?
2. Why do teams run at least three Sentinels?
3. What does quorum mean?
4. How does a client find the primary after failover?
5. Which command makes a replica a standalone primary?

#### Easy practical tasks

1. Write four sentences: Sentinel vs one primary with no watcher.
2. Open the Sentinel docs. Write the default port and one `SENTINEL` command.
3. Draw three Sentinels, one primary, two replicas.
4. Write the five manual promotion steps in your own words.

#### Medium practical tasks

1. Start a primary, a replica, and three Sentinels in Docker (follow an official or well-known compose example). Run `SENTINEL masters`. Save a summary.
2. Point `redis-cli` at Sentinel and resolve the primary host and port.
3. In a two-node lab, promote the replica by hand. `SET` a new key. Rebuild the old node as a replica. `GET` the key on both.

#### Advanced practical tasks

1. Perform a planned failover (`SENTINEL failover`). Record the new primary and the time until a `SET` works.
2. Automate a failover drill: write a counter, kill the primary, wait for Sentinel, read the counter. Report loss and downtime.

---

## Cluster hash slots and hash tags

Redis Cluster splits keys across many primary instances. Redis Cluster divides the keyspace into 16384 slots. Every key maps to one slot. Redis computes `CRC16(key) mod 16384` (with hash-tag rules). The result is the slot.

Each primary owns a set of slots. Together the primaries own all 16384 slots. If a slot has no owner, the cluster is not healthy for that slot.

A write of a key goes to the primary that owns the key's slot. Other primaries do not store that key. This is sharding. Memory and CPU scale out. A single key still lives on one primary. A huge key does not split.

```text
CLUSTER KEYSLOT user:42
CLUSTER SLOTS
CLUSTER NODES
```

`CLUSTER KEYSLOT` works on a cluster node and shows the slot for a key.

The cluster bus (a second port, often `16379` when the client port is `6379`) carries gossip and fail-over messages between nodes. Clients use the client port.

Logical DBs do not exist in Cluster. `SELECT` is not available. Use key prefixes.

A cluster needs several primaries (three is a common minimum for production). Each primary can have replicas for failover inside the cluster. That failover is not Sentinel. The cluster nodes vote.

`INFO cluster` and `CLUSTER INFO` show `cluster_state`. `ok` means slots are covered. `fail` means the cluster is not ready.

If the key contains a `{...}` substring, Redis uses only the characters between the first `{` and the matching `}` to compute the slot. That substring is the hash tag.

```text
{user:42}:profile
{user:42}:cart
{user:42}:orders
```

These three keys map to the same slot. They can live on the same primary. Multi-key commands that require one slot can run on them.

Rules:

- The first `{` starts the tag
- The first `}` after that ends the tag
- If `{}` is empty, Redis ignores the tag and hashes the full key
- If there is no `{`, Redis hashes the full key

Hash tags are a deliberate hot-spot tool. If all keys use `{global}`, one slot and one primary take all traffic. That defeats sharding.

Use hash tags when a transaction, a Lua script, or `MGET` must touch several keys. Do not tag everything.

Cluster-safe Lua must access only keys in `KEYS[]`, and those keys must share a slot (same tag in practice).

Resharding moves slots from one primary to another. You do this to add a primary, to remove a primary, or to balance memory and CPU. During migration a slot can be in a transitional state. Clients can receive `ASK` redirects. `ASK` means "try that node once for this key". `MOVED` means "update your slot map; that node owns the slot".

Resharding does not split a single large key. If one key is too big, you must change the data model.

Removing a primary requires moving its slots away first. Do not kill a primary that still owns slots.

### Questions

#### Theoretical questions

1. How many hash slots does Redis Cluster use?
2. How does Redis map a key to a slot (high level)?
3. Which part of `{user:42}:cart` does Redis hash?
4. Why can `{user:42}:a` and `{user:42}:b` work in one `MGET`?
5. What is the difference between `ASK` and `MOVED`?

#### Easy practical tasks

1. Write four sentences: Cluster vs one Redis instance.
2. Write five key names for one user that share a tag.
3. Draw three primaries and boxes of slots (for example 0–5460, 5461–10922, 10923–16383).
4. Open the Cluster spec or tutorial. Write the number 16384 and one reason the page gives (if any).

#### Medium practical tasks

1. Start a cluster (Docker or `redis-cli --cluster create`) or use a hosted cluster. Run `CLUSTER INFO`. Write `cluster_state` and `cluster_slots_assigned`.
2. On a cluster, run `CLUSTER KEYSLOT` for `{user:1}:a`, `{user:1}:b`, and `{user:2}:a`. Write which pair matches.
3. Show an empty-tag key `foo{}bar` vs `foobar` with `CLUSTER KEYSLOT`.

#### Advanced practical tasks

1. Read how CRC16 and slots work in the official spec. Write the formula and the hash-tag exception in six sentences.
2. On a lab cluster with at least two primaries, reshard a few slots. Run `CLUSTER KEYSLOT` and `CLUSTER NODES` before and after for one key.

---

## Multi-key commands across slots

Many commands accept more than one key: `MGET`, `MSET`, `DEL key1 key2`, `SUNION`, `SINTER`, `RENAME`, `ZUNIONSTORE`, `BITOP`, and Lua with several `KEYS`.

In Cluster, all keys in one command must hash to the same slot. If they do not, Redis returns an error such as `CROSSSLOT Keys in request don't hash to the same slot`.

```text
MGET {u:1}:name {u:1}:city
MGET u:1:name u:2:name
```

The first `MGET` can work (same tag). The second fails if the slots differ.

`RENAME` requires the source and destination in the same slot. To rename across slots, the application `GET`s, `SET`s the new key, and `DEL`s the old key. That path is not atomic.

Transactions (`MULTI`/`EXEC`) and `WATCH` also require all keys on the same node and slot rules. You cannot span two primaries in one Redis transaction.

Pipelines can include commands for different slots only if the client splits the pipeline by node. A naive pipeline to one node fails or redirects.

Single-key commands (`GET`, `HGET`, `INCR`) always target one slot. They are the easy case. The client follows `MOVED` to the right node.

When you design a schema, put keys that you must read together under one hash tag, or avoid multi-key commands and do several single-key calls.

`SCAN` is per node. A full cluster scan must query every primary.

### Questions

#### Theoretical questions

1. What error name does Redis use when keys are in different slots?
2. Why can `MGET` of two user ids fail in Cluster?
3. Why is `RENAME` across slots not a single atomic command?
4. Can `MULTI`/`EXEC` include keys from two primaries?
5. How do you `SCAN` the whole cluster?

#### Easy practical tasks

1. List six multi-key commands from this section.
2. Write four sentences: same-slot vs cross-slot.
3. Make a table: command, safe with hash tag?, unsafe example keys.
4. Open the `MSET` command page. Write any Cluster note that the page shows.

#### Medium practical tasks

1. On a cluster, reproduce `CROSSSLOT` with `MGET` of two untagged keys. Save the error.
2. Fix the same read with a hash tag or with two `GET`s. Show both results.
3. Try `RENAME` to a name that hashes to another slot. Record the error. Do the GET/SET/DEL path.

#### Advanced practical tasks

1. Write a client helper that groups keys by `CLUSTER KEYSLOT` (or client slot API) and runs one `MGET` per slot. Test with 20 keys.
2. Read which commands are explicitly not available in Cluster (docs). Write five names and one reason (example: `SELECT`).

---

## Client support for cluster mode

A cluster-aware client:

- Connects to at least one node
- Loads the slot map (`CLUSTER SLOTS`)
- Sends each command to the node that owns the slot
- Follows `MOVED` and updates the map
- Follows `ASK` with `ASKING` when a slot is migrating
- Splits pipelines and multi-key commands by slot when the library can

A cluster-unaware client talks to one node like standalone Redis. That node returns `MOVED` for keys it does not own. The unaware client treats `MOVED` as an error. Do not use a standalone-only client against Cluster.

`redis-cli -c` enables cluster mode. The CLI then follows redirects.

```text
redis-cli -c -h 127.0.0.1 -p 7000
```

Connection pools must store the slot map and a connection per node (or a subset). After failover or reshard, the map changes. Good clients refresh on `MOVED`.

Timeouts and retries need care. A retry of a write after a network error can double-apply `INCR`. Use idempotent writes or accept the risk.

Libraries: many official and popular clients have a cluster option (cluster client, or a flag). Read that page. Do not only set a host list for Sentinel when the server is Cluster. Sentinel and Cluster are different topologies.

`CLUSTER SLOTS` output includes node addresses and slot ranges. Firewalls must allow the client to reach every primary (and replicas if you read from replicas). One VIP that hides all nodes does not work unless it is a smart proxy.

Test the client in a lab with `MOVED` (restart a node, reshard one slot) before production.

`READONLY` is a Cluster command for replica reads. Replica reads can be stale, the same as Sentinel replica reads.

### Questions

#### Theoretical questions

1. What does `redis-cli -c` enable?
2. What does a client do when it receives `MOVED`?
3. Why does a standalone client fail against Cluster?
4. Why must the client reach every primary address?
5. What extra risk do retries add for `INCR`?

#### Easy practical tasks

1. Run `redis-cli` without `-c` against a cluster key that is not local (if you have a cluster). Save the `MOVED` text. Then retry with `-c`.
2. Write four sentences: cluster-aware vs unaware client.
3. Open the cluster tutorial for your language. Write the class or option name.
4. Draw a client, three primaries, and a `MOVED` arrow.

#### Medium practical tasks

1. Use a cluster client to `SET` 100 keys without tags. Confirm `CLUSTER COUNTKEYSINSLOT` is nonzero on more than one slot.
2. Disconnect one node (lab). Record client errors and whether the client recovered when the node returned.
3. Compare Sentinel client config and Cluster client config in a two-column table.

#### Advanced practical tasks

1. Enable replica reads in a cluster client if the library supports it (`READONLY`). Write the stale-read warning again in Cluster terms.
2. Implement or sketch slot calculation in your language (CRC16 mod 16384 plus hash tags). Verify against `CLUSTER KEYSLOT` for 20 keys.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do `REPLICAOF`, async lag, and promotion combine when the primary dies?
2. What does Sentinel add that a manual `REPLICAOF NO ONE` does not add?
3. How do slots, hash tags, and `CROSSSLOT` errors shape key design?
4. When do you choose Cluster instead of Sentinel plus one primary?
5. Why does a multi-key Lua script need a hash tag in Cluster?

#### Easy practical tasks

1. Write a cheat sheet: `REPLICAOF`, `WAIT`, `SENTINEL get-master-addr-by-name`, 16384, `{tag}`, `CROSSSLOT`, `-c`, `MOVED`.
2. Draw the happy path (one primary, two replicas) and the path after promotion.
3. Draw a schema for user profile + cart with tags, and a global inventory key without that user tag.
4. Write five lab rules: no writes on replicas, three Sentinels, `-c`, no `SELECT`, do not kill a primary that still owns slots.

#### Medium practical tasks

1. Document a Docker topology (compose service names, ports 6379 and 26379) so that another beginner can start Sentinel.
2. Write a one-page key-design guide for Cluster for a shop app (user, cart, sku, leaderboard).
3. Document client settings: Sentinel vs Cluster, seed nodes, timeouts, retries, and whether the library splits `MGET`.

#### Advanced practical tasks

1. Kill the primary, wait for Sentinel, write, then start the old primary. Prove there is one primary and that the old node is a replica.
2. Compare Cluster failover of a primary with Sentinel failover in a table: sharding, `SELECT`, multi-key, client type, minimum nodes.
