# 11. Clustering

## Description

Redis Cluster splits keys across many primary instances. Each primary owns a subset of hash slots. This topic covers the 16384 slots, hash tags, multi-key command limits, resharding at a high level, and what a client must do in cluster mode.

Complete the earlier topics first. Cluster is not Sentinel. Cluster shards data. Sentinel fails over one primary.

Use one term for each concept. A slot is an integer from `0` to `16383`. A shard is a primary (and its replicas) that owns some slots. A hash tag is the `{...}` part of a key that Redis uses to compute the slot. A redirect is a `MOVED` or `ASK` reply that tells the client another node.

---

## Redis Cluster hash slots (16384)

Redis Cluster divides the keyspace into 16384 slots. Every key maps to one slot. Redis computes `CRC16(key) mod 16384` (with hash-tag rules). The result is the slot.

Each primary owns a contiguous or non-contiguous set of slots. Together the primaries own all 16384 slots. If a slot has no owner, the cluster is not healthy for that slot.

A write of a key goes to the primary that owns the key's slot. Other primaries do not store that key. This is sharding. Memory and CPU scale out. A single key still lives on one primary. A huge key does not split.

```text
CLUSTER KEYSLOT user:42
CLUSTER SLOTS
CLUSTER NODES
```

`CLUSTER KEYSLOT` works on a cluster node and shows the slot for a key. Use it in the lab to prove that two keys share a slot or not.

The cluster bus (a second port, often `16379` when the client port is `6379`) carries gossip and fail-over messages between nodes. Clients use the client port. You do not send `GET` to the bus port.

Logical DBs do not exist in Cluster. `SELECT` is not available. Use key prefixes.

A cluster needs several primaries (three is a common minimum for production). Each primary can have replicas for failover inside the cluster. That failover is not Sentinel. The cluster nodes vote.

Do not run Cluster mode for a tiny cache that fits on one instance unless you need the operational practice. Cluster adds redirects, slot limits, and more nodes.

`INFO cluster` and `CLUSTER INFO` show `cluster_state`. `ok` means slots are covered. `fail` means the cluster is not ready.

### Questions

#### Theoretical questions

1. How many hash slots does Redis Cluster use?
2. How does Redis map a key to a slot (high level)?
3. Does one key split across two primaries?
4. What command shows the slot of a key?
5. Why is `SELECT` not a Cluster feature?

#### Easy practical tasks

1. Write four sentences: Cluster vs one Redis instance.
2. Open the Cluster spec or tutorial. Write the number 16384 and one reason the page gives (if any).
3. Draw three primaries and boxes of slots (for example 0–5460, 5461–10922, 10923–16383).
4. List two `CLUSTER` commands from this section and one purpose each.

#### Medium practical tasks

1. Start a cluster (Docker or `redis-cli --cluster create`) or use a hosted cluster. Run `CLUSTER INFO`. Write `cluster_state` and `cluster_slots_assigned`.
2. Run `CLUSTER KEYSLOT` on five key names. Write the five slot numbers.
3. Run `CLUSTER NODES`. Identify one primary and one replica line (role words in the flags).

#### Advanced practical tasks

1. Read how CRC16 and slots work in the official spec. Write the formula and the hash-tag exception in six sentences.
2. Compare a 3-primary cluster with one large instance for a 6 GB dataset. Write memory per node and the huge-key risk.

---

## Key hash tags `{user:42}`

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

```text
user:{42}:profile
```

Here the tag is `42`. All keys with `{42}` share a slot. That can put many users' keys in one slot if you tag only `{42}` style ids that collide. Prefer a tag that includes a unique user id: `{user:42}`.

Hash tags are a deliberate hot-spot tool. If all keys use `{global}`, one slot and one primary take all traffic. That defeats sharding.

Use hash tags when a transaction, a Lua script, or `MGET` must touch several keys. Do not tag everything.

`CLUSTER KEYSLOT` proves the rule. Two names with the same tag must show the same slot.

Cluster-safe Lua must access only keys in `KEYS[]`, and those keys must share a slot (same tag in practice).

### Questions

#### Theoretical questions

1. Which part of `{user:42}:cart` does Redis hash?
2. What happens when the braces are empty `{}`?
3. Why can `{user:42}:a` and `{user:42}:b` work in one `MGET`?
4. Why is `{global}` a bad tag for all keys?
5. What must a Cluster-safe Lua script do with key names?

#### Easy practical tasks

1. Write five key names for one user that share a tag.
2. Write four sentences that explain hash tags.
3. Compute (on paper) whether `a{x}b` and `c{x}d` share a slot (yes). Write why.
4. Open the hash tag docs. Copy the official rule in your own short words.

#### Medium practical tasks

1. On a cluster, run `CLUSTER KEYSLOT` for `{user:1}:a`, `{user:1}:b`, and `{user:2}:a`. Write which pair matches.
2. `MGET` two same-tag keys that you `SET` on the cluster. Then `MGET` two keys with different tags. Record the error or redirect behavior.
3. Show an empty-tag key `foo{}bar` vs `foobar` with `CLUSTER KEYSLOT`.

#### Advanced practical tasks

1. Design tags for a cart (user hash, cart list, inventory). Write which keys share a tag and which must not (inventory is global).
2. Measure skew: 10,000 keys with no tag vs 10,000 keys with one tag. Use `CLUSTER COUNTKEYSINSLOT` on a few slots. Write the imbalance.

---

## Multi-key commands that fail across slots

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

## Resharding (high-level)

Resharding moves slots from one primary to another. You do this to add a primary, to remove a primary, or to balance memory and CPU.

The operator (or a managed service) runs a reshard. `redis-cli --cluster reshard` is the common tool for a self-managed cluster. Cloud consoles hide the same idea.

At a high level:

1. Pick source primary, destination primary, and a list of slots
2. For each slot, migrate keys in batches (`MIGRATE` / cluster setslot states)
3. Update slot ownership so clients `MOVED` to the new primary
4. Repeat until the plan is done

During migration a slot can be in a transitional state. Clients can receive `ASK` redirects. `ASK` means "try that node once for this key, do not update your slot map forever". `MOVED` means "update your slot map; that node owns the slot".

Resharding uses extra CPU, network, and memory. Do it in a window. Watch latency and `CLUSTER INFO`.

Resharding does not split a single large key. If one key is too big, you must change the data model.

Adding a replica is not resharding. A replica copies one primary. It does not take slots until you promote it or you rebalance.

Removing a primary requires moving its slots away first. Do not kill a primary that still owns slots.

Managed Redis Cluster often reshard for you when you change the shard count. Read the vendor runbook.

### Questions

#### Theoretical questions

1. What does resharding move between primaries?
2. Why do you reshard when you add a primary?
3. What is the difference between `ASK` and `MOVED`?
4. Can resharding split one huge key?
5. What must you do before you remove a primary that owns slots?

#### Easy practical tasks

1. Write four sentences that describe resharding to a beginner.
2. Draw slots moving from node A to node B (before and after).
3. Open the Cluster admin tutorial. Write the `redis-cli --cluster` subcommand for reshard.
4. Write two reasons not to reshard during a peak sale.

#### Medium practical tasks

1. On a lab cluster with at least two primaries, reshard a few slots (or follow a managed UI). Run `CLUSTER KEYSLOT` and `CLUSTER NODES` before and after for one key.
2. Read `ASKING` in the docs. Write when a client sends `ASKING`.
3. List `CLUSTER SETSLOT` states at a high level (importing, migrating, stable) from the docs.

#### Advanced practical tasks

1. Time a small reshard of a lab slot that contains 10,000 keys. Write duration and whether `PING` latency changed.
2. Write a reshard checklist: backup, replica health, slot coverage, client versions, rollback (move slots back).

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

1. How do slots, hash tags, and `CROSSSLOT` errors shape key design?
2. When do you choose Cluster instead of Sentinel plus one primary?
3. How do `MOVED` and `ASK` differ in the client slot map?
4. Why does a multi-key Lua script need a hash tag in Cluster?
5. What operational work does resharding add that a single instance never needs?

#### Easy practical tasks

1. Write a cheat sheet: 16384, `CLUSTER KEYSLOT`, `{tag}`, `CROSSSLOT`, `-c`, `MOVED`, `ASK`, `CLUSTER INFO`.
2. Draw a schema for user profile + cart with tags, and a global inventory key without that user tag.
3. Write five lab rules: `-c`, no `SELECT`, tag multi-key, cluster client, do not kill a primary that still owns slots.
4. Export `CLUSTER NODES` to a file. Mark one primary, one replica, and a slot range.

#### Medium practical tasks

1. Build a three-primary cluster in Docker (or Redis Stack / official create). `SET` keys, fail one replica (not the only copy of slots), and document `CLUSTER INFO`.
2. Write a one-page key-design guide for Cluster for a shop app (user, cart, sku, leaderboard).
3. Document client settings: seed nodes, timeouts, retries, and whether the library splits `MGET`.

#### Advanced practical tasks

1. Reshard 100 slots in a lab, run a write load during the move, and record `ASK`/`MOVED` counts if the client logs them (or redis-cli `-c` behavior).
2. Compare Cluster failover of a primary with Sentinel failover in a table: sharding, `SELECT`, multi-key, client type, minimum nodes.
