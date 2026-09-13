# 1. Getting Started

## Description

Apache Kafka is a distributed event log. Programs write events to topics. Other programs read those events. Kafka stores the events on disk for a configured time or size.

This topic shows what Kafka is, how Kafka is different from a message queue and from a database, and how you start a local cluster. Complete this topic before you study partitions, producers, and consumers.

Use one term for each concept. An event is a fact that occurred. A record is the unit that Kafka stores. A topic is a named log of records. Use KRaft for all new clusters. Do not start ZooKeeper for new work.

---

## What Kafka is (distributed event log)

Kafka is a distributed event log. A log is an append-only sequence of records. Producers append records. Consumers read records by position. Kafka does not delete a record when one consumer reads it.

Kafka is also a streaming platform. Many producers can write. Many consumer groups can read the same topic at the same time. Each group has its own position in the log. Kafka copies partitions across brokers so that the cluster can continue after one broker stops.

Kafka is not a request-reply service. Kafka is not a place to update one row by primary key as the default model. You append a new record. Readers replay the log or keep their own view.

Kafka runs as a cluster of brokers. New clusters use KRaft. KRaft stores cluster metadata in a Raft controller quorum. The controller quorum is part of Kafka. You do not add ZooKeeper to a new cluster.

Teams use Kafka to:

- move events between services
- collect logs and metrics
- feed stream jobs
- keep a durable history of changes

### Questions

#### Theoretical questions

1. What is an event log in Kafka?
2. Why does Kafka keep a record after one consumer reads it?
3. What does "distributed" mean for a Kafka cluster?
4. How is Kafka different from a request-reply API?
5. What role does KRaft have in a new cluster?

#### Easy practical tasks

1. Write five sentences that describe Kafka. Use only facts from this section.
2. Make a two-column table: "Kafka property" and "What it means". Add four rows.
3. List three system types that fit Kafka. List two system types that do not fit well. Give one reason for each choice.
4. Open [https://kafka.apache.org/intro](https://kafka.apache.org/intro). Write the definition of Kafka that the page uses in one sentence.

#### Medium practical tasks

1. Compare Kafka with one system that you know (a queue, a database, or a log file). Write six short sentences.
2. Draw a simple diagram: two producers, one topic, two consumer groups. Label the log.
3. Find the Apache Kafka release notes for the current major version. Write three changes that help a beginner.

#### Advanced practical tasks

1. Read the introduction on kafka.apache.org. Write a one-page timeline of Kafka with years and one fact per year. Use public sources.
2. Explain why an append-only log helps many independent readers. Give one numeric example (for example, two groups at different offsets).

---

## Kafka vs a message queue vs a database

A message queue such as RabbitMQ delivers a message to a consumer. After the consumer acknowledges the message, the broker can remove it. The queue is a buffer of work. Many queue systems do not keep a long history for new readers.

Kafka is a log. A consumer does not remove the record. Retention time, retention size, or compaction removes records. A new consumer group can read old records if the records are still on disk.

A database stores the current state. You update a row. You query by key. A database is the right tool when you need transactions on current records and rich queries. Kafka does not replace a database.

You can use both. A service writes a change to Kafka. Another service updates a database from those events. Kafka is the durable log of changes. The database is the queryable state.

Do not treat Kafka as a job queue without a design. Kafka has no built-in delay queue or priority queue like some brokers. Do not treat Kafka as the only source of truth for interactive queries.

Short contrast:

- Queue: consume and remove (typical default).
- Kafka: consume and keep (until retention).
- Database: store and update current state.

### Questions

#### Theoretical questions

1. What happens to a message in a typical queue after the consumer acknowledges it?
2. What happens to a record in Kafka after a consumer reads it?
3. When is a database a better tool than Kafka?
4. Can a new consumer group read old Kafka records? What limit applies?
5. Why is Kafka a poor default for a priority job queue?

#### Easy practical tasks

1. Make a three-column table: Queue, Kafka, Database. Add one row for "history", one row for "update in place", and one row for "many independent readers".
2. Write four sentences that say when you choose Kafka and when you choose a database.
3. Name two queue products and two database products. Write one sentence for each about a typical use.
4. Explain in three sentences why two consumer groups can share one Kafka topic.

#### Medium practical tasks

1. Read a short official tutorial for RabbitMQ or another queue. Write five differences from Kafka. Use facts, not opinions.
2. Draw two sequences for the same order event: one through a queue, one through Kafka. Mark when a late reader can still see the event.
3. Pick one need (search by field, replay last hour, exactly one worker). Choose Kafka, a queue, or a database. Write the reason.

#### Advanced practical tasks

1. Write a one-page design note for a shop: orders go to Kafka, stock lives in a database. List what each store owns. List one failure if you invert that choice.
2. Find the Kafka documentation on retention. Write the two main retention limits (time and size) in your own words. Do not copy a long passage.

---

## Events, producers, consumers, and topics

An event is a fact: "user 9 placed order 55". The producer is the program that writes the event to Kafka. The consumer is the program that reads the event. The topic is the named log that holds the records.

A record has a key, a value, a timestamp, and optional headers. Topic 2 describes the record in full. The producer serializes the key and the value to bytes. Kafka stores those bytes.

One topic has one or more partitions. A partition is an ordered sequence. The producer selects a partition. The consumer reads from one or more partitions. Topic 2 explains partitions.

Many producers can write to the same topic. Many consumers can read the same topic. Consumers that share work use a consumer group. Topic 4 explains groups.

Name topics by the event type or the stream, not by the one service that writes today. Example: `orders.placed` is clearer than `order-service-out`.

### Questions

#### Theoretical questions

1. What is an event in this handbook?
2. What does a producer do?
3. What does a consumer do?
4. What is a topic?
5. What parts does a record have at a high level?

#### Easy practical tasks

1. Write three example events from a shop. For each event, name a topic, a producer, and a consumer.
2. Draw boxes for producer, topic, and consumer. Draw arrows for write and read.
3. List five topic names for a blog system. Use a consistent style.
4. Explain in four sentences the difference between an event and a Kafka record.

#### Medium practical tasks

1. Take one user action (sign-up). Split it into two events that two teams can consume. Name the topics.
2. Make a table: column "Program", column "Producer, consumer, or both". Add four real program types (API, mailer, search indexer, audit).
3. Find the official glossary or introduction page. Write the official words for producer, consumer, and topic.

#### Advanced practical tasks

1. Design five topics for a payments system. For each topic, write the event meaning, one producer, and two consumers.
2. Argue in one page why a topic name must not include the consumer name. Give a counter-example that becomes painful.

---

## Installing Kafka in KRaft mode (or Docker / a managed service)

Use KRaft for every new install. KRaft is the Kafka metadata mode. The controller quorum stores topic and broker metadata. Do not install ZooKeeper for a new cluster.

**Apache Kafka on a machine.** Download the current Apache Kafka release from [https://kafka.apache.org/downloads](https://kafka.apache.org/downloads). Unpack the archive. Generate a cluster ID. Format the storage directories with that ID and the server property file. Start the Kafka process. The official documentation for your version shows the exact commands and the property file name.

Typical command names (the suffix is `.sh` on macOS and Linux, `.bat` on Windows):

```text
kafka-storage random-uuid
kafka-storage format -t <CLUSTER_ID> -c <server.properties>
kafka-server-start <server.properties>
```

After the broker starts, clients use a bootstrap address. The common local value is `localhost:9092`.

**Docker.** The official Apache Kafka container image starts a KRaft broker. Use Docker when you want a short local setup. Pin an image version. Do not use ZooKeeper compose files from old blogs.

**Managed service.** Confluent Cloud, Amazon MSK, and similar products speak the Kafka protocol. Cloud hides brokers. You still use topics, producers, and consumers. New MSK clusters can use KRaft. Follow the current vendor guide for the cluster mode.

**Redpanda.** Redpanda speaks the Kafka API. You can use many Kafka clients with Redpanda. The operations are not the same as Apache Kafka. This handbook teaches Apache Kafka. Use Redpanda only when you know that you want that product. Topic 12 compares products.

Check the install:

```text
kafka-topics --bootstrap-server localhost:9092 --list
```

An empty list or a list of internal topics means the client reached the cluster.

### Questions

#### Theoretical questions

1. Why must a new cluster use KRaft and not ZooKeeper?
2. What does `kafka-storage format` do?
3. What is a bootstrap address?
4. What is the same between Apache Kafka, MSK, and Redpanda from a client view?
5. Why must you pin a Docker image version?

#### Easy practical tasks

1. Start a single-node KRaft Kafka (local install or Docker). Run a topic list command. Save the full command and the output.
2. Create a topic `demo` with three partitions and replication factor 1. Describe the topic. Save the output.
3. Write a two-column table: "Install path" and "When you use it". Add rows for local Apache Kafka, Docker, and one managed service.
4. Find the `server.properties` (or container env) for your install. Write the values of the listener and the process roles if they appear.

#### Medium practical tasks

1. Start Kafka. Create `demo`. Use the console producer to send three lines. Use the console consumer with `--from-beginning`. Save the three lines that you see.
2. Stop the broker. Start it again. Consume `demo` from the beginning. Confirm that the three lines are still there.
3. Compare the official Apache Docker quick start with one managed-service quick start. Write five differences (commands, ports, extra services).

#### Advanced practical tasks

1. Start a three-node KRaft cluster on one machine (three config files, three ports) or with Compose. Create a topic with replication factor 3. Describe the topic. Record the leader and the replicas.
2. Take an old tutorial that starts ZooKeeper. Rewrite the steps for KRaft only. List every command that you remove and the KRaft command that replaces it.

---

## Official docs

The primary documentation is [https://kafka.apache.org/documentation/](https://kafka.apache.org/documentation/). It contains the design, the APIs, the configuration keys, and the operations notes.

The introduction is [https://kafka.apache.org/intro](https://kafka.apache.org/intro). Start there when the terms are new.

Use the configuration section when you need the exact meaning of a broker or client key. Use the API section when you write a client. Use the KRaft notes in the operations or configuration pages for metadata mode.

Confluent Developer ([https://developer.confluent.io/](https://developer.confluent.io/)) has tutorials. The tutorials often use Confluent Cloud or Schema Registry. The Kafka protocol is the same. The extra products are not required for this topic.

A book that many teams use is *Kafka: The Definitive Guide*. Use it as a second source. The official site is the source for the current version.

Use `kafka-topics --help` and the other script help texts in the `bin` folder. The help lists the flags for your installed version.

### Questions

#### Theoretical questions

1. What is the primary official documentation URL for Apache Kafka?
2. When do you open the configuration section instead of the introduction?
3. Why can a Confluent tutorial show extra services that Apache Kafka does not require?
4. Where do you find the meaning of `acks` or `linger.ms`?
5. Why must you check the docs for your Kafka version?

#### Easy practical tasks

1. Open the official documentation home. Write the titles of five top-level sections.
2. Open the introduction. Write one sentence about topics and one sentence about partitions.
3. Run `kafka-topics --help` (or the `.bat` form). Write three flags that you will use this week.
4. Bookmark the documentation home, the introduction, and the downloads page.

#### Medium practical tasks

1. Find the configuration entry for `log.retention.hours` in the official docs. Write the default and the meaning in your own words.
2. Find the producer API page or client section. List the four client roles the site names (producer, consumer, streams, connect) if they appear.
3. Compare one page on kafka.apache.org with the same topic on developer.confluent.io. Write three facts that match and one extra product that Confluent adds.

#### Advanced practical tasks

1. Make a one-page map of the official documentation for a beginner. Assign each of topics 1–12 in this path to one official page.
2. Read the KRaft section for your version. Write ten facts about controllers and metadata. Do not mention ZooKeeper except in one sentence that says you do not use it for new clusters.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe the path from an empty machine to a running one-node KRaft cluster and a listed topic. Name each kind of step (download or image, format, start, client).
2. Why does this handbook treat an event and a record as two terms?
3. A teammate wants to start ZooKeeper "because old blogs use it". Which facts do you use to refuse that plan?
4. How do the official docs and a console `--help` flag list work together in daily work?
5. What stays the same when you move from local Docker Kafka to a managed service from a producer view?

#### Easy practical tasks

1. Create topic `getting.started`. Produce two records with the console producer. Consume them from the beginning. Write the exact commands.
2. Write a one-page cheat sheet: bootstrap address, KRaft format idea, create topic, console produce, console consume, official doc URL.
3. Export or copy your broker property file (or Compose file). Highlight the listener, the KRaft role settings, and the log directory if present.
4. Draw one diagram that includes producer, topic, consumer, broker, and KRaft controller. Use one sentence under each box.

#### Medium practical tasks

1. Write a small script (PowerShell or bash) that waits for `localhost:9092` and then lists topics. Stop on failure.
2. Create two topics. Delete one. List topics. Record every command and the cluster state after each command.
3. Document your local Kafka setup in ten steps so that another beginner can copy it. Include KRaft. Do not include ZooKeeper.

#### Advanced practical tasks

1. Run two clients against the same local cluster: the console tools and one short program in a language that you know. Produce from one. Consume from the other.
2. Compare Apache Kafka in Docker with Redpanda in Docker (or a managed trial). Write a short report: start time, Kafka API compatibility, and one operation that is different.
