# 11. Application Patterns

## Description

Application patterns are ways to use Kafka in a system design. The broker does not choose a pattern. You choose what a record means and how a service writes and reads it.

This topic shows event notification versus event-carried state versus event sourcing, outbox, CQRS, ordering, duplicates, poison messages, and the poll-loop rule. Complete this topic after topics 3, 4, 5, and 7. Use a KRaft cluster. These patterns do not need ZooKeeper.

Use one term for each concept. An event notification is a small signal that something changed. Event-carried state transfer is an event that carries the data that a reader needs. Event sourcing stores every change as an event and builds state from that log. An outbox is a table (or equivalent) in the same database transaction as the business write. CQRS splits write models and read models. Backpressure is how you slow input when the consumer cannot keep up.

---

## Event notification vs event-carried state vs event sourcing

**Event notification** says that a fact occurred and gives a pointer. Example: `{ "orderId": "55" }`. The consumer must call an API or read a database to get the order. Many consumers create load on the source. If the source is down, the event is not enough.

**Event-carried state transfer (ECST)** puts the needed fields in the event. Example: order id, items, totals, customer id. Consumers can update their own store without a call back. The event is larger. Schema evolution matters (topic 7).

**Event sourcing** uses the event log as the source of truth for a stream of changes. Current state is a fold of events. You do not overwrite a single "current row" as the only truth. Kafka can store that log. You still need a way to snapshot and to define the event types. Event sourcing is a design, not a Kafka switch.

Do not mix the three without a name. A topic that sometimes has a pointer and sometimes a full document will break consumers.

Notification is fine when the payload is huge or private and readers are few. ECST is the common integration pattern. Event sourcing is for a bounded context that accepts replay as the model.

Retention must match the pattern. Event sourcing needs a long or compacted-with-care history. Notification topics can use short time retention if the pointer is only a hint.

KRaft does not change the meaning of an event. It only hosts the log.

### Questions

#### Theoretical questions

1. What does an event notification contain?
2. What does event-carried state transfer contain?
3. What is the source of truth in event sourcing?
4. Why can a notification flood the source API?
5. Why must a topic pick one pattern and keep it?

#### Easy practical tasks

1. Write one shop event three ways: notification, ECST, event-sourced command or fact.
2. Make a table: Pattern, payload size, consumer needs API (yes or no).
3. Write four sentences about when you choose ECST.
4. Draw notification (event + HTTP get) versus ECST (event only).

#### Medium practical tasks

1. Take an existing API you know. Design one ECST topic and one notification topic. Write who consumes each.
2. Write six sentences on retention for event sourcing versus a notification topic.
3. Find a public article on ECST (for example a known integration book idea). Rewrite three facts in STE. Do not copy a long passage.

#### Advanced practical tasks

1. Design five events for payments. Mark each as notification, ECST, or sourcing. Write the schema fields at a high level.
2. Write a one-page standard: default ECST, when notification is allowed, when event sourcing is a team decision.

---

## Outbox and CQRS

The **outbox pattern** solves this failure: a service writes a database row and then fails before it produces to Kafka (or the reverse). The two stores diverge.

The service writes the business row and an **outbox row** in the **same database transaction**. A separate publisher reads the outbox and produces to Kafka. After a successful produce (and the durability that you require, topic 3 and topic 5), the publisher marks the outbox row as sent or deletes it.

If the process stops after the commit and before the produce, the publisher retries. Kafka can see a duplicate if the produce succeeded and the mark did not. Use an idempotent producer and a deduplication key (topic 5).

Debezium can capture the outbox table as CDC (topic 8). Then the application does not run a poll loop on the table. The transaction log becomes the publisher.

Do not produce first and write the database second without a compensating plan. Do not use a two-phase commit across Kafka and a SQL database as the default.

The Kafka cluster is KRaft. The outbox lives in your database. Connect or a worker reads it.

**CQRS** (command query responsibility segregation) splits the write model from the read model. Commands change the write side. Queries read a model that is shaped for the query.

Kafka is often the pipe between them. The write service commits (database and/or events). Events go to Kafka (outbox or event sourcing). A consumer builds a read store: search index, cache, or SQL view. The read API does not query the write database.

The read model is **eventually consistent**. A query can miss a command that is still in the log or in the consumer. Tell the product team. Do not hide this.

You can have more than one read model from one topic. That is a reason to use Kafka instead of one shared database for all queries.

Do not put every query on Kafka consume-from-beginning. Use a stored read model. Kafka retention is not your query API.

**Request-reply over Kafka** is usually a poor default. HTTP and RPC fit a call that waits for one response. If you still use two topics and a correlation id, set short timeouts and do not block the poll thread (next section).

### Questions

#### Theoretical questions

1. What dual-write failure does the outbox prevent?
2. What two writes share one database transaction?
3. What does CQRS split?
4. Why is the read model eventually consistent?
5. Why is "query Kafka from the beginning on each HTTP request" a poor default?

#### Easy practical tasks

1. Write five sentences about the outbox. Use only facts from this section.
2. Draw: HTTP handler → DB transaction (order + outbox) → publisher → Kafka → consumer → read DB → query API.
3. Make a table: Step, can fail, result if no outbox, result with outbox.
4. Give one example where the UI must tolerate a short delay.

#### Medium practical tasks

1. Implement a small outbox table and a publisher in a language that you know. Use KRaft Kafka. Stop after commit, before produce. Confirm a retry.
2. Build a tiny CQRS lab: produce an ECST event, consume into a map or SQLite, query the map.
3. Compare outbox+CDC with a process that writes only to Kafka and builds a read model. Write five trade-offs.

#### Advanced practical tasks

1. Add a second read model (for example, counts per day) from the same topic. Show both views after three events.
2. Write an outbox and CQRS standard: required columns, who publishes, how you measure lag, KRaft bootstrap, no dual-write without outbox.

---

## Ordering, duplicates, and poison messages

**Ordering** in Kafka is per partition (topic 2). Records in one partition have a total order. Records in two partitions do not.

If you need a happens-before relation between two events, put them in the same partition (same key) or put a causal token in the payload (for example, a version per entity).

Clock timestamps are not enough. Two producers can have skewed clocks. Do not sort global events by timestamp and call that causality.

Consumer groups do not change partition order. Rebalances can pause a partition. They do not reorder the log. Unclean leader election can drop records (topic 6). Keep unclean election off.

**Duplicates** are normal in at-least-once designs (topic 5). Causes: producer retry without idempotence, process-then-crash-before-commit, application retry of a new send, rebalance in the middle of a batch.

Handle duplicates with an idempotent consumer, a deduplication key, or Kafka transactions for Kafka-to-Kafka paths. Do not treat a duplicate as a Kafka defect when you chose at-least-once.

**Poison messages** are records that your application cannot process (topic 3). A poison record can stop the partition if you fail the whole poll loop. Patterns:

- write the record to a dead-letter topic, then commit the original offset
- skip after a bound number of retries and alert
- fix the producer contract so that the poison stops

A dead-letter topic is a normal Kafka topic. Set retention. Consume it with an operator process. Connect has a DLQ feature (topic 8). Applications can use the same idea.

Do not use `acks=0` to hide produce failures. Do not drop poison records without a log.

Request-reply over Kafka adds more duplicates on timeout and retry. Prefer HTTP for user-facing calls.

### Questions

#### Theoretical questions

1. Where does Kafka guarantee order?
2. How do you keep two events in order for one entity?
3. Name three causes of duplicates in an at-least-once path.
4. What happens if a consumer stops on the first poison record?
5. Why are clocks a poor causality tool?

#### Easy practical tasks

1. Write five sentences about order, duplicates, and poison. Use only facts from this section.
2. Draw two partitions. Place events for user 9 and user 8. Mark what is ordered.
3. Make a table: Problem, tool (key, event id, DLQ, EOS).
4. Give one shop example that needs same-key order (balance) and one that does not (independent clicks).

#### Medium practical tasks

1. Produce three keyed records and three null-key records. Consume. Write the order you see per partition.
2. Write a consumer that sends unparseable records to `orders.dlq` and then commits. Send one good record and one bad record.
3. Write a causal chain: create order, pay order. Show a bad design (different keys) and a good design (same order id key).

#### Advanced practical tasks

1. Design a saga of two entities (order and stock). Write how you accept that they are not in one partition. Name the compensating event.
2. Write an application standard: key rules, required event id, DLQ topic naming, unclean election off, KRaft.

---

## Do not block the poll loop

The consumer **poll loop** (topic 4) must call poll often enough. `max.poll.interval.ms` limits how long you may work between polls. If you block the thread on a slow HTTP call or a full database, the group coordinator can leave the member. A rebalance starts. Other partitions pause.

**Backpressure** means you slow the input when the output cannot keep up. Correct tools:

- pause partitions (`pause` / `resume`) while you drain a local buffer
- use a bounded queue and pause when the queue is full
- lower `max.poll.records` so each poll is smaller
- scale consumers or fix the slow sink
- use quotas on other noisy clients (topic 10), not as the only fix

Incorrect tools:

- sleep in a loop without poll
- wait for a Kafka request-reply on the same thread
- process thousands of records, then poll

Connect and Streams have their own task threads. The same idea holds: do not block a task on unbounded external I/O without a timeout and a pause strategy.

A slow poll looks like a Kafka outage. It is often your code (topic 10).

KRaft does not change poll rules. The group coordinator still lives on a broker.

### Questions

#### Theoretical questions

1. What happens if you do not poll before `max.poll.interval.ms`?
2. What is backpressure in this section?
3. What do `pause` and `resume` do?
4. Why is sleep-without-poll the wrong delay?
5. Why can a slow sink look like a Kafka failure?

#### Easy practical tasks

1. Write four sentences about the poll-loop rule. Use only facts from this section.
2. Draw: poll → process with timeout → pause if sink full → poll.
3. Make a table: Action, safe (yes or no). Add sleep, pause, smaller poll, unbounded HTTP.
4. Find `pause` and `resume` in your client documentation.

#### Medium practical tasks

1. Write a consumer that pauses when a fake sink queue is full and resumes when it drains. Produce a burst. Confirm that the member stays in the group.
2. Set a short `max.poll.interval.ms` and block on a long HTTP call. Record the rebalance. Then add a timeout. Compare.
3. Read official pause/resume notes. Write three facts in STE.

#### Advanced practical tasks

1. Design a worker that processes records in a thread pool without breaking the poll timer. Write the rules (bound the pool, pause partitions, or both). Implement a small version.
2. Write a review checklist that rejects a consumer that blocks poll on unbounded I/O.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do ECST, outbox, and CQRS form one integration path from a command API to a query API?
2. Why does "exactly-once for the business" need both a pattern (outbox or EOS) and a poll policy?
3. How do key design, duplicates, and poison messages interact on one topic?
4. What stays the same when the cluster is KRaft and the outbox is in PostgreSQL?
5. A teammate wants request-reply on Kafka for the public HTTP API. Which facts do you use in a short reply?

#### Easy practical tasks

1. Write a one-page cheat sheet: three event styles, outbox, CQRS, order, duplicate, poison, poll rule.
2. Draw your preferred shop path: command → outbox → Kafka → two read models. Label KRaft.
3. Label three of your lab topics with notification, ECST, or sourcing.
4. Bookmark one official or well-known page for outbox and one for CQRS. Write the URLs.

#### Medium practical tasks

1. Implement a tiny path: outbox row, produce with event id, idempotent consume into SQLite, query the row. Use KRaft.
2. Add a DLQ and a pause on sink failure to that consumer. Write the failure test.
3. Map each subsection to one official or standard pattern URL.

#### Advanced practical tasks

1. Design a payments bounded context: event style, outbox columns, read models, key rules, DLQ, poll timeouts. One page plus a diagram.
2. Write a team standard: default ECST, outbox required for dual write, no Kafka RPC for users, poll must not block, KRaft only.
