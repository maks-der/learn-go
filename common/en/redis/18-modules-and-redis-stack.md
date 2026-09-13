# 18. Modules and Redis Stack

## Description

A module adds commands and types that core Redis does not ship. Redis Stack is a distribution that loads a common set of modules. This topic surveys RedisJSON, RediSearch, RedisTimeSeries, and RedisBloom, then covers when a module is the better path and what operations cost you pay.

Complete topic 7 (module survey) first. That topic introduced command prefixes. This topic is the choice and operations view.

Use one term for each concept. A module is a server plugin that registers commands. Redis Stack is Redis plus a bundle of modules. A schema is the index definition for Search. A filter is a probabilistic structure (Bloom and related). Compatibility is whether your Redis version, Cluster, persistence, and client can run the module.

---

## RedisJSON

RedisJSON stores JSON documents as a first-class value. You address paths. You do not rewrite the whole document for a small change.

```text
JSON.SET lab:user:1 $ '{"name":"Ada","n":1}'
JSON.GET lab:user:1 $.name
JSON.NUMINCRBY lab:user:1 $.n 1
```

`$` is the document root in JSONPath. Path syntax is on the RedisJSON command pages. Check your module version.

Why not `SET` a JSON string? A string `SET` replaces the whole blob. Concurrent updates of two fields need a lock or a Lua read-modify-write. RedisJSON can update one path. Memory layout is for documents, not for a raw string.

`TYPE` on a JSON key is not `string`. Commands such as `GET` do not apply. Clients need RedisJSON support or raw command calls.

Search can index JSON documents (next section) when you run both modules.

Cluster: the JSON key is still one key in one slot. A huge document is a big key (topic 16).

Persistence and replication include module data when the module supports it. Test restart (topic 9).

Use RedisJSON when documents are nested, updates are partial, and you already run Redis Stack or a hosted JSON feature. Use a hash when the model is a flat map of small fields. Use a string of JSON when the document is small and you always replace it.

### Questions

#### Theoretical questions

1. What problem does RedisJSON solve that a JSON string does not solve well?
2. What does `$` mean in `JSON.SET`?
3. Why can concurrent field updates be easier with RedisJSON than with `GET`/`SET` of a string?
4. Does `GET` work on a JSON key?
5. When is a hash enough instead of RedisJSON?

#### Easy practical tasks

1. If Stack is available, run the three commands in this section. Save the replies. If not, write the Docker image name `redis/redis-stack` and a `docker run` line.
2. Open the `JSON.SET` page. Write the time complexity note.
3. Write four sentences: RedisJSON vs `SET` of JSON text.
4. Write one JSONPath for `address.city` from a guessed document.

#### Medium practical tasks

1. `JSON.SET` a document with an array. `JSON.ARRAPPEND` one element. `JSON.GET` the array path.
2. Compare `MEMORY USAGE` of the same data as a hash, as a JSON string, and as RedisJSON (if available). Write three numbers.
3. Try `GET` on a JSON key. Record the error.

#### Advanced practical tasks

1. Design a user document with RedisJSON plus a cache-aside TTL. Write which paths you update in place vs replace.
2. Test AOF or RDB restart with one JSON key. Document whether the key survives.

---

## RediSearch

RediSearch (Search) builds an index on hashes or JSON. You query the index with `FT.SEARCH`. This is not SQL. You must create a schema.

```text
FT.CREATE idx:lab ON HASH PREFIX 1 lab:doc: SCHEMA title TEXT body TEXT score NUMERIC
FT.SEARCH idx:lab "hello"
```

`FT.CREATE` defines fields and types (`TEXT`, `TAG`, `NUMERIC`, `GEO`, and others by version). Documents are ordinary Redis keys that match the prefix. The module maintains inverted indexes.

You must keep the index in sync. Writes to hashes (or JSON) that match the prefix update the index. If you bypass the expected type, the document may not index.

`FT.DROPINDEX` removes an index. Read the `DD` option before you drop. You can delete the index only, or also the documents.

Search uses RAM for the index. Large text corpora cost memory beyond the documents.

Aggregations (`FT.AGGREGATE`) exist. They are not a warehouse. Heavy analytics still belong in a system that is built for scans.

Cluster and Search have specific deployment rules. Read the module Cluster notes. Some setups use a dedicated Search shard layout.

Use Search when you need full-text or filter queries that hashes and `SCAN` cannot do. Do not use Search as a substitute for PostgreSQL when you need joins, constraints, and ad-hoc SQL.

### Questions

#### Theoretical questions

1. What must you create before `FT.SEARCH` works?
2. What Redis types can Search index in this section?
3. Why does Search use extra RAM beyond the documents?
4. Is `FT.SEARCH` SQL?
5. What operational command removes an index?

#### Easy practical tasks

1. Open the `FT.CREATE` page. Write three schema field types.
2. Write four sentences: Search vs `SCAN MATCH`.
3. If Stack is available, run `FT._LIST` or `FT.CREATE` on a lab index. Save the reply. If not, write the missing module error.
4. Write a prefix and two fields for a product index.

#### Medium practical tasks

1. Create an index on three hashes. `FT.SEARCH` a word that appears in one title. Save the reply. Drop the index in the lab.
2. Read `FT.DROPINDEX` options. Write the difference when documents are deleted vs kept.
3. Estimate why an index on 1 million short titles needs a memory test. Write the test plan.

#### Advanced practical tasks

1. Index JSON documents with a JSON path schema (if your version supports it). Search one field. Save the query.
2. Write a decision: RediSearch vs an external search engine vs SQL `LIKE`. Include ops team and freshness.

---

## RedisTimeSeries

RedisTimeSeries stores timestamped samples. Each key is one series (one metric, one labeled stream of points).

```text
TS.CREATE lab:cpu
TS.ADD lab:cpu * 0.42
TS.RANGE lab:cpu - +
```

`*` uses the server time. You can pass an explicit timestamp.

`TS.CREATERULE` downsamples into another series (avg, min, max, and others). Retention (`RETENTION`) deletes old samples. Compaction keeps coarse history.

Why not a sorted set? A sorted set can store `score=timestamp` and `member=value`. At large volume, TimeSeries is more compact and has aggregation rules. Topic 6 windows are fine for small cases.

Labels (`TS.ADD` with `LABELS`) let you query many series (`TS.MRANGE`). Cardinality of labels must stay bounded. High-cardinality labels (one series per user id at huge scale) explode memory.

Use TimeSeries for infrastructure or product metrics that you already want in Redis. Use a dedicated metrics system (Prometheus, a warehouse) when retention and query volume exceed a memory store.

Do not write one TimeSeries sample per user click at internet scale without a plan.

### Questions

#### Theoretical questions

1. What does one TimeSeries key represent?
2. What does `TS.CREATERULE` do?
3. Why can TimeSeries beat a sorted set at large volume?
4. What is the risk of high-cardinality labels?
5. When do you choose Prometheus (or similar) instead?

#### Easy practical tasks

1. If Stack is available, `TS.ADD` two samples. `TS.RANGE`. Save the replies.
2. Write four sentences: TimeSeries vs sorted set window (topic 6).
3. Open `TS.CREATE`. Write the `RETENTION` meaning.
4. Write a key name for CPU on host `web1`.

#### Medium practical tasks

1. Create a rule that downsamples into `lab:cpu:1m`. Add samples. `TS.RANGE` the destination after a wait or with explicit times.
2. Compare `MEMORY USAGE` of 1000 `TS.ADD` samples vs 1000 `ZADD` members (same numbers). Write both.
3. List three labels you would allow and one label you would forbid (high cardinality).

#### Advanced practical tasks

1. Design retention: raw 2 hours, 1-minute avg for 7 days. Write keys and rules.
2. Write when metrics stay in Redis vs when they ship to a dedicated TSDB. Include RAM math.

---

## RedisBloom

RedisBloom adds probabilistic structures: Bloom filters, Cuckoo filters, Count-Min sketches, Top-K, and related types. Exact names depend on the module version.

A Bloom filter answers "definitely not in the set" or "possibly in the set." It does not store the items. False positives exist. False negatives do not (in the standard Bloom filter).

```text
BF.RESERVE lab:bf 0.01 10000
BF.ADD lab:bf user:42
BF.EXISTS lab:bf user:42
```

`0.01` is a target error rate. `10000` is a capacity hint. After too many inserts, the error rate rises unless you reserved enough capacity.

Use a Bloom filter to skip a slower check (SQL lookup, disk) when the answer is "definitely missing." Topic 13 negative caching stores a key. A Bloom filter uses less memory for huge id spaces, with error.

HyperLogLog (topic 7) estimates unique counts. A Bloom filter tests membership. Do not mix the two.

A Count-Min sketch estimates frequencies. Top-K estimates heavy hitters. Read the command pages before you treat the numbers as exact.

Do not use a Bloom filter as an access-control list. A false positive can mean "possibly allowed." That is the wrong tool for security.

### Questions

#### Theoretical questions

1. What two answers can a Bloom filter give?
2. Does a standard Bloom filter have false negatives?
3. What happens if you insert far more items than the reserved capacity?
4. How does a Bloom filter differ from HyperLogLog?
5. Why is a Bloom filter a poor ACL?

#### Easy practical tasks

1. If the module is present, run `BF.RESERVE`, `BF.ADD`, `BF.EXISTS`. Save replies. If not, write the command names from the docs.
2. Write four sentences: Bloom filter vs a Redis set.
3. Write four sentences: Bloom filter vs negative cache keys.
4. Open `BF.RESERVE`. Write the two numeric arguments.

#### Medium practical tasks

1. Insert 1000 items. Test 100 missing items. Count how many `BF.EXISTS` return 1 (false positives). Write the count.
2. Compare `MEMORY USAGE` of a Bloom filter vs a set of the same 1000 ids.
3. Read Count-Min or Top-K in the docs. Write one command prefix and one use.

#### Advanced practical tasks

1. Design a "possibly seen URL" filter in front of SQL. Write the error-rate choice and the rebuild plan when you exceed capacity.
2. Write a table: set, HLL, Bloom, Count-Min. Rows: exact?, membership?, count?, memory, error type.

---

## When a module is better than building it yourself

Build it yourself (core types + application) when:

- A hash, stream, or sorted set already matches the access pattern
- The extra dependency is not worth the ops cost
- You must run vanilla Redis in a constrained platform
- The logic is small (one Lua script)

Use a module when:

- The module implements a known algorithm well (inverted index, compact time series, Bloom)
- Partial JSON updates are frequent and large
- Your platform already ships Redis Stack or the module
- Correctness of a homemade index would be worse than the module

Homemade search with `SCAN` and string match is a trap at scale. Homemade time series with unbounded sorted sets is a memory trap. Homemade Bloom with bitmaps is possible (topic 7) but easy to get wrong at huge offsets.

A module is not always faster. A module command still runs on the main thread (unless the module docs say otherwise). A huge `FT.SEARCH` can stall.

Licensing and vendor lock-in are part of the choice. Redis Stack, modules, and forks (topic 20) have different licenses. Read the current license before you ship.

Do not load an unmaintained module on a production primary.

### Questions

#### Theoretical questions

1. When is a core type plus Lua enough?
2. When is homemade `SCAN` search a trap?
3. Does a module command avoid the main thread by default?
4. Why does platform support matter?
5. Why do you read the license before you ship?

#### Easy practical tasks

1. Make a two-column table: "Need" and "Module or core". Add JSON partial update, full text, daily active bits, unique count.
2. Write four sentences: RedisJSON vs hash.
3. List two homemade designs this section calls traps.
4. Open the Redis Stack page. Write the module list that the page shows today.

#### Medium practical tasks

1. For your last project, pick three features. Mark module vs core. Give one reason each.
2. Read Redis bitmap Bloom vs RedisBloom. Write when a bitmap is enough (dense ids).
3. Find the license page for Redis Stack or RedisJSON. Write the license name and date you read it.

#### Advanced practical tasks

1. Write a decision record: adopt Search vs keep PostgreSQL full-text. Include ops, RAM, and query shapes.
2. Prototype the same feature with core types and with a module. Measure RAM and one query time. Write limits of the test.

---

## Compatibility and ops cost

A module changes the operations story.

Version matrix: Redis version, module version, Redis Stack image tag, and client library must match. A command that exists in Stack 7 may not exist in your old server.

Persistence: RDB and AOF must restore module values. Test `BGSAVE`, restart, and replica sync with module keys. A replica without the module cannot load the data.

Cluster: some modules have limits on multi-key commands and on how indexes shard. Read the module Cluster page. Do not assume every module is Cluster-safe.

Backup: a dump of module data is not a text JSON export unless you build one. Backup RDB/AOF as usual (topic 19). Test restore on a machine that has the same modules.

Memory: indexes and filters are extra. `INFO modules` and module-specific `INFO` sections help. `MEMORY USAGE` may not show a full index cost on another key.

Security: `MODULE LOAD` is a dangerous command (topic 17). Load only signed or vendor-supported modules. Restrict who can load.

Clients: the application must send module commands. Older clients still can use generic command APIs.

Observability: add module-specific metrics (index size, filter capacity) to the Redis dashboard.

Cost: more RAM, more upgrade steps, more failure modes, fewer people who have seen the incident class.

Use Redis Stack in Docker for learning:

```text
docker run --name stack-learn -p 6379:6379 redis/redis-stack
```

Pin an image digest in production. Do not run `latest` without a pin.

### Questions

#### Theoretical questions

1. What four versions must you keep in a matrix?
2. Why must a replica load the same module?
3. Why is `MODULE LOAD` restricted?
4. What extra backup test do module keys need?
5. Why pin a Redis Stack image in production?

#### Easy practical tasks

1. Run `INFO modules` or `MODULE LIST`. Write the module names (empty is a valid result).
2. Write four sentences: vanilla Redis vs Redis Stack ops.
3. Write a `docker run` line for Redis Stack on port `6380`.
4. List three extra failure modes that modules add.

#### Medium practical tasks

1. Start Stack and vanilla Redis on two ports. Compare `COMMAND LIST` counts or try `JSON.SET` on both. Write the difference.
2. Restart Stack after `JSON.SET` (or another module write). Confirm the key. Write persistence mode from `INFO persistence`.
3. Read Cluster notes for one module. Write one limit in your own words.

#### Advanced practical tasks

1. Write an upgrade plan: Stack image A to B, replica first, module compatibility check, rollback.
2. Build a restore drill: copy RDB from Stack, start a new Stack container with the same modules, confirm one JSON and one index (rebuild if needed). Document whether the index survived.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do RedisJSON, Search, TimeSeries, and Bloom map to "document", "query", "metric", and "maybe-seen"?
2. Which module costs are RAM-heavy even when key count looks small?
3. When do you refuse a module even if the feature fits?
4. How do topic 7 core types (bitmap, HLL, sorted set) still replace a module in small designs?
5. What must be true of persistence, replicas, and clients before you ship a module?

#### Easy practical tasks

1. Write a cheat sheet: `JSON.`, `FT.`, `TS.`, `BF.`, `MODULE LIST`, Redis Stack image.
2. Draw four boxes (JSON, Search, TimeSeries, Bloom) and one example key each.
3. Check your instance for modules. Write present or absent for each of the four.
4. Write five lab rules: pin images, test restart, no random `MODULE LOAD`, least privilege, memory test.

#### Medium practical tasks

1. If Stack is available, run one command from each module. Save replies. If not, write a compose file that starts Stack.
2. Write a compatibility table: feature, module, vanilla alternative, Cluster note, persist test.
3. Document client support in your language for JSON and Search (library names or "raw commands").

#### Advanced practical tasks

1. Write a production module policy: allowed modules, owners, license review, upgrade cadence, restore drill.
2. Compare Redis Stack with "core Redis + one external search + one TSDB". Write cost, latency, and team skills.
