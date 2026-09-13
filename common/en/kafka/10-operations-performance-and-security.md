# 10. Operations, Performance, and Security

## Description

Operations keep a Kafka cluster healthy. Performance is the rate of records and the delay of each record. Security protects the network, authenticates clients, and limits who may read or write.

This topic shows under-replicated partitions, disk, rolling restarts, throughput versus latency, batch and fetch tunables, TLS, SASL, ACLs, and quotas. Complete this topic after topics 3, 4, 6, and 7. Test on a KRaft cluster. Do not use ZooKeeper for new work.

Use one term for each concept. An under-replicated partition (URP) is a partition whose ISR is smaller than the replica list. Throughput is records or bytes per second. Latency is the time from produce send to a defined end (ack, or consume). TLS encrypts the TCP connection. SASL authenticates a user name or token. An ACL is a rule that allows or denies an operation on a resource. A quota limits how much a client may produce, fetch, or request.

---

## Under-replicated partitions, disk, and rolling restarts

Topic 6 defines the ISR. The replica list is the assigned copies. The ISR is the copies that are in sync. When the ISR is smaller than the replica list, the partition is **under-replicated**.

A URP alert means a follower is slow, stopped, or cut off. Durability is lower until the follower returns to the ISR. `acks=all` and `min.insync.replicas` can start to reject writes if the ISR shrinks too far (topic 6).

Monitor URP count as a cluster metric. A short URP during a rolling restart can be normal. A URP that stays after the broker is up is an incident.

`kafka-topics --describe` shows Replicas and ISR. Tools and metrics use names such as `UnderReplicatedPartitions`.

Do not ignore URP because "clients still produce". Producers can succeed on a smaller ISR. The risk is loss if more brokers fail.

Kafka stores partition logs in `log.dirs` (topic 6). Disk fills when produce rate is higher than retention delete, or when you add partitions and topics.

Retention delete and compaction run on closed segments. Many segments can become eligible at the same time. Then Kafka deletes a large amount of data in a short period. That is a **bursty delete**. Disk use drops in a step. Disk I/O can spike.

Do not delete segment files by hand. Do not fill the disk to 100%. Kafka and the OS need free space. Plan a headroom percent. Alert before the disk is full.

KRaft metadata logs also use disk on controller nodes. They are smaller than data logs in most clusters. Still monitor controller disks.

A **rolling restart** stops one broker at a time. The controller elects new leaders for partitions that lived on the stopped broker. Clients retry. Then you start the broker and wait until URP returns to zero before you stop the next broker.

Do not stop all brokers at once. Do not restart the next broker while URP is still high from the previous stop.

You can **reassign** partitions when you add a broker or replace a disk. The official tool is `kafka-reassign-partitions` (name can differ by version). Cruise Control is a tool that can propose and run partition moves. Read the product you use.

**Rack awareness** places replicas on different failure domains (racks or zones). A rack is not disaster recovery for a full region. Topic 12 covers multi-cluster copy.

Page cache is the operating-system cache of disk pages. Kafka reads often come from the page cache. Do not run other disk-heavy jobs on the same disks without a plan.

KRaft does not remove URP. KRaft manages metadata and leader election. Partition data still replicates between brokers.

### Questions

#### Theoretical questions

1. What is an under-replicated partition?
2. When can a short URP be expected?
3. What is a bursty delete?
4. What is a rolling restart?
5. Why must you wait for URP to return to zero before you stop the next broker?

#### Easy practical tasks

1. Describe a topic. Write Replicas and ISR for each partition. Say if any partition is under-replicated.
2. Write five sentences about URP, disk, and rolling restart. Use only facts from this section.
3. On a lab broker, find `log.dirs` and write the path and the free space.
4. Find the metric name for URP in official monitoring docs.

#### Medium practical tasks

1. On a multi-broker KRaft cluster, stop one broker. Describe a replicated topic. Write the URP state. Start the broker. Write when ISR recovers.
2. Create a small-retention topic. Produce enough data to roll segments. Watch disk and segment files before and after delete.
3. Draw a timeline: rolling restart of one broker, URP rises, URP returns to zero.

#### Advanced practical tasks

1. Measure time from broker stop to URP and time from start to ISR full. Write both times and the topic settings.
2. Write a one-page operations standard: URP alerts, disk headroom, rolling restart order, KRaft controllers, no ZooKeeper.

---

## Throughput vs latency; batch and fetch tunables

**Throughput** is how much data the system moves in a period. Teams measure records per second or bytes per second. You must say if the number is produce, consume, or both.

**Latency** is how long one record waits. Produce latency often means time to the ack. End-to-end latency means time to a consumer process. Those two numbers are not the same.

A large batch raises throughput. The first record in the batch waits for `linger.ms` or for the batch to fill (topic 3). That wait raises latency.

A small batch or `linger.ms=0` can lower latency and lower throughput. More requests hit the broker. CPU and request queues grow.

Producer tunables (topic 3): `batch.size`, `linger.ms`, `compression.type`, `buffer.memory`, `acks`. Idempotence and `acks=all` add a small cost and protect data.

Consumer tunables (topic 4): `fetch.min.bytes`, `fetch.max.wait.ms`, `max.partition.fetch.bytes`, `max.poll.records`. A larger fetch can raise throughput. The consumer waits up to `fetch.max.wait.ms` if `fetch.min.bytes` is not met. That wait can raise latency.

Many partitions can raise parallelism and throughput. They can also raise metadata, memory, and end-to-end delay if each partition has little data (topic 2).

Do not optimize without a number. Write the target: for example, 50 000 records per second produce, p99 produce ack under 20 ms. Then measure. Change one setting at a time.

Zero-copy transfer on the broker can send data from page cache to the socket with fewer copies. You do not configure this as a beginner switch. You keep disks and page cache healthy.

KRaft metadata is not the usual produce bottleneck on a stable cluster. A saturated controller can still slow admin operations and new leadership. Measure brokers and clients first for data-plane tests.

A slow poll looks like a Kafka outage. It is often your code (topic 4 and topic 11).

### Questions

#### Theoretical questions

1. What is throughput in this handbook?
2. What two latency ends must you not mix?
3. How can a larger batch raise latency?
4. What does `fetch.min.bytes` wait for?
5. Why must you write a numeric target before you tune?

#### Easy practical tasks

1. Write five sentences about throughput and latency. Use only facts from this section.
2. Make a table: Change, throughput effect, latency effect. Add linger, fetch wait, more partitions.
3. Write one numeric target for a lab topic (rate and p99).
4. Find official notes that mention batching and `fetch.min.bytes`. Write one sentence.

#### Medium practical tasks

1. Produce with a small client at two linger values. Estimate records per second. Write both numbers and the method.
2. Consume with two `fetch.max.wait.ms` values. Write whether end-to-end delay changed.
3. List three metrics you will collect (producer, broker, consumer) from official monitoring names.

#### Advanced practical tasks

1. Measure produce p50 and p99 and records per second for two batch settings on KRaft. Write a short report with limits of the test.
2. Write a one-page SLO draft: throughput, produce latency, end-to-end latency, and how you measure each.

---

## TLS, SASL, ACLs, and quotas

**TLS** (Transport Layer Security) encrypts bytes on the network between clients and brokers, and between brokers. Without TLS, a listener on a plain port sends records in the clear.

You configure a listener with a protocol such as `SSL` or `SASL_SSL`. You give the broker a keystore (the broker certificate and key) and a truststore (the CAs that the broker trusts). Clients get a truststore that trusts the broker certificate.

TLS does not replace ACLs. TLS protects the pipe. ACLs decide the API operations.

Certificate expiry stops clients. Monitor not-after dates. Rotate certificates before they expire. Do not copy a private key into a git repository.

KRaft controllers need TLS on the controller listener if you encrypt that path. Combined listeners and controller listeners can differ. Read the configuration for your version.

`kafka-topics` and other tools must use the same security protocol and the same bootstrap port. Example shape:

```text
kafka-topics --bootstrap-server localhost:9093 --command-config client.properties --list
```

**SASL** is a family of authentication mechanisms. Kafka uses SASL on a listener such as `SASL_SSL` (SASL on TLS) or `SASL_PLAINTEXT` (SASL without TLS). Do not use `SASL_PLAINTEXT` on a network that you do not trust.

**PLAIN** uses a user name and a password. Use TLS with PLAIN.

**SCRAM** (for example SCRAM-SHA-256) stores a hashed credential. You create users with `kafka-configs` or the Admin API.

**OAUTHBEARER** uses an OAuth 2.0 access token. Cloud products often use this path. You must handle token refresh in the client.

**mTLS** (mutual TLS) authenticates the client with a client certificate. The broker maps the certificate identity to a user.

Inter-broker authentication also uses SASL or mTLS. KRaft controllers authenticate on the controller listener.

Authentication gives a **principal**. **ACLs** use that principal. An ACL allows or denies an operation (read, write, create, describe) on a resource (topic, group, cluster). A **super user** bypasses ACL checks. Use super users only for break-glass admin.

Do not leave a cluster with no authorizer and a public listener. Enable an authorizer (official name for your version) and write ACLs before you expose the cluster.

**Quotas** limit produce byte rate, fetch byte rate, or request rate per client or per user. Quotas protect the cluster from one noisy client. Quotas are not a substitute for backpressure in your consumer (topic 11).

Encryption of data at rest is usually a disk or volume feature. Kafka does not replace disk encryption.

Old ZooKeeper ACLs are not the model for new work. Apply security on KRaft clusters.

### Questions

#### Theoretical questions

1. What does TLS protect on a Kafka listener?
2. Why is TLS not enough without ACLs?
3. What does SASL add that TLS alone does not add?
4. What is a principal in ACL terms?
5. What do quotas limit?

#### Easy practical tasks

1. Write five sentences about TLS, SASL, ACLs, and quotas. Use only facts from this section.
2. Make a table: Mechanism, PLAIN, SCRAM, OAUTHBEARER, mTLS. Add one row: needs TLS (yes or recommended).
3. Find official ACL operation names (Read, Write, and others). Write five names.
4. Draw client → TLS → broker, and principal → ACL → topic.

#### Medium practical tasks

1. Enable TLS on a local KRaft broker (self-signed in a lab). List topics with a command-config file. Save the property keys, not the passwords, in your notes.
2. Create a SCRAM user if your cluster supports it. Produce as that user. Deny write with an ACL. Record the error. Then allow write.
3. Set a produce quota on a lab user or client id. Produce a burst. Write whether you see throttle errors or delay.

#### Advanced practical tasks

1. Enable inter-broker TLS and a client `SASL_SSL` listener. Produce and consume. Write which ports you used. Include the KRaft controller listener if you encrypt it.
2. Write a security standard: TLS required, SASL mechanism, ACL defaults (deny), quota defaults, certificate rotation, no secrets in git, no ZooKeeper.

---

## Topic review

These questions do not repeat the questions in the sections above. They cover the full topic.

### Questions

#### Theoretical questions

1. How do URP, disk headroom, and a rolling restart form one operations story?
2. Why can a change that raises throughput also raise p99 latency?
3. How do TLS, SASL, and ACLs stack: pipe, identity, permission?
4. What stays the same when you move a lab from PLAINTEXT to `SASL_SSL` on KRaft?
5. A teammate wants to restart all brokers at once "to go faster". Which facts do you use in the reply?

#### Easy practical tasks

1. Write a one-page cheat sheet: URP, bursty delete, rolling restart, throughput versus latency, linger and fetch keys, TLS, SASL, ACL, quota.
2. Describe one topic. Highlight ISR versus replicas. Write free disk for `log.dirs`.
3. Bookmark official pages for monitoring, producer performance, and security.
4. List every configuration key named in this topic in one column and a six-word meaning in the second.

#### Medium practical tasks

1. Run a rolling restart of one broker on a multi-broker lab. Record URP before, during, and after. Produce during the restart. Write whether clients recovered.
2. Measure one produce test and one consume test. Change one tunable. Write the before and after numbers.
3. Map each subsection to one official URL.

#### Advanced practical tasks

1. Combine a rolling restart with TLS clients. Write a runbook: drain, restart, wait for ISR, next broker, certificate check.
2. Write a production standard: KRaft only, URP and disk alerts, SLO numbers, `SASL_SSL`, ACL review, quota for shared clusters.
