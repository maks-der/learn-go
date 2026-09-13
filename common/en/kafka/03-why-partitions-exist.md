# 3. Why Partitions Exist

## Description

A topic is a set of partitions. Partitions exist so that Kafka can write and read in parallel. Partitions also set the unit of order. This topic explains parallel work, order rules, keyed records, null keys, and how you choose a partition count.

Complete this topic after topic 2. Complete this topic before you tune producers. Use one term for each concept. Order in Kafka is per partition. The partitioner is the producer logic that selects a partition. The default partitioner uses the key when the key is not null.

You can add partitions to a topic. You cannot shrink the partition count in Apache Kafka. Choose the count with care.

---

## Parallel writes and reads

One partition is one ordered log on one leader broker. The leader handles all writes for that partition. The leader handles the usual reads for that partition. One partition has a throughput limit: disk, network, and CPU on that leader.

Many partitions let many leaders work at the same time. Partition 0 can live on broker 1. Partition 1 can live on broker 2. Producers can write to both partitions at the same time. Consumers in one group can read both partitions at the same time.

Parallel write does not mean one producer process is always faster. A single producer still batches and sends to several partitions. The cluster can accept more total load when the load is spread.

Parallel read in one group is limited by the partition count. You cannot have more active readers than partitions in that group. Topic 6 covers that rule again.

Do not create thousands of partitions without a reason. Each partition uses memory and files on the broker. Each partition adds work for the controller and for consumer rebalances.

### Questions

#### Theoretical questions

1. Why does one partition have a throughput limit?
2. How do many partitions increase total write capacity?
3. How do many partitions increase total read capacity in one group?
4. What limits the number of active consumers in one group?
5. Why is a very large partition count a problem?

#### Easy practical tasks

1. Write five sentences that explain why partitions exist. Use only facts from this section.
2. Draw two brokers. Place two partitions. Draw two producers and two consumers in one group.
3. Make a table: "Partition count" and "Max active consumers in one group". Add rows for 1, 3, and 6.
4. List three costs of each extra partition.

#### Medium practical tasks

1. Create topic `par.one` with one partition and topic `par.three` with three partitions. Produce 30 records to each with a small program or a loop. Compare how long the produce step takes. Write the two times and the method.
2. Start three console consumers in one group on `par.three`. Produce records. Confirm that more than one consumer prints records.
3. Find an official or Confluent note on partition count. Write three facts about broker load.

#### Advanced practical tasks

1. Measure produce rate to 1, 3, and 12 partitions on your laptop cluster. Write a small table. State the limit that you hit (CPU, disk, or client).
2. Read about partition limits in operations docs. Write a one-page note: when more partitions help, and when they only add cost.

---

## Ordering is per partition, not per topic

Kafka guarantees order only inside one partition. If record A is appended before record B in the same partition, a consumer of that partition sees A before B.

Kafka does not guarantee order across partitions. Record A in partition 0 and record B in partition 1 can arrive at consumers in any interleaving.

If you need a total order for a stream, use one partition. That choice removes parallel consumers in one group. If you need order per entity, use a key so that all records for that entity go to the same partition.

Do not assume that "the topic is ordered" because you produced in a sequence from one thread. The default partitioner can send consecutive records with different keys to different partitions.

Time on the producer machine is not the order. The order is the append order on the leader of that partition.

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
3. Make a two-column table: "Need" and "Partition strategy". Add rows for total order, order per user, and no order need.
4. Explain in three sentences why timestamps are not the partition order.

#### Medium practical tasks

1. Produce ten numbered records with different keys to a three-partition topic. Consume all partitions. Write the numbers in the consume order. Mark where topic order breaks.
2. Produce ten numbered records with the same key. Consume. Write whether the numbers stay in order.
3. Find the official sentence about order. Rewrite it in STE in two sentences.

#### Advanced practical tasks

1. Write a one-page design: payment events must stay in order per account, but the cluster must use many partitions. Name the key and one failure if you use a random key.
2. Build a small consumer that merges two partitions by timestamp and show a case where that merge is not the produce causal order. Explain the result.

---

## Keyed records and the default partitioner

A partitioner selects the partition for a record. The producer runs the partitioner before it sends the record.

The default partitioner uses this idea: if the record has a key, hash the key, then map the hash to a partition. All records with the same key go to the same partition, as long as the partition count does not change.

The hash is not something you compute by hand in daily work. You choose a key that identifies the entity that must stay ordered.

If you add partitions later, the hash maps to a new partition number for some keys. Records for the same key can go to a new partition after the change. Old records stay on the old partition. Order per key can break across the change. Plan the partition count before you rely on key order for a long time.

You can write a custom partitioner in some clients. Do that only when the default key hash is not enough. Document the rule. Test it.

### Questions

#### Theoretical questions

1. What does a partitioner do?
2. How does the default partitioner use a key?
3. Why do two records with the same key go to the same partition?
4. What happens to key-to-partition mapping when you add partitions?
5. When do you write a custom partitioner?

#### Easy practical tasks

1. Produce five records with key `user-9` to a three-partition topic. Consume with partition print. Write the partition number. Confirm that all five share it.
2. Produce five records with keys `user-1` through `user-5`. Write the partition for each key.
3. Write five sentences about the default partitioner. Use only facts from this section.
4. Make a table: Key, Expected property (same partition or not). Add three key pairs.

#### Medium practical tasks

1. Write a small producer that sends 100 keys and prints a histogram: partition → count. Describe the spread.
2. Add partitions to the topic (for example from 3 to 6). Send the same keys again. Write which keys moved.
3. Read the client documentation for the default partitioner in your language. Write the class or function name and the key rule.

#### Advanced practical tasks

1. After you add partitions, consume the full topic. For one key that moved, write the old partition records and the new partition records. Explain the order risk.
2. Design a custom partitioner rule for "region + id" without writing production code. Write the algorithm in steps. List one test that must pass.

---

## Null key and sticky / round-robin partitioners

A record can have a null key. The default behavior for a null key is not "hash the key". Modern Java clients use a sticky partitioner for null keys. The producer sends a batch to one partition, then moves to another partition when the batch is full or the linger time ends. That design improves batching.

Older clients and some tools use round-robin for null keys. Round-robin sends consecutive records to consecutive partitions. Batches stay small.

Console producers often send lines with no key. Those records spread across partitions. Do not use the console with empty keys when you test key order.

If you need spread and you do not need key order, a null key is valid. If you need order per entity, set a key. Do not mix "order per entity" with null keys.

Some clients let you set the partition field on the record. That choice skips the partitioner. Use it only when you have a clear assignment rule.

### Questions

#### Theoretical questions

1. What is a null key?
2. What is a sticky partitioner?
3. What is a round-robin partitioner?
4. Why does a sticky partitioner help batching?
5. When is a null key the wrong choice?

#### Easy practical tasks

1. Produce six lines with the console producer and no key to a three-partition topic. Write the partition of each record.
2. Write four sentences that compare sticky and round-robin. Use only facts from this section.
3. Make a table: "Key type" and "Partitioner behavior". Add rows for non-null key and null key.
4. List two tests where you must set a key and two tests where a null key is enough.

#### Medium practical tasks

1. Write a producer that sends 50 records with null keys and prints the partition sequence. Say whether the sequence looks sticky or round-robin.
2. Compare the Java client docs and the client you use (Go, Python, or other). Write the null-key behavior for each.
3. Produce with an explicit partition number if the API allows it. Confirm that the key hash is not used.

#### Advanced practical tasks

1. Measure batch size or request count for 10 000 null-key records with linger set to 0 and linger set to 10 ms. Write the two results.
2. Write a one-page note: when to use null keys in a metrics stream versus a user-event stream.

---

## Choosing a partition count (hard to shrink)

Choose a partition count when you create the topic. Use these questions:

- How much write load will this topic take?
- How many consumers do you want in one group at peak?
- How many brokers do you have?
- Do you need a total order? If yes, use one partition.

A common learning default is 3 partitions on a laptop. A common production start is a small number that still lets you add consumers, then you add partitions if load grows.

You can increase the partition count. You cannot decrease it in Apache Kafka. Extra empty partitions still cost memory and files. Extra partitions change the key hash mapping.

Replication factor is a separate choice. Partition count is about parallelism and order. Replication factor is about copies. Do not raise partitions when you only need more copies. Add brokers and replicas instead.

Review the count with the team that consumes the topic. A connector or a stream job may have its own limits per partition.

### Questions

#### Theoretical questions

1. Which questions help you choose a partition count?
2. Why is a one-partition topic the right choice for total order?
3. Can you decrease the partition count in Apache Kafka?
4. Why is an increase of partitions a risk for key order?
5. How is partition count different from replication factor?

#### Easy practical tasks

1. Write a partition count and a reason for these topics: audit log (total order), user clicks (high load), feature flags (low load).
2. Create a topic with 3 partitions. Add partitions to 5. Describe the topic before and after.
3. Try to find a shrink command in `kafka-topics --help`. Write whether a decrease option exists.
4. Make a table: "Need" and "Change partitions or change replicas". Add three rows.

#### Medium practical tasks

1. Write a one-page decision sheet for a topic that will have 20 MB/s writes and 8 consumers in one group. Choose a count. Show the arithmetic.
2. After you add partitions, rerun a key histogram producer. Write how the distribution changed.
3. Read official docs on altering partitions. Write the command shape and the warning about keys.

#### Advanced practical tasks

1. Design a migration: a 3-partition keyed topic must become 12 partitions. List the order risk and one application-level fix (for example a new topic and a dual-write window).
2. On a three-broker cluster, choose partition counts that spread leaders. Create the topic. Fill a leader table. Adjust the count if all leaders sit on one broker and explain why.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. State the two jobs of partitions in one sentence each: parallelism and order.
2. Why does a sticky null-key strategy exist if keys already spread load?
3. How do you explain "the topic was out of order" to a teammate who produced numbered lines with different keys?
4. What cannot you undo after you create too many partitions?
5. How do partition count, consumer group size, and key design work together?

#### Easy practical tasks

1. Create `order.demo` with three partitions. Send five records with key `acct-1` and five records with no key. Write a table: record number, key, partition.
2. Write a one-page cheat sheet: order rule, default partitioner, null key, add vs shrink.
3. Draw the path of one keyed record from producer to partition leader.
4. List five topic names from topic 1 ideas. Write a partition count for each and one reason.

#### Medium practical tasks

1. Write a script that creates a topic, produces a fixed key set, prints a histogram, adds partitions, and prints a new histogram.
2. Run two consumers in one group on a one-partition topic. Write why one consumer is idle.
3. Find three official pages (order, partitioner, alter topic). Write one fact from each that this topic uses.

#### Advanced practical tasks

1. Build a small test that fails if two records with the same key land in different partitions (before any partition change). Then add partitions and show the test fail for new records.
2. Write a design review checklist (ten items) that a team must complete before they set `partitions=64` on a new topic.
