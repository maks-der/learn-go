# 4. Producers

## Description

A producer is a client that appends records to topic partitions. The producer selects a partition, batches records, sends them to the partition leader, and handles success or error.

This topic covers the producer API, acknowledgements, the idempotent producer, retries, batching, compression, serialization, and poison messages. Complete this topic after topic 3. Complete this topic before topic 5.

Use one term for each concept. An acknowledgement (`acks`) is the broker reply that the producer waits for. The idempotent producer is a producer mode that prevents duplicates from retries. A poison message is a record that the producer or the consumer cannot process. Use KRaft clusters for all new work. The producer talks to brokers. The producer does not talk to ZooKeeper.

---

## Producer API

The producer API is the client interface that sends records. In Java the type is `KafkaProducer`. Other languages have a producer type with the same idea: configure, send, flush or close.

A typical sequence:

1. Build a configuration. Set `bootstrap.servers`. Set serializers for key and value.
2. Create the producer.
3. Build a record (`topic`, optional `key`, `value`, optional headers and partition).
4. Call `send`. `send` is often asynchronous. The client returns a future or calls a callback.
5. Call `flush` when you must wait for in-flight records. Call `close` when the process ends.

`bootstrap.servers` is a list of brokers. The producer fetches metadata. The producer then sends each batch to the current leader of the target partition.

Do not create a new producer for every record. A producer holds connections and buffers. Share one producer per process when the configuration is the same.

Handle the result of `send`. Ignore of the callback is a common cause of silent data loss.

The console producer is a small producer for tests. Use it to learn. Use the client API for applications.

### Questions

#### Theoretical questions

1. What is the producer API?
2. What configuration key names the brokers?
3. Why is `send` often asynchronous?
4. Why must you close or flush a producer?
5. Why must you not create one producer per record?

#### Easy practical tasks

1. Write five sentences that describe the producer API. Use only facts from this section.
2. Send one line with the console producer. Write the command.
3. Make a table: "API step" and "What it does". Add five rows.
4. Find the producer class or function name in the client for your language. Write the package or module name.

#### Medium practical tasks

1. Write a small program that produces three records and prints the partition and offset from the send result.
2. Produce without a `close` or `flush` in a short program that exits immediately. Write whether all records arrive. Then add `close` and compare.
3. Read the official producer API page for your client. List five configuration keys that you must understand.

#### Advanced practical tasks

1. Write a producer that sends 1000 records and waits for all callbacks. Count successes and errors. Print the counts.
2. Compare the Java `KafkaProducer.send` callback with the async API in your language. Write a one-page map of the same steps.

---

## Acks: `0`, `1`, `all`

`acks` controls when the producer treats a send as successful.

- `acks=0`. The producer does not wait for a broker reply. The broker can drop the record. The producer can also lose the record on a network error. Throughput can be high. Durability is low.
- `acks=1`. The leader writes the record to its local log and replies. Followers may not have the record yet. If the leader fails before followers copy the record, the record can be lost.
- `acks=all` (or `acks=-1`). The leader waits until the in-sync replicas (ISR) meet `min.insync.replicas`. Topic 9 defines ISR. This setting is the durable choice for important data.

The idempotent producer requires `acks=all`. Current Java clients enable idempotence by default and set `acks=all`.

Do not use `acks=0` for orders, payments, or other data that you cannot lose. Use `acks=0` only when loss is acceptable (some metrics).

`acks` is a producer setting. `min.insync.replicas` is a broker or topic setting. Both matter for `acks=all`.

### Questions

#### Theoretical questions

1. What does `acks=0` mean?
2. What does `acks=1` mean?
3. What does `acks=all` mean?
4. Why can `acks=1` still lose a record?
5. How does `min.insync.replicas` relate to `acks=all`?

#### Easy practical tasks

1. Make a three-row table: `acks` value, what the producer waits for, loss risk.
2. Write four sentences that compare `acks=1` and `acks=all`.
3. Find the default `acks` in your client documentation. Write the value.
4. List three event types that need `acks=all` and two that can use `acks=0`.

#### Medium practical tasks

1. Produce with `acks=all` and print success. Produce with `acks=0` if the client allows it. Write how you detect success in each case.
2. Read the official configuration for `acks`. Rewrite the three values in STE.
3. On a topic with replication factor 1, explain in five sentences what `acks=all` can still wait for.

#### Advanced practical tasks

1. On a multi-broker cluster, set `min.insync.replicas=2` and `acks=all`. Stop one replica. Produce. Record the error or the success. Restore the replica.
2. Write a one-page durability note for a payments topic: `acks`, replication factor, and `min.insync.replicas`. Do not use ZooKeeper in the note.

---

## Idempotent producer

The idempotent producer prevents duplicates that retries create. Without idempotence, a producer can send a record, lose the acknowledgement, retry, and append the same record twice.

With idempotence, the producer has a producer id (PID) and a sequence number per partition. The broker rejects a duplicate sequence. The log keeps one copy.

Enable the feature with `enable.idempotence=true`. Current Java clients set this default. Idempotence requires `acks=all` and a bounded set of in-flight requests per connection.

Idempotence does not make a whole business operation exactly-once. It removes producer retry duplicates. Topic 7 covers exactly-once semantics and transactions.

Do not disable idempotence to "make it simpler". The simple durable default is idempotence on.

The idempotent producer is not the same as a consumer that ignores duplicate business keys. Both can exist in one system.

### Questions

#### Theoretical questions

1. What duplicate problem does the idempotent producer solve?
2. What are a producer id and a sequence number for?
3. What `acks` value does idempotence require?
4. What problem does idempotence not solve?
5. Why is `enable.idempotence=true` the default in current Java clients?

#### Easy practical tasks

1. Write five sentences about the idempotent producer. Use only facts from this section.
2. Find `enable.idempotence` in your client docs. Write the default.
3. Make a table: "Without idempotence" and "With idempotence" for a retry after a lost ack.
4. Draw a sequence: send, lost ack, retry. Mark the broker check of the sequence number.

#### Medium practical tasks

1. Write a producer with idempotence enabled. Print the configuration the client actually uses if the API allows it.
2. Read the official idempotent producer documentation. Write the in-flight request limit that the feature needs.
3. Explain to a teammate in six sentences why two application sends with the same value can still create two records.

#### Advanced practical tasks

1. Design a test that would create a retry duplicate if idempotence were off (describe the fault). Do not disable safety on a shared cluster. Write the expected log contents with the feature on.
2. Compare idempotent producer and transactional producer in a one-page table. Topic 7 will go deeper. Use official names only.

---

## Retries and `delivery.timeout.ms`

Networks fail. Leaders move. The producer retries a failed send when the error is transient. `retries` (or an equivalent) sets how many times the client retries. Modern clients use a large retry count and bound the time instead.

`delivery.timeout.ms` is the total time budget for a record. The budget includes linger time, send time, and retries. When the budget ends, the producer fails the send. You handle that error.

`request.timeout.ms` is the timeout of one request. `linger.ms` is a wait before the batch leaves. The delivery timeout must be larger than linger plus request timeout.

Retry of a non-idempotent producer can create duplicates. Retry of an idempotent producer does not create those duplicates.

Do not retry forever in application code around `send` without a bound. Use the client timeout. Log the failed key and value (or a safe summary) so that you can repair.

Some errors are not transient. Examples: record too large, unknown topic (when auto-create is off), serialization error. Retries do not fix those errors.

### Questions

#### Theoretical questions

1. Why does a producer retry?
2. What does `delivery.timeout.ms` bound?
3. How do `linger.ms` and `request.timeout.ms` relate to the delivery timeout?
4. Which errors must you not expect retries to fix?
5. How does idempotence change the meaning of a retry?

#### Easy practical tasks

1. Make a table: Configuration key, meaning. Add `retries`, `delivery.timeout.ms`, `request.timeout.ms`, `linger.ms`.
2. Write four sentences about a time budget for one record.
3. Find the defaults for those keys in your client.
4. List three transient errors and three permanent errors (from docs or this section).

#### Medium practical tasks

1. Set a very small `delivery.timeout.ms` and a large `linger.ms` on a local producer. Produce. Record the error. Fix the values.
2. Read the official configuration page for `delivery.timeout.ms`. Rewrite the relationship to other timeouts in STE.
3. Write application error handling in a small producer: log and increment a counter on a failed send. Do not retry without a bound.

#### Advanced practical tasks

1. Draw a time line for one record: linger, first request, retry, success or expire. Put the delivery timeout as a bar.
2. Write a one-page runbook: what an on-call person does when produce error rates rise (check topics, leaders, timeouts, record size).

---

## Batching: `linger.ms`, `batch.size`

The producer groups records that go to the same partition into a batch. A batch is one produce request payload for that partition.

`batch.size` is the maximum batch size in bytes. The producer sends earlier when the batch is full.

`linger.ms` is the extra wait for more records. `linger.ms=0` sends as soon as possible. A small linger (for example 5–20 ms) often increases throughput because each request carries more records.

Batching increases latency by up to the linger time. Batching decreases the number of requests. Choose linger from the latency budget of the application.

The sticky partitioner for null keys exists to fill batches. Topic 3 covers that behavior.

`buffer.memory` is the total memory the producer can use for records that are not yet sent. When the buffer is full, `send` blocks or throws, depending on `max.block.ms`.

### Questions

#### Theoretical questions

1. What is a batch in the producer?
2. What does `batch.size` limit?
3. What does `linger.ms` wait for?
4. How does linger change latency and throughput?
5. What happens when `buffer.memory` is full?

#### Easy practical tasks

1. Write five sentences about batching. Use only facts from this section.
2. Make a table: `linger.ms=0` versus `linger.ms=10` for latency and request count (qualitative).
3. Find defaults for `linger.ms` and `batch.size` in your client.
4. Draw one partition batch that contains three records.

#### Medium practical tasks

1. Produce 10 000 small records with `linger.ms=0` and then with `linger.ms=10`. Compare elapsed time. Write both times.
2. Set a very small `batch.size`. Produce. Compare request behavior or elapsed time with a larger batch size.
3. Read official docs for `buffer.memory` and `max.block.ms`. Write what the producer does when the buffer is full.

#### Advanced practical tasks

1. Build a small bench: vary linger and batch size. Write a table of throughput. State the record size.
2. Explain in one page why a sticky partitioner plus linger fills batches better than round-robin null keys.

---

## Compression

The producer can compress a batch. Common codecs are `gzip`, `snappy`, `lz4`, and `zstd`. You set `compression.type` on the producer (or on the topic).

Compression reduces network and disk size. Compression uses CPU. `lz4` and `zstd` are common defaults for new systems. `gzip` is slower. `snappy` is common in older setups.

The broker stores the compressed batch. Consumers decompress. All clients that read the topic must understand the codec. Current Kafka clients support the common codecs.

Do not compress twice without a reason. If the value is already a compressed blob, producer compression may not help.

Measure. A JSON batch often compresses well. A batch of random bytes does not.

### Questions

#### Theoretical questions

1. What does producer compression compress?
2. Name four codecs that Kafka uses.
3. What resource does compression save? What resource does it use?
4. Who decompresses the batch?
5. Why might compression not help a random binary value?

#### Easy practical tasks

1. Make a table: Codec, typical note (speed or ratio) from docs or this section. Add four rows.
2. Find `compression.type` in your client. Write the allowed values.
3. Write four sentences about compression. Use only facts from this section.
4. List two payloads that compress well and two that do not.

#### Medium practical tasks

1. Produce the same 10 000 JSON records with `compression.type=none` and with `lz4` or `zstd`. Compare produce time and, if you can, topic size on disk.
2. Set compression on the topic instead of the producer if your version supports it. Write which setting won in a describe.
3. Read official codec notes. Write one sentence per codec.

#### Advanced practical tasks

1. Write a one-page choice: `lz4` versus `zstd` for a 50 MB/s topic. Use public benchmarks or your own measure.
2. Inspect a log segment (or a metric) to confirm compressed batches. Document the method.

---

## Serialization (JSON, Avro, Protobuf)

The producer must convert the key and the value to bytes. A serializer does that work. The consumer uses a matching deserializer.

JSON is text bytes. JSON is easy to debug. JSON has no built-in schema check on the broker. Producers can write different shapes. Consumers break.

Avro and Protobuf are binary encodings with a schema. Teams often store the schema in a Schema Registry. Topic 11 covers the registry and compatibility.

Set `key.serializer` and `value.serializer` (Java names) or the equivalent. A common learning pair is string or JSON for both key and value. Do not use a JSON serializer for a key that you want to hash as a stable id unless you control the JSON field order. A string id is a simpler key.

Kafka does not validate your JSON. A produce success means the bytes are in the log, not that the value is a valid business event.

### Questions

#### Theoretical questions

1. What does a serializer do?
2. Why is JSON easy and also risky?
3. What do Avro and Protobuf add that JSON does not?
4. Why is a raw string often a better key than JSON?
5. Does a successful produce prove that the value is a valid event?

#### Easy practical tasks

1. Produce a JSON line with the console producer. Consume it. Write the exact bytes you see as text.
2. Make a table: Format, human readable, schema on the wire. Add JSON, Avro, Protobuf.
3. Find the serializer class names in your client for string and JSON.
4. Write four sentences about serialization. Use only facts from this section.

#### Medium practical tasks

1. Write a producer that sends a small JSON object. Write a consumer that parses it. Then send a JSON object with a missing field. Record the consumer error.
2. Send a key as a string id and a value as JSON. Print the partition. Confirm the key is the id, not the JSON.
3. Read topic 11 headings in `kafka.topics.md`. Write three questions you will answer later about Avro.

#### Advanced practical tasks

1. Produce the same logical event as JSON and as Protobuf (or Avro). Compare sizes. Write the two sizes.
2. Write a one-page rule: when the team may use JSON in Kafka and when the team must use a schema format.

---

## Error handling and poison messages

A poison message is a record that blocks progress. On the produce side, a poison record can be a value that does not serialize, a record that is larger than `max.request.size` or the topic `max.message.bytes`, or a record that the broker rejects.

On the consume side (topic 5), a poison record can be bytes that do not deserialize or a value that fails a business rule. This section focuses on the producer.

Handle serializer errors before `send`. Do not retry a record that cannot serialize. Fix the application or route the event to a repair path.

Handle broker errors in the callback. Transient errors belong to the client retry budget. Permanent errors belong to metrics, logs, and a dead-letter plan.

A dead-letter topic is a topic that stores records that failed. Write the original topic, the error, and enough payload to repair. Do not create an infinite loop between the main topic and the dead-letter topic.

Bound record size. Large records hurt brokers and consumers. Store large blobs in object storage. Put a pointer in the Kafka value.

### Questions

#### Theoretical questions

1. What is a poison message on the produce path?
2. Why must you not retry a serializer error?
3. What is a dead-letter topic?
4. Why are large records a problem?
5. What belongs in a produce callback error handler?

#### Easy practical tasks

1. Write five sentences about produce errors. Use only facts from this section.
2. Make a table: Error type, retry or not. Add serialization, timeout, record too large.
3. Find `max.request.size` or `max.message.bytes` in the docs. Write the meaning.
4. List three fields that a dead-letter record must contain.

#### Medium practical tasks

1. Try to produce a record that is larger than the topic limit (raise a small `max.message.bytes` on a test topic). Record the error.
2. Write a producer callback that logs a permanent error and writes one line to a local file (a stand-in for a dead-letter topic).
3. Read official error codes or exception types for the producer. Classify five errors as transient or permanent.

#### Advanced practical tasks

1. Design a dead-letter topic for a payments producer. Include headers, key, and how a human replays a fixed record.
2. Write a one-page policy: max value size, pointer-to-S3 pattern, and who owns repair of poison records.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe the path of one record from `send` to a durable offset when `acks=all` and idempotence is on.
2. How do batching, compression, and linger work together in one produce request?
3. Why is a console producer not enough to learn `acks` and callbacks?
4. What configuration pair protects you from retry duplicates and from lost acks?
5. How does serialization failure differ from a broker timeout in the repair path?

#### Easy practical tasks

1. Write a producer that sends two JSON records with keys. Print topic, partition, offset. Use `acks=all`.
2. Write a one-page cheat sheet: `acks`, idempotence, `delivery.timeout.ms`, `linger.ms`, `batch.size`, `compression.type`, serializers.
3. From your client config object, highlight bootstrap, acks, and serializers.
4. Draw a producer with a buffer, a batch per partition, and connections to two leaders.

#### Medium practical tasks

1. Write a script or program that produces 1000 records, then a second run that uses linger and compression. Write throughput for both runs.
2. Create a topic with a small `max.message.bytes`. Produce a valid small record and a too-large record. Save both outcomes.
3. Map each subsection of this topic to one official configuration key or API type. Write the list.

#### Advanced practical tasks

1. Build a producer with idempotence, `acks=all`, `lz4`, and a 10 ms linger. Fail the broker during a send (stop the container). Record retry behavior and whether duplicates appear.
2. Write a production checklist (15 items) for a new producer. Include KRaft bootstrap, not ZooKeeper.
