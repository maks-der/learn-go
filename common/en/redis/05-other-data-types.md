# 5. Other Data Types

## Description

Redis has more types than strings, hashes, lists, sets, and sorted sets. This topic introduces bitmaps, HyperLogLog, streams, geospatial indexes, and Redis Stack modules (JSON, Search, TimeSeries). The goal is a correct first picture, not a full operations guide.

Complete this topic before you study pipelines, transactions, and Lua.

Use one term for each concept. A bitmap is a string that you address by bit offset. A HyperLogLog is an approximate counter of unique items. A stream is an append-only log of entries. A geo key stores positions. A module adds commands that core Redis does not have.

---

## Bitmaps

A bitmap is a string that Redis treats as a bit array. Bit `0` is the first bit of the first byte. You can set millions of bits. Redis grows the string to the highest offset that you write.

`SETBIT key offset value` sets bit `offset` to `0` or `1`. The reply is the old bit.

`GETBIT key offset` returns `0` or `1`. A missing key behaves as all zeros.

`BITCOUNT key` returns the number of bits that are `1`. You can pass a byte range.

`BITOP` runs `AND`, `OR`, `XOR`, or `NOT` across bitmaps and writes a destination key. `BITPOS` finds the first bit that is `0` or `1`.

```text
SETBIT lab:bits 0 1
SETBIT lab:bits 7 1
GETBIT lab:bits 7
BITCOUNT lab:bits
```

Typical use: one bit per user id for "was active today". Key `active:2026-09-13`, offset = user id. `BITCOUNT` is the daily active count. This works when ids are dense and not huge. A sparse 64-bit id space is a bad bitmap (the string would be enormous).

Bitmaps are not a new type in `TYPE`. `TYPE` returns `string`. Bitmap commands fail on a hash or list (`WRONGTYPE`).

`SETBIT` at a huge offset allocates RAM up to that offset. Do not test offset `2^32` on a small machine.

`GETRANGE` and `SETRANGE` also work because the value is a string. Prefer bit commands for bit logic.

### Questions

#### Theoretical questions

1. What Redis type does `TYPE` report for a bitmap key?
2. What does `SETBIT` return?
3. What does `BITCOUNT` count?
4. When is a bitmap a bad unique-user store?
5. What does `BITOP AND` do?

#### Easy practical tasks

1. `SETBIT lab:b 2 1`. `GETBIT lab:b 2` and `GETBIT lab:b 1`. Save the replies.
2. `BITCOUNT lab:b`. Save the reply.
3. Write four sentences: bitmap vs set of user ids.
4. Open the `SETBIT` command page. Write the complexity note about the offset.

#### Medium practical tasks

1. Mark users `1`, `3`, and `10` active. `BITCOUNT` the key. Flip user `3` to `0`. `BITCOUNT` again.
2. Use `BITOP OR` on two daily keys into `lab:or`. `BITCOUNT lab:or`.
3. `BITPOS lab:b 1`. Explain the index.

#### Advanced practical tasks

1. Estimate RAM for 10 million bits (one bit per user). Show the math. Compare with a set of 10 million integer members (order of magnitude).
2. Design daily + weekly active users with seven `BITOP OR` keys. Write the key names and the weekly job.

---

## HyperLogLog

A HyperLogLog (HLL) estimates the number of unique values that you added. It does not store the values. The memory is small (about 12 KB per key in Redis). The error is typically around 1 percent (standard HLL error; Redis documents the exact bound).

`PFADD key element [element ...]` adds elements. The reply is `1` if the estimate changed, else `0`.

`PFCOUNT key [key ...]` returns the estimate. Several keys give the estimate of the union.

`PFMERGE dest src [src ...]` writes a merged HLL into `dest`.

```text
PFADD lab:uv user:1 user:2 user:1
PFCOUNT lab:uv
```

The count is not exact. `PFADD` of the same element twice does not grow the true unique count. The estimate can still move by a small amount in edge cases. Do not use HLL for billing that needs an exact integer.

Use HLL for unique visitors, unique search terms, and other large cardinalities where a set would not fit.

`TYPE` of an HLL key is `string`. The encoding is special. Do not `GET` and parse it. Do not `SET` over an HLL key if you want to keep the sketch.

HLL does not tell you which users were unique. If you need the ids, use a set or a bitmap (when ids fit).

### Questions

#### Theoretical questions

1. What does HyperLogLog estimate?
2. Does Redis store each unique element in an HLL?
3. What is a typical error range for the estimate?
4. What does `PFCOUNT` of two keys estimate?
5. When must you not use HLL?

#### Easy practical tasks

1. `PFADD lab:h a b c a`. `PFCOUNT lab:h`. Save the reply.
2. `PFADD lab:h d`. `PFCOUNT lab:h`. Save the reply.
3. Write four sentences: HLL vs set.
4. Open the HyperLogLog page. Write the memory size that the page names.

#### Medium practical tasks

1. Add 1,000 sequential ids with a script. Compare `PFCOUNT` with 1000. Write the error.
2. `PFADD lab:h2` with an overlapping set. `PFCOUNT lab:h lab:h2`. Then `PFMERGE lab:m lab:h lab:h2`. `PFCOUNT lab:m`.
3. `TYPE lab:h`. Confirm it is `string`.

#### Advanced practical tasks

1. Add 1,000,000 ids (pipeline). Record `PFCOUNT`, `MEMORY USAGE`, and duration. Delete the key.
2. Read the official error description. Write a policy: "We accept HLL for analytics dashboards but not for invoices."

---

## Streams (`XADD`, `XREAD`) — intro

A stream is an append-only log. Each entry has an id and one or more field-value pairs.

`XADD key * field value [field value ...]` appends an entry. `*` asks Redis to create an id. The id looks like `timestamp-sequence`. You can pass an explicit id when you need it.

```text
XADD lab:s * type signup user 42
XLEN lab:s
XRANGE lab:s - +
```

`XREAD COUNT 10 STREAMS lab:s 0-0` reads from an id. Use `$` to read only new entries after the current end. `XREAD BLOCK 5000 STREAMS lab:s $` waits for a new entry.

`XREAD` is a fan-out read. Each client tracks its own id. Redis does not record that a client processed an entry.

A consumer group tracks a last-delivered id and pending entries per consumer. Commands (preview):

- `XGROUP CREATE lab:s group1 0 MKSTREAM`
- `XREADGROUP GROUP group1 consumerA COUNT 1 STREAMS lab:s >`
- `XACK lab:s group1 <id>`

`>` means "new messages for the group". Pending entries that a consumer did not acknowledge stay in the PEL (pending entries list). Another consumer can `XCLAIM` them. Topic 10 covers the full model.

Use a stream when you need a log, many consumers, and replay. Use Pub/Sub when you need fire-and-forget broadcast and you accept lost messages. Use a list queue when a simple competing pop is enough.

`XTRIM` or `MAXLEN` on `XADD` caps the stream. Unbounded streams fill memory.

### Questions

#### Theoretical questions

1. What does `XADD` with `*` set for the entry id?
2. What is the difference between `XREAD` and `XREADGROUP`?
3. What does `XACK` mean?
4. What problem does a consumer group solve that a list `RPOP` does not solve?
5. Why must you trim a stream?

#### Easy practical tasks

1. `XADD lab:s * msg hello`. Save the id.
2. `XRANGE lab:s - +`. Save the reply.
3. `XLEN lab:s`. Save the reply.
4. Write four sentences: stream vs list queue.

#### Medium practical tasks

1. `XREAD COUNT 2 STREAMS lab:s 0-0`. Then `XADD` one more entry. `XREAD` from the last id you saw.
2. Create a group with `XGROUP CREATE`. `XREADGROUP` as `c1`. Read `XPENDING`. Do not `XACK` yet. Write what you see.
3. `XADD lab:s MAXLEN 5` in a loop of 10 adds (or `XTRIM`). Confirm `XLEN` is at most 5.

#### Advanced practical tasks

1. Run two consumers in one group. Add four entries. Write how the entries split. `XACK` from each consumer.
2. Draw a diagram: stream, group, two consumers, PEL, `XACK`. Label `>` vs an explicit id.

---

## Geospatial commands

A geo key stores members with longitude and latitude. Redis stores them in a sorted set. The score is a geohash. `TYPE` may show `zset`. Use geo commands, not raw `ZADD`, unless you know the encoding.

`GEOADD key longitude latitude member` adds or updates a point. Longitude is first. The range is valid Earth coordinates. An invalid coordinate returns an error.

```text
GEOADD lab:geo 13.405 52.52 berlin
GEOADD lab:geo -0.1276 51.5074 london
```

`GEOSEARCH` (Redis 6.2+) finds members near a point or inside a box. You set radius, unit (`m`, `km`, `mi`, `ft`), and optional `ASC` / `DESC` and `COUNT`.

```text
GEOSEARCH lab:geo FROMLONLAT 13.4 52.5 BYRADIUS 300 km ASC
```

Older commands `GEORADIUS` and `GEORADIUSBYMEMBER` still exist. New code uses `GEOSEARCH` and `GEOSEARCHSTORE`.

`GEODIST` returns the distance between two members. `GEOPOS` returns coordinates. `GEOENCODE` / `GEODECODE` helpers exist on some versions.

Geo is a 2D index on a sphere model. It is not a full GIS. You do not get polygons of a country, projections, or road routing. Accuracy is enough for "stores near the user", not for surveying land.

Very large radii and huge sets are expensive. Use `COUNT` to cap results.

If you need queries by city name, store a hash for metadata and keep the geo key for distance only.

### Questions

#### Theoretical questions

1. In what order does `GEOADD` take coordinates?
2. What older type does Redis use under a geo key?
3. What command replaces `GEORADIUS` in new code?
4. What does `GEODIST` return?
5. Why is Redis geo not a full GIS?

#### Easy practical tasks

1. `GEOADD` two cities. `GEOPOS` one city. Save the reply.
2. `GEODIST lab:geo berlin london km`. Save the reply.
3. `GEOSEARCH` from Berlin with a 1000 km radius. Save the members.
4. Write four sentences: geo key vs two hash fields `lat` and `lon`.

#### Medium practical tasks

1. Add five points. Search `BYRADIUS` with `COUNT 2`. Confirm you get at most two members.
2. Use `FROMMEMBER berlin` if the syntax allows it on your version. Compare with `FROMLONLAT`.
3. Try an invalid longitude (for example `200`). Record the error.

#### Advanced practical tasks

1. `GEOSEARCHSTORE` results into `lab:near`. `ZRANGE lab:near 0 -1`. Explain the destination type.
2. Design keys for "shops in a city" plus geo search. Write how you filter by city without scanning the world.

---

## Redis Stack modules (survey)

Core Redis does not include a JSON document type, a full text index, or a compact time-series type. Redis Stack and Redis modules add those features.

RedisJSON (`JSON.SET`, `JSON.GET`, `JSON.ARRAPPEND`) stores JSON and updates paths. You can change `$.user.city` without a full `GET` / `SET` of a string. Use it when documents are large and updates are partial.

RediSearch (`FT.CREATE`, `FT.SEARCH`) indexes hashes or JSON. You can query text, numbers, and tags. This is a search engine beside Redis, not SQL. You must define a schema. You must maintain the index.

RedisTimeSeries (`TS.ADD`, `TS.RANGE`, `TS.CREATERULE`) stores timestamped samples. You get aggregation (avg, min, max) and downsampling. Use it for metrics. Do not store metrics as one sorted-set member per sample at large scale without a plan.

RedisBloom and related filters add Bloom filters, Cuckoo filters, and Count-Min sketches. They answer "possibly in the set" with less memory than a set.

Modules load into the server. Commands appear in `COMMAND`. A client must support the extra commands. Cluster and persistence must match module rules. Operations cost goes up (versions, RAM, backup of module data).

Use a module when the module matches the access pattern and your platform ships it (Redis Stack, Redis Cloud, or a managed module). Do not load random modules on a production core Redis without a review.

You can complete topics 1 through 11 with core Redis. Use Docker image `redis/redis-stack` when you practice this survey. Topic 12 covers module choice and production use.

### Questions

#### Theoretical questions

1. What problem does RedisJSON solve that a string of JSON does not solve well?
2. What does RediSearch require before `FT.SEARCH` works?
3. What extra operations does RedisTimeSeries give over a sorted set window?
4. What kind of answer does a Bloom filter give?
5. Why does a module increase operations cost?

#### Easy practical tasks

1. Open the Redis Stack page. List four module names that the page shows.
2. Write a two-column table: module and one command prefix (`JSON.`, `FT.`, `TS.`).
3. Check your instance: run `JSON.SET` or `FT._LIST`. Write whether the module is present.
4. Write four sentences: when to stay on core Redis types.

#### Medium practical tasks

1. If Stack is available, `JSON.SET lab:j $ '{"n":1}'` and `JSON.GET lab:j $.n`. Save the replies. If not, write the Docker command to start Redis Stack.
2. Read one `FT.CREATE` example in the docs. Write the index name, the prefix, and two fields.
3. Compare `TS.ADD` with `ZADD` for 100 samples in a short paragraph.

#### Advanced practical tasks

1. Install Redis Stack in Docker. Run one JSON, one Search, and one TimeSeries command. Save the outputs.
2. Write a decision page: core type vs module for documents, search, metrics, and unique counts (HLL vs Bloom vs set).

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Which of these features are core Redis types or string encodings, and which need a module?
2. How do bitmap, HyperLogLog, and set differ for "unique users today"?
3. How do streams improve on list queues for acknowledgement?
4. What coordinate order and command names do you use for a first geo query?
5. When do you accept an approximate answer (HLL, Bloom) instead of an exact set?

#### Easy practical tasks

1. Create one bitmap key, one HLL key, and one stream key. Run `TYPE` on each. Save the three types.
2. Write a cheat sheet: `SETBIT`, `PFADD`, `XADD`, `GEOADD`, and three module prefixes.
3. Draw a table: type, exact?, stores ids?, typical size.
4. `DEL` the lab keys from this topic that you still have.

#### Medium practical tasks

1. Solve "daily active users" three ways (set, bitmap, HLL) for 20 fake ids. Write counts and a size guess.
2. Write a 15-line stream consumer-group runbook: create group, read, ack, what to do if no ack.
3. Document when your team may enable Redis Stack and which module you would try first.

#### Advanced practical tasks

1. Build a small demo: geo shop search plus a stream of "shop viewed" events and an HLL of unique viewers. List all keys.
2. Read compatibility notes for modules on Redis Cluster. Write three constraints (slots, commands, or persistence).
