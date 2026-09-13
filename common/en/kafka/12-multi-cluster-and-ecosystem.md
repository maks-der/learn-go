# 12. Multi-Cluster and Ecosystem

## Description

One Kafka cluster lives in one failure domain even when you use racks (topic 10). A second cluster in another region or another account is a different system. The Kafka ecosystem is the set of products that speak the Kafka protocol or that sit next to a Kafka cluster.

This topic shows MirrorMaker 2, disaster recovery, Confluent, Redpanda, managed Kafka, Schema Registry, REST Proxy, UIs, and a learning order: produce and consume in your language, then add a Streams job. Complete this topic after topics 1, 4, 7, 8, and 9. Each new cluster uses KRaft. Do not add ZooKeeper. Replication tools talk to each cluster with `bootstrap.servers`.

Use one term for each concept. A source cluster is the cluster that already has the records. A target cluster is the cluster that receives the copy. MirrorMaker 2 (MM2) is the Apache Connect-based replicator. Disaster recovery (DR) is the plan to continue after a cluster or region fails. Apache Kafka is the ASF project. A Kafka-API compatible system accepts Kafka clients. A managed service runs brokers for you.

---

## MirrorMaker 2 and disaster recovery

**MirrorMaker 2** is the Apache Kafka replicator. It runs as Kafka Connect connectors. It copies topics from a source cluster to a target cluster. It can copy consumer group offsets (with translation). It can create topics on the target with a name pattern (often a prefix or a suffix).

MM2 uses source and checkpoint connectors. Workers need network access to both clusters. Configure two `bootstrap.servers` values. Both clusters can be KRaft. MM2 does not need ZooKeeper.

The copy is asynchronous. The target lags the source. That lag is your replication delay. Measure it.

MM2 is not a transaction that spans two clusters. Exactly-once inside one cluster (topic 5) does not make two clusters one atomic log.

Typical name pattern: source topic `orders.placed` becomes `sourceCluster.orders.placed` on the target. Consumers must use the target name.

**Offset translation** maps a group’s committed offset on the source to an offset on the target. Offset 1000 on cluster A is not the same record as offset 1000 on cluster B. Without translation, a common mistake is to start the group on the target at latest (lose data) or at earliest (reprocess everything). Prefer the same partition count on both clusters.

**Active-passive** means one cluster receives all produces for a topic (the active). The other cluster is a copy (the passive). On disaster, you switch clients to the passive cluster after a promote step.

**Active-active** means two clusters accept produces on the "same" business topic at the same time. You get conflicts. Most DR designs start as active-passive. Do not run two MM2 loops on the same topic in both directions without a filter, or you create a loop.

**Cluster Linking** is a Confluent product that links topics. It is not part of a plain Apache tarball. If you have only Apache Kafka, use MM2.

**Disaster recovery** is the documented path when a cluster or a region fails. You define **RPO** (how much recent data you can lose, related to lag) and **RTO** (how long until clients write and read again).

Failover steps (active-passive, shape only):

1. Stop producers to the failed source (or confirm they fail).
2. Confirm target lag and what you accept as RPO.
3. Promote the target (MM2 or Cluster Linking promote, or change it to the writer).
4. Apply offset translation for groups that must continue.
5. Point clients to the target bootstrap. Use TLS and SASL for that cluster (topic 10).
6. When the old region returns, do not start two actives without a failback plan.

Single-cluster rack awareness (topic 10) is not DR for a region loss. Test DR. An untested runbook fails. Use KRaft on both sides so that the runbook has no ZooKeeper restore step.

```text
# Conceptual MM2 keys (names follow current MM2 docs)
source.cluster.bootstrap.servers=...
target.cluster.bootstrap.servers=...
topics=orders.placed
```

### Questions

#### Theoretical questions

1. What does MirrorMaker 2 copy?
2. Why is offset 1000 not portable between clusters?
3. What is active-passive versus active-active?
4. What do RPO and RTO mean?
5. Why is rack awareness not region disaster recovery?

#### Easy practical tasks

1. Write five sentences about MM2 and DR. Use only facts from this section.
2. Make a table: Cluster, bootstrap, MM2 needs (read or write).
3. Find official MM2 documentation. Write the connector class names that the page lists.
4. Draw: source KRaft cluster → MM2 workers → target KRaft cluster.

#### Medium practical tasks

1. Read an official MM2 quick start. Write every process it starts (Connect, two Kafka clusters).
2. Make a table: Strategy on target (earliest, latest, translated), skip risk, duplicate risk.
3. Pick a shop: orders must be unique. Choose active-passive or active-active. Write the reason in six sentences.

#### Advanced practical tasks

1. Run MM2 between two local KRaft clusters (two Compose stacks or two ports). Copy one topic. Confirm records on the target. Write the target topic name.
2. Write a DR runbook: RPO, RTO, promote steps, offset translation, failback, ACLs, KRaft only, no ZooKeeper.

---

## Confluent, Redpanda, and managed Kafka

**Apache Kafka** is the project at [https://kafka.apache.org/](https://kafka.apache.org/). You download brokers, client libraries, Kafka Streams, and Kafka Connect. You operate them. New clusters use KRaft.

**Confluent Platform** is a distribution from Confluent. It includes Apache Kafka (or a compatible broker) plus extra components: Schema Registry, ksqlDB, Control Center, connectors, and commercial features. **Confluent Cloud** is the hosted form.

The Kafka protocol is the same idea. A producer that uses the Apache client can talk to Confluent Cloud if you set bootstrap, security, and sometimes a client package that Confluent documents.

Apache Kafka docs are the source for broker keys and KRaft. Confluent tutorials often start extra containers. Those are not required to learn topics 1–6. Read licenses for extra connectors.

**Redpanda** is a streaming system that speaks the Kafka API. Many Kafka clients work. The engine is not the Kafka JVM broker. Operations, disk layout, and extra features differ. This handbook teaches Apache Kafka. Use Redpanda when the team chooses that product and accepts the differences.

**Apache Pulsar** is a different pub/sub system. It uses a separate protocol by default. Pulsar is not "Kafka with another name". Do not run a Pulsar tutorial and think you learned Kafka partitions and KRaft.

**Amazon MSK** is Kafka as a managed service on AWS. New MSK clusters can use KRaft. You still use topics, ACLs, and clients.

**Confluent Cloud**, **Aiven**, and similar products give a bootstrap and a project. **WarpStream** stores data in object storage. The operations model differs from classic `log.dirs`. Clients still use the Kafka protocol. Test compatibility for your features.

Shared managed facts:

- You do not run `kafka-server-start` on your laptop for production.
- Choose KRaft when the vendor still lists ZooKeeper-based old versions.
- Networking (private link, VPC) is part of the design.
- Do not assume transactions, Connect, or Streams run inside the vendor without a vendor page.

If a vendor says "Kafka compatible", test: produce, consume, groups, transactions (if you need EOS), and admin.

KRaft is an Apache Kafka metadata mode. Redpanda and Pulsar have their own metadata.

### Questions

#### Theoretical questions

1. What do you get from Apache Kafka alone?
2. What extra kinds of products does Confluent Platform add?
3. What does Redpanda share with Kafka from a client view?
4. What does MSK manage for you?
5. Why do you choose KRaft when a vendor still lists ZooKeeper versions?

#### Easy practical tasks

1. Write five sentences that contrast Apache Kafka and Confluent Platform. Use only facts from this section.
2. Make a table: Item, Apache Kafka, Redpanda, one managed service.
3. Open kafka.apache.org and one cloud getting-started page. Write the bootstrap idea each uses.
4. List three Kafka features you would test on a clone (groups, transactions, ACLs).

#### Medium practical tasks

1. Compare an Apache Docker quick start with a Confluent compose file. Write five extra services.
2. Optional: run Redpanda in Docker. Point a Kafka console consumer at it. Write what worked.
3. Read a vendor Kafka API compatibility page. Write three supported and one limited feature if listed.

#### Advanced practical tasks

1. Write a cloud choice note: protocol compatibility, KRaft, registry, Connect, price drivers, VPC, exit plan.
2. Write a platform standard: core is Apache Kafka protocol on KRaft, allowed add-ons, forbidden ZooKeeper SKUs.

---

## Schema Registry, REST Proxy, and UIs

**Schema Registry** (Confluent or Karapace, topic 7) stores schemas and compatibility rules. It is an HTTP service. It is not the broker. It often uses a compacted `_schemas` topic on Kafka.

**REST Proxy** (Confluent REST Proxy and similar projects) exposes HTTP for produce, consume, and some admin. Use it when a language has no good client or when a legacy system can only call HTTP. Prefer a native Kafka client when you can. The proxy is another process to secure (topic 10).

**Control Center** is a Confluent UI for topics, Connect, ksql, and metrics. It is not part of Apache Kafka.

**AKHQ** and **Kafdrop** are open-source browser UIs. They list topics, browse records, and inspect groups. They need bootstrap and credentials. A UI that consumes a topic is a consumer. It can change lag or use a group if you configure it that way. Use a dedicated group id. Do not attach a UI to a production application group.

UIs do not replace `kafka-topics` and metrics. They help beginners see partitions and offsets.

All of these talk to a Kafka cluster. That cluster must be KRaft for new work. None of them require you to install ZooKeeper.

Pin UI versions. Do not expose a UI to the public internet without auth.

### Questions

#### Theoretical questions

1. What problem does Schema Registry solve?
2. When is a REST Proxy useful?
3. Why can a browse UI break a production consumer group?
4. Is Control Center part of Apache Kafka?
5. Why must a UI have its own group id?

#### Easy practical tasks

1. Make a table: Tool, protocol (HTTP or Kafka), writes to cluster (yes or no or maybe).
2. Write five sentences about UIs and groups. Use only facts from this section.
3. Find AKHQ or Kafdrop documentation. Write the bootstrap setting name.
4. Draw: browser → UI → Kafka, and browser → REST Proxy → Kafka.

#### Medium practical tasks

1. Run one UI (AKHQ, Kafdrop, or Control Center) against local KRaft. Browse a topic. Write the group id that the UI used if you can see it.
2. Produce with a native client. Consume with REST Proxy (if you run it) or refuse and write why you skipped it.
3. Compare registry UI (if any) with `curl` to the registry API. Write two operations.

#### Advanced practical tasks

1. Secure a UI with the same SASL/TLS as clients (topic 10). Write how you prevented public access.
2. Write an ecosystem-tool standard: allowed UIs, group id policy, registry product, no ZooKeeper.

---

## Produce and consume in your language; then add a Streams job

Learn in this order. Do not start with a stream cluster before you can produce and consume.

1. Start a local **KRaft** Kafka (topic 1). Create a topic with three partitions (topic 2).
2. Use the **console producer and consumer**. Confirm records and `--from-beginning`.
3. Write a **producer** in a language that you know (Java, Go, Python, or another official client). Set `bootstrap.servers`, serializers, `acks=all`, and the idempotent producer (topic 3).
4. Write a **consumer** in the same language. Set `group.id`, deserializers, auto-commit off, and a poll loop (topic 4). Commit after you finish the work.
5. Run **two consumers in one group**. Stop one. Confirm that partitions move. Watch lag.
6. Add an **event id** and an idempotent sink, or a small transaction if you need EOS (topic 5).
7. Use **JSON first**. Then add a Schema Registry and Avro or JSON Schema if your team needs a contract (topic 7).
8. Optional: add **Connect** (JDBC or Debezium) when you must copy a database (topic 8).
9. Then add a **Streams job**: a Java Kafka Streams topology, or a ksqlDB persistent query, or a small Flink job (topic 9). Start with a filter. Then add a count per key. Then add a window if you need time buckets.
10. Optional: add **TLS and SASL** on a local compose stack (topic 10).

A Streams job is a second client. It needs its own `application.id` (or Flink job id). It creates internal topics. Do not reuse the consumer `group.id` from step 4 for a different topology.

The official resources stay the same:

- [Apache Kafka Documentation](https://kafka.apache.org/documentation/)
- [Introduction](https://kafka.apache.org/intro)
- [Confluent Developer](https://developer.confluent.io/)

Practice on KRaft. Do not start ZooKeeper for this path.

### Questions

#### Theoretical questions

1. Why must you produce and consume before you start a Streams job?
2. What configuration must a first producer include for important data?
3. What configuration must a first consumer include for important data?
4. Why must a Streams `application.id` not reuse a random old group id?
5. Which official site is the source for broker keys and KRaft?

#### Easy practical tasks

1. Write the ten-step order from this section as a personal checklist. Mark what you already did.
2. Produce and consume in your language on topic `lang.demo`. Write the client library name and version.
3. Make a table: Step, topic number in this path, official URL.
4. Bookmark the three official resources named in this section.

#### Medium practical tasks

1. Add a second consumer in the same group. Stop one process. Save `kafka-consumer-groups --describe` before and after.
2. Add a header `event-id`. Consume and skip a duplicate id. Write the test.
3. Write a one-page plan for your first Streams job: input topic, filter or count, output topic, `application.id`, KRaft bootstrap.

#### Advanced practical tasks

1. Implement the plan: a small Streams (or ksqlDB) job on KRaft. List internal topics. Consume the output.
2. Write a six-month learning plan: language client, Connect optional, Streams, security, one MM2 lab. No ZooKeeper.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. What is "Kafka" in this handbook versus what is an add-on in the ecosystem?
2. How do MM2 lag, RPO, and offset translation form one failover story?
3. How do Redpanda and Confluent Cloud both "speak Kafka" but differ in who runs the broker code?
4. Why can a UI, a REST Proxy, and a Schema Registry each create a security review (topic 10)?
5. How does the produce-consume-then-Streams order protect a beginner from a stuck stream job?

#### Easy practical tasks

1. Write a one-page cheat sheet: MM2, DR terms, Apache vs Confluent vs Redpanda vs managed, registry, REST, UI, learning order.
2. Draw your preferred learning stack: KRaft Kafka + language client + optional registry + optional UI. No ZooKeeper.
3. Bookmark Apache Kafka, one UI, one cloud, one MM2 page.
4. List every product named in this topic in one column and "learn now / later / skip" in the second.

#### Medium practical tasks

1. Install one add-on (UI or registry) on local KRaft. Write the exact bootstrap and one thing the add-on cannot do.
2. Map each subsection to one official URL.
3. Write a table: Feature (EOS, Connect, Streams, ACLs). Mark Apache, one cloud, Redpanda as yes / no / check.

#### Advanced practical tasks

1. Produce from an Apache client to two backends if you can (local KRaft and one trial cloud or Redpanda). Write config diffs only (no secrets).
2. Write an ecosystem and DR standard for your company: core protocol on KRaft, allowed add-ons, MM2 or vendor link, forbidden ZooKeeper, review for every new UI.
