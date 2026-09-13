# 7. Messaging

## Description

Messaging is communication through messages that a broker or a log stores and delivers. This topic covers commands, events, and queries; brokers versus event logs; outbox and inbox; the dual-write problem; ordering, duplicates, and poison messages.

Messaging is not a default. A function call in a modular monolith is simpler (Topic 3). Add a message when you need a time split, a burst buffer, or a boundary between owners (Topics 4 and 6). Pair event-log details with `kafka.topics.md`.

Complete Topics 1 to 6 before this topic. Use timeouts and idempotency from Topic 2. Stay defensive. Do not use brokers to hide a missing module boundary.

Use one term for each concept. A command is a request. An event is a fact. A query is a read. A broker queue and an event log are not the same tool.

---

## Commands vs events vs queries

A command asks the system to do work. The name is an imperative: `PlaceOrder`, `CapturePayment`. A command can fail. The sender often wants a result (accepted or rejected).

An event states that something happened. The name is past tense: `OrderPlaced`, `PaymentCaptured`. An event is a fact. Subscribers react. The producer of a fact does not wait for all subscribers.

A query asks for data and does not change the source of truth. `GetOrder` is a query. Queries usually stay on HTTP or RPC (Topic 5). A message-based query is possible but is often slower and harder to debug. Prefer a request-response API for queries.

Do not mix the three in one name. `UpdateOrder` that also means "order was updated" confuses senders. Send a command. After success, publish an event.

Who is responsible:

- The command handler validates and changes the aggregate (Topic 4).
- The event subscriber must not be required for the command to commit, unless you use one transaction with an outbox (later section).
- The query handler reads an approved model (Topic 6).

A message type must include a name, a version, and a correlation identifier. The correlation identifier joins logs across processes (Topic 11).

Publish/subscribe (pub/sub) fits events. A publisher sends a message to a topic. Subscribers receive a copy. The publisher does not name each subscriber. A command usually has one owner. A queue with one consumer group is clearer for "do this job".

The publisher must not wait for subscriber success inside the same user request unless the product is a request-response bus. That wait recreates tight coupling.

### Questions

#### Theoretical questions

1. What is a command?
2. What is an event?
3. What is a query?
4. Why do queries usually stay on HTTP or RPC?
5. Why must you not mix a command name and an event name?

#### Easy practical tasks

1. Classify twelve names as command, event, or query.
2. Make a table: "Sender wants" and "Message kind". Add "please pay", "paid", and "show receipt".
3. Write four sentences that explain a failed command versus an event.
4. List five fields that every message header must include.

#### Medium practical tasks

1. Design the flow: HTTP `POST /orders` → command `PlaceOrder` → event `OrderPlaced` → mail subscriber.
2. Rewrite a mixed API `POST /update-and-notify` into a command and an event.
3. Write when a worker must send a command to another owner versus publish an event.

#### Advanced practical tasks

1. Write a one-page catalog: ten messages, kind, owner, and success or fail behavior.
2. Compare choreography (events only) with a command to an orchestrator at a high level. Topic 9 covers this later.

---

## Broker vs event log

A message broker (for example RabbitMQ-style queues) delivers a message to consumers. After a consumer acknowledges, the broker can remove the message from the queue (depending on the product). The focus is work to do. Competing consumers share a queue. One message is processed by one worker in the group.

An event log (for example Apache Kafka) appends records to a durable, ordered log. Consumers keep an offset. They can replay. The focus is a history of facts. See `kafka.topics.md`.

Use a broker queue when:

- You have a job (send mail, make a thumbnail).
- You want competing workers.
- You do not need a long replay of all facts.

Use an event log when:

- Many independent consumers need the same facts.
- You need replay and a longer retention.
- You accept more operational cost.

A log is not a database for arbitrary queries. A queue is not a full audit history if you delete messages.

Do not run both products on day one. A small team can start with a managed queue for one job. Add a log when you have several subscribers and a replay need.

Both tools are extra failure domains. You must design what happens when the broker or the log is down (Topic 8).

Delivery properties vary by product:

- At-most-once: a subscriber can miss a message.
- At-least-once: a subscriber can see duplicates. This case is common.
- Exactly-once: rare as an end-to-end business guarantee. You still design idempotent handlers.

### Questions

#### Theoretical questions

1. What is a message broker in this handbook?
2. What is an event log?
3. When do you use a queue?
4. When do you use a log?
5. Why is a log not a general query database?

#### Easy practical tasks

1. Write five sentences that compare queue and log.
2. Make a table: "Need" and "Queue or log". Add six needs (thumbnail, audit replay, one worker).
3. List four operational tasks that either product adds.
4. Write a pointer to `kafka.topics.md` and three terms you will learn there (topic, partition, offset).

#### Medium practical tasks

1. Choose a tool for "send registration email" and for "billing consumes all orders". Write reasons.
2. Write an ADR that starts with one managed queue and rejects a log for a student app.
3. Draw consumer groups on a queue versus two independent readers on a log.

#### Advanced practical tasks

1. Write a one-page operations comparison: retention, replay, and people cost for a five-person team.
2. Read the Kafka getting-started idea in `kafka.topics.md`. Map "topic" and "offset" to this section in ten sentences.

---

## Outbox, inbox, and dual write

A dual write is two writes to two stores that you cannot commit as one transaction. Example: write an order row to the database, then publish `OrderPlaced` to a broker. If the process stops between the two writes, the system is wrong. Either the event is missing or the event exists without the row.

The transactional outbox pattern stores the business row and an outbox row in one database transaction. A relay reads the outbox and publishes to the broker. The write of the fact and the write of the "please publish" record stay together.

An inbox (or processed-message store) records message identifiers that a consumer already handled. On a duplicate delivery, the consumer skips the side effect. This is idempotency for messages (Topic 2).

Rules:

- Do not dual-write and hope.
- Use an outbox when the source of truth and the broker are different products.
- Use an inbox or another idempotency store on the consumer.
- The relay must be restartable. The outbox row stays until publish succeeds.
- Bound outbox growth. A stuck relay is an operations incident.

In a modular monolith, an in-process event can still use an outbox if a later subscriber runs in another process. If all subscribers run in the same transaction as the write, you do not need a broker. Prefer that simple path when you can.

The outbox does not make subscribers transactional with the writer. Subscribers still see the fact later (Topic 4). The outbox only prevents a missing or extra publish relative to the source write.

### Questions

#### Theoretical questions

1. What is a dual write?
2. What does a transactional outbox store in one transaction?
3. What does an inbox store?
4. Why does an outbox not make subscribers transactional with the writer?
5. When can you skip a broker and still publish in process?

#### Easy practical tasks

1. Write five sentences that define dual write, outbox, and inbox.
2. Draw a timeline where the process stops after the database write and before the broker publish.
3. Make a table: "Pattern" and "Problem it reduces". Add outbox and inbox.
4. List four fields of an outbox row (id, payload, created time, published flag).

#### Medium practical tasks

1. Design an outbox for `OrderPlaced`. Write the transaction contents and the relay steps.
2. Write an inbox check for a mail worker that can see the same event twice.
3. Explain in eight sentences why "write DB then publish" is not safe.

#### Advanced practical tasks

1. Write a one-page outbox standard: table shape, relay, retry, and what happens if the broker is down.
2. Compare polling the outbox with a database notify feature at a high level. Write operations cost. Do not require a specific product.

---

## Ordering, duplicates, and poison messages

Ordering is the rule that says which message a consumer sees first. Many brokers do not give a global order for all messages. An event log can give order inside a partition or a key. If two facts must stay in order, you must put them on the same ordered stream (same key, same partition, or the same queue with one consumer).

Do not assume global order across all topics. Design the consumer so that an out-of-order pair is safe, or put the pair on one stream.

Duplicates are normal with at-least-once delivery. A consumer can see the same message twice after a crash or a timeout. Handlers must be idempotent. Use an inbox, a unique business key, or a stored idempotency key.

A poison message is a message that fails every time. Examples: invalid JSON, a missing required field, or a rule that always rejects. Infinite retry of a poison message blocks the queue or wastes workers.

Rules for poison messages:

- Bound the retry count.
- After the bound, move the message to a dead-letter store.
- Alert an operator (Topic 11).
- Do not drop a business fact without a record.
- Fix the producer or the handler. Then replay from the dead-letter store if the fact is still valid.

A validation error is often not worth a retry. A timeout of a dependency can be worth a retry with backoff (Topic 8).

Write the order key in the contract. Example: all events for one `orderId` stay in one partition. Events for different orders can interleave.

### Questions

#### Theoretical questions

1. Why is global order across all topics a false default?
2. How do you keep two facts in order?
3. Why are duplicates normal?
4. What is a poison message?
5. What do you do after the retry bound?

#### Easy practical tasks

1. Write five sentences about order, duplicates, and poison messages.
2. Make a table: "Failure" and "Retry? (yes/no)". Add invalid JSON, timeout, and 400 validation.
3. List four idempotency tools for a consumer (inbox, unique key, and similar).
4. Draw one partition key for all events of one order.

#### Medium practical tasks

1. Design a dead-letter path for a mail worker. Include operator notice and replay.
2. Write an ADR: at-least-once plus inbox, no claim of exactly-once business delivery.
3. Explain in eight sentences how a timeout plus retry can create a duplicate charge without a key.

#### Advanced practical tasks

1. Write a one-page consumer standard: order key, idempotency, retry, dead letter, and metrics.
2. Compare "one consumer for a queue" with "many consumers and a partition key" for order of events. Write when each fits.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do commands, events, and queries map to HTTP APIs versus brokers?
2. When is a function call in a monolith a better default than any message tool?
3. How do outbox and inbox together address dual write and duplicates?
4. Why is "exactly-once" a rare end-to-end business guarantee?
5. What Topic 6 ownership rule do you break if two services consume and rewrite the same table after an event?

#### Easy practical tasks

1. Write a one-page cheat sheet: command/event/query, queue versus log, outbox/inbox, order, duplicates, poison.
2. For a to-do app, name one command, one event, one query path (HTTP), and no broker.
3. Draw the unsafe dual write and the outbox fix as two sequence diagrams.
4. Bookmark `kafka.topics.md`. Write three terms that this topic uses.

#### Medium practical tasks

1. Write a short messaging brief for a campus shop: one queue for mail, no event log, outbox on orders.
2. Take a teammate design that publishes from the HTTP handler after a separate commit. Rewrite it with an outbox.
3. Write a twelve-week plan: weeks for one async job, and the measure that would reopen an event log.

#### Advanced practical tasks

1. Write a failure matrix: broker down, relay down, consumer down. Write user-visible behavior for place-order.
2. Read the Enterprise Integration Patterns site index at [https://www.enterpriseintegrationpatterns.com/](https://www.enterpriseintegrationpatterns.com/). Map five pattern names to this topic in your own words.
