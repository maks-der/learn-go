# 8. Storage and Retention

## Description

Kafka stores each partition as a log on disk. The log is a sequence of segment files. Retention rules delete old data. Compaction keeps the latest value per key. Tombstones mark a key as deleted.

This topic covers log segments, retention by time and by size, compaction, compact-plus-delete, tombstones, and `log.dirs`. Complete this topic after topic 2. Topic 9 covers copies of this data across brokers.

Use one term for each concept. A segment is one file (plus index files) in a partition log. Retention is the rule that removes old records. Compaction is the rule that keeps the last record per key. A tombstone is a record with a null value. Use KRaft for metadata. KRaft does not replace partition logs. Partition data still lives in `log.dirs`.

---

## Log segments

A partition log is not one endless file. Kafka splits the log into segments. The active segment receives new appends. When the active segment reaches a size or a time limit, Kafka closes it and starts a new active segment.

Each segment has a base offset. The file name uses that offset. Kafka also writes index files. The offset index maps an offset to a position in the `.log` file. The time index maps a timestamp to an offset. You do not edit these files by hand.

Retention and compaction run on closed segments. The active segment is not a candidate for delete in the same way. Very long active segments delay cleanup. Configuration such as `log.segment.bytes` and `log.roll.ms` controls when a new segment starts.

Small segments create many files. Large segments delay delete and make recovery heavier. Use the defaults unless you have a measured reason.

Do not delete segment files in the file system to "free disk". Use retention, or remove the topic with the admin API. Manual delete corrupts the partition.

### Questions

#### Theoretical questions

1. What is a log segment?
2. Which segment receives new records?
3. What do the index files do at a high level?
4. Why does a very large active segment delay cleanup?
5. Why must you not delete `.log` files by hand?

#### Easy practical tasks

1. Write five sentences about segments. Use only facts from this section.
2. Find `log.segment.bytes` in the official docs. Write the default.
3. Make a table: File kind, role. Add log, offset index, time index.
4. Draw a partition as three closed segments and one active segment.

#### Medium practical tasks

1. Locate `log.dirs` on your broker. List files for one partition. Write the file names that you see.
2. Set a small `segment.bytes` on a test topic if the config exists. Produce enough data to roll a segment. List the new files.
3. Read official log segment documentation. Write how a fetch uses the index.

#### Advanced practical tasks

1. Compare disk file count for the same data with a small segment size and with the default. Write a table.
2. Write a one-page note on why rolling by time (`log.roll.ms`) helps time-based retention on a quiet topic.

---

## Retention by time and by size

The default cleanup policy is `delete`. Kafka removes old records from closed segments when a retention limit is true.

**Time.** `retention.ms` (or `log.retention.hours` at broker level) is the maximum age of a record. When the time limit expires, Kafka can delete the segment.

**Size.** `retention.bytes` is the maximum size of the partition log. When the log is larger, Kafka deletes the oldest segments.

If you set both limits, Kafka can delete when either limit applies (check the exact rule in the docs for your version; the usual idea is that both bounds exist and delete can occur from either).

Retention is per partition, not per topic as one blob. A hot partition hits the size limit first.

`retention.ms=-1` and a disabled size limit can keep data forever. Disk will fill. Use infinite retention only when you have a disk plan.

Consumers with a committed offset that is now deleted have a problem. `auto.offset.reset` applies. You can lose the logical position. Monitor lag and disk.

### Questions

#### Theoretical questions

1. What does `delete` cleanup do?
2. What does `retention.ms` limit?
3. What does `retention.bytes` limit?
4. Is retention per partition or per topic as one size?
5. What is the risk of infinite retention?

#### Easy practical tasks

1. Write four sentences about time and size retention.
2. Make a table: Config key, unit, meaning. Add `retention.ms` and `retention.bytes`.
3. Find the broker defaults for retention on your install.
4. List two topics that need short retention and two that need long retention.

#### Medium practical tasks

1. Create a topic with `retention.ms` of a few minutes and a small segment size. Produce. Wait. Describe offsets. Write whether the earliest offset moved.
2. Create a topic with a small `retention.bytes`. Produce more than that size. Write what happens to old records.
3. Read official retention configuration. Rewrite the two limits in STE.

#### Advanced practical tasks

1. Measure time from "segment closed" to "segment gone" on a short-retention topic. Write the steps and the time.
2. Write a disk plan for a 100 MB/s topic with 7-day retention and replication factor 3. Show the arithmetic.

---

## Compaction (`cleanup.policy=compact`)

Compaction keeps the latest record for each key in a partition. Older records with the same key can be removed. Records with different keys stay.

Use compaction for a changelog: the current state of an entity is the last value for that key. Examples: user profile, account status, stream table state.

Compaction requires keys. Null keys do not compact as an entity changelog. Use a stable key.

Compaction is not instant. A cleaner thread reads closed segments and writes compacted segments. The latest value can exist more than once until the cleaner runs.

Compaction does not make a topic a database. You cannot query by key through the Kafka protocol as you query SQL. A consumer still reads the log. Kafka Streams and compacted topics together can build a table (topic 14).

Set `cleanup.policy=compact` on the topic. Do not assume the broker default is compact. The usual default is `delete`.

### Questions

#### Theoretical questions

1. What does compaction keep?
2. What records can compaction remove?
3. Why does compaction need keys?
4. Why can two records for the same key still exist for a time?
5. Why is a compacted topic not a SQL database?

#### Easy practical tasks

1. Write five sentences about compaction. Use only facts from this section.
2. Make a table: Topic type, policy. Add metrics stream, user profile changelog.
3. Find `cleanup.policy` in the docs. Write the allowed values.
4. Draw three keys with several values. Mark what compaction keeps.

#### Medium practical tasks

1. Create a compacted topic. Produce three values for key `u1` and one value for `u2`. Wait for compaction if you can trigger it (or produce enough data). Consume from the beginning. Write the remaining records.
2. Produce records with null keys to a compacted topic. Write why this is the wrong use.
3. Read official log compaction documentation. Write the role of the cleaner thread.

#### Advanced practical tasks

1. Tune `min.cleanable.dirty.ratio` or related keys on a test topic (official names). Document how soon compaction runs.
2. Write a one-page design for a compacted `users.profile` topic: key, value fields, and a consumer that builds an in-memory map.

---

## Compact + delete

A topic can use both policies: `cleanup.policy=compact,delete`. Kafka compact-cleans keys and also deletes records that are older than the retention time (or that exceed size, depending on configuration).

Use compact+delete when you want a changelog that does not grow forever. Old keys that no one updates still leave a record in a compact-only topic. The delete part can remove those old records after `retention.ms`.

Understand the product rule: compaction keeps the latest value per key. Delete removes data by age or size. Together they bound disk and still collapse duplicates of a hot key.

This combination is common for changelog topics that must not retain unused keys for years.

Test the policy on a non-production topic. Confirm that a quiet key disappears after the retention time. Confirm that a hot key still collapses to one recent value.

Do not set a retention time that is shorter than the time a consumer needs to rebuild state, unless the consumer can accept a missing key.

### Questions

#### Theoretical questions

1. What does `compact,delete` mean?
2. What problem does compact-only have with old keys?
3. What does the delete part bound?
4. Why must you test this policy?
5. How can a short retention break a state rebuild?

#### Easy practical tasks

1. Write four sentences about compact+delete.
2. Make a table: Policy, hot key duplicates, unused old keys.
3. Find an official mention of `compact,delete`. Write it in your own words.
4. List two topics that fit compact+delete and one that must be delete-only.

#### Medium practical tasks

1. Create a topic with `cleanup.policy=compact,delete` and a short `retention.ms`. Produce a key once. Wait. Consume. Write whether the key remains.
2. Produce many updates for one key. After cleanup, write how many records remain for that key.
3. Compare `describe --topics` config for a delete topic and a compact,delete topic.

#### Advanced practical tasks

1. Write a one-page choice: compact versus compact,delete for Kafka Streams changelog topics (high-level). Use official Streams notes if you read them.
2. Design a test that fails if an unused key stays longer than retention on a compact,delete topic.

---

## Tombstones

A tombstone is a record with a key and a null value. On a compacted topic, a tombstone means "this key is deleted". The cleaner can remove earlier values for that key. After a tombstone retention time, the tombstone itself can go away.

Producers send a tombstone when the entity is gone. Example: user deleted. Consumers that build a table remove the key from their table when they see a null value.

If you only stop sending updates, compaction keeps the last non-null value. The key stays. You must send a tombstone to delete the key from a compact changelog.

A null value on a `delete`-only topic is just a record with a null value. It does not have the same cleanup meaning as a compact tombstone.

Do not use a tombstone as a normal payload. Use a structured "deleted" event on a delete-policy topic if consumers must keep a history of the delete. Use a tombstone when the changelog must drop the key.

`delete.retention.ms` controls how long tombstones stay (official name). Readers need that window to see the delete.

### Questions

#### Theoretical questions

1. What is a tombstone?
2. What does a tombstone mean on a compacted topic?
3. What happens if you never send a tombstone for a removed entity?
4. Is a null value a tombstone on a delete-only topic?
5. Why must tombstones stay for a time?

#### Easy practical tasks

1. Write five sentences about tombstones. Use only facts from this section.
2. Draw a key with two values and one tombstone. Mark what a table consumer does.
3. Find `delete.retention.ms` in the docs. Write the meaning.
4. List three entities that need tombstones on a changelog.

#### Medium practical tasks

1. On a compacted topic, produce `u1=A`, then `u1` with a null value (console or client). Consume. Write what you see.
2. Write a small table consumer (in-memory map). Apply a tombstone. Confirm that the key is gone.
3. Read official tombstone documentation. Rewrite the lifecycle in STE.

#### Advanced practical tasks

1. Wait for tombstone retention on a test topic (short config). Show that a late consumer can miss the delete. Write the risk.
2. Write a one-page producer standard: when to send a tombstone versus a `UserDeleted` event on a history topic.

---

## Disk and `log.dirs`

`log.dirs` (or `log.dir`) is the broker configuration that names the directories for partition logs. You can set more than one directory. Kafka assigns partitions to directories. This is not a RAID manager. Use a proper disk setup under the directories.

Disk full is a serious failure. The broker can stop writes. Monitor disk. Retention is your main tool to free space. Compaction can reduce space for changelogs. Delete of a topic frees space after the delete process finishes.

Copies multiply disk. A replication factor of 3 stores about three times the data. Include replicas in disk plans.

Page cache matters. Kafka reads often come from the operating system cache. Do not run other disk-heavy jobs on the same disks without a plan. Topic 15 in the path covers operations.

KRaft metadata has its own storage directory in the controller configuration. Do not confuse the metadata log with `log.dirs` for topic partitions. Both need disk. Format both when you create a cluster (official format command).

Never point `log.dirs` at a network file system that does not meet Kafka requirements. Follow official disk guidance.

### Questions

#### Theoretical questions

1. What is `log.dirs`?
2. Why does replication factor 3 change disk plans?
3. What happens when the disk is full?
4. How is KRaft metadata storage different from `log.dirs`?
5. Why is a network file system a risk?

#### Easy practical tasks

1. Find `log.dirs` in your `server.properties` or container env. Write the path.
2. Write four sentences about disk and Kafka. Use only facts from this section.
3. Make a table: Disk user, directory kind. Add partition logs and KRaft metadata.
4. Estimate size: 10 MB/s, 2 days, 3 replicas. Write the arithmetic.

#### Medium practical tasks

1. List disk usage of your local `log.dirs`. Produce a large batch. List again. Write the delta.
2. Set a short retention. Wait. Write how usage changes.
3. Read official `log.dirs` documentation. Write how multiple directories work.

#### Advanced practical tasks

1. Fill a small test volume (or a quota) until produce fails. Record the error. Restore space. Write a monitor you would add.
2. Write a one-page disk standard: filesystem type, one Kafka role per disk set, and alerts at 60% and 80%.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe the life of a record from append to the active segment to delete or compact.
2. How do segment roll, retention, and the cleaner thread depend on each other?
3. When do you choose delete, compact, or compact+delete?
4. How do tombstones and `delete.retention.ms` protect table consumers?
5. Why do `log.dirs` and KRaft metadata directories both appear in a healthy install?

#### Easy practical tasks

1. Create `store.review` with a short time retention. Produce five records. Describe the topic config.
2. Write a one-page cheat sheet: segment, retention.ms, retention.bytes, compact, compact+delete, tombstone, log.dirs.
3. From a topic describe, highlight cleanup.policy and retention keys.
4. Draw compact+delete on a time axis for one unused key and one hot key.

#### Medium practical tasks

1. Write a script that creates a compacted topic, produces keys, produces a tombstone, and consumes from the beginning.
2. Compare disk size of a JSON changelog before and after many updates to the same keys (after compaction if it runs).
3. Map each subsection to an official configuration key.

#### Advanced practical tasks

1. Build a small "table from compacted topic" program. Restart it. Confirm that it rebuilds from the log. Then send a tombstone and confirm the key is gone.
2. Write a production storage standard: defaults for event streams versus changelogs, disk formula with replicas, and a ban on manual file delete.
