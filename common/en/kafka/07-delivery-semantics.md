# 7. Delivery Semantics

## Description

Delivery semantics describe what the system guarantees when a failure occurs. Kafka and your code together create the guarantee. The broker alone does not make a business action exactly-once.

This topic covers at-most-once, at-least-once, exactly-once (EOS), idempotent consumers, and deduplication keys. Complete this topic after topics 4, 5, and 6.

Use one term for each concept. At-most-once means a record can be lost and is not processed twice. At-least-once means a record is not lost and can be processed twice. Exactly-once means a processed record has one visible effect in the target system. EOS in Kafka is the transactional protocol. An idempotent consumer is application logic that can see a record twice and still produce one effect. Prefer KRaft clusters. Transactions run on Kafka brokers, not on ZooKeeper.

---

## At-most-once

At-most-once delivery means each record is processed zero times or one time. A record is never processed twice. A record can be lost.

A typical Kafka path that gives at-most-once:

1. The consumer commits the offset.
2. The consumer then processes the record.
3. The process stops after the commit and before the process step.

The restart skips the record. Auto-commit can create this path when the client commits before your work finishes (topic 5).

On the produce side, `acks=0` can lose records. That is also a form of loss. Combined with a consumer that never retries, you get at-most-once end-to-end.

Use at-most-once only when loss is acceptable. Examples: some metrics, some debug traces. Do not use at-most-once for payments or stock.

At-most-once is not a Kafka switch with that name. It is the result of commit-before-process and of produce without durability.

### Questions

#### Theoretical questions

1. What does at-most-once mean?
2. How can a consumer skip a record?
3. How can a producer lose a record?
4. When is at-most-once acceptable?
5. Is at-most-once a single Kafka configuration key?

#### Easy practical tasks

1. Write five sentences about at-most-once. Use only facts from this section.
2. Draw the commit-then-process sequence. Mark the crash point that skips a record.
3. Make a table: Setting or order, how it loses data. Add two rows.
4. List three event types that may use at-most-once and three that must not.

#### Medium practical tasks

1. Write a consumer that commits and then sleeps. Kill it during the sleep. Restart. Write which records never appear again.
2. Produce with `acks=0` if the client allows it. Stop the broker during a burst. Write whether all records exist after restart.
3. Read any official note that describes at-most-once. Rewrite it in STE.

#### Advanced practical tasks

1. Write a one-page policy that forbids at-most-once for a named list of topics. Include how you detect commit-before-process in review.
2. Compare at-most-once in Kafka with at-most-once in a typical queue acknowledgement. Write five differences.

---

## At-least-once

At-least-once delivery means each record is processed one time or more. A record is not lost (if produce is durable and retention still holds the record). A record can be processed twice.

A typical Kafka path that gives at-least-once:

1. The consumer processes the record (for example a database write).
2. The consumer then commits the offset.
3. The process stops after the write and before the commit.

The restart processes the record again. A producer retry without idempotence can also append a duplicate (topic 4).

At-least-once is the common default for careful applications: `acks=all`, idempotent producer, manual commit after work.

Duplicates are not a Kafka bug in this mode. Your consumer or your database must tolerate them, or you must add EOS or idempotent writes.

Lag and retries increase the chance of duplicates. A rebalance in the middle of a batch can also repeat records that were processed but not committed.

### Questions

#### Theoretical questions

1. What does at-least-once mean?
2. How does process-then-commit create a duplicate?
3. How can the producer side create a duplicate?
4. Why is at-least-once the common careful default?
5. Why must the consumer tolerate duplicates in this mode?

#### Easy practical tasks

1. Write four sentences about at-least-once. Use only facts from this section.
2. Draw process-then-commit. Mark the crash point that duplicates a record.
3. Make a table: At-most-once vs at-least-once. Add rows for loss and for duplicates.
4. List three downstream writes that become wrong if they see a duplicate (for example "add 10 to balance").

#### Medium practical tasks

1. Write a consumer that writes a line to a file, then commits. Kill it after the write and before the commit. Restart. Count duplicate lines.
2. Enable the idempotent producer. Retry a send. Write why the log should not contain a producer-retry duplicate.
3. Find official text about at-least-once. Rewrite it in STE.

#### Advanced practical tasks

1. Measure duplicate rate: kill the consumer 20 times during a 1000-record load. Count extra file lines. Write the rate.
2. Write a one-page note: which parts of the stack you set so that loss is unlikely and only duplicates remain.

---

## Exactly-once (EOS) — transactions, `isolation.level`, idempotent producer

Exactly-once semantics (EOS) in Kafka means you can write records and commit offsets in one transaction. A downstream consumer with `isolation.level=read_committed` does not see records from a transaction that did not commit.

Core pieces:

- **Idempotent producer.** Required. Removes retry duplicates (topic 4).
- **Transactions.** Set `transactional.id` to a stable unique id per producer instance. Call init, begin, send, send offsets to the transaction, commit or abort.
- **`isolation.level=read_committed`.** The consumer hides aborted and in-flight transactional records. The default `read_uncommitted` can show them.

EOS is not magic for every side effect. A transaction covers Kafka writes and the committed offsets that you send in that transaction. A write to a database outside the transaction is not covered unless you use a pattern such as outbox (topic 19) or an idempotent database write.

`transactional.id` must not be shared by two live processes. A fencing token fences an old producer after a restart.

Use EOS when you consume-transform-produce on Kafka and you cannot accept duplicate output records. Kafka Streams can enable EOS for you. Topic 14 covers that path.

Transactions use internal topics. The cluster must have a correct KRaft setup and enough replication for transaction state. Follow official production recommendations.

### Questions

#### Theoretical questions

1. What does EOS mean in Kafka?
2. What does `transactional.id` identify?
3. What does `isolation.level=read_committed` hide?
4. What side effect does a Kafka transaction not cover by itself?
5. What is fencing at a high level?

#### Easy practical tasks

1. Write five sentences about Kafka EOS. Use only facts from this section.
2. Make a table: Feature, role. Add idempotent producer, transaction, `read_committed`.
3. Find `isolation.level` in your client. Write the default.
4. Draw begin → send → send offsets → commit.

#### Medium practical tasks

1. Write a transactional consume-transform-produce if your client supports it. Abort once. Confirm that a `read_committed` consumer does not see the aborted records.
2. Read the same output topic with `read_uncommitted` if you can. Write the difference during an open transaction.
3. Read official transactional producer documentation. Write the required configuration keys.

#### Advanced practical tasks

1. Crash after `send` and before `commit`. Restart with the same `transactional.id`. Write whether the output contains duplicates for a `read_committed` consumer.
2. Write a one-page design: EOS inside Kafka versus idempotent writes to a database. Choose one for a payments projector.

---

## Idempotent consumers (the practical exact-once)

An idempotent consumer is a consumer whose process step can run twice and still leave one correct effect.

This is the practical exact-once for many systems. You use at-least-once Kafka (duplicates possible) and a sink that ignores a repeat.

Examples:

- store the event id in a processed-events table. Skip if the id exists.
- `UPSERT` a row by primary key to the latest state.
- set a value instead of increment a value, when the latest value is the truth.

The consumer still uses manual commit after the sink confirms the write. If you commit before the sink write, you can skip. If you write and then crash before commit, the second run hits the idempotent sink and does nothing extra.

Idempotent consumers work with systems that Kafka transactions cannot include. They are often simpler than EOS when the sink is a database.

You must define the idempotency key. The next section covers keys. A bad key (only the timestamp) fails.

Idempotent consumers do not remove the need for `acks=all` on important produces. Loss is still a separate problem.

### Questions

#### Theoretical questions

1. What is an idempotent consumer?
2. Why is this the practical exact-once for many teams?
3. Why must you commit after the sink write?
4. Why can `UPSERT` by key be idempotent when "add 1" is not?
5. Does an idempotent consumer replace durable produce?

#### Easy practical tasks

1. Write five sentences about idempotent consumers. Use only facts from this section.
2. Make a table: Sink operation, idempotent or not. Add increment, set, insert-if-absent.
3. Draw write-sink-then-commit with a second run that hits a duplicate key.
4. List three sinks in a shop that can be idempotent.

#### Medium practical tasks

1. Write a consumer that writes event ids to a local SQLite or a file set. Skip ids that exist. Kill and restart. Count sink rows versus input records.
2. Implement an increment counter and a set-latest counter. Send the same event twice. Write the two results.
3. Find a public article or official pattern on idempotent consumers. Write three facts. Do not copy the text.

#### Advanced practical tasks

1. Add a unique constraint on (consumer_name, event_id) in a real database. Prove that a duplicate poll cannot create two effects.
2. Write a one-page comparison: Kafka EOS versus idempotent consumer for a job that writes Postgres. Choose one and defend it.

---

## Deduplication keys

A deduplication key is the value that identifies one logical event. The consumer or the sink stores that key.

Good keys:

- a unique event id that the producer sets (UUID or ULID)
- a natural key plus a version (`order-55` + `seq=3`)
- the Kafka topic, partition, and offset, when you only need "this log record once"

Topic-partition-offset is unique in one cluster for that record. It does not stay stable if you copy records to a new topic or a new cluster without mapping. An event id in the payload survives a copy.

Bad keys:

- only the current time
- only the user id when the user can emit many events
- a JSON blob that changes field order

Put the event id in the value or in a header. The producer must set it once and must not change it on retry. The idempotent producer can retry the same bytes. Application retry of a new send must reuse the same event id.

Retention of the dedup store must be longer than the maximum retry window. If you forget a key after one day, a late duplicate can apply again.

### Questions

#### Theoretical questions

1. What is a deduplication key?
2. When is topic-partition-offset a valid key?
3. When is topic-partition-offset a poor key?
4. Why must a producer reuse the event id on an application retry?
5. Why must the dedup store retain keys long enough?

#### Easy practical tasks

1. Write four sentences about deduplication keys. Use only facts from this section.
2. Make a table: Key choice, good or bad, reason. Add five rows.
3. Design a header name for an event id. Write the format.
4. List two systems where you would store the keys (database table, cache with TTL).

#### Medium practical tasks

1. Produce records with an `event-id` header. Write a consumer that skips duplicate ids. Send the same id twice on purpose.
2. Use topic-partition-offset as the key in a file. Explain in five sentences when this breaks after a mirror.
3. Read about MirrorMaker or cluster linking at a high level (topic 18). Write one sentence about offset translation and keys.

#### Advanced practical tasks

1. Design a dedup table: columns, unique constraint, TTL or purge job, and the retry window. Implement a small version.
2. Write a one-page standard for producers: who creates the event id, which header, and what happens if the header is missing.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Put at-most-once, at-least-once, Kafka EOS, and idempotent consumers on one scale from "can lose" to "one effect". Use one sentence each.
2. Why is "exactly-once" a property of the full path, not of `poll` alone?
3. How do `acks=all`, idempotent produce, and process-then-commit combine?
4. When do you choose Kafka transactions instead of a processed-events table?
5. How does `read_committed` change what a consumer can see during an open transaction?

#### Easy practical tasks

1. Write a one-page cheat sheet with a sequence diagram for each of the four approaches in this topic.
2. Label three of your learning topics (`demo` events) with the semantics you would use in production.
3. Make a table: Failure (kill consumer, kill broker, retry produce), risk (loss, duplicate, none) for at-least-once with idempotent produce.
4. Draw fencing of an old transactional producer in three boxes.

#### Medium practical tasks

1. Implement two consumers on the same input: one commit-first, one process-first. Crash both. Write skip versus duplicate evidence.
2. If your client supports transactions, write one transactional produce and one non-transactional produce to the same topic. Consume with both isolation levels.
3. Map each subsection to official documentation (delivery semantics, transactions, isolation).

#### Advanced practical tasks

1. Build a small pipeline: consume A, produce B, with either EOS or an event-id table. Prove with a crash test that B does not double-apply a payment-like effect.
2. Write a production decision tree (one page) that a team follows when they pick semantics. Include KRaft and exclude ZooKeeper.
