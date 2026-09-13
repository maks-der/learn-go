# 19. Application Patterns

## Description

Application patterns are ways to use Kafka in a system design. The broker does not choose a pattern. You choose what a record means and how a service writes and reads it.

This topic covers event notification versus event-carried state transfer versus event sourcing, the outbox pattern, CQRS with Kafka, request-reply over Kafka, ordering and causality, and backpressure on the poll loop. Complete this topic after topics 4, 5, 7, 8, and 11. Use a KRaft cluster. These patterns do not need ZooKeeper.

Use one term for each concept. An event notification is a small signal that something changed. Event-carried state transfer is an event that carries the data that a reader needs. Event sourcing stores every change as an event and builds state from that log. An outbox is a table (or equivalent) in the same database transaction as the business write. CQRS splits write models and read models. Backpressure is how you slow input when the consumer cannot keep up.

---

## Event notification vs event-carried state transfer vs event sourcing

**Event notification** says that a fact occurred and gives a pointer. Example: `{ "orderId": "55" }`. The consumer must call an API or read a database to get the order. Many consumers create load on the source. If the source is down, the event is not enough.

**Event-carried state transfer (ECST)** puts the needed fields in the event. Example: order id, items, totals, customer id. Consumers can update their own store without a call back. The event is larger. Schema evolution matters (topic 11).

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
3. Find a public article on ECST (for example Greg Young or a known integration book idea). Rewrite three facts in STE. Do not copy a long passage.

#### Advanced practical tasks

1. Design five events for payments. Mark each as notification, ECST, or sourcing. Write the schema fields at a high level.
2. Write a one-page standard: default ECST, when notification is allowed, when event sourcing is a team decision.

---

## Outbox pattern

The **outbox pattern** solves this failure: a service writes a database row and then fails before it produces to Kafka (or the reverse). The two stores diverge.

The service writes the business row and an **outbox row** in the **same database transaction**. A separate publisher reads the outbox and produces to Kafka. After a successful produce (and the durability that you require, topic 4 and 7), the publisher marks the outbox row as sent or deletes it.

If the process stops after the commit and before the produce, the publisher retries. Kafka can see a duplicate if the produce succeeded and the mark did not. Use an idempotent producer and a deduplication key (topic 7).

Debezium can capture the outbox table (or a dedicated outbox schema) as CDC (topic 12). Then the application does not run a poll loop on the table. The transaction log becomes the publisher.

Do not produce first and write the database second without a compensating plan. Do not use a two-phase commit across Kafka and a SQL database as the default.

The Kafka cluster is KRaft. The outbox lives in your database. Connect or a worker reads it.

### Questions

#### Theoretical questions

1. What dual-write failure does the outbox prevent?
2. What two writes share one database transaction?
3. Who produces to Kafka in this pattern?
4. Why can Kafka still see a duplicate?
5. How can Debezium replace an application poll of the outbox?

#### Easy practical tasks

1. Write five sentences about the outbox. Use only facts from this section.
2. Draw: HTTP handler → DB transaction (order + outbox) → publisher → Kafka.
3. Make a table: Step, can fail, result if no outbox, result with outbox.
4. List two ways to publish the outbox (app poll, CDC).

#### Medium practical tasks

1. Implement a small outbox table and a publisher in a language that you know. Use KRaft Kafka. Crash after commit, before produce. Confirm a retry.
2. Add an event id. Consume with an idempotent handler (topic 7). Write a duplicate produce test.
3. Read a Debezium outbox router note if you use Debezium. Write how the topic name is chosen.

#### Advanced practical tasks

1. Compare outbox+CDC with a process that writes only to Kafka and builds a read model (no SQL write on the command). Write five trade-offs.
2. Write an outbox standard: required columns, who publishes, EOS or idempotence, KRaft bootstrap, no dual-write without outbox.

---

## CQRS with Kafka

**CQRS** (command query responsibility segregation) splits the write model from the read model. Commands change the write side. Queries read a model that is shaped for the query.

Kafka is often the pipe between them. The write service commits (database and/or events). Events go to Kafka (outbox or event sourcing). A consumer builds a read store: search index, cache, or SQL view. The read API does not query the write database.

The read model is **eventually consistent**. A query can miss a command that is still in the log or in the consumer. Tell the product team. Do not hide this.

You can have more than one read model from one topic. That is a reason to use Kafka instead of one shared database for all queries.

Interactive queries in Kafka Streams (topic 14) are a CQRS-style read of stream state. A database read model is the other common form.

Do not put every query on Kafka consume-from-beginning. Use a stored read model. Kafka retention is not your query API.

KRaft hosts the event topics. The read database is separate.

### Questions

#### Theoretical questions

1. What does CQRS split?
2. What role does Kafka play in a common CQRS path?
3. Why is the read model eventually consistent?
4. Why can one topic feed two read models?
5. Why is "query Kafka from the beginning on each HTTP request" a poor default?

#### Easy practical tasks

1. Write five sentences about CQRS and Kafka. Use only facts from this section.
2. Draw: command API → write model → Kafka → consumer → read DB → query API.
3. Make a table: Question, write model or read model.
4. Give one example where the UI must tolerate a short delay.

#### Medium practical tasks

1. Build a tiny CQRS lab: produce an ECST event, consume into a map or SQLite, query the map. Use KRaft.
2. Write six sentences on what the query returns if the consumer lags (topic 16).
3. Compare Streams interactive queries with a SQLite read model. Write four differences.

#### Advanced practical tasks

1. Add a second read model (for example, counts per day) from the same topic. Show both views after three events.
2. Write a CQRS standard: event pattern (ECST), how you measure lag, how the API talks about delay, KRaft only.

---

## Request-reply over Kafka (usually avoid)

**Request-reply** means a client sends a request and waits for a response that matches that request. HTTP and RPC do that well. Kafka is a log. It does not give you a socket per call.

People implement request-reply with two topics: requests and replies. The request carries a correlation id. The client waits for a reply with that id, often on a temporary partition or a dedicated reply topic.

Problems:

- Timeout and retry create duplicates. You need idempotence.
- Ordering of replies is not the same as HTTP connection lifetime.
- A stuck consumer looks like a slow API (poll loop, topic 5).
- You build a correlation store.

**Usually avoid** this pattern for user-facing request-reply. Use HTTP or gRPC to the service. Use Kafka for events that do not need a synchronous answer.

If you still do it (for example, a backend integration that already is Kafka-only), set short timeouts, bound the wait, and do not block a shared poll thread while you wait for a reply (next section).

KRaft does not make request-reply a good default.

### Questions

#### Theoretical questions

1. What extra field does Kafka request-reply need?
2. Why is Kafka a poor default for user HTTP-style calls?
3. What does a retry of a request create?
4. When might a team still use request-reply on Kafka?
5. Why does KRaft not fix this pattern?

#### Easy practical tasks

1. Write four sentences about why this handbook says "usually avoid".
2. Draw HTTP request-reply versus Kafka two-topic correlation.
3. Make a table: Need, use HTTP, use Kafka events.
4. List three failure modes (timeout, duplicate, stuck consumer).

#### Medium practical tasks

1. Write a sequence of five steps for a Kafka request-reply. Mark where you wait. Mark where you must not block poll.
2. Compare this pattern with CQRS (async result). Write when the UI polls a read model instead.
3. Find an official or well-known note that discourages using Kafka as an RPC bus. Rewrite three facts in STE.

#### Advanced practical tasks

1. If you implement a lab: one request topic, one reply topic, correlation id. Measure p99 wait. Write why you would still prefer HTTP.
2. Write a team rule: Kafka is not the user-facing RPC. List the two exceptions that need a written design.

---

## Ordering and causality

**Ordering** in Kafka is per partition (topic 3). Records in one partition have a total order. Records in two partitions do not.

If you need a happens-before relation between two events, put them in the same partition (same key) or put a causal token in the payload (for example, a version per entity).

**Causality** is the rule: if event A caused event B, readers must be able to see A before they treat B as complete. Same-key partitioning gives that for one entity. Cross-entity causality needs a design (sagas, versions, or a single writer).

Clock timestamps are not enough. Two producers can have skewed clocks. Kafka record timestamps can be broker time or create time. Do not sort global events by timestamp and call that causality.

Consumer groups do not change partition order. Two groups can read at different offsets. One group still reads one partition in order on one consumer.

Rebalances can pause a partition. They do not reorder the log.

Geo copies (topic 18) can delay a partition. A reader of the target must not assume it is caught up with the source.

KRaft leader election does not reorder a committed log. Unclean leader election can (topic 9). Keep unclean election off.

### Questions

#### Theoretical questions

1. Where does Kafka guarantee order?
2. How do you keep two events in order for one entity?
3. Why are clocks a poor causality tool?
4. Do consumer groups change partition order?
5. How can unclean leader election break order or durability?

#### Easy practical tasks

1. Write five sentences about order and causality. Use only facts from this section.
2. Draw two partitions. Place events for user 9 and user 8. Mark what is ordered.
3. Make a table: Need, same key, extra version field, single writer.
4. Give one shop example that needs same-key order (balance) and one that does not (independent clicks).

#### Medium practical tasks

1. Produce three keyed records and three null-key records. Consume. Write the order you see per partition.
2. Write a causal chain: create order, pay order. Show a bad design (different keys) and a good design (same order id key).
3. Read official notes on timestamps. Write create time versus log append time in STE.

#### Advanced practical tasks

1. Design a saga of two entities (order and stock). Write how you accept that they are not in one partition. Name the compensating event.
2. Write an ordering standard: key rules, no global timestamp sort, unclean election off, KRaft.

---

## Backpressure: do not block the poll loop

The consumer **poll loop** (topic 5) must call poll often enough. `max.poll.interval.ms` limits how long you may work between polls. If you block the thread on a slow HTTP call or a full database, the group coordinator can leave the member. A rebalance starts. Other partitions pause.

**Backpressure** means you slow the input when the output cannot keep up. Correct tools:

- pause partitions (`pause` / `resume`) while you drain a local buffer
- use a bounded queue and pause when the queue is full
- lower `max.poll.records` so each poll is smaller
- scale consumers or fix the slow sink
- use quotas on other noisy clients (topic 17), not as the only fix

Incorrect tools:

- sleep in a loop without poll
- wait for a Kafka request-reply on the same thread
- process thousands of records, then poll

Connect and Streams have their own task threads. The same idea holds: do not block a task on unbounded external I/O without a timeout and a pause strategy.

A slow poll looks like a Kafka outage. It is often your code (topic 16).

KRaft does not change poll rules. The group coordinator still lives on a broker.

### Questions

#### Theoretical questions

1. What happens if you do not poll before `max.poll.interval.ms`?
2. What is backpressure in this section?
3. What do `pause` and `resume` do?
4. Why is sleep-without-poll the wrong delay?
5. Why can a slow sink look like a Kafka failure?

#### Easy practical tasks

1. Write five sentences about the poll loop and backpressure. Use only facts from this section.
2. Make a table: Action, safe (yes or no), reason.
3. Find `pause` in the consumer API docs for your language. Write the method name.
4. Draw: poll → bounded queue → workers → pause when full.

#### Medium practical tasks

1. Write a consumer that processes too slowly (sleep). Watch a rebalance (topic 22 also practices this). Write the log lines.
2. Change the same program to pause partitions instead of blocking poll. Write the difference.
3. Read `max.poll.interval.ms` and `max.poll.records` docs. Write how they work together.

#### Advanced practical tasks

1. Implement a bounded queue with pause/resume. Load the sink slowly. Confirm the group stays stable and lag grows on purpose.
2. Write a consumer standard: max work per poll, pause rules, no request-reply on the poll thread, KRaft cluster.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do ECST, outbox, and CQRS form one integration path from a command to a query?
2. Why does this handbook reject Kafka request-reply as a default but accept Kafka as the event log?
3. How do ordering keys and the poll loop both protect correctness under load?
4. Which pattern needs topic 7 (idempotence) the most, and why?
5. What stays the same when these patterns run on KRaft instead of an old cluster?

#### Easy practical tasks

1. Write a one-page cheat sheet: three event styles, outbox, CQRS, avoid RPC, order, pause.
2. Draw a full shop path: HTTP command, outbox, Kafka, CQRS read, query HTTP. No reply topic.
3. Bookmark one outbox and one CQRS explanation that you trust. Write the URLs.
4. Label six real features in a product you know with a pattern from this topic.

#### Medium practical tasks

1. Implement a thin path: outbox or direct ECST produce, consumer read model, no request-reply. Use KRaft.
2. Map each subsection to one official or classic-pattern URL.
3. Write a review checklist for a pull request that adds a Kafka topic (pattern, key, poll, schema).

#### Advanced practical tasks

1. Add a slow sink and pause-based backpressure to the CQRS consumer. Show lag without a rebalance.
2. Write an application-pattern standard for your team: default ECST+outbox+CQRS, bans on dual-write and Kafka RPC, KRaft only.
