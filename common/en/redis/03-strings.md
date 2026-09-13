# 3. Strings

## Description

The string type is the basic Redis value. A string is a binary-safe blob. You can store text, JSON, a number, or raw bytes. This topic shows the write and read commands, the `SET` options, integer counters, and range edits. You also see a preview of counters and simple locks.

Complete this topic before you study expiry in depth and before you study hashes and lists.

Use one term for each concept. A string key holds one blob. `SET` writes the blob. `GET` reads the blob. `INCR` treats the blob as a base-10 integer. A lock preview uses `SET` with `NX` and a TTL. Do not mix string commands with hash commands on the same key.

---

## `SET`, `GET`, `MSET`, `MGET`

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

`GET` does not change TTL by itself. `SET` without TTL options removes any old TTL on that key. Topic 4 covers `KEEPTTL` and expiry options.

### Questions

#### Theoretical questions

1. What does `SET` do when the key already exists?
2. What does `GET` return for a missing key?
3. Why is `MSET` useful compared with many `SET` commands?
4. What is the order of values in an `MGET` reply?
5. What happens to an existing TTL when you run a plain `SET`?

#### Easy practical tasks

1. `SET` a key and `GET` it. Save both replies.
2. `MSET` two keys. `MGET` those two keys plus one missing name. Save the array.
3. Run `GET` on a name that you never set. Save the reply.
4. Write four sentences that describe `SET` and `MSET`.

#### Medium practical tasks

1. Time 100 separate `SET` calls vs one `MSET` of 100 pairs from a script. Write the two times.
2. `SET` a key. Run `TYPE` and `STRLEN`. Explain both replies.
3. Try `MGET` where one key is a hash (create the hash with `HSET` first). Record the result.

#### Advanced practical tasks

1. Find the documented maximum string size. Write a policy for the maximum value size in your app (with a number).
2. Write a small client program that `MSET`s 1,000 keys and `MGET`s them in batches of 50. Record total time.

---

## `SETNX`, `GETSET`, `SET` options (`EX`, `NX`, `XX`)

Older Redis versions shipped `SETNX` and `GETSET` as separate commands. Current Redis still accepts them. New code prefers `SET` with options.

`SET key value NX` writes only when the key does not exist. The reply is `OK` on success or null when the key exists. `SETNX` is the same idea.

`SET key value XX` writes only when the key already exists. The reply is `OK` or null.

`SET key value EX seconds` sets the value and a TTL in seconds. `PX` uses milliseconds. `EXAT` and `PXAT` use a Unix time (when your version supports them).

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

1. What does `NX` require before `SET` writes?
2. What does `XX` require before `SET` writes?
3. What does `EX` add to `SET`?
4. How does `SET ... GET` relate to `GETSET`?
5. What does `KEEPTTL` prevent?

#### Easy practical tasks

1. `SET flag 1 NX`. Run the same command again. Save both replies.
2. `SET flag 2 XX`. Save the reply.
3. `SET session:1 data EX 60` then `TTL session:1`. Save the TTL.
4. Open the `SET` command page. List the options that your page shows.

#### Medium practical tasks

1. Compare `SETNX` with `SET ... NX` on the same key. Write if the replies match.
2. `SET k v EX 120`. Then `SET k v2 KEEPTTL` (if your version supports it). Compare `TTL` before and after.
3. Use `SET ... GET` or `GETSET` to rotate a token. Record the old and new values.

#### Advanced practical tasks

1. Build a truth table: key exists or not, `NX` or `XX` or neither, expected reply. Test each row on your instance.
2. Read the history section of `SET`. Write which Redis version added `GET` and `KEEPTTL`.

---

## `INCR` / `DECR` / `INCRBY`

`INCR key` reads the string as a base-10 integer, adds `1`, and writes the new integer. The reply is the new value. If the key is missing, Redis treats the old value as `0`.

`DECR` subtracts `1`. `INCRBY key n` adds the integer `n`. `n` can be negative. `DECRBY` subtracts `n`. `INCRBYFLOAT` adds a floating-point number.

The value must look like an integer (or a float for `INCRBYFLOAT`). A string such as `hello` causes an error. A hash key causes `WRONGTYPE`.

```text
INCR views:post:10
INCR views:post:10
INCRBY views:post:10 5
```

Each command is atomic. Two clients can `INCR` the same key. Redis does not lose a count from an interleaving of two single commands.

`INCR` does not set a TTL. If you need a counter that resets, set a TTL on the first increment (see the lock and counter preview, and Lua in topic 8). A common pattern is a key per time window: `views:post:10:2026-09-13`.

Integers have limits. Redis uses 64-bit signed integers for `INCR`. An overflow returns an error.

Do not use `GET`, add one in the client, and `SET`. That path is not atomic. Two clients can write the same number. Use `INCR`.

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
2. Implement a rate-limit preview: `INCR` a key `rl:user:1` and set `EX 60` only when the reply is `1`. Write the commands you used (Lua can wait until topic 8).

---

## `APPEND`, `GETRANGE`, `SETRANGE`

`APPEND key value` adds bytes to the end of the string. If the key is missing, `APPEND` creates it. The reply is the new length.

`STRLEN key` returns the length in bytes. A missing key returns `0`.

`GETRANGE key start end` returns a substring. Indexes are zero-based. A negative index counts from the end. The end index is inclusive.

```text
SET word HELLO
GETRANGE word 0 1
GETRANGE word -2 -1
```

`SETRANGE key offset value` overwrites bytes that start at `offset`. If the string is shorter, Redis pads with zero bytes. The reply is the new length.

These commands edit a blob without a full `GET` and `SET`. They are useful for fixed-width records and for small append-only logs. They are a poor fit for JSON that you parse in the client. A hash or JSON module is clearer for objects.

`SETRANGE` can create a large string if `offset` is large. Redis then allocates RAM. Do not use a huge offset on a lab instance without care.

`GETRANGE` of a huge range copies many bytes to the client. Prefer small ranges.

Bit commands (`SETBIT`, `GETBIT`) also treat a string as a byte array. Topic 7 covers bitmaps.

### Questions

#### Theoretical questions

1. What does `APPEND` do when the key is missing?
2. Is the end index of `GETRANGE` inclusive or exclusive?
3. What does Redis write in the gap when `SETRANGE` uses an offset past the current end?
4. What does `STRLEN` return for a missing key?
5. Why is `SETRANGE` risky with a very large offset?

#### Easy practical tasks

1. `SET msg Hi` then `APPEND msg !`. `GET msg` and `STRLEN msg`. Save the replies.
2. `GETRANGE msg 0 0`. Save the reply.
3. `SETRANGE msg 0 h`. `GET msg`. Save the reply.
4. Write a two-column table: "Command" and "What it changes".

#### Medium practical tasks

1. Build the string `ABCDEF` with `SET` and two `SETRANGE` calls. Show `GET` after each step.
2. `GETRANGE` with start greater than end, and with indexes past the length. Record Redis behavior.
3. Compare `APPEND` plus `GET` with a client-side concatenate and `SET`. Write which path is atomic for the write.

#### Advanced practical tasks

1. Store a fixed-width record (for example 8-byte id + 8-byte count) with `SETRANGE`. Read each field with `GETRANGE`. Document the layout.
2. Read the complexity of `APPEND` and `SETRANGE` on the command pages. Write when a full `SET` is simpler.

---

## Counters and simple locks (preview)

A counter is a string that you change with `INCR`, `DECR`, or `INCRBY`. Typical uses are view counts, ticket numbers, and rate-limit windows.

A simple lock is a key that only one worker can create. The common pattern is:

```text
SET lock:resource-name holder-id NX EX 30
```

If the reply is `OK`, this worker holds the lock. If the reply is null, another worker holds the lock. The TTL ends the lock if the worker dies.

This pattern is a preview. It has limits:

- The worker must delete only its own lock. A naive `DEL lock:resource-name` can remove another worker's lock after expiry and reuse.
- One Redis instance is not a full consensus system. Multi-instance lock algorithms are a later debate (Redlock).
- Clock and pause issues exist in distributed systems.

For learning, use one instance, a short TTL, and a unique holder id. Release the lock with a check that the value still equals your holder id. Topic 8 shows Lua for that check in one atomic script.

Do not use a lock when `INCR` or a single Redis command already makes the update safe. Locks are for work that happens outside Redis (files, payments, slow APIs).

Counters and locks both depend on atomic `SET` and `INCR`. They are string patterns, not new types.

### Questions

#### Theoretical questions

1. Which commands implement a Redis counter?
2. Which `SET` options implement a simple lock acquire?
3. Why does a lock key need a TTL?
4. Why is a bare `DEL` on the lock key unsafe at release time?
5. When should you avoid a lock and use `INCR` instead?

#### Easy practical tasks

1. Build a counter `lab:tickets` with three `INCR` calls. Write the final number.
2. Acquire `lock:lab` with `SET ... NX EX 20` and a holder name. Save the reply.
3. Try to acquire the same lock from a second `redis-cli`. Save the reply.
4. Write four sentences: counter vs lock.

#### Medium practical tasks

1. After the lock expires, acquire it again with a different holder. Record `GET lock:lab`.
2. Write a release sequence that `GET`s the lock and `DEL`s only if the value matches. Note the race (two steps). 
3. Design counter keys for "likes per post" and "likes per post per day". Write example names.

#### Advanced practical tasks

1. Read a short public critique of simple Redis locks (or Redlock). Write three risks that this preview does not solve.
2. Write a Lua script on paper (topic 8 can run it) that releases the lock only if the value matches. List the `KEYS` and `ARGV`.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Which string commands write many keys, and which read many keys?
2. How do `NX`, `XX`, and `EX` change the meaning of `SET`?
3. When do you choose `INCR` instead of storing a number with `SET`?
4. How do `GETRANGE` and `SETRANGE` differ from `GET` and `SET` of the full value?
5. How do a counter and a simple lock both use the string type but solve different problems?

#### Easy practical tasks

1. Run this sequence on a fresh key prefix `rv:`: `SET`, `GET`, `INCR`, `APPEND`, `TTL` after `SET ... EX 90`. Save all replies.
2. Write a cheat sheet: `SET` options, `MGET`, `INCR`, `APPEND`, `GETRANGE`, `STRLEN`.
3. Store a small JSON object as a string. `GET` it. Write one limit of this approach vs a hash.
4. Create a lock key and a counter key. Run `TYPE` on both. Confirm both are `string`.

#### Medium practical tasks

1. Write a script that increments a counter 10,000 times with a pipeline or a loop. Print the final `GET` and the duration.
2. Implement a lab "unique email hold": `SET email:ada@example.com 1 NX EX 120`. Document success and conflict replies.
3. Document ten string commands in a table: command, atomic, changes TTL or not.

#### Advanced practical tasks

1. Use `DEBUG SLEEP` only on a private lab (or skip if disabled) vs a slow `GETRANGE` of a large string to feel blocking. Write a safer way to observe blocking (slow log). Do not use `DEBUG` on a shared server.
2. Build a tiny client that acquires a lock, increments a counter, and releases the lock with a value check. List the race you still have without Lua.
