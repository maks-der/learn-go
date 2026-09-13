# 4. Sorted Sets

## Description

A sorted set stores unique members. Each member has a score. Redis keeps the members in score order. This topic shows the score model, the main write and range commands, leaderboards, and simple sliding windows.

Complete this topic before you study streams, bitmaps, and other types.

Use one term for each concept. A member is a unique string. A score is a double-precision number. A rank is the position in the ordered set (zero-based). Do not treat a sorted set as a list. You address members by name and by score, not only by push and pop.

---

## Score and member

A sorted set is a set plus a number per member. The member is unique. The score can repeat. Two members can share the same score. Redis then orders those members by the member string (lexicographic order).

The score is a 64-bit floating-point number. Large integers can lose precision. For time, many teams store Unix time in seconds or milliseconds as the score. For ranks, they store points.

You cannot store two copies of the same member. A new `ZADD` of the same member updates the score.

A missing member has no score. `ZSCORE` returns null.

The sorted set is one key. Typical names:

```text
lb:game:42
z:recent:user:9
z:delay:emails
```

Pick the score so that range queries match the problem. If you need "top 10 by points", the score is points. If you need "events in the last 60 seconds", the score is a timestamp.

Memory cost is higher than a set. Redis keeps a dict and a skip list (or a compact encoding for small keys). Use a sorted set when you need order or ranges. Use a set when you need only uniqueness.

Do not use the member as a huge JSON blob. Store an id as the member. Store the object in a hash.

### Questions

#### Theoretical questions

1. What two pieces of data does each sorted-set entry have?
2. Can two members have the same score?
3. What happens when you `ZADD` an existing member with a new score?
4. What does `ZSCORE` return for a missing member?
5. Why store an id as the member instead of a large JSON string?

#### Easy practical tasks

1. Write five sentences that describe score and member.
2. Make a three-column table: member, score, meaning. Add four rows for a game.
3. Propose a score for "latest login time" and a score for "total spend".
4. Open the sorted set data-type page. Write one sentence from that page.

#### Medium practical tasks

1. Explain lexicographic order when scores are equal. Give three member names and their order at score `1`.
2. Find the score type in the docs (IEEE double). Write one limit for large integer ids as scores.
3. Draw a small skip-list idea (boxes and arrows) at a beginner level. Label "ordered by score".

#### Advanced practical tasks

1. Read why Redis uses a skip list plus a hash table for sorted sets. Write six sentences.
2. Design members and scores for flights (price and time). Write why you may need two sorted sets.

---

## `ZADD`, `ZRANGE`, `ZRANK`, `ZSCORE`, range by score

`ZADD key score member [score member ...]` adds or updates members. The reply is the number of new members (not the number of score updates), unless you use options such as `CH`.

Options include `NX` (add only if missing), `XX` (update only if present), `GT` / `LT` (update score only if greater or less), and `INCR` (add to the score, like `ZINCRBY`). Read the command page for your version.

`ZSCORE key member` returns the score as a string, or null.

`ZRANK key member` returns the rank from lowest score to highest. The lowest score has rank `0`. A missing member returns null.

`ZREVRANK` returns the rank from highest score to lowest. The highest score has rank `0`.

`ZRANGE key start stop` returns members by rank. `0` is the lowest score. `-1` is the highest. Add `WITHSCORES` to include scores.

```text
ZADD lb 100 ada 80 bea 80 cam
ZRANGE lb 0 -1 WITHSCORES
ZRANK lb ada
ZSCORE lb bea
```

`ZREVRANGE` returns members from high rank to low rank. Redis 6.2+ prefers `ZRANGE ... REV` instead of `ZREVRANGE`. Both exist. New code can use `ZRANGE` with `REV`.

`ZCARD` returns the member count. `ZCOUNT key min max` counts members in a score range.

`ZREM key member` removes members. `ZINCRBY key increment member` adds to a score.

Ranks change when scores change. Do not cache a rank in the application for a long time without a refresh.

`ZRANGEBYSCORE key min max` returns members whose scores lie between `min` and `max`. Both bounds are inclusive by default. Prefix a bound with `(` to make it exclusive.

```text
ZRANGEBYSCORE events 1700000000 1700003600
ZRANGEBYSCORE events (10 (20
```

Use `-inf` and `+inf` for open ends.

`LIMIT offset count` pages through the range. Always use `LIMIT` on large ranges.

Redis 6.2+ also allows:

```text
ZRANGE key min max BYSCORE
ZRANGE key min max BYSCORE REV
```

Prefer `ZRANGE ... BYSCORE` in new code when your server supports it. `ZRANGEBYSCORE` still works.

`ZREMRANGEBYSCORE key min max` deletes members in that score range. The reply is the number of removed members.

`ZREMRANGEBYRANK key start stop` deletes by rank. `ZPOPMIN` and `ZPOPMAX` remove one or more members from the ends.

These commands implement time windows: score is Unix time, `ZRANGEBYSCORE` reads the window, `ZREMRANGEBYSCORE -inf old-cut` drops old members.

Do not run `ZRANGEBYSCORE -inf +inf` on a huge key without `LIMIT`. The command can block and return a large array.

`ZCOUNT` is cheaper than fetching all members when you only need the size of a window.

### Questions

#### Theoretical questions

1. What does `ZADD` return by default?
2. What is rank `0` for `ZRANK`?
3. What is rank `0` for `ZREVRANK`?
4. What does `ZRANGEBYSCORE` select on?
5. How do you make a score bound exclusive?

#### Easy practical tasks

1. `ZADD` three members with different scores. `ZRANGE 0 -1 WITHSCORES`. Save the reply.
2. `ZSCORE` one member. `ZRANK` the same member. Save both.
3. `ZADD tw 10 a 20 b 30 c`. `ZRANGEBYSCORE tw 10 20`. Save the reply.
4. `ZREMRANGEBYSCORE tw -inf 15`. `ZRANGE tw 0 -1`. Save the reply.

#### Medium practical tasks

1. Add two members with the same score. Show `ZRANGE` order and `ZRANK` of each.
2. Use `LIMIT` to read the first two members in a score range. Then read the next two (offset 2).
3. Compare `ZRANGEBYSCORE` with `ZRANGE ... BYSCORE` on your server. Write if both work.

#### Advanced practical tasks

1. Build 1,000 members with a pipeline. Time `ZRANK` of one member vs `ZRANGE 0 -1` in the client to find the rank. Write the lesson.
2. Insert 10,000 timestamp scores. Remove everything older than a cutoff with `ZREMRANGEBYSCORE`. Time the delete.

---

## Leaderboards

A leaderboard is a sorted set of player ids and scores. The score is points (or time, if a lower time is better).

Typical operations:

- Submit a score: `ZADD lb:season1 2500 user:42` or `ZADD ... GT` so that Redis keeps only a better score
- Read the top 10: `ZREVRANGE lb:season1 0 9 WITHSCORES` or `ZRANGE ... REV`
- Read one player rank: `ZREVRANK lb:season1 user:42` (add 1 for a 1-based place)
- Read a score: `ZSCORE lb:season1 user:42`
- Read neighbors: get the rank, then `ZREVRANGE` around that rank

Display names and avatars do not belong in the member if they change. Store `user:42` as the member. Load the profile from a hash.

When two players share a score, document the tie rule (lexicographic member). If you need "first to reach the score wins", encode time in the score (a composite number) or use a second sorted set.

Sharding a huge leaderboard is hard. One sorted set lives on one instance (or one Cluster slot). For very large games, teams split boards (per region, per week) or use an application tier.

Expiry of a whole season: set a TTL on the key `lb:season1`, or delete the key at season end. You cannot TTL one member in core Redis. Remove a player with `ZREM`.

Do not `ZRANGE 0 -1` to sort in the client. Use ranks and ranges.

### Questions

#### Theoretical questions

1. Which command updates a player's points?
2. How do you read the top 10 by highest score?
3. How do you get a 1-based place from `ZREVRANK`?
4. Why is the member an id and not a display name?
5. Why is one giant global leaderboard hard to shard?

#### Easy practical tasks

1. Build a board with five players. Print the top 3 with scores.
2. `ZREVRANK` one player. Write place as rank plus one.
3. `ZINCRBY` a player by 50. Show the new top 3.
4. Draw a table: place, member, score.

#### Medium practical tasks

1. Implement "keep best score only" with `ZADD GT` or a `ZSCORE` compare in the client. Document the race of the two-command path.
2. Show a window of two players above and below a given player using rank math.
3. Create `lb:week1` with `EXPIRE` of 8 days. Write how a new week uses a new key.

#### Advanced practical tasks

1. Encode score and time in one double so that ties break by earlier time. Show the formula and two examples. Test order with `ZREVRANGE`.
2. Write a leaderboard API list (submit, top N, me, around me) and the Redis command for each. Note Cluster hash tags if two keys must stay together.

---

## Sliding windows

A sliding window keeps members that belong to the last N seconds (or the last N events). A sorted set can model the simple case.

Time window (unique event ids):

1. `ZADD window:user:9 <now> event-id`
2. `ZREMRANGEBYSCORE window:user:9 -inf <now - 60>`
3. `ZCARD window:user:9` — count in the last 60 seconds

That pattern is a rate-limit or "actions in the last minute" counter when each event has a unique id. The same member cannot appear twice. If the same id can repeat, add a unique suffix.

Score-only windows without ids can use the member as a random nonce or as `timestamp-seq`.

Event cap (last N items): use rank deletes. After `ZADD`, `ZREMRANGEBYRANK key 0 -<N+1>` if you use time as score and want the newest N. Check the rank direction. Many teams use `ZADD` plus `ZCARD` plus `ZPOPMIN` until the size is N.

This model is not RedisTimeSeries. You do not get aggregation, downsampling, or compact samples. You get exact members in a window until you delete them.

Limits:

- Each event uses memory for member and score
- `ZREMRANGEBYSCORE` on every request has a cost
- Cluster: all `ZADD` and remove commands must hit the same key

For heavy time series, use RedisTimeSeries (topic 5 survey) or a purpose-built store. For simple rate limits and "recent ids", a sorted set is enough.

You can store the window key with a TTL a bit longer than the window so that idle keys disappear.

### Questions

#### Theoretical questions

1. What do you store as the score in a time window?
2. Which command drops members that are too old?
3. How do you count events in the window after the trim?
4. Why must event members be unique in a sorted set?
5. When do you use RedisTimeSeries instead of this pattern?

#### Easy practical tasks

1. `ZADD win 100 e1 120 e2 200 e3`. Remove scores below 150. `ZRANGE 0 -1`. Save the reply.
2. `ZCARD win` after the trim. Save the reply.
3. Write the three-step window algorithm in your own words.
4. Propose key names for per-user 60-second windows.

#### Medium practical tasks

1. In `redis-cli` or a script, add five events with `TIME` or wall-clock scores. Trim to 60 seconds. Print `ZCARD`.
2. Implement last-5-events with `ZADD` and `ZREMRANGEBYRANK` or `ZPOPMIN`. Prove `ZCARD` is at most 5.
3. Set a TTL on the window key. Write why the TTL is longer than the window.

#### Advanced practical tasks

1. Script 10,000 `ZADD` plus trim per call. Measure time per operation. Write when the pattern becomes too slow.
2. Compare this window with a list of timestamps and with a string counter plus TTL. Write a three-row table: exact unique ids, memory, commands.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do rank ranges and score ranges differ, and which commands use each?
2. How does a leaderboard use `ZREVRANK` and `ZRANGE` together?
3. How does a sliding window use `ZADD` and `ZREMRANGEBYSCORE` together?
4. What does uniqueness of members change in a rate-limit window?
5. When do you pick a sorted set instead of a list or a set?

#### Easy practical tasks

1. Create `lab:z` with four members. Run `ZRANGE`, `ZREVRANGE` (or `REV`), `ZSCORE`, `ZRANK`. Save a table.
2. Write a cheat sheet: `ZADD`, `ZRANGE`, `ZRANGEBYSCORE`, `ZREM`, `ZREMRANGEBYSCORE`, `ZINCRBY`, `ZCARD`.
3. Draw one leaderboard and one time window. Label score meaning on each drawing.
4. `ZPOPMAX` once on `lab:z`. Show the remaining members.

#### Medium practical tasks

1. Build a weekly leaderboard key and a 60-second action window key. Run five commands on each. Document them.
2. Write a script that pages a score range with `LIMIT` until the range is empty.
3. Document tie-breaking for your game in ten lines (same score, member order, or composite score).

#### Advanced practical tasks

1. Implement submit, top 10, and around-me in a small program. Use one sorted set and one hash for names.
2. Read `ZUNIONSTORE` / `ZINTERSTORE` in the docs. Merge two weekly boards into `lab:month` with a weight of 1. Explain the result.
