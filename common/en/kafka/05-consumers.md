# 5. Consumers

## Description

A consumer is a client that reads records from topic partitions. The common path is a group, a subscription, and a poll loop. The consumer commits offsets so that it can continue after a restart.

This topic covers the consumer API, subscribe versus assign, the poll loop, commit modes, `enable.auto.commit` risks, poll and session timeouts, rebalance types, and seek. Complete this topic after topic 4. Topic 6 covers consumer groups in more detail. Topic 7 covers delivery semantics.

Use one term for each concept. Poll is the call that fetches records and drives the consumer. A commit stores an offset for a group and a partition. Subscribe lets the group assign partitions. Assign sets partitions in the client. Use KRaft clusters. The consumer talks to brokers, not to ZooKeeper.

---

## Consumer API

The consumer API is the client interface that reads records. In Java the type is `KafkaConsumer`. Other languages have a consumer type with the same idea: configure, subscribe or assign, poll, commit, close.

A typical sequence:

1. Build a configuration. Set `bootstrap.servers`. Set `group.id` for the group path. Set deserializers for key and value.
2. Create the consumer.
3. Subscribe to topic names, or assign partition objects.
4. Loop: call `poll` with a timeout. Process the records. Commit when your policy requires it.
5. Call `close` when the process ends. `close` can trigger a group leave and a rebalance.

`group.id` identifies the consumer group. Two processes with the same `group.id` share partitions. A new `group.id` reads independently.

`auto.offset.reset` applies when the group has no committed offset for a partition. `earliest` starts at the oldest retained record. `latest` starts at the log end. `none` fails if no offset exists.

Do not share one consumer object across threads unless the client documents that use. The common Java consumer is not thread-safe.

### Questions

#### Theoretical questions

1. What is the consumer API?
2. What does `group.id` identify?
3. When does `auto.offset.reset` apply?
4. What is the difference between `earliest` and `latest`?
5. Why must you close a consumer?

#### Easy practical tasks

1. Write five sentences that describe the consumer API. Use only facts from this section.
2. Consume a topic from the beginning with the console consumer and a group id. Write the command.
3. Make a table: Configuration key, meaning. Add `bootstrap.servers`, `group.id`, `auto.offset.reset`.
4. Find the consumer class name in your language client. Write the module path.

#### Medium practical tasks

1. Write a small program that subscribes, polls once or twice, prints each key and value, and closes.
2. Run two programs with different `group.id` values on the same topic. Confirm that both receive the records.
3. Read the official consumer API page. List five configuration keys that you must understand.

#### Advanced practical tasks

1. Write a consumer that logs `auto.offset.reset` behavior: new group with `earliest` versus new group with `latest` on a topic that already has data.
2. Compare the Java poll-and-commit loop with your language client. Write a one-page map of the same steps.

---

## Subscribe vs assign

`subscribe` registers topic names (or a pattern). The group coordinator and the assignor select partitions. When members join or leave, Kafka reassigns partitions. This is the common application path.

`assign` sets an explicit list of topic partitions on the consumer. The consumer does not join a group rebalance for that assignment. You own the list. If you run two processes that assign the same partition, both read the same records. There is no automatic share.

Use `subscribe` when workers must share a topic. Use `assign` for tools, replay jobs, and tests that must pin a partition.

Do not call `subscribe` and `assign` in a conflicting way on the same consumer. Pick one mode.

A subscription can use a topic pattern in some clients. Use a tight pattern. A wide pattern can subscribe to topics that you do not want.

### Questions

#### Theoretical questions

1. What does `subscribe` do?
2. What does `assign` do?
3. Who selects partitions after `subscribe`?
4. What happens if two processes `assign` the same partition?
5. When is `assign` the better choice?

#### Easy practical tasks

1. Write four sentences that compare `subscribe` and `assign`.
2. Make a table: Mode, group rebalance, who picks partitions.
3. Consume with the console consumer (subscribe path). Write that the console uses a group.
4. List three jobs for `subscribe` and three jobs for `assign`.

#### Medium practical tasks

1. Write a consumer that assigns partition 0 of a two-partition topic. Produce to both partitions. Confirm that you only see partition 0.
2. Write a consumer that subscribes to two topic names. Produce to both. Print the topic of each record.
3. Read your client docs for `subscribe` pattern support. Write one safe pattern and one unsafe pattern.

#### Advanced practical tasks

1. Run one `subscribe` group of two members and one `assign` tool on the same topic. Write how their offset storage differs.
2. Write a replay tool design: assign, seek to a time, read to a time, exit. Implement a small version.

---

## Poll loop

The poll loop is the heart of a consumer. Each `poll` call:

- joins the group and heartbeats as the client requires
- fetches records from assigned partitions
- returns a batch of records to your code

You must call `poll` often enough. If you do not, the group can consider the member dead. The next sections cover the timers.

Process records that `poll` returned before you call `poll` again. Do not start unbounded work in other threads and then poll without a bound. If processing is slow, you can exceed `max.poll.interval.ms`.

`max.poll.records` limits how many records one `poll` returns. A smaller value shortens each processing slice. A larger value increases throughput and the risk of a long slice.

Do not block forever inside the loop. Bound external calls. Topic 19 in the path repeats this rule as backpressure.

An empty poll is normal. The timeout on `poll` is the maximum wait for records, not an error.

### Questions

#### Theoretical questions

1. What three kinds of work does `poll` do at a high level?
2. Why must you call `poll` often?
3. What does `max.poll.records` limit?
4. Why is a long block inside the loop a problem?
5. Is an empty poll an error?

#### Easy practical tasks

1. Write five sentences about the poll loop. Use only facts from this section.
2. Draw a loop: poll → process → commit → poll.
3. Find `max.poll.records` in your client. Write the default.
4. List three tasks that are too slow to run inside an unbounded poll loop.

#### Medium practical tasks

1. Write a consumer that prints the record count of each poll. Produce 100 records. Write how many polls you see.
2. Set `max.poll.records=1`. Consume 10 records. Write how the loop looks.
3. Add a `sleep` inside the loop that is longer than `max.poll.interval.ms` on a test group. Record what happens (rebalance or revoke). Use a short interval on a local cluster.

#### Advanced practical tasks

1. Design a worker that processes records in a thread pool without breaking the poll timer. Write the rules (bound the pool, pause partitions, or both).
2. Measure poll interval and processing time for 10 000 records. Write p50 and p99 if you can. State whether you are inside the poll interval.

---

## Auto-commit vs manual commit

A commit writes the group offset for a partition to Kafka (the `__consumer_offsets` topic). After a restart, the consumer continues from the committed position.

**Auto-commit.** The client commits on a timer when `enable.auto.commit` is true. The timer is `auto.commit.interval.ms`. The client commits offsets that the consumer has returned from `poll`, not offsets that your business logic has finished, unless you finish the work before the next commit point.

**Manual commit.** Your code calls `commitSync` or `commitAsync` (Java names) after you finish the work. Synchronous commit waits for the broker. Asynchronous commit does not wait. Combine them with a clear policy.

Manual commit gives you control. You can commit after a database write. You can commit each record or each batch.

At-least-once processing is the common result: if you crash after the work and before the commit, you process again. Topic 7 covers that case.

Do not mix policies in one team without a written rule.

### Questions

#### Theoretical questions

1. What does a commit store?
2. When does auto-commit write offsets?
3. What is the difference between `commitSync` and `commitAsync`?
4. Why does a crash before commit cause a repeat?
5. Why do teams choose manual commit for important work?

#### Easy practical tasks

1. Make a table: Auto-commit, manual sync, manual async. Add one row for control and one row for typical risk.
2. Write four sentences about a commit. Use only facts from this section.
3. Find the default of `enable.auto.commit` in your client version.
4. Draw a time line: poll, process, crash, restart. Mark a missing commit.

#### Medium practical tasks

1. Write a consumer with manual `commitSync` after each poll batch. Restart it. Confirm that it continues after the last commit.
2. Describe a consumer group with `kafka-consumer-groups`. Write the committed offset and the lag.
3. Read official commit API docs. Write when `commitAsync` can lose a commit relative to a later `commitSync`.

#### Advanced practical tasks

1. Implement batch commit every N records. Document N. Show lag in the consumer-groups tool.
2. Write a one-page policy: when auto-commit is allowed (low-value metrics) and when it is forbidden (payments).

---

## `enable.auto.commit` pitfalls

`enable.auto.commit=true` is convenient. It is also easy to get wrong.

Pitfall 1: you process a record after the client has already committed it. A crash then skips the record. This is at-most-once. It happens when you still work on records while the auto-commit timer fires, or when you pass records to another thread.

Pitfall 2: you crash after you process but before the next auto-commit. You process again. This is at-least-once. Duplicates appear.

Pitfall 3: you think auto-commit means "commit when I finish my function". It does not. It means "commit on the client timer and poll cycle".

Pitfall 4: a rebalance can occur while auto-commit and processing overlap. You can commit or skip in surprising ways if other threads hold records.

For learning, turn auto-commit off. Commit in your code after you finish the work. If you keep auto-commit on, finish all work on the polled records before you call `poll` again, and do not hand records to other threads.

### Questions

#### Theoretical questions

1. How can auto-commit skip a record?
2. How can auto-commit duplicate a record?
3. Does auto-commit wait for your business function to finish?
4. Why are extra threads a risk with auto-commit?
5. What is the safer learning default?

#### Easy practical tasks

1. Write five sentences about auto-commit pitfalls. Use only facts from this section.
2. Make a table: Pitfall, skip or duplicate.
3. Find `auto.commit.interval.ms` in the docs. Write the default.
4. List three application types that must not use auto-commit.

#### Medium practical tasks

1. Write a consumer with auto-commit on and a long sleep after print. Kill it during the sleep. Restart. Write whether you see duplicates or skips.
2. Write the same consumer with auto-commit off and a commit after print. Repeat the kill test. Compare.
3. Read an official warning about auto-commit. Rewrite it in STE.

#### Advanced practical tasks

1. Demonstrate a skip: process in another thread, auto-commit on, crash after commit and before process. Write the evidence.
2. Write a code review checklist (eight items) that rejects unsafe auto-commit.

---

## `max.poll.interval.ms` and session timeout

The consumer has two important timers.

`session.timeout.ms` (with heartbeat interval) controls liveness of the process in the group. If heartbeats stop, the coordinator removes the member. Heartbeats can run in a background thread in current Java clients.

`max.poll.interval.ms` controls how long the group allows between two `poll` calls. If your processing takes longer than this interval, the coordinator removes the member even if heartbeats still run. The partition is assigned to another member. Your old member can then commit or process records that it no longer owns. That is a serious bug.

Set `max.poll.interval.ms` above your worst processing time for one poll batch. Or reduce `max.poll.records` so that each batch is smaller.

`heartbeat.interval.ms` must be lower than `session.timeout.ms`. The official docs give a ratio. Do not set random values.

These timers are not the same as `delivery.timeout.ms` on the producer.

### Questions

#### Theoretical questions

1. What does `session.timeout.ms` detect?
2. What does `max.poll.interval.ms` detect?
3. Why can heartbeats continue while the poll interval expires?
4. What can go wrong if a revoked member still processes records?
5. How does `max.poll.records` help a slow processor?

#### Easy practical tasks

1. Make a table: Timer, what it watches, typical failure.
2. Find the defaults for the three keys in this section in your client.
3. Write four sentences that separate session timeout from poll interval.
4. Draw a time line: poll, long process, poll interval expires, rebalance.

#### Medium practical tasks

1. Set `max.poll.interval.ms` to 5000 and sleep 8000 ms in the loop. Record the group describe output and the logs.
2. Lower `max.poll.records` so that the same work fits in the interval. Confirm that the member stays.
3. Read official docs for the heartbeat to session timeout ratio. Write the rule.

#### Advanced practical tasks

1. Write a one-page tuning note for a consumer that calls a slow HTTP API. Include pause/resume or a smaller poll batch.
2. Capture coordinator logs or client logs for a poll-interval kick. Annotate five lines.

---

## Rebalance: eager vs cooperative

A rebalance is a change of partition assignment in a group. Members join, leave, or fail. The assignor computes a new assignment.

**Eager rebalance** (old default protocol). All members revoke all partitions first. Consumption stops. Then members receive a new assignment. This is a stop-the-world pause. Topic 6 covers the group view of that pause.

**Cooperative rebalance** (incremental). Members revoke only the partitions that they must give away. They keep the rest. Consumption continues on the partitions that do not move.

Use a cooperative sticky assignor for new applications when the client supports it. Example Java name: `CooperativeStickyAssignor`. Set `partition.assignment.strategy` accordingly.

During a rebalance, commit offsets for partitions that you revoke. Current clients have revoke callbacks. Commit in the revoke callback when you use manual commit.

Do not take a long lock in a rebalance callback. You can exceed timers and cause another rebalance.

### Questions

#### Theoretical questions

1. What is a rebalance?
2. What does an eager rebalance revoke?
3. What does a cooperative rebalance revoke?
4. Why is cooperative rebalance better for pause time?
5. What must you commit in a revoke callback?

#### Easy practical tasks

1. Write five sentences about rebalance types. Use only facts from this section.
2. Make a table: Eager vs cooperative. Add rows for revoke scope and consume pause.
3. Find the default assignor in your client version.
4. List three events that trigger a rebalance.

#### Medium practical tasks

1. Start two consumers in one group. Add a third. Watch logs for revoke and assign. Write whether all partitions revoked.
2. Set the cooperative sticky assignor. Repeat the join. Compare the revoke set.
3. Read official incremental cooperative rebalancing notes. Write three facts in STE.

#### Advanced practical tasks

1. Measure the consume pause (time with no records) during a join for eager versus cooperative on a topic with many partitions.
2. Implement a revoke callback that commits sync. Kill a member. Show that the new member does not duplicate the last committed batch.

---

## Seek, beginning, end

Seek moves the consumer position in a partition. The next `poll` reads from that offset.

- Seek to an offset. You must assign the partition (or wait until it is assigned). Some clients require a poll before seek.
- Seek to the beginning. The next read is the earliest retained offset.
- Seek to the end. The next read is the log end. You skip existing records.

Seek does not always commit. If you seek and then crash before a commit, the group can restart at the old committed offset. If you want the seek to survive a restart, commit after the seek (or use the admin reset tool).

The consumer-groups tool can reset offsets for a group that is idle. Use that tool for operations. Use seek in a running consumer for replay in process.

Do not seek on a partition that you do not own in a subscribe group. The assignment can change.

`--from-beginning` on the console consumer is a seek-to-beginning for that start, often with a new group.

### Questions

#### Theoretical questions

1. What does seek do?
2. What is seek to beginning?
3. What is seek to end?
4. Why does a seek not always survive a crash?
5. When do you use the consumer-groups reset tool instead of seek?

#### Easy practical tasks

1. Write four sentences about seek. Use only facts from this section.
2. Consume with `--from-beginning`. Then start a new group without that flag. Write the difference.
3. Make a table: Operation, next record.
4. Find `seek`, `seekToBeginning`, and `seekToEnd` (or equivalents) in your client.

#### Medium practical tasks

1. Write a program that assigns a partition, seeks to offset 1, and prints the next record.
2. Seek to end, produce two new records, and confirm that you only see those two.
3. Reset a stopped group with `kafka-consumer-groups` to earliest. Start the consumer. Write the commands.

#### Advanced practical tasks

1. Seek to a timestamp if the API supports it (`offsetsForTimes`). Replay one minute of a topic. Document the API calls.
2. Write a one-page warning: seek plus auto-commit can undo the seek. Give a sequence.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe the safe manual-commit loop from subscribe to close in ordered steps.
2. How do subscribe, poll timers, and cooperative rebalance work together in a live group?
3. Why is assign the wrong default for a horizontally scaled worker service?
4. What is the difference between "position in this process" and "committed offset for the group"?
5. A teammate sets a huge `max.poll.interval.ms` to hide slow processing. What problems remain?

#### Easy practical tasks

1. Write a consumer with auto-commit off, subscribe, poll, print, `commitSync`, close. Run it on `demo`.
2. Write a one-page cheat sheet: subscribe vs assign, auto vs manual commit, two timers, seek, cooperative assignor.
3. Describe your consumer group. Highlight assignment, committed offset, and lag.
4. Draw eager revoke-all versus cooperative revoke-some.

#### Medium practical tasks

1. Write a script that starts a consumer, produces 20 records, and prints lag until lag is 0.
2. Demonstrate seek to beginning and a group reset with the admin tool as two different methods. Write when you pick each method.
3. Map each subsection to one official documentation heading.

#### Advanced practical tasks

1. Build a consumer that pauses a partition when a downstream API fails, retries, then resumes. Do not exceed the poll interval.
2. Break a consumer with a slow poll. Watch the rebalance. Fix it with a smaller poll batch. Record before and after group describe output.
