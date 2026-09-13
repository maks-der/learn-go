# 4. Expiration and Eviction

## Description

Redis can delete a key after a time. Redis can also delete keys when memory is full. This topic covers expiry commands, eviction policies, and `maxmemory`. You learn what happens when the instance cannot store a new write.

Complete this topic before you rely on Redis as a cache in later designs.

Use one term for each concept. Expiry is a timer on a key. TTL is the remaining time. Eviction is removal under memory pressure. `maxmemory` is the memory limit. A policy is the rule that selects keys to evict. Do not mix `expired_keys` and `evicted_keys`.

---

## `EXPIRE`, `PEXPIRE`, `TTL`, `PTTL`

`EXPIRE key seconds` sets a lifetime in whole seconds. Redis deletes the key when the time ends. The reply is `1` if Redis set the expiry, or `0` if the key does not exist.

`PEXPIRE key milliseconds` uses milliseconds. Use `PEXPIRE` when one-second steps are too coarse.

`TTL key` returns the remaining seconds. The reply is:

- a non-negative integer — seconds left
- `-1` — the key exists and has no expiry
- `-2` — the key does not exist (or already expired)

`PTTL key` returns remaining milliseconds, or the same `-1` and `-2` rules.

```text
SET session:abc data
EXPIRE session:abc 300
TTL session:abc
```

`PERSIST key` removes the expiry and keeps the value. `TTL` then returns `-1`.

`EXPIRE` on a missing key does nothing useful (`0`). Set the value first, or use `SET key value EX seconds` in one command.

Redis 7 adds `EXPIRE` options such as `NX`, `XX`, `GT`, and `LT` on some versions. Those options set expiry only when a condition holds. Read the command page for your version.

Expiry uses server time. A large jump of the system clock can change when keys expire. Do not depend on millisecond precision for business-critical legal deadlines. Use expiry for caches, sessions, and temporary holds.

A read of an expired key looks like a missing key. The application must handle null.

### Questions

#### Theoretical questions

1. What is the difference between `EXPIRE` and `PEXPIRE`?
2. What does `TTL` return when the key has no expiry?
3. What does `TTL` return when the key is missing?
4. What does `PERSIST` do?
5. Why is `SET ... EX` often better than `SET` plus `EXPIRE`?

#### Easy practical tasks

1. `SET k 1` then `EXPIRE k 20` then `TTL k`. Save the replies.
2. Run `TTL nosuch`. Save the reply.
3. `PERSIST k` then `TTL k`. Save the reply.
4. `SET m 1 PX 2500` or `PEXPIRE` then `PTTL m`. Save the reply.

#### Medium practical tasks

1. Compare `SET k v EX 10` with `SET k v` plus `EXPIRE k 10` from two clients at once. Write the race on the two-command path.
2. If your server has `EXPIRE ... GT`, test it. Otherwise write the options from the command page.
3. Watch `TTL` once per second for five seconds on a key with `EX 8`. Write the values.

#### Advanced practical tasks

1. Read how Redis expires keys (active cycle plus lazy delete). Write six sentences with no extra sources beyond the official docs.
2. Measure whether `TTL` of `-2` can appear while `EXISTS` was `1` a moment earlier. Explain with expiry, not with a bug story.

---

## `EXPIREAT`

`EXPIREAT key unix-seconds` sets expiry at an absolute Unix time in seconds. Redis deletes the key when the server clock reaches that time.

`PEXPIREAT key unix-milliseconds` uses milliseconds since the Unix epoch.

```text
EXPIREAT session:abc 1893456000
```

Use `EXPIREAT` when many keys must die at the same wall-clock time (end of a day, end of a sale). Use `EXPIRE` when you want "live for 15 minutes from now".

The argument is UTC Unix time, not a local date string. Compute the timestamp in the application. Confirm the timezone of that computation.

If the timestamp is in the past, Redis deletes the key (or the command behaves as an immediate expiry). Test the exact reply on your version. Do not store past times by mistake.

`SET key value EXAT timestamp` is the one-command form on versions that support `EXAT`. Prefer one command when you create the key and the deadline together.

`TTL` after `EXPIREAT` still shows remaining seconds. `TTL` is relative. `EXPIREAT` is absolute.

A replica uses the expiry time that the primary stored. Clock skew between hosts can still surprise you in operations. Keep NTP healthy on Redis hosts.

### Questions

#### Theoretical questions

1. How is `EXPIREAT` different from `EXPIRE`?
2. What unit does `PEXPIREAT` use?
3. When do you choose `EXPIREAT` instead of `EXPIRE`?
4. Who computes the Unix timestamp — Redis or the application?
5. What should you check if the timestamp is already in the past?

#### Easy practical tasks

1. Compute "now + 120 seconds" as a Unix time. `SET e 1` then `EXPIREAT e <that time>`. Run `TTL e`.
2. Open the `EXPIREAT` command page. Write the time complexity and the return type.
3. Write four sentences that compare relative TTL and absolute expiry.
4. Run `TIME` on Redis if the command exists. Write the server Unix time.

#### Medium practical tasks

1. Set two keys to expire at the same `EXPIREAT` timestamp. Confirm both `TTL` values stay close.
2. Use `SET ... EXAT` if available. Compare with `SET` plus `EXPIREAT`.
3. Set `EXPIREAT` to a past time on a lab key. Record whether the key remains.

#### Advanced practical tasks

1. Write a small program that expires all "day" keys at the next midnight UTC with `EXPIREAT`. Print the timestamp and three sample keys.
2. Read about expire accuracy and replicas in the docs. Write three operational risks (clock, lag, past timestamps).

---

## Eviction policies: `noeviction`, `allkeys-lru`, `volatile-lru`, LFU, random

When `maxmemory` is set and Redis needs more RAM for a write, Redis applies an eviction policy. The policy is a `CONFIG` value: `maxmemory-policy`.

Common policies:

- `noeviction` — Redis does not evict. Writes that need memory return an error. Reads still work.
- `allkeys-lru` — Redis evicts any key, least recently used first (approximate LRU).
- `volatile-lru` — Redis evicts only keys that have an expiry, LRU among those keys.
- `allkeys-lfu` — Redis evicts any key, least frequently used first (approximate LFU).
- `volatile-lfu` — LFU among keys that have an expiry.
- `allkeys-random` — Redis evicts a random key.
- `volatile-random` — random among keys that have an expiry.
- `volatile-ttl` — among keys with expiry, evict the key with the shortest TTL first.

LRU means "least recently used". A key that no client touched for a long time is a better eviction candidate than a hot key. Redis uses an approximation, not a perfect LRU list for every key.

LFU means "least frequently used". A key with few reads is a better candidate than a key with many reads. Redis stores a counter with decay.

`volatile-*` policies do nothing useful if no key has an expiry and memory is full. Writes then fail like `noeviction`. If you choose `volatile-lru`, give cache keys a TTL.

`allkeys-lru` is a common cache policy. Redis can evict keys that have no TTL. Durable keys without TTL can disappear. Do not store must-keep data on an instance that uses `allkeys-*` unless you accept loss.

Read `CONFIG GET maxmemory-policy`. Change a lab instance with `CONFIG SET maxmemory-policy allkeys-lru`. Persist the setting in `redis.conf` for a real service.

### Questions

#### Theoretical questions

1. What does `noeviction` do when memory is full?
2. What is the difference between `allkeys-lru` and `volatile-lru`?
3. What does LFU use that LRU does not use?
4. Why can `volatile-lru` behave like `noeviction`?
5. Why is `allkeys-lru` dangerous for data that must survive?

#### Easy practical tasks

1. Run `CONFIG GET maxmemory-policy`. Save the value.
2. Make a two-column table: policy name and one sentence.
3. Write which policy you would pick for a pure cache. Give one reason.
4. Open the eviction docs on redis.io. Write the policy list from the page.

#### Medium practical tasks

1. On a lab instance, set `allkeys-lru`, then set `volatile-lru`. Write how you would prove the difference in a later `maxmemory` test.
2. Explain approximate LRU in four sentences from the official docs (sampling).
3. Draw a decision tree: must-keep keys vs cache keys vs mixed instance.

#### Advanced practical tasks

1. Read `maxmemory-samples` (or the current sample setting). Write how sample size changes LRU quality vs CPU.
2. Compare LFU decay settings in the docs. Propose values for a session cache vs a catalog cache.

---

## `maxmemory`

`maxmemory` is the upper bound of memory that Redis should use. The value is in bytes. `CONFIG GET maxmemory` shows it. `0` on a 64-bit build means no limit. Redis then grows until the operating system refuses RAM.

Set `maxmemory` in production. Pick a value below the host RAM so that the OS, the AOF rewrite, and `BGSAVE` copy-on-write have room. A common pattern is 50 to 70 percent of RAM for a host that also persists. Topic 9 and operations topics cover forks.

```text
CONFIG SET maxmemory 64mb
```

`CONFIG SET` changes the running instance. Write the same value in `redis.conf` (`maxmemory 64mb`) so a restart keeps the limit.

`INFO memory` shows `used_memory`, `used_memory_rss`, and `maxmemory`. `used_memory` is the allocator view of data. `used_memory_rss` is the process size that the OS sees. They differ.

Replicas and clients also use memory. Buffers can grow. `maxmemory` is not only "sum of key sizes".

A replica can have a different `maxmemory` than the primary. If you evict on a replica, the replica can diverge. Many teams set replicas so that they do not evict, or they keep extra RAM.

Do not set `maxmemory` below the current `used_memory` without a plan. Redis may evict at once or reject writes, depending on the policy.

### Questions

#### Theoretical questions

1. What does `maxmemory` of `0` mean on a 64-bit Redis?
2. Why must `maxmemory` stay below total host RAM?
3. What is the difference between `used_memory` and `used_memory_rss`?
4. Does `CONFIG SET maxmemory` survive a restart by itself?
5. Why can replica `maxmemory` cause divergence?

#### Easy practical tasks

1. Run `CONFIG GET maxmemory` and `INFO memory`. Write `used_memory_human` and `maxmemory`.
2. Convert 128 MB to bytes. Write the number that you would pass to `maxmemory`.
3. Find `maxmemory` in a sample `redis.conf` on the web or in your image. Copy the line.
4. Write four sentences: why a cache instance needs `maxmemory`.

#### Medium practical tasks

1. On a lab container, set `maxmemory` to a small value (for example `10mb`). Confirm with `CONFIG GET`.
2. Compare `used_memory` before and after 10,000 small keys. Write the delta.
3. Read about `maxmemory` and output buffers in the docs. Write one risk of large client buffers.

#### Advanced practical tasks

1. Size `maxmemory` for a 4 GB host that also runs AOF. Show the arithmetic and the fork/COW warning.
2. Inspect `MEMORY STATS` or `INFO memory` for fragmentation. Write two fields and what you think they mean.

---

## What happens when Redis is full

"Full" means `used_memory` reached `maxmemory` and a command needs more memory.

If the policy is `noeviction`, Redis returns an error to write commands that grow memory. Typical text includes `OOM command not allowed when used memory > 'maxmemory'`. Read commands still run. `DEL` can free memory.

If the policy evicts keys, Redis removes one or more keys, then retries the write. You lose those keys. The write can succeed. `INFO stats` increases `evicted_keys`.

If Redis cannot evict enough (for example `volatile-lru` and no volatile keys), the write fails with an OOM error.

The application must handle write errors. A cache-aside app can skip the cache and use the primary database. An app that treats Redis as the only store can fail the user request. Design that path on purpose.

Clients that pipeline many writes can see a burst of errors. A monitor should alert on `evicted_keys` rate and on OOM errors in logs.

Redis does not pause forever to compact memory. Fragmentation can keep RSS high after deletes. `MEMORY PURGE` or a restart can help in some allocator setups. Do not rely on that in the application.

Do not wait for the OS OOM killer. Set `maxmemory` and a policy. The OS killer is harder to control than a Redis error.

### Questions

#### Theoretical questions

1. What error can a write receive when Redis is full and the policy is `noeviction`?
2. What happens to existing keys under `allkeys-lru` when a new `SET` needs space?
3. Can reads succeed when writes fail for memory?
4. What should a cache-aside application do when `SET` fails?
5. Why is the OS OOM killer a worse limit than `maxmemory`?

#### Easy practical tasks

1. Write the OOM error text from the docs or from a lab test.
2. Make a flowchart: full Redis, policy branch, success or error.
3. Run `INFO stats` and write `evicted_keys`.
4. List three application reactions to a failed cache write.

#### Medium practical tasks

1. Fill a tiny `maxmemory` instance with `allkeys-lru` until `evicted_keys` grows. Use small `SET`s. Record `INFO memory` and `INFO stats`.
2. Repeat with `noeviction`. Record the first write error. Delete keys to recover.
3. Repeat with `volatile-lru` and keys that have no TTL. Record whether writes fail.

#### Advanced practical tasks

1. Script a fill test: report `used_memory`, `evicted_keys`, and the first error. Restore `maxmemory` after the test.
2. Write a runbook page: alerts, `INFO` fields, and how the app behaves when Redis is full (cache vs database-of-record).

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do `EXPIRE`, `EXPIREAT`, and `SET ... EX` differ in the way you express time?
2. How do you tell expiry deletion from eviction using `INFO`?
3. Why do `volatile-*` policies require TTLs on keys that Redis may drop?
4. What configuration pair (`maxmemory` and policy) fits a cache, and what pair fits a store that must not drop keys?
5. What does a full Redis instance mean for the next `SET` vs the next `GET`?

#### Easy practical tasks

1. Create `lab:e` with `EX 15`, `lab:p` with `PERSIST` after an expire, and `lab:n` with no TTL. Fill a table: key, `TTL`, meaning.
2. Write a cheat sheet: `EXPIRE`, `PEXPIRE`, `EXPIREAT`, `TTL`, `PTTL`, `PERSIST`, `maxmemory-policy` names.
3. Export `INFO memory` and `INFO stats`. Mark `used_memory`, `maxmemory`, `expired_keys`, `evicted_keys`.
4. Draw one timeline: key created, TTL set, key expires, client `GET`s null.

#### Medium practical tasks

1. Write a script that sets 100 keys with random TTLs between 10 and 30 seconds. Print how many remain after 20 seconds.
2. Change policy from `noeviction` to `allkeys-lru` on a lab instance. Document the two `CONFIG` commands and a restart note.
3. Document a cache policy page: default TTL, jitter idea (later topic), `maxmemory`, and eviction policy.

#### Advanced practical tasks

1. Compare `EXPIRE` of many keys vs `SET` with `EX` at write time in a load script. Write duration and a race note.
2. Read `latency doctor` or the slow log after a fill-and-evict test. Write one observation about eviction cost.
