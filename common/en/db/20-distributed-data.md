# 20. Distributed Data

## Description

This topic shows how you split data across nodes and how you keep writes safe when more than one service is involved. You learn sharding, consistent hashing, why distributed transactions are expensive, saga and outbox patterns, eventual consistency, idempotent writes, and clock problems.

Use one term for each concept. A shard is a subset of the data on one node (or one pair). A saga is a sequence of local transactions with compensating actions. Complete this topic after you understand replication, transactions, and CAP in practical words. Distribution adds failure modes. It does not remove the need for a system of record.

This path is vendor-neutral. Pattern names are standard. Product syntax differs.

---

## Sharding / partitioning

Sharding (horizontal partitioning) splits the rows of a dataset by a shard key. Each shard holds a disjoint subset of keys. A router or the application sends a request to the shard that owns the key.

```text
shard key = customer_id
customers 1..999    --> shard A
customers 1000..1999 --> shard B
```

Use sharding when one primary cannot hold the data or the write rate. Sharding is an operational and design cost. Do not shard a small database.

A good shard key:

- appears in almost every query
- spreads writes (no single hot key)
- keeps related rows that you join on the same shard when you can

A poor shard key creates a hot shard (all new orders go to "today" if you shard only on date) or forces cross-shard joins.

Cross-shard queries must fan out. They are slower and harder to isolate. Transactions that touch two shards are distributed transactions (later section) or they are application sagas.

Vertical partitioning (split columns into more tables or more stores) is a different idea. This section means row split by key.

Range shards are simple and can unbalance. Hash shards spread better and make range scans harder. A directory (lookup table of key to shard) is flexible and is one more system to run.

Resharding (move keys to new nodes) is a planned project. Design for more shards than you need on day one, or use a scheme that can split (consistent hashing, next section).

Do not shard "for high availability." Replication gives copies. Sharding gives capacity. You often need both.

### Questions

#### Theoretical questions

1. What does a shard key decide?
2. When do you shard?
3. What makes a shard key poor?
4. Why are cross-shard joins expensive?
5. Why is sharding not a substitute for replication?

#### Easy practical tasks

1. Pick a shard key for `orders` in a shop. Write one reason.
2. Draw two shards and a router.
3. Write one query that stays on one shard and one that fans out.
4. Write one hot-key example (a celebrity user id).

#### Medium practical tasks

1. Design range shards versus hash shards for `customer_id`. Write a range report that hash shards make hard.
2. Estimate: 10 million customers, 4 shards. Write how you assign ids. Show a skew if ids are sequential and you use range shards badly.
3. Write a reshard step list at a high level: add node, move keys, update router, verify counts.

#### Advanced practical tasks

1. Write a one-page shard plan for a multi-tenant app (tenant_id as key). Include a large tenant that needs its own shard.
2. Table-top a cross-shard order: customer on A, inventory on B. Write why a single local transaction is not enough.

---

## Consistent hashing (idea)

Consistent hashing places keys on a ring. Nodes occupy positions on the ring. A key hashes to a point. The key belongs to the next node along the ring (or to a set of replicas on the ring).

When you add or remove a node, only the keys that map to the changed arc move. In a simple modulo scheme (`hash(key) % N`), almost all keys move when N changes.

```text
Ring:  0 ---- N1 ---- N2 ---- N3 ---- 0
Key k hashes to a point; owner is the next node clockwise
Add N4: only keys in N4's new arc change owner
```

Virtual nodes (many positions per physical node) spread load. A large machine gets more virtual nodes.

Clients or a proxy compute the hash. They must use the same hash function and the same ring membership. A membership change must propagate. During the change, some keys can be on the old node and the new node. You need a move protocol: copy, switch, delete.

Consistent hashing is an idea. Products add buckets, vnodes, and token ranges. The goal is the same: cheap cluster resize and even spread.

Do not treat the hash as a secret. Do not put a low-cardinality field (boolean) as the only hash input. You get two piles, not a ring of many keys.

Do not assume that consistent hashing gives transactions. It only answers "which node owns this key."

### Questions

#### Theoretical questions

1. What problem does consistent hashing reduce when N changes?
2. What does modulo N do to keys when N grows by one?
3. What is a virtual node?
4. What must all clients share to find the same owner?
5. Does consistent hashing provide a transaction? Explain.

#### Easy practical tasks

1. Draw a ring with three nodes and five keys.
2. Add a fourth node. Mark which keys move.
3. Write why a boolean shard key defeats the ring.
4. Write modulo versus ring in two sentences.

#### Medium practical tasks

1. Simulate 20 keys and 3 nodes on paper with a simple hash (example: `id % 360` as an angle). Assign owners. Add a node. Count moves.
2. Read a product token-range or vnode page. Map the terms to this section.
3. Write a move protocol: copy keys, dual-write, switch reads, stop old.

#### Advanced practical tasks

1. Write a one-page note: consistent hashing, vnodes, and a membership-change incident (some clients on the old ring).
2. Compare consistent hashing with a central directory. Write three benefits of each.

---

## Distributed transactions and why they are expensive

A distributed transaction commits changes on more than one node (or more than one product) as one atomic unit. The usual protocol is two-phase commit (2PC): a coordinator asks participants to prepare, then it asks them to commit.

```text
1. Coordinator: PREPARE on A and B
2. A and B: flush enough state to commit or abort later; answer yes/no
3. Coordinator: COMMIT or ABORT on A and B
```

Costs:

- extra round trips
- locks or prepared state held until the coordinator finishes
- a blocking period if the coordinator crashes after prepare
- all participants must implement the protocol
- latency is the slowest participant plus the network

If A commits and B aborts without a protocol, you have a partial write. Business data is wrong.

Many teams avoid 2PC across products. They use a single local transaction in one system of record, then they propagate (saga, outbox).

Some products offer distributed transactions inside one cluster. They still cost more than a single-shard transaction. Prefer a shard key that keeps one business event on one shard.

Do not start a 2PC for a log line and a cache invalidate. Use a local transaction plus a retryable side effect.

Do not hold user-facing connections open while a coordinator waits on a slow second store.

### Questions

#### Theoretical questions

1. What does a distributed transaction try to guarantee?
2. What are the two phases in 2PC?
3. Why can a crashed coordinator block participants?
4. Why is a single-shard local transaction cheaper?
5. When is 2PC a poor fit?

#### Easy practical tasks

1. Number the 2PC steps.
2. Draw A prepared, B prepared, coordinator down.
3. Write one business event that should stay on one shard.
4. Write one side effect that should not join 2PC (example: send email).

#### Medium practical tasks

1. Write a sequence where A commits and B fails without 2PC. Name the repair.
2. Read whether your DBMS supports 2PC or XA. Write the name and one warning from the docs.
3. Compare latency: one local `COMMIT` versus two prepares and two commits on a 20 ms link. Add the times on paper.

#### Advanced practical tasks

1. Write a one-page decision: reject 2PC between DBMS and a search index. Propose outbox instead.
2. Table-top coordinator recovery: where the decision log lives, and what a participant does if it never gets phase 2.

---

## Saga and outbox patterns

A saga is a sequence of local transactions. Each step has a compensating action that undoes the business effect if a later step fails.

```text
1. Reserve inventory   (compensate: release)
2. Charge payment      (compensate: refund)
3. Create shipment     (compensate: cancel shipment)
```

The saga is not ACID isolation across steps. Other readers can see a reserved item before the charge finishes. You design status fields (`reserved`, `paid`, `failed`).

An **orchestration** saga has a coordinator process that calls each step. A **choreography** saga lets each service react to events. Orchestration is easier to see. Choreography avoids one boss process and is harder to debug.

The **outbox** pattern writes the business row and an "event to send" row in the same local transaction. A publisher reads the outbox and sends the event to a bus or to another store. Then it marks the outbox row as sent.

```text
BEGIN
  UPDATE orders SET status = 'paid'
  INSERT INTO outbox (event) VALUES ('OrderPaid')
COMMIT
-- publisher sends OrderPaid; at-least-once
```

Without an outbox, you can commit the order and crash before you send the event. The other service never hears. Dual-write without a transaction is the common bug.

Compensations must be safe to retry. A refund that runs twice is a money bug unless the refund is idempotent (next sections).

Do not use a saga to hide a missing local constraint. Do not use an outbox as a chat log for humans. It is a delivery mechanism.

### Questions

#### Theoretical questions

1. What does a compensating action do?
2. Why is a saga not one ACID transaction?
3. What two writes does an outbox put in one local transaction?
4. What dual-write bug does the outbox prevent?
5. What is the difference between orchestration and choreography?

#### Easy practical tasks

1. Write a three-step saga for a bookstore order. Name each compensate.
2. Draw order table plus outbox table plus publisher.
3. Write one status value that other services can see in the middle of a saga.
4. Label a flow as orchestration or choreography.

#### Medium practical tasks

1. Implement an outbox table in a practice schema. Insert an order and an event in one transaction. Write a publisher loop in pseudo-code.
2. Write a compensate that fails (payment refund API down). Write how you retry and what you show the user.
3. Compare "send HTTP in the same request after COMMIT" with outbox. Write two failure cases.

#### Advanced practical tasks

1. Write a one-page saga for checkout with timeouts and a manual repair queue.
2. Design outbox cleanup: sent rows, retention, and poison events. Include idempotent consumers.

---

## Eventual consistency

Eventual consistency means copies can differ for a time. If writes stop, the copies converge to the same values (under the product's rules). Readers can see stale data during the delay.

```text
Write on node A at t=0
Read on node B at t=50ms  -->  old value
Read on node B at t=2s    -->  new value (if replication caught up)
```

This model appears in:

- asynchronous replicas
- search indexes after a write
- caches with TTL
- saga steps that have not finished
- DNS and many object stores

Use eventual consistency when the business can accept a short stale read. Example: a view count, a search result, a recommendation. Do not use it as the only story for a bank balance unless you have a precise conflict rule and a ledger.

Give users a truth they can understand: "Your change is saved. Search can take a minute." A silent stale read looks like a lost write.

**Read your writes** is a stronger need. After a save, read the primary or wait for a version. **Monotonic reads** mean you do not go backward in time in one session.

Do not say "eventual consistency" to mean "we do not know." Write the delay budget and the repair if convergence fails (a stuck replica).

Causal stories matter: if comment B replies to comment A, a reader should not see B without A. That requirement is more than "eventually both exist."

### Questions

#### Theoretical questions

1. What does eventual consistency promise if writes stop?
2. Name four places this model appears.
3. Why can a stale read look like a lost write?
4. What does read-your-writes require after a save?
5. Why is a reply-without-parent a causal problem?

#### Easy practical tasks

1. Write two features that can be eventually consistent and two that cannot for a shop.
2. Write a user-facing sentence for a delayed search index.
3. Draw a stale read on a replica.
4. Define a delay budget (example: 5 seconds) for cart count versus for payment status.

#### Medium practical tasks

1. Demo a write on a primary and a read on an async replica (or simulate with a delay). Record the stale window.
2. Design a "saved" page that reads the primary while the list page reads a replica.
3. Write a monitor: if replica lag exceeds the budget, stop sending reads to the replica.

#### Advanced practical tasks

1. Write a one-page consistency menu: strong local txn, read-your-writes, eventual search, saga-visible states.
2. Table-top a stuck replica that never converges. Write detection, user impact, and repair.

---

## Idempotent writes

An idempotent write can run more than once and leave the same final state as one successful run. Networks retry. Sagas retry. Outbox consumers retry. Without idempotency, a retry charges a card twice.

```text
POST /payments { idempotency_key: "ord-1001-pay" }
-- first call: charges
-- retry: returns the same payment id, no second charge
```

Methods:

1. A client idempotency key stored with the result.
2. A natural unique key (`order_id` unique on `payments`).
3. A state machine: `pending` to `paid` only once.
4. Upsert with a version number.

At-least-once delivery (outbox, message bus) requires idempotent consumers. At-most-once delivery loses messages. Exactly-once is a combination of at-least-once plus idempotency (and careful side effects).

`UPDATE balances SET n = n + 10` is not idempotent. Two retries add 20. `UPDATE balances SET n = 90 WHERE account = 1 AND version = 3` can be idempotent if the version check fails on the second apply.

Do not rely on "the user will not double-click." The timeout retry is enough to double-apply.

Log the key. Reject a second different body for the same key (conflict).

### Questions

#### Theoretical questions

1. What does an idempotent write guarantee on retry?
2. Why do outbox consumers need idempotency?
3. Why is `n = n + 10` not idempotent?
4. What is exactly-once in practical terms in this section?
5. Why is "users will not double-click" a weak control?

#### Easy practical tasks

1. Add an `idempotency_key` column idea to a `payments` table. Write the unique constraint.
2. Write a retry story for "request timed out, user does not know if the charge ran."
3. Label three statements as idempotent or not: `INSERT` without key, `DELETE FROM t WHERE id = 1`, `n = n + 1`.
4. Write the conflict rule: same key, different body.

#### Medium practical tasks

1. Implement a unique payment per `order_id` in SQL. Run the insert twice. Record the second result.
2. Write a consumer that stores `event_id` as processed. Skip a duplicate. Prove with two deliveries.
3. Rewrite an increment as a versioned set. Show a double retry that does not double-add.

#### Advanced practical tasks

1. Write a one-page idempotency standard for HTTP writes and for queue consumers. Include key lifetime.
2. Design a refund that is idempotent and a compensate that uses it in a saga. Write the keys.

---

## Clock and ordering problems

Each machine has a clock. Clocks skew. NTP steps a clock. A virtual machine can pause. You cannot trust `NOW()` on two hosts as a global order.

```text
Host A thinks t=10:00:00.200  writes row
Host B thinks t=10:00:00.050  writes row
-- B's timestamp is earlier; B did not happen first
```

Problems:

- last-write-wins with wall clocks loses the later real write
- unique "created_at" order is wrong across nodes
- TTL computed from a future clock expires at once
- a certificate or token looks expired because the clock jumped

Better tools for order:

- a sequence or log position on one primary
- a logical clock or version vector for concurrent updates
- a causal id (parent event id) in the payload
- Hybrid logical clocks in some products (still not magic)

Do not sort global events only by `timestamptz` from many writers. Use a single sequencer if you need a total order.

Do not measure lag as `replica_now - primary_now` without care. Prefer log-sequence difference.

In tests, fake time. In production, alert on clock skew between nodes.

For humans, show local time. For order, store UTC plus a monotonic source when the product gives one.

### Questions

#### Theoretical questions

1. Why can two host clocks disagree on order?
2. What does last-write-wins lose when clocks skew?
3. What is a safer order source than wall time on many hosts?
4. Why is a log-sequence lag measure better than a clock difference?
5. Why can a clock jump break TTL?

#### Easy practical tasks

1. Write one story of a skewed last-write-wins on a document `updated_at`.
2. Name two better order sources from this section.
3. Draw two hosts and two timestamps that invert real order.
4. Write why the handbook still stores instants in UTC for display.

#### Medium practical tasks

1. Compare `ORDER BY created_at` on rows from two services. Write a case that misorders.
2. Read how your DBMS exposes a WAL or LSN number. Write how a replica lag can use it.
3. Write a version-vector idea for two replicas that update the same key (two counters). Keep it short.

#### Advanced practical tasks

1. Write a one-page ordering rule: local txn order, cross-service causal ids, ban on LWW wall clocks for money.
2. Table-top a 5-second NTP step on a primary. List features that break (TTL, tokens, LWW). Write detections.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do a shard key, consistent hashing, and a local transaction work together to avoid 2PC?
2. When do you choose a saga and an outbox instead of a distributed commit?
3. How do eventual consistency, idempotent consumers, and clock skew appear in one checkout retry?
4. Which order source do you use inside a shard versus across services?
5. A teammate shards on `created_at` date, uses 2PC to search, and sorts merges by wall clock. Which facts do you use in the reply?

#### Easy practical tasks

1. Write a cheat sheet: shard, hash ring, 2PC cost, saga, outbox, eventual, idempotent, clocks.
2. Pick a shard key and an idempotency key for `PlaceOrder`.
3. Draw outbox plus a retrying consumer with a processed-event table.
4. Write one user sentence for a stale derived store.

#### Medium practical tasks

1. Design a two-shard shop that still charges once (idempotency) and updates search later (outbox). Write the tables.
2. Table-top a partition during a saga. Write what the user sees and which compensate runs.
3. Compare hash resize with directory resize in a short table.

#### Advanced practical tasks

1. Write an end-to-end design: shard map, one-shard checkout transaction, outbox, saga compensate, idempotent pay, no wall-clock LWW.
2. Review a public distributed-system post. Replace slogans with the terms from this topic. List missing repair paths.
