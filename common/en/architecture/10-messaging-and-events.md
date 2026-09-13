# 10. Messaging and Events

## Description

Messaging is communication through messages that a broker or a log stores and delivers. This topic covers commands, events, and queries; brokers versus event logs; publish/subscribe; outbox and inbox; the dual-write problem; ordering and duplicates; and poison messages.

Messaging is not a default. A function call in a modular monolith is simpler (Topic 5). Add a message when you need a time split, a burst buffer, or a boundary between owners (Topics 7 and 9). Pair event-log details with `kafka.topics.md`.

Complete Topics 1 to 9 before this topic. Use timeouts and idempotency from Topic 4. Stay defensive. Do not use brokers to hide a missing module boundary.

Use one term for each concept. A command is a request. An event is a fact. A query is a read. A broker queue and an event log are not the same tool.

---

## Commands vs events vs queries

A command asks the system to do work. The name is an imperative: `PlaceOrder`, `CapturePayment`. A command can fail. The sender often wants a result (accepted or rejected).

An event states that something happened. The name is past tense: `OrderPlaced`, `PaymentCaptured`. An event is a fact. Subscribers react. The producer of a fact does not wait for all subscribers.

A query asks for data and does not change the source of truth. `GetOrder` is a query. Queries usually stay on HTTP or RPC (Topic 8). A message-based query is possible but is often slower and harder to debug. Prefer a request-response API for queries.

Do not mix the three in one name. `UpdateOrder` that also means "order was updated" confuses senders. Send a command. After success, publish an event.

Who is responsible:

- The command handler validates and changes the aggregate (Topic 7).
- The event subscriber must not be required for the command to commit, unless you use one transaction with an outbox (later section).
- The query handler reads an approved model (Topic 9).

A message type must include a name, a version, and a correlation identifier. The correlation identifier joins logs across processes (Topic 2).

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
2. Compare choreography (events only) with a command to an orchestrator at a high level. Topic 12 covers this later.

---

## Message broker vs event log (`kafka.topics.md`)

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

Both tools are extra failure domains. You must design what happens when the broker or the log is down (Topic 11).

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

## Pub/sub

Publish/subscribe (pub/sub) is a pattern where a publisher sends a message to a topic (or exchange). Subscribers receive a copy if they subscribe. The publisher does not name each subscriber.

Pub/sub fits events. `OrderPlaced` can go to mail, search, and analytics. The order module does not import those modules (Topic 3: cycles).

Pub/sub does not fit every command. A command usually has one owner. A queue with one consumer group is clearer for "do this job".

Delivery properties vary by product:

- At-most-once: a subscriber can miss a message.
- At-least-once: a subscriber can see duplicates. This case is common.
- Exactly-once: rare as an end-to-end business guarantee. You still design idempotent handlers.

Do not assume that all subscribers are up to date at the same time. Pub/sub is eventually consistent across subscribers (Topic 7).

Filter on the subscriber side or with product features. Do not publish a giant "everything" topic if you can avoid it. Name topics after the event type or the bounded context.

The publisher must not wait for subscriber success inside the same user request unless the product is a request-response bus. That wait recreates tight coupling.

### Questions

#### Theoretical questions

1. What is pub/sub?
2. Why does pub/sub fit events?
3. Why is a command often a poor pub/sub citizen?
4. What is at-least-once delivery?
5. Why must the publisher not wait for all subscribers in the user request?

#### Easy practical tasks

1. Draw a publisher and three subscribers for `LoanReturned`.
2. Make a table: "Delivery" and "Handler duty". Add at-most-once and at-least-once.
3. Write four topic names that match a ubiquitous language.
4. List three defects of one global `events` topic.

#### Medium practical tasks

1. Design subscriptions for search and mail. Write what happens if mail is down and search is up.
2. Rewrite a design where place-order HTTP waits for three HTTP side calls. Use pub/sub instead. Write the user-visible change.
3. Write a test plan for a subscriber that must ignore events that it does not own.

#### Advanced practical tasks

1. Write a one-page pub/sub standard: topic naming, payload version, and no circular events.
2. Compare broker fan-out with a log that several consumer groups read. Use terms from the previous section.

---

## Outbox and inbox

The outbox pattern writes the business change and the outgoing message in the same database transaction. A relay process later reads the outbox table and publishes to the broker. This pattern avoids a dual write (next section).

Typical outbox steps:

1. The handler starts a transaction.
2. The handler updates the aggregate.
3. The handler inserts an outbox row (event payload, destination, status).
4. The handler commits.
5. A publisher reads unpublished rows and sends them.
6. The publisher marks rows as sent (or deletes them after a policy).

The inbox pattern stores an incoming message identifier in the consumer database before (or in the same transaction as) the side effect. A duplicate delivery finds the inbox row and skips the side effect. This pattern is idempotency for consumers.

Outbox and inbox need the same store as the business data, or a transaction that can include both. If the outbox is a different database without a distributed transaction, you did not solve dual write.

Relay publishers must be idempotent toward the broker. Brokers can also duplicate. Consumers still need inbox or another idempotency key.

Start simple: one outbox table, one relay, metrics for lag (oldest unsent row).

### Questions

#### Theoretical questions

1. What is the outbox pattern?
2. What is the inbox pattern?
3. Why must the outbox row commit with the business change?
4. How does an inbox stop a duplicate side effect?
5. Why does a second database for the outbox fail the goal?

#### Easy practical tasks

1. Number the outbox steps in your own words.
2. Draw tables: `orders` and `outbox`. Show one transaction.
3. Make a table: "Inbox column" and "Purpose". Add message id, processed time, and hash.
4. Write four sentences on relay lag as a quality measure.

#### Medium practical tasks

1. Design outbox fields for `OrderPlaced`. Include version and correlation id.
2. Write the consumer inbox algorithm in eight steps. Include a crash after side effect and before inbox write (and why that order is wrong).
3. Write an ADR: outbox in PostgreSQL, one worker publishes.

#### Advanced practical tasks

1. Write a one-page relay design: batch size, backoff, and what happens if the broker is down for one hour.
2. Compare polling the outbox with a database notify feature at a high level. Write operational trade-offs.

---

## Dual write problem

A dual write is two writes to two systems that are not in one transaction. Example: commit an order in SQL, then publish `OrderPlaced` to a broker. Either write can succeed alone.

If the database commits and the publish fails, subscribers never see the order. If the publish succeeds and the database rolls back, subscribers see an order that does not exist.

Retries do not fully fix the problem. You can still fail between the two writes. You can also duplicate the message.

The outbox pattern is the usual repair when the message must follow a database fact. The inbox pattern is the usual repair when a received message must become a database fact once.

Other repairs:

- Change the design so that there is one write (no message).
- Use a transaction that the product truly supports (rare across DB and broker).
- Accept the gap and run a periodic reconciler that compares stores (Topic 9 copies).

Do not ignore the gap. A demo that "usually works" is not an architecture.

Logs help you detect dual-write gaps. Metrics on outbox lag and on "orders without events" find silent loss.

### Questions

#### Theoretical questions

1. What is a dual write?
2. What are the two mismatch cases (commit without publish, publish without commit)?
3. Why do retries not fully solve dual write?
4. How does the outbox pattern address dual write?
5. What is a reconciler?

#### Easy practical tasks

1. Draw the two failure timelines for order plus event.
2. Write five sentences that define dual write. Use only facts from this section.
3. List three designs that avoid dual write (one write, outbox, accept and reconcile).
4. Make a table: "Symptom" and "Likely gap". Add missing mail and phantom search document.

#### Medium practical tasks

1. Take a handler that saves then publishes. Rewrite it as outbox steps on paper.
2. Write a reconciler query idea: orders from the last hour with no outbox row. Write the repair.
3. Explain why "publish then save" is also a dual write.

#### Advanced practical tasks

1. Write a one-page incident story (fictional) of a dual write in production. Include detection and the fix.
2. Compare outbox with a transaction coordinator across two products. Write why small teams prefer outbox.

---

## Ordering and duplicates

Ordering means that messages arrive in a defined sequence. Products differ. A queue can deliver out of order under retry or with more than one consumer. An event log often orders records inside a partition, not across the whole log. See `kafka.topics.md` for partition order.

Do not assume global order. If the business needs order, pick a key that maps to one partition or one queue, and keep one consumer for that key when the product requires it. Example: all events for `orderId=42` in one partition.

Duplicates happen with at-least-once delivery, with retries, and with a client that republishes. Handlers must be idempotent (Topic 4). Inbox keys, natural keys (`paymentId`), and "set state" operations reduce duplicate harm.

Last-write-wins is a choice. It can drop an older but still valid update if a late message arrives. Version numbers or event time plus a rule can help. State the rule.

Reordering plus duplicates can invert business meaning (`cancelled` then `placed` if you apply blindly). Prefer handlers that load the aggregate and apply a rule, not handlers that blindly overwrite.

Document the order guarantee that you actually have. A slide that says "ordered" without a key is a wish.

### Questions

#### Theoretical questions

1. Why is global order a wish on many brokers?
2. How can a partition key protect order for one entity?
3. Why do duplicates appear?
4. What is last-write-wins and when can it harm?
5. Why must a handler load the aggregate instead of blind overwrite?

#### Easy practical tasks

1. Write four sentences on per-entity order versus global order.
2. Make a table: "Cause" and "Duplicate? / Reorder?". Add retry, two consumers, and republish.
3. List three idempotency strategies for `PaymentCaptured`.
4. Draw two events for one order that arrive in the wrong order. Write a safe rule.

#### Medium practical tasks

1. Choose a Kafka-style partition key for order events. Write what is ordered and what is not.
2. Design a version field on an aggregate that rejects an old event.
3. Write tests: duplicate event, old event, and gap (event 3 before event 2).

#### Advanced practical tasks

1. Write a one-page ordering standard: keys, consumer count, and what the UI may assume.
2. Analyze a public Kafka ordering FAQ at a high level. Map it to a shop checkout. Do not copy long text.

---

## Poison messages

A poison message is a message that a consumer cannot process successfully, even after retries. Causes include a bad payload, a missing field after a version change, a business rejection that the handler treats as a crash, and a bug that throws on that shape.

If you retry a poison message forever, the consumer stalls. The queue grows. Healthy messages wait. Availability drops (Topic 2).

A typical handling path:

1. Retry a few times with backoff (Topic 11).
2. If the error is permanent (validation), do not retry without a bound.
3. Move the message to a dead-letter queue (DLQ) or a poison table.
4. Alert a human or a repair job.
5. Keep the original payload and the error.

Do not drop poison messages in silence. You lose data and you lose the signal.

Do not send every transient error to the DLQ on the first failure. A short outage of a dependency is not poison. Timeouts and `5xx` from a dependency are often transient.

Version mismatches are a common poison source. Use additive event fields (Topic 8 ideas). A consumer that requires a new field that old producers do not send will poison the stream.

Repair can be: fix the consumer, fix the payload, or skip with an explicit record. Skipping is a decision. Write it.

### Questions

#### Theoretical questions

1. What is a poison message?
2. Why does infinite retry harm availability?
3. What is a dead-letter queue?
4. How is a transient dependency error different from poison?
5. Why must you not drop poison messages in silence?

#### Easy practical tasks

1. Write five causes of a poison message.
2. Number a five-step poison path in your own words.
3. Make a table: "Error" and "Retry or DLQ". Add six errors.
4. List four fields to store with a dead letter.

#### Medium practical tasks

1. Design a DLQ process for mail workers: alert, inspect, replay, or skip.
2. Write a consumer policy: max retries, backoff, and what is permanent.
3. Show how an event version change can poison old consumers. Write the compatible fix.

#### Advanced practical tasks

1. Write a one-page poison-message runbook for on-call: dashboards, DLQ size, and replay steps.
2. Design a safe replay: idempotent handlers, rate limit, and a stop switch.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do commands, events, and queries map to HTTP, queues, and logs in one small system?
2. Why do outbox and inbox exist if the broker "is reliable"?
3. How do ordering, duplicates, and poison messages change handler design together?
4. When is a function call still better than pub/sub?
5. How does Topic 9 (ownership) decide who may publish which event?

#### Easy practical tasks

1. Write a one-page cheat sheet of all terms in this topic.
2. Design one async job (thumbnail) with a queue and one domain event (`PhotoUploaded`) with pub/sub. Keep them distinct.
3. Write three ADR titles: no log yet, outbox for order events, DLQ for mail.
4. Bookmark `kafka.topics.md`. Write five questions you will answer when you start that path.

#### Medium practical tasks

1. Write a two-page messaging design for a library: which facts are events, which jobs are commands, and how dual write is avoided.
2. Draw the full path: HTTP write → transaction + outbox → relay → topic → inbox → SQL.
3. Role-play a request to "just emit everything to Kafka". Write a smaller plan.

#### Advanced practical tasks

1. Write a messaging standard (three pages) for a team of four: message kinds, outbox, idempotency, DLQ.
2. Design a drill: broker down for 30 minutes. Write user impact, outbox growth, and recovery. Topic 11 will reuse this drill.
