# 6. Pipelines, Transactions, and Scripts

## Description

A Redis client sends commands and reads replies. This topic explains one command at a time versus a pipeline, Redis transactions (`MULTI`, `EXEC`, `WATCH`), and server-side scripts. You learn what "atomic" means in Redis and how that differs from SQL rollback.

Complete this topic before you study persistence and replication.

Use one term for each concept. A round trip is one network send and one wait for a reply. A pipeline is many commands in one send, with replies later. A transaction is a queued block that `EXEC` runs without interleaving other commands. A script is Lua or a Redis Function that runs on the server.

---

## Pipeline vs one command at a time

The default mode is request/response. The client sends `GET k`. The client waits. Redis replies. The client sends the next command. Each command pays a network delay (RTT).

A pipeline sends many commands without a wait after each command. The client then reads all replies in order. Redis still runs the commands one by one on the main thread. The gain is less time spent waiting on the network.

```text
# conceptual, not a special Redis command
INCR views
INCR views
INCR views
```

In a pipeline, those three `INCR`s can sit in one TCP write. The replies are `1`, `2`, `3` if the key started missing.

`redis-cli` can pipeline from a file:

```text
redis-cli --pipe < commands.txt
```

Client libraries have a pipeline API. Use it for many independent commands (thousands of `SET`s, many `HGET`s).

A pipeline is not a transaction. Other clients can run commands between your pipelined commands. There is no rollback. If the connection drops, you may not know which commands ran.

Keep pipeline batches bounded (for example hundreds or a few thousand commands). A huge pipeline uses client and server buffers.

Do not pipeline commands that must see the previous reply (unless you accept that you cannot branch in the middle). For "read, then write if", use `WATCH`, a Lua script, or a Function.

On a 0.5 ms RTT, 10,000 `GET`s wait about 5 seconds for the network alone. A pipeline of 10,000 `GET`s pays a few RTTs for the whole batch plus server time.

`MGET` and `MSET` reduce RTTs for strings on one slot. They are not a full substitute for a pipeline of mixed commands.

Cluster clients split batches by slot (topic 8). A pipeline of huge `GET`s can still stall the main thread and fill buffers.

### Questions

#### Theoretical questions

1. What does the client wait for in request/response mode?
2. What does a pipeline reduce?
3. Does a pipeline run commands in parallel on many cores?
4. Can another client run commands in the middle of your pipeline?
5. Why must a pipeline batch have a size limit?

#### Easy practical tasks

1. Write four sentences: request/response vs pipeline.
2. Draw a timeline of three `PING`s with waits, then three `PING`s in a pipeline.
3. Run `redis-cli --help` and find `--pipe`. Write the flag line.
4. Open your language's Redis client docs. Write the pipeline type or method name.

#### Medium practical tasks

1. Time 1,000 `INCR`s in a loop (wait each time) vs 1,000 `INCR`s in a pipeline. Write both times.
2. Pipeline `SET` of 100 keys and then `MGET` them. Confirm values.
3. Break a pipeline in the middle of a script (close the client). Write what you cannot prove without a read.

#### Advanced practical tasks

1. Find a recommended batch size in official or client docs. Test 100 vs 10,000 commands. Write buffer or time effects.
2. Pipeline commands to two different keys on Redis Cluster (topic 8 preview). Record whether the client had to split by slot.

---

## `MULTI` / `EXEC` / `WATCH`

`MULTI` starts a transaction. Redis queues the next commands on that connection. Redis does not run them yet. The immediate reply is `QUEUED`.

`EXEC` runs the queued commands in order as one unit on the main thread. No other connection runs a command in the middle of that `EXEC`. The reply is an array of replies, one per queued command.

`DISCARD` drops the queue and ends the transaction without `EXEC`.

```text
MULTI
INCR acct:1
INCR acct:2
EXEC
```

If a queued command has a syntax error, Redis can reject `EXEC` and run nothing (depends on when Redis detects the error). If a command fails at execution (`WRONGTYPE`, missing key for a type), later commands in the same `EXEC` still run. Redis does not roll back the earlier commands.

`WATCH key [key ...]` marks keys on the current connection. If any watched key changes before `EXEC`, Redis aborts the transaction. `EXEC` returns a null reply. The client then retries: `WATCH` again, read, `MULTI`, write, `EXEC`.

This is optimistic locking. You assume conflict is rare. You detect conflict at `EXEC`.

```text
WATCH acct:1
GET acct:1
MULTI
SET acct:1 40
EXEC
```

If another client writes `acct:1` after `WATCH` and before `EXEC`, `EXEC` fails. Your `SET` does not run.

`UNWATCH` clears watches. `EXEC` and `DISCARD` also clear watches.

`WATCH` is per connection. You must `GET` on the same connection after `WATCH` and use that connection for `MULTI` / `EXEC`. A connection pool must not mix these steps across connections.

`WATCH` does not lock other clients. They can still write. Your `EXEC` fails instead.

Do not `WATCH` a huge set of keys. Do not use `WATCH` as a distributed lock. Use `SET NX EX` for a lock preview, or a script for multi-key updates.

`WATCH` plus `MULTI` is enough for "read a balance, write if still the same" on one instance. On Cluster, all keys in the transaction must live in the same slot (topic 8).

Use `MULTI` / `EXEC` when you need several commands to run without interleaving, and you do not need to read a value in the middle. You cannot `GET` inside `MULTI` and use that value in the next queued command on the client in a useful way. The `GET` is queued. The client does not see the value until `EXEC`.

For logic that reads and then writes, use `WATCH` or Lua.

Transactions do not have isolation levels like SQL. They do not have savepoints. They do not wait for disk by themselves.

### Questions

#### Theoretical questions

1. What reply do you see after a command that follows `MULTI`?
2. What does `EXEC` return on success?
3. What does `WATCH` cause `EXEC` to do when a key changes?
4. Does `WATCH` stop other clients from writing?
5. Why can you not branch on a `GET` result inside `MULTI` on the client?

#### Easy practical tasks

1. Run `MULTI`, `SET lab:t 1`, `INCR lab:t`, `EXEC`. Save the `EXEC` array.
2. Run `MULTI`, `SET lab:t 2`, `DISCARD`. `GET lab:t`. Save the reply.
3. `WATCH lab:w`, `GET lab:w`, `MULTI`, `SET lab:w 1`, `EXEC` with no other writer. Save `EXEC`.
4. Write the retry loop in five numbered steps (no code required).

#### Medium practical tasks

1. From a second client, `GET` a key while the first client is inside `MULTI` but before `EXEC`. Then `EXEC`. Write what each client saw.
2. Two terminals: client A `WATCH`s and waits. Client B `SET`s the key. Client A `EXEC`s. Record the null `EXEC`.
3. Queue `INCR` on a hash key (create a hash first). `EXEC`. Show which replies are errors and whether other queued `SET`s applied.

#### Advanced practical tasks

1. Read the official transactions page. Write the exact rules for errors before `EXEC` vs errors inside `EXEC`.
2. Write a client program that retries a `WATCH` increment up to 5 times under parallel workers. Count successes and retries.

---

## Why Redis transactions are not SQL rollback

SQL transactions can `ROLLBACK`. The engine undoes writes when you abort, or when a constraint fails (depending on the engine and the statement).

A Redis `EXEC` does not undo. If command 1 in the transaction succeeds and command 2 fails, command 1 stays applied. Redis is not a relational engine. There are no row locks, no WAL rollback of the transaction, and no `ROLLBACK` command that restores values.

Redis transactions give:

- no interleaving of other commands during `EXEC`
- optional abort of the whole `EXEC` when `WATCH` detects a change
- queued execution on one connection

Redis transactions do not give:

- rollback of a failed command inside `EXEC`
- interactive read-your-writes inside the queue
- multi-instance atomic commit (no two-phase commit in core Redis)

If you need "all keys update or none" with logic, put the logic in a Lua script or a Function. The script runs as one command. If you `redis.call` a write and then hit an error, writes that already ran in that script stay applied unless you check errors before you write. Write scripts so that checks happen first.

Application-level compensation (write the inverse) is possible but easy to get wrong. Prefer one atomic script for money-like updates on one Redis instance.

Do not name Redis `MULTI` a "SQL transaction" in a design review. Say "queued atomic batch" or "Redis transaction".

### Questions

#### Theoretical questions

1. What does SQL `ROLLBACK` do that Redis `EXEC` does not do?
2. What happens if the second command in `EXEC` returns `WRONGTYPE`?
3. What does `WATCH` abort, and what does it not restore?
4. How does a Lua script improve "all or nothing" logic compared with `MULTI`?
5. What phrase should you use instead of "SQL transaction" for `MULTI`/`EXEC`?

#### Easy practical tasks

1. Make a two-column table: SQL transaction vs Redis `EXEC`. Add four rows.
2. Write four sentences for a teammate who expects rollback.
3. List three Redis guarantees during `EXEC`.
4. Find `DISCARD` on the command page. Write that it is not `ROLLBACK` of applied writes (because writes are not applied yet).

#### Medium practical tasks

1. Queue `SET lab:ok 1` and `INCR lab:hashfield` on a hash key. `EXEC`. `GET lab:ok`. Prove the first write stayed.
2. Write a compensation plan for a failed two-key update (inverse `INCR`). List two ways compensation can fail.
3. Read the official "Transactions" page section on errors. Quote the rule in your own words (two sentences).

#### Advanced practical tasks

1. Write a short design note: payment hold in Redis vs payment hold in PostgreSQL. Include rollback and crash.
2. Trace a failed `EXEC` after `WATCH` in a log: no keys changed vs keys changed. Write how you tell the two cases.

---

## Lua: `EVAL`, `EVALSHA`

`EVAL script numkeys key [key ...] arg [arg ...]` runs a Lua script on the server. The script uses `KEYS[]` and `ARGV[]`. Redis runs the script as one command. Other commands wait.

```text
EVAL "return redis.call('INCR', KEYS[1])" 1 lab:c
```

`redis.call` raises an error on a Redis error. `redis.pcall` returns an error value.

`EVALSHA sha1 numkeys ...` runs a script that the server already cached. `SCRIPT LOAD` returns the SHA1. `EVALSHA` saves bandwidth. If the server restarted and the cache is empty, the client receives `NOSCRIPT` and must `EVAL` or `SCRIPT LOAD` again.

Use a script when you must read and write several keys with logic and no interleaving. Example: increment a counter and set `EXPIRE` only when the new value is `1`. Example: release a lock if `GET` equals your token.

Rules:

- Pass every key that you access in the `KEYS` list. Do not build key names from `ARGV` in Cluster-safe scripts.
- Keep scripts short. A long script blocks the instance.
- Do not use Lua as a general application server. No long loops. No random sleeps.
- Scripts must be deterministic for replication. Avoid `TIME` in old modes unless you follow the documented pattern. Read the scripting docs for your version.

Redis replicates the script or its effects to replicas. A bug in a script can replicate too.

Prefer `EVALSHA` plus a load step in application startup.

### Questions

#### Theoretical questions

1. What does `numkeys` in `EVAL` mean?
2. Why is a Lua script atomic with respect to other commands?
3. When do you use `EVALSHA` instead of `EVAL`?
4. What error means the script cache does not have your SHA1?
5. Why must Cluster-safe scripts list all keys in `KEYS`?

#### Easy practical tasks

1. `EVAL` a script that `GET`s `KEYS[1]`. Pass one key. Save the reply.
2. `EVAL` `INCR` plus `EXPIRE` when the value is `1`. Test twice. Save TTLs.
3. `SCRIPT LOAD` a one-line script. Run `EVALSHA`. Save the SHA1.
4. Write four sentences: Lua vs `MULTI` for read-then-write.

#### Medium practical tasks

1. Write a lock-release script: `GET` then `DEL` if `ARGV[1]` matches. Test match and mismatch.
2. Call `EVALSHA` after `SCRIPT FLUSH` on a lab instance. Record `NOSCRIPT`. Reload.
3. Compare `redis.call` and `redis.pcall` on `INCR` of a hash key. Write the two behaviors.

#### Advanced practical tasks

1. Implement a two-key transfer that refuses a negative balance. Use `KEYS` for both accounts. Test success and reject.
2. Read replication rules for scripts on your Redis version. Write whether Redis replicates the script body or the resulting commands.

---

## Functions (Redis 7+)

Redis 7 adds Functions. A Function is a registered library that lives in the server. Clients call `FCALL name numkeys key [key ...] arg ...`.

You load a library with `FUNCTION LOAD` and a payload that the docs specify (Lua library syntax). `FUNCTION LIST` shows libraries. `FUNCTION DELETE` removes a library. `FUNCTION FLUSH` removes all (lab only).

Functions solve operational problems of raw `EVAL` strings in every client:

- The server stores the code
- You can update a library in one place
- `FCALL` uses a stable name

Functions still run on the main thread. They are still atomic. They still block if they are slow. Cluster key rules still apply.

`EVAL` still works. Many applications keep Lua scripts. New designs on Redis 7+ can prefer Functions when the team wants server-side versioning.

You must persist functions like other data if you need them after a restart. Check whether your persistence mode and replication restore functions. The official Functions page describes persistence and replication.

Do not mix ten copies of the same Lua in clients and a Function library without a single source of truth.

AOF and RDB can include function definitions so that a replica or a restart still has `FCALL` targets. Test this in a lab before production.

### Questions

#### Theoretical questions

1. Which Redis version introduced Functions?
2. What command invokes a Function?
3. What operational problem do Functions solve compared with `EVAL` in every client?
4. Do Functions run in parallel with other commands on the same instance?
5. Why must you test Functions across restart?

#### Easy practical tasks

1. Run `INFO server` and write whether the version is 7 or newer.
2. Open the Functions introduction page. Write `FUNCTION LOAD` and `FCALL` in one sentence each.
3. Run `FUNCTION LIST` on your instance. Save the reply (empty is fine).
4. Write four sentences: Function vs `EVALSHA`.

#### Medium practical tasks

1. If you have Redis 7+, load the smallest official example library. `FCALL` it. Save the reply. If you do not, write the upgrade or Docker tag you need.
2. `FUNCTION LIST` after load. Then `FUNCTION DELETE` (lab). Confirm `FCALL` fails.
3. Compare where the code lives: application repo (`EVAL`) vs Redis (`FUNCTION LOAD`). Write a team policy of five lines.

#### Advanced practical tasks

1. Port your lock-release Lua to a Function library. Call it with `FCALL`. Test match and mismatch.
2. Restart Redis with RDB or AOF. Check `FUNCTION LIST`. Document whether the library survived and which persistence flags you used.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Rank pipeline, `MULTI`/`EXEC`, `WATCH`, and Lua by "no other command in the middle" and by "can use a read result in later logic".
2. What fails across connection pools if you use `WATCH`?
3. How is Redis atomicity different from SQL rollback?
4. When do you choose `EVALSHA` or `FCALL` instead of a client-side pipeline?
5. What shared risk do long pipelines, long transactions, and long scripts have on a single-threaded Redis?

#### Easy practical tasks

1. Run one pipeline of `PING`s, one `MULTI`/`EXEC` of `PING`s, and one `EVAL` that returns `PONG`. Save the shapes of the three replies.
2. Write a cheat sheet: `--pipe`, `MULTI`, `EXEC`, `DISCARD`, `WATCH`, `EVAL`, `EVALSHA`, `FCALL`.
3. Draw a flowchart: independent writes → pipeline; all-or-nothing batch without reads → `MULTI`; read-then-write → `WATCH` or Lua.
4. List three commands you will never put in a hot script (`KEYS`, huge `LRANGE`, `DEBUG SLEEP`).

#### Medium practical tasks

1. Write a script that pipelines 5,000 `SET`s, then runs a Lua `INCR`+`EXPIRE` on one rate-limit key. Print times.
2. Document a connection-pool rule for transactions and scripts (one connection, timeout, retry).
3. Implement increment-if-less-than-N three ways (`WATCH`, Lua, Function if available). Write a table of races and retries.

#### Advanced practical tasks

1. Under parallel clients, compare lost updates for `GET`+`SET` vs `INCR` vs Lua compare-and-set. Report final values.
2. Read `SCRIPT KILL` and `FUNCTION KILL` (busy script). Write when you can kill a script and when you cannot (writes already started).
