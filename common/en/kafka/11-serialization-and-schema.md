# 11. Serialization and Schema

## Description

Kafka stores keys and values as bytes. A serializer writes those bytes. A deserializer reads them. A schema describes the shape of the data. A Schema Registry stores schemas and can enforce compatibility.

This topic covers bytes on the wire, Schema Registry (Confluent or Karapace), Avro, Protobuf, JSON Schema, compatibility modes, and subject naming. Complete this topic after topic 4. You do not need ZooKeeper. Schema Registry is a separate service. KRaft Kafka and Schema Registry work together through the Kafka protocol and HTTP.

Use one term for each concept. Bytes on the wire are the key and value payloads in a record. A schema is a contract for those bytes. A subject is the registry name that holds a sequence of schema versions. Compatibility is the rule that allows or rejects a new schema version.

---

## Bytes on the wire

A Kafka record key is a byte sequence or null. A Kafka record value is a byte sequence or null. Brokers do not parse JSON or Avro. Brokers do not check your fields.

The producer serializer creates the bytes. The consumer deserializer must understand the same format. If the pair does not match, the consumer throws or prints garbage.

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

1. Parse a Registry-encoded value (after you complete the next sections) and print the schema id. Confirm it in the registry API.
2. Write a one-page contract: topic name, key format, value format, and whether Registry is required.

---

## Schema Registry (Confluent or Karapace)

A Schema Registry is an HTTP service that stores schemas. Producers register a schema (or use a known id). Consumers load the schema by id.

**Confluent Schema Registry** is the common product in Confluent Platform and Confluent Cloud. **Karapace** is an open-source registry that speaks a compatible API. Other registries exist. This handbook uses "the registry" for the service role.

The registry is not Kafka. You start it next to the cluster. It often stores its own data in a Kafka topic (for example `_schemas`). That topic is a compacted changelog. Do not treat it as an application topic.

Clients need `schema.registry.url` (name can differ). Authentication on the registry is separate from Kafka SASL. Topic 17 covers security.

Without a registry, you can still use Avro or Protobuf files in your repo. Then every consumer must have the schema file. The registry helps when many teams share topics and when schemas evolve.

Do not run a production registry as a single un-backed process without a plan. Follow the product operations guide.

You can learn Kafka topics 1–10 without a registry. Use a registry when you need a shared contract.

### Questions

#### Theoretical questions

1. What does a Schema Registry store?
2. How does a consumer get the schema for a record?
3. What is Karapace?
4. How is the registry different from a Kafka broker?
5. What Kafka topic does a registry often use for its own data?

#### Easy practical tasks

1. Write four sentences about Schema Registry. Use only facts from this section.
2. Open the Confluent Schema Registry documentation. Write the base URL of the docs page.
3. Open the Karapace documentation. Write one sentence about API compatibility.
4. Make a table: With registry, without registry. Add schema share and evolution.

#### Medium practical tasks

1. Start a registry (Confluent container, Karapace, or a cloud trial). List subjects (empty is fine). Write the command or HTTP call.
2. Find the `_schemas` topic or the equivalent. Describe it. Write the cleanup policy if shown.
3. Point a client at a wrong registry URL. Record the error. Fix the URL.

#### Advanced practical tasks

1. Register one Avro or JSON Schema by HTTP. Fetch it by id. Write the requests (no secrets).
2. Write a one-page operations note: registry high availability, the `_schemas` topic RF, and who may register.

---

## Avro / Protobuf / JSON Schema

These three schema systems are common with a registry.

**Avro.** A compact binary encoding. The schema is JSON. Avro is common in Kafka and Hadoop-style systems. Readers need the writer schema. The registry stores it. Avro has a simple evolution model that matches registry compatibility well.

**Protobuf.** A compact binary encoding. You write `.proto` files. Many services already use Protobuf for RPC. Kafka clients can use Protobuf with a registry. Evolution rules differ from Avro. Field numbers matter.

**JSON Schema.** The payload can stay JSON text (or a binary form, depending on the serializer). JSON Schema describes fields and types. Debug is easier. Size is larger than Avro or Protobuf. Validation can occur in the serializer.

Pick one format per topic. Do not change format in place. Use a new topic if you must change the family (Avro to Protobuf).

JSON without a schema is not JSON Schema. Raw JSON has no registry check.

Code generation is common. Avro and Protobuf generate types. Generated types reduce field-name errors. Commit the schema source (`.avsc` or `.proto`) in the repository that owns the topic.

### Questions

#### Theoretical questions

1. What is Avro in one sentence?
2. What is Protobuf in one sentence?
3. What is JSON Schema in one sentence?
4. Why must one topic use one format family?
5. How is raw JSON different from JSON Schema?

#### Easy practical tasks

1. Make a three-column table: Avro, Protobuf, JSON Schema. Add rows for readable payload, typical size, schema source file.
2. Write five sentences that compare the three formats. Use only facts from this section.
3. Find a serializer class name for each format in your client docs (or Confluent docs).
4. List two topics that fit Avro and one topic that might keep JSON Schema for debug.

#### Medium practical tasks

1. Define a small `Order` schema in one format. Produce and consume one record with a registry (or a file schema if you cannot run a registry).
2. Add an optional field. Produce an old producer and a new consumer, or the reverse, as your compatibility mode allows. Write the result.
3. Read official or Confluent format guides. Write one evolution rule for Avro and one for Protobuf.

#### Advanced practical tasks

1. Encode the same logical order in Avro and in JSON. Compare byte sizes for 1000 records.
2. Write a team standard: default format, when you allow JSON Schema, and when you forbid raw JSON on shared topics.

---

## Compatibility: backward, forward, full

The registry can reject a new schema if it breaks the subject compatibility mode.

**Backward.** A new schema can read data that writers produced with the previous schema. New consumers can read old records. This is the common default. Example: you add an optional field with a default.

**Forward.** An old schema can read data that a new writer produces. Old consumers can read new records. Example: you delete a field that old readers never required.

**Full.** Both backward and forward hold.

**None.** The registry accepts any schema. Consumers can break. Avoid this on shared topics.

Transitive modes (`BACKWARD_TRANSITIVE`, `FORWARD_TRANSITIVE`, `FULL_TRANSITIVE`) check the new schema against all stored versions, not only the last one. Use transitive modes when consumers can be far behind.

Compatibility is not a substitute for tests. Test a new schema with an old consumer binary and with old files from production (anonymized).

A breaking change needs a new subject or a new topic and a migration plan. Do not force `NONE` to "unblock" a pipeline.

### Questions

#### Theoretical questions

1. What does backward compatibility allow?
2. What does forward compatibility allow?
3. What does full compatibility require?
4. What extra check does a transitive mode add?
5. Why is `NONE` risky on a shared topic?

#### Easy practical tasks

1. Write four sentences about compatibility. Use only facts from this section.
2. Make a table: Mode, old data + new reader, new data + old reader.
3. Find the default compatibility in Confluent or Karapace docs.
4. Give one schema change that is backward-safe and one that is not.

#### Medium practical tasks

1. Set a subject to BACKWARD. Register v1. Try to register a breaking v2. Record the error. Register a compatible v2.
2. Explain in six sentences why adding a required field without a default breaks backward compatibility.
3. Read official compatibility examples. Rewrite one example in STE.

#### Advanced practical tasks

1. Enable BACKWARD_TRANSITIVE. Show a change that passes versus the last version but fails versus v1. Write the schemas.
2. Write a one-page migration: incompatible change, new topic, dual publish window, consumer move.

---

## Subject naming strategies

A subject is the registry key for a sequence of versions. The naming strategy maps a Kafka record to a subject.

Common Confluent strategies:

- **TopicNameStrategy.** Subject is `{topic}-value` or `{topic}-key`. One schema family per topic for the value. This is the usual default.
- **RecordNameStrategy.** Subject is the fully qualified record name (Avro name or Protobuf message). Several topics can share one record schema. One topic can contain more than one record type.
- **TopicRecordNameStrategy.** Subject is `{topic}-{record name}`. The schema is specific to the topic and the record type.

Use TopicNameStrategy when one topic has one event type. That case is the simplest.

Use a record-based strategy when one topic is a stream of several types. Then consumers must handle more than one schema. This is harder.

Do not invent random subject names in each producer. Set the strategy in the client. Align all producers of a topic.

A wrong strategy creates a new subject and looks like "the schema is missing" on the consumer.

Keys often use a string id and no registry. If the key is Avro, the subject is `{topic}-key` under TopicNameStrategy.

### Questions

#### Theoretical questions

1. What is a subject?
2. What subject does TopicNameStrategy use for a value?
3. When do you use RecordNameStrategy?
4. What is TopicRecordNameStrategy?
5. Why must all producers of a topic use the same strategy?

#### Easy practical tasks

1. Write five sentences about subject names. Use only facts from this section.
2. Make a table: Strategy, example subject for topic `orders.placed` and record `Order`.
3. Find the default strategy in your client.
4. List two topics that fit TopicNameStrategy and one that might need record names.

#### Medium practical tasks

1. Register a schema with TopicNameStrategy. List subjects. Write the exact subject name.
2. Change only the strategy in a second producer (lab). Write the new subject and the consumer error if it occurs.
3. Read official subject name strategy documentation. Rewrite the three strategies in STE.

#### Advanced practical tasks

1. Design a topic that carries two event types. Choose a strategy. Write consumer dispatch rules.
2. Write a one-page standard: default TopicNameStrategy, who may approve RecordNameStrategy, and how you name topics so that subject names stay clear.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe the path from an object in a producer to a typed object in a consumer when a registry is in use.
2. Why can the broker say "success" when the consumer cannot deserialize?
3. How do compatibility mode and subject strategy work together?
4. When do you skip a registry and when must you use one?
5. How does this topic change the serializer section of topic 4?

#### Easy practical tasks

1. Write a one-page cheat sheet: wire bytes, registry, three formats, four compatibility ideas, three strategies.
2. Draw producer → registry HTTP → Kafka bytes → consumer → registry HTTP.
3. Name a subject, a schema id, and a topic as three different objects in one example.
4. Bookmark official Avro, Protobuf, and Schema Registry pages that you used.

#### Medium practical tasks

1. Build a small path: KRaft Kafka, a registry, one Avro or JSON Schema topic, produce, consume, list versions.
2. Fail a register on purpose (incompatible schema). Export the error text. Then register a compatible version.
3. Map each subsection to one official or Confluent URL.

#### Advanced practical tasks

1. Run two consumer versions (old schema, new schema) against a BACKWARD subject while you produce with the new schema. Write who succeeds.
2. Write a production schema standard: format, compatibility, strategy, CI check that registers schemas, and KRaft Kafka with no ZooKeeper.
