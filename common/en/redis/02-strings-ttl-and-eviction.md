# 2. Strings, TTL, and Eviction

## Description

The string type is the basic Redis value. A string is a binary-safe blob. You can store text, JSON, a number, or raw bytes. This topic shows write and read commands, `SET` options, integer counters, expiry commands, eviction policies, and `maxmemory`. You learn what happens when the instance cannot store a new write.

Complete this topic before you study hashes, lists, and sets.

Use one term for each concept. A string key holds one blob. `SET` writes the blob. `GET` reads the blob. `INCR` treats the blob as a base-10 integer. Expiry is a timer on a key. TTL is the remaining time. Eviction is removal under memory pressure. `maxmemory` is the memory limit. A policy is the rule that selects keys to evict. Do not mix `expired_keys` and `evicted_keys`.

---

## `SET`, `GET`, `MSET`, `MGET`, `SET` options (`EX`, `NX`, `XX`)

`SET key value` writes a string. Redis creates the key or replaces the old value. The type becomes string. The reply is `OK`.

`GET key` reads the string. A missing key returns a null reply. A key of another type returns `WRONGTYPE`.

```text
SET user:1:name Ada
GET user:1:name
```

`MSET` writes many string keys in one command. All pairs apply as one atomic command. Other clients do not see a partial `MSET`.

```text
MSET user:1:name Ada user:1:city London
```

`MGET` reads many string keys in one command. The reply is an array. Each slot is a value or null. The order matches the argument order.

```text
MGET user:1:name user:1:city user:1:missing
```

Use `MGET` and `MSET` to cut round trips. The keys must exist as strings or be missing. One wrong type in `MGET` can fail the command.

A string has a maximum size of 512 MB. Do not store large files in Redis. Keep values small. A few kilobytes is a common cache size. Measure `MEMORY USAGE` when you are not sure.

`GET` does not change TTL by itself. `SET` without TTL options removes any old TTL on that key.

`SET key value NX` writes only when the key does not exist. The reply is `OK` on success or null when the key exists. `SETNX` is the same idea.

`SET key value XX` writes only when the key already exists. The reply is `OK` or null.

`SET key value EX seconds` sets the value and a TTL in seconds. `PX` uses milliseconds. `EXAT` and `PXAT` use a Unix time when your version supports them.

You can combine options:

```text
SET lock:order:9 worker1 NX EX 30
```

That command writes `worker1` only if `lock:order:9` is missing, and it expires in 30 seconds.

`GETSET key value` returns the old string and then writes the new string. Redis 6.2 and later also allow `SET key value GET`, which can combine with `NX` or `XX`. Prefer `SET ... GET` on new servers.

`KEEPTTL` keeps the current TTL when you replace the value. A plain `SET` would remove the TTL.

Read the `SET` command page for your version. The option list grew over time. Do not guess. Check `redis_version`.

### Questions

#### Theoretical questions

1. What does `SET` do when the key already exists?
2. What does `GET` return for a missing key?
3. What does `NX` require before `SET` writes?
4. What does `XX` require before `SET` writes?
5. What happens to an existing TTL when you run a plain `SET`?

#### Easy practical tasks

1. `SET` a key and `GET` it. Save both replies.
2. `MSET` two keys. `MGET` those two keys plus one missing name. Save the array.
3. `SET flag 1 NX`. Run the same command again. Save both replies.
4. `SET session:1 data EX 60` then `TTL session:1`. Save the TTL.

#### Medium practical tasks

1. Time 100 separate `SET` calls vs one `MSET` of 100 pairs from a script. Write the two times.
2. `SET k v EX 120`. Then `SET k v2 KEEPTTL` if your version supports it. Compare `TTL` before and after.
3. Try `MGET` where one key is a hash (create the hash with `HSET` first). Record the result.

#### Advanced practical tasks

1. Build a truth table: key exists or not, `NX` or `XX` or neither, expected reply. Test each row on your instance.
2. Find the documented maximum string size. Write a policy for the maximum value size in your app (with a number).

---

## `INCR` and related commands

`INCR key` reads the string as a base-10 integer, adds `1`, and writes the new integer. The reply is the new value. If the key is missing, Redis treats the old value as `0`.

`DECR` subtracts `1`. `INCRBY key n` adds the integer `n`. `n` can be negative. `DECRBY` subtracts `n`. `INCRBYFLOAT` adds a floating-point number.

The value must look like an integer (or a float for `INCRBYFLOAT`). A string such as `hello` causes an error. A hash key causes `WRONGTYPE`.

```text
INCR views:post:10
INCR views:post:10
INCRBY views:post:10 5
```

Each command is atomic. Two clients can `INCR` the same key. Redis does not lose a count from an interleaving of two single commands.

`INCR` does not set a TTL. If you need a counter that resets, set a TTL on the first increment (Lua in topic 6 can do that in one step). A common pattern is a key per time window: `views:post:10:2026-09-13`.

Integers have limits. Redis uses 64-bit signed integers for `INCR`. An overflow returns an error.

Do not use `GET`, add one in the client, and `SET`. That path is not atomic. Two clients can write the same number. Use `INCR`.

`APPEND key value` adds bytes to the end of the string. If the key is missing, `APPEND` creates it. `STRLEN key` returns the length in bytes. `GETRANGE` and `SETRANGE` edit a substring. Those commands are useful for fixed-width records. They are a poor fit for JSON that you parse in the client.

### Questions

#### Theoretical questions

1. What does `INCR` do when the key is missing?
2. What error occurs when the string is not an integer?
3. Why is `INCR` safer than `GET` plus `SET` in the client?
4. What is the integer size that `INCR` uses?
5. Does `INCR` set a TTL by itself?

#### Easy practical tasks

1. `INCR lab:count` three times. Save each reply.
2. `INCRBY lab:count 10` then `DECR lab:count`. Save the replies.
3. `SET lab:bad hello` then `INCR lab:bad`. Save the error.
4. Write four sentences that describe `INCR` and `INCRBY`.

#### Medium practical tasks

1. From two terminals, run `INCR lab:race` 1,000 times each. Write the final value and whether it is 2000.
2. `INCRBYFLOAT lab:float 0.5` twice. Save the replies.
3. Design daily counter key names for one metric. Write three example keys for three days.

#### Advanced practical tasks

1. Try to overflow a counter near the 64-bit limit (use `SET` to a large integer, then `INCR`). Record the error. Reset the key.
2. Implement a rate-limit preview: `INCR` a key `rl:user:1` and set `EX 60` only when the reply is `1`. Write the commands you used.

---

## `EXPIRE`, `TTL`, `EXPIREAT`

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

`EXPIREAT key unix-seconds` sets expiry at an absolute Unix time in seconds. Redis deletes the key when the server clock reaches that time. `PEXPIREAT` uses milliseconds since the Unix epoch.

Use `EXPIREAT` when many keys must die at the same wall-clock time (end of a day, end of a sale). Use `EXPIRE` when you want "live for 15 minutes from now".

The argument is UTC Unix time, not a local date string. Compute the timestamp in the application.

If the timestamp is in the past, Redis deletes the key (or the command behaves as an immediate expiry). Test the exact reply on your version.

`SET key value EXAT timestamp` is the one-command form on versions that support `EXAT`. Prefer one command when you create the key and the deadline together.

`TTL` after `EXPIREAT` still shows remaining seconds. `TTL` is relative. `EXPIREAT` is absolute.

Redis 7 adds `EXPIRE` options such as `NX`, `XX`, `GT`, and `LT` on some versions. Those options set expiry only when a condition holds. Read the command page for your version.

Expiry uses server time. A large jump of the system clock can change when keys expire. Do not depend on millisecond precision for business-critical legal deadlines. Use expiry for caches, sessions, and temporary holds.

A read of an expired key looks like a missing key. The application must handle null.

A key with no expiry lives until you delete it or until eviction removes it. Eviction is different from TTL. Topic counters `expired_keys` and `evicted_keys` in `INFO stats` are not the same.

Active expiry is a background cycle. Redis does not wait for a read to delete every expired key. Redis also checks expiry on access.

### Questions

#### Theoretical questions

1. What is the difference between `EXPIRE` and `EXPIREAT`?
2. What does `TTL` return when the key has no expiry?
3. What does `TTL` return when the key is missing?
4. What does `PERSIST` do?
5. Why is `SET ... EX` often better than `SET` plus `EXPIRE`?

#### Easy practical tasks

1. `SET k 1` then `EXPIRE k 20` then `TTL k`. Save the replies.
2. Run `TTL nosuch`. Save the reply.
3. `PERSIST k` then `TTL k`. Save the reply.
4. Compute "now + 120 seconds" as a Unix time. `SET e 1` then `EXPIREAT e <that time>`. Run `TTL e`.

#### Medium practical tasks

1. Compare `SET k v EX 10` with `SET k v` plus `EXPIRE k 10` from two clients at once. Write the race on the two-command path.
2. Set `EXPIREAT` to a past time on a lab key. Record whether the key remains.
3. Watch `TTL` once per second for five seconds on a key with `EX 8`. Write the values.

#### Advanced practical tasks

1. Read how Redis expires keys (active cycle plus lazy delete). Write six sentences from the official docs.
2. Write a small program that expires all "day" keys at the next midnight UTC with `EXPIREAT`. Print the timestamp and three sample keys.

---

## Eviction policies and `maxmemory`

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

LFU means "least frequently used". A key with few reads is a better candidate than a key with many reads.

`volatile-*` policies do nothing useful if no key has an expiry and memory is full. Writes then fail like `noeviction`. If you choose `volatile-lru`, give cache keys a TTL.

`allkeys-lru` is a common cache policy. Redis can evict keys that have no TTL. Durable keys without TTL can disappear. Do not store must-keep data on an instance that uses `allkeys-*` unless you accept loss.

`maxmemory` is the upper bound of memory that Redis should use. The value is in bytes. `CONFIG GET maxmemory` shows it. `0` on a 64-bit build means no limit. Redis then grows until the operating system refuses RAM.

Set `maxmemory` in production. Pick a value below the host RAM so that the OS, the AOF rewrite, and `BGSAVE` copy-on-write have room. A common pattern is 50 to 70 percent of RAM for a host that also persists.

```text
CONFIG SET maxmemory 64mb
```

`CONFIG SET` changes the running instance. Write the same value in `redis.conf` (`maxmemory 64mb`) so a restart keeps the limit.

`INFO memory` shows `used_memory`, `used_memory_rss`, and `maxmemory`. `used_memory` is the allocator view of data. `used_memory_rss` is the process size that the OS sees. They differ.

Replicas and clients also use memory. Buffers can grow. `maxmemory` is not only "sum of key sizes".

A replica can have a different `maxmemory` than the primary. If you evict on a replica, the replica can diverge. Many teams set replicas so that they do not evict, or they keep extra RAM.

Read `CONFIG GET maxmemory-policy`. Change a lab instance with `CONFIG SET maxmemory-policy allkeys-lru`. Persist the setting in `redis.conf` for a real service.

TTL is a per-key timer. Eviction is a memory-pressure mechanism. A cache often uses both: a TTL so that data does not stay stale forever, and an eviction policy so that Redis can drop keys when RAM is full.

### Questions

#### Theoretical questions

1. What does `noeviction` do when memory is full?
2. What is the difference between `allkeys-lru` and `volatile-lru`?
3. What does `maxmemory` of `0` mean on a 64-bit Redis?
4. Why must `maxmemory` stay below total host RAM?
5. Why is `allkeys-lru` dangerous for data that must survive?

#### Easy practical tasks

1. Run `CONFIG GET maxmemory-policy`. Save the value.
2. Run `CONFIG GET maxmemory` and `INFO memory`. Write `used_memory_human` and `maxmemory`.
3. Make a two-column table: policy name and one sentence.
4. Write which policy you would pick for a pure cache. Give one reason.

#### Medium practical tasks

1. On a lab container, set `maxmemory` to a small value (for example `10mb`). Confirm with `CONFIG GET`.
2. Compare `used_memory` before and after 10,000 small keys. Write the delta.
3. Explain approximate LRU in four sentences from the official docs (sampling).

#### Advanced practical tasks

1. Size `maxmemory` for a 4 GB host that also runs AOF. Show the arithmetic and the fork/COW warning.
2. Read `maxmemory-samples`. Write how sample size changes LRU quality vs CPU.

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

1. How do `NX`, `XX`, and `EX` change the meaning of `SET`?
2. When do you choose `INCR` instead of storing a number with `SET`?
3. How do `EXPIRE`, `EXPIREAT`, and `SET ... EX` differ in the way you express time?
4. How do you tell expiry deletion from eviction using `INFO`?
5. What configuration pair (`maxmemory` and policy) fits a cache, and what pair fits a store that must not drop keys?

#### Easy practical tasks

1. Run this sequence on a fresh key prefix `rv:`: `SET`, `GET`, `INCR`, `TTL` after `SET ... EX 90`. Save all replies.
2. Write a cheat sheet: `SET` options, `MGET`, `INCR`, `EXPIRE`, `EXPIREAT`, `TTL`, `maxmemory-policy` names.
3. Create `lab:e` with `EX 15`, `lab:p` with `PERSIST` after an expire, and `lab:n` with no TTL. Fill a table: key, `TTL`, meaning.
4. Export `INFO memory` and `INFO stats`. Mark `used_memory`, `maxmemory`, `expired_keys`, `evicted_keys`.

#### Medium practical tasks

1. Write a script that increments a counter 10,000 times with a pipeline or a loop. Print the final `GET` and the duration.
2. Write a script that sets 100 keys with random TTLs between 10 and 30 seconds. Print how many remain after 20 seconds.
3. Document a cache policy page: default TTL, `maxmemory`, and eviction policy.

#### Advanced practical tasks

1. Compare `EXPIRE` of many keys vs `SET` with `EX` at write time in a load script. Write duration and a race note.
2. Read `latency doctor` or the slow log after a fill-and-evict test. Write one observation about eviction cost.
