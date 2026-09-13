# 14. Kafka Streams (Core)

## Description

Kafka Streams is a Java library that reads topics, processes records, and writes topics. You define a topology. The library runs the topology in your process. Instances that share an `application.id` share the work.

This topic covers topology, KStream, KTable, the Processor API versus the DSL, joins, exactly-once in Streams, interactive queries, and tests. Complete this topic after topics 5, 6, 7, 8, and 13. Run Streams against a KRaft Kafka cluster. Do not add ZooKeeper.

Use one term for each concept. A topology is the graph of sources, processors, and sinks. A KStream is a stream of events. A KTable is a changelog of the latest value per key. The DSL is the high-level API. The Processor API is the low-level API. Interactive queries read a state store from the running instances.

---

## Topology, KStream, KTable

A **topology** is a directed graph. Sources read Kafka topics. Processors transform records. Sinks write Kafka topics. Kafka Streams builds a task per thread from this graph. Each task owns a subset of partitions.

A **KStream** is a record stream. Each record is an event. Two records with the same key are two events. Example: order 1, then order 2 for the same user.

A **KTable** is a table view of a changelog. For each key, the table keeps the latest value. A null value is a tombstone and removes the key (topic 8). Example: user 9 now lives in city A, then city B. The table shows city B.

You can turn a KStream into a KTable (aggregate or latest-per-key). You can turn a KTable into a KStream (each change becomes an event).

Input topics for tables are often compacted. Input topics for streams often use delete retention. Match `cleanup.policy` to the meaning of the topic (topic 8).

`application.id` names the consumer group and prefixes internal topics (repartition, changelog). Do not reuse an `application.id` for a different topology without a reset plan.

```text
StreamsBuilder builder = new StreamsBuilder();
KStream<String, String> orders = builder.stream("orders.placed");
KTable<String, String> customers = builder.table("customers");
```

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
4. Find `StreamsBuilder` in the official docs. Write two methods that add a source.

#### Medium practical tasks

1. Write a short program that reads one topic as a KStream and writes a filtered topic. Run it on KRaft Kafka. Write the commands that create the topics.
2. Read the same data as a KTable in a second small program (or a second topology). Describe how a second record with the same key changes the table.
3. List internal topics after the job runs. Mark changelog and repartition names.

#### Advanced practical tasks

1. Build a topology that aggregates a KStream into a KTable (count per key) and writes the table to a sink topic. Show two input records and the sink values.
2. Write a one-page note: when you change a topology, how you choose a new `application.id` or a reset, and why you do not use ZooKeeper for that reset.

---

## Processor API vs DSL

The **DSL** (domain-specific language) is the high-level API: `stream`, `table`, `filter`, `map`, `groupByKey`, `aggregate`, `join`. Most jobs use the DSL. The DSL builds a topology for you.

The **Processor API** lets you implement a processor node. You receive records in `process`. You can schedule `punctuate` (time-based callbacks). You read and write state stores directly. You attach processors with `Topology.addSource`, `addProcessor`, and `addSink`.

Use the DSL when the operation exists. Use the Processor API when you need custom timing, custom store access, or a graph that the DSL does not express.

You can mix them. The DSL can add a custom processor. The Processor API can stay in a small part of a larger DSL topology.

The Processor API does not remove the need for correct keys, partitioning, and state restore. You still use changelog stores when the processor is stateful.

Do not write a Processor API job to avoid learning KStream and KTable. Learn the DSL first.

Both APIs run in the same Streams runtime. Both use KRaft Kafka as the log.

### Questions

#### Theoretical questions

1. What is the DSL in Kafka Streams?
2. What extra control does the Processor API give you?
3. What is `punctuate`?
4. When do you mix the two APIs?
5. Does the Processor API skip changelog restore for stateful processors?

#### Easy practical tasks

1. Make a two-column table: DSL operation, Processor API idea (process, punctuate, store).
2. Write four sentences about when you stay on the DSL.
3. Find an official Processor API example outline. Write the method names `init`, `process`, and `close` if they appear.
4. List five DSL methods from the official KStream Javadoc or docs page.

#### Medium practical tasks

1. Implement one filter in the DSL and describe the same filter as a Processor API `process` that forwards or skips. Do not need both in production.
2. Read official notes on punctuation (stream time versus wall clock). Write five STE sentences.
3. Draw a DSL topology and the same graph as Processor API nodes.

#### Advanced practical tasks

1. Write a Processor API node that counts records per key in a store and forwards the count every N seconds with punctuation. Test it with the test driver (later section) or a lab cluster.
2. Write a review rule: DSL by default, Processor API needs a design note, KRaft bootstrap only.

---

## Joins

A join in Streams combines records that share a key. Partitioning must use the same key and the same partition count on both sides (or Streams must repartition).

**Stream-stream join** joins two KStreams. The join is windowed. A record on one side finds records on the other side inside the join window. Without a window, the job cannot know how long to wait.

**Stream-table join** joins a KStream with a KTable. Each stream event looks up the current table value for that key. This is a fact enrichment. Example: order event plus current customer profile.

**Table-table join** joins two KTables. The result is a KTable. When either table changes, the join result for that key updates.

Foreign-key joins exist for tables when the join key is not the record key. Read the current Streams documentation for that case.

Null keys do not join. Tombstones on a table remove the join result for that key.

If keys are on the wrong partition, Streams writes a repartition topic. That topic is internal. It adds produce and consume cost.

Joins are stateful. They use stores and changelogs (topic 13).

### Questions

#### Theoretical questions

1. Why must join sides share a key and compatible partitioning?
2. Why is a stream-stream join windowed?
3. What does a stream-table join look up?
4. What does a table-table join produce?
5. What happens to a join when a table key receives a tombstone?

#### Easy practical tasks

1. Write one shop example for each join type (stream-stream, stream-table, table-table).
2. Make a table: Join type, window required, result type (stream or table).
3. Draw order stream + customer table → enriched order stream.
4. Write four sentences about repartition topics.

#### Medium practical tasks

1. Implement a stream-table join in the DSL on KRaft Kafka. Create both topics. Produce a customer, then an order. Consume the sink.
2. Produce the order before the customer. Write what the join emits (or does not emit). Then produce the customer and a second order.
3. Find official windowed join options. Write the window type you choose and why.

#### Advanced practical tasks

1. Implement a windowed stream-stream join (for example, click and view). Show a pair inside the window and a pair outside the window.
2. Write a one-page join standard: key design, partition count, compacted table topics, and how you detect a missing enrichment.

---

## Exactly-once in Streams

Exactly-once in Streams means the processing of an input record has one visible effect on the output topics and on the internal state that Streams manages. It uses the Kafka transactional protocol (topic 7).

You set `processing.guarantee` to `exactly_once_v2` (current name; older versions used `exactly_once`). Streams wraps consume, state update, and produce in a transaction. A restart does not apply the same input twice to the output.

`exactly_once_v2` needs brokers that support it. Current Apache Kafka on KRaft supports it. Read the version matrix in the official docs.

EOS does not cover side effects that you add in a processor (HTTP calls, extra database writes). Those side effects can run twice. Put external writes behind an idempotent consumer or an outbox (topic 19).

`at_least_once` is the other common guarantee. It is simpler and can write duplicates to the sink.

Idempotent producers are part of the EOS path. You still need to understand topic 7.

Isolation: transactional consumers of the sink must use `isolation.level=read_committed` if they must not see aborted writes (topic 7).

Do not enable EOS only on a sink producer while the Streams app uses at-least-once. Set the Streams guarantee.

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

1. Run a small Streams job with `at_least_once`. Crash it in the middle of a lab (kill the process). Write whether you can see a duplicate in the sink.
2. Run the same job with `exactly_once_v2` on KRaft. Repeat a crash test. Write what you observe. Do not claim a proof from one test.
3. Read official EOS requirements (broker version, settings). Write five facts.

#### Advanced practical tasks

1. Combine Streams EOS with a consumer of the sink that uses `read_committed`. Show that an aborted transaction is not visible.
2. Write a one-page design: Streams EOS for Kafka-to-Kafka, plus an idempotent write for one database side effect (topic 7 and topic 19).

---

## Interactive queries (high-level)

Interactive queries let you read a state store while the Streams instances run. The store is the same store that the topology updates. You do not need a separate database for a simple key lookup of that state.

Each instance holds only the partitions that it owns. A query for a key must go to the instance that owns that key. Kafka Streams exposes metadata: which instance hosts which store partition. You add RPC (HTTP or another protocol) between instances.

If you query the wrong instance, you get no key or you must forward the request. Build a thin query API in front of the instances.

Interactive queries are eventually consistent with the input log. A record that is not processed yet is not in the store. EOS and restore affect when a key appears.

Do not use interactive queries as a general-purpose database. There is no rich ad-hoc query language. Availability depends on the Streams instances. If all instances stop, queries stop.

This section is high-level. Read `ReadOnlyKeyValueStore` and `StreamsMetadata` in the official docs when you implement.

KRaft does not host the query API. Brokers host the changelog. Your process hosts the store and the query endpoint.

### Questions

#### Theoretical questions

1. What do interactive queries read?
2. Why can one instance not answer every key?
3. What metadata do you need to find the right instance?
4. Why are interactive queries not a full database?
5. Where does the query HTTP server run: broker or application?

#### Easy practical tasks

1. Write four sentences about interactive queries. Use only facts from this section.
2. Draw three Streams instances, one store, and a key that lives on instance 2. Draw a client request path.
3. Make a table: Query type, works on interactive queries (yes or no): get-by-key, full-text search, join across all keys.
4. Find `queryMetadataForKey` or the current metadata API name in the docs. Write it down.

#### Medium practical tasks

1. Read an official interactive-query example. Write the steps: expose store, find host, call RPC.
2. Explain in six sentences what happens to queries during a rebalance.
3. Compare interactive queries with a sink topic that another service loads into a database (topic 19 CQRS). Write four differences.

#### Advanced practical tasks

1. Add a get-by-key HTTP endpoint to a count topology. Query a key that exists and a key that does not. Write the responses.
2. Write an availability note: how many instances, how you route, and what the client sees when one instance stops.

---

## Testing topologies

You can test a topology without a Kafka cluster. `TopologyTestDriver` (Java test dependency) feeds input records and reads output records. You control time for windows.

Unit tests check map, filter, join, and aggregation logic. You pipe a list of records into the source topic name and assert the sink records.

Integration tests use a real cluster (EmbeddedKafka, Testcontainers, or a local KRaft broker). Those tests check Serdes, EOS, and restore. Do not skip integration tests if you use EOS or custom Serdes.

Test the topology class, not only the business function. Wrong Serdes or wrong join windows fail at runtime.

Do not use production `application.id` and production topics in tests. Use test names.

The test driver does not replace a load test (topic 16). It does not replace an operations test of rolling restart (topic 15).

Example test shape:

```text
TopologyTestDriver driver = new TopologyTestDriver(topology, props);
TestInputTopic<...> input = driver.createInputTopic(...);
TestOutputTopic<...> output = driver.createOutputTopic(...);
input.pipeInput(key, value);
assert output.readKeyValue() ...
```

### Questions

#### Theoretical questions

1. What does `TopologyTestDriver` replace in a unit test?
2. Why do you control time in window tests?
3. When do you still need a real KRaft cluster in tests?
4. Why must you test Serdes as part of the topology?
5. Why must tests avoid production `application.id`?

#### Easy practical tasks

1. Write five sentences about Streams testing. Use only facts from this section.
2. Make a table: Test type, uses cluster (yes or no), what it proves.
3. Find `TopologyTestDriver` in the official testing guide. Write two class names for input and output topics.
4. List three assertions you would write for a filter topology.

#### Medium practical tasks

1. Write a unit test for a filter topology with the test driver. One record passes. One record does not.
2. Write a unit test for a count aggregation. Pipe two keys. Assert two counts.
3. Add a window test that advances time. Assert that a late record is in or out of the window per your spec.

#### Advanced practical tasks

1. Write an integration test against Testcontainers or local KRaft: start job, produce, consume sink, then restart the job and produce again.
2. Write a test standard: unit tests for every topology, integration tests for EOS and Serdes, no ZooKeeper in test compose files.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe the path of one order record through a stream-table join into a sink under `exactly_once_v2`. Name topology parts and a store.
2. How do KStream, KTable, and interactive queries show three views of the same keys?
3. When do you choose the Processor API after you already know the DSL?
4. How do changelog topics, repartition topics, and application topics differ in a Streams app?
5. What from topics 6, 7, and 8 must you still set correctly for a Streams job to be safe?

#### Easy practical tasks

1. Write a one-page cheat sheet: topology, KStream, KTable, DSL, Processor API, three joins, EOS key, IQ, test driver.
2. Draw a full diagram: two source topics, join, store, changelog, sink, query arrow.
3. Bookmark the official Streams DSL, EOS, and testing pages.
4. Name every topic this sample job needs: `orders.placed`, `customers`, sink, plus two internal kinds.

#### Medium practical tasks

1. Build one job that uses a filter, a stream-table join, and a sink. Unit-test the topology. Run it on KRaft. Save configs (`application.id`, guarantee).
2. Map each subsection to one official URL.
3. Reset or replace `application.id` in a lab after a topology change. Write the topic cleanup steps you used.

#### Advanced practical tasks

1. Add interactive queries and EOS to the job. Crash one instance. Query a key and consume the sink. Write what stayed correct.
2. Write a production Streams standard: Java version, KRaft bootstrap, `processing.guarantee`, test layers, no ZooKeeper, naming of `application.id`.
