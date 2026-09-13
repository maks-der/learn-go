# 6. Consumer Groups

## Description

A consumer group is a set of consumers that share one `group.id`. The group shares the partitions of the subscribed topics. Kafka assigns each partition to at most one active member of the group.

This topic covers the group coordinator, the one-partition rule, scaling, static membership, the cooperative sticky assignor, and stop-the-world rebalances. Complete this topic after topic 5.

Use one term for each concept. A member is one consumer process in the group. The group coordinator is the broker that stores membership for that group. An assignor is the algorithm that maps partitions to members. Use KRaft clusters. Group metadata lives on brokers (the `__consumer_offsets` topic and the coordinator). You do not use ZooKeeper for consumer groups in current Kafka.

---

## Group coordinator

The group coordinator is a broker. Kafka maps each group id to a coordinator. Members send join, heartbeat, and commit requests to that coordinator.

The coordinator:

- tracks living members
- starts a rebalance when membership changes
- stores committed offsets (the data is in `__consumer_offsets`)

If the coordinator broker stops, another broker becomes the coordinator for that group. Members find the new coordinator through the cluster metadata. Short errors can appear during that move.

You do not pick the coordinator by hand. You set `bootstrap.servers` and `group.id`. The client discovers the coordinator.

Do not create a huge number of unused group ids. Each group uses storage in `__consumer_offsets` until offsets expire.

`kafka-consumer-groups.sh` (or `.bat`) talks to the cluster and shows state: members, assignment, lag. Use it when you debug a group.

### Questions

#### Theoretical questions

1. What is the group coordinator?
2. How does a consumer find its coordinator?
3. What three duties does this section list for the coordinator?
4. What happens when the coordinator broker stops?
5. Why is a large set of dead group ids a problem?

#### Easy practical tasks

1. Write five sentences about the group coordinator. Use only facts from this section.
2. Run `kafka-consumer-groups --list`. Write the group ids that you see.
3. Make a table: Request type, why the member sends it. Add join, heartbeat, commit.
4. Find the official description of the group coordinator. Write it in your own words.

#### Medium practical tasks

1. Describe one group with `--describe`. Write the coordinator broker id if the tool shows it (or write the member host list).
2. Stop and start a consumer. Watch `--describe` before, during, and after. Write the state names that you see (`Stable`, `Empty`, or similar).
3. Read how `__consumer_offsets` relates to the coordinator. Write five sentences.

#### Advanced practical tasks

1. On a three-broker cluster, create many groups and see if coordinators spread (if the tool shows coordinator ids). Write a small table.
2. Write a one-page debug guide: member cannot join. Include bootstrap, group id, coordinator move, and ACL preview.

---

## One partition → at most one active consumer in a group

Kafka assigns each partition to at most one active member in a group. That rule keeps one reader per partition. Order in the partition stays simple. The member does not share the partition with a peer in the same group.

If the group has more members than partitions, extra members are idle. They receive no partitions. They still join the group. They become useful when a member leaves or when you add partitions.

If the group has fewer members than partitions, some members receive more than one partition.

Two groups do not follow this rule with each other. The same partition can have one active reader in group A and one active reader in group B at the same time.

"Active" means the current assignee. During a rebalance, ownership moves. A slow revoked member can still process a record after a new member already reads that partition. Prevent that with revoke callbacks and a stop of work for revoked partitions.

### Questions

#### Theoretical questions

1. How many active members of one group read one partition?
2. Why does that rule exist?
3. What happens when you have more members than partitions?
4. What happens when you have fewer members than partitions?
5. How can two processes still read the same partition at the same time?

#### Easy practical tasks

1. Draw three partitions and two members. Assign partitions. Draw four members and three partitions. Mark the idle member.
2. Write four sentences about the one-partition rule.
3. Make a table: Members, partitions, idle members. Add three rows.
4. Explain in three sentences why two groups can duplicate processing on purpose.

#### Medium practical tasks

1. Create a two-partition topic. Start three consumers in one group. Produce records. Confirm that one consumer is idle.
2. Describe the group. Write the assignment. Confirm that no partition appears twice.
3. Start a second group on the same topic. Confirm that both groups receive records.

#### Advanced practical tasks

1. Write a test plan that detects double processing in one group (it must fail if the assignor is correct). Run a long consume and check keys per partition per member.
2. Write a one-page note on revoke: how a member must stop work on a partition that it no longer owns.

---

## Scaling consumers

You scale a group by adding members. Throughput can grow until the member count equals the partition count. After that, extra members do not increase parallelism for that topic in that group.

You scale a topic by adding partitions. That change can move keys (topic 3). Plan partitions before you need a large group.

You can also give each member more CPU and faster processing. Then each member can handle more partitions.

Scale independent work with a new group, not with more members in the same group. A billing service and an email service must not share a `group.id`.

Uneven keys create a hot partition. Extra members do not help a hot partition. Fix the key or split the workload.

Watch lag. Lag is the difference between the log end and the committed offset. If lag grows, you need faster members, more partitions, or less work per record.

### Questions

#### Theoretical questions

1. What is the parallel consume limit for one group on one topic?
2. How do you raise that limit?
3. When do you create a new group instead of a new member?
4. Why does a hot partition ignore extra members?
5. What is lag?

#### Easy practical tasks

1. Write five sentences about scaling consumers. Use only facts from this section.
2. Make a table: Action, what it scales. Add member, partition, new group.
3. Describe a group and write lag per partition.
4. List three causes of growing lag.

#### Medium practical tasks

1. Produce a steady stream. Start one consumer, then a second, in one group on a four-partition topic. Write how lag changes.
2. Use a key that always hashes to one partition (one key). Add members. Write why lag stays on one partition.
3. Find official or Confluent guidance on consumer scaling. Write three facts.

#### Advanced practical tasks

1. Measure records per second with 1, 2, and 4 members on a 4-partition topic. Write a table.
2. Write a capacity plan: 40 MB/s topic, 10 ms process time per 1 KB record. Estimate partitions and members.

---

## Static membership

Static membership gives a consumer a stable identity. You set `group.instance.id` to a unique string per process instance (for example a pod name).

Without static membership, a restart is a leave and a join. The group rebalances. Partitions move and then move back. That bounce costs time.

With static membership, the coordinator can keep the assignment during a short restart. The member comes back with the same instance id before the session timeout. The group can avoid a full rebalance.

`group.instance.id` must be unique in the group. Two processes with the same instance id are a serious misconfiguration.

Static membership does not replace cooperative rebalance. Use both. Static membership reduces rebalances on restart. Cooperative rebalance reduces the pause when a rebalance still occurs.

Remove the instance id when you retire the instance, or wait until the group expires the member.

### Questions

#### Theoretical questions

1. What is `group.instance.id`?
2. What problem does static membership reduce?
3. Why must the instance id be unique?
4. What must happen before the session timeout after a restart?
5. Does static membership replace a cooperative assignor?

#### Easy practical tasks

1. Write four sentences about static membership. Use only facts from this section.
2. Make a table: Restart without static id, restart with static id.
3. Find `group.instance.id` in your client docs.
4. Write three good instance id values and two bad values (duplicates or random every start).

#### Medium practical tasks

1. Run a consumer with a static instance id. Restart it quickly. Describe the group during the restart. Write whether partitions bounce to another member.
2. Run without a static id. Repeat the restart. Compare.
3. Read official static membership documentation. Write three configuration keys that work with it.

#### Advanced practical tasks

1. On Kubernetes-style names, design an instance id scheme (StatefulSet name vs random pod name). Write why a random name defeats the feature.
2. Write a one-page runbook: replace a consumer host. Include session timeout and when a rebalance still occurs.

---

## Cooperative sticky assignor

An assignor maps partitions to members. Kafka includes several assignors. Names differ by client. Common Java names:

- `RangeAssignor`
- `RoundRobinAssignor`
- `StickyAssignor`
- `CooperativeStickyAssignor`

**Sticky** means the assignor tries to keep a partition on the same member when it can. That reduces movement.

**Cooperative** means the rebalance protocol is incremental. Members revoke only partitions that must move. Topic 5 introduced eager versus cooperative.

`CooperativeStickyAssignor` combines both ideas. It is the usual choice for new applications.

All members of a group must use compatible assignors. A mix of eager-only and cooperative strategies can fail or fall back. Set the same strategy list on every member.

After you change the assignor, roll members with care. Follow the official upgrade notes for your client version.

### Questions

#### Theoretical questions

1. What does an assignor do?
2. What does sticky mean?
3. What does cooperative mean?
4. Why must all members use a compatible assignor?
5. Why is `CooperativeStickyAssignor` a good default for new work?

#### Easy practical tasks

1. Write five sentences about the cooperative sticky assignor. Use only facts from this section.
2. Make a table: Assignor name, sticky, cooperative. Add four rows from this section.
3. Find how to set the assignor in your client.
4. List two problems that a non-sticky assignor can cause on a join.

#### Medium practical tasks

1. Start a group with the cooperative sticky assignor. Add a member. Write which partitions moved.
2. If you can, run a range assignor group and add a member. Compare how many partitions moved.
3. Read official assignor documentation. Write the configuration key and the class list.

#### Advanced practical tasks

1. Count partition moves across five join/leave events for sticky versus range. Write a table.
2. Write an upgrade plan from an eager assignor to cooperative sticky. Use official steps. Do not mix strategies at random.

---

## Stop-the-world rebalances (old protocol)

The old eager protocol is a stop-the-world rebalance. Every member revokes every partition. No member consumes until the new assignment exists. On a large group, that pause can last seconds or more.

Causes of frequent stop-the-world rebalances:

- members restart without static membership
- `max.poll.interval.ms` expires
- session timeout expires
- you add or remove members often (scale up and down)
- a bad health check kills processes

Effects: lag spikes, timeout errors downstream, and more load when all members resume.

The new path is cooperative rebalance plus static membership plus healthy poll times. That path still rebalances, but the pause is smaller.

Some old clients only support eager rebalance. Upgrade the client. Do not keep the old protocol for a large production group without a reason.

This handbook does not use ZooKeeper group coordination. Current Kafka consumer groups are on the brokers.

### Questions

#### Theoretical questions

1. What is a stop-the-world rebalance?
2. When do members consume during an eager rebalance?
3. Name four causes of frequent rebalances.
4. What three features reduce the pause?
5. Why must you upgrade an old eager-only client for a large group?

#### Easy practical tasks

1. Write four sentences about stop-the-world rebalances.
2. Make a table: Cause, how you reduce it.
3. Draw eager rebalance as a bar where consume is zero for all members.
4. List three metrics or logs that show a rebalance storm.

#### Medium practical tasks

1. Trigger many rebalances (join and leave in a loop) with an eager assignor. Watch lag. Write the lag shape.
2. Repeat with cooperative sticky and static ids if you can. Compare the lag shape.
3. Read official incremental cooperative rebalancing background. Write why eager is stop-the-world.

#### Advanced practical tasks

1. Write a one-page incident report template for a rebalance storm. Include timers, assignor, instance ids, and deploy events.
2. Measure time-to-stable after a member kill for eager versus cooperative. Write both times and the partition count.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. Describe a group from join to stable assignment. Name the coordinator, the assignor, and the one-partition rule.
2. How do scaling, hot keys, and partition count interact?
3. Why do static membership and a cooperative sticky assignor complement each other?
4. What storage on the cluster remembers both membership-related offsets and group commits?
5. A teammate adds 20 consumers to a 3-partition topic to "go faster". What do you explain?

#### Easy practical tasks

1. Create a 3-partition topic. Run two members in group `g-scale`. Describe the group. Write the assignment.
2. Write a one-page cheat sheet: coordinator, one-partition rule, scale limit, `group.instance.id`, cooperative sticky, eager pause.
3. List groups on your cluster. Mark empty groups and stable groups.
4. Draw a group of three members and six partitions after one member stops.

#### Medium practical tasks

1. Write a script that describes a group every two seconds while you add a member. Save the assignment history.
2. Give two members static instance ids. Restart one. Save `--describe` output during the restart.
3. Map this topic to official pages for consumer groups, assignors, and static membership.

#### Advanced practical tasks

1. Build a small dashboard (even a CSV) of lag per member during a rolling restart with and without static membership.
2. Write a production group standard: required assignor, required instance id scheme, max poll interval policy, and KRaft bootstrap only.
