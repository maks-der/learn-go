# 11. Distributed Systems Basics

## Description

A distributed system is a system that uses more than one process that communicate over a network. This topic covers the fallacies of distributed computing, partial failure, timeouts and retries, idempotent handlers, circuit breakers, bulkheads, CAP in practical words, and the idea of consensus.

Distribution is not a goal. It is a cost that you pay when one process is not enough (Topics 5 and 9). Every network hop adds failure modes. Complete Topics 1 to 10 before this topic. Pair with `net.topics.md` if sockets and latency are new.

Use one term for each concept. A timeout is not a retry. A circuit breaker is not a bulkhead. CAP is a teaching model, not a product checkbox. You rarely implement Raft or Paxos yourself.

---

## Fallacies of distributed computing

The fallacies of distributed computing are a known list of false assumptions. L. Peter Deutsch and others recorded them. This handbook uses them as a checklist.

The usual list:

1. The network is reliable.
2. Latency is zero.
3. Bandwidth is infinite.
4. The network is secure.
5. Topology does not change.
6. There is one administrator.
7. Transport cost is zero.
8. The network is homogeneous.

Each item is false in production. Packets drop. Calls take milliseconds to hundreds of milliseconds. Links fill. Networks are hostile. Routes change. Many teams share the path. Data transfer costs money. Protocols and versions differ.

If you design as if the list were true, you omit timeouts, you ignore P99 latency, you skip authentication, and you treat a remote function like a local function.

A local function call can fail too, but the failure is usually the same process. A remote call can fail while the other process continues. That is partial failure (next section).

Use the list in reviews. For each new remote dependency, ask: which fallacy did we assume?

### Questions

#### Theoretical questions

1. What are the fallacies of distributed computing?
2. Why is "latency is zero" false for a user-facing budget?
3. Why is "the network is secure" a dangerous assumption?
4. How does "there is one administrator" fail in a company?
5. Why must a remote call not look like a local call in the design?

#### Easy practical tasks

1. Write the eight fallacies in your own short sentences.
2. Make a table: "Fallacy" and "Design control". Add timeout, TLS, and a latency budget for three rows. Then finish all eight.
3. List four defects that appear if you believe the network is reliable.
4. Write five sentences that compare a local call and a remote call.

#### Medium practical tasks

1. Review a homework client that calls an HTTP API with no timeout. Map each risk to a fallacy.
2. Write a one-page review checklist that quotes the eight items as questions.
3. Estimate a latency budget that includes three network hops. Show that "zero latency" would hide the problem.

#### Advanced practical tasks

1. Write a one-page history note: who wrote the fallacies and why they still apply. Use public sources. Do not copy long passages.
2. Apply the list to a public cloud architecture diagram. Mark three assumed fallacies.

---

## Partial failure

A partial failure is a state where some parts of the system work and other parts do not. In one process, a crash often stops all work. In a distributed system, the API can be up while the mail worker is down. The database can accept writes while the cache is down. One replica can be stale.

Partial failure is the normal failure mode. You must define behavior for each dependency:

- Fail the user request.
- Degrade (show cached data, hide a widget).
- Queue work for later.
- Serve a subset of features.

If you do not define the behavior, the user sees random errors and long waits.

Detection is hard. A process that does not reply can be dead, slow, or partitioned from you. You cannot always know. Timeouts turn "unknown" into a local decision (Topic 4). That decision can be wrong: the work may still complete.

Write a failure matrix: dependency × behavior × user message. Practice it. Topic 15 in the later path covers graceful degradation in more depth. This section is the start.

Health checks help, but a process can pass `/health` and still fail on the real query. Check the dependency that matters, with a bound.

### Questions

#### Theoretical questions

1. What is a partial failure?
2. Why is partial failure common in a distributed system?
3. What three behaviors can you choose when a dependency is down?
4. Why is "no reply" not the same as "did not run"?
5. Why can a simple health check lie?

#### Easy practical tasks

1. Write four sentences that define partial failure.
2. Make a matrix: four dependencies (DB, cache, mail, search) and the user-visible behavior if each one is down.
3. List five user messages that are honest ("search is delayed") and two that hide the truth.
4. Draw a system with three processes. Mark one process down. Write what still works.

#### Medium practical tasks

1. Design degrade rules for a shop home page when search is down and when the catalog database is down. They must differ.
2. Write a failure matrix for place-order: payments down versus email down.
3. Explain in eight sentences why a two-node system has more visible states than "up" and "down".

#### Advanced practical tasks

1. Write a one-page partial-failure catalog for six dependencies. Include detect method and max wait.
2. Design a game-day script (paper) that turns off one dependency. List expected signals. Do not attack systems that you do not own.

---

## Timeouts, retries, exponential backoff

A timeout bounds the wait (Topic 4). A retry repeats the call after a failure. Exponential backoff increases the wait between retries (for example 100 ms, 200 ms, 400 ms, plus a random jitter).

Retries help when the failure is transient: a short network drop, a brief overload, a rolling deploy. Retries harm when the failure is permanent (bad request) or when all clients retry at once (retry storm).

Rules:

- Retry only idempotent operations, or operations that you protected with an idempotency key (Topics 4, 8, and 10).
- Bound the retry count.
- Use backoff and jitter so that clients do not align.
- Do not retry `400` validation errors.
- Keep the total time inside the user budget.
- Apply a timeout on each attempt, not one timeout for the whole loop without a plan.

A retry of a non-idempotent `POST` can charge a card twice. That defect is worse than a user retry.

Jitter is a random extra delay. Without jitter, many workers wake at the same time and hit the dependency together.

Server-side retries plus client-side retries multiply load. Write who retries. Prefer one layer.

Metrics: retry count, timeout count, and success after retry. If success after retry is rare, fix the dependency. Do not hide a chronic fault with retries.

### Questions

#### Theoretical questions

1. What is a retry?
2. What is exponential backoff?
3. What is jitter and why do you add it?
4. Which errors must you not retry?
5. Why can client retries plus server retries multiply load?

#### Easy practical tasks

1. Write a retry policy: max 3 attempts, timeouts of 200 ms, backoff 100/200/400 ms, jitter 0–50 ms.
2. Make a table: "Error" and "Retry? (yes/no)". Add timeout, 400, 401, 429, and 503.
3. Show arithmetic: three retries that break a 1-second user budget.
4. Write five sentences on retry storms.

#### Medium practical tasks

1. Design retries for `GET` of a public catalog versus `POST` of a payment. They must differ.
2. Write an ADR: HTTP client timeouts and retries for all outbound calls.
3. Explain `429` with `Retry-After` and how a client must behave.

#### Advanced practical tasks

1. Write a one-page retry standard: who retries, idempotency, budget, and metrics.
2. Simulate on paper 1000 clients that retry at the same instant without jitter versus with jitter. Write the load shape.

---

## Idempotent handlers

An idempotent handler can process the same message or the same request more than once without a second business effect. Distributed systems deliver duplicates (Topic 10). Clients retry (previous section). Handlers must tolerate both.

Methods:

- Natural idempotency: set a state (`status=paid`) that is already true.
- Idempotency key or inbox id stored with the effect.
- Unique business key (`payment_provider_ref`) with a uniqueness rule in the database.
- Compare-and-set with a version.

The handler must also be safe if it crashes after the effect and before the ack. The message returns. The second run must see the stored key and stop.

Logging must not explode. A duplicate is expected. Log at a low rate or with a counter. Do not page a human on every duplicate.

Testing is required. Write tests for: first delivery, immediate duplicate, duplicate after a crash, and conflicting payload with the same id.

Idempotent handlers do not remove the need for poison-message handling. A bad payload that throws before the inbox write will retry and can become poison (Topic 10).

### Questions

#### Theoretical questions

1. Why do distributed handlers see duplicates?
2. Name three methods to make a handler idempotent.
3. Why must the effect and the inbox id use one transaction when you can?
4. Why is a crash after the effect and before the ack a normal case?
5. How do idempotent handlers and poison messages differ?

#### Easy practical tasks

1. Write an algorithm in eight steps for an inbox-backed handler.
2. Make a table: "Handler" and "Idempotency method". Add send mail, capture payment, and update search document.
3. List four tests for one handler.
4. Write four sentences on why "just retry" is incomplete without idempotency.

#### Medium practical tasks

1. Design a unique constraint for `PaymentCaptured` so that a duplicate insert fails safely.
2. Write the conflicting-payload case (same id, different amount). Choose `409` or a dead letter.
3. Draw a timeline: effect, crash, redelivery, no second charge.

#### Advanced practical tasks

1. Write a one-page handler template that a team copies: timeout, retry, inbox, metrics.
2. Compare at-least-once plus idempotent handlers with a vendor "exactly-once" claim. Write what still remains in your code.

---

## Circuit breaker

A circuit breaker is a guard that stops calls to a dependency after a threshold of failures. The breaker opens. Calls fail fast (or use a fallback) without waiting for a long timeout. After a sleep, the breaker allows a trial call (half-open). A success closes the breaker. A failure opens it again.

States:

- Closed: calls pass through. Failures increment a counter.
- Open: calls do not go to the dependency. They fail immediately.
- Half-open: a limited number of trial calls.

The goal is to protect your process and to give the dependency time to recover. If all handlers wait 30 seconds on a dead host, your threads exhaust. You fail even the paths that do not need that host.

A circuit breaker is not a retry. It is a stop. Combine them with care. Retries against an open breaker waste nothing if you fail fast. Retries that start when the breaker is closed can still storm. Use backoff.

Fallback must be safe. Cached data can be a fallback. A fake success for a payment is not a fallback.

Configure from measures: error rate, slow-call rate, window, and open duration. A breaker that never opens is decoration. A breaker that opens on one error is noise.

### Questions

#### Theoretical questions

1. What is a circuit breaker?
2. What are the three states?
3. Why does fail-fast help your process?
4. How is a breaker different from a retry?
5. What is a safe fallback?

#### Easy practical tasks

1. Draw the three states and the transitions.
2. Write four sentences that define open versus half-open.
3. Make a table: "Dependency" and "Fallback". Add search, recommendations, and payments.
4. List five settings (threshold, window, sleep) in one sentence each.

#### Medium practical tasks

1. Design a breaker for a payment vendor: when to open, what the user sees, and when to try again.
2. Explain in eight sentences how a breaker prevents thread exhaustion.
3. Write an ADR: breakers on all outbound HTTP, no breaker on the in-process function calls.

#### Advanced practical tasks

1. Write a one-page breaker standard: metrics, dashboards, and who may change thresholds.
2. Compare a client-side breaker with a mesh or gateway breaker. Write who owns the state.

---

## Bulkhead

A bulkhead is a partition that limits the resources that one dependency or one workload can consume. The name comes from ship compartments. A leak in one compartment must not fill the ship.

In software, bulkheads are separate thread pools, connection pools, process limits, or queues. If the report export uses all HTTP workers, login can still work if login has its own pool.

Without bulkheads, one slow client or one slow dependency can take the whole process. That is a common outage shape.

Examples:

- A connection pool per database or per tenant.
- A worker pool for "slow reports" and a pool for "interactive API".
- A queue length limit (backpressure) so that accept stops when full.
- Separate processes (Topic 4) for the interactive API and the worker.

Bulkheads have a cost. Idle capacity in one pool cannot serve the other pool. Size the pools from measures.

Bulkheads and circuit breakers work together. The breaker stops calling a dead dependency. The bulkhead stops that dependency from taking all threads even when the breaker is closed and the dependency is only slow.

Do not invent twenty pools for a student app. Start with one process and one timeout. Add a pool split when a measure shows starvation.

### Questions

#### Theoretical questions

1. What is a bulkhead?
2. Why does a slow dependency threaten the whole process?
3. Give three software forms of a bulkhead.
4. How do bulkheads and circuit breakers differ?
5. What is the cost of many pools?

#### Easy practical tasks

1. Write four sentences that define a bulkhead. Use the ship picture in one sentence only.
2. Draw one process with two pools: API and reports.
3. Make a table: "Workload" and "Resource cap". Add four workloads.
4. List three signals that you need a bulkhead (queue wait, pool wait, timeout on unrelated routes).

#### Medium practical tasks

1. Design connection-pool sizes for a database: interactive 20, batch 5. Write what happens when batch is full.
2. Write an ADR that splits the API process from the worker process as a bulkhead.
3. Explain backpressure: reject or slow accept when a queue is full. Write the user-visible error.

#### Advanced practical tasks

1. Write a one-page bulkhead map for a three-dependency API. Include threads, connections, and queues.
2. Compare process isolation with in-process pools. Use Topic 4 terms.

---

## CAP in practical words

CAP is a teaching model for a distributed data store during a network partition. The letters mean:

- C — consistency: all nodes that accept reads show the same latest write (a strong form).
- A — availability: every request to a non-failing node gets a response (not an error that says "I cannot answer").
- P — partition tolerance: the system continues even if the network splits the nodes.

During a partition, you cannot have both strong consistency and full availability in the sense of the model. You must choose a rule:

- Refuse some writes or reads until the partition heals (favor a consistency rule).
- Accept writes on both sides and resolve later (favor availability). You then have conflict.

Practical words for a beginner:

- If two replicas cannot talk, do you still take an order on both replicas?
- If yes, you can take conflicting orders for the last seat. You must repair.
- If no, one side tells the user "try later". You lose some availability.

CAP is not a sticker on a product ("we are AP"). Real systems mix rules per operation. A password change may refuse during a partition. A public read of a blog post may use a stale replica.

Latency is not in the three letters. PACELC is a later teaching extension: if there is a partition, trade A and C; else trade latency and consistency. You do not need the acronym to use the idea. Fast reads from a nearby replica can be stale.

Write the rule per use case. "The system is consistent" is empty.

### Questions

#### Theoretical questions

1. What do C, A, and P mean in practical words?
2. What choice appears during a partition?
3. Why is "we are AP" a weak product claim?
4. How can two operations in one product use different rules?
5. How does a nearby stale replica trade latency and consistency?

#### Easy practical tasks

1. Write five sentences that explain CAP without extra letters.
2. Make a table: "Use case" and "Refuse or accept during partition". Add pay, browse catalog, and post a comment.
3. Draw two nodes that cannot talk. Write the last-seat problem.
4. List three user messages for "we refuse the write until replicas talk".

#### Medium practical tasks

1. Write rules for a library hold: last copy of a book, two campuses, link down.
2. Explain in eight sentences why a single-node database is not a CAP example in the partition sense.
3. Write an ADR: one primary for writes, replicas for reads, stale read allowed on search only.

#### Advanced practical tasks

1. Write a one-page CAP explainer for stakeholders who are not engineers. Use the last-seat story.
2. Read a public CAP FAQ. Write three misunderstandings that beginners have. Do not copy long quotations.

---

## Consensus (idea: Raft / Paxos) — you rarely implement it

Consensus is the problem of several nodes that must agree on one value or one log of commands, even if some nodes fail or messages reorder. Databases and coordinators use consensus to elect a leader and to replicate a log.

Paxos is a family of consensus algorithms. Raft is a consensus algorithm that many people find easier to learn. Both are hard to implement correctly. Bugs are subtle. You rarely implement them in an application.

What you need to know as an architect:

- A consensus group has a majority (quorum). If you lose too many nodes, the group stops accepting writes.
- A leader can change. Clients can see a short interruption.
- Consensus is not free. Extra round trips increase latency.
- Managed databases and brokers already include this machinery. You configure it. You do not rewrite it.

Do not start a new product with a homemade cluster protocol. Use a managed store or a well-tested library that the team can operate.

TLA+ and formal methods help people who design protocols. That work is optional and late (Topic 21). For this path, remember: agreement across nodes is a specialized problem.

If a vendor says "we use Raft", ask practical questions: What is the quorum? What happens when a zone is down? What is the failover time? Those questions are architecture. The Raft paper is optional reading.

### Questions

#### Theoretical questions

1. What is consensus in one sentence?
2. Why do you rarely implement Raft or Paxos in an application?
3. What is a quorum in practical words?
4. How does consensus increase latency?
5. What questions do you ask when a product uses Raft?

#### Easy practical tasks

1. Write four sentences that define consensus for a beginner.
2. Make a table: "Node count" and "Majority". Add 3 and 5.
3. List three products that already solve consensus for you (managed DB, etcd-like, broker).
4. Write five sentences on why a homemade cluster is a risk for a small team.

#### Medium practical tasks

1. Explain why two nodes cannot form a safe majority if one node is down (split-brain idea in simple words).
2. Write an ADR: use a managed PostgreSQL primary with replicas. Do not build a custom consensus layer.
3. Draw a three-node group. Mark one node down. Write if writes can continue (assume majority of three).

#### Advanced practical tasks

1. Read a public Raft overview (not the full paper unless you want). Write ten sentences: leader, log, majority. Do not copy long text.
2. Compare "single primary plus replicas" with "multi-primary consensus" at a high level. Write operations cost for a five-person team.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do the fallacies, partial failure, and CAP describe the same world from three views?
2. How do timeout, retry, backoff, and idempotent handlers form one outbound-call policy?
3. How do circuit breakers and bulkheads protect different resources?
4. Why does this path place distributed basics after monoliths and messaging?
5. What does a beginner implement, and what does a beginner configure?

#### Easy practical tasks

1. Write a one-page cheat sheet of all terms in this topic.
2. For one outbound HTTP call, write timeout, retry, idempotency, breaker, and bulkhead settings.
3. Write three ADR titles: no homemade consensus, retry policy, degrade when search is down.
4. Draw a "dependency down" poster for a shop with four dependencies and four behaviors.

#### Medium practical tasks

1. Write a two-page distributed-call standard for a team of four. Include the failure matrix.
2. Take the broker-down drill from Topic 10. Add timeouts, outbox growth, and user messages from this topic.
3. Role-play a request to run the API in three regions on week two. Write the CAP and operations reply.

#### Advanced practical tasks

1. Write a tabletop exercise (90 minutes) that walks a partition between API and database. Include breaker, bulkhead, and user text.
2. Map this topic to later topics in `architecture.topics.md` (12–15). Write which quality each later topic deepens.
