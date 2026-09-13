# 9. Stream Processing

## Description

Stream processing reads records from Kafka, computes a result, and writes records or a queryable state. This topic names the common tools and the core Kafka Streams ideas.

This topic shows Kafka Streams versus ksqlDB versus Flink, KStream, KTable, topology, stateless versus stateful processing, windowing, and exactly-once in Streams. Complete this topic after topics 4, 5, 6, and 8. Stream jobs use a Kafka cluster. Use KRaft for new clusters. Stream processors do not need ZooKeeper.

Use one term for each concept. A stream job is a program that processes records continuously. A topology is the graph of sources, processors, and sinks. A KStream is a stream of events. A KTable is a changelog of the latest value per key. Stateless processing uses only the current record. Stateful processing keeps a store. A window groups records by time.

---

## Kafka Streams vs ksqlDB vs Flink

**Kafka Streams** is a Java library. You run it inside your application process. It is not a separate cluster. Each instance is a Kafka consumer and a Kafka producer (and more). Instances that share an `application.id` form a group. Kafka assigns partitions to those instances.

Kafka Streams stores local state on disk (often RocksDB) and backs that state with changelog topics on the Kafka cluster. After a restart or a rebalance, an instance can restore state from the changelog.

Use Kafka Streams when the team writes Java (or a JVM language), the job fits in the application, and you want to operate the job as a normal service. Do not use Kafka Streams if you need a non-JVM first-class API as the only option.

**ksqlDB** (also called KSQL in older material) is a stream-processing system that uses SQL-like statements. You run a ksqlDB server. Clients send statements: create a stream, create a table, filter, join, aggregate.

A stream in ksqlDB is a sequence of events (like a KStream). A table is the latest value per key (like a KTable). ksqlDB uses Kafka as the log and uses Kafka Streams internally in many deployments. You still need a Kafka cluster. ksqlDB is not part of the Apache Kafka broker. It is a separate product (Confluent and related distributions).

Use ksqlDB when analysts or platform teams want SQL and a server, not a Java topology in every service. Do not expect every Streams feature or every custom processor in SQL.

**Apache Flink** is a stream-processing engine. You run a Flink cluster (JobManager and TaskManagers) or a managed Flink service. A Flink job is a DAG of operators. Flink can use Kafka as a source and as a sink.

Flink is not a Kafka library inside your small service by default. It is a compute cluster. It has its own checkpoint and savepoint model. Kafka offsets are part of that state when you use a Kafka source.

Use Flink when you need large-scale stateful jobs, complex event time, or a team that already operates Flink. Use Kafka Streams when the job belongs in a microservice and the language is JVM. The two can share the same KRaft Kafka cluster as the log.

Connect copies and transforms records (topic 8). Streams, ksqlDB, and Flink compute. Some teams use Connect and a stream job together.

All of these tools talk to brokers with `bootstrap.servers`. They do not start ZooKeeper. Use KRaft for new Kafka clusters.

### Questions

#### Theoretical questions

1. Is Kafka Streams a cluster or a library?
2. What interface does ksqlDB give the user?
3. What kind of system is Apache Flink?
4. When do you choose Flink instead of Kafka Streams?
5. Does a Flink Kafka connector require ZooKeeper?

#### Easy practical tasks

1. Write five sentences that contrast the three tools. Use only facts from this section.
2. Make a table: Item, Kafka Streams, ksqlDB, Flink. Add rows for process model and language.
3. Find the official Kafka Streams page. Write the sentence that says it is a library or a client.
4. Draw: Kafka topic → one of the three tools → Kafka sink.

#### Medium practical tasks

1. Compare Kafka Streams with Connect (topic 8). Write four differences.
2. Read an official ksqlDB quick start. Write which extra processes it starts besides Kafka.
3. Find the Flink Kafka connector documentation. Write the option name for bootstrap servers.

#### Advanced practical tasks

1. Write a one-page decision table for your team: Connect, Streams, ksqlDB, Flink. One row each. Include "needs a Kafka cluster (KRaft)".
2. If you run a lab: start KRaft Kafka and one of Streams, ksqlDB, or Flink. Name the input topic, the output topic, and one extra process.

---

## KStream, KTable, and topology

A **topology** is a directed graph. Sources read Kafka topics. Processors transform records. Sinks write Kafka topics. Kafka Streams builds a task per thread from this graph. Each task owns a subset of partitions.

A **KStream** is a record stream. Each record is an event. Two records with the same key are two events. Example: order 1, then order 2 for the same user.

A **KTable** is a table view of a changelog. For each key, the table keeps the latest value. A null value is a tombstone and removes the key (topic 6). Example: user 9 now lives in city A, then city B. The table shows city B.

You can turn a KStream into a KTable (aggregate or latest-per-key). You can turn a KTable into a KStream (each change becomes an event).

Input topics for tables are often compacted. Input topics for streams often use delete retention. Match `cleanup.policy` to the meaning of the topic (topic 6).

`application.id` names the consumer group and prefixes internal topics (repartition, changelog). Do not reuse an `application.id` for a different topology without a reset plan.

```text
StreamsBuilder builder = new StreamsBuilder();
KStream<String, String> orders = builder.stream("orders.placed");
KTable<String, String> customers = builder.table("customers");
```

The **DSL** is the high-level API: `stream`, `table`, `filter`, `map`, `groupByKey`, `aggregate`, `join`. Most jobs use the DSL. The **Processor API** lets you implement a processor node. Learn the DSL first.

Joins combine records that share a key. Stream-stream joins are windowed. Stream-table joins look up the current table value. Table-table joins produce a table. Null keys do not join. If keys are on the wrong partition, Streams writes a repartition topic.

The cluster is KRaft. Streams uses `bootstrap.servers`. There is no ZooKeeper setting in current Streams.

### Questions

#### Theoretical questions

1. What is a topology?
2. How is a KStream different from a KTable?
3. What does a tombstone do in a KTable?
4. What does `application.id` name?
5. Why must a table input topic often use compaction?

#### Easy practical tasks

1. Write five sentences about KStream and KTable. Use only facts from this section.
2. Make a table: Topic meaning, KStream or KTable, typical cleanup policy.
3. Draw a topology: source `orders.placed` → filter → sink `orders.large`.
4. Write one shop example for stream-table join.

#### Medium practical tasks

1. Write a short program that reads one topic as a KStream and writes a filtered topic. Run it on KRaft Kafka. Write the commands that create the topics.
2. List internal topics after the job runs. Mark changelog and repartition names.
3. Implement a stream-table join in the DSL if you can. Produce a customer, then an order. Consume the sink.

#### Advanced practical tasks

1. Build a topology that aggregates a KStream into a KTable (count per key) and writes the table to a sink topic. Show two input records and the sink values.
2. Write a one-page note: when you change a topology, how you choose a new `application.id` or a reset, and why you do not use ZooKeeper for that reset.

---

## Stateless map/filter vs stateful aggregations

**Stateless** processing uses only the current record (and constants). `map` changes the record. `filter` keeps or drops the record. The job can lose local memory and still be correct if it reads the topic again from the last committed offset.

**Stateful** processing keeps a store keyed by a record key (and sometimes by a window). Examples: count, sum, join, unique user per day. The job must restore that store after a crash. Changelog topics exist for that reason.

A filter that drops small orders is stateless. A count of orders per customer is stateful. A join of orders and customers is stateful.

A **state store** holds the current aggregation or join state. Kafka Streams often uses RocksDB on the instance disk. A **changelog topic** is a Kafka topic that records each store update. The topic is usually compacted (topic 6). After a failure, the instance reads the changelog and rebuilds the store.

Tombstones in the changelog remove a key from the store. Do not write to changelog topics from a console producer. Do not set a short delete retention on a changelog unless you accept a wrong restore.

State grows with keys and with windows. Unbounded key sets fill disk.

Exactly-once and at-least-once (topic 5) matter more when the job has state. A duplicate process of a count can double the count unless the framework uses transactions or idempotent updates.

Do not hide a database write inside a "stateless" map without a plan. That write is external state.

KRaft stores cluster metadata. KRaft does not replace changelog topics. Partition data and changelogs still live in `log.dirs` on brokers.

### Questions

#### Theoretical questions

1. What makes a processor stateless?
2. What makes a processor stateful?
3. Why can a stateless job recover from an offset alone?
4. What is a changelog topic?
5. Does KRaft replace changelog topics?

#### Easy practical tasks

1. Label six operations as stateless or stateful: map, filter, count, sum, drop field, join.
2. Write four sentences about state growth.
3. Draw: record in → store update → changelog append → restore read.
4. Find an official sentence about Streams changelog topics. Rewrite it in STE.

#### Medium practical tasks

1. Take the official word-count example description. Mark the stateful step. Mark the stateless steps.
2. After you run any Streams or ksqlDB lab, list topics. Mark internal changelog names.
3. Explain in six sentences what happens if you delete a changelog topic while the job still runs.

#### Advanced practical tasks

1. Design a job: filter fraud flags (stateless), then count flags per account per hour (stateful). Name the store key and the output topic.
2. Restore test: run a stateful job, stop it, delete only the local state directory (not Kafka), start again. Write how the job uses the changelog.

---

## Windowing

A **window** groups records that share a key and that fall in a time range. Windowing is for stateful aggregation. Time can be event time (record timestamp) or processing time (clock on the job). Event time is the usual choice when records can arrive late.

**Tumbling** windows have a fixed size and do not overlap. Example: 5-minute windows at 10:00, 10:05, 10:10. Each record belongs to one tumbling window.

**Hopping** windows have a fixed size and a hop interval. The hop is smaller than the size when you want overlap. Example: size 5 minutes, hop 1 minute. One record can belong to more than one hopping window.

**Session** windows use a gap. If no record arrives for that key for the gap time, the session closes. The next record starts a new session. Session length is not fixed.

Late records need a grace period or allowed lateness in some engines. After the window closes, a late record can be ignored or can update a late result. Read the engine documentation.

Windows create more state than a simple count-by-key. Each window key is `(key, window-start)` or a session id. Retention of window state is a configuration.

Kafka Streams, ksqlDB, and Flink all implement these window types. The names are almost the same. The configuration keys differ.

Stream-stream joins are also windowed. A record on one side finds records on the other side inside the join window.

### Questions

#### Theoretical questions

1. What is a tumbling window?
2. What is a hopping window?
3. What is a session window?
4. What is the difference between event time and processing time?
5. Why do windows increase state size?

#### Easy practical tasks

1. Draw three tumbling windows of 5 minutes on a time line. Place four records.
2. Draw hopping windows of size 5 minutes and hop 1 minute for the same records. Mark overlap.
3. Write a session example: two clicks 2 minutes apart, gap 5 minutes, then a click 20 minutes later.
4. Make a table: Window type, overlap, length fixed (yes or no).

#### Medium practical tasks

1. Find tumbling, hopping, and session in Kafka Streams documentation. Write the method or operator names.
2. Find the same three ideas in Flink or ksqlDB docs. Write the names they use.
3. Write six sentences about late records and grace. Use only facts that you read in one official page.

#### Advanced practical tasks

1. Specify one metric: orders per store per 15-minute tumbling window, event time, 2-minute grace. Write the store key and when a result is final.
2. Compare hopping and tumbling for a 10-minute moving view. Write when hopping is required and the state cost.

---

## Exactly-once in Streams

Exactly-once in Streams means the processing of an input record has one visible effect on the output topics and on the internal state that Streams manages. It uses the Kafka transactional protocol (topic 5).

You set `processing.guarantee` to `exactly_once_v2` (current name; older versions used `exactly_once`). Streams wraps consume, state update, and produce in a transaction. A restart does not apply the same input twice to the output.

`exactly_once_v2` needs brokers that support it. Current Apache Kafka on KRaft supports it. Read the version matrix in the official docs.

EOS does not cover side effects that you add in a processor (HTTP calls, extra database writes). Those side effects can run twice. Put external writes behind an idempotent consumer or an outbox (topic 11).

`at_least_once` is the other common guarantee. It is simpler and can write duplicates to the sink.

Isolation: transactional consumers of the sink must use `isolation.level=read_committed` if they must not see aborted writes (topic 5).

Do not enable EOS only on a sink producer while the Streams app uses at-least-once. Set the Streams guarantee.

You can test a topology without a cluster with `TopologyTestDriver` (Java). Integration tests still need a real KRaft cluster for Serdes, EOS, and restore.

### Questions

#### Theoretical questions

1. What does exactly-once mean in Kafka Streams?
2. Which configuration key sets the guarantee?
3. What work does a Streams transaction include?
4. What side effects are outside EOS?
5. Why must a sink consumer use `read_committed` to avoid aborted records?

#### Easy practical tasks

1. Write five sentences about Streams EOS. Use only facts from this section.
2. Make a table: Guarantee, possible duplicate sink records, typical use.
3. Find `processing.guarantee` in the official configuration page. Write the values that the page lists.
4. List three operations that EOS covers and two that it does not cover.

#### Medium practical tasks

1. Run a small Streams job with `at_least_once`. Stop it in the middle of a lab. Write whether you can see a duplicate in the sink.
2. Run the same job with `exactly_once_v2` on KRaft. Repeat a crash test. Write what you observe. Do not claim a proof from one test.
3. Read official EOS requirements (broker version, settings). Write five facts.

#### Advanced practical tasks

1. Combine Streams EOS with a consumer of the sink that uses `read_committed`. Show that an aborted transaction is not visible.
2. Write a one-page design: Streams EOS for Kafka-to-Kafka, plus an idempotent write for one database side effect (topic 5 and topic 11).

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Place Connect, Kafka Streams, ksqlDB, and Flink on one line from "copy data" to "heavy compute". Give one reason for each place.
2. How do stateless filters, windows, and changelog topics work together in one job?
3. What is the same about all these tools from the Kafka cluster view?
4. Why must you still understand consumer groups and compaction before you run a stateful job?
5. How do `application.id`, changelog topics, and EOS share one failure story?

#### Easy practical tasks

1. Write a one-page cheat sheet: three tools, KStream versus KTable, stateless versus stateful, three window types, Streams EOS.
2. Draw one diagram that includes a KRaft cluster and two of the three tools that share topics.
3. Bookmark official pages for Streams, ksqlDB, and Flink Kafka connector.
4. Label five real shop rules as stateless, stateful, or windowed-stateful.

#### Medium practical tasks

1. Write a decision for three jobs: (a) mask a field into S3, (b) running balance per account, (c) 15-minute unique users. Choose Connect, Streams, ksqlDB, or Flink for each. Give one reason.
2. Map each subsection to one official URL.
3. List internal topic types named in this topic and in topic 8. Write who creates them.

#### Advanced practical tasks

1. Design one pipeline that uses Connect (source), Kafka Streams or ksqlDB (aggregate), and a sink. Name topics, keys, window, and changelog role. Use KRaft only.
2. Write a two-page comparison for your team: when you add Flink, when you stay on Streams, and how both use the same KRaft Kafka.
