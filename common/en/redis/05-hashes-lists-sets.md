# 5. Hashes, Lists, Sets

## Description

Redis has three common collection types besides strings. A hash stores field-value pairs under one key. A list stores an ordered sequence of strings. A set stores unique strings with no order. This topic shows the main commands and when to pick each type.

Complete this topic before you study sorted sets.

Use one term for each concept. A field is a name inside a hash. An element is a string in a list. A member is a string in a set. Do not use `GET` on these keys. Use the commands for that type.

---

## Hashes: `HSET`, `HGET`, `HGETALL`, `HINCRBY`, `HDEL`

A hash is a map inside one key. Use a hash for an object with a few fields (profile, settings, rate-limit buckets).

`HSET key field value` writes one field. You can pass many field-value pairs in one `HSET`. The reply is the number of fields that Redis added (not the number that Redis only updated). Older servers had `HMSET`. New code uses `HSET`.

`HGET key field` reads one field. A missing field returns null. A missing key returns null.

`HMGET key field [field ...]` reads many fields in one command. The array can contain nulls.

`HGETALL key` returns all fields and values as a flat array. On a large hash, `HGETALL` is expensive and blocks. Prefer `HSCAN` for large hashes.

`HINCRBY key field integer` adds an integer to a field. The field must hold an integer or be missing. `HINCRBYFLOAT` adds a float.

`HDEL key field [field ...]` removes fields. `HEXISTS` tests one field. `HKEYS` and `HVALS` return all names or all values. `HLEN` returns the field count.

```text
HSET user:42 name Ada city London
HGET user:42 name
HINCRBY user:42 logins 1
HDEL user:42 city
```

A hash is not a second Redis instance. You still address the hash by one key. You cannot expire one field with a core command in older Redis. Redis 7.4+ adds field-level expiry on hashes in some builds. Check your version. Until then, expire the whole key or split fields into keys.

Do not use a hash with unbounded fields (one field per event forever). Memory grows. Use a stream, a list with `LTRIM`, or many keys with TTLs.

### Questions

#### Theoretical questions

1. What does a hash store under one key?
2. What does `HGET` return when the field is missing?
3. Why can `HGETALL` be a problem on a large hash?
4. What does `HINCRBY` require about the field value?
5. Can you expire a single hash field on every Redis version?

#### Easy practical tasks

1. `HSET` two fields on `lab:user:1`. `HGET` one field. Save the replies.
2. `HGETALL lab:user:1`. Save the reply.
3. `HINCRBY lab:user:1 score 5` twice. `HGET` the field.
4. `HDEL` one field. `HEXISTS` on that field. Save the replies.

#### Medium practical tasks

1. Compare a user as many string keys (`user:1:name`) vs one hash. Write six sentences on `MGET`, `HGETALL`, and TTL.
2. Use `HMGET` for three fields, one missing. Record the array.
3. Run `HLEN` before and after `HSET` of an existing field. Explain the `HSET` reply.

#### Advanced practical tasks

1. Load 10,000 fields into one hash. Compare time of `HGETALL` vs `HSCAN` loops. Write the times. Delete the key.
2. Read hash field expiry docs for your version. Write whether you can expire one field and the command names.

---

## Lists: `LPUSH`, `RPUSH`, `LPOP`, `RPOP`, `LRANGE`, `LTRIM`

A list is a sequence of strings. The left end and the right end are both usable. Redis does not use a numeric index as the primary write API. You push and pop at the ends.

`LPUSH key element [element ...]` inserts at the left (head). `RPUSH` inserts at the right (tail). The reply is the length after the push.

`LPOP key` removes and returns the left element. `RPOP` uses the right. A missing key or an empty list returns null. You can pass a count on versions that support `LPOP key N`.

`LRANGE key start stop` returns elements in that index range. Indexes are zero-based. `0` is the left-most element after the current pushes. `-1` is the last element. The stop index is inclusive.

```text
RPUSH lab:list a b c
LRANGE lab:list 0 -1
LRANGE lab:list 0 1
```

`LTRIM key start stop` keeps only that range and deletes the rest. Use `LTRIM` to cap a list (latest 100 items).

`LLEN` returns the length. `LINDEX` reads one index. `LSET` writes one index.

A list can grow without bound. An unbounded list can fill memory. Always plan a cap (`LTRIM`) or a consumer that pops.

`LRANGE 0 -1` of a huge list blocks and copies a large payload. Do not do that in production.

Blocking variants `BLPOP` and `BRPOP` wait for an element. They are useful for queues. The next section covers the queue pattern.

### Questions

#### Theoretical questions

1. Where does `LPUSH` insert an element?
2. What does `LPOP` return on an empty list?
3. Is the stop index of `LRANGE` inclusive?
4. What does `LTRIM` delete?
5. Why is `LRANGE 0 -1` dangerous on a large list?

#### Easy practical tasks

1. `RPUSH lab:q x y z`. `LRANGE lab:q 0 -1`. Save the reply.
2. `LPOP lab:q`. `LRANGE lab:q 0 -1`. Save both.
3. `LLEN lab:q`. Save the reply.
4. `RPUSH` three items then `LTRIM lab:q 0 1`. `LRANGE` the list.

#### Medium practical tasks

1. Build `1 2 3 4 5` with only `LPUSH` or only `RPUSH`. Show `LRANGE` and explain the order.
2. Use `LINDEX` and `LSET` to change the middle element. Record before and after.
3. `LPUSH` 1,000 items. Time `LRANGE 0 9` vs `LRANGE 0 -1`. Write the two times.

#### Advanced practical tasks

1. Read list encoding notes (quicklist / listpack) at a high level. Write two sentences on why a huge `LRANGE` is costly.
2. Implement a "latest 20 events" key with `LPUSH` and `LTRIM 0 19` after each push. Prove the length stays 20.

---

## Lists as queues

A queue is a list plus a rule: one end receives work, the other end gives work.

A FIFO queue (first in, first out) uses `LPUSH` to enqueue and `RPOP` to dequeue, or `RPUSH` and `LPOP`. Pick one pair and keep it.

```text
LPUSH jobs {"id":1}
RPOP jobs
```

A LIFO stack uses push and pop on the same end (`LPUSH` and `LPOP`).

`BLPOP` and `BRPOP` block the client until an element exists or a timeout ends. A worker can wait without a busy loop.

```text
BRPOP jobs 5
```

The reply is the key name and the element, or null on timeout. One worker can wait on several lists: `BRPOP high normal 5`. Redis pops the first list that has data, in argument order.

Many workers can pop from the same list. Each job goes to one worker. That is competing consumers. Redis does not track "in progress" on a plain list. If a worker pops a job and then dies, the job is gone.

Reliable queues need more structure: a processing list, a timeout, and a retry. Redis Streams (topic 7 and later stream topics) give consumer groups and acknowledgements. Use Streams when you need that model. Use a list when a lost job is acceptable or when a single worker is enough.

Do not `LRANGE` the queue in the worker hot path. Pop.

`RPOPLPUSH` or `LMOVE` (newer) can move an element to a processing list atomically. That is a better list-based reliability pattern than a single pop.

### Questions

#### Theoretical questions

1. Which command pair gives FIFO on a list?
2. What does `BRPOP` do when the list is empty?
3. What happens to a job if a worker `RPOP`s and then the worker dies?
4. What is a competing consumer on a Redis list?
5. Which newer type improves acknowledgement of jobs?

#### Easy practical tasks

1. Enqueue three jobs with `LPUSH`. Dequeue with `RPOP` three times. Write the order.
2. Run `BRPOP lab:empty 2` and wait. Save the reply.
3. Draw the list with left and right labels and arrows for enqueue and dequeue.
4. Write four sentences: list queue vs stack.

#### Medium practical tasks

1. Open two `redis-cli` windows. `BRPOP` in one. `LPUSH` in the other. Record both sides.
2. Use `LMOVE` or `RPOPLPUSH` from `jobs` to `jobs:busy`. Document the two lists after the move.
3. Start two blocking pop clients on the same list. Push one job. Write which client received it.

#### Advanced practical tasks

1. Design a list-based retry: busy list, `LRANGE` of stuck items, requeue. Write the failure cases. Compare with Streams in five sentences.
2. Time 10,000 `LPUSH` plus 10,000 `RPOP` from a script. Write throughput and the final `LLEN`.

---

## Sets: `SADD`, `SREM`, `SMEMBERS`, `SISMEMBER`, `SINTER`, `SUNION`

A set is a collection of unique members. Order is not defined. A duplicate `SADD` does not add a second copy.

`SADD key member [member ...]` adds members. The reply is the number of members that were new.

`SREM key member [member ...]` removes members.

`SISMEMBER key member` returns `1` or `0`. `SMISMEMBER` tests many members on versions that support it.

`SMEMBERS key` returns all members. On a large set, `SMEMBERS` is expensive. Use `SSCAN` for large sets.

`SCARD` returns the size.

`SINTER key [key ...]` returns the intersection. `SUNION` returns the union. `SDIFF` returns members in the first set and not in the following sets. `SINTERSTORE`, `SUNIONSTORE`, and `SDIFFSTORE` write the result to a destination key.

```text
SADD lab:tags:1 redis cache
SADD lab:tags:2 cache linux
SINTER lab:tags:1 lab:tags:2
```

Sets are good for tags, unique visitors (when the set is not huge), and relationship indexes (`user:42:roles`). They are a poor unique-visitor store at internet scale. Use HyperLogLog (topic 7) when you need an approximate count and not the exact ids.

A set member is a string. You cannot store a nested set. You can store ids and keep the objects in hashes.

Do not use a set as an unbounded event log. Use a list with `LTRIM` or a stream.

### Questions

#### Theoretical questions

1. What happens when you `SADD` a member that already exists?
2. What does `SISMEMBER` return?
3. Why is `SMEMBERS` risky on a large set?
4. What does `SINTER` return?
5. When is HyperLogLog a better fit than a set?

#### Easy practical tasks

1. `SADD lab:s a b c a`. Save the two replies if you split the adds, or `SCARD`.
2. `SISMEMBER lab:s b` and `SISMEMBER lab:s z`. Save both.
3. `SREM lab:s c`. `SMEMBERS lab:s`. Save the reply.
4. `SADD lab:t b d`. `SINTER lab:s lab:t`. Save the reply.

#### Medium practical tasks

1. Build two tag sets for two posts. Compute `SUNION` and `SDIFF`. Explain the meaning in the domain.
2. `SINTERSTORE lab:out lab:s lab:t`. `SMEMBERS lab:out`. Then `DEL lab:out`.
3. Compare `SISMEMBER` vs `SMEMBERS` plus a client search for one id. Write which path to use.

#### Advanced practical tasks

1. Add 50,000 members with a pipeline. Time `SISMEMBER` vs `SMEMBERS`. Delete the set after the test.
2. Design a "users who like post X" set and a "posts liked by user Y" set. Write how you keep both in sync on like and unlike.

---

## When to use each type

Pick the type from the access pattern, not from the domain word "object".

Use a string when you have one blob, a counter, or a lock. Use `GET` / `SET` / `INCR`.

Use a hash when you have one id and many small fields that you read or update one at a time. A user profile is the usual example. You can `HINCRBY` one field. You can `HGET` one field without the full object.

Use a list when order and ends matter: queues, stacks, and "latest N" feeds. You push, pop, and trim. You do not need random membership tests.

Use a set when uniqueness and membership matter: tags, roles, unique ids. You need `SISMEMBER`, intersection, or union. You do not need scores (that is a sorted set).

Examples:

- Session blob → string with TTL
- User profile → hash
- Job queue → list (or stream)
- Latest 50 notifications → list plus `LTRIM`
- Tags on a post → set
- All post ids with tag `redis` → set (secondary index)

You can combine types. A post hash `post:9` plus a set `tag:redis:posts` that contains `9` is a normal design.

If you need rank or a numeric score, use a sorted set (topic 6). If you need an append-only log with ids, use a stream (topic 7).

Move a wrong type only with a migration: read, write the new key, delete the old key. Redis does not convert a hash into a list in place.

### Questions

#### Theoretical questions

1. Which type fits a counter?
2. Which type fits a profile with five fields?
3. Which type fits a FIFO job queue in this topic?
4. Which type fits tag membership tests?
5. What type do you need if members must have scores?

#### Easy practical tasks

1. Make a four-row table: problem, type, one command.
2. Name one bad fit for each type (hash, list, set).
3. Sketch keys for a blog post: hash for body fields, set for tags.
4. Write four sentences that explain why a queue is not a set.

#### Medium practical tasks

1. Model a shopping cart three ways (string JSON, hash of sku→qty, list of events). Write a winner and why.
2. Implement the winner in `redis-cli` with three operations (add, read, remove).
3. Find a `WRONGTYPE` by using `LPUSH` on a hash key. Record the error. Delete the key.

#### Advanced practical tasks

1. Write a one-page type-selection guide with eight access patterns. Include TTL notes.
2. Review a public Redis schema (a blog post or module docs). List three keys and whether you agree with the types.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do `HGET`, `LINDEX`, and `SISMEMBER` differ in what they identify inside the key?
2. Which commands in this topic can block a client on purpose, and which can block the server because the payload is huge?
3. How do you cap memory for a hash, a list, and a set?
4. When do you store a result with `SINTERSTORE` instead of `SINTER`?
5. How do you model one user with a profile, a session, roles, and a notification feed using these types?

#### Easy practical tasks

1. Create `lab:u` (hash), `lab:feed` (list), `lab:roles` (set). Run `TYPE` on each. Save the three types.
2. Write a cheat sheet: five hash commands, six list commands, six set commands.
3. Draw the three types with a tiny example each (two fields, three elements, two members).
4. `DEL` the three lab keys. Run `EXISTS` on each.

#### Medium practical tasks

1. Build a mini app in the CLI: user hash, role set, notification list with trim to 5. Document every command.
2. Write a script that fills a list, a set, and a hash with 100 items each and prints `MEMORY USAGE` for each key.
3. Document a queue worker in ten steps using `BRPOP` and a second list for failures (`LPUSH jobs:dead`).

#### Advanced practical tasks

1. Implement like/unlike with a set and `HINCRBY` on a post hash so that `likes` and membership stay consistent. Note the two-command race. Mention Lua (topic 8) as the fix.
2. Compare `HSCAN`, `LRANGE` paging, and `SSCAN` for export of large keys. Write a safe export plan.
