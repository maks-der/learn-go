# 1. Getting Started

## Description

Redis is an in-memory data store. Programs use Redis as a cache, as a message broker, and as a store for keys in RAM. This topic shows what Redis is, how Redis compares with two common alternatives, and how you start an instance. You also learn `redis-cli`, the key-value model, single-thread command execution, and key names.

Complete this topic before you study strings, TTL, and other types.

Use one term for each concept. An instance is one Redis process that runs. A client is a program that sends commands to that instance. A command is one request with a name and arguments. A key is a string name. A value is the data that Redis stores for that key. A type is the kind of the value.

---

## What Redis is (cache, broker, and in-memory store)

Redis keeps data in RAM. A client sends a command. Redis applies the command to keys in memory. Redis then sends a reply. The path is short. Latency is often less than one millisecond on a local network.

Redis is a data store. You write a key. You read that key. Redis also supports more than a plain string value. Later topics cover hashes, lists, sets, sorted sets, streams, and other types.

Teams use Redis in three common roles:

- Cache. The application writes a result that is expensive to compute. The next read uses Redis. A TTL can remove the key after a time.
- Broker. The application publishes messages or writes stream entries. Other processes read those messages.
- In-memory store. The application stores state in Redis and treats Redis as the source of truth. Persistence and replication become important.

Redis is not a replacement for every database. Redis does not replace a relational engine for complex joins and multi-row constraints. Redis is not a full message platform with long retention when you compare it with a log system such as Apache Kafka.

Redis is fast because data lives in memory and because command execution uses one main thread. Persistence to disk is optional. You can run Redis with no files on disk. You can also enable snapshots or an append-only file. Topic 7 covers those modes.

A typical production setup has one or more Redis instances, a client library in the application, and a clear policy for key names and TTLs. A typical learning setup is one instance on your machine and `redis-cli`.

### Questions

#### Theoretical questions

1. What does "in-memory" mean for data that Redis stores?
2. Name the three common roles of Redis in an application.
3. Why is Redis often faster than a disk-first database for a simple key read?
4. When is Redis a poor replacement for a relational database?
5. What is the difference between a Redis instance and a Redis client?

#### Easy practical tasks

1. Write five sentences that describe Redis. Use only facts from this section.
2. Make a two-column table: "Role" and "What Redis stores". Add one row for cache, one row for broker, and one row for in-memory store.
3. List three application features that fit Redis. List two features that do not fit well. Give one reason for each choice.
4. Open [https://redis.io/docs/](https://redis.io/docs/). Write the sentence that the site uses to describe Redis.

#### Medium practical tasks

1. Compare Redis with one database that you already know. Write six short sentences. Cover memory, types, and durability.
2. Draw a simple diagram: application, Redis instance, and disk. Label the path of a `GET` when data is only in RAM.
3. Find the current stable Redis version on redis.io. Write the version number and the date of the release notes page.

#### Advanced practical tasks

1. Read the Redis introduction pages on redis.io. Write a one-page timeline with years and one product fact per year.
2. Interview a teammate or read a public post-mortem. Write how that team used Redis (cache, broker, or store) and one risk they named.

---

## Redis vs Memcached vs a SQL cache table

Memcached is an in-memory cache. Memcached stores keys and values. Memcached has a small command set. Memcached does not give Redis data types such as lists, sets, and sorted sets. Memcached does not give Redis-style persistence, replication, or streams.

A SQL cache table is a table in a relational database. The application writes cache rows into that table. The same engine stores the cache and the durable data. SQL gives transactions, constraints, and query tools. A cache row in SQL is slower than a Redis key for a simple get. The cache also uses the same disk and the same locks as other tables.

Choose Redis when you need low latency, rich types, and a separate cache or session store. Choose Memcached when you need a simple cache and your team already runs Memcached. Choose a SQL cache table when the cache must join with other tables, when you cannot add a new service, or when the cache must use the same transaction as the durable write.

Redis and Memcached both evict keys when memory is full if you configure eviction. A SQL table grows until the disk is full unless the application deletes rows.

Redis can persist data. Memcached is a cache. A process restart of Memcached loses the cache. A process restart of Redis can reload data if you enable persistence.

Do not treat these three options as the same product. Measure the read path. Measure the operational cost. Use one primary store for durable business records unless you design Redis as that store on purpose.

### Questions

#### Theoretical questions

1. What can Redis store that Memcached does not store as first-class types?
2. Why can a SQL cache table be slower than Redis for a single-key read?
3. When is a SQL cache table a better choice than Redis?
4. What happens to Memcached data after a process restart?
5. Do Redis and Memcached both support eviction when memory is full?

#### Easy practical tasks

1. Make a three-column table: "Property", "Redis", "Memcached". Add four rows (types, persistence, typical role, command set).
2. Write four sentences that explain a SQL cache table to a beginner.
3. List two risks of putting the cache in the same SQL database as the durable data.
4. Open the Memcached protocol or project page. Write one command name that Memcached and Redis both use in a similar way.

#### Medium practical tasks

1. Design a session store three times: Redis, Memcached, and a SQL table. Write six sentences on expiry and lookup.
2. Find official or vendor docs that compare Redis and Memcached. Write three facts. Mark each fact as "official" or "opinion".
3. Estimate the extra network hop when the application uses Redis instead of a SQL cache table on the same database host. Write the hops.

#### Advanced practical tasks

1. Read a public benchmark that compares Redis and Memcached. Record the test shape, the metric, and one limit of the test.
2. Write a decision page for your team: when to add Redis, when to keep SQL-only cache rows, and when to keep Memcached.

---

## Installing Redis or using Docker

You can run Redis on your machine, in a container, or in a hosted service.

On Linux, install the Redis package from your distribution or download the official build. After the install, start the `redis-server` process. The default port is `6379`.

On Windows, the Redis project does not ship a first-class native Windows server for current versions. Use Docker Desktop, Windows Subsystem for Linux (WSL), or a hosted Redis service.

Docker is a common learning path. Pull an official image. Run a container that publishes port `6379`:

```text
docker run --name redis-learn -p 6379:6379 redis
```

The container starts `redis-server`. Stop the container when you finish the session. Start the same container again to keep data if the container still exists and you did not remove the volume.

A hosted Redis service gives a host, a port, and a password. You connect with `redis-cli` or a client library over the network. Use TLS when the service requires TLS.

Redis Stack is a distribution that adds modules such as JSON and Search. Use Redis Stack when a later topic needs those modules. Use plain Redis for topics 1 through 11 unless a task names a module.

After the instance starts, confirm that a client can connect. The next section uses `redis-cli` for that check.

Do not expose a learning instance to the public internet. Bind to `127.0.0.1` on a local machine. Set a password on any instance that is not only on localhost.

### Questions

#### Theoretical questions

1. What is the default TCP port of Redis?
2. Why does a current Windows desktop often use Docker or WSL for Redis?
3. What does `docker run ... -p 6379:6379 redis` do?
4. What three values does a hosted Redis service usually give you?
5. What is Redis Stack, and when do you need it in this learning path?

#### Easy practical tasks

1. Start Redis with Docker or a local install. Write the exact command that you used.
2. Confirm that port `6379` listens on your machine. Write the command and the result.
3. Stop the instance. Start it again. Write whether existing keys are still present (yes, no, or not tested).
4. Open a hosted Redis sign-up page. Write the connection fields that the form shows. You do not need to create an account.

#### Medium practical tasks

1. Run Redis with a named Docker volume. Restart the container. Document whether data survives.
2. Start Redis on a non-default port (for example `6380`). Connect a client to that port. Record both commands.
3. Compare a local Docker instance and a hosted instance in a four-row table: cost, persistence, network, password.

#### Advanced practical tasks

1. Build Redis from the official source on Linux or WSL. Record `redis-server --version` and the compile commands.
2. Run Redis Stack in Docker and plain Redis in Docker on two ports. Write which extra commands appear only in Redis Stack (`COMMAND LIST` or docs).

---

## `redis-cli`, `PING`, `INFO`

`redis-cli` is the official command-line client. The tool sends Redis commands and prints replies. Use `redis-cli` for learning, for health checks, and for one-off operations.

Start an interactive session against the default host and port:

```text
redis-cli
```

The prompt shows the host and the logical DB number, for example `127.0.0.1:6379>`. Type a command. Press Enter. Read the reply.

Connect to another host or port:

```text
redis-cli -h 127.0.0.1 -p 6379
```

Send one command without an interactive session:

```text
redis-cli PING
```

If the instance requires a password, pass `-a` or use `AUTH` after you connect. Prefer an environment variable or a prompt over a password in the shell history. Hosted services often require a user name and a password. Use `--user` and `--pass` or `AUTH` as the service docs show.

`PING` checks the connection. Redis replies `PONG` when the instance accepts the command. You can send a message. `PING hello` echoes `hello`. Client libraries also use `PING` as a health check.

`INFO` returns text sections about the instance. Run `INFO` for the full report. Run `INFO server`, `INFO memory`, `INFO clients`, `INFO replication`, or `INFO keyspace` for one section.

Useful lines for beginners:

- `redis_version` — the server version
- `used_memory_human` — memory in use
- `connected_clients` — open client connections
- `db0` in `INFO keyspace` — key count and expiry count for logical DB 0

The output is not JSON. Parse it as lines of `key:value`.

`SELECT` changes the logical DB for the current connection. The default is DB `0`. Keys in DB `0` do not appear in DB `1`. Modern applications use one logical DB, usually DB `0`. They split data with key prefixes, not with `SELECT`. Redis Cluster does not support `SELECT`. Do not design a new system that depends on many logical DBs.

`redis-cli` prints reply types in a readable form. A simple string, an integer, a bulk string, an array, and an error look different. An error starts with a message from the server.

Leave the interactive session with `QUIT` or `exit`, or close the terminal.

Do not run `KEYS *` on a production instance. That command can block the server on a large keyspace. This learning path uses `KEYS` only on a local empty instance. Topic 9 introduces `SCAN`.

The primary documentation is [https://redis.io/docs/](https://redis.io/docs/). The command reference is [https://redis.io/commands/](https://redis.io/commands/). When a page and this handbook disagree on a new flag, trust the command page for your Redis version. Check `INFO server` for `redis_version`.

### Questions

#### Theoretical questions

1. What is `redis-cli`?
2. What reply does `PING` return with no argument?
3. What kind of output does `INFO` return?
4. Why do modern applications use one logical DB?
5. Why is a password on the command line a risk?

#### Easy practical tasks

1. Open `redis-cli` against your instance. Run `PING`. Save the full reply.
2. Run `INFO server`. Write `redis_version` and `tcp_port`.
3. Run `INFO memory`. Write `used_memory_human`.
4. Run `redis-cli --help`. Write three flags and one short purpose for each flag.

#### Medium practical tasks

1. Connect `redis-cli` to a Docker Redis container with `docker exec` or with `-p` from the host. Document both paths if you can.
2. Run `INFO keyspace` before and after you add three keys in DB `0`. Explain the `keys` and `expires` fields.
3. Run `redis-cli --raw GET somekey` and `redis-cli GET somekey`. Write how the two outputs differ.

#### Advanced practical tasks

1. Use `redis-cli --tls` against a hosted Redis that requires TLS, or start a local TLS setup from the official docs. Record the flags.
2. Write a small shell or PowerShell loop that runs `redis-cli PING` ten times and prints the time for each call.

---

## Key-value model: keys are strings; values have a type

Redis stores pairs. Each pair has one key and one value. The client names the key in the command. Redis finds the value for that key. There is no SQL `FROM` clause and no document collection name in the core model.

A Redis key is a binary-safe string. You usually use UTF-8 text. There is a maximum key size (512 MB). Do not use large keys. Short, stable names are the rule.

A missing key is not an error for most read commands. `GET` of a missing key returns a null reply. Write commands create the key when the key does not exist, unless the command says otherwise.

One key holds one value. That value has one type. You cannot store a hash and a string under the same key at the same time. You can delete the key and create it again with a new type.

The main types in this learning path are:

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

Redis does not query by value in the basic model. You do not say "find all keys whose string equals `alice`". You design keys so that the application already knows the name, or you keep secondary indexes on purpose.

Keys live in one keyspace per logical DB. The instance can hold millions of keys when memory allows. Each key has a small overhead in RAM. Many tiny keys cost more memory than one larger structure.

Pick the type from the access pattern. A string fits a blob or a counter. A hash fits a small object with fields. A list fits a queue or a recent-item log. A set fits unique membership. A sorted set fits ranking and ranges.

### Questions

#### Theoretical questions

1. What are the two parts of a Redis pair?
2. What is the type of a Redis key?
3. What does `GET` return when the key does not exist?
4. What error does Redis return when a command does not match the type?
5. Name three commands that work on any key.

#### Easy practical tasks

1. Run `SET a 1` then `TYPE a`. Save the reply.
2. Run `TYPE missing`. Save the reply.
3. Write a two-column table: "Type" and "One good use". Add five rows from this section.
4. Run `GET nosuchkey` on an empty key name. Save the reply.

#### Medium practical tasks

1. Create a string key. Run `HGET` on that key. Record the error. Delete the key.
2. Compare a SQL row `users.id = 42` with a Redis key `user:42`. Write five sentences on lookup and on "find by email".
3. Read the `TYPE` command page. Write all type names that the page lists for your version.

#### Advanced practical tasks

1. Find the documented maximum key size. Write the number and the URL. Explain why you will never approach that limit.
2. Write a type-selection guide of one page: six access patterns and the Redis type you choose for each pattern.

---

## Single-threaded command execution (high-level)

Redis runs commands on one main thread. The instance takes one command, runs it to completion, then takes the next command. Two `INCR` commands on the same key do not interleave. You do not need a lock for a single command.

This model makes many operations simple. `INCR` is atomic. `LPUSH` is atomic. A Lua script or a Redis Function that runs on the server is also atomic with respect to other commands.

The cost is also clear. A slow command blocks other commands. `KEYS *` on a large keyspace, `LRANGE` of a huge list, or a heavy `SORT` can stall the instance. Keep commands small. Split large work.

Since Redis 6, optional I/O threads can read and write sockets. Those threads do not run the command logic. The main thread still executes commands. Do not think that Redis is a multi-core command engine because I/O threads exist.

Pipelines send many commands in one network round trip. The server still runs those commands one by one. Pipelines reduce wait time on the network. They do not run commands in parallel on many cores. Topic 6 covers pipelines.

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

## Key naming conventions

Redis has one flat keyspace per logical DB. There are no folders. You simulate folders with a convention in the key string.

A common pattern uses `:` as a separator:

```text
user:42:profile
user:42:cart
session:abc123
cache:product:1001
```

Some teams use `#` or `/`. Pick one separator and keep it. The character has no special meaning to Redis except in Cluster hash tags, which use `{` and `}`. Topic 8 covers hash tags.

A good name is stable, unique, and readable. Include the object type, the id, and the field role when you need more than one key per object. Do not put secrets in key names. Key names appear in `MONITOR`, in logs, and in support tools.

Use a product or tenant prefix in a shared instance:

```text
shopA:user:42
shopB:user:42
```

Those keys do not collide. They still share memory and `maxmemory`. Multi-tenant isolation on one instance is a convention, not a security boundary. Topic 12 covers tenant policy.

Document the scheme. Write a one-page key catalog. Avoid ad-hoc names such as `data` or `tmp`. For temporary keys, include a purpose and a unique suffix, and set a TTL.

Prefer letters, digits, `:`, `-`, and `.`. A key can contain spaces, but spaces make CLI work harder. A key such as `user:42:profile` is clear. A key that holds a full JSON document in the name is wrong.

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

1. Describe the path from an empty machine to a working `PING`. Name the install method and the client.
2. How do cache, broker, and store roles change the need for persistence?
3. Why is a key prefix a better isolation tool than `SELECT` in a new application?
4. How do the key-value model and typed values work together when you choose a command?
5. A teammate wants Memcached because "Redis is only a cache." Which facts do you use in a short reply?

#### Easy practical tasks

1. Start Redis. Run `PING`, `INFO server`, and `SET learn:start 1` then `GET learn:start`. Save all replies.
2. Write a one-page cheat sheet: install command, `redis-cli` flags, `PING`, `INFO`, `TYPE`, and three doc URLs.
3. Create three keys of different names. Run `TYPE`, `TTL`, and `EXISTS` on each. Save a small table.
4. Write ten key names for a future app. Check that each name has a type word and an id.

#### Medium practical tasks

1. Write a small script (PowerShell or bash) that checks Redis with `redis-cli PING` and exits with a non-zero status on failure.
2. Put one key in Docker Redis and one key in a second instance (other port or hosted). Document how you point `redis-cli` at each instance.
3. Document your Redis lab setup in ten steps so that another beginner can copy it (Docker or install, port, `redis-cli`, first `PING`).

#### Advanced practical tasks

1. Run two Redis containers on ports `6379` and `6380`. Write a table of `INFO server` fields that differ (`tcp_port`, `process_id`, `run_id`).
2. Connect a client library in a language that you know. Send `PING` from code and from `redis-cli`. Record both replies and the library name.
