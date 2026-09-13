# 2. Topics, Partitions, and Records

## Description

Kafka has a small set of core objects. You must know each name before you write producers and consumers. This topic defines topic, partition, offset, broker, cluster, record, replica, and leader. It also explains order, keys, and partition count.

Complete this topic after topic 1. Complete this topic before topic 3. Use one term for each concept. A topic is a named log. A partition is an ordered slice of that log. An offset is a position in one partition. A record is the stored unit. A broker is a Kafka server process. A cluster is a set of brokers. A replica is a copy of a partition. The leader is the replica that handles reads and writes.

You can add partitions to a topic. You cannot shrink the partition count in Apache Kafka. Choose the count with care. Use KRaft for all new clusters.

---

## Topic, partition, offset, broker, cluster

A **topic** is a named, append-only log of records. Producers write to a topic by name. Consumers subscribe to a topic by name. The name is a string. Use a clear, stable name. Example: `orders.placed`.

A topic has a partition count. A topic has a replication factor. You set those values when you create the topic, unless the broker defaults apply. You can add partitions later. You cannot shrink the partition count in Apache Kafka.

A topic has configuration. Examples: retention time, cleanup policy, compression type. You can set defaults on the broker. You can override values on one topic. Topic 6 and topic 7 cover storage and admin.

A **partition** is an ordered sequence of records inside a topic. Kafka appends each new record to one partition. Records in one partition have increasing offsets.

Each partition has a number. The numbers start at 0. A topic with three partitions has partitions 0, 1, and 2.

An **offset** is the position of a record in one partition. The first record in a new partition often has offset 0. The next append gets the next offset. Offset is not a clock. Offset is not unique across partitions. Partition 0 offset 5 and partition 1 offset 5 are two different records.

A **broker** is one Kafka server process. The broker stores partition replicas. The broker serves produce and fetch requests. Clients find brokers through a bootstrap address.

A **cluster** is a set of brokers that share metadata. New clusters use KRaft. The KRaft controller quorum stores topic and broker metadata. You do not add ZooKeeper to a new cluster.

Internal topics exist for Kafka itself. Examples: `__consumer_offsets` and transaction topics. Do not write application events to those names. Topic 7 describes them.

A topic is not a folder of files that you edit. Kafka stores topic data as log segments on the brokers that hold the partitions. Topic 6 covers segments.

### Questions

#### Theoretical questions

1. What is a topic?
2. What is a partition?
3. What is an offset, and why is it not unique across partitions?
4. What is a broker, and what is a cluster?
5. Can you reduce the number of partitions of a topic in Apache Kafka?

#### Easy practical tasks

1. Write five topic names for a library system. Use one naming style.
2. Create a topic `core.topic` with three partitions. Describe it. Write the partition count, the replication factor, and the partition numbers.
3. Make a two-column table: "Object" and "What it means". Add rows for topic, partition, offset, broker, and cluster.
4. List three topic names that are poor (too vague or tied to one consumer). Rewrite each name.

#### Medium practical tasks

1. Create two topics with different partition counts. Describe both. Write how `kafka-topics --describe` shows each partition.
2. Find the official documentation page that defines a topic and a partition. Write each official idea in your own words.
3. Try to create a topic with a name that contains a space or a forbidden character. Record the error. Create a valid name.

#### Advanced practical tasks

1. Design a topic list for an online course platform (enroll, watch, complete). For each topic, set a partition count and a retention reason.
2. On a multi-broker KRaft cluster, create a topic with replication factor 3. Describe leaders and replicas. Write which broker holds which role.

---

## Record: key, value, timestamp, headers

A **record** is the unit that Kafka stores. A record has a key, a value, a timestamp, and optional headers. The key and the value are bytes or null. Brokers do not parse JSON or Avro. Topic 7 covers serialization and schemas.

The **key** selects the partition when you use the default partitioner and the key is not null. The key is also the identity for compaction. Topic 6 covers compaction. Use a stable key when records for one entity must stay in order.

The **value** is the payload. Example: the fields of an order. A null value on a compacted topic is a tombstone. Topic 6 defines tombstones.

The **timestamp** is a time that Kafka stores with the record. The producer can set create time. The broker can set log append time. Configuration chooses which time Kafka uses. Do not treat the timestamp as the only order. The append order in the partition is the order.

**Headers** are optional key-value pairs on the record. Teams use headers for a content type, a trace id, or an event id. Headers are bytes. Consumers must agree on header names. Do not hide the only copy of important business data in a header that some consumers ignore.

A record also has a topic name, a partition number, and an offset after Kafka appends it. The producer send result returns partition and offset when the send succeeds.

Do not mix an event and a record. An event is the fact. A record is the stored unit. One event becomes one record when the producer writes it.

### Questions

#### Theoretical questions

1. What four parts of a record does this section name?
2. What type are the key and the value on the broker?
3. What is the difference between create time and log append time?
4. What are headers for?
5. How is an event different from a record?

#### Easy practical tasks

1. Write five sentences about a record. Use only facts from this section.
2. Make a table: Part, type or meaning, example. Add key, value, timestamp, header.
3. Produce a line with the console producer. Consume it. Write which parts you can see in the console output.
4. Draw one record box. Label key, value, timestamp, and one header.

#### Medium practical tasks

1. Write a small program that produces a record with a key, a value, and one header. Print the partition and offset from the send result.
2. Find the official record or `ProducerRecord` documentation for your client. List the fields that you can set.
3. Produce two records with the same key and different values. Consume from the beginning. Write the order and the offsets.

#### Advanced practical tasks

1. Compare create-time and log-append-time on a test topic if your install lets you set the timestamp type. Write how the stored timestamps differ.
2. Write a one-page contract for a topic: key meaning, value fields, required headers, and timestamp type.

---

## Replica and leader

A **replica** is a copy of a partition on one broker. The replication factor is the number of replicas that you request when you create the topic. A learning cluster with one broker uses replication factor 1. A production cluster often uses replication factor 3.

The **leader** is the replica that handles writes and the usual reads for that partition. **Followers** copy the leader. Producers send to the leader. Consumers read from the leader in the common case.

If the leader broker stops, the controller elects a new leader from the in-sync replicas. Topic 6 defines the in-sync replica set (ISR) and unclean leader election.

Each partition has its own leader. Partition 0 can have a leader on broker 1. Partition 1 can have a leader on broker 2. That spread lets the cluster share load.

KRaft does not replace replicas. KRaft stores metadata and elects leaders. Partition data still lives on brokers in `log.dirs`.

Do not treat "replication factor 3" as "three copies that are always in sync". Followers can fall behind. Topic 6 explains ISR size and `min.insync.replicas`.

### Questions

#### Theoretical questions

1. What is a replica?
2. What does the leader do?
3. What do followers do?
4. Who elects a new leader when the leader broker stops?
5. Does KRaft store the partition log instead of replicas?

#### Easy practical tasks

1. Describe a topic. Write the leader, the replica list, and the ISR for each partition if the tool shows them.
2. Write four sentences about replica and leader. Use only facts from this section.
3. Make a table: Replication factor 1 versus 3. Add rows for brokers needed and failure tolerance at a high level.
4. Draw one partition with one leader and two followers.

#### Medium practical tasks

1. On a multi-broker KRaft cluster, create a topic with replication factor 3. Stop the leader broker if you can. Describe the topic. Write the new leader.
2. Find official documentation on replicas and leaders. Write three facts in your own words.
3. Compare a one-broker lab topic with a three-broker topic in six sentences. Cover leader and copies.

#### Advanced practical tasks

1. Write a one-page note: when you use replication factor 1 (laptop) and when you use 3 (production). Include KRaft controller count as a separate choice.
2. Measure time from leader stop to a successful produce on the new leader. Write the topic settings and the time.

---

## Ordering is per partition, not per topic

Kafka guarantees order only inside one partition. If record A is appended before record B in the same partition, a consumer of that partition sees A before B.

Kafka does not guarantee order across partitions. Record A in partition 0 and record B in partition 1 can arrive at consumers in any interleaving.

If you need a total order for a stream, use one partition. That choice removes parallel consumers in one group. Topic 4 covers the one-active-consumer rule. If you need order per entity, use a key so that all records for that entity go to the same partition.

Do not assume that "the topic is ordered" because you produced in a sequence from one thread. The default partitioner can send consecutive records with different keys to different partitions.

Time on the producer machine is not the order. The order is the append order on the leader of that partition.

KRaft leader election does not reorder a committed log. Unclean leader election can drop records. Topic 6 covers that danger. Keep unclean election off for important data.

### Questions

#### Theoretical questions

1. What order does Kafka guarantee?
2. What order does Kafka not guarantee?
3. How do you get a total order for one topic?
4. How do you get order per entity?
5. Why is produce order in one thread not the same as topic order?

#### Easy practical tasks

1. Write four sentences that describe Kafka order. Use only facts from this section.
2. Draw two partitions. Place events for user 1 and user 2. Mark a pair of events that a consumer can see "out of topic order".
3. Make a table: Need, one partition, same key, different keys.
4. List three event types that need per-entity order and two that do not.

#### Medium practical tasks

1. Create a topic with three partitions. Produce ten records with different keys from one thread. Consume. Write whether the consume order matches the produce order.
2. Produce three records with the same key. Consume. Confirm that those three stay in produce order.
3. Find official text about partition order. Rewrite it in STE in four sentences.

#### Advanced practical tasks

1. Write a one-page design: order for a customer balance versus order for independent click events. Name the key and the partition count for each.
2. Explain in one page why a global sort by timestamp is not a Kafka order guarantee. Use clock skew as one reason.

---

## Keys and the default partitioner

The **partitioner** is the producer logic that selects a partition. The **default partitioner** uses the key when the key is not null. Records with the same key go to the same partition (for a fixed partition count). That rule gives per-key order.

If the key is null, the default partitioner spreads records across partitions. Current clients often use a sticky strategy: they send a batch to one partition, then they pick another partition. Null-key records do not stay in one entity order.

You can set the partition number in the producer record. Then the partitioner does not choose. Use that path for tools and tests. Application code usually sets a key and lets the partitioner run.

If you add partitions to a topic, the map from key to partition can change. Old records stay in the old partitions. New records with the same key can go to a new partition. Per-key order can break across the change. Plan the partition count before you have a large history.

A bad key (a constant, or only a timestamp) creates a hot partition. Extra consumers in one group do not help a hot partition. Topic 4 covers that limit.

Do not use a random key when you need order. Do not use a user id as the only key when one user can emit many independent streams that you must process in parallel.

### Questions

#### Theoretical questions

1. What does the default partitioner do when the key is not null?
2. What does the default partitioner do when the key is null?
3. What happens to key-to-partition mapping when you add partitions?
4. What is a hot partition?
5. When do you set the partition number yourself?

#### Easy practical tasks

1. Write five sentences about keys and the partitioner. Use only facts from this section.
2. Make a table: Key type (stable entity, null, random, constant). Add partition behavior and order.
3. Produce five records with key `user-9` to a three-partition topic. Describe or consume. Write the partition.
4. Find the default partitioner name in your client documentation.

#### Medium practical tasks

1. Produce 30 records with null keys and 30 records with keys `a` and `b` only. Write how records spread across partitions.
2. Add a partition to a topic that already has keyed records. Produce the same keys again. Write whether the partition for a key changed.
3. Read official partitioner documentation. Write three facts in STE.

#### Advanced practical tasks

1. Write a custom partitioner only if your client supports it. Send even keys to partition 0 and odd keys to partition 1 on a two-partition topic. Prove it with consume output. Then delete the experiment topic.
2. Write a key standard: who chooses the key, what happens on a partition increase, and how you detect a hot key.

---

## Choosing a partition count

The partition count is the unit of parallel write and parallel read in one consumer group. One partition is one ordered log on one leader. That leader has a throughput limit: disk, network, and CPU.

Many partitions let many leaders work at the same time. Consumers in one group can read different partitions at the same time. You cannot have more active readers than partitions in that group. Topic 4 repeats that rule.

Choose a count that fits:

- expected produce rate and consume rate
- number of consumers that you plan in one group
- order needs (one partition for total order)
- broker count (spread leaders)

You can add partitions later. You cannot shrink the count. Extra partitions use memory and files on the broker. Extra partitions add work for the controller and for consumer rebalances. Do not create thousands of partitions without a measured reason.

A common beginner lab uses 3 partitions. A production count needs arithmetic: target rate, record size, and headroom. Topic 10 covers performance. Replication multiplies disk. Topic 6 covers retention and replicas.

On a one-broker KRaft lab, a high partition count still runs. It does not prove production capacity.

### Questions

#### Theoretical questions

1. Why does one partition have a throughput limit?
2. What limits the number of active consumers in one group?
3. Why is a very large partition count a problem?
4. Why must you choose the count with care if you cannot shrink it?
5. Why does a laptop cluster not prove a production partition count?

#### Easy practical tasks

1. Write five sentences that explain why partitions exist. Use only facts from this section.
2. Make a table: "Partition count" and "Max active consumers in one group". Add rows for 1, 3, and 6.
3. List three costs of each extra partition.
4. Create topic `par.lab` with three partitions. Describe it. Write why 3 is a reasonable lab default.

#### Medium practical tasks

1. Create topic `par.one` with one partition and topic `par.three` with three partitions. Produce 30 records to each. Compare how long the produce step takes. Write the two times and the method.
2. Find an official or Confluent note on partition count. Write three facts about broker load.
3. Plan a count for 20 MB/s and 1 KB records. Write the arithmetic and the assumptions.

#### Advanced practical tasks

1. Measure produce rate to 1, 3, and 12 partitions on your laptop cluster. Write a small table. State the limit that you hit (CPU, disk, or client).
2. Write a one-page note: when more partitions help, when they only add cost, and how you add partitions without a surprise key move.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe the path of one keyed record from producer to disk. Name partitioner, leader, offset, and replica.
2. How do key, partition count, and order interact when you add a partition?
3. Why are offset, timestamp, and produce-thread order three different ideas?
4. What stays on brokers versus what stays in the KRaft metadata log?
5. A teammate wants 1000 partitions "for speed" on a three-broker lab. What facts do you use in the reply?

#### Easy practical tasks

1. Create `parts.review` with three partitions. Produce two keyed records and two null-key records. Consume from the beginning. Write partition and offset for each record.
2. Write a one-page cheat sheet: topic, partition, offset, broker, cluster, record parts, replica, leader, order rule, default partitioner, count choice.
3. Draw one cluster: three brokers, one topic, three partitions, leaders spread, two replicas each if your RF allows it.
4. From a topic describe, highlight partition numbers, leaders, and offsets if shown.

#### Medium practical tasks

1. Write a small program that prints key, partition, offset, timestamp, and headers for every consumed record on a three-partition topic.
2. Increase partitions on a lab topic. Produce the same keys before and after. Save a table of key to partition.
3. Map each subsection to one official documentation heading.

#### Advanced practical tasks

1. Design keys and partition counts for three topics: payments (strict per-account order), metrics (null keys), and a compacted user profile. Write why each choice fits.
2. Write a production standard: naming, default RF, how you pick partition count, and a ban on ZooKeeper flags in create commands.
