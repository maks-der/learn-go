# 20. Advanced and Internals

## Description

This topic is awareness of how Redis runs inside one process and how alternative servers relate to it. You learn the event loop, I/O threads, active expiry, the Cluster bus, Functions versus Lua, and names such as Valkey and KeyDB.

Complete topics 2, 8, 11, and 16 first. You do not need to read the Redis C source to finish this topic. Official docs and `INFO` fields are enough.

Use one term for each concept. The event loop is the wait-and-dispatch cycle that handles sockets and time events. I/O threads read and write sockets; they do not run command logic. The active expire cycle is a background job that deletes expired keys. The Cluster bus is the node-to-node channel. A Function is a named library on the server. An alternative implementation is a different server that speaks Redis protocols or commands.

---

## Event loop (ae) awareness

Redis uses an event loop (historically the `ae` library). The process waits for a socket that is ready, then reads a command, then runs it on the main thread, then writes the reply when the socket is ready.

Time events sit in the same loop: cron-like work such as expiry sampling, replication housekeeping, and statistics.

One command runs to completion. A slow command delays the next socket and the next time event. That is why `KEYS` and a huge `HGETALL` stall clients (topic 2, topic 16).

The loop is not a promise that Redis uses only one core for all work. Forked children write RDB files. I/O threads can move socket read/write. Modules can add work. Command execution for the default engine remains one main thread.

You do not configure `ae` by name in `redis.conf`. You see the effect in latency, the slow log, and `INFO`.

Blocked commands (`BRPOP`, `XREAD BLOCK`) park that client. The loop continues to serve other clients. `blocked_clients` (topic 19) counts those waiters.

Do not design an application command that needs 100 ms of CPU on the main thread. Split the work.

### Questions

#### Theoretical questions

1. What does the event loop wait for?
2. What happens to the loop during a slow command?
3. What are time events for?
4. Do blocked clients stop the whole loop?
5. Why do you still call command execution single-threaded?

#### Easy practical tasks

1. Write four sentences: event loop vs "one OS thread for everything including RDB".
2. Draw socket ready → command → reply.
3. List three time-event jobs from this section.
4. Open a Redis internals page or source overview. Write the name `ae` and one sentence.

#### Medium practical tasks

1. Run a blocking `BLPOP` and many `PING`s from another client. Confirm `PING` still works. Write the observation.
2. Generate a slow command in a lab. Watch other `PING` latency. Write both times.
3. Read `INFO` commandstats if available (`INFO commandstats`). Write the two slowest commands.

#### Advanced practical tasks

1. Read a short description of `beforeSleep` / cron in Redis docs or a reliable article. Write six sentences on what runs between commands.
2. Write a latency budget: max command time you allow, and how you enforce it (slow log, deny `KEYS`).

---

## I/O threading

Since Redis 6, you can enable I/O threads. Those threads read incoming bytes and write outgoing bytes. The main thread still parses and executes commands.

```text
CONFIG GET io-threads
CONFIG GET io-threads-do-reads
```

Defaults are often off or `1` (no extra I/O threads). A small instance on a local socket may not gain. A high-connection, small-command workload on a multi-core host can gain throughput.

I/O threads do not make `INCR` run in parallel on the same key. They do not remove the big-key problem.

Tuning:

- Set `io-threads` to a small number of cores, not to "all cores"
- Measure `INFO` and application p99
- Read the version notes; behavior changed across releases

Do not enable I/O threads and then assume Cluster is unnecessary. Scale of command CPU still needs shards when one core is the limit.

Threads for I/O are not KeyDB-style multi-threaded command execution (later section).

Leave I/O threads off until you have a measurement that shows a win.

### Questions

#### Theoretical questions

1. What do I/O threads do?
2. What do I/O threads not do?
3. Which Redis version introduced this feature?
4. Why can a local lab see no gain?
5. Why do I/O threads not replace Cluster?

#### Easy practical tasks

1. `CONFIG GET io-threads`. Save the value.
2. Write four sentences: I/O threads vs main-thread commands.
3. Open Redis 6 release notes or I/O thread docs. Write one official sentence in your own words.
4. List a workload that might benefit (many small `GET`s, many connections).

#### Medium practical tasks

1. On a lab multi-core host, compare `PING` or `GET` throughput with `io-threads` 1 vs 4 (lab only). Write both numbers and a caveat.
2. Read `INFO` fields related to I/O threads if present. Write the names.
3. Write a policy: when your team allows `io-threads` and who measures.

#### Advanced practical tasks

1. Profile or document CPU: main thread busy vs extra cores idle before and after I/O threads. Write what you can conclude.
2. Compare I/O threads with "start a second Redis instance" for the same host. Write contention and keyspace split.

---

## Active expire cycle

Keys with a TTL must disappear. Redis does not only wait for the next `GET`.

Two mechanisms:

- Lazy expire: when a command touches the key, Redis deletes it if the TTL is in the past
- Active expire: a cycle samples keys that have TTLs and deletes expired ones

The active cycle runs as a time event in the loop. It uses a fraction of time so that expiry does not monopolize the main thread. If many keys expire at once (topic 13 stampede, topic 4), the cycle can take more work. `expired_keys` rises.

`INFO` can show expire-related fields (`expired_stale_perc` and others by version). A high stale percent means many keys in the expire set are already dead but not yet deleted. Memory stays high until the cycle catches up.

A huge expire storm after a bulk load with the same TTL is an operations event. Jitter (topic 13) spreads that storm.

`EXPIRE` of millions of keys in one second also creates future work.

Do not confuse `expired_keys` with `evicted_keys`. Eviction is `maxmemory` (topic 4).

### Questions

#### Theoretical questions

1. What is lazy expire?
2. What is the active expire cycle?
3. Why can memory stay high after keys "should" have expired?
4. How does TTL jitter help the expire cycle?
5. How do expired and evicted counters differ?

#### Easy practical tasks

1. `SET lab:e 1 EX 2`. Wait 3 seconds. `GET lab:e`. Save the reply.
2. Run `INFO stats`. Write `expired_keys`.
3. Write four sentences: lazy vs active expire.
4. Open the Redis expire docs. Write one sentence about the sampling cycle.

#### Medium practical tasks

1. Create 10,000 keys with `EX 5` in a pipeline. Watch `DBSIZE` and `used_memory` for 15 seconds. Write when the count drops.
2. Repeat with jitter of 0–10 seconds. Compare the drop curve.
3. Find expire fields in `INFO`. Write two field names and values.

#### Advanced practical tasks

1. Read the official expired-key deletion page. Write the high-level algorithm in six sentences (no need for C).
2. Design a bulk load that avoids an expire storm. Write TTL policy and batch size.

---

## Cluster bus

Topic 11 mentioned the Cluster bus. Each Cluster node listens on a second port for node-to-node messages. The default is the client port plus `10000` (client `6379` → bus `16379`).

The bus carries:

- Gossip (which nodes exist, which slots they own, flags)
- Failover votes
- Pub/Sub fan-out across nodes
- Migration messages related to slots

Clients do not send `GET` to the bus port. Firewalls must allow bus traffic between nodes, not from the public internet.

```text
CLUSTER INFO
CONFIG GET cluster-port
```

If the bus is blocked, nodes mark each other as unreachable. The cluster can go `fail` even if clients can reach every client port.

The bus protocol is not RESP application commands. Do not point `redis-cli` at the bus port and expect a normal prompt to be useful.

Security: protect the bus like the client port (private network, no public bind). A stranger on the bus is a cluster-member problem.

Do not publish the bus port on a LoadBalancer.

### Questions

#### Theoretical questions

1. What is the default bus port if the client port is `6379`?
2. Name three message kinds on the bus.
3. Why must nodes reach each other on the bus?
4. Do clients use the bus port for `GET`?
5. Why must the bus stay off the public internet?

#### Easy practical tasks

1. Write the default client port and bus port pair.
2. Write four sentences: client port vs bus port.
3. Open Cluster spec or tutorial. Write the +10000 rule.
4. Draw three nodes, two ports each, and client arrows only to client ports.

#### Medium practical tasks

1. On a lab cluster, `CLUSTER INFO` and `CLUSTER NODES`. Write one bus-related field if present.
2. Read `cluster-port` / `cluster-announce-bus-port` docs. Write when you set them (NAT, Docker).
3. Write a firewall table: app → client ports, node → bus ports, internet → none.

#### Advanced practical tasks

1. Simulate a blocked bus in a lab (or read a post-mortem). Write symptoms (`cluster_state`, node flags) and the fix.
2. Write a Cluster network checklist: client, bus, TLS if used, announce IPs.

---

## Redis Functions vs Lua

Topic 8 covered `EVAL` / `EVALSHA` and Redis 7 Functions (`FUNCTION LOAD`, `FCALL`).

Lua scripts:

- Travel from the client (or from `SCRIPT LOAD`)
- Use SHA1 cache that empties on restart unless you reload
- Are easy to copy-paste into every service
- Create version drift (three SHA1s of "the same" lock script)

Functions:

- Live as named libraries on the server
- Update in one place
- Persist with RDB/AOF when configured (test it)
- Still run on the main thread and still block if slow
- Still must list keys for Cluster safety

Both are atomic with respect to other commands. Both must be deterministic enough for replication. Both are not a general application server.

Choose Functions for Redis 7+ when many clients must share one server-side version. Choose `EVALSHA` plus a deploy-time `SCRIPT LOAD` when you are on older Redis or when the team already has that pipeline.

Do not maintain both a Function and a Lua copy of the same lock without one source of truth.

`FUNCTION KILL` and `SCRIPT KILL` exist for a runaway script that the server allows you to kill (read the rules; some scripts are unkillable while they run certain commands).

### Questions

#### Theoretical questions

1. Where does Function code live after `FUNCTION LOAD`?
2. What operational problem do Functions solve versus `EVAL` in every client?
3. Do Functions run in parallel with other commands?
4. Why must you test Functions across restart?
5. When do you keep `EVALSHA`?

#### Easy practical tasks

1. Write `redis_version`. State whether Functions exist (7+).
2. Write four sentences: `FCALL` vs `EVALSHA`.
3. Run `FUNCTION LIST`. Save the reply.
4. Open the Functions page. Write `FUNCTION LOAD` and `FCALL` in one sentence each.

#### Medium practical tasks

1. If Redis 7+, load a tiny official example and `FCALL` it. If not, write the image tag you need.
2. `SCRIPT LOAD` a one-liner. `EVALSHA`. `SCRIPT FLUSH`. `EVALSHA` again. Record `NOSCRIPT`.
3. Write a team policy: one store for server-side code (git path + load job).

#### Advanced practical tasks

1. Port a lock-release script to a Function. Test match and mismatch. Restart with persistence. Confirm the library.
2. Read kill rules for scripts vs functions. Write when a kill works and when you must wait.

---

## Alternative implementations (Valkey, KeyDB) — awareness

The Redis protocol and command set inspired other servers.

Valkey is a Linux Foundation project that forked from Redis after a license change in the Redis project. Many commands stay familiar. It is an alternative you may see in vendors and distros. Read Valkey docs for Cluster, modules, and license. Do not assume every Redis module loads.

KeyDB is a fork that emphasizes multi-threaded command execution. The concurrency model is not the Redis main-thread model. Race assumptions that are true on Redis (`INCR` atomic on one thread) need a careful read of KeyDB docs. Modules and Cluster behavior can differ.

Other proxies and compatible caches exist (vendor caches, `redis-stack` images, cloud engines). Compatibility is a spectrum: common `GET`/`SET` work; `MODULE`, `FUNCTION`, and exact `INFO` fields may not.

Awareness rules:

- Speak the server name in runbooks ("Valkey 8" not "Redis" if it is Valkey)
- Test the client against that server
- Read the license for your use
- Do not mix replicas of different engines unless a vendor supports it

This handbook remains a Redis learning path. Commands here target Redis. When you use an alternative, verify each command on that engine.

Do not treat a fork as Redis for compliance or for support contracts without a check.

### Questions

#### Theoretical questions

1. What is Valkey in one sentence?
2. What design emphasis is KeyDB known for?
3. Why might a Redis module fail on a fork?
4. Why do you name the engine in a runbook?
5. Why is this handbook still written for Redis?

#### Easy practical tasks

1. Open the Valkey site. Write the project one-line description in your own words.
2. Open a KeyDB page. Write one sentence about threads.
3. Write four sentences: protocol compatible vs fully compatible.
4. Write the license names you see today for Redis and Valkey (date your note).

#### Medium practical tasks

1. If you can run a Valkey container, `PING` and `INFO server`. Write the server name field. If not, write a `docker` image name from the docs.
2. Compare one feature (Functions, ACL, Cluster) in Redis docs vs Valkey docs. Write match or differ.
3. Write a client test list: 10 commands you would run before you switch an app.

#### Advanced practical tasks

1. Write a migration note: Redis to Valkey (or the reverse) for a cache instance. Include license, image, modules, and rollback.
2. Read a reliable comparison of KeyDB threading vs Redis I/O threads. Write six sentences on what stays atomic and what you must re-read.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do the event loop, I/O threads, and the expire cycle share one process without becoming a multi-core command engine?
2. What breaks if the Cluster bus is isolated but client ports are open?
3. When do Functions, Lua, and an alternative server each change your deploy story?
4. How does an expire storm connect topic 13 jitter to this topic's cycle?
5. What do you tell a teammate who says "KeyDB is Redis with more cores, so our Lua races go away"?

#### Easy practical tasks

1. Write a cheat sheet: `ae`, `io-threads`, active expire, bus +10000, `FCALL`, Valkey, KeyDB.
2. Draw one process: main thread, optional I/O threads, child fork for RDB.
3. Export `INFO server` and `INFO cluster` (or write "no cluster"). Highlight version and cluster state.
4. Write five awareness rules for forks and modules.

#### Medium practical tasks

1. On one lab instance, collect: `io-threads`, expire fields, `FUNCTION LIST`, `redis_version`. Write a one-page internals snapshot.
2. Write a table: feature, Redis behavior, what to verify on Valkey, what to verify on KeyDB.
3. Document a "slow command" incident using loop + expire + big key language from this topic and topic 16.

#### Advanced practical tasks

1. Write an internals primer for new ops: 2 pages, no C, with `INFO` commands to show each idea.
2. Compare Redis Cluster bus with a vendor-managed cluster. Write which internals you no longer configure and which you still need (slots, tags).
