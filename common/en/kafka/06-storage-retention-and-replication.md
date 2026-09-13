# 6. Storage, Retention, and Replication

## Description

Kafka stores each partition as a log on disk. The log is a sequence of segment files. Retention rules delete old data. Compaction keeps the latest value per key. Replication copies each partition to more than one broker.

This topic shows log segments, retention by time and by size, compaction, tombstones, ISR, `min.insync.replicas`, unclean leader election, and the KRaft controller. Complete this topic after topics 2 and 3. Prefer KRaft for all new clusters. Do not start new clusters on ZooKeeper.

Use one term for each concept. A segment is one file (plus index files) in a partition log. Retention is the rule that removes old records. Compaction is the rule that keeps the last record per key. A tombstone is a record with a null value. A replica is a copy of a partition. The ISR is the set of replicas that are up to date. The controller is the component that manages cluster metadata and partition leaders. In new Kafka, the controller is a KRaft quorum.

---

## Log segments

A partition log is not one endless file. Kafka splits the log into **segments**. The active segment receives new appends. When the active segment reaches a size or a time limit, Kafka closes it and starts a new active segment.

Each segment has a base offset. The file name uses that offset. Kafka also writes index files. The offset index maps an offset to a position in the `.log` file. The time index maps a timestamp to an offset. You do not edit these files by hand.

Retention and compaction run on closed segments. The active segment is not a candidate for delete in the same way. Very long active segments delay cleanup. Configuration such as `log.segment.bytes` and `log.roll.ms` controls when a new segment starts.

Small segments create many files. Large segments delay delete and make recovery heavier. Use the defaults unless you have a measured reason.

`log.dirs` (or `log.dir`) is the broker configuration that names the directories for partition logs. You can set more than one directory. Kafka assigns partitions to directories. This is not a RAID manager. Use a proper disk setup under the directories.

Do not delete segment files in the file system to "free disk". Use retention, or remove the topic with the admin API. Manual delete corrupts the partition.

KRaft metadata has its own storage directory in the controller configuration. Do not confuse the metadata log with `log.dirs` for topic partitions. Both need disk. Format both when you create a cluster (official format command).

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
4. Find `log.dirs` in your `server.properties` or container env. Write the path.

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

If you set both limits, Kafka can delete when either limit applies. Check the exact rule in the docs for your version.

Retention is per partition, not per topic as one blob. A hot partition hits the size limit first.

`retention.ms=-1` and a disabled size limit can keep data forever. Disk will fill. Use infinite retention only when you have a disk plan.

Consumers with a committed offset that is now deleted have a problem. `auto.offset.reset` applies (topic 4). You can lose the logical position. Monitor lag and disk.

Copies multiply disk. A replication factor of 3 stores about three times the data. Include replicas in disk plans.

Disk full is a serious failure. The broker can stop writes. Monitor disk. Retention is your main tool to free space. Do not fill the disk to 100%.

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
4. Estimate size: 10 MB/s, 2 days, 3 replicas. Write the arithmetic.

#### Medium practical tasks

1. Create a topic with `retention.ms` of a few minutes and a small segment size. Produce. Wait. Describe offsets. Write whether the earliest offset moved.
2. Create a topic with a small `retention.bytes`. Produce more than that size. Write what happens to old records.
3. Read official retention configuration. Rewrite the two limits in STE.

#### Advanced practical tasks

1. Measure time from "segment closed" to "segment gone" on a short-retention topic. Write the steps and the time.
2. Write a disk plan for a 100 MB/s topic with 7-day retention and replication factor 3. Show the arithmetic.

---

## Compaction and tombstones

**Compaction** keeps the latest record for each key in a partition. Older records with the same key can be removed. Records with different keys stay.

Use compaction for a changelog: the current state of an entity is the last value for that key. Examples: user profile, account status, stream table state.

Compaction requires keys. Null keys do not compact as an entity changelog. Use a stable key.

Compaction is not instant. A cleaner thread reads closed segments and writes compacted segments. The latest value can exist more than once until the cleaner runs.

Compaction does not make a topic a database. You cannot query by key through the Kafka protocol as you query SQL. A consumer still reads the log.

Set `cleanup.policy=compact` on the topic. Do not assume the broker default is compact. The usual default is `delete`.

A topic can use both policies: `cleanup.policy=compact,delete`. Kafka compact-cleans keys and also deletes records that are older than the retention time. Use compact+delete when you want a changelog that does not grow forever. Old keys that no one updates still leave a record in a compact-only topic.

A **tombstone** is a record with a key and a null value. On a compacted topic, a tombstone means "this key is deleted". The cleaner can remove earlier values for that key. After a tombstone retention time, the tombstone itself can go away.

Producers send a tombstone when the entity is gone. Example: user deleted. Consumers that build a table remove the key from their table when they see a null value.

If you only stop sending updates, compaction keeps the last non-null value. The key stays. You must send a tombstone to delete the key from a compact changelog.

A null value on a `delete`-only topic is just a record with a null value. It does not have the same cleanup meaning as a compact tombstone.

`delete.retention.ms` controls how long tombstones stay (official name). Readers need that window to see the delete.

### Questions

#### Theoretical questions

1. What does compaction keep?
2. Why does compaction need keys?
3. What is a tombstone?
4. What happens if you never send a tombstone for a removed entity?
5. What does `compact,delete` add that compact-only does not add?

#### Easy practical tasks

1. Write five sentences about compaction and tombstones. Use only facts from this section.
2. Make a table: Topic type, policy. Add metrics stream, user profile changelog.
3. Find `cleanup.policy` and `delete.retention.ms` in the docs. Write the meaning of each.
4. Draw a key with two values and one tombstone. Mark what a table consumer does.

#### Medium practical tasks

1. Create a compacted topic. Produce three values for key `u1` and one value for `u2`. Consume from the beginning after you wait or trigger cleanup if you can. Write the remaining records.
2. On a compacted topic, produce `u1=A`, then `u1` with a null value. Consume. Write what you see.
3. Create a topic with `cleanup.policy=compact,delete` and a short `retention.ms`. Produce a key once. Wait. Consume. Write whether the key remains.

#### Advanced practical tasks

1. Write a small table consumer (in-memory map). Apply updates and a tombstone. Confirm that the key is gone.
2. Write a one-page producer standard: when to send a tombstone versus a `UserDeleted` event on a history topic.

---

## ISR, `min.insync.replicas`, unclean leader election

The **in-sync replica set (ISR)** is the set of replicas that have caught up with the leader. The leader is always in the ISR. A follower stays in the ISR when it fetches new records within the configured lag limit.

When a follower is too slow or stops, the leader removes it from the ISR. The replica is still assigned, but it is not in-sync. `describe` output shows replicas and ISR as two lists.

For `acks=all`, the leader waits for the current ISR (and for `min.insync.replicas`). A produce that meets that wait is durable on those in-sync brokers.

Do not treat "replication factor 3" as "three in-sync copies at all times". The ISR can shrink to 1 if two followers fall behind.

`min.insync.replicas` is the minimum ISR size that the leader requires for a successful produce when the producer uses `acks=all`. If the ISR is smaller than this value, the leader rejects the produce.

Example: replication factor 3, `min.insync.replicas=2`. The cluster can lose one replica and still accept durable writes. If two replicas are out, writes with `acks=all` fail. That failure is better than a silent write to a single disk when you required two copies.

A common production pair is replication factor 3 and `min.insync.replicas=2`. A one-broker learning cluster must use `min.insync.replicas=1` and replication factor 1.

**Replication lag** is how far a follower is behind the leader. High lag causes ISR shrink. Do not fix lag by enabling unclean election.

**Unclean leader election** means Kafka can elect a leader that is not in the ISR. That replica is behind. The records that existed only on the old leader can disappear from the log. Consumers can see a rewind.

The safe default is to refuse unclean election. If no ISR member is alive, the partition stays offline until an in-sync replica returns. That outage is better than silent data loss for important data.

The configuration name is `unclean.leader.election.enable`. Set it to `false` for important clusters. KRaft does not change this danger. The controller can still elect an out-of-ISR leader if you allow it.

When the ISR is smaller than the replica list, the partition is **under-replicated**. Topic 10 covers operations for that alert.

### Questions

#### Theoretical questions

1. What is the ISR?
2. What does `min.insync.replicas` require?
3. Why is replication factor not the same as the current ISR size?
4. What is unclean leader election?
5. What data can disappear if unclean election is on?

#### Easy practical tasks

1. Describe a topic. Write Replicas and ISR for each partition.
2. Write five sentences about ISR and min ISR. Use only facts from this section.
3. Find `unclean.leader.election.enable` in the docs. Write the meaning of `true` and `false`.
4. Make a table: RF 3 and min ISR 2 versus min ISR 1. Add a durability note.

#### Medium practical tasks

1. On a multi-broker cluster, stop a follower. Describe the topic. Write the new ISR. Start the follower. Write when it returns.
2. On a three-broker cluster, set `min.insync.replicas=2`. Stop two followers. Produce with `acks=all`. Record the error.
3. Read the official warning about unclean leader election. Rewrite it in STE in six sentences.

#### Advanced practical tasks

1. Write a durability matrix: rows `acks` 1 and all, columns min ISR 1 and 2. Fill each cell with "can lose on leader death" or "blocked write" as appropriate.
2. Write a production durability standard: RF, min ISR, unclean election flag, and a ban on ZooKeeper for new clusters.

---

## KRaft controller (do not start new clusters on ZooKeeper)

The **controller** manages cluster metadata: topics, partition assignments, leader elections, and broker membership.

**KRaft (use this).** Kafka stores metadata in a Raft log on controller nodes. You format storage with a cluster id. You start one or more controllers. A production cluster uses an odd number of dedicated controllers (for example 3). A laptop cluster can use combined broker+controller processes.

KRaft does not use ZooKeeper. New Kafka versions do not start ZooKeeper. Kafka 4.0 and later remove ZooKeeper. Learn KRaft only.

**Old ZooKeeper (do not use for new work).** Old clusters stored some metadata in ZooKeeper. Tutorials from that time start a ZooKeeper process first. Ignore those steps when you build a new cluster. If you maintain an old cluster, follow the official migration guide to KRaft. Do not design new systems on ZooKeeper.

The controller is not the same as a partition leader. A partition leader handles produce and fetch for one partition. The controller assigns that leader.

When the active controller fails, the KRaft quorum elects a new controller. Metadata operations pause briefly. Partition data on brokers stays.

KRaft does not replace partition logs. Partition data still lives in `log.dirs`. KRaft does not replace ISR or replication.

Admin tools use `--bootstrap-server`. Do not use old ZooKeeper flags on new clusters.

### Questions

#### Theoretical questions

1. What does the controller manage?
2. Where does KRaft store metadata?
3. Why do production clusters use an odd number of controllers?
4. What must you not start for a new cluster?
5. How is a controller different from a partition leader?

#### Easy practical tasks

1. Write five sentences about the KRaft controller. Use only facts from this section.
2. Find the process roles in your config (`process.roles` or the equivalent). Write them.
3. Make a table: KRaft vs old ZooKeeper. Add one row: new setups.
4. Open the official KRaft documentation. Write the URL and one fact.

#### Medium practical tasks

1. In a combined-mode local cluster, find log lines that name the active controller. Write one line (paraphrase if needed).
2. Read the official migration overview (ZooKeeper to KRaft) at a high level. Write why a new learner can skip it.
3. Draw a quorum of three controllers and three brokers. Label metadata versus partition logs.

#### Advanced practical tasks

1. Run a three-controller KRaft cluster (or document a Compose file). Stop the active controller. Confirm that produce to a topic still works after the new controller is active.
2. Write a one-page "new cluster" install sheet that never mentions a ZooKeeper start command except in a footnote "do not do this".

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe the life of a record from append to the active segment to delete or compact.
2. How do segment roll, retention, and the cleaner thread depend on each other?
3. Describe a durable produce: leader, ISR, `min.insync.replicas`, `acks=all`.
4. Why is unclean leader election a data-loss feature, not a high-availability feature?
5. What two storage systems exist in a KRaft cluster (metadata versus partition logs)?

#### Easy practical tasks

1. Create `store.review` with a short time retention. Produce five records. Describe the topic config. Highlight cleanup policy and retention keys.
2. Write a one-page cheat sheet: segment, retention, compact, tombstone, ISR, min ISR, unclean election, KRaft controller.
3. Describe a replicated topic. Highlight leader, replicas, ISR, and min ISR if shown.
4. Draw compact+delete on a time axis for one unused key and one hot key.

#### Medium practical tasks

1. Write a script that creates a compacted topic, produces keys, produces a tombstone, and consumes from the beginning.
2. Create a topic with RF 3 and min ISR 2 on a three-broker cluster if you have one. Produce. Stop one broker. Produce again. Describe ISR after each step.
3. Map each subsection to an official configuration key.

#### Advanced practical tasks

1. Build a small "table from compacted topic" program. Restart it. Confirm that it rebuilds from the log. Then send a tombstone and confirm the key is gone.
2. Write a production storage standard: defaults for event streams versus changelogs, disk formula with replicas, unclean election off, KRaft only, and a ban on manual file delete.
