# 2. Core Abstractions

## Description

Kafka has a small set of core objects. You must know each name before you write producers and consumers. This topic defines topic, partition, offset, record, broker, cluster, replica, leader, and consumer group.

Complete this topic after topic 1. Complete this topic before topic 3. Topic 3 explains why partitions exist. Topic 6 explains consumer groups in more detail. Topic 9 explains replicas and durability.

Use one term for each concept. A topic is a named log. A partition is an ordered slice of that log. An offset is a position in one partition. A record is the stored unit. A broker is a Kafka server process. A cluster is a set of brokers. A replica is a copy of a partition. The leader is the replica that handles reads and writes. A consumer group is a set of consumers that share partitions of a topic.

---

## Topic

A topic is a named, append-only log of records. Producers write to a topic by name. Consumers subscribe to a topic by name. The name is a string. Use a clear, stable name. Example: `orders.placed`.

A topic has a partition count. A topic has a replication factor. You set those values when you create the topic, unless the broker defaults apply. You can add partitions later. You cannot shrink the partition count in Apache Kafka. Topic 3 covers that limit.

A topic has configuration. Examples: retention time, cleanup policy, compression type. You can set defaults on the broker. You can override values on one topic. Topic 8 and topic 10 cover storage and admin.

Internal topics exist for Kafka itself. Examples: `__consumer_offsets` and transaction topics. Do not write application events to those names. Topic 10 describes them.

A topic is not a folder of files that you edit. Kafka stores topic data as log segments on the brokers that hold the partitions.

### Questions

#### Theoretical questions

1. What is a topic?
2. How do producers and consumers refer to a topic?
3. What two numeric properties does a topic have at create time?
4. Can you reduce the number of partitions of a topic in Apache Kafka?
5. What is an internal topic at a high level?

#### Easy practical tasks

1. Write five topic names for a library system. Use one naming style.
2. Create a topic `core.topic`. Describe it. Write the partition count and the replication factor from the output.
3. Make a two-column table: "Topic property" and "What it means". Add four rows.
4. List three topic names that are poor (too vague or tied to one consumer). Rewrite each name.

#### Medium practical tasks

1. Create two topics with different retention settings if your install allows per-topic config. Describe both. Write the difference.
2. Find the official documentation page that defines a topic. Write the official sentence in your own words.
3. Try to create a topic with a name that contains a space or a forbidden character. Record the error. Create a valid name.

#### Advanced practical tasks

1. Design a topic list for an online course platform (enroll, watch, complete). For each topic, set a partition count and a retention reason.
2. Read about internal topics in the official docs. Write what `__consumer_offsets` stores in five sentences. Do not produce to that topic.

---

## Partition

A partition is an ordered sequence of records inside a topic. Kafka appends each new record to one partition. Records in one partition have increasing offsets.

A topic with one partition is a single ordered log. A topic with many partitions is a set of ordered logs. Kafka does not keep a single global order across all partitions. Topic 3 explains that rule.

Each partition has a number. The numbers start at 0. A topic with three partitions has partitions 0, 1, and 2.

Each partition has a leader replica. Producers write to the leader. Consumers read from the leader in the common case. Follower replicas copy the leader. Topic 9 covers the copy process.

Partitions let Kafka scale. Different partitions can live on different brokers. Different consumers in one group can read different partitions at the same time.

### Questions

#### Theoretical questions

1. What is a partition?
2. How does Kafka number partitions in a topic?
3. Does Kafka guarantee order across all partitions of a topic?
4. What is the leader of a partition?
5. How do partitions help a cluster scale?

#### Easy practical tasks

1. Create a topic with three partitions. Describe the topic. Write the partition numbers.
2. Draw one topic box that contains three partition boxes. Label offsets 0, 1, 2 in partition 0.
3. Write four sentences that describe a partition. Use only facts from this section.
4. Explain why a topic with one partition cannot use two active consumers in the same group for parallel read. Use one sentence from topic 6 if you already know it, or reason from "one ordered log".

#### Medium practical tasks

1. Produce six console records with no key to a three-partition topic. Consume from the beginning and print partition and offset if your consumer flag allows it. Write how the records spread.
2. Describe the topic and name the broker that leads each partition.
3. Compare a one-partition topic and a three-partition topic for the same event type. Write three sentences about order and three sentences about parallelism.

#### Advanced practical tasks

1. Create a topic with six partitions on a single broker. Then, if you have three brokers, create the same shape with replication factor 2. Draw where leaders live.
2. Read the official text on partitions. Write how a partition relates to a log directory on disk in five sentences.

---

## Offset

An offset is a number that identifies a record inside one partition. The first record in a partition has offset 0 in a new partition (the starting offset can differ after delete and some operations, but the offset still identifies the record).

Offsets increase as Kafka appends records. Offset 7 is after offset 6 in the same partition. Offsets in partition 0 have no meaning in partition 1. Always name both the partition and the offset when you point to a record.

A consumer stores a committed offset for each assigned partition. The committed offset is the next position to read, or the last processed position, depending on how the client commits. Topic 5 covers commit rules. Do not mix those two conventions in one team without a written rule.

The log end offset (LEO) is the next offset that the leader will assign. The high watermark is the offset that followers have confirmed. Consumers that read only committed data do not read past the high watermark. Topic 20 in the path covers that pair in more detail. For this topic, remember: the offset is the address of a record in a partition.

### Questions

#### Theoretical questions

1. What is an offset?
2. Why is an offset not enough without a partition number?
3. What is a committed offset for a consumer?
4. What is the log end offset at a high level?
5. Do two partitions share one offset sequence?

#### Easy practical tasks

1. Consume a small topic from the beginning with partition and offset print enabled. Write three lines as `partition:offset`.
2. Make a table with columns Partition, Offset, Value for six imaginary records in two partitions.
3. Write four sentences that explain why "offset 10" is incomplete.
4. Draw a partition as a row of boxes. Number the offsets. Mark one box as "next read".

#### Medium practical tasks

1. Produce three records. Consume them. Produce two more. Consume again without `--from-beginning` in a new group, then with `--from-beginning`. Write what you see.
2. Find `kafka-get-offsets` or the describe command that shows offsets. Write the earliest and latest offset for one partition.
3. Reset a consumer group to the beginning with the consumer-groups tool (or the equivalent). Consume again. Record the commands.

#### Advanced practical tasks

1. Explain the difference between "last processed offset" and "next offset to read" in a one-page note. Pick one convention for a team.
2. After retention deletes old records, describe a partition and read the earliest offset. Write why that number may not be 0.

---

## Record: key, value, timestamp, headers

A record is the unit that Kafka stores in a partition. A record has these parts:

- **Key.** Bytes, or null. The default partitioner uses the key to select a partition. Topic 3 covers that rule.
- **Value.** Bytes, or null. The value is the payload. A null value can be a tombstone on a compacted topic. Topic 8 covers tombstones.
- **Timestamp.** A time that the producer or the broker sets. Broker configuration chooses the timestamp type.
- **Headers.** Optional key/value pairs on the record. Use headers for small metadata (trace id, content type). Do not put the main payload in a header.

Kafka does not parse the key or the value. Kafka stores bytes. Your serializer writes the bytes. Your deserializer reads the bytes. Topic 11 covers schemas.

The record also has the partition and the offset after the broker appends it. The producer receives that metadata when the write succeeds (unless you use `acks=0`). Topic 4 covers acknowledgements.

Keep the key stable for the same entity. Example: order id as the key for all events about that order. Then those records go to the same partition and stay in order.

### Questions

#### Theoretical questions

1. What four parts of a record does this section name?
2. What type does Kafka store for the key and the value?
3. What is a header for?
4. Who can set the timestamp?
5. Why does a stable key matter for order?

#### Easy practical tasks

1. Produce two console records with keys (`key:value` mode if the tool supports it). Consume them. Write the keys that you see.
2. Make a table: Key, Value, Timestamp meaning, Header example. Add three rows for shop events.
3. Write five sentences that describe a record. Use only facts from this section.
4. List three data items that belong in the value and two items that can live in a header.

#### Medium practical tasks

1. Write a small producer in a language that you know. Set a key, a value, one header, and a timestamp if the API allows it. Read the record back.
2. Produce a record with a null key and a record with a key. Write which partition each record used (print metadata).
3. Find the official record or message documentation. Write the official field list in your own words.

#### Advanced practical tasks

1. Design a header convention for trace id and schema version. Write the header names and the value format. Produce and consume one record that uses them.
2. Compare record timestamp types (create time vs log append time) in the official docs. Write when each type is the better choice.

---

## Broker

A broker is one Kafka server process. The broker stores partition replicas on disk. The broker accepts produce and fetch requests from clients. The broker is a member of a cluster.

Each broker has a broker id. The id is unique in the cluster. Clients do not pick a broker for every request by hand. Clients use a bootstrap address. The cluster tells the client which broker leads each partition.

A broker can run controller work in KRaft mode. Some nodes are controllers, some nodes are brokers, or one process can combine both roles. Combined mode is common for a laptop cluster. Dedicated controller nodes are common for a large cluster.

A broker uses listeners. The listener is the host and port that clients use. Local work often uses `localhost:9092`. A wrong advertised address is a common cause of "I can list topics but I cannot produce" on Docker.

### Questions

#### Theoretical questions

1. What is a broker?
2. What is a broker id?
3. What is a bootstrap address for, if the cluster has many brokers?
4. What is combined KRaft mode?
5. Why must the advertised listener be reachable from the client?

#### Easy practical tasks

1. Find the broker id in your local cluster (describe cluster or broker logs). Write the id.
2. Write four sentences that describe a broker. Use only facts from this section.
3. Make a table: "Broker duty" and "Example". Add four rows.
4. Write the bootstrap address that your console tools use.

#### Medium practical tasks

1. Describe the cluster (or list brokers). Write each broker id and the listener if the tool shows it.
2. Change nothing in production. On a local Docker setup, read `advertised.listeners` (or the equivalent). Write how a client on the host reaches the broker.
3. Find the official broker configuration page. Write the meaning of `log.dirs` and `listeners` in your own words.

#### Advanced practical tasks

1. Run two brokers (or two containers) in one KRaft cluster. Create a topic with two partitions. Write which broker leads each partition.
2. Break the advertised listener on a local container on purpose (wrong hostname). Record the client error. Restore the value. Write the cause in five sentences.

---

## Cluster

A cluster is a set of brokers that work as one Kafka system. The cluster shares metadata: topics, partitions, replica assignments, and configuration. KRaft stores that metadata in the controller quorum.

Clients see one cluster through the bootstrap address. Any broker in the bootstrap list can send metadata. You list more than one bootstrap broker so that a start still works when one broker is down.

A cluster has a cluster id. You set the id when you format storage. All brokers in the cluster must use the same cluster id and the same metadata.

Do not mix brokers from two clusters in one client configuration. A second cluster is a second bootstrap list and a second set of topics. Topic 18 in the path covers copies between clusters.

A one-broker cluster is valid for learning. A production cluster uses more than one broker so that replicas can live on different machines.

### Questions

#### Theoretical questions

1. What is a cluster?
2. What metadata does a cluster share?
3. Why do you list more than one bootstrap broker?
4. What is a cluster id?
5. Why is a one-broker cluster a poor production default?

#### Easy practical tasks

1. Find the cluster id of your local install (logs or a cluster describe command). Write it.
2. Draw three brokers and one controller quorum box. Label the cluster.
3. Write five sentences that describe a cluster. Use only facts from this section.
4. Write your bootstrap list. If it has one address, write why a second address would help.

#### Medium practical tasks

1. Compare a one-broker describe-cluster output with the official picture of a multi-broker cluster. Write five differences you would expect.
2. Read the KRaft format documentation. Write why two brokers with different cluster ids cannot join.
3. List topics on your cluster. Mark which names are application topics and which names look internal.

#### Advanced practical tasks

1. Write a one-page plan for a three-broker learning cluster: ports, roles (combined or split), and a topic with replication factor 3.
2. Start that three-broker cluster, or document why your machine cannot, and complete the plan on paper with expected `describe` output.

---

## Replica and leader

A replica is a copy of a partition on a broker. The replication factor is the number of replicas for each partition. A replication factor of 3 means three copies.

One replica is the leader. The other replicas are followers. Producers write to the leader. Followers fetch from the leader and append to their local log. Consumers read from the leader in the usual setup.

If the leader broker stops, a follower in the in-sync set can become the new leader. Topic 9 defines the in-sync replica set (ISR) and the danger of unclean leader election.

Replicas are per partition, not per topic as a single blob. Partition 0 can have its leader on broker 1. Partition 1 can have its leader on broker 2.

For a one-broker learning cluster, the replication factor must be 1. A higher factor fails because there are not enough brokers.

### Questions

#### Theoretical questions

1. What is a replica?
2. What is the replication factor?
3. What does the leader do that a follower does not do in the usual produce path?
4. What happens when the leader broker stops, at a high level?
5. Why can a one-broker cluster not use replication factor 3?

#### Easy practical tasks

1. Describe a topic. Write the replica list and the leader for each partition.
2. Make a table: Partition, Leader, Replicas. Fill it from a real `describe` output.
3. Write four sentences that describe a replica and a leader.
4. Create a topic with replication factor 1. Try replication factor 3 on a one-broker cluster. Record the error.

#### Medium practical tasks

1. On a multi-broker cluster, create a topic with replication factor 2 or 3. Stop the leader broker if you can do that safely. Describe the topic again. Write the new leader.
2. Find the official replication documentation. Write the words leader and follower in your own sentences.
3. Draw partition 0 with three replicas. Mark the leader. Draw arrows for produce and for follower fetch.

#### Advanced practical tasks

1. Write a one-page note: replication factor 2 versus 3. Cover disk cost and the number of brokers that can stop.
2. Read ahead in topic 9 terms (ISR, `min.insync.replicas`) in the official docs. Write three sentences that connect those terms to the leader. Do not copy the later handbook text.

---

## Consumer group

A consumer group is a set of consumers that share one group id. Kafka assigns each partition of the subscribed topics to at most one active consumer in that group. Two consumers in the same group do not read the same partition at the same time.

Different groups are independent. Two groups can read the same topic at the same time. Each group has its own committed offsets.

The group coordinator is a broker that stores group membership and helps assignment. Topic 6 covers the coordinator, static membership, and rebalance types.

Use a new group id when you want a new independent reader. Reuse a group id when you want to share the work or continue from the stored offsets.

A consumer that uses manual assign can work without a group in some clients. The common application path is subscribe plus a group id. Topic 5 covers subscribe and assign.

### Questions

#### Theoretical questions

1. What is a consumer group?
2. How many active consumers in one group read the same partition at the same time?
3. How do two groups share one topic?
4. What is the group coordinator at a high level?
5. When do you choose a new group id?

#### Easy practical tasks

1. Start one console consumer with `--group g1` on a topic that has data. Write one record that you see.
2. Start a second console consumer with `--group g2` on the same topic from the beginning. Confirm that both groups can read the same records.
3. Write five sentences that describe a consumer group. Use only facts from this section.
4. Make a table: Group id, Topic, Meaning. Add three rows (billing, audit, search).

#### Medium practical tasks

1. Create a topic with two partitions. Start two console consumers with the same group id. Produce six records. Write which consumer prints which records (label the terminals).
2. Use `kafka-consumer-groups --describe --group g1`. Write the assigned partitions and the committed offsets.
3. Kill one consumer in a two-member group. Watch the remaining consumer. Write what happens to the partitions.

#### Advanced practical tasks

1. Draw the assignment for three partitions and two consumers in one group. Then draw three consumers. Mark the idle consumer if one exists.
2. Read the official consumer group introduction. Write how offsets are stored (topic name) and why a group id is part of that storage key.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Use one sentence for each term: topic, partition, offset, record, broker, cluster, replica, leader, consumer group.
2. Why must a learner learn partition and offset together?
3. How does a record key connect to a partition without a custom partitioner?
4. What objects exist only on brokers (not in the client process)?
5. A teammate says "the cluster offset of the message". Which terms do you correct, and why?

#### Easy practical tasks

1. Create `core.review` with three partitions. Produce one keyed record and one record without a key. Describe the topic. Consume from the beginning with offsets printed.
2. Write a one-page cheat sheet with a one-line definition for each core term in this topic.
3. From a `describe` output, highlight topic name, partitions, leader, replicas, and ISR if shown.
4. Draw one figure that contains all core terms. Use arrows for produce, replicate, and consume.

#### Medium practical tasks

1. Write a script that creates a topic, describes it, and lists consumer groups. Run it against your local cluster.
2. Start two groups on the same topic. Describe both groups. Write how the committed offsets can differ after you consume only in one group.
3. Map each core term to one official documentation heading. Write the URL and the heading text.

#### Advanced practical tasks

1. On a three-broker KRaft cluster, create a topic with six partitions and replication factor 3. Fill a table: partition, leader, replicas. Predict the table before you describe, then compare.
2. Write a short oral script (one minute) that teaches these terms to a person who knows databases only. Do not use ZooKeeper.
