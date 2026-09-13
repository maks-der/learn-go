# 8. Distributed Systems

## Description

A distributed system is a system that uses more than one process that communicate over a network. This topic covers the fallacies of distributed computing, partial failure, retries, backoff, circuit breakers, bulkheads, CAP in practical words, and consensus.

Distribution is not a goal. It is a cost that you pay when one process is not enough (Topics 3 and 9). Every network hop adds failure modes. Complete Topics 1 to 7 before this topic. Pair with `net.topics.md` if sockets and latency are new.

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

## Partial failure, retries, backoff

A partial failure is a state where some parts of the system work and other parts do not. In one process, a crash often stops all work. In a distributed system, the API can be up while the mail worker is down. The database can accept writes while the cache is down. One replica can be stale.

Partial failure is the normal failure mode. You must define behavior for each dependency:

- Fail the user request.
- Degrade (show cached data, hide a widget).
- Queue work for later.
- Serve a subset of features.

If you do not define the behavior, the user sees random errors and long waits.

Detection is hard. A process that does not reply can be dead, slow, or partitioned from you. You cannot always know. Timeouts turn "unknown" into a local decision (Topic 2). That decision can be wrong: the work may still complete.

A timeout bounds the wait. A retry repeats the call after a failure. Exponential backoff increases the wait between retries (for example 100 ms, 200 ms, 400 ms, plus a random jitter).

Retries help when the failure is transient: a short network drop, a brief overload, a rolling deploy. Retries harm when the failure is permanent (bad request) or when all clients retry at once (retry storm).

Rules:

- Retry only idempotent operations, or operations that you protected with an idempotency key (Topics 2 and 5).
- Bound the retry count.
- Use backoff and jitter so that clients do not align.
- Do not retry `400` validation errors.
- Keep the total time inside the user budget.
- Apply a timeout on each attempt.

A retry of a non-idempotent `POST` can charge a card twice. That defect is worse than a user retry.

Jitter is a random extra delay. Without jitter, many workers wake at the same time and hit the dependency together.

Server-side retries plus client-side retries multiply load. Write who retries. Prefer one layer.

Health checks help, but a process can pass `/health` and still fail on the real query. Check the dependency that matters, with a bound.

Metrics: retry count, timeout count, and success after retry. If success after retry is rare, fix the dependency. Do not hide a chronic fault with retries.

### Questions

#### Theoretical questions

1. What is a partial failure?
2. Why is "no reply" not the same as "did not run"?
3. What is exponential backoff?
4. What is jitter and why do you add it?
5. Which errors must you not retry?

#### Easy practical tasks

1. Write four sentences that define partial failure.
2. Make a matrix: four dependencies (DB, cache, mail, search) and the user-visible behavior if each one is down.
3. Write a retry policy: max 3 attempts, timeouts of 200 ms, backoff 100/200/400 ms, jitter 0–50 ms.
4. List five user messages that are honest ("search is delayed") and two that hide the truth.

#### Medium practical tasks

1. Design degrade rules for a shop home page when search is down and when the catalog database is down. They must differ.
2. Write a failure matrix for place-order: payments down versus email down.
3. Explain in eight sentences why client retries plus server retries multiply load.

#### Advanced practical tasks

1. Write a one-page partial-failure catalog for six dependencies. Include detect method and max wait.
2. Design a game-day script (paper) that turns off one dependency. List expected signals. Do not attack systems that you do not own.

---

## Circuit breaker and bulkhead

A circuit breaker is a guard that stops calls to a failing dependency for a time. After a threshold of failures, the breaker opens. Calls fail fast. After a cool-down, a small number of trial calls run. If they succeed, the breaker closes. If they fail, the breaker stays open.

The breaker protects the caller. The caller does not fill all threads on a dead dependency. The user gets a fast error or a degrade path instead of a long wait.

A bulkhead is isolation of resources. One dependency gets its own thread pool, connection pool, or process. A storm on search must not take all connections that the catalog database needs.

Rules:

- A circuit breaker needs a timeout. Without a timeout, the breaker sees hangs, not failures.
- Fail fast is not a retry. Combine with a degrade rule from the previous section.
- Tune thresholds with measures. A breaker that opens on one blip is noise. A breaker that never opens is decoration.
- Bulkheads cost capacity. You hold unused slots so that one fault stays local.

A circuit breaker does not repair the dependency. It buys time and protects the rest of the system. Operators still need alerts (Topic 11).

Do not put a breaker on a local function call in one process unless you have a clear isolation goal. Breakers belong on remote or unstable dependencies.

A process or a container is a coarse bulkhead. A modular monolith can still use pools per dependency. Topic 9 adds more processes. Each process is a bulkhead and a new failure domain.

### Questions

#### Theoretical questions

1. What is a circuit breaker?
2. What happens when the breaker is open?
3. What is a bulkhead?
4. Why does a breaker need a timeout?
5. Why does a bulkhead cost capacity?

#### Easy practical tasks

1. Write five sentences that define circuit breaker and bulkhead.
2. Draw closed, open, and half-open states of a breaker.
3. Make a table: "Dependency" and "Own pool? (yes/no)". Add database, search, and mail.
4. List four user-visible results of fail-fast versus a long hang.

#### Medium practical tasks

1. Design a breaker for a payment HTTP call: threshold, cool-down, and degrade behavior.
2. Write how a search storm can exhaust a shared connection pool. Then draw a bulkhead.
3. Explain in eight sentences when a breaker hides a chronic fault (link to retry metrics).

#### Advanced practical tasks

1. Write a one-page standard: where breakers live, who owns thresholds, and how you test them.
2. Compare thread-pool bulkheads with separate processes as two isolation grades. Write operations cost.

---

## CAP in practical words

CAP is a teaching model for a distributed data store during a network partition. The usual letters are:

- C: consistency — readers see a single up-to-date value (a strong reading).
- A: availability — every working node returns a response.
- P: partition tolerance — the system continues even when nodes cannot talk.

On a real network, partitions happen. You cannot ignore P. During a partition you choose a side: refuse some writes or reads to keep one value (lean toward C), or accept answers that can disagree (lean toward A).

CAP is not a product checkbox. A database marketing page that says "CA" or "AP" is not a design. You must name the operation, the replica set, and the user-visible behavior when a replica is down.

Practical words:

- If the primary is unreachable, does the app fail the write or write to another node?
- If two replicas disagree, which value wins?
- How long may a reader see an old value (Topic 4)?

Most student systems use one primary database. That is a simple choice: writes go to one place. You still have a partition between the app and the database. That is a timeout and a fail, not a CAP exam question.

When you add replicas or more than one writer, write the behavior in an ADR. Topic 6 already said the source of truth must win. CAP language helps you state what happens when the network splits the copies.

Do not use CAP to reject a useful replica. Use it to write the degrade rule.

### Questions

#### Theoretical questions

1. What do C, A, and P stand for in this handbook?
2. Why can you not ignore partitions on a real network?
3. What do you choose during a partition?
4. Why is a marketing "AP" label not a design?
5. What is the usual student-system choice for writes?

#### Easy practical tasks

1. Write five sentences that explain CAP in practical words.
2. Make a table: "Event" and "Fail write or accept a copy?". Add three events.
3. List four questions you ask a vendor about replica behavior.
4. Draw two nodes that cannot talk. Write what a user of each node sees in a C-leaning design and in an A-leaning design.

#### Medium practical tasks

1. Write an ADR for a campus catalog: one primary, fail writes if the primary is down.
2. Explain in eight sentences how a read replica lag is a consistency choice, not only a scale choice.
3. Compare "fail the checkout" with "accept an order on a local node" during a split. Write the repair problem.

#### Advanced practical tasks

1. Write a one-page note that maps CAP language to your source of truth and one copy (cache or replica).
2. Read a public database document on failover. Write ten sentences in your own words. Do not copy long passages.

---

## Consensus (Raft / Paxos) — you rarely implement it

Consensus is the problem of getting a set of nodes to agree on a value or on a log of values even when some nodes fail. Leader election and a replicated log often use a consensus algorithm.

Paxos is a family of consensus algorithms. Raft is a consensus algorithm that many people find easier to study. Both are hard to implement correctly. Small bugs cause split brains or lost writes.

You rarely implement Raft or Paxos yourself. You use a product that already embeds a proven implementation: a consensus-backed coordination service, a database failover mechanism, or a broker controller. Your job is to operate the product and to understand the failure story.

When you need consensus:

- A cluster must elect one leader.
- Replicas must agree on the order of writes.
- A configuration change must be safe.

When you do not need to write Raft:

- A student monolith with one primary database.
- Application business rules (use a transaction in one store).
- A "we might need a cluster" slide.

If you study Raft, read a careful source and draw the leader, followers, and the log. Do not ship a homework Raft as the source of truth for money.

Topic 3 warned against over-engineering. A homemade consensus layer is over-engineering for almost every product in this path.

### Questions

#### Theoretical questions

1. What is consensus in this handbook?
2. What is Raft at a high level?
3. Why is a correct implementation hard?
4. When do products use consensus?
5. Why must you rarely implement Raft yourself?

#### Easy practical tasks

1. Write five sentences about consensus. Use only facts from this section.
2. Make a table: "Need" and "Use a product or write Raft?". Add six needs.
3. List four defects of a buggy leader election (two leaders, lost write, and similar).
4. Write four sentences to a teammate who wants to "just write Paxos" for a shop.

#### Medium practical tasks

1. Write an ADR: use the managed database failover, do not implement Raft.
2. Draw a leader and two followers. Write what a client does if the leader is down (product-level, not your algorithm).
3. Explain in eight sentences why business invariants belong in an aggregate, not in a homemade consensus module.

#### Advanced practical tasks

1. Read a public Raft overview (paper abstract or a teaching page). Write a one-page map of terms in your own words. Do not copy long passages.
2. Compare "one primary database" with "a consensus-backed store" for a five-person team. Write operations cost.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do the eight fallacies explain the need for timeouts, retries, and breakers together?
2. Why is partial failure a stronger idea than "the system is up or down"?
3. How do bulkheads and circuit breakers protect different things?
4. When is CAP language useful, and when is a single primary plus a timeout enough?
5. Why does this path tell you to use consensus products and not to write Raft for a campus app?

#### Easy practical tasks

1. Write a one-page cheat sheet: eight fallacies, partial failure, retry rules, breaker states, bulkhead, CAP in one paragraph, consensus warning.
2. For a to-do API that calls one database, write one timeout, no retry on `400`, and no breaker yet.
3. Draw three processes. Mark one down. Write what still works.
4. Bookmark `net.topics.md`. Write three network facts that this topic assumes.

#### Medium practical tasks

1. Write a short distributed-systems brief for a campus shop with API, database, and mail worker: failures, retries, and degrade.
2. Take a teammate client with no timeout and infinite retries. Rewrite the policy.
3. Write a twelve-week plan: weeks for one remote call with timeout, and the measure that would add a breaker.

#### Advanced practical tasks

1. Write a game-day plan (paper) for database down, broker down, and disk full. List signals. Do not run it on systems that you do not own.
2. Map this topic to later Topic 11 (SLOs and alerts) in a one-page table: failure mode, user-visible result, signal.
