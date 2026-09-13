# 21. Ecosystem

## Description

The Kafka ecosystem is the set of products that speak the Kafka protocol or that sit next to a Kafka cluster. Apache Kafka is the open-source broker and the protocol. Other products add servers, clouds, and user interfaces.

This topic covers Confluent Platform versus Apache Kafka, Redpanda and Apache Pulsar as comparison points, Schema Registry, REST Proxy, Control Center / AKHQ / Kafdrop, and cloud services (MSK, Confluent Cloud, Aiven, WarpStream). Complete this topic after topics 1, 11, 12, and 14. Prefer KRaft for Apache Kafka and for managed Kafka that offers it. Do not start ZooKeeper for new work.

Use one term for each concept. Apache Kafka is the ASF project. A Kafka-API compatible system accepts Kafka clients. A managed service runs brokers for you. A UI is a browser tool that calls admin and consumer APIs. A registry stores schemas (topic 11). A REST Proxy exposes HTTP for produce and consume.

---

## Confluent Platform vs Apache Kafka

**Apache Kafka** is the project at [https://kafka.apache.org/](https://kafka.apache.org/). You download brokers, client libraries, Kafka Streams, and Kafka Connect. You operate them. New clusters use KRaft.

**Confluent Platform** is a distribution from Confluent. It includes Apache Kafka (or a compatible broker) plus extra components that Confluent maintains: Schema Registry, ksqlDB, Control Center, connectors, and commercial features. **Confluent Cloud** is the hosted form.

The Kafka protocol is the same idea. A producer that uses the Apache client can talk to Confluent Cloud if you set bootstrap, security, and sometimes a client package that Confluent documents.

Differences that matter for a beginner:

- Apache Kafka docs are the source for broker keys and KRaft.
- Confluent tutorials often start extra containers (registry, ksql). Those are not required to learn topics 1–10.
- Licenses differ. Some Confluent connectors and UI features are not Apache 2.0. Read the license before you install.

Do not treat "Confluent" and "Apache Kafka" as two protocols. Treat them as a core plus optional products.

You can run Apache Kafka on KRaft and add Karapace (topic 11) instead of Confluent Schema Registry. That is a valid Apache-centered path.

### Questions

#### Theoretical questions

1. What do you get from Apache Kafka alone?
2. What extra kinds of products does Confluent Platform add?
3. Why can an Apache client talk to Confluent Cloud?
4. Why must you read licenses for extra connectors?
5. Must you run ksqlDB to learn Kafka topics 1–10?

#### Easy practical tasks

1. Write five sentences that contrast Apache Kafka and Confluent Platform. Use only facts from this section.
2. Make a table: Component, in Apache tarball (yes or no), typical Confluent add-on.
3. Open kafka.apache.org and confluent.io (or developer.confluent.io). Write three facts that match and one extra product.
4. List two Confluent components you already met in this path (registry, ksql, Connect).

#### Medium practical tasks

1. Compare an Apache Docker quick start with a Confluent Platform compose file. Write five extra services.
2. Find the license page for one Confluent connector you care about. Write the license name.
3. Write a team sentence: "Our source of truth for broker config is Apache docs; we add Confluent X for Y."

#### Advanced practical tasks

1. Run Apache KRaft Kafka only, then add one Confluent or Karapace registry. Write what stayed Apache.
2. Write a platform standard: Apache Kafka on KRaft, which Confluent products you allow, and a license review step.

---

## Redpanda, Apache Pulsar (comparison points)

**Redpanda** is a streaming system that speaks the Kafka API. Many Kafka clients work. The engine is not the Kafka JVM broker. Operations, disk layout, and extra features differ. This handbook teaches Apache Kafka. Use Redpanda when the team chooses that product and accepts the differences.

**Apache Pulsar** is a different pub/sub and streaming system. It uses a separate protocol by default. It has a different model (brokers plus bookies, segmented storage). Some tools bridge Pulsar and Kafka. Pulsar is not "Kafka with another name".

Comparison points (high level):

- Protocol: Redpanda ≈ Kafka API. Pulsar ≠ Kafka API (unless a proxy).
- Process model: Kafka JVM brokers (KRaft). Redpanda is its own binary. Pulsar is brokers + storage layer.
- Streams: Kafka Streams is for Kafka. Pulsar has Functions. Redpanda has its own extras.
- Operations: tools (`kafka-topics`) work on Kafka and often on Redpanda. They do not manage Pulsar.

Do not run a Pulsar tutorial and think you learned Kafka partitions and KRaft. Do not assume every Kafka Connect connector works on every Kafka-API clone.

If a vendor says "Kafka compatible", test: produce, consume, groups, transactions (if you need EOS), and admin.

KRaft is an Apache Kafka metadata mode. Redpanda and Pulsar have their own metadata. Do not look for `kafka-storage format` on Pulsar.

### Questions

#### Theoretical questions

1. What does Redpanda share with Kafka from a client view?
2. Why is Pulsar not a Kafka clone?
3. What must you test when a vendor says "compatible"?
4. Do `kafka-topics` manage Pulsar?
5. Why does this handbook still use KRaft language for Apache Kafka only?

#### Easy practical tasks

1. Write four sentences about Redpanda and four about Pulsar. Use only facts from this section.
2. Make a table: Item, Apache Kafka, Redpanda, Pulsar.
3. Open the Redpanda Kafka-compatibility page and the Pulsar home. Write one sentence each.
4. List three Kafka features you would test on a clone (groups, transactions, ACLs).

#### Medium practical tasks

1. Optional: run Redpanda in Docker. Point a Kafka console consumer at it. Write what worked.
2. Read a Pulsar versus Kafka comparison from official Pulsar or Kafka docs (not a random blog). Write five factual differences.
3. Write when you would refuse to call a system "just Kafka".

#### Advanced practical tasks

1. Write a one-page comparison for your team: API, EOS, Connect, Streams, operations, cloud. Mark unknown as unknown.
2. Write a rule: learning path stays on Apache Kafka KRaft; clones are an optional lab after topic 22.

---

## Schema Registry, REST Proxy, Control Center / AKHQ / Kafdrop

**Schema Registry** (Confluent or Karapace, topic 11) stores schemas and compatibility rules. It is an HTTP service. It is not the broker.

**REST Proxy** (Confluent REST Proxy and similar projects) exposes HTTP for produce, consume, and some admin. Use it when a language has no good client or when a legacy system can only call HTTP. Prefer a native Kafka client when you can. The proxy is another process to secure (topic 17).

**Control Center** is a Confluent UI for topics, Connect, ksql, and metrics. It is not part of Apache Kafka.

**AKHQ** and **Kafdrop** are open-source browser UIs. They list topics, browse records, and inspect groups. They need bootstrap and credentials. A UI that consumes a topic is a consumer. It can change lag or use a group if you configure it that way. Use a dedicated group id. Do not attach a UI to a production application group.

UIs do not replace `kafka-topics` and metrics. They help beginners see partitions and offsets.

All of these talk to a Kafka cluster. That cluster should be KRaft for new work. None of them require you to install ZooKeeper.

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

1. Secure a UI with the same SASL/TLS as clients (topic 17). Write how you prevented public access.
2. Write an ecosystem-tool standard: allowed UIs, group id policy, registry product, no ZooKeeper.

---

## Cloud: MSK, Confluent Cloud, Aiven, WarpStream

**Amazon MSK** is Kafka as a managed service on AWS. You choose version, size, and networking. New MSK clusters can use KRaft. You still use topics, ACLs, and clients. You do not SSH to a broker as the default.

**Confluent Cloud** is Confluent’s hosted Kafka plus registry, ksql, and other cloud products. Bootstrap is a hostname that Confluent gives you. Auth is often API keys or OAUTHBEARER (topic 17).

**Aiven** is a managed data platform that includes Kafka (and other stores). You get a bootstrap and a project. Features follow the Aiven Kafka offering.

**WarpStream** is a Kafka-protocol service that stores data in object storage. The operations model differs from a classic broker disk in `log.dirs`. Clients still use the Kafka protocol. Test compatibility for your features.

Shared cloud facts:

- You do not run `kafka-server-start` on your laptop for production.
- You still need KRaft awareness when the vendor offers ZooKeeper-based old versions. Choose KRaft.
- Networking (private link, VPC) is part of the design.
- Price includes partitions, storage, and cross-zone traffic (topic 16).
- Vendors add their own UI. The protocol stays Kafka.

Do not assume transactions, exactly-once, or a specific ACL model without a vendor page. Do not assume Connect or Streams run inside the vendor (sometimes you run them yourself).

This handbook stays on Apache Kafka concepts. Cloud is a deployment of those concepts.

### Questions

#### Theoretical questions

1. What does MSK manage for you?
2. What extra cloud products does Confluent Cloud often include?
3. Why must you still learn topics and consumer groups on a managed service?
4. How is WarpStream different at a storage level (high level)?
5. Why do you choose KRaft when a vendor still lists ZooKeeper versions?

#### Easy practical tasks

1. Make a table: Service, cloud vendor or model, you run brokers (yes or no).
2. Write four sentences about what stays the same on every managed Kafka.
3. Open the current MSK and Confluent Cloud "getting started" pages. Write the bootstrap idea each uses.
4. List three costs (partitions, storage, network) that topic 16 already prepared you to ask.

#### Medium practical tasks

1. If you have a trial: create a KRaft (or vendor-current) cluster. Produce and consume from your machine. Write the security mechanism.
2. Compare MSK Connect or Confluent Cloud connectors with self-managed Connect (topic 12). Write five differences.
3. Read WarpStream or Aiven Kafka docs for Kafka API compatibility. Write three supported and one limited feature if listed.

#### Advanced practical tasks

1. Write a cloud choice note: protocol compatibility, KRaft, registry, Connect, price drivers, VPC, exit plan (how you leave).
2. Write a standard: new environments use managed KRaft Kafka or Apache KRaft; ZooKeeper SKUs are forbidden; client settings you lock.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. What is "Kafka" in this handbook versus what is an add-on in the ecosystem?
2. How do Redpanda and Confluent Cloud both "speak Kafka" but differ in who runs the broker code?
3. Why can a UI, a REST Proxy, and a Schema Registry each create a security review (topic 17)?
4. What do you check before you move a lab from Apache Docker to MSK?
5. How does KRaft appear (or not appear) in Apache, Confluent, Redpanda, and Pulsar?

#### Easy practical tasks

1. Write a one-page cheat sheet: Apache vs Confluent, two comparison systems, five side tools, four clouds.
2. Draw your preferred learning stack: KRaft Kafka + one UI + optional registry. No ZooKeeper.
3. Bookmark Apache Kafka, one UI, one cloud, one compatibility page.
4. List every product named in this topic in one column and "learn now / later / skip" in the second.

#### Medium practical tasks

1. Install one add-on (UI or registry) on local KRaft. Write the exact bootstrap and one thing the add-on cannot do.
2. Map each subsection to one official URL.
3. Write a table: Feature (EOS, Connect, Streams, ACLs). Mark Apache, one cloud, Redpanda as yes / no / check.

#### Advanced practical tasks

1. Produce from an Apache client to two backends if you can (local KRaft and one trial cloud or Redpanda). Write config diffs only (no secrets).
2. Write an ecosystem standard for your company: core is Apache Kafka protocol on KRaft, allowed add-ons, forbidden ZooKeeper, review for every new UI.
