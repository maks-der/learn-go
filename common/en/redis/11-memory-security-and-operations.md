# 11. Memory, Security, and Operations

## Description

Redis is fast when data fits in RAM and commands stay small. Redis does not start as a public internet service. Operations keep Redis correct after the application works. This topic shows how you read memory, how you find large keys, how the slow log works, how you harden access, and how you handle `BGSAVE`, disk fill, and upgrades.

Complete topic 2 (`maxmemory`), topic 7 (persistence), and topic 8 (replication) first. Use a lab instance for every practical task that changes auth or command names. Do not practice on a shared production instance.

Use one term for each concept. Used memory is the bytes that Redis reports for data and overhead. RSS is the process size that the operating system sees. A big key is one value that is large in bytes or in element count. The slow log is a list of commands that exceeded a time threshold. An ACL user is a named account with a command list and a key pattern. TLS encrypts the connection. Copy-on-write (COW) is the extra RAM during a fork.

This topic is hardening and operations. It does not describe attacks.

---

## `INFO memory`, `MEMORY USAGE`, big keys

`INFO memory` prints memory fields as `key:value` lines. Useful fields (names can vary slightly by version):

- `used_memory` / `used_memory_human` — allocator view of Redis data
- `used_memory_rss` / `used_memory_rss_human` — operating system resident size
- `used_memory_peak_human` — peak since start
- `maxmemory` — configured limit (`0` often means no limit on 64-bit)
- `mem_fragmentation_ratio` — RSS divided by used memory (approximate)

```text
INFO memory
```

`used_memory` and RSS are not the same. Fragmentation, file mappings, and copy-on-write during `BGSAVE` raise RSS.

A ratio much above `1` can mean fragmentation or a recent fork. A ratio below `1` can appear with shared pages or allocator quirks. Do not treat one number as a verdict. Use `MEMORY DOCTOR` and the docs for your version.

`INFO stats` adds `expired_keys`, `evicted_keys`, `keyspace_hits`, and `keyspace_misses`. Watch `used_memory` against `maxmemory`. When they meet, eviction or `OOM` errors begin (topic 2).

`MEMORY USAGE key` estimates the bytes of one key, including overhead. It is an estimate. It is enough to find outliers.

```text
MEMORY USAGE lab:big
```

The reply is an integer (bytes) or null if the key is missing.

`MEMORY DOCTOR` returns text advice about the instance. Treat it as a hint. `MEMORY STATS` returns a map of allocator and dataset numbers.

`MEMORY PURGE` asks the allocator to release free pages to the operating system when the allocator supports it. Effect varies. Use it in a lab first.

Do not run `MEMORY USAGE` on every key in production in a tight loop without `SCAN` pacing. Each call is a command on the main thread.

A big key is a key with a large byte size or a large number of elements. Examples: a 50 MB string, a list of 10 million items, a hash with 1 million fields.

Risks:

- `GET`, `HGETALL`, `LRANGE 0 -1`, `SMEMBERS` copy or walk the whole value on the main thread
- Replication and RDB grow
- Network buffers grow
- Cluster cannot split one key across slots (topic 8)

Find big keys:

```text
redis-cli --bigkeys
redis-cli --memkeys
```

These tools use `SCAN`. They sample. They can miss a key. They still add load. Run them on a replica when you can (topic 8). Accept stale reads.

Fixes:

- Split one huge collection into many keys (shards, time buckets)
- Store blobs in object storage; keep a pointer in Redis
- Use `HSCAN`, `LRANGE` pages, `SSCAN` instead of full dumps
- `UNLINK` instead of `DEL` for large values (lazy free)

Do not ignore a big key because `used_memory` still has headroom. The next `HGETALL` can stall the instance.

`--bigkeys` is a lab and a replica tool. Do not schedule it on the primary every minute.

`OBJECT ENCODING key` shows the current in-memory layout (`listpack`, `hashtable`, `int`, `embstr`, and others). Redis converts small collections to larger structures when they grow. Benchmarks on 10 fields do not predict 100,000 fields.

### Questions

#### Theoretical questions

1. What is the difference between `used_memory` and `used_memory_rss`?
2. What does `MEMORY USAGE` return?
3. What two senses of "big" apply to a key?
4. Why is `HGETALL` on a huge hash dangerous?
5. Why run a big-key scan on a replica when you can?

#### Easy practical tasks

1. Run `INFO memory`. Write `used_memory_human` and `used_memory_rss_human`.
2. `SET lab:m hello`. `MEMORY USAGE lab:m`. Save the number.
3. Run `MEMORY DOCTOR`. Write one sentence from the reply (your own words).
4. Run `redis-cli --bigkeys` on a lab instance. Write the summary lines.

#### Medium practical tasks

1. `SET` 10,000 keys. Record `used_memory` before and after. Delete the keys. Record again.
2. Compare `MEMORY USAGE` of a 1-field hash and a 1000-field hash. Write both numbers.
3. `HSET` 50,000 fields. Time `HGETALL` vs an `HSCAN` loop. Write both times. `UNLINK` the key.

#### Advanced practical tasks

1. `SCAN` 200 keys and print `MEMORY USAGE` for each. List the top five. Delete only lab keys.
2. Write a big-key policy: max string bytes, max hash fields, how you alert, and the split pattern.

---

## Slow log

The slow log records commands that exceed a threshold. It lives in memory on the instance.

```text
SLOWLOG GET 10
SLOWLOG LEN
SLOWLOG RESET
```

Config:

- `slowlog-log-slower-than` — microseconds (1000 microseconds = 1 millisecond). A negative value can disable logging (check your version).
- `slowlog-max-len` — maximum entries

```text
CONFIG GET slowlog-log-slower-than
CONFIG GET slowlog-max-len
```

Each entry includes an id, a timestamp, an execution time, and the command arguments (truncated). Use it to find `KEYS`, big `HGETALL`, heavy `SORT`, or a slow Lua script.

The time is command execution on the server. It does not include client RTT.

`SLOWLOG RESET` clears the log. Use it after you copy the entries. Do not reset on a shared production instance without a reason.

A command can be slow because the key is big, because the script loops, or because the host is overloaded. Pair the slow log with `MEMORY USAGE` and `INFO`.

Do not set the threshold to `0` on a busy production instance for a long time. The log fills with normal commands.

Latency monitor (`LATENCY DOCTOR`, `LATENCY LATEST`) is a related tool for spikes. Use it when the slow log is not enough.

### Questions

#### Theoretical questions

1. What does the slow log record?
2. What unit is `slowlog-log-slower-than`?
3. Does the slow log include network RTT?
4. What does `SLOWLOG RESET` do?
5. Why is threshold `0` risky on a busy instance?

#### Easy practical tasks

1. `CONFIG GET slowlog-log-slower-than`. Save the value.
2. `SLOWLOG GET 5`. Write how many entries you see (zero is fine).
3. Write four sentences: slow log vs `INFO memory`.
4. Open the `SLOWLOG` command page. Write the reply fields.

#### Medium practical tasks

1. On a lab instance, set `slowlog-log-slower-than` to a small number. Run a large `KEYS` or a big `HGETALL`. `SLOWLOG GET 3`. Save one entry. Restore the threshold.
2. `SLOWLOG LEN` before and after. Write the two numbers.
3. Enable a short Lua loop (lab). Confirm it appears in the slow log.

#### Advanced practical tasks

1. Collect 20 slow log entries (lab load). Group them by command name. Write the top two causes.
2. Write an ops playbook: threshold, max len, how often you scrape, when you use `LATENCY DOCTOR`.

---

## ACLs, bind address, TLS

Older Redis used one password for the whole instance:

```text
CONFIG SET requirepass labpass
AUTH labpass
```

`requirepass` is a single shared secret. Every client that knows it has full power. Redis 6 adds Access Control Lists (ACLs). Prefer ACLs on Redis 6 and newer.

ACL ideas:

- A default user exists. You can set a password on it or disable it.
- `ACL SETUSER` creates or changes a user.
- A user has rules: commands (`+GET`, `-FLUSHALL`), categories (`+@read`), and key patterns (`~cache:*`).
- `AUTH username password` selects the user.

```text
ACL LIST
ACL WHOAMI
ACL SETUSER app on >apppass ~cache:* +@read +SET +DEL -@dangerous
```

Syntax details are on the `ACL` command pages. Check your version. The `>` prefix in `SETUSER` sets a password. `on` enables the user.

Application clients use a least-privilege user. A human operator uses a stronger user for `CONFIG` and `ACL` work. Do not give the application `FLUSHALL`, `DEBUG`, or `KEYS`.

Store passwords outside the repository. Use a secret manager or the environment. Do not put production passwords in shell history if you can use a prompt or a file with tight permissions.

`ACL SAVE` writes the ACL file when you use an ACL file. `CONFIG REWRITE` can persist `redis.conf` changes. A `CONFIG SET` password that you do not persist can disappear after restart. Test restart.

`bind` sets the addresses that Redis listens on.

```text
CONFIG GET bind
```

`bind 127.0.0.1` (and `::1` for IPv6 localhost) accepts only local clients. That is the correct default for a laptop lab.

`bind 0.0.0.0` listens on all IPv4 addresses. Any host that can reach the port can attempt a connection. Combine that with no auth and you have an open data store on the network.

Protected mode: when Redis has no bind configuration that you intended and no password, Redis may refuse remote commands. Protected mode is a safety net. It is not a substitute for a firewall and for ACLs.

In Docker, publishing `-p 6379:6379` can expose Redis on the host. Bind inside the container and publish only to `127.0.0.1:6379` on the host when you learn:

```text
docker run --name redis-learn -p 127.0.0.1:6379:6379 redis
```

A cloud security group or firewall must allow only the application subnets. Do not rely on "nobody knows the IP."

TLS encrypts bytes on the network. Without TLS, a packet capture on the path can show commands and values. Auth passwords also travel in clear text without TLS (RESP `AUTH`).

Redis can listen with TLS (`tls-port`, certificates, and related `tls-*` settings). Hosted Redis often requires TLS. Use `redis-cli --tls` and client library TLS options.

TLS does not replace ACLs. TLS protects the path. ACLs decide what the client may do after it connects.

Certificates expire. Calendar the expiry. Do not disable certificate verification in production to "make it work." Fix the CA and the host name.

Do not expose Redis to the public internet. Bots scan port `6379`. Bind to private addresses. Use a firewall. Use ACL users. Use TLS on any path that leaves a trusted host.

If you already exposed a lab, treat all keys as compromised. Rotate passwords. Take the instance off the public address.

### Questions

#### Theoretical questions

1. What is the limit of a single `requirepass`?
2. Which Redis version introduced ACLs?
3. What does `bind 127.0.0.1` restrict?
4. What does TLS protect on the Redis path?
5. Does TLS replace ACL users?

#### Easy practical tasks

1. Run `ACL LIST` on your instance. Write the user names.
2. `CONFIG GET bind`. Save the value.
3. `CONFIG GET protected-mode`. Save the value.
4. Write four sentences: TLS vs ACL.

#### Medium practical tasks

1. On a lab instance, create a user that can only `PING` and `GET` keys that start with `lab:`. Prove `SET` fails. Then delete the user.
2. Write a Docker publish flag that listens only on the host loopback.
3. Open the official Redis TLS page. Write three `tls-*` setting names.

#### Advanced practical tasks

1. Design three users: `app`, `readonly`, `ops`. Write command categories and key patterns for each.
2. Start Redis with TLS from the official example in a lab. Connect with `redis-cli --tls`. Record the flags. Use only lab certificates.

---

## Dangerous commands (`FLUSHALL`, `KEYS`, `DEBUG`)

Some commands are correct in a lab and harmful on a shared instance.

`FLUSHALL` deletes all keys on the instance (all logical DBs). `FLUSHDB` deletes the current logical DB. There is no undo.

`KEYS` walks the keyspace on the main thread. On a large instance it can stall Redis (topic 9).

`DEBUG` includes subcommands that can crash the process or block it (`DEBUG SEGFAULT`, `DEBUG SLEEP`). Those are developer tools. They are not application commands.

`CONFIG GET` can reveal settings. `CONFIG SET` can change bind, persistence, and directories. A client that can `CONFIG SET` can weaken security or break persistence.

Other commands in the same class: `SHUTDOWN`, `REPLICAOF` / `SLAVEOF`, `MODULE LOAD`, `SCRIPT DEBUG`. The `@dangerous` ACL category groups many of them. Check `ACL CAT dangerous` on your version.

The application user must not have these commands. An operator uses them in a controlled window with a change ticket.

Do not run `FLUSHALL` as a deploy step on a durable Redis. Topic 9 and topic 12 explain cache warming instead.

Older setups renamed dangerous commands in `redis.conf`:

```text
rename-command FLUSHALL ""
rename-command KEYS ""
```

An empty new name disables the command. Redis 7+ deprecates `rename-command` in favor of ACLs. Prefer ACLs. A secret rename is a weak control. The name can leak.

Do not rename `AUTH` or `ACL` in a way that locks you out of a remote instance. Keep a local console path that you tested.

Application code must use `SCAN`, not `KEYS`. If you disable `KEYS`, tutorials that still show `KEYS *` will fail. That is what you want in production.

### Questions

#### Theoretical questions

1. What does `FLUSHALL` delete?
2. Why is `KEYS` dangerous on a large instance?
3. Why is `DEBUG` not an application command?
4. What can `CONFIG SET` change that affects security?
5. What is `@dangerous` in ACLs?

#### Easy practical tasks

1. Run `ACL CAT dangerous` if the command exists. Write five command names.
2. Write four sentences: `FLUSHALL` vs `DEL` of one prefix with `SCAN`.
3. Open the `FLUSHALL` command page. Write the time complexity note.
4. Write one sentence each: `FLUSHALL`, `KEYS`, `DEBUG`, `CONFIG SET`.

#### Medium practical tasks

1. On an empty lab instance only, run `FLUSHDB` after you `SET` one key. Confirm `DBSIZE` is 0. Do not do this on shared data.
2. In your client, try to call `FLUSHALL` as a non-ops user (or write the ACL that would block it). Record the error.
3. In a disposable lab config, disable `FLUSHALL` with the method your version documents (ACL or rename). Prove `FLUSHALL` fails. Restore the lab.

#### Advanced practical tasks

1. Write an ops policy: which roles may run `CONFIG`, `FLUSHALL`, `DEBUG`, and `KEYS`. Include a ban on application use.
2. Build an ACL for `app` that includes `@read`, `@write`, and `@hash` (as needed) and excludes `@dangerous`, `KEYS`, and `FLUSHALL`. Test on lab keys.

---

## Fork/COW during `BGSAVE`, disk fill, upgrades

Topic 7 explained `BGSAVE`. The parent forks a child. The child writes the RDB. The operating system shares pages until the parent writes a page. Then COW copies that page. Many writes during a long save can raise RSS toward almost two copies of the dataset.

Effects:

- The host can swap or the Linux OOM killer can stop Redis
- Latency can rise during the fork on some kernels
- AOF rewrite (`BGREWRITEAOF`) also forks

Leave free RAM for the peak. A common planning rule is to keep headroom for a full extra copy if the write rate is high. Measure RSS during save. Do not guess only from `used_memory`.

Watch `used_memory_rss` while `rdb_bgsave_in_progress` is `1`.

Linux `overcommit_memory` notes appear in Redis warnings. A `fork` error in the log means the kernel refused the child. Fix RAM, overcommit, or dataset size.

`SAVE` does not fork. `SAVE` blocks the main thread. Do not use `SAVE` in production.

RDB and AOF write files under `dir`. If the volume is full, `BGSAVE` fails, AOF append fails, or rewrite fails. Redis can stop writes (`stop-writes-on-bgsave-error`, AOF error behavior).

AOF rewrite and RDB write need extra space. A rewrite can need room for the new file beside the old file. A volume at 95 percent can fail even if the final file would fit.

`INFO persistence` fields such as `rdb_last_bgsave_status` and `aof_last_write_status` tell you the last result.

If the disk is full:

- Free space (old dumps, logs) on a lab or with a known backup policy
- Do not delete the only live AOF without a plan
- Fix the volume size before you raise traffic

Replication does not replace a full disk on the primary. Do not put `dir` on a tiny container overlay disk without a mount.

An upgrade changes the Redis binary or image. Plan it.

Typical replica-first path for a primary plus replicas:

1. Upgrade a replica. Check replication and `INFO`.
2. Upgrade the other replicas.
3. Fail over (Sentinel or Cluster) so that an upgraded replica becomes primary.
4. Upgrade the old primary (now a replica).
5. Test the application against new command behavior.

Read the release notes. Some versions change defaults (ACL, AOF, replica names). Modules: upgrade Redis Stack as a bundle.

Rollback: keep the previous image and a backup from before the upgrade. A data-format change can make rollback harder.

Do not upgrade the only primary in place during peak without a replica and a backup.

`CONFIG GET` / `CONFIG SET` change live settings. `CONFIG REWRITE` writes them into `redis.conf`. Without a rewrite, a restart can restore old values. Restrict `CONFIG SET` with ACLs. Do not experiment on a durable production primary.

A replica is not a backup. A backup is a file copy that you can restore after a bad `FLUSHALL` or a lost host. Trigger `BGSAVE`, wait until it finishes, copy the RDB off the host, checksum the copy, and test restore on a lab instance.

Watch `INFO` for hit rate, `evicted_keys`, `connected_clients`, `blocked_clients`, `instantaneous_ops_per_sec`, and last save status. Alert on eviction on a durable instance, failed saves, `rejected_connections`, and memory near `maxmemory`.

### Questions

#### Theoretical questions

1. Why does RSS grow during `BGSAVE` if the parent writes many keys?
2. Why can a volume at 95 percent still fail a rewrite?
3. Why upgrade a replica before the primary?
4. Why is a replica not a backup?
5. What does a `fork` error in the log often mean?

#### Easy practical tasks

1. Run `BGSAVE`. Watch `rdb_bgsave_in_progress`. Write `LASTSAVE` after.
2. `CONFIG GET dir`. Save the path.
3. Write `redis_version` from `INFO server`.
4. Write four sentences: fork vs `SAVE` on the main thread.

#### Medium practical tasks

1. During `BGSAVE`, rewrite a large set of keys in a lab. Record RSS peak vs idle.
2. On a disposable volume, fill the disk until `BGSAVE` fails (or set a tiny quota). Record the error and `rdb_last_bgsave_status`. Then free space.
3. `BGSAVE`. Wait. Copy the RDB to a lab folder. Restore on a second Docker port. Confirm a known key.

#### Advanced practical tasks

1. Size a host for 8 GB `used_memory` and a busy write rate. Write RAM and disk numbers with a COW caveat.
2. Write an upgrade matrix: server, Stack modules, your client, OS. Fill current and target. Include backup and rollback.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do `INFO memory`, `MEMORY USAGE`, and `--bigkeys` answer three different questions?
2. How do bind, firewall, TLS, and ACLs each stop a different failure (reach, sniff, act)?
3. Which dangerous commands hurt availability (`KEYS`, `DEBUG SLEEP`) versus durability (`FLUSHALL`)?
4. How do `CONFIG SET`, disk fill, and COW interact during one `BGSAVE` on a small host?
5. Why do upgrade order and backup exist in the same change window?

#### Easy practical tasks

1. Write a cheat sheet: `INFO memory`, `MEMORY USAGE`, `--bigkeys`, `SLOWLOG GET`, `ACL LIST`, `bind`, `--tls`, `FLUSHALL`, `BGSAVE`.
2. Export `INFO memory` and highlight `used_memory_human`, `used_memory_rss_human`, `maxmemory`.
3. Draw RAM (used + COW headroom) and disk (RDB + rewrite room).
4. Write five lab rules: local bind, no production `FLUSHALL`, no `KEYS`, no `SAVE` in prod, copy files before repair.

#### Medium practical tasks

1. Create a lab dataset with one big hash and 1000 small keys. Find the big hash with `--bigkeys` and with `SCAN`+`HLEN`. Write both paths.
2. Harden a disposable lab: ACL app user, bind localhost, `FLUSHALL` denied. Connect as app and as ops. Record both sessions.
3. Run a mini drill: change a harmless config, `BGSAVE`, copy RDB, restore on another port, read a key. Write the times.

#### Advanced practical tasks

1. Write a monthly Redis health script outline: memory, big keys on a replica, slow log snapshot, hit rate. No production `KEYS`.
2. Produce a one-page Redis security and operations baseline for your team. Include versions (ACL, TLS), network, commands, COW headroom, and backup restore.
