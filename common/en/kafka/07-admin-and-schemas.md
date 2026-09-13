# 7. Admin and Schemas

## Description

Admin work creates and changes topics, sets configuration, and inspects groups. Kafka stores keys and values as bytes. A schema describes the shape of those bytes. A Schema Registry stores schemas and can enforce compatibility.

This topic shows create, describe, and delete of topics, internal topics, bytes on the wire, Schema Registry, Avro, Protobuf, JSON Schema, and compatibility levels. Complete this topic after topics 2, 3, and 6. You do not need ZooKeeper. Admin tools use `--bootstrap-server`. Schema Registry is a separate service. KRaft Kafka and Schema Registry work together through the Kafka protocol and HTTP.

Use one term for each concept. An admin operation changes cluster metadata or topic configuration. An internal topic is a topic that Kafka uses for its own state. Bytes on the wire are the key and value payloads in a record. A schema is a contract for those bytes. A subject is the registry name that holds a sequence of schema versions. Compatibility is the rule that allows or rejects a new schema version.

---

## Create, describe, and delete topics

You create a topic with a name, a partition count, and a replication factor. Example command shape:

```text
kafka-topics --bootstrap-server localhost:9092 --create --topic orders.placed --partitions 3 --replication-factor 1
```

On Windows, use the `.bat` file. On macOS and Linux, use the `.sh` file. Some installs put the tools on `PATH` without a suffix.

`--describe` shows partitions, leaders, replicas, ISR, and configuration when you ask for configs.

`--list` shows topic names.

`--delete` removes a topic. The cluster must allow delete (`delete.topic.enable` or the current equivalent). Delete is not instant. Data leaves disk after the delete process. Do not delete a topic in production without a change process.

Auto-create can exist on a broker (`auto.create.topics.enable`). Do not rely on auto-create in production. Create topics with an explicit partition count and replication factor.

Topic names have rules. Use letters, digits, `.`, `_`, and `-`. Avoid spaces. Avoid names that collide with internal prefixes.

Broker defaults apply to all topics. You can override keys on one topic. Examples: `retention.ms`, `cleanup.policy`, `min.insync.replicas`, `max.message.bytes`. Set overrides at create time (`--config`) or later with `kafka-configs`.

Describe with configs to see the effective values and whether a value is a default or a topic override. Use overrides for special topics. Do not create a unique config for every topic without a reason.

The Admin API in your language does the same operations from a program or a pipeline.

Do not use old ZooKeeper flags on new KRaft clusters.

### Questions

#### Theoretical questions

1. What three properties do you set when you create a topic?
2. What does `--describe` show?
3. Why must you not rely on auto-create in production?
4. Why is delete not instant?
5. What address flag do new tools use instead of ZooKeeper?

#### Easy practical tasks

1. Create `admin.demo` with 3 partitions. Describe it. List topics. Write the commands.
2. Write four sentences about topic admin. Use only facts from this section.
3. Find `--help` for `kafka-topics`. Write five flags.
4. Try an invalid topic name. Record the error.

#### Medium practical tasks

1. Delete `admin.demo`. List topics until the name is gone. Write how long it took.
2. Create a topic with `retention.ms=60000`. Describe the configs. Then alter the value with `kafka-configs` if you can. Write before and after.
3. Write a small Admin API program that creates and describes a topic.

#### Advanced practical tasks

1. Create a topic with RF 3 on a three-broker cluster. Describe. Then delete. Confirm that log directories drop the partition folders after delete.
2. Write a one-page change process: who can create topics, required RF, required min ISR, and how you record the change.

---

## Internal topics

Kafka creates **internal topics** for its own state.

`__consumer_offsets` stores committed offsets for consumer groups. Do not produce application events to this topic. Do not delete it. The group coordinator uses it (topic 4).

Transaction state uses internal topics (names include transaction state in official docs, for example `__transaction_state`). EOS producers need them (topic 5). Do not delete them.

Other internal names can appear. Cluster metadata in KRaft uses a metadata log on controllers. That log is not a normal application topic. Connect and Streams can add more topics when you use those products (topics 8 and 9).

`--list` can show internal topics. Some tools hide them unless you ask. Treat every name that starts with `__` as reserved unless the documentation says otherwise.

Internal topics still use partitions, replicas, and retention. They need disk and a correct replication factor. On a one-broker lab, they use replication factor 1. On production, they need a production RF.

Do not compact or delete internal topics as an experiment. You can break groups and transactions.

### Questions

#### Theoretical questions

1. What does `__consumer_offsets` store?
2. Why must you not produce application events to an internal topic?
3. What product features need transaction internal topics?
4. How is the KRaft metadata log different from an application topic?
5. Why does a production RF still apply to internal topics?

#### Easy practical tasks

1. List topics including internals if your tool supports that flag. Write the names that start with `__`.
2. Write four sentences about internal topics. Use only facts from this section.
3. Describe `__consumer_offsets` if it exists. Write partition count and RF.
4. Make a table: Name or kind, who writes, who reads.

#### Medium practical tasks

1. After you run a consumer group, describe `__consumer_offsets` again. Write what changed if anything is visible.
2. Find official documentation on internal topics. Write three names and one sentence each.
3. Read how Connect or Streams names internal topics (preview). Write one prefix or pattern.

#### Advanced practical tasks

1. Write a one-page operations note: which internal topics you monitor, which you never delete, and how RF is set on a new KRaft cluster.
2. On a three-broker cluster, confirm RF of `__consumer_offsets`. Write whether it matches your production standard.

---

## Bytes on the wire

A Kafka record key is a byte sequence or null. A Kafka record value is a byte sequence or null. Brokers do not parse JSON or Avro. Brokers do not check your fields.

The producer serializer creates the bytes (topic 3). The consumer deserializer must understand the same format. If the pair does not match, the consumer throws or prints garbage.

Plain formats:

- string (UTF-8 text)
- raw JSON text
- integer or long in a fixed binary layout (client serializers)

With Confluent-style Schema Registry, the value often uses this layout:

1. a magic byte
2. a 4-byte schema id
3. the Avro, Protobuf, or JSON payload

The broker still sees bytes. The registry maps the id to a schema. The deserializer fetches the schema (and caches it) and then decodes the payload.

Do not mix a raw JSON producer and a Registry Avro consumer on the same topic. Create a clear contract per topic.

Headers can carry a content type. Headers are optional. Do not rely on a header if the deserializer does not read it.

### Questions

#### Theoretical questions

1. What does Kafka store for a key and a value?
2. Does the broker validate JSON fields?
3. What happens when the serializer and the deserializer do not match?
4. What is the Confluent-style wire layout at a high level?
5. Why is a content-type header not enough by itself?

#### Easy practical tasks

1. Write five sentences about bytes on the wire. Use only facts from this section.
2. Produce a JSON string. Consume as a string. Write the exact text.
3. Make a table: Format, human readable on consume, registry id on the wire.
4. Draw magic byte + schema id + payload.

#### Medium practical tasks

1. Produce UTF-8 text and consume with a binary dump (hex) if you can. Write the first bytes.
2. Find the official or Confluent wire-format documentation. Write the size of the schema id field.
3. Write a producer and consumer with mismatched deserializers on purpose. Record the error. Fix the pair.

#### Advanced practical tasks

1. After you complete the registry sections, parse a Registry-encoded value and print the schema id. Confirm it in the registry API.
2. Write a one-page contract: topic name, key format, value format, and whether Registry is required.

---

## Schema Registry, Avro / Protobuf / JSON Schema

A **Schema Registry** is an HTTP service that stores schemas. Producers register a schema (or use a known id). Consumers load the schema by id.

**Confluent Schema Registry** is the common product in Confluent Platform and Confluent Cloud. **Karapace** is an open-source registry that speaks a compatible API. Other registries exist. This handbook uses "the registry" for the service role.

The registry is not Kafka. You start it next to the cluster. It often stores its own data in a Kafka topic (for example `_schemas`). That topic is a compacted changelog. Do not treat it as an application topic.

Clients need `schema.registry.url` (name can differ). Authentication on the registry is separate from Kafka SASL. Topic 10 covers security.

Without a registry, you can still use Avro or Protobuf files in your repo. Then every consumer must have the schema file. The registry helps when many teams share topics and when schemas evolve.

You can learn Kafka topics 1–6 without a registry. Use a registry when you need a shared contract.

**Avro.** A compact binary encoding. The schema is JSON. Avro is common in Kafka and Hadoop-style systems. Readers need the writer schema. The registry stores it. Avro has a simple evolution model that matches registry compatibility well.

**Protobuf.** A compact binary encoding. You write `.proto` files. Many services already use Protobuf for RPC. Kafka clients can use Protobuf with a registry. Evolution rules differ from Avro. Field numbers matter.

**JSON Schema.** A schema for JSON documents. Payloads stay readable as text (plus the registry wrapper if you use the registry). Evolution and validation rules follow JSON Schema and the registry mode.

Pick one format per topic. Do not mix Avro and JSON on the same subject without a plan.

Do not run a production registry as a single un-backed process without a plan. Follow the product operations guide.

### Questions

#### Theoretical questions

1. What does a Schema Registry store?
2. How does a consumer get the schema for a record?
3. What is Karapace?
4. How is Avro different from JSON Schema at a high level?
5. Why do Protobuf field numbers matter?

#### Easy practical tasks

1. Write four sentences about Schema Registry. Use only facts from this section.
2. Make a table: Format, binary or text, typical schema file.
3. Open the Confluent Schema Registry documentation. Write the base URL of the docs page.
4. Open the Karapace documentation. Write one sentence about API compatibility.

#### Medium practical tasks

1. Start a registry (Confluent container, Karapace, or a cloud trial). List subjects (empty is fine). Write the command or HTTP call.
2. Find the `_schemas` topic or the equivalent. Describe it. Write the cleanup policy if shown.
3. Register one Avro, Protobuf, or JSON Schema by HTTP if you can. Fetch it by id. Write the requests (no secrets).

#### Advanced practical tasks

1. Produce and consume one Avro or JSON Schema record through the registry on KRaft Kafka. Write the subject name and the schema id.
2. Write a one-page operations note: registry high availability, the `_schemas` topic RF, and who may register.

---

## Compatibility levels

**Compatibility** is the rule that the registry uses when you register a new schema version for a subject. The registry accepts or rejects the new version.

Common levels (names follow Confluent-style registries; confirm on your product):

- **BACKWARD.** A new schema can read data that was written with the previous schema. Add optional fields. Do not remove fields that old writers still send unless the rules allow it. This is a common default.
- **FORWARD.** Old readers can read data written with the new schema. The new schema is more careful about required fields.
- **FULL.** Both backward and forward for the pair of versions that the mode defines.
- **NONE.** The registry does not check. You can break readers.

Some products add `*_TRANSITIVE` modes. Those modes check the new schema against all previous versions, not only the last one.

**Subject naming** is how the registry names the sequence of versions. Common strategies: topic name plus `-value` or `-key`, or a record name. All producers of one topic must agree.

A rejected register is a success for safety. Do not set `NONE` in production to "unblock" a deploy. Fix the schema or add a new topic.

Consumers that cache schemas must handle a new id. The wire format carries the id. The consumer fetches the new schema.

Compatibility does not replace tests. Write a consumer test with an old payload and a new payload.

### Questions

#### Theoretical questions

1. What does a compatibility level control?
2. What does BACKWARD allow a new reader to do?
3. What does FORWARD protect?
4. What is a subject?
5. Why is `NONE` a risk in production?

#### Easy practical tasks

1. Write five sentences about compatibility. Use only facts from this section.
2. Make a table: Level, new reader vs old data, old reader vs new data.
3. Find the default compatibility in the registry docs that you use.
4. Write a subject name for topic `orders.placed` value under the topic-name strategy.

#### Medium practical tasks

1. Register schema v1. Change a field (add optional, or remove a field). Try to register v2 under BACKWARD. Write accept or reject.
2. Repeat under NONE or a weaker mode only in a lab. Write the difference. Restore a safe mode.
3. Read official compatibility examples for Avro. Rewrite two examples in STE.

#### Advanced practical tasks

1. Write a team standard: default BACKWARD (or FULL), subject naming, who may set NONE, and how you test old consumers.
2. Plan a breaking change: new topic versus a compatibility break. Write when you create `orders.placed.v2`.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe the path from `kafka-topics --create` to a consumer that decodes a Registry Avro value. Name admin, bytes, registry, and compatibility.
2. Why can an internal topic and an application topic share the same cluster but must not share the same naming style?
3. What stays the same when you move from raw JSON to Avro with a registry?
4. How do topic config overrides and schema compatibility both protect readers, in different ways?
5. A blog uses `--zookeeper` to create topics. What do you change on a KRaft cluster?

#### Easy practical tasks

1. Create `schema.review` with 3 partitions. Describe configs. Produce a JSON string. Consume it. Write the commands.
2. Write a one-page cheat sheet: create/describe/delete, internal topics, wire bytes, registry, three formats, compatibility levels.
3. List topics. Mark internal names. Mark your application names.
4. Bookmark official admin docs and one Schema Registry docs home.

#### Medium practical tasks

1. Alter retention on a lab topic. Register a schema if you have a registry. Write both change procedures.
2. Map each subsection to one official URL.
3. Compare `kafka-topics` with the Admin API for create and describe. Write five equivalent operations.

#### Advanced practical tasks

1. Build a small pipeline: Admin API creates a topic, a producer writes Avro or JSON Schema through a registry, a consumer decodes. Record subject, id, and compatibility mode.
2. Write a platform standard: topic create checklist, forbidden deletes, registry product, default compatibility, KRaft bootstrap only, no ZooKeeper flags.
