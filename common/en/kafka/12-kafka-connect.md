# 12. Kafka Connect

## Description

Kafka Connect is a framework that moves data between Kafka and other systems. A source connector writes records into Kafka. A sink connector reads records from Kafka and writes them to another system.

This topic covers source and sink connectors, tasks and converters, single message transforms (SMT), dead letter queues, and common connectors (JDBC, S3, Elasticsearch, Debezium). Complete this topic after topics 4, 5, and 11. Connect workers talk to a Kafka cluster. Use KRaft for that cluster. Do not add ZooKeeper for new work.

Use one term for each concept. A connector is a plugin that defines how data moves. A task is a unit of work that the worker starts for a connector. A converter turns bytes into Connect data and back. An SMT changes one record in the worker. A dead letter queue is a topic that stores records that failed.

---

## Source and sink connectors

A source connector reads an external system and produces records to Kafka topics. Examples: a database table, a file directory, a message API.

A sink connector consumes Kafka topics and writes to an external system. Examples: an object store, a search index, a warehouse.

Connect runs as one or more workers. Standalone mode uses one process and a local file for offsets. Distributed mode uses a group of workers. Distributed mode stores connector configuration, offsets, and status in Kafka topics. Those topics are internal Connect topics. Do not use them as application topics.

You submit a connector configuration (JSON through the REST API, or a file in standalone mode). The configuration names the connector class, the topics or tables, and the converters. The workers start the connector.

Connect is not a replacement for every custom producer or consumer. Use Connect when a maintained connector exists and you want start, stop, and scale in one place.

The Kafka cluster can be KRaft. Connect does not need ZooKeeper. Workers use `bootstrap.servers`. Admin tools for the cluster use `--bootstrap-server`.

Typical worker properties:

```text
bootstrap.servers=localhost:9092
group.id=connect-cluster
key.converter=...
value.converter=...
config.storage.topic=connect-configs
offset.storage.topic=connect-offsets
status.storage.topic=connect-status
```

Create those internal topics with the replication factor that your cluster can support. On a one-broker learning cluster, use replication factor 1.

### Questions

#### Theoretical questions

1. What does a source connector do?
2. What does a sink connector do?
3. What is the difference between standalone mode and distributed mode?
4. Where does distributed Connect store connector configuration and offsets?
5. Why does Connect not need ZooKeeper on a KRaft cluster?

#### Easy practical tasks

1. Write five sentences about source and sink connectors. Use only facts from this section.
2. Make a two-column table: "Direction" and "Example system". Add two source rows and two sink rows.
3. Find the official Connect documentation. Write the REST path that lists connectors (for your version).
4. List the three internal topic roles (config, offset, status) and one sentence for each.

#### Medium practical tasks

1. Start a KRaft Kafka cluster and a Connect worker. Call the REST endpoint that shows worker status. Save the command and the response.
2. Compare standalone and distributed in the official docs. Write five differences.
3. Draw boxes: external system, source connector, Kafka topic, sink connector, second external system.

#### Advanced practical tasks

1. Run a two-worker distributed Connect cluster against KRaft Kafka. Submit one file or datagen source. Confirm both workers are in the same `group.id`.
2. Write a one-page operations note: how you create Connect internal topics on KRaft, how you back them up, and why you do not treat them as application topics.

---

## Tasks and converters

A connector can start one or more tasks. A task is the process unit that copies a slice of the work. For a source, tasks can split tables or query ranges. For a sink, tasks can split topic partitions.

`tasks.max` in the connector configuration is an upper limit. The connector can start fewer tasks. One task cannot be larger than the work that the connector can split. Example: a sink cannot use more tasks than the number of assigned partitions in a useful way.

A converter serializes and deserializes record keys and values at the Connect boundary. Common converters:

- `StringConverter` for plain text
- `JsonConverter` for JSON
- Avro, Protobuf, or JSON Schema converters that talk to a Schema Registry (topic 11)

The converter on the source output must match what consumers expect. The converter on the sink input must match what is on the topic. A mismatch produces errors or garbage fields.

Connect has a key converter and a value converter. Header converters can exist. Set them in the worker defaults and override them per connector when needed.

Do not confuse a converter with an SMT. A converter changes the encoding. An SMT changes the record content or routing after the converter produces Connect data.

Internal Connect topics often use a JSON converter. Application topics can use a Registry converter. Read the connector documentation for the required pair.

### Questions

#### Theoretical questions

1. What is a Connect task?
2. What does `tasks.max` mean?
3. What does a converter do?
4. Why must the sink converter match the topic encoding?
5. How is a converter different from an SMT?

#### Easy practical tasks

1. Write four sentences about tasks. Use only facts from this section.
2. Make a table: Converter class, typical payload, registry required (yes or no).
3. Find `tasks.max` in a sample connector configuration in the official docs. Copy the key name only and write its meaning in your own words.
4. Draw a source path: external row → converter → Kafka bytes.

#### Medium practical tasks

1. Run one connector with `tasks.max=1` and then with a higher value. Write how many tasks actually start (REST or logs).
2. Produce JSON with a console producer. Point a sink at that topic with a mismatched converter. Record the error. Then match the converter.
3. Read converter documentation. Write when you set converters on the worker and when you set them on the connector.

#### Advanced practical tasks

1. Use an Avro converter with a Schema Registry and a KRaft Kafka cluster. Produce through a source or a producer. Consume with a compatible client. Write the subject name.
2. Write a standard: default converters for internal Connect topics versus application topics, and who may change `tasks.max`.

---

## SMT (single message transforms)

A single message transform (SMT) is a function that the worker applies to one record. The SMT runs in Connect, not in your application.

Typical SMT actions:

- rename or drop a field
- change a topic name (route)
- insert a static field
- mask a value
- flatten a nested structure

You chain SMTs in the connector configuration. The order in the list is the order of execution. A later SMT sees the output of the earlier SMT.

SMTs do not replace stream processing. An SMT must stay simple. It has no window and no join across records. For aggregations, use Kafka Streams, ksqlDB, or Flink (topic 13).

A bad SMT can drop required fields or break the converter contract. Test SMTs in a lab topic first.

Do not put secrets in SMT static fields. Do not use an SMT as the only access-control mechanism. Topic 17 covers ACLs and encryption.

Example configuration keys (names can differ by version):

```text
transforms=MaskCard
transforms.MaskCard.type=...
transforms.MaskCard.fields=card_number
```

Read the official SMT list for your Connect version. Connector-specific transforms can exist. Prefer maintained transforms.

### Questions

#### Theoretical questions

1. What is an SMT?
2. Where does an SMT run?
3. Why is an SMT not a stream processor?
4. Why does SMT order matter?
5. Why is an SMT a poor access-control tool?

#### Easy practical tasks

1. Write five sentences about SMTs. Use only facts from this section.
2. Make a table: SMT action, example use, risk if wrong.
3. Find three built-in SMT class names in the official documentation.
4. Write one SMT chain of two steps for an order record (example: drop a field, then rename a field). Do not implement it yet.

#### Medium practical tasks

1. Apply one official SMT (for example, insert a field or change a topic) on a lab connector. Show a record before and after.
2. Break a record on purpose with an SMT that drops a required field. Record the worker error.
3. Compare SMT documentation with a one-line Kafka Streams map (topic 13). Write three differences.

#### Advanced practical tasks

1. Build a two-SMT chain that masks one field and routes to a second topic. Write the configuration and two sample records.
2. Write a review checklist: when to use an SMT, when to use Streams, and how you test the chain on KRaft Kafka.

---

## Dead letter queues

A dead letter queue (DLQ) in Connect is a Kafka topic that receives records that the connector cannot process. Typical causes: a converter error, an SMT error, or a reject from the external system.

Without a DLQ, a poison record can block a sink task. The task retries and does not move forward. With a DLQ, the worker can write the failed record (and error headers) to the DLQ topic and continue.

Enable the DLQ only when the connector and your version support it. Configuration names include `errors.tolerance` and a DLQ topic name. Read the current Connect error-handling page for the exact keys.

The DLQ is a normal Kafka topic. Set retention. Set a replication factor that matches your durability needs on KRaft. Consume the DLQ with a separate process. Fix the record or the mapping. Do not ignore a growing DLQ.

A DLQ is not a retry delay queue with a schedule. Connect retries are a separate setting. Do not treat DLQ volume as success.

Source connectors can also fail. Some sources write errors to logs only. Read the connector documentation.

Topic 4 describes poison messages on the producer path. The Connect DLQ is the operational tool on the Connect path.

### Questions

#### Theoretical questions

1. What is a Connect DLQ?
2. What happens to a sink task when a poison record has no DLQ?
3. Is a DLQ the same as a scheduled retry queue?
4. What must you set on the DLQ topic besides the name?
5. Who reads the DLQ in a healthy operation?

#### Easy practical tasks

1. Write four sentences about the Connect DLQ. Use only facts from this section.
2. Make a table: Failure type, DLQ useful (yes or no), one reason.
3. Find `errors.tolerance` (or the current equivalent) in the official docs. Write the meaning in your own words.
4. Draw: sink task → error → DLQ topic → operator consumer.

#### Medium practical tasks

1. Create a DLQ topic on KRaft Kafka. Configure a sink that fails on bad JSON. Send one good record and one bad record. Confirm the bad record in the DLQ.
2. Read DLQ headers if the worker writes them. List the header names that you see.
3. Set retention on the DLQ. Describe the topic. Write the retention value.

#### Advanced practical tasks

1. Write a one-page DLQ runbook: alert when the DLQ rate rises, who consumes it, and how you replay a fixed record to the original topic.
2. Compare Connect DLQ with application poison-message handling from topic 4. Write five similarities and three differences.

---

## Common connectors: JDBC, S3, Elasticsearch, Debezium (CDC)

**JDBC source** reads tables or queries through JDBC and writes rows to Kafka. You choose incrementing mode, timestamp mode, or bulk mode. Incrementing mode needs a rising id. Timestamp mode needs a reliable timestamp column. Bulk mode copies the full result on each poll. Do not use bulk mode as the default for large tables.

**JDBC sink** writes Kafka records to tables. You must map fields to columns. Primary key configuration controls upsert versus insert.

**S3 sink** writes records to object storage (Amazon S3 or an S3-compatible API). The connector groups records into files. You set format (JSON, Avro, Parquet) and a flush size or interval. Object storage is not a message bus. Readers of S3 see files, not a live consumer group.

**Elasticsearch sink** writes records to an index. The document id often comes from the record key. Mapping and index settings belong to Elasticsearch. A converter or SMT must produce a document that the index accepts.

**Debezium** is a CDC (change data capture) project. A Debezium connector reads the database transaction log (WAL, binlog, or equivalent) and writes change events to Kafka. CDC is not the same as a JDBC poll. CDC can capture deletes and the order of commits. You must set database permissions and a snapshot policy. Debezium connectors often run on Kafka Connect.

Install connector plugins in the Connect plugin path. Pin plugin versions. Read the connector license. Some connectors are not part of Apache Kafka.

All of these connectors need a running KRaft Kafka cluster (or another Kafka-API cluster that you accept). They do not start ZooKeeper.

### Questions

#### Theoretical questions

1. What is the difference between JDBC incrementing mode and bulk mode?
2. What does an S3 sink write: records in a log, or files in a bucket?
3. What does CDC read that a JDBC poll does not read?
4. Why does a Debezium connector need database log access?
5. Are these connectors part of the Apache Kafka broker process?

#### Easy practical tasks

1. Make a table: Connector, source or sink, external system.
2. Write three sentences that contrast JDBC source and Debezium.
3. Find the plugin installation page for one connector. Write the directory or property that points to plugins.
4. Name one risk of JDBC bulk mode on a large table.

#### Medium practical tasks

1. Run one of: JDBC source (or a file source if you have no database), S3-compatible sink, or a Debezium tutorial against KRaft Kafka. Write the topic names that appear.
2. For JDBC, write a three-row comparison: incrementing, timestamp, bulk.
3. Read a Debezium event example in official docs. Label before-image, after-image, and operation type if they appear.

#### Advanced practical tasks

1. Design a path: Debezium or JDBC source → Kafka topic → S3 sink or Elasticsearch sink. Write converters, a possible SMT, and a DLQ topic. Do not invent product keys that you did not read.
2. Write an operations page: plugin version pin, Connect worker image, KRaft bootstrap, and how you restart one connector without a ZooKeeper step.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe the path of one database row from a source connector to a sink file or document. Name connector, task, converter, optional SMT, topic, and sink.
2. When do you write a custom consumer instead of a sink connector?
3. How do Connect internal topics, application topics, and a DLQ topic differ in purpose?
4. What stays the same when Connect moves from a local KRaft broker to MSK or another managed Kafka?
5. How does topic 11 (schema) change the converter choice for Connect?

#### Easy practical tasks

1. Write a one-page cheat sheet: source, sink, task, converter, SMT, DLQ, four common connectors.
2. Draw a distributed Connect worker group and a KRaft Kafka cluster. Label `bootstrap.servers` and the three internal topics.
3. Bookmark the official Connect page and one connector page that you will use.
4. List every configuration key named in this topic in one column and a six-word meaning in the second column.

#### Medium practical tasks

1. Start KRaft Kafka, Connect, and one source. Create a matching sink or a console consumer. Record every URL and topic name.
2. Fail a converter on purpose. Enable a DLQ (if the connector supports it). Export one DLQ record and the worker log line.
3. Map each subsection of this topic to one official URL.

#### Advanced practical tasks

1. Run Debezium or JDBC CDC-style capture into Kafka, then a sink to S3 or a file. Write a recovery test: stop the sink, write more source changes, start the sink, confirm catch-up.
2. Write a production Connect standard: KRaft only, plugin pins, converter policy, SMT review, DLQ runbook, and no ZooKeeper.
