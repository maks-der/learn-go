# 3. Producers

## Description

A producer is a client that appends records to topic partitions. The producer selects a partition, batches records, sends them to the partition leader, and handles success or error.

This topic shows the producer API, acknowledgements, the idempotent producer, retries, delivery timeout, batching, compression, serialization, and poison messages. Complete this topic after topic 2. Complete this topic before topic 4.

Use one term for each concept. An acknowledgement (`acks`) is the broker reply that the producer waits for. The idempotent producer is a producer mode that prevents duplicates from retries. A poison message is a record that the producer or the consumer cannot process. Use KRaft clusters for all new work. The producer talks to brokers. The producer does not talk to ZooKeeper.

---

## Producer API

The producer API is the client interface that sends records. In Java the type is `KafkaProducer`. Other languages have a producer type with the same idea: configure, send, flush or close.

A typical sequence:

1. Build a configuration. Set `bootstrap.servers`. Set serializers for key and value.
2. Create the producer.
3. Build a record (topic, optional key, value, optional headers and partition).
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
- `acks=all` (or `acks=-1`). The leader waits until the in-sync replicas (ISR) meet `min.insync.replicas`. Topic 6 defines ISR. This setting is the durable choice for important data.

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

1. On a three-broker KRaft cluster, produce with `acks=all` while you stop a follower. Write whether produces succeed and what describe shows for ISR.
2. Write a durability matrix: rows `acks` 0, 1, and all. Add one column for "can lose on leader death". Fill each cell from facts in this topic and topic 6.

---

## Idempotent producer, retries, and delivery timeout

The **idempotent producer** prevents duplicates that retries would create on the same partition. The broker stores a producer id and a sequence number. A retry of the same batch does not append a second copy.

Current Java clients enable idempotence by default. Idempotence requires `acks=all` and a bounded `max.in.flight.requests.per.connection` (the client sets a safe value). Do not turn idempotence off without a written reason.

**Retries** happen when the producer gets a retriable error (network, leader change, not enough replicas). The producer sends the batch again. Without idempotence, a retry after a lost ack can create a duplicate in the log.

`retries` and `delivery.timeout.ms` bound how long the producer tries. `delivery.timeout.ms` is the maximum time from send to success or failure. `request.timeout.ms` is the time for one request. `linger.ms` is time spent in the batch before the first send. Those times must fit together. The official docs show the relation.

`retry.backoff.ms` waits between retries. A very small backoff can overload a recovering broker.

Idempotence does not make the full application exactly-once. A new `send` of the same business event with a new sequence is a new record. Use an event id in the payload or a transaction (topic 5) when the application retries a new send.

`transactional.id` is a further mode. Topic 5 covers transactions. You do not need transactions for this section.

### Questions

#### Theoretical questions

1. What duplicate problem does the idempotent producer prevent?
2. What configuration does idempotence require for `acks`?
3. How can a retry create a duplicate if idempotence is off?
4. What does `delivery.timeout.ms` limit?
5. Why does a new application `send` of the same event still create a new record?

#### Easy practical tasks

1. Write five sentences about the idempotent producer. Use only facts from this section.
2. Make a table: Key (`enable.idempotence`, `retries`, `delivery.timeout.ms`, `request.timeout.ms`). Add one meaning each.
3. Find the defaults for those keys in your client version.
4. Draw: first send, lost ack, retry, one record on the log (idempotent) versus two records (not idempotent).

#### Medium practical tasks

1. Write a producer with idempotence on. Print the configuration the client actually uses if the API allows it.
2. Read official notes on `delivery.timeout.ms` versus `linger.ms` and `request.timeout.ms`. Write the inequality in your own words.
3. Force a retriable error (stop a leader on a multi-broker lab, or use a short timeout). Write whether the callback reports success after retry.

#### Advanced practical tasks

1. Compare produce with idempotence on and off during retries if you can create a lost-ack case. Count records in the topic. Write the method and the limits of the test.
2. Write a one-page producer standard: idempotence on, `acks=all`, how you set delivery timeout, and when you add `transactional.id`.

---

## Batching and compression

The producer groups records into **batches** per partition. `batch.size` is the maximum batch size in bytes. `linger.ms` is the extra wait to add more records before send.

When the batch is full, the producer sends even if linger time is not over. When linger time is over, the producer sends even if the batch is not full.

A larger batch raises throughput. The first record in the batch waits. That wait raises latency. Topic 10 covers that trade-off again.

`buffer.memory` limits how much the producer can hold before send. If the buffer is full, `send` blocks or fails per configuration.

**Compression** runs on the batch. Common types: `none`, `lz4`, `snappy`, `gzip`, `zstd`. A larger batch often compresses better. Compression uses CPU on the client. The broker stores the compressed batch. Consumers decompress.

Do not set `batch.size` larger than what the broker accepts. Broker `message.max.bytes` and topic `max.message.bytes` limit the record batch.

For low latency, keep linger small. For high throughput, raise linger and batch size until the broker or the network is the limit. Measure after each change.

KRaft does not change these producer keys.

### Questions

#### Theoretical questions

1. What does `batch.size` limit?
2. What does `linger.ms` wait for?
3. When does the producer send before linger ends?
4. Why can compression work better on a large batch?
5. What happens when `buffer.memory` is full?

#### Easy practical tasks

1. Write four sentences about batching. Use only facts from this section.
2. Make a table: Key, unit, effect on throughput and latency. Add `batch.size` and `linger.ms`.
3. Find the default compression type in your client.
4. List two compression types and one reason to pick each type.

#### Medium practical tasks

1. Produce a stream with `linger.ms=0` and with `linger.ms=20`. Estimate records per second. Write both numbers and the method.
2. Enable `lz4` or `snappy`. Produce the same payload. Write whether produce time or CPU changed.
3. Read official batch and compression configuration. Write three facts in STE.

#### Advanced practical tasks

1. Measure produce p50 and p99 and records per second for two batch settings on KRaft. Write a short report with limits of the test.
2. Write a tuning note: start values for a low-latency topic and for a high-throughput topic. Include `max.message.bytes`.

---

## Serialization

A **serializer** writes the key or the value as bytes. A **deserializer** on the consumer reads those bytes. The pair must match. Topic 7 covers schemas and Schema Registry.

Kafka brokers do not parse fields. They store bytes. If you produce JSON text and consume with an Avro deserializer, the consumer fails or prints garbage.

Common client serializers:

- string (UTF-8)
- integer or long (fixed binary layout)
- byte array (you build the bytes)
- JSON (as text or through a library)
- Avro, Protobuf, or JSON Schema (often with a registry)

Set `key.serializer` and `value.serializer` (Java names). Other languages use the same idea.

Do not change the format of a live topic without a plan. Old records stay in the old format until retention removes them.

Headers can name a content type. The deserializer must still understand the bytes. A header is not a substitute for a matching serializer pair.

Null key and null value are valid. A null value is a tombstone on a compacted topic (topic 6). The serializer must allow null when you send a tombstone.

### Questions

#### Theoretical questions

1. What does a serializer create?
2. Does the broker validate JSON fields?
3. What happens when the serializer and the deserializer do not match?
4. Why is a content-type header not enough by itself?
5. When must the value serializer allow null?

#### Easy practical tasks

1. Write five sentences about serialization. Use only facts from this section.
2. Produce a JSON string. Consume as a string. Write the exact text.
3. Make a table: Format, human readable on consume, typical serializer name.
4. Find the serializer class names in your client for string and byte array.

#### Medium practical tasks

1. Write a producer and consumer with mismatched deserializers on purpose. Record the error. Fix the pair.
2. Produce UTF-8 text and dump the value as hex if you can. Write the first bytes.
3. Read official Serde or serializer documentation for your language. List five built-in serializers.

#### Advanced practical tasks

1. Write a one-page contract: topic name, key format, value format, and whether a Schema Registry is required (topic 7).
2. Produce a tombstone (null value) with a key. Consume. Write how your client prints the null value.

---

## Poison messages

A **poison message** is a record that a client cannot process. On the produce path, serialization can fail, the record can exceed `max.message.bytes`, or the broker can reject the produce. On the consume path, deserialization can fail, or the business logic can reject the payload.

If you ignore produce errors, you lose data and you may not see the poison. If a consumer stops on the first bad record, the partition does not move forward. Topic 4 covers the poll loop. Topic 8 covers a dead letter queue in Connect.

Handle produce failures in the callback or the future. Log the topic, partition, and error. Do not retry forever without a bound (`delivery.timeout.ms`).

For a record that is too large, split the payload, store the blob outside Kafka, or raise the size limit with a plan. Do not raise `max.message.bytes` on every topic without a disk and fetch plan.

For a payload that consumers cannot parse, fix the producer contract. A dead-letter topic is an application pattern: write the bad record to a second topic and commit the original offset. Topic 5 and topic 11 cover duplicates and poison handling in applications.

Do not use `acks=0` to "skip" poison errors. You only hide the failure.

### Questions

#### Theoretical questions

1. What is a poison message on the produce path?
2. What is a poison message on the consume path?
3. Why must you handle the send callback?
4. What are three responses to a record that is too large?
5. Why is `acks=0` the wrong fix for produce errors?

#### Easy practical tasks

1. Write four sentences about poison messages. Use only facts from this section.
2. Make a table: Failure, produce or consume, typical action.
3. Find `max.message.bytes` (or `message.max.bytes`) in the official docs. Write the default.
4. Draw: callback error → log → retry or dead-letter.

#### Medium practical tasks

1. Try to produce a value larger than the topic max if you can set a small `max.message.bytes`. Record the error.
2. Produce invalid bytes for your consumer deserializer. Write whether the consumer stops or skips. Do not leave the group broken; fix the consumer or the topic.
3. Read official error handling notes for your producer client. Write three error types that are retriable and two that are not.

#### Advanced practical tasks

1. Design a produce-side dead-letter topic: which errors go there, which headers you copy, and how you alert. Implement a small version.
2. Write a review checklist that rejects a producer that ignores send errors.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe a durable produce from `send` to ack. Name bootstrap, partitioner, batch, `acks=all`, and idempotence.
2. How do linger, batch size, and compression change both rate and delay?
3. Why is "the producer retried" not the same as "the application sent the event twice"?
4. What must stay the same between a producer serializer and a consumer deserializer?
5. A teammate sets `acks=0` and ignores callbacks to "go faster". Which facts do you use in a short reply?

#### Easy practical tasks

1. Write a producer that sends three keyed records with `acks=all`, then flushes and closes. Print partition and offset.
2. Write a one-page cheat sheet: API steps, three `acks` values, idempotence, delivery timeout, batch keys, compression, serializers, poison.
3. From your client defaults, write `acks`, idempotence, and linger.
4. Draw the path: record → serializer → batch → leader → ack.

#### Medium practical tasks

1. Change only `linger.ms` on a 1000-record produce. Record duration and whether callbacks all succeed.
2. Map each subsection to one official configuration key or API type.
3. Produce to a missing topic with auto-create off. Record the error. Create the topic. Retry.

#### Advanced practical tasks

1. Build a small producer that counts callback errors, retries at the application level with a reused event-id header, and stops after `delivery.timeout.ms` style budget.
2. Write a production producer standard: KRaft bootstrap only, `acks=all`, idempotence, serializer contract, and callback policy.
