# 16. Memory and Performance

## Description

Redis is fast when data fits in RAM and commands stay small. This topic shows how you read memory, how you find large keys, how Redis encodes small values, how pipelines cut wait time, and how the slow log records expensive commands.

Complete topic 2 (single-threaded execution), topic 4 (`maxmemory`), and topic 8 (pipelines) first.

Use one term for each concept. Used memory is the bytes that Redis reports for data and overhead. RSS is the process size that the operating system sees. A big key is one value that is large in bytes or in element count. An encoding is the in-memory layout of a value. The slow log is a list of commands that exceeded a time threshold.

---

## `INFO memory`

`INFO memory` prints memory fields as `key:value` lines. Useful fields (names can vary slightly by version):

- `used_memory` / `used_memory_human` — allocator view of Redis data
- `used_memory_rss` / `used_memory_rss_human` — operating system resident size
- `used_memory_peak_human` — peak since start
- `maxmemory` — configured limit (`0` often means no limit on 64-bit)
- `mem_fragmentation_ratio` — RSS divided by used memory (approximate)
- `allocator_*` — allocator detail on some builds

```text
INFO memory
```

`used_memory` and RSS are not the same. Fragmentation, file mappings, and copy-on-write during `BGSAVE` (topic 9, topic 19) raise RSS.

A ratio much above `1` can mean fragmentation or a recent fork. A ratio below `1` can appear with shared pages or allocator quirks. Do not treat one number as a verdict. Use `MEMORY DOCTOR` and the docs for your version.

`INFO stats` adds `expired_keys`, `evicted_keys`, `keyspace_hits`, and `keyspace_misses`. Hit rate is a derived number (topic 19).

Watch `used_memory` against `maxmemory`. When they meet, eviction or `OOM` errors begin (topic 4).

Do not sample `INFO` in a tight loop on the main thread of a huge instance without a need. The command is usually cheap. The habit of scraping every 10 seconds is normal.

### Questions

#### Theoretical questions

1. What is the difference between `used_memory` and `used_memory_rss`?
2. What does `maxmemory` `0` often mean on a 64-bit build?
3. What is `mem_fragmentation_ratio` in simple terms?
4. Which `INFO stats` counters relate to cache hits and evictions?
5. Why can `BGSAVE` change RSS?

#### Easy practical tasks

1. Run `INFO memory`. Write `used_memory_human` and `used_memory_rss_human`.
2. Run `CONFIG GET maxmemory`. Write the value.
3. Write four sentences: used memory vs RSS.
4. Run `INFO stats`. Write `evicted_keys` and `expired_keys`.

#### Medium practical tasks

1. `SET` 10,000 keys. Record `used_memory` before and after. Delete the keys. Record again.
2. Compute a simple hit rate from `keyspace_hits` and `keyspace_misses` (if both are nonzero). Write the formula and the number.
3. Compare `INFO memory` on idle vs during a `BGSAVE`. Write RSS change.

#### Advanced practical tasks

1. Read the official `INFO` field list for your version. Add three fields that this section omitted and one sentence each.
2. Graph or table `used_memory` every 5 seconds during a load and a `BGSAVE`. Write the peak and the cause you infer.

---

## `MEMORY USAGE`, `MEMORY DOCTOR`

`MEMORY USAGE key` estimates the bytes of one key, including overhead. It is an estimate. It is enough to find outliers.

```text
MEMORY USAGE lab:big
```

The reply is an integer (bytes) or null if the key is missing.

`MEMORY DOCTOR` returns text advice about the instance (fragmentation, peak, allocator). Treat it as a hint, not as a ticket that you close blindly.

`MEMORY STATS` returns a map of allocator and dataset numbers. Use it when `INFO memory` is not enough.

`MEMORY PURGE` asks the allocator to release free pages to the operating system when the allocator supports it. Effect varies. Use it in a lab first. Do not expect it to shrink a full dataset.

`MEMORY MALLOC-SIZE` exists on some builds for allocator debugging. You do not need it for application work.

Do not run `MEMORY USAGE` on every key in production in a tight loop without `SCAN` pacing. Each call is a command on the main thread. Batch your investigation.

A small `MEMORY USAGE` with a huge `HLEN` can still be a problem at `HGETALL` time. Size in bytes and size in element count are both "big" (next section).

### Questions

#### Theoretical questions

1. What does `MEMORY USAGE` return?
2. What kind of output does `MEMORY DOCTOR` return?
3. Why is `MEMORY USAGE` an estimate?
4. What does `MEMORY PURGE` try to do?
5. Why can a small byte size still be a dangerous key?

#### Easy practical tasks

1. `SET lab:m hello`. `MEMORY USAGE lab:m`. Save the number.
2. Run `MEMORY DOCTOR`. Write one sentence from the reply (your own words).
3. `MEMORY USAGE` a missing key. Save the reply.
4. Write four sentences: `INFO memory` vs `MEMORY USAGE` for one key.

#### Medium practical tasks

1. Compare `MEMORY USAGE` of a 1-field hash and a 1000-field hash. Write both numbers.
2. Run `MEMORY STATS`. Write two field names that you recognize.
3. After a large `DEL`, try `MEMORY PURGE` in a lab. Record `used_memory_rss` before and after (may not change).

#### Advanced practical tasks

1. `SCAN` 200 keys and print `MEMORY USAGE` for each. List the top five. Delete only lab keys.
2. Read how `MEMORY USAGE` samples nested structures in the docs. Write one limit of the estimate.

---

## Big keys (`--bigkeys`, `SCAN`)

A big key is a key with a large byte size or a large number of elements. Examples: a 50 MB string, a list of 10 million items, a hash with 1 million fields.

Risks:

- `GET`, `HGETALL`, `LRANGE 0 -1`, `SMEMBERS` copy or walk the whole value on the main thread
- Replication and RDB grow
- Network buffers grow
- Cluster cannot split one key across slots (topic 11)

Find big keys:

```text
redis-cli --bigkeys
redis-cli --memkeys
```

These tools use `SCAN`. They sample. They can miss a key. They still add load. Run them on a replica when you can (topic 10). Accept stale reads.

You can also `SCAN` and call `MEMORY USAGE`, `STRLEN`, `HLEN`, `LLEN`, `SCARD`, `ZCARD`, `XLEN` per key. Pace the loop.

Fixes:

- Split one huge collection into many keys (shards, time buckets)
- Store blobs in object storage; keep a pointer in Redis
- Use `HSCAN`, `LRANGE` pages, `SSCAN` instead of full dumps
- `UNLINK` instead of `DEL` for large values (lazy free)

Do not ignore a big key because `used_memory` still has headroom. The next `HGETALL` can stall the instance.

`--bigkeys` is a lab and a replica tool. Do not schedule it on the primary every minute.

### Questions

#### Theoretical questions

1. What two senses of "big" apply to a key?
2. Why is `HGETALL` on a huge hash dangerous?
3. How does `redis-cli --bigkeys` walk the keyspace?
4. Why run a big-key scan on a replica when you can?
5. Why can Cluster not split one huge key?

#### Easy practical tasks

1. Create a string of about 100 KB. `STRLEN` and `MEMORY USAGE`. Save both.
2. Run `redis-cli --bigkeys` on a lab instance. Write the summary lines.
3. Write four sentences: big key vs many small keys.
4. List five commands that walk a whole collection.

#### Medium practical tasks

1. `HSET` 50,000 fields. Time `HGETALL` vs an `HSCAN` loop. Write both times. `UNLINK` the key.
2. Compare `--bigkeys` and `--memkeys` output. Write what each ranks.
3. Design a split: one daily list vs 24 hourly lists. Write key names and a read plan.

#### Advanced practical tasks

1. On a replica (or a copy), run `--bigkeys` during a light write load. Write whether the command affected `INFO latency` or `slowlog` on the replica.
2. Write a big-key policy: max string bytes, max hash fields, how you alert, and the split pattern.

---

## Encoding internals: ziplist / listpack / hashtable (awareness)

Redis stores small collections in compact encodings to save RAM. When the collection grows past thresholds, Redis converts the value to a larger structure.

Historical and current names (awareness):

- ziplist — older compact list of entries (many versions used it for small hashes, lists, zsets)
- listpack — newer compact encoding that replaces ziplist in current Redis
- quicklist — a list of compact nodes (list implementation)
- hashtable (dict) — hash table for larger hashes and sets
- intset — small set of integers
- skiplist plus dict — larger sorted sets

You do not pick the encoding in the application. Redis converts when limits in `redis.conf` are crossed (`hash-max-listpack-entries`, `hash-max-ziplist-entries` on old versions, and similar keys).

`OBJECT ENCODING key` shows the current encoding:

```text
OBJECT ENCODING lab:h
```

Replies look like `listpack`, `hashtable`, `int`, `embstr`, `raw`, `quicklist`, `intset`, `skiplist`.

Why this matters:

- A small hash can use far less RAM than you expect from a naive pointer model
- After conversion, memory jumps
- Benchmarks on 10 fields do not predict 100,000 fields

Do not tune listpack limits in production without a memory test. A higher limit keeps compact encoding longer and can make a single command that rewrites the whole compact blob more expensive.

This section is awareness. You do not write a ziplist by hand.

### Questions

#### Theoretical questions

1. What command shows the encoding of a key?
2. What problem do ziplist and listpack solve?
3. What happens when a collection exceeds a configured limit?
4. Why can a 10-field hash mislead a 100,000-field plan?
5. Should the application choose the encoding?

#### Easy practical tasks

1. `HSET lab:e a 1`. `OBJECT ENCODING lab:e`. Save the reply.
2. Open `CONFIG GET *listpack*` or `*ziplist*`. Write two setting names.
3. Write four sentences: compact encoding vs hashtable.
4. `SET lab:i 1`. `OBJECT ENCODING lab:i`. Save the reply (`int` is common).

#### Medium practical tasks

1. Add fields to a hash until `OBJECT ENCODING` changes. Write the field count and the two encodings.
2. Compare `MEMORY USAGE` just before and just after the conversion. Write both numbers.
3. Read the config help for one `hash-max-*` setting. Write the two limits (entries and value size) if present.

#### Advanced practical tasks

1. Read the Redis documentation on listpack vs ziplist history. Write a six-sentence timeline.
2. Estimate RAM for 1 million tiny hashes vs one hash with 1 million fields. Use a lab sample of 10,000 and scale with a stated caveat.

---

## Pipelines vs many round-trips

Topic 8 explained pipelines. This section is the performance view.

Each request/response pays one network round-trip (RTT). On a 0.5 ms RTT, 10,000 `GET`s wait about 5 seconds for the network alone. A pipeline of 10,000 `GET`s pays a few RTTs for the whole batch (plus server time).

```text
redis-cli --pipe < commands.txt
```

Rules:

- Pipeline independent commands
- Bound the batch (hundreds to a few thousands)
- Do not assume atomicity (other clients interleave)
- Cluster clients split batches by slot (topic 11, topic 12)
- A pipeline of huge `GET`s can still stall the main thread and fill buffers

`MGET` and `MSET` reduce RTTs for strings on one slot. They are not a full substitute for a pipeline of mixed commands.

Measure:

- Wall time of a loop vs a pipeline
- `INFO stats` `instantaneous_ops_per_sec`
- Client memory during a huge pipeline

Do not pipeline a command that must see the previous reply before you choose the next command.

### Questions

#### Theoretical questions

1. What does a pipeline reduce?
2. Why must a pipeline batch have a size limit?
3. Does a pipeline make commands atomic as a group?
4. What extra split does a Cluster client do?
5. When is `MGET` enough instead of a pipeline?

#### Easy practical tasks

1. Write four sentences: loop of `GET` vs pipeline of `GET`.
2. Open your client pipeline API. Write the method names to start and exec.
3. Draw 3 RTTs for 3 loops and 1 RTT for 3 pipelined commands.
4. Run `redis-cli --help` and find `--pipe`. Write the flag line.

#### Medium practical tasks

1. Time 5,000 `INCR`s in a loop vs 5,000 in a pipeline. Write both times.
2. Pipeline `SET` of 100 keys, then `MGET` those keys. Confirm values.
3. Pipeline on Cluster across untagged keys (if you have Cluster). Record whether the client split the batch.

#### Advanced practical tasks

1. Test batch sizes 50, 500, and 10,000. Write time and any error or memory pressure.
2. Compare pipeline vs a Lua script that does 100 `INCR`s. Write atomicity and time differences.

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

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do `INFO memory`, `MEMORY USAGE`, and `--bigkeys` answer three different questions?
2. Why can a compact encoding hide a future memory jump at the same time that a pipeline hides network wait?
3. What makes a command "slow" on Redis even when the client has a large timeout?
4. How do big keys, single-threaded execution, and Cluster slots constrain a design together?
5. When do you prefer a replica for memory investigation?

#### Easy practical tasks

1. Write a cheat sheet: `INFO memory`, `MEMORY USAGE`, `OBJECT ENCODING`, `--bigkeys`, `--pipe`, `SLOWLOG GET`.
2. Export `INFO memory` and highlight `used_memory_human`, `used_memory_rss_human`, `maxmemory`.
3. List three commands you will not run on a huge collection in production.
4. Draw used memory vs RSS during idle, load, and `BGSAVE`.

#### Medium practical tasks

1. Create a lab dataset with one big hash and 1000 small keys. Find the big hash with `--bigkeys` and with `SCAN`+`HLEN`. Write both paths.
2. Pipeline 10,000 `PING`s. Compare time to 10,000 sequential `PING`s. Add one `SLOWLOG GET`.
3. Document encoding of a hash as it grows through the conversion point. Include `MEMORY USAGE`.

#### Advanced practical tasks

1. Write a monthly Redis health script outline: memory, big keys on a replica, slow log snapshot, hit rate. No production `KEYS`.
2. Reproduce a stall in a lab (big `HGETALL`). Show the slow log entry and the fix (HSCAN or split). Record times.
