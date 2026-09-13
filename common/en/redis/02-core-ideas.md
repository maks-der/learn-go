# 2. Core Ideas

## Description

Redis is a key-value store with typed values. This topic explains the data model, the execution model, and the memory model. You also learn TTL, eviction, and key name conventions.

Complete this topic before you study string commands and other types.

Use one term for each concept. A key is a string name. A value is the data that Redis stores for that key. A type is the kind of the value (string, hash, list, and others). TTL is the remaining time to live of a key. Eviction is the removal of keys when memory is full.

---

## Key-value model

Redis stores pairs. Each pair has one key and one value. The client names the key in the command. Redis finds the value for that key. There is no SQL `FROM` clause and no document collection name in the core model.

A missing key is not an error for most read commands. `GET` of a missing key returns a null reply. Write commands create the key when the key does not exist, unless the command says otherwise.

One key holds one value. That value has one type. You cannot store a hash and a string under the same key at the same time. You can delete the key and create it again with a new type.

Redis does not query by value in the basic model. You do not say "find all keys whose string equals `alice`". You design keys so that the application already knows the name, or you keep secondary indexes on purpose (sets, sorted sets, or a search module).

Keys live in one keyspace per logical DB. The instance can hold millions of keys when memory allows. Each key has a small overhead in RAM. Many tiny keys cost more memory than one larger structure.

The key-value model is simple. The application owns the naming scheme and the access paths. Redis stores and retrieves. Redis does not replace your domain model. You map the domain to keys.

### Questions

#### Theoretical questions

1. What are the two parts of a Redis pair?
2. What does `GET` return when the key does not exist?
3. Can one key hold two types at the same time?
4. Why does Redis not offer a general "search by value" in the core model?
5. Who designs the key names in a Redis application?

#### Easy practical tasks

1. Draw a table with three rows: key, type, example value. Use only imaginary names.
2. Write four sentences that explain the key-value model to a person who knows SQL tables.
3. Run `GET nosuchkey` on an empty key name. Save the reply.
4. List three domain objects (user, session, cart). Propose one key name for each object.

#### Medium practical tasks

1. Compare a SQL row `users.id = 42` with a Redis key `user:42`. Write five sentences on lookup and on "find by email".
2. Create a key, read it, delete it with `DEL`, and read it again. Record each reply.
3. Try to `GET` a key after you create it as a different type in a later experiment, or write the error that Redis returns when types clash (`WRONGTYPE`).

#### Advanced practical tasks

1. Estimate RAM for one million keys that each hold a 16-byte string. Use public Redis overhead notes or `MEMORY USAGE` on a sample. Show the math.
2. Design key names for a blog (post, author, tag). Write which lookups are O(1) and which need a secondary key.

---

## Keys are strings; values have a type

A Redis key is a binary-safe string. You usually use UTF-8 text. A key can contain spaces, but spaces make CLI work harder. Prefer letters, digits, `:`, `-`, and `.`.

There is a maximum key size (512 MB). Do not use large keys. Short, stable names are the rule. A key such as `user:42:profile` is clear. A key that holds a full JSON document in the name is wrong.

The value has a type. The main types in this learning path are:

- string
- hash
- list
- set
- sorted set
- stream
- bitmap (a string used as bits)
- HyperLogLog (a string used as a sketch)
- geospatial (a sorted set with scores)

`TYPE key` returns the type name or `none` when the key is missing.

Commands belong to types. `HGET` works on a hash. `HGET` on a string key returns a `WRONGTYPE` error. `GET` works on a string. `GET` on a hash returns `WRONGTYPE`.

Some commands work on any key: `DEL`, `EXISTS`, `EXPIRE`, `TTL`, `TYPE`, `RENAME`. These commands do not read the inner structure.

Pick the type from the access pattern. A string fits a blob or a counter. A hash fits a small object with fields. A list fits a queue or a recent-item log. A set fits unique membership. A sorted set fits ranking and ranges.

Do not store every object as JSON in a string when you need field updates. A hash or a JSON module can update one field without a full rewrite. Later topics show those types.

### Questions

#### Theoretical questions

1. What is the type of a Redis key?
2. What command reports the type of a value?
3. What error does Redis return when a command does not match the type?
4. Name three commands that work on any key.
5. Why is a very large key name a bad design?

#### Easy practical tasks

1. Run `SET a 1` then `TYPE a`. Save the reply.
2. Run `TYPE missing`. Save the reply.
3. Write a two-column table: "Type" and "One good use". Add five rows from this section.
4. List four characters that you will allow in key names in your lab.

#### Medium practical tasks

1. Create a string key. Run `HGET` on that key. Record the error. Delete the key.
2. Measure `MEMORY USAGE` for a short key name and a 200-byte key name with the same value. Write the two numbers.
3. Read the `TYPE` command page. Write all type names that the page lists for your version.

#### Advanced practical tasks

1. Find the documented maximum key size. Write the number and the URL. Explain why you will never approach that limit.
2. Write a type-selection guide of one page: six access patterns and the Redis type you choose for each pattern.

---

## Single-threaded command execution (high-level; I/O threads exist)

Redis runs commands on one main thread. The instance takes one command, runs it to completion, then takes the next command. Two `INCR` commands on the same key do not interleave. You do not need a lock for a single command.

This model makes many operations simple. `INCR` is atomic. `LPUSH` is atomic. A Lua script or a Redis Function that runs on the server is also atomic with respect to other commands.

The cost is also clear. A slow command blocks other commands. `KEYS *` on a large keyspace, `LRANGE` of a huge list, or a heavy `SORT` can stall the instance. Keep commands small. Split large work.

Since Redis 6, optional I/O threads can read and write sockets. Those threads do not run the command logic. The main thread still executes commands. Do not think that Redis is a multi-core command engine because I/O threads exist.

Pipelines send many commands in one network round trip. The server still runs those commands one by one. Pipelines reduce wait time on the network. They do not run commands in parallel on many cores.

For one instance, CPU scale-out of command execution means more shards (Cluster) or more instances, not more threads on the same keyspace. Vertical speed still comes from fast commands and enough RAM.

### Questions

#### Theoretical questions

1. How many threads execute Redis commands on one instance?
2. Why is `INCR` safe without an application lock?
3. What happens to other clients when one command is very slow?
4. What do I/O threads do, and what do they not do?
5. Does a pipeline run commands in parallel on many cores?

#### Easy practical tasks

1. Write four sentences that explain single-threaded command execution.
2. Make a two-column table: "Feature" and "Atomic on one instance?". Add `INCR`, `GET`, and `PING`.
3. Open the Redis 6 release notes or I/O threads docs. Write one sentence about I/O threads.
4. List three commands that can be slow on large data (from this section or docs).

#### Medium practical tasks

1. Time a loop of 10,000 `INCR` calls from `redis-cli` or a script. Write the wall time.
2. Read about `UNLINK` vs `DEL` for large values. Write when a delete can still cost a lot of time.
3. Draw a sequence diagram: two clients send `INCR` on the same key. Show that the server runs one command, then the other.

#### Advanced practical tasks

1. Enable I/O threads in a local `redis.conf` (`io-threads`). Compare `INFO` fields before and after a load of `PING`s. Write what you can and cannot conclude.
2. Find a public incident caused by a slow Redis command (`KEYS`, big `HGETALL`, or similar). Write the command and the lesson in five sentences.

---

## In-memory first; persistence is optional

Redis reads and writes RAM first. A successful `SET` changes memory. The client receives `OK`. Disk is not in that path unless you configured a wait for durable write (that wait is not the default for every setup).

Persistence is optional. You can run Redis with no RDB and no AOF. A process stop then loses all keys. That mode is valid for a pure cache.

When you need data after a restart, you enable persistence. Redis can write snapshots (RDB). Redis can append each change to a file (AOF). You can combine both. Topic 9 covers the modes and the loss window.

Memory is the capacity limit. The instance size is the size of keys, values, and overhead. Disk size matters only when persistence is on.

Treat Redis as memory-first in application design. Do not assume that every `OK` is already on disk. If Redis is the database of record, you must choose a persistence mode and test restart. If Redis is a cache, you must accept loss and rebuild from the source of truth.

`INFO persistence` shows whether RDB or AOF is active and when the last save finished.

### Questions

#### Theoretical questions

1. Where does Redis apply a `SET` first?
2. What happens to keys if Redis stops and persistence is off?
3. Name the two persistence families that later topics cover.
4. What resource limits a Redis instance when persistence is off?
5. When must you test a restart of Redis?

#### Easy practical tasks

1. Run `INFO persistence`. Write whether AOF is enabled (`aof_enabled`).
2. Write three sentences: Redis as cache vs Redis as database of record.
3. Draw the path of `SET` from client to RAM. Put disk on the side as optional.
4. List two application types that can run Redis with no persistence.

#### Medium practical tasks

1. On a local instance with no persistence, `SET` a key, restart the container or process, and `GET` the key. Record the result.
2. Read the first page of the official persistence docs. Write the names of RDB and AOF in your own words (two sentences).
3. Compare `INFO memory` and `INFO persistence`. Write one fact from each section.

#### Advanced practical tasks

1. Start two local instances: persistence off, and AOF on. Restart both. Document which keys survive. Use only lab data.
2. Write a policy paragraph for your team: "Redis is a cache" or "Redis is durable." Include restart tests and a source of truth.

---

## TTL and eviction

TTL is time to live. You can attach an expiry to a key. When the TTL reaches zero, Redis deletes the key. Commands such as `EXPIRE` and `SET ... EX` set expiry. Topic 4 covers the full command set.

A key with no expiry lives until you delete it or until eviction removes it. `TTL key` returns the remaining seconds, `-1` if the key exists and has no expiry, or `-2` if the key is missing.

Eviction is different from TTL. Eviction runs when the instance reaches `maxmemory` and a write needs space. Redis then removes keys by a policy (`allkeys-lru`, `volatile-lru`, `noeviction`, and others). Topic 4 covers policies.

TTL is a per-key timer. Eviction is a memory-pressure mechanism. A cache often uses both: a TTL so that data does not stay stale forever, and an eviction policy so that Redis can drop keys when RAM is full.

Do not confuse "the key expired" with "Redis evicted the key". `INFO stats` has `expired_keys` and `evicted_keys`. Those counters are not the same.

Active expiry is a background cycle. Redis does not wait for a read to delete every expired key. Redis also checks expiry on access. You can see a key count drop without a client `DEL`.

### Questions

#### Theoretical questions

1. What does TTL mean for a Redis key?
2. What is the difference between TTL expiry and eviction?
3. What does `TTL` return when the key exists and has no expiry?
4. When does eviction run?
5. Why does a cache often use both a TTL and an eviction policy?

#### Easy practical tasks

1. Run `SET temp 1 EX 30` then `TTL temp`. Save the number (it can be 30 or 29).
2. Run `SET stay 1` then `TTL stay`. Save the reply.
3. Run `INFO stats` and write `expired_keys` and `evicted_keys`.
4. Write four sentences that explain expiry vs eviction to a beginner.

#### Medium practical tasks

1. Set a key with `EX 5`. Wait six seconds. `GET` the key. Record the reply.
2. Read `maxmemory` with `CONFIG GET maxmemory`. Write the value and what `0` means on a 64-bit build.
3. Draw a flowchart: write arrives, memory full or not, eviction policy, success or error.

#### Advanced practical tasks

1. Read the official page on expired key deletion (active expire cycle). Write how Redis finds expired keys without a client read.
2. Design TTL and eviction for a session cache and for a product catalog cache. Write why the two designs differ.

---

## Namespace keys with `:` or `#` conventions

Redis has one flat keyspace per logical DB. There are no folders. You simulate folders with a convention in the key string.

A common pattern uses `:` as a separator:

```text
user:42:profile
user:42:cart
session:abc123
cache:product:1001
```

Some teams use `#` or `/`. Pick one separator and keep it. The character has no special meaning to Redis except in Cluster hash tags, which use `{` and `}`. Topic 11 covers hash tags.

A good name is stable, unique, and readable. Include the object type, the id, and the field role when you need more than one key per object. Do not put secrets in key names. Key names appear in `MONITOR`, in logs, and in support tools.

Use a product or tenant prefix in a shared instance:

```text
shopA:user:42
shopB:user:42
```

Those keys do not collide. They still share memory and `maxmemory`. Multi-tenant isolation on one instance is a convention, not a security boundary.

Document the scheme. Write a one-page key catalog. Avoid ad-hoc names such as `data` or `tmp`. For temporary keys, include a purpose and a unique suffix, and set a TTL.

`SCAN` with a `MATCH` pattern can list keys by prefix on a lab instance. Do not use `KEYS user:*` in production.

### Questions

#### Theoretical questions

1. Does Redis have folders for keys?
2. Why do teams put `:` in key names?
3. Which characters are special for Cluster hash tags?
4. Why is a tenant prefix not a security boundary?
5. Why must you not put secrets in key names?

#### Easy practical tasks

1. Write five key names for a shop (user, cart, session, sku, cache). Use `:`.
2. Rewrite the same five names with `#`. State which style you prefer and why (one sentence).
3. Create `lab:user:1` and `lab:user:2` with `SET`. Run `SCAN 0 MATCH lab:user:*`. Save the output.
4. Write three bad key names and one reason each (secret, too long, no type).

#### Medium practical tasks

1. Draft a one-page key catalog for a todo app: pattern, type, TTL, owner.
2. Show a collision: two features use `user:1` for different meanings. Fix the names with a role suffix.
3. Use `SCAN` with `COUNT` and `MATCH` on 100 lab keys. Write why `SCAN` is safer than `KEYS` (from docs or this path).

#### Advanced practical tasks

1. Read the Cluster hash tag rules. Write how `{user:42}.profile` and `{user:42}.cart` relate to slots (preview is enough).
2. Design prefixes for two environments (`dev`, `prod`) on separate instances vs on one instance. Recommend one design and give two reasons.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do the key-value model and typed values work together when you choose a command?
2. Why can a single slow command make a fast instance look "down" to other clients?
3. How do optional persistence, TTL, and eviction change the meaning of "the data is in Redis"?
4. What parts of a key name are convention, and what parts does Redis interpret?
5. A teammate stores all objects as JSON strings. Which core ideas do you use to discuss types?

#### Easy practical tasks

1. Create three keys of different names. Run `TYPE`, `TTL`, and `EXISTS` on each. Save a small table.
2. Write a cheat sheet: `TYPE`, `TTL`, `EXISTS`, `DEL`, `SCAN`, and the `WRONGTYPE` error.
3. Draw one diagram that includes RAM, optional disk, TTL expiry, and eviction.
4. Write ten key names for a future app. Check that each name has a type word and an id.

#### Medium practical tasks

1. Write a short lab script that sets five keys, prints `INFO keyspace`, deletes them, and prints `INFO keyspace` again.
2. Configure a tiny `maxmemory` on a local instance (topic 4 preview). Write the current `maxmemory` and a one-sentence plan for a later test.
3. Document your key naming rules in ten lines so that another beginner can name keys the same way.

#### Advanced practical tasks

1. Use `MEMORY DOCTOR` or `MEMORY STATS` on a lab instance. Write three lines that you understand and two lines that you will study later.
2. Compare Redis key-value access with a document store lookup and a SQL primary-key lookup. Write nine short sentences (three per system).
