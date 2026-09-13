# 15. Operations

## Description

Operations keep a Kafka cluster healthy. You watch replicas, disk, memory, and restarts. You move partitions when brokers or disks change. You use KRaft for metadata. You do not use ZooKeeper for new clusters.

This topic covers under-replicated partitions, disk use and bursty deletes, JVM and page cache, `kafka-reassign-partitions`, Cruise Control awareness, rolling restarts, and rack awareness. Complete this topic after topics 8, 9, and 10.

Use one term for each concept. An under-replicated partition (URP) is a partition whose ISR is smaller than the replica list. Page cache is the operating-system cache of disk pages. A rolling restart stops one broker at a time. Rack awareness places replicas on different failure domains. Cruise Control is a tool that can propose and run partition moves.

---

## ISR, under-replicated partitions

Topic 9 defines the ISR. The replica list is the assigned copies. The ISR is the copies that are in sync. When the ISR is smaller than the replica list, the partition is **under-replicated**.

A URP alert means a follower is slow, stopped, or cut off. Durability is lower until the follower returns to the ISR. `acks=all` and `min.insync.replicas` can start to reject writes if the ISR shrinks too far (topic 9).

Monitor URP count as a cluster metric. A short URP during a rolling restart can be normal. A URP that stays after the broker is up is an incident.

`kafka-topics --describe` shows Replicas and ISR. Tools and metrics use names such as `UnderReplicatedPartitions`.

Do not ignore URP because "clients still produce". Producers can succeed on a smaller ISR. The risk is loss if more brokers fail.

KRaft does not remove URP. KRaft manages metadata and leader election. Partition data still replicates between brokers.

### Questions

#### Theoretical questions

1. What is an under-replicated partition?
2. How do you see URP in a topic describe?
3. When can a short URP be expected?
4. Why can producers still succeed while URP is greater than zero?
5. Does KRaft prevent URP?

#### Easy practical tasks

1. Describe a topic. Write Replicas and ISR for each partition. Say if any partition is under-replicated.
2. Write five sentences about URP. Use only facts from this section.
3. Make a table: ISR size versus replica count, URP yes or no, durability note.
4. Find the metric name for URP in official monitoring docs.

#### Medium practical tasks

1. On a multi-broker KRaft cluster, stop one broker. Describe a replicated topic. Write the URP state. Start the broker. Write when ISR recovers.
2. Compare URP with "offline partitions" (no leader) in the docs. Write three differences.
3. Draw a timeline: rolling restart of one broker, URP rises, URP returns to zero.

#### Advanced practical tasks

1. Measure time from broker stop to URP and time from start to ISR full. Write both times and the topic settings.
2. Write a one-page alert standard: URP warning versus critical, and when a restart URP is ignored.

---

## Disk usage and bursty deletes

Kafka stores partition logs in `log.dirs` (topic 8). Disk fills when produce rate is higher than retention delete, or when you add partitions and topics.

Retention delete and compaction run on closed segments. Many segments can become eligible at the same time. Then Kafka deletes a large amount of data in a short period. That is a **bursty delete**. Disk use drops in a step. Disk I/O can spike. Consumers and page cache can feel that spike.

Causes of a burst: a retention change, a time boundary that many segments share, a compact run after a large load.

Do not delete segment files by hand. Do not fill the disk to 100%. Kafka and the OS need free space. Plan a headroom percent. Alert before the disk is full.

Compaction can keep keys forever if keys never get tombstones. A compacted topic can grow without a time delete unless you use compact+delete (topic 8).

KRaft metadata logs also use disk on controller nodes. They are smaller than data logs in most clusters. Still monitor controller disks.

After a bursty delete, `df` looks better. Confirm that `log.dirs` is the volume that you monitor.

### Questions

#### Theoretical questions

1. Where does Kafka store partition logs?
2. What is a bursty delete?
3. Why can many segments expire at the same time?
4. Why must you not delete `.log` files by hand?
5. Why can a compacted topic grow for a long time?

#### Easy practical tasks

1. Write four sentences about disk and retention. Use only facts from this section.
2. On a lab broker, find `log.dirs` and write the path and the free space.
3. Make a table: Event, disk effect (grow, drop, spike I/O).
4. Describe a topic. Write retention bytes or hours.

#### Medium practical tasks

1. Create a small-retention topic. Produce enough data to roll segments. Watch disk and segment files before and after delete.
2. Change retention to a shorter time on a lab topic. Observe whether delete happens at once or over time. Write what you see.
3. Find official notes on log retention and deletion. Write three facts about when a segment can be deleted.

#### Advanced practical tasks

1. Produce a burst, then let retention delete. Record disk percent, I/O, and consumer lag if any during the delete.
2. Write a capacity plan: daily produce bytes, retention, replica factor, headroom percent, and controller disk on KRaft.

---

## JVM and page cache

The Kafka broker runs on the JVM. The JVM heap holds metadata, request objects, and some caches. Partition log data is on disk. The operating system keeps recent disk pages in **page cache** (RAM that is not the JVM heap).

Kafka produce and consume can use page cache. A consumer that reads recent data often reads from RAM, not from the disk spindle or SSD firmware path in the slow case.

If you give almost all RAM to the JVM heap, page cache becomes small. Throughput can drop. A common practice is a moderate heap and a large remaining RAM for page cache. Follow current Apache Kafka operations guidance for heap size. Do not copy a heap size from an old blog without a check.

GC pauses on a large heap can delay requests. That delay can cause follower lag and URP.

KRaft controllers are also JVM processes when they run in Kafka. Combined broker+controller processes (`process.roles`) share the same JVM. Size the heap for that combined role. Isolated controllers still need a stable heap.

Monitor heap use, GC time, and OS page cache or available memory. Monitor disk I/O. A sudden I/O rise with a full cache miss pattern can mean readers left the hot set.

Do not tune JVM flags that you do not understand. Change one flag in a lab first.

### Questions

#### Theoretical questions

1. What lives on the JVM heap versus in the page cache for a broker?
2. Why can a too-large heap hurt Kafka?
3. How can a long GC pause cause URP?
4. Why does a combined KRaft broker+controller process need heap planning?
5. Why must you not copy old heap flags without a check?

#### Easy practical tasks

1. Write five sentences about JVM and page cache. Use only facts from this section.
2. Find the heap settings in your lab broker start script or environment. Write the values.
3. Make a table: Memory region, typical content, too small effect.
4. Find official or current Apache notes on broker memory. Write one recommended idea in your own words.

#### Medium practical tasks

1. Record broker RSS, heap use, and disk I/O during a consume of old data versus new data (lab). Write the difference.
2. Read a GC log or metric for a lab broker during load. Write pause times if you can enable them safely.
3. Compare a combined KRaft process with a controller-only process in the docs. Write three memory notes.

#### Advanced practical tasks

1. Run two lab configs: large heap versus moderate heap with more free RAM. Measure consume throughput of a hot topic. Write the result and limits of the test.
2. Write a memory standard: heap range, what you monitor, and a ban on ad-hoc JVM flags in production.

---

## `kafka-reassign-partitions`

`kafka-reassign-partitions` is the command-line tool that moves partition replicas between brokers. You use it when you add a broker, remove a broker, or balance disk.

The usual path:

1. Generate a reassignment JSON (current assignment plus a proposed move).
2. Edit or accept the proposed JSON.
3. Execute the reassignment.
4. Verify that the move finished.

On Windows the file is `.bat`. On macOS and Linux the file is `.sh`. Some installs put the tool on `PATH` without a suffix.

Use `--bootstrap-server`. Do not use old ZooKeeper connection flags on a KRaft cluster.

A move copies bytes over the network. Disk and network load rise. Throttle settings exist. A large move during peak traffic can cause lag and URP.

Do not move all partitions at once on a large cluster without a throttle and a plan. Do not remove a broker from the cluster before its replicas have left it.

The Admin API can also alter assignments. Cruise Control (next section) can generate safer plans. The manual tool is the built-in mechanism that you must understand.

Example shape:

```text
kafka-reassign-partitions --bootstrap-server localhost:9092 --reassignment-json-file move.json --execute
kafka-reassign-partitions --bootstrap-server localhost:9092 --reassignment-json-file move.json --verify
```

### Questions

#### Theoretical questions

1. What does `kafka-reassign-partitions` do?
2. What are the generate, execute, and verify steps?
3. Which address flag do you use on KRaft?
4. Why does a reassignment increase network and disk load?
5. Why must you not decommission a broker before replicas move away?

#### Easy practical tasks

1. Run the tool with `--help`. Write five flags.
2. Write four sentences about reassignment. Use only facts from this section.
3. Make a table: Step, input file, danger if skipped.
4. Find official reassignment documentation. Write the JSON field names that you see (topics, partitions, replicas).

#### Medium practical tasks

1. On a three-broker KRaft lab, create a topic with replicas on two brokers. Generate a move that adds the third broker. Execute and verify. Describe the topic after.
2. Find throttle-related flags or configs in the docs. Write how you would limit a move.
3. Export the current assignment JSON for one topic. Keep it as a backup note.

#### Advanced practical tasks

1. Move a partition off a broker, then remove that broker from the lab. Write every command and the URP timeline.
2. Write a change procedure: who approves a JSON, how you throttle, how you verify, and KRaft-only flags.

---

## Cruise Control (awareness)

Cruise Control is an open-source system (from LinkedIn and the community) that analyzes a Kafka cluster and can propose partition reassignments. It can balance disk, network, and leadership. It can run anomaly detection (for example, self-healing after a broker failure).

This handbook requires **awareness**, not a full install. You must know that the tool exists. You must know that it talks to Kafka and to metrics. You must know that it still uses the same reassignment idea as `kafka-reassign-partitions`.

Cruise Control is not part of the Apache Kafka tarball by default. You deploy it as a separate service. It needs metrics (often from a reporter on the brokers).

Do not let two automation systems execute conflicting reassignments at the same time. If you use Cruise Control, make it the one executor, or lock manual moves.

Cruise Control does not replace KRaft. It does not replace monitoring of URP and disk. It does not need ZooKeeper on current Kafka.

Managed Kafka services can have their own balancer. The idea is the same: propose a move, execute, verify.

### Questions

#### Theoretical questions

1. What problem does Cruise Control solve?
2. Is Cruise Control included in Apache Kafka by default?
3. Why must only one system execute reassignments?
4. What does Cruise Control still depend on in the cluster?
5. Does Cruise Control replace KRaft?

#### Easy practical tasks

1. Write five sentences about Cruise Control. Use only facts from this section.
2. Make a table: Manual `kafka-reassign-partitions`, Cruise Control. Add rows for propose, execute, schedule.
3. Open the Cruise Control project README or docs. Write three goals that the page names.
4. List two managed-Kafka features that are similar (from vendor docs you already use).

#### Medium practical tasks

1. Read one official or project page on Cruise Control goals (disk, network, leader). Rewrite them in STE.
2. Draw: metrics → Cruise Control → reassignment JSON or API → brokers.
3. Write a policy sentence: when a human may run `kafka-reassign-partitions` if Cruise Control is installed.

#### Advanced practical tasks

1. Optional lab: run Cruise Control against a multi-broker KRaft cluster. Take one proposal. Write whether you executed it.
2. Write a one-page awareness note for your team: what Cruise Control is, what it is not, and the no-ZooKeeper rule.

---

## Rolling restarts

A **rolling restart** restarts one broker at a time so that the cluster stays available. Clients use other leaders. Followers catch up. Then you restart the next broker.

Typical steps for each broker:

1. Confirm URP is zero (or understood).
2. Stop the broker process.
3. Wait until leadership moves and the cluster is stable.
4. Start the broker.
5. Wait until that broker is in the ISR for its replicas.
6. Go to the next broker.

On KRaft, some nodes are controllers. If you run combined roles, each restart is a broker and a controller voter. Restart combined nodes more carefully. Keep the controller quorum available. Do not stop a majority of controllers at the same time.

If you run dedicated controllers, you can roll brokers and controllers with separate rules. Always keep a Raft majority.

Do not set `unclean.leader.election.enable` to true to "make restarts easier" (topic 9). Do not restart all brokers at once.

After a configuration change, a rolling restart applies the new broker settings that need a restart. Some dynamic configs do not need a restart (topic 10).

### Questions

#### Theoretical questions

1. What is a rolling restart?
2. Why do you wait for ISR after each start?
3. Why must you not stop a majority of KRaft controllers at once?
4. Why is unclean leader election the wrong tool for restarts?
5. Which configs need a restart and which can be dynamic?

#### Easy practical tasks

1. Write the six-step list for one broker in your own words.
2. Make a table: Combined KRaft node, dedicated controller. Add a restart rule for each.
3. Find official rolling upgrade or restart notes. Write three facts.
4. Draw three brokers. Mark one stopped. Mark where leaders go.

#### Medium practical tasks

1. On a three-node KRaft lab, perform a rolling restart. After each node, describe a topic and write ISR.
2. Time the full roll. Write how long URP was greater than zero.
3. Compare a rolling restart with a rolling upgrade in the official upgrade guide. Write five shared steps.

#### Advanced practical tasks

1. Change one broker config that needs a restart. Roll the cluster. Confirm the new value with `kafka-configs` or describe.
2. Write a restart runbook: checks before, KRaft quorum rule, URP wait, client bootstrap note, no ZooKeeper.

---

## Rack awareness

**Rack awareness** places replicas so that they do not all live in one failure domain. A rack can be a rack, an availability zone, or a site label. You set `broker.rack` on each broker.

When Kafka assigns replicas, it tries to spread replicas of a partition across racks. Then one rack failure does not take all copies.

Rack awareness is not a full geo-replication design. Topic 18 covers multi-cluster and geo. Rack awareness is inside one cluster.

If all brokers have the same rack value, there is no spread. If you have three replicas and two racks, two replicas can share a rack. Plan replica factor and rack count together.

Clients can use rack-aware replica selection so that a consumer fetches from a replica in the same rack when it is in the ISR. That can reduce cross-zone network cost. Correctness still uses the leader for produces.

KRaft voters also have locations. Spread controllers across racks when you can. Loss of one rack must not remove the Raft majority if you have enough voters.

Do not invent rack names that do not match real failure domains.

### Questions

#### Theoretical questions

1. What does `broker.rack` represent?
2. How does replica assignment use racks?
3. Why is rack awareness not the same as multi-cluster geo?
4. How can a consumer use rack awareness?
5. Why must KRaft controllers also spread across racks when possible?

#### Easy practical tasks

1. Write five sentences about rack awareness. Use only facts from this section.
2. Make a table: Three brokers in three racks, RF=3. Where can replicas of one partition go?
3. Find `broker.rack` in official configuration docs. Write the meaning in your own words.
4. Draw two racks, three brokers, one partition with RF=3. Mark a bad assignment (all in one rack) and a good assignment.

#### Medium practical tasks

1. Set different `broker.rack` values on a three-broker KRaft lab. Create a topic with RF=3. Describe replica placement.
2. Find replica selector or client rack settings in the docs. Write the key names.
3. Write what happens if you create the topic before you set racks. Do you need a reassignment?

#### Advanced practical tasks

1. Simulate a rack loss (stop all brokers in one rack). Write which partitions stay online and the URP or offline state.
2. Write a placement standard: rack names = AZ names, RF versus rack count, controller spread, KRaft only.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do URP, disk bursts, and JVM pauses relate to each other in one incident story?
2. When do you use `kafka-reassign-partitions` by hand and when do you use a balancer such as Cruise Control?
3. What extra rule does KRaft add to a rolling restart that a "brokers only" view misses?
4. How do rack awareness and reassignment work together when you add a broker in a new rack?
5. What from topics 8 and 9 must you re-check after every large operations change?

#### Easy practical tasks

1. Write a one-page operations cheat sheet: URP, disk, heap versus page cache, reassign, Cruise Control, roll, racks.
2. Draw a KRaft three-node cluster with racks, `log.dirs`, and a client bootstrap list.
3. Bookmark official operations, reassignment, and KRaft pages.
4. List the commands named in this topic and one purpose sentence for each.

#### Medium practical tasks

1. Run a lab day: describe URP, check disk, roll one broker, reassign one partition. Write a timeline.
2. Map each subsection to one official URL.
3. Write a monitoring list: five metrics and the action you take when each is high.

#### Advanced practical tasks

1. Add a broker with a new rack, reassign, roll all nodes, and record URP and disk through the change. Use KRaft only.
2. Write a production operations standard: alerts, restart quorum rules, reassignment ownership, memory policy, no ZooKeeper.
