# 1. Getting Started

## Description

Redis is an in-memory data store. Programs use Redis as a cache, as a message broker, and as a database. This topic shows what Redis is, how Redis compares with two common alternatives, and how you start a Redis instance. You also learn `redis-cli` and the first commands.

Complete this topic before you study keys, types, and persistence.

Use one term for each concept. An instance is one running Redis process. A client is a program that sends commands to that instance. A command is one request with a name and arguments. Do not mix the words "database" and "logical DB" without care. A logical DB is a numbered namespace inside one instance.

---

## What Redis is (in-memory data store, used as cache, broker, and database)

Redis keeps data in RAM. A client sends a command. Redis applies the command to keys in memory. Redis then sends a reply. The path is short. Latency is often less than one millisecond on a local network.

Redis is a data store. You write a key. You read that key. Redis also supports more than a plain string value. Later topics cover hashes, lists, sets, sorted sets, streams, and other types.

Teams use Redis in three common roles:

- Cache. The application writes a result that is expensive to compute. The next read uses Redis. A TTL can remove the key after a time.
- Broker. The application publishes messages or writes stream entries. Other processes read those messages.
- Database. The application stores state in Redis and treats Redis as the source of truth. Persistence and replication become important.

Redis is not a replacement for every database. Redis does not replace a relational engine for complex joins and multi-row constraints. Redis is not a full message platform with long retention when you compare it with a log system such as Apache Kafka.

Redis is fast because data lives in memory and because command execution uses one main thread. Persistence to disk is optional. You can run Redis with no files on disk. You can also enable snapshots or an append-only file. Later topics cover those modes.

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
2. Make a two-column table: "Role" and "What Redis stores". Add one row for cache, one row for broker, and one row for database.
3. List three application features that fit Redis. List two features that do not fit well. Give one reason for each choice.
4. Open [https://redis.io/docs/](https://redis.io/docs/). Write the sentence that the site uses to describe Redis.

#### Medium practical tasks

1. Compare Redis with one database that you already know. Write six short sentences. Cover memory, types, and durability.
2. Draw a simple diagram: application, Redis instance, and disk. Label the path of a `GET` when data is only in RAM.
3. Find the current stable Redis version on redis.io. Write the version number and the date of the release notes page.

#### Advanced practical tasks

1. Read the Redis introduction pages on redis.io. Write a one-page timeline with years and one product fact per year.
2. Interview a teammate or read a public post-mortem. Write how that team used Redis (cache, broker, or database) and one risk they named.

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

## Installing Redis or using Docker / Redis Cloud

You can run Redis on your machine, in a container, or in a hosted service.

On Linux, install the Redis package from your distribution or download the official build. After the install, start the `redis-server` process. The default port is `6379`.

On Windows, the Redis project does not ship a first-class native Windows server for current versions. Use Docker Desktop, Windows Subsystem for Linux (WSL), or a hosted Redis service.

Docker is a common learning path. Pull an official image. Run a container that publishes port `6379`:

```text
docker run --name redis-learn -p 6379:6379 redis
```

The container starts `redis-server`. Stop the container when you finish the session. Start the same container again to keep data if the container still exists and you did not remove the volume.

Redis Cloud is a hosted service. You create a database in the browser. The service gives a host, a port, and a password. You connect with `redis-cli` or a client library over the network. Use TLS when the service requires TLS.

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
4. Open a Redis Cloud or other hosted sign-up page. Write the connection fields that the form shows. You do not need to create an account.

#### Medium practical tasks

1. Run Redis with a named Docker volume. Restart the container. Document whether data survives.
2. Start Redis on a non-default port (for example `6380`). Connect a client to that port. Record both commands.
3. Compare a local Docker instance and a hosted instance in a four-row table: cost, persistence, network, password.

#### Advanced practical tasks

1. Build Redis from the official source on Linux or WSL. Record `redis-server --version` and the compile commands.
2. Run Redis Stack in Docker and plain Redis in Docker on two ports. Write which extra commands appear only in Redis Stack (`COMMAND LIST` or docs).

---

## `redis-cli`

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

`redis-cli` prints reply types in a readable form. A simple string, an integer, a bulk string, an array, and an error look different. Learn to read those shapes. An error starts with a message from the server.

Leave the interactive session with `QUIT` or `exit`, or close the terminal.

Do not run `KEYS *` on a production instance. That command can block the server on a large keyspace. This learning path uses `KEYS` only on a local empty instance. Later topics introduce `SCAN`.

### Questions

#### Theoretical questions

1. What is `redis-cli`?
2. Which flags set the host and the port?
3. How do you send one command and then exit without an interactive prompt?
4. Why is a password on the command line a risk?
5. What does the interactive prompt `127.0.0.1:6379>` tell you?

#### Easy practical tasks

1. Open `redis-cli` against your instance. Run `PING`. Save the full reply.
2. Run `redis-cli -h 127.0.0.1 -p 6379 PING` as a one-shot command. Save the output.
3. Run `redis-cli --help`. Write three flags and one short purpose for each flag.
4. In an interactive session, type a command name that does not exist. Record the error text.

#### Medium practical tasks

1. Connect `redis-cli` to a Docker Redis container with `docker exec -it redis-learn redis-cli` or with `-p` from the host. Document both paths if you can.
2. Enable a password with `CONFIG SET requirepass learnpass` on a local instance. Connect with `AUTH`. Then remove the password. Use only a local instance.
3. Run `redis-cli --raw GET somekey` and `redis-cli GET somekey`. Write how the two outputs differ.

#### Advanced practical tasks

1. Use `redis-cli --tls` against a hosted Redis that requires TLS, or start a local TLS setup from the official docs. Record the flags.
2. Write a small shell or PowerShell loop that runs `redis-cli PING` ten times and prints the time for each call.

---

## `PING`, `INFO`, `SELECT` (logical DBs: use one DB in modern apps)

`PING` checks the connection. Redis replies `PONG` when the instance accepts the command. You can send a message:

```text
PING
PING hello
```

The second form echoes `hello`. Use `PING` after you connect. Client libraries also use `PING` as a health check.

`INFO` returns text sections about the instance. Run `INFO` for the full report. Run `INFO server`, `INFO memory`, `INFO clients`, `INFO replication`, or `INFO keyspace` for one section.

Useful lines for beginners:

- `redis_version` — the server version
- `used_memory_human` — memory in use
- `connected_clients` — open client connections
- `db0` in `INFO keyspace` — key count and expiry count for logical DB 0

The output is not JSON. Parse it as lines of `key:value`.

`SELECT` changes the logical DB for the current connection. The default is DB `0`. A default install has 16 logical DBs (`0` through `15`). Keys in DB `0` do not appear in DB `1`.

```text
SELECT 0
SET learn:flag 1
SELECT 1
GET learn:flag
```

The `GET` in DB `1` does not see the key from DB `0`.

Modern applications use one logical DB, usually DB `0`. They split data with key prefixes, not with `SELECT`. Redis Cluster does not support `SELECT`. Many hosted services expose only DB `0`. Do not design a new system that depends on many logical DBs.

`SELECT` is still useful in a local lab when you want a clean namespace. Reset the connection to DB `0` when you finish.

### Questions

#### Theoretical questions

1. What reply does `PING` return with no argument?
2. What kind of output does `INFO` return?
3. What does `SELECT 1` change?
4. Why do modern applications use one logical DB?
5. Does Redis Cluster support `SELECT`? What does that fact mean for key design?

#### Easy practical tasks

1. Run `PING` and `PING learner`. Save both replies.
2. Run `INFO server`. Write `redis_version` and `tcp_port`.
3. Run `INFO memory`. Write `used_memory_human`.
4. Run `SELECT 1`, then `PING`, then `SELECT 0`. Confirm that both selects succeed on a standalone instance.

#### Medium practical tasks

1. Put a key in DB `0` and a different value on the same key name in DB `1`. Read both with `SELECT`. Record the two values.
2. Run `INFO keyspace` before and after you add three keys in DB `0`. Explain the `keys` and `expires` fields.
3. Try `SELECT 1` against a Redis Cluster or a service that forbids it, or read the error in the docs. Write the error or the documented limit.

#### Advanced practical tasks

1. Change `databases` in a local `redis.conf` to `4`, restart, and test `SELECT 3` and `SELECT 4`. Record which command fails.
2. Write a one-page note for a team: "We use only DB 0." Include Cluster, hosted Redis, and prefix design.

---

## Official docs: redis.io/docs

The primary documentation is [https://redis.io/docs/](https://redis.io/docs/). The site covers data types, commands, persistence, replication, and Cluster.

The command reference is [https://redis.io/commands/](https://redis.io/commands/). Each command page shows the syntax, the time complexity, the return type, and examples. Open the command page when you learn a new command.

Data type pages start at [https://redis.io/docs/data-types/](https://redis.io/docs/data-types/). Read the type page before you memorize every command.

Redis University is [https://university.redis.io/](https://university.redis.io/). The courses are optional extras. This handbook is the path for the topics in `redis.topics.md`.

Use `HELP` in `redis-cli` for a short in-server hint when the server supports it. The web command page is the complete reference.

When a page and this handbook disagree on a new flag, trust the command page for your Redis version. Check `INFO server` for `redis_version`.

Bookmark the docs home, the command index, and the data type index. You will open these pages in every later topic.

### Questions

#### Theoretical questions

1. What site is the primary Redis documentation?
2. What does one command page on redis.io usually include?
3. Where do you read an overview of Redis data types?
4. Why must you check `redis_version` when you read a command page?
5. What is Redis University in relation to this handbook?

#### Easy practical tasks

1. Open the `PING` command page. Write the time complexity and the return type.
2. Open the data types index. List five type names that the page shows.
3. Bookmark redis.io/docs, redis.io/commands, and redis.io/docs/data-types.
4. Use the command search. Find `INFO`. Write the section list that the page names.

#### Medium practical tasks

1. Compare the `SET` command page for two Redis versions (or the "history" notes on the page). Write one option that a new version added.
2. Open the persistence docs. Write the names of the two main persistence modes in one sentence each.
3. Find the page that explains Redis vs other databases or "What is Redis". Write three claims and the URL.

#### Advanced practical tasks

1. Download or open the command.json / command reference data if the site provides it. Count how many commands start with `H`. Write the method.
2. Make a personal index: topic number from this path, official URL, and one sentence. Cover topics 1 through 5.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe the path from an empty machine to a working `PING`. Name the install method and the client.
2. How do cache, broker, and database roles change the need for persistence?
3. Why is a key prefix a better isolation tool than `SELECT` in a new application?
4. What is the difference between `redis-cli` interactive mode and a one-shot `redis-cli` command?
5. A teammate wants Memcached because "Redis is only a cache." Which facts do you use in a short reply?

#### Easy practical tasks

1. Start Redis. Run `PING`, `INFO server`, and `SET learn:start 1` then `GET learn:start`. Save all replies.
2. Write a one-page cheat sheet: install command, `redis-cli` flags, `PING`, `INFO`, `SELECT`, and three doc URLs.
3. Export `INFO` to a text file. Highlight `redis_version`, `tcp_port`, `used_memory_human`, and `connected_clients`.
4. Create a folder note that records host, port, and "DB 0 only" as your lab rules.

#### Medium practical tasks

1. Write a small script (PowerShell or bash) that starts nothing but checks Redis with `redis-cli PING` and exits with a non-zero status on failure.
2. Put one key in Docker Redis and one key in a second instance (other port or hosted). Document how you point `redis-cli` at each instance.
3. Document your Redis lab setup in ten steps so that another beginner can copy it (Docker or install, port, `redis-cli`, first `PING`).

#### Advanced practical tasks

1. Run two Redis containers on ports `6379` and `6380`. Write a table of `INFO server` fields that differ (`tcp_port`, `process_id`, `run_id`).
2. Connect a client library in a language that you know. Send `PING` from code and from `redis-cli`. Record both replies and the library name.
