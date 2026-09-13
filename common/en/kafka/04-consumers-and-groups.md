# 4. Consumers and Groups

## Description

A consumer is a client that reads records from topic partitions. The common path is a group, a subscription, and a poll loop. The consumer commits offsets so that it can continue after a restart.

This topic shows subscribe versus assign, the poll loop, commit modes, poll and session timeouts, consumer groups, rebalance, static membership, cooperative assignors, and seek. Complete this topic after topic 3. Topic 5 covers delivery semantics.

Use one term for each concept. Poll is the call that fetches records and drives the consumer. A commit stores an offset for a group and a partition. A consumer group is a set of consumers that share one `group.id`. A member is one consumer process in the group. Subscribe lets the group assign partitions. Assign sets partitions in the client. Use KRaft clusters. The consumer talks to brokers, not to ZooKeeper.

---

## Subscribe vs assign and the poll loop

The consumer API is the client interface that reads records. In Java the type is `KafkaConsumer`. Other languages have a consumer type with the same idea: configure, subscribe or assign, poll, commit, close.

A typical sequence:

1. Build a configuration. Set `bootstrap.servers`. Set `group.id` for the group path. Set deserializers for key and value.
2. Create the consumer.
3. Subscribe to topic names, or assign partition objects.
4. Loop: call `poll` with a timeout. Process the records. Commit when your policy requires it.
5. Call `close` when the process ends. `close` can trigger a group leave and a rebalance.

`subscribe` registers topic names (or a pattern). The group coordinator and the assignor select partitions. When members join or leave, Kafka reassigns partitions. This is the common application path.

`assign` sets an explicit list of topic partitions on the consumer. The consumer does not join a group rebalance for that assignment. You own the list. If you run two processes that assign the same partition, both read the same records. There is no automatic share.

Use `subscribe` when workers must share a topic. Use `assign` for tools, replay jobs, and tests that must pin a partition. Do not call `subscribe` and `assign` in a conflicting way on the same consumer. Pick one mode.

The **poll loop** is the main loop of a consumer. Each `poll` call joins the group and sends heartbeats as the client requires, fetches records from assigned partitions, and returns a batch of records to your code.

You must call `poll` often enough. If you do not, the group can consider the member dead. Process records that `poll` returned before you call `poll` again. `max.poll.records` limits how many records one `poll` returns.

An empty poll is normal. The timeout on `poll` is the maximum wait for records, not an error.

`auto.offset.reset` applies when the group has no committed offset for a partition. `earliest` starts at the oldest retained record. `latest` starts at the log end. `none` fails if no offset exists.

Do not share one consumer object across threads unless the client documents that use. The common Java consumer is not thread-safe. Topic 11 repeats: do not block the poll loop.

### Questions

#### Theoretical questions

1. What does `subscribe` do?
2. What does `assign` do?
3. What three kinds of work does `poll` do at a high level?
4. When does `auto.offset.reset` apply?
5. Why must you call `poll` often?

#### Easy practical tasks

1. Write five sentences that compare `subscribe`, `assign`, and the poll loop. Use only facts from this section.
2. Consume a topic from the beginning with the console consumer and a group id. Write the command.
3. Make a table: Mode, group rebalance, who picks partitions. Add subscribe and assign.
4. Find `max.poll.records` in your client. Write the default.

#### Medium practical tasks

1. Write a consumer that assigns partition 0 of a two-partition topic. Produce to both partitions. Confirm that you only see partition 0.
2. Write a consumer that subscribes, polls, and prints each key and value. Use two topic names if you can.
3. Set `max.poll.records=1`. Consume 10 records. Write how the loop looks.

#### Advanced practical tasks

1. Run one `subscribe` group of two members and one `assign` tool on the same topic. Write how their offset storage differs.
2. Write a replay tool design: assign, seek to a time, read to a time, exit. Implement a small version.

---

## Auto-commit vs manual commit

A commit writes the group offset for a partition to Kafka (the `__consumer_offsets` topic). After a restart, the consumer continues from the committed position.

**Auto-commit.** The client commits on a timer when `enable.auto.commit` is true. The timer is `auto.commit.interval.ms`. The client commits offsets that the consumer has returned from `poll`, not offsets that your business logic has finished, unless you finish the work before the next commit point.

**Manual commit.** Your code calls `commitSync` or `commitAsync` (Java names) after you finish the work. Synchronous commit waits for the broker. Asynchronous commit does not wait. Combine them with a clear policy.

Manual commit gives you control. You can commit after a database write. You can commit each record or each batch.

At-least-once processing is the common result: if you crash after the work and before the commit, you process again. Topic 5 covers that case.

Auto-commit is easy to get wrong. You can process a record after the client has already committed it. A crash then skips the record. That is at-most-once. Extra threads that hold records increase that risk.

For learning, turn auto-commit off. Commit in your code after you finish the work. If you keep auto-commit on, finish all work on the polled records before you call `poll` again, and do not hand records to other threads.

Do not mix policies in one team without a written rule.

### Questions

#### Theoretical questions

1. What does a commit store?
2. When does auto-commit write offsets?
3. What is the difference between `commitSync` and `commitAsync`?
4. How can auto-commit skip a record?
5. Why do teams choose manual commit for important work?

#### Easy practical tasks

1. Make a table: Auto-commit, manual sync, manual async. Add one row for control and one row for typical risk.
2. Write four sentences about a commit. Use only facts from this section.
3. Find the default of `enable.auto.commit` in your client version.
4. Draw a time line: poll, process, crash, restart. Mark a missing commit.

#### Medium practical tasks

1. Write a consumer with manual `commitSync` after each poll batch. Restart it. Confirm that it continues after the last commit.
2. Describe a consumer group with `kafka-consumer-groups`. Write the committed offset and the lag.
3. Write a consumer with auto-commit on and a long sleep after print. Stop it during the sleep. Restart. Write whether you see duplicates or skips.

#### Advanced practical tasks

1. Implement batch commit every N records. Document N. Show lag in the consumer-groups tool.
2. Write a one-page policy: when auto-commit is allowed (low-value metrics) and when it is forbidden (payments).

---

## `max.poll.interval.ms` and session timeout

The consumer has two important timers.

`session.timeout.ms` (with heartbeat interval) controls liveness of the process in the group. If heartbeats stop, the coordinator removes the member. Heartbeats can run in a background thread in current Java clients.

`max.poll.interval.ms` controls how long the group allows between two `poll` calls. If your processing takes longer than this interval, the coordinator removes the member even if heartbeats still run. The partition is assigned to another member. Your old member can then commit or process records that it no longer owns. That is a serious defect.

Set `max.poll.interval.ms` above your worst processing time for one poll batch. Or reduce `max.poll.records` so that each batch is smaller.

`heartbeat.interval.ms` must be lower than `session.timeout.ms`. The official docs give a ratio. Do not set random values.

These timers are not the same as `delivery.timeout.ms` on the producer.

Lag is the difference between the log end and the committed offset. If you raise `max.poll.interval.ms` only to hide slow work, lag still grows. Fix the sink or pause partitions (topic 11).

### Questions

#### Theoretical questions

1. What does `session.timeout.ms` detect?
2. What does `max.poll.interval.ms` detect?
3. Why can heartbeats continue while the poll interval expires?
4. What can go wrong if a revoked member still processes records?
5. How does `max.poll.records` help a slow processor?

#### Easy practical tasks

1. Make a table: Timer, what it watches, typical failure.
2. Find the defaults for `session.timeout.ms`, `max.poll.interval.ms`, and `heartbeat.interval.ms` in your client.
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

## Consumer groups: one active consumer per partition

A **consumer group** is a set of consumers that share one `group.id`. The group shares the partitions of the subscribed topics. Kafka assigns each partition to at most one active member of the group.

That rule keeps one reader per partition. Order in the partition stays simple. The member does not share the partition with a peer in the same group.

If the group has more members than partitions, extra members are idle. They receive no partitions. They still join the group. They become useful when a member leaves or when you add partitions.

If the group has fewer members than partitions, some members receive more than one partition.

Two groups do not follow this rule with each other. The same partition can have one active reader in group A and one active reader in group B at the same time.

The **group coordinator** is a broker. Kafka maps each group id to a coordinator. Members send join, heartbeat, and commit requests to that coordinator. You do not pick the coordinator by hand. You set `bootstrap.servers` and `group.id`.

You scale a group by adding members. Throughput can grow until the member count equals the partition count. After that, extra members do not increase parallelism for that topic in that group. Scale independent work with a new group. A billing service and an email service must not share a `group.id`.

Uneven keys create a hot partition. Extra members do not help a hot partition. Watch lag.

Do not create a huge number of unused group ids. Each group uses storage in `__consumer_offsets` until offsets expire.

`kafka-consumer-groups` talks to the cluster and shows state: members, assignment, lag. Use it when you debug a group.

### Questions

#### Theoretical questions

1. How many active members of one group read one partition?
2. What happens when you have more members than partitions?
3. How can two processes still read the same partition at the same time?
4. What is the group coordinator?
5. What is the parallel consume limit for one group on one topic?

#### Easy practical tasks

1. Draw three partitions and two members. Assign partitions. Draw four members and three partitions. Mark the idle member.
2. Write four sentences about the one-partition rule.
3. Run `kafka-consumer-groups --list`. Write the group ids that you see.
4. Make a table: Action, what it scales. Add member, partition, new group.

#### Medium practical tasks

1. Create a two-partition topic. Start three consumers in one group. Produce records. Confirm that one consumer is idle.
2. Start a second group on the same topic. Confirm that both groups receive records.
3. Produce a steady stream. Start one consumer, then a second, in one group on a four-partition topic. Write how lag changes.

#### Advanced practical tasks

1. Use a key that always hashes to one partition. Add members. Write why lag stays on one partition.
2. Write a one-page note on revoke: how a member must stop work on a partition that it no longer owns.

---

## Rebalance, static membership, and cooperative assignors

A **rebalance** is a change of partition assignment in a group. Members join, leave, or fail. The assignor computes a new assignment.

**Eager rebalance** (old protocol). All members revoke all partitions first. Consumption stops. Then members receive a new assignment. This is a pause for the whole group.

**Cooperative rebalance** (incremental). Members revoke only the partitions that they must give away. They keep the rest. Consumption continues on the partitions that do not move.

An **assignor** maps partitions to members. Common Java names: `RangeAssignor`, `RoundRobinAssignor`, `StickyAssignor`, `CooperativeStickyAssignor`.

**Sticky** means the assignor tries to keep a partition on the same member when it can. **Cooperative** means the rebalance protocol is incremental. `CooperativeStickyAssignor` combines both ideas. It is the usual choice for new applications.

All members of a group must use compatible assignors. Set the same strategy list on every member.

**Static membership** gives a consumer a stable identity. You set `group.instance.id` to a unique string per process instance (for example a pod name). Without static membership, a restart is a leave and a join. The group rebalances. With static membership, the coordinator can keep the assignment during a short restart if the member returns before the session timeout.

`group.instance.id` must be unique in the group. Two processes with the same instance id are a serious misconfiguration. Static membership does not replace a cooperative assignor. Use both.

During a rebalance, commit offsets for partitions that you revoke. Current clients have revoke callbacks. Commit in the revoke callback when you use manual commit. Do not take a long lock in a rebalance callback.

Causes of frequent rebalances: members restart without static membership, `max.poll.interval.ms` expires, session timeout expires, you add or remove members often.

This handbook does not use ZooKeeper group coordination. Current Kafka consumer groups are on the brokers.

### Questions

#### Theoretical questions

1. What is a rebalance?
2. What does an eager rebalance revoke?
3. What does a cooperative rebalance revoke?
4. What problem does `group.instance.id` reduce?
5. Why must all members use a compatible assignor?

#### Easy practical tasks

1. Write five sentences about rebalance types and static membership. Use only facts from this section.
2. Make a table: Eager vs cooperative. Add rows for revoke scope and consume pause.
3. Find the default assignor in your client version.
4. Write three good instance id values and two bad values (duplicates or random every start).

#### Medium practical tasks

1. Start two consumers in one group. Add a third. Watch logs for revoke and assign. Write whether all partitions revoked.
2. Set the cooperative sticky assignor. Repeat the join. Compare the revoke set.
3. Run a consumer with a static instance id. Restart it quickly. Describe the group during the restart. Write whether partitions bounce to another member.

#### Advanced practical tasks

1. Measure the consume pause (time with no records) during a join for eager versus cooperative on a topic with many partitions.
2. Write a production group standard: required assignor, required instance id scheme, max poll interval policy, and KRaft bootstrap only.

---

## Seek

Seek moves the consumer position in a partition. The next `poll` reads from that offset.

- Seek to an offset. You must assign the partition (or wait until it is assigned). Some clients require a poll before seek.
- Seek to the beginning. The next read is the earliest retained offset.
- Seek to the end. The next read is the log end. You skip existing records.

Seek does not always commit. If you seek and then crash before a commit, the group can restart at the old committed offset. If you want the seek to survive a restart, commit after the seek (or use the admin reset tool).

The consumer-groups tool can reset offsets for a group that is idle. Use that tool for operations. Use seek in a running consumer for replay in process.

Do not seek on a partition that you do not own in a subscribe group. The assignment can change.

`--from-beginning` on the console consumer is a seek-to-beginning for that start, often with a new group.

Some clients support seek by timestamp (`offsetsForTimes`). The broker maps a time to an offset through the time index (topic 6).

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

1. Seek to a timestamp if the API supports it. Replay one minute of a topic. Document the API calls.
2. Write a one-page warning: seek plus auto-commit can undo the seek. Give a sequence.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe the safe manual-commit loop from subscribe to close in ordered steps.
2. How do subscribe, poll timers, and a cooperative sticky assignor work together in a live group?
3. Why is assign the wrong default for a horizontally scaled worker service?
4. What is the difference between "position in this process" and "committed offset for the group"?
5. A teammate adds 20 consumers to a 3-partition topic to "go faster". What do you explain?

#### Easy practical tasks

1. Write a consumer with auto-commit off, subscribe, poll, print, `commitSync`, close. Run it on `demo`.
2. Write a one-page cheat sheet: subscribe vs assign, auto vs manual commit, two timers, one-partition rule, static id, cooperative assignor, seek.
3. Describe your consumer group. Highlight assignment, committed offset, and lag.
4. Draw eager revoke-all versus cooperative revoke-some.

#### Medium practical tasks

1. Write a script that starts a consumer, produces 20 records, and prints lag until lag is 0.
2. Demonstrate seek to beginning and a group reset with the admin tool as two different methods. Write when you pick each method.
3. Give two members static instance ids. Restart one. Save `--describe` output during the restart.

#### Advanced practical tasks

1. Build a consumer that pauses a partition when a downstream API fails, retries, then resumes. Do not exceed the poll interval.
2. Break a consumer with a slow poll. Watch the rebalance. Fix it with a smaller poll batch. Record before and after group describe output.
