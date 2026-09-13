# 18. Multi-Cluster and Geo

## Description

One Kafka cluster lives in one failure domain even when you use racks (topic 15). A second cluster in another region or another account is a different system. You copy topics between clusters when you need disaster recovery, geo read locality, or a split of teams.

This topic covers MirrorMaker 2, Cluster Linking (Confluent), active-passive versus active-active, offset translation, and disaster recovery. Complete this topic after topics 5, 6, 9, 12, and 15. Each cluster uses KRaft. Do not add ZooKeeper. Replication tools talk to each cluster with `bootstrap.servers`.

Use one term for each concept. A source cluster is the cluster that already has the records. A target cluster is the cluster that receives the copy. MirrorMaker 2 (MM2) is the Apache Connect-based replicator. Cluster Linking is a Confluent product that links topics. Offset translation maps a consumer group position from one cluster to another. Disaster recovery (DR) is the plan to continue after a cluster or region fails.

---

## MirrorMaker 2

MirrorMaker 2 is the Apache Kafka replicator. It runs as Kafka Connect connectors. It copies topics from a source cluster to a target cluster. It can copy consumer group offsets (with translation). It can create topics on the target with a name pattern (often a prefix or a suffix).

MM2 uses source and checkpoint connectors. Workers need network access to both clusters. Configure two `bootstrap.servers` values. Both clusters can be KRaft. MM2 does not need ZooKeeper.

The copy is asynchronous. The target lags the source. That lag is your replication delay. Measure it.

MM2 is not a transaction that spans two clusters. Exactly-once inside one cluster (topic 7) does not make two clusters one atomic log.

Do not run two uncontrolled MM2 flows that write the same target topic from two sources without a design. You get duplicates or conflicts.

Identity and ACLs (topic 17) apply on both clusters. The MM2 principal needs read on the source and write on the target.

Typical name pattern: source topic `orders.placed` becomes `sourceCluster.orders.placed` on the target. Change the pattern only with a plan. Consumers must use the target name.

```text
# Conceptual worker / connector keys (names follow current MM2 docs)
source.cluster.bootstrap.servers=...
target.cluster.bootstrap.servers=...
topics=orders.placed
```

### Questions

#### Theoretical questions

1. What does MirrorMaker 2 copy?
2. On what framework does MM2 run?
3. Why is the target always behind the source?
4. Does Streams EOS on the source make the two clusters atomic?
5. Why does the target topic name often include a prefix?

#### Easy practical tasks

1. Write five sentences about MM2. Use only facts from this section.
2. Make a table: Cluster, bootstrap, MM2 needs (read or write).
3. Find official MM2 documentation. Write the connector class names that the page lists.
4. Draw: source KRaft cluster → MM2 workers → target KRaft cluster.

#### Medium practical tasks

1. Read an official MM2 quick start. Write every process it starts (Connect, two Kafka clusters).
2. List three MM2 configuration keys and their meaning in your own words.
3. Compare MM2 with a custom consume-and-produce program. Write four differences.

#### Advanced practical tasks

1. Run MM2 between two local KRaft clusters (two Compose stacks or two ports). Copy one topic. Confirm records on the target. Write the target topic name.
2. Write an operations note: ACLs for MM2, how you measure lag, and how you restart Connect without ZooKeeper.

---

## Cluster Linking (Confluent)

**Cluster Linking** is a Confluent Platform and Confluent Cloud feature. It creates a link from a source cluster to a destination cluster. A **mirror topic** on the destination follows a source topic.

Cluster Linking is not part of Apache Kafka Apache-only tarball. You use it when you run Confluent products that include it. This handbook treats it as a comparison point.

Differences that Confluent documents (high level; confirm on the current product page):

- The destination can keep the same topic name (no prefix) in many setups.
- Offset preservation can be closer to the source than a naive consume-produce copy.
- The link is a first-class cluster object, not only a Connect job that you wrote.

You still have two clusters. You still have async copy and lag. You still need ACLs and TLS on both sides. Destination and source can be KRaft.

Do not assume Cluster Linking exists on Amazon MSK or on a plain Apache install. If you have only Apache Kafka, use MM2 or another replicator that you support.

Cluster Linking does not remove the need for a DR runbook (later section). A link can stop. You must know how to promote a mirror or how to fail over.

### Questions

#### Theoretical questions

1. What product family includes Cluster Linking?
2. What is a mirror topic?
3. Why can Cluster Linking keep the same topic name when MM2 often adds a prefix?
4. Does Cluster Linking run on every Apache Kafka install?
5. Why do you still need a failover plan?

#### Easy practical tasks

1. Write four sentences about Cluster Linking. Use only facts from this section.
2. Make a table: MM2, Cluster Linking. Add rows for Apache-only, topic names, where it runs.
3. Open the current Confluent Cluster Linking page. Write three claims the page makes. Do not copy a long passage.
4. Draw source topic → link → destination mirror topic.

#### Medium practical tasks

1. Compare official MM2 offset behavior with Confluent Cluster Linking offset notes. Write five differences or "not stated".
2. List two Confluent Cloud and two Apache-only options for geo copy.
3. Write when your team would refuse Cluster Linking (no Confluent) and what you use instead.

#### Advanced practical tasks

1. If you have a Confluent trial: create a link and a mirror topic. Produce on the source. Consume on the destination. Write lag if shown.
2. Write a decision page: Apache MM2 versus Cluster Linking, license, operations skill, KRaft on both sides.

---

## Active-passive vs active-active

**Active-passive** means one cluster receives all produces for a topic (the active). The other cluster is a copy (the passive). On disaster, you switch clients to the passive cluster after a promote step. There is one writer. Conflict is rare.

**Active-active** means two (or more) clusters accept produces on the "same" business topic at the same time. Each cluster replicates to the other. The same key can change in two places. You get conflicts. You must design merge rules or you must partition writes by region (users of region A write only to cluster A).

Active-active is not "twice the availability for free". It is a data-design problem. Last-write-wins by timestamp can lose a write. Event sourcing (topic 19) can help if every write is an event with a unique id, but you still order per key.

Most DR designs start as active-passive. Active-active is for low-latency writes in two regions when the business can isolate keys or can merge.

Consumers in a passive region can read the local copy if you accept lag. That is a read-local, write-home pattern.

Do not run two MM2 loops on the same topic in both directions without a filter (for example, exclude the replicated prefix) or you create a loop.

Both clusters use KRaft. Failover is a client bootstrap and ACL change, not a ZooKeeper migrate.

### Questions

#### Theoretical questions

1. What is active-passive?
2. What is active-active?
3. Why does active-active create conflicts?
4. What is a write-home, read-local pattern?
5. How can two-way replication create a loop?

#### Easy practical tasks

1. Write five sentences that contrast the two modes. Use only facts from this section.
2. Make a table: Mode, writers, conflict risk, typical use.
3. Draw active-passive: all produces to cluster A, copy to B.
4. Draw a loop: A→B and B→A on the same topic name with no filter.

#### Medium practical tasks

1. Pick a shop: orders must be unique. Choose active-passive or active-active. Write the reason in six sentences.
2. Design a key rule that makes active-active safer (for example, region in the key). Write an example key.
3. Find MM2 documentation on preventing replication loops. Write the idea in STE.

#### Advanced practical tasks

1. Write a one-page design: two regions, read-local, write-home, MM2 one way, how you fail over writes.
2. Write a conflict example with two updates to the same customer profile in two clusters. Show why last-write-wins can be wrong.

---

## Offset translation

A consumer group offset is a position in a **specific cluster** (topic 5 and 6). Offset 1000 on cluster A is not the same record as offset 1000 on cluster B. The target log can have different partition counts, different start offsets, and extra control records.

**Offset translation** maps a group’s committed offset on the source to an offset on the target so that a consumer that fails over does not skip a large gap or repeat the whole topic.

MM2 checkpoint connectors write translated offsets into topics on the target. You then apply those offsets to a group on the target (procedure depends on version). Cluster Linking has its own offset features. Read the product you use.

Without translation, a common mistake is to start the group on the target at latest (lose data) or at earliest (reprocess everything).

Translation is per partition. If the target has a different partition count, mapping is harder or impossible for a simple offset number. Prefer the same partition count and the same key partitioner on both clusters.

Compacted topics and truncated logs make translation harder. Test the failover of a group in a lab.

KRaft does not store your application group offsets in a special geo table. `__consumer_offsets` is per cluster.

### Questions

#### Theoretical questions

1. Why is offset 1000 not portable between clusters?
2. What does offset translation do?
3. What goes wrong if you start at latest on the target after failover?
4. Why must partition counts match for a simple mapping?
5. Where do group offsets live on each cluster?

#### Easy practical tasks

1. Write five sentences about offset translation. Use only facts from this section.
2. Make a table: Strategy on target (earliest, latest, translated), skip risk, duplicate risk.
3. Find MM2 checkpoint or offset-sync documentation. Write the topic or connector role name.
4. Draw source partition offset 50 → mapped target offset (unknown number) → consumer group.

#### Medium practical tasks

1. Consume on the source to offset N. Replicate. Try to find the matching record on the target by key. Write the two offsets.
2. Read official limits of offset translation. Write three cases where it fails or is approximate.
3. Write a failover checklist step: "apply translated offsets" versus "reset to earliest".

#### Advanced practical tasks

1. In a two-cluster lab, fail over a group with MM2 checkpoints (or the current procedure). Measure duplicates or gaps on a numbered payload.
2. Write an offset-translation standard: same partition count, how you test, who may reset a group on the DR cluster.

---

## Disaster recovery

**Disaster recovery** is the documented path when a cluster or a region fails. You define:

- **RPO** (recovery point objective): how much recent data you can lose (related to replication lag).
- **RTO** (recovery time objective): how long until clients write and read again.

A multi-cluster copy is necessary but not sufficient. You also need DNS or bootstrap change, ACLs on the target, schema registry availability (topic 11), Connect and Streams `application.id` behavior, and a decision maker.

Failover steps (active-passive, shape only):

1. Stop producers to the failed source (or confirm they fail).
2. Confirm target lag and what you accept as RPO.
3. Promote the target (MM2 or Cluster Linking promote, or change it to the writer).
4. Apply offset translation for groups that must continue.
5. Point clients to the target bootstrap. Use TLS and SASL for that cluster (topic 17).
6. When the old region returns, do not start two actives without a failback plan.

Failback is a second project. You must copy or catch up the old cluster and switch again.

Test DR. An untested runbook fails. Use KRaft on both sides so that the runbook has no ZooKeeper restore step.

Single-cluster rack awareness (topic 15) is not DR for a region loss.

### Questions

#### Theoretical questions

1. What is RPO in a Kafka geo design?
2. What is RTO?
3. Why is a running MM2 not a complete DR plan?
4. Why must you test failover?
5. Why does this handbook refuse a ZooKeeper step in a new DR runbook?

#### Easy practical tasks

1. Write five sentences about Kafka DR. Use only facts from this section.
2. Make a table: Step, owner, depends on (MM2, DNS, ACL).
3. Write an example RPO and RTO for a lab shop (numbers you choose).
4. List five systems besides Kafka brokers that the runbook must name (clients, registry, Connect).

#### Medium practical tasks

1. Write a one-page active-passive runbook for two local KRaft clusters. Leave no ZooKeeper step.
2. Add a failback paragraph. Write the risk if both clusters accept writes.
3. Compare your runbook with official MM2 or Cluster Linking failover notes. Write three gaps.

#### Advanced practical tasks

1. Rehearse a failover in a lab: kill the source, promote, consume on the target with translated offsets. Time RTO. Estimate RPO from lag.
2. Write a production DR standard: RPO/RTO, test calendar, promote authority, schema and security, KRaft only.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do MM2, Cluster Linking, and a custom consumer-producer differ, and when is each honest?
2. How do active-passive, offset translation, and RPO fit in one sentence each?
3. What must be equal on source and target for a clean failover (partitions, keys, security, schemas)?
4. How does topic 15 rack awareness differ from this topic?
5. What KRaft fact stays true when you have two clusters?

#### Easy practical tasks

1. Write a one-page cheat sheet: MM2, Cluster Linking, two active modes, offset translation, RPO/RTO.
2. Draw two KRaft clusters, MM2, a consumer group, and a failover arrow.
3. Bookmark Apache MM2 and one Confluent Cluster Linking page.
4. List every bootstrap and ACL need for a replicator principal.

#### Medium practical tasks

1. Design (on paper or lab) two KRaft clusters and one MM2 flow for one topic. Fill names, prefixes, and a group failover step.
2. Map each subsection to one official URL.
3. Write a table: failure (AZ, region, bad ACL, MM2 down). What still works?

#### Advanced practical tasks

1. Run a full lab: two clusters, MM2, produce, fail source, translate offsets, produce to target. Write numbers for lag, duplicates, and time.
2. Write a geo standard for your team: default active-passive, when active-active is allowed, product choice, KRaft, no ZooKeeper.
