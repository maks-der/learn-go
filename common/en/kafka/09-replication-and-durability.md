# 9. Replication and Durability

## Description

Replication copies each partition to more than one broker. Durability is the chance that a committed record stays after a broker fails. Kafka uses a leader, followers, and an in-sync replica set (ISR).

This topic covers ISR, `min.insync.replicas`, unclean leader election, replication lag, and the KRaft controller. Complete this topic after topics 2 and 4. Prefer KRaft for all new clusters. Do not use ZooKeeper for new work. This section mentions ZooKeeper only as old history.

Use one term for each concept. A replica is a copy of a partition. The ISR is the set of replicas that are up to date. The controller is the component that manages cluster metadata and partition leaders. In new Kafka, the controller is a KRaft quorum.

---

## ISR (in-sync replicas)

The in-sync replica set (ISR) is the set of replicas that have caught up with the leader. The leader is always in the ISR. A follower stays in the ISR when it fetches new records within the configured lag limit.

When a follower is too slow or stops, the leader removes it from the ISR. The replica is still assigned, but it is not in-sync. `describe` output shows replicas and ISR as two lists.

For `acks=all`, the leader waits for the current ISR (and for `min.insync.replicas`). A produce that meets that wait is durable on those in-sync brokers.

A follower that catches up can join the ISR again. Replication lag (next section) explains the delay.

ISR is per partition. Partition 0 and partition 1 have separate ISR lists.

Do not treat "replication factor 3" as "three in-sync copies at all times". The ISR can shrink to 1 if two followers fall behind. Then `acks=all` only waits for one broker unless `min.insync.replicas` blocks the write.

### Questions

#### Theoretical questions

1. What is the ISR?
2. Who is always in the ISR?
3. When does a follower leave the ISR?
4. How does `acks=all` use the ISR?
5. Why is replication factor not the same as the current ISR size?

#### Easy practical tasks

1. Describe a topic. Write Replicas and ISR for each partition.
2. Write five sentences about ISR. Use only facts from this section.
3. Make a table: Replication factor 3, ISR size 3 versus ISR size 1. Add a durability note.
4. Draw a leader and two followers. Mark one follower out of ISR.

#### Medium practical tasks

1. On a multi-broker cluster, stop a follower. Describe the topic. Write the new ISR.
2. Start the follower again. Wait. Describe. Write when it returns to the ISR.
3. Find official ISR documentation. Write the configuration that controls lag (for example `replica.lag.time.max.ms`) in your own words.

#### Advanced practical tasks

1. Measure time to leave ISR and time to rejoin after you stop and start a follower. Write both times.
2. Write a one-page note: how you monitor under-replicated partitions (ISR smaller than the replica list).

---

## `min.insync.replicas`

`min.insync.replicas` is the minimum ISR size that the leader requires for a successful produce when the producer uses `acks=all`. If the ISR is smaller than this value, the leader rejects the produce.

Example: replication factor 3, `min.insync.replicas=2`. The cluster can lose one replica and still accept durable writes. If two replicas are out, writes with `acks=all` fail. That failure is better than a silent write to a single disk when you required two copies.

`min.insync.replicas=1` with `acks=all` only guarantees one in-sync copy. That copy is the leader. You can still lose data if that broker dies before followers catch up.

Set this value on the broker default and override it on important topics. The producer `acks` setting does not replace this broker check.

A common production pair is replication factor 3 and `min.insync.replicas=2`. A one-broker learning cluster must use `min.insync.replicas=1` and replication factor 1.

### Questions

#### Theoretical questions

1. What does `min.insync.replicas` require?
2. When does the leader reject a produce because of this setting?
3. Why is `min.insync.replicas=2` a common choice with three replicas?
4. What does `min.insync.replicas=1` still allow?
5. Does the producer `acks` key replace `min.insync.replicas`?

#### Easy practical tasks

1. Write four sentences about `min.insync.replicas`.
2. Make a table: RF, min ISR, brokers that can fail and still write with `acks=all`.
3. Find the default `min.insync.replicas` on your cluster.
4. List two topics that need min ISR 2 and one learning topic that cannot.

#### Medium practical tasks

1. On a three-broker cluster, set `min.insync.replicas=2`. Stop two followers (or two brokers that hold replicas). Produce with `acks=all`. Record the error.
2. Set `min.insync.replicas=1`. Repeat. Write the difference.
3. Read official documentation for `min.insync.replicas`. Rewrite the produce rule in STE.

#### Advanced practical tasks

1. Write a durability matrix: rows `acks` 1 and all, columns min ISR 1 and 2. Fill each cell with "can lose on leader death" or "blocked write" as appropriate.
2. Write a topic create standard for production and for laptop. Include RF and min ISR.

---

## Unclean leader election (danger)

Unclean leader election means Kafka can elect a leader that is not in the ISR. That replica is behind. The records that existed only on the old leader can disappear from the log. Consumers can see a rewind. New readers never see those records.

The safe default is to refuse unclean election. If no ISR member is alive, the partition stays offline until an in-sync replica returns. Writes stop. That outage is better than silent data loss for important data.

The configuration name is `unclean.leader.election.enable`. Set it to `false` for important clusters. Only consider `true` when availability of stale data is more important than correctness, and the business accepts loss. Document that choice.

Unclean election is not a repair tool. Do not enable it to "make the topic work" during a failure without a written decision.

KRaft does not change this danger. The controller can still elect an out-of-ISR leader if you allow it.

### Questions

#### Theoretical questions

1. What is unclean leader election?
2. What data can disappear?
3. What is the safe default when no ISR member is alive?
4. Why is an outage sometimes better than an election?
5. Does KRaft remove this risk?

#### Easy practical tasks

1. Write five sentences about unclean leader election. Use only facts from this section.
2. Find `unclean.leader.election.enable` in the docs. Write the meaning of `true` and `false`.
3. Make a table: Setting false vs true. Add availability and data loss.
4. Draw an ISR of {1} and a dead replica 2. Mark what unclean election of 2 would drop.

#### Medium practical tasks

1. Read the official warning about unclean leader election. Rewrite it in STE in six sentences.
2. Check the value on your local cluster. Write it.
3. Write a short incident choice: 40 minutes of partition offline versus possible loss of the last 200 ms of writes. Pick one and defend it.

#### Advanced practical tasks

1. In a lab only, enable unclean election, create a failure where the only ISR member dies and a lagging replica remains, and observe offsets. Restore safe config after the lab. Write the evidence of loss if you see it.
2. Write a one-page policy that forbids unclean election on payment topics and names who can change the flag.

---

## Replication lag

Replication lag is how far a follower is behind the leader. You can express lag as offsets (leader log end minus follower log end) or as time.

Lag grows when the follower is slow: disk, network, CPU, or a long garbage-collection pause. Lag also grows when the leader takes a large produce burst.

High lag causes ISR shrink. Then durability and `min.insync.replicas` interact. High lag also delays the moment when a follower is a safe leader.

Consumers that read from the leader do not wait for all followers. They wait for the high watermark (committed data). Topic 20 covers high watermark versus log end offset. For this topic: a record is safe for `acks=all` when the ISR has it.

Monitor under-replicated partitions and follower fetch metrics. A permanent lag on one broker is a hardware or config problem.

Do not fix lag by enabling unclean election.

### Questions

#### Theoretical questions

1. What is replication lag?
2. What causes lag to grow?
3. How does lag relate to ISR membership?
4. Why is a high-lag follower a poor leader candidate?
5. Why is unclean election the wrong fix for lag?

#### Easy practical tasks

1. Write four sentences about replication lag.
2. Make a table: Cause of lag, what you check.
3. Find a metric or command that shows under-replicated partitions.
4. Draw leader LEO and follower LEO. Mark the lag.

#### Medium practical tasks

1. Produce a large burst on a three-replica topic. Watch ISR or lag if you have tools (JMX, describe, or a UI). Write what you see.
2. Throttle a follower (small Docker CPU limit) if you can. Write whether the follower leaves ISR.
3. Read official replica fetch configuration. Write two keys that affect how fast a follower catches up.

#### Advanced practical tasks

1. Measure catch-up time after a follower was down for N minutes of produce load. Write N, data volume, and catch-up time.
2. Write a one-page runbook for under-replicated partitions: check disk, network, broker logs, then add capacity.

---

## Controller (KRaft vs old ZooKeeper)

The controller manages cluster metadata: topics, partition assignments, leader elections, and broker membership.

**KRaft (use this).** Kafka stores metadata in a Raft log on controller nodes. You format storage with a cluster id. You start one or more controllers. A production cluster uses an odd number of dedicated controllers (for example 3). A laptop cluster can use combined broker+controller processes.

KRaft does not use ZooKeeper. New Kafka versions do not start ZooKeeper. Kafka 4.0 and later remove ZooKeeper. Learn KRaft only.

**Old ZooKeeper (do not use for new work).** Old clusters stored some metadata in ZooKeeper. Tutorials from that time start a ZooKeeper process first. Ignore those steps when you build a new cluster. If you maintain an old cluster, follow the official migration guide to KRaft. Do not design new systems on ZooKeeper.

The controller is not the same as a partition leader. A partition leader handles produce and fetch for one partition. The controller assigns that leader.

When the active controller fails, the KRaft quorum elects a new controller. Metadata operations pause briefly. Partition data on brokers stays.

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
2. Read the official migration overview (ZK to KRaft) at a high level. Write why a new learner can skip it.
3. Draw a quorum of three controllers and three brokers. Label metadata versus partition logs.

#### Advanced practical tasks

1. Run a three-controller KRaft cluster (or document a Compose file). Stop the active controller. Confirm that produce to a topic still works after the new controller is active.
2. Write a one-page "new cluster" install sheet that never mentions a ZooKeeper start command except in a footnote "do not do this".

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe a durable produce: leader, ISR, `min.insync.replicas`, `acks=all`.
2. How do lag, ISR shrink, and blocked writes protect data?
3. Why is unclean leader election a data-loss feature, not a high-availability feature?
4. What two storage systems exist in a KRaft cluster (metadata vs partition logs)?
5. A blog starts ZooKeeper and two brokers. What do you change for a new install?

#### Easy practical tasks

1. Describe a replicated topic. Highlight leader, replicas, ISR, and min ISR if shown.
2. Write a one-page cheat sheet: ISR, min.insync.replicas, unclean election, lag, KRaft controller.
3. Write your cluster's RF default and min ISR default.
4. Draw a fail of one follower and a fail of the leader (ISR still has a member).

#### Medium practical tasks

1. Create a topic with RF 3 and min ISR 2. Produce. Stop one broker. Produce again. Describe ISR after each step.
2. Find official pages for ISR, min.insync.replicas, unclean election, and KRaft. Write one sentence from each in STE.
3. Compare a one-node laptop durability story with a three-node production story in six sentences.

#### Advanced practical tasks

1. Run a controlled leader failover (stop the leader broker). Measure time until produce succeeds again. Write ISR and leader before and after.
2. Write a production durability standard: RF, min ISR, unclean election flag, controller count, and a ban on ZooKeeper for new clusters.
