# 13. Stream Processing Survey

## Description

Stream processing reads records from Kafka, computes a result, and writes records or a queryable state. This topic is a survey. It names the common tools and the shared ideas. Topic 14 covers Kafka Streams in more detail.

This topic covers Kafka Streams, KSQL / ksqlDB, Apache Flink, stateless versus stateful processing, window types, changelog topics, and state stores. Complete this topic after topics 2, 5, 7, and 8. Stream jobs use a Kafka cluster. Use KRaft for new clusters. Stream processors do not need ZooKeeper.

Use one term for each concept. A stream job is a program that processes records continuously. Stateless processing uses only the current record. Stateful processing keeps a store. A window groups records by time. A changelog topic is a compacted Kafka topic that records store updates. A state store is the local (or queryable) state of the job.

---

## Kafka Streams

Kafka Streams is a Java library. You run it inside your application process. It is not a separate cluster. Each instance is a Kafka consumer and a Kafka producer (and more). Instances that share an `application.id` form a group. Kafka assigns partitions to those instances.

You build a topology: a graph of source topics, processors, and sink topics. Topic 14 explains KStream, KTable, the DSL, and the Processor API.

Kafka Streams stores local state on disk (often RocksDB) and backs that state with changelog topics on the Kafka cluster. After a restart or a rebalance, an instance can restore state from the changelog.

Exactly-once processing is available through the Kafka transactional protocol (topic 7 and topic 14). You set a processing guarantee in the Streams configuration.

Use Kafka Streams when the team writes Java (or a JVM language), the job fits in the application, and you want to operate the job as a normal service. Do not use Kafka Streams if you need a non-JVM first-class API as the only option.

Kafka Streams talks to brokers. Use KRaft brokers for new work. There is no ZooKeeper client in current Streams.

### Questions

#### Theoretical questions

1. Is Kafka Streams a cluster or a library?
2. What does `application.id` do?
3. Where does Kafka Streams keep local state, and how does it recover that state?
4. When is Kafka Streams a good fit?
5. Does Kafka Streams require ZooKeeper on a KRaft cluster?

#### Easy practical tasks

1. Write five sentences about Kafka Streams. Use only facts from this section.
2. Make a table: Role, Kafka Streams, a separate stream cluster (empty cells for the next sections).
3. Find the official Kafka Streams page. Write the sentence that says it is a library or a client.
4. Draw one application process, one `application.id` group of two instances, and two topic partitions.

#### Medium practical tasks

1. List the Streams configuration keys for bootstrap servers and `application.id` from the official docs.
2. Compare Kafka Streams with Connect (topic 12). Write four differences.
3. Read the official word-count example description (not a full copy). Write the input topic, the grouping, and the output topic in your own words.

#### Advanced practical tasks

1. After topic 14, run a small Streams job on KRaft Kafka. Write the changelog topic names that Kafka creates.
2. Write a one-page choice note: Streams as a library versus a cluster product. Include operations (deploy, scale, restore).

---

## KSQL / ksqlDB

ksqlDB (also called KSQL in older material) is a stream-processing system that uses SQL-like statements. You run a ksqlDB server. Clients send statements: create a stream, create a table, filter, join, aggregate.

A **stream** in ksqlDB is a sequence of events (like a KStream). A **table** is the latest value per key (like a KTable). Topic 14 defines those ideas in Kafka Streams. The names are similar on purpose.

ksqlDB uses Kafka as the log and uses Kafka Streams internally in many deployments. You still need a Kafka cluster. Use KRaft for that cluster. ksqlDB is not part of the Apache Kafka broker. It is a separate product (Confluent and related distributions).

Use ksqlDB when analysts or platform teams want SQL and a server, not a Java topology in every service. Do not expect every Streams feature or every custom processor in SQL.

Statement examples (shape only; follow current ksqlDB syntax):

```text
CREATE STREAM orders (... ) WITH (KAFKA_TOPIC='orders.placed', VALUE_FORMAT='JSON');
CREATE STREAM large_orders AS SELECT * FROM orders WHERE amount > 100;
```

Formats and Schema Registry follow topic 11. Persistent queries write results to Kafka topics.

### Questions

#### Theoretical questions

1. What interface does ksqlDB give the user?
2. What is a stream versus a table in ksqlDB at a high level?
3. Does ksqlDB replace the Kafka broker?
4. When do you choose ksqlDB instead of a Java Streams app?
5. Why can old documents say KSQL and new documents say ksqlDB?

#### Easy practical tasks

1. Write four sentences about ksqlDB. Use only facts from this section.
2. Make a table: Statement type, what it creates (stream, table, or topic).
3. Find the current ksqlDB documentation home. Write the product name that the page uses.
4. Translate one English rule ("orders with amount greater than 100") into a SELECT-style sentence. Do not run it yet.

#### Medium practical tasks

1. Read an official ksqlDB quick start. Write which extra processes it starts besides Kafka.
2. Compare one filter in ksqlDB SQL with the same filter described as a Streams map/filter (next sections). Write three similarities.
3. List two ksqlDB data formats that the docs name and how they relate to topic 11.

#### Advanced practical tasks

1. If you run a lab: start KRaft Kafka and ksqlDB. Create a stream from a topic. Run a persistent query. Name the result topic.
2. Write a platform standard: who may run persistent queries, how you name result topics, and that Kafka stays on KRaft.

---

## Apache Flink (often with Kafka)

Apache Flink is a stream-processing engine. You run a Flink cluster (JobManager and TaskManagers) or a managed Flink service. A Flink job is a DAG of operators. Flink can use Kafka as a source and as a sink.

Flink is not a Kafka library inside your small service by default. It is a compute cluster. It has its own checkpoint and savepoint model. Kafka offsets are part of that state when you use a Kafka source.

Use Flink when you need large-scale stateful jobs, complex event time, or a team that already operates Flink. Use Kafka Streams when the job belongs in a microservice and the language is JVM. The two can share the same KRaft Kafka cluster as the log.

Flink also has batch and stream unification. This handbook only needs the idea: Flink reads Kafka, computes, writes Kafka or another sink.

Connectors for Kafka in Flink must match the Kafka protocol version. They do not start ZooKeeper. Point them at `bootstrap.servers`.

Do not treat Flink as a drop-in replacement for Kafka Connect. Connect copies and transforms records. Flink computes windows, joins, and large state. Some teams use both.

### Questions

#### Theoretical questions

1. What kind of system is Apache Flink?
2. How does Flink use Kafka in a typical path?
3. Where do Kafka offsets live in a Flink Kafka source at a high level?
4. When do you choose Flink instead of Kafka Streams?
5. Does a Flink Kafka connector require ZooKeeper?

#### Easy practical tasks

1. Write five sentences that contrast Flink and Kafka Streams.
2. Make a table: Item, Kafka Streams, Flink. Add rows for process model, language, state restore idea.
3. Find the Flink Kafka connector documentation. Write the option name for bootstrap servers.
4. Draw: Kafka topic → Flink source → operators → Kafka sink.

#### Medium practical tasks

1. Read an official Flink checkpoint paragraph. Rewrite it in four STE sentences. Relate it to Kafka offsets in your own words.
2. List two Flink APIs or layers that the docs name (for example DataStream). Do not learn all of them in this topic.
3. Compare Flink with ksqlDB: who writes SQL versus who writes a cluster job. Write five sentences.

#### Advanced practical tasks

1. Run a Flink Kafka source/sink tutorial against KRaft Kafka if you have time. Write job start, topic names, and one checkpoint observation.
2. Write a one-page decision table for your team: Connect, Streams, ksqlDB, Flink. One row each. Include "needs a Kafka cluster (KRaft)".

---

## Stateless map/filter vs stateful aggregations

**Stateless** processing uses only the current record (and constants). `map` changes the record. `filter` keeps or drops the record. The job can lose local memory and still be correct if it reads the topic again from the last committed offset.

**Stateful** processing keeps a store keyed by a record key (and sometimes by a window). Examples: count, sum, join, unique user per day. The job must restore that store after a crash. Changelog topics and checkpoints exist for that reason.

A filter that drops small orders is stateless. A count of orders per customer is stateful. A join of orders and customers is stateful.

State grows with keys and with windows. Unbounded key sets fill disk. Plan compaction, retention, and tombstones on changelog topics (topic 8).

Exactly-once and at-least-once (topic 7) matter more when the job has state. A duplicate process of a stateless map can still write a duplicate output record. A duplicate process of a count can double the count unless the framework uses transactions or idempotent updates.

Do not hide a database write inside a "stateless" map without a plan. That write is external state.

### Questions

#### Theoretical questions

1. What makes a processor stateless?
2. What makes a processor stateful?
3. Why can a stateless job recover from an offset alone?
4. Why does a count per key need a store?
5. How can a duplicate processing of a count be wrong?

#### Easy practical tasks

1. Label six operations as stateless or stateful: map, filter, count, sum, drop field, join.
2. Write four sentences about state growth.
3. Give one shop example of a stateless rule and one of a stateful rule.
4. Draw a key `user-9` and a running count that updates on each order.

#### Medium practical tasks

1. Write a table: Operation, state key, what the store holds.
2. Explain in six sentences why changelog topics (next section) exist for stateful jobs and not for a pure filter.
3. Take the word-count example. Mark the stateful step. Mark the stateless steps.

#### Advanced practical tasks

1. Design a job: filter fraud flags (stateless), then count flags per account per hour (stateful). Name the store key and the output topic.
2. Write a one-page risk note: external database writes inside a map versus a Kafka sink with EOS (topics 7 and 14).

---

## Windowing: tumbling, hopping, session

A window groups records that share a key and that fall in a time range. Windowing is for stateful aggregation. Time can be event time (record timestamp) or processing time (clock on the job). Event time is the usual choice when records can arrive late.

**Tumbling** windows have a fixed size and do not overlap. Example: 5-minute windows at 10:00, 10:05, 10:10. Each record belongs to one tumbling window.

**Hopping** windows have a fixed size and a hop interval. The hop is smaller than the size when you want overlap. Example: size 5 minutes, hop 1 minute. One record can belong to more than one hopping window.

**Session** windows use a gap. If no record arrives for that key for the gap time, the session closes. The next record starts a new session. Session length is not fixed.

Late records need a grace period or allowed lateness in some engines. After the window closes, a late record can be ignored or can update a late result. Read the engine documentation.

Windows create more state than a simple count-by-key. Each window key is `(key, window-start)` or a session id. Retention of window state is a configuration.

Kafka Streams, ksqlDB, and Flink all implement these window types. The names are almost the same. The configuration keys differ.

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

## Changelog topics and state stores

A **state store** holds the current aggregation or join state. Kafka Streams often uses RocksDB on the instance disk. Flink uses its own state backends. ksqlDB uses Streams stores under the SQL layer.

A **changelog topic** is a Kafka topic that records each store update as a record (key = store key, value = new store value). The topic is usually compacted (topic 8). After a failure, the instance reads the changelog and rebuilds the store.

Tombstones in the changelog remove a key from the store (topic 8). Retention must keep the compacted history that restore needs. Do not set a short delete retention on a changelog unless you accept a wrong restore.

Changelog topics are internal. Kafka Streams names them from `application.id` and the store name. Do not write to them from a console producer.

Interactive queries (topic 14) read the state store. The changelog is the durable copy on the cluster.

KRaft stores cluster metadata. KRaft does not replace changelog topics. Partition data and changelogs still live in `log.dirs` on brokers.

If you delete an `application.id` or a Flink job without a plan, you can leave orphan internal topics. Use admin procedures.

### Questions

#### Theoretical questions

1. What is a state store?
2. What is a changelog topic?
3. Why is compaction common on changelog topics?
4. What does a tombstone mean in a changelog?
5. Does KRaft replace changelog topics?

#### Easy practical tasks

1. Write five sentences about changelog topics. Use only facts from this section.
2. Draw: record in → store update → changelog append → restore read.
3. Make a table: Store, changelog, application topic. Add one purpose row.
4. Find an official sentence about Streams changelog topics. Rewrite it in STE.

#### Medium practical tasks

1. After you run any Streams or ksqlDB lab, list topics. Mark internal changelog names.
2. Describe one changelog topic. Write cleanup policy and partition count if shown.
3. Explain in five sentences what happens if you delete a changelog topic while the job still runs.

#### Advanced practical tasks

1. Restore test: run a stateful job, stop it, delete only the local state directory (not Kafka), start again. Write how the job uses the changelog.
2. Write a one-page standard: naming of `application.id`, who may delete internal topics, compaction on changelogs, KRaft brokers, no ZooKeeper.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Place Connect, Kafka Streams, ksqlDB, and Flink on one line from "copy data" to "heavy compute". Give one reason for each place.
2. How do stateless filters, windows, and changelog topics work together in one job?
3. What is the same about all these tools from the Kafka cluster view?
4. Why does this handbook teach the survey before Kafka Streams core (topic 14)?
5. What must you still understand from topics 5–8 before you run a stateful job?

#### Easy practical tasks

1. Write a one-page cheat sheet: four tools, stateless versus stateful, three window types, store and changelog.
2. Draw one diagram that includes a KRaft cluster and two of the four tools that share topics.
3. Bookmark official pages for Streams, ksqlDB, and Flink Kafka connector.
4. Label five real shop rules as stateless, stateful, or windowed-stateful.

#### Medium practical tasks

1. Write a decision for three jobs: (a) mask a field into S3, (b) running balance per account, (c) 15-minute unique users. Choose Connect, Streams, ksqlDB, or Flink for each. Give one reason.
2. Map each subsection to one official URL.
3. List all internal topic types named in this topic and in topic 12. Write who creates them.

#### Advanced practical tasks

1. Design one pipeline that uses Connect (source), Kafka Streams or ksqlDB (aggregate), and a sink. Name topics, keys, window, and changelog role. Use KRaft only.
2. Write a two-page comparison for your team: when you add Flink, when you stay on Streams, and how both use the same KRaft Kafka.
